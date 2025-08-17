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
func (p *PerformaServiceImpl) GetHosts(startDate string, endDate string) ([]*params.PerformaHostResponse, error) {
	// Set default value
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
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
				durationMillis := duration.Milliseconds()

				totalDuration += durationMillis
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
func (p *PerformaServiceImpl) GetHostByID(id string, startTime string, endTime string) (*params.PerformaHostDetailResponse, error) {
	start, end, err := timehandler.ParseDateRange(startTime, endTime)
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
		log.Printf("attendance id: %d", att.ID)
		log.Printf("attendance name: %s", att.Name)
	}

	log.Printf("length of attendances: %d", len(hostAttendances))

	var totalDuration int64
	var totalSales uint
	var totalPaid uint
	var detailList []params.PerformaHostItemResponse

	for _, att := range hostAttendances {
		log.Printf("host att %d", att.ID)
		accountSessions, err := p.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to get account sessions for attendance %d: %w", att.ID, err)
		}

		log.Printf("count session: %d", len(accountSessions))

		for _, session := range accountSessions {
			if session.CheckIn == nil || session.CheckOut == nil {
				continue
			}

			detailList = append(detailList, params.PerformaHostItemResponse{
				AccountName: session.AccountName,
				Duration:    time.Time{},
				Sales:       session.GMVSales,
				Paid:        session.GMVPaid,
			})

			log.Printf("account: %s", session.AccountName)

			// Hitung durasi
			duration := session.CheckOut.Sub(*session.CheckIn)
			durationMillis := duration.Milliseconds()

			totalDuration += durationMillis
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

// GetStudios implements PerformaService.
func (p *PerformaServiceImpl) GetStudios(startTime string, endTime string) ([]*params.PerformaStudioResponse, error) {
	// Set default value
	if startTime == "" {
		startTime = time.Now().Format("2006-01-02")
	}
	if endTime == "" {
		endTime = time.Now().Format("2006-01-02")
	}

	start, end, err := timehandler.ParseDateRange(startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Get Attendances by date range
	attendances, err := p.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	// Get All Studio
	studios, err := p.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.PerformaStudioResponse

	for _, studio := range studios {
		var totalDuration int64
		var totalSales uint
		var totalPaid uint

		// Get attendances by studio ID
		for _, att := range attendances {
			if att.StudioID == studio.ID {
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

				totalDuration += att.Duration
				totalSales += session.GMVSales
				totalPaid += session.GMVPaid
			}
		}

		results = append(results, &params.PerformaStudioResponse{
			ID:            fmt.Sprintf("%d", studio.ID),
			Name:          studio.Name,
			TotalDuration: totalDuration,
			TotalSales:    totalSales,
		})
	}

	return results, nil
}

// GetStudioByID implements PerformaService.
func (p *PerformaServiceImpl) GetStudioByID(id string, startTime string, endTime string) (*params.PerformaStudioResponse, error) {
	panic("unimplemented")
}
