package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

func NewMetric(curr, prev int64) params.Metric {
	return params.Metric{
		Total: curr,
		Diff:  curr - prev,
		Ratio: CalcRatio(curr, prev),
	}
}

func CalcRatio(curr, prev int64) int64 {
	if prev == 0 {
		if curr > 0 {
			return 100
		}
		return 0
	}
	return (curr - prev) * 100 / prev
}

// Hitung ACOS (%)
func CalcACOS(ads, revenue int64) float64 {
	if revenue == 0 {
		return 0
	}
	return (float64(ads) / float64(revenue)) * 100
}

// Hitung ROAS (rasio)
func CalcROAS(ads, revenue int64) float64 {
	if ads == 0 {
		return 0
	}
	return float64(revenue) / float64(ads)
}
