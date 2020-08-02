package discord

import (
	"context"
	"fmt"
)

func (bot *Bot) cmdAuth(ctxt CommandContext) bool {
	var greeting string
	err := bot.db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
	}

	ctxt.Respond(greeting)

	return true
}
