package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
URL短链接生成器
	使用Go基础语法实现URL短链接服务
	支持自定义短链接别名
	包含链接过期时间设置
	提供链接访问统计功能
	支持批量生成和管理
	展示Go语言结构体、映射、时间处理等特性
*/

// URLEntry URL条目结构体
type URLEntry struct {
	ID          string     // 短链接ID
	OriginalURL string     // 原始URL
	ShortCode   string     // 短链接代码
	CreatedAt   time.Time  // 创建时间
	ExpiresAt   *time.Time // 过期时间（可选）
	AccessCount int        // 访问次数
	LastAccess  *time.Time // 最后访问时间
	CustomAlias string     // 自定义别名
	Description string     // 描述信息
}

// URLShortener URL短链接服务结构体
type URLShortener struct {
	URLs       map[string]*URLEntry // 存储URL条目 (shortCode -> URLEntry)
	BaseURL    string               // 基础URL
	CodeLength int                  // 短链接代码长度
}

// ShortenerConfig 短链接服务配置
type ShortenerConfig struct {
	BaseURL    string // 基础URL
	CodeLength int    // 短链接代码长度
	DefaultTTL int    // 默认过期时间（小时）
}

// 字符集用于生成短链接代码
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 创建新的URL短链接服务
func NewURLShortener(config *ShortenerConfig) *URLShortener {
	return &URLShortener{
		URLs:       make(map[string]*URLEntry),
		BaseURL:    config.BaseURL,
		CodeLength: config.CodeLength,
	}
}

// 生成随机短链接代码
func (us *URLShortener) generateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, us.CodeLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 生成基于URL的哈希代码
func (us *URLShortener) generateHashCode(originalURL string) string {
	hash := md5.Sum([]byte(originalURL + strconv.FormatInt(time.Now().Unix(), 10)))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr[:us.CodeLength]
}

// 验证URL格式
func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

// 创建短链接
func (us *URLShortener) CreateShortURL(originalURL, customAlias, description string, ttlHours int) (*URLEntry, error) {
	// 验证URL格式
	if !isValidURL(originalURL) {
		return nil, fmt.Errorf("无效的URL格式: %s", originalURL)
	}

	var shortCode string

	// 如果提供了自定义别名，检查是否已存在
	if customAlias != "" {
		if _, exists := us.URLs[customAlias]; exists {
			return nil, fmt.Errorf("自定义别名 '%s' 已存在", customAlias)
		}
		shortCode = customAlias
	} else {
		// 生成短链接代码，确保唯一性
		for {
			shortCode = us.generateShortCode()
			if _, exists := us.URLs[shortCode]; !exists {
				break
			}
		}
	}

	// 创建URL条目
	entry := &URLEntry{
		ID:          fmt.Sprintf("url_%d", time.Now().Unix()),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		AccessCount: 0,
		CustomAlias: customAlias,
		Description: description,
	}

	// 设置过期时间
	if ttlHours > 0 {
		expiresAt := time.Now().Add(time.Duration(ttlHours) * time.Hour)
		entry.ExpiresAt = &expiresAt
	}

	// 存储URL条目
	us.URLs[shortCode] = entry

	return entry, nil
}

// 解析短链接
func (us *URLShortener) ResolveShortURL(shortCode string) (*URLEntry, error) {
	entry, exists := us.URLs[shortCode]
	if !exists {
		return nil, fmt.Errorf("短链接 '%s' 不存在", shortCode)
	}

	// 检查是否过期
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return nil, fmt.Errorf("短链接 '%s' 已过期", shortCode)
	}

	// 更新访问统计
	entry.AccessCount++
	now := time.Now()
	entry.LastAccess = &now

	return entry, nil
}

// 获取短链接统计信息
func (us *URLShortener) GetStats(shortCode string) (*URLEntry, error) {
	entry, exists := us.URLs[shortCode]
	if !exists {
		return nil, fmt.Errorf("短链接 '%s' 不存在", shortCode)
	}
	return entry, nil
}

// 列出所有短链接
func (us *URLShortener) ListURLs() []*URLEntry {
	var entries []*URLEntry
	for _, entry := range us.URLs {
		entries = append(entries, entry)
	}
	return entries
}

// 删除短链接
func (us *URLShortener) DeleteShortURL(shortCode string) error {
	if _, exists := us.URLs[shortCode]; !exists {
		return fmt.Errorf("短链接 '%s' 不存在", shortCode)
	}
	delete(us.URLs, shortCode)
	return nil
}

