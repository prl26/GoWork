/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-13 11:33:59
 * @LastEditTime: 2023-05-04 10:27:10
 */
package metamask

import (
	"log"

	"github.com/szyhf/go-excel"
)

type ExcelInfo struct {
	NetWorkName string `xlsx:"column(网络名称)"`
	NewRpcUrl   string `xlsx:"column(新的RPC URL)"`
	ChainID     string `xlsx:"column(链ID)"`
	Currency    string `xlsx:"column(货币符号)"`
	BlockUrl    string `xlsx:"column(区块浏览器 URL)"`
	BitId       string `xlsx:"column(窗口ID)"`
	MetaMaskPwd string `xlsx:"column(metamask密码)"`
	MetaKey     string `xlsx:"column(metamask私钥)"`
}

type ExcelInfoTransfer struct {
	OKXAddr      string `xlsx:"column(OKX地址)"`
	MetaMaskPwd  string `xlsx:"column(MetaMask密码)"`
	MetaKey      string `xlsx:"column(metamask私钥)"`
	BitId        string `xlsx:"column(窗口ID)"`
	Count        string `xlsx:"column(金额)"`
	Tokens       string `xlsx:"column(币种)"`
	ContarctAddr string `xlsx:"column(代币合约地址)"`
}

func getExcelInfos(fileUrl string) []ExcelInfo {
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

func getExcelInfoTransfer(fileUrl string) []ExcelInfoTransfer {
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

	var infoList []ExcelInfoTransfer
	rd.ReadAll(&infoList)

	return infoList
}
