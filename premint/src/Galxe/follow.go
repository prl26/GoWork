package Galxe

import (
	"errors"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

func GalxeFollow(wd selenium.WebDriver, url selenium.WebElement) error {
	url.Click()
	time.Sleep(2 * time.Second)
	//切换到新打开的页面
	handles1, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles1)-1])
	CurrentHandle2, _ := wd.CurrentWindowHandle()

	log.Println("打开detail第二次的handle长度", len(handles1))
	log.Println("当前handle", CurrentHandle2)
	//关闭弹窗
	pop, _ := wd.FindElements(selenium.ByCSSSelector, ".iconfont.icon-close.text-20-regular.clickable.m-popup-close-icon")
	if len(pop) > 0 {
		for _, v := range pop {
			v.Click()
		}
	}
	time.Sleep(1 * time.Second)
	//寻找follow元素并点击
	err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			_, err := wd.FindElement(selenium.ByCSSSelector, ".spine-player-canvas")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("失败")
	}, 15*time.Second)
	if err != nil {
		log.Println("查找follow按钮出错了")
		return err
	} else {
		followButton, _ := wd.FindElement(selenium.ByCSSSelector, ".spine-player-canvas")
		followButton.Click()
		time.Sleep(1 * time.Second)

	}
	return err
}
