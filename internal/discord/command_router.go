package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/persistence"
	"log"
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

type Middleware func(cmd Command, ctxt CommandContext) error

type CommandRouter struct {
	commands    map[string]Command
	middlewares []Middleware
}

func (router *CommandRouter) register(cmd Command) error {
	if _, ok := router.commands[cmd.Command]; ok {
		return fmt.Errorf("Command already registered: %v", cmd.Command)
	}

	router.commands[cmd.Command] = cmd

	return nil
}

func (router *CommandRouter) registerMiddleware(m Middleware) {
	router.middlewares = append(router.middlewares, m)
}

func (router *CommandRouter) onMessageCreate(sess *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message

	if strings.HasPrefix(msg.Content, commandPrefix) {
		router.processCommand(sess, msg)
	}
}

func (router *CommandRouter) processCommand(sess *discordgo.Session, msg *discordgo.Message) {
	log.Printf("Processing command: %v", msg.Content)
	// Get rid of prefix
	cmdString := strings.Replace(msg.Content, commandPrefix, "", 1)
	command := strings.Split(cmdString, " ")
	args := command[1:]

	cmd, ok := router.commands[command[0]]
	if !ok {
		log.Printf("Unknown command: %v", command[0])
		return
	}

	ctxt := CommandContext{
		Sess: sess,
		Msg:  msg,
		Args: args,
	}

	if cmd.MinArgs != -1 && len(args) < cmd.MinArgs {
		log.Printf("Command %v: Too few arguments.", command[0])
		ctxt.Respond(usage(&cmd))
		return
	}

	if cmd.MaxArgs != -1 && len(args) > cmd.MaxArgs {
		log.Printf("Command %v: Too many arguments.", command[0])
		ctxt.Respond(usage(&cmd))
		return
	}

	for _, m := range router.middlewares {
		err := m(cmd, ctxt)
		if err != nil {
			log.Printf("Middleware %v failed: %v. Aborting command.\n", m, err)
			ctxt.InternalError(err)
			return
		}
	}

	if ok := cmd.F(ctxt); !ok {
		ctxt.Respond(usage(&cmd))
	}
}

type CommandContext struct {
	Sess    *discordgo.Session
	Msg     *discordgo.Message
	Args    []string
	Guild   *persistence.DiscordGuild
	Channel *persistence.DiscordChannel
	User    *persistence.DiscordUser
}

func (ctxt *CommandContext) Respond(msg string) error {
	_, err := ctxt.Sess.ChannelMessageSend(
		ctxt.Msg.ChannelID,
		msg,
	)

	if err != nil {
		log.Printf("Error while responding with message: %v", err)
	}

	return err
}

func (ctxt *CommandContext) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := ctxt.Sess.ChannelMessageSendEmbed(
		ctxt.Msg.ChannelID,
		embed,
	)

	if err != nil {
		log.Printf("Error while responding with embed: %v", err)
	}

	return err
}

// Get channel details from API
func (ctxt *CommandContext) GetChannel() (*discordgo.Channel, error) {
	return ctxt.Sess.Channel(ctxt.Msg.ChannelID)
}

func (ctxt *CommandContext) IsDM() (bool, error) {
	channel, err := ctxt.GetChannel()
	if err != nil {
		return false, err
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func (ctxt *CommandContext) InternalError(err error) {
	msg := fmt.Sprintf(
		"An internal error has happened while performing this operation.\nPlease report the following to 'Morrolan#3163':\n`%v`",
		err,
	)
	log.Printf("Internal error: %v", err)
	ctxt.Respond(msg)
}

func usage(cmd *Command) string {
	return fmt.Sprintf("Usage: %v%v\n", commandPrefix, cmd.Usage)
}
