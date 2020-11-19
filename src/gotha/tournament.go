package gotha

import (
	"sort"
	"tournament_pair/src/parameter_set"
)

type Tournament struct {
	tournamentParameterSet *parameter_set.TournamentParameterSet
	/**
	 * HashMap of Players The key is the getKeyString
	 */
	hmPlayers map[string]*Player
	/**
	 * HashMap of Games The key is (roundNumber * Gotha.MAX_NUMBER_OF_TABLES +
	 * tableNumber)
	 */
	hmGames map[string]*Game

	hmScoredPlayers map[string]*ScoredPlayer

	byePlayers []*Player
}

func NewTournament() *Tournament {
	t := new(Tournament)
	t.tournamentParameterSet = parameter_set.NewTournamentParameterSet()
	// 添加报名人员
	t.hmPlayers = make(map[string]*Player)

	t.hmGames = make(map[string]*Game)
	t.byePlayers = make([]*Player, MAX_NUMBER_OF_ROUNDS)

	return t
}

func (t *Tournament) GetTournamentSet() *parameter_set.TournamentParameterSet {
	return t.tournamentParameterSet
}
func (t *Tournament) SetTournamentSet(set *parameter_set.TournamentParameterSet) {
	t.tournamentParameterSet = set
}

func (t *Tournament) AddPlayer(players []*Player) {
	for _, p := range players {
		t.hmPlayers[p.SetKeyString()] = p
	}
}

func (t *Tournament) FillPairingInfo(roundNumber int) {
	gps := t.tournamentParameterSet.GetGeneralParameterSet()
	pps := t.tournamentParameterSet.GetPlacementParameterSet()
	paiPs := t.tournamentParameterSet.GetPairingParameterSet()

	mainScoreMin := 0
	mainScoreMax := 0
	groupNumber := 0 // 有几轮就有几个groupNumber

	for cat := 0; cat < gps.GetNumberOfCategories(); cat++ {
		for mainScore := mainScoreMax; mainScore >= mainScoreMin; mainScore-- {
			alSPGroup := make([]*ScoredPlayer, 0)
			for _, sp := range t.hmScoredPlayers {
				if sp.Category(gps) != cat {
					continue
				}
				// 获取到本轮为止的全胜人员
				// 第一轮全体选手
				if sp.GetCritValue(100, roundNumber-1)/2 != mainScore {
					continue
				}

				alSPGroup = append(alSPGroup, sp)
			}
			if len(alSPGroup) <= 0 {
				continue
			}
			// 压入一个规则，（第一轮的时候增加rating规则，其余情况不压入规则！！）
			crit := pps.GetPlaCriteria()
			additionalCrit := paiPs.GetPaiMaAdditionalPlacementCritSystem1()
			if roundNumber > paiPs.GetPaiMaLastRoundForSeedSystem1() {
				additionalCrit = paiPs.GetPaiMaAdditionalPlacementCritSystem2()
			}

			paiCrit := make([]int, len(crit)+1)
			copy(crit, paiCrit)
			paiCrit[len(paiCrit)-1] = additionalCrit
			// 根据规则排序
			spc := NewScoredPlayerComparator(alSPGroup, paiCrit, roundNumber-1)
			sort.Sort(spc)

			for index, sp := range alSPGroup {
				sp.GroupNumber = groupNumber // 所记录的就是最后一轮全胜轮次
				sp.GroupSize = len(alSPGroup)
				sp.InnerPlacement = index
			}
			groupNumber++
		}
	}

	numberOfGroups := groupNumber

	for _, sp := range t.hmScoredPlayers {
		sp.NumberOfGroups = numberOfGroups
	}

	// 重置上下调
	for _, sp := range t.hmScoredPlayers {
		sp.NbDU = 0
		sp.NbDD = 0
	}

	if roundNumber > 1 {
		// prepare an Array of scores before round r
		// 计算上下调
	}
}

