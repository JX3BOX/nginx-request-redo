package redorequest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/satyrius/gonx"
)

type NginxConf struct {
	LogFilePath      string              `json:"logFilePath"`      // nginx日志文件路径
	RouterField      string              `json:"routerFiled"`      // 路由对应的日志字段名，当使用默认日志格式时，该值可选， 否则必填
	LogFormat        string              `json:"logFormat"`        // 可选，有默认值， 日志格式参考: github.com/satyrius/gonx
	Filter           map[string][]string `json:"filter"`           // 需要过滤的参数。 {[变量名]:Array<匹配值>}
	Server           string              `json:"server"`           // 服务器 地址, http(s)://ip:port 或者 http(s)://domain.com
	ExtraHeaders     map[string]string   `json:"extraHeaders"`     // 需要增加的header头
	ExtraQueryParams map[string]string   `json:"extraQueryParams"` // 需要增加的查询参数
}

func matchFilter(value string, rules []string) bool {
	if rules == nil || len(rules) == 0 {
		return true
	}
	for _, router := range rules {
		if router == value {
			return true
		}
	}
	return false
}

func CheckConf(conf *NginxConf) string {

	if conf.Server == "" || conf.LogFilePath == "" {
		return "关键server, logFilePath不能为空"
	}

	if conf.LogFormat == "" {
		conf.LogFormat = `$remote_addr - $remote_user [$time_local] "$request_method $request_uri $http_version" $status $bytes_sent "-" "$http_user_agent"`
		if conf.RouterField == "" {
			conf.RouterField = "request_uri"
		}
	}

	return ""
}

var client = http.Client{}

func redo(uri string, conf NginxConf) {
	server := conf.Server + uri
	urlObject, _ := url.Parse(server)

	v := urlObject.Query()
	if conf.ExtraQueryParams != nil {
		for k, value := range conf.ExtraQueryParams {
			v.Set(k, value)
		}
	}

	urlObject.RawQuery = v.Encode()

	req, err := http.NewRequest(urlObject.String(), "GET", nil)
	if err != nil {
		log.Println(err)
		return
	}
	for k, value := range conf.ExtraHeaders {
		req.Header.Add(k, value)
	}

	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	msgBody, err := ioutil.ReadAll(response.Body)
	log.Println(string(msgBody))
}

func RedoRequest(conf NginxConf) {
	file, err := os.Open(conf.LogFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := gonx.NewReader(file, conf.LogFormat)
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		for field, value := range conf.Filter {
			if v, err := rec.Field(field); err != nil || !matchFilter(v, value) {
				continue
			}
		}
		if uri, err := rec.Field(conf.RouterField); err == nil && uri != "" {
			redo(uri, conf)
		}
	}
}
