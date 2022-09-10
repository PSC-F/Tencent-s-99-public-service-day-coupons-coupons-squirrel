package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var client http.Client

const (
	LOGIN_URL   = "https://alliance.yunzhanxinxi.com/login"
	CONTENT_TYPE = "application/x-www-form-urlencoded"
)

func main() {
	login()
	getDashBoard()
}

// 获取后台看板数据
func getDashBoard() {
	c, l := getDate()
	url := "https://alliance.yunzhanxinxi.com/order/list/statistics-info"
	contentType := "application/x-www-form-urlencoded"
	bodyStr := "channelId=127223&promotionId=&startTime=" + l.String() + "&endTime=" + c.String() + "&orderType="
	t, _ := client.Post(url, contentType, strings.NewReader(bodyStr)) 
	defer t.Body.Close()
	body, _ := ioutil.ReadAll(t.Body)
	fmt.Println(string(body))
}

// 获取七天内数据时间
func getDate() (time.Time, time.Time) {
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -7)
	currentTime.Format("2006-01-02")
	oldTime.Format("2006-01-02")
	return currentTime, oldTime
}

// 登陆后台接口、缓存登陆信息
func login() {
	b := "" 
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
	t, _ := client.Post(LOGIN_URL, CONTENT_TYPE, strings.NewReader(b))
	body, _ := ioutil.ReadAll(t.Body)
	fmt.Println(string(body))
}
