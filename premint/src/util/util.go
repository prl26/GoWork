/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-06 11:08:57
 * @LastEditTime: 2023-06-06 11:35:31
 */
package util

import (
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/tebeka/selenium"
)

// 获取第几个窗口
func GetWind(i int, wd selenium.WebDriver) {
	// 获取所有的窗口句柄
	handles, err := wd.WindowHandles()
	if err != nil {
		log.Println(err)
	}
	log.Println("窗口代码：", handles)
	newHandle := handles[i]
	if err := wd.SwitchWindow(newHandle); err != nil {
		log.Println(err)
	}
	log.Println("选择窗口：", i)
}

func GetWindByName(name string, wd selenium.WebDriver) {
	// 获取所有的窗口句柄
	handles, err := wd.WindowHandles()
	if err != nil {
		log.Println(err)
	}
	for _, handle := range handles {
		if err := wd.SwitchWindow(handle); err != nil {
			log.Println(err)
		}

		title, _ := wd.Title()
		log.Println("窗口：", title)
		if strings.EqualFold(name, title) {
			break
		}
	}

}

func GetCurrentWindow(wd selenium.WebDriver) {
	handle, err := wd.CurrentWindowHandle()
	if err != nil {
		log.Println("选择当前窗口出错", err)
	}
	wd.SwitchWindow(handle)
}
func GetCurrentWindowAndReturn(wd selenium.WebDriver) string {
	handle, err := wd.CurrentWindowHandle()
	if err != nil {
		log.Println("选择当前窗口出错", err)
	}
	wd.SwitchWindow(handle)
	return handle
}
func GetLastWindow(wd selenium.WebDriver) {
	handles, err := wd.WindowHandles()
	if err != nil {
		log.Println(err)
	}
	log.Println("窗口数量：", handles)
	newHandle := handles[len(handles)-1]
	if err := wd.SwitchWindow(newHandle); err != nil {
		log.Println(err)
	}
	log.Println("选择窗口：", newHandle)
}

func OpenNewWin(wd selenium.WebDriver, numWin int) {

	// tabScript := `window.open();`
	// if _, err := wd.ExecuteScript(tabScript, nil); err != nil {
	// 	log.Println(err)
	// }
	body, _ := wd.FindElement(selenium.ByTagName, "body")
	body.SendKeys(selenium.F12Key)
	// err := wd.SendModifier(selenium.ControlKey, true)
	// if err != nil {
	// 	log.Println(err)
	// }
	// wd.SendModifier("t", true)
	// numWin += 1
	// GetWind(numWin, wd)
}

func ChromePath() string {
	f, err := os.Open("./config.txt")
	if err != nil {
		log.Println("err:", err)

	}
	defer f.Close()

	content, err := readAllContent(f)
	if err != nil {
		log.Println("err:", err)
	}
	log.Println("content:", content)
	return content[0]
}

func readAllContent(r io.Reader) ([]string, error) {
	var b = make([]byte, 4096)
	_, err := r.Read(b)
	if err != nil {
		return nil, err
	}

	l := strings.Split(string(b), "\r\n")
	return l, nil
}

// func CheckRole() {
// 	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2023-04-15 08:08:08", time.Local)

// 	tl := t.Unix() - time.Now().Unix()

// 	if tl <= 0 {
// 		panic(nil)
// 	}
// }

func SetLog(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			// 在此处处理或记录 panic 的信息
			log.Println(r)
			log.Println(string(debug.Stack()))
		}
	}()

	//defer func() {
	//	err := recover()
	//	if err != nil {
	//		log.Println(err)
	//		log.Println(string(debug.Stack()))
	//	}
	//}()
	fn()
}
