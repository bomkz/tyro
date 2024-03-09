package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func tickTime() {
	for {
		tick <- true
		time.Sleep(0 * time.Millisecond)
	}
}

func readLog() {

	// Start log reading ticker
	go tickTime()

	/*
	*Set readlog to false,
	*prevents unwanted media inputs caused by log being read for first time
	*when player is already in game and has pressed buttons.
	 */
	readlog := true

	// Reads log every tick.
	for {
		<-tick

		// Reads log line by line
		logFile := readLogFile()

		// Appends each log line to a variable.
		logLinesTmp := strings.Split(logFile, "\r\n")

		// Only takes action on log lines that are new, and ignores old ones.
		for y, x := range logLinesTmp {
			if y > (len(logLines) - 1) {
				if readlog {
					fmt.Println(y)
					// Handles new log lines.
					logHandler(x)
				}

			}
		}
		logLines = logLinesTmp

		// Set to true after the first for loop passes
		readlog = true
	}
}

// Runs the appropriate function depending on the log line contents.
func logHandler(newline string) bool {
	if strings.Contains(newline, "Attempting to join lobby") {

		go onLobbyJoin()
		InLobby = true
		fmt.Println("bwehhhh")
		Message <- newline
		<-done

	} else if InLobby {
		Message <- newline
		fmt.Println("bwahhhh")
		<-done
	}
	return true

}

var InLobby bool

func UpdateStatus(statusType int) {
	fmt.Println("")
}

func onLobbyJoin() {

	var currentLobby LobbyStruct

	currentLobby.Lobby.ID = uuid.New()
	LobbyHistory = append(LobbyHistory, currentLobby)

	currentLobby.Lobby.WinningTeam = "Invalid"

	for {
		currentMessage := <-Message
		if !currentLobby.Lobby.PreLobby.LoadedIn {

			currentLobby = preLobbyHandler(currentMessage, currentLobby)

		}
		switch {
		case strings.Contains(currentMessage, "- Slot designation for "):
			currentLobby = onSlotDefine(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "- Info:"):
			currentLobby = onNewPlayer(currentMessage, currentLobby)
			currentLobby = onPlayerUpdate(currentLobby)
		case strings.Contains(currentMessage, "RPC_BeginObjective"):
			currentLobby = onBeginObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "RPC_CompleteObjective"):
			currentLobby = onCompleteObjective(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "has disconnected.") && strings.Contains(currentMessage, "FlightLogger: "):
			currentLobby = onPlayerLeave(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Lobby message from ") && strings.Contains(currentMessage, " killed ") && !strings.Contains(currentMessage, " was killed by "):
			currentLobby = onKill(currentMessage, currentLobby)
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && !strings.Contains(currentMessage, "$log_EF-24G") && !strings.Contains(currentMessage, "$log_T-55"):
			currentLobby = onEnvDeath(currentMessage, currentLobby)
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && (strings.Contains(currentMessage, "$log_EF-24G") || strings.Contains(currentMessage, "$log_T-55")):
			currentLobby = onEnvDeathMC(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "OnBriefingSeatUpdated("):
			currentLobby = onBriefingSeatUpdated(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "Setting up slot UI: "):
			currentLobby = onSlotUISetup(currentMessage, currentLobby)
		case currentMessage == "LeaveLobby()":
			InLobby = false
			currentLobby.Lobby.LeaveTime = time.Now()

			for x, y := range LobbyHistory {
				if y.Lobby.ID == currentLobby.Lobby.ID {
					LobbyHistory[x] = currentLobby
					done <- true
					return
				}
			}
		}

		for x, y := range LobbyHistory {
			if y.Lobby.ID == currentLobby.Lobby.ID {
				LobbyHistory[x] = currentLobby
			}
		}
		done <- true

	}

}

