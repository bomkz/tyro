package main

import (
	"time"

	"github.com/google/uuid"
)

var logLines []string

var tick = make(chan bool)

type LobbyStruct struct {
	Players []LobbyPlayerStruct
	Lobby   LobbyInfoStruct
	Slots   []SlotDefineStruct
}

type LobbyPlayerStruct struct {
	Name       string
	JoinedAt   time.Time
	LeftAt     time.Time
	ID64       string
	Aircraft   string
	Copilot    bool
	Team       string
	KillCount  int
	Active     bool
	DeathCount int
	Kills      []KillStruct
	Deaths     []DeathStruct
}

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
	Weapon     string
	Time       time.Time
	Copilot    bool
	KilledBy   string
	Killed     string
	KilledName string
	KilledID   string
	PlayerTeam string
	UserTeam   string
}

type DeathStruct struct {
	Weapon       string
	Time         time.Time
	DiedWith     string
	KilledBy     string
	KilledByName string
	KilledByID   string
	PlayerTeam   string
	UserTeam     string
}

type LobbyInfoStruct struct {
	PreLobby         LobbyJoinInfoStruct
	ID               uuid.UUID
	Name             string
	ID64             string
	HostName         string
	HostID64         string
	TotalLobbyKills  int
	TotalLobbyDeaths int
	WinningTeam      string
	JoinTime         time.Time
	LeaveTime        time.Time
	Objectives       []ObjectiveStruct
}

type ObjectiveStruct struct {
	ID          int
	Completed   bool
	CompletedAt time.Time
}
type LobbyJoinInfoStruct struct {
	LoadedIn      bool
	JoinAttempted bool
	JoinReqStatus bool
	LobbyInfo     string
	ScenarioInfo  string
}

type SlotDefineStruct struct {
	Team     string
	ID       string
	Aircraft string
	Copilot  bool
}

type TrueSlotStruct struct {
	ID       string
	Aircraft string
}
