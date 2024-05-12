package main

import (
	"fmt"
	"strings"
	"time"
)

func onIlstUpdate(currentMessage string, currentLobby LobbyStruct, host bool) LobbyStruct {

	ilst, found := ilstUnwraps(currentMessage)

	if !found {
		return currentLobby
	}

	for _, x := range ilst {

		matches, found := separateIlst(x)
		if !found {
			continue
		}

		if strings.Contains(matches[1], "[") {
			matches[1], _, _ = strings.Cut(matches[1], " [")
		}
		if host {

			for x, y := range currentLobby.Players {

				if y.Name == matches[1] && y.Active && y.ID64 == "" {
					y.ID64 = matches[0]
					currentLobby.Players[x] = y
				}

			}

		}

		var foundMatch bool

		foundMatch, currentLobby = checkIfUserExists(matches, currentLobby)

		if foundMatch {
			continue
		}
		newPlayer := createPlayer(matches[1], matches[0], matches[2])

		currentLobby.Players = append(currentLobby.Players, newPlayer)
	}

	return currentLobby
}

/*
* onBeginObjective is called everytime a new objective is beginned, and adds onto the lobby objective counter.
 */
func onBeginObjective(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	// appends new objective with id equal to the length of the amount of current lobby objectives.
	objectiveName := strings.TrimPrefix(currentMessage, "Setting up objective ")
	newObjective := ObjectiveStruct{
		Name:    objectiveName,
		BeganAt: time.Now().Unix(),
	}
	currentLobby.Lobby.Objectives = append(currentLobby.Lobby.Objectives, newObjective)
	return currentLobby
}

func onResetObjective(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	// appends new objective with id equal to the length of the amount of current lobby objectives.
	objectiveName := strings.TrimPrefix(currentMessage, "Resetting objective ")
	newObjective := ObjectiveStruct{
		Name:    objectiveName,
		BeganAt: time.Now().Unix(),
	}
	currentLobby.Lobby.Objectives = append(currentLobby.Lobby.Objectives, newObjective)
	return currentLobby
}

func onCompleteObjective(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	objectiveName := strings.TrimPrefix(currentMessage, "FlightLogger: Completed objective: ")
	for x, y := range currentLobby.Lobby.Objectives {
		if y.Name == objectiveName && y.Result == "" {
			y.Result = "Completed"
			y.ResultedAt = time.Now().Unix()
			currentLobby.Lobby.Objectives[x] = y
		}
	}
	return currentLobby
}

func onFailObjective(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	objectiveName := strings.TrimPrefix(currentMessage, "FlightLogger: Failed objective: ")
	for x, y := range currentLobby.Lobby.Objectives {
		if y.Name == objectiveName && y.Result == "" {
			y.Result = "Failed"
			y.ResultedAt = time.Now().Unix()
			currentLobby.Lobby.Objectives[x] = y
		}
	}
	return currentLobby
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
		killedName = append(killedName, killed)
	}

	if strings.Contains(killed, "[") {
		killed, _, _ = strings.Cut(killed, " [")
	}
	for x, y := range killedName {
		if strings.Contains(y, "[") {
			killedName[x], _, _ = strings.Cut(y, " [")
		}
	}
	if strings.Contains(killer, "[") {
		killed, _, _ = strings.Cut(killed, " [")
	}

	weapon, _ = strings.CutSuffix(trimmedMessage, ".")
	newKill := KillStruct{
		Weapon: weapon,
		Time:   time.Now().Unix(),
		Killed: killed,
	}

	for _, y := range killedName {
		for _, u := range currentLobby.Players {
			if u.Name == y && u.Active {
				newKill.KilledID += "(" + u.ID64 + ")"
				newKill.KilledName += "(" + u.Name + ")"
				newKill.PlayerTeam = u.Team
				newKill.Killed = u.Aircraft
				break
			}
		}
	}

	PlayerKilling := false
	AlivePlayer := LobbyPlayerStruct{}
	var aircraft string
	var killerid string
	for x, y := range currentLobby.Players {
		if y.Name == killer && y.Active {
			PlayerKilling = true
			y.KillCount += 1
			aircraft = y.Aircraft
			newKill.UserTeam = y.Team
			killerid = y.ID64
			newKill.KilledBy = y.Aircraft
			newKill.UserTeam = y.Team
			if newKill.PlayerTeam == "" && newKill.KilledID == "" && newKill.KilledName == "" {
				newKill.KilledID = "<environment>"
				newKill.PlayerTeam = "<environment>"
				newKill.KilledName = "<environment>"
			}
			newKill.Copilot = y.Copilot
			y.Kills = append(y.Kills, newKill)
			currentLobby.Players[x] = y
			AlivePlayer = y
			break
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
		for x, y := range killedName {
			if strings.Contains(y, "[") {
				killedName[x], _, _ = strings.Cut(y, " [")
			}
		}

		PlayerKilled := false
		DeadPlayer := LobbyPlayerStruct{}

		for x, y := range currentLobby.Players {
			for _, h := range killedName {
				if y.Name == h && y.Active {
					y.DeathCount += 1
					newDeath.DiedWith = y.Aircraft
					y.Deaths = append(y.Deaths, newDeath)
					currentLobby.Players[x] = y
					PlayerKilled = true
					DeadPlayer = y
					break
				}
			}

		}

		if PlayerKilled && PlayerKilling {
			if AlivePlayer.Kills[len(AlivePlayer.Kills)-1].Weapon == "GUN" {
				fmt.Println("HUMILIATED")
			}
		}

	}

	return currentLobby
}

