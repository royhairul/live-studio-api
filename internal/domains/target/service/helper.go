package service

func CalcRatio(real, target int64) float64 {
	if target == 0 {
		return 0
	}
	return (float64(real) / float64(target)) * 100
}
