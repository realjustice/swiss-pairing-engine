package main

import (
	"flag"
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
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
	// add player black
	player1 := gotha.NewPlayer()
	player1.SetName("Karadaban")
	player1.SetFirstName("Denis")

	player1.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排
	player1.SetRank(50)
	tournament.AddPlayer(player1)

	// white
	player2 := gotha.NewPlayer()
	player2.SetName("Wu")
	player2.SetFirstName("Beilun")
	player2.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排

	player2.SetRank(2556)
	tournament.AddPlayer(player2)
	// bye player
	byePlayer := gotha.NewPlayer()
	byePlayer.SetName("Granger")
	byePlayer.SetFirstName("Alban")
	byePlayer.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排
	byePlayer.SetRank(100)
	tournament.AddPlayer(byePlayer)

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

	t.GetByePlayer(1)
	t.SetByePlayer(1, "GRANGERALBAN")

	for _, game := range t.SortGameByTableNumber() {
		fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
	}
}
