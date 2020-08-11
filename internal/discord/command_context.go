package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/persistence"
	"log"
)

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
