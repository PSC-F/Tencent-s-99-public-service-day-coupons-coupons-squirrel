package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type PublicData struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}
type Detail struct {
	StatisticsDate string `json:"statistics_date"`
	ClickTimes     int    `json:"click_times"`
	EstimateAmount string `json:"estimate_amount"`
	PayNum         int    `json:"pay_num"`
	SettleAmount   string `json:"settle_amount"`
}
type Statistics struct {
	ClickTimes     int    `json:"click_times"`
	EstimateAmount string `json:"estimate_amount"`
	PayNum         int    `json:"pay_num"`
	SettleAmount   string `json:"settle_amount"`
}
type List struct {
	ChannelID      string `json:"channel_id"`
	ClickTimes     string `json:"click_times"`
	StatisticsDate string `json:"statistics_date"`
	PayNum         string `json:"pay_num"`
	EstimateAmount string `json:"estimate_amount"`
	SettleAmount   string `json:"settle_amount"`
	ChannelName    string `json:"channel_name"`
	SortTime       int    `json:"sort_time"`
}
type Data struct {
	Detail     []Detail   `json:"detail"`
	Statistics Statistics `json:"statistics"`
	List       []List     `json:"list"`
}

var client http.Client

const (
	LOGIN_URL   = "https://alliance.yunzhanxinxi.com/login"
	ContentType = "application/x-www-form-urlencoded"
)

func main() {
	http.HandleFunc("/contributions", getPublicData)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("HTTP SERVER failed,err:", err)
		return
	}
}

func getPublicData(writer http.ResponseWriter, request *http.Request) {
	login()
	res := getDashBoard()
	tmpl, err := template.ParseFiles("./publicData.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板, 并将结果写入w
	tmpl.Execute(writer, res)
}

// 获取后台看板数据
func getDashBoard() PublicData {
	var PublicData PublicData
	c, l := getDate()
	url := "https://alliance.yunzhanxinxi.com/order/list/statistics-info"
	contentType := "application/x-www-form-urlencoded"
	bodyStr := "channelId=127223&promotionId=&startTime=" + l.String() + "&endTime=" + c.String() + "&orderType="
	t, _ := client.Post(url, contentType, strings.NewReader(bodyStr)) // 在这里登陆
	defer t.Body.Close()
	body, _ := ioutil.ReadAll(t.Body)
	errJson := json.Unmarshal([]byte(string(body)), &PublicData)
	if errJson != nil {
		fmt.Println(errJson) // 错误写进日志文件
	}
	return PublicData
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
	t, _ := client.Post(LOGIN_URL, ContentType, strings.NewReader(b))
	body, _ := ioutil.ReadAll(t.Body)
	fmt.Println(string(body))
}
