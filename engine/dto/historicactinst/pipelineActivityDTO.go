package historicactinst

type PipelineActivityDTO struct {
	BusinessResult string         `json:"businessResult"`
	RequestBody    map[string]any `json:"requestBody"`
	RequestUrl     string         `json:"requestUrl"`
}
