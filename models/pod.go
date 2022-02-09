package models

type PodDto struct {
	Namespace      string            `json:"namespace"`
	Name           string            `json:"name"`
	Label          map[string]string `json:"label"`
	ContainerName  string            `json:"containerName"`
	ContainerImage string            `json:"containerImage"`
}
