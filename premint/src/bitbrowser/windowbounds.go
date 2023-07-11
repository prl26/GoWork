/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-06 15:30:59
 * @LastEditTime: 2023-04-25 13:47:08
 */
package bitbrowser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RequestDataWindowbounds struct {
	ID             string   `json:"id"`
	LoadExtensions bool     `json:"loadExtensions"`
	Args           []string `json:"args"`
	ExtractIp      bool     `json:"extractIp"`
}

type ResponseDataWindowbounds struct {
	Success bool    `json:"success"`
	Data    DataRes `json:"data"`
}

type DataResWindowbounds struct {
	Ws          string `json:"ws"`
	Http        string `json:"http"`
	CoreVersion string `json:"coreVersion"`
	Driver      string `json:"driver"`
}

var para = map[string]interface{}{
	"type":    "box",
	"startX":  0,
	"startY":  0,
	"width":   500,
	"height":  300,
	"col":     4,
	"spaceX":  0,
	"spaceY":  0,
	"offsetX": 50,
	"offsetY": 50,
}

func WindowboundsByPara() DataRes {
	// 打开 JSON 文件
	//configFile, err := os.Open("./窗口排列.json")
	//if err != nil {
	//	log.Println("读取窗口排列文件出错", err)
	//}
	//defer configFile.Close()
	//
	//// 读取 JSON 文件内容
	//reqDataJson, err := ioutil.ReadAll(configFile)
	//if err != nil {
	//	log.Println("读取窗口排列文件出错", err)
	//}
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("请求BIt:", string(para))
	Para, _ := json.Marshal(para)
	postReq, _ := http.NewRequest(http.MethodPost, BASE_URL+"/windowbounds", bytes.NewBuffer(Para))
	postReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: &http.Transport{}}

	response, err := client.Do(postReq)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("bit返回", string(body))
	responseBody := &ResponseData{}
	json.Unmarshal(body, responseBody)
	return responseBody.Data
}

func Windowbounds() DataRes {
	// 打开 JSON 文件
	configFile, err := os.Open("./窗口排列.json")
	if err != nil {
		log.Println("读取窗口排列文件出错", err)
	}
	defer configFile.Close()

	// 读取 JSON 文件内容
	reqDataJson, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println("读取窗口排列文件出错", err)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("请求BIt:", string(reqDataJson))
	postReq, _ := http.NewRequest(http.MethodPost, BASE_URL+"/windowbounds", bytes.NewBuffer(reqDataJson))
	postReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: &http.Transport{}}

	response, err := client.Do(postReq)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("bit返回", string(body))
	responseBody := &ResponseData{}
	json.Unmarshal(body, responseBody)
	return responseBody.Data
}
