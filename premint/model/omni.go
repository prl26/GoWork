package model

type OMNIExcelInfo struct {
	Address string `xlsx:"column(地址)"`
	Type    string `xlsx:"column(类型)"`
	BitId   string `xlsx:"column(窗口ID)"`
	MetaPwd string `xlsx:"column(MetaMask密码)"`
}
