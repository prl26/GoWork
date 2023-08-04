package main

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"sync"
	"time"
)

// https://twitter.com/settings/push_notifications
var wg sync.WaitGroup

func main() {
	url := "https://twitter.com/settings/push_notifications"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-419.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\Check-7.31.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\Check-7.31.txt"
	TxtSuccessOut := "D:\\GoWork\\resource\\SuccessInfos\\Check-7.31.txt"
	dstFile, err := os.OpenFile(TxtfileOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer dstFile.Close()
	successFile, err := os.OpenFile(TxtSuccessOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer successFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"助记词", "私钥", "公钥", "地址", "类型", "窗口ID", "MetaMask密码"}
	//创建新的excel文件
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &title)

	//定义一次开多少线程
	fmt.Println("数据长度--------", len(excelInfos))
	// 获取内容并写入Excel
	go func() {
		for t := 0; t < len(excelInfos); t++ {
			data := <-chs
			log.Println("接受到一条错误信息：", data)
			axis := fmt.Sprintf("A%d", t+2)
			excel.SetSheetRow("Sheet1", axis, &data)
		}
	}()

	////单个打开
	for k, v := range excelInfos {
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		if wd != nil {
			handle, _ := wd.WindowHandles()
			if len(handle) > 1 {
				handle1 := util.GetCurrentWindowAndReturn(wd)
				//关闭多余标签页
				bitbrowser.CloseOtherLabels(wd, handle1)
				wd.SwitchWindow(handle1)
			}
			time.Sleep(1 * time.Second)
			wg.Add(1)
			go util.SetLog(func() {
				defer wg.Done()
				err := notification(v, k, chs, wd, url, dstFile, successFile)
				if err != nil {
					log.Println("!-------!", v.BitId, "失败")
				}
				defer bitbrowser.CloseBrower(v.BitId)
			})
		}
		wg.Wait()
	}
	close(chs)
	err = excel.SaveAs(filepout)
	if err != nil {
		log.Println("excel 保存失败----", err)
	} else {
		log.Println("excel 保存成功----", err)

	}
}
func notification(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, successFile *os.File) (err error) {
	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	//打开银河链接
	err = wd.Get(url)
	if err != nil {
		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v---%v", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("银河打开成功")
	}
	time.Sleep(2 * time.Second)

	return err
}
