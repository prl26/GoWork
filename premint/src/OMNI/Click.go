package main

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/tebeka/selenium"
	"log"
	"strings"
	"time"
)

func OmniClick(wd selenium.WebDriver, pwd string) {
	//进入领水页面
	fmt.Println("*********进入领水页面********")
	wd.Get("https://faucet.omni.network/")
	util.GetCurrentWindow(wd)

	button, err := wd.FindElement(selenium.ByClassName, "Button_root__3vtfo.FaucetForm_btn__l_QG1")
	if err != nil {
		log.Println("Claim查找失败", err)
	}
	input_element, err := wd.FindElement(selenium.ByXPATH, "/html/body/main/form/label/input")
	if err != nil {
		log.Println("输入框查找失败", err)
	}
	input_element.SendKeys(pwd)
	button.Click()

	//弹窗捕捉
	err = CaptureCongrats(wd, pwd)
	if err != nil {

	}

}
func CaptureCongrats(wd selenium.WebDriver, value string) error {
	err := wd.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByClassName, "FaucetForm_header__HZ85J")
		if err != nil {
			//没有获取到成功
			element1, err := wd.FindElement(selenium.ByClassName, "Notification_message__TODVe")
			if err != nil {
				return false, err
			} else {
				if text, _ := element1.Text(); strings.Contains(text, "Address Limit: you've already claimed funds today") {
					log.Println(value, "warning的值：", text)
					time.Sleep(1)
					return true, nil
				}
				return false, err
			}
		} else {
			log.Println(fmt.Printf("%v--领取成功", value))
			time.Sleep(1)
			return true, nil
		}
		return false, err
	}, 25, 3)
	return err
}
