package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "localhost", "要扫描的主机地址")
	startPort := flag.Int("start", 1, "要扫描的起始端口")
	endPort := flag.Int("end", 1024, "要扫描的结束端口")
	timeout := flag.Int("timeout", 500, "扫描超时时间（毫秒）")
	concurrency := flag.Int("con", 100, "扫描的并发数")
	flag.Parse()

	// 验证端口范围
	if *startPort < 1 || *startPort > *endPort || *endPort > 65535 {
		fmt.Println("无效的端口范围，请确保 1 <= 起始端口 <= 结束端口 <= 65535")
		os.Exit(0)
	}

	fmt.Printf("开始扫描 %s 的端口范围 %d-%d...\n", *host, *startPort, *endPort)
	fmt.Printf("扫描超时为 %d 毫秒，并发数为 %d\n", *timeout, *concurrency)

	startTime := time.Now()

	// 创建带缓冲的通道控制并发数量
	semaphore := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	var openPorts []int
	var mu sync.Mutex

	// 遍历端口范围
	for port := *startPort; port <= *endPort; port++ {
		semaphore <- struct{}{} // 获取信号量
		wg.Add(1)

		go func(port int) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			address := fmt.Sprintf("%s:%d", *host, port)
			conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Millisecond)

			if err == nil {
				conn.Close()
				mu.Lock()
				openPorts = append(openPorts, port)
				mu.Unlock()
				fmt.Printf("端口 %d 已开放\n", port)
			}
		}(port)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// 排序开放端口
	sort.Ints(openPorts)

	fmt.Printf("扫描完成，耗时 %s\n", duration)
	fmt.Printf("共扫描 %d 个端口\n", *endPort-*startPort+1)

	if len(openPorts) > 0 {
		fmt.Println("开放端口列表：")
		for _, port := range openPorts {
			fmt.Println(port)
		}
	} else {
		fmt.Println("没有开放端口")
	}
}
