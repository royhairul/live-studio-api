package service

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/royhairul/live-studio-api/helpers/timehandler"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"

	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	accountsessionparams "github.com/royhairul/live-studio-api/internal/domains/accountsession/params"
	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	hostrepository "github.com/royhairul/live-studio-api/internal/domains/host/repository"
	liveservice "github.com/royhairul/live-studio-api/internal/domains/live/service"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	schedulerepository "github.com/royhairul/live-studio-api/internal/domains/schedule/repository"

	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
)

type AttendanceServiceImpl struct {
	repository        repository.AttendanceRepository
	hostRepo          hostrepository.HostRepository
	scheduleRepo      schedulerepository.ScheduleRepository
	accountSvc        accountservice.AccountService
	accountSessionSvc accountsessionservice.AccountsessionService
	liveSvc           liveservice.LiveService
	shopeeSvc         shopeeservice.ShopeeLiveService
}

func NewAttendanceService(
	repository repository.AttendanceRepository,
	hostRepo hostrepository.HostRepository,
	scheduleRepo schedulerepository.ScheduleRepository,
	accountSvc accountservice.AccountService,
	accountSessionSvc accountsessionservice.AccountsessionService,
	liveSvc liveservice.LiveService,
	shopeeSvc shopeeservice.ShopeeLiveService,
) AttendanceService {
	return &AttendanceServiceImpl{repository, hostRepo, scheduleRepo, accountSvc, accountSessionSvc, liveSvc, shopeeSvc}
}

func (s *AttendanceServiceImpl) FindAll() ([]*params.AttendanceResponse, error) {
	var results []*params.AttendanceResponse

	attendances, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	for _, attendance := range attendances {
		results = append(results, params.NewAttendanceResponse(attendance))
	}

	return results, nil
}

// FindAllByHostID implements AttendanceService.
func (s *AttendanceServiceImpl) FindAllByHostID(id string) ([]*params.AttendanceResponse, error) {
	panic("unimplemented")
}

// FindUncheckedOut implements AttendanceService.
func (s *AttendanceServiceImpl) FindUncheckedOut() ([]*params.AttendanceResponse, error) {
	var results []*params.AttendanceResponse

	attendances, err := s.repository.FindUncheckedOutByHost()
	if err != nil {
		return nil, err
	}

	for _, attendance := range attendances {
		results = append(results, &params.AttendanceResponse{
			ID:       attendance.ID,
			HostID:   *attendance.Host.ID,
			Name:     attendance.Host.Name,
			Date:     attendance.Date,
			CheckIn:  attendance.CheckedInAt,
			CheckOut: attendance.CheckedOutAt,

			ShiftID:   attendance.ShiftID,
			ShiftName: attendance.Shift.Name,

			Note: attendance.Note,
		})
	}

	return results, nil
}

// FindByDateRange implements AttendanceService.
func (s *AttendanceServiceImpl) FindByDateRange(startTime *time.Time, endTime *time.Time) ([]*params.AttendanceResponse, error) {
	attendances, err := s.repository.FindAllByDateRange(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attendances: %w", err)
	}

	var results []*params.AttendanceResponse
	for _, attendance := range attendances {
		results = append(results, params.NewAttendanceResponse(attendance))
	}

	return results, nil
}

