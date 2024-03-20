package model

import (
	"encoding/json"
)

type CostDetail struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Cost     []struct {
		Region                string  `json:"region"`
		InputCostPer1KTokens  float64 `json:"input_cost_per_1k_tokens"`
		OutputCostPer1KTokens float64 `json:"output_cost_per_1k_tokens"`
	} `json:"cost"`
}

type CostEstimator struct {
	modelCostDetails map[string]*CostDetail
}

func NewCostEstimator(modelsPriceDetails []byte) (*CostEstimator, error) {
	var modelCostDetails map[string]*CostDetail
	err := json.Unmarshal(modelsPriceDetails, &modelCostDetails)
	if err != nil {
		return nil, err
	}

	return &CostEstimator{modelCostDetails: modelCostDetails}, nil
}

func (m *CostEstimator) EstimateModelInvocationCost(metadata *InvocationLogMetadata) *InvocationLogMetadata {
	modelCostDetail, ok := m.modelCostDetails[metadata.ModelID]
	if ok {
		for _, costByRegion := range modelCostDetail.Cost {
			if costByRegion.Region == "any" || costByRegion.Region == metadata.Region {
				metadata.InputTokenCostUSD = (costByRegion.InputCostPer1KTokens / 1000) * float64(metadata.InputTokenCount)
				metadata.OutputTokenCostUSD = (costByRegion.OutputCostPer1KTokens / 1000) * float64(metadata.OutputTokenCount)
				metadata.ModelProvider = modelCostDetail.Provider
				metadata.ModelName = modelCostDetail.Name
				break
			}
		}
	}

	return metadata
}
