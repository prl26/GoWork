/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-02-22 14:14:43
 * @LastEditTime: 2023-06-06 11:36:51
 */
package premint

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/metamask"
	"github.com/JianLinWei1/premint-selenium/src/twitter"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
)

func Start() {
	excelInfos := util.GetExcelInfos("./Premint脚本点击模版.xlsx")
	chs := make(chan string, len(excelInfos))

	for j, bro := range excelInfos {
		time.Sleep(1 * time.Second)
		wd, _ := wdservice.InitWd(j, bro.BitId)
		if wd != nil {
			go util.SetLog(func() {
				ToRegister(excelInfos[j], j, chs, wd)
			})
		} else {
			log.Println("第" + strconv.Itoa(j+1) + "条数据浏览器初始化失败****")
		}

	}
	bitbrowser.Windowbounds()
	for t := 0; t < len(excelInfos); t++ {
		log.Println(<-chs)
	}
}

func Test() {
	excelInfos := util.GetExcelInfos("./Premint脚本点击模版.xlsx")
	chs := make(chan string, len(excelInfos))
	for j, _ := range excelInfos {
		time.Sleep(1 * time.Second)
		wd, _ := wdservice.InitWd(j, "f7f6a055803548f493b3b45ab5e062db")
		go ToRegister(excelInfos[j], j, chs, wd)
	}
	for t := 0; t < len(excelInfos); t++ {
		log.Println(<-chs)
	}
}

func ToRegister(excelInfo util.ExcelInfo, i int, ch chan<- string, wd selenium.WebDriver) {
	log.Println("*********************开始处理第" + strconv.Itoa(i+1) + "条数据******************")
	metamask.MetaMaskLogin(wd, excelInfo.MetaPwd)
	time.Sleep(1 * time.Second)
	if err := wd.Get(excelInfo.PremintUrl); err != nil {
		log.Println(err)
	}
	log.Println("********已打开浏览器********")
	log.Println(wd)
	wd.SetImplicitWaitTimeout(3 * time.Second)
	isLogin := loginToRegisterClick(wd)
	log.Println("*********等待10秒**********")
	if !isLogin {
		time.Sleep(1 * time.Second)
		clickMetaMask(wd)
		PostSign(wd, excelInfo.PremintUrl)
	}
	time.Sleep(3 * time.Second)
	//util.GetWind(0, wd)
	twitter.FollowAndTweeTwitter(wd)
	clickToRegister(wd)
	ch <- "第" + strconv.Itoa(i+1) + "条数据执行结束"
}

// 点击去登录
func loginToRegisterClick(wd selenium.WebDriver) bool {
	log.Println("*********点击loginToRegister********")
	isLogin := false
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		loginBtn, err := wd.FindElements(selenium.ByCSSSelector, ".btn.btn-styled.btn-success.btn-shadow.btn-xl.btn-block.mt-3")
		if err != nil {
			log.Println(err)
		}
		if len(loginBtn) <= 0 {
			log.Println("没有找到登录按钮")
			isLogin = true
			return false, err
		}
		loginBtn[0].Click()
		return true, err
	}, 10*time.Second)

	return isLogin
}

// 点击metamask授权
func clickMetaMask(wd selenium.WebDriver) {
	log.Println("*********点击metamask授权********")
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		metaMaskBtn, err := wd.FindElements(selenium.ByCSSSelector, ".btn.btn-styled.btn-block.bg-muted.border.btn-circle.mb-3.d-flex.align-items-center.justify-content-center")
		if err != nil {
			log.Println(err)
		}
		if len(metaMaskBtn) <= 0 {
			log.Println("没有找到metamask按钮")
			return false, err
		}
		metaMaskBtn[0].Click()
		return true, err
	}, 30*time.Second)
}

// 去签名
func PostSign(wd selenium.WebDriver, url string) {
	log.Println("*********去签名********")
	//util.GetWind(1, wd)
	time.Sleep(3 * time.Second)
	handle, _ := wd.CurrentWindowHandle()
	util.GetWindByName("MetaMask Notification", wd)

	//wd.Refresh()
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		log.Println(wd.Title())
		signBtn, err := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.btn--large.request-signature__footer__sign-button")
		if err != nil {
			log.Println(err)
		}
		if len(signBtn) <= 0 {
			log.Println("没有找到签名按钮>>>>去点击下一步")
			wd.Get("chrome-extension://nkbihfbeogaeaoehlefnkodbefgpgknn/home.html")

			flag := signNext(wd)
			if flag {
				return true, err
			}
			return false, err
		}
		signBtn[0].Click()
		return true, err
	}, 30*time.Second)

	wd.Close()
	//util.GetWind(0, wd)
	wd.SwitchWindow(handle)
	wd.Get(url)

}

// 签名下一步
func signNext(wd selenium.WebDriver) bool {
	nextBtn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	if err != nil {
		log.Println("未找到下一步按钮")
		return false
	}
	nextBtn.Click()
	conBtn, _ := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
	if conBtn != nil {
		conBtn.Click()
	}

	return true
}

func clickToRegister(wd selenium.WebDriver) {
	log.Println("*********最后的注册********")
	var continued string
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		btn, err := wd.FindElement(selenium.ByID, "register-submit")
		if err != nil {
			log.Println("*****1注册按钮不可点击****可能需要DC验证请手动操作*****")

			fmt.Println("输入y继续：")
			fmt.Scanln(&continued)
			if continued == "y" {
				return false, err
			}
			return false, err
		}
		err1 := btn.Click()
		if err1 != nil {
			log.Println("*****2注册按钮不可点击****可能需要DC验证请手动操作*****")
			return false, err
		}
		return true, err
	})

}
