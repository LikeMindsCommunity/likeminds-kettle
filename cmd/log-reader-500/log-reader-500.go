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

	Distinct500Map := make(map[string]int)

	for i := 0; i < len(logs); i++ {
		TextPayload := logs[i].TextPayload

		DateTime := strings.Split(TextPayload, "|")[0]
		Date := strings.Split(DateTime, "-")[0]
		IP := strings.Split(logs[i].TextPayload, "|")[3]
		API := strings.Split(logs[i].TextPayload, "|")[4]

		Distinct500MapKey := Date + IP + API
		Distinct500Map[Distinct500MapKey] = Distinct500Map[Distinct500MapKey] + 1
	}

	fmt.Print("Msg: ", len(Distinct500Map))
}