func (t *Tournament) MakeAutomaticPairing(alPlayersToPair []*Player, roundNumber int) ([]*Game, bool) {
	// not even
	if len(alPlayersToPair)%2 != 0 {
		return nil, false
	}

	t.fillBaseScoringInfoIfNecessary()

	// todo getGamesListBefore()
	alPreviousGames := t.gamesListBefore(roundNumber)

	// fill pairing info
	t.FillPairingInfo(roundNumber)

	// todo 大于300人的比赛暂未处理

	alRemainingPlayers := make([]*Player, len(alPlayersToPair))
	copy(alRemainingPlayers, alPlayersToPair)

	alg := t.pairAGroup(alRemainingPlayers, roundNumber, t.hmScoredPlayers, alPreviousGames)
	alGames := make([]*Game, len(alg))
	// fill game
	copy(alGames, alg)

	return alGames, true
}

func (t *Tournament) fillBaseScoringInfoIfNecessary() {
	// 0) Preparation
	// **************
	gps := t.tournamentParameterSet.GetGeneralParameterSet()
	t.hmScoredPlayers = make(map[string]*ScoredPlayer)
	for _, p := range t.hmPlayers {
		sp := NewScoredPlayer(gps, p)
		t.hmScoredPlayers[p.GetKeyString()] = sp
	}
	numberOfRoundsToCompute := gps.GetNumberOfRounds()

	// 1) participation 编排状态
	// ****************

	// 初始化编排状态
	for _, sp := range t.hmScoredPlayers {
		for r := 0; r < numberOfRoundsToCompute; r++ {
			if !sp.GetParticipating()[r] {
				sp.SetParticipation(r, ABSENT) //  弃权选手
			} else {
				sp.SetParticipation(r, NOT_ASSIGNED) // As an initial status 初始状态
			}
		}
	}

	// 更新编排状态为已编排
	for _, g := range t.hmGames {
		wp := g.GetWhitePlayer()
		bp := g.GetBlackPlayer()
		if wp == nil {
			continue
		}

		if bp == nil {
			continue
		}
		r := g.GetRoundNumber()
		wSP := t.hmScoredPlayers[wp.GetKeyString()]
		bSP := t.hmScoredPlayers[bp.GetKeyString()]
		// 更新编排状态
		wSP.SetParticipation(r, PAIRED)
		bSP.SetParticipation(r, PAIRED)
		// 将game放入scoredPlayer中
		wSP.SetGame(r, g)
		bSP.SetGame(r, g)
	}

	// 设置轮空人员编排状态
	for r := 0; r < numberOfRoundsToCompute; r++ {
		p := t.byePlayers[r]
		if p != nil {
			sp := t.hmScoredPlayers[p.GetKeyString()]
			sp.SetParticipation(r, BYE)
		}
	}

	// 2) nbwX2  (计算总分)
	for r := 0; r < numberOfRoundsToCompute; r++ {
		// Initialize
		for _, sp := range t.hmScoredPlayers {
			if r == 0 {
				sp.SetNBWX2(r, 0)
			} else {
				sp.SetNBWX2(r, sp.GetNBWX2(r-1))
			}
		}

		// Points from games 根据对局计算出每个选手每一轮的大分
		for _, g := range t.hmGames {
			if g.RoundNumber != r {
				continue
			}
			wP := g.GetWhitePlayer()
			bP := g.GetBlackPlayer()
			if wP == nil {
				continue
			}
			if bP == nil {
				continue
			}

			wSP := t.hmScoredPlayers[wP.GetKeyString()]
			bSP := t.hmScoredPlayers[bP.GetKeyString()]

			// select result
			selectResult(bSP, wSP, g, r)
		}

		for _, sp := range t.hmScoredPlayers {
			// 弃权或者轮空
			nbPtsNBW2AbsentOrBye := 0
			for r := 0; r < numberOfRoundsToCompute; r++ {
				if sp.GetParticipation(r) == ABSENT { // 弃权
					nbPtsNBW2AbsentOrBye += gps.GetGenNBW2ValueAbsent()
				}
				if sp.GetParticipation(r) == BYE {
					nbPtsNBW2AbsentOrBye += gps.GetGenNBW2ValueAbsent()
				}

				sp.SetNBWX2(r, sp.GetNBWX2(r)+nbPtsNBW2AbsentOrBye)
			}

		}

		// 4.1) SOSW, SOSWM1, SOSWM2,SODOSW
		// soswX2 Sum of Opponents nbw2 对手总分
		// soswM1X2 Sum of (n-1) Opponents nbw2 上一轮对手大分总和
		// soswM2X2 Sum of (n-2) Opponents nbw2 上上轮对手大分总和
		// sdswX4  Sum of Defeated Opponents nbw2 击败的对手分*2

		for r := 0; r < numberOfRoundsToCompute; r++ {
			for _, sp := range t.hmScoredPlayers {
				oswX2 := make([]int, numberOfRoundsToCompute)  // 对手每一轮的大分
				doswX4 := make([]int, numberOfRoundsToCompute) // Defeated opponents score
				for rr := 0; rr <= r; rr++ {
					if sp.GetParticipation(rr) != PAIRED {
						oswX2[rr] = 0
						doswX4[rr] = 0
					} else {
						g := sp.GetGame(rr)
						opp := t.opponent(g, sp.Player)
						// 如果选手胜，result=2 ，和棋 result=1 else 0
						result := getWX2(g, sp.Player)

						sOpp := t.hmScoredPlayers[opp.GetKeyString()]
						oswX2[rr] = sOpp.GetNBWX2(r)    // 对手的总分
						doswX4[rr] = oswX2[rr] * result // 对手的总分*result
					}
				}
				// 计算 soswM2X2，sdswX4
				sosX2 := 0
				sdsX4 := 0
				// 为什么要再开一个循环？？
				for rr := 0; rr <= r; rr++ {
					sosX2 += oswX2[rr]
					sdsX4 += doswX4[rr]
				}
				sp.SetSOSWX2(r, sosX2)
				sp.SetSDSWX4(r, sdsX4)

				sosM1X2 := 0
				sosM2X2 := 0
				// soswM1X2 Sum of (n-1) Opponents nbw2 上一轮对手大分总和
				// soswM2X2 Sum of (n-2) Opponents nbw2 上上轮对手大分总和

				if r == 0 {
					sosM1X2 = 0
					sosM2X2 = 0
				} else if r == 1 {
					sosM1X2 = max(oswX2[0], oswX2[1]) // oswX2>0 sosM1X2=oswX2 else sosM1X2=0
					sosM2X2 = 0
				} else {
					rMin := 0
					for rr := 1; rr <= r; rr++ {
						if oswX2[rr] < oswX2[rMin] {
							rMin = rr
						}
					}
					rMin2 := 0
					if rMin == 0 {
						rMin2 = 1
					}
					for rr := 0; rr <= r; rr++ {
						if rr == rMin {
							continue
						}
						if oswX2[rr] < oswX2[rMin2] {
							rMin2 = rr
						}
					}
					// 一顿看不懂的操作之后，获取到了上轮的对手总分，和上上轮对手的总分
					sosM1X2 = sp.GetSOSWX2(r) - oswX2[rMin]
					sosM2X2 = sosM1X2 - oswX2[rMin2]
				}
				sp.SetSOSWM1X2(r, sosM1X2)
				sp.SetSOSWM2X2(r, sosM2X2)
			}
		}

		// 5) ssswX2(对应sososwX2)  Sum of opponents sosw2 * 2 对手的SOSW总和
		for r := 0; r < numberOfRoundsToCompute; r++ {
			for _, sp := range t.hmScoredPlayers {
				sososwX2 := 0
				for rr := 0; rr <= r; rr++ {
					//！！ 只有选手在该轮次下处于匹配状态的时候才需要进行ssswX2的计算
					if sp.GetParticipation(rr) != PAIRED {
						sososwX2 += 0
					} else {
						g := sp.GetGame(rr)
						// 计算对手分
						opp := NewPlayer()
						if g.GetWhitePlayer().HasSameKeyString(sp.Player) {
							opp = g.GetBlackPlayer()
						} else {
							opp = g.GetWhitePlayer()
						}
						sOpp := t.hmScoredPlayers[opp.GetKeyString()]
						sososwX2 += sOpp.GetSOSWX2(r)
					}
				}
				sp.SetSSSWX2(r, sososwX2)
			}
		}
	}
}

