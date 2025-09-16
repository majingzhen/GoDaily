package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
文件后缀批量修改工具
	flag包处理命令行参数
	os和path/filepath包处理文件系统操作
	递归遍历目录结构
	错误处理和状态报告
	字符串处理函数
*/
// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("文件后缀批量修改工具")
	fmt.Println("用法: ext_changer [选项]")
	fmt.Println("选项:")
	fmt.Println("  -dir			要处理的目录路径(默认：当前目录)")
	fmt.Println("  -old			要修改的文件后缀(例如：.txt)")
	fmt.Println("  -new			修改后的文件后缀(例如：.md)")
	fmt.Println("  -recurse		是否递归处理子目录(true/false, 默认：false)")
	fmt.Println("  -dry			试运行，不实际修改文件(true/false, 默认：false)")
	fmt.Println("\n示例:")
	fmt.Println("  将当前目录下所有.txt文件改为.md")
	fmt.Println("  ext_changer -old .txt -new .md")
	fmt.Println("  将当前目录下所有.jpg文件改为.png，并递归处理子目录")
	fmt.Println("  ext_changer -old .jpg -new .png -recurse true")
}

// 检查后缀是否符合格式（不含点）
func normalizeExtension(ext string) string {
	// TrimPrefix 从字符串的开头移除指定的前缀
	return strings.TrimPrefix(ext, ".")
}

// 处理文件重命名
func processFile(filePath string, oldExt string, newExt string, dryRun bool) (bool, error) {
	// 获取文件名和当前后缀
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	// 检查文件是否有指定的旧后缀
	if !strings.HasSuffix(fileName, "."+oldExt) {
		return false, nil
	}

	// 生成新文件名
	newFileName := strings.TrimSuffix(fileName, "."+oldExt) + "." + newExt
	newFilePath := filepath.Join(dir, newFileName)

	// 检查新文件是否已存在
	if _, err := os.Stat(newFilePath); err == nil {
		return false, fmt.Errorf("文件已存在: %s", newFilePath)
	}

	// 显示操作信息
	fmt.Printf("将 ’%s‘ 改为 '%s'\n", fileName, newFileName)
	fmt.Println(dryRun)
	if !dryRun {
		// 重命名文件
		if err := os.Rename(filePath, newFilePath); err != nil {
			return false, err
		}
	}

	return true, nil
}

// 处理目录中的文件
func processDirectory(rootDir string, oldExt, newExt string, recurse, dryRun bool) (int, int, error) {
	var total, changed int
	// 遍历目录
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("访问路径失败: %s, 错误: %v", path, err)
		}
		// 如果是目录且不递归处理，则跳过
		if info.IsDir() && path != rootDir && !recurse {
			return filepath.SkipDir
		}
		// 处理文件
		if !info.IsDir() {
			total++
			ok, err := processFile(path, oldExt, newExt, dryRun)
			if err != nil {
				fmt.Printf("处理文件失败: %s, 错误: %v\n", path, err)
				return nil // 继续处理下一个文件
			}
			if ok {
				changed++
			}
		}
		return nil
	})
	return total, changed, err
}

// 主函数
func main() {
	// 解析命令行参数
	dir := flag.String("dir", ".", "要处理的目录路径")
	oldExt := flag.String("old", "", "要修改的文件后缀")
	newExt := flag.String("new", "", "新的文件后缀")
	recurse := flag.Bool("recurse", false, "是否递归处理子目录")
	dryRun := flag.Bool("dry", false, "试运行，不实际修改文件")
	help := flag.Bool("help", false, "显示帮助信息")
	flag.Parse()

	// 显示帮助信息
	if *help || *oldExt == "" || *newExt == "" {
		printHelp()
		if *oldExt == "" || *newExt == "" {
			os.Exit(1) // 参数缺失，异常退出
		}
		os.Exit(0) // 显示帮助后正常退出
	}
	// 标准化后缀格式（去掉开头的点）
	normalizedOldExt := normalizeExtension(*oldExt)
	normalizedNewExt := normalizeExtension(*newExt)

	// 显示操作信息
	fmt.Printf("正在处理目录: %s\n", *dir)
	fmt.Printf("将所有的 .%s 文件改为 .%s\n", normalizedOldExt, normalizedNewExt)
	if *recurse {
		fmt.Println("递归处理子目录")
	}
	if *dryRun {
		fmt.Println("试运行，不修改文件")
	}
	fmt.Println("--------------------------------")

	// 处理目录
	total, changed, err := processDirectory(*dir, normalizedOldExt, normalizedNewExt, *recurse, *dryRun)
	// 显示结果摘要
	fmt.Println("------------------------")
	fmt.Printf("处理完成。共检查 %d 个文件，", total)
	if *dryRun {
		fmt.Printf("将修改 %d 个文件\n", changed)
	} else {
		fmt.Printf("成功修改 %d 个文件\n", changed)
	}

	if err != nil {
		fmt.Printf("处理过程中出现错误: %v\n", err)
		os.Exit(1)
	}
}
