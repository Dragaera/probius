package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Bot struct {
	ClientID  string
	Token     string
	Session   *discordgo.Session
	cmdRouter *CommandRouter
}

func (bot *Bot) Run() error {
	if bot.Session == nil {
		return fmt.Errorf("Bot not initiated, be sure to use discord.Create(...)")
	}

	err := bot.Session.Open()
	if err != nil {
		return fmt.Errorf("Error connecting to Discord:", err)
	}

	fmt.Println("Bot is running.")
	fmt.Println("Invite me:", bot.InviteURL())

	// Terminate on ^c or SIGTERM
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Session.Close()

	return nil
}

func Create(bot *Bot) (*Bot, error) {
	if len(bot.ClientID) == 0 {
		return bot, fmt.Errorf("ClientID must not be nil.")
	}

	if len(bot.Token) == 0 {
		return bot, fmt.Errorf("Token must not be nil.")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		return bot, fmt.Errorf("Error creating Discord session:", err)
	}
	bot.Session = dg

	// Specify intents, limiting data we receive.
	dg.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages,
	)

	// Prepare command router
	router := CommandRouter{}
	router.commands = make(map[string]Command)
	bot.cmdRouter = &router
	dg.AddHandler(bot.cmdRouter.onMessageCreate)

	// And hook up commands
	bot.registerCommands()

	return bot, nil
}

func (bot *Bot) registerCommands() {
	bot.cmdRouter.register(
		Command{
			Command:     "help",
			Description: "Show help about commands",
			Usage:       "help",
			F:           bot.cmdHelp,
		},
	)

	bot.cmdRouter.register(
		Command{
			Command:     "auth",
			Description: "Authorize the bot to access the SC2replaystats.com API on your behalf",
			Usage:       "auth <api_key>",
			F:           bot.cmdAuth,
		},
	)
}

func (bot *Bot) InviteURL() string {
	return fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%v&scope=bot", bot.ClientID)
}

func (bot *Bot) cmdHelp(sess *discordgo.Session, msg *discordgo.Message, args []string) bool {
	out := strings.Builder{}

	if len(args) == 0 {
		// Command list
		out.WriteString("Available commands:\n")
		for _, cmd := range bot.cmdRouter.commands {
			fmt.Fprintf(
				&out,
				"\t`%v%v`: %v\n",
				commandPrefix,
				cmd.Command,
				cmd.Description,
			)
		}
	} else if len(args) == 1 {
		// Help about one command
		if cmd, ok := bot.cmdRouter.commands[args[0]]; ok {
			fmt.Fprintf(&out, "%v: %v\n", cmd.Command, cmd.Description)
			fmt.Fprintf(&out, "\tUsage: `%v%v`\n", commandPrefix, cmd.Usage)
		} else {
			fmt.Fprintf(&out, "Unknown command: `%v`\n", args[0])
		}
	} else {
		return false
	}

	sess.ChannelMessageSend(msg.ChannelID, out.String())
	return true
}

func (bot *Bot) cmdAuth(sess *discordgo.Session, msg *discordgo.Message, args []string) bool {
	return true
}