func selectResult(bSP *ScoredPlayer, wSP *ScoredPlayer, g *Game, r int) {
	switch g.GetResult() {
	case RESULT_BOTHLOSE, RESULT_BOTHLOSE_BYDEF, RESULT_UNKNOWN:
	case RESULT_WHITEWINS, RESULT_WHITEWINS_BYDEF:
		wSP.SetNBWX2(r, wSP.GetNBWX2(r)+2)
	case RESULT_BLACKWINS, RESULT_BLACKWINS_BYDEF:
		bSP.SetNBWX2(r, bSP.GetNBWX2(r)+2)
	case RESULT_EQUAL, RESULT_EQUAL_BYDEF:
		wSP.SetNBWX2(r, wSP.GetNBWX2(r)+1)
		bSP.SetNBWX2(r, bSP.GetNBWX2(r)+1)
	case RESULT_BOTHWIN, RESULT_BOTHWIN_BYDEF:
		wSP.SetNBWX2(r, wSP.GetNBWX2(r)+2)
		bSP.SetNBWX2(r, bSP.GetNBWX2(r)+2)
	}
}

func (t *Tournament) opponent(g *Game, p *Player) *Player {
	if g == nil {
		return nil
	}
	wP := g.GetWhitePlayer()
	bP := g.GetBlackPlayer()
	if wP.HasSameKeyString(p) {
		return bP
	} else if bP.HasSameKeyString(p) {
		return wP

	}
	return nil
}

