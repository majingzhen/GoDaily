# URL短链接生成器

一个使用Go语言实现的URL短链接服务工具，支持创建、管理和统计短链接。

## 功能特性

- ✅ **URL缩短**: 将长URL转换为短链接
- ✅ **自定义别名**: 支持自定义短链接别名
- ✅ **过期时间**: 可设置链接过期时间
- ✅ **访问统计**: 记录链接访问次数和时间
- ✅ **批量管理**: 列出、删除、清理过期链接
- ✅ **交互模式**: 提供友好的交互式命令行界面
- ✅ **链接验证**: 验证URL格式有效性

## 技术特点

- 使用Go标准库实现
- 支持MD5哈希和随机字符串生成算法
- 内存存储(可扩展为数据库存储)
- 命令行参数解析
- 时间处理和过期管理
- 错误处理和输入验证

## 使用方法

### 编译程序
```bash
go build url_shortener.go
```

### 基本用法

#### 1. 创建短链接
```bash
# 创建基本短链接
./url_shortener -url https://www.example.com

# 创建带别名的短链接
./url_shortener -url https://www.github.com -alias github

# 创建带描述和过期时间的短链接
./url_shortener -url https://www.google.com -alias google -desc "谷歌搜索" -ttl 24
```

#### 2. 交互模式
```bash
# 启动交互模式
./url_shortener -interactive
```

### 交互模式命令

进入交互模式后，可以使用以下命令：

- `create <URL> [别名] [描述] [过期小时]` - 创建短链接
- `resolve <代码>` - 解析短链接
- `list` - 列出所有短链接
- `stats <代码>` - 查看链接统计
- `delete <代码>` - 删除短链接
- `cleanup` - 清理过期链接
- `help` - 显示帮助
- `exit` - 退出程序

### 命令行选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-url` | 要缩短的URL | 无 |
| `-alias` | 自定义短链接别名 | 自动生成 |
| `-desc` | 链接描述 | 无 |
| `-ttl` | 过期时间(小时) | 0(永不过期) |
| `-base` | 基础URL | http://short.ly |
| `-length` | 短链接代码长度 | 6 |
| `-interactive` | 交互模式 | false |
| `-help` | 显示帮助信息 | false |

## 使用示例

### 示例1: 创建简单短链接
```bash
$ ./url_shortener -url https://www.baidu.com

正在创建短链接...

✅ 短链接创建成功!
========================================
🔗 短链接信息:
  ID: url_1640995200
  短链接: http://short.ly/aB3xY9
  原始URL: https://www.baidu.com
  创建时间: 2024-01-01 12:00:00
  过期时间: 永不过期
  访问次数: 0
========================================
创建完成!
```

### 示例2: 交互模式使用
```bash
$ ./url_shortener -interactive

🔗 URL短链接生成器 - 交互模式
输入 'help' 查看可用命令

> create https://www.github.com github 代码托管平台 48
✅ 短链接创建成功!
🔗 短链接信息:
  ID: url_1640995260
  短链接: http://short.ly/github
  原始URL: https://www.github.com
  自定义别名: github
  描述: 代码托管平台
  创建时间: 2024-01-01 12:01:00
  过期时间: 2024-01-03 12:01:00
  剩余时间: 48.0 小时
  访问次数: 0

> list
📋 短链接列表 (共 2 条):
========================================
1. http://short.ly/aB3xY9
   -> https://www.baidu.com
   访问: 0 次 | 状态: 正常
   ----
2. http://short.ly/github
   -> https://www.github.com
   访问: 0 次 | 状态: 正常
   描述: 代码托管平台
   ----

> resolve github
🎯 重定向到: https://www.github.com
🔗 短链接信息:
  ID: url_1640995260
  短链接: http://short.ly/github
  原始URL: https://www.github.com
  自定义别名: github
  描述: 代码托管平台
  创建时间: 2024-01-01 12:01:00
  过期时间: 2024-01-03 12:01:00
  剩余时间: 47.9 小时
  访问次数: 1
  最后访问: 2024-01-01 12:02:30
```

## 代码结构

### 主要结构体

- `URLEntry`: 存储单个URL的完整信息
- `URLShortener`: 短链接服务主要逻辑
- `ShortenerConfig`: 服务配置信息

### 核心功能

- `CreateShortURL()`: 创建短链接
- `ResolveShortURL()`: 解析短链接并更新统计
- `GetStats()`: 获取链接统计信息
- `CleanupExpired()`: 清理过期链接

## 扩展建议

1. **持久化存储**: 集成数据库(如SQLite、MySQL)存储链接数据
2. **Web界面**: 添加HTTP服务器提供Web管理界面
3. **API接口**: 提供RESTful API供其他应用调用
4. **访问日志**: 记录详细的访问日志和分析
5. **批量导入**: 支持从文件批量导入URL
6. **QR码生成**: 为短链接生成二维码
7. **链接备份**: 导出/导入链接数据功能

## 学习要点

这个项目展示了以下Go语言特性：

- 结构体和方法
- 映射(map)的使用
- 时间处理和格式化
- 命令行参数解析
- 错误处理
- 字符串处理
- 文件I/O操作
- 随机数生成
- 加密哈希(MD5)
- 交互式命令行程序设计