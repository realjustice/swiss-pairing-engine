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
	} else if result < 0 {
		return true
	} else if result == 0 {
		if compareTo(spc.ScoredPlayers[i].Name, spc.ScoredPlayers[j].Name) > 0 {
			return false
		} else if compareTo(spc.ScoredPlayers[i].Name, spc.ScoredPlayers[j].Name) < 0 {
			return true
		}

		if compareTo(spc.ScoredPlayers[i].FirstName, spc.ScoredPlayers[j].FirstName) > 0 {
			return false
		} else if compareTo(spc.ScoredPlayers[i].FirstName, spc.ScoredPlayers[j].FirstName) < 0 {
			return true
		}
	}

	return false
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

func compareTo(a string, b string) int {
	// 将两个字符串变成等长的char数组
	charA := []byte(a)
	charB := []byte(b)

	// 使用strings.Repeat来填充字符串到相同长度
	if len(a) < len(b) {
		charA = append(charA, make([]byte, len(b)-len(a))...)
	} else if len(a) > len(b) {
		charB = append(charB, make([]byte, len(a)-len(b))...)
	}

	// 之后比较
	for i := 0; i < len(charA); i++ { // 使用扩展后的字符数组长度
		if charA[i] > charB[i] {
			return 1
		} else if charA[i] < charB[i] {
			return -1
		}
	}
	return 0
}
