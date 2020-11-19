package gotha

import (
	"tournament_pair/src/parameter_set"
)

const (
	UNKNOWN      = 0
	ABSENT       = -3 // 缺席
	NOT_ASSIGNED = -2 // 未匹配
	BYE          = -1
	PAIRED       = 1
)

type ScoredPlayers []*ScoredPlayer

type ScoredPlayer struct {
	*Player
	NumberOfGroups int // 分组数量 Very redundant
	GroupNumber    int
	GroupSize      int // redundant
	InnerPlacement int // placement in homogeneous group (category and mainScore) beteen 0 and size(group) - 1
	NbDU           int // 上调次数
	NbDD           int // 下调次数
	/** generalParameterSet is a part of ScoredPlayer because mms is dependent on McMahon bars and floors */
	generalParameterSet *parameter_set.GeneralParameterSet

	/** for each round, participation[r] can be : ABSENT, NOT_ASSIGNED, BYE or PAIRED
	/** 记录每一轮的参赛情况
	*/
	participation []int
	gameArray     []*Game

	nbwX2    []int // number of wins * 2 获胜数量*2
	cuswX2   []int // Sum of successive nbw2
	soswX2   []int // Sum of Opponents nbw2
	soswM1X2 []int // Sum of (n-1) Opponents nbw2
	soswM2X2 []int // Sum of (n-2) Opponents nbw2
	sdswX4   []int // Sum of Defeated Opponents nbw2 X2 击败的对手分*2  *2
	extX2    []int // Exploits tentes (based on nbw2, with a weight factor)
	exrX2    []int // Exploits reussis(based on nbw2, with a weight factor)
	ssswX2   []int // Sum of opponents sosw2 * 2

	dc  int // Direct Confrontation
	sdc int // Simplified Direct Confrontation
}

func NewScoredPlayer(gps *parameter_set.GeneralParameterSet, player *Player) *ScoredPlayer {
	this := &ScoredPlayer{}
	this.Player = deepCopyPlayer(player)
	this.generalParameterSet = gps
	numberOfRounds := gps.GetNumberOfRounds()
	this.participation = make([]int, numberOfRounds)
	this.gameArray = make([]*Game, numberOfRounds)

	// First level scores
	this.nbwX2 = make([]int, numberOfRounds) // number of wins * 2 获胜数量*2 大分

	// Second level scores
	this.cuswX2 = make([]int, numberOfRounds) // Sum of successive nbw2 总分
	this.soswX2 = make([]int, numberOfRounds) // Sum of Opponents nbw2 对手
	this.soswM1X2 = make([]int, numberOfRounds)
	this.soswM2X2 = make([]int, numberOfRounds)
	this.sdswX4 = make([]int, numberOfRounds)

	this.extX2 = make([]int, numberOfRounds)
	this.exrX2 = make([]int, numberOfRounds)

	this.ssswX2 = make([]int, numberOfRounds)

	// dc and sdc are defined for the current round number
	this.dc = 0
	this.sdc = 0

	return this
}

// 获取标准数据
func (sp *ScoredPlayer) GetCritValue(criterion int, roundNumber int) int {
	switch criterion {
	case parameter_set.PLA_CRIT_NUL:
		return 0
	case parameter_set.PLA_CRIT_RATING:
		return sp.GetRating()
	case parameter_set.PLA_CRIT_NBW:
		if roundNumber >= 0 {
			return sp.nbwX2[roundNumber]
		} else {
			return 0
		}
	case parameter_set.PLA_CRIT_SOSW:
		if roundNumber >= 0 {
			return sp.soswX2[roundNumber]
		} else {
			return 0
		}
	case parameter_set.PLA_CRIT_SOSOSW:
		if roundNumber >= 0 {
			return sp.ssswX2[roundNumber]
		} else {
			return 0
		}

	default:
		return 0
	}
}

func (sp *ScoredPlayer) GetRating() int {
	return sp.Rating
}

func (sp *ScoredPlayer) SetParticipation(rn int, participation int) {
	if sp.isValidRoundNumber(rn) {
		sp.participation[rn] = participation
	} else {
		sp.participation[rn] = UNKNOWN
	}
}

func (sp *ScoredPlayer) isValidRoundNumber(rn int) bool {
	if rn < 0 || rn > len(sp.participation) {
		return false
	}
	return true
}

func (sp *ScoredPlayer) SetGame(rn int, g *Game) {
	if sp.isValidRoundNumber(rn) {
		sp.gameArray[rn] = g
	}
}

func (sp *ScoredPlayer) SetNBWX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.nbwX2[rn] = value
	}
}

func (sp *ScoredPlayer) GetNBWX2(rn int) int {
	if sp.isValidRoundNumber(rn) {
		return sp.nbwX2[rn]
	}
	return 0

}

func (sp *ScoredPlayer) GetParticipation(rn int) int {
	if sp.isValidRoundNumber(rn) {
		return sp.participation[rn]
	}
	return 0
}

func (sp *ScoredPlayer) SetCUSWX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.cuswX2[rn] = value
	}
}

func (sp *ScoredPlayer) GetCUSWX2(rn int) int {
	if sp.isValidRoundNumber(rn) {
		return sp.cuswX2[rn]
	}
	return 0
}

func (sp *ScoredPlayer) GetGame(rn int) *Game {
	if sp.isValidRoundNumber(rn) {
		return sp.gameArray[rn]
	}
	return nil
}

func (sp *ScoredPlayer) GetSOSWX2(rn int) int {
	if sp.isValidRoundNumber(rn) {
		return sp.soswX2[rn]
	}
	return 0
}

func (sp *ScoredPlayer) SetSSSWX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.ssswX2[rn] = value
	}
}

func (sp *ScoredPlayer) SetEXTX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.extX2[rn] = value
	}
}

func (sp *ScoredPlayer) SetEXRX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.exrX2[rn] = value
	}
}

func (sp *ScoredPlayer) SetSOSWX2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.soswX2[rn] = value
	}
}

func (sp *ScoredPlayer) SetSDSWX4(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.sdswX4[rn] = value
	}
}

func (sp *ScoredPlayer) SetSOSWM1X2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.soswM1X2[rn] = value
	}
}

func (sp *ScoredPlayer) SetSOSWM2X2(rn int, value int) {
	if sp.isValidRoundNumber(rn) {
		sp.soswM2X2[rn] = value
	}
}
