package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
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
		Message <- newline
		<-done

	} else if InLobby {
		Message <- newline
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
		case strings.Contains(currentMessage, "Setting up slot UI: "):
			currentLobby = onSlotUISetup(currentMessage, currentLobby)
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
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && !strings.Contains(currentMessage, "$log_EF-24G") && !strings.Contains(currentMessage, "$log_T-55") && !strings.Contains(currentMessage, "$log_AH-94"):
			currentLobby = onEnvDeath(currentMessage, currentLobby)
		case strings.Contains(currentMessage, " was killed by") && strings.Contains(currentMessage, "Lobby message from") && !strings.Contains(currentMessage, "environment") && (strings.Contains(currentMessage, "$log_EF-24G") || strings.Contains(currentMessage, "$log_T-55")) || strings.Contains(currentMessage, "$log_AH-94"):
			currentLobby = onEnvDeathMC(currentMessage, currentLobby)
		case strings.Contains(currentMessage, "identity updated: "):
			currentLobby = onIdentityUpdate(currentMessage, currentLobby)
		case currentMessage == "LeaveLobby()":
			InLobby = false
			currentLobby.Lobby.LeaveTime = time.Now()

			for _, x := range LobbyHistory {
				if x.Lobby.ID == currentLobby.Lobby.ID {
					x = currentLobby
					done <- true
					return
				}
			}
		}

		for _, x := range LobbyHistory {
			if x.Lobby.ID == currentLobby.Lobby.ID {
				x = currentLobby
			}
		}
		done <- true

	}

}

/*
* Function called to handle multicrew player deaths against environment.
* Environment counts as AI, or controlled flight into terrain.
* TODO: tidy up this function and replace spaghetti code with RegEx.
* TODO: capture string examples using breaker point.
 */
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

	// Fill in missing death information using currentLobby.Players information
	for _, x := range currentLobby.Players {
		// If x.Name matches either of tmpname1 or tmpname2, and x.Active is true, then.
		if (x.Name == tmpname1 || x.Name == tmpname2) && x.Active {

			//  Fill in missing information using x
			newDeath.UserTeam = x.Team
			newDeath.DiedWith = x.Aircraft

			// Set killer team as environment (i.e. AI killed user.)
			newDeath.PlayerTeam = "<environment>"

			// Set time of death as time.Now()
			newDeath.Time = time.Now()

			// Append death to death array, and increase player death count.
			x.Deaths = append(x.Deaths, newDeath)
			x.DeathCount += 1
		}
	}

	// return currentLobby to main function.
	return currentLobby
}

/*
*	Function called to handle player deaths against environment.
*	Environment counts as AI, or controlled flight into terrain.
* 	TODO: tidy up function and replace spaghetti code with RegEx.
*	TODO: capture string examples by using breakpoints.
 */
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

/*
*	Track identity updates by matching examples with RegEx, and update player's current aircraft.
*	TODO: Test function in both single crew and multicrew aircraft.
*	TODO: Obtain proper multicrew string example.
 */
func onIdentityUpdate(currentMessage string, currentLobby LobbyStruct) LobbyStruct {

	// Look for matching aircraft in current message string, value returns as string and...
	newAircraft, found := matchAircraft(currentMessage)
	// If not found, return current lobby to main function.
	if !found {
		return currentLobby
	}
	// Look for pilot(s) username in current message, value returns as []string and...
	newCrew, found := matchUsername(currentMessage)
	// If none are found, return currentLobby to main function.
	if !found {
		return currentLobby
	}

	// Range over newCrew []string and...
	for _, x := range newCrew {
		// Range over currentLobby.Players []LobbyPlayerStruct and...
		for _, y := range currentLobby.Players {
			// If current currentLobby.Players.Name index matches current newCrew.Name index,
			// and current currentLobby.Players.Active equals true, then...
			if y.Name == x && y.Active {
				// current currentLobby.Players.Aircraft index equals newAircraft
				y.Aircraft = newAircraft
			}
		}
	}

	// return currentLobby to main function
	return currentLobby
}

/*
*	matchUsername matches a pilot(s) usernames using RegEx from a given string of type 'identity updated'.
*   If any amount of matches are found, then return matchesFound are returned, and found is set to true.
*	If no matches are found then matchesFound is set to nil, and found is set to false.
 */
