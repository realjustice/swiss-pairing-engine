package main

import (
	"flag"
	"fmt"
	"github.com/realjustice/swiss-pairing-engine/src/gotha"
	"io/ioutil"
	"log"
	"os"
)

var (
	round  = flag.Int("round", 1, "The round number")
	system = flag.String("system", "SWISS", "the pair system")
)

func main() {
	flag.Parse()

	// Step1 init the pair engine
	g := gotha.NewGotha()

	// Step2 import the data resource
	// you can import from the xml file
	//pwd, _ := os.Getwd()
	filePath := `/Users/justice/Desktop/tour_43_all.xml` //  your file path
	importFromXMLFile(filePath, g)

	// or from bytes
	// importFromBytes(filePath, g)

	// Step3 chose the pair system
	// Currently only supports Swiss-made arrangements

	g.SelectSystem(*system)
	t := g.GetTournament()
	t.GetTournamentSet().GetGeneralParameterSet().SetNumberOfRounds(10)

	// Step4 choose the players （via keyString）
	// By default, all players participate in this round
	//t.SetSelectedPlayers([]string{"KARADABANDENIS", "WUBEILUN"})

	// Step5 pair
	newGames := make([]*gotha.Game, 0)
	//t.SetSelectedPlayers([""])
	t.Pair(*round).Walk(func(game *gotha.Game) (isStop bool) {
		fmt.Printf("white : %s  <> black : %s  tableNumber:%d \n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName, game.GetTableNumber())
		newGames = append(newGames, game)
		return false
	})

	//  Step6 you will get a io.reader
	//rd, err := g.IO.FlushGameToXML(newGames)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// overwrite your xml file (Optional operation)
	//file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func() { _ = file.Close() }()
	//_, err = io.Copy(file, rd)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func importFromXMLFile(filePath string, g *gotha.Gotha) {
	if err := g.ImportFromXMLFile(filePath); err != nil {
		log.Fatal(err)
	}
}

func importFromBytes(filePath string, g *gotha.Gotha) {
	tempXML, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = tempXML.Close() }()
	fd, err := ioutil.ReadAll(tempXML)
	if err != nil {
		log.Fatal(err)
	}
	if err = g.ImportFromBytes(fd); err != nil {
		log.Fatal(err)
	}
}
