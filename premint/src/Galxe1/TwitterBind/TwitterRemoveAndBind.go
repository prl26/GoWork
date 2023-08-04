package TwitterBind

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/atotto/clipboard"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// https://galxe.com/accountSetting?tab=SocialLinlk
var wg sync.WaitGroup

func RemoveAndBind() {
	url := "https://galxe.com/accountSetting?tab=SocialLinlk"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\Bind-fail.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\Bind-fail.txt"
	TxtSuccessOut := "D:\\GoWork\\resource\\SuccessInfos\\Bind.txt"
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
				err := removeAndBind(v, k+10, chs, wd, url, dstFile, successFile)
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
func removeAndBind(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, successFile *os.File) (err error) {
	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	//打开银河链接
	err = wd.Get(url)
	if err != nil {
		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v---%v\n", excelInfo.BitId, i))
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
		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("小狐狸登陆成功")
	}
	time.Sleep(1 * time.Second)
	//打开个人首页
	handleNow, _ := wd.WindowHandles()
	wd.SwitchWindow(handleNow[0])
	err = openHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString(fmt.Sprintf("进入个人主页失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
	}
	time.Sleep(1 * time.Second)

	//先保持未登陆状态
	//1.点击DisConnect
	err = ClickDisconnect(wd)
	if err != nil {
		log.Println("点击DisConnect失败")
		dstFile.WriteString(fmt.Sprintf("点击DisConnect失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("点击DisConnect成功")
	}
	time.Sleep(1 * time.Second)

	//这里resize一下
	wd.ResizeWindow(handleNow[0], 1500, 1440)

	//然后连接小狐狸
	err = ConnectMetamask(wd)
	if err != nil {
		log.Println("连接小狐狸失败")
		dstFile.WriteString(fmt.Sprintf("连接小狐狸失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("连接小狐狸成功")
	}
	time.Sleep(2 * time.Second)
	handleNow1, _ := wd.WindowHandles()
	if len(handleNow1) > 1 {
		wd.SwitchWindow(handleNow1[1])
	}
	//处理小狐狸登陆
	err = MetamaskLogin(wd)
	if err != nil {
		log.Println("处理小狐狸登陆失败")
		dstFile.WriteString(fmt.Sprintf("处理小狐狸登陆失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("处理小狐狸登陆成功")
		time.Sleep(2 * time.Second)
	}
	//如果还有弹窗
	err = CloseSignInPop(wd)
	if err != nil {
		log.Println("关闭小狐狸弹窗失败")
		dstFile.WriteString(fmt.Sprintf("处理小狐狸登陆失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("关闭小狐狸弹窗成功")
		time.Sleep(1 * time.Second)
	}
	//判断是否绑定了findConnectTwitter
	count, err := findConnectTwitter(wd)
	if err != nil {
		log.Println("判断是否绑定失败")
		dstFile.WriteString(fmt.Sprintf("判断是否绑定失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("判断是否绑定成功")
		time.Sleep(1 * time.Second)
	}
	//如果已经绑定了就解绑
	if count == 1 {
		log.Println("开始解绑")
		err = ClickRemove(wd)
		if err != nil {
			log.Println("点击remove失败")
			dstFile.WriteString(fmt.Sprintf("点击remove失败-----%v---%v\n", excelInfo.BitId, i))
			ch <- wrongData
			return err
		} else {
			log.Println("点击remove成功")
		}
		time.Sleep(2 * time.Second)
	}
	//开始绑定Twitter  ClickConnectTwitter
	err = ClickConnectTwitter(wd)
	if err != nil {
		log.Println("开始绑定Twitter失败")
		dstFile.WriteString(fmt.Sprintf("开始绑定Twitter失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("开始绑定Twitter成功")
		time.Sleep(2 * time.Second)
	}
	err = CloseSignInPop(wd)
	if err != nil {
		log.Println("关闭小狐狸弹窗失败")
		dstFile.WriteString(fmt.Sprintf("处理小狐狸登陆失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("关闭小狐狸弹窗成功")
		time.Sleep(1 * time.Second)
	}
	failinfo, err := TweetTwitter(wd)
	if strings.Contains(failinfo, "Something went wrong, but don’t fret — let’s give it another shot.") {
		log.Println("Something went wrong, but don’t fret — let’s give it another shot.")
		dstFile.WriteString(fmt.Sprintf("Something went wrong, but don’t fret — let’s give it another shot.-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else if err != nil {
		log.Println("Tweet Twitter失败")
		dstFile.WriteString(fmt.Sprintf("Tweet Twitter失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("Tweet Twitter成功")
		hanNow, _ := wd.WindowHandles()
		wd.SwitchWindow(hanNow[1])
		wd.MaximizeWindow(hanNow[1])
		time.Sleep(2 * time.Second)
	}
	//先复制个空的
	clipboard.WriteAll("")
	userName, err := GetInProfile(wd)
	if err != nil {
		log.Println("获取Link失败")
		dstFile.WriteString(fmt.Sprintf("获取Link失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("获取Link成功,--username", userName)
		time.Sleep(2 * time.Second)
		handowNow, _ := wd.WindowHandles()
		wd.SwitchWindow(handowNow[0])
	}

	//判定成功的弹窗是否出现
	detail, err := Sendkey(wd)
	if strings.Contains(detail, "Twitter account verified successfully") {
		log.Println("第一次Verify成功--info", detail)
	} else {
		log.Println("第一次Verify失败,info--", detail)
		dstFile.WriteString(fmt.Sprintf("第一次Verify失败-----%v-----%v---%v\n", excelInfo.BitId, i, detail))
		ch <- wrongData
		return err
	}
	time.Sleep(2 * time.Second)

	//进入setting 二次判定
	err = ClickSocialLink(wd)
	if err != nil {
		log.Println("第一次Verify成功后进去setting失败")
		dstFile.WriteString(fmt.Sprintf("第一次Verify成功后进去setting失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("第一次Verify成功后进入setting成功，开始第二次判定")
		time.Sleep(2 * time.Second)
	}
	//开始username和绑定的账号校验
	BindUsername, err := findTwitterUsername(wd)
	if err != nil || BindUsername == "" {
		log.Println("获取绑定的账号失败")
		dstFile.WriteString(fmt.Sprintf("获取绑定的账号失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("获取绑定的账号成功--", BindUsername)
		time.Sleep(1 * time.Second)
	}

	if userName == BindUsername {
		log.Println("第二次verify成功")
		successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
		return nil
	} else {
		log.Println("第二次verify失败")
		dstFile.WriteString(fmt.Sprintf("第二次verify失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	}

	return err
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

	}
	return
}
func openHomePage(wd selenium.WebDriver) (err error) {
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

// 点击disconnect
func ClickDisconnect(wd selenium.WebDriver) (err error) {
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
		sbutton, _ := settingButton.FindElement(selenium.ByXPATH, "//div[contains(text(), 'Disconnect')]")
		err = sbutton.Click()
	}
	return err
}

// 点击连接小狐狸 connect-btn clickable text-14-bold
func ConnectMetamask(wd selenium.WebDriver) (err error) {
	var ConnectWallet selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			ConnectWallet, err = wd.FindElement(selenium.ByCSSSelector, ".connect-btn.clickable.text-14-bold")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				err = ConnectWallet.Click()
				if err == nil {
					return true, nil
				}
				time.Sleep(1 * time.Second)
				continue
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err != nil {
		return err
	} else {
		//找到具有setting属性的div
		//wallet-option-item clickable padding-item
		var MeatMask []selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				MeatMask, err = wd.FindElements(selenium.ByCSSSelector, ".wallet-option-item.clickable.padding-item")
				log.Println(len(MeatMask))
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				} else {
					err = MeatMask[1].Click()
					if err == nil {
						return true, nil
					}
				}
			}
			return false, errors.New("Fail")
		}, 5*time.Second)
	}
	return err
}
func MetamaskLogin(wd selenium.WebDriver) (err error) {
	var Password selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			handleNow1, _ := wd.WindowHandles()
			if len(handleNow1) > 1 {
				wd.SwitchWindow(handleNow1[1])
			}
			Password, err = wd.FindElement(selenium.ByID, "password")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, err
	}, 2*time.Second)
	if err == nil {
		Password.SendKeys("SHIfeng0615")
		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
		if err == nil {
			UnLock.Click()
		}
		var SignIn1 selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				SignIn1, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 2*time.Second)
		if err != nil {
			Next, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
			if err == nil {
				Next.Click()
				time.Sleep(1 * time.Second)
				Connect, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
				if err == nil {
					Connect.Click()
				}
				var SignIn selenium.WebElement
				err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < 10; i++ {
						log.Println("寻找signIn")
						mHandle, _ := wd.WindowHandles()
						log.Println(len(mHandle))
						if len(mHandle) > 1 {
							wd.SwitchWindow(mHandle[1])
						}
						SignIn, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
						if err != nil {
							time.Sleep(500 * time.Millisecond)
							continue
						} else {
							return true, nil
						}
					}
					return false, err
				}, 3*time.Second)
				if err == nil {
					SignIn.Click()
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			err = SignIn1.Click()
			if err != nil {
				log.Println("metamask登陆出错")
				return err
			}
		}
	}
	return nil
}
func closeSignInPop(wd selenium.WebDriver, handle string) (err error) {
	wd.SwitchWindow(handle)
	gotIt, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	if len(gotIt) > 1 {
		gotIt[len(gotIt)-1].Click()
	}
	time.Sleep(500 * time.Millisecond)
	SignIn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
	if err == nil {
		err = SignIn.Click()
		if err != nil {
			return err
		}
	}
	return err
}
func CloseSignInPop(wd selenium.WebDriver) (err error) {
	ShandleNow, _ := wd.WindowHandles()
	if len(ShandleNow) > 1 {
		err = closeSignInPop(wd, ShandleNow[1])
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		wd.SwitchWindow(ShandleNow[0])
	}
	return nil
}
func findConnectTwitter(wd selenium.WebDriver) (count int, err error) {
	count = 0
	var SocialLinks []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			SocialLinks, err = wd.FindElements(selenium.ByCSSSelector, ".social-account-link")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err != nil {
		return
	} else {
		//找到具有setting属性的div
		for _, v := range SocialLinks {
			linktext, _ := v.FindElement(selenium.ByCSSSelector, ".social-account-link-text")
			text, _ := linktext.Text()
			if strings.Contains(text, "Connect Twitter Account") {
				return 0, nil
			} else {
				continue
			}
		}
		return 1, nil
	}
	return 1, nil
}
func ClickRemove(wd selenium.WebDriver) (err error) {
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
func ClickConnectTwitter(wd selenium.WebDriver) (err error) {
	var SocialLinks []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			SocialLinks, err = wd.FindElements(selenium.ByCSSSelector, ".social-account-link")
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
		return
	} else {
		//找到具有setting属性的div
		for _, v := range SocialLinks {
			linktext, err := v.FindElement(selenium.ByCSSSelector, ".social-account-link-text")
			if err != nil {
				return err
			}
			text, err := linktext.Text()
			if err != nil {
				return err
			}
			if strings.Contains(text, "Connect Twitter Account") {
				err = v.Click()
			} else {
				continue
			}
		}
	}
	return err
}
func TweetTwitter(wd selenium.WebDriver) (tweetInfo string, err error) {
	//开始发推特
	TweetButton, err := wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
	if err == nil {
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				TweetButton, err = wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
				if len(TweetButton) == 0 {
					time.Sleep(1 * time.Millisecond)
					continue
				} else {
					log.Println(len(TweetButton))
					err = TweetButton[len(TweetButton)-2].Click()
					time.Sleep(1 * time.Second)
					if err != nil {
						handle, _ := wd.WindowHandles()
						if len(handle) == 1 {
							err1 := wd.Refresh()
							log.Println("刷新", err1)
							continue
						}
					}
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 8*time.Second)
		if err != nil {
			log.Println("进入tweet失败", err)
			return "", err
		} else {
			log.Println("进入tweet成功")
		}
		//<span class="css-901oao css-16my406 r-poiln3 r-bcqeeo r-qvutc0">Tweet</span>
		var Tweet []selenium.WebElement
		log.Println("寻找Tweet按钮")
		time.Sleep(1 * time.Second)
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				handle, _ := wd.WindowHandles()
				if len(handle) > 1 {
					wd.SwitchWindow(handle[1])
				}
				Tweet, _ = wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
				if len(Tweet) == 0 {
					time.Sleep(1 * time.Second)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 4*time.Second)
		if err != nil {
			log.Println("没有找到Tweet按钮")
			return "", err
		} else {
			for _, v := range Tweet {
				text, _ := v.Text()
				if strings.Contains(text, "Twitter Circle") {
					_, err = wd.ExecuteScript("window.scrollTo(0, 1000);", nil)
					if err != nil {
						fmt.Printf("无法滚动页面：%v\n", err)
						return
					}
				}
				if strings.Contains(text, "Tweet") && !strings.Contains(text, "Twitter Circle") {
					log.Println("开始推特发文", text)
					err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
						for i := 0; i < 10; i++ {
							time.Sleep(1 * time.Second)
							err = v.Click()
							if err != nil {
								_, _ = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight/document.body.scrollHeight);", nil)
								time.Sleep(1 * time.Second)
								continue
							} else {
								return true, nil
							}
						}
						return false, errors.New("Fail")
					}, 6*time.Second)
					if err == nil {
						time.Sleep(1 * time.Second)
						//
						//
						var failInfo []selenium.WebElement
						_ = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
							for i := 0; i < 10; i++ {
								failInfo, _ = wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
								for _, v1 := range failInfo {
									text1, err := v1.Text()
									if err == nil {
										if strings.Contains(text1, "Something went wrong, but don’t fret — let’s give it another shot.") {
											tweetInfo = text1
										}
									}
								}
							}
							return false, errors.New("f")
						}, 2*time.Second)
						if strings.Contains(tweetInfo, "Something went wrong, but don’t fret — let’s give it another shot.") {
							return tweetInfo, errors.New("fial")
						}
						GotIt, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-1mnahxq.r-19yznuf.r-64el8z.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
						if err == nil {
							err = GotIt.Click()
							if err == nil {
								log.Println("转发推特成功")
								return "", nil
							}
						} else {
							Close, err := wd.FindElement(selenium.ByCSSSelector, ".r-18jsvk2.r-4qtqp9.r-yyyyoo.r-z80fyv.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-19wmn03")
							if err == nil {
								err = Close.Click()
								if err == nil {
									log.Println("转发推特成功-----你已经发过了")
									return "", nil
								}
							}
						}
					} else {
						log.Println("tweet的时候出错了")
					}
				}
				if strings.Contains(text, "Post") && !strings.Contains(text, "Twitter Circle") {
					log.Println("开始推特发文", text)
					err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
						for i := 0; i < 10; i++ {
							time.Sleep(1 * time.Second)
							err = v.Click()
							if err != nil {
								_, _ = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight/document.body.scrollHeight);", nil)
								time.Sleep(1 * time.Second)
								continue
							} else {
								return true, nil
							}
						}
						return false, errors.New("Fail")
					}, 6*time.Second)
					if err == nil {
						time.Sleep(1 * time.Second)
						//
						var failInfo []selenium.WebElement
						_ = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
							for i := 0; i < 10; i++ {
								failInfo, _ = wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
								for _, v1 := range failInfo {
									text1, err := v1.Text()
									if err == nil {
										if strings.Contains(text1, "Something went wrong, but don’t fret — let’s give it another shot.") {
											tweetInfo = text1
											return true, nil
										}
									}
								}
							}
							return false, errors.New("f")
						}, 2*time.Second)
						if strings.Contains(tweetInfo, "Something went wrong, but don’t fret — let’s give it another shot.") {
							return tweetInfo, errors.New("fial")
						}
						GotIt, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-1mnahxq.r-19yznuf.r-64el8z.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
						if err == nil {
							err = GotIt.Click()
							if err == nil {
								log.Println("转发推特成功")
								return "", nil
							}
						} else {
							Close, err := wd.FindElement(selenium.ByCSSSelector, ".r-18jsvk2.r-4qtqp9.r-yyyyoo.r-z80fyv.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-19wmn03")
							if err == nil {
								err = Close.Click()
								if err == nil {
									log.Println("转发推特成功-----你已经发过了")
									return "", nil
								}
							}
						}
					} else {
						log.Println("tweet的时候出错了")
					}
				}

			}
		}
	} else {
		log.Println("点击Tweet1失败")
		return "", err
	}
	return "", err
}

// fire的电脑配置
func GetInProfile(wd selenium.WebDriver) (userName string, err error) {
	//点击进入profile
	var Profile selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			Profile, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"react-root\"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[9]")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err != nil {
		log.Println("进入Profile失败")
	} else {
		err = Profile.Click()
		if err != nil {
			log.Println("进入Profile失败")
		}
		//获取推特链接
		time.Sleep(2 * time.Second)
		var articles []selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				articles, _ = wd.FindElements(selenium.ByTagName, "article")
				if len(articles) == 0 {
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 2*time.Second)
		url, _ := wd.CurrentURL()
		log.Println("当前url", url)
		urlParts := strings.Split(url, "/")
		username := urlParts[len(urlParts)-1]
		log.Println("查找到的article数量:", len(username))
		//随便follow几个
		//css-18t94o4 css-1dbjc4n r-1ny4l3l r-ymttw5 r-1f1sjgu r-o7ynqc r-6416eg
		if len(articles) > 0 {
			var LinkButtons1 selenium.WebElement
			err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
				for i := 0; i < 10; i++ {
					LinkButtons1, err = articles[0].FindElement(selenium.ByCSSSelector, ".css-1dbjc4n.r-1kbdv8c.r-18u37iz.r-1wtj0ep.r-1s2bzr4.r-hzcoqn")
					if err != nil {
						time.Sleep(500 * time.Millisecond)
						continue
					} else {
						return true, nil
					}
				}
				return false, errors.New("Fail")
			}, 5*time.Second)
			if err != nil {
				return "", err
			}
			LinkButtons, _ := LinkButtons1.FindElements(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-1777fci.r-bt1l66.r-1ny4l3l.r-bztko3.r-lrvibr")
			log.Println("查找到的LinkButtons数量:", len(LinkButtons))
			time.Sleep(1 * time.Second)
			if len(LinkButtons) > 0 {
				err = LinkButtons[len(LinkButtons)-1].Click()
				if err != nil {
					log.Println("点击获取link链接按钮失败")
				} else {
					//这里会出现四个弹窗,点击第一个
					time.Sleep(1 * time.Second)
					Links, _ := wd.FindElements(selenium.ByCSSSelector, ".css-901oao.r-18jsvk2.r-37j5jr.r-a023e6.r-b88u0q.r-rjixqe.r-bcqeeo.r-qvutc0")
					if len(Links) > 0 {
						err = Links[0].Click()
						if err != nil {
							log.Println("获取链接失败")
						} else {
							//follows, _ := wd.FindElements(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-1ny4l3l.r-ymttw5.r-1f1sjgu.r-o7ynqc.r-6416eg")
							//log.Println(len(follows))
							//for k, v := range follows {
							//	if k >= len(follows)-2 {
							//		follow, err := v.FindElement(selenium.ByCSSSelector, ".css-901oao.css-16my406.css-1hf3ou5.r-poiln3.r-1b43r9.r-1cwl3u0.r-bcqeeo.r-qvutc0")
							//		if err == nil {
							//			follow.Click()
							//		}
							//	}
							//}
							return username, nil
						}
					}
				}
			}
		} else {
			return username, errors.New("fail")
		}
	}

	return "", err
}
func Sendkey(wd selenium.WebDriver) (detail string, err error) {
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		log.Println("读取剪贴板内容失败：", err)
	} else {
		log.Println(clipboardContent)
		Input, err := wd.FindElement(selenium.ByTagName, "input")
		if err != nil {
			return "", err
		} else {
			err = Input.SendKeys(clipboardContent)
			if err != nil {
				log.Println("粘贴链接失败", err)
				return "", err
			} else {
				time.Sleep(1 * time.Second)
				Verifys, _ := wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
				if len(Verifys) > 0 {
					log.Println("点击verify")
					Verifys[len(Verifys)-1].Click()
				}
			}
		}
	}
	time.Sleep(2 * time.Second)
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		var isSuccess []selenium.WebElement
		for i := 0; i < 10; i++ {
			isSuccess, _ = wd.FindElements(selenium.ByCSSSelector, ".text-lighten-1")
			if len(isSuccess) == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				text, err := isSuccess[len(isSuccess)-1].Text()
				//if strings.Contains(text, "Twitter account verified successfully") {
				//	log.Println("Twitter account verified successfully")
				//	detail = text
				//	return true, nil
				//}
				if err == nil && text != "" {
					detail = text
					return true, nil
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
		return false, errors.New("Fail")
	}, 6*time.Second)
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
		err = button.Click()
	}
	return
}
func findTwitterUsername(wd selenium.WebDriver) (username string, err error) {
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
		return "", err
	} else {
		//找到twitter在的div,并点击remove
		if len(removeTwitter) > 0 {
			for _, v := range removeTwitter {
				SaName, _ := v.FindElement(selenium.ByCSSSelector, ".sa-name")
				text, _ := SaName.Text()
				if strings.Contains(text, "Twitter") {
					actions, _ := v.FindElement(selenium.ByCSSSelector, ".sa-info-id")
					text1, err := actions.Text()
					if err != nil {
						return "", err
					}
					err = actions.Click()
					if err != nil {
						log.Println("点击进去推特失败")
						return "", err
					}
					return text1, nil
				} else {
					continue
				}
			}
		} else {
			log.Println("没有找到对应的box")
			return "", err
		}
	}
	return
}
