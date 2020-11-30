package gotha

import (
	"bytes"
	"errors"
	"github.com/beevik/etree"
	"io"
	"os"
	"strconv"
)

type Input struct {
	Root     *etree.Element
	bPlayers bool
	bGames   bool
}

type OptionFunc func(input *Input)

type IOption interface {
	apply(*Input)
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
func (i *Input) ImportTournamentFromXMLFile(filePath string, tournament *Tournament) error {
	file, err := os.Open(filePath)
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}
	return i.ImportFromReader(file, tournament)
}

/***
  导入bytes[]
*/
func (i *Input) ImportTournamentFromBytes(bs []byte, tournament *Tournament) error {
	return i.ImportFromReader(bytes.NewReader(bs), tournament)
}

func (i *Input) ImportFromReader(ri io.Reader, t *Tournament) error {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(ri)
	if err != nil {
		return err
	}
	root := doc.SelectElement("Tournament")
	if root == nil {
		return errors.New("xml 导入选手失败！")
	}
	i.Root = root
	// 导入所有比赛选手
	if i.bPlayers {
		players, err := i.importPlayersFromXML()
		if err != nil {
			return err
		}
		for _, p := range players {
			t.AddPlayer(p)
		}
	}

	// 导入所有比赛对局
	if i.bGames {
		games, err := i.importGamesFromXML(t)
		if err != nil {
			return err
		}
		if len(games) <= 0 || games == nil {
			return nil
		}

		for _, g := range games {
			t.AddGame(g)
		}
	}
	return nil
}

func (i *Input) importPlayersFromXML() (players []*Player, err error) {
	participating := make([]bool, MAX_NUMBER_OF_ROUNDS)
	players = make([]*Player, 0)
	playersXML := i.Root.SelectElement("Players")
	if playersXML == nil {
		err = errors.New("player 导入失败！！")
		return players, err
	}
	for _, playerXML := range playersXML.SelectElements("Player") {
		p := NewPlayer()
		p.Name = playerXML.SelectAttrValue("name", "")
		p.FirstName = playerXML.SelectAttrValue("firstName", "")
		ratingStr := playerXML.SelectAttrValue("rating", "")
		rating, err := strconv.Atoi(ratingStr)
		if err != nil {
			return players, err
		}
		p.Rating = rating
		p.ParticipatingStr = playerXML.SelectAttrValue("participating", "")
		p.SetRank(p.Rating)
		for i := 0; i < len(p.ParticipatingStr); i++ {
			if string(p.ParticipatingStr[i]) == "0" {
				participating[i] = false
			} else {
				participating[i] = true
			}
		}
		// 设置本轮是否参与编排
		p.SetParticipating(participating)
		players = append(players, p)
	}

	return players, err
}

func (i *Input) importGamesFromXML(tournament *Tournament) (games []*Game, err error) {
	games = make([]*Game, 0)
	gamesXML := i.Root.SelectElement("Games")
	if gamesXML == nil {
		return games, err
	}
	for _, gameXML := range gamesXML.SelectElements("Game") {
		g := NewGame()
		roundNumberXML := gameXML.SelectAttrValue("roundNumber", "0")
		roundNumber, err := strconv.Atoi(roundNumberXML)
		if err != nil {
			return games, err
		}
		tableNumberXML := gameXML.SelectAttrValue("tableNumber", "0")
		tableNumber, err := strconv.Atoi(tableNumberXML)
		if err != nil {
			return games, err
		}
		g.SetRoundNumber(roundNumber - 1)
		g.SetBlackPlayer(tournament.GetPlayerByKeyString(gameXML.SelectAttrValue("blackPlayer", "")))
		g.SetWhitePlayer(tournament.GetPlayerByKeyString(gameXML.SelectAttrValue("whitePlayer", "")))
		g.SetKnownColor(true)
		g.SetResult(SelectResult(gameXML.SelectAttrValue("result", "")))
		g.SetTableNumber(tableNumber - 1)
		games = append(games, g)
	}
	return games, err
}
