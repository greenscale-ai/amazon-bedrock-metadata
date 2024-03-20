package model

import "github.com/greenscale-ai/genai-carbon-footprint/carbonfootprint"

type CarbonFootprintEstimator struct {
	tpd             int
	mem             int
	carbonIntensity int
}

func NewCarbonFootprintEstimator(tpd, mem, carbonIntensity int) *CarbonFootprintEstimator {
	return &CarbonFootprintEstimator{
		tpd:             tpd,
		mem:             mem,
		carbonIntensity: carbonIntensity,
	}
}

func (m *CarbonFootprintEstimator) EstimateModelInvocationCarbonFootprint(metadata *InvocationLogMetadata) *InvocationLogMetadata {
	params := carbonfootprint.Params{
		CarbonIntensity:       float64(m.carbonIntensity),
		TotalInferenceLatency: float64(metadata.InvocationLatency),
		TokenSize:             metadata.OutputTokenCount + metadata.InputTokenCount,
		TDP:                   m.tpd,
		Mem:                   float64(m.mem),
	}

	metadata.EnergyConsumptionkWh, metadata.CarbonEmissiongCO2e, _ = carbonfootprint.CalculateUsageAndEmission(params)

	return metadata
}
