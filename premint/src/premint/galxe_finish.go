package premint

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"log"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// 通过excel获取数据打开比特浏览器
// 银河链接
func GalxeFinish() {
	//url := "https://galxe.com/EchoDEX/campaign/GCDsmUSvqd"
	url := "https://galxe.com/OmniNetwork/campaign/GCSmgUW7Fo"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-1.xlsx")
	fmt.Println("数据长度--------", len(excelInfos))
	chs := make(chan string, len(excelInfos))
	for k, v := range excelInfos {
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		if wd != nil {
			wg.Add(1)
			//go util.SetLog(func() {
			go StartGalxe(v, k, chs, wd, url)
			//})
		}
	}
	//bitbrowser.WindowboundsByPara()

	wg.Wait()
}

func StartGalxe(excelInfo model.OMNIExcelInfo, i int, ch chan<- string, wd selenium.WebDriver, url string) {
	//log.Println("*********************开始处理第" + strconv.Itoa(i+1) + "条数据******************")
	///*	打开网址登陆小狐狸
	// */
	//metamask.MetaMaskLogin(wd, excelInfo.MetaPwd)
	//time.Sleep(1 * time.Second)
	//
	//log.Println("打开银河链接")
	//err := wd.Get(url)
	//if err != nil {
	//	log.Println("打开银河链接出错了")
	//} else {
	//	log.Println("银河打开成功")
	//
	//}
	//time.Sleep(2 * time.Second)
	//
	//handle := util.GetCurrentWindowAndReturn(wd)
	////关闭多余标签页
	//bitbrowser.CloseOtherLabels(wd, handle)
	//wd.SwitchWindow(handle)
	//time.Sleep(5 * time.Second)
	//
	//ChooseNetwork(wd, "Polygon")
	//time.Sleep(2 * time.Second)
	//main_handle, err := wd.WindowHandles()
	////如果打开了小狐狸
	//if len(main_handle) > 1 {
	//	err = ConfirmMeta(wd, main_handle)
	//}
	////打开所需下拉框
	//time.Sleep(2 * time.Second)
	//wd.SwitchWindow(handle)

	FindAllDropDownBox(wd, 8)
	time.Sleep(2 * time.Second)
	//一个一个打开链接并完成任务detail-text text-14-regular clickable
	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
	if err != nil {
		log.Println("查找任务失败")
	} else {
		fmt.Println("当前页面url总数量", len(aUrl))
		nowHandle, _ := wd.CurrentWindowHandle()
		log.Println("nowHandle:", nowHandle)
		for k, v := range aUrl {
			//if k < 3 || k == len(aUrl)-1 {
			//	continue
			//} else {
			//	log.Println("开始第", k, "条url的打开")
			//	v.Click()
			//	time.Sleep(5 * time.Second)
			//}
			//handles, _ := wd.WindowHandles()
			//wd.SwitchWindow(handles[0])
			if k > len(aUrl)-3 || len(aUrl) < 3 {
				continue
			} else {
				v.Click()
				//wd.Get("https://galxe.com/credential/293670416934412288")
				//util.GetCurrentWindow(wd)
				//var s []interface{}
				//s = append(s, "https://galxe.com/credential/293670416934412288")
				//wd.ExecuteScript("window.open(arguments[0])", s)
				time.Sleep(3 * time.Second)
				handles, _ := wd.WindowHandles()
				wd.SwitchWindow(handles[1])
				CurrentHandle1, _ := wd.CurrentWindowHandle()
				log.Println("打开detail的", CurrentHandle1)
				url, err := wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
				if err != nil {
					log.Println("查找url失败")
				} else {
					url.Click()
					handles, _ := wd.WindowHandles()
					fmt.Println("打开twitter后handle的长度", len(handles))
					time.Sleep(2 * time.Second)
					wd.SwitchWindow(handles[len(handles)-1])
					CurrentHandle1, _ := wd.CurrentWindowHandle()

					log.Println("打开twitter的", CurrentHandle1)
					//css-18t94o4 css-1dbjc4n r-42olwf r-sdzlij r-1phboty r-rs99b7 r-16y2uox r-6gpygo r-peo1c r-1ps3wis r-1ny4l3l r-1udh08x r-1guathk r-1udbk01 r-o7ynqc r-6416eg r-lrvibr r-3s2u2q
					//css-18t94o4 css-1dbjc4n r-42olwf r-sdzlij r-1phboty r-rs99b7 r-16y2uox r-6gpygo r-peo1c r-1ps3wis r-1ny4l3l r-1udh08x r-1guathk r-1udbk01 r-o7ynqc r-6416eg r-lrvibr r-3s2u2q
					follow, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-16y2uox.r-6gpygo.r-peo1c.r-1ps3wis.r-1ny4l3l.r-1udh08x.r-1guathk.r-1udbk01.r-o7ynqc.r-6416eg.r-lrvibr.r-3s2u2q")

					//follow, err := wd.FindElements(selenium.ByCSSSelector, ".css-901oao.r-1awozwy.r-jwli3a.r-6koalj.r-18u37iz.r-16y2uox.r-37j5jr.r-a023e6.r-b88u0q.r-1777fci.r-rjixqe.r-bcqeeo.r-q4m81j.r-qvutc0")
					if err != nil {
						log.Println("没有找到follow按钮")
					} else {
						//log.Println("follow查找到的长度", len(follow))
						follow.Click()
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
		bitbrowser.CloseOtherLabels(wd, nowHandle)
	}
	wg.Done()
	fmt.Println("处理完毕")

}

// 选择网络
func ChooseNetwork(wd selenium.WebDriver, str string) {
	button, err := wd.FindElements(selenium.ByCSSSelector, ".text-14-regular.text-overline-ellipsis-1")

	if err != nil {
		log.Println("查找元素出错了")
	} else {
		//fmt.Println(button.Text())
		log.Println("查找成功了")
		time.Sleep(2 * time.Second)

		err = button[1].Click()
		if err != nil {
			log.Println("点击失败")

		}

	}
	time.Sleep(1 * time.Second)
	d1, err := wd.FindElements(selenium.ByCSSSelector, ".wallet-option-item.text-16-bold")
	if err != nil {
		log.Println("查找元素出错了")
	} else {
		log.Println("查找成功了")
		time.Sleep(1 * time.Second)
		for k, v := range d1 {
			text, _ := v.Text()
			if strings.Contains(text, str) {
				err = d1[k].Click()
				if err != nil {
					log.Println("点击失败")
				}
			}
		}
	}
}

// 小狐狸确认
func ConfirmMeta(wd selenium.WebDriver, main_handle []string) error {
	wd.SwitchWindow(main_handle[1])
	url, err := wd.CurrentURL()
	log.Println(url)
	time.Sleep(1 * time.Second)
	knows, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	if len(knows) > 0 {
		for _, know := range knows {
			know.Click()
		}
	} else {
		log.Println("未找到明白了按钮")
	}
	button, err := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	fmt.Println(len(button))
	if err != nil {
		log.Println("查找切换网络失败")
		return err
	} else {
		button[0].Click()
	}
	return nil
}
func FindAllDropDownBox(wd selenium.WebDriver, num int) {
	dropdownElements, _ := wd.FindElements(selenium.ByCSSSelector, ".v-expansion-panel")

	length := len(dropdownElements) - num
	log.Println(len(dropdownElements))
	for k, dropdownElement := range dropdownElements {
		if k < length {
			continue
		}
		log.Println("开始第", k, "个下拉框的处理")
		// 判断下拉框状态
		time.Sleep(1 * time.Second)
		ariaExpanded, err := dropdownElement.GetAttribute("aria-expanded")
		if err != nil {
			log.Println(err)
			continue
		}
		if ariaExpanded == "false" {
			// 点击下拉框头部元素以展开
			err = dropdownElement.Click()
			if err != nil {
				log.Println(err)
				continue
			} else {
				// 等待一段时间以确保下拉框展开完全
				time.Sleep(1 * time.Second)
				//处理detail
			}
		}
	}
}
func StartTasks(wd selenium.WebDriver) {
	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
	if err != nil {
		log.Println("查找任务失败")
	} else {
		for k, v := range aUrl {
			fmt.Println("第", k, "条url的值为:", v)
		}
	}
}
