package gotha

import . "tournament_pair/src/parameter_set"

const (
	MAX_NUMBER_OF_ROUNDS = 20
)

type Gotha struct {
	tournament *Tournament
}

func NewGotha() *Gotha {
	gotha := new(Gotha)
	return gotha
}

// 选择编排模式
func (g *Gotha) SelectSystem(system int) {
	tps := NewTournamentParameterSet()
	switch system {
	case TYPE_UNDEFINED:
	case TYPE_SWISS:
		tps.InitForSwiss()
	}
	g.tournament.SetTournamentSet(tps)
}

func (g *Gotha) GetFromXMLFile(filePath string) {
	input := NewInput()
	input.WithOption(WithPlayers())
	g.tournament = NewTournament()
	input.ImportTournamentFromXMLFile(filePath, g.tournament)
}

func (g *Gotha) GetTournament() *Tournament {
	return g.tournament
}
