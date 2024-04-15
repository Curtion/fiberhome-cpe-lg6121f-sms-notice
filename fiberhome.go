package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	iv = []byte(intAesIV())
)

type session struct {
	SessionId string `json:"sessionid"`
}

func getSessionId() (string, error) {
	content, err := requestGet("/api/tmp/FHNCAPIS?ajaxmethod=get_refresh_sessionid")
	if err != nil {
		return "", err
	}
	var data = new(session)
	err = json.Unmarshal([]byte(content), data)
	if err != nil {
		return "", err
	}
	return data.SessionId, nil
}

func requestGet(path string) (string, error) {
	resp, err := http.Get(url + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("请求:%v错误: HTTP Code: %v", path, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type body struct {
	DataObj    interface{} `json:"dataObj"`
	AjaxMethod string      `json:"ajaxmethod"`
	SessionId  string      `json:"sessionid"`
}

func requestPost(dataObj interface{}, path string, ajaxmethod string) (string, error) {
	sessionId, err := getSessionId()
	if err != nil {
		log.Print(err)
	}

	data := body{
		DataObj:    dataObj,
		AjaxMethod: ajaxmethod,
		SessionId:  sessionId,
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	log.Print("请求路径: ", url+path, " 请求数据: ", string(dataJson), " 请求方法: ", ajaxmethod)

	postData, err := encryptFunc(string(dataJson), []byte(sessionId[0:16]), iv)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url+path, "application/json", strings.NewReader(fmt.Sprintf("%x", postData)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("请求:%v错误: HTTP Code: %v", ajaxmethod, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var stringData string

	if ajaxmethod == "DO_WEB_LOGIN" {
		stringData = string(body)
	} else {
		byteData, err := hex.DecodeString(string(body))
		if err != nil {
			return "", err
		}
		stringData, err = decryptFunc(byteData, []byte(sessionId[0:16]), iv)
		if err != nil {
			return "", err
		}
	}

	return stringData, nil
}
