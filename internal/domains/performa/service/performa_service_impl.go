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

func (p *PerformaServiceImpl) fetchAttendances(start, end *time.Time) ([]attendanceparams.AttendanceResponse, error) {
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, err
	}

	res := make([]attendanceparams.AttendanceResponse, 0, len(attendances))
	for _, att := range attendances {
		res = append(res, *att)
	}
	return res, nil
}

// GetAccounts implements PerformaService.
func (p *PerformaServiceImpl) GetAccounts(startDate, endDate string) (*params.PerformaAccountResponse, error) {
	// Default ke hari ini
	today := *timehandler.DateNow()
	if startDate == "" {
		startDate = today
	}
	if endDate == "" {
		endDate = today
	}

	// Parse range tanggal
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Hitung durasi (days) dan periode sebelumnya
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// === Ambil attendances ===
	currAttendances, err := p.fetchAttendances(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	prevAttendances, err := p.fetchAttendances(&prevStart, &prevEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous attendances: %w", err)
	}

	// === Current Period ===
	currList, currGMV, currAds, currCommissionPaid, currCommissionPending, currIncome, err := p.buildPerformaAccountDetailList(currAttendances, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to build current performa list: %w", err)
	}

	// === Previous Period ===
	_, prevGMV, prevAds, prevCommissionPaid, prevCommissionPending, prevIncome, err := p.buildPerformaAccountDetailList(prevAttendances, &prevStart, &prevEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to build previous performa list: %w", err)
	}

	// === Gabungkan current & previous list ===
	accountMap := make(map[uint]*params.PerformaAccountDetailItemResponse)

	// Masukkan data current
	for _, item := range currList {
		newItem := item
		accountMap[item.AccountID] = &newItem
	}

	// Convert map -> slice
	list := make([]params.PerformaAccountDetailItemResponse, 0, len(accountMap))
	for _, v := range accountMap {
		list = append(list, *v)
	}

	var totalCommission, totalGMV, totalAds, totalIncome int64
	for _, item := range list {
		totalCommission += item.Commission
		totalGMV += item.GMV
		totalAds += item.Ads
		totalIncome += item.Income
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
			Commission: NewMetric(currCommissionPaid+currCommissionPending, prevCommissionPaid+prevCommissionPending),
			GMV:        NewMetric(currGMV, prevGMV),
			Ads:        NewMetric(int64(currAds), int64(prevAds)),
			Income:     NewMetric(currIncome, prevIncome),
		},
		List: list,
	}

	return results, nil
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
		// Reset per studio
		var studioGMV, studioPrevGMV int64
		var studioCommissionPaid, studioPrevCommissionPaid int64
		var studioCommissionPending, studioPrevCommissionPending int64
		var studioAds, studioPrevAds uint
		var studioIncome, studioPrevIncome int64

		// Get Account in this studio
		accounts, err := p.accountSvc.FindByStudio(fmt.Sprintf("%d", studio.ID))
		if err != nil {
			return nil, err
		}

		for _, account := range accounts {
			// current transactions
			transactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", account.ID), start, end)
			if err != nil {
				return nil, err
			}
			for _, tx := range transactions {
				studioCommissionPaid += int64(tx.Commission.Paid)
				studioCommissionPending += int64(tx.Commission.Pending)
			}

			// previous transactions
			prevTransactions, err := p.transactionSvc.FindAllByDate(fmt.Sprintf("%d", account.ID), &prevStart, &prevEnd)
			if err != nil {
				return nil, err
			}
			for _, tx := range prevTransactions {
				studioPrevCommissionPaid += int64(tx.Commission.Paid)
				studioPrevCommissionPending += int64(tx.Commission.Pending)
			}

			// ads
			allAds, err := p.accountAdsSvc.FindByDateAndAccounts(start, end, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range allAds {
				studioAds += a.Ads
			}

			prevAllAds, err := p.accountAdsSvc.FindByDateAndAccounts(&prevStart, &prevEnd, fmt.Sprintf("%d", account.ID))
			if err != nil {
				return nil, err
			}
			for _, a := range prevAllAds {
				studioPrevAds += a.Ads
			}
		}

		// GMV current
		for _, att := range attendances {
			accountsessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, err
			}
			for _, session := range accountsessions {
				studioGMV += int64(session.GMVPaid)
			}
		}

		// GMV previous
		for _, att := range prevAttendances {
			accountsessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, err
			}
			for _, session := range accountsessions {
				studioPrevGMV += int64(session.GMVPaid)
			}
		}

		studioIncome = (studioCommissionPaid + studioCommissionPending) - int64(studioAds)
		studioPrevIncome = (studioPrevCommissionPaid + studioPrevCommissionPending) - int64(studioPrevAds)

		// Tambahkan ke list
		list = append(list, params.PerformaStudioItemResponse{
			StudioID:   fmt.Sprintf("%d", studio.ID),
			StudioName: studio.Name,
			Commission: studioCommissionPaid + studioCommissionPending,
			GMV:        studioGMV,
			Ads:        int64(studioAds),
			Income:     studioIncome,
		})

		// Akumulasi ke total metrics
		currGMV += studioGMV
		prevGMV += studioPrevGMV
		currCommissionPaid += studioCommissionPaid
		prevCommissionPaid += studioPrevCommissionPaid
		currCommissionPending += studioCommissionPending
		prevCommissionPending += studioPrevCommissionPending
		currAds += studioAds
		prevAds += studioPrevAds
		currIncome += studioIncome
		prevIncome += studioPrevIncome
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
