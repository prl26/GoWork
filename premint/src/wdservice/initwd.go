/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-20 09:18:13
 * @LastEditTime: 2023-04-27 10:43:48
 */
package wdservice

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	port      = 9000
	chromeDri = "D:\\bitbrowser\\resources\\chromedriver\\112\\chromedriver.exe"
)

// 打开比特浏览器，跳转到最后一个窗口
func InitWd(i int, bitID string) (selenium.WebDriver, int) {
	time.Sleep(3 * time.Second)
	bitData := bitbrowser.OpenBrowser(bitID)
	pport, _ := GetFreePort()
	//pport, _ := strconv.Atoi(strings.Split(bitData.Http, ":")[1])
	ops := []selenium.ServiceOption{}
	chromeCaps := chrome.Capabilities{
		DebuggerAddr: bitData.Http,
	}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddChrome(chromeCaps)
	selenium.NewChromeDriverService(bitData.Driver, pport, ops...)
	log.Println(fmt.Sprintf("http://127.0.0.1:%d/wd/hub", pport))
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", pport))
	if err != nil {
		log.Println("初始化浏览器失败", err)
		return nil, 0
	}
	log.Println(wd)
	handles, _ := wd.WindowHandles()
	numWin := len(handles) - 1
	log.Println("窗口数量：", numWin, handles)
	//util.GetWind(0, wd)
	return wd, numWin
}

// 测试
func InitWdTest(i int, bitID string) (selenium.WebDriver, int) {
	//bitData := bitbrowser.OpenBrowser(bitID)
	pport := port + i
	//pport, _ := strconv.Atoi(strings.Split(bitData.Http, ":")[1])
	ops := []selenium.ServiceOption{}
	chromeCaps := chrome.Capabilities{
		DebuggerAddr: "127.0.0.1:50772",
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddChrome(chromeCaps)
	selenium.NewChromeDriverService("C:\\Users\\31972\\AppData\\Roaming\\bitbrowser\\chromedriver\\104\\chromedriver.exe", pport, ops...)

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://127.0.0.1:%d/wd/hub", pport))
	if err != nil {
		log.Println("初始化浏览器", err)
	}
	log.Println(wd)
	handles, _ := wd.WindowHandles()
	numWin := len(handles) - 1
	util.GetWind(numWin, wd)
	return wd, numWin
}

func InitCmd() int {
	fmt.Println("请选择要操作的内容：")
	fmt.Println("1: Premint注册")
	fmt.Println("2: 钱包注册")
	fmt.Println("3: metamask添加网络")
	fmt.Println("4: 币安操作")
	fmt.Println("5: MetaMask-OKX操作")
	fmt.Println("6: Twitter关注转发")
	fmt.Println("7: 银河操作")
	fmt.Println("8: Discord接受邀请")
	fmt.Println("9: 表单操作")
	fmt.Println("10: Claim")
	fmt.Printf("请输入选项:")
	var cmd int = 0
	fmt.Scanln(&cmd)
	//util.CheckRole()
	return cmd
}

func InitWalletCmd() int {
	fmt.Printf("请输入生成的钱包数量:")
	var count int = 0
	fmt.Scanln(&count)

	return count
}

func InitWalletTypeCmd() int {
	fmt.Println("请选择要生成的类型：")
	fmt.Println("1: ETH")
	fmt.Println("2: Apots")
	fmt.Println("3: Cosmos")
	fmt.Println("4: Sqlana")
	fmt.Println("5: Sui")
	fmt.Printf("请输入选项:")
	var count int = 0
	fmt.Scanln(&count)

	return count
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
