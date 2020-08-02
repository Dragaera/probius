package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) cmdAuth(sess *discordgo.Session, msg *discordgo.Message, args []string) bool {
	var greeting string
	err := bot.db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
	}

	sess.ChannelMessageSend(msg.ChannelID, greeting)
	return true
}
