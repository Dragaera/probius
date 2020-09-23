package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dragaera/probius/internal/sc2replay"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (bot *Bot) cmdSupply(ctxt CommandContext) bool {
	ts := ctxt.Args()[0]
	seconds, err := timestampToSeconds(ts)
	if err != nil {
		ctxt.Respond(err.Error())
		return true
	}

	attachments := ctxt.Msg().Attachments
	if len(attachments) == 0 {
		ctxt.Respond("Replay must be attached to message")
		return true
	}

	// I have not managed to have more than one attachment per message with
	// the Discord client, but it is possible according to the API spec -
	// potentially by other API consumers.
	for _, att := range attachments {
		file, err := downloadFile(att.URL)
		if err != nil {
			ctxt.InternalError(err)
			return true
		}
		defer os.Remove(file)

		report, err := generateReport(file, seconds)
		if err != nil {
			ctxt.Respond(fmt.Sprintf("Error while processing replay: %v", err))
			return true
		}

		embed := buildSupplyEmbed(&report, ts)
		ctxt.RespondEmbed(&embed)
	}

	return true
}

func timestampToSeconds(timestamp string) (int, error) {
	invalidTimestampErr := fmt.Errorf("Timestamp must be of format [HH]:MM:SS")

	parts := strings.Split(timestamp, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, invalidTimestampErr
	}

	hours := 0
	minutes := 0
	seconds := 0
	offset := 0

	if len(parts) == 3 {
		var err error
		hours, err = strconv.Atoi(parts[0])
		if err != nil || hours < 0 {
			return 0, invalidTimestampErr
		}
		// If hours are in position 0, minutes and seconds are in
		// positions 1 and 2 respectively.
		offset = 1
	}

	minutes, err := strconv.Atoi(parts[offset])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, invalidTimestampErr
	}

	seconds, err = strconv.Atoi(parts[offset+1])
	if err != nil || seconds < 0 || seconds > 59 {
		return 0, invalidTimestampErr
	}

	return hours*3600 + minutes*60 + seconds, nil
}

func downloadFile(URL string) (string, error) {
	response, err := http.Get(URL)
	if err != nil {
		return "", fmt.Errorf("Unable to download replay: %v", err)
	}
	defer response.Body.Close()

	// Leaving the target directory empty will default to os.Tempdir
	file, err := ioutil.TempFile("", "probius_supply_replay_")
	if err != nil {
		return "", fmt.Errorf("Unable to create new file: %v", err)
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", fmt.Errorf("Unable to write to file: %v", err)
	}

	return file.Name(), nil
}

func generateReport(file string, duration int) (sc2replay.Report, error) {
	var report = sc2replay.Report{}

	replay, err := sc2replay.FromFile(file)
	if err != nil {
		return report, fmt.Errorf("Unable to load replay: %v\n", err)
	}
	defer replay.Close()

	ticks, err := replay.TicksUntilSeconds(float64(duration))
	if err != nil {
		return report, fmt.Errorf("Unable to determine amount of ticks until 7:00: %v\n", err)
	}

	ownerID, err := replay.OwnerPlayerID()
	if err != nil {
		return report, fmt.Errorf("Unable to determine owner player ID: %v", err)
	}

	report = sc2replay.Report{
		PlayerID: ownerID,
		Replay:   &replay,
	}
	report.At(ticks)

	return report, nil
}

func buildSupplyEmbed(report *sc2replay.Report, timestamp string) discordgo.MessageEmbed {
	timestampField := discordgo.MessageEmbedField{
		Name:   "Timestamp",
		Value:  timestamp,
		Inline: true,
	}

	supplyField := discordgo.MessageEmbedField{
		Name:   "Supply",
		Value:  strconv.Itoa(report.IngameSupply()),
		Inline: true,
	}

	unitField := discordgo.MessageEmbedField{
		Name:   "Units",
		Value:  buildUnitList(report),
		Inline: false,
	}

	fields := []*discordgo.MessageEmbedField{
		&timestampField,
		&supplyField,
		&unitField,
	}

	embed := discordgo.MessageEmbed{
		Title:  "Supply report",
		Fields: fields,
	}

	return embed
}

func buildUnitList(report *sc2replay.Report) string {
	out := strings.Builder{}

	for name, count := range report.UnitCount {
		fmt.Fprintf(
			&out,
			"- %v: %v\n",
			name,
			count,
		)
	}

	return out.String()
}
