package performa

import (
	"context"
	"fmt"
	"time"

	accountadsparam "github.com/royhairul/live-studio-api/internal/domains/accountads/params"
	accountsessionparam "github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	attendanceparam "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"
	transactionparam "github.com/royhairul/live-studio-api/internal/domains/transaction/params"

	accountadsservice "github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

type PerformaAggregatorImpl struct {
	hostSvc           hostservice.HostService
	attendanceSvc     attendanceservice.AttendanceService
	accountsessionSvc accountsessionservice.AccountsessionService
	accountadsSvc     accountadsservice.AccountadsService
	transactionSvc    transactionservice.TransactionService
}

func NewPerformaAggregator(
	hostSvc hostservice.HostService,
	attendanceSvc attendanceservice.AttendanceService,
	accountsessionSvc accountsessionservice.AccountsessionService,
	accountadsSvc accountadsservice.AccountadsService,
	transactionSvc transactionservice.TransactionService,
) PerformaAggregator {
	return &PerformaAggregatorImpl{
		hostSvc:           hostSvc,
		attendanceSvc:     attendanceSvc,
		accountsessionSvc: accountsessionSvc,
		accountadsSvc:     accountadsSvc,
		transactionSvc:    transactionSvc,
	}
}

// Calculate implements PerformaAggregator.
func (p *PerformaAggregatorImpl) Calculate(ctx context.Context, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	// Ambil attendance
	attendances, err := p.attendanceSvc.WithDateRange(*startDate, *endDate).FindAll(ctx)
	if err != nil {
		return nil, TotalPerformaAccount{}, err
	}

	return p.aggregateByAttendances(ctx, attendances, startDate, endDate)
}

func (p *PerformaAggregatorImpl) CalculateByHosts(ctx context.Context, startDate, endDate *time.Time) ([]*performaparam.PerformaHostSummaryResponse, error) {
	hosts, err := p.hostSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*performaparam.PerformaHostSummaryResponse, 0, len(hosts))

	for _, host := range hosts {
		// Gunakan CalculateByHost untuk dapatkan total per host
		detail, err := p.CalculateByHost(ctx, host.ID.String(), startDate, endDate)
		if err != nil {
			continue // skip host bermasalah tanpa hentikan seluruh proses
		}

		results = append(results, &performaparam.PerformaHostSummaryResponse{
			ID:            detail.ID,
			Name:          detail.Name,
			TotalDuration: detail.TotalDuration,
			TotalSales:    detail.TotalSales,
			TotalPaid:     detail.TotalPaid,
		})
	}

	return results, nil
}

func (p *PerformaAggregatorImpl) CalculateByHost(ctx context.Context, hostID string, startDate, endDate *time.Time) (performaparam.PerformaHostDetailResponse, error) {
	attendances, err := p.attendanceSvc.
		WithHostID(hostID).
		WithDateRange(*startDate, *endDate).
		FindAll(ctx)
	if err != nil {
		return performaparam.PerformaHostDetailResponse{}, err
	}

	total := &performaparam.TotalPerformaHost{}
	list := make([]performaparam.PerformaHostItemResponse, 0)

	for _, att := range attendances {
		sessions, err := p.accountsessionSvc.
			WithAttendanceID(fmt.Sprint(att.ID)).
			FindAll(ctx)
		if err != nil {
			continue
		}

		for _, session := range sessions {
			item := performaparam.PerformaHostItemResponse{
				AccountName: session.AccountName,
				Duration:    session.Duration,
				Sales:       int64(session.GMVSales),
				Paid:        int64(session.GMVPaid),
			}

			total.Duration += item.Duration
			total.Sales += item.Sales
			total.Paid += item.Paid
			list = append(list, item)
		}
	}

	host, err := p.hostSvc.FindByID(ctx, hostID)
	if err != nil {
		return performaparam.PerformaHostDetailResponse{}, err
	}

	var avgSales, avgPaid int64
	if total.Duration > 0 {
		avgSales = total.Sales / total.Duration
		avgPaid = total.Paid / total.Duration
	}

	return performaparam.PerformaHostDetailResponse{
		PerformaHostSummaryResponse: performaparam.PerformaHostSummaryResponse{
			ID:            host.ID.String(),
			Name:          host.Name,
			TotalDuration: total.Duration,
			TotalSales:    total.Sales,
			TotalPaid:     total.Paid,
		},
		AvgSales: avgSales,
		AvgPaid:  avgPaid,
		Total:    total,
		List:     list,
	}, nil
}

