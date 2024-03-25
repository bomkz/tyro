package main

import (
	"time"

	"github.com/google/uuid"
)

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
	Weapon     string `json:"weapon"`
	Time       int64  `json:"time"`
	Copilot    bool   `json:"copilot"`
	KilledBy   string `json:"killedby"`
	Killed     string `json:"killed"`
	KilledName string `json:"killedbyname"`
	KilledID   string `json:"killedbyid"`
	PlayerTeam string `json:"playerteam"`
	UserTeam   string `json:"userteam"`
}

type DeathStruct struct {
	Weapon       string `json:"weapon"`
	Time         int64  `json:"time"`
	DiedWith     string `json:"diedwith"`
	KilledBy     string `json:"killedby"`
	KilledByName string `json:"killedbyname"`
	KilledByID   string `json:"killedbyid"`
	PlayerTeam   string `json:"playerteam"`
	UserTeam     string `json:"userteam"`
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
	JoinTime         int64               `json:"jointime"`
	LeaveTime        int64               `json:"leavetime"`
	Objectives       []ObjectiveStruct   `json:"objectives"`
}

type ObjectiveStruct struct {
	Name       string `json:"name"`
	BeganAt    int64  `json:"beganat"`
	Result     string `json:"result"`
	ResultedAt int64  `json:"resultedat"`
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

var info = `TYRO ` + Version + `
Telemetry Yield Real-time Observations
Utility to display discord rich presence, and dump your statistics afterwards.

Credits to the following for the:

@kentuckyfrieda10wallsimper - AH-94 Photo on the default rich presence token. 
@dubyaaa - T-55 Photo on the default rich presence token.
@toast2812 - EF-24G and F-45A Photo on the default rich presence token.
@joespeed52 - F/A-26B Photo on the default rich presence token.
@romanian_wallet_inspector - A/V-42C Photo on the default rich presence token.

https://discord.gg/caw8 - Amazing liveries displayed in these photos.

------------------------------------------------------------------------------------------------------
Licensed under MIT
------------------------------------------------------------------------------------------------------

Upon closure, the program will save your entire gameplay session statistics to a timestamp tagged json file,
You can delete or keeep it, however, in the future, I have some projects planned where you can use these files
and output graphs of various varieties as you choose, and also view mission summaries in an easily readable way.
This file is filled with a lot of useful information.
------------------------------------------------------------------------------------------------------
! ! TYRO is starting up, please make sure you are currently not spawned in an aircraft if VTOL VR is already running ! !
------------------------------------------------------------------------------------------------------


`
