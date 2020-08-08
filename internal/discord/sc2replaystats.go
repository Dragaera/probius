package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/persistence"
	sc2r "github.com/dragaera/probius/internal/sc2replaystats"
	"log"
	"strconv"
	"strings"
	"time"
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

	api := sc2r.API{APIKey: apiKey}
	if !api.Verify() {
		ctxt.Respond("Unable to verify your API key. Please double-check it is correct.")
		return true
	}

	discordId := ctxt.Msg.Author.ID
	user, err := persistence.GetOrCreateSC2ReplayStatsUser(bot.db, ctxt.Msg.Author.ID, apiKey)
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

func (bot *Bot) cmdLast(ctxt CommandContext) bool {
	user, err := persistence.GetSC2ReplayStatsUser(bot.db, ctxt.Msg.Author.ID)
	if err != nil {
		ctxt.Respond("You have not yet granted the bot access to the SC2Replaystats API. Please do so - **in a DM** - with the `!auth` command.")
		return true
	}

	api := sc2r.API{APIKey: user.APIKey}
	replay, err := api.LastReplay()
	if err != nil {
		ctxt.Respond(fmt.Sprintf("An error has happened while contacting the SC2Replaystats API: %v", err))
		return true
	}

	embed := buildReplayEmbed(api, replay)
	ctxt.RespondEmbed(&embed)
	if err != nil {
		ctxt.Respond(fmt.Sprintf("Unable to embed replay: %v", err))
	}

	return true
}

func (bot *Bot) cmdReplay(ctxt CommandContext) bool {
	user, err := persistence.GetSC2ReplayStatsUser(bot.db, ctxt.Msg.Author.ID)
	if err != nil {
		ctxt.Respond("You have not yet granted the bot access to the SC2Replaystats API. Please do so - **in a DM** - with the `!auth` command.")
		return true
	}

	replayId, err := strconv.Atoi(ctxt.Args[0])
	if err != nil {
		ctxt.Respond(fmt.Sprintf("Invalid ID: %v. Must be numeric", ctxt.Args[0]))
		return true
	}

	api := sc2r.API{APIKey: user.APIKey}

	replay, err := api.Replay(replayId)
	if err != nil {
		ctxt.Respond(fmt.Sprintf("An error has happened while contacting the SC2Replaystats API: %v", err))
		return true
	}

	embed := buildReplayEmbed(api, replay)
	err = ctxt.RespondEmbed(&embed)
	if err != nil {
		ctxt.Respond(fmt.Sprintf("Unable to embed replay: %v", err))
	}

	return true
}

func buildReplayEmbed(api sc2r.API, replay sc2r.Replay) discordgo.MessageEmbed {
	mapField := discordgo.MessageEmbedField{
		Name:   "Map",
		Value:  replay.MapName,
		Inline: true,
	}

	winnerField := discordgo.MessageEmbedField{
		Name:   "Winner",
		Value:  fmt.Sprintf("||%v||", replay.WinningPlayer),
		Inline: true,
	}

	gameLengthField := discordgo.MessageEmbedField{
		Name:   "Game Length",
		Value:  fmt.Sprintf("%.0f min", replay.GameLength.Minutes()),
		Inline: false,
	}

	fields := []*discordgo.MessageEmbedField{
		&mapField,
		&winnerField,
		&gameLengthField,
	}

	mapThumbnail := discordgo.MessageEmbedThumbnail{
		URL: mapThumbnailURL(replay.MapName),
	}

	// time.String() returns some wonky go-specific format, while the API
	// obviously expects ISO8601 / RFC3339.
	ts := replay.ReplayDate.Format(time.RFC3339)
	embed := discordgo.MessageEmbed{
		URL:       replay.ReplayURL,
		Title:     constructReplayTitle(api, replay),
		Timestamp: ts,
		Thumbnail: &mapThumbnail,
		Fields:    fields,
	}

	return embed
}

func constructReplayTitle(api sc2r.API, replay sc2r.Replay) string {
	playersByTeam := replay.PlayersByTeam()

	teamMonikers := make([]string, len(playersByTeam))

	for teamId, replayPlayers := range playersByTeam {
		playerMonikers := make([]string, len(replayPlayers))

		for playerIdx, replayPlayer := range replayPlayers {
			playerName := replayPlayer.Player.Name
			if len(playerName) == 0 {
				// `/last-replay` endpoint exposes per-player
				// information, whereas `/replay/$id` endpoint
				// only exposes replay-player information.

				log.Printf(
					"API response did not contain player names, querying API for details of player with ID = %v",
					replayPlayer.ID,
				)

				player, err := api.Player(replayPlayer.ID)
				if err == nil {
					playerName = player.Name
				} else {
					playerName = "Unknown player"
				}
			}

			playerMonikers[playerIdx] = fmt.Sprintf("[%v] %v", replayPlayer.Race.Shorthand(), playerName)
		}

		// Team IDs are 1-based
		teamMonikers[teamId-1] = strings.Join(playerMonikers, ", ")
	}

	return strings.Join(teamMonikers, " vs ")
}

func mapThumbnailURL(mapName string) string {
	// Replace all occurences of space with underscore
	mapName = strings.Replace(mapName, " ", "_", -1)
	// Hosted on SC2ReplaysStats website
	return fmt.Sprintf("https://sc2replaystats.com/images/maps/large/%v.jpg", mapName)
}
