package Galxe

import (
	"errors"
	"github.com/tebeka/selenium"
	"log"
	"strings"
	"time"
)

func TwitterBind(wd selenium.WebDriver) (err error) {
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
		return err
	} else {
		//找到具有setting属性的div
		for _, v := range SocialLinks {
			linktext, _ := v.FindElement(selenium.ByCSSSelector, ".social-account-link-text")
			text, _ := linktext.Text()
			if strings.Contains(text, "Connect Twitter Account") {

			} else {
				continue
			}
		}
	}
	return
}
func ClickConnectTwitter(wd selenium.WebDriver, Connectbutton selenium.WebElement) (err error) {
	err = Connectbutton.Click()
	if err != nil {
		log.Println("点击connect Twitter Account失败")
		return err
	} else {

	}
	return
}
