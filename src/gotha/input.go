package gotha

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

type Input struct {
	bPlayers bool
	bGames   bool
}

type OptionFunc func(input *Input)

type IOption interface {
	apply(*Input)
}

type TourXML struct {
	Tournament xml.Name  `xml:"Tournament"`
	Version    string    `xml:"dataVersion,attr"`
	Player     XMLPlayer `xml:"Players"`
}

type XMLPlayer struct {
	xml.Name `xml:"Players"`
	Players  []*Player `xml:"Player"`
}

func (f OptionFunc) apply(input *Input) {
	f(input)
}

func WithGames() IOption {
	return OptionFunc(func(input *Input) {
		input.bGames = true
	})
}

func WithPlayers() IOption {
	return OptionFunc(func(input *Input) {
		input.bPlayers = true
	})
}

func NewInput() *Input {
	return new(Input)
}

func (i *Input) WithOption(options ...IOption) {
	for _, option := range options {
		option.apply(i)
	}
}

/***
  导入xml文件
*/
func (i *Input) ImportTournamentFromXMLFile(filePath string, tournament *Tournament) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	tourXml := new(TourXML)
	err = xml.Unmarshal(data, tourXml)
	if err != nil {
		log.Fatal(err)
	}

	// 导入所有比赛选手
	if i.bPlayers {
		players := importPlayersFromXMLFile(tourXml)
		tournament.AddPlayer(players)
	}

	// 导入所有比赛对局
	if i.bGames {
		i.importGamesFromXMLFile()
	}
}

func importPlayersFromXMLFile(tourXML *TourXML) []*Player {
	participating := make([]bool, MAX_NUMBER_OF_ROUNDS)
	players := tourXML.Player.Players
	for _, p := range players {
		for i := 0; i < len(p.ParticipatingStr); i++ {
			if string(p.ParticipatingStr[i]) == "0" {
				participating[i] = false
			} else {
				participating[i] = true
			}
		}
		// 设置本轮是否参与编排
		p.SetParticipating(participating)
		p.SetRank(p.Rating)
	}

	return players
}

func (i *Input) importGamesFromXMLFile() []*Game {
	return nil
}
