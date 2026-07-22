package performa

import (
	"context"
	"fmt"
	"time"

	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"

	hostaccountparam "github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
	liveentity "github.com/royhairul/live-studio-api/internal/domains/live/entity"
	liveparam "github.com/royhairul/live-studio-api/internal/domains/live/params"

	accountadsservice "github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"
	hostaccountservice "github.com/royhairul/live-studio-api/internal/domains/hostaccount/service"
	liverepository "github.com/royhairul/live-studio-api/internal/domains/live/repository"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

// Attendance and account sessions are deliberately absent: every performa
// figure now comes from lives, bridged to hosts by host_accounts.
type PerformaAggregatorImpl struct {
	hostSvc        hostservice.HostService
	accountadsSvc  accountadsservice.AccountadsService
	transactionSvc transactionservice.TransactionService

	// The repository rather than LiveService: the service pulls in the Shopee
	// HTTP client and its paginated response shape, neither of which is wanted
	// for a local aggregate.
	liveRepo liverepository.LiveRepository

	// The host↔account bridge for reporting, in place of attendance.
	hostAccountSvc hostaccountservice.HostAccountService
}

func NewPerformaAggregator(
	hostSvc hostservice.HostService,
	accountadsSvc accountadsservice.AccountadsService,
	transactionSvc transactionservice.TransactionService,
	liveRepo liverepository.LiveRepository,
	hostAccountSvc hostaccountservice.HostAccountService,
) PerformaAggregator {
	return &PerformaAggregatorImpl{
		hostSvc:        hostSvc,
		accountadsSvc:  accountadsSvc,
		transactionSvc: transactionSvc,
		liveRepo:       liveRepo,
		hostAccountSvc: hostAccountSvc,
	}
}

// Calculate implements PerformaAggregator.
func (p *PerformaAggregatorImpl) Calculate(ctx context.Context, startDate, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	return p.aggregateByLives(ctx, nil, startDate, endDate)
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

// CalculateByHost implements PerformaAggregator.
//
// Everything here comes from the lives table, bridged by host_accounts —
// attendance is no longer consulted at all. A host's numbers are the numbers of
// the accounts they held, whether or not anyone remembered to check in.
//
// Two contract details worth keeping in mind:
//
//   - Shopee reports duration in milliseconds, but total_duration has always
//     been seconds (the dashboard renders it with intToHumanTime), so it is
//     converted here rather than at the edges.
//   - list is now one row per account instead of one row per attendance
//     session, which is what its account_name column implied all along.
func (p *PerformaAggregatorImpl) CalculateByHost(ctx context.Context, hostID string, startDate, endDate *time.Time) (performaparam.PerformaHostDetailResponse, error) {
	host, err := p.hostSvc.FindByID(ctx, hostID)
	if err != nil {
		return performaparam.PerformaHostDetailResponse{}, err
	}

	sessions, err := p.hostLiveSessions(ctx, hostID, startDate, endDate)
	if err != nil {
		return performaparam.PerformaHostDetailResponse{}, err
	}

	total := &performaparam.TotalPerformaHost{}
	list := make([]performaparam.PerformaHostItemResponse, 0)

	// Satu baris per akun: satu akun bisa punya banyak sesi di dalam rentang.
	type accountTotals struct {
		name     string
		duration int64
		sales    float64
		paid     float64
	}

	perAccount := map[uint]*accountTotals{}
	liveDays := map[string]bool{}

	for _, live := range sessions {
		bucket, exists := perAccount[live.AccountID]
		if !exists {
			bucket = &accountTotals{name: live.Account.Name}
			perAccount[live.AccountID] = bucket
		}

		bucket.duration += live.Duration / 1000
		bucket.sales += live.PlacedSales
		bucket.paid += live.ConfirmedSales

		if live.StartTime != nil {
			liveDays[live.StartTime.Format("2006-01-02")] = true
		}
	}

	for _, bucket := range perAccount {
		item := performaparam.PerformaHostItemResponse{
			AccountName: bucket.name,
			Duration:    bucket.duration,
			Sales:       int64(bucket.sales),
			Paid:        int64(bucket.paid),
		}

		total.Duration += item.Duration
		total.Sales += item.Sales
		total.Paid += item.Paid
		list = append(list, item)
	}

	// Rata-rata dibagi jumlah hari yang benar-benar ada siarannya, bukan panjang
	// rentang laporan — rentang sebulan dengan siaran tiga hari tidak seharusnya
	// dibagi tiga puluh.
	var avgSales, avgPaid int64
	if days := int64(len(liveDays)); days > 0 {
		avgSales = total.Sales / days
		avgPaid = total.Paid / days
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

// hostLiveSessions returns the live sessions belonging to a host in the window:
// the sessions of the accounts they held, limited to the spans they held them.
//
// Shared by CalculateByHost and CalculateLiveMetricsByHost so the money figures
// and the engagement figures can never be drawn from different sets of sessions.
func (p *PerformaAggregatorImpl) hostLiveSessions(
	ctx context.Context,
	hostID string,
	startDate, endDate *time.Time,
) ([]*liveentity.Live, error) {
	spans, err := p.hostAccountSvc.AccountIDsForHost(ctx, hostID, *startDate, *endDate)
	if err != nil {
		return nil, err
	}
	if len(spans) == 0 {
		return nil, nil
	}

	// Satu akun bisa punya beberapa periode penugasan di dalam rentang laporan.
	spansByAccount := make(map[uint][]hostaccountparam.HostAccountSpan)
	for _, span := range spans {
		spansByAccount[span.AccountID] = append(spansByAccount[span.AccountID], span)
	}

	var matched []*liveentity.Live
	seen := make(map[int64]bool)

	for accountID, accountSpans := range spansByAccount {
		id := fmt.Sprint(accountID)

		// Page/PageSize sengaja nol: PageSize <= 0 berarti tanpa limit.
		lives, err := p.liveRepo.FindAll(ctx, liveparam.LiveFilter{
			AccountID: &id,
			StartTime: startDate,
			EndTime:   endDate,
		})
		if err != nil {
			return nil, err
		}

		for _, live := range lives {
			if live.StartTime == nil || seen[live.SessionID] {
				continue
			}
			for _, span := range accountSpans {
				if span.Contains(*live.StartTime) {
					seen[live.SessionID] = true
					matched = append(matched, live)
					break
				}
			}
		}
	}

	return matched, nil
}

// liveAccumulator sums engagement across live sessions.
//
// CTR and conversion rate are ratios, so they are accumulated as totals and
// divided once at the end — averaging per-session percentages would weight a
// 10-viewer session the same as a 10,000-viewer one.
type liveAccumulator struct {
	clicks    int
	views     int
	engagedUV int

	weightedCR float64
	sumCR      float64
	countCR    int

	// Guards against counting a session twice when spans overlap.
	seen map[int64]bool
}

func newLiveAccumulator() *liveAccumulator {
	return &liveAccumulator{seen: make(map[int64]bool)}
}

func (a *liveAccumulator) add(live *liveentity.Live) {
	if live == nil || a.seen[live.SessionID] {
		return
	}
	a.seen[live.SessionID] = true

	a.clicks += live.ProductClicks
	a.views += live.Views
	a.engagedUV += live.EngagedUV

	a.weightedCR += live.ConversionRate * float64(live.Views)
	a.sumCR += live.ConversionRate
	a.countCR++
}

func (a *liveAccumulator) result() performaparam.PerformaLiveMetrics {
	metrics := performaparam.PerformaLiveMetrics{ActiveViewers: a.engagedUV}

	if a.views > 0 {
		metrics.CTR = float64(a.clicks) / float64(a.views)
		metrics.ConversionRate = a.weightedCR / float64(a.views)
	} else if a.countCR > 0 {
		// Shopee kadang tidak melaporkan views sama sekali; pakai rata-rata biasa
		// supaya conversion rate tidak ikut hilang.
		metrics.ConversionRate = a.sumCR / float64(a.countCR)
	}

	return metrics
}

// CalculateLiveMetricsByHost implements PerformaAggregator.
//
// Rows in lives belong to an account, not to a host, so the bridge is the
// host_accounts assignment: a session counts for this host only when it started
// while the host actually held that account. Attendance is deliberately not
// consulted — presence and ownership are different questions, and a host with no
// check-in still owns their accounts' numbers.
//
// A host with no assignment yet reports zeros rather than an error.
func (p *PerformaAggregatorImpl) CalculateLiveMetricsByHost(
	ctx context.Context,
	hostID string,
	startDate, endDate *time.Time,
) (performaparam.PerformaLiveMetrics, error) {
	sessions, err := p.hostLiveSessions(ctx, hostID, startDate, endDate)
	if err != nil {
		return performaparam.PerformaLiveMetrics{}, err
	}

	acc := newLiveAccumulator()
	for _, live := range sessions {
		acc.add(live)
	}

	return acc.result(), nil
}

// CalculateLiveMetrics implements PerformaAggregator.
func (p *PerformaAggregatorImpl) CalculateLiveMetrics(
	ctx context.Context,
	studioID *string,
	startDate, endDate *time.Time,
) (performaparam.PerformaLiveMetrics, error) {
	lives, err := p.liveRepo.FindAll(ctx, liveparam.LiveFilter{
		StudioID:  studioID,
		StartTime: startDate,
		EndTime:   endDate,
	})
	if err != nil {
		return performaparam.PerformaLiveMetrics{}, err
	}

	acc := newLiveAccumulator()
	for _, live := range lives {
		acc.add(live)
	}

	return acc.result(), nil
}

// CalculateByStudio implements PerformaAggregator.
func (p *PerformaAggregatorImpl) CalculateByStudio(ctx context.Context, studio_id string, startDate *time.Time, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	return p.aggregateByLives(ctx, &studio_id, startDate, endDate)
}

// aggregateByLives builds the per-account performa list straight from the lives
// table. studioID nil means every studio.
//
// Attendance is not consulted at all: an account's sessions belong to that
// account whether or not anyone checked in. Under the old attendance-driven
// aggregation an account with live data but no check-in silently reported
// nothing, which is the bug this replaces.
//
// Note that GMV therefore now comes from Shopee's confirmed_sales rather than
// the manually recorded account_sessions GMV, so figures can differ from what
// the old aggregation produced for the same range.
func (p *PerformaAggregatorImpl) aggregateByLives(
	ctx context.Context,
	studioID *string,
	startDate, endDate *time.Time,
) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerformaAccount, error) {
	list := []performaparam.PerformaStudioDetailItemResponse{}
	total := TotalPerformaAccount{}

	lives, err := p.liveRepo.FindAll(ctx, liveparam.LiveFilter{
		StudioID:  studioID,
		StartTime: startDate,
		EndTime:   endDate,
	})
	if err != nil {
		return nil, TotalPerformaAccount{}, err
	}

	// Kelompokkan sesi per akun.
	type accountBucket struct {
		name  string
		sales float64
		live  *liveAccumulator
	}

	buckets := map[uint]*accountBucket{}
	for _, live := range lives {
		bucket, exists := buckets[live.AccountID]
		if !exists {
			bucket = &accountBucket{name: live.Account.Name, live: newLiveAccumulator()}
			buckets[live.AccountID] = bucket
		}

		bucket.sales += live.ConfirmedSales
		bucket.live.add(live)
	}

	totalLive := newLiveAccumulator()
	for _, live := range lives {
		totalLive.add(live)
	}

	for accountID, bucket := range buckets {
		id := fmt.Sprint(accountID)

		tx, err := p.transactionSvc.WithAccountID(id).WithDate(*startDate, *endDate).GetTotalCommission(ctx)
		if err != nil {
			return nil, TotalPerformaAccount{}, err
		}

		ads, err := p.accountadsSvc.WithAccountID(id).WithDateRange(*startDate, *endDate).GetTotalAds(ctx)
		if err != nil {
			return nil, TotalPerformaAccount{}, err
		}

		gmv := int64(bucket.sales)

		item := performaparam.PerformaStudioDetailItemResponse{
			AccountID:   accountID,
			AccountName: bucket.name,
			PerformaMetricItem: performaparam.PerformaMetricItem{
				GMV:                 gmv,
				Commission:          tx.Total,
				Ads:                 int64(ads.TotalAds),
				Income:              tx.Total - int64(ads.TotalAds),
				Acos:                calcACOS(int64(ads.TotalAds), gmv),
				Roas:                calcROAS(int64(ads.TotalAds), gmv),
				PerformaLiveMetrics: bucket.live.result(),
			},
		}

		list = append(list, item)

		total.GMV += item.GMV
		total.Ads += item.Ads
		total.CommissionPaid += tx.Paid
		total.CommissionPending += tx.Pending
		total.CommissionTotal += tx.Total
		total.Income += item.Income
	}

	// Rasio total dihitung ulang dari seluruh sesi, bukan dirata-rata dari nilai
	// per akun, supaya akun kecil tidak menarik angkanya.
	total.PerformaLiveMetrics = totalLive.result()

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