func (s *AttendanceServiceImpl) CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceResponse, error) {
	parsedDate, err := timehandler.ParseDate(req.Date)
	if err != nil {
		return nil, err
	}

	host, err := s.hostRepo.FindByID(req.HostID)
	if err != nil {
		return nil, err
	}

	existAttendance, err := s.repository.FindUncheckedOutByStudio(fmt.Sprintf("%d", req.StudioID), parsedDate)
	if err == nil && existAttendance != nil {
		if existAttendance.HostID != nil && *existAttendance.HostID == *host.ID {
			return nil, fmt.Errorf("host %s already checkin", host.Name)
		} else {
			log.Print("studio founded, and there's a host")
			_, err := s.CheckOut(params.AttendanceCheckOutRequest{ID: existAttendance.ID})
			if err != nil {
				return nil, fmt.Errorf("failed to auto-checkout previous host: %w", err)
			}
		}
	}

	note := s.GenerateNote(nil, *parsedDate, req.ShiftID)

	attendance := entity.Attendance{
		Date:        parsedDate,
		ShiftID:     req.ShiftID,
		HostID:      host.ID,
		StudioID:    req.StudioID,
		CheckedInAt: timehandler.TimeNow(),
		Status:      "present",
		Note:        note,
	}

	created, err := s.repository.Create(&attendance)
	if err != nil {
		return nil, err
	}

	// Create record for account session
	accounts, err := s.accountSvc.FindByStudio(strconv.FormatUint(uint64(req.StudioID), 10))
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {

		live, err := s.shopeeSvc.GetShopeeLiveRealTime(account.Cookie)
		if err != nil {
			log.Printf("Failed to get live data for account %s: %v", account.Name, err)
			continue
		}

		accountSessionReq := accountsessionparams.CreateAccountsessionRequest{
			AccountID:    account.ID,
			AttendanceID: created.ID,
		}

		if len(live) == 0 {
			log.Printf("No live data for account %s", account.Name)
			accountSessionReq.GMVSalesStart = 0
			accountSessionReq.GMVPaidStart = 0

		} else {
			accountSessionReq.GMVPaidStart = uint(live[0].PlacedSales)
			accountSessionReq.GMVSalesStart = uint(live[0].ConfirmedSales)
		}

		_, err = s.accountSessionSvc.Create(accountSessionReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create account session for account %s: %w", account.Name, err)
		}
	}

	result := params.NewAttendanceResponse(created)
	return result, nil
}

func (s *AttendanceServiceImpl) CheckOut(req params.AttendanceCheckOutRequest) (*params.AttendanceResponse, error) {
	attendance, err := s.repository.FindByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("Gagal menemukan attendance ID %d", req.ID)
	}

	attendance.CheckedOutAt = timehandler.TimeNow()
	if err := s.repository.Save(attendance); err != nil {
		return nil, fmt.Errorf("Gagal menyimpan attendance ID %d", req.ID)
	}

	// Update account session with checkout time
	accountSessions, err := s.accountSessionSvc.FindAllByAttendanceID(strconv.FormatUint(uint64(req.ID), 10))
	if err != nil {
		return nil, err
	}

	for _, session := range accountSessions {
		account, err := s.accountSvc.FindById(strconv.FormatUint(uint64(session.AccountID), 10))
		if err != nil {
			log.Printf("Failed to get account ID %d: %v", session.AccountID, err)
			continue
		}

		live, err := s.shopeeSvc.GetShopeeLiveRealTime(account.Cookie)
		if err != nil {
			log.Printf("Failed to get live data for account %s: %v", account.Name, err)
			continue
		}

		var updateReq accountsessionparams.UpdateEndSessionRequest

		if len(live) > 0 {
			updateReq.GMVSalesEnd = uint(live[0].ConfirmedSales)
			updateReq.GMVPaidEnd = uint(live[0].PlacedSales)
		} else {
			updateReq.GMVSalesEnd = 0
			updateReq.GMVPaidEnd = 0
		}

		_, err = s.accountSessionSvc.UpdateEndSession(strconv.FormatUint(uint64(session.ID), 10), updateReq)
		if err != nil {
			return nil, fmt.Errorf("failed to update account session for attendance ID %d: %w", req.ID, err)
		}
	}

	result := params.NewAttendanceResponse(attendance)
	return result, nil
}

func (s *AttendanceServiceImpl) GenerateNote(schedule *scheduleentity.Schedule, attendanceDate time.Time, shiftID uint) string {
	if schedule == nil || schedule.ID == 0 {
		return "Tidak ada jadwal"
	}

	expectedDate := schedule.Date.Format("2006-01-02")
	actualDate := attendanceDate.Format("2006-01-02")

	dateMatch := expectedDate == actualDate
	shiftMatch := schedule.ShiftID == shiftID

	if dateMatch && shiftMatch {
		return "sesuai jadwal"
	}
	if dateMatch && !shiftMatch {
		return "shift tidak sesuai"
	}
	if !dateMatch && shiftMatch {
		return "tanggal tidak sesuai"
	}
	return "tanggal dan shift tidak sesuai"
}
