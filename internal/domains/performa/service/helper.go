package service

import (
	"fmt"
	"time"

	attendanceparams "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

// buildPerformaAccountDetailList membangun detail performa akun + agregasi
func (p *PerformaServiceImpl) buildPerformaAccountDetailList(
	attendances []attendanceparams.AttendanceResponse,
	start, end *time.Time,
) ([]params.PerformaAccountDetailItemResponse, int64, int64, int64, int64, int64, error) {
	var (
		totalGMV               int64
		totalAds               int64
		totalCommissionPaid    int64
		totalCommissionPending int64
		totalIncome            int64
	)

	// Gunakan map supaya unik per account
	accountMap := make(map[uint]*params.PerformaAccountDetailItemResponse)

	// === Step 1: kumpulin GMV per account ===
	for _, att := range attendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, 0, 0, 0, 0, 0,
				fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			accountGMV := int64(session.GMVPaid)
			totalGMV += accountGMV

			// Pastikan account ada di map
			item, exists := accountMap[session.AccountID]
			if !exists {
				item = &params.PerformaAccountDetailItemResponse{
					AccountID:   session.AccountID,
					AccountName: session.AccountName,
				}
				accountMap[session.AccountID] = item
			}

			item.GMV += accountGMV
		}
	}

	// === Step 2: hitung Commission, Ads, Income per account ===
	for accID, item := range accountMap {
		// --- Commission ---
		var commissionPaid, commissionPending int64
		transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", accID), start, end)
		if err != nil {
			return nil, 0, 0, 0, 0, 0,
				fmt.Errorf("failed to get transactions for account %d: %w", accID, err)
		}
		for _, tx := range transactions {
			commissionPaid += int64(tx.Commission.Paid)
			commissionPending += int64(tx.Commission.Pending)
		}

		// --- Ads ---
		var accountAds int64
		ads, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", accID))
		if err != nil {
			return nil, 0, 0, 0, 0, 0,
				fmt.Errorf("failed to get ads for account %d: %w", accID, err)
		}
		for _, ad := range ads {
			accountAds += int64(ad.Ads)
		}

		// --- Income ---
		accountIncome := (commissionPaid + commissionPending) - accountAds

		// Update item
		item.Commission = commissionPaid + commissionPending
		item.Ads = accountAds
		item.Acos = calcACOS(accountAds, item.GMV)
		item.Roas = calcROAS(accountAds, item.GMV)
		item.Income = accountIncome

		// Update total
		totalCommissionPaid += commissionPaid
		totalCommissionPending += commissionPending
		totalAds += accountAds
		totalIncome += accountIncome
	}

	// === Step 3: convert map ke slice ===
	list := make([]params.PerformaAccountDetailItemResponse, 0, len(accountMap))
	for _, v := range accountMap {
		list = append(list, *v)
	}

	return list, totalGMV, totalAds, totalCommissionPaid, totalCommissionPending, totalIncome, nil
}

// normalizeDateRange memastikan start dan end mencakup full day
func normalizeDateRange(start, end *time.Time) (*time.Time, *time.Time) {
	if start == nil || end == nil {
		return start, end
	}

	dayStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	dayEnd := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, int(time.Nanosecond*999999999), end.Location())
	return &dayStart, &dayEnd
}

// buildPerformaDetailList membangun detail performa studio + agregasi
func (p *PerformaServiceImpl) buildPerformaDetailList(
	attendances []attendanceparams.AttendanceResponse,
	start, end *time.Time,
) ([]params.PerformaStudioDetailItemResponse, int64, int64, int64, int64, int64, error) {
	// normalisasi range tanggal
	start, end = normalizeDateRange(start, end)

	var (
		totalGMV               int64
		totalAds               int64
		totalCommissionPaid    int64
		totalCommissionPending int64
		totalIncome            int64
		list                   []params.PerformaStudioDetailItemResponse
	)

	for _, att := range attendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			// --- GMV ---
			accountGMV := int64(session.GMVPaid)
			totalGMV += accountGMV

			// --- Commission ---
			var commissionPaid, commissionPending int64
			transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", session.AccountID), start, end)
			if err != nil {
				return nil, 0, 0, 0, 0, 0, fmt.Errorf("failed to get transactions for account %d: %w", session.AccountID, err)
			}
			for _, ac := range transactions {
				commissionPaid += int64(ac.Commission.Paid)
				commissionPending += int64(ac.Commission.Pending)
			}
			totalCommissionPaid += commissionPaid
			totalCommissionPending += commissionPending

			// --- Ads ---
			var accountAds int64
			ads, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", session.AccountID))
			if err != nil {
				return nil, 0, 0, 0, 0, 0, fmt.Errorf("failed to get ads for account %d: %w", session.AccountID, err)
			}
			for _, a := range ads {
				accountAds += int64(a.Ads)
			}
			totalAds += accountAds

			// --- Income ---
			accountIncome := (commissionPaid + commissionPending) - accountAds
			totalIncome += accountIncome

			// --- Append detail ---
			list = append(list, params.PerformaStudioDetailItemResponse{
				AccountID:   session.AccountID,
				AccountName: session.AccountName,
				GMV:         accountGMV,
				Commission:  commissionPaid + commissionPending,
				Ads:         accountAds,
				Acos:        calcACOS(accountAds, accountGMV),
				Roas:        calcROAS(accountAds, accountGMV),
				Income:      accountIncome,
			})
		}
	}

	return list, totalGMV, totalAds, totalCommissionPaid, totalCommissionPending, totalIncome, nil
}

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
