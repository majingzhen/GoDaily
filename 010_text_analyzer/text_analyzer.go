package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

/**
文本分析工具
	使用Go基础语法实现文本文件分析
	支持统计文本的多种指标
	包括字符数、单词数、行数、段落数等
	支持词频分析和文本摘要
	展示Go语言文件操作、正则表达式、排序等特性
*/

// TextStats 文本统计结构体
type TextStats struct {
	Filename       string          // 文件名
	Characters     int             // 字符数（包含空格）
	CharactersNoWS int             // 字符数（不含空格）
	Words          int             // 单词数
	Lines          int             // 行数
	Paragraphs     int             // 段落数
	Sentences      int             // 句子数
	ChineseChars   int             // 中文字符数
	EnglishWords   int             // 英文单词数
	Numbers        int             // 数字字符数
	Punctuation    int             // 标点符号数
	WordFreq       map[string]int  // 词频统计
	TopWords       []WordFrequency // 高频词汇
	ReadingTime    float64         // 预估阅读时间（分钟）
}

// WordFrequency 词频结构体
type WordFrequency struct {
	Word  string
	Count int
}

// AnalysisConfig 分析配置结构体
type AnalysisConfig struct {
	ShowWordFreq  bool   // 显示词频分析
	TopWordsCount int    // 显示高频词数量
	MinWordLength int    // 最小单词长度
	IgnoreCase    bool   // 忽略大小写
	OutputFormat  string // 输出格式 (text/json)
	ReadingSpeed  int    // 阅读速度（字/分钟）
}

// 常用停用词（中英文）
var stopWords = map[string]bool{
	"的": true, "了": true, "在": true, "是": true, "我": true, "有": true, "和": true, "就": true,
	"不": true, "人": true, "都": true, "一": true, "一个": true, "上": true, "也": true, "很": true,
	"到": true, "说": true, "要": true, "去": true, "你": true, "会": true, "着": true, "没有": true,
	"看": true, "好": true, "自己": true, "这": true, "她": true, "他": true, "但是": true, "那": true,
	"the": true, "be": true, "to": true, "of": true, "and": true, "a": true, "in": true, "that": true,
	"have": true, "i": true, "it": true, "for": true, "not": true, "on": true, "with": true, "he": true,
	"as": true, "you": true, "do": true, "at": true, "this": true, "but": true, "his": true, "by": true,
	"from": true, "they": true, "we": true, "say": true, "her": true, "she": true, "or": true, "an": true,
}

// 分析文本文件
func analyzeText(filename string, config *AnalysisConfig) (*TextStats, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	text := string(content)
	stats := &TextStats{
		Filename: filename,
		WordFreq: make(map[string]int),
	}

	// 基础统计
	stats.Characters = utf8.RuneCountInString(text)
	stats.CharactersNoWS = utf8.RuneCountInString(strings.ReplaceAll(strings.ReplaceAll(text, " ", ""), "\t", ""))
	stats.Lines = countLines(text)
	stats.Paragraphs = countParagraphs(text)
	stats.Sentences = countSentences(text)

	// 字符类型统计
	for _, char := range text {
		switch {
		case unicode.Is(unicode.Han, char):
			stats.ChineseChars++
		case unicode.IsDigit(char):
			stats.Numbers++
		case unicode.IsPunct(char):
			stats.Punctuation++
		}
	}

	// 单词统计和词频分析
	words := extractWords(text, config)
	stats.Words = len(words)

	englishWordCount := 0
	for _, word := range words {
		if isEnglishWord(word) {
			englishWordCount++
		}

		// 词频统计
		if config.ShowWordFreq {
			key := word
			if config.IgnoreCase {
				key = strings.ToLower(word)
			}
			if len(key) >= config.MinWordLength && !stopWords[key] {
				stats.WordFreq[key]++
			}
		}
	}
	stats.EnglishWords = englishWordCount

	// 生成高频词列表
	if config.ShowWordFreq {
		stats.TopWords = getTopWords(stats.WordFreq, config.TopWordsCount)
	}

	// 计算预估阅读时间
	totalReadableChars := stats.ChineseChars + stats.EnglishWords
	stats.ReadingTime = float64(totalReadableChars) / float64(config.ReadingSpeed)

	return stats, nil
}

// 统计行数
func countLines(text string) int {
	scanner := bufio.NewScanner(strings.NewReader(text))
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}

// 统计段落数
func countParagraphs(text string) int {
	paragraphs := strings.Split(strings.TrimSpace(text), "\n\n")
	return len(paragraphs)
}

// 统计句子数
func countSentences(text string) int {
	// 使用正则表达式匹配句子结尾
	re := regexp.MustCompile(`[.!?。！？]+`)
	matches := re.FindAllString(text, -1)
	return len(matches)
}

