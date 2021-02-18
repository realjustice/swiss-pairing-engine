package gotha

import (
	"fmt"
	weighted_match_long "github.com/realjustice/maximum_weight_matching/src"
	"github.com/realjustice/swiss-pairing-engine/src/parameter_set"
	"github.com/realjustice/swiss-pairing-engine/src/weight"
	"math"
	"sort"
	"strings"
)

type Tournament struct {
	tournamentParameterSet *parameter_set.TournamentParameterSet
	selectedPlayers        []*Player
	/**
	 * HashMap of Players The key is the getKeyString
	 */
	hmPlayers map[string]*Player
	/**
	 * HashMap of Games The key is (roundNumber * Gotha.MAX_NUMBER_OF_TABLES +
	 * tableNumber)
	 */
	hmGames map[int]*Game

	hmScoredPlayers map[string]*ScoredPlayer

	byePlayers []*Player
}

type GameIterator struct {
	data  []*Game
	index int
}

func (ti *GameIterator) HasNext() bool {
	if ti.data == nil {
		return false
	}
	return ti.index < len(ti.data)
}

func (ti *GameIterator) Walk(f func(g *Game) (isStop bool)) {
	if ti == nil {
		return
	}

	for ti.HasNext() {
		isStop := f(ti.data[ti.index])
		if isStop {
			break
		}
		ti.index++
	}
}

func NewGameIterator(data []*Game) *GameIterator {
	return &GameIterator{data: data}
}

func NewTournament() *Tournament {
	t := new(Tournament)
	t.tournamentParameterSet = parameter_set.NewTournamentParameterSet()
	// 添加报名人员
	t.hmPlayers = make(map[string]*Player)

	t.hmGames = make(map[int]*Game)
	t.byePlayers = make([]*Player, MAX_NUMBER_OF_ROUNDS)

	return t
}

func (t *Tournament) GetTournamentSet() *parameter_set.TournamentParameterSet {
	return t.tournamentParameterSet
}
func (t *Tournament) SetTournamentSet(set *parameter_set.TournamentParameterSet) {
	t.tournamentParameterSet = set
}

func (t *Tournament) AddPlayer(p *Player) {
	if p.keyString == "" {
		p.keyString = strings.ToUpper(p.Name + p.FirstName)
	}

	t.hmPlayers[p.GetKeyString()] = p
	//t.selectedPlayers = append(t.selectedPlayers, p)
}

func (t *Tournament) AddGame(g *Game) bool {
	if g == nil {
		return false
	}
	wp := g.GetWhitePlayer()
	bp := g.GetBlackPlayer()
	if wp == nil || bp == nil {
		return false
	}
	r := g.GetRoundNumber()
	tt := g.GetTableNumber()

	key := r*MAX_NUMBER_OF_TABLES + tt
	t.hmGames[key] = g
	return true
}

