package bitbrowser

import (
	"github.com/tebeka/selenium"
	"log"
)

func CloseOtherLabels(wd selenium.WebDriver, target string) {
	handles, err := wd.WindowHandles()
	if err != nil {
		log.Println(err)
	}
	//log.Println("窗口代码：", handles)
	//log.Println("标签数量：", len(handles))
	//关闭除工作台的其他标签页
	log.Println("关闭除目标的其他标签页")
	for _, v := range handles {
		if v != target {
			wd.SwitchWindow(v)
			wd.Close()
		}
	}
	//for i := len(handles) - 2; i >= 0; i-- {
	//	newHandle := handles[i]
	//	wd.SwitchWindow(newHandle)
	//	wd.Close()
	//}
}
