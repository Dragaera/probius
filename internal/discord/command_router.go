package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const commandPrefix string = "!"

type Command struct {
	Command     string
	Description string
	Usage       string
	MinArgs     int
	MaxArgs     int
	F           func(sess *discordgo.Session, msg *discordgo.Message, args []string) bool
}

type CommandRouter struct {
	commands map[string]Command
}

func (router *CommandRouter) register(cmd Command) error {
	if _, ok := router.commands[cmd.Command]; ok {
		return fmt.Errorf("Command already registered: %v", cmd.Command)
	}

	router.commands[cmd.Command] = cmd

	return nil
}

func (router *CommandRouter) onMessageCreate(sess *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message

	if strings.HasPrefix(msg.Content, commandPrefix) {
		router.processCommand(sess, msg)
	}
}

func (router *CommandRouter) processCommand(sess *discordgo.Session, msg *discordgo.Message) {
	// Get rid of prefix
	cmdString := strings.Replace(msg.Content, commandPrefix, "", 1)
	command := strings.Split(cmdString, " ")
	args := command[1:]

	cmd, ok := router.commands[command[0]]
	if !ok {
		fmt.Println("Unknown command:", command[0])
		return
	}

	if cmd.MinArgs != -1 && len(args) < cmd.MinArgs {
		respond(sess, msg.ChannelID, usage(&cmd))
		return
	}

	if cmd.MaxArgs != -1 && len(args) > cmd.MaxArgs {
		respond(sess, msg.ChannelID, usage(&cmd))
		return
	}

	if ok := cmd.F(sess, msg, args); !ok {
		respond(sess, msg.ChannelID, usage(&cmd))
	}
}

func respond(sess *discordgo.Session, channelID string, msg string) {
	sess.ChannelMessageSend(
		channelID,
		msg,
	)
}

func usage(cmd *Command) string {
	return fmt.Sprintf("Usage: %v%v\n", commandPrefix, cmd.Usage)
}
