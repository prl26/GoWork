package Galxe

import (
	"errors"
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

var wg sync.WaitGroup

//1.打开银河
//2.打开个人首页

func Remove() {
	url := "https://galxe.com/OmniNetwork/campaign/GCSmgUW7Fo"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-419.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\removeFailInfos.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\removeFailInfos.txt"
	dstFile, _ := os.Create(TxtfileOut)
	defer dstFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"地址", "类型", "窗口ID", "MetaMask密码"}
	//创建新的excel文件
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &title)

	//定义一次开多少线程
	Size := 10
	counter := 0
	var Ids []string

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

	for k, v := range excelInfos {
		counter++
		Ids = append(Ids, v.BitId)
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		handle, _ := wd.CurrentWindowHandle()
		wd.ResizeWindow(handle, 500, 300)
		if wd != nil {
			wg.Add(1)
			go util.SetLog(func() {
				defer wg.Done()
				err := RemoveTwitter(v, k, chs, wd, url, dstFile)
				if err != nil {
					log.Println("!-------!", v.BitId, "失败")
				}
			})
		}
		if counter >= Size && (counter-1)%Size == 0 {
			bitbrowser.WindowboundsByPara()
			log.Println("counter------", counter)
			wg.Wait()
			bitbrowser.WindowboundsByPara()
			time.Sleep(2 * time.Second)
			//for _, v := range Ids {
			//	bitbrowser.CloseBrower(v)
			//}
			Ids = Ids[:0]
		}
	}

	//for k, v := range excelInfos {
	//	fmt.Println("----------", v.Address)
	//	//打开比特浏览器
	//	wd, _ := wdservice.InitWd(k, v.BitId)
	//
	//	if wd != nil {
	//		wg.Add(1)
	//		go util.SetLog(func() {
	//			defer wg.Done()
	//			StartOmniGalxe(v, k, chs, wd, url)
	//			//bitbrowser.CloseBrower(v.BitId)
	//		})
	//	}
	//	bitbrowser.WindowboundsByPara()
	//	wg.Wait()
	//}
	wg.Wait()
	close(chs)
	excel.SaveAs(filepout)

}
func RemoveTwitter(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
	//打开个人首页

	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	handleNow, _ := wd.CurrentWindowHandle()
	wd.MaximizeWindow(handleNow)
	err = OpenHomePage(wd, dstFile)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString("进入个人主页失败")
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
		return nil
	}

	return
}
func OpenHomePage(wd selenium.WebDriver, dstFile *os.File) (err error) {
	var button selenium.WebElement
	time.Sleep(1 * time.Second)
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			log.Println("开始第", i, "次")
			button, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div/div/div/div[1]/div[1]/div[1]/div/div[2]/div[2]/div[1]/div[2]/div/div/svg/g/rect")
			if err != nil {
				time.Sleep(1 * time.Second)
				return false, err
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 10*time.Second)
	if err != nil {
		log.Println("meiyou")
		return err
	} else {
		//button,_:= wd.FindElement(selenium.ByXPATH,"//*[@id=\"topNavbar\"]/div/div[2]/div[2]/div[1]/div[2]/div/div/svg/g/rect")
		button.Click()
		time.Sleep(3 * time.Second)
	}

	return err
}
