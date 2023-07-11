/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-17 15:19:18
 * @LastEditTime: 2023-03-17 16:13:49
 */
package bitbrowser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JianLinWei1/premint-selenium/model"
	"io/ioutil"
	"net/http"
)

func ListBrowser(reqData model.ListRequestData) []string {

	reqDataJson, err := json.Marshal(reqData)
	if err != nil {
		panic(err)
	}
	fmt.Println("请求BIt:", string(reqDataJson))
	postReq, _ := http.NewRequest(http.MethodPost, BASE_URL+"/browser/list", bytes.NewBuffer(reqDataJson))
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
	//fmt.Println("bit返回", string(body))
	responseBody := &model.ListResponseData{}
	json.Unmarshal(body, responseBody)
	//fmt.Println(responseBody)
	var list []string
	for i := 0; i < len(responseBody.Data.List); i++ {
		list = append(list, responseBody.Data.List[i].ID)
	}
	return list
}
