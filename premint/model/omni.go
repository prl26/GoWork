package model

type OMNIExcelInfo struct {
	HelpWords  string `xlsx:"column(助记词）"`
	PrivateKey string `xlsx:"column(私钥)"`
	PublicKey  string `xlsx:"column(公钥)"`
	Address    string `xlsx:"column(地址)"`
	Type       string `xlsx:"column(类型)"`
	BitId      string `xlsx:"column(窗口ID)"`
	MetaPwd    string `xlsx:"column(MetaMask密码)"`
}