func getWX2(g *Game, p *Player) int {
	wX2 := 0
	wP := g.GetWhitePlayer()
	bP := g.GetBlackPlayer()
	pIsWhite := true
	if wP.HasSameKeyString(p) {
		pIsWhite = true
	} else if bP.HasSameKeyString(p) {
		pIsWhite = false
	} else {
		return 0
	}
	switch g.GetResult() {
	case RESULT_BOTHLOSE, RESULT_BOTHLOSE_BYDEF, RESULT_UNKNOWN:
		wX2 = 0
	case RESULT_WHITEWINS, RESULT_WHITEWINS_BYDEF:
		if pIsWhite {
			wX2 = 2
		}
	case RESULT_BLACKWINS, RESULT_BLACKWINS_BYDEF:
		if !pIsWhite {
			wX2 = 2
		}
	case RESULT_EQUAL, RESULT_EQUAL_BYDEF:
		wX2 = 1
	case RESULT_BOTHWIN, RESULT_BOTHWIN_BYDEF:
		wX2 = 2
	}
	return wX2
}

func (t *Tournament) gamesListBefore(roundNumber int) []*Game {
	gL := make([]*Game, 0)
	for _, g := range t.hmGames {
		if g.GetRoundNumber() < roundNumber {
			gL = append(gL, g)
		}
	}
	return gL
}

func (t *Tournament) pairAGroup(alGroupedPlayers []*Player, roundNumber int, hmScoredPlayers map[string]*ScoredPlayer, alPreviousGames []*Game) []*Game {
	// 分组数量
	numberOfPlayersInGroup := len(alGroupedPlayers)

	//  Prepare infos about Score groups : sgSize, sgNumber and innerPosition
	//  And DUDD information
	costs := make([][]int64, 0)
	// 生成一个二维数组，替代java中 long[numberOfPlayersInGroup][numberOfPlayersInGroup]
	for i := 0; i < numberOfPlayersInGroup; i++ {
		costs = append(costs, make([]int64, numberOfPlayersInGroup))
	}

	for i := 0; i < numberOfPlayersInGroup; i++ {
		costs[i][i] = 0
		for j := i + 1; j < numberOfPlayersInGroup; j++ {
			p1 := alGroupedPlayers[i]
			p2 := alGroupedPlayers[j]
			sP1 := t.hmScoredPlayers[p1.GetKeyString()]
			sP2 := t.hmScoredPlayers[p2.GetKeyString()]
			costs[i][j] = t.costValue(sP1, sP2, roundNumber, alPreviousGames)
			costs[j][i] = costs[i][j]
		}
	}

	//mate := make([]int, 0)

	alG := make([]*Game, 0)
	for i := 1; i <= len(costs); i++ {
		//if i < mate[i] {
		//	p1 := alGroupedPlayers[i-1]
		//	p2 := alGroupedPlayers[mate[i]-1]
		//	sP1 := t.hmScoredPlayers[p1.GetKeyString()]
		//	sP2 := t.hmScoredPlayers[p2.GetKeyString()]
		//	g :=
		//	alg := append(alG, g)
		//}

	}

	return alG
}

