/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-13 11:33:59
 * @LastEditTime: 2023-04-22 12:30:31
 */
package util

import (
	"github.com/JianLinWei1/premint-selenium/model"
	"github.com/szyhf/go-excel"
	"log"
)

type ExcelInfo struct {
	PremintUrl  string `xlsx:"column(premint地址)"`
	MetaPwd     string `xlsx:"column(metamask登录密码)"`
	MetaWords   string `xlsx:"column(metamask助记词)"`
	MetaKey     string `xlsx:"column(metamask私钥)"`
	TwitterUser string `xlsx:"column(twitter账号)"`
	TwitterPwd  string `xlsx:"column(twitter密码)"`
	TwitterAt   string `xlsx:"column(twitter用户名)"`
	ProxyUrl    string `xlsx:"column(代理地址)"`
	BitId       string `xlsx:"column(窗口ID)"`
}

func GetExcelInfos(fileUrl string) []ExcelInfo {
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

	var infoList []ExcelInfo
	rd.ReadAll(&infoList)

	return infoList
}
func GetOMNIExcelInfos(fileUrl string) []model.OMNIExcelInfo {
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

	var infoList []model.OMNIExcelInfo
	rd.ReadAll(&infoList)

	return infoList
}
