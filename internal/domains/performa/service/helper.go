package service

import (
	"fmt"
	"time"

	attendanceparams "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
)

// Helper untuk build detail list & agregasi
func (p *PerformaServiceImpl) buildPerformaAccountDetailList(attendances []attendanceparams.AttendanceResponse, start, end *time.Time) ([]params.PerformaAccountDetailItemResponse, int64, int64, int64, int64, int64, error) {
	var totalGMV, totalAds, totalCommissionPaid, totalCommissionPending, totalIncome int64
	list := []params.PerformaAccountDetailItemResponse{}

	for _, att := range attendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			accountGMV := int64(session.GMVPaid)
			totalGMV += accountGMV

			// Transactions → Commission
			transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", session.AccountID), start, end)
			if err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}

			var commissionPaid, commissionPending int64
			for _, ac := range transactions {
				commissionPaid += int64(ac.Commission.Paid)
				commissionPending += int64(ac.Commission.Pending)
			}

			// Ads
			ads, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", session.AccountID))
			if err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}
			accountAds := int64(0)
			for _, a := range ads {
				accountAds += int64(a.Ads)
			}
			totalAds += accountAds

			// Income
			accountIncome := (commissionPaid + commissionPending) - accountAds
			totalCommissionPaid += commissionPaid
			totalCommissionPending += commissionPending

			totalIncome += accountIncome

			// Detail per account
			list = append(list, params.PerformaAccountDetailItemResponse{
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

// Helper untuk build detail list & agregasi
func (p *PerformaServiceImpl) buildPerformaDetailList(attendances []attendanceparams.AttendanceResponse, start, end *time.Time) ([]params.PerformaStudioDetailItemResponse, int64, int64, int64, int64, int64, error) {
	var totalGMV, totalAds, totalCommissionPaid, totalCommissionPending, totalIncome int64
	list := []params.PerformaStudioDetailItemResponse{}

	for _, att := range attendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			accountGMV := int64(session.GMVPaid)
			totalGMV += accountGMV

			// Transactions → Commission
			transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", session.AccountID), start, end)
			if err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}

			var commissionPaid, commissionPending int64
			for _, ac := range transactions {
				commissionPaid += int64(ac.Commission.Paid)
				commissionPending += int64(ac.Commission.Pending)
			}

			// Ads
			ads, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", session.AccountID))
			if err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}
			accountAds := int64(0)
			for _, a := range ads {
				accountAds += int64(a.Ads)
			}
			totalAds += accountAds

			// Income
			accountIncome := (commissionPaid + commissionPending) - accountAds
			totalCommissionPaid += commissionPaid
			totalCommissionPending += commissionPending

			totalIncome += accountIncome

			// Detail per account
			list = append(list, params.PerformaStudioDetailItemResponse{
				AccountID:   session.AccountID,
				AccountName: session.AccountName,
				GMV:         accountGMV,
				// CommissionPaid:    commissionPaid,
				// CommissionPending: commissionPending,
				Commission: commissionPaid + commissionPending,
				Ads:        accountAds,
				Acos:       calcACOS(accountAds, accountGMV),
				Roas:       calcROAS(accountAds, accountGMV),
				Income:     accountIncome,
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