func (t *Tournament) costValue(sP1 *ScoredPlayer, sP2 *ScoredPlayer, roundNumber int, alPreviousGames []*Game) int64 {
	gps := t.tournamentParameterSet.GetGeneralParameterSet()
	paiPS := t.tournamentParameterSet.GetPairingParameterSet()
	var cost int64 = 1 // 1 is minimum value because 0 means "no matching allowed"

	// 是否匹配过
	// Base Criterion 1 : Avoid Duplicating Game
	// Did p1 and p2 already play ?
	numberOfPreviousGamesP1P2 := 0
	for r := 0; r < roundNumber; r++ {
		g1 := sP1.GetGame(r)
		if g1 == nil {
			continue
		}
		if sP1.HasSameKeyString(g1.GetWhitePlayer()) && sP2.HasSameKeyString(g1.GetBlackPlayer()) {
			numberOfPreviousGamesP1P2++
		}
		if sP1.HasSameKeyString(g1.GetBlackPlayer()) && sP2.HasSameKeyString(g1.GetWhitePlayer()) {
			numberOfPreviousGamesP1P2++
		}
	}
	// 如果之前从未匹配过，则会加上一个非常大的权重
	if numberOfPreviousGamesP1P2 == 0 {
		cost += paiPS.GetPaiBaAvoidDuplGame()
	}

	// 增加随机因子(默认不使用)
	//Base Criterion 2 : Random
	var nR int64
	if paiPS.IsPaiBaDeterministic() {

	} else {

	}
	cost += nR

	// 黑白平衡
	// Base Criterion 3 : Balance W and B
	// This cost is never applied if potential Handicap != 0
	// It is fully applied if wbBalance(sP1) and wbBalance(sP2) are strictly of different signs
	// It is half applied if one of wbBalance is 0 and the other is >=2

	var bwBalanceCost int64 = 0
	g := t.gameBetween(sP1, sP2, roundNumber)
	poHd := g.GetHandicap()
	if poHd == 0 {
		wb1 := wbBalance(sP1, roundNumber-1)
		wb2 := wbBalance(sP2, roundNumber-1)
		// 双方都是黑居多
		if wb1*wb2 < 0 {
			bwBalanceCost = paiPS.GetPaiBaBalanceWB()
		} else if wb1 == 0 && abs(wb2) >= 2 {
			bwBalanceCost = paiPS.GetPaiBaBalanceWB() / 2
		} else if wb2 == 0 && abs(wb1) >= 2 {
			bwBalanceCost = paiPS.GetPaiBaBalanceWB() / 2
		}
	}
	cost += bwBalanceCost

	// Main Criterion 2 : Minimize score difference
	var scoCost int64 = 0
	scoRange := sP1.NumberOfGroups
	if sP1.Category(gps) == sP2.Category(gps) {
		x := float64(abs(sP1.GroupNumber-sP2.GroupNumber)) / float64(scoRange)
		k := paiPS.GetPaiStandardNX1Factor()
		scoCost = int64(float64(paiPS.GetPaiMaMinimizeScoreDifference()) * (1.0 - x) * (1.0 + k*x))
	}
	cost += scoCost
	// todo Main Criterion 3 : If different groups, make a directed Draw-up/Draw-down

	// Main Criterion 4 : Seeding
	var seedCost int64 = 0
	if sP1.GroupNumber == sP2.GroupNumber {
		groupSize := int64(sP1.GroupSize)
		cla1 := int64(sP1.InnerPlacement)
		cla2 := int64(sP2.InnerPlacement)
		maxSeedingWeight := paiPS.GetPaiMaMaximizeSeeding()
		currentSeedSystem := paiPS.GetPaiMaSeedSystem2()
		if roundNumber <= paiPS.GetPaiMaLastRoundForSeedSystem1() {
			currentSeedSystem = paiPS.GetPaiMaSeedSystem1()
		}
		// select seed system
		if currentSeedSystem == parameter_set.PAIMA_SEED_SPLITANDRANDOM {
			// todo 暂时用不到

		} else if currentSeedSystem == parameter_set.PAIMA_SEED_SPLITANDFOLD {
			// The best is to get cla1 + cla2 - (groupSize - 1) close to 0
			x := int64(cla1 + cla2 - (groupSize - 1))
			seedCost = maxSeedingWeight - (maxSeedingWeight * x / (groupSize - 1) * x / (groupSize - 1))
		} else if currentSeedSystem == parameter_set.PAIMA_SEED_SPLITANDSLIP {
			// The best is to get 2 * |Cla1 - Cla2| - groupSize    close to 0
			x := 2*abs64(cla1-cla2) - groupSize
			seedCost = maxSeedingWeight - (maxSeedingWeight * x / groupSize * x / groupSize)
		}
	}
	cost += seedCost

	return cost
}

