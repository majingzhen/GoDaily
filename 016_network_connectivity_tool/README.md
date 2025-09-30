# 网络连通性测试工具 (016)

一个功能强大的网络连通性测试工具，支持 PING 测试、TCP/UDP 端口检测和端口扫描功能。

## ✨ 功能特性

- **多种测试模式**: 支持 PING、TCP、UDP 和端口扫描
- **并发检测**: 支持多线程并发测试，提高检测效率
- **灵活端口配置**: 支持单个端口、端口列表和端口范围扫描
- **多种输出格式**: 支持控制台友好格式和 JSON 格式输出
- **结果保存**: 支持将测试结果保存到文件
- **详细信息**: 提供连接延迟、错误信息等详细数据
- **安全限制**: 内置扫描范围限制，防止过度扫描

## 🛠️ 技术特性

- **高并发**: 使用 goroutines 和信号量控制并发数
- **超时控制**: 可配置连接超时时间
- **错误处理**: 完善的错误处理和报告机制
- **跨平台**: 支持 Windows、Linux 和 macOS
- **内存友好**: 优化的并发模型，避免资源耗尽

## 🚀 使用方法

### 基本用法

```bash
# PING 测试 (默认)
network_connectivity_tool -host google.com

# 指定主机和模式
network_connectivity_tool -host example.com -mode ping
```

### TCP 端口测试

```bash
# 测试常用端口
network_connectivity_tool -host example.com -mode tcp

# 测试指定端口
network_connectivity_tool -host example.com -mode tcp -ports 80,443,22,25

# 设置超时和并发数
network_connectivity_tool -host example.com -mode tcp -timeout 5s -threads 100
```

### UDP 端口测试

```bash
# 测试 DNS 端口
network_connectivity_tool -host 8.8.8.8 -mode udp -ports 53

# 测试多个 UDP 服务
network_connectivity_tool -host example.com -mode udp -ports 53,123,161
```

### 端口扫描

```bash
# 扫描常用端口
network_connectivity_tool -host 192.168.1.1 -mode scan

# 扫描端口范围
network_connectivity_tool -host example.com -mode scan -range 1-1000

# 快速扫描
network_connectivity_tool -host example.com -mode scan -range 1-100 -threads 200
```

### 输出和保存

```bash
# JSON 格式输出
network_connectivity_tool -host example.com -mode scan -output json

# 保存结果到文件
network_connectivity_tool -host example.com -mode scan -output json -file scan_result.json

# 详细输出模式
network_connectivity_tool -host example.com -mode tcp -verbose
```

## 📋 命令行选项

| 选项 | 默认值 | 描述 |
|------|--------|------|
| `-host` | `8.8.8.8` | 目标主机地址 |
| `-mode` | `ping` | 检测模式：ping, tcp, udp, scan |
| `-ports` | `80,443,22,21,25,53,110,993,995` | 端口列表(逗号分隔) |
| `-range` | | 端口范围(如: 1-1000) |
| `-timeout` | `3s` | 连接超时时间 |
| `-output` | `console` | 输出格式：console, json |
| `-file` | | 保存结果到文件 |
| `-threads` | `50` | 并发线程数 |
| `-verbose` | `false` | 详细输出 |
| `-help` | `false` | 显示帮助信息 |

## 📊 输出示例

### 控制台输出 - PING 测试

```
========== 网络连通性测试结果 ==========
时间: 2024-01-15 10:30:45
主机: google.com
协议: PING
延迟: 15.234ms
状态: ✅ 连接成功
=====================================
```

### 控制台输出 - 端口测试

```
========== 端口连通性测试结果 ==========
时间: 2024-01-15 10:30:45
主机: example.com
测试端口数: 9
开放端口: 3
关闭端口: 6

详细结果:
✅ 端口 80/TCP - 延迟: 25.123ms
✅ 端口 443/TCP - 延迟: 30.456ms
❌ 端口 22/TCP - 延迟: 3.001s
✅ 端口 21/TCP - 延迟: 45.789ms
❌ 端口 25/TCP - 延迟: 3.001s
====================================
```

### 控制台输出 - 端口扫描

```
========== 端口扫描结果 ==========
时间: 2024-01-15 10:30:45
主机: 192.168.1.1
扫描端口总数: 1000
开放端口数: 5
关闭端口数: 995

开放端口:
  ✅ 22/tcp
  ✅ 53/tcp
  ✅ 80/tcp
  ✅ 443/tcp
  ✅ 8080/tcp

================================
```

### JSON 输出示例

```json
{
  "timestamp": "2024-01-15T10:30:45.123456789+08:00",
  "host": "example.com",
  "total_ports": 3,
  "open_ports": [80, 443],
  "closed_ports": [22],
  "results": [
    {
      "timestamp": "2024-01-15T10:30:45.123456789+08:00",
      "host": "example.com",
      "port": 80,
      "success": true,
      "latency": 25123000,
      "type": "tcp"
    },
    {
      "timestamp": "2024-01-15T10:30:45.123456789+08:00",
      "host": "example.com",
      "port": 443,
      "success": true,
      "latency": 30456000,
      "type": "tcp"
    },
    {
      "timestamp": "2024-01-15T10:30:45.123456789+08:00",
      "host": "example.com",
      "port": 22,
      "success": false,
      "latency": 3001000000,
      "error": "dial tcp example.com:22: i/o timeout",
      "type": "tcp"
    }
  ]
}
```

## 🏗️ 构建和运行

### 构建

```bash
cd 016_network_connectivity_tool
go build -o network_connectivity_tool network_connectivity_tool.go
```

### Windows

```bash
go build -o network_connectivity_tool.exe network_connectivity_tool.go
```

### 运行

```bash
# 显示帮助
./network_connectivity_tool -help

# 基本测试
./network_connectivity_tool -host google.com

# 端口扫描
./network_connectivity_tool -host 192.168.1.1 -mode scan -range 1-100
```

## 📈 性能说明

- **并发控制**: 默认并发数为 50，可根据系统性能调整
- **超时设置**: 默认超时 3 秒，可根据网络环境调整
- **扫描限制**: 单次扫描最多 10000 个端口，防止过度消耗资源
- **内存使用**: 每个并发连接占用约 8KB 内存

## 🔒 安全注意事项

- **合法使用**: 仅对自己拥有或有权限测试的主机进行扫描
- **扫描频率**: 避免高频扫描，可能被目标主机误认为攻击
- **端口范围**: 大范围扫描会消耗大量网络和系统资源
- **防火墙**: 某些防火墙可能会阻止或记录端口扫描行为

## 🛠️ 技术实现

- **网络库**: 使用 Go 标准库 `net` 包进行网络连接
- **并发控制**: 使用 goroutines 和信号量模式控制并发数
- **超时处理**: 使用 `DialTimeout` 实现连接超时控制
- **数据结构**: 使用结构体和 JSON 标签支持多格式输出
- **错误处理**: 完善的错误捕获和传播机制

## 📝 使用场景

- **网络诊断**: 检查网络连通性和服务可用性
- **服务监控**: 监控关键服务端口的开放状态
- **安全审计**: 检查服务器暴露的端口和服务
- **网络排查**: 排查网络连接问题和服务故障
- **批量检测**: 批量检测多个端口的连通性

## ⚠️ 限制说明

- **ICMP PING**: 由于权限限制，PING 模式使用 TCP 连接模拟
- **UDP 检测**: UDP 连接检测准确性有限，建议结合应用层测试
- **防火墙**: 防火墙规则可能影响检测结果的准确性
- **系统限制**: 操作系统的文件描述符限制可能影响大范围扫描