package gotha

const (
	RESULT_UNKNOWN         = 0
	RESULT_BYDEF           = 2 << 8 // 100000000
	RESULT_LOSE            = 2 << 3
	RESULT_WHITEWINS       = 17
	RESULT_WHITEWINS_BYDEF = RESULT_WHITEWINS + RESULT_BYDEF
	RESULT_BLACKWINS       = 18
	RESULT_BLACKWINS_BYDEF = RESULT_BLACKWINS + RESULT_BYDEF
	RESULT_EQUAL           = 19
	RESULT_EQUAL_BYDEF     = RESULT_EQUAL + RESULT_BYDEF
	RESULT_BOTHLOSE        = RESULT_LOSE + RESULT_LOSE
	RESULT_BOTHLOSE_BYDEF  = RESULT_BOTHLOSE + RESULT_BYDEF
	RESULT_BOTHWIN         = 35
	RESULT_BOTHWIN_BYDEF   = RESULT_BOTHWIN + RESULT_BYDEF
)

type Game struct {
	RoundNumber int // begin from 0
	TableNumber int // begin from 0
	KnownColor  bool
	Handicap    int
	result      int
	blackPlayer *Player
	whitePlayer *Player
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) GetBlackPlayer() *Player {
	return g.blackPlayer
}

func (g *Game) GetWhitePlayer() *Player {
	return g.whitePlayer
}

func (g *Game) GetResult() int {
	return g.result
}

func (g *Game) GetHandicap() int {
	return g.Handicap
}

func (g *Game) GetRoundNumber() int {
	return g.RoundNumber
}

func (g *Game) GetTableNumber() int {
	return g.TableNumber
}

func (g *Game) SetKnownColor(knownColor bool) {
	g.KnownColor = knownColor
}

func (g *Game) SetResult(result int) {
	g.result = result
}

func (g *Game) SetRoundNumber(roundNumber int) {
	g.RoundNumber = roundNumber
}

func (g *Game) SetTableNumber(tableNumber int) {
	g.TableNumber = tableNumber
}

func (g *Game) SetWhitePlayer(player *Player) {
	g.whitePlayer = player
}

func (g *Game) SetBlackPlayer(player *Player) {
	g.blackPlayer = player
}

func (g *Game) SetHandicap(val int) {
	g.Handicap = val
	if g.Handicap < 0 {
		g.Handicap = 0
	}
	if g.Handicap > 9 {
		g.Handicap = 9
	}
}

func SelectResult(resultDesc string) int {
	switch resultDesc {
	case "RESULT_WHITEWINS":
		return RESULT_WHITEWINS
	case "RESULT_BLACKWINS":
		return RESULT_BLACKWINS
	case "RESULT_EQUAL":
		return RESULT_EQUAL
	case "RESULT_BOTHLOSE":
		return RESULT_BOTHLOSE
	case "RESULT_BOTHWIN":
		return RESULT_BOTHWIN
	case "RESULT_WHITEWINS_BYDEF":
		return RESULT_WHITEWINS_BYDEF
	case "RESULT_BLACKWINS_BYDEF":
		return RESULT_BLACKWINS_BYDEF
	case "RESULT_EQUAL_BYDEF":
		return RESULT_EQUAL_BYDEF
	case "RESULT_BOTHLOSE_BYDEF":
		return RESULT_BOTHLOSE_BYDEF
	case "RESULT_BOTHWIN_BYDEF":
		return RESULT_BOTHWIN_BYDEF
	default:
		return RESULT_UNKNOWN
	}
}

func ConvertResult(result int) string {
	switch result {
	case RESULT_WHITEWINS:
		return "RESULT_WHITEWINS"
	case RESULT_BLACKWINS:
		return "RESULT_BLACKWINS"
	case RESULT_EQUAL:
		return "RESULT_EQUAL"
	case RESULT_BOTHLOSE:
		return "RESULT_BOTHLOSE"
	case RESULT_BOTHWIN:
		return "RESULT_BOTHWIN"
	case RESULT_WHITEWINS_BYDEF:
		return "RESULT_WHITEWINS_BYDEF"
	case RESULT_BLACKWINS_BYDEF:
		return "RESULT_BLACKWINS_BYDEF"
	case RESULT_EQUAL_BYDEF:
		return "RESULT_EQUAL_BYDEF"
	case RESULT_BOTHLOSE_BYDEF:
		return "RESULT_BOTHLOSE_BYDEF"
	case RESULT_BOTHWIN_BYDEF:
		return "RESULT_BOTHWIN_BYDEF"
	default:
		return "RESULT_UNKNOWN"
	}
}
