package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
文本搜索工具
	支持在文件和目录中搜索指定文本
	可选择是否区分大小写
	支持递归搜索子目录
	显示匹配内容所在的文件名和行号
*/

// 搜索文件内容
func searchInFile(filePath, query string, caseSensitive bool) ([]int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var lines []int
	text := string(content)
	linesText := strings.Split(text, "\n")

	for i, line := range linesText {
		lineToCheck := line
		queryToCheck := query

		// 如果不区分大小写，统一转换为小写
		if !caseSensitive {
			lineToCheck = strings.ToLower(line)
			queryToCheck = strings.ToLower(query)
		}

		if strings.Contains(lineToCheck, queryToCheck) {
			lines = append(lines, i+1) // 行号从1开始
		}
	}
	return lines, nil
}

// 处理目录搜索
func searchInDirectory(rootDir, query string, caseSensitive, recurse bool) {
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("访问路径失败: %s, 错误: %v", path, err)
			return nil
		}

		// 如果是目录且不递归处理，则跳过
		if info.IsDir() && path != rootDir && !recurse {
			return filepath.SkipDir
		}

		// 只处理文件
		if !info.IsDir() {
			lines, err := searchInFile(path, query, caseSensitive)
			if err != nil {
				fmt.Printf("读取 %s 失败：%v\n", path, err)
				return nil
			}

			if len(lines) > 0 {
				fmt.Printf("在 %s 中找到匹配项：\n", path)
				for _, line := range lines {
					fmt.Printf("第 %d 行\n", line)
				}
			}
		}
		return nil
	})
}

func main() {
	// 解析命令行参数
	dir := flag.String("dir", ".", "搜索目录")
	query := flag.String("query", "", "要搜索的文本")
	caseSensitive := flag.Bool("case", false, "是否区分大小写")
	recurse := flag.Bool("recurse", false, "是否递归搜索子目录")
	help := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助
	if *help || *query == "" {
		fmt.Println("文本搜索工具")
		fmt.Println("用法: text_search -query 搜索文本 [选项]")
		fmt.Println("选项:")
		fmt.Println("  -dir      搜索目录 (默认: 当前目录)")
		fmt.Println("  -case     是否区分大小写 (true/false, 默认: false)")
		fmt.Println("  -recurse  是否递归搜索子目录 (true/false, 默认: false)")
		fmt.Println("  -help     显示帮助信息")
		os.Exit(0)
	}

	fmt.Printf("搜索文本: %q\n", *query)
	fmt.Printf("搜索目录: %s\n", *dir)
	fmt.Printf("区分大小写: %v\n", *caseSensitive)
	fmt.Printf("递归搜索: %v\n", *recurse)
	fmt.Println("------------------------")

	searchInDirectory(*dir, *query, *caseSensitive, *recurse)
	fmt.Println("------------------------")
	fmt.Println("搜索完成")
}
