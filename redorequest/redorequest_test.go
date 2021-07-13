package redorequest

import "testing"

func TestRedoRequest(t *testing.T) {
	RedoRequest(NginxConf{
		LogFilePath:  "/home/ec/mrver/redorequest/request.txt",
		TargetRouter: []string{"/api/v1"},
		LogFormat:    `$remote_addr - $remote_user [$time_local] "$request_method $request_uri $http_version" $status $bytes_sent "-" "$http_user_agent"`,
		StatusCode:   []int{200, 499},
		Server:       "http://localhost:8000",
	})
}
