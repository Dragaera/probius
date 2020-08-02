package discord

import (
	"fmt"
	"github.com/dragaera/probius/internal/persistence"
)

func (bot *Bot) cmdAuth(ctxt CommandContext) bool {
	apiKey := ctxt.Args[0]

	isDM, err := ctxt.IsDM()
	if err != nil {
		ctxt.InternalError(err)
		return true
	}
	if !isDM {
		ctxt.Respond("**Only use this command via DM**.\nIf you used this command in a public channel, you might have just exposed your API key. If so, please reset it on the profile page and try again via direct message.")
		return true
	}

	discordId := ctxt.Msg.Author.ID
	user, err := persistence.GetOrCreateSC2ReplayStatsUser(bot.db, discordId, apiKey)
	if err != nil {
		ctxt.InternalError(err)
		return true
	}

	if user.APIKey != apiKey {
		err = user.UpdateAPIKey(bot.db, apiKey)
		if err != nil {
			ctxt.InternalError(err)
			return true
		}
	}

	ctxt.Respond(fmt.Sprintf("Successfully set API key of Discord user %v to %v", discordId, apiKey))
	return true
}
