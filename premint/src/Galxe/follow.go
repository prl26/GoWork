package Galxe

import (
	"errors"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

func GalxeFollow(wd selenium.WebDriver) (err error) {
	time.Sleep(1 * time.Second)
	main_handle, _ := wd.WindowHandles()
	if len(main_handle) > 3 {
		LoginRequest(wd, main_handle)
	}
	//切换到新打开的页面
	time.Sleep(1 * time.Second)
	handles1, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles1)-1])
	//关闭弹窗
	time.Sleep(2 * time.Second)

	pop, _ := wd.FindElements(selenium.ByCSSSelector, ".iconfont.icon-close.text-20-regular.clickable.m-popup-close-icon")
	log.Println("pop长度：", len(pop))
	if len(pop) > 0 {
		for _, v := range pop {
			v.Click()
		}
	}
	//寻找follow元素并点击
	time.Sleep(2 * time.Second)
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			_, err = wd.FindElement(selenium.ByCSSSelector, ".spine-player-canvas")
			if err != nil {
				ClosePop(wd)
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("失败")
	}, 10*time.Second)
	if err != nil {
		log.Println("查找follow按钮出错了")
		return err
	} else {
		followButton, _ := wd.FindElement(selenium.ByCSSSelector, ".spine-player-canvas")
		err = followButton.Click()
		if err != nil {
			err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
				for i := 0; i < 10; i++ {
					_, err = wd.FindElement(selenium.ByCSSSelector, ".dialog-classic-content")
					if err != nil {
						time.Sleep(1 * time.Second)
						continue
					} else {
						return true, nil
					}
				}
				return false, errors.New("失败")
			}, 5*time.Second)
			if err != nil {
				log.Println("Galxe Task 失败")
				return err
			} else {
				log.Println("----------Galxe follow成功-------------")
			}
		} else {
			log.Println("----------Galxe follow成功-------------")
		}
	}
	return err
}
func LoginRequest(wd selenium.WebDriver, main_handle []string) {
	wd.SwitchWindow(main_handle[3])
	time.Sleep(1 * time.Second)
	knows, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	if len(knows) > 0 {
		for _, know := range knows {
			know.Click()
		}
		LoginButon, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		log.Println("处理第一次小狐狸登陆,button长度：", len(LoginButon))
		if len(LoginButon) > 0 {
			for _, v := range LoginButon {
				v.Click()
			}
		}
	}
}
func ClosePop(wd selenium.WebDriver) {
	pop, _ := wd.FindElements(selenium.ByCSSSelector, ".iconfont.icon-close.text-20-regular.clickable.m-popup-close-icon")
	log.Println("pop长度：", len(pop))
	if len(pop) > 0 {
		for _, v := range pop {
			v.Click()
		}
	}
}
