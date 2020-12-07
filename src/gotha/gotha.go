package gotha

import (
	"bytes"
	. "github.com/realjustice/swiss-pairing-engine/src/parameter_set"
	"io"
	"os"
)

const (
	MAX_NUMBER_OF_ROUNDS   = 20
	MAX_NUMBER_OF_PLAYERS  = 1200
	MAX_NUMBER_OF_TABLES   = MAX_NUMBER_OF_PLAYERS / 2
	PAIRING_GROUP_MIN_SIZE = 100
	PAIRING_GROUP_MAX_SIZE = 3 * PAIRING_GROUP_MIN_SIZE
)

type Gotha struct {
	tournament *Tournament
	IO         *InputOutput
}

func NewGotha() *Gotha {
	gotha := new(Gotha)
	return gotha
}

// 选择编排模式
func (g *Gotha) SelectSystem(systemStr string) {
	tps := g.tournament.tournamentParameterSet
	system := ConvertSystem(systemStr)
	switch system {
	case TYPE_UNDEFINED:
	case TYPE_SWISS:
		tps.InitForSwiss()
	}
	g.tournament.SetTournamentSet(tps)
}

func (g *Gotha) ImportFromXMLFile(filePath string) error {
	file, err := os.Open(filePath)
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}
	return g.ImportFromReader(file)
}

func (g *Gotha) FlushGameToXML() (io.Reader, error) {
	return g.IO.FlushGameToXML(g.tournament.SortGameByTableNumber())
}

func (g *Gotha) ImportFromBytes(bs []byte) error {
	return g.ImportFromReader(bytes.NewReader(bs))
}

func (g *Gotha) ImportFromReader(io io.Reader) error {
	gothaIO := NewInputOutput()
	g.IO = gothaIO
	gothaIO.WithOption(WithPlayers(), WithGames())
	g.tournament = NewTournament()
	return g.IO.ImportFromReader(io, g.tournament)
}

func (g *Gotha) GetTournament() *Tournament {
	return g.tournament
}

func (g *Gotha) SetTournament(t *Tournament) {
	g.tournament = t
}
