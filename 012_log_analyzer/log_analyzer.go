package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/**
日志分析器
	使用Go语言实现的日志文件分析工具
	支持多种日志格式解析
	提供统计分析和报告生成功能
	支持错误日志过滤和分类
	展示Go语言文件处理、正则表达式、数据统计等特性
*/

// LogEntry 日志条目结构体
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	IP        string    `json:"ip,omitempty"`
	Method    string    `json:"method,omitempty"`
	URL       string    `json:"url,omitempty"`
	Status    int       `json:"status,omitempty"`
	Size      int       `json:"size,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Source    string    `json:"source"`
}

// LogStats 日志统计结构体
type LogStats struct {
	TotalLines   int            `json:"total_lines"`
	ValidLines   int            `json:"valid_lines"`
	ErrorLines   int            `json:"error_lines"`
	LevelCounts  map[string]int `json:"level_counts"`
	IPCounts     map[string]int `json:"ip_counts"`
	StatusCounts map[string]int `json:"status_counts"`
	MethodCounts map[string]int `json:"method_counts"`
	HourlyCounts map[string]int `json:"hourly_counts"`
	TopIPs       []string       `json:"top_ips"`
	TopURLs      []string       `json:"top_urls"`
	TopErrors    map[string]int `json:"top_errors"`
	StartTime    *time.Time     `json:"start_time"`
	EndTime      *time.Time     `json:"end_time"`
	TimeRange    string         `json:"time_range"`
}

// LogAnalyzer 日志分析器结构体
type LogAnalyzer struct {
	Entries  []LogEntry
	Stats    LogStats
	Patterns map[string]*regexp.Regexp
	Config   AnalyzerConfig
}

// AnalyzerConfig 分析器配置
type AnalyzerConfig struct {
	LogFormat     string // 日志格式类型
	TimeFormat    string // 时间格式
	FilterLevel   string // 过滤日志级别
	FilterPattern string // 过滤模式
	OutputFormat  string // 输出格式
	TopN          int    // 显示前N项统计
	ShowDetails   bool   // 显示详细信息
}

// 预定义的日志格式正则表达式
var logPatterns = map[string]string{
	"apache": `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) ([^"]*)" (\d+) (\d+|-) "([^"]*)" "([^"]*)"`,
	"nginx":  `^(\S+) - - \[([^\]]+)\] "(\S+) ([^"]*)" (\d+) (\d+|-) "([^"]*)" "([^"]*)"`,
	"common": `^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (.*)`,
	"json":   `^\{.*\}$`,
	"syslog": `^(\w{3} \d{1,2} \d{2}:\d{2}:\d{2}) (\S+) (\S+): (.*)`,
}

// 时间格式映射
var timeFormats = map[string]string{
	"apache":  "02/Jan/2006:15:04:05 -0700",
	"nginx":   "02/Jan/2006:15:04:05 -0700",
	"common":  "2006-01-02 15:04:05",
	"iso8601": "2006-01-02T15:04:05Z07:00",
	"rfc3339": time.RFC3339,
	"syslog":  "Jan 2 15:04:05",
}

// NewLogAnalyzer 创建新的日志分析器
func NewLogAnalyzer(config AnalyzerConfig) *LogAnalyzer {
	analyzer := &LogAnalyzer{
		Entries: make([]LogEntry, 0),
		Stats: LogStats{
			LevelCounts:  make(map[string]int),
			IPCounts:     make(map[string]int),
			StatusCounts: make(map[string]int),
			MethodCounts: make(map[string]int),
			HourlyCounts: make(map[string]int),
			TopErrors:    make(map[string]int),
		},
		Patterns: make(map[string]*regexp.Regexp),
		Config:   config,
	}

	// 编译正则表达式
	for name, pattern := range logPatterns {
		if compiled, err := regexp.Compile(pattern); err == nil {
			analyzer.Patterns[name] = compiled
		}
	}

	return analyzer
}

// ParseLogFile 解析日志文件
func (la *LogAnalyzer) ParseLogFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	return la.ParseLogReader(file, filename)
}

// ParseLogReader 解析日志读取器
func (la *LogAnalyzer) ParseLogReader(reader io.Reader, source string) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		la.Stats.TotalLines++

		entry, err := la.parseLine(line, source)
		if err != nil {
			la.Stats.ErrorLines++
			continue
		}

		// 应用过滤器
		if la.shouldFilterEntry(entry) {
			continue
		}

		la.Entries = append(la.Entries, entry)
		la.Stats.ValidLines++
		la.updateStats(entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取日志文件时出错: %v", err)
	}

	la.calculateDerivedStats()
	return nil
}

