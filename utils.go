package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

func ilstUnwraps(currentMessage string) (matchesFound []string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(ilstRegex, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match := regexp2FindAllString(re, currentMessage)

	if match == nil {
		return
	}
	return match, true

}

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func separateIlst(ilst string) (matchesFound []string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(ilstSeparateRegex, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match := regexp2FindAllString(re, ilst)

	if match == nil {
		return
	}
	return match, true
}

func checkIfUserExists(matches []string, currentLobby LobbyStruct) (userExists bool, lobby LobbyStruct) {

	id := matches[0]
	username := matches[1]

	var foundUser bool
	for x, y := range currentLobby.Players {
		if y.Name == username && y.ID64 == id && !y.InGame {
			foundUser = true

			y.InGame = true
			currentLobby.Players[x] = y
		} else if y.Name == username && y.ID64 == id && y.InGame {

			foundUser = true
		}
	}

	if foundUser {
		return true, currentLobby
	}

	return false, currentLobby

}

func createPlayer(username string, id64 string) LobbyPlayerStruct {

	newPlayer := LobbyPlayerStruct{
		Name:     username,
		JoinedAt: time.Now(),
		Active:   true,
		InGame:   true,
		Team:     "Allied",
		ID64:     id64,
	}

	return newPlayer
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

func matchUnit(currentMessage string) (matchesFound []string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(UnitRegex, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match := regexp2FindAllString(re, currentMessage)
	if match == nil {
		return
	}
	var matches []string
	for _, x := range match {
		if x == "T-55" || x == "F-45A" || x == "AH-94" || x == "AV-42C" || x == "EF-24G" || x == "F/A-26B" {
			return
		}
		matches = append(matches, x)
	}
	if matches == nil {
		return
	}
	return matches, true

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
	matches := regexp2FindAllString(re, currentMessage)
	if matches == nil {
		return nil, false
	}

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

func matchHostedMap(currentMessage string) (matchesFound string, found bool) {

	// Compile the regular expression
	re := regexp2.MustCompile(lobbyHostRegex, regexp2.DefaultUnmarshalOptions)

	// Find all matches in the string
	match, err := re.FindStringMatch(currentMessage)
	if err != nil {
		return
	}
	var matches []string

	// If no matches found, then...
	if match == nil {
		// Return nil slice, and set found to false
		return
	}
	for _, x := range match.Captures {
		matches = append(matches, x.String())
	}

	return matches[0], true
}
