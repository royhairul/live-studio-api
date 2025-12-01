package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"

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
	options           params.AttendanceFilter
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
	return &AttendanceServiceImpl{
		repository:        repository,
		hostRepo:          hostRepo,
		scheduleRepo:      scheduleRepo,
		accountSvc:        accountSvc,
		accountSessionSvc: accountSessionSvc,
		liveSvc:           liveSvc,
		shopeeSvc:         shopeeSvc,
		options:           params.AttendanceFilter{},
	}
}

// WithAccountID implements AttendanceService.
func (s *AttendanceServiceImpl) WithAccountID(accountID string) AttendanceService {
	instance := *s
	instance.options.AccountID = &accountID
	return &instance
}

// WithDateRange implements AttendanceService.
func (s *AttendanceServiceImpl) WithDateRange(startTime time.Time, endTime time.Time) AttendanceService {
	instance := *s
	instance.options.StartTime = &startTime
	instance.options.EndTime = &endTime
	return &instance
}

// WithHostID implements AttendanceService.
func (s *AttendanceServiceImpl) WithHostID(hostID string) AttendanceService {
	instance := *s
	instance.options.HostID = &hostID
	return &instance
}

// WithStudioID implements AttendanceService.
func (s *AttendanceServiceImpl) WithStudioID(studioID string) AttendanceService {
	instance := *s
	instance.options.StudioID = &studioID
	return &instance
}

func (s *AttendanceServiceImpl) FindAll(ctx context.Context) ([]*params.AttendanceResponse, error) {
	attendances, err := s.repository.FindAll(ctx, s.options)
	if err != nil {
		return nil, err
	}

	var results []*params.AttendanceResponse
	for _, attendance := range attendances {
		results = append(results, params.NewAttendanceResponse(attendance))
	}

	return results, nil
}

// FindUncheckedOut implements AttendanceService.
func (s *AttendanceServiceImpl) FindUncheckedOut(ctx context.Context) ([]*params.AttendanceResponse, error) {
	var results []*params.AttendanceResponse

	attendances, err := s.repository.FindUncheckedOutByHost(ctx)
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

func (s *AttendanceServiceImpl) CheckIn(ctx context.Context, req params.AttendanceCheckInRequest) (*params.AttendanceResponse, error) {
	parsedDate, err := timehandler.ParseDate(*timehandler.DateNow())
	if err != nil {
		return nil, err
	}

	host, err := s.hostRepo.FindByID(ctx, req.HostID)
	if err != nil {
		return nil, err
	}

	existAttendance, err := s.repository.FindUncheckedOutByStudio(ctx, fmt.Sprintf("%d", req.StudioID), parsedDate)
	if err == nil && existAttendance != nil {
		if existAttendance.HostID != nil && *existAttendance.HostID == *host.ID {
			return nil, fmt.Errorf("host %s already checkin", host.Name)
		} else {
			_, err := s.CheckOut(ctx, params.AttendanceCheckOutRequest{ID: []uint{existAttendance.ID}})
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
		Status:      "active",
		Note:        note,
	}

	created, err := s.repository.Create(ctx, &attendance)
	if err != nil {
		return nil, err
	}

	// Create record for account session
	accounts, err := s.accountSvc.WithStudioID(fmt.Sprintf("%d", req.StudioID)).FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {

		live, err := s.shopeeSvc.GetLiveSessionRT(account.Cookie)
		if err != nil {
			log.Printf("Failed to get live data for account %s: %v", account.Name, err)
			continue
		}

		accountSessionReq := accountsessionparams.CreateAccountsessionRequest{
			AccountID:    account.ID,
			AttendanceID: created.ID,
			StudioID:     created.StudioID,
		}

		if len(live) == 0 {
			accountSessionReq.GMVSalesStart = 0
			accountSessionReq.GMVPaidStart = 0

		} else {
			accountSessionReq.GMVPaidStart = uint(live[0].PlacedSales)
			accountSessionReq.GMVSalesStart = uint(live[0].ConfirmedSales)
		}

		_, err = s.accountSessionSvc.Create(ctx, accountSessionReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create account session for account %s: %w", account.Name, err)
		}
	}

	result := params.NewAttendanceResponse(created)
	return result, nil
}

func (s *AttendanceServiceImpl) CheckOut(ctx context.Context, req params.AttendanceCheckOutRequest) ([]*params.AttendanceResponse, error) {
	var responses []*params.AttendanceResponse

	for _, id := range req.ID {
		attendance, err := s.repository.FindByID(ctx, id)
		if err != nil {
			log.Printf("Gagal menemukan attendance ID %d: %v", id, err)
			continue
		}

		attendance.CheckedOutAt = timehandler.TimeNow()
		attendance.Status = "inactive"

		if err := s.repository.Save(ctx, attendance); err != nil {
			log.Printf("Gagal menyimpan attendance ID %d: %v", id, err)
			continue
		}

		// Update account session with checkout time
		accountSessions, err := s.accountSessionSvc.
			WithAttendanceID(fmt.Sprintf("%d", attendance.ID)).
			FindAll(ctx)
		if err != nil {
			log.Printf("Failed to get account sessions for attendance ID %d: %v", id, err)
			continue
		}

		for _, session := range accountSessions {
			account, err := s.accountSvc.WithID(fmt.Sprintf("%d", session.AccountID)).FindOne(ctx)
			if err != nil {
				log.Printf("Failed to get account ID %d: %v", session.AccountID, err)
				continue
			}

			live, err := s.shopeeSvc.GetLiveSessionRT(account.Cookie)
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

			if _, err := s.accountSessionSvc.UpdateEndSession(
				ctx,
				strconv.FormatUint(uint64(session.ID), 10),
				updateReq,
			); err != nil {
				log.Printf("Failed to update account session for attendance ID %d: %v", id, err)
				continue
			}
		}

		responses = append(responses, params.NewAttendanceResponse(attendance))
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("tidak ada attendance yang berhasil di-checkout")
	}

	return responses, nil
}

func (s *AttendanceServiceImpl) GenerateNote(schedule *scheduleentity.Schedule, attendanceDate time.Time, shiftID uint) string {
	if schedule == nil || schedule.ID == 0 {
		return "Tidak ada jadwal"
	}

	expectedDate := schedule.Date.Format(constants.LayoutYYMMDD)
	actualDate := attendanceDate.Format(constants.LayoutYYMMDD)

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
