/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-24 14:52:34
 * @LastEditTime: 2023-03-27 14:38:41
 */
package twitter

import (
	"log"

	"github.com/szyhf/go-excel"
)

type ExcelInfo struct {
	Email string `xlsx:"column(邮箱)"`
	BitId string `xlsx:"column(窗口ID)"`
	Pwd   string `xlsx:"column(密码)"`
}

type ExcelInfoTweet struct {
	Link  string `xlsx:"column(tweet链接)"`
	Tp    int    `xlsx:"column(操作类型 1:关注2:喜欢加转发)"`
	BitId string `xlsx:"column(窗口ID)"`
}

func getExcelInfos[T ExcelInfo | ExcelInfoTweet](fileUrl string) []T {
	conn := excel.NewConnecter()
	err := conn.Open(fileUrl)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	rd, err := conn.NewReader("Sheet1")
	if err != nil {
		log.Println(err)
	}
	defer rd.Close()

	var infoList []T
	rd.ReadAll(&infoList)

	return infoList
}