// parseLine 解析单行日志
func (la *LogAnalyzer) parseLine(line, source string) (LogEntry, error) {
	entry := LogEntry{Source: source}

	// 尝试JSON格式
	if strings.HasPrefix(line, "{") {
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			return entry, nil
		}
	}

	// 尝试预定义格式
	switch la.Config.LogFormat {
	case "apache", "nginx":
		return la.parseWebLog(line, source)
	case "common":
		return la.parseCommonLog(line, source)
	case "syslog":
		return la.parseSyslog(line, source)
	default:
		return la.parseGenericLog(line, source)
	}
}

// parseWebLog 解析Web服务器日志
func (la *LogAnalyzer) parseWebLog(line, source string) (LogEntry, error) {
	pattern := la.Patterns[la.Config.LogFormat]
	if pattern == nil {
		return LogEntry{}, fmt.Errorf("未找到格式模式: %s", la.Config.LogFormat)
	}

	matches := pattern.FindStringSubmatch(line)
	if len(matches) < 8 {
		return LogEntry{}, fmt.Errorf("日志格式不匹配")
	}

	entry := LogEntry{
		IP:        matches[1],
		Method:    matches[3],
		URL:       matches[4],
		UserAgent: matches[8],
		Source:    source,
	}

	// 解析时间
	if timestamp, err := time.Parse(timeFormats[la.Config.LogFormat], matches[2]); err == nil {
		entry.Timestamp = timestamp
	}

	// 解析状态码
	if status, err := strconv.Atoi(matches[5]); err == nil {
		entry.Status = status
		entry.Level = la.getLogLevelFromStatus(status)
	}

	// 解析大小
	if matches[6] != "-" {
		if size, err := strconv.Atoi(matches[6]); err == nil {
			entry.Size = size
		}
	}

	entry.Message = fmt.Sprintf("%s %s - %d", entry.Method, entry.URL, entry.Status)

	return entry, nil
}

// parseCommonLog 解析通用日志格式
func (la *LogAnalyzer) parseCommonLog(line, source string) (LogEntry, error) {
	pattern := la.Patterns["common"]
	matches := pattern.FindStringSubmatch(line)
	if len(matches) < 4 {
		return LogEntry{}, fmt.Errorf("日志格式不匹配")
	}

	entry := LogEntry{
		Level:   strings.ToUpper(matches[2]),
		Message: matches[3],
		Source:  source,
	}

	// 解析时间
	if timestamp, err := time.Parse(timeFormats["common"], matches[1]); err == nil {
		entry.Timestamp = timestamp
	}

	return entry, nil
}

// parseSyslog 解析系统日志格式
func (la *LogAnalyzer) parseSyslog(line, source string) (LogEntry, error) {
	pattern := la.Patterns["syslog"]
	matches := pattern.FindStringSubmatch(line)
	if len(matches) < 5 {
		return LogEntry{}, fmt.Errorf("日志格式不匹配")
	}

	entry := LogEntry{
		Level:   "INFO",
		Message: matches[4],
		Source:  source,
	}

	// 解析时间（系统日志通常没有年份，使用当前年份）
	timeStr := fmt.Sprintf("%d %s", time.Now().Year(), matches[1])
	if timestamp, err := time.Parse("2006 Jan 2 15:04:05", timeStr); err == nil {
		entry.Timestamp = timestamp
	}

	return entry, nil
}

