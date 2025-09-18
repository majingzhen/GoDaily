# 日志分析器

一个使用Go语言实现的强大日志文件分析工具，支持多种日志格式解析、统计分析和报告生成。

## 功能特性

- ✅ **多格式支持**: 支持Apache、Nginx、系统日志、JSON等多种格式
- ✅ **智能解析**: 自动识别日志格式和时间戳
- ✅ **统计分析**: 提供全面的日志统计和趋势分析
- ✅ **过滤功能**: 支持按级别、模式过滤日志条目
- ✅ **报告生成**: 生成详细的文本和JSON格式报告
- ✅ **错误分析**: 识别和统计错误日志模式
- ✅ **性能统计**: HTTP状态码、响应时间等性能指标
- ✅ **IP分析**: 访问IP统计和排名
- ✅ **时间分析**: 按小时统计访问趋势

## 技术特点

- 使用Go标准库实现文件处理和正则表达式
- 支持大文件流式处理，内存占用低
- 灵活的配置系统和命令行参数
- 完整的错误处理和输入验证
- 支持JSON格式数据导出
- 高性能的统计算法和数据结构

## 使用方法

### 编译程序
```bash
go build log_analyzer.go
```

### 基本用法

#### 1. 分析Web服务器日志
```bash
# 分析Apache访问日志
./log_analyzer -file access.log -format apache

# 分析Nginx访问日志  
./log_analyzer -file access.log -format nginx

# 自动检测格式
./log_analyzer -file access.log -format auto
```

#### 2. 过滤和分析
```bash
# 只分析错误日志
./log_analyzer -file app.log -level ERROR

# 按模式过滤
./log_analyzer -file app.log -pattern "database|sql"

# 显示前20项统计
./log_analyzer -file access.log -top 20
```

#### 3. 导出报告
```bash
# 导出文本报告
./log_analyzer -file access.log -out report.txt

# 导出JSON报告
./log_analyzer -file access.log -output json -out report.json

# 包含详细日志条目
./log_analyzer -file app.log -details -output json -out detailed.json
```

### 支持的日志格式

#### Apache/Nginx访问日志
```
192.168.1.100 - - [25/Dec/2023:10:15:30 +0800] "GET /index.html HTTP/1.1" 200 1234 "https://example.com" "Mozilla/5.0"
```

#### 通用应用日志
```
2023-12-25 10:15:30 [ERROR] Database connection failed
2023-12-25 10:15:31 [INFO] Retrying connection...
```

#### JSON格式日志
```json
{"timestamp":"2023-12-25T10:15:30Z","level":"ERROR","message":"Database error","ip":"192.168.1.100"}
```

#### 系统日志
```
Dec 25 10:15:30 server01 nginx: connection timed out
```

### 命令行选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-file` | 日志文件路径 | 必需 |
| `-format` | 日志格式 (apache/nginx/common/syslog/json/auto) | auto |
| `-level` | 过滤日志级别 (ERROR/WARN/INFO/DEBUG) | 无 |
| `-pattern` | 过滤模式 (正则表达式) | 无 |
| `-output` | 输出格式 (text/json) | text |
| `-out` | 输出文件路径 | 标准输出 |
| `-top` | 显示前N项统计 | 10 |
| `-details` | 包含详细日志条目 | false |
| `-help` | 显示帮助信息 | false |

## 使用示例

### 示例1: 分析Web服务器访问日志
```bash
$ ./log_analyzer -file access.log -format nginx

🔍 开始分析日志文件: access.log
📋 使用格式: nginx
✅ 解析完成! 处理了 10000 行日志

📊 日志分析报告
========================================

📈 基础统计:
  总行数: 10000
  有效行数: 9876
  错误行数: 124
  时间范围: 2023-12-25 00:00:01 ~ 2023-12-25 23:59:58
  持续时间: 23h59m57s

🎯 日志级别统计:
  INFO: 8234 (83.4%)
  WARN: 1234 (12.5%)
  ERROR: 408 (4.1%)

🌐 Top IP地址:
  1. 192.168.1.100: 1234 (12.5%)
  2. 192.168.1.101: 876 (8.9%)
  3. 10.0.0.15: 654 (6.6%)
  
📊 HTTP状态码统计:
  2xx: 8234 (83.4%)
  3xx: 876 (8.9%)
  4xx: 654 (6.6%)
  5xx: 112 (1.1%)

🔧 HTTP方法统计:
  GET: 7890 (79.9%)
  POST: 1234 (12.5%)
  PUT: 456 (4.6%)
  DELETE: 296 (3.0%)

========================================
分析完成! 🎉
```

### 示例2: 错误日志分析
```bash
$ ./log_analyzer -file app.log -level ERROR -pattern "database"

📊 日志分析报告
========================================

📈 基础统计:
  总行数: 50000
  有效行数: 234
  错误行数: 12
  时间范围: 2023-12-25 08:30:15 ~ 2023-12-25 18:45:32

❌ Top错误信息:
  1. Database connection timeout (89次)
  2. SQL query execution failed (45次)
  3. Connection pool exhausted (23次)
  4. Database deadlock detected (12次)
  5. Transaction rollback failed (8次)

========================================
分析完成! 🎉
```

### 示例3: JSON格式导出
```bash
$ ./log_analyzer -file access.log -output json -details -out analysis.json

✅ JSON报告已保存到: analysis.json
```

生成的JSON报告结构：
```json
{
  "stats": {
    "total_lines": 10000,
    "valid_lines": 9876,
    "error_lines": 124,
    "level_counts": {
      "INFO": 8234,
      "ERROR": 408,
      "WARN": 1234
    },
    "ip_counts": {
      "192.168.1.100": 1234,
      "192.168.1.101": 876
    },
    "status_counts": {
      "2xx": 8234,
      "4xx": 654
    },
    "time_range": "23h59m57s"
  },
  "entries": [
    {
      "timestamp": "2023-12-25T10:15:30Z",
      "level": "INFO",
      "message": "GET /index.html - 200",
      "ip": "192.168.1.100",
      "method": "GET",
      "url": "/index.html",
      "status": 200,
      "size": 1234
    }
  ]
}
```

## 代码结构

### 主要结构体

- `LogEntry`: 存储单个日志条目的完整信息
- `LogStats`: 统计信息汇总
- `LogAnalyzer`: 日志分析器主要逻辑
- `AnalyzerConfig`: 分析器配置参数

### 核心功能

- `ParseLogFile()`: 解析日志文件
- `parseLine()`: 解析单行日志
- `updateStats()`: 更新统计信息
- `GenerateReport()`: 生成分析报告
- `ExportJSON()`: 导出JSON格式数据

## 扩展建议

1. **实时监控**: 支持tail -f模式实时分析日志
2. **数据库存储**: 将分析结果存储到数据库
3. **Web界面**: 提供Web界面展示分析结果
4. **告警功能**: 基于阈值的智能告警
5. **图表生成**: 生成统计图表和趋势图
6. **分布式分析**: 支持多文件并行分析
7. **机器学习**: 异常检测和模式识别
8. **插件系统**: 支持自定义解析器和分析器

## 性能特点

- 支持大文件处理（GB级别）
- 内存使用优化，流式处理
- 高效的正则表达式匹配
- 快速的统计算法
- 并发安全的数据结构

## 学习要点

这个项目展示了以下Go语言特性：

- 文件I/O和流处理
- 正则表达式使用
- 时间处理和格式化
- JSON序列化和反序列化
- 命令行参数解析
- 数据结构和算法
- 错误处理模式
- 接口设计
- 性能优化技巧
- 大数据处理技术