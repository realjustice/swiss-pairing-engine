package gotha

import (
	"github.com/realjustice/swiss-pairing-engine/src/parameter_set"
	"strings"
)

type Player struct {
	Name             string
	FirstName        string
	Rating           int
	Rank             int
	ParticipatingStr string
	keyString        string
	// 记录每一轮是否参与编排
	participating []bool
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Category(gps *parameter_set.GeneralParameterSet) int {
	return 0
}

func (p *Player) SetKeyString(keyStr string) string {
	if keyStr == "" {
		p.keyString = strings.ToUpper(p.Name + p.FirstName)
	}
	p.keyString = keyStr
	return p.keyString
}

func deepCopyPlayer(player *Player) *Player {
	newPlayer := *player
	return &newPlayer
}

func (p *Player) GetParticipating() []bool {
	copyParticipating := make([]bool, len(p.participating))
	copy(copyParticipating, p.participating)
	return copyParticipating
}

/**
 * 2 players never have the same key string.
 * hasSameKeyString is, thus a way to test if 2 references refer to the same player
 **/
func (p *Player) HasSameKeyString(player *Player) bool {
	if player == nil {
		return false
	}
	if p.keyString == player.keyString {
		return true
	}

	return false
}

func (p *Player) GetKeyString() string {
	return p.keyString
}

func (p *Player) SetRank(rating int) {
	rank := rankFromRating(rating)
	p.Rank = rank
}

func (p *Player) GetRank() int {
	return p.Rank
}

func (p *Player) SetParticipating(val []bool) {
	newVal := make([]bool, len(val))
	copy(newVal, val)
	p.participating = newVal
}

func (p *Player) SetParticipatingStr(val string) {
	p.ParticipatingStr = val
	participating := make([]bool, MAX_NUMBER_OF_ROUNDS)
	for i := 0; i < len(p.ParticipatingStr); i++ {
		if string(p.ParticipatingStr[i]) == "0" {
			participating[i] = false
		} else {
			participating[i] = true
		}
	}
	// 设置本轮是否参与编排
	p.SetParticipating(participating)
}

func (p *Player) SetFirstName(firstName string) {
	p.FirstName = firstName
}

func (p *Player) SetName(name string) {
	p.Name = name
}
