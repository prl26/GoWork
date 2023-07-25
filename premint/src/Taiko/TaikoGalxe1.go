package Taiko

//
//import (
//	"errors"
//	"fmt"
//	"github.com/JianLinWei1/premint-selenium/model"
//	"github.com/JianLinWei1/premint-selenium/src/Galxe"
//	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
//	"github.com/JianLinWei1/premint-selenium/src/twitter"
//	"github.com/JianLinWei1/premint-selenium/src/util"
//	"github.com/JianLinWei1/premint-selenium/src/wdservice"
//	"github.com/tebeka/selenium"
//	"github.com/xuri/excelize/v2"
//	"log"
//	"os"
//	"strings"
//	"time"
//)
//
//func Taiko1Start() {
//	url := "https://galxe.com/taiko/campaign/GC8U8U5sYm"
//	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
//	filepout := "D:\\GoWork\\resource\\FailInfos\\TaikoGalxe1-jinniu200.xlsx"
//	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\TaikoGalxe1-jinniu200.txt"
//	dstFile, _ := os.Create(TxtfileOut)
//	defer dstFile.Close()
//	chs := make(chan []string, len(excelInfos))
//	var title = []string{"助记词", "私钥", "公钥", "地址", "类型", "窗口ID", "MetaMask密码"}
//	//创建新的excel文件
//	excel := excelize.NewFile()
//	excel.SetSheetRow("Sheet1", "A1", &title)
//
//	//定义一次开多少线程
//
//	fmt.Println("数据长度--------", len(excelInfos))
//
//	// 获取内容并写入Excel
//	go func() {
//		for t := 0; t < len(excelInfos); t++ {
//			data := <-chs
//			log.Println("接受到一条错误信息：", data)
//			axis := fmt.Sprintf("A%d", t+2)
//			excel.SetSheetRow("Sheet1", axis, &data)
//		}
//	}()
//
//	////单个打开
//	for k, v := range excelInfos {
//		fmt.Println("----------", v.Address)
//		//打开比特浏览器
//		wd, _ := wdservice.InitWd(k, v.BitId)
//		if wd != nil {
//			handle, _ := wd.WindowHandles()
//			if len(handle) > 1 {
//				handle1 := util.GetCurrentWindowAndReturn(wd)
//				//关闭多余标签页
//				bitbrowser.CloseOtherLabels(wd, handle1)
//				wd.SwitchWindow(handle1)
//			}
//			time.Sleep(1 * time.Second)
//			wg.Add(1)
//			go util.SetLog(func() {
//				defer wg.Done()
//				err := TaikoGalxe1(v, k, chs, wd, url, dstFile)
//				if err != nil {
//					log.Println("!-------!", v.BitId, "失败")
//				}
//				defer bitbrowser.CloseBrower(v.BitId)
//			})
//		}
//		wg.Wait()
//	}
//	close(chs)
//	err := excel.SaveAs(filepout)
//	if err != nil {
//		log.Println("excel 保存失败----", err)
//	}
//}
//func TaikoGalxe1(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
//	wrongData := []string{excelInfo.HelpWords, excelInfo.PrivateKey, excelInfo.PublicKey, excelInfo.Address, excelInfo.Type, excelInfo.BitId, excelInfo.MetaPwd}
//	wd.Get(url)
//	err = wd.Get(url)
//	if err != nil {
//		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("银河打开成功")
//	}
//	time.Sleep(2 * time.Second)
//
//	//小狐狸登陆
//	nowHandle, _ := wd.WindowHandles()
//	if len(nowHandle) > 1 {
//		wd.SwitchWindow(nowHandle[1])
//		err = SmallFoxLoginNoSign(wd)
//	}
//	if err != nil {
//		log.Println("小狐狸登陆出错-----", excelInfo.BitId)
//		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("小狐狸登陆成功")
//		time.Sleep(1 * time.Second)
//		hanNow, _ := wd.WindowHandles()
//		wd.SwitchWindow(hanNow[0])
//		wd.MaximizeWindow(hanNow[0])
//		time.Sleep(1 * time.Second)
//	}
//
//	texts, err := TaikoFindAllDropDownBox(wd, 12)
//	if err != nil {
//		ch <- wrongData
//		log.Println("打开下拉框失败了------", excelInfo.Address)
//		return err
//	} else {
//		log.Println("打开所有下拉框成功")
//		time.Sleep(2 * time.Second)
//	}
//	fmt.Println("text的值：", texts)
//
//	handle, _ := wd.WindowHandles()
//	err = TaikoUrlTask(wd, 1, handle[0], texts)
//	if err != nil {
//		ch <- wrongData
//		log.Println("处理url------", excelInfo.Address)
//		return err
//	} else {
//		log.Println("处理所有url成功")
//		time.Sleep(2 * time.Second)
//	}
//	RefreshALl(wd)
//	time.Sleep(10 * time.Second)
//	return err
//}
//func SmallFoxLoginNoSign(wd selenium.WebDriver) (err error) {
//	var Password selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			Password, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"password\"]")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, err
//	}, 3*time.Second)
//	if err != nil {
//		log.Println("没有找到登陆按钮")
//		return err
//	} else {
//		Password.SendKeys("SHIfeng0615")
//		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
//		if err == nil {
//			UnLock.Click()
//		}
//	}
//	return
//}
//func Authorize(wd selenium.WebDriver) {
//	ShandleNow, _ := wd.WindowHandles()
//	if len(ShandleNow) > 1 {
//		closeSignInPop(wd, ShandleNow[1])
//		time.Sleep(2 * time.Second)
//		wd.SwitchWindow(ShandleNow[0])
//		wd.MaximizeWindow(ShandleNow[0])
//	}
//}
//func authorize(wd selenium.WebDriver, handle string) {
//	wd.SwitchWindow(handle)
//	SignIn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
//	if err == nil {
//		SignIn.Click()
//		log.Println("关闭小狐狸弹窗")
//	}
//}
//func RefreshALl(wd selenium.WebDriver) {
//	refresh, _ := wd.FindElements(selenium.ByCSSSelector, ".clickable.refresh.icon.responsive")
//	log.Println("查找到刷新按钮的数量：", len(refresh))
//	if len(refresh) > 0 {
//		for _, v := range refresh {
//			err := v.Click()
//			if err != nil {
//				log.Println("点击刷新按钮失败", err)
//			} else {
//				log.Println("点击刷新按钮成功")
//			}
//		}
//	} else if refresh1, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(refresh1) > 0 {
//		for _, v := range refresh1 {
//			text, _ := v.Text()
//			if strings.Contains(text, "Verify") {
//				err := v.Click()
//				if err != nil {
//					log.Println("点击刷新按钮失败", err)
//				} else {
//					log.Println("点击刷新按钮成功")
//				}
//			} else {
//				continue
//			}
//		}
//	}
//}
//func CloseSignInPop(wd selenium.WebDriver) {
//	ShandleNow, _ := wd.WindowHandles()
//	if len(ShandleNow) > 1 {
//		closeSignInPop(wd, ShandleNow[1])
//		time.Sleep(2 * time.Second)
//		wd.SwitchWindow(ShandleNow[0])
//		wd.MaximizeWindow(ShandleNow[0])
//	}
//}
//func closeSignInPop(wd selenium.WebDriver, handle string) {
//	wd.SwitchWindow(handle)
//	SignIn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
//	if err == nil {
//		SignIn.Click()
//		log.Println("关闭小狐狸弹窗")
//	}
//}
//
//func TaikoFindAllDropDownBox(wd selenium.WebDriver, num int) (texts []string, err error) {
//	dropdownElements, _ := wd.FindElements(selenium.ByCSSSelector, ".v-expansion-panel")
//	//length := len(dropdownElements) - num
//	//dropLen := len(dropdownElements)
//	log.Println(len(dropdownElements))
//	for k, dropdownElement := range dropdownElements {
//		//跳过一定数量的下拉框
//		//if k < length && dropLen > num {
//		//	continue
//		//}
//		// 判断下拉框状态
//		ariaExpanded, err := dropdownElement.GetAttribute("aria-expanded")
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		if ariaExpanded == "false" {
//			log.Println("开始第", k, "个下拉框的处理")
//			// 点击下拉框头部元素以展开   这里不点击整个div，而是点击里面具体的一个div
//			expand, err := dropdownElement.FindElement(selenium.ByCSSSelector, ".expand-icon")
//			if err == nil {
//				// 等待一段时间以确保下拉框展开完全
//				text, _ := dropdownElement.FindElement(selenium.ByCSSSelector, ".cred-name.usual-text.text-overline-ellipsis-1")
//				textDetail, _ := text.Text()
//				fmt.Println("text的值：", textDetail)
//				texts = append(texts, textDetail)
//				//点击下拉框
//				wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//					for i := 0; i < 5; i++ {
//						err1 := expand.Click()
//						if err1 != nil {
//							time.Sleep(1 * time.Second)
//							continue
//						} else {
//							return true, nil
//						}
//					}
//					return false, errors.New("失败")
//				}, 3*time.Second)
//				if err != nil {
//					log.Println("打开下拉框失败")
//					return nil, err
//				}
//			} else if text, err := dropdownElement.FindElement(selenium.ByCSSSelector, ".text-16-semi-bold.usual-text"); err == nil {
//				textDetail1, _ := text.Text()
//				texts = append(texts, textDetail1)
//				fmt.Println("text的值：", textDetail1)
//				wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//					for i := 0; i < 5; i++ {
//						err1 := dropdownElement.Click()
//						if err1 != nil {
//							time.Sleep(1 * time.Second)
//							continue
//						} else {
//							return true, nil
//						}
//					}
//					return false, errors.New("失败")
//				}, 3*time.Second)
//			} else {
//				continue
//			}
//		}
//	}
//	return
//}
//func TaikoUrlTask(wd selenium.WebDriver, num int, handle string, texts []string) error {
//	aUrl, err := wd.FindElements(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
//	if len(aUrl) == 0 {
//		log.Println("没有找到任务")
//		return nil
//	} else {
//		log.Println("当前页面url总数量", len(aUrl))
//		for k, v := range aUrl {
//			text, _ := v.Text()
//			log.Println("第", k, "个url的text---", text)
//			if strings.Contains(text, "Detail") {
//				log.Println("开始第", k, "个url的点击")
//				err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//					for i := 0; i < 10; i++ {
//						err := v.Click()
//						if err != nil {
//							time.Sleep(1 * time.Second)
//							continue
//						} else {
//							return true, nil
//						}
//					}
//					return false, errors.New("失败")
//				}, 8*time.Second)
//				if err != nil {
//					log.Println("点击第", k, "个url失败：", err)
//					return err
//				} else {
//					time.Sleep(1 * time.Second)
//					//权柄切换
//					handles, _ := wd.WindowHandles()
//					wd.SwitchWindow(handles[1])
//					err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//						for i := 0; i < 10; i++ {
//							_, err = wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
//							if err != nil {
//								time.Sleep(1 * time.Second)
//								continue
//							} else {
//								return true, nil
//							}
//						}
//						return false, errors.New("失败")
//					}, 5*time.Second)
//					if err != nil {
//						log.Println("查找第二重url失败")
//						return err
//					} else {
//						//第一个任务
//						url, _ := wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
//						switch {
//						case strings.Contains(texts[k], "Omni") && strings.Contains(texts[k], "Users"):
//							err1 := Galxe.GalxeFollow(wd, url)
//							//选择主页面并关闭其他页面
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//							time.Sleep(3 * time.Second)
//							if err1 != nil {
//								return err1
//							}
//						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Follow"):
//							log.Println("进入twitter Follow")
//							err1 := twitter.TwitterFollow(wd, url)
//							time.Sleep(3 * time.Second)
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//							if err1 != nil {
//								return err1
//							}
//						case strings.Contains(texts[k], "Tweet Retweeters"):
//							log.Println("进入Retweeters")
//							err1 := twitter.TwitterReweet(wd, url)
//							time.Sleep(3 * time.Second)
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//							if err1 != nil {
//								return err1
//							}
//						case strings.Contains(texts[k], "Tweet Liker"):
//							log.Println("进入Liker")
//							err1 := twitter.TwitterLike(wd, url)
//							time.Sleep(3 * time.Second)
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//							if err1 != nil {
//								return err1
//							}
//						case strings.Contains(texts[k], "Twitter") && strings.Contains(texts[k], "Tweet"):
//							err1 := twitter.TwitterTweet(wd, url)
//							time.Sleep(3 * time.Second)
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//							if err1 != nil {
//								return err1
//							}
//						case strings.Contains(texts[k], "Discord") && strings.Contains(texts[k], "verified"):
//						default:
//							log.Println("进入deFault")
//							text, _ := url.Text()
//							log.Println(text)
//							wd.Get(text)
//							time.Sleep(3 * time.Second)
//							bitbrowser.CloseOtherLabels(wd, handle)
//							wd.SwitchWindow(handle)
//						}
//					}
//				}
//			} else {
//				continue
//			}
//		}
//	}
//	return nil
//}
