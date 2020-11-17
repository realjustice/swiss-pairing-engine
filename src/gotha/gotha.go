package gotha

import . "tournament_pair/src/parameter_set"

type Gotha struct {
}

func NewGotha(tournament *Tournament) *Gotha {
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
}
