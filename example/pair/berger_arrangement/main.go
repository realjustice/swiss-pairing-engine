package main

import (
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	// init tournament
	tournament := gotha.NewTournament()

	for i := 1; i <= 7; i++ {
		player := gotha.NewPlayer()
		player.SetName(strconv.Itoa(i))
		player.SetRank(getRandScore())
		tournament.AddPlayer(player)
	}

	tournament.BergerArrange().Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("black : %s  <> white : %s tableNumber:%d roundNumber:%d \n", game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName, game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.TableNumber, game.RoundNumber)

		return false
	})

	for i := 1; i <= 7; i++ {
		fmt.Printf("本轮轮空人员：%s\n", tournament.GetByePlayer(i).Name)
	}

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