func onEnvDeathMC(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	var newDeath DeathStruct
	_, trimmedMessage, _ := strings.Cut(currentMessage, "(")
	tmpname, trimmedMessage, _ := strings.Cut(trimmedMessage, ")")
	tmpname1, tmpname2, found := strings.Cut(tmpname, ", ")
	if !found {
		tmpname1 = tmpname
	}
	_, trimmedMessage, _ = strings.Cut(trimmedMessage, " was killed by ")
	newDeath.Weapon, trimmedMessage, _ = strings.Cut(trimmedMessage, " (")
	newDeath.KilledBy, _ = strings.CutSuffix(trimmedMessage, ".")
	newDeath.KilledBy = "(" + newDeath.KilledBy
	newDeath.KilledByName = "<Environment>"
	for x, y := range currentLobby.Players {
		if (y.Name == tmpname1 || y.Name == tmpname2) && y.Active {
			newDeath.UserTeam = y.Team
			newDeath.DiedWith = y.Aircraft
			newDeath.PlayerTeam = "<environment>"
			newDeath.Time = time.Now()
			currentLobby.Players[x].Deaths = append(currentLobby.Players[x].Deaths, newDeath)
			currentLobby.Players[x].DeathCount += 1
		}
	}

	return currentLobby
}

func onEnvDeath(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	var newDeath DeathStruct
	_, trimmedMessage, _ := strings.Cut(currentMessage, "$log_")
	name, trimmedMessage, _ := strings.Cut(trimmedMessage, " was killed by ")
	newDeath.Weapon, trimmedMessage, _ = strings.Cut(trimmedMessage, " (")
	newDeath.KilledBy, _ = strings.CutSuffix(trimmedMessage, ".")
	newDeath.KilledBy = "(" + newDeath.KilledBy
	newDeath.KilledByName = "<Environment>"
	for x, y := range currentLobby.Players {
		if y.Name == name && y.Active {
			newDeath.UserTeam = y.Team
			newDeath.DiedWith = y.Aircraft
			newDeath.PlayerTeam = "<environment>"
			newDeath.Time = time.Now()
			currentLobby.Players[x].Deaths = append(currentLobby.Players[x].Deaths, newDeath)
			currentLobby.Players[x].DeathCount += 1
		}
	}

	return currentLobby

}

func onSlotDefine(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	slot := strings.TrimPrefix(currentMessage, "- Slot designation for ")
	var newSlot SlotDefineStruct
	newSlot.Team, slot, _ = strings.Cut(slot, ".")
	tmpdata, slot, _ := strings.Cut(slot, "[id=")
	if tmpdata == "0" {
		newSlot.Copilot = true
	}
	newSlot.ID, slot, _ = strings.Cut(slot, "] (")
	newSlot.Aircraft = strings.TrimSuffix(slot, ")")

	if newSlot.ID == "0" {
		newSlot.Aircraft = "Invalid"
		newSlot.Copilot = false
	}

	newSlot.ID = fmt.Sprint(len(currentLobby.Slots))

	currentLobby.Slots = append(currentLobby.Slots, newSlot)

	return currentLobby
}

func onSlotUISetup(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	slot := strings.TrimPrefix(currentMessage, "Setting up slot UI: ")
	team, slot, _ := strings.Cut(slot, " ")
	id, slot, _ := strings.Cut(slot, " ")
	_, name, _ := strings.Cut(slot, ") = ")

	var aircraft string
	var copilot bool
	for _, x := range currentLobby.Slots {
		if x.Team == team && x.ID == id {
			aircraft = x.Aircraft
			copilot = x.Copilot
		}
	}

	for x, y := range currentLobby.Players {
		if y.Name == name && y.Active {
			currentLobby.Players[x].Aircraft = aircraft
			currentLobby.Players[x].Copilot = copilot
		}
	}

	return currentLobby
}

