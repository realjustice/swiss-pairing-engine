package parameter_set

const (
	TYPE_UNDEFINED = 0
	TYPE_MCMAHON   = 1
	TYPE_SWISS     = 2
	TYPE_SWISSCAT  = 3
)

type TournamentParameterSet struct {
	generalParameterSet   *GeneralParameterSet
	pairingParameterSet   *PairingParameterSet
	placementParameterSet *PlacementParameterSet
}

func NewTournamentParameterSet() *TournamentParameterSet {
	set := &TournamentParameterSet{}
	set.generalParameterSet = NewGeneralParameterSet()
	set.pairingParameterSet = NewPairingParameterSet()
	set.placementParameterSet = NewPlacementParameterSet()
	return set
}

func (t *TournamentParameterSet) GetPlacementParameterSet() *PlacementParameterSet {
	return t.placementParameterSet
}

func (t *TournamentParameterSet) GetGeneralParameterSet() *GeneralParameterSet {
	return t.generalParameterSet
}

func (t *TournamentParameterSet) GetPairingParameterSet() *PairingParameterSet {
	return t.pairingParameterSet
}

func (t *TournamentParameterSet) InitForSwiss() {
	t.placementParameterSet.InitForSwiss()
}
