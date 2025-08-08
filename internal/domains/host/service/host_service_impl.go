package service

import (
	"fmt"
	"log"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"

	"github.com/royhairul/live-studio-api/internal/domains/host/params"
	"github.com/royhairul/live-studio-api/internal/domains/host/repository"

	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceparams "github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
)

type HostServiceImpl struct {
	repository        repository.HostRepository
	attendanceSvc     attendanceservice.AttendanceService
	accountSessionSvc accountsessionservice.AccountsessionService
}

func NewHostService(
	repository repository.HostRepository,
	attendanceSvc attendanceservice.AttendanceService,
	accountSessionSvc accountsessionservice.AccountsessionService,
) HostService {
	return &HostServiceImpl{repository, attendanceSvc, accountSessionSvc}
}

// Create implements HostService.
func (h *HostServiceImpl) Create(req params.CreateHostRequest) (*params.HostResponse, error) {
	host := entity.Host{
		Name:     req.Name,
		Phone:    req.Phone,
		StudioID: req.StudioID,
	}

	created, err := h.repository.Create(&host)
	if err != nil {
		return nil, err
	}

	result := params.NewHostResponse(created)

	return result, nil
}

// Update implements HostService.
func (h *HostServiceImpl) Update(id string, hostReq params.UpdateHostRequest) (*params.HostResponse, error) {
	host, err := h.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if hostReq.Name != nil {
		host.Name = *hostReq.Name
	}

	if hostReq.Phone != nil {
		host.Phone = *hostReq.Phone
	}

	if hostReq.StudioID != nil {
		host.StudioID = *hostReq.StudioID
		host.Studio = studioentity.Studio{}
	}

	saved, err := h.repository.Update(host)
	if err != nil {
		return nil, fmt.Errorf("failed to update host : %w", err)
	}

	log.Println(saved)

	result := params.NewHostResponse(saved)

	return result, nil
}

// FindAll implements HostService.
func (h *HostServiceImpl) FindAll() ([]*params.HostResponse, error) {
	hosts, err := h.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.HostResponse
	for _, host := range hosts {
		results = append(results, params.NewHostResponse(host))
	}

	return results, nil
}

// FindByID implements HostService.
func (h *HostServiceImpl) FindByID(id string) (*params.HostResponse, error) {
	host, err := h.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewHostResponse(host)

	return result, nil
}

// Delete implements HostService.
func (h *HostServiceImpl) Delete(id string) error {
	if err := h.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// FindAllGroupedByStudio implements HostService.
func (h *HostServiceImpl) FindAllGroupedByStudio() ([]*params.HostGroupedByStudioResponse, error) {
	hosts, err := h.repository.FindAll()
	if err != nil {
		return nil, err
	}

	groupMap := make(map[uint]*params.HostGroupedByStudioResponse)
	for _, host := range hosts {
		group, exists := groupMap[host.StudioID]
		if !exists {
			group = &params.HostGroupedByStudioResponse{
				StudioID:   host.StudioID,
				StudioName: host.Studio.Name,
				Hosts:      []params.HostResponse{},
			}
			groupMap[host.StudioID] = group
		}

		group.Hosts = append(group.Hosts, *params.NewHostResponse(host))
	}

	var results []*params.HostGroupedByStudioResponse
	for _, group := range groupMap {
		results = append(results, group)
	}

	return results, nil
}

// FindAllPerform implements HostService.
func (h *HostServiceImpl) FindAllPerform(startTime string, endTime string) ([]*params.HostPerformaResponse, error) {
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
	attendances, err := h.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	// Mengambil semua data host terlebih dahulu
	hosts, err := h.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.HostPerformaResponse

	for _, host := range hosts {
		// Get attendances by Host ID
		var hostAttendances []attendanceparams.AttendanceResponse
		for _, att := range attendances {
			if att.HostID == *host.ID {
				hostAttendances = append(hostAttendances, *att)
			}
		}

		var totalDuration int64
		var totalSales uint
		var totalPaid uint
		var detailList []params.HostPerformaDetailResponse

		for _, att := range hostAttendances {
			accountSessions, err := h.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
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
				detailList = append(detailList, params.HostPerformaDetailResponse{
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

		results = append(results, &params.HostPerformaResponse{
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

// FindByIDPerform implements HostService.
func (h *HostServiceImpl) FindByIDPerform(id string, startTime string, endTime string) (*params.HostPerformaResponse, error) {
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
	attendances, err := h.attendanceSvc.FindByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %v", err)
	}

	host, err := h.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Get attendances by Host ID
	var hostAttendances []attendanceparams.AttendanceResponse
	for _, att := range attendances {
		if att.HostID == *host.ID {
			hostAttendances = append(hostAttendances, *att)
		}
	}

	var totalDuration int64
	var totalSales uint
	var totalPaid uint
	var detailList []params.HostPerformaDetailResponse

	for _, att := range hostAttendances {
		accountSessions, err := h.accountSessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
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
			detailList = append(detailList, params.HostPerformaDetailResponse{
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

	result := &params.HostPerformaResponse{
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
