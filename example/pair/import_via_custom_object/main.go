package main

import (
	"flag"
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"github.com/realjustice/swiss-pairing-engine/src/parameter_set"
	"strings"
)

var (
	round  = flag.Int("round", 1, "The round number")
	system = flag.String("system", "SWISS", "the pair system")
)

func main() {
	flag.Parse()
	g := gotha.NewGotha()
	// 如果不想通过xml导入的方式来初始化编排数据，可以通过
	tournament := gotha.NewTournament()
	// add player black
	player1 := gotha.NewPlayer()
	player1.SetName("Karadaban")
	player1.SetFirstName("Denis")
	player1.SetKeyString(strings.ToUpper(player1.Name + player1.FirstName))
	player1.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排
	player1.SetRank(2253)
	tournament.AddPlayer(player1)

	// white
	player2 := gotha.NewPlayer()
	player2.SetName("Wu")
	player2.SetFirstName("Beilun")
	player2.SetKeyString(strings.ToUpper(player2.Name + player2.FirstName))
	player2.SetParticipatingStr("111111111111111") // 是否参与每一轮的编排

	player2.SetRank(1961)
	tournament.AddPlayer(player2)

	// add game
	game := gotha.NewGame()
	game.SetRoundNumber(0)
	game.SetResult(gotha.SelectResult("RESULT_UNKNOWN"))
	game.SetTableNumber(1)
	game.SetHandicap(0)
	game.SetKnownColor(true)
	bPlayerKeyString := "KARADABANDENIS"
	wPlayerKeyString := "WUBEILUN"
	game.SetBlackPlayer(tournament.GetPlayerByKeyString(bPlayerKeyString))
	game.SetWhitePlayer(tournament.GetPlayerByKeyString(wPlayerKeyString))
	tournament.AddGame(game)

	gps := parameter_set.NewGeneralParameterSet()
	gps.SetNumberOfRounds(10)

	g.SetTournament(tournament)
	t := g.GetTournament()
	g.SelectSystem(*system)

	// Step4 choose the players （via keyString）
	// By default, all players participate in this round
	t.SetSelectedPlayers([]string{"KARADABANDENIS", "WUBEILUN"})
	// Step5 pair
	t.Pair(*round)
	for _, game := range t.SortGameByTableNumber() {
		fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
	}
}
