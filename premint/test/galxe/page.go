package galxe

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Remove() {
	url := "https://galxe.com/twitterConnect"
	excelInfos := util.GetOMNIExcelInfos("D:\\GoWork\\resource\\测试数据1-425.xlsx")

	fmt.Println("数据长度--------", len(excelInfos))

	for k, v := range excelInfos {
		fmt.Println("----------", v.Address)
		//打开比特浏览器
		wd, _ := wdservice.InitWd(k, v.BitId)
		if wd != nil {
			wg.Add(1)
			go util.SetLog(func() {
				defer wg.Done()
				testPage(wd, url)
				//bitbrowser.CloseBrower(v.BitId)
			})
		}
		//bitbrowser.WindowboundsByPara()
		wg.Wait()
	}

}
func testPage(wd selenium.WebDriver, url string) {
	wd.Get(url)

	time.Sleep(2 * time.Second)
	button, _ := wd.FindElement(selenium.ByCSSSelector, ".tc-tweet-button")
	button.Click()

	time.Sleep(1 * time.Second)
	handle, _ := wd.WindowHandles()
	log.Println("当前页面的标签页数量；；---", len(handle))

	time.Sleep(3 * time.Second)
}