// parseGenericLog 解析通用日志格式
func (la *LogAnalyzer) parseGenericLog(line, source string) (LogEntry, error) {
	// 简单的通用解析，提取时间戳和消息
	entry := LogEntry{
		Message:   line,
		Level:     "INFO",
		Source:    source,
		Timestamp: time.Now(),
	}

	// 尝试提取时间戳
	timeRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}[\sT]\d{2}:\d{2}:\d{2})`)
	if matches := timeRegex.FindStringSubmatch(line); len(matches) > 1 {
		if timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1]); err == nil {
			entry.Timestamp = timestamp
		} else if timestamp, err := time.Parse("2006-01-02T15:04:05", matches[1]); err == nil {
			entry.Timestamp = timestamp
		}
	}

	// 尝试提取日志级别
	levelRegex := regexp.MustCompile(`\[(ERROR|WARN|INFO|DEBUG|TRACE|FATAL)\]`)
	if matches := levelRegex.FindStringSubmatch(line); len(matches) > 1 {
		entry.Level = matches[1]
	}

	return entry, nil
}

// getLogLevelFromStatus 根据状态码确定日志级别
func (la *LogAnalyzer) getLogLevelFromStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "INFO"
	case status >= 300 && status < 400:
		return "WARN"
	case status >= 400 && status < 500:
		return "ERROR"
	case status >= 500:
		return "FATAL"
	default:
		return "INFO"
	}
}

// shouldFilterEntry 检查是否应该过滤该条目
func (la *LogAnalyzer) shouldFilterEntry(entry LogEntry) bool {
	// 按级别过滤
	if la.Config.FilterLevel != "" && entry.Level != la.Config.FilterLevel {
		return true
	}

	// 按模式过滤
	if la.Config.FilterPattern != "" {
		matched, _ := regexp.MatchString(la.Config.FilterPattern, entry.Message)
		return !matched
	}

	return false
}

// updateStats 更新统计信息
func (la *LogAnalyzer) updateStats(entry LogEntry) {
	// 更新级别统计
	la.Stats.LevelCounts[entry.Level]++

	// 更新IP统计
	if entry.IP != "" {
		la.Stats.IPCounts[entry.IP]++
	}

	// 更新状态码统计
	if entry.Status > 0 {
		statusRange := fmt.Sprintf("%dxx", entry.Status/100)
		la.Stats.StatusCounts[statusRange]++
	}

	// 更新方法统计
	if entry.Method != "" {
		la.Stats.MethodCounts[entry.Method]++
	}

	// 更新小时统计
	hourKey := entry.Timestamp.Format("2006-01-02 15")
	la.Stats.HourlyCounts[hourKey]++

	// 更新错误统计
	if entry.Level == "ERROR" || entry.Level == "FATAL" {
		la.Stats.TopErrors[entry.Message]++
	}

	// 更新时间范围
	if la.Stats.StartTime == nil || entry.Timestamp.Before(*la.Stats.StartTime) {
		la.Stats.StartTime = &entry.Timestamp
	}
	if la.Stats.EndTime == nil || entry.Timestamp.After(*la.Stats.EndTime) {
		la.Stats.EndTime = &entry.Timestamp
	}
}

// calculateDerivedStats 计算派生统计信息
func (la *LogAnalyzer) calculateDerivedStats() {
	// 计算Top IPs
	la.Stats.TopIPs = la.getTopItems(la.Stats.IPCounts, la.Config.TopN)

	// 计算时间范围
	if la.Stats.StartTime != nil && la.Stats.EndTime != nil {
		duration := la.Stats.EndTime.Sub(*la.Stats.StartTime)
		la.Stats.TimeRange = duration.String()
	}
}

// getTopItems 获取前N项统计
func (la *LogAnalyzer) getTopItems(counts map[string]int, n int) []string {
	type item struct {
		key   string
		count int
	}

	items := make([]item, 0, len(counts))
	for k, v := range counts {
		items = append(items, item{k, v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})

	if n > len(items) {
		n = len(items)
	}

	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = items[i].key
	}

	return result
}

// GenerateReport 生成分析报告
func (la *LogAnalyzer) GenerateReport() string {
	var report strings.Builder

	report.WriteString("📊 日志分析报告\n")
	report.WriteString("========================================\n\n")

	// 基础统计
	report.WriteString("📈 基础统计:\n")
	report.WriteString(fmt.Sprintf("  总行数: %d\n", la.Stats.TotalLines))
	report.WriteString(fmt.Sprintf("  有效行数: %d\n", la.Stats.ValidLines))
	report.WriteString(fmt.Sprintf("  错误行数: %d\n", la.Stats.ErrorLines))
	if la.Stats.StartTime != nil && la.Stats.EndTime != nil {
		report.WriteString(fmt.Sprintf("  时间范围: %s ~ %s\n",
			la.Stats.StartTime.Format("2006-01-02 15:04:05"),
			la.Stats.EndTime.Format("2006-01-02 15:04:05")))
		report.WriteString(fmt.Sprintf("  持续时间: %s\n", la.Stats.TimeRange))
	}

	// 日志级别统计
	if len(la.Stats.LevelCounts) > 0 {
		report.WriteString("\n🎯 日志级别统计:\n")
		for level, count := range la.Stats.LevelCounts {
			percentage := float64(count) / float64(la.Stats.ValidLines) * 100
			report.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", level, count, percentage))
		}
	}

	// IP统计
	if len(la.Stats.IPCounts) > 0 {
		report.WriteString("\n🌐 Top IP地址:\n")
		for i, ip := range la.Stats.TopIPs {
			if i >= 10 { // 限制显示前10个
				break
			}
			count := la.Stats.IPCounts[ip]
			percentage := float64(count) / float64(la.Stats.ValidLines) * 100
			report.WriteString(fmt.Sprintf("  %d. %s: %d (%.1f%%)\n", i+1, ip, count, percentage))
		}
	}

	// HTTP状态码统计
	if len(la.Stats.StatusCounts) > 0 {
		report.WriteString("\n📊 HTTP状态码统计:\n")
		for status, count := range la.Stats.StatusCounts {
			percentage := float64(count) / float64(la.Stats.ValidLines) * 100
			report.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", status, count, percentage))
		}
	}

	// HTTP方法统计
	if len(la.Stats.MethodCounts) > 0 {
		report.WriteString("\n🔧 HTTP方法统计:\n")
		for method, count := range la.Stats.MethodCounts {
			percentage := float64(count) / float64(la.Stats.ValidLines) * 100
			report.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", method, count, percentage))
		}
	}

	// 错误统计
	if len(la.Stats.TopErrors) > 0 {
		report.WriteString("\n❌ Top错误信息:\n")
		type errorItem struct {
			message string
			count   int
		}

		errors := make([]errorItem, 0, len(la.Stats.TopErrors))
		for msg, count := range la.Stats.TopErrors {
			errors = append(errors, errorItem{msg, count})
		}

		sort.Slice(errors, func(i, j int) bool {
			return errors[i].count > errors[j].count
		})

		for i, err := range errors {
			if i >= 5 { // 限制显示前5个错误
				break
			}
			report.WriteString(fmt.Sprintf("  %d. %s (%d次)\n", i+1, err.message, err.count))
		}
	}

	report.WriteString("\n========================================\n")
	report.WriteString("分析完成! 🎉\n")

	return report.String()
}

// ExportJSON 导出JSON格式报告
func (la *LogAnalyzer) ExportJSON(filename string) error {
	data, err := json.MarshalIndent(struct {
		Stats   LogStats   `json:"stats"`
		Entries []LogEntry `json:"entries,omitempty"`
	}{
		Stats: la.Stats,
		Entries: func() []LogEntry {
			if la.Config.ShowDetails {
				return la.Entries
			}
			return nil
		}(),
	}, "", "  ")

	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// 主函数
func main() {
	// 命令行参数
	var (
		logFile       = flag.String("file", "", "日志文件路径")
		logFormat     = flag.String("format", "auto", "日志格式 (apache/nginx/common/syslog/json/auto)")
		filterLevel   = flag.String("level", "", "过滤日志级别 (ERROR/WARN/INFO/DEBUG)")
		filterPattern = flag.String("pattern", "", "过滤模式 (正则表达式)")
		outputFormat  = flag.String("output", "text", "输出格式 (text/json)")
		outputFile    = flag.String("out", "", "输出文件路径")
		topN          = flag.Int("top", 10, "显示前N项统计")
		showDetails   = flag.Bool("details", false, "显示详细信息")
		showHelp      = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *showHelp || *logFile == "" {
		fmt.Println("🔍 日志分析器 - 使用帮助")
		fmt.Println("========================================")
		fmt.Println("用法: log_analyzer [选项]")
		fmt.Println("\n选项:")
		flag.PrintDefaults()
		fmt.Println("\n支持的日志格式:")
		fmt.Println("  apache  - Apache访问日志")
		fmt.Println("  nginx   - Nginx访问日志")
		fmt.Println("  common  - 通用日志格式 (时间 [级别] 消息)")
		fmt.Println("  syslog  - 系统日志格式")
		fmt.Println("  json    - JSON格式日志")
		fmt.Println("  auto    - 自动检测格式")
		fmt.Println("\n示例:")
		fmt.Println("  log_analyzer -file access.log -format nginx")
		fmt.Println("  log_analyzer -file app.log -level ERROR -output json")
		fmt.Println("  log_analyzer -file system.log -pattern \"database\" -top 5")
		return
	}

	// 创建分析器配置
	config := AnalyzerConfig{
		LogFormat:     *logFormat,
		FilterLevel:   *filterLevel,
		FilterPattern: *filterPattern,
		OutputFormat:  *outputFormat,
		TopN:          *topN,
		ShowDetails:   *showDetails,
	}

	// 创建分析器
	analyzer := NewLogAnalyzer(config)

	fmt.Printf("🔍 开始分析日志文件: %s\n", *logFile)
	fmt.Printf("📋 使用格式: %s\n", *logFormat)

	// 解析日志文件
	if err := analyzer.ParseLogFile(*logFile); err != nil {
		fmt.Printf("❌ 解析失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 解析完成! 处理了 %d 行日志\n\n", analyzer.Stats.TotalLines)

	// 生成报告
	var output string
	if *outputFormat == "json" {
		if *outputFile != "" {
			if err := analyzer.ExportJSON(*outputFile); err != nil {
				fmt.Printf("❌ 导出JSON失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ JSON报告已保存到: %s\n", *outputFile)
			return
		} else {
			data, _ := json.MarshalIndent(analyzer.Stats, "", "  ")
			output = string(data)
		}
	} else {
		output = analyzer.GenerateReport()
	}

	// 输出结果
	if *outputFile != "" && *outputFormat != "json" {
		if err := os.WriteFile(*outputFile, []byte(output), 0644); err != nil {
			fmt.Printf("❌ 保存报告失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 报告已保存到: %s\n", *outputFile)
	} else {
		fmt.Print(output)
	}
}
