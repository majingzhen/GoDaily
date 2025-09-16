// 基于命令行代办事项管理器
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// Todo 待办事项
type Todo struct {
	Id          int       `json:"id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Completed   bool      `json:"completed"`
}

var (
	todos        []Todo
	filePath     = "todo.json"
	addFlag      string
	listFlag     bool
	delFlag      int
	completeFlag int
)

func init() {
	// 解析命令行参数
	flag.StringVar(&addFlag, "add", "", "添加新的待办事项")
	flag.BoolVar(&listFlag, "list", false, "列出所有待办事项")
	flag.IntVar(&delFlag, "del", 0, "删除指定编号的待办事项")
	flag.IntVar(&completeFlag, "complete", 0, "完成指定编号的待办事项")
	flag.Parse()

	// 加载待办事项
	loadTodos()
}

func main() {
	switch {
	case addFlag != "":
		addTodo(addFlag)
	case listFlag:
		listTodos()
	case delFlag != 0:
		delTodo(delFlag)
	case completeFlag != 0:
		completeTodo(completeFlag)
	default:
		fmt.Println("使用方法:")
		fmt.Println(" - 添加待办: todo -add '要做的事情'")
		fmt.Println(" - 列出所有待办: todo -list")
		fmt.Println(" - 删除待办: todo -del [Id]")
		fmt.Println(" - 完成待办: todo -complete [Id]")
	}
}

// loadTodos 加载待办事项
func loadTodos() {
	file, err := os.Open(filePath)
	if err != nil {
		// 如果文件不存在，初始化空切片
		if os.IsNotExist(err) {
			todos = []Todo{}
			return
		}
		fmt.Printf("加载待办事项失败: %v\n\n", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	if err != nil {
		fmt.Printf("加载待办事项失败: %v\n\n", err)
	}
}

// saveTodos 保存待办事项
func saveTodos() {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("保存待办事项失败: %v\n\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(todos)
	if err != nil {
		fmt.Printf("保存待办事项失败: %v\n\n", err)
	}
}

// addTodo 添加待办事项
func addTodo(content string) {
	id := 1
	if len(todos) > 0 {
		id = todos[len(todos)-1].Id + 1
	}

	todo := Todo{
		Id:        id,
		Content:   content,
		CreatedAt: time.Now(),
		Completed: false,
	}
	todos = append(todos, todo)
	saveTodos()
	fmt.Printf("添加待办事项成功(Id: %d)\n", todo.Id)
}

// listTodos 列出所有待办事项
func listTodos() {
	if len(todos) == 0 {
		fmt.Println("没有待办事项")
		return
	}
	fmt.Println("待办列表:")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, todo := range todos {
		status := " "
		if todo.Completed {
			status = "√"
		}
		fmt.Printf("[%s] %d. %s (创建于: %s)\n",
			status,
			todo.Id,
			todo.Content,
			todo.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

// delTodo 删除指定编号的待办事项
func delTodo(id int) {
	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos()
			fmt.Printf("删除待办事项成功(Id: %d)\n\n", todo.Id)
			return
		}
	}
	fmt.Println("待办事项不存在")
}

// completeTodo 完成指定编号的待办事项
func completeTodo(id int) {
	for i, todo := range todos {
		if todo.Id == id {
			todo.Completed = true
			todo.CompletedAt = time.Now()
			todos[i] = todo
			saveTodos()
			fmt.Printf("完成待办事项(Id: %d)\n\n", todo.Id)
			return
		}
	}
	fmt.Println("待办事项不存在")
}