// 清理过期链接
func (us *URLShortener) CleanupExpired() int {
	var expired []string
	now := time.Now()

	for shortCode, entry := range us.URLs {
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
			expired = append(expired, shortCode)
		}
	}

	for _, shortCode := range expired {
		delete(us.URLs, shortCode)
	}

	return len(expired)
}

// 显示URL条目详细信息
func displayURLEntry(entry *URLEntry, baseURL string) {
	fmt.Printf("🔗 短链接信息:\n")
	fmt.Printf("  ID: %s\n", entry.ID)
	fmt.Printf("  短链接: %s/%s\n", baseURL, entry.ShortCode)
	fmt.Printf("  原始URL: %s\n", entry.OriginalURL)
	if entry.CustomAlias != "" {
		fmt.Printf("  自定义别名: %s\n", entry.CustomAlias)
	}
	if entry.Description != "" {
		fmt.Printf("  描述: %s\n", entry.Description)
	}
	fmt.Printf("  创建时间: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))

	if entry.ExpiresAt != nil {
		fmt.Printf("  过期时间: %s\n", entry.ExpiresAt.Format("2006-01-02 15:04:05"))
		remaining := time.Until(*entry.ExpiresAt)
		if remaining > 0 {
			fmt.Printf("  剩余时间: %.1f 小时\n", remaining.Hours())
		} else {
			fmt.Printf("  状态: 已过期\n")
		}
	} else {
		fmt.Printf("  过期时间: 永不过期\n")
	}

	fmt.Printf("  访问次数: %d\n", entry.AccessCount)
	if entry.LastAccess != nil {
		fmt.Printf("  最后访问: %s\n", entry.LastAccess.Format("2006-01-02 15:04:05"))
	}
}

// 显示所有URL列表
func displayURLList(entries []*URLEntry, baseURL string) {
	if len(entries) == 0 {
		fmt.Println("暂无短链接记录")
		return
	}

	fmt.Printf("📋 短链接列表 (共 %d 条):\n", len(entries))
	fmt.Println("========================================")

	for i, entry := range entries {
		status := "正常"
		if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
			status = "已过期"
		}

		fmt.Printf("%d. %s/%s\n", i+1, baseURL, entry.ShortCode)
		fmt.Printf("   -> %s\n", entry.OriginalURL)
		fmt.Printf("   访问: %d 次 | 状态: %s\n", entry.AccessCount, status)
		if entry.Description != "" {
			fmt.Printf("   描述: %s\n", entry.Description)
		}
		fmt.Println("   ----")
	}
}

// 交互式模式
func runInteractiveMode(shortener *URLShortener) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("🔗 URL短链接生成器 - 交互模式")
	fmt.Println("输入 'help' 查看可用命令")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "help", "h":
			showInteractiveHelp()

		case "create", "c":
			if len(parts) < 2 {
				fmt.Println("用法: create <URL> [别名] [描述] [过期小时]")
				continue
			}

			originalURL := parts[1]
			var alias, description string
			var ttlHours int

			if len(parts) > 2 {
				alias = parts[2]
			}
			if len(parts) > 3 {
				description = strings.Join(parts[3:len(parts)-1], " ")
			}
			if len(parts) > 4 {
				if hours, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
					ttlHours = hours
				}
			}

			entry, err := shortener.CreateShortURL(originalURL, alias, description, ttlHours)
			if err != nil {
				fmt.Printf("创建失败: %v\n", err)
			} else {
				fmt.Printf("✅ 短链接创建成功!\n")
				displayURLEntry(entry, shortener.BaseURL)
			}

		case "resolve", "r":
			if len(parts) < 2 {
				fmt.Println("用法: resolve <短链接代码>")
				continue
			}

			entry, err := shortener.ResolveShortURL(parts[1])
			if err != nil {
				fmt.Printf("解析失败: %v\n", err)
			} else {
				fmt.Printf("🎯 重定向到: %s\n", entry.OriginalURL)
				displayURLEntry(entry, shortener.BaseURL)
			}

		case "list", "l":
			entries := shortener.ListURLs()
			displayURLList(entries, shortener.BaseURL)

		case "stats", "s":
			if len(parts) < 2 {
				fmt.Println("用法: stats <短链接代码>")
				continue
			}

			entry, err := shortener.GetStats(parts[1])
			if err != nil {
				fmt.Printf("获取统计失败: %v\n", err)
			} else {
				displayURLEntry(entry, shortener.BaseURL)
			}

		case "delete", "d":
			if len(parts) < 2 {
				fmt.Println("用法: delete <短链接代码>")
				continue
			}

			err := shortener.DeleteShortURL(parts[1])
			if err != nil {
				fmt.Printf("删除失败: %v\n", err)
			} else {
				fmt.Printf("✅ 短链接 '%s' 已删除\n", parts[1])
			}

		case "cleanup":
			count := shortener.CleanupExpired()
			fmt.Printf("✅ 已清理 %d 个过期链接\n", count)

		case "exit", "quit", "q":
			fmt.Println("再见!")
			return

		default:
			fmt.Printf("未知命令: %s\n", command)
			fmt.Println("输入 'help' 查看可用命令")
		}
	}
}

