package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// 游戏配置
const (
	maxTries = 10  // 最大尝试次数
	minNum   = 1   // 最小数字
	maxNum   = 100 // 最大数字
)

// 初始化随机数种子
func init() {
	rand.NewSource(time.Now().UnixNano())
}

// 生成目标数字
func generateTarget() int {
	return rand.Intn(maxNum-minNum+1) + minNum
}

// 解析用户输入
func parseInput(input string) (int, error) {
	number, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("请输入有效数字")
	}
	if number < minNum || number > maxNum {
		return 0, fmt.Errorf("请输入%d到%d之间的数字", minNum, maxNum)
	}
	return number, nil
}

// 显示游戏帮助信息
func showHelp() {
	fmt.Println("猜数字游戏规则:")
	fmt.Printf("1.系统会生成一个%d到%d之间的随机整数\n", minNum, maxNum)
	fmt.Printf("2.你有%d次机会猜出这个数字\n", maxTries)
	fmt.Println("3.每次猜测后，系统会提示你猜大了还是猜小了")
	fmt.Println("4.输入 'quit' 可以退出游戏")
	fmt.Println("--------------------------------------------------------")
}

// 猜数字游戏
func main() {
	// 游戏初始化
	tartget := generateTarget()
	tries := 0
	won := false

	// 欢迎信息
	fmt.Println("欢迎来到猜数字游戏! ")
	showHelp()

	// 循环进行游戏
	for tries < maxTries {
		remaining := maxTries - tries
		fmt.Printf("\n请输入你的猜测 (还剩%d次机会): ", remaining)

		// 读取用户输入
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("输入错误，请重新输入")
			continue
		}

		if input == "quit" {
			fmt.Println("游戏结束")
			os.Exit(0)
		}

		// 解析并验证输入
		guess, err := parseInput(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 增加尝试次数
		tries++

		// 判断猜测结果
		if guess == tartget {
			fmt.Println("恭喜你，猜对了！")
			won = true
			break
		} else if guess < tartget {
			fmt.Println("猜小了")
		} else {
			fmt.Println("猜大了")
		}
	}

	if won {
		fmt.Printf("游戏结束，你猜了%d次，恭喜你猜对了！", tries)
	} else {
		fmt.Printf("游戏结束，你猜了%d次，没有猜对，游戏结束！", tries)
	}
}
