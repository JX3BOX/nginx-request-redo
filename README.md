## RedoNginxRequest

根据nginx的访问日志，将请求重新发送一遍。


## 快速开始

`git clone https://github.com/JX3BOX/nginx-request-redo.git`

配置说明:

```
type NginxConf struct {
	LogFilePath      string            `json:"logFilePath"`      // nginx日志文件路径
	TargetRouter     []string          `json:"targetRouter"`     // 需要重新发起的接口路由
	LogFormat        string            `json:"logFormat"`        // 日志格式
	StatusCode       []int             `json:"statusCode"`       // 请求的状态码要求
	Server           string            `json:"server"`           // 服务器 地址, http(s)://ip:port 或者 http(s)://domain.com
	ExtraHeaders     map[string]string `json:"extraHeaders"`     // 需要增加的header头
	ExtraQueryParams map[string]string `json:"extraQueryParams"` // 需要增加的查询参数
}

```

根据注释说明 修改 `config.json` 

然后 执行

`go run main.go` 即可