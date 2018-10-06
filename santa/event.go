package santa

type EventLine struct {
	MachineID string                 `json:"machine_id"`
	Timestamp string                 `json:"timestamp"`
	Event     map[string]interface{} `json:"event"`
}

type EventsList struct {
	Events []map[string]interface{} `json:"events"`
}
