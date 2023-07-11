package main

import (
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/JianLinWei1/premint-selenium/src/bitbrowser"
	"log"
)

func Update() {
	ids := []string{"9c7e4790ee804e509e6b5dd9deb20bd1"}
	reqData := model.BitDetailStruct{
		Ids:                ids,
		GroupId:            "",
		Platform:           "",
		Host:               "154.213.165.207",
		Port:               "",
		BrowserFingerPrint: model.BrowserFingerPrint{},
	}
	res := bitbrowser.UpdataBrower(reqData)
	fmt.Print()
	println(res)

}

//批量更新比特浏览器信息
//func main() {
//	file := "D:\\Go Work\\premint\\Premint脚本点击模版.xlsx"
//	res := util.GetExcelInfos(file)
//	fmt.Println(res[0].BitId)
//	for i := 0; i < len(res); i++ {
//		var ids []string
//		ids = append(ids, res[0].BitId)
//		reqData := model.BitDetailStruct{
//			Ids:                ids,
//			GroupId:            "",
//			Platform:           "",
//			Host:               "154.213.165.207",
//			Port:               "",
//			BrowserFingerPrint: model.BrowserFingerPrint{},
//		}
//		res := bitbrowser.UpdataBrower(reqData)
//		fmt.Print()
//		println(res)
//	}
//}

func main() {
	req := model.ListRequestData{
		Page:     0,
		PageSize: 20,
		GroupId:  "2c9bc04788f5e6a20188f5f4cc0f1ba0",
	}
	res := bitbrowser.ListBrowser(req)
	log.Println("成功了")
	log.Println(res)
}
