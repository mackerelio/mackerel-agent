package mackerel

// Host XXX
type Host struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // TODO ENUM
	Status string `json:"status"`
}

type HostSpec struct {
	Name          string                   `json:"name"`
	Meta          map[string]interface{}   `json:"meta"`
	Interfaces    []map[string]interface{} `json:"interfaces"`
	RoleFullnames []string                 `json:"roleFullnames"`
	Checks        []string                 `json:"checks,omitempty"`
	DisplayName   string                   `json:"displayName,omitempty"`
}
