package OMNI

import (
	"errors"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/Galxe"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/premint"
	"github.com/JianLinWei1/premint-selenium/src/twitter"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func OmniStart() {
	//file, err := os.Open("config.txt")
	//if err != nil {
	//	fmt.Println("无法打开文件:", err)
	//	return
	//}
	//defer file.Close()
	//
	//// 使用bufio包创建一个新的Scanner用于读取文件内容
	//scanner := bufio.NewScanner(file)
	//var lines []string
	//// 循环读取文件的每一行
	//for scanner.Scan() {
	//	line := scanner.Text()
	//	lines = append(lines, line)
	//}
	//// 检查是否有错误发生在scanner.Scan()过程中
	//if err := scanner.Err(); err != nil {
	//	fmt.Println("读取文件时发生错误:", err)
	//}
	//url := lines[0]
	//isHaveFollow := lines[1]
	//workingDir, err := os.Getwd()
	//path1 := fmt.Sprintf("%v\\测试数据-200.xlsx", workingDir)
	//path2 := fmt.Sprintf("%v\\failInfo.xlsx", workingDir)
	//path3 := fmt.Sprintf("%v\\failInfo.txt", workingDir)
	//path4 := fmt.Sprintf("%v\\successInfo.txt", workingDir)
	//path1 := "./resource/测试数据-200.xlsx"
	//path2 := "./resource/failInfo-OmniStart.xlsx"
	//path3 := "./resource/failInfo-OmniStart.txt"
	//path4 := "./resource/successInfo-OmniStart.txt"
	isHaveFollow := "false"
	url := "https://galxe.com/altlayer/campaign/GC9tiUeiq3"
	//url := "https://galxe.com/OmniNetwork/campaign/GCSrxU7M8K"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-20.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\failInfo-测试数据-20.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\failInfo-测试数据-20.txt"
	TxtSuccessOut := "D:\\GoWork\\resource\\SuccessInfos\\successInfo-测试数据-20.txt"
	//excelInfos := util.GetOMNIExcelInfos(path1)
	//filepout := path2
	//TxtfileOut := path3
	//TxtSuccessOut := path4
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
				err := TaikoGalxe1(v, k, chs, wd, url, dstFile, isHaveFollow, successFile)
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
func TaikoGalxe1(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, isHaveFollow string, successFile *os.File) (err error) {
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

	//小狐狸登陆
	nowHandle, _ := wd.WindowHandles()
	if len(nowHandle) > 1 {
		wd.SwitchWindow(nowHandle[1])
		log.Println(excelInfo.MetaPwd)
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
	//获取所有选择
	handle, _ := wd.WindowHandles()
	err = TaikoUrlTask(wd, 10, handle[0], isHaveFollow)
	if err != nil {
		ch <- wrongData
		log.Println("处理url出错------", excelInfo.Address)
		dstFile.WriteString(fmt.Sprintf("处理url出错-----%v----%v\n", excelInfo.BitId, i))
		return err
	} else {
		log.Println("处理所有url成功")
		time.Sleep(2 * time.Second)
	}
	RefreshALl(wd)
	//如果有anth
	time.Sleep(4 * time.Second)
	handle, _ = wd.WindowHandles()
	if len(handle) > 1 {
		for i := 0; i < len(handle); i++ {
			if i == len(handle)-1 {
				break
			}
			wd.SwitchWindow(handle[len(handle)-1-i])
			button, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-sdzlij.r-1phboty.r-rs99b7.r-19yznuf.r-64el8z.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
			if err != nil {
				return err
			}
			button.Click()
			time.Sleep(4 * time.Second)
		}
	}
	//successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
	//开始Claim,刷新，开始判定领取
	wd.SwitchWindow(handle[0])
	err = wd.Refresh()
	if err != nil {
		log.Println(excelInfo.BitId, "打开银河链接出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("打开银河链接出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("银河打开成功")
	}
	time.Sleep(2 * time.Second)
	nowHandle, _ = wd.WindowHandles()
	if len(nowHandle) > 1 {
		wd.SwitchWindow(nowHandle[1])
		err = SmallFoxLoginNoSign(wd, excelInfo.MetaPwd)
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
		if err != nil {
			log.Println("小狐狸关闭弹窗失败-----", excelInfo.Address)
			ch <- wrongData
			return err
		}
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
			successFile.WriteString(fmt.Sprintf("老版--成功-----%v-----%v\n", excelInfo.BitId, i))
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
			successFile.WriteString(fmt.Sprintf("新版--成功-----%v-----%v\n", excelInfo.BitId, i))
			return nil
		}
	} else if isClaim, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(isClaim) > 0 {
		for _, v := range isClaim {
			text, _ := v.Text()
			if strings.Contains(text, "Claim") {
				log.Println("你已经领取过了")
				successFile.WriteString(fmt.Sprintf("你已经领取过了--成功-----%v-----%v\n", excelInfo.BitId, i))
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
func SmallFoxLoginNoSign(wd selenium.WebDriver, passwd string) (err error) {
	var Password selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			Password, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"password\"]")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, err
	}, 3*time.Second)
	if err != nil {
		log.Println("没有找到登陆按钮")
		return err
	} else {
		Password.SendKeys(passwd)
		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
		if err == nil {
			UnLock.Click()
		}
	}
	return
}

func RefreshALl(wd selenium.WebDriver) {
	refresh, _ := wd.FindElements(selenium.ByCSSSelector, ".clickable.refresh.icon.responsive")
	log.Println("查找到刷新按钮的数量：", len(refresh))
	if len(refresh) > 0 {
		for _, v := range refresh {
			err := v.Click()
			if err != nil {
				log.Println("点击刷新按钮失败", err)
			} else {
				log.Println("点击刷新按钮成功")
			}
		}
	} else if refresh1, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(refresh1) > 0 {
		for _, v := range refresh1 {
			text, _ := v.Text()
			if strings.Contains(text, "Verify") {
				err := v.Click()
				if err != nil {
					log.Println("点击刷新按钮失败", err)
				} else {
					log.Println("点击刷新按钮成功")
				}
			} else {
				continue
			}
		}
	}
	time.Sleep(1 * time.Second)
}
func CloseSignInPop(wd selenium.WebDriver) {
	ShandleNow, _ := wd.WindowHandles()
	if len(ShandleNow) > 1 {
		closeSignInPop(wd, ShandleNow[1])
		time.Sleep(2 * time.Second)
		wd.SwitchWindow(ShandleNow[0])
		wd.MaximizeWindow(ShandleNow[0])
	}
}
func closeSignInPop(wd selenium.WebDriver, handle string) {
	wd.SwitchWindow(handle)
	SignIn, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
	if err == nil {
		SignIn.Click()
		log.Println("关闭小狐狸弹窗")
	}
}

// d-flex height-100 width-100 click-area
// v-btn__content
func TaikoUrlTask(wd selenium.WebDriver, num int, handle string, isHaveFollow string) error {
	Url, _ := wd.FindElements(selenium.ByCSSSelector, ".d-flex.height-100.width-100.click-area")
	if len(Url) > 0 {
		log.Println(len(Url))
		err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				_, err1 := wd.FindElement(selenium.ByCSSSelector, ".text-12-semi-bold.text-white")
				if err1 != nil {
					time.Sleep(1 * time.Second)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 10*time.Second)
		if err != nil {
			log.Println("查找follow失败")
			return err
		}
		followButton, _ := wd.FindElement(selenium.ByCSSSelector, ".text-12-semi-bold.text-white")
		text, _ := followButton.Text()
		if !strings.Contains(text, "Following") {
			followButton.Click()
			log.Println("点击follow成功")
		}
		time.Sleep(1 * time.Second)
		for k, v := range Url {
			if k > 2 {
				_, err = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight/document.body.scrollHeight);", nil)
				log.Println("向下滑动", err)
			}
			text, err := v.FindElement(selenium.ByCSSSelector, ".cred-name.usual-text.text-overline-ellipsis-1")
			if err != nil {
				log.Println("查找url属性失败")
				return err
			} else {
				detail, _ := text.Text()
				err = v.Click()
				if err != nil {
					log.Println("点击任务失败")
					return err
				}
				err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < num; i++ {
						handle1, _ := wd.WindowHandles()
						if len(handle1) > 1 {
							return true, nil
						} else {
							time.Sleep(500 * time.Millisecond)
							continue
						}
					}
					return false, err
				}, 2*time.Second)
				if err != nil {
					log.Println("第", k, "条任务已经执行过了")
				} else {
					switch {
					case strings.Contains(detail, "Omni") && strings.Contains(detail, "Users"):
						err1 := Galxe.GalxeFollow(wd)
						//选择主页面并关闭其他页面
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						time.Sleep(2 * time.Second)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Twitter") && strings.Contains(detail, "Follow"):
						log.Println("进入twitter Follow")
						v.Click()
						err1 := twitter.TwitterFollow(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Tweet Retweeters"):
						log.Println("进入Retweeters")
						v.Click()
						err1 := twitter.TwitterReweet(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Tweet Liker"):
						log.Println("进入Liker")
						v.Click()
						err1 := twitter.TwitterLike(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Twitter") && strings.Contains(detail, "Tweet"):
						err1 := twitter.TwitterTweet(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					default:
						log.Println("进入deFault")
						v.Click()
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
					}
				}
			}
		}
	} else if Url, _ := wd.FindElements(selenium.ByCSSSelector, ".v-btn__content"); len(Url) > 0 {
		log.Println("老版本，url长度", len(Url))
		length := len(Url) - 2
		//taskNum := length / 2
		//log.Println("计算出的taskNum：", taskNum)
		for i := 0; i < length; i++ {
			v := Url[len(Url)-1-i]
			detail, err := v.Text()
			if err != nil {
				log.Println("查看url属性失败")
			} else {
				if !strings.Contains(detail, "Verify") {
					switch {
					case strings.Contains(detail, "Go"):
						log.Println("进入go")
						v.Click()
						err1 := Galxe.GalxeVisit(wd)
						//选择主页面并关闭其他页面
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						time.Sleep(2 * time.Second)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Follow Now"):
						log.Println("进入twitter Follow")
						v.Click()
						err1 := twitter.TwitterFollow(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Retweet Now"):
						log.Println("进入Retweeters")
						v.Click()
						err1 := twitter.TwitterReweet(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Like Now"):
						log.Println("进入Liker")
						v.Click()
						err1 := twitter.TwitterLike(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					case strings.Contains(detail, "Twitter") && strings.Contains(detail, "Tweet"):
						v.Click()
						err1 := twitter.TwitterTweet(wd)
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
						if err1 != nil {
							return err1
						}
					default:
						log.Println("进入deFault")
						v.Click()
						time.Sleep(2 * time.Second)
						bitbrowser.CloseOtherLabels(wd, handle)
						wd.SwitchWindow(handle)
					}
				}
			}
		}
		if isHaveFollow == "true" {
			err := FindAllDropDownBox(wd, 1)
			if err != nil {
				log.Println("omni follow失败")
			}
		}
	} else {
		log.Println("查找任务失败")
		return errors.New("Fail")
	}

	return nil
}
func FindAllDropDownBox(wd selenium.WebDriver, num int) error {
	dropdownElements, _ := wd.FindElements(selenium.ByCSSSelector, ".v-expansion-panel")
	log.Println(len(dropdownElements))
	for k, dropdownElement := range dropdownElements {
		if k == 7 {
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
					aUrl, err := dropdownElement.FindElement(selenium.ByCSSSelector, ".detail-text.text-14-regular.clickable")
					if err != nil {
						log.Println("点击第一重url失败")
					} else {
						aUrl.Click()
						time.Sleep(1 * time.Second)
						var button selenium.WebElement
						err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
							for i := 0; i < 10; i++ {
								handle, _ := wd.WindowHandles()
								wd.SwitchWindow(handle[len(handle)-1])
								button, err = wd.FindElement(selenium.ByCSSSelector, "a.credential-link")
								if err != nil {
									time.Sleep(1 * time.Second)
									continue
								} else {
									return true, nil
								}
							}
							return false, errors.New("失败")
						}, 10*time.Second)
						button.Click()
						err1 := Galxe.GalxeFollow(wd)
						//选择主页面并关闭其他页面
						handle, _ := wd.WindowHandles()
						bitbrowser.CloseOtherLabels(wd, handle[0])
						wd.SwitchWindow(handle[0])
						time.Sleep(2 * time.Second)
						if err1 != nil {
							return err1
						}
					}

				}
			}
		} else {
			continue
		}
	}
	return nil
}
