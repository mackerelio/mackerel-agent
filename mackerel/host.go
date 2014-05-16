package mackerel

type Host struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // TODO ENUM
	Status string `json:"status"`
}

type Metrics struct {
	Id     string `json:"id"`
	HostId string `json:"hostId"`
	Name   string `json:"name"`
}
