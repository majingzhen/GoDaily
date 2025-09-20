package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SystemInfo struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       CPUInfo   `json:"cpu"`
	Memory    MemInfo   `json:"memory"`
	Disk      DiskInfo  `json:"disk"`
	System    SysInfo   `json:"system"`
}

type CPUInfo struct {
	Usage   float64 `json:"usage"`
	Cores   int     `json:"cores"`
	LoadAvg float64 `json:"load_avg"`
}

type MemInfo struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Usage     float64 `json:"usage"`
}

type DiskInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
}

type SysInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	NumCPU   int    `json:"num_cpu"`
	Hostname string `json:"hostname"`
}

func main() {
	var (
		interval = flag.Duration("interval", 5*time.Second, "监控间隔时间")
		count    = flag.Int("count", 1, "监控次数 (-1 表示持续监控)")
		output   = flag.String("output", "console", "输出格式: console, json")
		file     = flag.String("file", "", "输出到文件")
		help     = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	var outputFile *os.File
	var err error

	if *file != "" {
		outputFile, err = os.Create(*file)
		if err != nil {
			fmt.Printf("创建输出文件失败: %v\n", err)
			return
		}
		defer outputFile.Close()
	}

	monitorCount := 0
	for {
		if *count > 0 && monitorCount >= *count {
			break
		}

		info, err := getSystemInfo()
		if err != nil {
			fmt.Printf("获取系统信息失败: %v\n", err)
			time.Sleep(*interval)
			continue
		}

		switch *output {
		case "json":
			outputJSON(info, outputFile)
		default:
			outputConsole(info, outputFile)
		}

		monitorCount++
		if *count == -1 || (*count > 1 && monitorCount < *count) {
			time.Sleep(*interval)
		}
	}
}

func showHelp() {
	fmt.Println("系统监控工具 - 实时监控系统资源使用情况")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  system_monitor [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -interval duration  监控间隔时间 (默认: 5s)")
	fmt.Println("  -count int          监控次数, -1表示持续监控 (默认: 1)")
	fmt.Println("  -output string      输出格式: console, json (默认: console)")
	fmt.Println("  -file string        输出到文件")
	fmt.Println("  -help               显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  system_monitor                                    # 单次监控")
	fmt.Println("  system_monitor -count -1 -interval 2s             # 持续监控，2秒间隔")
	fmt.Println("  system_monitor -output json -file monitor.json    # JSON格式输出到文件")
	fmt.Println("  system_monitor -count 10 -interval 1s             # 监控10次，1秒间隔")
}

func getSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		Timestamp: time.Now(),
	}

	var err error

	info.CPU, err = getCPUInfo()
	if err != nil {
		return nil, fmt.Errorf("获取CPU信息失败: %v", err)
	}

	info.Memory, err = getMemoryInfo()
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %v", err)
	}

	info.Disk, err = getDiskInfo()
	if err != nil {
		return nil, fmt.Errorf("获取磁盘信息失败: %v", err)
	}

	info.System = getBasicSystemInfo()

	return info, nil
}

func getCPUInfo() (CPUInfo, error) {
	cpu := CPUInfo{
		Cores: runtime.NumCPU(),
	}

	if runtime.GOOS == "linux" {
		usage, err := getCPUUsageLinux()
		if err == nil {
			cpu.Usage = usage
		}

		loadAvg, err := getLoadAverageLinux()
		if err == nil {
			cpu.LoadAvg = loadAvg
		}
	} else if runtime.GOOS == "windows" {
		usage, err := getCPUUsageWindows()
		if err == nil {
			cpu.Usage = usage
		}
	}

	return cpu, nil
}

func getCPUUsageLinux() (float64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return 0, fmt.Errorf("无法读取CPU统计信息")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 8 {
		return 0, fmt.Errorf("CPU统计信息格式错误")
	}

	var total, idle uint64
	for i := 1; i < len(fields); i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return 0, err
		}
		total += val
		if i == 4 {
			idle = val
		}
	}

	if total == 0 {
		return 0, nil
	}

	usage := float64(total-idle) / float64(total) * 100
	return usage, nil
}

