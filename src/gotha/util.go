package gotha

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func rankFromRating(rating int) int {
	rk := (rating+950)/100 - 30
	if rk > 8 {
		rk = 8
	}
	if rk < -30 {
		rk = -30
	}
	return rk
}
