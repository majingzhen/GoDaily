package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/**
文件同步工具
	使用Go语言实现的文件同步工具
	支持双向同步和单向同步模式
	提供文件差异检测和冲突解决
	支持增量同步和文件过滤
	展示Go语言文件操作、并发处理和错误处理等特性
*/

// SyncConfig 同步配置
type SyncConfig struct {
	SourceDir      string        // 源目录
	TargetDir      string        // 目标目录
	SyncMode       string        // 同步模式: unidirectional(单向), bidirectional(双向)
	CheckInterval  time.Duration // 检查间隔
	MaxFileSize    int64         // 最大文件大小
	IncludePattern string        // 包含模式
	ExcludePattern string        // 排除模式
	DryRun         bool          // 干运行模式
	Verbose        bool          // 详细输出
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mod_time"`
	Checksum string    `json:"checksum,omitempty"`
	IsDir    bool      `json:"is_dir"`
}

// FileSyncTool 文件同步工具
type FileSyncTool struct {
	Config    SyncConfig
	FileCache map[string]FileInfo
	mutex     sync.RWMutex
}

// NewFileSyncTool 创建新的文件同步工具
func NewFileSyncTool(config SyncConfig) *FileSyncTool {
	return &FileSyncTool{
		Config:    config,
		FileCache: make(map[string]FileInfo),
	}
}

// CalculateChecksum 计算文件校验和
func (fst *FileSyncTool) CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ScanDirectory 扫描目录
func (fst *FileSyncTool) ScanDirectory(dirPath string) (map[string]FileInfo, error) {
	files := make(map[string]FileInfo)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if path == dirPath {
			return nil
		}

		// 应用文件过滤
		relPath, _ := filepath.Rel(dirPath, path)
		if !fst.shouldIncludeFile(relPath, info) {
			return nil
		}

		fileInfo := FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		// 计算非目录文件的校验和
		if !info.IsDir() && info.Size() <= fst.Config.MaxFileSize {
			checksum, err := fst.CalculateChecksum(path)
			if err == nil {
				fileInfo.Checksum = checksum
			}
		}

		files[relPath] = fileInfo
		return nil
	})

	return files, err
}

// shouldIncludeFile 检查文件是否应该包含
func (fst *FileSyncTool) shouldIncludeFile(relPath string, info os.FileInfo) bool {
	// 排除目录
	if info.IsDir() {
		return true
	}

	// 应用包含模式
	if fst.Config.IncludePattern != "" {
		matched, _ := filepath.Match(fst.Config.IncludePattern, filepath.Base(relPath))
		if !matched {
			return false
		}
	}

	// 应用排除模式
	if fst.Config.ExcludePattern != "" {
		matched, _ := filepath.Match(fst.Config.ExcludePattern, filepath.Base(relPath))
		if matched {
			return false
		}
	}

	// 文件大小限制
	if info.Size() > fst.Config.MaxFileSize {
		return false
	}

	return true
}

// CompareDirectories 比较两个目录
func (fst *FileSyncTool) CompareDirectories() (toCopy, toDelete, conflicts []string) {
	sourceFiles, err := fst.ScanDirectory(fst.Config.SourceDir)
	if err != nil {
		log.Printf("扫描源目录失败: %v", err)
		return
	}

	targetFiles, err := fst.ScanDirectory(fst.Config.TargetDir)
	if err != nil {
		log.Printf("扫描目标目录失败: %v", err)
		return
	}

	fst.mutex.Lock()
	defer fst.mutex.Unlock()

	// 找出需要复制的文件
	for relPath, sourceFile := range sourceFiles {
		targetFile, exists := targetFiles[relPath]
		if !exists {
			toCopy = append(toCopy, relPath)
			continue
		}

		// 文件存在，比较差异
		if !sourceFile.IsDir && !targetFile.IsDir {
			if sourceFile.Size != targetFile.Size ||
				sourceFile.ModTime.After(targetFile.ModTime.Add(time.Second)) ||
				sourceFile.Checksum != targetFile.Checksum {
				toCopy = append(toCopy, relPath)
			}
		}
	}

	// 找出需要删除的文件
	for relPath := range targetFiles {
		if _, exists := sourceFiles[relPath]; !exists {
			toDelete = append(toDelete, relPath)
		}
	}

	// 检测冲突
	if fst.Config.SyncMode == "bidirectional" {
		for relPath, sourceFile := range sourceFiles {
			if targetFile, exists := targetFiles[relPath]; exists {
				if !sourceFile.IsDir && !targetFile.IsDir &&
					sourceFile.Checksum != targetFile.Checksum &&
					sourceFile.ModTime.After(targetFile.ModTime) &&
					targetFile.ModTime.After(sourceFile.ModTime.Add(-time.Second)) {
					conflicts = append(conflicts, relPath)
				}
			}
		}
	}

	return toCopy, toDelete, conflicts
}

