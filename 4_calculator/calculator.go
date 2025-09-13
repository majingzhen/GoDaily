package main

import (
	"fmt"
)

// 加法
func add(a, b float64) float64 {
	return a + b
}

// 减法
func sub(a, b float64) float64 {
	return a - b
}

// 乘法
func mul(a, b float64) float64 {
	return a * b
}

// 除法
func div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("除数不能为0")
	}
	return a / b, nil
}

// 计算器
func calculator() {
	fmt.Println("=== 简单计算器 ===")
	fmt.Println("支持的操作: +, -, *, /")
	fmt.Println("输入 'exit' 退出")

	var a, b float64
	var op string
	for {
		fmt.Print("请输入表达式 (例如: 3 + 4)")
		fmt.Scan(&a, &op, &b)

		if op == "exit" {
			fmt.Println("退出计算器")
			break
		}

		switch op {
		case "+":
			fmt.Printf("结果: %.2f\n", add(a, b))
		case "-":
			fmt.Printf("结果: %.2f\n", sub(a, b))
		case "*":
			fmt.Printf("结果: %.2f\n", mul(a, b))
		case "/":
			result, err := div(a, b)
			if err != nil {
				fmt.Println("错误:", err)
			} else {
				fmt.Printf("结果: %.2f\n", result)
			}
		default:
			fmt.Println("无效的操作符")
		}
	}
}

func main() {
	calculator()
}
