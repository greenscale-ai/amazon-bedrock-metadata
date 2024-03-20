package model

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestCostEstimator_EstimateModelInvocationCostAnyRegion(t *testing.T) {
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

	costEstimator, err := NewCostEstimator(modelsPriceDetails)

	if err != nil {
		fmt.Println(err)
	}
	metaData := &InvocationLogMetadata{
		Region:           "us-east-1",
		ModelID:          "ai21.j2-mid-v1",
		InputTokenCount:  1000,
		OutputTokenCount: 1000,
	}

	metaData = costEstimator.EstimateModelInvocationCost(metaData)
	if metaData.InputTokenCostUSD != 0.0125 {
		t.Errorf("got %f, wanted %f", metaData.InputTokenCostUSD, 0.0125)
	}

	if metaData.OutputTokenCostUSD != 0.0125 {
		t.Errorf("got %f, wanted %f", metaData.OutputTokenCostUSD, 0.0125)
	}
}

func TestCostEstimator_EstimateModelInvocationCostSpecificRegion(t *testing.T) {
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

	costEstimator, err := NewCostEstimator(modelsPriceDetails)

	if err != nil {
		fmt.Println(err)
	}
	metadata := &InvocationLogMetadata{
		Region:           "us-east-1",
		ModelID:          "amazon.titan-text-express-v1",
		InputTokenCount:  1000,
		OutputTokenCount: 1000,
	}

	metadata = costEstimator.EstimateModelInvocationCost(metadata)
	if metadata.InputTokenCostUSD != 0.0008 {
		t.Errorf("got %f, wanted %f", metadata.InputTokenCostUSD, 0.0008)
	}

	if metadata.OutputTokenCostUSD != 0.0016 {
		t.Errorf("got %f, wanted %f", metadata.OutputTokenCostUSD, 0.0016)
	}
}
