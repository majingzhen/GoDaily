package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/**
 	老虎机游戏
	结构体（SlotMachine）存储游戏状态
	数组和切片管理图标和转轮
	随机数生成实现老虎机的随机性
	循环和条件判断控制游戏流程
	函数和方法组织游戏逻辑
	错误处理确保输入有效性
*/

// 老虎机图标
const (
	Cherry  = "🍒"
	Lemon   = "🍋"
	Orange  = "🍊"
	Bell    = "🔔"
	Bar     = "🍫"
	Seven   = "7️⃣"
	Diamond = "💎"
)

// SlotMachine  老虎机结构体
type SlotMachine struct {
	Balance int       // 玩家余额
	Reels   [3]string // 三个转轮
	Symbols []string  // 符号列表
}

// NewSlotMachine 创建新的老虎机
func NewSlotMachine(initialBalance int) *SlotMachine {
	return &SlotMachine{
		Balance: initialBalance,
		Symbols: []string{Cherry, Lemon, Orange, Bell, Bar, Seven, Diamond},
	}
}

// Spin 旋转老虎机
func (sm *SlotMachine) Spin(bet int) int {
	// 检查余额是否足够
	if sm.Balance < bet {
		fmt.Println("余额不足! ")
		return 0
	}
	// 扣除积分
	sm.Balance -= bet
	// 随机生成三个图标
	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		index := rand.Intn(len(sm.Symbols))
		sm.Reels[i] = sm.Symbols[index]
	}

	// 显示旋转结果
	sm.DisplayReels()

	// 判断中奖情况并计算积分
	winAmount := sm.calculateWin(bet)

	if winAmount > 0 {
		fmt.Printf("恭喜! 你赢了 %d 积分\n", winAmount)
		sm.Balance += winAmount
	} else {
		fmt.Println("很遗憾，未中奖！")
	}
	return winAmount
}

// DisplayReels 显示旋转结果
func (sm *SlotMachine) DisplayReels() {
	fmt.Println("\n==========")
	fmt.Printf("| %s %s %s |\n", sm.Reels[0], sm.Reels[1], sm.Reels[2])
	fmt.Println("==========")
}

// calculateWin 计算中奖
func (sm *SlotMachine) calculateWin(bet int) int {
	// 三个相同的7是最高奖
	if sm.Reels[0] == Seven && sm.Reels[1] == Seven && sm.Reels[2] == Seven {
		return bet * 100
	}
	// 三个相同的钻石
	if sm.Reels[0] == Diamond && sm.Reels[1] == Diamond && sm.Reels[2] == Diamond {
		return bet * 50
	}
	// 三个相同的其他图标
	if sm.Reels[0] == sm.Reels[1] && sm.Reels[1] == sm.Reels[2] {
		return bet * 10
	}
	// 两个相同
	if sm.Reels[0] == sm.Reels[1] || sm.Reels[1] == sm.Reels[2] || sm.Reels[0] == sm.Reels[2] {
		return bet * 5
	}
	// 至少包含一个7
	if sm.Reels[0] == Seven || sm.Reels[1] == Seven || sm.Reels[2] == Seven {
		return bet * 1
	}
	// 未中奖
	return 0
}

// 显示游戏帮助
func showHelp() {
	fmt.Println("\n===== 游戏帮助 =====")
	fmt.Println("1. 输入赌注金额进行游戏")
	fmt.Println("2. 输入0退出游戏")
	fmt.Println("3. 输入h查看帮助")
	fmt.Println("中奖规则:")
	fmt.Println("- 三个7️⃣: 100倍奖励")
	fmt.Println("- 三个💎: 50倍奖励")
	fmt.Println("- 三个相同其他图标: 10倍奖励")
	fmt.Println("- 两个相邻相同图标: 5倍奖励")
	fmt.Println("- 至少一个7️⃣: 1倍奖励")
}

// 主流程
func main() {
	fmt.Println("===== 欢迎来到老虎机游戏! =====")
	fmt.Println("祝你好运!")

	// 初始积分
	slotMachine := NewSlotMachine(100)

	// 显示帮助
	showHelp()

	for {
		fmt.Printf("\n当前余额: %d 币\n", slotMachine.Balance)
		fmt.Print("请输入赌注金额(输入0退出, h帮助): ")

		var input string
		fmt.Scan(&input)

		// 处理帮助请求
		if input == "h" || input == "H" {
			showHelp()
			continue
		}

		// 转换输入为数字
		bet, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("请输入有效的数字!")
			continue
		}
		if bet == 0 {
			fmt.Println("感谢游玩! ")
			fmt.Printf("最终余额: %d 币\n", slotMachine.Balance)
			os.Exit(0)
		}
		// 检查输入是否有效
		if bet < 0 {
			fmt.Println("投入积分不能为负数!")
			continue
		}
		// 旋转老虎机
		slotMachine.Spin(bet)

		// 检查是否破产
		if slotMachine.Balance <= 0 {
			fmt.Println("游戏结束，你破产了!")
			os.Exit(0)
		}
	}
}
