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
	"strings"
	"sync"
	"time"
)

var wg1 sync.WaitGroup
var wg sync.WaitGroup

//1.打开银河
//2.打开个人首页

func Remove() {
	url := "https://galxe.com/OmniNetwork/campaign/GCSmgUW7Fo"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\remove-7.24.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\remove-7.24.txt"
	dstFile, _ := os.Create(TxtfileOut)
	defer dstFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"地址", "类型", "窗口ID", "MetaMask密码"}
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

	//批量打开
	//var Ids []string
	//var wds []selenium.WebDriver
	//var infos []model.OMNIExcelInfo
	//for k, v := range excelInfos {
	//	Ids = append(Ids, v.BitId)
	//	fmt.Println("----------", v.Address)
	//	//打开比特浏览器
	//	if len(Ids) == 10 {
	//		//先批量打开浏览器
	//		for _, id := range Ids {
	//			time.Sleep(2 * time.Second)
	//			wg1.Add(1)
	//			go util.SetLog(func() {
	//				defer wg1.Done()
	//				wd, _ := wdservice.InitWd(k, id)
	//				if wd != nil {
	//					handle, _ := wd.CurrentWindowHandle()
	//					bitbrowser.CloseOtherLabels(wd, handle)
	//					wds = append(wds, wd)
	//					wg.Add(1)
	//					infos = append(infos, v)
	//				} else {
	//					log.Println(v.BitId, "-----窗口打开失败")
	//					dstFile.WriteString(fmt.Sprintf("%v-----窗口打开失败", v.BitId))
	//					wrongData := []string{v.Address, v.Type, v.BitId, v.MetaPwd}
	//					chs <- wrongData
	//				}
	//			})
	//		}
	//		//等待一下
	//		wg1.Wait()
	//		//开始处理
	//		for k1, v1 := range infos {
	//			wg.Add(1)
	//			go util.SetLog(func() {
	//				defer wg.Done()
	//				wd, _ := wdservice.InitWd(k1, v1.BitId)
	//				err := RemoveTwitter(v, k, chs, wd, url, dstFile)
	//				if err != nil {
	//					log.Println("!-------!", v.BitId, "失败")
	//				}
	//			})
	//		}
	//		wg.Wait()
	//		for _, v2 := range infos {
	//			bitbrowser.CloseBrower(v2.BitId)
	//		}
	//		infos = infos[:0]
	//	}
	//
	//}

	//单个打开
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
				err := RemoveTwitter(v, k, chs, wd, url, dstFile)
				if err != nil {
					log.Println("!-------!", v.BitId, "失败")
				}
				defer bitbrowser.CloseBrower(v.BitId)
			})
		}
		wg.Wait()
	}
	close(chs)
	err := excel.SaveAs(filepout)
	if err != nil {
		log.Println("excel 保存失败----", err)
	}
}
func RemoveTwitter(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	//打开银河链接

	err = wd.Get(url)
	if err != nil {
		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("银河打开成功")
	}
	time.Sleep(2 * time.Second)
	//小狐狸登陆
	nowHandle, _ := wd.WindowHandles()
	if len(nowHandle) > 1 {
		wd.SwitchWindow(nowHandle[1])
		err = SmallFoxLogin(wd)
	}
	if err != nil {
		log.Println("小狐狸登陆出错-----", excelInfo.BitId)
		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("小狐狸登陆成功")
	}
	time.Sleep(1 * time.Second)
	//打开个人首页
	handleNow, _ := wd.WindowHandles()
	wd.SwitchWindow(handleNow[0])
	log.Println(len(handleNow))
	err = wd.MaximizeWindow(handleNow[0])
	if err != nil {
		log.Println("窗口方法失败", err)
	}

	err = OpenHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString(fmt.Sprintf("进入个人主页失败-----%v", excelInfo.BitId))
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
		dstFile.WriteString(fmt.Sprintf("进入setting失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("进入setting成功")
	}
	time.Sleep(3 * time.Second)
	//这里可能有小狐狸弹窗
	ShandleNow, _ := wd.WindowHandles()
	if len(ShandleNow) > 1 {
		closeSignInPop(wd, ShandleNow[1])
		time.Sleep(2 * time.Second)
		wd.SwitchWindow(handleNow[0])
	}
	//进入了 Edit Profile页面 -->找到Social Link 点击
	err = ClickSocialLink(wd, handleNow[0])
	if err != nil {
		log.Println("进入ClickSocialLink失败")
		dstFile.WriteString(fmt.Sprintf("进入ClickSocialLink失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("进入ClickSocialLink成功")
	}
	time.Sleep(2 * time.Second)

	//找到twiteer,再找到remove
	err = ClickRemovek(wd)
	if err != nil {
		log.Println("点击remove失败")
		dstFile.WriteString(fmt.Sprintf("点击remove失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("点击remove成功")
	}
	time.Sleep(2 * time.Second)
	return
}
func closeSignInPop(wd selenium.WebDriver, handle string) {
	wd.SwitchWindow(handle)
	SignIn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
	if err == nil {
		SignIn.Click()
		log.Println("关闭小狐狸弹窗")
	}
}
func SmallFoxLogin(wd selenium.WebDriver) (err error) {
	var Password selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			Password, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"password\"]")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, err
	}, 4*time.Second)
	if err != nil {
		log.Println("没有找到登陆按钮")
		return err
	} else {
		Password.SendKeys("SHIfeng0615")
		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
		if err == nil {
			UnLock.Click()
		}
		//var SignIn1 selenium.WebElement
		//err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		//	for i := 0; i < 10; i++ {
		//		SignIn1, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		//		if err != nil {
		//			time.Sleep(1 * time.Second)
		//			continue
		//		} else {
		//			return true, nil
		//		}
		//	}
		//	return false, errors.New("Fail")
		//}, 2*time.Second)
		//if err != nil {
		//	Next, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
		//	if err == nil {
		//		Next.Click()
		//		time.Sleep(1 * time.Second)
		//		Connect, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		//		if err == nil {
		//			Connect.Click()
		//		}
		//		var SignIn selenium.WebElement
		//		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		//			for i := 0; i < 10; i++ {
		//				log.Println("寻找signIn")
		//				mHandle, _ := wd.WindowHandles()
		//				log.Println(len(mHandle))
		//				if len(mHandle) > 1 {
		//					wd.SwitchWindow(mHandle[1])
		//				}
		//				SignIn, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		//				if err != nil {
		//					time.Sleep(500 * time.Millisecond)
		//					continue
		//				} else {
		//					return true, nil
		//				}
		//			}
		//			return false, err
		//		}, 3*time.Second)
		//		if err == nil {
		//			SignIn.Click()
		//		}
		//		time.Sleep(1 * time.Second)
		//		SignIn2, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		//		if err == nil {
		//			SignIn2.Click()
		//		}
		//	}
		//} else {
		//	err = SignIn1.Click()
		//	if err != nil {
		//		log.Println("metamask登陆出错")
		//		return err
		//	}
		//	time.Sleep(2 * time.Second)
		//	SignIn2, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		//	if err == nil {
		//		SignIn2.Click()
		//	}
		//}

	}
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
		err = sbutton.Click()
	}
	return err
}
func ClickSocialLink(wd selenium.WebDriver, handle string) (err error) {
	var SocialLink selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			ShandleNow, _ := wd.WindowHandles()
			if len(ShandleNow) > 1 {
				closeSignInPop(wd, ShandleNow[1])
				time.Sleep(2 * time.Second)
				wd.SwitchWindow(handle)
			}
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
func ClickRemovek(wd selenium.WebDriver) (err error) {
	var removeTwitter []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			removeTwitter, err = wd.FindElements(selenium.ByCSSSelector, ".social-account-box")
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
		//找到twitter在的div,并点击remove
		if len(removeTwitter) > 0 {
			for _, v := range removeTwitter {
				SaName, _ := v.FindElement(selenium.ByCSSSelector, ".sa-name")
				text, _ := SaName.Text()
				if strings.Contains(text, "Twitter") {
					actions, _ := v.FindElements(selenium.ByCSSSelector, ".sa-action")
					if len(actions) > 0 {
						err = actions[0].Click()
						if err != nil {
							log.Println("点击remove失败")
						} else {
							log.Println("点击remove成功")
							return nil
						}
					}
				} else {
					continue
				}
			}
		} else {
			log.Println("没有找到对应的box,已经remove了")
			return nil
		}
	}
	return
}
