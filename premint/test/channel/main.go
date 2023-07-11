package main

import (
	"fmt"
)

func main() {
	// 无缓冲通道
	ch1 := make(chan int)

	// 有缓冲通道，缓冲区大小为 3
	ch2 := make(chan int, 3)

	// 向无缓冲通道发送数据（会阻塞）
	go func() {
		ch1 <- 1
		fmt.Println("Sent 1 to unbuffered channel")
	}()

	// 从无缓冲通道接收数据（会阻塞）
	go func() {

		val := <-ch1
		fmt.Println("Received", val, "from unbuffered channel")
	}()

	// 向有缓冲通道发送数据（不会阻塞）
	go func() {
		ch2 <- 2
		fmt.Println("Sent 2 to buffered channel")
		close(ch2)
	}()

	// 从有缓冲通道接收数据（不会阻塞）
	go func() {
		val := <-ch2
		fmt.Println("Received", val, "from buffered channel")
	}()

	// 等待协程执行完成
	fmt.Scanln()
}
