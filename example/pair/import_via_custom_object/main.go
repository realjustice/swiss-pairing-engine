package main

import (
	"flag"
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"math/rand"
	"strings"
	"time"
)

var (
	round  = flag.Int("round", 1, "The round number")
	system = flag.String("system", "SWISS", "the pair system")
)

func main() {

	flag.Parse()
	g := gotha.NewGotha()
	// Step1 init pair object
	tournament := gotha.NewTournament()

	blackKeyStrings, whiteKeyStrings := make([]string, 0), make([]string, 0)

	for i := 0; i < 12; i++ {
		player := gotha.NewPlayer()
		player.SetName(getRandName(8))
		player.SetFirstName(getRandName(5) + "")
		if i%2 == 0 {
			blackKeyStrings = append(blackKeyStrings, strings.ToUpper(player.Name+player.FirstName))
		} else {
			whiteKeyStrings = append(whiteKeyStrings, strings.ToUpper(player.Name+player.FirstName))
		}
		player.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排
		player.SetRank(getRandScore())
		tournament.AddPlayer(player)
	}

	// bye player
	//byePlayer := gotha.NewPlayer()
	//byePlayer.SetName("Granger")
	//byePlayer.SetFirstName("Alban")
	//byePlayer.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排
	//byePlayer.SetRank(100)
	//tournament.AddPlayer(byePlayer)

	// add game
	for i := 0; i < 6; i++ {
		game := gotha.NewGame()
		game.SetRoundNumber(0)
		game.SetResult(gotha.SelectResult("RESULT_WHITEWINS"))
		game.SetTableNumber(i)
		game.SetHandicap(0)
		game.SetKnownColor(true)
		bPlayerKeyString := blackKeyStrings[i]
		wPlayerKeyString := whiteKeyStrings[i]
		game.SetBlackPlayer(tournament.GetPlayerByKeyString(bPlayerKeyString))
		game.SetWhitePlayer(tournament.GetPlayerByKeyString(wPlayerKeyString))
		tournament.AddGame(game)
	}

	gps := tournament.GetTournamentSet().GetGeneralParameterSet()
	gps.SetNumberOfRounds(10)
	g.SetTournament(tournament)
	t := g.GetTournament()
	g.SelectSystem(*system)

	// Step2 choose the players （via keyString）
	// By default, all players participate in this round
	//t.SetSelectedPlayers()

	// Step3 pair and you will get a game iterator（current roundNumber）
	t.Pair(*round).Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
		fmt.Println(game.TableNumber)
		return false
	})
}

func getRandName(length int) string {
	rand.Seed(time.Now().UnixNano())
	str := make([]byte, length)
	for i := 0; i < length; i++ {
		char := rand.Intn(26)
		str[i] = "a"[0] + uint8(char)
	}

	return strings.ToUpper(string(str))
}

func getRandScore() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2050)
}
