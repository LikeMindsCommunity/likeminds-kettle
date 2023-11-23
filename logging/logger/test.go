// package main

// import (
// 	"context"

// 	"cloud.google.com/go/logging"
// )

// func main() {

// 	println("Starting")

// 	ctx := context.Background()
// 	logClient, err := logging.NewClient(ctx, "likeminds-nonprod-prj-24e1")
// 	if err != nil {
// 		println("Failed to create log client with error: ", err.Error())
// 		return
// 	}

// 	defer logClient.Close()

// 	logger := logClient.Logger("test")
// 	logger.Log(logging.Entry{Payload: "Hello World"})

// 	println("Done")
// }