// CopyFile 复制文件
func (fst *FileSyncTool) CopyFile(relPath string) error {
	sourcePath := filepath.Join(fst.Config.SourceDir, relPath)
	targetPath := filepath.Join(fst.Config.TargetDir, relPath)

	// 确保目标目录存在
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	if fst.Config.DryRun {
		if fst.Config.Verbose {
			fmt.Printf("[DRY RUN] 复制: %s -> %s\n", sourcePath, targetPath)
		}
		return nil
	}

	// 复制文件
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}

	// 保持文件属性
	sourceInfo, err := os.Stat(sourcePath)
	if err == nil {
		os.Chmod(targetPath, sourceInfo.Mode())
		os.Chtimes(targetPath, time.Now(), sourceInfo.ModTime())
	}

	if fst.Config.Verbose {
		fmt.Printf("复制: %s -> %s\n", sourcePath, targetPath)
	}

	return nil
}

// DeleteFile 删除文件
func (fst *FileSyncTool) DeleteFile(relPath string) error {
	targetPath := filepath.Join(fst.Config.TargetDir, relPath)

	if fst.Config.DryRun {
		if fst.Config.Verbose {
			fmt.Printf("[DRY RUN] 删除: %s\n", targetPath)
		}
		return nil
	}

	if err := os.Remove(targetPath); err != nil {
		return err
	}

	if fst.Config.Verbose {
		fmt.Printf("删除: %s\n", targetPath)
	}

	return nil
}

// ResolveConflict 解决文件冲突
func (fst *FileSyncTool) ResolveConflict(relPath string, choice string) error {
	sourcePath := filepath.Join(fst.Config.SourceDir, relPath)
	targetPath := filepath.Join(fst.Config.TargetDir, relPath)

	switch choice {
	case "source":
		return fst.CopyFile(relPath)
	case "target":
		// 将目标文件复制回源目录
		if fst.Config.DryRun {
			if fst.Config.Verbose {
				fmt.Printf("[DRY RUN] 保留目标文件: %s\n", targetPath)
			}
			return nil
		}

		sourceFile, err := os.Open(targetPath)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		targetFile, err := os.Create(sourcePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			return err
		}

		if fst.Config.Verbose {
			fmt.Printf("保留目标文件: %s\n", targetPath)
		}

	case "both":
		// 重命名目标文件
		newName := relPath + ".conflict_" + time.Now().Format("20060102_150405")
		newPath := filepath.Join(fst.Config.TargetDir, newName)

		if fst.Config.DryRun {
			if fst.Config.Verbose {
				fmt.Printf("[DRY RUN] 重命名冲突文件: %s -> %s\n", targetPath, newPath)
			}
			return nil
		}

		if err := os.Rename(targetPath, newPath); err != nil {
			return err
		}

		// 复制源文件
		if err := fst.CopyFile(relPath); err != nil {
			return err
		}

		if fst.Config.Verbose {
			fmt.Printf("处理冲突: %s -> %s (保留两个版本)\n", targetPath, newPath)
		}
	}

	return nil
}

