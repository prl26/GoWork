package premint

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/metamask"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// 通过excel获取数据打开比特浏览器
// 银河链接
func GalxeFinish() {

	excelInfos := util.GetOMNIExcelInfos("D:\\Go Work\\resource\\测试数据1条.xlsx")
	fmt.Println("数据长度--------", len(excelInfos))
	chs := make(chan string, len(excelInfos))
	for k, v := range excelInfos {
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		if wd != nil {
			wg.Add(1)
			//go util.SetLog(func() {
			go StartGalxe(v, k, chs, wd)
			//})
		}
	}
	//bitbrowser.WindowboundsByPara()

	wg.Wait()
}

func StartGalxe(excelInfo model.OMNIExcelInfo, i int, ch chan<- string, wd selenium.WebDriver) {
	log.Println("*********************开始处理第" + strconv.Itoa(i+1) + "条数据******************")
	/*	打开网址登陆小狐狸
	 */
	metamask.MetaMaskLogin(wd, excelInfo.MetaPwd)
	time.Sleep(1 * time.Second)

	log.Println("打开银河链接")
	err := wd.Get("https://galxe.com/EchoDEX/campaign/GCDsmUSvqd")
	if err != nil {
		log.Println("打开银河链接出错了")
	} else {
		log.Println("银河打开成功")

	}

	handle := util.GetCurrentWindowAndReturn(wd)
	time.Sleep(5 * time.Second)
	//关闭多余标签页
	bitbrowser.CloseOtherLabels(wd, handle)

	ChooseNetwork(wd, "polygon")
	main_handle, err := wd.WindowHandles()
	//如果打开了小狐狸

	if len(main_handle) > 1 {
		err = ConfirmMeta(wd, main_handle)
	}
	//打开所需下拉框
	FindAllDropDownBox(wd, 8)

	//一个一个打开链接并完成任务
	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
	if err != nil {
		log.Println("查找任务失败")
	} else {
		for k, v := range aUrl {
			fmt.Println("第", k, "条url的值为:", v)
		}
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
