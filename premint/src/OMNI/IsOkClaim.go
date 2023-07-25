package OMNI

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
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
	url := "https://galxe.com/OmniNetwork/campaign/GCSmgUW7Fo"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\OmniGalxeCheck-7.25.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\OmniGalxeCheck-7.25.txt"
	dstFile, _ := os.Create(TxtfileOut)
	defer dstFile.Close()
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
				err := CheckIsOk(v, k, chs, wd, url, dstFile)
				if err != nil {
					log.Println("!-------!", v.BitId, "失败")
				}
				defer bitbrowser.CloseBrower(v.BitId)
			})
		}
		wg.Wait()
	}
	close(chs)
	err := excel.SaveAs(filepout)
	if err != nil {
		log.Println("excel 保存失败----", err)
	}
}

func CheckIsOk(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
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
		err = SmallFoxLoginNoSign(wd)
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
	//开始Claim
	log.Println("开始Claim")
	isOK := false
	Claim1, err1 := wd.FindElement(selenium.ByCSSSelector, ".g-btn.width-max-100.v-btn.v-btn--block.v-btn--is-elevated.v-btn--has-bg.theme--dark.v-size--default.primary.text-16-bold")
	if err1 == nil {
		err = Claim1.Click()
		if err == nil {
			log.Println("点击Claim成功")
			time.Sleep(4 * time.Second)
			return nil
		}
	} else if claim2, err2 := wd.FindElement(selenium.ByCSSSelector, ".g-btn.width-max-100.v-btn.v-btn--block.v-btn--is-elevated.v-btn--has-bg.theme--dark.v-size--x-large.primary.text-16-bold"); err2 == nil {
		//g-btn width-max-100 v-btn v-btn--block v-btn--is-elevated v-btn--has-bg theme--dark v-size--x-large primary text-16-bold  新版的
		err = claim2.Click()
		if err == nil {
			log.Println("点击Claim成功")
			time.Sleep(4 * time.Second)
			return nil
		}
	} else if isClaim, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(isClaim) > 0 {
		log.Println(len(isClaim))
		for _, v := range isClaim {
			text, _ := v.Text()
			if strings.Contains(text, "Claim") {
				log.Println("你已经领取过了")
				isOK = true
				time.Sleep(3 * time.Second)
				return nil
			}
		}
	}
	if isOK == false {
		log.Println("点击Claim失败")
		dstFile.WriteString(fmt.Sprintf("点击Claim失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return errors.New("Fail")
	}
	return err
}
