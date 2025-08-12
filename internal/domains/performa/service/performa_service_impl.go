package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/performa/params"
	"github.com/royhairul/live-studio-api/internal/domains/performa/repository"

	hostservice "github.com/royhairul/live-studio-api/internal/domains/host/service"

	attendanceparams "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
)

type PerformaServiceImpl struct {
	repository        repository.PerformaRepository
	hostSvc           hostservice.HostService
	attendanceSvc     attendanceservice.AttendanceService
	accountSessionSvc accountsessionservice.AccountsessionService
	accountSvc        accountservice.AccountService
	studioSvc         studioservice.StudioService
}

func NewPerformaService(
	repository repository.PerformaRepository,
	hostSvc hostservice.HostService,
	attendanceSvc attendanceservice.AttendanceService,
	accountSessionSvc accountsessionservice.AccountsessionService,
	accountSvc accountservice.AccountService,
	studioSvc studioservice.StudioService,
) PerformaService {
	return &PerformaServiceImpl{repository, hostSvc, attendanceSvc, accountSessionSvc, accountSvc, studioSvc}
}

// GetHosts implements PerformaService.
func (p *PerformaServiceImpl) GetHosts(startTime string, endTime string) ([]*params.PerformaHostResponse, error) {
	if len(startTime) == 10 { // format: "YYYY-MM-DD"
		startTime += " 00:00:00"
	}
	if len(endTime) == 10 {
		endTime += " 23:59:59"
	}

	start, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid startTime format: %v", err)
	}
	end, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid endTime format: %v", err)
	}

	// Mencari attendance berdasarkan startTime and endTime
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	// Mengambil semua data host terlebih dahulu
	hosts, err := p.hostSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.PerformaHostResponse

	for _, host := range hosts {
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
		var detailList []params.PerformaHostDetailResponse

		for _, att := range hostAttendances {
			accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
			if err != nil {
				return nil, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
			}

			for _, session := range accountSessions {
				if session.CheckIn == nil || session.CheckOut == nil {
					continue
				}

				// Hitung durasi
				duration := session.CheckOut.Sub(*session.CheckIn)
				durationMillis := duration.Milliseconds()

				// Hitung sales dan paid
				sales := session.GMVSalesEnd - session.GMVSalesStart
				paid := session.GMVPaidEnd - session.GMVPaidStart

				// Tambahkan detail
				detailList = append(detailList, params.PerformaHostDetailResponse{
					AccountName: session.AccountName,
					Duration:    time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(duration),
					Sales:       sales,
					Paid:        paid,
				})

				totalDuration += durationMillis
				totalSales += sales
				totalPaid += paid
			}
		}

		var avgSales uint
		var avgPaid uint
		count := len(detailList)
		if count > 0 {
			avgSales = totalSales / uint(count)
			avgPaid = totalPaid / uint(count)
		}

		results = append(results, &params.PerformaHostResponse{
			ID:            host.ID.String(),
			Name:          host.Name,
			TotalDuration: totalDuration,
			TotalSales:    totalSales,
			AvgSales:      avgSales,
			AvgPaid:       avgPaid,
			List:          detailList,
		})
	}

	return results, nil
}

// GetHostByID implements PerformaService.
func (p *PerformaServiceImpl) GetHostByID(id string, startTime string, endTime string) (*params.PerformaHostResponse, error) {
	// Get attendances by date range
	if len(startTime) == 10 { // format: "YYYY-MM-DD"
		startTime += " 00:00:00"
	}
	if len(endTime) == 10 {
		endTime += " 23:59:59"
	}

	start, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid startTime format: %v", err)
	}

	end, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid endTime format: %v", err)
	}

	// Mencari attendance berdasarkan startTime and endTime
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
	var detailList []params.PerformaHostDetailResponse

	for _, att := range hostAttendances {
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			// Hitung durasi
			duration := session.CheckOut.Sub(*session.CheckIn)
			durationMillis := duration.Milliseconds()

			// Hitung sales dan paid
			sales := session.GMVSalesEnd - session.GMVSalesStart
			paid := session.GMVPaidEnd - session.GMVPaidStart

			// Tambahkan detail
			detailList = append(detailList, params.PerformaHostDetailResponse{
				AccountName: session.AccountName,
				Duration:    time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(duration),
				Sales:       sales,
				Paid:        paid,
			})

			totalDuration += durationMillis
			totalSales += sales
			totalPaid += paid
		}
	}

	var avgSales uint
	var avgPaid uint
	count := len(detailList)
	if count > 0 {
		avgSales = totalSales / uint(count)
		avgPaid = totalPaid / uint(count)
	}

	result := &params.PerformaHostResponse{
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

// GetStudios implements PerformaService.
func (p *PerformaServiceImpl) GetStudios(startTime string, endTime string) ([]*params.PerformaStudioResponse, error) {
	panic("unimplemented")
}

// GetStudioByID implements PerformaService.
func (p *PerformaServiceImpl) GetStudioByID(id string, startTime string, endTime string) (*params.PerformaStudioResponse, error) {
	panic("unimplemented")
}