// 提取单词
func extractWords(text string, config *AnalysisConfig) []string {
	// 使用正则表达式提取单词（支持中英文）
	re := regexp.MustCompile(`[\p{L}\p{N}]+`)
	words := re.FindAllString(text, -1)

	var filteredWords []string
	for _, word := range words {
		if len(word) >= config.MinWordLength {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}

// 判断是否为英文单词
func isEnglishWord(word string) bool {
	for _, char := range word {
		if char > 127 {
			return false
		}
	}
	return true
}

// 获取高频词
func getTopWords(wordFreq map[string]int, topCount int) []WordFrequency {
	var words []WordFrequency
	for word, count := range wordFreq {
		words = append(words, WordFrequency{Word: word, Count: count})
	}

	// 按频率排序
	sort.Slice(words, func(i, j int) bool {
		return words[i].Count > words[j].Count
	})

	// 返回前N个
	if len(words) < topCount {
		topCount = len(words)
	}
	return words[:topCount]
}

// 显示统计结果
func displayStats(stats *TextStats, config *AnalysisConfig) {
	fmt.Printf("文本分析结果: %s\n", stats.Filename)
	fmt.Println("========================================")

	// 基础统计
	fmt.Printf("📄 基础统计:\n")
	fmt.Printf("  总字符数: %d\n", stats.Characters)
	fmt.Printf("  有效字符数: %d (不含空格)\n", stats.CharactersNoWS)
	fmt.Printf("  单词数: %d\n", stats.Words)
	fmt.Printf("  行数: %d\n", stats.Lines)
	fmt.Printf("  段落数: %d\n", stats.Paragraphs)
	fmt.Printf("  句子数: %d\n", stats.Sentences)
	fmt.Println()

	// 字符类型统计
	fmt.Printf("🔤 字符类型统计:\n")
	fmt.Printf("  中文字符: %d\n", stats.ChineseChars)
	fmt.Printf("  英文单词: %d\n", stats.EnglishWords)
	fmt.Printf("  数字字符: %d\n", stats.Numbers)
	fmt.Printf("  标点符号: %d\n", stats.Punctuation)
	fmt.Println()

	// 阅读时间
	fmt.Printf("⏱️  预估阅读时间: %.1f 分钟\n", stats.ReadingTime)
	fmt.Println()

	// 词频分析
	if config.ShowWordFreq && len(stats.TopWords) > 0 {
		fmt.Printf("📊 高频词汇 (前 %d 个):\n", len(stats.TopWords))
		for i, word := range stats.TopWords {
			fmt.Printf("%2d. %-15s %d 次\n", i+1, word.Word, word.Count)
		}
		fmt.Println()
	}

	// 文本复杂度评估
	fmt.Printf("📈 文本复杂度评估:\n")
	avgWordsPerSentence := float64(stats.Words) / float64(stats.Sentences)
	avgCharsPerWord := float64(stats.CharactersNoWS) / float64(stats.Words)

	fmt.Printf("  平均句长: %.1f 个单词\n", avgWordsPerSentence)
	fmt.Printf("  平均词长: %.1f 个字符\n", avgCharsPerWord)

	// 复杂度评级
	complexity := "简单"
	if avgWordsPerSentence > 20 || avgCharsPerWord > 6 {
		complexity = "复杂"
	} else if avgWordsPerSentence > 15 || avgCharsPerWord > 4 {
		complexity = "中等"
	}
	fmt.Printf("  复杂度等级: %s\n", complexity)
}

// 显示帮助信息
func showHelp() {
	fmt.Println("文本分析工具")
	fmt.Println("用法: text_analyzer [选项] <文件路径>")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -freq          显示词频分析 (默认: false)")
	fmt.Println("  -top           显示高频词数量 (默认: 10)")
	fmt.Println("  -minlen        最小单词长度 (默认: 2)")
	fmt.Println("  -ignore-case   忽略大小写 (默认: true)")
	fmt.Println("  -speed         阅读速度(字/分钟) (默认: 200)")
	fmt.Println("  -help          显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  分析文本文件基础信息:")
	fmt.Println("  text_analyzer document.txt")
	fmt.Println()
	fmt.Println("  分析并显示词频统计:")
	fmt.Println("  text_analyzer -freq -top 20 document.txt")
	fmt.Println()
	fmt.Println("  自定义分析参数:")
	fmt.Println("  text_analyzer -freq -top 15 -minlen 3 -speed 180 document.txt")
}

func main() {
	// 解析命令行参数
	showWordFreq := flag.Bool("freq", false, "显示词频分析")
	topWordsCount := flag.Int("top", 10, "显示高频词数量")
	minWordLength := flag.Int("minlen", 2, "最小单词长度")
	ignoreCase := flag.Bool("ignore-case", true, "忽略大小写")
	readingSpeed := flag.Int("speed", 200, "阅读速度(字/分钟)")
	help := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助
	if *help {
		showHelp()
		os.Exit(0)
	}

	// 检查文件参数
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("错误: 请指定要分析的文件路径")
		fmt.Println("使用 -help 查看使用说明")
		os.Exit(1)
	}

	filename := args[0]

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("错误: 文件 '%s' 不存在\n", filename)
		os.Exit(1)
	}

	// 创建分析配置
	config := &AnalysisConfig{
		ShowWordFreq:  *showWordFreq,
		TopWordsCount: *topWordsCount,
		MinWordLength: *minWordLength,
		IgnoreCase:    *ignoreCase,
		ReadingSpeed:  *readingSpeed,
	}

	// 分析文本
	fmt.Printf("正在分析文件: %s\n\n", filename)
	stats, err := analyzeText(filename, config)
	if err != nil {
		fmt.Printf("分析失败: %v\n", err)
		os.Exit(1)
	}

	// 显示结果
	displayStats(stats, config)

	fmt.Println("========================================")
	fmt.Println("分析完成!")
}
