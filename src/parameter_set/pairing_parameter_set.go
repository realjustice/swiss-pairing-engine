package parameter_set

const (
	PAIBA_MAX_AVOIDDUPLGAME             = 5 * 1e14 // 5 * 10^14
	PAIBA_MAX_BALANCEWB                 = 1e6
	PAIMA_MAX_MINIMIZE_SCORE_DIFFERENCE = 1e11
	PAIMA_MAX_MAXIMIZE_SEEDING          = PAIMA_MAX_MINIMIZE_SCORE_DIFFERENCE / 20000
	PAIMA_MAX_DUDD_WEIGHT               = PAIMA_MAX_MINIMIZE_SCORE_DIFFERENCE / 1000
	PAIMA_SEED_SPLITANDRANDOM           = 1
	PAIMA_SEED_SPLITANDFOLD             = 2
	PAIMA_SEED_SPLITANDSLIP             = 3
)

// 匹配参数
type PairingParameterSet struct {
	// 额外规则
	paiMaAdditionalPlacementCritSystem1 int
	paiMaAdditionalPlacementCritSystem2 int

	// minimize score 系数
	paiStandardNX1Factor float64

	// minimize score 分差
	paiMaMinimizeScoreDifference int64

	// 种子因
	paiMaMaximizeSeeding int64

	paiBaBalanceWB int64

	// 额外规则所适用的最后一个轮次
	paiMaLastRoundForSeedSystem1 int
	paiMaSeedSystem1             int
	paiMaSeedSystem2             int

	paiBaAvoidDuplGame int64
	paiBaDeterministic bool
}

func NewPairingParameterSet() *PairingParameterSet {
	set := new(PairingParameterSet)
	set.paiMaAdditionalPlacementCritSystem1 = PLA_CRIT_RATING
	set.paiMaAdditionalPlacementCritSystem2 = PLA_CRIT_NUL
	set.paiBaAvoidDuplGame = PAIBA_MAX_AVOIDDUPLGAME
	set.paiBaDeterministic = true
	set.paiMaLastRoundForSeedSystem1 = 1
	set.paiBaBalanceWB = PAIBA_MAX_BALANCEWB
	set.paiStandardNX1Factor = 0.5
	set.paiMaMinimizeScoreDifference = PAIMA_MAX_MINIMIZE_SCORE_DIFFERENCE
	set.paiMaMaximizeSeeding = PAIMA_MAX_MAXIMIZE_SEEDING
	set.paiMaSeedSystem1 = PAIMA_SEED_SPLITANDRANDOM
	set.paiMaSeedSystem2 = PAIMA_SEED_SPLITANDFOLD

	return set
}

func (p *PairingParameterSet) InitForSwiss() {
	p.paiBaAvoidDuplGame = PAIBA_MAX_AVOIDDUPLGAME
	//p.paiBaRandom = 0
	p.paiBaDeterministic = true
	p.paiBaBalanceWB = PAIBA_MAX_BALANCEWB
	p.paiMaMinimizeScoreDifference = PAIMA_MAX_MINIMIZE_SCORE_DIFFERENCE

	p.paiMaSeedSystem1 = PAIMA_SEED_SPLITANDSLIP
	p.paiMaSeedSystem2 = PAIMA_SEED_SPLITANDSLIP
	//p.paiSeDefSecCrit = PAIMA_MAX_AVOID_MIXING_CATEGORIES
	//p.paiMaDUDDWeight = PAIMA_MAX_DUDD_WEIGHT
	//p.paiSeMinimizeHandicap = 0
	//p.paiSeAvoidSameGeo = 0

}

func (p *PairingParameterSet) GetPaiBaAvoidDuplGame() int64 {
	return p.paiBaAvoidDuplGame
}

func (p *PairingParameterSet) IsPaiBaDeterministic() bool {
	return p.paiBaDeterministic
}

func (p *PairingParameterSet) GetPaiMaAdditionalPlacementCritSystem1() int {
	return p.paiMaAdditionalPlacementCritSystem1
}

func (p *PairingParameterSet) GetPaiMaAdditionalPlacementCritSystem2() int {
	return p.paiMaAdditionalPlacementCritSystem2
}

func (p *PairingParameterSet) GetPaiMaLastRoundForSeedSystem1() int {
	return p.paiMaLastRoundForSeedSystem1
}

func (p *PairingParameterSet) GetPaiBaBalanceWB() int64 {
	return p.paiBaBalanceWB
}

func (p *PairingParameterSet) GetPaiStandardNX1Factor() float64 {
	return p.paiStandardNX1Factor
}

func (p *PairingParameterSet) GetPaiMaMinimizeScoreDifference() int64 {
	return p.paiMaMinimizeScoreDifference
}

func (p *PairingParameterSet) GetPaiMaMaximizeSeeding() int64 {
	return p.paiMaMaximizeSeeding
}

func (p *PairingParameterSet) GetPaiMaSeedSystem1() int {
	return p.paiMaSeedSystem1
}

func (p *PairingParameterSet) GetPaiMaSeedSystem2() int {
	return p.paiMaSeedSystem2
}
