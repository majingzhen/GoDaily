# 文件完整性检查工具 (015)

一个用于监控文件变化和计算校验和的工具，支持多种哈希算法和持续监控模式。

## 功能特性

- ✅ 支持 MD5、SHA1、SHA256 校验和算法
- ✅ 递归扫描子目录
- ✅ 实时监控模式，检测文件变化
- ✅ 支持控制台和 JSON 输出格式
- ✅ 结果保存到文件
- ✅ 检测文件新增、修改、删除、大小变化

## 使用方法

### 基本使用

```bash
# 检查当前目录
file_integrity_checker

# 检查指定目录，使用 MD5 算法
file_integrity_checker -path /tmp -algo md5

# 递归检查子目录
file_integrity_checker -recursive

# 输出 JSON 格式
file_integrity_checker -output json

# 保存结果到文件
file_integrity_checker -file result.json
```

### 监控模式

```bash
# 每30秒监控一次目录变化
file_integrity_checker -monitor 30s

# 使用 SHA1 算法，每1分钟监控一次
file_integrity_checker -algo sha1 -monitor 1m
```

## 命令行选项

| 选项 | 默认值 | 描述 |
|------|--------|------|
| `-path` | `.` | 要检查的目录路径 |
| `-algo` | `sha256` | 校验和算法：md5, sha1, sha256 |
| `-recursive` | `false` | 递归检查子目录 |
| `-output` | `console` | 输出格式：console, json |
| `-file` | | 保存结果到文件 |
| `-monitor` | | 监控模式间隔（如：5s, 1m） |
| `-help` | `false` | 显示帮助信息 |

## 输出示例

### 控制台输出

```
========== 文件完整性检查报告 ==========
时间: 2024-01-15 10:30:45
目录: 3 个
文件: 25 个
总大小: 15.2 MB

文件列表:
main.go                                   2.1 KB  e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
config.yaml                               1.5 KB  a1b2c3d4e5f67890123456789012345678901234
README.md                                 5.2 KB  f1e2d3c4b5a67890123456789012345678901234
================================
```

### JSON 输出

```json
{
  "timestamp": "2024-01-15T10:30:45.123456789+08:00",
  "total_files": 25,
  "total_dirs": 3,
  "total_size": 15200000,
  "files": [
    {
      "path": "/path/to/main.go",
      "size": 2156,
      "modified": "2024-01-15T09:15:23.456789+08:00",
      "mode": "-rw-r--r--",
      "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "is_dir": false,
      "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    }
  ]
}
```

## 监控模式

在监控模式下，工具会持续扫描指定目录并报告文件变化：

```
开始监控模式，间隔: 30s
按 Ctrl+C 停止监控

[10:35:30] 检测到变化:
  新增: /path/to/newfile.txt
  修改: /path/to/updated.yaml
  大小变化: /path/to/growing.log

[10:36:00] 检测到变化:
  删除: /path/to/deleted.txt
  校验和变化: /path/to/modified.exe
```

## 构建和运行

```bash
# 构建
cd 015_file_integrity_checker
go build -o file_integrity_checker file_integrity_checker.go

# 运行
./file_integrity_checker -help
```

## 技术实现

- 使用 Go 标准库的 `crypto/md5`, `crypto/sha1`, `crypto/sha256` 计算校验和
- 通过 `filepath.Walk` 遍历文件系统
- 支持跨平台运行（Windows/Linux）
- 内存高效，可处理大目录结构

## 注意事项

- 对于大文件，计算校验和可能需要较长时间
- 监控模式会持续占用系统资源
- 建议在生产环境中设置合理的监控间隔
- 某些系统目录可能需要管理员权限才能访问