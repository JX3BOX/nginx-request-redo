package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"redonginx/redorequest"
)

func main() {
	body, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var conf redorequest.NginxConf

	err = json.Unmarshal(body, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if msg := redorequest.CheckConf(&conf); msg != "" {
		log.Fatal(msg)
	}
	redorequest.RedoRequest(conf)
}
