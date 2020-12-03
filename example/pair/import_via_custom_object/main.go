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

	for i := 0; i < 12; i++ {
		player := gotha.NewPlayer()
		player.SetName(getRandName(8))
		player.SetFirstName(getRandName(5) + "")

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
	//game := gotha.NewGame()
	//game.SetRoundNumber(0)
	//game.SetResult(gotha.SelectResult("RESULT_UNKNOWN"))
	//game.SetTableNumber(1)
	//game.SetHandicap(0)
	//game.SetKnownColor(true)
	//bPlayerKeyString := "KARADABANDENIS"
	//wPlayerKeyString := "WUBEILUN"
	//game.SetBlackPlayer(tournament.GetPlayerByKeyString(bPlayerKeyString))
	//game.SetWhitePlayer(tournament.GetPlayerByKeyString(wPlayerKeyString))
	//tournament.AddGame(game)

	gps := tournament.GetTournamentSet().GetGeneralParameterSet()
	gps.SetNumberOfRounds(10)
	g.SetTournament(tournament)
	fmt.Println(g.GetTournament().GetTournamentSet().GetGeneralParameterSet().GetNumberOfRounds())
	t := g.GetTournament()
	g.SelectSystem(*system)

	// Step2 choose the players （via keyString）
	// By default, all players participate in this round
	t.SetSelectedPlayers()
	// Step3 pair
	t.Pair(*round)

	for _, game := range t.SortGameByTableNumber() {
		fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
	}
	t.SortGameByTableNumberFromRn(1)
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
