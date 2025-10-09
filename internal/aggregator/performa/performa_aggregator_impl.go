package performa

import (
	"fmt"
	"time"

	accountadsparam "github.com/royhairul/live-studio-api/internal/domains/accountads/params"
	accountsessionparam "github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	performaparam "github.com/royhairul/live-studio-api/internal/domains/performa/params"
	transactionparam "github.com/royhairul/live-studio-api/internal/domains/transaction/params"

	accountadsservice "github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

type PerformaAggregatorImpl struct {
	attendanceSvc     attendanceservice.AttendanceService
	accountsessionSvc accountsessionservice.AccountsessionService
	accountadsSvc     accountadsservice.AccountadsService
	transactionSvc    transactionservice.TransactionService
}

func NewPerformaAggregator(
	attendanceSvc attendanceservice.AttendanceService,
	accountsessionSvc accountsessionservice.AccountsessionService,
	accountadsSvc accountadsservice.AccountadsService,
	transactionSvc transactionservice.TransactionService,
) PerformaAggregator {
	return &PerformaAggregatorImpl{
		attendanceSvc:     attendanceSvc,
		accountsessionSvc: accountsessionSvc,
		accountadsSvc:     accountadsSvc,
		transactionSvc:    transactionSvc,
	}
}

// Calculate implements PerformaAggregator.
func (p *PerformaAggregatorImpl) Calculate(
	startDate, endDate *time.Time,
) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerforma, error) {
	list := []performaparam.PerformaStudioDetailItemResponse{}
	total := TotalPerforma{}

	// Ambil attendance
	attendances, err := p.attendanceSvc.WithDateRange(*startDate, *endDate).FindAll()
	if err != nil {
		return nil, TotalPerforma{}, err
	}

	// Kumpulkan semua session
	var allSessions []*accountsessionparam.AccountsessionResponse
	for _, att := range attendances {
		sessions, _ := p.accountsessionSvc.WithAttendanceID(fmt.Sprint(att.ID)).FindAll()
		allSessions = append(allSessions, sessions...)
	}

	// Ambil unique accountIDs
	accountIDs := make(map[string]*accountsessionparam.AccountsessionResponse)
	for _, session := range allSessions {
		key := fmt.Sprint(session.AccountID)
		// Simpan session pertama sebagai referensi AccountName
		if _, exists := accountIDs[key]; !exists {
			accountIDs[key] = session
		} else {
			// Akumulasi GMV kalau sudah ada
			accountIDs[key].GMVPaid += session.GMVPaid
		}
	}

	// Ambil transaksi dan ads per account
	txMap := make(map[string]transactionparam.TransactionCommission)
	adsMap := make(map[string]accountadsparam.AccountadsTotalResponse)

	for id := range accountIDs {
		// Commission
		transaction, err := p.transactionSvc.
			WithAccountID(id).
			WithDate(*startDate, *endDate).
			GetTotalCommission()
		if err != nil {
			return nil, TotalPerforma{}, err
		}
		txMap[id] = *transaction

		// Ads
		adsTotal, err := p.accountadsSvc.
			WithAccountID(id).
			WithDateRange(*startDate, *endDate).
			GetTotalAds()
		if err != nil {
			return nil, TotalPerforma{}, err
		}
		adsMap[id] = *adsTotal
	}

	// Bangun hasil per account
	for id, session := range accountIDs {
		tx := txMap[id]
		ads := adsMap[id]

		item := performaparam.PerformaStudioDetailItemResponse{
			AccountID:   session.AccountID,
			AccountName: session.AccountName,
			GMV:         int64(session.GMVPaid),
			Commission:  tx.CommissionTotal,
			Ads:         int64(ads.TotalAds),
			Income:      tx.CommissionTotal - int64(ads.TotalAds),
			Acos:        calcACOS(int64(ads.TotalAds), int64(session.GMVPaid)),
			Roas:        calcROAS(int64(ads.TotalAds), int64(session.GMVPaid)),
		}

		list = append(list, item)

		// Update global total
		total.GMV += item.GMV
		total.Ads += item.Ads
		total.CommissionPaid += tx.CommissionPaid
		total.CommissionPending += tx.CommissionPending
		total.CommissionTotal += tx.CommissionTotal
		total.Income += item.Income
	}

	return list, total, nil
}

// CalculateByStudio implements PerformaAggregator.
func (p *PerformaAggregatorImpl) CalculateByStudio(studio_id string, startDate *time.Time, endDate *time.Time) ([]performaparam.PerformaStudioDetailItemResponse, TotalPerforma, error) {
	list := []performaparam.PerformaStudioDetailItemResponse{}
	total := TotalPerforma{}

	// Get Attendances
	attendances, err := p.attendanceSvc.WithStudioID(studio_id).WithDateRange(*startDate, *endDate).FindAll()
	if err != nil {
		return nil, TotalPerforma{}, err
	}

	// Kumpulkan semua session
	var allSessions []*accountsessionparam.AccountsessionResponse
	for _, att := range attendances {
		sessions, _ := p.accountsessionSvc.WithAttendanceID(fmt.Sprint(att.ID)).FindAll()
		allSessions = append(allSessions, sessions...)
	}

	// Get unique accountIDs
	accountIDs := make(map[string]*accountsessionparam.AccountsessionResponse)
	for _, session := range allSessions {
		key := fmt.Sprint(session.AccountID)
		// Save first session as reference
		if _, exists := accountIDs[key]; !exists {
			accountIDs[key] = session
		} else {
			// Acumulate GMV
			accountIDs[key].GMVPaid += session.GMVPaid
		}
	}

	// Get tranaction per account
	txMap := make(map[string]transactionparam.TransactionCommission)
	adsMap := make(map[string]accountadsparam.AccountadsTotalResponse)

	for id := range accountIDs {
		// Commission
		transaction, err := p.transactionSvc.
			WithAccountID(id).
			WithDate(*startDate, *endDate).
			GetTotalCommission()
		if err != nil {
			return nil, TotalPerforma{}, err
		}
		txMap[id] = *transaction

		// Ads
		adsTotal, err := p.accountadsSvc.
			WithAccountID(id).
			WithDateRange(*startDate, *endDate).
			GetTotalAds()
		if err != nil {
			return nil, TotalPerforma{}, err
		}
		adsMap[id] = *adsTotal
	}

	// Create result per account
	for id, session := range accountIDs {
		tx := txMap[id]
		ads := adsMap[id]

		item := performaparam.PerformaStudioDetailItemResponse{
			AccountID:   session.AccountID,
			AccountName: session.AccountName,
			GMV:         int64(session.GMVPaid),
			Commission:  tx.CommissionTotal,
			Ads:         int64(ads.TotalAds),
			Income:      tx.CommissionTotal - int64(ads.TotalAds),
			Acos:        calcACOS(int64(ads.TotalAds), int64(session.GMVPaid)),
			Roas:        calcROAS(int64(ads.TotalAds), int64(session.GMVPaid)),
		}

		list = append(list, item)

		// Update global total
		total.GMV += item.GMV
		total.Ads += item.Ads
		total.CommissionPaid += tx.CommissionPaid
		total.CommissionPending += tx.CommissionPending
		total.CommissionTotal += tx.CommissionTotal
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
