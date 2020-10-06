package common

type Spec struct {
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	Properties map[string]string `json:"properties"`
}
