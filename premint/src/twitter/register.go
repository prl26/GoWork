package twitter

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/JianLinWei1/premint-selenium/src/wdservice"
	"github.com/tebeka/selenium"
)

func StartRegister() {
	eis := getExcelInfos[ExcelInfo]("./Twitter批量注册模板.xlsx")
	chs := make(chan string, len(eis))
	for i, ei := range eis {
		if ei.BitId == "" {
			chs <- "跳过"
			continue
		}
		wd, _ := wdservice.InitWd(i, ei.BitId)
		go toRegister(chs, i, ei, wd)
	}

	for i := 0; i < len(eis); i++ {
		log.Println(<-chs)
	}
}

func toRegister(ch chan<- string, i int, ei ExcelInfo, wd selenium.WebDriver) {
	if err := wd.Get("https://twitter.com/"); err != nil {
		log.Println(err)
	}
	clickRegisterBtn(wd)
	clickCreateAcount(wd, ei)
	verify(wd, ei)
}

func clickRegisterBtn(wd selenium.WebDriver) {

	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(5 * time.Second)
		registerBtnParent, err := wd.FindElements(selenium.ByCSSSelector, ".css-1dbjc4n.r-l5o3uw.r-1upvrn0")
		if err != nil {
			log.Println(err)
			return false, err
		}
		if len(registerBtnParent) <= 0 {
			log.Println("未找到twitter注册父节点")
			wd.Refresh()
			time.Sleep(10 * time.Second)
			return false, err
		}
		regiterBtn, err := registerBtnParent[0].FindElements(selenium.ByTagName, "a")
		if err != nil {
			log.Println(err)
		}
		if len(regiterBtn) <= 0 {
			log.Println("未找到twitter注册按钮")
			wd.Refresh()
			time.Sleep(5 * time.Second)
			return false, err
		}
		regiterBtn[1].Click()
		return true, err
	})
}

func clickCreateAcount(wd selenium.WebDriver, ei ExcelInfo) string {
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		time.Sleep(2 * time.Second)
		btn, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-ywje51.r-usiww2.r-2yi16.r-1qi8awa.r-1ny4l3l.r-ymttw5.r-o7ynqc.r-6416eg.r-lrvibr.r-13qz1uu")
		log.Println(err)
		if err != nil {
			return false, nil
		}
		if err := btn.Click(); err != nil {
			log.Println(err)
		}
		return true, err
	})
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		// email
		updateEmail, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-901oao.r-1cvl2hr.r-37j5jr.r-a023e6.r-16dba41.r-rjixqe.r-bcqeeo.r-1ff274t.r-qvutc0")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		if err := updateEmail.Click(); err != nil {
			log.Println(err)
		}
		return true, err
	})

	//input name
	inputs, err := wd.FindElements(selenium.ByCSSSelector, ".r-30o5oe.r-1niwhzg.r-17gur6a.r-1yadl64.r-deolkf.r-homxoj.r-poiln3.r-7cikom.r-1ny4l3l.r-t60dpp.r-1dz5y72.r-fdjqy7.r-13qz1uu")
	log.Println(err)
	name := randomString(rand.Intn(3) + 4)
	if err := inputs[0].SendKeys(name); err != nil {
		log.Println(err)
	}
	//input email
	if err := inputs[1].SendKeys(ei.Email); err != nil {
		log.Println(err)
	}
	//selete birth
	selects, err := wd.FindElements(selenium.ByCSSSelector, ".r-30o5oe.r-1niwhzg.r-17gur6a.r-1yadl64.r-18jsvk2.r-1loqt21.r-37j5jr.r-1inkyih.r-rjixqe.r-crgep1.r-1wzrnnt.r-1ny4l3l.r-t60dpp.r-xd6kpl.r-1pn2ns4.r-ttdzmv")

	if err != nil {
		log.Println(err)
	}
	selects[0].SendKeys(strconv.Itoa(rand.Intn(12)))
	selects[1].SendKeys(strconv.Itoa(rand.Intn(30)))
	selects[2].SendKeys(strconv.Itoa(rand.Intn(20) + 1970))
	//next

	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		nextBtn, err := wd.FindElement(selenium.ByCSSSelector, ".css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-19yznuf.r-64el8z.r-icoktb.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		time.Sleep(2 * time.Second)
		if err := nextBtn.Click(); err != nil {
			log.Println(err)
		}
		return true, err
	})
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		nextBtn2, err := wd.FindElement(selenium.ByCSSSelector, ".css-18t94o4.css-1dbjc4n.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-19yznuf.r-64el8z.r-1ny4l3l.r-1dye5f7.r-o7ynqc.r-6416eg.r-lrvibr")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		time.Sleep(2 * time.Second)
		if err := nextBtn2.Click(); err != nil {
			log.Println(err)
		}

		return true, err
	})

	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		nextBtn2, err := wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.css-1hf3ou5.r-poiln3.r-1inkyih.r-rjixqe.r-bcqeeo.r-qvutc0")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		time.Sleep(2 * time.Second)
		if err := nextBtn2[0].Click(); err != nil {
			log.Println(err)
		}

		return true, err
	})

	return name
}

func verify(wd selenium.WebDriver, ei ExcelInfo) {
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		ver, err := wd.FindElement(selenium.ByCSSSelector, ".sc-1io4bok-0.gdVRUf.heading.text")
		if ver != nil || err == nil {
			log.Println("*****检测到人工验证请手动操作*****")
			var continued string
			fmt.Println("输入y继续：")
			fmt.Scanln(&continued)
			if continued == "y" {
				return true, err
			}

			return false, nil
		} else {
			log.Println("*****检测到人工验证请手动操作*****")
			var continued string
			fmt.Println("输入y继续：")
			fmt.Scanln(&continued)
			if continued == "y" {
				return true, err
			}
			return false, nil
		}

		return true, err
	})

	//code
	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		inputCode, err := wd.FindElement(selenium.ByCSSSelector, ".r-30o5oe.r-1niwhzg.r-17gur6a.r-1yadl64.r-deolkf.r-homxoj.r-poiln3.r-7cikom.r-1ny4l3l.r-t60dpp.r-1dz5y72.r-fdjqy7.r-13qz1uu")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		code := GetTwitterCode(ei.Email)
		fmt.Println("sdfsd", code)
		if err := inputCode.SendKeys(code); err != nil {
			log.Println(err)
		}
		return true, err
	})

	wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		nextBtn2, err := wd.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.css-1hf3ou5.r-poiln3.r-1inkyih.r-rjixqe.r-bcqeeo.r-qvutc0")
		if err != nil {
			log.Println(err)
			return false, nil
		}
		time.Sleep(2 * time.Second)
		if err := nextBtn2[0].Click(); err != nil {
			log.Println(err)
		}

		return true, err
	})
	//pwd

}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
