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
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-425.xlsx")
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
		//handle, _ := wd.CurrentWindowHandle()
		//wd.ResizeWindow(handle, 500, 300)
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
	err = OpenHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString("进入个人主页失败")
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
	}
	time.Sleep(2 * time.Second)

	//进入setting
	err = ClickSetting(wd)
	if err != nil {
		log.Println("进入setting失败")
		dstFile.WriteString("进入setting失败")
		ch <- wrongData
		return err
	} else {
		log.Println("进入setting成功")
	}
	time.Sleep(2 * time.Second)

	//进入了 Edit Profile页面 -->找到Social Link 点击
	err = ClickSocialLink(wd)
	if err != nil {
		log.Println("进入ClickSocialLink失败")
		dstFile.WriteString("进入ClickSocialLink失败")
		ch <- wrongData
		return err
	} else {
		log.Println("进入ClickSocialLink成功")
	}
	time.Sleep(2 * time.Second)
	return
}
func OpenHomePage(wd selenium.WebDriver) (err error) {
	var HomePage selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			HomePage, err = wd.FindElement(selenium.ByCSSSelector, ".campaign-avatar-inner")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 10*time.Second)
	if err != nil {
		return err
	} else {
		//button,_:= wd.FindElement(selenium.ByXPATH,"//*[@id=\"topNavbar\"]/div/div[2]/div[2]/div[1]/div[2]/div/div/svg/g/rect")
		script := "var element = arguments[0];" +
			"var mouseEvent = document.createEvent('MouseEvents');" +
			"mouseEvent.initMouseEvent('mouseover', true, true, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);" +
			"element.dispatchEvent(mouseEvent);"
		_, err = wd.ExecuteScript(script, []interface{}{HomePage})
	}
	return err
}
func ClickSetting(wd selenium.WebDriver) (err error) {
	var settingButton selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			settingButton, err = wd.FindElement(selenium.ByCSSSelector, ".account-setting-menu")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 10*time.Second)
	if err != nil {
		return err
	} else {
		//找到具有setting属性的div
		sbutton, _ := settingButton.FindElement(selenium.ByXPATH, "//div[contains(text(), 'Setting')]")
		sbutton.Click()
	}
	return
}
func ClickSocialLink(wd selenium.WebDriver) (err error) {
	var SocialLink selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			SocialLink, err = wd.FindElement(selenium.ByCSSSelector, ".options-mobile-wrap")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 10*time.Second)
	if err != nil {
		return err
	} else {
		//找到具有setting属性的div
		button, _ := SocialLink.FindElement(selenium.ByXPATH, "//div[contains(text(), 'Social Link')]")
		button.Click()
	}
	return
}
