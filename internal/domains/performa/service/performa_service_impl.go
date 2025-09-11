package service

import (
	"fmt"
	"log"
	"time"

	"github.com/royhairul/live-studio-api/helpers/timehandler"
	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
	"github.com/royhairul/live-studio-api/internal/domains/performa/repository"

	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"

	attendanceparams "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	accountadsservice "github.com/royhairul/live-studio-api/internal/domains/accountads/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

type PerformaServiceImpl struct {
	repository        repository.PerformaRepository
	hostSvc           hostservice.HostService
	attendanceSvc     attendanceservice.AttendanceService
	accountSessionSvc accountsessionservice.AccountsessionService
	accountSvc        accountservice.AccountService
	transactionSvc    transactionservice.TransactionService
	studioSvc         studioservice.StudioService
	accountAdsSvc     accountadsservice.AccountadsService
}

func NewPerformaService(
	repository repository.PerformaRepository,
	hostSvc hostservice.HostService,
	attendanceSvc attendanceservice.AttendanceService,
	accountSessionSvc accountsessionservice.AccountsessionService,
	accountSvc accountservice.AccountService,
	transactionSvc transactionservice.TransactionService,
	accountAdsSvc accountadsservice.AccountadsService,
	studioSvc studioservice.StudioService,
) PerformaService {
	return &PerformaServiceImpl{
		repository,
		hostSvc,
		attendanceSvc,
		accountSessionSvc,
		accountSvc,
		transactionSvc,
		studioSvc,
		accountAdsSvc,
	}
}

// GetHosts implements PerformaService.
func (p *PerformaServiceImpl) GetHosts(startDate string, endDate string) ([]*params.PerformaHostResponse, error) {
	// Set default value
	if startDate == "" {
		startDate = *timehandler.DateNow()
	}
	if endDate == "" {
		endDate = *timehandler.DateNow()
	}

	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get Attendances by date range
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	// Get all host
	hosts, err := p.hostSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.PerformaHostResponse

	for _, host := range hosts {
		var totalDuration int64
		var totalSales uint
		var totalPaid uint

		// Get attendances by Host ID
		for _, att := range attendances {
			if att.HostID == host.ID {
				continue
			}

			accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
			}

			for _, session := range accountSessions {
				if session.CheckIn == nil || session.CheckOut == nil {
					continue
				}

				// Calculate duration
				duration := session.CheckOut.Sub(*session.CheckIn)
				durationSeconds := duration.Seconds()

				totalDuration += int64(durationSeconds)
				totalSales += session.GMVSales
				totalPaid += session.GMVPaid
			}
		}

		results = append(results, &params.PerformaHostResponse{
			ID:            host.ID.String(),
			Name:          host.Name,
			TotalDuration: totalDuration,
			TotalSales:    totalSales,
		})
	}

	return results, nil
}

// GetHostByID implements PerformaService.
func (p *PerformaServiceImpl) GetHostByID(id string, startDate string, endDate string) (*params.PerformaHostDetailResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get attendance by startTime and endTime
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	// Get Host By ID
	host, err := p.hostSvc.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Get attendances by Host ID
	var hostAttendances []attendanceparams.AttendanceResponse
	for _, att := range attendances {
		if att.HostID == host.ID {
			hostAttendances = append(hostAttendances, *att)
		}
	}

	var totalDuration int64
	var totalSales uint
	var totalPaid uint
	detailList := []params.PerformaHostItemResponse{}

	for _, att := range hostAttendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			detailList = append(detailList, params.PerformaHostItemResponse{
				AccountName: session.AccountName,
				Duration:    session.Duration,
				Sales:       session.GMVSales,
				Paid:        session.GMVPaid,
			})

			log.Printf("account: %s", session.AccountName)

			// Hitung durasi
			duration := session.CheckOut.Sub(*session.CheckIn)
			durationSeconds := duration.Seconds()

			totalDuration += int64(durationSeconds)
			totalSales += session.GMVSales
			totalPaid += session.GMVPaid

		}
	}

	var avgSales uint
	var avgPaid uint
	count := len(detailList)
	if count > 0 {
		avgSales = totalSales / uint(count)
		avgPaid = totalPaid / uint(count)
	}

	result := &params.PerformaHostDetailResponse{
		ID:            host.ID.String(),
		Name:          host.Name,
		TotalDuration: totalDuration,
		TotalSales:    totalSales,
		AvgSales:      avgSales,
		AvgPaid:       avgPaid,
		List:          detailList,
	}

	return result, nil
}

