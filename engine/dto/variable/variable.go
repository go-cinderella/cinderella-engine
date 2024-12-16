package variable

type VariableSearchRequest struct {
	Name      string `json:"name"`
	Value     any    `json:"value"`
	Operation string `json:"operation"`
	Type      string `json:"type"`
}

type VariableRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`            // java的基础类型 integer | string | ......
	Scope string `json:"scope,omitempty"` // task 使用
	Value any    `json:"value,omitempty"`
}