func matchUsername(currentMessage string) (matchesFound []string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(crewRegex, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match, err := re.FindStringMatch(currentMessage)
	if err != nil {
		return nil, false
	} else if match == nil {
		return nil, false
	}
	var matches []string
	for _, x := range match.Captures {
		matches = append(matches, x.String())
	}

	// Get the rightmost match (last element in the slice)
	if matches != nil {
		rightmostMatch := matches[len(matches)-1]

		// Check if multiple pilots are in the match
		crew1, crew2, found := strings.Cut(rightmostMatch, ", ")
		// If only one pilot is found, then...
		if !found {

			// Return single pilot as []string
			return []string{rightmostMatch}, true
		} else
		// Else if multiple pilots found,
		{

			// Return pilots as []string
			return []string{crew1, crew2}, true
		}
	} else
	// If no matches found, then...
	{

		// return []string as nil, and return bool as false to signal no match found.
		return nil, false
	}
}

/*
*	matchAircraft matches an aircraft name using RegEx from a given string of type 'identity updated'.
*	If a match is found, then return matchFound is returned, and found is set to true.
*	If no match is found then matchFound is set to "", and found is set to false.
 */
func matchAircraft(currentMessage string) (matchFound string, found bool) {
	// Compile the regular expression.
	re := regexp.MustCompile(craftRegex)

	// Find all matches in the string.
	matches := re.FindAllString(currentMessage, -1)

	// Get the rightmost match (last element in the slice), and...
	if len(matches) > 0 {
		rightmostMatch := matches[len(matches)-1]

		// Return match, and set found to true.

		return rightmostMatch, true
	} else
	// If no matches are found, then...
	{
		// Return "", and set found to false.

		return "", false
	}
}

/*
*	matchID64 matches all 17 digit numbers in a string,
*	which is the exact length of SteamID64 numbers, using RegEx.
 */
func matchID64(currentMessage string) (matchesFound []string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(id64RegEx, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match, err := re.FindStringMatch(currentMessage)
	if err != nil {
		return nil, false
	}
	var matches []string

	// If no matches found, then...
	if match == nil {
		// Return nil slice, and set found to false
		return nil, false
	}
	for _, x := range match.Captures {
		matches = append(matches, x.String())
	}

	return matches, true
}

/*
*	onSlotUISetup is called whenever a slot is updated.
*	Checks if a player changed teams, and calls on a function that
*	creates new player with said team on array if not exists.
 */
func onSlotUISetup(currentMessage string, currentLobby LobbyStruct) LobbyStruct {

	player, found := matchUsername(currentMessage)
	if !found {
		return currentLobby
	}

	var team string

	if strings.Contains(currentMessage, "Allied") {
		team = "Allied"

	} else if strings.Contains(currentMessage, "Enemy") {
		team = "Enemy"
	}

	for _, x := range player {
		for _, y := range currentLobby.Players {
			if y.Name == x && y.Active && y.Team != team {
				currentLobby = switchPlayerTeam(y, currentLobby)
				break
			}
		}
	}

	return currentLobby
}

/*
*	switchPlayerTeam is called upon when a function need a player's team swapped,
*	and creates new player with said team on array if not exists.
 */
func switchPlayerTeam(player LobbyPlayerStruct, currentLobby LobbyStruct) LobbyStruct {
	var found bool

	// Set old player as inactive.
	for _, x := range currentLobby.Players {
		if x.ID64 == player.ID64 && x.Active && x.Team == player.Team {
			x.Active = false
		}
	}

	//	Checks whether inactive player with certain ID64 exists, and...
	for _, x := range currentLobby.Players {
		if x.ID64 == player.ID64 && !x.Active && x.Team != player.Team {
			// If exists, set found to true, then..
			found = true
		}
	}

	// If opposite team player found, do, otherwise do...
	if found {
		for _, x := range currentLobby.Players {
			if x.ID64 == player.ID64 && x.Active && x.Team == player.Team {
				// Set old player as inactive
				x.Active = false
			} else if x.ID64 == player.ID64 && !x.Active && x.Team != player.Team {
				// Set new player as active
				x.Active = true
			}
		}
	} else {
		// Create new active player
		newPlayer := LobbyPlayerStruct{
			Name:     player.Name,
			ID64:     player.ID64,
			JoinedAt: time.Now(),
			Active:   true,
		}
		// Set proper team
		if player.Team == "Allied" {
			newPlayer.Team = "Enemy"
		} else {
			newPlayer.Team = "Allied"
		}
		// Append to currentLobby array and...
		currentLobby.Players = append(currentLobby.Players, newPlayer)
	}

	// Return updated lobby to calling function.
	return currentLobby
}

/*
* 	onPlayerLeave is called everytime a player leaves, and sets the leave time to time.Now()
*	TODO: Replace spaghetti code with RegEx
*	TODO: Obtain string sample
 */
func onPlayerLeave(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	trimmedMessage, _ := strings.CutPrefix(currentMessage, "FlightLogger: ")
	name, _ := strings.CutSuffix(trimmedMessage, " has disconnected.")
	for _, x := range currentLobby.Players {
		if x.Name == name && x.LeftAt.IsZero() {
			x.LeftAt = time.Now()
		}
	}
	return currentLobby
}

/*
*	onNewPlayer is called everytime a player joins, then calls unwrapInfo to figure out if player is new,
*	or player rejoined lobby, and act accordingly
*	TODO: Replace spaghetti code with RegEx.
*	TODO: Obtain string sample
 */
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

/*
* onBeginObjective is called everytime a new objective is beginned, and adds onto the lobby objective counter.
 */
func onBeginObjective(_ string, currentLobby LobbyStruct) LobbyStruct {
	// appends new objective with id equal to the length of the amount of current lobby objectives.
	currentLobby.Lobby.Objectives = append(currentLobby.Lobby.Objectives, ObjectiveStruct{
		ID: len(currentLobby.Lobby.Objectives),
	})
	return currentLobby
}

/*
*	onCompleteObjective is called everytime an objective is completed, and adds onto the lobby objective counter.
*	TODO: replace logic with RegEx.
*	TODO: Obtain string sample
 */
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

/*
*	preLobbyHandler handles the prelobby joining mechanism, and fills in information about the lobby,
*	into the array that will be needed later once joined and in game
 */
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

/*
*	unwrapInfo extracts player info from strings of type " - info"
*	TODO: Replace spaghetti code with RegEx.
*	TODO: Obtain string sample
 */
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

/*
*	onKill is called everytime a player kills something or someone,
*	and unwraps the information to save into an array.
*	TODO: Replace spaghetti code with RegEx.
*	TODO: Obtain string sample
 */
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
	for _, x := range currentLobby.Players {
		if x.Name == killer && x.Active {
			x.KillCount += 1
			aircraft = x.Aircraft
			newKill.UserTeam = x.Team
			killerid = x.ID64
			newKill.Copilot = x.Copilot
			x.Kills = append(x.Kills, newKill)
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
		for _, x := range currentLobby.Players {
			for _, h := range killedName {
				if x.Name == h && x.Active {
					x.DeathCount += 1
					x.Deaths = append(x.Deaths, newDeath)

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
