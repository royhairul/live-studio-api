package service

import "github.com/royhairul/live-studio-api/internal/domains/dashboard/params"

func NewMetric(curr, prev int64) params.Metric {
	return params.Metric{
		Total: curr,
		Diff:  curr - prev,
		Ratio: calcRatio(curr, prev),
	}
}

func calcRatio(curr, prev int64) int64 {
	if prev == 0 {
		if curr > 0 {
			return 100
		}
		return 0
	}
	return (curr - prev) * 100 / prev
}
