// Sample logging-quickstart writes a log entry to Cloud Logging.
package main

import (
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

func main() {
	getEntries("likeminds-prod-prj-5464", fmt.Sprintf(
		`resource.type="k8s_container" 
AND resource.labels.project_id="likeminds-prod-prj-5464" AND resource.labels.location="asia-south1" 
AND resource.labels.cluster_name="likeminds-prod-autopilot-cluster"
AND resource.labels.namespace_name="app-deploy"
AND labels.k8s-pod/app="kettle" severity>=DEFAULT
AND "Client.Timeout"
AND timestamp > "%s"`, time.Now().Add(-744*time.Hour).Format(time.RFC3339)))
}

func getEntries(projectID string, query string) ([]*logging.Entry, error) {
	ctx := context.Background()
	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	var entries []*logging.Entry

	iter := adminClient.Entries(ctx,
		logadmin.Filter(query),
		logadmin.NewestFirst(),
	)

	for {
		entry, err := iter.Next()
		if err == iterator.Done {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	fmt.Printf("Total entries: %d", len(entries))
	return entries, nil
}