// CalculateByStudio implements PerformaAggregator.
func (p *PerformaAggregatorImpl) CalculateByStudio(ctx context.Context, studio_id string, startDate *time.Time, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	// Get Attendances
	attendances, err := p.attendanceSvc.WithStudioID(studio_id).WithDateRange(*startDate, *endDate).FindAll(ctx)
	if err != nil {
		return nil, TotalPerformaAccount{}, err
	}

	return p.aggregateByAttendances(ctx, attendances, startDate, endDate)
}

func (p *PerformaAggregatorImpl) aggregateByAttendances(
	ctx context.Context,
	attendances []*attendanceparam.AttendanceResponse,
	startDate, endDate *time.Time,
) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	list := []performaparam.PerformaStudioDetailItemResponse{}
	total := TotalPerformaAccount{}

	// Kumpulkan semua session
	var allSessions []*accountsessionparam.AccountsessionResponse
	for _, att := range attendances {
		sessions, _ := p.accountsessionSvc.WithAttendanceID(fmt.Sprint(att.ID)).FindAll(ctx)
		allSessions = append(allSessions, sessions...)
	}

	// Gabungkan berdasarkan AccountID
	accountIDs := make(map[string]*accountsessionparam.AccountsessionResponse)
	for _, session := range allSessions {
		key := fmt.Sprint(session.AccountID)
		if acc, exists := accountIDs[key]; exists {
			acc.GMVPaid += session.GMVPaid
		} else {
			accountIDs[key] = session
		}
	}

	// Ambil data transaksi & ads per akun
	txMap := map[string]transactionparam.Commission{}
	adsMap := map[string]accountadsparam.AccountadsTotalResponse{}

	for id := range accountIDs {
		tx, err := p.transactionSvc.WithAccountID(id).WithDate(*startDate, *endDate).GetTotalCommission(ctx)
		if err != nil {
			return nil, TotalPerformaAccount{}, err
		}
		txMap[id] = *tx

		ads, err := p.accountadsSvc.WithAccountID(id).WithDateRange(*startDate, *endDate).GetTotalAds(ctx)
		if err != nil {
			return nil, TotalPerformaAccount{}, err
		}
		adsMap[id] = *ads
	}

	// Bangun hasil akhir per akun
	for id, session := range accountIDs {
		txCommission := txMap[id]
		ads := adsMap[id]

		item := performaparam.PerformaStudioDetailItemResponse{
			AccountID:   session.AccountID,
			AccountName: session.AccountName,
			PerformaMetricItem: performaparam.PerformaMetricItem{
				GMV:        int64(session.GMVPaid),
				Commission: txCommission.Total,
				Ads:        int64(ads.TotalAds),
				Income:     txCommission.Total - int64(ads.TotalAds),
				Acos:       calcACOS(int64(ads.TotalAds), int64(session.GMVPaid)),
				Roas:       calcROAS(int64(ads.TotalAds), int64(session.GMVPaid)),
			},
		}

		list = append(list, item)

		total.GMV += item.GMV
		total.Ads += item.Ads
		total.CommissionPaid += txCommission.Paid
		total.CommissionPending += txCommission.Pending
		total.CommissionTotal += txCommission.Total
		total.Income += item.Income
	}

	return list, total, nil
}

// Hitung ACOS (%)
func calcACOS(ads, revenue int64) float64 {
	if revenue == 0 {
		return 0
	}
	return (float64(ads) / float64(revenue)) * 100
}

// Hitung ROAS (rasio)
func calcROAS(ads, revenue int64) float64 {
	if ads == 0 {
		return 0 // Tidak ada biaya iklan → anggap ROAS = 0
	}
	return float64(revenue) / float64(ads)
}
