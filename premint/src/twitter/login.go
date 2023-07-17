/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-13 15:16:05
 * @LastEditTime: 2023-03-17 14:50:55
 */
package twitter

import (
	"log"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/tebeka/selenium"
)

// 登录twitter
func LoginTwitter(userName string, pwd string, userAt string, wd selenium.WebDriver) {
	tabScript := `window.open();`
	if _, err := wd.ExecuteScript(tabScript, nil); err != nil {
		log.Println(err)
	}
	//选择第二个窗口
	util.GetWind(1, wd)
	wd.Get("https://twitter.com/")
	loginClick(wd)
	inputUserName(userName, wd)
	inputUserAt(userAt, wd)
	inputPwd(pwd, wd)
	time.Sleep(10 * time.Second)
	wd.Close()
	util.GetWind(0, wd)

}
func loginClick(wd selenium.WebDriver) {
	log.Println("*******twitter点击登录*******")
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		loginBtnParent, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-l5o3uw.r-1upvrn0")
		if err != nil {
			log.Println(err)
			return false, err
		}
		if len(loginBtnParent) <= 0 {
			log.Println("未找到twitter登录父节点")
			wd.Refresh()
			time.Sleep(10 * time.Second)
			return false, err
		}
		loginBtn, err := loginBtnParent[0].FindElements(selenium.ByCSSSelector, ".css-4rbku5.css-18t94o4.css-1dbjc4n.r-1niwhzg.r-sdzlij.r-1phboty.r-rs99b7.r-1loqt21.r-2yi16.r-1qi8awa.r-1ny4l3l.r-ymttw5.r-o7ynqc.r-6416eg.r-lrvibr")
		if err != nil {
			log.Println(err)
		}
		if len(loginBtn) <= 0 {
			log.Println("未找到twitter登录按钮")
			wd.Refresh()
			time.Sleep(5 * time.Second)
			return false, err
		}
		loginBtn[0].Click()
		return true, err
	})
}

func inputUserName(userName string, wd selenium.WebDriver) {
	log.Println("*****输入twitter账户")
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		inputNode, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-mk0yit.r-1f1sjgu.r-13qz1uu")
		if len(inputNode) <= 0 {
			log.Println("未找到输入框")
			return false, err
		}
		input, err := inputNode[0].FindElement(selenium.ByTagName, "input")
		// 输入并回车
		if err := input.SendKeys(userName); err != nil {
			log.Println(err)
		}
		if err := input.SendKeys(selenium.EnterKey); err != nil {
			log.Println(err)
		}
		return true, err
	})
}

func inputUserAt(userAt string, wd selenium.WebDriver) {
	log.Println("********输入用户名********")
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		inputNode, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-mk0yit.r-1f1sjgu")
		if len(inputNode) <= 0 {
			log.Println("未找到输入框")
			return false, err
		}
		input, err := inputNode[0].FindElement(selenium.ByTagName, "input")
		// 输入并回车
		if err := input.SendKeys("@" + userAt); err != nil {
			log.Println(err)
		}
		if err := input.SendKeys(selenium.EnterKey); err != nil {
			log.Println(err)
		}
		return true, err
	})

}

func inputPwd(pwd string, wd selenium.WebDriver) {
	log.Println("*******输入密码******")
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		inputNode, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-1wzrnnt.r-1udh08x.r-xd6kpl.r-1pn2ns4.r-ttdzmv")
		if len(inputNode) <= 0 {
			log.Println("未找到输入框")
			return false, err
		}
		input, err := inputNode[1].FindElement(selenium.ByTagName, "input")
		// 输入并回车
		if err := input.SendKeys(pwd); err != nil {
			log.Println(err)
		}
		if err := input.SendKeys(selenium.EnterKey); err != nil {
			log.Println(err)
		}
		return true, err
	})

}
