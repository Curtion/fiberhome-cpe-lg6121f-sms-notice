package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	username string
	password string
	waitTime int
	brakKey  string
	url      string
)

func init() {
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.IntVar(&waitTime, "w", 10, "wait time")
	flag.StringVar(&brakKey, "b", "", "bark key")
	flag.StringVar(&url, "url", "http://192.168.8.1", "5g cpe url")
}

func main() {
	flag.Parse()
	logout := make(chan bool)
	go login(logout)

	for {
		<-logout
		log.Printf("账号在其它地方登录, 等待%d分钟后重新登录", waitTime)
		<-time.After(time.Duration(waitTime) * time.Minute)
		go login(logout)
	}

}

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(logout chan bool) {
	user := userLogin{
		Username: username,
		Password: password,
	}

	loginInfo, err := requestPost(user, "/api/sign/DO_WEB_LOGIN", "DO_WEB_LOGIN")
	if err != nil {
		log.Print(err)
	}
	status := strings.Split(loginInfo, "|")[0]
	if status == "1" {
		panic("当前已有用户在别处登录，请稍后登录")
	} else if status == "2" {
		panic("您的连续错误登录次数已经达到3次, 请1分钟后再试")
	} else if status == "3" {
		panic("管理帐号已被禁用，请另选帐号登录")
	} else if status == "4" {
		panic("用户名或密码错误，请重试")
	} else if status == "5" {
		panic("未知错误")
	} else if status == "0" {
		log.Print("登录成功")
	}
	cancel := make(chan bool)
	go watchSms(cancel)
	for {
		islogin, err := requestGet("/api/tmp/IS_LOGGED_IN")
		if err != nil {
			log.Print(err)
		}
		log.Print("是否登录: ", strings.TrimSpace(islogin))

		hb, err := requestGet("/api/tmp/heartbeat")
		if err != nil {
			log.Print(err)
		}
		log.Print("心跳请求: ", strings.TrimSpace(hb))

		if strings.TrimSpace(islogin) == "0" || strings.TrimSpace(hb) != "true" {
			close(cancel)
			logout <- true
			return
		}

		<-time.After(3 * time.Second)
	}
}

type NewSmsFlag struct {
	NewSmsFlag string `json:"new_sms_flag"`
}

func watchSms(cancel chan bool) {
	for {
		select {
		case <-cancel:
			log.Print("退出监听短信")
			return
		default:
		}
		smsFlag, err := requestGet("/api/tmp/FHAPIS?ajaxmethod=get_new_sms")
		if err != nil {
			log.Print(err)
		}
		var data = new(NewSmsFlag)
		err = json.Unmarshal([]byte(smsFlag), data)
		if err != nil {
			log.Print(err)
		}
		log.Print("是否有新短信: ", strings.TrimSpace(data.NewSmsFlag))
		if strings.TrimSpace(data.NewSmsFlag) == "true" {
			smsNotice()
		}
		<-time.After(3 * time.Second)
	}
}

func smsNotice() {
	sms, err := requestPost(nil, "/api/tmp/FHAPIS", "get_sms_data")
	if err != nil {
		log.Print(err)
	}

	var msg map[string]interface{}
	err = json.Unmarshal([]byte(sms), &msg)
	if err != nil {
		log.Print(err)
	}

	for _, v := range msg {
		for _, vv := range v.(map[string]interface{}) {
			if m, ok := vv.(map[string]interface{}); ok {
				if m["rcvorsend"] == "recv" {
					if m["isOpened"] == "0" {
						log.Print("--------------------新短信--------------------")
						log.Print("短信号码: ", v.(map[string]interface{})["session_phone"])
						log.Print("短信内容: ", m["msg_content"])
						log.Print("短信时间: ", m["time"])
						log.Print("短信ID: ", m["childnode"])
						log.Print("--------------------------------------------")
						barkNotice(v.(map[string]interface{})["session_phone"].(string), m["msg_content"].(string))
						readSms(m["childnode"].(string))
					}
				}
			}
		}
	}
}

func readSms(id string) {
	var data = map[string]interface{}{
		"url": map[string]interface{}{
			"smsIsopend" + id: "InternetGatewayDevice.X_FH_MobileNetwork.SMS_Recv.SMS_RecvMsg." + id + ".isOpened",
		},
		"value": map[string]interface{}{
			"smsIsopend" + id: "1",
		},
	}
	res, err := requestPost(data, "/api/tmp/FHAPIS", "set_value_by_xmlnode")
	if err != nil {
		log.Print(err)
	}
	log.Printf("短信[%s]已读: %s", id, res)
}

func barkNotice(title, content string) (string, error) {
	if brakKey == "" {
		return "", fmt.Errorf("bark key is empty")
	}
	url := fmt.Sprintf("https://api.day.app/%s/%s/%s?level=timeSensitive", brakKey, title, content)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
