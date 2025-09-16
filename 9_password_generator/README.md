# 密码生成器 (Password Generator)

一个功能强大的命令行密码生成工具，使用Go语言编写，支持多种自定义选项。

## 功能特性

- 🔐 生成安全的随机密码
- 📏 自定义密码长度
- 🔤 支持多种字符类型（大小写字母、数字、特殊字符）
- 🚫 可选择排除易混淆字符
- 📊 密码强度评估
- 🔢 批量生成多个密码
- 💡 智能确保包含所选字符类型

## 使用方法

### 基本用法

```bash
# 生成一个默认12位密码
go run password_generator.go

# 生成指定长度的密码
go run password_generator.go -length 16

# 生成多个密码
go run password_generator.go -count 5
```

### 高级选项

```bash
# 生成包含特殊字符的强密码
go run password_generator.go -length 16 -symbols true

# 生成不包含易混淆字符的密码
go run password_generator.go -exclude true

# 只使用字母和数字
go run password_generator.go -symbols false

# 自定义所有选项
go run password_generator.go -length 20 -count 3 -upper true -lower true -numbers true -symbols true -exclude true
```

## 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-length` | int | 12 | 密码长度 |
| `-count` | int | 1 | 生成密码数量 |
| `-lower` | bool | true | 包含小写字母 |
| `-upper` | bool | true | 包含大写字母 |
| `-numbers` | bool | true | 包含数字 |
| `-symbols` | bool | false | 包含特殊字符 |
| `-exclude` | bool | false | 排除易混淆字符 |
| `-help` | bool | false | 显示帮助信息 |

## 字符集说明

- **小写字母**: a-z
- **大写字母**: A-Z  
- **数字**: 0-9
- **特殊字符**: !@#$%^&*()_+-=[]{}|;:,.<>?
- **易混淆字符**: 0, O, o, 1, l, I, i

## 密码强度评估

程序会自动评估生成密码的强度：

- **很弱**: 分数 < 2
- **弱**: 分数 2-3
- **中等**: 分数 4-5  
- **强**: 分数 ≥ 6

评估标准：
- 长度 ≥ 8位 (+1分)
- 长度 ≥ 12位 (+1分)
- 包含小写字母 (+1分)
- 包含大写字母 (+1分)
- 包含数字 (+1分)
- 包含特殊字符 (+1分)

## 示例输出

```
密码配置:
  长度: 16
  数量: 3
  小写字母: true
  大写字母: true
  数字: true
  特殊字符: true
  排除易混淆字符: false
------------------------
密码 1: K9m#Xp2@vBnQ7wRt (强度: 强)
密码 2: Fy8$Zc3!hLmN6qWe (强度: 强)
密码 3: Dj5%Gk9&uPsT2rVx (强度: 强)
------------------------
密码生成完成
```

## 技术实现

本工具使用了以下Go语言基础语法和概念：

- **包管理**: `package main` 和标准库导入
- **结构体**: `PasswordConfig` 存储配置信息
- **常量**: 定义字符集常量
- **函数和方法**: 模块化代码组织
- **命令行参数**: 使用 `flag` 包处理参数
- **字符串操作**: `strings` 包进行字符串处理
- **随机数生成**: `math/rand` 生成随机密码
- **切片和数组**: 管理字符集和密码字符
- **循环和条件**: 控制密码生成逻辑
- **错误处理**: Go语言标准错误处理模式

## 编译和运行

```bash
# 直接运行
go run password_generator.go -help

# 编译后运行
go build -o password_generator password_generator.go
./password_generator -length 16 -symbols true
```

## 安全提示

- 生成的密码请妥善保存
- 建议使用密码管理器存储密码
- 定期更换重要账户密码
- 不要在不安全的环境中生成密码