// GetAccounts implements PerformaService.
func (p *PerformaServiceImpl) GetAccounts(startDate string, endDate string) ([]*params.PerformaStudioDetailItemResponse, error) {
	// Set default value
	if startDate == "" {
		startDate = *timehandler.DateNow()
	}
	if endDate == "" {
		endDate = *timehandler.DateNow()
	}

	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get Attendances (current + previous)
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	var currAttendances []attendanceparams.AttendanceResponse
	for _, att := range attendances {
		currAttendances = append(currAttendances, *att)
	}

	// ===== Current Period =====
	currList, _, _, _, _, _, err := p.buildPerformaDetailList(currAttendances, start, end)
	if err != nil {
		return nil, err
	}

	var list []*params.PerformaStudioDetailItemResponse

	for _, item := range currList {
		found := false
		for _, existing := range list {
			if existing.AccountID == item.AccountID {
				// Update semua field kecuali AccountID & AccountName
				existing.GMV += item.GMV
				// existing.CommissionPaid += item.CommissionPaid
				// existing.CommissionPending += item.CommissionPending
				existing.Commission += item.Commission
				existing.Ads += item.Ads
				existing.Income += item.Income

				// untuk field float (ACOS, ROAS) biasanya dihitung ulang rata-rata atau ratio
				// ini contoh: ambil nilai terbaru saja (overwrite)
				existing.Acos = item.Acos
				existing.Roas = item.Roas

				found = true
				break
			}
		}

		if !found {
			newItem := item
			list = append(list, &newItem)
		}
	}

	return list, nil
}

// GetAccountByID implements PerformaService.
func (p *PerformaServiceImpl) GetAccountByID() {
	panic("unimplemented")
}

// GetStudios implements PerformaService.
func (p *PerformaServiceImpl) GetStudios(startDate string, endDate string) (*params.PerformaStudioResponse, error) {
	// Set default value
	if startDate == "" {
		startDate = *timehandler.DateNow()
	}
	if endDate == "" {
		endDate = *timehandler.DateNow()
	}

	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous period dihitung mundur dengan panjang hari yang sama
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get Attendances (current + previous)
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	prevAttendances, err := p.attendanceSvc.FindByDateRange(&prevStart, &prevEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get prev attendances: %v", err)
	}

	// Get All Studio
	studios, err := p.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	list := []params.PerformaStudioItemResponse{}

	var currGMV, prevGMV int64
	var currCommissionPaid, prevCommissionPaid int64
	var currCommissionPending, prevCommissionPending int64
	var currAds, prevAds uint
	var currIncome, prevIncome int64

	for _, studio := range studios {

		// Get Account in this studio
		accounts, err := p.accountSvc.FindByStudio(fmt.Sprintf("%d", studio.ID))
		if err != nil {
			return nil, err
		}

		// Get Transactions per account
		for _, account := range accounts {
			// current period
			transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", account.ID), start, end)
			if err != nil {
				return nil, err
			}

			// previous period
			prevTransactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", account.ID), &prevStart, &prevEnd)
			if err != nil {
				return nil, err
			}

			// accumulate current
			for _, tx := range transactions {
				currCommissionPaid += int64(tx.Commission.Paid)
				currCommissionPending += int64(tx.Commission.Pending)
			}

			// // accumulate previous
			for _, tx := range prevTransactions {
				prevCommissionPaid += int64(tx.Commission.Paid)
				prevCommissionPending += int64(tx.Commission.Pending)
			}

			allAds, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range allAds {
				currAds += a.Ads
			}

			prevAllAds, err := p.accountAdsSvc.FindByDateAndAccounts(&prevStart, &prevEnd, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range prevAllAds {
				prevAds += a.Ads
			}

		}

		// GMV
		for _, att := range attendances {
			accountsessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, err
			}

			for _, session := range accountsessions {
				currGMV += int64(session.GMVPaid)
			}
		}

		// GMV
		for _, att := range prevAttendances {
			accountsessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, err
			}

			for _, session := range accountsessions {
				prevGMV += int64(session.GMVPaid)
			}
		}

		currIncome += (currCommissionPaid + currCommissionPending) - int64(currAds)
		prevIncome += (prevCommissionPaid + prevCommissionPending) - int64(prevAds)

		list = append(list, params.PerformaStudioItemResponse{
			StudioID:   fmt.Sprintf("%d", studio.ID),
			StudioName: studio.Name,
			Commission: currCommissionPaid + currCommissionPending,
			GMV:        currGMV,
			Ads:        int64(currAds),
			Income:     currIncome,
		})
	}

	// currTotalIncome := (currCommissionPaid + currCommissionPaid) - int64(currAds)
	// prevTotalIncome := (prevCommissionPaid + prevCommissionPaid) - int64(prevAds)

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
			// CommissionPaid:    NewMetric(currCommissionPaid, prevCommissionPaid),
			// CommissionPending: NewMetric(currCommissionPending, prevCommissionPending),
			Commission: NewMetric((currCommissionPaid + currCommissionPending), (prevCommissionPaid + prevCommissionPending)),
			GMV:        NewMetric(currGMV, prevGMV),
			Ads:        NewMetric(int64(currAds), int64(prevAds)),
			Income:     NewMetric(currIncome, prevIncome),
		},
		List: list,
	}

	return results, nil
}

