package redorequest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/satyrius/gonx"
)

type NginxConf struct {
	LogFilePath      string            `json:"logFilePath"`      // nginx日志文件路径
	TargetRouter     []string          `json:"targetRouter"`     // 需要重新发起的接口路由
	LogFormat        string            `json:"logFormat"`        // 日志格式
	StatusCode       []int             `json:"statusCode"`       // 请求的状态码要求
	Server           string            `json:"server"`           // 服务器 地址, http(s)://ip:port 或者 http(s)://domain.com
	ExtraHeaders     map[string]string `json:"extraHeaders"`     // 需要增加的header头
	ExtraQueryParams map[string]string `json:"extraQueryParams"` // 需要增加的查询参数
}

func inStatusCode(status string, codeList []int) bool {
	if codeList == nil || len(codeList) == 0 {
		return true
	}
	for _, code := range codeList {
		if strconv.Itoa(code) == status {
			return true
		}
	}
	return false
}
func matchRouter(uri string, routerList []string) bool {
	if routerList == nil || len(routerList) == 0 {
		return true
	}

	for _, router := range routerList {
		// TODO 待改进，支持路由通配符等
		if router == uri {
			return true
		}
	}
	return false
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
		if status, err := rec.Field("status"); err != nil || !inStatusCode(status, conf.StatusCode) {
			continue
		}
		if uri, err := rec.Field("request_uri"); err == nil && uri != "" {
			if matchRouter(uri, conf.TargetRouter) {
				redo(uri, conf)
			}
		}
	}
}
