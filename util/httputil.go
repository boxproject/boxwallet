package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpUtil struct {
	HttpClient *http.Client
	BaseURL    string
}

func Get(urlStr string, query map[string]string) (result string, err error) {
	str := ""
	i := 0
	for k, v := range query {
		if i > 0 {
			str += "&"
		}
		str += fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		i++
	}
	response, err := http.Get(urlStr + "?" + str)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Request(url string, method string, data interface{}) (result string, err error) {
	client := &http.Client{}
	jsonData, _ := json.Marshal(data)
	reqest, _ := http.NewRequest(method, url, strings.NewReader(string(jsonData)))
	reqest.Header.Set("Content-Type", "application/json")

	response, err := client.Do(reqest)
	if err != nil {
		return "", err
	}
	if response.StatusCode == 200 {
		body, err_ := ioutil.ReadAll(response.Body)
		if err_ != nil {
			return "", err_
		}
		bodystr := string(body)
		return bodystr, nil
	}
	return "", err

}
