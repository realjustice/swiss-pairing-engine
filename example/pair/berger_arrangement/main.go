package main

import (
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"math/rand"
	"strings"
	"time"
)

func main() {
	// init tournament
	tournament := gotha.NewTournament()

	blackKeyStrings, whiteKeyStrings := make([]string, 0), make([]string, 0)

	for i := 0; i < 51; i++ {
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

	tournament.BergerArrange().Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("white : %s  <> black : %s tableNumber:%d roundNumber:%d \n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName, game.TableNumber, game.RoundNumber)
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
