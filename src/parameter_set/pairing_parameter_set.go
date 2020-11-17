package parameter_set

const (
	PAIBA_MAX_AVOIDDUPLGAME = 5 * 10e14 // 5 * 10^14
)

// 匹配参数
type PairingParameterSet struct {
	PaiMaAdditionalPlacementCritSystem1 int
	PaiBaAvoidDuplGame                  int64
	paiBaDeterministic                  bool
}

func NewPairingParameterSet() *PairingParameterSet {
	set := new(PairingParameterSet)
	set.PaiMaAdditionalPlacementCritSystem1 = PLA_CRIT_RATING
	set.PaiBaAvoidDuplGame = PAIBA_MAX_AVOIDDUPLGAME
	set.paiBaDeterministic = true
	return set
}

func (p *PairingParameterSet) GetPaiBaAvoidDuplGame() int64 {
	return p.PaiBaAvoidDuplGame
}

func (p *PairingParameterSet) IsPaiBaDeterministic() bool {
	return p.paiBaDeterministic
}