// RunSync 执行同步
func (fst *FileSyncTool) RunSync() error {
	fmt.Printf("🔍 开始同步: %s -> %s\n", fst.Config.SourceDir, fst.Config.TargetDir)
	fmt.Printf("📋 模式: %s\n", fst.Config.SyncMode)

	toCopy, toDelete, conflicts := fst.CompareDirectories()

	fmt.Printf("📊 检测到变化:\n")
	fmt.Printf("  需要复制: %d 个文件\n", len(toCopy))
	fmt.Printf("  需要删除: %d 个文件\n", len(toDelete))
	fmt.Printf("  冲突文件: %d 个\n", len(conflicts))

	// 处理删除
	for _, file := range toDelete {
		if err := fst.DeleteFile(file); err != nil {
			log.Printf("删除失败 %s: %v", file, err)
		}
	}

	// 处理复制
	for _, file := range toCopy {
		if err := fst.CopyFile(file); err != nil {
			log.Printf("复制失败 %s: %v", file, err)
		}
	}

	// 处理冲突
	if len(conflicts) > 0 {
		fmt.Printf("\n⚠️  检测到文件冲突:\n")
		reader := bufio.NewReader(os.Stdin)

		for _, file := range conflicts {
			fmt.Printf("冲突文件: %s\n", file)
			fmt.Printf("请选择处理方式 (source/target/both): ")

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if err := fst.ResolveConflict(file, choice); err != nil {
				log.Printf("处理冲突失败 %s: %v", file, err)
			}
		}
	}

	fmt.Printf("✅ 同步完成!\n")
	return nil
}

// ContinuousSync 持续同步模式
func (fst *FileSyncTool) ContinuousSync() {
	fmt.Printf("🔄 启动持续同步模式，检查间隔: %v\n", fst.Config.CheckInterval)

	ticker := time.NewTicker(fst.Config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n⏰ %s - 执行定期同步检查\n", time.Now().Format("2006-01-02 15:04:05"))
			fst.RunSync()
		}
	}
}

func main() {
	var (
		sourceDir      = flag.String("source", "", "源目录路径")
		targetDir      = flag.String("target", "", "目标目录路径")
		syncMode       = flag.String("mode", "unidirectional", "同步模式: unidirectional(单向)/bidirectional(双向)")
		checkInterval  = flag.Duration("interval", 30*time.Second, "检查间隔(持续模式)")
		maxFileSize    = flag.Int64("maxsize", 100*1024*1024, "最大文件大小(字节)")
		includePattern = flag.String("include", "", "包含文件模式")
		excludePattern = flag.String("exclude", "", "排除文件模式")
		dryRun         = flag.Bool("dryrun", false, "干运行模式(不实际执行)")
		continuous     = flag.Bool("continuous", false, "持续同步模式")
		verbose        = flag.Bool("verbose", false, "详细输出")
		showHelp       = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *showHelp || *sourceDir == "" || *targetDir == "" {
		fmt.Println("🔄 文件同步工具 - 使用帮助")
		fmt.Println("========================================")
		fmt.Println("用法: file_sync_tool [选项]")
		fmt.Println("\n选项:")
		flag.PrintDefaults()
		fmt.Println("\n同步模式说明:")
		fmt.Println("  unidirectional - 单向同步: 源目录 -> 目标目录")
		fmt.Println("  bidirectional  - 双向同步: 检测并解决冲突")
		fmt.Println("\n示例:")
		fmt.Println("  file_sync_tool -source ./src -target ./backup -mode unidirectional")
		fmt.Println("  file_sync_tool -source ./docs -target ./backup -include *.txt -dryrun")
		fmt.Println("  file_sync_tool -source ./data -target ./sync -continuous -interval 1m")
		return
	}

	// 验证目录存在
	if _, err := os.Stat(*sourceDir); os.IsNotExist(err) {
		fmt.Printf("❌ 源目录不存在: %s\n", *sourceDir)
		os.Exit(1)
	}

	if _, err := os.Stat(*targetDir); os.IsNotExist(err) {
		fmt.Printf("⚠️  目标目录不存在，将创建: %s\n", *targetDir)
		if err := os.MkdirAll(*targetDir, 0755); err != nil {
			fmt.Printf("❌ 创建目标目录失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建配置
	config := SyncConfig{
		SourceDir:      *sourceDir,
		TargetDir:      *targetDir,
		SyncMode:       *syncMode,
		CheckInterval:  *checkInterval,
		MaxFileSize:    *maxFileSize,
		IncludePattern: *includePattern,
		ExcludePattern: *excludePattern,
		DryRun:         *dryRun,
		Verbose:        *verbose,
	}

	// 创建同步工具
	syncTool := NewFileSyncTool(config)

	// 执行同步
	if err := syncTool.RunSync(); err != nil {
		fmt.Printf("❌ 同步失败: %v\n", err)
		os.Exit(1)
	}

	// 持续同步模式
	if *continuous {
		syncTool.ContinuousSync()
	}
}
