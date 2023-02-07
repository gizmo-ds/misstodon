package models

type MkStreamMessage struct {
	Type string `json:"type"`
}

func (m MkStreamMessage) ToStreamEvent() StreamEvent {
	return StreamEvent{
		Event: m.Type,
	}
}
