package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Log struct {
	JsonPayload JsonPayload `json:"jsonPayload"`
}

type JsonPayload struct {
	Msg string `json:"msg"`
}

type Msg struct {
	Request Request `json:"request"`
}

type Request struct {
	Headers Headers `json:"headers"`
}

type Headers struct {
	Authorization []interface{} `json:"Authorization"`
}

func main() {
	fileContent, err := os.Open("logs.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var logs []Log
	json.Unmarshal(byteResult, &logs)

	AuthorizationMap := make(map[string]int)

	for i := 0; i < len(logs); i++ {
		var msg Msg
		json.Unmarshal([]byte(logs[i].JsonPayload.Msg), &msg)
		AuthorizationValue := fmt.Sprint(msg.Request.Headers.Authorization)
		AuthorizationMap[AuthorizationValue] = AuthorizationMap[AuthorizationValue] + 1
	}

	fmt.Print("Msg: ", len(AuthorizationMap))
}
