package main

import (
	"fmt"
	"os"
)

// 定义文件路径

func main() {
	dstFile := "config.txt"

	// 创建或打开文件，如果文件不存在则创建，如果文件存在则截断内容
	//file, err := os.Create(dstFile)
	//追加模式
	file, err := os.OpenFile(dstFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString("这是一个示例文本2。\n")
	if err != nil {
		fmt.Println("写入文件出错:", err)
		return
	}

	fmt.Println("文件已创建或截断，内容写入成功:", dstFile)
}