func onBriefingSeatUpdated(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	trimmedMessage, _ := strings.CutPrefix(currentMessage, "BriefingAvatarSync OnBriefingSeatUpdated(")
	ID, trimmedMessage, _ := strings.Cut(trimmedMessage, ", ")
	for x, y := range currentLobby.Players {

		switch {
		case y.ID64 == ID && strings.Contains(trimmedMessage, "Allied") && y.Team == "Enemy":
			var playerExist bool
			for u, z := range currentLobby.Players {
				if z.ID64 == ID && z.Team == "Allied" {
					currentLobby.Players[x].Active = false
					currentLobby.Players[u].Active = true
					playerExist = true
				}
			}
			currentLobby.Players[x].Active = false
			if !playerExist {
				newPlayer := LobbyPlayerStruct{
					ID64:     y.ID64,
					Name:     y.Name,
					JoinedAt: time.Now(),
					Team:     "Allied",
					Active:   true,
				}
				currentLobby.Players = append(currentLobby.Players, newPlayer)

			}
		case y.ID64 == ID && strings.Contains(trimmedMessage, "Enemy") && y.Team == "Allied":
			var playerExist bool
			for u, z := range currentLobby.Players {
				if z.ID64 == ID && z.Team == "Enemy" {
					currentLobby.Players[x].Active = false
					currentLobby.Players[u].Active = true
					playerExist = true
				}
			}
			currentLobby.Players[x].Active = false
			if !playerExist {
				newPlayer := LobbyPlayerStruct{
					ID64:     y.ID64,
					Name:     y.Name,
					JoinedAt: time.Now(),
					Team:     "Enemy",
					Active:   true,
				}
				currentLobby.Players = append(currentLobby.Players, newPlayer)

			}
		}

	}
	return currentLobby
}
func onPlayerLeave(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	trimmedMessage, _ := strings.CutPrefix(currentMessage, "FlightLogger: ")
	name, _ := strings.CutSuffix(trimmedMessage, " has disconnected.")
	for x, y := range currentLobby.Players {
		if y.Name == name && y.LeftAt.IsZero() {
			currentLobby.Players[x].LeftAt = time.Now()
		}
	}
	return currentLobby
}

func onNewPlayer(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	newPlayer := unwrapInfo(strings.TrimPrefix(currentMessage, "- Info: "))
	foundPlayer := false
	for _, x := range currentLobby.Players {
		if x.ID64 == newPlayer.ID64 && x.Team == newPlayer.Team {
			foundPlayer = true
			break
		}
	}
	if !foundPlayer {
		currentLobby.Players = append(currentLobby.Players, newPlayer)
	}
	return currentLobby
}
func onBeginObjective(_ string, currentLobby LobbyStruct) LobbyStruct {
	currentLobby.Lobby.Objectives = append(currentLobby.Lobby.Objectives, ObjectiveStruct{
		ID: len(currentLobby.Lobby.Objectives),
	})
	return currentLobby
}
func onCompleteObjective(currentMessage string, currentLobby LobbyStruct) LobbyStruct {

	_, cutString, _ := strings.Cut(currentMessage, "RPC_CompleteObjective")
	number := strings.TrimSuffix(cutString, ")")
	ID, err := strconv.Atoi(number)
	if err != nil {
		ID = 0
	}
	currentLobby.Lobby.Objectives[ID].Completed = true
	currentLobby.Lobby.Objectives[ID].CompletedAt = time.Now()
	return currentLobby
}

var done = make(chan bool)

func preLobbyHandler(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	switch {
	case strings.Contains(currentMessage, "Attempting to join lobby"):
		var name string
		currentLobby.Lobby.PreLobby.JoinAttempted = true
		currentLobby.Lobby.PreLobby.LobbyInfo, _ = strings.CutPrefix(currentMessage, "Attempting to join lobby ")
		currentLobby.Lobby.ID64, name, _ = strings.Cut(currentLobby.Lobby.PreLobby.LobbyInfo, " (")
		currentLobby.Lobby.ID64 = strings.TrimPrefix(currentLobby.Lobby.ID64, "VTMPMainMenu: Attempting to join lobby ")
		currentLobby.Lobby.Name, _ = strings.CutSuffix(name, ")")
	case strings.Contains(currentMessage, "Join request accepted!"):
		currentLobby.Lobby.PreLobby.JoinReqStatus = true
	case strings.Contains(currentMessage, "getting scenario"):
		currentLobby.Lobby.PreLobby.ScenarioInfo, _ = strings.CutPrefix(currentMessage, "getting scenario ")
	case strings.Contains(currentMessage, "Connecting to host: "):
		var trimmedMessage string
		var ID string
		_, trimmedMessage, _ = strings.Cut(currentMessage, "Creating socket client. Connecting to host: ")
		currentLobby.Lobby.HostName, ID, _ = strings.Cut(trimmedMessage, " (")
		currentLobby.Lobby.HostID64, _, _ = strings.Cut(ID, ")")
	case currentMessage == "Connected":
		currentLobby.Lobby.PreLobby.LoadedIn = true
		currentLobby.Lobby.JoinTime = time.Now()
	}
	return currentLobby
}
func unwrapInfo(currentMessage string) LobbyPlayerStruct {
	var player LobbyPlayerStruct
	var err error
	info := strings.TrimPrefix(currentMessage, " - Info: ")
	player.ID64, info, _ = strings.Cut(info, ",")
	player.Name, info, _ = strings.Cut(info, ",")
	player.Team, info, _ = strings.Cut(info, ",")
	player.Active = true
	_, info, _ = strings.Cut(info, ",")
	_, info, _ = strings.Cut(info, ",")
	pK, info, _ := strings.Cut(info, ",")
	pD, _, _ := strings.Cut(info, ",")

	player.KillCount, err = strconv.Atoi(pK)
	if err != nil {
		player.KillCount = 0
	}
	player.DeathCount, err = strconv.Atoi(pD)
	if err != nil {
		player.KillCount = 0
	}
	player.JoinedAt = time.Now()
	return player
}

