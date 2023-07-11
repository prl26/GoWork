/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-06 15:30:59
 * @LastEditTime: 2023-04-26 15:03:32
 */
package bitbrowser

import (
	"bytes"
	"encoding/json"
	"github.com/JianLinWei1/premint-selenium/model"
	"io/ioutil"
	"log"
	"net/http"
)

const BASE_URL = "http://127.0.0.1:54345"

// openBrower参数结构体
type RequestData struct {
	ID             string   `json:"id"`
	LoadExtensions bool     `json:"loadExtensions"`
	Args           []string `json:"args"`
	ExtractIp      bool     `json:"extractIp"`
}

type ResponseData struct {
	Success bool    `json:"success"`
	Data    DataRes `json:"data"`
}

type DataRes struct {
	Ws          string `json:"ws"`
	Http        string `json:"http"`
	CoreVersion string `json:"coreVersion"`
	Driver      string `json:"driver"`
}

// updateBrower 参数结构体

func OpenBrowser(id string) DataRes {
	reqData := RequestData{
		ID:             id,
		LoadExtensions: true,
		Args:           []string{
			// "-headless",
		},
		ExtractIp: true,
	}
	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}
	log.Println("请求BIt:", string(reqDataJson))

	postReq, _ := http.NewRequest(http.MethodPost, BASE_URL+"/browser/open", bytes.NewBuffer(reqDataJson))
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
	log.Println("bit返回", string(body))
	responseBody := &ResponseData{}
	json.Unmarshal(body, responseBody)
	return responseBody.Data
}
func UpdataBrower(reqData model.BitDetailStruct) bool {

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}
	log.Println("更新bit:", string(reqDataJson))

	postReq, _ := http.NewRequest(http.MethodPost, BASE_URL+"/browser/update/partial", bytes.NewBuffer(reqDataJson))
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
	log.Println("bit返回", string(body))
	responseBody := &model.UpdateResponseData{}
	json.Unmarshal(body, responseBody)
	return responseBody.Success
}