// 显示交互模式帮助
func showInteractiveHelp() {
	fmt.Println("\n📖 可用命令:")
	fmt.Println("  create <URL> [别名] [描述] [过期小时] - 创建短链接")
	fmt.Println("  resolve <代码>                      - 解析短链接")
	fmt.Println("  list                               - 列出所有短链接")
	fmt.Println("  stats <代码>                       - 查看链接统计")
	fmt.Println("  delete <代码>                      - 删除短链接")
	fmt.Println("  cleanup                            - 清理过期链接")
	fmt.Println("  help                               - 显示帮助")
	fmt.Println("  exit                               - 退出程序")
	fmt.Println("\n💡 示例:")
	fmt.Println("  create https://www.google.com")
	fmt.Println("  create https://www.github.com github")
	fmt.Println("  create https://www.baidu.com baidu 搜索引擎 24")
}

// 显示帮助信息
func showHelp() {
	fmt.Println("URL短链接生成器")
	fmt.Println("用法: url_shortener [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -url         要缩短的URL")
	fmt.Println("  -alias       自定义短链接别名")
	fmt.Println("  -desc        链接描述")
	fmt.Println("  -ttl         过期时间(小时) (默认: 0, 永不过期)")
	fmt.Println("  -base        基础URL (默认: http://short.ly)")
	fmt.Println("  -length      短链接代码长度 (默认: 6)")
	fmt.Println("  -interactive 交互模式")
	fmt.Println("  -help        显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  创建短链接:")
	fmt.Println("  url_shortener -url https://www.example.com")
	fmt.Println()
	fmt.Println("  创建带别名的短链接:")
	fmt.Println("  url_shortener -url https://www.github.com -alias github")
	fmt.Println()
	fmt.Println("  创建24小时后过期的短链接:")
	fmt.Println("  url_shortener -url https://www.google.com -ttl 24 -desc 搜索引擎")
	fmt.Println()
	fmt.Println("  启动交互模式:")
	fmt.Println("  url_shortener -interactive")
}

func main() {
	// 解析命令行参数
	urlToShorten := flag.String("url", "", "要缩短的URL")
	customAlias := flag.String("alias", "", "自定义短链接别名")
	description := flag.String("desc", "", "链接描述")
	ttlHours := flag.Int("ttl", 0, "过期时间(小时)")
	baseURL := flag.String("base", "http://short.ly", "基础URL")
	codeLength := flag.Int("length", 6, "短链接代码长度")
	interactive := flag.Bool("interactive", false, "交互模式")
	help := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助
	if *help {
		showHelp()
		os.Exit(0)
	}

	// 创建短链接服务配置
	config := &ShortenerConfig{
		BaseURL:    *baseURL,
		CodeLength: *codeLength,
	}

	// 创建短链接服务
	shortener := NewURLShortener(config)

	// 交互模式
	if *interactive {
		runInteractiveMode(shortener)
		return
	}

	// 如果没有提供URL，显示帮助
	if *urlToShorten == "" {
		fmt.Println("错误: 请提供要缩短的URL")
		fmt.Println("使用 -help 查看使用说明，或使用 -interactive 启动交互模式")
		os.Exit(1)
	}

	// 创建短链接
	fmt.Printf("正在创建短链接...\n\n")
	entry, err := shortener.CreateShortURL(*urlToShorten, *customAlias, *description, *ttlHours)
	if err != nil {
		fmt.Printf("创建失败: %v\n", err)
		os.Exit(1)
	}

	// 显示结果
	fmt.Printf("✅ 短链接创建成功!\n")
	fmt.Println("========================================")
	displayURLEntry(entry, shortener.BaseURL)
	fmt.Println("========================================")
	fmt.Println("创建完成!")
}
