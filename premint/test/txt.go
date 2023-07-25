package main

import (
	"bufio"
	"fmt"
	"os"
)

//	func main() {
//		fmt.Println("------------")
//		//        cmd := exec.Command("F:\\最新版本游戏\\test1\\test\\test.exe", arg...)
//		cmd := exec.Command("C:\\Program Files (x86)\\AutoIt3\\SciTE\\post.exe", "C:\\image\\1.png")
//		_, err := cmd.CombinedOutput()
//		if err != nil {
//			fmt.Println("Error:", err)
//			return
//		}
//	}
func main() {
	filePath := "C:\\name1600.txt"
	res := getDescribeTxt(filePath)
	for k, _ := range res {
		fmt.Print(res[k])
	}
}
func getDescribeTxt(filePath string) (infos []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败 = ", err)
	}
	defer file.Close()              // 关闭文本流
	reader := bufio.NewReader(file) // 读取文本数据
	for {
		line, err := reader.ReadString('\n') // 读取直到遇到换行符
		if err != nil {
			break // 文件读取完毕或发生错误时退出循环
		}

		//if strings.Contains(line, " ") {
		//	line = line[strings.Index(line, " ")+1:] // 去除序号和空格
		infos = append(infos, line)

	}
	return
}