func (t *Tournament) GetPlayers() []*Player {
	players := make([]*Player, 0)
	for _, p := range t.hmPlayers {
		players = append(players, p)
	}
	return players
}

/**
 * builds and return a new Game with everything defined except tableNumber
 */
func (t *Tournament) gameBetween(sP1 *ScoredPlayer, sP2 *ScoredPlayer, roundNumber int) *Game {
	hdPS := t.tournamentParameterSet.GetHandicapParameterSet()
	g := NewGame()
	hd := 0
	pseudoRank1 := sP1.GetRank()
	pseudoRank2 := sP2.GetRank()
	// 取player段位和阈值
	pseudoRank1 = min(pseudoRank1, hdPS.GetHdNoHdRankThreshold())
	pseudoRank2 = min(pseudoRank2, hdPS.GetHdNoHdRankThreshold())
	hd = pseudoRank1 - pseudoRank2
	if hd > 0 {
		hd = hd - hdPS.GetHdCorrection()
		if hd < 0 {
			hd = 0
		}
	}

	if hd < 0 {
		hd = hd + hdPS.GetHdCorrection()
		if hd > 0 {
			hd = 0
		}
	}

	if hd > hdPS.GetHdCeiling() {
		hd = hdPS.GetHdCeiling()
	}
	if hd < -hdPS.GetHdCeiling() {
		hd = -hdPS.GetHdCeiling()
	}

	p1 := t.GetPlayerByKeyString(sP1.GetKeyString())
	p2 := t.GetPlayerByKeyString(sP2.GetKeyString())

	if hd > 0 {
		g.SetWhitePlayer(p1)
		g.SetBlackPlayer(p2)
		g.SetHandicap(hd)
	} else if hd < 0 {
		g.SetWhitePlayer(p2)
		g.SetBlackPlayer(p1)
		g.SetHandicap(-hd)
	} else {
		g.SetHandicap(0)
		if wbBalance(sP1, roundNumber-1) > wbBalance(sP2, roundNumber-1) {
			g.SetWhitePlayer(p2)
			g.SetBlackPlayer(p1)
		} else if wbBalance(sP1, roundNumber-1) > wbBalance(sP2, roundNumber-1) {
			g.SetWhitePlayer(p1)
			g.SetBlackPlayer(p2)
		} else {
			g.SetWhitePlayer(p2)
			g.SetBlackPlayer(p1)
		}
	}

	g.SetKnownColor(true)
	g.SetResult(RESULT_UNKNOWN)
	g.SetRoundNumber(roundNumber)

	return g
}

func (t *Tournament) GetPlayerByKeyString(strNaFi string) *Player {
	return t.hmPlayers[strNaFi]
}
