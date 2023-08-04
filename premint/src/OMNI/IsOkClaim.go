package OMNI

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/premint"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
	"time"
)

func OmniCheckStart() {
	file, err := os.Open("config.txt")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 使用bufio包创建一个新的Scanner用于读取文件内容
	scanner := bufio.NewScanner(file)
	var lines []string
	// 循环读取文件的每一行
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	// 检查是否有错误发生在scanner.Scan()过程中
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时发生错误:", err)
	}
	url := lines[0]
	//workingDir, err := os.Getwd()
	//url := "https://galxe.com/OmniNetwork/campaign/GCXHgUWBKg"
	//path1 := fmt.Sprintf("%v\\测试数据-200.xlsx", workingDir)
	//path2 := fmt.Sprintf("%v\\failInfo.xlsx", workingDir)
	//path3 := fmt.Sprintf("%v\\failInfo.txt", workingDir)
	//path4 := fmt.Sprintf("%v\\successInfo.txt", workingDir)
	path1 := "./resource/测试数据-200.xlsx"
	path2 := "./resource/failInfo-CheckIsOk.xlsx"
	path3 := "./resource/failInfo-CheckIsOk.txt"
	path4 := "./resource/successInfo-CheckIsOk.txt"
	//excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	//filepout := "D:\\GoWork\\resource\\FailInfos\\OmniGalxe1-200.xlsx"
	//TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\OmniGalxe1-200.txt"

	excelInfos := util.GetOMNIExcelInfos(path1)
	filepout := path2
	TxtfileOut := path3
	TxtSuccessOut := path4
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
				err := CheckIsOk(v, k, chs, wd, url, dstFile, successFile)
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
	}
}

func CheckIsOk(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, successFile *os.File) (err error) {
	wrongData := []string{excelInfo.HelpWords, excelInfo.PrivateKey, excelInfo.PublicKey, excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
	wd.Get(url)
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

	nowHandle, _ := wd.WindowHandles()
	if len(nowHandle) > 1 {
		wd.SwitchWindow(nowHandle[1])
		err = SmallFoxLoginNoSign(wd, excelInfo.MetaPwd)
	}
	if err != nil {
		log.Println("小狐狸登陆出错-----", excelInfo.BitId)
		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("小狐狸登陆成功")
		time.Sleep(1 * time.Second)
		hanNow, _ := wd.WindowHandles()
		wd.SwitchWindow(hanNow[0])
		wd.MaximizeWindow(hanNow[0])
		time.Sleep(1 * time.Second)
	}

	err = premint.ChooseNetwork(wd, "Polygon")
	if err != nil {
		log.Println("切换网址失败-----", excelInfo.Address)
		ch <- wrongData
		return err
	}
	time.Sleep(3 * time.Second)
	//关闭小狐狸弹窗
	main_handle, _ := wd.WindowHandles()
	if len(main_handle) > 1 {
		premint.LoginRequest(wd, main_handle)
		err = premint.ConfirmMeta(wd, main_handle)
	}
	if err != nil {
		log.Println("小狐狸关闭弹窗失败-----", excelInfo.Address)
		ch <- wrongData
		return err
	}
	time.Sleep(1 * time.Second)
	//开始Claim
	hanNow1, _ := wd.WindowHandles()
	wd.SwitchWindow(hanNow1[0])
	log.Println("开始Claim")
	isOK := false
	Claim1, err1 := wd.FindElement(selenium.ByCSSSelector, ".g-btn.width-max-100.v-btn.v-btn--block.v-btn--is-elevated.v-btn--has-bg.theme--dark.v-size--default.primary.text-16-bold")
	if err1 == nil {
		Claim1.Click()
		//text-left trx-title
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				text, err := wd.FindElement(selenium.ByCSSSelector, ".text-left.trx-title")
				if err == nil {
					detail, err1 := text.Text()
					if err1 == nil {
						if strings.Contains(detail, "Transaction has been submitted") {
							return true, nil
						}
					}
				} else {
					time.Sleep(1 * time.Second)
					continue
				}
			}
			return false, errors.New("fail")
		}, 15*time.Second)
		if err == nil {
			log.Println("点击Claim1成功 --老版")
			time.Sleep(2 * time.Second)
			successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
			return nil
		}
	} else if claim2, err2 := wd.FindElement(selenium.ByCSSSelector, ".g-btn.width-max-100.v-btn.v-btn--block.v-btn--is-elevated.v-btn--has-bg.theme--dark.v-size--x-large.primary.text-16-bold"); err2 == nil {
		//g-btn width-max-100 v-btn v-btn--block v-btn--is-elevated v-btn--has-bg theme--dark v-size--x-large primary text-16-bold  新版的
		claim2.Click()
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				text, err := wd.FindElement(selenium.ByCSSSelector, ".text-left.trx-title")
				if err == nil {
					detail, err1 := text.Text()
					if err1 == nil {
						if strings.Contains(detail, "Transaction has been submitted") {
							return true, nil
						}
					}
				} else {
					time.Sleep(1 * time.Second)
					continue
				}
			}
			return false, errors.New("fail")
		}, 15*time.Second)
		if err == nil {
			log.Println("点击Claim2成功--新版")
			time.Sleep(2 * time.Second)
			successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
			return nil
		}
	} else if isClaim, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(isClaim) > 0 {
		for _, v := range isClaim {
			text, _ := v.Text()
			if strings.Contains(text, "Claim") {
				log.Println("你已经领取过了")
				successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
				isOK = true
				time.Sleep(5 * time.Second)
				return nil
			}
		}
	}
	if isOK == false {
		log.Println("Claim失败")
		time.Sleep(3 * time.Second)
		dstFile.WriteString(fmt.Sprintf("Claim失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return errors.New("Fail")
	}
	return err
}
