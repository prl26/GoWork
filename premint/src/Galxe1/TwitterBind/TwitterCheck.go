package TwitterBind

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/atotto/clipboard"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
	"time"
)

// https://galxe.com/accountSetting?tab=SocialLinlk

func Check() {
	url := "https://galxe.com/accountSetting?tab=SocialLinlk"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\Check-8.2.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\Check-8.2.txt"
	TxtSuccessOut := "D:\\GoWork\\resource\\SuccessInfos\\Check-8.2.txt"
	dstFile, err := os.OpenFile(TxtfileOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer dstFile.Close()
	successFile, err := os.OpenFile(TxtSuccessOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return
	}
	defer successFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"助记词", "私钥", "公钥", "地址", "类型", "窗口ID", "MetaMask密码"}
	//创建新的excel文件
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &title)

	//定义一次开多少线程
	fmt.Println("数据长度--------", len(excelInfos))
	// 获取内容并写入Excel
	go func() {
		for t := 0; t < len(excelInfos); t++ {
			data := <-chs
			log.Println("接受到一条错误信息：", data)
			axis := fmt.Sprintf("A%d", t+2)
			excel.SetSheetRow("Sheet1", axis, &data)
		}
	}()

	////单个打开
	for k, v := range excelInfos {
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		if wd != nil {
			handle, _ := wd.WindowHandles()
			if len(handle) > 1 {
				handle1 := util.GetCurrentWindowAndReturn(wd)
				//关闭多余标签页
				bitbrowser.CloseOtherLabels(wd, handle1)
				wd.SwitchWindow(handle1)
			}
			time.Sleep(1 * time.Second)
			wg.Add(1)
			go util.SetLog(func() {
				defer wg.Done()
				err := check(v, k, chs, wd, url, dstFile, successFile)
				if err != nil {
					log.Println("!-------!", v.BitId, "失败")
				}
				defer bitbrowser.CloseBrower(v.BitId)
			})
		}
		wg.Wait()
	}
	close(chs)
	err = excel.SaveAs(filepout)
	if err != nil {
		log.Println("excel 保存失败----", err)
	} else {
		log.Println("excel 保存成功----", err)

	}
}
func check(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, successFile *os.File) (err error) {
	wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	//打开银河链接
	err = wd.Get(url)
	if err != nil {
		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("银河打开成功")
	}
	time.Sleep(2 * time.Second)
	//小狐狸登陆
	nowHandle, _ := wd.WindowHandles()
	if len(nowHandle) > 1 {
		wd.SwitchWindow(nowHandle[1])
		err = SmallFoxLogin(wd)
	}
	if err != nil {
		log.Println("小狐狸登陆出错-----", excelInfo.BitId)
		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("小狐狸登陆成功")
	}
	time.Sleep(1 * time.Second)
	//打开个人首页
	handleNow, _ := wd.WindowHandles()
	wd.SwitchWindow(handleNow[0])
	err = openHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString(fmt.Sprintf("进入个人主页失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
	}
	time.Sleep(1 * time.Second)

	//先保持未登陆状态
	//1.点击DisConnect
	err = ClickDisconnect(wd)
	if err != nil {
		log.Println("点击DisConnect失败")
		dstFile.WriteString(fmt.Sprintf("点击DisConnect失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("点击DisConnect成功")
	}
	time.Sleep(1 * time.Second)

	//这里resize一下
	wd.ResizeWindow(handleNow[0], 1500, 1440)

	//然后连接小狐狸
	err = ConnectMetamask(wd)
	if err != nil {
		log.Println("连接小狐狸失败")
		dstFile.WriteString(fmt.Sprintf("连接小狐狸失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("连接小狐狸成功")
	}
	time.Sleep(2 * time.Second)

	//处理小狐狸登陆
	err = MetamaskLogin(wd)
	if err != nil {
		log.Println("处理小狐狸登陆失败")
		dstFile.WriteString(fmt.Sprintf("处理小狐狸登陆失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("处理小狐狸登陆成功")
		time.Sleep(2 * time.Second)
	}
	//如果还有弹窗
	err = CloseSignInPop(wd)
	if err != nil {
		log.Println("关闭小狐狸弹窗失败")
		dstFile.WriteString(fmt.Sprintf("处理小狐狸登陆失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("关闭小狐狸弹窗成功")
		time.Sleep(1 * time.Second)
	} //获取绑定的账号
	BindUsername, err := findTwitterUsername(wd)
	if err != nil {
		log.Println("获取绑定的账号失败")
		dstFile.WriteString(fmt.Sprintf("获取绑定的账号失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("获取绑定的账号成功--", BindUsername)
		time.Sleep(2 * time.Second)
	}
	hanNow, _ := wd.WindowHandles()
	wd.SwitchWindow(hanNow[len(hanNow)-1])
	wd.MaximizeWindow(hanNow[len(hanNow)-1])
	clipboard.WriteAll("")
	userName, err := GetInProfile1(wd)
	if err != nil {
		log.Println("获取Link失败")
		dstFile.WriteString(fmt.Sprintf("获取Link失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("获取Link成功,--username", userName)
		time.Sleep(2 * time.Second)
		handowNow, _ := wd.WindowHandles()
		wd.SwitchWindow(handowNow[0])
	}
	//url, _ = wd.CurrentURL()
	//urlParts := strings.Split(url, "/")
	//
	//// Get the last part of the URL, which is the username
	//username := urlParts[len(urlParts)-1]
	//log.Println(username)
	if BindUsername == userName {
		log.Println("第二次verify成功")
		successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
		return nil
	} else {
		log.Println("第二次verify失败")
		dstFile.WriteString(fmt.Sprintf("第二次verify失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		time.Sleep(2 * time.Second)
	}
	return err
}
func GetInProfile1(wd selenium.WebDriver) (userName string, err error) {
	//点击进入profile
	var Profile selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			Profile, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"react-root\"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[9]")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err != nil {
		log.Println("进入Profile失败")
	} else {
		err = Profile.Click()
		if err != nil {
			log.Println("进入Profile失败")
			return "", err
		}
		//获取推特链接
		time.Sleep(2 * time.Second)
		url, _ := wd.CurrentURL()
		log.Println("当前url", url)
		urlParts := strings.Split(url, "/")
		username := urlParts[len(urlParts)-1]
		return username, err
	}
	return "", err
}
