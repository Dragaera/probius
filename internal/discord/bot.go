package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/config"
	"github.com/dragaera/probius/internal/persistence"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Bot struct {
	Config    config.Config
	Session   *discordgo.Session
	cmdRouter *CommandRouter
	db        *pgxpool.Pool
}

func (bot *Bot) Run(db *pgxpool.Pool) error {
	if bot.Session == nil {
		return fmt.Errorf("Bot not initiated, be sure to use discord.Create(...)")
	}

	bot.db = db
	defer bot.db.Close()

	err := bot.Session.Open()
	if err != nil {
		return fmt.Errorf("Error connecting to Discord:", err)
	}

	log.Print("Bot is running.")
	log.Printf("Invite me: %v", bot.InviteURL())

	// Terminate on ^c or SIGTERM
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Session.Close()

	return nil
}

func Create(bot *Bot) (*Bot, error) {
	if len(bot.Config.Discord.ClientID) == 0 {
		return bot, fmt.Errorf("ClientID must not be nil.")
	}

	if len(bot.Config.Discord.Token) == 0 {
		return bot, fmt.Errorf("Token must not be nil.")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + bot.Config.Discord.Token)
	if err != nil {
		return bot, fmt.Errorf("Error creating Discord session:", err)
	}
	bot.Session = dg

	// Specify intents, limiting data we receive.
	dg.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages,
	)

	if err = bot.initializeCommands(); err != nil {
		return nil, err
	}

	return bot, nil
}

func (bot *Bot) initializeCommands() error {
	// Prepare command router
	router := CommandRouter{}
	router.commands = make(map[string]Command)
	router.middlewares = make([]Middleware, 0)

	bot.cmdRouter = &router
	bot.Session.AddHandler(bot.cmdRouter.onMessageCreate)

	// And hook up commands and middlewares
	bot.registerCommands()
	bot.registerMiddlewares()

	return nil
}

func (bot *Bot) registerMiddlewares() {
	bot.cmdRouter.registerMiddleware(bot.enrichContext)
}

func (bot *Bot) registerCommands() {
	bot.cmdRouter.register(
		Command{
			Command:     "help",
			Description: "Show help about commands",
			Usage:       "help [command]",
			MinArgs:     0,
			MaxArgs:     1,
			F:           bot.cmdHelp,
		},
	)

	bot.cmdRouter.register(
		Command{
			Command:     "auth",
			Description: "Authorize the bot to access the SC2replaystats.com API on your behalf",
			Usage:       "auth <api_key>",
			MinArgs:     1,
			MaxArgs:     1,
			F:           bot.cmdAuth,
		},
	)

	bot.cmdRouter.register(
		Command{
			Command:     "last",
			Description: "Embeds the most-recently uploaded replay",
			Usage:       "last",
			MinArgs:     0,
			MaxArgs:     0,
			F:           bot.cmdLast,
		},
	)

	bot.cmdRouter.register(
		Command{
			Command:     "replay",
			Description: "Embds the replay with the given ID",
			Usage:       "replay <id>",
			MinArgs:     1,
			MaxArgs:     1,
			F:           bot.cmdReplay,
		},
	)
}

func (bot *Bot) InviteURL() string {
	return fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%v&scope=bot", bot.Config.Discord.ClientID)
}

func (bot *Bot) cmdHelp(ctxt CommandContext) bool {
	out := strings.Builder{}

	switch len(ctxt.Args) {
	case 0:
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
	case 1:
		cmdIdentifier := ctxt.Args[0]
		// Help about one command
		if cmd, ok := bot.cmdRouter.commands[cmdIdentifier]; ok {
			fmt.Fprintf(&out, "%v: %v\n", cmd.Command, cmd.Description)
			fmt.Fprintf(&out, "\tUsage: `%v%v`\n", commandPrefix, cmd.Usage)
		} else {
			fmt.Fprintf(&out, "Unknown command: `%v`\n", cmdIdentifier)
		}
	default:
		return false
	}

	ctxt.Sess.ChannelMessageSend(ctxt.Msg.ChannelID, out.String())
	return true
}

func (bot *Bot) enrichContext(cmd Command, ctxt CommandContext) error {
	user, err := persistence.DiscordUserFromDgo(bot.db, ctxt.Msg.Author)
	if err != nil {
		return fmt.Errorf("Unable to enrich context with user: %v", err)
	}
	ctxt.User = &user

	return nil
}
