package premint

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/Galxe"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/metamask"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"strings"
	"time"
)

func OmniGalxe() {
	url := "https://galxe.com/OmniNetwork/campaign/GCSmgUW7Fo"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据5-1.xlsx")
	filepout := "D:\\GoWork\\resource\\failInfos.xlsx"

	fmt.Println("数据长度--------", len(excelInfos))
	chs := make(chan []string, len(excelInfos))
	var title = []string{"地址", "类型", "窗口ID", "MetaMask密码"}
	//创建新的excel文件
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &title)

	// 获取内容并写入Excel
	go func() {
		for t := 0; t < len(excelInfos); t++ {
			data := <-chs
			log.Println("接受到一条错误信息：", data)
			axis := fmt.Sprintf("A%d", t+2)
			excel.SetSheetRow("Sheet1", axis, &data)
		}
	}()
	Size := 10
	counter := 0
	var Ids []string
	for k, v := range excelInfos {
		counter++

		Ids = append(Ids, v.BitId)
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		handle, _ := wd.CurrentWindowHandle()
		wd.ResizeWindow(handle, 500, 300)
		if wd != nil {
			wg.Add(1)
			go util.SetLog(func() {
				defer wg.Done()
				err := StartOmniGalxe(v, k, chs, wd, url)
				if err != nil {
					log.Println(v.BitId, "失败")
				}
			})
		}
		if counter >= Size && (counter-1)%Size == 0 {
			bitbrowser.WindowboundsByPara()
			log.Println("counter------", counter)
			wg.Wait()
			bitbrowser.WindowboundsByPara()
			time.Sleep(3 * time.Second)
			//for _, v := range Ids {
			//	bitbrowser.CloseBrower(v)
			//}
			Ids = Ids[:0]
		}
	}

	//for k, v := range excelInfos {
	//	fmt.Println("----------", v.Address)
	//	//打开比特浏览器
	//	wd, _ := wdservice.InitWd(k, v.BitId)
	//
	//	if wd != nil {
	//		wg.Add(1)
	//		go util.SetLog(func() {
	//			defer wg.Done()
	//			StartOmniGalxe(v, k, chs, wd, url)
	//			//bitbrowser.CloseBrower(v.BitId)
	//		})
	//	}
	//	bitbrowser.WindowboundsByPara()
	//	wg.Wait()
	//}
	wg.Wait()
	close(chs)
	excel.SaveAs(filepout)

}
func StartOmniGalxe(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string) error {
	log.Println("*********************开始处理第" + strconv.Itoa(i+1) + "条数据******************")
	/*	打开网址登陆小狐狸
	 */
	metamask.MetaMaskLogin(wd, excelInfo.MetaPwd)
	time.Sleep(1 * time.Second)

	handle := util.GetCurrentWindowAndReturn(wd)
	////关闭多余标签页
	bitbrowser.CloseOtherLabels(wd, handle)
	wd.SwitchWindow(handle)
	time.Sleep(2 * time.Second)

	log.Println("打开银河链接")
	err := wd.Get(url)
	if err != nil {
		log.Println("打开银河链接出错了-----", excelInfo.Address)
		wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
		ch <- wrongData
		return err
	} else {
		log.Println("银河打开成功")
	}
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)
	//关闭小狐狸，第一次会有登陆
	main_handle, _ := wd.WindowHandles()
	//如果打开了小狐狸
	if len(main_handle) > 1 {
		LoginRequest(wd, main_handle)
	}

	time.Sleep(1 * time.Second)
	bitbrowser.CloseOtherLabels(wd, handle)
	wd.SwitchWindow(handle)
	time.Sleep(6 * time.Second)
	handleNow, _ := wd.CurrentWindowHandle()
	wd.MaximizeWindow(handleNow)
	err = ChooseNetwork(wd, "Polygon")
	if err != nil {
		log.Println("切换网址失败-----", excelInfo.Address)
		wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
		ch <- wrongData
		return err
	}
	//如果打开了小狐狸
	time.Sleep(1 * time.Second)
	main_handle, _ = wd.WindowHandles()
	if len(main_handle) > 1 {
		LoginRequest(wd, main_handle)
		err = ConfirmMeta(wd, main_handle)
	}
	if err != nil {
		log.Println("小狐狸关闭弹窗失败-----", excelInfo.Address)
		wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
		ch <- wrongData
		return err
	}
	time.Sleep(1 * time.Second)
	wd.SwitchWindow(handle)
	time.Sleep(1 * time.Second)

	//打开所需下拉框
	texts, err := OmniFindAllDropDownBox(wd, 3)
	if err != nil {
		wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
		ch <- wrongData
		log.Println("打开下拉框失败了------", excelInfo.Address)
		return err
	}
	bitbrowser.CloseOtherLabels(wd, handle)

	//处理打开的url
	time.Sleep(1 * time.Second)
	err = UrlTask(wd, 3, handle, texts)
	//bitbrowser.WindowboundsByPara()
	if err != nil {
		wrongData := []string{excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
		ch <- wrongData
		log.Println("url处理失败了----", excelInfo.Address)
		return err
	}
	//	点击刷新按钮 并刷新整个界面
	bitbrowser.CloseOtherLabels(wd, handle)
	wd.Get(url)
	time.Sleep(1 * time.Second)
	//缩小窗口会挡住
	//Refresh(wd)
	time.Sleep(1 * time.Second)

	fmt.Println("处理完毕")
	wd.ResizeWindow(handle, 500, 300)
	return nil
}
func Refresh(wd selenium.WebDriver) {
	//*[@id="ga-data-campaign-model-2"]/div[2]/div[1]/div[5]/div/div[1]/div/div/div[4]/div[1]/div/button/div[2]/div/div/div/div
	refresh, _ := wd.FindElements(selenium.ByCSSSelector, ".clickable.refresh.icon.responsive")
	log.Println("查找到刷新按钮的数量：", len(refresh))
	if len(refresh) > 0 {
		for _, v := range refresh {
			err := v.Click()
			if err != nil {
				log.Println(err)
			} else {
				log.Println("点击失败")
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		log.Println("查找刷新按钮失败")
	}
}
func UrlTaskNull(wd selenium.WebDriver, num int, handle string, texts []string) error {
	return nil
}
func UrlTask(wd selenium.WebDriver, num int, handle string, texts []string) error {
	handleNow, _ := wd.CurrentWindowHandle()
	wd.MaximizeWindow(handleNow)
	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
	if err != nil {
		log.Println("查找任务失败")
		//bitbrowser.WindowboundsByPara()
		return err
	} else {
		//bitbrowser.WindowboundsByPara()
		fmt.Println("当前页面url总数量", len(aUrl))
		nowHandle, _ := wd.CurrentWindowHandle()
		log.Println("nowHandle:", nowHandle)
		for k, v := range aUrl {
			if k > len(aUrl)-num && len(aUrl) > num {
				continue
			} else {
				handleNow, _ := wd.CurrentWindowHandle()
				wd.MaximizeWindow(handleNow)
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
					log.Println("点击第", k, "个url失败：", err)
					//bitbrowser.WindowboundsByPara()
					return err
				} else {
					//bitbrowser.WindowboundsByPara()
					handles, _ := wd.WindowHandles()
					wd.SwitchWindow(handles[1])
					CurrentHandle1, _ := wd.CurrentWindowHandle()
					log.Println("打开detail第一次的", CurrentHandle1)
					//10秒等待元素点击
					handleNow1, _ := wd.CurrentWindowHandle()
					wd.MaximizeWindow(handleNow1)
					time.Sleep(1 * time.Second)
					err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
						for i := 0; i < 10; i++ {
							_, err = wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
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
						//bitbrowser.WindowboundsByPara()
						log.Println("查找第二重url失败")
						return err
					} else {
						//第一个任务
						url, _ := wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
						//bitbrowser.WindowboundsByPara()
						switch {
						case strings.Contains(texts[k], "Omni") && strings.Contains(texts[k], "Users"):
							err1 := Galxe.GalxeFollow(wd, url)
							//选择主页面并关闭其他页面
							bitbrowser.CloseOtherLabels(wd, handle)
							wd.SwitchWindow(handle)
							time.Sleep(2 * time.Second)
							if err1 != nil {
								return err1
							}
						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Follow"):
							//err1 := twitter.TwitterFollow(wd, url)
							//bitbrowser.CloseOtherLabels(wd, handle)
							//wd.SwitchWindow(handle)
							//time.Sleep(2 * time.Second)
							//if err1 != nil {
							//	return err1
							//}
						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Tweet"):
							//err1 := twitter.TwitterTweet(wd, url)
							//bitbrowser.CloseOtherLabels(wd, handle)
							//wd.SwitchWindow(handle)
							//time.Sleep(2 * time.Second)
							//if err1 != nil {
							//	return err1
							//}
						case strings.Contains(texts[k], "Discord") && strings.Contains(texts[k], "verified"):

						default:
							fmt.Println("未知方法")
						}
					}
				}
			}
		}
	}
	return nil
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
func OmniFindAllDropDownBox(wd selenium.WebDriver, num int) (texts []string, err error) {
	handleNow, _ := wd.CurrentWindowHandle()
	wd.MaximizeWindow(handleNow)
	dropdownElements, _ := wd.FindElements(selenium.ByCSSSelector, ".v-expansion-panel")
	length := len(dropdownElements) - num
	dropLen := len(dropdownElements)
	log.Println(len(dropdownElements))
	//bitbrowser.WindowboundsByPara()

	for k, dropdownElement := range dropdownElements {
		//跳过一定数量的下拉框
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
				//点击下拉框
				wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < 5; i++ {
						handleNow, _ := wd.CurrentWindowHandle()
						wd.MaximizeWindow(handleNow)
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
					//bitbrowser.WindowboundsByPara()
					return nil, err
				}
				//bitbrowser.WindowboundsByPara()
			}
		}
	}
	//bitbrowser.WindowboundsByPara()

	return
}
func LoginRequest(wd selenium.WebDriver, main_handle []string) {
	wd.SwitchWindow(main_handle[1])
	time.Sleep(1 * time.Second)
	knows, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
	if len(knows) > 0 {
		for _, know := range knows {
			know.Click()
		}
		LoginButon, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
		log.Println("处理第一次小狐狸登陆,button长度：", len(LoginButon))
		if len(LoginButon) > 0 {
			for _, v := range LoginButon {
				v.Click()
			}
		}
	}
}
func AllowAdd(wd selenium.WebDriver, main_handle []string) {
	wd.SwitchWindow(main_handle[1])
	time.Sleep(1 * time.Second)
	LoginButon, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
	log.Println("处理第一次小狐狸登陆,button长度：", len(LoginButon))
	if len(LoginButon) > 0 {
		for _, v := range LoginButon {
			v.Click()
		}
	}
}
