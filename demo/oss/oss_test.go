package oss

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"testing"
)

func TestOSS(t *testing.T) {
	//cosId, ok := os.LookupEnv("COS_APP_ID")
	//if !ok {
	//	panic("none COS_APP_ID")
	//}
	//cosKey, ok := os.LookupEnv("COS_APP_SECRET")
	//if !ok {
	//	panic("none COS_APP_SECRET")
	//}

	// Using the SDK's default configuration, load additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Build the request with its input parameters
	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}
}
