package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"tournament_pair/src/gotha"
)

var (
	round  = flag.Int("round", 1, "The round number")
	system = flag.String("system", "SWISS", "the pair system")
)

func main() {
	flag.Parse()

	// step1 init
	g := gotha.NewGotha()
	// Demo 1
	filePath := `../demo.xml` //  your file path
	if err := g.GetFromXMLFile(filePath); err != nil {
		log.Fatal(err)
	}

	// Demo2

	// step1 chose the pair system
	// Currently only supports Swiss-made arrangements
	g.SelectSystem(*system)
	t := g.GetTournament()

	// step2 choose the player （via keyString）
	t.SetSelectedPlayers([]string{"ARNAUDANCELIN",
		"AVENELAUGUSTIN",
		"BILLOIRECLEMENT",
		"BILLOUETSIMON",
		"BLANCHARDBENJAMIN",
		"BUFFARDEMMANUEL",
		"CANCEPHILIPPE",
		"COQUELETLAURENT",
		"CORNUEJOLSDOMINIQUE",
		"CRUBELLIERETIENNE",
		"DOUSSOTPATRICE",
		"FEVRIERLOUIS",
		"GAUTHIERHENRI",
		"GRANGERALBAN",
		"GUENNOUMORAN",
		"GUEVELBRENDAN",
		"HENRYYANNIS",
		"HWANGIN-SEONG",
		"IMAMURA-CORNUEJOLSTORU",
		"KARADABANDENIS",
		"KUNNESTÉPHAN",
		"LE_BROUSTERREMI",
		"LE_CALVÉTANGUY",
		"LEESEMI",
		"LIHAOHAN",
		"MASSONFABIEN",
		"MIZESSYNFRANÇOIS",
		"MOSCATELLIALDO",
		"NADDEFJEAN_LOUP",
		"NESMEVINCENT",
		"NGUYENHUU_PHUOC",
		"PAPAZOGLOUBENJAMIN",
		"PARCOITDAVID",
		"PUYAUBREAUNICOLAS",
		"ROBERTLUDWIG",
		"TECCHIOPIERRE",
		"TURLOTQUENTIN",
		//"VANNIERRÉMI",
		"WUBEILUN",
		"WURZINGERRALF"})

	// step3 pair
	t.Pair(*round)
	for _, game := range t.SortGameByTableNumber() {
		fmt.Printf("white : %s  <> black : %s\n", game.GetWhitePlayer().Name+" "+game.GetWhitePlayer().FirstName, game.GetBlackPlayer().Name+" "+game.GetBlackPlayer().FirstName)
	}

	//  step4 you will get a io.reader
	rd, err := g.IO.FlushGameToXML(t.SortGameByTableNumber())
	if err != nil {
		log.Fatal(err)
	}

	// step5 overwrite your xml file (Optional operation)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	_, err = io.Copy(file, rd)
	if err != nil {
		log.Fatal(err)
	}
}
