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
	F           func(ctxt CommandContext) bool
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

	ctxt := CommandContext{
		Sess: sess,
		Msg:  msg,
		Args: args,
	}

	if cmd.MinArgs != -1 && len(args) < cmd.MinArgs {
		ctxt.Respond(usage(&cmd))
		return
	}

	if cmd.MaxArgs != -1 && len(args) > cmd.MaxArgs {
		ctxt.Respond(usage(&cmd))
		return
	}

	if ok := cmd.F(ctxt); !ok {
		ctxt.Respond(usage(&cmd))
	}
}

type CommandContext struct {
	Sess *discordgo.Session
	Msg  *discordgo.Message
	Args []string
}

func (ctxt *CommandContext) Respond(msg string) error {
	_, err := ctxt.Sess.ChannelMessageSend(
		ctxt.Msg.ChannelID,
		msg,
	)

	return err
}

func usage(cmd *Command) string {
	return fmt.Sprintf("Usage: %v%v\n", commandPrefix, cmd.Usage)
}
