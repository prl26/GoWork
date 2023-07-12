package premint

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/Galxe"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/twitter"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"log"
	"strconv"
	"strings"
	"time"
)

func OmniGalxe() {
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
			go StartOmniGalxe(v, k, chs, wd, url)
			//})
		}
	}
	//bitbrowser.WindowboundsByPara()

	wg.Wait()
}
func StartOmniGalxe(excelInfo model.OMNIExcelInfo, i int, ch chan<- string, wd selenium.WebDriver, url string) {
	log.Println("*********************开始处理第" + strconv.Itoa(i+1) + "条数据******************")
	/*	打开网址登陆小狐狸
	 */
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
	handle := util.GetCurrentWindowAndReturn(wd)
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
	texts := OmniFindAllDropDownBox(wd, 3)

	//处理打开的url
	time.Sleep(1 * time.Second)
	UrlTask(wd, 3, handle, texts)

	bitbrowser.CloseOtherLabels(wd, handle)
	fmt.Println("处理完毕")
	wg.Done()
}

func UrlTask(wd selenium.WebDriver, num int, handle string, texts []string) {
	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
	if err != nil {
		log.Println("查找任务失败")
	} else {
		fmt.Println("当前页面url总数量", len(aUrl))
		nowHandle, _ := wd.CurrentWindowHandle()
		log.Println("nowHandle:", nowHandle)
		for k, v := range aUrl {
			if k > len(aUrl)-num && len(aUrl) < num {
				continue
			} else {
				log.Println("开始第", k, "个url的点击")
				err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < 10; i++ {
						err := v.Click()
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
					log.Println("点击url失败：", err)
				} else {
					handles, _ := wd.WindowHandles()
					wd.SwitchWindow(handles[1])
					CurrentHandle1, _ := wd.CurrentWindowHandle()
					log.Println("打开detail第一次的", CurrentHandle1)
					time.Sleep(2 * time.Second)
					//10秒等待元素点击
					err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
						for i := 0; i < 10; i++ {
							_, err := wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
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
						log.Println("查找url失败")
					} else {
						//第一个任务
						url, _ := wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
						switch {
						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Follow"):
							Galxe.GalxeFollow(wd, url)
							//选择主页面并关闭其他页面
							log.Println("切换为主页", handle)
							bitbrowser.CloseOtherLabels(wd, handle)
							wd.SwitchWindow(handle)
							time.Sleep(2 * time.Second)
						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Liker"):
							twitter.TwitterFollow(wd, url)
							bitbrowser.CloseOtherLabels(wd, handle)
							wd.SwitchWindow(handle)
							time.Sleep(2 * time.Second)
						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Retweeters"):

						case strings.Contains(texts[k], "Discord") && strings.Contains(texts[k], "verified"):

						default:
							fmt.Println("未知方法")
						}
					}
				}
			}
		}
	}
}
func UrlTaskToOne(wd selenium.WebDriver) {
	time.Sleep(3 * time.Second)
	handles, _ := wd.WindowHandles()

	wd.SwitchWindow(handles[1])
	CurrentHandle1, _ := wd.CurrentWindowHandle()
	urlNow, _ := wd.CurrentURL()
	log.Println("打开detail的", urlNow)
	log.Println("打开detail的", CurrentHandle1)

}
func OmniFindAllDropDownBox(wd selenium.WebDriver, num int) (texts []string) {
	dropdownElements, _ := wd.FindElements(selenium.ByCSSSelector, ".v-expansion-panel")

	length := len(dropdownElements) - num
	dropLen := len(dropdownElements)
	log.Println(len(dropdownElements))
	for k, dropdownElement := range dropdownElements {
		if k < length && dropLen > num {
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
			// 点击下拉框头部元素以展开   这里不点击整个div，而是点击里面具体的一个div
			//err = dropdownElement.Click()
			expand, err := dropdownElement.FindElement(selenium.ByCSSSelector, ".expand-icon")
			if err != nil {
				log.Println(err)
				continue
			} else {
				// 等待一段时间以确保下拉框展开完全
				time.Sleep(1 * time.Second)
				text, _ := dropdownElement.FindElement(selenium.ByCSSSelector, ".cred-name.usual-text.text-overline-ellipsis-1")
				textDetail, _ := text.Text()
				fmt.Println("text的值：", textDetail)
				texts = append(texts, textDetail)
				wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < 5; i++ {
						err1 := expand.Click()
						if err1 != nil {
							time.Sleep(1 * time.Second)
							continue
						} else {
							return true, nil
						}
					}
					return false, errors.New("失败")
				}, 6*time.Second)
				if err != nil {
					log.Println("打开下拉框失败")
				}
				//判断此次点击是否打开了url --- 处理detail
				//if nowHadle, _ := wd.WindowHandles(); len(nowHadle) != 1 {
				//	log.Println(len(nowHadle))
				//	log.Println(nowHadle)
				//	UrlTask(wd)
				//	time.Sleep(2 * time.Second)
				//}
			}
		}
	}
	return
}
