package gotha

import (
	"bytes"
	"errors"
	"github.com/beevik/etree"
	"io"
	"strconv"
	"strings"
)

type InputOutput struct {
	Doc        *etree.Document
	Root       *etree.Element
	bPlayers   bool
	bGames     bool
	bByePlayer bool
}

type OptionFunc func(input *InputOutput)

type IOption interface {
	apply(*InputOutput)
}

func (f OptionFunc) apply(io *InputOutput) {
	f(io)
}

func WithGames() IOption {
	return OptionFunc(func(io *InputOutput) {
		io.bGames = true
	})
}

func WithPlayers() IOption {
	return OptionFunc(func(io *InputOutput) {
		io.bPlayers = true
	})
}

func WithByePlayers() IOption {
	return OptionFunc(func(io *InputOutput) {
		io.bByePlayer = true
	})
}

func NewInputOutput() *InputOutput {
	return new(InputOutput)
}

func (i *InputOutput) WithOption(options ...IOption) {
	for _, option := range options {
		option.apply(i)
	}
}

func (i *InputOutput) ImportFromReader(ri io.Reader, t *Tournament) error {
	doc := etree.NewDocument()
	i.Doc = doc
	_, err := doc.ReadFrom(ri)
	if err != nil {
		return err
	}
	root := doc.SelectElement("Tournament")
	if root == nil {
		return errors.New("xml 导入选手失败！")
	}
	i.Root = root
	// 导入比赛初始配置
	i.importGeneralParameterSetFromXML(t)
	// 导入所有比赛选手
	if i.bPlayers {
		players, err := i.importPlayersFromXML()
		if err != nil {
			return err
		}
		//t.selectedPlayers = make([]*Player, len(players))
		for _, p := range players {
			t.AddPlayer(p)
		}
	}
	if i.bByePlayer {
		i.importByePlayersFromXML(t)
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

func (i *InputOutput) importPlayersFromXML() (players []*Player, err error) {
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
		// 如果XML中不导入，则根据firstName+name 生成
		p.SetKeyString(playerXML.SelectAttrValue("keyString", strings.ToUpper(p.Name+p.FirstName)))
		rating, err := strconv.Atoi(ratingStr)
		if err != nil {
			return players, err
		}
		p.Rating = rating
		p.SetParticipatingStr(playerXML.SelectAttrValue("participating", ""))
		p.SetRank(p.Rating)

		players = append(players, p)
	}

	return players, err
}

func (i *InputOutput) importGamesFromXML(tournament *Tournament) (games []*Game, err error) {
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
		g.setRoundNumber(roundNumber - 1)
		g.SetBlackPlayer(tournament.GetPlayerByKeyString(strings.ToUpper(gameXML.SelectAttrValue("blackPlayer", ""))))
		g.SetWhitePlayer(tournament.GetPlayerByKeyString(strings.ToUpper(gameXML.SelectAttrValue("whitePlayer", ""))))
		g.SetKnownColor(true)
		g.SetResult(SelectResult(gameXML.SelectAttrValue("result", "")))
		g.SetTableNumber(tableNumber - 1)
		games = append(games, g)
	}
	return games, err
}

func (i *InputOutput) importByePlayersFromXML(tournament *Tournament) {
	byePlayersXML := i.Root.SelectElement("ByePlayers")
	if byePlayersXML == nil {
		return
	}
	byePlayerXML := byePlayersXML.SelectElements("ByePlayer")
	if byePlayerXML == nil {
		return
	}
	for _, v := range byePlayerXML {
		keyString := v.SelectAttrValue("player", "")
		roundNumberStr := v.SelectAttrValue("roundNumber", "")
		roundNumber, _ := strconv.Atoi(roundNumberStr)
		tournament.SetByePlayer(roundNumber, keyString)
	}
}

func (i *InputOutput) importGeneralParameterSetFromXML(tournament *Tournament) {
	gps := tournament.tournamentParameterSet.GetGeneralParameterSet()
	tprXML := i.Root.SelectElement("TournamentParameterSet")
	if tprXML == nil {
		return
	}
	gpsXML := tprXML.SelectElement("GeneralParameterSet")
	if gpsXML == nil {
		return
	}
	rnStr := gpsXML.SelectAttrValue("numberOfRounds", "5")
	rn, _ := strconv.Atoi(rnStr)
	gps.SetNumberOfRounds(rn)
}

func (i *InputOutput) FlushGameToXML(games []*Game) (io.Reader, error) {
	gamesXML := i.Root.SelectElement("Games")
	if gamesXML == nil {
		gamesXML = i.Root.CreateElement("Games")
	}
	count := len(gamesXML.ChildElements())

	for index, g := range games {
		gameXML := gamesXML.CreateElement("Game")
		gameXML.CreateAttr("blackPlayer", g.blackPlayer.keyString)
		gameXML.CreateAttr("handicap", strconv.Itoa(g.GetHandicap()))
		gameXML.CreateAttr("blackPlayer", g.blackPlayer.keyString)
		gameXML.CreateAttr("result", ConvertResult(g.GetResult()))
		gameXML.CreateAttr("roundNumber", strconv.Itoa(g.GetRoundNumber()+1))
		gameXML.CreateAttr("tableNumber", strconv.Itoa(g.GetTableNumber()))
		gameXML.CreateAttr("whitePlayer", g.whitePlayer.keyString)
		gamesXML.InsertChildAt(index+count, gameXML)
	}
	i.Doc.Indent(2)

	bsw, err := i.Doc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bsw), nil
}

func createOrUpdateXMLAttr(element *etree.Element, key string, value string) {
	attr := element.SelectAttr(key)
	if attr == nil {
		element.CreateAttr(key, value)
	} else {
		attr.Value = value
	}
}
