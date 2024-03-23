package main

import (
	"time"

	"github.com/google/uuid"
)

var logLines []string

var tick = make(chan bool)

type LobbyStruct struct {
	Players []LobbyPlayerStruct `json:"players"`
	Lobby   LobbyInfoStruct     `json:"lobby"`
	Slots   []SlotDefineStruct  `json:"slots"`
}

type LobbyPlayerStruct struct {
	Name       string        `json:"name"`
	JoinedAt   time.Time     `json:"joinedat"`
	LeftAt     time.Time     `json:"leftat"`
	InGame     bool          `json:"ingame"`
	ID64       string        `json:"id64"`
	Aircraft   string        `json:"aircraft"`
	Copilot    bool          `json:"copilot"`
	Team       string        `json:"team"`
	KillCount  int           `json:"killcount"`
	Active     bool          `json:"active"`
	DeathCount int           `json:"deathcount"`
	Kills      []KillStruct  `json:"kills"`
	Deaths     []DeathStruct `json:"deaths"`
}

var currentPilot string

/* Weapon Type
*	Kill count with weapons
* 	Time of kill
* 	What got killed
*   Who got killed (player name or AI)
*	Who got killed (ID)
* 	Kill info (against AI, or Player)
*   Team the player is on
*	Team the Player got a kill against
 */
type KillStruct struct {
	Weapon     string    `json:"weapon"`
	Time       time.Time `json:"time"`
	Copilot    bool      `json:"copilot"`
	KilledBy   string    `json:"killedby"`
	Killed     string    `json:"killed"`
	KilledName string    `json:"killedbyname"`
	KilledID   string    `json:"killedbyid"`
	PlayerTeam string    `json:"playerteam"`
	UserTeam   string    `json:"userteam"`
}

type DeathStruct struct {
	Weapon       string    `json:"weapon"`
	Time         time.Time `json:"time"`
	DiedWith     string    `json:"diedwith"`
	KilledBy     string    `json:"killedby"`
	KilledByName string    `json:"killedbyname"`
	KilledByID   string    `json:"killedbyid"`
	PlayerTeam   string    `json:"playerteam"`
	UserTeam     string    `json:"userteam"`
}

type LobbyInfoStruct struct {
	PreLobby         LobbyJoinInfoStruct `json:"prelobby"`
	ID               uuid.UUID           `json:"id"`
	Name             string              `json:"name"`
	ID64             string              `json:"id64"`
	HostName         string              `json:"hostname"`
	HostID64         string              `json:"hostid64"`
	TotalLobbyKills  int                 `json:"totallobbykills"`
	TotalLobbyDeaths int                 `json:"totallobbydeaths"`
	WinningTeam      string              `json:"winningteams"`
	JoinTime         time.Time           `json:"jointime"`
	LeaveTime        time.Time           `json:"leavetime"`
	Objectives       []ObjectiveStruct   `json:"objectives"`
}

type ObjectiveStruct struct {
	Name       string    `json:"name"`
	BeganAt    time.Time `json:"beganat"`
	Result     string    `json:"result"`
	ResultedAt time.Time `json:"resultedat"`
}
type LobbyJoinInfoStruct struct {
	LoadedIn      bool   `json:"loadedin"`
	JoinAttempted bool   `json:"joinattempted"`
	JoinReqStatus bool   `json:"joinreqstatus"`
	LobbyInfo     string `json:"lobbyinfo"`
	ScenarioInfo  string `json:"scenarioinfo"`
}

type SlotDefineStruct struct {
	Team     string `json:"team"`
	ID       string `json:"id"`
	Aircraft string `json:"aircraft"`
	Copilot  bool   `json:"copilot"`
}

type TrueSlotStruct struct {
	ID       string `json:"id"`
	Aircraft string `json:"aircraft"`
}

var done = make(chan bool)

var Message = make(chan string)

var LobbyHistory []LobbyStruct

const crewRegex = `(?<=\()\w+(?:,\s*\w+(?:\s+\w+)*)*(?=\))`

const UnitRegex = `[A-Z]+(?:\/[A-Z]+)?-\d+[A-Z]?`

const craftRegex = `(?:AH-94|AV-42C|F-45A|F\/A-26B|EF-24G|T-55)`

const ilstSeparateRegex = `[^,]+`

// const id64RegEx = `\d{17}`

const lobbyHostRegex = `(?<=:)[^()]+(?=\s*\()(?<=\S)`

const ilstRegex = `(\d+,[^,]+,[^,]+,-?\d+,[^,]+,\d+,\d+)`
