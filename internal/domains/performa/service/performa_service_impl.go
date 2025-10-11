package service

import (
	"fmt"
	"log"

	performaagg "github.com/royhairul/live-studio-api/internal/aggregator/performa"
	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
	"github.com/royhairul/live-studio-api/internal/domains/performa/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"

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
	aggregator        performaagg.PerformaAggregator
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
	aggregator performaagg.PerformaAggregator,
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
		aggregator,
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
	attendances, err := p.attendanceSvc.WithDateRange(*start, *end).FindAll()
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

			accountSessions, err := p.accountSessionSvc.WithAttendanceID(fmt.Sprintf("%d", att.ID)).FindAll()
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
	attendances, err := p.attendanceSvc.WithDateRange(*start, *end).FindAll()
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
		accountSessions, err := p.accountSessionSvc.WithAttendanceID(fmt.Sprintf("%d", att.ID)).FindAll()
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
func (p *PerformaServiceImpl) GetAccounts(startDate, endDate string) (*params.PerformaAccountResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	// Hitung durasi (days) dan periode sebelumnya
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// === Current Period ===
	currList, currTotal, err := p.aggregator.Calculate(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to build current performa list: %w", err)
	}

	// === Previous Period ===
	_, prevTotal, err := p.aggregator.Calculate(start, end)
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
	}

	return results, nil
}

// GetStudios implements PerformaService.
func (p *PerformaServiceImpl) GetStudios(startDate string, endDate string) (*params.PerformaStudioResponse, error) {
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous period dihitung mundur dengan panjang hari yang sama
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get All Studio
	studios, err := p.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var (
		currGMV, prevGMV                             int64
		currCommissionPaid, prevCommissionPaid       int64
		currCommissionPending, prevCommissionPending int64
		currAds, prevAds                             int64
		currIncome, prevIncome                       int64
	)

	list := []params.PerformaStudioItemResponse{}
	for _, studio := range studios {
		_, currTotal, err := p.aggregator.CalculateByStudio(fmt.Sprint(studio.ID), start, end)
		if err != nil {
			return nil, err
		}

		_, prevTotal, err := p.aggregator.CalculateByStudio(fmt.Sprint(studio.ID), &prevStart, &prevEnd)
		if err != nil {
			return nil, err
		}

		item := params.PerformaStudioItemResponse{
			StudioID:   fmt.Sprint(studio.ID),
			StudioName: studio.Name,
			GMV:        currTotal.GMV,
			Commission: currTotal.CommissionPaid + currTotal.CommissionPending,
			Ads:        currTotal.Ads,
			Income:     currTotal.Income,
		}

		list = append(list, item)

		// Calculate metrics
		currGMV += currTotal.GMV
		prevGMV += prevTotal.GMV
		currCommissionPaid += currTotal.CommissionPaid
		prevCommissionPaid += prevTotal.CommissionPaid
		currCommissionPending += currTotal.CommissionPending
		prevCommissionPending += prevTotal.CommissionPending
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
	start, end, err := timehandler.ParseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// previous range
	days := int(end.Sub(*start).Hours()/24) + 1
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.AddDate(0, 0, -days+1)

	// Get Studio by ID
	studio, err := p.studioSvc.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Current Performa
	currList, currTotal, err := p.aggregator.CalculateByStudio(fmt.Sprint(studio.ID), start, end)
	if err != nil {
		return nil, err
	}

	// Previous Performa
	_, prevTotal, err := p.aggregator.CalculateByStudio(fmt.Sprint(studio.ID), &prevStart, &prevEnd)
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
			Commission: NewMetric((currTotal.CommissionPaid + currTotal.CommissionPending), (prevTotal.CommissionPaid + prevTotal.CommissionPending)),
			Income:     NewMetric(currTotal.Income, prevTotal.Income),
		},
	}

	return result, nil
}
