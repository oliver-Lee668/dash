package models

// DynamicListDto dynamic 列表参数
type DynamicListDto struct {
	Namespace string `json:"namespace" form:"namespace"`
	Group     string `json:"group" form:"group"`
	Version   string `json:"version" form:"version"`
	Resource  string `json:"resource" form:"resource"`
	Name      string `json:"name" form:"name"`
	Kind      string `json:"kind" form:"kind"`
}
