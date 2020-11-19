package gotha

type ScoredPlayerComparator struct {
	ScoredPlayers
	criterion   []int
	roundNumber int
}

func NewScoredPlayerComparator(scoredPlayers ScoredPlayers, criterion []int, roundNumber int) *ScoredPlayerComparator {
	comparator := new(ScoredPlayerComparator)
	comparator.ScoredPlayers = scoredPlayers
	newCriterion := make([]int, len(criterion))
	copy(newCriterion, criterion)
	comparator.criterion = newCriterion
	comparator.roundNumber = roundNumber
	return comparator
}

func (spc *ScoredPlayerComparator) Len() int {
	return len(spc.ScoredPlayers)
}

func (spc *ScoredPlayerComparator) Less(i, j int) bool {
	result := spc.BetterScore(spc.ScoredPlayers[i], spc.ScoredPlayers[j])
	if result > 0 {
		return false
	}

	return true
}

func (spc ScoredPlayerComparator) Swap(i, j int) {
	spc.ScoredPlayers[i], spc.ScoredPlayers[j] = spc.ScoredPlayers[j], spc.ScoredPlayers[i]
}

func (spc *ScoredPlayerComparator) BetterScore(sp1 *ScoredPlayer, sp2 *ScoredPlayer) int {
	for cr := 0; cr < len(spc.criterion); cr++ {
		if sp1.GetCritValue(spc.criterion[cr], spc.roundNumber) < sp2.GetCritValue(spc.criterion[cr], spc.roundNumber) {
			return 1
		} else if sp1.GetCritValue(spc.criterion[cr], spc.roundNumber) > sp2.GetCritValue(spc.criterion[cr], spc.roundNumber) {
			return -1
		}
	}
	return 0
}
