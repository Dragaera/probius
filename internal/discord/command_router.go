package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

const commandPrefix string = "!"

type Command struct {
	Command     string
	Aliases     []string
	Description string
	Usage       string
	MinArgs     int
	MaxArgs     int
	Middleware  []Middleware
	F           func(ctxt CommandContext) bool
}

type Middleware func(cmd Command, ctxt CommandContext) (CommandContext, error)

type CommandRouter struct {
	commands    map[string]Command
	middlewares []Middleware
}

func (router *CommandRouter) register(cmd Command) error {
	if err := router.checkCommandConflict(cmd.Command); err != nil {
		return fmt.Errorf("Command '%v': %v", cmd.Command, err)
	}
	for _, alias := range cmd.Aliases {
		if err := router.checkCommandConflict(alias); err != nil {
			return fmt.Errorf("Alias '%v': %v", alias, err)
		}
	}

	router.commands[cmd.Command] = cmd
	for _, alias := range cmd.Aliases {
		if _, ok := router.commands[alias]; ok {
			// As conflicts with other commands were checked
			// before, this has to be a duplicate alias within the
			// command itself.
			log.Printf(
				"Warning: Ignoring duplicate alias '%v' within command: %v",
				cmd.Command,
				alias,
			)
			continue
		}
		router.commands[alias] = cmd
	}

	return nil
}

func (router *CommandRouter) checkCommandConflict(cmd string) error {
	if conflictCmd, ok := router.commands[cmd]; ok {
		if conflictCmd.Command == cmd {
			return fmt.Errorf(
				"Already in use by other command: %v",
				conflictCmd.Command,
			)
		} else {
			// Conflict must be an alias
			return fmt.Errorf(
				"Already in use as alias of other command: %v",
				conflictCmd.Command,
			)
		}
	}

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

	var ctxt CommandContext
	ctxt = &BaseCommandContext{
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
		// Assigning directly to `ctxt` will lead to a 'declared but not used' error
		newCtxt, err := m(cmd, ctxt)
		ctxt = newCtxt
		if err != nil {
			log.Printf("Bot middleware %v failed: %v. Aborting command.\n", m, err)
			return
		}
	}

	for _, m := range cmd.Middleware {
		// Assigning directly to `ctxt` will lead to a 'declared but not used' error
		newCtxt, err := m(cmd, ctxt)
		ctxt = newCtxt
		if err != nil {
			log.Printf("Command middleware %v failed: %v. Aborting command.\n", m, err)
			return
		}
	}

	if ok := cmd.F(ctxt); !ok {
		ctxt.Respond(usage(&cmd))
	}
}

func usage(cmd *Command) string {
	return fmt.Sprintf("Usage: %v%v\n", commandPrefix, cmd.Usage)
}
