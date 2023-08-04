package Galxe

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

func Bind() {
	url := "https://galxe.com/accountSetting?tab=SocialLinlk"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\bind-7.24.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\bind-7.24.txt"
	dstFile, _ := os.Create(TxtfileOut)
	defer dstFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"地址", "类型", "窗口ID", "MetaMask密码"}
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
				err := TwitterBind(v, k, chs, wd, url, dstFile)
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

func TwitterBind(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File) (err error) {
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
		err = SmallFoxLoginNoSign(wd)
	}
	if err != nil {
		log.Println("小狐狸登陆出错-----", excelInfo.BitId)
		dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("小狐狸登陆成功")
		time.Sleep(2 * time.Second)
		hanNow, _ := wd.WindowHandles()
		wd.SwitchWindow(hanNow[0])
		_, err = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight/document.body.scrollHeight);", nil)
		if err == nil {
			log.Println("点击connect时向上滑动一下")
		}
		time.Sleep(1 * time.Second)
	}
	//如果还有弹窗
	CloseSignInPop(wd)
	count, err := ClickConnectTwiier(wd)
	if err != nil {
		log.Println("点击Connect Twitter失败")
		dstFile.WriteString(fmt.Sprintf("点击Connect Twitter失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("点击Connect Twitter成功")
		time.Sleep(2 * time.Second)
	}
	if count == 1 {
		log.Println("点击Connect Twitter失败")
		dstFile.WriteString(fmt.Sprintf("点击Connect Twitter失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	}

	err = TweetTwitter(wd)
	if err != nil {
		log.Println("Tweet Twitter失败")
		dstFile.WriteString(fmt.Sprintf("Tweet Twitter失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("Tweet Twitter成功")
		hanNow, _ := wd.WindowHandles()
		wd.SwitchWindow(hanNow[1])
		wd.MaximizeWindow(hanNow[1])
		time.Sleep(2 * time.Second)
	}
	//先复制个空的
	clipboard.WriteAll("")
	err = GetInProfile(wd)
	if err != nil {
		log.Println("获取Link失败")
		dstFile.WriteString(fmt.Sprintf("获取Link-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("获取Link成功")
		time.Sleep(2 * time.Second)
		handowNow, _ := wd.WindowHandles()
		wd.SwitchWindow(handowNow[0])
	}
	err = Sendkey(wd)
	if err != nil {
		log.Println("Verify失败")
		dstFile.WriteString(fmt.Sprintf("Verify失败-----%v-----%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("Verify成功")
	}
	handowNow, _ := wd.WindowHandles()
	wd.SwitchWindow(handowNow[0])
	wd.MaximizeWindow(handowNow[0])
	time.Sleep(1 * time.Second)
	err = OpenHomePage(wd)
	if err != nil {
		log.Println("进入个人主页失败")
		dstFile.WriteString(fmt.Sprintf("进入个人主页失败-----%v-----%v", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("进入个人主页成功")
	}
	time.Sleep(1 * time.Second)
	//进入setting
	err = ClickSetting(wd)
	if err != nil {
		log.Println("进入setting失败")
		dstFile.WriteString(fmt.Sprintf("进入setting失败-----%v-----%v", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("进入setting成功")
		time.Sleep(1 * time.Second)

	}
	err = ClickSocialLink(wd, handowNow[0])
	if err != nil {
		log.Println("进入ClickSocialLink失败")
		dstFile.WriteString(fmt.Sprintf("进入ClickSocialLink失败-----%v-----%v", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("进入ClickSocialLink成功")
	}
	time.Sleep(2 * time.Second)
	return err
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
func ClickConnectTwiier(wd selenium.WebDriver) (count int, err error) {
	count = 0
	var SocialLinks []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			SocialLinks, err = wd.FindElements(selenium.ByCSSSelector, ".social-account-link")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 10*time.Second)
	if err != nil {
		return
	} else {
		//找到具有setting属性的div
		for _, v := range SocialLinks {
			linktext, _ := v.FindElement(selenium.ByCSSSelector, ".social-account-link-text")
			text, _ := linktext.Text()
			if strings.Contains(text, "Connect Twitter Account") {
				return 1, nil
			} else {
				continue
			}
		}
		return 0, errors.New("Fail")
	}
	return 0, errors.New("Fail")
}
func SmallFoxLoginNoSign(wd selenium.WebDriver) (err error) {
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
		Password.SendKeys("SHIfeng0615")
		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
		if err == nil {
			UnLock.Click()
		}
	}
	return
}
func TweetTwitter(wd selenium.WebDriver) (err error) {
	//开始发推特
	TweetButton, err := wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
	if err == nil {
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				TweetButton, err = wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
				if len(TweetButton) == 0 {
					time.Sleep(1 * time.Millisecond)
					continue
				} else {
					log.Println(len(TweetButton))
					err = TweetButton[len(TweetButton)-2].Click()
					if err != nil {
						handle, _ := wd.WindowHandles()
						if len(handle) == 1 {
							err1 := wd.Refresh()
							log.Println("刷新", err1)
							time.Sleep(1 * time.Second)
							continue
						}
					}
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 8*time.Second)
		if err != nil {
			log.Println("进入tweet失败", err)
			return err
		} else {
			log.Println("进入tweet成功")
		}
		//<span class="css-901oao css-16my406 r-poiln3 r-bcqeeo r-qvutc0">Tweet</span>
		var Tweet []selenium.WebElement
		log.Println("寻找Tweet按钮")
		time.Sleep(1 * time.Second)
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				handle, _ := wd.WindowHandles()
				if len(handle) > 1 {
					wd.SwitchWindow(handle[1])
				}
				Tweet, _ = wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
				if len(Tweet) == 0 {
					time.Sleep(1 * time.Second)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 4*time.Second)
		if err != nil {
			log.Println("没有找到Tweet按钮")
			return err
		} else {
			for _, v := range Tweet {
				text, _ := v.Text()
				if strings.Contains(text, "Twitter Circle") {
					_, err = wd.ExecuteScript("window.scrollTo(0, 1000);", nil)
					if err != nil {
						fmt.Printf("无法滚动页面：%v\n", err)
						return
					}
				}
				if strings.Contains(text, "Tweet") && !strings.Contains(text, "Twitter Circle") {
					log.Println("开始推特发文", text)
					err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
						for i := 0; i < 10; i++ {
							time.Sleep(1 * time.Second)
							err = v.Click()
							if err != nil {
								_, _ = wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight/document.body.scrollHeight);", nil)
								time.Sleep(1 * time.Second)
								continue
							} else {
								return true, nil
							}
						}
						return false, errors.New("Fail")
					}, 6*time.Second)
					if err == nil {
						time.Sleep(1 * time.Second)
						GotIt, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-1mnahxq.r-19yznuf.r-64el8z.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
						if err == nil {
							err = GotIt.Click()
							if err == nil {
								log.Println("转发推特成功")
								return nil
							}
						} else {
							Close, err := wd.FindElement(selenium.ByCSSSelector, ".r-18jsvk2.r-4qtqp9.r-yyyyoo.r-z80fyv.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-19wmn03")
							if err == nil {
								err = Close.Click()
								if err == nil {
									log.Println("转发推特成功-----你已经发过了")
									return nil
								}
							}
						}
					} else {
						log.Println("tweet的时候出错了")
					}
				}
			}
		}
	} else {
		log.Println("点击Tweet1失败")
		return err
	}
	return err
}

func GetInProfile(wd selenium.WebDriver) (err error) {
	//*[@id="id__f6eimqzspac"]/div[5]/div/div/div
	//*[@id="react-root"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[9]
	//*[@id="react-root"]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[8]
	//*[@id="id__71tlehdomw7"]/div[5]/div/div/div/div/div
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
		}

		//获取推特链接
		time.Sleep(2 * time.Second)
		var articles []selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				articles, _ = wd.FindElements(selenium.ByTagName, "article")
				if len(articles) == 0 {
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 2*time.Second)

		log.Println("查找到的article数量:", len(articles))
		if len(articles) > 0 {
			LinkButtons1, _ := articles[0].FindElement(selenium.ByCSSSelector, ".css-1dbjc4n.r-1kbdv8c.r-18u37iz.r-1wtj0ep.r-1s2bzr4.r-1mdbhws")
			LinkButtons, _ := LinkButtons1.FindElements(selenium.ByCSSSelector, ".r-4qtqp9.r-yyyyoo.r-1xvli5t.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-1hdv0qi")
			log.Println("查找到的LinkButtons数量:", len(LinkButtons))
			time.Sleep(1 * time.Second)
			if len(LinkButtons) > 0 {
				err = LinkButtons[len(LinkButtons)-1].Click()
				if err != nil {
					log.Println("点击获取link链接按钮失败")
				} else {
					//这里会出现四个弹窗,点击第一个
					time.Sleep(1 * time.Second)
					Links, _ := wd.FindElements(selenium.ByCSSSelector, ".css-901oao.r-18jsvk2.r-37j5jr.r-a023e6.r-b88u0q.r-rjixqe.r-bcqeeo.r-qvutc0")
					if len(Links) > 0 {
						err = Links[0].Click()
						if err != nil {
							log.Println("获取链接失败")
						} else {
							return nil
						}
					}
				}
			}
		} else {
			return errors.New("fail")
		}
	}
	//err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
	//	for i := 0; i < 10; i++ {
	//		LinkButton, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"id__sqguip3j0f\"]/div[5]/div")
	//		if err != nil {
	//			time.Sleep(500 * time.Millisecond)
	//			continue
	//		} else {
	//			return true, nil
	//		}
	//	}
	//	return false, errors.New("Fail")
	//}, 2*time.Second)
	//if err != nil {
	//	log.Println("没有找到链接按钮")
	//} else {
	//}
	return err
}
func Sendkey(wd selenium.WebDriver) (err error) {
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		log.Println("读取剪贴板内容失败：", err)
	} else {
		log.Println(clipboardContent)
		time.Sleep(3 * time.Second)
		Input, err := wd.FindElement(selenium.ByTagName, "input")
		if err != nil {
			return err
		} else {
			err = Input.SendKeys(clipboardContent)
			if err != nil {
				log.Println("粘贴链接失败", err)
				return err
			} else {
				time.Sleep(1 * time.Second)
				Verifys, _ := wd.FindElements(selenium.ByCSSSelector, ".tc-tweet-button")
				if len(Verifys) > 0 {
					err = Verifys[len(Verifys)-1].Click()
					if err == nil {
						return nil
					} else {
						return err
					}
				}
			}
		}
	}

	return
}