func getCPUUsageWindows() (float64, error) {
	return 0, fmt.Errorf("Windows CPU使用率监控需要额外的系统API")
}

func getLoadAverageLinux() (float64, error) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return 0, fmt.Errorf("无法读取负载平均值")
	}

	fields := strings.Fields(scanner.Text())
	if len(fields) < 1 {
		return 0, fmt.Errorf("负载平均值格式错误")
	}

	loadAvg, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return loadAvg, nil
}

func getMemoryInfo() (MemInfo, error) {
	mem := MemInfo{}

	if runtime.GOOS == "linux" {
		return getMemoryInfoLinux()
	} else if runtime.GOOS == "windows" {
		return getMemoryInfoWindows()
	}

	return mem, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}

func getMemoryInfoLinux() (MemInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemInfo{}, err
	}
	defer file.Close()

	mem := MemInfo{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		value *= 1024

		switch key {
		case "MemTotal":
			mem.Total = value
		case "MemAvailable":
			mem.Available = value
		}
	}

	mem.Used = mem.Total - mem.Available
	if mem.Total > 0 {
		mem.Usage = float64(mem.Used) / float64(mem.Total) * 100
	}

	return mem, nil
}

func getMemoryInfoWindows() (MemInfo, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	mem := MemInfo{
		Used: m.Sys,
	}

	return mem, fmt.Errorf("Windows内存监控需要额外的系统API")
}

func getDiskInfo() (DiskInfo, error) {
	disk := DiskInfo{}

	if runtime.GOOS == "linux" {
		return getDiskInfoLinux()
	} else if runtime.GOOS == "windows" {
		return getDiskInfoWindows()
	}

	return disk, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
}

func getDiskInfoLinux() (DiskInfo, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return DiskInfo{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "/" {
			return getDiskUsage("/")
		}
	}

	return getDiskUsage("/")
}

func getDiskInfoWindows() (DiskInfo, error) {
	return getDiskUsage("C:\\")
}

func getDiskUsage(path string) (DiskInfo, error) {
	disk := DiskInfo{}

	file, err := os.Open(path)
	if err != nil {
		return disk, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return disk, err
	}

	if !stat.IsDir() {
		return disk, fmt.Errorf("%s 不是目录", path)
	}

	return disk, fmt.Errorf("磁盘使用率监控需要系统调用支持")
}

func getBasicSystemInfo() SysInfo {
	hostname, _ := os.Hostname()

	return SysInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		NumCPU:   runtime.NumCPU(),
		Hostname: hostname,
	}
}

func outputConsole(info *SystemInfo, file *os.File) {
	output := fmt.Sprintf(`
========== 系统监控报告 ==========
时间: %s
主机: %s
系统: %s/%s

CPU 信息:
  核心数: %d
  使用率: %.2f%%
  负载平均值: %.2f

内存信息:
  总内存: %s
  已用内存: %s
  可用内存: %s
  使用率: %.2f%%

磁盘信息:
  总容量: %s
  已用容量: %s
  可用容量: %s
  使用率: %.2f%%

================================
`,
		info.Timestamp.Format("2006-01-02 15:04:05"),
		info.System.Hostname,
		info.System.OS,
		info.System.Arch,
		info.CPU.Cores,
		info.CPU.Usage,
		info.CPU.LoadAvg,
		formatBytes(info.Memory.Total),
		formatBytes(info.Memory.Used),
		formatBytes(info.Memory.Available),
		info.Memory.Usage,
		formatBytes(info.Disk.Total),
		formatBytes(info.Disk.Used),
		formatBytes(info.Disk.Available),
		info.Disk.Usage,
	)

	if file != nil {
		file.WriteString(output)
	} else {
		fmt.Print(output)
	}
}

func outputJSON(info *SystemInfo, file *os.File) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	if file != nil {
		file.Write(data)
		file.WriteString("\n")
	} else {
		fmt.Println(string(data))
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
