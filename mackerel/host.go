package mackerel

// Host XXX
type Host struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // TODO ENUM
	Status string `json:"status"`
}
