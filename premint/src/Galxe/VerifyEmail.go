package Galxe

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
	"regexp"
	"strings"
	"time"
)

//1.打开银河
//2.打开个人首页

func VerifyEmail() {
	url := "https://galxe.com/credentials"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\Verify-7.28.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\Verify-7.28.txt"
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
				err := verifyEmail(v, k, chs, wd, url, dstFile)
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
	} else {
		log.Println("excel 保存成功----")
	}

}
func verifyEmail(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	url1 := "https://mail.rambler.ru/folder/INBOX"
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
	//if err != nil {
	//	log.Println("窗口方法失败", err)
	//}
	err = openHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString(fmt.Sprintf("进入个人主页失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
	}
	time.Sleep(2 * time.Second)

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
	change, err := wd.FindElements(selenium.ByCSSSelector, ".setting-email-change-action")
	log.Println(len(change))
	if len(change) > 0 {
		text, _ := change[len(change)-1].Text()
		log.Println(text)
		if strings.Contains(text, "Change") {
			log.Println("你已经绑定过了")
			return nil
		}
	}
	CloseSignInPopNoMax(wd)
	//打开邮箱界面
	wd.SwitchWindow(handleNow[0])
	_, err = wd.ExecuteScript("window.open(arguments[0], '_blank');", []interface{}{url1})
	//err = wd.Get(url1)
	if err != nil {
		log.Println(excelInfo.BitId, "打开邮箱界面出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("邮箱界面成功")
	}
	handow, _ := wd.WindowHandles()
	err = wd.SwitchWindow(handow[len(handow)-1])
	if err != nil {
		log.Println("切换窗口权柄失败")
		return err
	}
	time.Sleep(2 * time.Second)
	//获取邮箱地址
	clipboard.WriteAll("")

	err = copyEmailAddress(wd)
	if err != nil {
		log.Println("获取邮箱地址失败")
		dstFile.WriteString(fmt.Sprintf("获取邮箱地址失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	}
	clipboardContent, err := clipboard.ReadAll()
	log.Println("获取邮箱地址成功:", clipboardContent)

	time.Sleep(2 * time.Second)

	//email-verify-code-input

	err = wd.SwitchWindow(handow[0])
	if err != nil {
		log.Println("切换窗口权柄失败")
		return err
	}
	err = sendAndClickInputs(wd, clipboardContent)
	if err != nil {
		log.Println("发送邮箱验证码失败")
		dstFile.WriteString(fmt.Sprintf("发送邮箱验证码失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("发送邮箱验证码成功")
	}
	//获取到验证码
	time.Sleep(10 * time.Second)

	err = wd.SwitchWindow(handow[len(handow)-1])
	if err != nil {
		log.Println("切换窗口权柄失败")
		return err
	}
	time.Sleep(5 * time.Second)
	emailCode, err := getVerifyCode(wd)
	log.Println(emailCode)
	if err != nil {
		log.Println(excelInfo.BitId, "获取验证码出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("获取验证码出错了-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("获取验证码成功")
		log.Println("获取到的验证码:", emailCode)
		time.Sleep(2 * time.Second)
		err := wd.Close()
		if err != nil {
			log.Println("关闭邮箱失败")
			return err
		}
		err = wd.SwitchWindow(handow[0])
		if err != nil {
			log.Println("切换原窗口权柄失败")
			return err
		}
	}
	//填写code并进行Veify
	time.Sleep(2 * time.Second)
	err = sendAndClickVerify(wd, emailCode)
	if err != nil {
		log.Println("填写code并进行Veify失败")
		dstFile.WriteString(fmt.Sprintf("填写code并进行Veify失败-----%v", excelInfo.BitId))
		ch <- wrongData
		return err
	} else {
		log.Println("填写code并进行Veify成功")
	}
	time.Sleep(5 * time.Second)
	change1, err := wd.FindElements(selenium.ByCSSSelector, ".setting-email-change-action")
	if len(change1) > 0 {
		text1, _ := change1[len(change1)-1].Text()
		if !strings.Contains(text1, "Change") {
			dstFile.WriteString(fmt.Sprintf("进行Veify失败-----%v", excelInfo.BitId))
			log.Println("绑定出错")
			return err
		}
	} else {
		log.Println("绑定出错")
		dstFile.WriteString(fmt.Sprintf("进行Veify失败-----%v", excelInfo.BitId))
		return errors.New("fail")
	}
	return nil

}
func sendAndClickVerify(wd selenium.WebDriver, emailCode string) (err error) {
	var inputs []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			inputs, _ = wd.FindElements(selenium.ByCSSSelector, ".email-verify-code-input")
			if len(inputs) > 0 {
				return true, nil
			} else {
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err == nil {
		err = inputs[len(inputs)-1].SendKeys(emailCode)
		if err == nil {
			sendCode, _ := wd.FindElements(selenium.ByCSSSelector, ".email-verify-btn")
			err = sendCode[len(sendCode)-1].Click()
			if err == nil {
				return nil
			}
		} else {
			log.Println("点击Verify失败")
			return err
		}
	}
	return err
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
func sendAndClickInputs(wd selenium.WebDriver, address string) (err error) {
	var inputs []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			inputs, _ = wd.FindElements(selenium.ByCSSSelector, ".email-verify-code-input")
			if len(inputs) > 0 {
				return true, nil
			} else {
				time.Sleep(500 * time.Millisecond)
				continue
			}
		}
		return false, errors.New("Fial")
	}, 5*time.Second)
	if err == nil {
		err = inputs[0].SendKeys(address)
		if err == nil {
			sendCode, _ := wd.FindElements(selenium.ByCSSSelector, ".email-send-verify-code-btn")
			err = sendCode[len(sendCode)-1].Click()
			if err == nil {
				return nil
			}
		} else {
			log.Println("邮箱填写失败")
			return err
		}
	}
	return err
}
func getVerifyCode(wd selenium.WebDriver) (emailCode string, err error) {
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			div, err := wd.FindElement(selenium.ByCSSSelector, ".MailList-list-2L")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			Code, _ := div.FindElements(selenium.ByTagName, "a")
			if len(Code) == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			span, _ := Code[0].FindElement(selenium.ByCSSSelector, ".ListItem-snippet-1a")
			text, err2 := span.Text()
			if err2 != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			emailCode = text[:6]
			isOk := isNumeric(emailCode)
			if !isOk {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return true, nil
		}
		return false, errors.New("fail")
	}, 20*time.Second)
	return
}
func copyEmailAddress(wd selenium.WebDriver) (err error) {
	var profile []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			profile, _ = wd.FindElements(selenium.ByCSSSelector, ".rc__528vB")
			if len(profile) > 0 {
				err = profile[len(profile)-1].Click()
				if err != nil {
					time.Sleep(500 * time.Millisecond)
					continue
				}
				return true, nil
			} else {
				time.Sleep(500 * time.Millisecond)
				continue
			}

		}
		return false, errors.New("fail")
	}, 5*time.Second)
	if err == nil {
		var emailAddress []selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				emailAddress, _ = wd.FindElements(selenium.ByCSSSelector, ".rc__AXUNk")
				if len(emailAddress) > 0 {
					err = emailAddress[len(profile)-1].Click()
					if err != nil {
						time.Sleep(500 * time.Millisecond)
						continue
					}
					return true, nil
				} else {
					time.Sleep(500 * time.Millisecond)
					continue
				}

			}
			return false, errors.New("fail")
		}, 5*time.Second)
	}
	err = profile[len(profile)-1].Click()
	return err
}
func isNumeric(str string) bool {
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	return numericRegex.MatchString(str)
}
func CloseSignInPopNoMax(wd selenium.WebDriver) {
	ShandleNow, _ := wd.WindowHandles()
	if len(ShandleNow) > 1 {
		closeSignInPop(wd, ShandleNow[1])
		time.Sleep(2 * time.Second)
	}
}
