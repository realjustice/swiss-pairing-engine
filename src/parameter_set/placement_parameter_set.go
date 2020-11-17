package parameter_set

// 人员配置参数
type PlacementParameterSet struct {
	plaCriteria []int
}

const (
	PLA_CRIT_NUL               = 0   // Null criterion
	PLA_MAX_NUMBER_OF_CRITERIA = 6   // 最大的规则数量
	PLA_CRIT_NBW               = 100 // Number of Wins (swiss)
	PLA_CRIT_SOSW              = 110 // Sum of Opponents NbW
	PLA_CRIT_SOSOSW            = 130 // Sum of opponents SOS
	PLA_CRIT_RATING            = 12  // Rating
)

func NewPlacementParameterSet() *PlacementParameterSet {
	return &PlacementParameterSet{}
}

func (set *PlacementParameterSet) GetPlaCriteria() []int {
	plaC := make([]int, len(set.plaCriteria))
	// plaCriteria 赋值
	copy(plaC, set.plaCriteria)
	return plaC
}

// 初始swiss规则
func (set *PlacementParameterSet) InitForSwiss() {
	set.plaCriteria = make([]int, PLA_MAX_NUMBER_OF_CRITERIA)
	set.plaCriteria[0] = PLA_CRIT_NBW
	set.plaCriteria[1] = PLA_CRIT_SOSW
	set.plaCriteria[2] = PLA_CRIT_SOSOSW
	set.plaCriteria[3] = PLA_CRIT_NUL
	set.plaCriteria[4] = PLA_CRIT_NUL
	set.plaCriteria[5] = PLA_CRIT_NUL
}
