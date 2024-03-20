package model

import (
	"testing"
)

func TestEstimateModelInvocationCarbonFootprint(t *testing.T) {
	metadata := &InvocationLogMetadata{
		Region:            "us-east-1",
		ModelID:           "amazon.titan-text-express-v1",
		InputTokenCount:   1000,
		OutputTokenCount:  1000,
		InvocationLatency: 2600,
		FirstByteLatency:  1000,
	}

	modelCarbonFootprint := NewCarbonFootprintEstimator(400, 768000, 450)
	modelCarbonFootprint.EstimateModelInvocationCarbonFootprint(metadata)

	if metadata.EnergyConsumptionkWh != 0.00034463104 {
		t.Errorf("got %f, wanted %f", metadata.EnergyConsumptionkWh, 0.0008)
	}

	if metadata.CarbonEmissiongCO2e != 0.155083968 {
		t.Errorf("got %f, wanted %f", metadata.CarbonEmissiongCO2e, 0.0016)
	}
}