func onMCSeatOccupy(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	username, aircraft, _ := strings.Cut(currentMessage, " has entered a multicrew seat in ")
	if strings.Contains(username, "[") {
		username, _, _ = strings.Cut(username, " [")
	}
	craft := ""
	switch {
	case strings.Contains(aircraft, "EF-24"):
		craft = "EF-24G"
	case strings.Contains(aircraft, "T-55") || strings.Contains(aircraft, "Test55"):
		craft = "T-55"
	case strings.Contains(aircraft, "AH-94"):
		craft = "AH-94"
	}

	for x, y := range currentLobby.Players {
		if y.Name == username && y.Active && craft != "" {
			currentLobby.Players[x].Aircraft = craft
		}
	}

	return currentLobby
}

func onMissionEnd(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	if strings.Contains(currentMessage, "Allied") {
		currentLobby.Lobby.WinningTeam = "Allied"

	} else if strings.Contains(currentMessage, "Enemy") {
		currentLobby.Lobby.WinningTeam = "Enemy"
	}
	return currentLobby
}
func updateLobbyCount(currentLobby LobbyStruct) LobbyStruct {
	currentLobby.Lobby.TotalLobbyDeaths = 0
	currentLobby.Lobby.TotalLobbyKills = 0
	for _, y := range currentLobby.Players {
		currentLobby.Lobby.TotalLobbyKills += len(y.Kills)
		currentLobby.Lobby.TotalLobbyDeaths += len(y.Deaths)
	}
	return currentLobby
}

// func onPlayerUpdate(currentLobby LobbyStruct) LobbyStruct {
//	return currentLobby
//}

/*
* 	onPlayerLeave is called everytime a player leaves, and sets the leave time to time.Now()
*	TODO: Replace spaghetti code with RegEx
*	TODO: Obtain string sample
 */
func onPlayerLeave(currentMessage string, currentLobby LobbyStruct) LobbyStruct {
	trimmedMessage, _ := strings.CutPrefix(currentMessage, "FlightLogger: ")
	name, _ := strings.CutSuffix(trimmedMessage, " has disconnected.")
	if strings.Contains(name, "[") {
		name, _, _ = strings.Cut(name, " [")
	}
	for x, y := range currentLobby.Players {
		if y.Name == name && y.LeftAt.IsZero() {
			y.LeftAt = time.Now()
			y.InGame = false
			currentLobby.Players[x] = y
		}
	}

	return currentLobby
}

