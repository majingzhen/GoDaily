package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

/**
 * 单词计数程序
 */

// 统计字符串中每个单词出现的次数
func countWords(s string) map[string]int {
	words := make(map[string]int)

	// 将字符串转为小写
	lowercase := strings.ToLower(s)

	// 构建单词
	var currentWord strings.Builder

	for _, r := range lowercase {
		// 检查是否为字母或数字
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			currentWord.WriteRune(r)
		} else if currentWord.Len() > 0 {
			// 如果遇到非字母数字且当前有单词，则添加到映射
			word := currentWord.String()
			words[word]++
			currentWord.Reset()
		}
	}
	// 处理最后一个单词
	if currentWord.Len() > 0 {
		word := currentWord.String()
		words[word]++
	}
	return words
}

// 单词计数程序
func wordCount() {
	fmt.Println("=== 单词计数程序 ===")
	fmt.Println("请输入一段文本（输入空行结束）:")

	scanner := bufio.NewScanner(os.Stdin)
	var input strings.Builder

	// 读取多行输入，直到空行
	for {
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			break
		}
		input.WriteString(line + " ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取错误：", err)
		return
	}

	text := input.String()
	if text == "" {
		fmt.Println("输入为空")
		return
	}

	wordCounts := countWords(text)

	for word, count := range wordCounts {
		fmt.Printf("%s: %d\n", word, count)
	}

}

func main() {
	wordCount()
}
