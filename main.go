package main

import (
	"flag"
	"log"
	"strings"
	"time"
)

var (
	username string
	password string
	waitTime int
)

func init() {
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.IntVar(&waitTime, "w", 10, "wait time")
}

func main() {
	// {"dataObj":{"url":{"smsIsopend11":"InternetGatewayDevice.X_FH_MobileNetwork.SMS_Recv.SMS_RecvMsg.11.isOpened"},"value":{"smsIsopend11":"1"}},"ajaxmethod":"set_value_by_xmlnode","sessionid":"7IdwD0laV80U5Mqp3mDl65QtcLRHO5z0"}
	// {"success":"true"}
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
		log.Fatal(err)
	}
	status := strings.Split(loginInfo, "|")[0]
	if status == "1" {
		panic("当前已有用户在别处登录，请稍后登录")
	} else if status == "2" {
		panic("您的连续错误登录次数已经达到3次，请1分钟后再试")
	} else if status == "3" {
		panic("管理帐号已被禁用，请另选帐号登录")
	} else if status == "4" {
		panic("用户名或密码错误，请重试")
	} else if status == "5" {
		panic("未知错误")
	} else if status == "0" {
		log.Print("登录成功")
	}
	go getSms()
	for {
		islogin, err := requestGet("/api/tmp/IS_LOGGED_IN")
		if err != nil {
			log.Fatal(err)
		}
		log.Print("是否登录: ", strings.TrimSpace(islogin))

		hb, err := requestGet("/api/tmp/heartbeat")
		if err != nil {
			log.Fatal(err)
		}
		log.Print("心跳请求: ", strings.TrimSpace(hb))

		if strings.TrimSpace(islogin) == "0" || strings.TrimSpace(hb) != "true" {
			logout <- true
			break
		}

		<-time.After(3 * time.Second)
	}
}

func getSms() {
	for {
		smsFlag, err := requestGet("/api/tmp/FHAPIS?ajaxmethod=get_new_sms")
		if err != nil {
			log.Fatal(err)
		}
		log.Print("短信提醒: ", strings.TrimSpace(smsFlag))
		<-time.After(3 * time.Second)
	}
	// sms, err := requestPost(nil, "/api/tmp/FHAPIS", "get_sms_data")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print("短信: ", sms)
}

// func testDecrypt() {
// 	byteSlice, err := hex.DecodeString("B1BFD5D292A1A4C0BFE0838BA88E8C5AD2F49595E143C090E9F9590726A93226")
// 	if err != nil {
// 		panic("Failed to decode hex string: " + err.Error())
// 	}
// 	test, err := decryptFunc(byteSlice, []byte("7IdwD0laV80U5Mqp3mDl65QtcLRHO5z0"[0:16]), iv)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("测试解密: %s\n", test)
// }
