package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ConnectivityResult struct {
	Timestamp time.Time     `json:"timestamp"`
	Host      string        `json:"host"`
	Port      int           `json:"port,omitempty"`
	Success   bool          `json:"success"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	Type      string        `json:"type"` // "ping", "tcp", "udp"
}

type ScanResult struct {
	Timestamp   time.Time            `json:"timestamp"`
	Host        string               `json:"host"`
	TotalPorts  int                  `json:"total_ports"`
	OpenPorts   []int                `json:"open_ports"`
	ClosedPorts []int                `json:"closed_ports"`
	Results     []ConnectivityResult `json:"results"`
}

type NetworkTool struct {
	host    string
	ports   []int
	timeout time.Duration
	output  string
	file    string
	mode    string
	threads int
	verbose bool
}

func main() {
	tool := &NetworkTool{}

	flag.StringVar(&tool.host, "host", "8.8.8.8", "目标主机地址")
	flag.StringVar(&tool.mode, "mode", "ping", "检测模式: ping, tcp, udp, scan")
	portsStr := flag.String("ports", "80,443,22,21,25,53,110,993,995", "端口列表(逗号分隔)")
	portRange := flag.String("range", "", "端口范围(如: 1-1000)")
	flag.DurationVar(&tool.timeout, "timeout", 3*time.Second, "连接超时时间")
	flag.StringVar(&tool.output, "output", "console", "输出格式: console, json")
	flag.StringVar(&tool.file, "file", "", "保存结果到文件")
	flag.IntVar(&tool.threads, "threads", 50, "并发线程数")
	flag.BoolVar(&tool.verbose, "verbose", false, "详细输出")
	help := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 解析端口
	if *portRange != "" {
		tool.ports = parsePortRange(*portRange)
	} else {
		tool.ports = parsePorts(*portsStr)
	}

	switch tool.mode {
	case "ping":
		result := tool.pingHost()
		tool.outputResult(result)
	case "tcp":
		results := tool.testTCPPorts()
		tool.outputResults(results)
	case "udp":
		results := tool.testUDPPorts()
		tool.outputResults(results)
	case "scan":
		scanResult := tool.scanPorts()
		tool.outputScanResult(scanResult)
	default:
		fmt.Printf("错误: 不支持的模式 '%s'\n", tool.mode)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`网络连通性测试工具 (016)

使用方法:
  network_connectivity_tool [选项]

选项:
  -host string        目标主机地址 (默认: "8.8.8.8")
  -mode string        检测模式: ping, tcp, udp, scan (默认: "ping")
  -ports string       端口列表，逗号分隔 (默认: "80,443,22,21,25,53,110,993,995")
  -range string       端口范围，如: 1-1000
  -timeout duration   连接超时时间 (默认: 3s)
  -output string      输出格式: console, json (默认: "console")
  -file string        保存结果到文件
  -threads int        并发线程数 (默认: 50)
  -verbose            详细输出
  -help               显示此帮助信息

示例:
  # PING 测试
  network_connectivity_tool -host google.com -mode ping

  # TCP 端口测试
  network_connectivity_tool -host example.com -mode tcp -ports 80,443,22

  # 端口扫描
  network_connectivity_tool -host 192.168.1.1 -mode scan -range 1-1000

  # UDP 测试
  network_connectivity_tool -host 8.8.8.8 -mode udp -ports 53,123

  # JSON 输出并保存到文件
  network_connectivity_tool -host example.com -mode scan -output json -file result.json`)
}

func (nt *NetworkTool) pingHost() ConnectivityResult {
	start := time.Now()

	// 对于 ICMP ping，我们使用 TCP 连接测试作为替代
	// 因为 ICMP 需要特殊权限
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(nt.host, "80"), nt.timeout)

	result := ConnectivityResult{
		Timestamp: time.Now(),
		Host:      nt.host,
		Type:      "ping",
		Latency:   time.Since(start),
	}

	if err != nil {
		// 尝试其他常用端口
		ports := []string{"443", "22", "21"}
		for _, port := range ports {
			conn, err = net.DialTimeout("tcp", net.JoinHostPort(nt.host, port), nt.timeout)
			if err == nil {
				break
			}
		}
	}

	if err == nil {
		result.Success = true
		if conn != nil {
			conn.Close()
		}
	} else {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

func (nt *NetworkTool) testTCPPorts() []ConnectivityResult {
	var results []ConnectivityResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, nt.threads)

	for _, port := range nt.ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := nt.testTCPPort(p)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

func (nt *NetworkTool) testTCPPort(port int) ConnectivityResult {
	start := time.Now()
	address := net.JoinHostPort(nt.host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, nt.timeout)

	result := ConnectivityResult{
		Timestamp: time.Now(),
		Host:      nt.host,
		Port:      port,
		Type:      "tcp",
		Latency:   time.Since(start),
	}

	if err == nil {
		result.Success = true
		conn.Close()
	} else {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

func (nt *NetworkTool) testUDPPorts() []ConnectivityResult {
	var results []ConnectivityResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	semaphore := make(chan struct{}, nt.threads)

	for _, port := range nt.ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := nt.testUDPPort(p)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return results
}

func (nt *NetworkTool) testUDPPort(port int) ConnectivityResult {
	start := time.Now()
	address := net.JoinHostPort(nt.host, strconv.Itoa(port))

	conn, err := net.DialTimeout("udp", address, nt.timeout)

	result := ConnectivityResult{
		Timestamp: time.Now(),
		Host:      nt.host,
		Port:      port,
		Type:      "udp",
		Latency:   time.Since(start),
	}

	if err == nil {
		result.Success = true
		// UDP 连接不会立即返回错误，所以我们假设成功
		// 实际的 UDP 测试需要发送数据包并等待响应
		conn.Close()
	} else {
		result.Success = false
		result.Error = err.Error()
	}

	return result
}

func (nt *NetworkTool) scanPorts() ScanResult {
	results := nt.testTCPPorts()

	var openPorts, closedPorts []int

	for _, result := range results {
		if result.Success {
			openPorts = append(openPorts, result.Port)
		} else {
			closedPorts = append(closedPorts, result.Port)
		}
	}

	return ScanResult{
		Timestamp:   time.Now(),
		Host:        nt.host,
		TotalPorts:  len(nt.ports),
		OpenPorts:   openPorts,
		ClosedPorts: closedPorts,
		Results:     results,
	}
}

func (nt *NetworkTool) outputResult(result ConnectivityResult) {
	if nt.output == "json" {
		nt.outputJSON(result)
	} else {
		nt.outputConsole(result)
	}
}

func (nt *NetworkTool) outputResults(results []ConnectivityResult) {
	if nt.output == "json" {
		nt.outputJSON(results)
	} else {
		nt.outputConsoleResults(results)
	}
}

func (nt *NetworkTool) outputScanResult(result ScanResult) {
	if nt.output == "json" {
		nt.outputJSON(result)
	} else {
		nt.outputConsoleScan(result)
	}
}

func (nt *NetworkTool) outputConsole(result ConnectivityResult) {
	fmt.Println("========== 网络连通性测试结果 ==========")
	fmt.Printf("时间: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("主机: %s\n", result.Host)

	if result.Port > 0 {
		fmt.Printf("端口: %d\n", result.Port)
	}

	fmt.Printf("协议: %s\n", strings.ToUpper(result.Type))
	fmt.Printf("延迟: %v\n", result.Latency)

	if result.Success {
		fmt.Println("状态: ✅ 连接成功")
	} else {
		fmt.Println("状态: ❌ 连接失败")
		if result.Error != "" {
			fmt.Printf("错误: %s\n", result.Error)
		}
	}
	fmt.Println("=====================================")
}

func (nt *NetworkTool) outputConsoleResults(results []ConnectivityResult) {
	fmt.Println("========== 端口连通性测试结果 ==========")
	fmt.Printf("时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("主机: %s\n", nt.host)
	fmt.Printf("测试端口数: %d\n", len(results))

	openCount := 0
	for _, result := range results {
		if result.Success {
			openCount++
		}
	}

	fmt.Printf("开放端口: %d\n", openCount)
	fmt.Printf("关闭端口: %d\n", len(results)-openCount)
	fmt.Println()

	fmt.Println("详细结果:")
	for _, result := range results {
		status := "❌"
		if result.Success {
			status = "✅"
		}

		fmt.Printf("%s 端口 %d/%s - 延迟: %v", status, result.Port,
			strings.ToUpper(result.Type), result.Latency)

		if !result.Success && nt.verbose {
			fmt.Printf(" (%s)", result.Error)
		}
		fmt.Println()
	}
	fmt.Println("====================================")
}

func (nt *NetworkTool) outputConsoleScan(result ScanResult) {
	fmt.Println("========== 端口扫描结果 ==========")
	fmt.Printf("时间: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("主机: %s\n", result.Host)
	fmt.Printf("扫描端口总数: %d\n", result.TotalPorts)
	fmt.Printf("开放端口数: %d\n", len(result.OpenPorts))
	fmt.Printf("关闭端口数: %d\n", len(result.ClosedPorts))
	fmt.Println()

	if len(result.OpenPorts) > 0 {
		fmt.Println("开放端口:")
		for _, port := range result.OpenPorts {
			fmt.Printf("  ✅ %d/tcp\n", port)
		}
		fmt.Println()
	}

	if nt.verbose && len(result.ClosedPorts) > 0 {
		fmt.Println("关闭端口:")
		for _, port := range result.ClosedPorts {
			fmt.Printf("  ❌ %d/tcp\n", port)
		}
		fmt.Println()
	}

	fmt.Println("================================")
}

func (nt *NetworkTool) outputJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("JSON 序列化错误: %v\n", err)
		return
	}

	output := string(jsonData)

	if nt.file != "" {
		err = os.WriteFile(nt.file, jsonData, 0644)
		if err != nil {
			fmt.Printf("文件写入错误: %v\n", err)
		} else {
			fmt.Printf("结果已保存到: %s\n", nt.file)
		}
	} else {
		fmt.Println(output)
	}
}

func parsePorts(portsStr string) []int {
	var ports []int

	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if port, err := strconv.Atoi(part); err == nil {
			if port > 0 && port <= 65535 {
				ports = append(ports, port)
			}
		}
	}

	return ports
}

func parsePortRange(rangeStr string) []int {
	var ports []int

	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return ports
	}

	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

	if err1 != nil || err2 != nil || start < 1 || end > 65535 || start > end {
		return ports
	}

	// 限制扫描范围避免过大
	if end-start > 10000 {
		fmt.Printf("警告: 端口范围过大，限制为前 10000 个端口\n")
		end = start + 9999
	}

	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}

	return ports
}
