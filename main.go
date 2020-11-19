package main

import (
	"tournament_pair/src/gotha"
	"tournament_pair/src/parameter_set"
)

func main() {
	g := gotha.NewGotha()
	g.GetFromXMLFile(`/Users/justice/Desktop/clm2014B.xml`)
	// 选择编排类型
	g.SelectSystem(parameter_set.TYPE_SWISS)
	t := g.GetTournament()

	t.MakeAutomaticPairing(t.GetPlayers(), 0)
}
