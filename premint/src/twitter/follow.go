/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-06 09:36:24
 * @LastEditTime: 2023-05-09 10:18:23
 */
package twitter

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"github.com/JianLinWei1/premint-selenium/src/util"
	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
)

func StartFllowAndTweet() {
	eis := getExcelInfos[ExcelInfoTweet]("./Twitter操作模板.xlsx")

	chs := make(chan string, len(eis))
	for i, ei := range eis {
		if ei.BitId == "" {
			chs <- "跳过"
			continue
		}
		time.Sleep(1 + time.Second)
		wd, _ := wdservice.InitWd(i, ei.BitId)
		if wd != nil {
			go util.SetLog(func() { toFllowAndTweet(chs, i, ei, wd) })
		} else {
			log.Println("第" + strconv.Itoa(i+1) + "条数据浏览器初始化失败****")
		}

	}
	bitbrowser.Windowbounds()
	for i := 0; i < len(eis); i++ {
		log.Println(<-chs)
	}

}
func toFllowAndTweet(ch chan string, i int, ei ExcelInfoTweet, wd selenium.WebDriver) {
	wd.Get(ei.Link)
	if ei.Tp == 1 {
		findFollowBtnAndClick(wd)
	} else if ei.Tp == 2 {
		Retweet(wd)
	}

	ch <- "第" + strconv.Itoa(i) + "条数据执行完成"
}

func FollowAndTweeTwitter(wd selenium.WebDriver) {
	log.Println("*********Twitter一系列操作********")
	followAndTweetBtns, err := wd.FindElements(selenium.ByCSSSelector, ".mb-2.border.bg-muted.p-2.rounded.text-md")
	if err != nil {
		log.Println(err)
	}
	if len(followAndTweetBtns) <= 0 {
		return
	}
	//第一个去follow
	FollowTwitter(wd, followAndTweetBtns[0])
	if len(followAndTweetBtns) >= 2 {
		//第二个tweet
		TweetTwitter(wd, followAndTweetBtns[1])
	}

}

func FollowTwitter(wd selenium.WebDriver, div selenium.WebElement) {
	//当前页面handle
	hanle, _ := wd.CurrentWindowHandle()
	follows, err := div.FindElements(selenium.ByTagName, "a")
	if err != nil {
		log.Println(err)
	}
	for _, follwDiv := range follows {
		wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
			//点击
			follwDiv.Click()
			time.Sleep(1 * time.Second)
			//util.GetWind(1, wd)
			util.GetLastWindow(wd)
			findFollowBtnAndClick(wd)
			//wd.Close()
			wd.SwitchWindow(hanle)
			return true, err
		}, 30*time.Second)
	}
	wd.SwitchWindow(hanle)

}

func findFollowBtnAndClick(wd selenium.WebDriver) {
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		title, _ := wd.Title()
		log.Println("当前Tiltle:", title)
		closeSvg, _ := wd.FindElement(selenium.ByCSSSelector, ".r-jwli3a.r-4qtqp9.r-yyyyoo.r-z80fyv.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-19wmn03")
		if closeSvg != nil {
			log.Println("检测到弹出****")
			closeSvg.Click()
		}
		divs, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-1habvwh.r-18u37iz.r-1w6e6rj.r-1wtj0ep")
		if err != nil {
			log.Println(err)
		}
		if len(divs) <= 0 {
			log.Println("未找到follwButton父节点")
			return false, nil
		}
		divs[0].Click()
		time.Sleep(3 * time.Second)
		divsTwo := findFllowBtn(divs)
		if len(divsTwo) <= 0 {
			log.Println("未找到follwButton")
			return false, nil
		}

		//关注检测
		text, err := divsTwo[0].Text()
		log.Println("获取Ttext", text)
		if err != nil {
			log.Println("获取关注内容错误", err)
		} else {
			if !strings.EqualFold(text, "Following") && !strings.EqualFold(text, "正在关注") {
				divsTwo[0].Click()
			}
		}
		time.Sleep(5 * time.Second)
		return true, err

	}, 30*time.Second)

}