// GetStudioByID implements PerformaService.
func (p *PerformaServiceImpl) GetStudioByID(id string, startDate string, endDate string) (*params.PerformaStudioDetailResponse, error) {
	// Default date
	if startDate == "" {
		startDate = *timehandler.DateNow()
	}
	if endDate == "" {
		endDate = *timehandler.DateNow()
	}

	// Parse date range
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Hitung previous range (durasi sama, mundur ke belakang)
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Ambil attendances current & prev
	currAttendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}
	prevAttendances, err := p.attendanceSvc.FindByDateRange(&prevStart, &prevEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous attendances: %v", err)
	}

	// Ambil studio
	studio, err := p.studioSvc.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Filter attendances by Studio ID
	var studioCurrAtt, studioPrevAtt []attendanceparams.AttendanceResponse
	for _, att := range currAttendances {
		if att.StudioID == studio.ID {
			studioCurrAtt = append(studioCurrAtt, *att)
		}
	}
	for _, att := range prevAttendances {
		if att.StudioID == studio.ID {
			studioPrevAtt = append(studioPrevAtt, *att)
		}
	}

	// ===== Current Period =====
	currList, currGMV, currAds, currCommissionPaid, currCommissionPending, currIncome, err := p.buildPerformaDetailList(studioCurrAtt, start, end)
	if err != nil {
		return nil, err
	}

	// ===== Previous Period =====
	_, prevGMV, prevAds, prevCommissionPaid, prevCommissionPending, prevIncome, err := p.buildPerformaDetailList(studioPrevAtt, &prevStart, &prevEnd)
	if err != nil {
		return nil, err
	}

	// Build response
	result := &params.PerformaStudioDetailResponse{
		StudioID:   studio.ID,
		StudioName: studio.Name,
		List:       currList, // untuk per account tampilkan yang current

		// Aggregate metrics
		Metrics: params.Metrics{
			GMV: NewMetric(currGMV, prevGMV),
			Ads: NewMetric(currAds, prevAds),
			// CommissionPaid:    NewMetric(currCommissionPaid, prevCommissionPaid),
			// CommissionPending: NewMetric(currCommissionPending, prevCommissionPending),
			Commission: NewMetric((currCommissionPaid + currCommissionPending), (prevCommissionPaid + prevCommissionPending)),
			Income:     NewMetric(currIncome, prevIncome),
		},
	}

	return result, nil
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
