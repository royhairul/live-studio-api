package service

import (
	"context"
	"fmt"

	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"

	// aggregator
	performaagg "github.com/royhairul/live-studio-api/internal/aggregator/performa"

	// service
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
)

type PerformaServiceImpl struct {
	hostSvc    hostservice.HostService
	accountSvc accountservice.AccountService
	studioSvc  studioservice.StudioService
	aggregator performaagg.PerformaAggregator
}

func NewPerformaService(
	hostSvc hostservice.HostService,
	accountSvc accountservice.AccountService,
	studioSvc studioservice.StudioService,
	aggregator performaagg.PerformaAggregator,
) PerformaService {
	return &PerformaServiceImpl{
		hostSvc,
		accountSvc,
		studioSvc,
		aggregator,
	}
}

// GetHosts implements PerformaService.
func (p *PerformaServiceImpl) GetHosts(ctx context.Context, startDate string, endDate string) ([]*params.PerformaHostSummaryResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	results, err := p.aggregator.CalculateByHosts(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetHostByID implements PerformaService.
func (p *PerformaServiceImpl) GetHostByID(ctx context.Context, id string, startDate string, endDate string) (*params.PerformaHostDetailResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	results, err := p.aggregator.CalculateByHost(ctx, id, start, end)
	if err != nil {
		return nil, err
	}

	// Angka engagement berasal dari tabel lives, bukan account_sessions, jadi
	// dihitung terpisah. Hanya dipakai halaman detail — GetHosts tidak ikut.
	metrics, err := p.aggregator.CalculateLiveMetricsByHost(ctx, id, start, end)
	if err != nil {
		return nil, err
	}
	results.PerformaLiveMetrics = metrics

	return &results, nil
}

// GetAccounts implements PerformaService.
func (p *PerformaServiceImpl) GetAccounts(ctx context.Context, startDate, endDate string) (*params.PerformaAccountResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Hitung durasi (days) dan periode sebelumnya
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// === Current Period ===
	currList, currTotal, err := p.aggregator.Calculate(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to build current performa list: %w", err)
	}

	// === Previous Period ===
	// Rentang sebelumnya, bukan rentang berjalan — sebelumnya keduanya memakai
	// start/end yang sama sehingga diff dan ratio selalu nol.
	_, prevTotal, err := p.aggregator.Calculate(ctx, &prevStart, &prevEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to build previous performa list: %w", err)
	}

	// === Build response ===
	results := &params.PerformaAccountResponse{
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
			Commission: NewMetric(currTotal.CommissionTotal, prevTotal.CommissionTotal),
			GMV:        NewMetric(currTotal.GMV, prevTotal.GMV),
			Ads:        NewMetric(currTotal.Ads, prevTotal.Ads),
			Income:     NewMetric(currTotal.Income, prevTotal.Income),
		},
		List: currList,

		// Rasio dihitung dari seluruh sesi sekaligus, bukan rata-rata per akun.
		PerformaLiveMetrics: currTotal.PerformaLiveMetrics,
	}

	return results, nil
}

// GetStudios implements PerformaService.
func (p *PerformaServiceImpl) GetStudios(ctx context.Context, startDate string, endDate string) (*params.PerformaStudioResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous period dihitung mundur dengan panjang hari yang sama
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get All Studio
	studios, err := p.studioSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var (
		currGMV, prevGMV                         int64
		currCommissionTotal, prevCommissionTotal int64
		currAds, prevAds                         int64
		currIncome, prevIncome                   int64
	)

	list := []params.PerformaStudioItemResponse{}
	for _, studio := range studios {
		_, currTotal, err := p.aggregator.CalculateByStudio(ctx, fmt.Sprint(studio.ID), start, end)
		if err != nil {
			return nil, err
		}

		_, prevTotal, err := p.aggregator.CalculateByStudio(ctx, fmt.Sprint(studio.ID), &prevStart, &prevEnd)
		if err != nil {
			return nil, err
		}

		item := params.PerformaStudioItemResponse{
			StudioID:   fmt.Sprint(studio.ID),
			StudioName: studio.Name,
			PerformaMetricItem: params.PerformaMetricItem{
				GMV: currTotal.GMV,
				// Total, bukan paid+pending: Income diturunkan dari total, jadi
				// memakai paid+pending di sini membuat "Komisi 0" bisa muncul
				// berdampingan dengan "Pendapatan 542,5 JT".
				Commission:          currTotal.CommissionTotal,
				Ads:                 currTotal.Ads,
				Income:              currTotal.Income,
				PerformaLiveMetrics: currTotal.PerformaLiveMetrics,
			},
		}

		list = append(list, item)

		// Calculate metrics
		currGMV += currTotal.GMV
		prevGMV += prevTotal.GMV
		currCommissionTotal += currTotal.CommissionTotal
		prevCommissionTotal += prevTotal.CommissionTotal
		currAds += currTotal.Ads
		prevAds += prevTotal.Ads
		currIncome += currTotal.Income
		prevIncome += prevTotal.Income
	}

	results := &params.PerformaStudioResponse{
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
			Commission: NewMetric(currCommissionTotal, prevCommissionTotal),
			GMV:        NewMetric(currGMV, prevGMV),
			Ads:        NewMetric(int64(currAds), int64(prevAds)),
			Income:     NewMetric(currIncome, prevIncome),
		},
		List: list,
	}

	// Rasio lintas studio tidak bisa dijumlahkan dari angka per studio, jadi
	// dihitung ulang sekali dari seluruh sesi.
	liveMetrics, err := p.aggregator.CalculateLiveMetrics(ctx, nil, start, end)
	if err != nil {
		return nil, err
	}
	results.PerformaLiveMetrics = liveMetrics

	return results, nil
}

// GetStudioByID implements PerformaService.
func (p *PerformaServiceImpl) GetStudioByID(ctx context.Context, id string, startDate string, endDate string) (*params.PerformaStudioDetailResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous range
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get Studio by ID
	studio, err := p.studioSvc.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Current Performa
	currList, currTotal, err := p.aggregator.CalculateByStudio(ctx, fmt.Sprint(studio.ID), start, end)
	if err != nil {
		return nil, err
	}

	// Previous Performa
	_, prevTotal, err := p.aggregator.CalculateByStudio(ctx, fmt.Sprint(studio.ID), &prevStart, &prevEnd)
	if err != nil {
		return nil, err
	}

	// Build response
	result := &params.PerformaStudioDetailResponse{
		StudioID:   studio.ID,
		StudioName: studio.Name,
		List:       currList,

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

		// Aggregate metrics
		Metrics: params.Metrics{
			GMV:        NewMetric(currTotal.GMV, prevTotal.GMV),
			Ads:        NewMetric(currTotal.Ads, prevTotal.Ads),
			Commission: NewMetric(currTotal.CommissionTotal, prevTotal.CommissionTotal),
			Income:     NewMetric(currTotal.Income, prevTotal.Income),
		},

		PerformaLiveMetrics: currTotal.PerformaLiveMetrics,
	}

	return result, nil
}
