package processor

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/greenscale-ai/amazon-bedrock-metadata/pkg/model"
	"log"
	"sync"
)

type Processor struct {
	modelInvocation                *model.MetadataGenerator
	s3ClientRead                   *s3.S3
	s3ClientWrite                  *s3.S3
	modelInvocationLogsInputBucket string
	metadataLogsOutputBucket       string
}

func NewProcessor(s3ClientRead, s3ClientWrite *s3.S3, modelInvocation *model.MetadataGenerator, modelInvocationLogsInputBucket, metadataLogsOutputBucket string) *Processor {
	return &Processor{
		s3ClientRead:                   s3ClientRead,
		s3ClientWrite:                  s3ClientWrite,
		modelInvocation:                modelInvocation,
		modelInvocationLogsInputBucket: modelInvocationLogsInputBucket,
		metadataLogsOutputBucket:       metadataLogsOutputBucket,
	}
}

func (p *Processor) ProcessModelInvocationLogs(accountID, region, modelInvocationLogsInputBucketPrefix string, year, month, day, hour int) error {
	s3Objects, err := p.listObjectsInDateRange(accountID, region, modelInvocationLogsInputBucketPrefix, year, month, day, hour)
	if err != nil {
		return err
	}

	workerPool := make(chan struct{}, 10)

	var wg sync.WaitGroup
	for _, obj := range s3Objects {
		wg.Add(1)
		workerPool <- struct{}{}

		go func(obj *s3.Object) {
			defer wg.Done()
			defer func() { <-workerPool }()

			processedLogs, err := p.ProcessModelInvocationLogS3Object(*obj.Key, p.processLog)
			if err != nil {
				log.Printf("Error processing S3 object: %s, error:%v\n", *obj.Key, err)
				return
			}

			err = p.uploadS3Object(*obj.Key, processedLogs)
			if err != nil {
				log.Printf("Error uploading S3 object: %s, error:%v\n", *obj.Key, err)
				return
			}
		}(obj)
	}

	wg.Wait()

	return nil
}

func (p *Processor) ProcessModelInvocationLogS3Object(sourceKey string, logProcessorFunc func([]byte) ([]byte, error)) ([]byte, error) {
	result, err := p.s3ClientRead.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(p.modelInvocationLogsInputBucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object %q from bucket %q, %v", sourceKey, p.modelInvocationLogsInputBucket, err)
	}
	defer result.Body.Close()

	var processedLog bytes.Buffer
	scanner := bufio.NewScanner(result.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		logMetadata, err := logProcessorFunc(line)
		if err != nil {
			log.Printf("Error processing log: %v\n", err)
			continue
		}
		processedLog.Write(logMetadata)
		processedLog.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from S3 object: %v", err)
	}

	return processedLog.Bytes(), nil
}

func (p *Processor) listObjectsInDateRange(accountID, region, modelInvocationLogsInputBucketPrefix string, year, month, day, hour int) ([]*s3.Object, error) {
	var datePrefix string
	if modelInvocationLogsInputBucketPrefix == "" {
		datePrefix = fmt.Sprintf("AWSLogs/%s/BedrockModelInvocationLogs/%s/%d/%02d/%02d/%02d", accountID, region, year, month, day, hour)
	} else {
		datePrefix = fmt.Sprintf("%s/AWSLogs/%s/BedrockModelInvocationLogs/%s/%d/%02d/%02d/%02d", modelInvocationLogsInputBucketPrefix, accountID, region, year, month, day, hour)
	}

	paginator := func(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
		return p.s3ClientRead.ListObjects(input)
	}

	params := &s3.ListObjectsInput{
		Bucket: aws.String(p.modelInvocationLogsInputBucket),
		Prefix: aws.String(datePrefix),
	}

	var objects []*s3.Object
	for {
		output, err := paginator(params)
		if err != nil {
			return nil, fmt.Errorf("error listing objects: %w", err)
		}

		for _, obj := range output.Contents {
			objects = append(objects, obj)
		}

		if !*output.IsTruncated {
			break
		}

		params.Marker = output.NextMarker
	}

	return objects, nil
}

func (p *Processor) uploadS3Object(objectKey string, body []byte) error {

	gzippedContent, err := p.gzipContent(body)
	if err != nil {
		return err
	}
	_, err = p.s3ClientWrite.PutObject(&s3.PutObjectInput{
		Bucket:          aws.String(p.metadataLogsOutputBucket),
		Key:             aws.String(objectKey),
		Body:            bytes.NewReader(gzippedContent),
		ContentEncoding: aws.String("gzip"),
		ContentType:     aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("unable to upload object %q to bucket %q, %v", objectKey, p.metadataLogsOutputBucket, err)
	}
	return nil
}

func (p *Processor) processLog(line []byte) ([]byte, error) {
	var modelInvocationLog model.InvocationLog

	err := json.Unmarshal(line, &modelInvocationLog)
	if err != nil {
		return nil, err
	}
	metadata, err := p.modelInvocation.GenerateModelInvocationLogMetadata(&modelInvocationLog)
	if err != nil {
		return nil, err
	}

	transformedContent, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	return transformedContent, nil
}

func (p *Processor) gzipContent(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(input)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
