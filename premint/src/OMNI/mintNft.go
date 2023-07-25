package OMNI

import (
	"bufio"
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
	"os/exec"
	"strings"
	"time"
)

func MintNft() {
	url := "https://dripverse.org/"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-425.xlsx")
	filepout := "D:\\GoWork\\resource\\FailInfos\\测试点击.xlsx"
	TxtfileOut := "D:\\GoWork\\resource\\FailInfos\\测试点击.txt"
	dstFile, _ := os.Create(TxtfileOut)
	defer dstFile.Close()
	chs := make(chan []string, len(excelInfos))
	var title = []string{"助记词", "私钥", "公钥", "地址", "类型", "窗口ID", "MetaMask密码"}
	//创建新的excel文件
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &title)
	//定义一次开多少线程
	fmt.Println("数据长度--------", len(excelInfos))
	describeFilePath := "C:\\describe.txt"
	nameFilePath := "C:\\name1600.txt"
	describes := getDescribeTxt(describeFilePath)
	names := getNameTxt(nameFilePath)
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
				filePath := fmt.Sprintf("C:\\img\\img (%v).jpg", k+1)
				err := mintNft(v, k, chs, wd, url, dstFile, names[k], describes[k], filePath)
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
	} else {
		log.Println("excel 保存成功----")

	}
}
func mintNft(excelInfo model.OMNIExcelInfo, i int, ch chan<- []string, wd selenium.WebDriver, url string, dstFile *os.File, name string, describe string, filePath string) (err error) {
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

	err = ClickBuild(wd)
	if err != nil {
		log.Println(excelInfo.BitId, "点击build出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("点击build出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("点击build成功")
	}
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			//build, err = wd.FindElement(selenium.ByCSSSelector, ".bg-glass-light.dark:bg-glass-dark.relative.m-4.h-20.w-20.cursor-pointer.overflow-hidden.rounded-xl.lg:h-28.lg:w-28")
			_, err = wd.FindElement(selenium.ByCSSSelector, ".h-full.w-full.object-contain.p-5")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 2*time.Second)
	if err == nil {
		log.Println("进入到这里了")
		err = ClickMetamask(wd)
		if err != nil {
			log.Println(excelInfo.BitId, "点击连接metamask出错了-----", err)
			dstFile.WriteString(fmt.Sprintf("点击metamask出错了-----%v---%v\n", excelInfo.BitId, i))
			ch <- wrongData
			return err
		} else {
			log.Println("点击连接metamask成功")
			time.Sleep(2 * time.Second)
		}
		nowHandle, _ := wd.WindowHandles()
		if len(nowHandle) > 1 {
			wd.SwitchWindow(nowHandle[1])
			err = SmallFoxLoginAndNext(wd)
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
				wd.ResizeWindow(hanNow[0], 1024, 1080)
				time.Sleep(1 * time.Second)
			}
		} else {
			log.Println("小狐狸登陆出错-----", excelInfo.BitId)
			dstFile.WriteString(fmt.Sprintf("小狐狸登陆出错-----%v----%v\n", excelInfo.BitId, i))
			ch <- wrongData
			return err
		}
		//第二次点击build
		//err = ClickBuild(wd)
		//if err != nil {
		//	log.Println(excelInfo.BitId, "第二次点击build出错了-----", err)
		//	dstFile.WriteString(fmt.Sprintf("第二次点击build出错了-----%v---%v\n", excelInfo.BitId, i))
		//	ch <- wrongData
		//	return err
		//} else {
		//	log.Println("第二次点击build成功")
		//	time.Sleep(2 * time.Second)
		//}
	}
	err = ClickMint(wd)
	if err != nil {
		log.Println(excelInfo.BitId, "点击Mint出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("点击Mint出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("点击Mint成功")
		hanNow, _ := wd.WindowHandles()
		wd.SwitchWindow(hanNow[0])
		wd.ResizeWindow(hanNow[0], 1024, 1440)
		time.Sleep(2 * time.Second)
	}

	//上传图片
	err = SendNameAndDescribition(wd)
	if err != nil {
		log.Println(excelInfo.BitId, "上传图片出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("上传图片出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("上传图片成功")
		time.Sleep(2 * time.Second)
	}
	//处理页面输入
	err = SendNameAndDescribition(wd)
	if err != nil {
		log.Println(excelInfo.BitId, "处理页面输入出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("处理页面输入出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("处理页面输入成功")
		time.Sleep(2 * time.Second)
	}
	//选择网络
	err = ChooseNetwork(wd)
	if err != nil {
		log.Println(excelInfo.BitId, "选择网络出错了-----", err)
		dstFile.WriteString(fmt.Sprintf("选择网络出错了-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
		return err
	} else {
		log.Println("选择网络成功")
		time.Sleep(2 * time.Second)
	}
	submit, err := wd.FindElement(selenium.ByCSSSelector, ".text-xl.text-white")
	if err == nil {
		err = submit.Click()
		if err != nil {
			log.Println("点击submit失败")
			dstFile.WriteString(fmt.Sprintf("点击submit失败-----%v---%v\n", excelInfo.BitId, i))
			ch <- wrongData
		}
	} else {
		log.Println("查找submit失败")
		dstFile.WriteString(fmt.Sprintf("查找submit失败-----%v---%v\n", excelInfo.BitId, i))
		ch <- wrongData
	}
	//confirm
	//button btn--rounded btn-primary page-container__footer-button
	//button btn--rounded btn-primary
	time.Sleep(2 * time.Second)

	//关闭小狐狸弹窗

	return err
}

func ClickBuild(wd selenium.WebDriver) (err error) {
	var build selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			//build, err = wd.FindElement(selenium.ByCSSSelector, ".bg-glass-light-button.dark:bg-glass-dark.hover:dark:bg-glass-dark-hover.flex.cursor-pointer.items-center.justify-center.rounded-full.text-gray-800.shadow-xl.transition-all.duration-300.hover:scale-105.hover:tracking-wider.dark:text-gray-100.gradient-button.h-10.w-48.lg:h-12.2xl:h-14.undefined")
			build, err = wd.FindElement(selenium.ByCSSSelector, ".text-xl.text-white")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				build.Click()
				return true, nil
			}
		}
		return false, err
	}, 2*time.Second)
	if err == nil {
		return nil
	}
	return err
}
func ClickMetamask(wd selenium.WebDriver) (err error) {
	var build selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			//build, err = wd.FindElement(selenium.ByCSSSelector, ".bg-glass-light.dark:bg-glass-dark.relative.m-4.h-20.w-20.cursor-pointer.overflow-hidden.rounded-xl.lg:h-28.lg:w-28")
			build, err = wd.FindElement(selenium.ByCSSSelector, ".h-full.w-full.object-contain.p-5")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			} else {
				build.Click()
				return true, nil
			}
		}
		return false, err
	}, 5*time.Second)
	if err == nil {
		return nil
	}
	return err
}
func SmallFoxLoginAndNext(wd selenium.WebDriver) (err error) {
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
		var Next selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				Next, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
				if err == nil {
					return true, nil
				} else {
					time.Sleep(500 * time.Millisecond)
					continue
				}
			}
			return false, err
		},
			2*time.Second)
		if err == nil {
			err1 := Next.Click()
			if err1 == nil {
				log.Println("点击Next成功")
				var Connect selenium.WebElement
				err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
					for i := 0; i < 10; i++ {
						Connect, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
						if err == nil {
							return true, nil
						} else {
							time.Sleep(500 * time.Millisecond)
							continue
						}
					}
					return false, err
				},
					2*time.Second)
				if err == nil {
					err2 := Connect.Click()
					if err2 == nil {
						log.Println("点击Connect成功")
						return nil
					} else {
						log.Println("点击Connect失败", err)
						return err
					}
				}
			} else {
				log.Println("点击next失败", err)
				return err
			}
		}
	}
	return
}
func ClickMint(wd selenium.WebDriver) (err error) {
	log.Println("点击Mint")
	var build selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			build, err = wd.FindElement(selenium.ByCSSSelector, ".absolute.h-full.w-full.object-cover")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				build.Click()
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 5*time.Second)
	if err == nil {
		return nil
	}
	return err
}
func SendNameAndDescribition(wd selenium.WebDriver) (err error) {
	var name []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 5; i++ {
			name, err = wd.FindElements(selenium.ByTagName, "input")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 2*time.Second)
	if err == nil {
		log.Println(len(name))
		err = name[1].SendKeys("test")
		if err == nil {
			log.Println("发送name成功")
		} else {
			log.Println("name发送失败", err)
			return err
		}
	} else {
		log.Println("Name发送失败", err)
		return err
	}
	var dercribe selenium.WebElement

	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 5; i++ {
			dercribe, err = wd.FindElement(selenium.ByTagName, "textarea")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 2*time.Second)
	if err == nil {
		err = dercribe.SendKeys("test")
		if err == nil {
			log.Println("发送describe成功")
		} else {
			log.Println("发送describe失败", err)
			return err
		}
	} else {
		log.Println("发送describe失败", err)
		return err
	}
	return nil
}
func ChooseNetwork(wd selenium.WebDriver) (err error) {
	var name []selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 5; i++ {
			name, err = wd.FindElements(selenium.ByTagName, "input")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 2*time.Second)
	if err == nil {
		log.Println(len(name))
		script := fmt.Sprintf(`arguments[0].value = '%s'`, "Omni Testnet")
		_, err = wd.ExecuteScript(script, []interface{}{name[len(name)-1]})

		//err = name[len(name)-1].SendKeys("Omni Testnet")
		if err == nil {
			log.Println("点击network成功")
		} else {
			log.Println("点击network失败", err)
			return err
		}
	}
	return err
}
func closeMetamaskPop(wd selenium.WebDriver) (err error) {
	//button btn--rounded btn-primary   approve
	//button btn--rounded btn-primary  switch network
	// confirm
	//confirm,err = wd.FindElement(selenium.ByCSSSelector,".button.btn--rounded.btn-primary.page-container__footer-button")
	var approve selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			approve, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 4*time.Second)
	// 如果找到了approve，就先approve，然后在切换network，再confirm
	if err == nil {
		err = approve.Click()
		if err != nil {
			log.Println("点击approve失败")
			return err
		}
		log.Println("点击approve成功，开始点击switchNwtwork")
		var switchNwtwork selenium.WebElement
		err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			for i := 0; i < 10; i++ {
				switchNwtwork, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary")
				if err != nil {
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					return true, nil
				}
			}
			return false, errors.New("Fail")
		}, 4*time.Second)
		// 如果找到了approve，就先approve，然后在切换network，再confirm
		if err == nil {
			err = switchNwtwork.Click()
			if err != nil {
				log.Println("点击switchNwtwork失败")
				return err
			}
			log.Println("点击switchNwtwork成功")
		} else {
			log.Println("点击switchNwtwork失败")
			return err
		}
	}
	var confirm selenium.WebElement
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		for i := 0; i < 10; i++ {
			confirm, err = wd.FindElement(selenium.ByCSSSelector, ".button.btn--rounded.btn-primary.page-container__footer-button")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				return true, nil
			}
		}
		return false, errors.New("Fail")
	}, 4*time.Second)
	// 如果找到了approve，就先approve，然后在切换network，再confirm
	if err == nil {
		err = confirm.Click()
		if err != nil {
			log.Println("点击confirm失败")
			return err
		}
		log.Println("点击confirm成功")
	} else {
		log.Println("点击confirm失败")
		return err
	}
	return err
}
func getDescribeTxt(filePath string) (infos []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败 = ", err)
	}
	defer file.Close()              // 关闭文本流
	reader := bufio.NewReader(file) // 读取文本数据
	for {
		line, err := reader.ReadString('\n') // 读取直到遇到换行符
		if err != nil {
			break // 文件读取完毕或发生错误时退出循环
		}
		if strings.Contains(line, " ") {
			line = line[strings.Index(line, " ")+1:] // 去除序号和空格
			infos = append(infos, line)              // 打印读取的行内容
		}
	}
	return
}
func getNameTxt(filePath string) (infos []string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败 = ", err)
	}
	defer file.Close()              // 关闭文本流
	reader := bufio.NewReader(file) // 读取文本数据
	for {
		line, err := reader.ReadString('\n') // 读取直到遇到换行符
		if err != nil {
			break // 文件读取完毕或发生错误时退出循环
		}
		infos = append(infos, line)
	}
	return
}
func postImage(path string) error {
	fmt.Println("------------")
	cmd := exec.Command("C:\\post.exe", path)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("post image Error:", err)
		return err
	}
	return nil
}
