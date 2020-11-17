package parameter_set

// 通用配置参数
type GeneralParameterSet struct {
	NumberOfCategories int
	numberOfRounds     int
	genNBW2ValueAbsent int
	genNBW2ValueBye    int
}

func NewGeneralParameterSet() *GeneralParameterSet {
	set := new(GeneralParameterSet)
	// 设置默认值
	set.NumberOfCategories = 1 // 默认分类数量
	set.numberOfRounds = 5     // 默认轮次数量
	set.genNBW2ValueBye = 2    // 默认轮空积分
	return set
}

func (p *GeneralParameterSet) SetNumberOfRounds(roundNumber int) {
	p.numberOfRounds = roundNumber
}

func (p *GeneralParameterSet) GetNumberOfRounds() int {
	return p.numberOfRounds
}

func (p *GeneralParameterSet) GetGenNBW2ValueAbsent() int {
	return p.genNBW2ValueAbsent
}

func (p *GeneralParameterSet) GetGenNBW2ValueBye() int {
	return p.genNBW2ValueBye
}
