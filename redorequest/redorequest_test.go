package redorequest

import "testing"

func TestRedoRequest(t *testing.T) {
	RedoRequest(NginxConf{
		LogFilePath: "/home/ec/mrver/redorequest/request.txt",
		Filter: map[string][]string{
			"status":      {"200", "499"},
			"request_uri": {"/api/helloworld"},
		},
		RouterField: "request_uri",
		Server:      "http://localhost:8000",
	})
}
