/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-02-22 09:41:15
 * @LastEditTime: 2023-02-22 09:41:23
 */
package twitter

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

const (
	driverPath = "./driver/chromedriver.exe"
	port       = 9515
)

func BeLiek(link string) string {
	fmt.Println("测试测试", link)
	//link = "https://twitter.com/brgridiron/status/1624948323416436736"
	//1.开启selenium服务
	//设置selium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(driverPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//延迟关闭服务
	defer service.Stop()
	//2.调用浏览器
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	//wd.MaximizeWindow("")
	//延迟退出chrome
	//defer wd.Quit()

	//3.对页面元素进行操作
	//获取link的页面
	if err := wd.Get(link); err != nil {
		panic(err)
	}
	//等待页面加载完成
	wd.SetImplicitWaitTimeout(30 * time.Second)
	//找到点赞的div
	we, err := wd.FindElements(selenium.ByCSSSelector, ".r-4qtqp9.r-yyyyoo.r-50lct3.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-1srniue")
	//we, err := wd.FindElements(selenium.ByXPATH, "//*[contains(@class,'r-4qtqp9 r-yyyyoo r-50lct3 r-dnmrzs r-bnwqim r-1plcrui r-lrvibr r-1srniue')]/..")
	if err != nil {
		panic(err)
	}

	if len(we) <= 0 {
		fmt.Println("未获取到元素")
		return ""
	}

	parentLike, err := we[2].FindElement(selenium.ByXPATH, "./..")
	if err != nil {
		fmt.Println("查找父节点失败", err)
	}
	fmt.Println(parentLike.Text())
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		//点击
		err = parentLike.Click()
		return true, err
	}, 30*time.Second)
	// //点击
	// err = parentLike.Click()
	// if err != nil {
	// 	fmt.Println("错误：", err)
	// 	//panic(err)
	// }

	//睡眠20秒后退出
	time.Sleep(20 * time.Minute)

	return ""
}
func TwitterTweet(wd selenium.WebDriver) error {

	//切换到新打开的页面
	handles1, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles1)-1])
	CurrentHandle2, _ := wd.CurrentWindowHandle()

	log.Println("打开detail第二次的handle长度", len(handles1))
	log.Println("当前handle", CurrentHandle2)

	err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			//css-901oao r-1awozwy r-jwli3a r-6koalj r-18u37iz r-16y2uox r-37j5jr r-a023e6 r-b88u0q r-1777fci r-rjixqe r-bcqeeo r-q4m81j r-qvutc0
			//_, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-16y2uox.r-6gpygo.r-peo1c.r-1ps3wis.r-1ny4l3l.r-1udh08x.r-1guathk.r-1udbk01.r-o7ynqc.r-6416eg.r-lrvibr.r-3s2u2q")
			handles1, _ := wd.WindowHandles()
			wd.SwitchWindow(handles1[len(handles1)-1])
			_, err := wd.FindElement(selenium.ByCSSSelector, ".css-901oao.r-1awozwy.r-jwli3a.r-6koalj.r-18u37iz.r-16y2uox.r-37j5jr.r-a023e6.r-b88u0q.r-1777fci.r-rjixqe.r-bcqeeo.r-q4m81j.r-qvutc0")

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
		log.Println("查找follow失败")
		return err
	} else {
		handleNow, _ := wd.CurrentWindowHandle()
		wd.MaximizeWindow(handleNow)
		button, _ := wd.FindElement(selenium.ByCSSSelector, ".css-901oao.r-1awozwy.r-jwli3a.r-6koalj.r-18u37iz.r-16y2uox.r-37j5jr.r-a023e6.r-b88u0q.r-1777fci.r-rjixqe.r-bcqeeo.r-q4m81j.r-qvutc0")
		err := button.Click()
		if err != nil {
			log.Println("twitter tweet 点击失败")
			return err
		}
	}
	return err
}
func TwitterLike(wd selenium.WebDriver) error {
	//切换到新打开的页面
	handles1, _ := wd.WindowHandles()
	wd.SwitchWindow(handles1[len(handles1)-1])
	err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			handles1, _ := wd.WindowHandles()
			wd.SwitchWindow(handles1[len(handles1)-1])
			_, err := wd.FindElement(selenium.ByCSSSelector, ".css-901oao.r-1awozwy.r-jwli3a.r-6koalj.r-18u37iz.r-16y2uox.r-37j5jr.r-a023e6.r-b88u0q.r-1777fci.r-rjixqe.r-bcqeeo.r-q4m81j.r-qvutc0")
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
		log.Println("查找Like失败")
		return err
	} else {
		button, _ := wd.FindElement(selenium.ByCSSSelector, ".css-901oao.r-1awozwy.r-jwli3a.r-6koalj.r-18u37iz.r-16y2uox.r-37j5jr.r-a023e6.r-b88u0q.r-1777fci.r-rjixqe.r-bcqeeo.r-q4m81j.r-qvutc0")
		err := button.Click()
		if err != nil {
			log.Println("twitter Like 点击失败")
			return err
		} else {
			log.Println("twitter Like 点击成功")
		}
	}
	return err
}
