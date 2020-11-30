package parameter_set

const (
	TYPE_UNDEFINED = 0
	TYPE_MCMAHON   = 1
	TYPE_SWISS     = 2
	TYPE_SWISSCAT  = 3
)

func ConvertSystem(system string) int {
	switch system {
	case "SWISS":
		return TYPE_SWISS
	default:
		return TYPE_SWISS
	}
}

type TournamentParameterSet struct {
	generalParameterSet   *GeneralParameterSet
	pairingParameterSet   *PairingParameterSet
	placementParameterSet *PlacementParameterSet
	handicapParameterSet  *HandicapParameterSet
}

func NewTournamentParameterSet() *TournamentParameterSet {
	set := &TournamentParameterSet{}
	set.generalParameterSet = NewGeneralParameterSet()
	set.pairingParameterSet = NewPairingParameterSet()
	set.placementParameterSet = NewPlacementParameterSet()
	set.handicapParameterSet = NewHandicapParameterSet()
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

// 瑞士制编排初始化
func (t *TournamentParameterSet) InitForSwiss() {
	t.generalParameterSet.InitSwiss()
	t.placementParameterSet.InitForSwiss()
	t.handicapParameterSet.InitForSwiss()
	t.pairingParameterSet.InitForSwiss()
}

func (t *TournamentParameterSet) SetGeneralParameterSet(set *GeneralParameterSet) {
	t.generalParameterSet = set
}

func (t *TournamentParameterSet) GetHandicapParameterSet() *HandicapParameterSet {
	return t.handicapParameterSet
}
