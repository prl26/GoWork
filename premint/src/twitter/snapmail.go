/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-02-27 17:49:44
 * @LastEditTime: 2023-03-24 17:39:55
 */
package twitter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const BASEURL = "https://snapmail.cc/emaillist/"

type Body struct {
}
type Result struct {
	Subject string `json:"subject"`
}

func GetTwitterCode(em string) string {

	req, err := http.NewRequest(http.MethodGet, BASEURL+em, nil)
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{Transport: &http.Transport{}}
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	result := []Result{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
	}
	var reply string
	if len(result) > 0 {
		log.Println(result[0].Subject)
		re := regexp.MustCompile("[0-9]+")
		reply = re.FindString(result[0].Subject)
		log.Println(reply)
	}
	return reply
}
