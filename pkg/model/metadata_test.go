package model

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"io"
	"log"
	"os"
	"testing"
)

func TestMetadataGenerator_GenerateModelInvocationLogMetadataInvoke(t *testing.T) {
	input, err := os.Open("test_data/input_invoke.json")
	if err != nil {
		log.Println(err)
		return
	}

	defer input.Close()

	var invocationLog InvocationLog
	invocation, _ := io.ReadAll(input)
	_ = json.Unmarshal(invocation, &invocationLog)

	iamSess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	iamClient := iam.New(iamSess)
	identityTagsBuilder := NewIdentityTagsBuilder(iamClient)

	modelPriceFile, err := os.Open("../../models.json")
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

	modelCostEstimator, err := NewCostEstimator(modelsPriceDetails)
	if err != nil {
		log.Println(err)
		return
	}

	modelCarbonFootprint := NewCarbonFootprintEstimator(400, 768000, 450)

	modelMetaDataGenerator := NewMetadataGenerator(modelCostEstimator, modelCarbonFootprint, identityTagsBuilder)
	metadata, err := modelMetaDataGenerator.GenerateModelInvocationLogMetadata(&invocationLog)

	if fmt.Sprintf("%2f", metadata.InputTokenCostUSD) != "0.000003" {
		t.Errorf("got %2f, wanted %s", metadata.InputTokenCostUSD, "0.000003")
	}

	if fmt.Sprintf("%2f", metadata.OutputTokenCostUSD) != "0.000160" {
		t.Errorf("got %2f, wanted %s", metadata.OutputTokenCostUSD, "0.000160")
	}
}

func TestMetadataGenerator_GenerateModelInvocationLogMetadataInvokeStream(t *testing.T) {
	input, err := os.Open("test_data/input_invoke_stream.json")
	if err != nil {
		log.Println(err)
		return
	}

	defer input.Close()

	var invocationLog InvocationLog
	invocation, _ := io.ReadAll(input)
	_ = json.Unmarshal(invocation, &invocationLog)

	iamSess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	iamClient := iam.New(iamSess)
	identityTagsBuilder := NewIdentityTagsBuilder(iamClient)

	modelPriceFile, err := os.Open("../../models.json")
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

	modelCostEstimator, err := NewCostEstimator(modelsPriceDetails)
	if err != nil {
		log.Println(err)
		return
	}

	modelCarbonFootprint := NewCarbonFootprintEstimator(400, 768000, 450)

	modelMetaDataGenerator := NewMetadataGenerator(modelCostEstimator, modelCarbonFootprint, identityTagsBuilder)
	metadata, err := modelMetaDataGenerator.GenerateModelInvocationLogMetadata(&invocationLog)

	if fmt.Sprintf("%2f", metadata.InputTokenCostUSD) != "0.000006" {
		t.Errorf("got %2f, wanted %s", metadata.InputTokenCostUSD, "0.000006")
	}

	if fmt.Sprintf("%2f", metadata.OutputTokenCostUSD) != "0.000142" {
		t.Errorf("got %2f, wanted %s", metadata.OutputTokenCostUSD, "0.000142")
	}

	if fmt.Sprintf("%2f", metadata.EnergyConsumptionkWh) != "0.001078" {
		t.Errorf("got %2f, wanted %s", metadata.EnergyConsumptionkWh, "0.001078")
	}

	if fmt.Sprintf("%2f", metadata.CarbonEmissiongCO2e) != "0.484876" {
		t.Errorf("got %2f, wanted %s", metadata.CarbonEmissiongCO2e, "0.484876")
	}
}
