package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	Mode     string    `json:"mode"`
	MD5      string    `json:"md5,omitempty"`
	SHA1     string    `json:"sha1,omitempty"`
	SHA256   string    `json:"sha256,omitempty"`
	IsDir    bool      `json:"is_dir"`
	Checksum string    `json:"checksum,omitempty"`
}

type FileStats struct {
	Timestamp  time.Time  `json:"timestamp"`
	TotalFiles int        `json:"total_files"`
	TotalDirs  int        `json:"total_dirs"`
	TotalSize  int64      `json:"total_size"`
	Files      []FileInfo `json:"files,omitempty"`
}

func main() {
	var (
		path      = flag.String("path", ".", "要检查的目录路径")
		algo      = flag.String("algo", "sha256", "校验和算法: md5, sha1, sha256")
		recursive = flag.Bool("recursive", false, "递归检查子目录")
		output    = flag.String("output", "console", "输出格式: console, json")
		file      = flag.String("file", "", "保存结果到文件")
		monitor   = flag.Duration("monitor", 0, "监控模式间隔 (如: 5s, 1m)")
		help      = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	stats, err := scanDirectory(*path, *algo, *recursive)
	if err != nil {
		fmt.Printf("扫描目录失败: %v\n", err)
		return
	}

	outputResults(stats, *output, *file)

	if *monitor > 0 {
		startMonitoring(*path, *algo, *recursive, *output, *file, *monitor)
	}
}

func scanDirectory(path, algo string, recursive bool) (*FileStats, error) {
	stats := &FileStats{
		Timestamp: time.Now(),
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !recursive && filePath != path {
			relPath, _ := filepath.Rel(path, filePath)
			if strings.Contains(relPath, string(filepath.Separator)) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		fileInfo := FileInfo{
			Path:     filePath,
			Size:     info.Size(),
			Modified: info.ModTime(),
			Mode:     info.Mode().String(),
			IsDir:    info.IsDir(),
		}

		if !info.IsDir() {
			checksum, err := calculateChecksum(filePath, algo)
			if err == nil {
				fileInfo.Checksum = checksum
				switch algo {
				case "md5":
					fileInfo.MD5 = checksum
				case "sha1":
					fileInfo.SHA1 = checksum
				case "sha256":
					fileInfo.SHA256 = checksum
				}
			}
		}

		if info.IsDir() {
			stats.TotalDirs++
		} else {
			stats.TotalFiles++
			stats.TotalSize += info.Size()
			stats.Files = append(stats.Files, fileInfo)
		}

		return nil
	})

	return stats, err
}

func calculateChecksum(filePath, algo string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var hash []byte
	switch algo {
	case "md5":
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	case "sha1":
		h := sha1.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	case "sha256":
		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	default:
		return "", fmt.Errorf("不支持的算法: %s", algo)
	}

	return hex.EncodeToString(hash), nil
}

func outputResults(stats *FileStats, output, filePath string) {
	var outputFile *os.File
	var err error

	if filePath != "" {
		outputFile, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("创建输出文件失败: %v\n", err)
			return
		}
		defer outputFile.Close()
	}

	switch output {
	case "json":
		outputJSON(stats, outputFile)
	default:
		outputConsole(stats, outputFile)
	}
}

func outputConsole(stats *FileStats, file *os.File) {
	output := fmt.Sprintf(`
========== 文件完整性检查报告 ==========
时间: %s
目录: %d 个
文件: %d 个
总大小: %s

文件列表:
`,
		stats.Timestamp.Format("2006-01-02 15:04:05"),
		stats.TotalDirs,
		stats.TotalFiles,
		formatBytes(stats.TotalSize),
	)

	for _, fileInfo := range stats.Files {
		output += fmt.Sprintf("%-50s %10s %s\n",
			filepath.Base(fileInfo.Path),
			formatBytes(fileInfo.Size),
			fileInfo.Checksum,
		)
	}

	output += "================================\n"

	if file != nil {
		file.WriteString(output)
	} else {
		fmt.Print(output)
	}
}

func outputJSON(stats *FileStats, file *os.File) {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	if file != nil {
		file.Write(data)
		file.WriteString("\n")
	} else {
		fmt.Println(string(data))
	}
}

func startMonitoring(path, algo string, recursive bool, output, filePath string, interval time.Duration) {
	fmt.Printf("开始监控模式，间隔: %v\n", interval)
	fmt.Printf("按 Ctrl+C 停止监控\n\n")

	previousStats := make(map[string]FileInfo)

	for {
		stats, err := scanDirectory(path, algo, recursive)
		if err != nil {
			fmt.Printf("扫描失败: %v\n", err)
			time.Sleep(interval)
			continue
		}

		currentFiles := make(map[string]FileInfo)
		for _, file := range stats.Files {
			currentFiles[file.Path] = file
		}

		var changes []string

		for path, currentFile := range currentFiles {
			if prevFile, exists := previousStats[path]; exists {
				if currentFile.Modified != prevFile.Modified {
					changes = append(changes, fmt.Sprintf("修改: %s", path))
				}
				if currentFile.Size != prevFile.Size {
					changes = append(changes, fmt.Sprintf("大小变化: %s", path))
				}
				if currentFile.Checksum != prevFile.Checksum {
					changes = append(changes, fmt.Sprintf("校验和变化: %s", path))
				}
			} else {
				changes = append(changes, fmt.Sprintf("新增: %s", path))
			}
		}

		for path := range previousStats {
			if _, exists := currentFiles[path]; !exists {
				changes = append(changes, fmt.Sprintf("删除: %s", path))
			}
		}

		if len(changes) > 0 {
			fmt.Printf("[%s] 检测到变化:\n", time.Now().Format("15:04:05"))
			for _, change := range changes {
				fmt.Printf("  %s\n", change)
			}
			fmt.Println()
		}

		previousStats = currentFiles
		time.Sleep(interval)
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func showHelp() {
	fmt.Println("文件完整性检查工具 - 监控文件变化和计算校验和")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  file_integrity_checker [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -path string        要检查的目录路径 (默认: .)")
	fmt.Println("  -algo string        校验和算法: md5, sha1, sha256 (默认: sha256)")
	fmt.Println("  -recursive          递归检查子目录")
	fmt.Println("  -output string      输出格式: console, json (默认: console)")
	fmt.Println("  -file string        保存结果到文件")
	fmt.Println("  -monitor duration   监控模式间隔 (如: 5s, 1m)")
	fmt.Println("  -help               显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  file_integrity_checker                            # 检查当前目录")
	fmt.Println("  file_integrity_checker -path /tmp -algo md5       # 检查/tmp目录使用MD5")
	fmt.Println("  file_integrity_checker -recursive -output json    # 递归检查并输出JSON")
	fmt.Println("  file_integrity_checker -monitor 30s              # 每30秒监控一次")
	fmt.Println("  file_integrity_checker -file result.json          # 保存结果到文件")
}
