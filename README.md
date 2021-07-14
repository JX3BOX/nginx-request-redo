## RedoNginxRequest

根据nginx的访问日志，将请求重新发送一遍。由于`nginx`日志无法记录出请求内容，因此仅支持 `GET`请求


## 快速开始

`git clone https://github.com/JX3BOX/nginx-request-redo.git`


根据注释说明 修改 `config.json` 

然后 执行

`go run main.go` 即可.

附：配置说明

```json
{
    "logFilePath": "/www/wwwlogs/nginx/access.log", // nginx日志位置
	"logFormat": "$remote_addr - $remote_user [$time_local] \"$request_method $request_uri $http_version\" $status $bytes_sent \"-\" \"$http_user_agent\"",  // 可选，有默认值， 日志格式参考: github.com/
    "routerFiled": "request_uri", // 路由对应的日志字段名，当使用默认日志格式时，该值可选
    "filter": { // 过滤的参数。 {[变量名]:Array<匹配值>}
        "status": [
            "200"
        ],
        "request_method": [
            "GET"
        ]
    },
    "server": "http://localhost:8080", // 服务器地址, 必填。 http(s)://ip:port 或者 http(s)://domain.com
    "extraHeaders": { // 需要增加的header头
        "Token": "1"
    },
    "extextraQueryParams": { // 需要增加的查询参数
        "v": "111"
    }
}
```

nginx日志默认格式为:

```
$remote_addr - $remote_user [$time_local] "$request_method $request_uri $http_version" $status $bytes_sent "-" "$http_user_agent"
```