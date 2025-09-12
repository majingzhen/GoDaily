# 猜数字游戏

基于 Go 语言的控制台猜数字游戏

## Go 语法特性实现

### 1. 包管理和导入
```go
package main

import (
    "fmt"        // 格式化输入输出
    "math/rand"  // 随机数生成
    "os"         // 操作系统接口
    "strconv"    // 字符串转换
    "time"       // 时间处理
)
```

### 2. 常量声明
```go
const (
    maxTries = 10  // 常量组声明
    minNum   = 1   
    maxNum   = 100
)
```

### 3. init 函数
```go
func init() {
    rand.NewSource(time.Now().UnixNano())  // 包初始化时设置随机种子
}
```

### 4. 函数定义和返回值
```go
// 单返回值函数
func generateTarget() int {
    return rand.Intn(maxNum-minNum+1) + minNum
}

// 多返回值函数 - Go 特色语法
func parseInput(input string) (int, error) {
    number, err := strconv.Atoi(input)
    if err != nil {
        return 0, fmt.Errorf("请输入有效数字")  // 错误包装
    }
    if number < minNum || number > maxNum {
        return 0, fmt.Errorf("请输入%d到%d之间的数字", minNum, maxNum)
    }
    return number, nil
}
```

### 5. 错误处理模式
```go
// Go 标准错误处理模式
guess, err := parseInput(input)
if err != nil {
    fmt.Println(err)
    continue  // 跳过本次循环
}

// 输入读取错误处理
_, err := fmt.Scanln(&input)
if err != nil {
    fmt.Println("输入错误，请重新输入")
    continue
}
```

### 6. 变量声明和初始化
```go
// 短变量声明
tartget := generateTarget()
tries := 0
won := false

// 变量声明
var input string

// 多重赋值
remaining := maxTries - tries
```

### 7. 控制结构

#### for 循环
```go
// for 循环作为 while 使用
for tries < maxTries {
    // 游戏逻辑
}
```

#### 条件判断
```go
// if-else if-else 链式判断
if guess == tartget {
    fmt.Println("恭喜你，猜对了！")
    won = true
    break
} else if guess < tartget {
    fmt.Println("猜小了")
} else {
    fmt.Println("猜大了")
}

// 字符串比较
if input == "quit" {
    fmt.Println("游戏结束")
    os.Exit(0)  // 程序退出
}
```

### 8. 格式化输出
```go
// Printf 格式化输出
fmt.Printf("请输入你的猜测 (还剩%d次机会): ", remaining)
fmt.Printf("请输入%d到%d之间的数字", minNum, maxNum)
fmt.Printf("游戏结束，你猜了%d次，恭喜你猜对了！", tries)
```

## 核心技术点

### 随机数生成
```go
// 设置随机种子 - 确保每次运行结果不同
rand.NewSource(time.Now().UnixNano())

// 生成指定范围随机数
return rand.Intn(maxNum-minNum+1) + minNum
```

### 字符串转换
```go
// 字符串转整数，返回值和错误
number, err := strconv.Atoi(input)
```

### 错误创建
```go
// 使用 fmt.Errorf 创建格式化错误
return 0, fmt.Errorf("请输入%d到%d之间的数字", minNum, maxNum)
```

## 编译运行
```bash
go run guess_number.go
# 或
go build guess_number.go && ./guess_number
```