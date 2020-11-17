package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"tournament_pair/src/gotha"
)

type ToutXML struct {
	Tournament xml.Name  `xml:"Tournament"`
	Version    string    `xml:"dataVersion,attr"`
	Player     XMLPlayer `xml:"Players"`
}

type XMLPlayer struct {
	xml.Name `xml:"Players"`
	Players  []gotha.Player `xml:"Player"`
}



var hmScoredPlayers map[string]gotha.ScoredPlayer

func main() {
	file, err := os.Open(`/Users/justice/Desktop/clm2014B.xml`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	tourXml := new(ToutXML)
	err = xml.Unmarshal(data, tourXml)
	if err != nil {
		log.Fatal(err)
	}
	players := tourXml.Player.Players
	t := gotha.NewTournament(tourXml.Player.Players)
	t.MakeAutomaticPairing(players, 0)
}
