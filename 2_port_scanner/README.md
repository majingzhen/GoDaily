# 端口扫描器 (Port Scanner)

一个高效的多线程端口扫描工具，用于检测目标主机的开放端口。

## 🔍 功能特性

- **并发扫描**：使用goroutines实现高并发扫描，提高扫描效率
- **可配置参数**：支持自定义主机地址、端口范围、超时时间和并发数
- **信号量控制**：通过带缓冲通道控制并发数量，避免资源耗尽
- **实时输出**：发现开放端口时立即显示结果
- **性能统计**：显示扫描耗时和端口统计信息
- **结果排序**：开放端口按数字顺序排列显示

## 📋 技术特点

### 核心概念
- **并发编程**：使用goroutines和sync.WaitGroup实现并发控制
- **网络编程**：使用net.DialTimeout进行TCP连接测试
- **信号量模式**：通过带缓冲通道限制并发数量
- **互斥锁**：使用sync.Mutex保护共享资源(openPorts切片)
- **命令行参数**：使用flag包处理命令行选项

### 程序流程
1. 解析命令行参数
2. 验证端口范围有效性
3. 创建信号量通道控制并发
4. 启动goroutines并发扫描端口
5. 收集开放端口结果
6. 排序并显示扫描结果

## 🚀 使用方法

### 基本用法
```bash
# 编译程序
go build -o port_scanner port_scanner.go

# 扫描本地主机的常用端口
./port_scanner

# 扫描指定主机
./port_scanner -host=192.168.1.1

# 扫描指定端口范围
./port_scanner -host=google.com -start=80 -end=443
```

### 参数说明
```bash
-host string
    要扫描的主机地址 (默认: "localhost")
    
-start int
    要扫描的起始端口 (默认: 1)
    
-end int
    要扫描的结束端口 (默认: 1024)
    
-timeout int
    扫描超时时间，单位毫秒 (默认: 500)
    
-con int
    扫描的并发数 (默认: 100)
```

### 使用示例

**示例1：快速扫描本地常用端口**
```bash
./port_scanner -start=1 -end=1000 -concurrency=50
```

**示例2：扫描远程服务器的Web端口**
```bash
./port_scanner -host=example.com -start=80 -end=8080 -timeout=2000
```

**示例3：高并发扫描大范围端口**
```bash
./port_scanner -host=192.168.1.1 -start=1 -end=65535 -con=500 -timeout=500
```

## 📊 输出示例

```
开始扫描 localhost 的端口范围 1-1024...
扫描超时为 1000 毫秒，并发数为 100
端口 22 已开放
端口 80 已开放
端口 443 已开放
端口 3306 已开放
扫描完成，耗时 2.354s
共扫描 1024 个端口
开放端口列表：
22
80
443
3306
```

## ⚡ 性能优化

### 并发控制
- 使用信号量模式限制并发数量，避免创建过多goroutines
- 通过sync.WaitGroup确保所有goroutines完成后再退出
- 使用互斥锁保护共享资源，避免竞态条件

### 参数调优建议
- **并发数**：根据系统资源和网络带宽调整，通常50-500之间
- **超时时间**：本地网络500-1000ms，远程网络1000-3000ms
- **端口范围**：按需扫描，避免扫描不必要的端口范围

## 🔧 代码解析

### 关键代码段

**1. 并发控制和信号量**
```go
// 创建带缓冲的通道控制并发数量
semaphore := make(chan struct{}, *concurrency)

// 在goroutine中获取和释放信号量
semaphore <- struct{}{}        // 获取信号量
defer func() { <-semaphore }() // 释放信号量
```

**2. TCP端口检测**
```go
address := fmt.Sprintf("%s:%d", *host, port)
conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Millisecond)
if err == nil {
    conn.Close()
    // 端口开放
}
```

**3. 线程安全的结果收集**
```go
var mu sync.Mutex
mu.Lock()
openPorts = append(openPorts, port)
mu.Unlock()
```

## 🛡️ 安全提醒

- 仅用于合法的网络诊断和安全测试
- 不要对未授权的主机进行扫描
- 高并发扫描可能被目标系统识别为恶意行为
- 建议在测试环境或本地网络中使用

## 📈 学习价值

通过这个项目可以学习到：

1. **Go并发编程**：goroutines、channels、sync包的使用
2. **网络编程基础**：TCP连接、超时处理
3. **系统编程**：命令行参数处理、错误处理
4. **性能优化**：并发控制、资源管理
5. **安全意识**：网络扫描的合理使用

## 🔄 扩展思路

- 添加UDP端口扫描支持
- 实现端口服务识别功能
- 添加扫描结果导出功能
- 支持从文件读取目标列表
- 添加进度条显示
- 实现更智能的超时策略