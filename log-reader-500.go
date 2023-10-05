package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type LogError struct {
	TextPayload string `json:"textPayload"`
}

func main() {
	fileContent, err := os.Open("logs-500.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var logs []LogError
	json.Unmarshal(byteResult, &logs)

	IPMap := make(map[string]int)

	for i := 0; i < len(logs); i++ {
		IPMapValue := strings.Split(logs[i].TextPayload, "|")[3]
		IPMap[IPMapValue] = IPMap[IPMapValue] + 1
	}

	fmt.Print("Msg: ", len(IPMap))
}
