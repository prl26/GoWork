package premint

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"log"
	"strconv"
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
	//metamask.MetaMaskLogin(wd, excelInfo.MetaPwd)
	//time.Sleep(1 * time.Second)
	//
	//log.Println("打开银河链接")
	//err := wd.Get("https://galxe.com/EchoDEX/campaign/GCDsmUSvqd")
	//if err != nil {
	//	log.Println("打开银河链接出错了")
	//} else {
	//	log.Println("银河打开成功")
	//
	//}
	//
	//handle := util.GetCurrentWindowAndReturn(wd)
	//time.Sleep(5 * time.Second)
	////关闭多余标签页
	//bitbrowser.CloseOtherLabels(wd, handle)

	//打开页面上所有下拉框
	//*[@id="topNavbar"]/div/div[2]/div[2]/div[1]/div[1]/div[1]/div/div/span
	button, err := wd.FindElement(selenium.ByCSSSelector, ".text-14-regular.text-overline-ellipsis-1")

	if err != nil {
		log.Println("查找元素出错了")
	} else {
		//fmt.Println(button.Text())
		time.Sleep(1 * time.Second)

		err = button.Click()
		if err != nil {
			log.Println("点击失败")

		}
		log.Println("查找成功了")

	}
	time.Sleep(1 * time.Second)
	d1, err := wd.FindElements(selenium.ByCSSSelector, ".flex-auto")
	if err != nil {
		log.Println("查找元素出错了")
	} else {
		time.Sleep(1 * time.Second)
		fmt.Println(d1[0].Text())
		d1[0].Click()
		log.Println("查找成功了")
	}

	wg.Done()
	fmt.Println("处理完毕")

}
