/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-21 14:19:00
 * @LastEditTime: 2023-04-23 13:54:09
 */
package metamask

import (
	"fmt"
	"log"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/tebeka/selenium"
)

func MetaMaskLogin(wd selenium.WebDriver, pwd string) {
	fmt.Println("*********插件登录********")
	wd.Get("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html#unlock")
	util.GetCurrentWindow(wd)
	count := 0
	//选择登录框
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(2 * time.Second)
		input, err := wd.FindElement(selenium.ByTagName, "input")
		time.Sleep(2 * time.Second)
		//we, err := wd.FindElements(selenium.ByXPATH, "//*[contains(@class,'r-4qtqp9 r-yyyyoo r-50lct3 r-dnmrzs r-bnwqim r-1plcrui r-lrvibr r-1srniue')]/..")
		if err != nil {
			log.Println("输入框查找失败", err)
			count += 1
			if count > 3 {
				return true, err
			}
		}
		if input == nil {
			return false, nil
		}
		// 输入并回车
		if err := input.SendKeys(pwd); err != nil {
			log.Println(err)
		}
		if err := input.SendKeys(selenium.EnterKey); err != nil {
			log.Println(err)
		}

		return true, err

	})
	//检测弹窗
	closePopup(wd)

}

func closePopup(wd selenium.WebDriver) {
	time.Sleep(3 * time.Second)
	log.Println("检测弹窗")
	pop, err := wd.FindElement(selenium.ByCSSSelector, ".popover-container")
	if err != nil {
		log.Println("未找到弹窗", err)
	} else {
		closeIcons, _ := pop.FindElements(selenium.ByCSSSelector, ".box.mm-button-icon.mm-button-icon--size-sm.box--display-inline-flex.box--flex-direction-row.box--justify-content-center.box--align-items-center.box--color-icon-default.box--background-color-transparent.box--rounded-lg")
		if len(closeIcons) <= 0 {
			log.Println("未找到关闭按钮")
		} else {
			closeIcons[0].Click()
		}
	}

	networkPop, err := wd.FindElement(selenium.ByCSSSelector, "popover-container")
	if err != nil {
		log.Println("未找到弹窗2", err)
	} else {
		knows, _ := networkPop.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
		if len(knows) > 0 {
			for _, know := range knows {
				know.Click()
			}
		} else {
			log.Println("未找到明白了按钮")
		}
	}

}
func closePopup1(wd selenium.WebDriver) {
	time.Sleep(3 * time.Second)
	log.Println("检测弹窗")
	pop, err := wd.FindElement(selenium.ByCSSSelector, ".popover-container")
	if err != nil {
		log.Println("未找到弹窗", err)
	} else {
		closeIcons, _ := pop.FindElements(selenium.ByCSSSelector, ".box.mm-button-icon.mm-button-icon--size-sm.box--display-inline-flex.box--flex-direction-row.box--justify-content-center.box--align-items-center.box--color-icon-default.box--background-color-transparent.box--rounded-lg")
		if len(closeIcons) <= 0 {
			log.Println("未找到关闭按钮")
		} else {
			closeIcons[0].Click()
		}
	}

	networkPop, err := wd.FindElement(selenium.ByCSSSelector, "popover-container")
	if err != nil {
		log.Println("未找到弹窗2", err)
	} else {
		knows, _ := networkPop.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
		if len(knows) > 0 {
			for _, know := range knows {
				know.Click()
			}
		} else {
			log.Println("未找到明白了按钮")
		}
	}

}
