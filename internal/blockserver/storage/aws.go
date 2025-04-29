package storage

import (
	"context"
	"fmt"
	"skybox-backend/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Statically save the AWS S3 client for reuse
var s3Client *s3.Client

// GetS3Client returns the AWS S3 client
func GetS3Client() *s3.Client {
	if s3Client == nil {
		s3Client = NewAWSClient()
	}
	return s3Client
}

func NewAWSClient() *s3.Client {
	fmt.Println("Connecting to AWS S3...")

	if configs.Config.AWSEnabled == false {
		fmt.Println("AWS S3 is not enabled in the configuration.")
		return nil
	}

	// Load AWS configuration from environment variables or config file
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(configs.Config.AWSRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     configs.Config.AWSKey,
				SecretAccessKey: configs.Config.AWSSecret,
				SessionToken:    configs.Config.AWSSessionToken,
			},
		}),
	)
	if err != nil {
		panic(fmt.Errorf("unable to load AWS configuration, %v", err))
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Test the connection
	_, err = s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		panic(fmt.Errorf("unable to connect to AWS S3, %v", err))
	}

	fmt.Println("Connected to AWS S3")
	return s3Client
}

// CloseAWSClient is a placeholder for cleaning up AWS resources
func CloseAWSClient() {
	fmt.Println("Closing AWS resources (if any)...")
	// Add any necessary cleanup code here
}
