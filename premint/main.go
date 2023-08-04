/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-02-20 15:26:12
 * @LastEditTime: 2023-06-06 11:33:01
 */
package main

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/src/Galxe"
	"github.com/JianLinWei1/premint-selenium/src/Galxe1/TwitterBind"
	"github.com/JianLinWei1/premint-selenium/src/OMNI"
	"github.com/JianLinWei1/premint-selenium/src/Taiko"
	"github.com/JianLinWei1/premint-selenium/src/premint"
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
	util.SetLog(func() {
		toMain()
	})
}
func toMain() {
	cmd := wdservice.InitCmd()
	switch cmd {
	case 1:
		premint.Start()
		break
	case 6:
		Taiko.Taiko1Start()
	case 7:
		OMNI.OmniStart()
	case 8:
		TwitterBind.RemoveAndBind()
	case 9:
		Galxe.Remove()
	case 11:
		OMNI.OmniCheckStart()
	case 13:
		Galxe.VerifyEmail()
	case 14:
		TwitterBind.Check()
	}
}

func isExit() {
	fmt.Println("是否继续操作(y/n)：")
	var isExit string
	fmt.Scanln(&isExit)
	if strings.EqualFold(isExit, "y") {
		toMain()
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