func findFllowBtn(divs []selenium.WebElement) []selenium.WebElement {
	divsTwo, err := divs[0].FindElements(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-2yi16.r-1qi8awa.r-1ny4l3l.r-ymttw5.r-o7ynqc.r-6416eg.r-lrvibr")

	if err != nil {
		log.Println(err)
	}
	log.Println(divs[0].GetAttribute("class"))
	if len(divsTwo) <= 0 {
		log.Println("第一次未找到follwButton")
	} else {
		return divsTwo
	}
	divsTwo, err = divs[0].FindElements(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-sdzlij.r-1phboty.r-rs99b7.r-2yi16.r-1qi8awa.r-1ny4l3l.r-ymttw5.r-o7ynqc.r-6416eg.r-lrvibr")
	if len(divsTwo) <= 0 {
		log.Println("第二次未找到follwButton")
	} else {
		return divsTwo
	}
	return nil
}

func TweetTwitter(wd selenium.WebDriver, div selenium.WebElement) {
	//util.GetWind(0, wd)
	hanle, _ := wd.CurrentWindowHandle()
	tweetes, err := div.FindElements(selenium.ByTagName, "a")
	if err != nil {
		log.Println(err)
	}
	for _, tweetDiv := range tweetes {
		tweetDiv.Click()
		time.Sleep(1 * time.Second)
		util.GetLastWindow(wd)
		Retweet(wd)
		//wd.Close()
		wd.SwitchWindow(hanle)

	}
	//util.GetWind(0, wd)
	wd.SwitchWindow(hanle)
}

func Retweet(wd selenium.WebDriver) {
	wd.Refresh()
	//检测弹窗
	findDiag(wd)
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		we, err := wd.FindElements(selenium.ByCSSSelector, ".r-4qtqp9.r-yyyyoo.r-50lct3.r-dnmrzs.r-bnwqim.r-1plcrui.r-lrvibr.r-1srniue")
		if err != nil {
			log.Println(err)
		}

		if len(we) <= 0 {
			log.Println("未获取到Retweet元素")
			return false, err
		}

		retweet, err := we[1].FindElement(selenium.ByXPATH, "./..")
		if err != nil {
			log.Println("查找父节点失败", err)
		}
		//点击
		retweet.Click()

		reweetBtns, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-1loqt21.r-18u37iz.r-1ny4l3l.r-ymttw5.r-1f1sjgu.r-o7ynqc.r-6416eg.r-13qz1uu")
		if err != nil || len(reweetBtns) <= 0 {
			log.Println(err)
			return false, nil
		}
		//检测是否已经转发
		time.Sleep(2 * time.Second)
		text, err := reweetBtns[0].Text()
		if err != nil {
			log.Println("检测是否转发出错", err)
		} else {
			if !strings.EqualFold(text, "Undo Retweet") && !strings.EqualFold(text, "撤销转推") {
				reweetBtns[0].Click()
			}
		}

		time.Sleep(1 * time.Second)
		parentLike, _ := we[2].FindElement(selenium.ByXPATH, "./..")
		//检测是否喜欢
		parentLikeParent, _ := parentLike.FindElement(selenium.ByXPATH, "./..")
		parentLikeParent, _ = parentLikeParent.FindElement(selenium.ByXPATH, "./..")
		labelText, err := parentLikeParent.GetAttribute("aria-label")

		if err != nil {
			log.Println("检测是否喜欢出错", err)
		} else {
			if !strings.EqualFold(labelText, "已喜欢") && !strings.EqualFold(labelText, "Liked") {
				parentLike.Click()
			}
		}

		time.Sleep(1 * time.Second)
		return true, err
	})

}

func findDiag(wd selenium.WebDriver) {
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		we, err := wd.FindElement(selenium.ByCSSSelector, ".css-1dbjc4n.r-z6ln5t.r-14lw9ot.r-1867qdf.r-1jgb5lz.r-pm9dpa.r-1rnoaur.r-494qqr.r-13qz1uu")
		if we == nil || err != nil {
			log.Println(err)
			return false, nil
		}
		sp, err := we.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
		for _, span := range sp {
			txt, _ := span.Text()
			log.Println(txt)
			if strings.EqualFold(txt, "Cancel") || strings.EqualFold(txt, "取消") {
				span.Click()
			}
		}
		return true, nil
	}, 20*time.Second)

}
