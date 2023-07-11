/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-02-22 09:41:15
 * @LastEditTime: 2023-02-22 09:41:23
 */
package twitter

import (
	"fmt"
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
