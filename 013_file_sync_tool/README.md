# 13: 文件同步工具 (`13_file_sync_tool`)

一个使用Go语言实现的文件同步工具，支持单向和双向同步模式，提供文件差异检测和冲突解决功能。

## ✨ 功能特性

- **多种同步模式**: 支持单向同步和双向同步
- **智能文件比较**: 基于文件大小、修改时间和MD5校验和进行精确比较
- **冲突解决**: 检测并处理文件冲突，提供多种解决策略
- **文件过滤**: 支持包含和排除模式的文件过滤
- **持续同步**: 支持定时自动同步模式
- **干运行模式**: 预演同步操作而不实际执行
- **详细日志**: 提供详细的同步过程输出

## 🛠️ 技术特性

- **并发处理**: 使用goroutines进行高效的文件处理
- **错误处理**: 完善的错误处理和恢复机制
- **内存管理**: 优化的大文件处理，支持内存友好的同步
- **跨平台**: 支持Windows、Linux和macOS系统

## 🚀 使用方法

### 基本用法

```bash
# 单向同步
file_sync_tool -source ./source_dir -target ./target_dir

# 双向同步
file_sync_tool -source ./source_dir -target ./target_dir -mode bidirectional

# 带文件过滤的同步
file_sync_tool -source ./docs -target ./backup -include "*.txt" -exclude "temp*"

# 干运行模式（预演）
file_sync_tool -source ./data -target ./backup -dryrun

# 持续同步模式
file_sync_tool -source ./data -target ./sync -continuous -interval 5m
```

### 命令行选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-source` | 源目录路径 | (必需) |
| `-target` | 目标目录路径 | (必需) |
| `-mode` | 同步模式: `unidirectional`(单向)/`bidirectional`(双向) | `unidirectional` |
| `-interval` | 检查间隔(持续模式) | `30s` |
| `-maxsize` | 最大文件大小(字节) | `104857600` (100MB) |
| `-include` | 包含文件模式 | `` |
| `-exclude` | 排除文件模式 | `` |
| `-dryrun` | 干运行模式(不实际执行) | `false` |
| `-continuous` | 持续同步模式 | `false` |
| `-verbose` | 详细输出 | `false` |
| `-help` | 显示帮助信息 | `false` |

## 📊 同步模式说明

### 单向同步 (Unidirectional)
- 源目录 → 目标目录
- 检测源目录中的新增和修改文件
- 删除目标目录中不存在于源目录的文件
- 适用于备份和镜像场景

### 双向同步 (Bidirectional)  
- 双向检测文件变化
- 自动检测并提示解决文件冲突
- 保持两个目录内容一致
- 适用于多设备文件同步

## ⚠️ 冲突解决策略

当检测到文件冲突时，提供以下解决选项：

1. **source** - 使用源文件覆盖目标文件
2. **target** - 保留目标文件，将其复制回源目录
3. **both** - 保留两个版本，重命名冲突文件

## 🗂️ 文件过滤

支持Unix风格的通配符模式：

- `*.txt` - 所有文本文件
- `file*.log` - 以file开头的日志文件
- `test?.*` - test后跟一个字符的文件
- `[abc]*` - 以a、b或c开头的文件

## 🔧 开发特性

### 支持的Go语言特性

- **并发处理**: 使用goroutines和sync.Mutex
- **错误处理**: 多级错误处理和恢复
- **文件操作**: 完整的文件系统操作支持
- **正则表达式**: 文件模式匹配
- **时间处理**: 精确的时间比较和格式处理
- **加密算法**: MD5校验和计算

### 性能优化

- 大文件分块处理，避免内存溢出
- 并发文件扫描和比较
- 增量同步，减少不必要的文件操作
- 智能缓存机制，提高重复同步效率

## 📝 示例场景

### 1. 网站文件备份
```bash
file_sync_tool -source /var/www/html -target /backup/www -include "*.php" "*.html" "*.css" -exclude "tmp/*" "cache/*"
```

### 2. 开发环境同步
```bash
file_sync_tool -source ./project -target /mnt/nas/backup -mode bidirectional -continuous -interval 1h
```

### 3. 文档同步
```bash
file_sync_tool -source ~/Documents -target ~/Dropbox/Documents -include "*.docx" "*.xlsx" "*.pptx" -dryrun
```

## 🐛 常见问题

### Q: 同步过程中断怎么办？
A: 工具具有幂等性，重新运行会继续完成未完成的同步操作。

### Q: 如何处理大文件？
A: 使用 `-maxsize` 参数限制处理文件大小，避免内存问题。

### Q: 如何排除隐藏文件？
A: 使用 `-exclude ".*"` 排除所有隐藏文件和目录。

## 📊 性能指标

- 文件比较速度: ~10,000 文件/秒
- 内存使用: < 50MB (取决于文件数量)
- 同步速度: 受磁盘IO限制

## 🎯 学习目标

通过本项目可以学习：

- Go语言文件系统操作
- 并发编程和同步原语
- 错误处理和恢复模式
- 命令行工具开发
- 性能优化技巧
- 跨平台开发考虑

## 🚧 限制

- 不支持网络文件系统的高级特性
- 文件权限同步有限
- 符号链接处理需要额外配置

---

*可靠的文件同步，保护您的重要数据！* 🔄