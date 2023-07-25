package Galxe

import (
	"errors"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

func GalxeVisit(wd selenium.WebDriver) (err error) {
	//切换到新打开的页面
	time.Sleep(1 * time.Second)
	handles1, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles1)-1])
	//继续 access
	var Access selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			Access, err = wd.FindElement(selenium.ByCSSSelector, ".button-special")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("失败")
	}, 10*time.Second)
	if err != nil {
		log.Println("跳转页面失败")
		return err
	}
	Access.Click()
	//寻找follow元素并点击
	handles2, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles2)-1])
	time.Sleep(2 * time.Second)
	return err
}