func (t *Tournament) FillPairingInfo(roundNumber int) {
	gps := t.tournamentParameterSet.GetGeneralParameterSet()
	pps := t.tournamentParameterSet.GetPlacementParameterSet()
	paiPs := t.tournamentParameterSet.GetPairingParameterSet()
	mainCrit := pps.MainCriterion()
	mainScoreMin := 0
	mainScoreMax := roundNumber
	groupNumber := 0 // group 的数量（比如第四轮结束一共有大分 8分，6分，4分，2分，0分的五类选手，则分为5个小组）

	for cat := 0; cat < gps.GetNumberOfCategories(); cat++ {
		for mainScore := mainScoreMax; mainScore >= mainScoreMin; mainScore-- {
			alSPGroup := make([]*ScoredPlayer, 0)
			for _, sp := range t.hmScoredPlayers {
				if sp.Category(gps) != cat {
					continue
				}
				// 获取到本轮为止的全胜人员
				// 第一轮为全体选手
				if sp.GetCritValue(mainCrit, roundNumber-1)/2 != mainScore {
					continue
				}

				alSPGroup = append(alSPGroup, sp)
			}
			if len(alSPGroup) <= 0 {
				continue
			}
			// 压入一个规则，（会在规则数组的尾部压入rating规则，用来保证选手的排名是有序的）
			crit := pps.GetPlaCriteria()
			additionalCrit := paiPs.GetPaiMaAdditionalPlacementCritSystem1()
			if roundNumber > paiPs.GetPaiMaLastRoundForSeedSystem1() {
				additionalCrit = paiPs.GetPaiMaAdditionalPlacementCritSystem2()
			}

			paiCrit := make([]int, len(crit)+1)
			copy(paiCrit, crit)
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

	//ps := make(map[string]*ScoredPlayer)
	//readFile, err := os.Open("gob")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//dec := gob.NewDecoder(readFile)
	//err2 := dec.Decode(&ps)
	//
	//if err2 != nil {
	//	fmt.Println(err2)
	//	return
	//}
	//
	//
	//writeFile, err := os.Create("gob")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//enc := gob.NewEncoder(writeFile)
	//err2 = enc.Encode(t.hmScoredPlayers)
	//fmt.Println(err2)
	//a := make([]string, len(t.hmScoredPlayers))
	//for k, v := range t.hmScoredPlayers {
	//	a[t.hmScoredPlayers[k].InnerPlacement] = v.FirstName + v.FirstName
	//}
	//for index, v := range a {
	//	fmt.Printf("%s序号:%d\n", v, index)
	//}

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
		alTempScoredPlayers := make([]*ScoredPlayer, 0)
		for _, p := range t.hmScoredPlayers {
			alTempScoredPlayers = append(alTempScoredPlayers, p)
		}

		nbP := len(alTempScoredPlayers)
		scoreBefore := make([][]int, 0)
		for r := 0; r < roundNumber; r++ {
			scoreBefore = append(scoreBefore, make([]int, nbP))
		}

		for r := 0; r < roundNumber; r++ {
			for iSP := 0; iSP < nbP; iSP++ {
				sp := alTempScoredPlayers[iSP]
				scoreBefore[r][iSP] = sp.GetCritValue(mainCrit, r-1) / 2
			}
		}

		for r := 0; r < roundNumber; r++ {
			for iSP := 0; iSP < nbP; iSP++ {
				sp := alTempScoredPlayers[iSP]
				g := sp.GetGame(r)
				if g == nil {
					continue
				}
				wP := g.GetWhitePlayer()
				bP := g.GetBlackPlayer()
				var opp *Player
				if sp.HasSameKeyString(wP) {
					opp = bP
				} else {
					opp = wP
				}
				sOpp := t.hmScoredPlayers[opp.GetKeyString()]

				f := func() int {
					for i, v := range alTempScoredPlayers {
						if v == sOpp {
							return i
						}
					}
					return math.MinInt64
				}

				iSOpp := f()
				if scoreBefore[r][iSP] < scoreBefore[r][iSOpp] {
					sp.NbDU++
				}
				if scoreBefore[r][iSP] > scoreBefore[r][iSOpp] {
					sp.NbDD++
				}
			}
		}
	}
}

func (t *Tournament) getScoredPlayers() ScoredPlayers {
	sps := make(ScoredPlayers, len(t.hmScoredPlayers))
	i := 0
	for _, sp := range t.hmScoredPlayers {
		sps[i] = sp
		i++
	}
	return sps
}

func (t *Tournament) MakeAutomaticPairing(roundNumber int) ([]*Game, bool) {
	alPlayersToPair := t.selectedPlayers
	// not even
	if len(alPlayersToPair)%2 != 0 {
		return nil, false
	}

	gps := t.tournamentParameterSet.GetGeneralParameterSet()
	pps := t.tournamentParameterSet.GetPlacementParameterSet()

	t.fillBaseScoringInfoIfNecessary()

	// fill pairing info
	t.FillPairingInfo(roundNumber)
	// todo 当前只支持swiss
	mainCrit := pps.MainCriterion()
	// todo getGamesListBefore()
	alPreviousGames := t.gamesListBefore(roundNumber)

	mainScoreMin := 0
	mainScoreMax := roundNumber

	alRemainingPlayers := make([]*Player, len(alPlayersToPair))
	copy(alRemainingPlayers, alPlayersToPair)
	var alg []*Game
	alGames := make([]*Game, 0)
	// todo 大于300人的比赛暂未处理
	for len(alRemainingPlayers) > PAIRING_GROUP_MAX_SIZE {
		bGroupReady := false

		alGroupedPlayers := struct {
			data []*Player
		}{}
		alGroupedPlayers.data = make([]*Player, 0)

		for cat := 0; cat < gps.GetNumberOfCategories(); cat++ {
			for mainScore := mainScoreMax; mainScore >= mainScoreMin; mainScore-- {
				for it := NewPlayerIterator(&alRemainingPlayers); it.HasNext(); {
					p := it.Next()
					if p.Category(gps) > cat {
						continue
					}
					sp := t.hmScoredPlayers[p.GetKeyString()]
					if sp.GetCritValue(mainCrit, roundNumber-1)/2 < mainScore {
						continue
					}
					alGroupedPlayers.data = append(alGroupedPlayers.data, p)
					it.Remove()
					// 2 Emergency breaks
					if len(alGroupedPlayers.data) >= PAIRING_GROUP_MAX_SIZE {
						bGroupReady = true
						break
					}
					if len(alRemainingPlayers) <= PAIRING_GROUP_MIN_SIZE {
						bGroupReady = true
						break
					}
				}
				// Is the group ready for pairing ?
				if len(alGroupedPlayers.data) >= PAIRING_GROUP_MIN_SIZE && len(alGroupedPlayers.data)%2 == 0 {
					bGroupReady = true
				}
				if bGroupReady {
					break
				}
			}
			if bGroupReady {
				break
			}
		}
		alg = t.pairAGroup(alGroupedPlayers.data, roundNumber, alPreviousGames)
		alGames = append(alGames, alg...)
	}

	t.FillPairingInfo(roundNumber)

	alg = t.pairAGroup(alRemainingPlayers, roundNumber, alPreviousGames)

	// fill game
	alGames = append(alGames, alg...)

	return alGames, true
}

func (t *Tournament) Pair(roundNumber int) *GameIterator {
	if t.selectedPlayers == nil {
		t.selectedPlayers = getSelectedPlayerFromSPs(t.orderScoredPlayersList(roundNumber))
	}
	roundNumber--
	if len(t.selectedPlayers)%2 != 0 {
		// set bye player
		t.ChooseAByePlayer(t.selectedPlayers, roundNumber)
		// remove bye player from alPlayersToPair
		byeP := t.getByePlayer(roundNumber)
		var pToRemove *Player
		for _, p := range t.selectedPlayers {
			if p.HasSameKeyString(byeP) {
				pToRemove = p
			}
		}
		t.removePlayer(pToRemove)
	}

	alNewGames, isSucceed := t.MakeAutomaticPairing(roundNumber)

	if !isSucceed {
		return nil
	}
	tN := 0
	for _, g := range alNewGames {
		stop := true

		f := func() {
			stop = true
			oldGames := t.gamesList(roundNumber)
			for _, oldG := range oldGames {
				if oldG.GetRoundNumber() != roundNumber {
					continue
				}
				if oldG.GetTableNumber() == tN {
					tN++
					stop = false
				}
			}
		}
		f()
		for !stop {
			f()
		}
		tN++
		g.SetTableNumber(tN)
	}
	for _, g := range alNewGames {
		t.AddGame(g)
	}
	return NewGameIterator(alNewGames)
}

func (t *Tournament) BergerArrange() *GameIterator {
	games := make([]*Game, 0)
	var lastPlayer *Player
	n := len(t.selectedPlayers)
	moves := (n+n%2-4)/2 + 1
	round := n
	if n%2 == 0 {
		round--
	}
	if n%2 == 0 {
		lastPlayer = t.selectedPlayers[n-1]
		t.selectedPlayers = t.selectedPlayers[:n-1]
		n--
	}

	head, tail := 0, n

	ringBuffer := func(nums []*Player, head int, tail int, n int) []*Player {
		players := make([]*Player, 0)
		for head != tail {
			players = append(players, nums[head%n])
			head++
		}
		return players
	}

	for i := 1; i <= round; i++ {
		newPlayers := ringBuffer(t.selectedPlayers, head, tail, n)
		newPlayers = append(newPlayers, lastPlayer)
		if round%2 == 0 {
			// swap
			newPlayers[0], newPlayers[len(newPlayers)-1] = newPlayers[len(newPlayers)-1], newPlayers[0]
		}

		// two point
		l, r := 0, len(newPlayers)-1
		tableNumber := 0
		for l < r {
			tableNumber++
			if newPlayers[l] == nil {
				t.SetByePlayer(i, newPlayers[r].GetKeyString())
				l++
				r--
				continue
			}
			if newPlayers[r] == nil {
				t.SetByePlayer(i, newPlayers[l].GetKeyString())
				l++
				r--
				continue
			}

			game := &Game{RoundNumber: i, TableNumber: tableNumber, KnownColor: true, Handicap: 0, result: UNKNOWN, blackPlayer: newPlayers[l], whitePlayer: newPlayers[r]}
			games = append(games, game)
			l++
			r--
		}

		head += n - moves
		tail += n - moves
	}
	return NewGameIterator(games)
}

func (t *Tournament) orderScoredPlayersList(roundNumber int) ScoredPlayers {
	roundNumber--
	t.fillBaseScoringInfoIfNecessary()
	crit := t.tournamentParameterSet.GetPlacementParameterSet().GetPlaCriteria()
	primaryCrit := make([]int, len(crit))
	for iC := 0; iC < len(crit); iC++ {
		primaryCrit[iC] = crit[iC]
	}
	alOrderedScoredPlayers := t.getScoredPlayers()
	spc := NewScoredPlayerComparator(alOrderedScoredPlayers, primaryCrit, roundNumber)
	sort.Sort(spc)
	return alOrderedScoredPlayers
}

func getSelectedPlayerFromSPs(sps ScoredPlayers) []*Player {
	players := make([]*Player, len(sps))
	for i, p := range sps {
		players[i] = p.Player
	}
	return players
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
	}
	for _, sp := range t.hmScoredPlayers {
		// 弃权或者轮空
		nbPtsNBW2AbsentOrBye := 0
		for r := 0; r < numberOfRoundsToCompute; r++ {
			if sp.GetParticipation(r) == ABSENT { // 弃权
				nbPtsNBW2AbsentOrBye += gps.GetGenNBW2ValueAbsent()
			}
			if sp.GetParticipation(r) == BYE {
				nbPtsNBW2AbsentOrBye += gps.GetGenNBW2ValueBye()
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

func (t *Tournament) gamesList(roundNumber int) []*Game {
	gL := make([]*Game, 0)
	for _, g := range t.hmGames {
		if g.GetRoundNumber() == roundNumber {
			gL = append(gL, g)
		}
	}
	return gL
}

func (t *Tournament) pairAGroup(alGroupedPlayers []*Player, roundNumber int, alPreviousGames []*Game) []*Game {
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
	w := weight.NewWeightedMatchLong()
	var total int64
	for _, v1 := range costs {
		for _, v2 := range v1 {
			total += v2
		}
	}
	fmt.Println(total)

	mate := w.WeightedMatchLong(costs, weighted_match_long.MAXIMIZE)
	alG := make([]*Game, 0)

	for i := 1; i <= len(costs); i++ {
		if i < mate[i] {
			p1 := alGroupedPlayers[i-1]
			p2 := alGroupedPlayers[mate[i]-1]
			sP1 := t.hmScoredPlayers[p1.GetKeyString()]
			sP2 := t.hmScoredPlayers[p2.GetKeyString()]
			g := t.gameBetween(sP1, sP2, roundNumber)
			alG = append(alG, g)
		}
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
	// 假设现在有三个group （group1 ，group2 ，group3 其中group1总分<group2<group3），我在group2这个分组中
	// 那么我最容易遇到的是group1中"上调次数最多,且上调次数大于下调次数的选手"
	// 以及group3中"下调次数最多,且下调次数大于上调次数的选手"
	var duddCost int64
	if abs(sP1.GroupNumber-sP2.GroupNumber) < 4 && sP1.GroupNumber != sP2.GroupNumber {
		// 4 scenarii
		// scenario = 0 : One of the players has already been drawn in the same sense
		// scenario = 1 : normal conditions (does not correct anything and no previous drawn in the same sense)
		// scenario = 2 : it corrects a previous DU/DD
		// scenario = 3 : it corrects a previous DU/DD for both
		scenario := 1
		if sP1.NbDU > 0 && sP1.GroupNumber > sP2.GroupNumber {
			scenario = 0
		}
		if sP1.NbDD > 0 && sP1.GroupNumber < sP2.GroupNumber {
			scenario = 0
		}
		if sP2.NbDU > 0 && sP2.GroupNumber > sP1.GroupNumber {
			scenario = 0
		}
		if sP2.NbDD > 0 && sP2.GroupNumber < sP1.GroupNumber {
			scenario = 0
		}

		if scenario != 0 && sP1.NbDU > 0 && sP1.GroupNumber < sP2.GroupNumber {
			scenario++
		}
		if scenario != 0 && sP1.NbDD > 0 && sP1.GroupNumber > sP2.GroupNumber {
			scenario++
		}
		if scenario != 0 && sP2.NbDU > 0 && sP2.NbDD < sP2.NbDU && sP2.GroupNumber < sP1.GroupNumber {
			scenario++
		}
		if scenario != 0 && sP2.NbDD > 0 && sP2.NbDU < sP2.NbDD && sP2.GroupNumber > sP1.GroupNumber {
			scenario++
		}

		// Modifs V3.33.04

		duddWeight := paiPS.GetPaiMaDUDDWeight() / 4
		upperSP := func() *ScoredPlayer {
			if sP1.GroupNumber < sP2.GroupNumber {
				return sP1
			} else {
				return sP2
			}
		}()

		lowerSP := func() *ScoredPlayer {
			if sP1.GroupNumber < sP2.GroupNumber {
				return sP2
			} else {
				return sP1
			}
		}()

		if paiPS.GetPaiMaDUDDUpperMode() == parameter_set.PAIMA_DUDD_TOP {
			duddCost += duddWeight / 2 * int64(upperSP.GroupSize-1-upperSP.InnerPlacement) / int64(upperSP.GroupSize)
		} else if paiPS.GetPaiMaDUDDUpperMode() == parameter_set.PAIMA_DUDD_MID {
			duddCost += duddWeight / 2 * int64(upperSP.GroupSize-1-abs(2*upperSP.InnerPlacement-upperSP.GroupSize+1)) / int64(upperSP.GroupSize)
		} else if paiPS.GetPaiMaDUDDUpperMode() == parameter_set.PAIMA_DUDD_BOT {
			duddCost += duddWeight / 2 * int64(upperSP.InnerPlacement) / int64(upperSP.GroupSize)
		}
		if paiPS.GetPaiMaDUDDLowerMode() == parameter_set.PAIMA_DUDD_TOP {
			duddCost += duddWeight / 2 * int64(lowerSP.GroupSize-1-lowerSP.InnerPlacement) / int64(lowerSP.GroupSize)
		} else if paiPS.GetPaiMaDUDDLowerMode() == parameter_set.PAIMA_DUDD_MID {
			duddCost += duddWeight / 2 * int64(lowerSP.GroupSize-1-abs(2*lowerSP.InnerPlacement-lowerSP.GroupSize+1)) / int64(lowerSP.GroupSize)
		} else if paiPS.GetPaiMaDUDDLowerMode() == parameter_set.PAIMA_DUDD_BOT {
			duddCost += duddWeight / 2 * int64(lowerSP.InnerPlacement) / int64(lowerSP.GroupSize)
		}

		if scenario == 1 || (scenario >= 1 && !paiPS.IsPaiMaCompensateDUDD()) {
			duddCost += duddWeight
		} else if scenario == 2 {
			duddCost += 2 * duddWeight
		} else if scenario == 3 {
			duddCost += 3 * duddWeight
		}
	}
	cost += duddCost
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

func (t *Tournament) GetGamesFromRound(rn int) []*Game {
	games := make([]*Game, 0)
	for _, g := range t.hmGames {
		if g.RoundNumber == rn {
			games = append(games, g)
		}
	}

	return games
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
		} else if wbBalance(sP1, roundNumber-1) < wbBalance(sP2, roundNumber-1) {
			g.SetWhitePlayer(p1)
			g.SetBlackPlayer(p2)
		} else {
			//if getRand(){
			//	g.SetWhitePlayer(p2)
			//	g.SetBlackPlayer(p1)
			//}else {
			//	g.SetWhitePlayer(p1)
			//	g.SetBlackPlayer(p2)
			//}

			g.SetWhitePlayer(p1)
			g.SetBlackPlayer(p2)
		}
	}

	g.SetKnownColor(true)
	g.SetResult(RESULT_UNKNOWN)
	g.setRoundNumber(roundNumber)

	return g
}

func (t *Tournament) GetPlayerByKeyString(strNaFi string) *Player {
	return t.hmPlayers[strNaFi]
}

func (t *Tournament) SetSelectedPlayers(strNaFis []string) {
	for _, keyString := range strNaFis {
		t.selectedPlayers = append(t.selectedPlayers, t.GetPlayerByKeyString(keyString))
	}
}

func (t *Tournament) getGameByRoundAndTable(rn int, tableNumber int) *Game {
	key := rn*MAX_NUMBER_OF_TABLES + tableNumber
	return t.hmGames[key]
}

func (t *Tournament) SortGameByTableNumber() []*Game {
	games := make([]*Game, 0)
	for _, g := range t.hmGames {
		games = append(games, g)
	}

	// 冒泡排序
	n := len(games)
	hasChanged := false
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if games[j].TableNumber > games[j+1].TableNumber {
				games[j], games[j+1] = games[j+1], games[j]
				hasChanged = true
			}
		}
		if !hasChanged {
			break
		}
	}

	hasChanged = false
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if games[j].RoundNumber > games[j+1].RoundNumber {
				games[j], games[j+1] = games[j+1], games[j]
				hasChanged = true
			}
		}
		if !hasChanged {
			break
		}
	}

	return games
}

