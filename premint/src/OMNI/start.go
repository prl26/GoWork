package main

import (
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"log"
	"strconv"
	"time"
)

func Start() {
	excelInfos := util.GetOMNIExcelInfos("D:\\Go Work\\resource\\1600个钱包地址2023.07.04.xlsx")
	chs := make(chan string, len(excelInfos))

	for i := 0; i < len(excelInfos); i++ {
		time.Sleep(1 * time.Second)
		//打开浏览器窗口
		wd, _ := wdservice.InitWd(i, excelInfos[i].BitId)
		if wd != nil {
			go util.SetLog(func() {
				OmniClick(wd, excelInfos[i].Address)
			})
		} else {
			log.Println("第" + strconv.Itoa(i+1) + "条数据浏览器初始化失败****")
		}

	}
	bitbrowser.Windowbounds()
	for t := 0; t < len(excelInfos); t++ {
		log.Println(<-chs)
	}
}
func testStart() {
	time.Sleep(1 * time.Second)
	//打开浏览器窗口
	wd, _ := wdservice.InitWd(0, "96d8d0676da04a2bb4e1138b5052085d")
	if wd != nil {
		go util.SetLog(func() {
			OmniClick(wd, "0x243749Bf42346bB9A122034E81281819AA2CFbb6")
		})
	} else {
		log.Println("第" + strconv.Itoa(1) + "条数据浏览器初始化失败****")
	}

	bitbrowser.Windowbounds()

}
