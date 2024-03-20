package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/greenscale-ai/amazon-bedrock-metadata/pkg/model"
	"github.com/greenscale-ai/amazon-bedrock-metadata/pkg/processor"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Define environment variable names
const (
	modelInvocationLogsInputBucketEnv       = "MODEL_INVOCATION_LOGS_INPUT_BUCKET"
	modelInvocationLogsInputBucketPrefixEnv = "MODEL_INVOCATION_LOGS_INPUT_BUCKET_PREFIX"
	modelInvocationLogsInputBucketRegionEnv = "MODEL_INVOCATION_LOGS_INPUT_BUCKET_REGION"
	metadataLogsOutputBucketEnv             = "METADATA_LOGS_OUTPUT_BUCKET"
	metadataLogsOutputBucketRegionEnv       = "METADATA_LOGS_OUTPUT_BUCKET_REGION"
	awsAccountIDEnv                         = "AWS_ACCOUNT_ID"
	pickLastHourEnv                         = "PICK_LAST_HOUR"
	yearEnv                                 = "YEAR"
	monthEnv                                = "MONTH"
	dayEnv                                  = "DAY"
	hourEnv                                 = "HOUR"
)

var Version = "number missing"

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("Starting amazon-bedrock-metadata process, build", Version)

	if runningFromLambda() {
		lambda.Start(processLogs)
	} else {
		processLogs()
	}
}

func runningFromLambda() bool {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Running from Lambda")
		return true
	}
	return false
}

func processLogs() {
	modelInvocationLogsInputBucket := os.Getenv(modelInvocationLogsInputBucketEnv)
	if modelInvocationLogsInputBucket == "" {
		log.Println("Error: Model Invocation Logs Input S3 bucket name is required.")
		return
	}

	modelInvocationLogsInputBucketPrefix := os.Getenv(modelInvocationLogsInputBucketPrefixEnv)

	modelInvocationLogsInputBucketRegion := os.Getenv(modelInvocationLogsInputBucketRegionEnv)
	if modelInvocationLogsInputBucketRegion == "" {
		log.Println("Error: Model Invocation Logs Input S3 bucket region is required.")
		return
	}

	metadataLogsOutputBucket := os.Getenv(metadataLogsOutputBucketEnv)
	if metadataLogsOutputBucket == "" {
		log.Println("Error: Metadata Logs Output S3 bucket name is required.")
		return
	}

	metadataLogsOutputBucketRegion := os.Getenv(metadataLogsOutputBucketRegionEnv)
	if metadataLogsOutputBucketRegion == "" {
		log.Println("Error: Metadata Logs Output S3 bucket region is required.")
		return
	}

	awsAccountID := os.Getenv(awsAccountIDEnv)
	if awsAccountID == "" {
		log.Println("Error: AWS Account ID is required.")
		return
	}

	pickLastHour := os.Getenv(pickLastHourEnv) != "" && os.Getenv(pickLastHourEnv) != "false"
	year, _ := strconv.Atoi(os.Getenv(yearEnv))
	month, _ := strconv.Atoi(os.Getenv(monthEnv))
	day, _ := strconv.Atoi(os.Getenv(dayEnv))
	hour, _ := strconv.Atoi(os.Getenv(hourEnv))

	if !pickLastHour {
		if year == 0 || year < 1970 {
			log.Println("Error: Invalid year. Please provide a valid year.")
			return
		}
		if month < 1 || month > 12 {
			log.Println("Error: Invalid month. Please provide a value between 1 and 12.")
			return
		}
		if day < 1 || day > 31 {
			log.Println("Error: Invalid day. Please provide a value between 1 and 31.")
			return
		}
		if hour < 0 || hour > 23 {
			log.Println("Error: Invalid hour. Please provide a value between 0 and 23.")
			return
		}
	} else {
		// Pick current time
		now := time.Now()
		hourAgo := now.Add(-1 * time.Hour)
		utcNow := hourAgo.In(time.UTC)

		year = utcNow.Year()
		month = int(utcNow.Month())
		day = utcNow.Day()
		hour = utcNow.Hour()
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(modelInvocationLogsInputBucketRegion),
	})
	s3ClientRead := s3.New(sess)

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(metadataLogsOutputBucketRegion),
	})
	s3ClientWrite := s3.New(sess)

	iamSess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	iamClient := iam.New(iamSess)

	identityTagsBuilder := model.NewIdentityTagsBuilder(iamClient)
	pwd, _ := os.Getwd()
	modelsFilePath := fmt.Sprintf("%s/models.json", pwd)
	modelPriceFile, err := os.Open(modelsFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	defer func(modelPriceFile *os.File) {
		err := modelPriceFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(modelPriceFile)

	modelsPriceDetails, err := io.ReadAll(modelPriceFile)
	if err != nil {
		log.Println(err)
		return
	}

	modelCostEstimator, err := model.NewCostEstimator(modelsPriceDetails)
	if err != nil {
		log.Println(err)
		return
	}

	// Estimating carbon footprint based on configuration for AWS Inferentia2 instance types and average global carbon intensity
	modelCarbonFootprint := model.NewCarbonFootprintEstimator(400, 768000, 450)

	modelMetaDataGenerator := model.NewMetadataGenerator(modelCostEstimator, modelCarbonFootprint, identityTagsBuilder)

	modelLogsProcessor := processor.NewProcessor(s3ClientRead, s3ClientWrite, modelMetaDataGenerator, modelInvocationLogsInputBucket, metadataLogsOutputBucket)

	err = modelLogsProcessor.ProcessModelInvocationLogs(awsAccountID, modelInvocationLogsInputBucketRegion, modelInvocationLogsInputBucketPrefix, year, month, day, hour)
	if err != nil {
		log.Println(err)
		return
	}
}
