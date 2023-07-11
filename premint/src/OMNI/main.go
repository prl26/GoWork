package main

import (
	"github.com/JianLinWei1/premint-selenium/src/util"
	"io"
	"log"
	"os"
	"runtime/debug"
)

var (
	Count        int = 20
	WarningCount int = 0
	RoundCount   int = 0
	AddressCount int = 0
	ClaimCount   int = 0
	RepatedCount int = 4
	RepatedCheck int = 3
)

func main() {

	util.SetLog(func() {
		toMain()
	})
}
func toMain() {

	Start()

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
