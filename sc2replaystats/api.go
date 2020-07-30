package sc2replaystats

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const baseURL string = "https://api.sc2replaystats.com"

type Replay struct {
	ReplayURL     string         `json:"replay_url"`
	ReplayId      int            `json:"replay_id"`
	MapName       string         `json:"map_name"`
	Format        string         `json:"format"`
	GameType      GameType       `json:"game_type"`
	Players       []ReplayPlayer `json:"players"`
	SeasonId      int            `json:"seasons_id"`
	ReplayDate    time.Time      `json:"replay_date"`
	ReplayVersion string         `json:"replay_version"`
}

func (replay *Replay) Winner() *ReplayPlayer {
	if replay.Players[0].Winner {
		return &(replay.Players[0])
	} else {
		return &(replay.Players[1])
	}
}

type ReplayPlayer struct {
	Id         int         `json:"players_id"`
	ClanTag    string      `json:"clan"`
	Race       Race        `json:"race"`
	Mmr        int         `json:"mmr"`
	Division   string      `json:"division"`
	ServerRank int         `json:"server_rank"`
	GlobalRank int         `json:"global_rank"`
	Apm        int         `json:"apm"`
	Team       int         `json:"team"`
	Winner     SC2RBool    `json:"winner"`
	Color      PlayerColor `json:"color"`
}

type Player struct {
	Id           int    `json:"players_id"`
	Name         string `json:"players_name"`
	BattleNetURL string `json:"battle_net_url"`
	Clan         Clan   `json:"clan"`
}

type Clan struct {
	Id   int    `json:"clans_id"`
	Name string `json:"clan_name"`
	Tag  string `json:"clan_tag"`
}

type PlayerColor color.RGBA

func (playerColor *PlayerColor) UnmarshalJSON(data []byte) error {
	var colorString string
	err := json.Unmarshal(data, &colorString)
	if err != nil {
		return err
	}

	components := strings.Split(colorString, ",")
	if len(components) != 3 {
		return fmt.Errorf("Invalid color string: %v", colorString)
	}

	r, err := strconv.ParseUint(components[0], 10, 8)
	if err != nil {
		return fmt.Errorf("Invalid color string (R): %v", colorString)
	}

	g, err := strconv.ParseUint(components[1], 10, 8)
	if err != nil {
		return fmt.Errorf("Invalid color string (G): %v", colorString)
	}

	b, err := strconv.ParseUint(components[2], 10, 8)
	if err != nil {
		return fmt.Errorf("Invalid color string (B): %v", colorString)
	}

	// Numbers are uint64, but guaranteed (third parameter of ParseUint) to
	// fit in 8bit uint
	*playerColor = PlayerColor{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}

	return nil
}

type Race int

const (
	Protoss Race = iota
	Terran
	Zerg
)

func (race Race) String() string {
	switch race {
	case Protoss:
		return "Protoss"
	case Terran:
		return "Terran"
	case Zerg:
		return "Zerg"
	default:
		return "Unknown"
	}
}

func (race *Race) UnmarshalJSON(b []byte) error {
	var raceString string
	err := json.Unmarshal(b, &raceString)
	if err != nil {
		return err
	}

	switch raceString {
	case "P":
		*race = Protoss
	case "T":
		*race = Terran
	case "Z":
		*race = Zerg
	default:
		return fmt.Errorf("Invalid race: %v", raceString)
	}

	return nil
}

// SC2Replaystats uses `0` for false, `1` for true. Using a custom type allows
// us to have a custom UnmarshalJSON method.
type SC2RBool bool

func (winner *SC2RBool) UnmarshalJSON(b []byte) error {
	var winnerInt int
	err := json.Unmarshal(b, &winnerInt)
	if err != nil {
		return err
	}

	switch winnerInt {
	case 0:
		*winner = false
	case 1:
		*winner = true
	default:
		return fmt.Errorf("Invalid winner: %v", winnerInt)
	}

	return nil
}

type GameType int

const (
	AutoMM GameType = iota
	AutoMMAI
	AutoMMArchon
	Private
)

func (gameType GameType) String() string {
	switch gameType {
	case AutoMM:
		return "AutoMM"
	case AutoMMAI:
		return "AutoMM AI"
	case AutoMMArchon:
		return "AutoMM Archon"
	case Private:
		return "Private"
	default:
		return "Unknown"
	}
}

func (gameType *GameType) UnmarshalJSON(b []byte) error {
	var gameTypeString string
	err := json.Unmarshal(b, &gameTypeString)
	if err != nil {
		return err
	}

	switch gameTypeString {
	case "AutoMM":
		*gameType = AutoMM
	case "AutoMM_AI":
		*gameType = AutoMMAI
	case "AutoMM_Archon":
		*gameType = AutoMMArchon
	case "Private":
		*gameType = Private
	default:
		return fmt.Errorf("Invalid game type: %v", gameTypeString)
	}

	return nil
}

type API struct {
	APIKey string
	client *http.Client
}

func (api *API) LastReplay() (Replay, error) {
	var replay Replay

	body, err := api.call("account/last-replay")
	if err != nil {
		// No need to wrap, errors returned by call() should be plenty
		// descriptive
		return replay, err
	}

	err = json.Unmarshal(body, &replay)
	if err != nil {
		return replay, fmt.Errorf("Error while unmarshalling JSON: %v", err)
	}

	return replay, nil
}

func (api *API) Player(playerId int) (Player, error) {
	var player Player

	body, err := api.call(fmt.Sprintf("player/%v", playerId))
	if err != nil {
		// No need to wrap, errors returned by call() should be plenty
		// descriptive
		return player, err
	}

	err = json.Unmarshal(body, &player)
	if err != nil {
		return player, fmt.Errorf("Error while unmarshalling JSON: %v", err)
	}

	return player, nil
}

func (api *API) call(path string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v", baseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("Error while creating HTTP request: %v", err)
	}

	req.Header.Add("Authorization", api.APIKey)
	resp, err := api.getClient().Do(req)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("Error while performing HTTP request: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("Error while reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return make([]byte, 0), fmt.Errorf("Error while calling API: Status code = %v, body = %v", resp.StatusCode, string(body))
	}

	return body, nil
}

func (api *API) getClient() *http.Client {
	if api.client == nil {
		api.client = &http.Client{}
	}

	return api.client
}
