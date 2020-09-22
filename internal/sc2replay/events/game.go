package events

type GameUserLeave struct {
	BaseEvent
	LeaveReason int    `json:"leaveReason"`
	UserID      UserID `json:"userid"`
}

// No, I don't know why it's nested like this. Ask Blizzard.
type UserID struct {
	UserID int64 `json:"userId"`
}
