package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"
	hostrepository "github.com/royhairul/live-studio-api/internal/domains/host/repository"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	schedulerepository "github.com/royhairul/live-studio-api/internal/domains/schedule/repository"
)

type AttendanceServiceImpl struct {
	repository   repository.AttendanceRepository
	hostRepo     hostrepository.HostRepository
	scheduleRepo schedulerepository.ScheduleRepository
}

func NewAttendanceService(repository repository.AttendanceRepository, hostRepo hostrepository.HostRepository, scheduleRepo schedulerepository.ScheduleRepository) AttendanceService {
	return &AttendanceServiceImpl{repository, hostRepo, scheduleRepo}
}

func (s *AttendanceServiceImpl) FindAll() ([]*params.AttendanceResponse, error) {
	var results []*params.AttendanceResponse

	attendances, err := s.repository.FindAll()
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

			ShiftStartTime: attendance.Shift.StartTime,
			ShiftEndTime:   attendance.Shift.EndTime,

			Note: attendance.Note,
		})
	}

	return results, nil
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

			ShiftStartTime: attendance.Shift.StartTime,
			ShiftEndTime:   attendance.Shift.EndTime,

			Note: attendance.Note,
		})
	}

	return results, nil
}

func (s *AttendanceServiceImpl) CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceCheckInSummary, error) {
	var results []params.AttendanceCheckInResult
	successCount := 0
	failedCount := 0

	tx := s.repository.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	repoWithTx := s.repository.WithTx(tx)

	parsedShiftID, err := strconv.ParseUint(req.ShiftID, 10, 0)
	if err != nil {
		// handle error parsing, misalnya return error ke caller
		return nil, fmt.Errorf("invalid shift ID: %v", err)
	}

	for _, hostID := range req.HostIDs {

		host, err := s.hostRepo.FindByID(hostID.String())
		if err != nil {
			return nil, err
		}

		// Cari attendance existing
		existingAttendance, err := s.repository.FindByHostShiftAndDate(hostID, req.ShiftID, req.Date)
		if err == nil && existingAttendance != nil {
			if existingAttendance.CheckedOutAt == nil {
				// Sudah check-in & belum checkout → tolak check-in
				results = append(results, params.AttendanceCheckInResult{
					HostName: host.Name,
					Message:  fmt.Sprintf("Host %s sudah check-in dan belum checkout", host.Name),
					Status:   "failed",
				})
				failedCount++
				continue
			}
			// else: sudah check-in & sudah checkout → boleh check-in lagi
		}

		note := s.GenerateNote(nil, req.Date, uint(parsedShiftID))

		attendance := entity.Attendance{
			Date:    &req.Date,
			ShiftID: uint(parsedShiftID),

			HostID: &hostID,

			CheckedInAt: helpers.TimeNow(),

			Status: "present",
			Note:   note,
		}

		_, err = repoWithTx.Create(&attendance)
		if err != nil {
			results = append(results, params.AttendanceCheckInResult{
				HostName: host.Name,
				Message:  "Gagal menyimpan attendance",
				Status:   "failed",
			})
			failedCount++
			continue
		}

		results = append(results, params.AttendanceCheckInResult{
			HostName: host.Name,
			Message:  "Berhasil check-in",
			Status:   "success",
		})
		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.GenerateSummary(successCount, failedCount, results), nil
}

func (s *AttendanceServiceImpl) CheckOut(req params.AttendanceCheckOutRequest) error {
	var errs []error
	for _, ids := range req.AttendanceIDs {
		attendance, err := s.repository.FindByID(ids)
		if err != nil {
			errs = append(errs, fmt.Errorf("Gagal menemukan attendance ID %d", ids))
			continue
		}

		attendance.CheckedOutAt = helpers.TimeNow()
		if err := s.repository.Save(attendance); err != nil {
			errs = append(errs, fmt.Errorf("Gagal menyimpan attendance ID %d", ids))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Terdapat %d error saat check-out: %v", len(errs), errs)
	}

	return nil
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

func (s *AttendanceServiceImpl) GenerateSummary(successCount, failedCount int, results []params.AttendanceCheckInResult) *params.AttendanceCheckInSummary {
	finalMessage := "Berhasil check-in semua host"
	if successCount == 0 {
		finalMessage = "Gagal check-in semua host"
	} else if failedCount > 0 {
		finalMessage = "Sebagian berhasil check-in"
	}

	return &params.AttendanceCheckInSummary{
		Message:      finalMessage,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
	}
}