func onSetTeam(currentMessage string, currentLobby LobbyStruct) LobbyStruct {

	leftside, rightside, found := strings.Cut(currentMessage, "(")
	if !found {
		return currentLobby
	}
	username, rightestside, found := strings.Cut(rightside, ")")
	if username == "Clone" {
		_, username, _ = strings.Cut(rightestside, "(")
		username, _, _ = strings.Cut(username, ")")
	}
	if strings.Contains(username, "T-55") || strings.Contains(username, "AH-94") || strings.Contains(username, "EF-24G") {
		leftside, username, _ = strings.Cut(username, "(")
	}
	var multicrew bool
	if strings.Contains(username, ",") {
		multicrew = true
	}

	if strings.Contains(username, "(") {
		_, username, _ = strings.Cut(username, "(")
		username, _, _ = strings.Cut(username, ")")
	}

	if !found {
		return currentLobby
	}

	var aircraft string

	switch {
	case strings.Contains(leftside, "A/V-42C") || strings.Contains(leftside, "vtol4") || strings.Contains(leftside, "AV-42C"):
		aircraft = "A/V-42C"
	case strings.Contains(leftside, "F/A-26B") || strings.Contains(leftside, "FA-26B") || strings.Contains(leftside, "afighter"):
		aircraft = "F/A-26B"
	case strings.Contains(leftside, "EF-24G") || strings.Contains(leftside, "EF-24"):
		aircraft = "EF-24G"
	case strings.Contains(leftside, "F-45A") || strings.Contains(leftside, "SEVTF"):
		aircraft = "F-45A"
	case strings.Contains(leftside, "T-55") || strings.Contains(currentMessage, "Test55"):
		aircraft = "T-55"
	case strings.Contains(leftside, "AH-94"):
		aircraft = "AH-94"
	}

	if strings.Contains(username, "[") && !multicrew {
		username, _, _ = strings.Cut(username, " [")
	}

	if multicrew {
		username1, username2, _ := strings.Cut(username, ", ")

		if strings.Contains(username1, "[") {
			username1, _, _ = strings.Cut(username1, " [")
		}
		if strings.Contains(currentMessage, "") {
			username2, _, _ = strings.Cut(username2, " [")
		}
		for x, y := range currentLobby.Players {
			if y.Name == username1 && y.Active && y.Aircraft != "" {
				currentLobby.Players[x].Aircraft = aircraft

			}
		}
		for x, y := range currentLobby.Players {
			if y.Name == username2 && y.Active && aircraft != "" {
				currentLobby.Players[x].Aircraft = aircraft

			}
		}
	}
	for x, y := range currentLobby.Players {
		if y.Name == username && y.Active && aircraft != "" && !multicrew {
			currentLobby.Players[x].Aircraft = aircraft

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

	for x, y := range newCrew {
		if strings.Contains(y, "[") {
			newCrew[x], _, _ = strings.Cut(y, " [")
		}
	}

	// Range over newCrew []string and...
	for _, y := range newCrew {
		// Range over currentLobby.Players []LobbyPlayerStruct and...
		for z, u := range currentLobby.Players {
			// If current currentLobby.Players.Name index matches current newCrew.Name index,
			// and current currentLobby.Players.Active equals true, then...
			if u.Name == y && u.Active {
				// current currentLobby.Players.Aircraft index equals newAircraft
				u.Aircraft = newAircraft

				currentLobby.Players[z] = u
			}
		}
	}

	// return currentLobby to main function
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
	username, found := matchUsername(currentMessage)
	if !found {
		return currentLobby
	}
	for x, y := range username {
		if strings.Contains(y, "[") {
			username[x], _, _ = strings.Cut(y, " [")
		}

	}

	unit, found := matchUnit(currentMessage)
	if !found {
		return currentLobby
	}

	if len(unit) == 1 {
		newDeath.Weapon = unit[0]
	}
	newDeath.KilledBy = "<environment>"

	for _, h := range username {

		for x, y := range currentLobby.Players {
			if y.Name == h && y.Active {
				newDeath.UserTeam = y.Team
				newDeath.DiedWith = y.Aircraft
				newDeath.PlayerTeam = "<environment>"
				newDeath.Time = time.Now().Unix()
				currentLobby.Players[x].Deaths = append(currentLobby.Players[x].Deaths, newDeath)
				currentLobby.Players[x].DeathCount += 1
			}
		}
	}

	return currentLobby

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

	if strings.Contains(tmpname1, "[") {
		tmpname1, _, _ = strings.Cut(tmpname1, " [")
	}
	if strings.Contains(tmpname2, "[") {
		tmpname2, _, _ = strings.Cut(tmpname2, " [")
	}

	// Fill in missing death information using currentLobby.Players information
	for x, y := range currentLobby.Players {
		// If x.Name matches either of tmpname1 or tmpname2, and x.Active is true, then.
		if (y.Name == tmpname1 || y.Name == tmpname2) && y.Active {

			//  Fill in missing information using x
			newDeath.UserTeam = y.Team
			newDeath.DiedWith = y.Aircraft

			// Set killer team as environment (i.e. AI killed user.)
			newDeath.PlayerTeam = "<environment>"

			// Set time of death as time.Now()
			newDeath.Time = time.Now().Unix()

			// Append death to death array, and increase player death count.
			y.Deaths = append(y.Deaths, newDeath)
			y.DeathCount += 1
			currentLobby.Players[x] = y
			break
		}
	}

	// return currentLobby to main function.
	return currentLobby
}

/*
*	onSlotUISetup is called whenever a slot is updated. Unverified works as host mode.
*	Checks if a player changed teams, and calls on a function that
*	creates new player with said team on array if not exists.
 */
func onSlotUISetup(currentMessage string, currentLobby LobbyStruct) LobbyStruct {

	player, found := matchUsername(currentMessage)
	if !found {
		return currentLobby
	}

	for x, y := range player {
		if strings.Contains(y, "[") {
			player[x], _, _ = strings.Cut(y, " [")
		}
	}
	var team string

	if strings.Contains(currentMessage, "Allied") {
		team = "Allied"

	} else if strings.Contains(currentMessage, "Enemy") {
		team = "Enemy"
	}

	for _, y := range player {
		for _, u := range currentLobby.Players {
			if u.Name == y && u.Active && u.Team != team {
				currentLobby = switchPlayerTeam(u, currentLobby)
				break
			}
		}
	}

	return currentLobby
}
