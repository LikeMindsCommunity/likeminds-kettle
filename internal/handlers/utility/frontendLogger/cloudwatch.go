package frontendLogger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	"github.com/nateshr/likeminds-authentication/internal/logging"
)

type CloudwatchPayloadEntry struct {
	JsonPayload map[string]interface{} `json:"jsonPayload"`
	Timestamp   time.Time              `json:"timestamp"`
}

func GetCloudwatchClient() (*cloudwatchlogs.Client, error) {
	// Load AWS SDK config
	cloudwatchIAMUserKey := environment.GoDotEnvVariable(environment.CloudwatchIAMUserKey)
	cloudwatchIAMUserSecret := environment.GoDotEnvVariable(environment.CloudwatchIAMUserSecret)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cloudwatchIAMUserKey, cloudwatchIAMUserSecret, "")),
	)
	if err != nil {
		return nil, err
	}

	// Create CloudWatchLogs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	return client, nil
}

func LogToCloudWatch(client *cloudwatchlogs.Client, logGroupName string, logStreamName string, entries []CloudwatchPayloadEntry) error {

	// Put log events
	for _, entry := range entries {
		timestamp := entry.Timestamp.UnixMilli()
		jsonString, err := json.Marshal(entry.JsonPayload)
		if err != nil {
			logging.Error(fmt.Sprint("Error marshalling JSON: ", err.Error()))
		}

		output, err := client.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  &logGroupName,
			LogStreamName: &logStreamName,
			LogEvents: []types.InputLogEvent{{
				Message:   aws.String(string(jsonString)),
				Timestamp: aws.Int64(timestamp),
			}},
		})
		if err != nil {
			logging.Info(fmt.Sprint("PutLogEvents function output: ", output))
			return err
		}
	}
	return nil
}

func CreateLogGroupIfNotExist(client *cloudwatchlogs.Client, logGroupName string) error {
	_, err := client.CreateLogGroup(context.TODO(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		return err
	}

	retentionDays := RetentionDays
	_, err = client.PutRetentionPolicy(context.TODO(), &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    &logGroupName,
		RetentionInDays: &retentionDays,
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateLogStreamIfNotExist(client *cloudwatchlogs.Client, logGroupName, logStreamName string) error {
	_, err := client.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &logStreamName,
	})
	return err
}
