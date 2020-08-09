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
	Middleware  []Middleware
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

	ctxt := BaseCommandContext{
		sess: sess,
		msg:  msg,
		args: args,
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
		err := m(cmd, &ctxt)
		if err != nil {
			log.Printf("Bot middleware %v failed: %v. Aborting command.\n", m, err)
			return
		}
	}

	for _, m := range cmd.Middleware {
		err := m(cmd, &ctxt)
		if err != nil {
			log.Printf("Command middleware %v failed: %v. Aborting command.\n", m, err)
			return
		}
	}

	if ok := cmd.F(&ctxt); !ok {
		ctxt.Respond(usage(&cmd))
	}
}

type CommandContext interface {
	SetSess(*discordgo.Session)
	Sess() *discordgo.Session

	SetMsg(*discordgo.Message)
	Msg() *discordgo.Message

	SetArgs([]string)
	Args() []string

	SetGuild(*persistence.DiscordGuild)
	Guild() *persistence.DiscordGuild

	SetChannel(*persistence.DiscordChannel)
	Channel() *persistence.DiscordChannel

	SetUser(*persistence.DiscordUser)
	User() *persistence.DiscordUser

	Respond(string) error
	RespondEmbed(*discordgo.MessageEmbed) error
	InternalError(error) error
}

type BaseCommandContext struct {
	sess    *discordgo.Session
	msg     *discordgo.Message
	args    []string
	guild   *persistence.DiscordGuild
	channel *persistence.DiscordChannel
	user    *persistence.DiscordUser
}

func (ctxt *BaseCommandContext) SetSess(sess *discordgo.Session) {
	ctxt.sess = sess
}
func (ctxt *BaseCommandContext) Sess() *discordgo.Session {
	return ctxt.sess
}

func (ctxt *BaseCommandContext) SetMsg(msg *discordgo.Message) {
	ctxt.msg = msg
}
func (ctxt *BaseCommandContext) Msg() *discordgo.Message {
	return ctxt.msg
}

func (ctxt *BaseCommandContext) SetArgs(args []string) {
	ctxt.args = args
}
func (ctxt *BaseCommandContext) Args() []string {
	return ctxt.args
}

func (ctxt *BaseCommandContext) SetGuild(guild *persistence.DiscordGuild) {
	ctxt.guild = guild
}
func (ctxt *BaseCommandContext) Guild() *persistence.DiscordGuild {
	return ctxt.guild
}

func (ctxt *BaseCommandContext) SetChannel(channel *persistence.DiscordChannel) {
	ctxt.channel = channel
}
func (ctxt *BaseCommandContext) Channel() *persistence.DiscordChannel {
	return ctxt.channel
}

func (ctxt *BaseCommandContext) SetUser(user *persistence.DiscordUser) {
	ctxt.user = user
}
func (ctxt *BaseCommandContext) User() *persistence.DiscordUser {
	return ctxt.user
}

func (ctxt *BaseCommandContext) Respond(msg string) error {
	_, err := ctxt.Sess().ChannelMessageSend(
		ctxt.Msg().ChannelID,
		msg,
	)

	if err != nil {
		log.Printf("Error while responding with message: %v", err)
	}

	return err
}

func (ctxt *BaseCommandContext) RespondEmbed(embed *discordgo.MessageEmbed) error {
	_, err := ctxt.Sess().ChannelMessageSendEmbed(
		ctxt.Msg().ChannelID,
		embed,
	)

	if err != nil {
		log.Printf("Error while responding with embed: %v", err)
	}

	return err
}

func (ctxt *BaseCommandContext) InternalError(err error) error {
	msg := fmt.Sprintf(
		"An internal error has happened while performing this operation.\nPlease report the following to 'Morrolan#3163':\n`%v`",
		err,
	)
	log.Printf("Internal error: %v", err)
	ctxt.Respond(msg)

	return err
}

func usage(cmd *Command) string {
	return fmt.Sprintf("Usage: %v%v\n", commandPrefix, cmd.Usage)
}
