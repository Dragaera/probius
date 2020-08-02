package persistence

import (
	"time"
)

type DiscordUser struct {
	ID        int
	DiscordID int
	CreatedAt time.Time
}
