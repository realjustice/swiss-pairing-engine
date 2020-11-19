package parameter_set

type HandicapParameterSet struct {
	hdBasedOnMMS        bool
	hdNoHdRankThreshold int
	hdCorrection        int
	hdCeiling           int
}

func NewHandicapParameterSet() *HandicapParameterSet {
	set := new(HandicapParameterSet)
	set.hdNoHdRankThreshold = 0
	set.hdCorrection = 1

	return set
}

func (set *HandicapParameterSet) InitForSwiss() {
	set.hdBasedOnMMS = false
	set.hdNoHdRankThreshold = -30
	set.hdCorrection = 0
	set.hdCeiling = 0
}

func (set *HandicapParameterSet) GetHdNoHdRankThreshold() int {
	return set.hdNoHdRankThreshold
}

func (set *HandicapParameterSet) GetHdCorrection() int {
	return set.hdCorrection
}

func (set *HandicapParameterSet) GetHdCeiling() int {
	return set.hdCeiling
}
