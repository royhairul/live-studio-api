package service

import (
	"context"
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/dashboard/params"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"

	performaparams "github.com/royhairul/live-studio-api/internal/domains/performa/params"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"

	performaagg "github.com/royhairul/live-studio-api/internal/aggregator/performa"
)

type DashboardServiceImpl struct {
	// TODO: add repository dependency
	hostSvc     hostservice.HostService
	accountSvc  accountservice.AccountService
	studioSvc   studioservice.StudioService
	performaAgg performaagg.PerformaAggregator
}

func NewDashboardService(
	hostSvc hostservice.HostService,
	accountSvc accountservice.AccountService,
	studioSvc studioservice.StudioService,
	performaAgg performaagg.PerformaAggregator,
) DashboardService {
	return &DashboardServiceImpl{
		hostSvc,
		accountSvc,
		studioSvc,
		performaAgg,
	}
}

// DashboardAdmin implements DashboardService.
func (d *DashboardServiceImpl) DashboardAdmin(ctx context.Context, startDate string, endDate string) (*params.DashboardResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous period dihitung mundur dengan panjang hari yang sama
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get All Studio
	studios, err := d.studioSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var charts []params.Chart
	for day := *start; !day.After(*end); day = day.AddDate(0, 0, 1) {
		startDay := timehandler.StartOfDay(day)
		endDay := timehandler.EndOfDay(day)

		_, total, err := d.performaAgg.Calculate(ctx, startDay, endDay)
		if err != nil {
			return nil, err
		}
		charts = append(charts, params.Chart{
			Date:       timehandler.FormatDate(&day),
			GMV:        total.GMV,
			Ads:        total.Ads,
			Commission: total.CommissionTotal,
			Income:     total.Income,
		})
	}

	list := []performaparams.PerformaStudioItemResponse{}

	var (
		currGMV, prevGMV               int64
		currCommission, prevCommission int64
		currAds, prevAds               int64
		currIncome, prevIncome         int64
	)

	for _, studio := range studios {
		_, currTotal, _ := d.performaAgg.CalculateByStudio(ctx, fmt.Sprint(studio.ID), start, end)
		_, prevTotal, _ := d.performaAgg.CalculateByStudio(ctx, fmt.Sprint(studio.ID), &prevStart, &prevEnd)

		// Tambahkan ke list
		list = append(list, performaparams.PerformaStudioItemResponse{
			StudioID:   fmt.Sprintf("%d", studio.ID),
			StudioName: studio.Name,
			PerformaMetricItem: performaparams.PerformaMetricItem{
				Commission: currTotal.CommissionTotal,
				GMV:        currTotal.GMV,
				Ads:        currTotal.Ads,
				Income:     currTotal.Income,
			},
		})

		// Akumulasi ke total metrics
		currGMV += currTotal.GMV
		prevGMV += prevTotal.GMV
		currCommission += currTotal.CommissionTotal
		prevCommission += prevTotal.CommissionTotal
		currAds += currTotal.Ads
		prevAds += prevTotal.Ads
		currIncome += currTotal.Income
		prevIncome += prevTotal.Income
	}

	accounts, err := d.accountSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	hosts, err := d.hostSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	results := &params.DashboardResponse{
		CurrentPeriod: params.PeriodInfo{
			Start: timehandler.FormatDate(start),
			End:   timehandler.FormatDate(end),
			Days:  days,
		},
		PreviousPeriod: params.PeriodInfo{
			Start: timehandler.FormatDate(&prevStart),
			End:   timehandler.FormatDate(&prevEnd),
			Days:  days,
		},
		Metrics: params.Metrics{
			Commission: NewMetric(currCommission, prevCommission),
			GMV:        NewMetric(currGMV, prevGMV),
			Ads:        NewMetric(int64(currAds), int64(prevAds)),
			Income:     NewMetric(currIncome, prevIncome),
			Account:    int64(len(accounts)),
			Host:       int64(len(hosts)),
			Studio:     int64(len(studios)),
		},
		Charts: charts,
		List:   list,
	}

	return results, nil
}