func onKill(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	var trimmedMessage string
	var killer string
	var weapon string
	var killed string
	var killedName []string
	var tmpname string
	_, trimmedMessage, _ = strings.Cut(currentMessage, "$log_")
	killer, trimmedMessage, _ = strings.Cut(trimmedMessage, " killed ")
	if strings.Contains(currentMessage, "(") {

		killed, trimmedMessage, _ = strings.Cut(trimmedMessage, " (")
		tmpname, trimmedMessage, _ = strings.Cut(trimmedMessage, ") with ")
		tmpname2, tmpname1, found := strings.Cut(tmpname, ", ")
		if found {
			killedName = append(killedName, tmpname1, tmpname2)
		} else {
			killedName = append(killedName, tmpname)
		}

	} else {
		killed, trimmedMessage, _ = strings.Cut(trimmedMessage, " with ")
	}

	weapon, _ = strings.CutSuffix(trimmedMessage, ".")
	newKill := KillStruct{
		Weapon: weapon,
		Time:   time.Now(),
		Killed: killed,
	}

	for _, h := range killedName {
		for _, y := range currentLobby.Players {
			if y.Name == h && y.Active {
				newKill.KilledID += "(" + y.ID64 + ")"
				newKill.KilledName += "(" + y.Name + ")"
				newKill.PlayerTeam = y.Team
			}
		}
	}
	var aircraft string
	var killerid string
	for x, y := range currentLobby.Players {
		if y.Name == killer && y.Active {
			currentLobby.Players[x].KillCount += 1
			aircraft = y.Aircraft
			newKill.UserTeam = y.Team
			killerid = y.ID64
			newKill.Copilot = y.Copilot
			currentLobby.Players[x].Kills = append(currentLobby.Players[x].Kills, newKill)
		}
	}

	if killedName != nil {
		newDeath := DeathStruct{
			Weapon:       newKill.Weapon,
			Time:         newKill.Time,
			KilledBy:     aircraft,
			KilledByName: killer,
			KilledByID:   killerid,
			PlayerTeam:   newKill.UserTeam,
			UserTeam:     newKill.PlayerTeam,
		}
		for x, y := range currentLobby.Players {
			for _, h := range killedName {
				if y.Name == h && y.Active {
					currentLobby.Players[x].DeathCount += 1
					currentLobby.Players[x].Deaths = append(currentLobby.Players[x].Deaths, newDeath)

				}
			}

		}
	}

	return currentLobby
}

func onPlayerUpdate(currentLobby LobbyStruct) LobbyStruct {
	currentLobby.Lobby.TotalLobbyDeaths = 0
	currentLobby.Lobby.TotalLobbyKills = 0

	for _, x := range currentLobby.Players {
		currentLobby.Lobby.TotalLobbyDeaths += x.DeathCount
		currentLobby.Lobby.TotalLobbyKills += x.KillCount
	}
	return currentLobby
}

var Message = make(chan string)

var LobbyHistory []LobbyStruct
