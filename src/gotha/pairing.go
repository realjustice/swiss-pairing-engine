package gotha

/***
  黑白平衡 balance>0 白多 balance<0 黑多
*/
func wbBalance(sP *ScoredPlayer, rn int) int {
	if rn < 0 {
		return 0
	}
	balance := 0
	for r := 0; r <= rn; r++ {
		g := sP.GetGame(r)
		if g == nil {
			continue
		}
		if g.GetHandicap() != 0 {
			continue
		}
		if sP.HasSameKeyString(g.GetWhitePlayer()) {
			balance++
		} else {
			balance--
		}
	}
	return balance
}
