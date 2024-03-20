package model

import "time"

type InvocationLog struct {
	SchemaType    string    `json:"schemaType"`
	SchemaVersion string    `json:"schemaVersion"`
	Timestamp     time.Time `json:"timestamp"`
	AccountID     string    `json:"accountId"`
	Identity      struct {
		Arn string `json:"arn"`
	} `json:"identity"`
	Region    string `json:"region"`
	RequestID string `json:"requestId"`
	Operation string `json:"operation"`
	ModelID   string `json:"modelId"`
	Input     struct {
		InputContentType string `json:"inputContentType"`
		InputTokenCount  int    `json:"inputTokenCount"`
	} `json:"input"`
	Output struct {
		OutputContentType string `json:"outputContentType"`
		OutputTokenCount  int    `json:"outputTokenCount"`
		OutputBodyJSON    any    `json:"outputBodyJson"`
	} `json:"output"`
}

type AmazonBedrockInvocationMetrics struct {
	InputTokenCount   int `json:"inputTokenCount"`
	OutputTokenCount  int `json:"outputTokenCount"`
	InvocationLatency int `json:"invocationLatency"`
	FirstByteLatency  int `json:"firstByteLatency"`
}

type InvocationLogOutputBodyJSON struct {
	OutputText                     string                         `json:"outputText"`
	Index                          int                            `json:"index"`
	TotalOutputTextTokenCount      any                            `json:"totalOutputTextTokenCount"`
	CompletionReason               any                            `json:"completionReason"`
	InputTextTokenCount            int                            `json:"inputTextTokenCount"`
	AmazonBedrockInvocationMetrics AmazonBedrockInvocationMetrics `json:"amazon-bedrock-invocationMetrics,omitempty"`
}

type IdentityTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InvocationLogMetadata struct {
	Timestamp time.Time `json:"timestamp"`
	AccountID string    `json:"accountId"`
	Identity  struct {
		Arn string `json:"arn"`
	} `json:"identity"`
	IdentityTags         []IdentityTag `json:"identityTags"`
	Region               string        `json:"region"`
	RequestID            string        `json:"requestId"`
	Operation            string        `json:"operation"`
	ModelID              string        `json:"modelId"`
	ModelName            string        `json:"modelName"`
	ModelProvider        string        `json:"modelProvider"`
	InputContentType     string        `json:"inputContentType"`
	OutputContentType    string        `json:"outputContentType"`
	InputTokenCount      int           `json:"inputTokenCount"`
	OutputTokenCount     int           `json:"outputTokenCount"`
	InputTokenCostUSD    float64       `json:"inputTokenCostUSD"`
	OutputTokenCostUSD   float64       `json:"outputTokenCostUSD"`
	InvocationLatency    int           `json:"invocationLatency,omitempty"`
	FirstByteLatency     int           `json:"firstByteLatency,omitempty"`
	EnergyConsumptionkWh float64       `json:"energyConsumptionkWh,omitempty"`
	CarbonEmissiongCO2e  float64       `json:"carbonEmissiongCO2e,omitempty"`
}
