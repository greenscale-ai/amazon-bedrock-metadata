package model

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/arn"
	"log"
	"strings"
)

type MetadataGenerator struct {
	modelCost       *CostEstimator
	carbonFootprint *CarbonFootprintEstimator
	identity        *IdentityTagsBuilder
}

func NewMetadataGenerator(modelCost *CostEstimator, carbonFootprint *CarbonFootprintEstimator, identity *IdentityTagsBuilder) *MetadataGenerator {
	return &MetadataGenerator{
		modelCost:       modelCost,
		carbonFootprint: carbonFootprint,
		identity:        identity,
	}
}

func (m *MetadataGenerator) GenerateModelInvocationLogMetadata(modelInvocationLog *InvocationLog) (modelInvocationLogMetadata *InvocationLogMetadata, err error) {
	var modelId string
	if arn.IsARN(modelInvocationLog.ModelID) {
		modelIdARN, err := arn.Parse(modelInvocationLog.ModelID)
		if err != nil {
			return modelInvocationLogMetadata, err
		}
		modelId = strings.Split(modelIdARN.Resource, "/")[1]
		modelId = strings.Split(modelId, ":")[0]
	} else {
		modelId = modelInvocationLog.ModelID
	}

	modelInvocationLogMetadata = &InvocationLogMetadata{
		ModelID:           modelId,
		Timestamp:         modelInvocationLog.Timestamp,
		AccountID:         modelInvocationLog.AccountID,
		Region:            modelInvocationLog.Region,
		RequestID:         modelInvocationLog.RequestID,
		Operation:         modelInvocationLog.Operation,
		Identity:          modelInvocationLog.Identity,
		InputContentType:  modelInvocationLog.Input.InputContentType,
		OutputContentType: modelInvocationLog.Output.OutputContentType,
		InputTokenCount:   modelInvocationLog.Input.InputTokenCount,
		OutputTokenCount:  modelInvocationLog.Output.OutputTokenCount,
	}

	modelInvocationLogMetadata = m.modelCost.EstimateModelInvocationCost(modelInvocationLogMetadata)

	// parse additional latency metrics related to streaming operation
	if modelInvocationLog.Operation == "InvokeModelWithResponseStream" {
		outputBody, err := json.Marshal(modelInvocationLog.Output.OutputBodyJSON)
		if err == nil {
			outputBodyPayload := make([]InvocationLogOutputBodyJSON, 0)
			err = json.Unmarshal(outputBody, &outputBodyPayload)
			if err == nil {
				for _, output := range outputBodyPayload {
					if output.AmazonBedrockInvocationMetrics.InvocationLatency != 0 && output.AmazonBedrockInvocationMetrics.FirstByteLatency != 0 {
						modelInvocationLogMetadata.InvocationLatency = output.AmazonBedrockInvocationMetrics.InvocationLatency
						modelInvocationLogMetadata.FirstByteLatency = output.AmazonBedrockInvocationMetrics.FirstByteLatency
						modelInvocationLogMetadata = m.carbonFootprint.EstimateModelInvocationCarbonFootprint(modelInvocationLogMetadata)
						break
					}
				}
			}
		}
	}

	// Get IAM identity tags
	identityTags, err := m.identity.GetIdentityTags(modelInvocationLog.Identity.Arn)
	if err != nil {
		log.Printf("unable to get tags for %s\n", modelInvocationLog.Identity.Arn)
	}

	for _, identityTag := range identityTags {
		modelInvocationLogMetadata.IdentityTags = append(modelInvocationLogMetadata.IdentityTags, IdentityTag{Key: *identityTag.Key, Value: *identityTag.Value})
	}

	return modelInvocationLogMetadata, nil
}
