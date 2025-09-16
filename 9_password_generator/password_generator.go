package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

/**
密码生成器工具
	使用Go基础语法实现安全密码生成
	支持自定义密码长度和字符集
	包含大小写字母、数字、特殊字符
	可选择包含或排除特定字符类型
	批量生成多个密码
*/

// 字符集常量
const (
	LowerCase = "abcdefghijklmnopqrstuvwxyz"
	UpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers   = "0123456789"
	Symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// PasswordConfig 密码配置结构体
type PasswordConfig struct {
	Length           int  // 密码长度
	IncludeLower     bool // 包含小写字母
	IncludeUpper     bool // 包含大写字母
	IncludeNumbers   bool // 包含数字
	IncludeSymbols   bool // 包含特殊字符
	ExcludeAmbiguous bool // 排除易混淆字符
	Count            int  // 生成密码数量
}

// 易混淆字符
var ambiguousChars = []string{"0", "O", "o", "1", "l", "I", "i"}

// 生成字符集
func (pc *PasswordConfig) buildCharset() string {
	var charset strings.Builder

	if pc.IncludeLower {
		charset.WriteString(LowerCase)
	}
	if pc.IncludeUpper {
		charset.WriteString(UpperCase)
	}
	if pc.IncludeNumbers {
		charset.WriteString(Numbers)
	}
	if pc.IncludeSymbols {
		charset.WriteString(Symbols)
	}

	result := charset.String()

	// 排除易混淆字符
	if pc.ExcludeAmbiguous {
		for _, char := range ambiguousChars {
			result = strings.ReplaceAll(result, char, "")
		}
	}

	return result
}

// 生成单个密码
func generatePassword(config *PasswordConfig) (string, error) {
	charset := config.buildCharset()

	if len(charset) == 0 {
		return "", fmt.Errorf("字符集为空，请至少选择一种字符类型")
	}

	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	var password strings.Builder

	// 确保密码包含每种选中的字符类型至少一个
	if config.IncludeLower && len(LowerCase) > 0 {
		password.WriteByte(LowerCase[rand.Intn(len(LowerCase))])
	}
	if config.IncludeUpper && len(UpperCase) > 0 {
		password.WriteByte(UpperCase[rand.Intn(len(UpperCase))])
	}
	if config.IncludeNumbers && len(Numbers) > 0 {
		password.WriteByte(Numbers[rand.Intn(len(Numbers))])
	}
	if config.IncludeSymbols && len(Symbols) > 0 {
		password.WriteByte(Symbols[rand.Intn(len(Symbols))])
	}

	// 生成剩余字符
	for password.Len() < config.Length {
		randomIndex := rand.Intn(len(charset))
		password.WriteByte(charset[randomIndex])
	}

	// 打乱密码字符顺序
	passwordBytes := []byte(password.String())
	for i := len(passwordBytes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		passwordBytes[i], passwordBytes[j] = passwordBytes[j], passwordBytes[i]
	}

	return string(passwordBytes), nil
}

// 评估密码强度
func evaluatePasswordStrength(password string) string {
	var score int
	var hasLower, hasUpper, hasNumber, hasSymbol bool

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	// 计算分数
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSymbol {
		score++
	}

	// 根据分数返回强度
	switch {
	case score >= 6:
		return "强"
	case score >= 4:
		return "中等"
	case score >= 2:
		return "弱"
	default:
		return "很弱"
	}
}

// 显示帮助信息
func showHelp() {
	fmt.Println("密码生成器工具")
	fmt.Println("用法: password_generator [选项]")
	fmt.Println("选项:")
	fmt.Println("  -length     密码长度 (默认: 12)")
	fmt.Println("  -count      生成密码数量 (默认: 1)")
	fmt.Println("  -lower      包含小写字母 (默认: true)")
	fmt.Println("  -upper      包含大写字母 (默认: true)")
	fmt.Println("  -numbers    包含数字 (默认: true)")
	fmt.Println("  -symbols    包含特殊字符 (默认: false)")
	fmt.Println("  -exclude    排除易混淆字符 (默认: false)")
	fmt.Println("  -help       显示帮助信息")
	fmt.Println("\n示例:")
	fmt.Println("  生成一个12位包含所有字符类型的密码:")
	fmt.Println("  password_generator -length 12 -symbols true")
	fmt.Println("  生成5个8位不包含特殊字符的密码:")
	fmt.Println("  password_generator -length 8 -count 5 -symbols false")
}

func main() {
	// 解析命令行参数
	length := flag.Int("length", 12, "密码长度")
	count := flag.Int("count", 1, "生成密码数量")
	includeLower := flag.Bool("lower", true, "包含小写字母")
	includeUpper := flag.Bool("upper", true, "包含大写字母")
	includeNumbers := flag.Bool("numbers", true, "包含数字")
	includeSymbols := flag.Bool("symbols", false, "包含特殊字符")
	excludeAmbiguous := flag.Bool("exclude", false, "排除易混淆字符")
	help := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助
	if *help {
		showHelp()
		os.Exit(0)
	}

	// 验证参数
	if *length < 1 {
		fmt.Println("错误: 密码长度必须大于0")
		os.Exit(1)
	}

	if *count < 1 {
		fmt.Println("错误: 生成数量必须大于0")
		os.Exit(1)
	}

	// 创建密码配置
	config := &PasswordConfig{
		Length:           *length,
		IncludeLower:     *includeLower,
		IncludeUpper:     *includeUpper,
		IncludeNumbers:   *includeNumbers,
		IncludeSymbols:   *includeSymbols,
		ExcludeAmbiguous: *excludeAmbiguous,
		Count:            *count,
	}

	// 显示配置信息
	fmt.Printf("密码配置:\n")
	fmt.Printf("  长度: %d\n", config.Length)
	fmt.Printf("  数量: %d\n", config.Count)
	fmt.Printf("  小写字母: %v\n", config.IncludeLower)
	fmt.Printf("  大写字母: %v\n", config.IncludeUpper)
	fmt.Printf("  数字: %v\n", config.IncludeNumbers)
	fmt.Printf("  特殊字符: %v\n", config.IncludeSymbols)
	fmt.Printf("  排除易混淆字符: %v\n", config.ExcludeAmbiguous)
	fmt.Println("------------------------")

	// 生成密码
	for i := 0; i < config.Count; i++ {
		password, err := generatePassword(config)
		if err != nil {
			fmt.Printf("生成密码失败: %v\n", err)
			os.Exit(1)
		}

		strength := evaluatePasswordStrength(password)
		fmt.Printf("密码 %d: %s (强度: %s)\n", i+1, password, strength)
	}

	fmt.Println("------------------------")
	fmt.Println("密码生成完成")
}
