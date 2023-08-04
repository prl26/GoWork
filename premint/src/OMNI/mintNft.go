package OMNI

//
//import (
//	"bufio"
//	"errors"
//	"fmt"
//	"github.com/JianLinWei1/premint-selenium/model"
//	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
//	"github.com/JianLinWei1/premint-selenium/src/util"
//	"github.com/JianLinWei1/premint-selenium/src/wdservice"
//	"github.com/go-vgo/robotgo"
//	"github.com/tebeka/selenium"
//	"github.com/xuri/excelize/v2"
//	"log"
//	"os"
//	"os/exec"
//	"strings"
//	"time"
//)
//
//func MintNft() {
//	file, err := os.Open("config.txt")
//	if err != nil {
//		fmt.Println("无法打开文件:", err)
//		return
//	}
//	defer file.Close()
//
//	// 使用bufio包创建一个新的Scanner用于读取文件内容
//	scanner := bufio.NewScanner(file)
//	var lines []string
//	// 循环读取文件的每一行
//	for scanner.Scan() {
//		line := scanner.Text()
//		lines = append(lines, line)
//	}
//	// 检查是否有错误发生在scanner.Scan()过程中
//	if err := scanner.Err(); err != nil {
//		fmt.Println("读取文件时发生错误:", err)
//	}
//	url := lines[0]
//	//workingDir, err := os.Getwd()
//	//url := "https://galxe.com/OmniNetwork/campaign/GCXHgUWBKg"
//	//path1 := fmt.Sprintf("%v\\测试数据-200.xlsx", workingDir)
//	//path2 := fmt.Sprintf("%v\\failInfo.xlsx", workingDir)
//	//path3 := fmt.Sprintf("%v\\failInfo.txt", workingDir)
//	//path4 := fmt.Sprintf("%v\\successInfo.txt", workingDir)
//	path1 := "./resource/测试数据-200.xlsx"
//	path2 := "./resource/failInfo-MintNft.xlsx"
//	path3 := "./resource/failInfo-MintNft.txt"
//	path4 := "./resource/successInfo-MintNft.txt"
//	//excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据-200.xlsx")
//	//filepout := "D:\\GoWork\\resource\\FailInfos\\OmniGalxe1-200.xlsx"
//	//TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\OmniGalxe1-200.txt"
//	excelInfos := util.GetOMNIExcelInfos(path1)
//	filepout := path2
//	TxtfileOut := path3
//	TxtSuccessOut := path4
//	dstFile, err := os.OpenFile(TxtfileOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
//	if err != nil {
//		fmt.Println("无法创建文件:", err)
//		return
//	}
//	defer dstFile.Close()
//	successFile, err := os.OpenFile(TxtSuccessOut, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
//	if err != nil {
//		fmt.Println("无法创建文件:", err)
//		return
//	}
//	defer successFile.Close()
//	chs := make(chan []string, len(excelInfos))
//	var title = []string{"助记词", "私钥", "公钥", "地址", "类型", "窗口ID", "MetaMask密码"}
//	//创建新的excel文件
//	excel := excelize.NewFile()
//	excel.SetSheetRow("Sheet1", "A1", &title)
//	//定义一次开多少线程
//	fmt.Println("数据长度--------", len(excelInfos))
//	describeFilePath := "C:\\describe.txt"
//	nameFilePath := "C:\\name1600.txt"
//	describes := GetDescribeTxt(describeFilePath)
//	names := GetNameTxt(nameFilePath)
//	// 获取内容并写入Excel
//	go func() {
//		for t := 0; t < len(excelInfos); t++ {
//			data := <-chs
//			log.Println("接受到一条错误信息：", data)
//			axis := fmt.Sprintf("A%d", t+2)
//			excel.SetSheetRow("Sheet1", axis, &data)
//		}
//	}()
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
//				filePath := fmt.Sprintf("C:\\img\\img (%v).jpg", k+201)
//				err := mintNft(v, k, chs, wd, url, dstFile, names[k], describes[k], filePath, successFile)
//				if err != nil {
//					log.Println("!-------!", v.BitId, "失败")
//				}
//				defer bitbrowser.CloseBrower(v.BitId)
//			})
//		}
//		wg.Wait()
//	}
//	close(chs)
//	err = excel.SaveAs(filepout)
//	if err != nil {
//		log.Println("excel 保存失败----", err)
//	} else {
//		log.Println("excel 保存成功----")
//
//	}
//}
//func mintNft(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, name string, describe string, filePath string, successFile *os.File) (err error) {
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
//	//开始前保持未登录状态
//	err = reLogin(wd)
//	if err != nil {
//		log.Println(excelInfo.BitId, "开始前保持未登录状态出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("开始前保持未登录状态出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("开始前保持未登录状态成功")
//		time.Sleep(2 * time.Second)
//	}
//	//开始build
//	err = ClickBuild(wd)
//	if err != nil {
//		log.Println(excelInfo.BitId, "点击build出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("点击build出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("点击build成功")
//	}
//
//	//开始点击clickMeatmask
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			_, err = wd.FindElement(selenium.ByCSSSelector, ".h-full.w-full.object-contain.p-5")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 2*time.Second)
//	if err == nil {
//		err = ClickMetamask(wd)
//		if err != nil {
//			log.Println(excelInfo.BitId, "点击连接metamask出错了-----", err)
//			dstFile.WriteString(fmt.Sprintf("点击metamask出错了-----%v---%v\n", excelInfo.BitId, i))
//			ch <- wrongData
//			return err
//		} else {
//			log.Println("点击连接metamask成功")
//			time.Sleep(2 * time.Second)
//		}
//
//		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//			for i := 0; i < 10; i++ {
//				nowHandle, _ := wd.WindowHandles()
//				if len(nowHandle) == 1 {
//					time.Sleep(500 * time.Millisecond)
//					continue
//				} else {
//					return true, nil
//				}
//			}
//			return false, errors.New("Fail")
//		}, 3*time.Second)
//		if err == nil {
//			nowHandle, _ := wd.WindowHandles()
//			wd.SwitchWindow(nowHandle[1])
//			gotit, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//			if len(gotit) > 0 {
//				log.Println("开始点击gotIt")
//				gotit[len(gotit)-1].Click()
//				err = SmallFoxNext(wd)
//				if err != nil {
//					return err
//				}
//			} else {
//				err = SmallFoxLogin(wd)
//				SmallFoxNext(wd)
//				if err != nil {
//					log.Println("小狐狸登陆出错-----", excelInfo.BitId)
//					dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
//					ch <- wrongData
//					return err
//				} else {
//					log.Println("小狐狸登陆成功")
//					time.Sleep(1 * time.Second)
//				}
//			}
//		} else {
//			log.Println("小狐狸登陆出错-----", excelInfo.BitId)
//			dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
//			ch <- wrongData
//			return err
//		}
//	}
//
//	//判断是不是Mint过了
//	hanNow, _ := wd.WindowHandles()
//	wd.SwitchWindow(hanNow[0])
//	err = isMinted(wd)
//	if err == nil {
//		log.Println("你已经Mint过了")
//		return nil
//	}
//
//	wd.ResizeWindow(hanNow[0], 1024, 1080)
//	time.Sleep(1 * time.Second)
//	//点击Mint
//	err = ClickMint(wd)
//	if err != nil {
//		log.Println(excelInfo.BitId, "点击Mint出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("点击Mint出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("点击Mint成功")
//		hanNow, _ := wd.WindowHandles()
//		wd.SwitchWindow(hanNow[0])
//		wd.ResizeWindow(hanNow[0], 1024, 1440)
//		time.Sleep(2 * time.Second)
//	}
//
//	//上传图片
//	err = postImage(wd, filePath)
//	if err != nil {
//		log.Println(excelInfo.BitId, "上传图片出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("上传图片出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("上传图片成功")
//		time.Sleep(2 * time.Second)
//	}
//	//处理页面输入
//	err = SendNameAndDescribition(wd, name, describe)
//	if err != nil {
//		log.Println(excelInfo.BitId, "处理页面输入出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("处理页面输入出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("处理页面输入成功")
//		time.Sleep(2 * time.Second)
//	}
//	//选择网络
//	err = ChooseNetwork(wd)
//	if err != nil {
//		log.Println(excelInfo.BitId, "选择网络出错了-----", err)
//		dstFile.WriteString(fmt.Sprintf("选择网络出错了-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//		return err
//	} else {
//		log.Println("选择网络成功")
//		time.Sleep(2 * time.Second)
//	}
//	submit, err := wd.FindElement(selenium.ByCSSSelector, ".text-xl.text-white")
//	if err == nil {
//		err = submit.Click()
//		if err != nil {
//			log.Println("点击submit失败")
//			dstFile.WriteString(fmt.Sprintf("点击submit失败-----%v---%v\n", excelInfo.BitId, i))
//			ch <- wrongData
//		} else {
//			log.Println("点击submit成功")
//		}
//	} else {
//		log.Println("查找submit失败")
//		dstFile.WriteString(fmt.Sprintf("查找submit失败-----%v---%v\n", excelInfo.BitId, i))
//		ch <- wrongData
//	}
//	time.Sleep(2 * time.Second)
//	//关闭小狐狸弹窗
//	chanNow, _ := wd.WindowHandles()
//	if len(chanNow) > 1 {
//		wd.SwitchWindow(chanNow[1])
//		log.Println("开始关闭小狐狸弹窗")
//		err = closeMetamaskPop(wd)
//		if err != nil {
//			log.Println(excelInfo.BitId, "关闭小狐狸弹窗出错了-----", err)
//			dstFile.WriteString(fmt.Sprintf("关闭小狐狸弹窗出错了-----%v---%v\n", excelInfo.BitId, i))
//			ch <- wrongData
//			return err
//		} else {
//			log.Println("关闭小狐狸弹窗成功")
//			time.Sleep(2 * time.Second)
//		}
//	}
//	//开始等待confirm
//	log.Println("开始45s等待confirm")
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 30; i++ {
//			hanNow1, _ := wd.WindowHandles()
//			if len(hanNow1) > 1 {
//				wd.SwitchWindow(hanNow1[1])
//				_, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
//				if err == nil {
//					return true, nil
//				}
//			} else {
//				time.Sleep(1 * time.Second)
//				continue
//			}
//		}
//		return false, errors.New("Fail")
//	}, 45*time.Second)
//	if err == nil {
//		//点击confirm
//		var confirm selenium.WebElement
//		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//			for i := 0; i < 30; i++ {
//				hanNow1, _ := wd.WindowHandles()
//				if len(hanNow1) > 1 {
//					wd.SwitchWindow(hanNow1[1])
//					confirm, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
//					if err == nil {
//						err = confirm.Click()
//						if err != nil {
//							return true, nil
//						}
//					}
//				} else {
//					time.Sleep(500 * time.Millisecond)
//					continue
//				}
//			}
//			return false, errors.New("Fail")
//		}, 5*time.Second)
//		//如果点击成功
//		if err == nil {
//			log.Println("点击confirm成功")
//			log.Println("confirm成功后再等待判定是否成功")
//			//ml-4 text-lg text-white
//			err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//				for i := 0; i < 30; i++ {
//					hanNow1, _ := wd.WindowHandles()
//					wd.SwitchWindow(hanNow1[0])
//					_, err = wd.FindElement(selenium.ByCSSSelector, ".ml-4.text-lg.text-white")
//					if err == nil {
//						return true, nil
//					} else {
//						time.Sleep(1 * time.Second)
//						continue
//					}
//				}
//				return false, errors.New("Fail")
//			}, 40*time.Second)
//			if err == nil {
//				log.Println("判定成功")
//				successFile.WriteString(fmt.Sprintf("成功-----%v-----%v\n", excelInfo.BitId, i))
//			} else {
//				log.Println("判定失败")
//				dstFile.WriteString(fmt.Sprintf("判定失败-----%v---%v\n", excelInfo.BitId, i))
//				ch <- wrongData
//			}
//			time.Sleep(3 * time.Second)
//		} else {
//			log.Println("点击confirm失败")
//			dstFile.WriteString(fmt.Sprintf("进行confirm失败-----%v---%v\n", excelInfo.BitId, i))
//			ch <- wrongData
//			return err
//		}
//	}
//	return nil
//}
//
//func isMinted(wd selenium.WebDriver) (err error) {
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			elementsWithSameClass, _ := wd.FindElements(selenium.ByCSSSelector, ".mantine-1akf1zi.mantine-Checkbox-input")
//			log.Println("element", len(elementsWithSameClass))
//			if len(elementsWithSameClass) > 0 {
//				script := "return window.getComputedStyle(arguments[0]).background;"
//				clolr, _ := wd.ExecuteScript(script, []interface{}{elementsWithSameClass[0]})
//				clolr1, _ := wd.ExecuteScript(script, []interface{}{elementsWithSameClass[1]})
//				if clolr == clolr1 {
//					return true, nil
//				} else {
//					return false, errors.New("fail")
//				}
//			} else {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			}
//		}
//		return false, errors.New("fail")
//	}, 4*time.Second)
//	return err
//}
//func reLogin(wd selenium.WebDriver) (err error) {
//	_, err = wd.FindElement(selenium.ByCSSSelector, ".text-lg.text-white")
//	if err != nil {
//		log.Println("开始重新login")
//		meta, err := wd.FindElement(selenium.ByCSSSelector, ".flex.w-full.items-center.justify-center")
//		if err != nil {
//			log.Println("没有找到下拉框按钮")
//			return err
//		} else {
//			meta.Click()
//			log.Println("点击下拉框")
//			err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//				for i := 0; i < 10; i++ {
//					loginOut, _ := meta.FindElements(selenium.ByTagName, "p")
//					log.Println("长度", len(loginOut))
//					if len(loginOut) > 0 {
//						err = loginOut[len(loginOut)-1].Click()
//						if err == nil {
//							return true, nil
//						}
//					} else {
//						time.Sleep(500 * time.Millisecond)
//						continue
//					}
//				}
//				return false, errors.New("fail")
//			}, 2*time.Second)
//			if err != nil {
//				return err
//			}
//			err1 := wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//				for i := 0; i < 10; i++ {
//					loginOut1, _ := wd.FindElements(selenium.ByCSSSelector, ".text-lg.text-white")
//					log.Println("长度", len(loginOut1))
//					if len(loginOut1) > 0 {
//						err = loginOut1[len(loginOut1)-1].Click()
//						if err == nil {
//							return true, nil
//						}
//					} else {
//						time.Sleep(500 * time.Millisecond)
//						continue
//					}
//				}
//				return false, errors.New("fail")
//			}, 2*time.Second)
//			if err1 != nil {
//				return err1
//			}
//		}
//	}
//	return nil
//}
//func ClickBuild(wd selenium.WebDriver) (err error) {
//	var build selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			//build, err = wd.FindElement(selenium.ByCSSSelector, ".bg-glass-light-button.dark:bg-glass-dark.hover:dark:bg-glass-dark-hover.flex.cursor-pointer.items-center.justify-center.rounded-full.text-gray-800.shadow-xl.transition-all.duration-300.hover:scale-105.hover:tracking-wider.dark:text-gray-100.gradient-button.h-10.w-48.lg:h-12.2xl:h-14.undefined")
//			build, err = wd.FindElement(selenium.ByCSSSelector, ".text-xl.text-white")
//			if err != nil {
//				time.Sleep(1 * time.Second)
//				continue
//			} else {
//				build.Click()
//				return true, nil
//			}
//		}
//		return false, err
//	}, 2*time.Second)
//	if err == nil {
//		return nil
//	}
//	return err
//}
//func ClickMetamask(wd selenium.WebDriver) (err error) {
//	var build selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			//build, err = wd.FindElement(selenium.ByCSSSelector, ".bg-glass-light.dark:bg-glass-dark.relative.m-4.h-20.w-20.cursor-pointer.overflow-hidden.rounded-xl.lg:h-28.lg:w-28")
//			build, err = wd.FindElement(selenium.ByCSSSelector, ".h-full.w-full.object-contain.p-5")
//			if err != nil {
//				time.Sleep(1 * time.Second)
//				continue
//			} else {
//				build.Click()
//				return true, nil
//			}
//		}
//		return false, err
//	}, 5*time.Second)
//	if err == nil {
//		return nil
//	}
//	return err
//}
//func SmallFoxLogin(wd selenium.WebDriver) (err error) {
//	var Password selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			Password, err = wd.FindElement(selenium.ByID, "password")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, err
//	}, 2*time.Second)
//	if err != nil {
//		log.Println("没有找到登陆按钮")
//		return err
//	} else {
//		Password.SendKeys("SHIfeng0615")
//		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
//		if err == nil {
//			UnLock.Click()
//			err = SmallFoxNext(wd)
//		}
//		gotit, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//		if len(gotit) > 0 {
//			err = gotit[len(gotit)-1].Click()
//		}
//	}
//	return err
//}
//
//func SmallFoxNext(wd selenium.WebDriver) (err error) {
//	var Next selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			hanNow1, _ := wd.WindowHandles()
//			if len(hanNow1) > 1 {
//				wd.SwitchWindow(hanNow1[1])
//			}
//			Next, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//			if err == nil {
//				return true, nil
//			} else {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			}
//		}
//		return false, err
//	},
//		2*time.Second)
//	if err == nil {
//		err1 := Next.Click()
//		if err1 == nil {
//			log.Println("点击Next成功")
//			var Connect selenium.WebElement
//			err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//				for i := 0; i < 10; i++ {
//					Connect, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//					if err == nil {
//						return true, nil
//					} else {
//						time.Sleep(500 * time.Millisecond)
//						continue
//					}
//				}
//				return false, err
//			},
//				2*time.Second)
//			if err == nil {
//				err2 := Connect.Click()
//				if err2 == nil {
//					log.Println("点击Connect成功")
//					return nil
//				} else {
//					log.Println("点击Connect失败", err)
//					return err
//				}
//			}
//		} else {
//			log.Println("点击next失败", err1)
//			return err1
//		}
//	}
//	return err
//}
//func ClickMint(wd selenium.WebDriver) (err error) {
//	log.Println("点击Mint")
//	var build selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			build, err = wd.FindElement(selenium.ByCSSSelector, ".absolute.h-full.w-full.object-cover")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				build.Click()
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 5*time.Second)
//	if err == nil {
//		return nil
//	}
//	return err
//}
//func SendNameAndDescribition(wd selenium.WebDriver, Name string, Describe string) (err error) {
//	var name []selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 5; i++ {
//			name, err = wd.FindElements(selenium.ByTagName, "input")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 2*time.Second)
//	if err == nil {
//		log.Println(len(name))
//		err = name[1].SendKeys(Name)
//		if err == nil {
//			log.Println("发送name成功")
//		} else {
//			log.Println("name发送失败", err)
//			return err
//		}
//	} else {
//		log.Println("Name发送失败", err)
//		return err
//	}
//	var dercribe selenium.WebElement
//
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 5; i++ {
//			dercribe, err = wd.FindElement(selenium.ByTagName, "textarea")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 2*time.Second)
//	if err == nil {
//		err = dercribe.SendKeys(Describe)
//		if err == nil {
//			log.Println("发送describe成功")
//		} else {
//			log.Println("发送describe失败", err)
//			return err
//		}
//	} else {
//		log.Println("发送describe失败", err)
//		return err
//	}
//	return nil
//}
//func ChooseNetwork(wd selenium.WebDriver) (err error) {
//	var name []selenium.WebElement
//	_, err = wd.ExecuteScript("window.scrollTo(0, 0.8*document.body.scrollHeight);", nil)
//	if err == nil {
//		log.Println("滑动到页面底部")
//	}
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 5; i++ {
//			name, err = wd.FindElements(selenium.ByTagName, "input")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 2*time.Second)
//	if err == nil {
//		err = name[len(name)-1].Click()
//		if err != nil {
//			log.Println("点击network失败", err)
//			return err
//		}
//		//140，520
//		robotgo.Move(140, 520)
//		robotgo.Click("left", true)
//		time.Sleep(1 * time.Second)
//		log.Println("点击network成功")
//	}
//	return err
//}
//func closeMetamaskPop(wd selenium.WebDriver) (err error) {
//	var Password selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			Password, err = wd.FindElement(selenium.ByID, "password")
//			gotit, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//			if len(gotit) > 0 {
//				log.Println("开始点击gotIt")
//				gotit[len(gotit)-1].Click()
//			}
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, err
//	}, 2*time.Second)
//	if err == nil {
//		Password.SendKeys("SHIfeng0615")
//		UnLock, err := wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-default")
//		if err == nil {
//			UnLock.Click()
//		}
//	}
//	var approve selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			approve, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//			gotit, _ := wd.FindElements(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//			if len(gotit) > 0 {
//				log.Println("开始点击gotIt")
//				gotit[len(gotit)-1].Click()
//			}
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("Fail")
//	}, 3*time.Second)
//	// 如果找到了approve，就先approve，然后在切换network，再confirm
//	if err == nil {
//		err = approve.Click()
//		if err != nil {
//			log.Println("点击approve失败")
//			return err
//		}
//		log.Println("点击approve成功，开始点击switchNwtwork")
//		time.Sleep(1 * time.Second)
//		var switchNwtwork selenium.WebElement
//		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//			for i := 0; i < 10; i++ {
//				hanNow, _ := wd.WindowHandles()
//				wd.SwitchWindow(hanNow[1])
//				switchNwtwork, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
//				if err != nil {
//					time.Sleep(500 * time.Millisecond)
//					continue
//				} else {
//					return true, nil
//				}
//			}
//			return false, errors.New("Fail")
//		}, 3*time.Second)
//		// 如果找到了approve，就先approve，然后在切换network，再confirm
//		if err == nil {
//			err = switchNwtwork.Click()
//			if err != nil {
//				log.Println("点击switchNwtwork失败")
//				return err
//			}
//			log.Println("点击switchNwtwork成功")
//		} else {
//			log.Println("点击switchNwtwork失败")
//			return err
//		}
//	}
//	return nil
//}
//func GetDescribeTxt(filePath string) (infos []string) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println("文件打开失败 = ", err)
//	}
//	defer file.Close()              // 关闭文本流
//	reader := bufio.NewReader(file) // 读取文本数据
//	for {
//		line, err := reader.ReadString('\n') // 读取直到遇到换行符
//		if err != nil {
//			break // 文件读取完毕或发生错误时退出循环
//		}
//		if strings.Contains(line, " ") {
//			line = line[strings.Index(line, " ")+1:] // 去除序号和空格
//			infos = append(infos, line)              // 打印读取的行内容
//		}
//	}
//	return
//}
//func GetNameTxt(filePath string) (infos []string) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println("文件打开失败 = ", err)
//	}
//	defer file.Close()              // 关闭文本流
//	reader := bufio.NewReader(file) // 读取文本数据
//	for {
//		line, err := reader.ReadString('\n') // 读取直到遇到换行符
//		if err != nil {
//			break // 文件读取完毕或发生错误时退出循环
//		}
//		infos = append(infos, line)
//	}
//	return
//}
//func postImage(wd selenium.WebDriver, path string) (err error) {
//	var imageClick selenium.WebElement
//	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
//		for i := 0; i < 10; i++ {
//			imageClick, err = wd.FindElement(selenium.ByCSSSelector, ".flex.h-44.flex-col.items-center.justify-center")
//			if err != nil {
//				time.Sleep(500 * time.Millisecond)
//				continue
//			} else {
//				return true, nil
//			}
//		}
//		return false, errors.New("fail")
//	}, 3*time.Second)
//	if err == nil {
//		err = imageClick.Click()
//		if err != nil {
//			log.Println("点击上传图片失败")
//			return err
//		} else {
//			fmt.Println("------------开始上传图片")
//			cmd := exec.Command("C:\\post.exe", path)
//			_, err = cmd.CombinedOutput()
//			if err != nil {
//				fmt.Println("post image Error:", err)
//				return err
//			}
//		}
//	} else {
//		log.Println("没有找到上传图片按钮")
//		return err
//	}
//
//	return nil
//}
