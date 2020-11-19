package gotha

import (
	"encoding/xml"
	"strings"
	"tournament_pair/src/parameter_set"
)

type Player struct {
	Player           xml.Name `xml:"Player"`
	Name             string   `xml:"name,attr"`
	Rating           int      `xml:"rating,attr"`
	Rank             int
	ParticipatingStr string `xml:"participatingStr,attr"`
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

func (p *Player) SetKeyString() string {
	p.keyString = strings.ToUpper(p.Name)
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