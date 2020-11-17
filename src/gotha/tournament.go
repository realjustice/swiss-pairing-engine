package gotha

import (
	"tournament_pair/src"
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

func NewTournament(players []Player) *Tournament {
	t := &Tournament{}
	// 添加报名人员
	t.hmPlayers = make(map[string]*Player)
	t.AddPlayer(players)

	t.hmGames = make(map[string]*Game)
	t.byePlayers = make([]*Player, src.MAX_NUMBER_OF_ROUNDS)

	return t
}

func (t *Tournament) GetTournamentSet() *parameter_set.TournamentParameterSet {
	return t.tournamentParameterSet
}

func (t *Tournament) AddPlayer(players []Player) {
	for _, player := range players {
		t.hmPlayers[player.SetKeyString()] = &player
	}
}

func (t *Tournament) FillPairingInfo(roundNumber int) {
	gps := parameter_set.NewGeneralParameterSet()
	pps := parameter_set.NewPlacementParameterSet()
	paiPs := parameter_set.NewPairingParameterSet()
	mainScoreMin := 0
	mainScoreMax := 0
	groupNumber := 0 // 有几轮就有几个groupNumber

	for cat := 0; cat < gps.NumberOfCategories; cat++ {
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
			crit := pps.GetPlaCriteria()
			additionalCrit := paiPs.PaiMaAdditionalPlacementCritSystem1
			// todo 种子

			paiCrit := make([]int, len(crit)+1)
			copy(crit, paiCrit)
			paiCrit[len(paiCrit)-1] = additionalCrit
			// todo sort very important

			// groupNumber 所记录的就是最后一次的获胜轮次
			//
			for index, sp := range alSPGroup {
				sp.GroupNumber = groupNumber
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

func (t *Tournament) MakeAutomaticPairing(alPlayersToPair []Player, roundNumber int) ([]Game, bool) {
	// not even
	if len(alPlayersToPair)%2 != 0 {
		return nil, false
	}

	t.fillBaseScoringInfoIfNecessary()

	// todo getGamesListBefore()
	//alPreviousGames := t.gamesListBefore(roundNumber)

	// fill pairing info
	t.FillPairingInfo(roundNumber)

	// todo 大于300人的比赛暂未处理

	// fillPairingInfo

	alGames := make([]Game, 0)
	alg := make([]Game, 0)

	// fill game
	copy(alGames, alg)

	return alGames, false
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

	// 2) nbwX2  (计算总分？)
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

		// 3) CUSSW 计算每一轮的总分之和
		for _, sp := range t.hmScoredPlayers {
			sp.SetCUSWX2(0, sp.GetNBWX2(0))
			for r := 1; r < numberOfRoundsToCompute; r++ {
				// 之前轮次的总分+ 本轮大分
				sp.SetCUSWX2(r, sp.GetCUSWX2(r-1)+sp.GetNBWX2(r))
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

		// 5) ssswX2 Sum of opponents sosw2 * 2 对手的SOSW总和
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

		// 6)  EXT EXR 额外字段？
		for r := 0; r < numberOfRoundsToCompute; r++ {
			for _, sp := range t.hmScoredPlayers {
				extX2 := 0
				exrX2 := 0
				for rr := 0; rr <= r; rr++ {
					// 未匹配状态下不用计算
					if sp.GetParticipation(rr) != PAIRED {
						continue
					}
					g := sp.GetGame(rr)
					opp := NewPlayer()
					spWasWhite := false // 是否执白
					if g.GetWhitePlayer().HasSameKeyString(sp.Player) {
						opp = g.GetBlackPlayer()
						spWasWhite = true
					} else {
						opp = g.GetWhitePlayer()
						spWasWhite = false
					}
					sOpp := t.hmScoredPlayers[opp.GetKeyString()]
					// 保存handicap
					realHd := g.GetHandicap()
					if !spWasWhite {
						realHd = -realHd
					}
					naturalHd := sp.GetRank() - sOpp.GetRank()
					coef := 0
					if realHd-naturalHd <= 0 {
						coef = 0
					}
					if realHd-naturalHd == 0 {
						coef = 1
					}
					if realHd-naturalHd == 1 {
						coef = 2
					}
					if realHd-naturalHd >= 2 {
						coef = 3
					}

					extX2 += sOpp.GetNBWX2(r) * coef
					bWin := false
					if spWasWhite && (g.GetResult() == RESULT_WHITEWINS || g.GetResult() == RESULT_WHITEWINS_BYDEF ||
						g.GetResult() == RESULT_BOTHWIN || g.GetResult() == RESULT_BOTHWIN_BYDEF) {
						bWin = true
					}
					if !spWasWhite && (g.GetResult() == RESULT_BLACKWINS || g.GetResult() == RESULT_BLACKWINS_BYDEF ||
						g.GetResult() == RESULT_BOTHWIN || g.GetResult() == RESULT_BOTHWIN_BYDEF) {
						bWin = true
					}
					if bWin {
						exrX2 += sOpp.GetNBWX2(r) * coef
					}
				}
				sp.SetEXTX2(r, extX2)
				sp.SetEXTX2(r, exrX2)
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
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

func (t *Tournament) pairAGroup(alGroupedPlayers []*Player, roundNumber int, hmScoredPlayers map[string]*ScoredPlayer, alPreviousGames []*Game) {
	// 分组数量
	numberOfPlayersInGroup := len(alGroupedPlayers)
	//      Prepare infos about Score groups : sgSize, sgNumber and innerPosition
	//      And DUDD information

	costs := make([][]int64, 0)
	// 生成一个二维数组，替代java中 long[numberOfPlayersInGroup][numberOfPlayersInGroup]
	for i := 0; i <= numberOfPlayersInGroup; i++ {
		costs = append(costs, make([]int64, numberOfPlayersInGroup))
	}

	for i := 0; i < numberOfPlayersInGroup; i++ {
		costs[i][i] = 0
		for j := i + 1; j < numberOfPlayersInGroup; j++ {
			//p1 := alGroupedPlayers[i]
			//p2 := alGroupedPlayers[j]
			//sP1 := t.hmScoredPlayers[p1.GetKeyString()]
			//sP2 := t.hmScoredPlayers[p2.GetKeyString()]

		}
	}
}

func (t *Tournament) costValue(sP1 *ScoredPlayer, sP2 *ScoredPlayer, roundNumber int, alPreviousGames []*Game) {
	//gps := t.tournamentParameterSet.GetGeneralParameterSet()
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
		// 如果之前从未匹配过，则会加上一个非常大的权重
		if numberOfPreviousGamesP1P2 == 0 {
			cost += paiPS.GetPaiBaAvoidDuplGame()
		}

		// 增加随机因子
		//Base Criterion 2 : Random
		var nR int64
		if paiPS.IsPaiBaDeterministic() {

		} else {

		}
		cost += nR

		// Base Criterion 3 : Balance W and B
		// This cost is never applied if potential Handicap != 0
		// It is fully applied if wbBalance(sP1) and wbBalance(sP2) are strictly of different signs
		// It is half applied if one of wbBalance is 0 and the other is >=2

		//var bwBalanceCost int64= 0

		//potHd:=
	}

}
