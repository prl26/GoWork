package main

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var wg sync.WaitGroup

//go:generate goversioninfo -64 -gofilepackage="main"
func main() {
	count := 0
	util.SetLog(func() {
		toMain(count)
	})
}
func toMain(count int) {
	cmd := wdservice.InitCmd()
	switch cmd {

	case 9:
		fmt.Println("9")
	case 11:
		fmt.Println("10")
	}
	if count == 0 {
		count++
		fmt.Println("1")
	}
}

func isExit() {
	fmt.Println("是否继续操作(y/n)：")
	var isExit string
	fmt.Scanln(&isExit)
	if strings.EqualFold(isExit, "y") {
		toMain(1)
	}
}
func init() {

	//util.CheckRole()

	file := "./" + "log" + ".log"
	os.Remove(file)
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	// 设置log同时输出到文件和http接口
	//log.SetOutput(io.MultiWriter(file, &httpClient{apiUrl}))

	//log.SetPrefix("[qSkipTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
			log.Println(string(debug.Stack()))
		}
	}()
}
