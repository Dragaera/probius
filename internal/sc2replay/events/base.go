package events

// Used only for composition
type BaseEvent struct {
	ID   int `json:"id"`
	Loop int `json:"loop"`
	// EventType EventType `json:"evtTypeName"`
}