func (t *Tournament) SortGameByTableNumberFromRn(rn int) []*Game {
	rn--
	allGames := t.SortGameByTableNumber()
	games := make([]*Game, 0)
	for _, v := range allGames {
		if v.GetRoundNumber() == rn {
			games = append(games, v)
		}

	}
	return games
}

func (t *Tournament) SetGameResult(rn int, tableNumber int, result string) {
	rn--
	game := t.getGameByRoundAndTable(rn, tableNumber)

	game.SetResult(SelectResult(result))
}

func (t *Tournament) getByePlayer(roundNumber int) *Player {
	if t.byePlayers == nil {
		t.byePlayers = make([]*Player, MAX_NUMBER_OF_ROUNDS)
	}
	return t.byePlayers[roundNumber]
}

func (t *Tournament) SetByePlayer(roundNumber int, keyString string) {
	roundNumber--
	t.setByePlayer(roundNumber, keyString)
}

func (t *Tournament) setByePlayer(roundNumber int, keyString string) {
	if t.byePlayers == nil {
		t.byePlayers = make([]*Player, MAX_NUMBER_OF_ROUNDS)
	}
	t.byePlayers[roundNumber] = t.GetPlayerByKeyString(keyString)
}

func (t *Tournament) GetByePlayer(roundNumber int) *Player {
	roundNumber--
	return t.getByePlayer(roundNumber)
}

func (t *Tournament) ChooseAByePlayer(alPlayers []*Player, roundNumber int) {
	// The weight allocated to each player is 1000 * number of previous byes + rank
	// The chosen player will be the player with the minimum weight
	var bestPlayerForBye *Player
	minWeight := 1000*(MAX_NUMBER_OF_ROUNDS-1) + 38 + 1 // Nobody can have such a weight neither more

	for _, p := range alPlayers {
		weightForBye := p.GetRank()
		for r := 0; r < roundNumber; r++ {
			if t.byePlayers[r] == nil {
				continue
			}
			if p.HasSameKeyString(t.byePlayers[r]) {
				weightForBye += 1000
			}
		}
		if weightForBye < minWeight {
			minWeight = weightForBye
			bestPlayerForBye = p
		}
	}

	t.byePlayers[roundNumber] = bestPlayerForBye
}

func (t *Tournament) removePlayer(rmP *Player) {
	for i, p := range t.selectedPlayers {
		if p == rmP {
			t.selectedPlayers = append(t.selectedPlayers[:i], t.selectedPlayers[i+1:]...)
			break
		}
	}
}
