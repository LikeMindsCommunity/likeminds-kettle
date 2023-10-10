package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Log struct {
	JsonPayload JsonPayload `json:"jsonPayload"`
	TimeStamp   string      `json:"timestamp"`
}

type JsonPayload struct {
	Msg string `json:"msg"`
}

type Msg struct {
	Request Request `json:"request"`
}

type Request struct {
	AbsoluteURI string  `json:"absolute_uri"`
	Headers     Headers `json:"headers"`
}

type Headers struct {
	Authorization []interface{} `json:"Authorization"`
}

func main() {
	fileContent, err := os.Open("./cmd/log-reader-timeout/logs-timeout.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var logs []Log
	json.Unmarshal(byteResult, &logs)

	DistinctTimeoutMap := make(map[string]int)

	for i := 0; i < len(logs); i++ {
		log := logs[i]
		var msg Msg
		json.Unmarshal([]byte(log.JsonPayload.Msg), &msg)

		Authorization := fmt.Sprint(msg.Request.Headers.Authorization)
		Date := strings.Split(log.TimeStamp, "T")[0]
		fmt.Println(Date)
		API := msg.Request.AbsoluteURI
		fmt.Println(API)

		DistinctTimeoutMapKey := Authorization + Date + API
		DistinctTimeoutMap[DistinctTimeoutMapKey] = DistinctTimeoutMap[DistinctTimeoutMapKey] + 1
	}

	fmt.Print(len(DistinctTimeoutMap))
}
