package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"
	"github.com/royhairul/live-studio-api/models"
)

type AttendanceServiceImpl struct {
	repository repository.AttendanceRepository
}

func NewAttendanceService(repository repository.AttendanceRepository) AttendanceService {
	return &AttendanceServiceImpl{repository}
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

	for _, hostID := range req.HostIDs {
		var note string
		var scheduleID *uint

		schedule, err := repoWithTx.FindScheduleByHostShiftAndDate(hostID, req.ShiftID, req.Date)
		hostName := fmt.Sprintf("Host ID %d", hostID) // default fallback host name
		if schedule != nil && schedule.Host != (models.Host{}) {
			hostName = schedule.Host.Name
		}

		// Jika ada schedule, cek duplikat attendance
		if err == nil && schedule != nil {
			scheduleID = &schedule.ID
			note = s.GenerateNote(schedule, req.Date, req.ShiftID)

			_, err := s.repository.FindByScheduleID(schedule.ID)
			if err == nil {
				results = append(results, params.AttendanceCheckInResult{
					HostName: hostName,
					Message:  "Host sudah check-in sebelumnya",
					Status:   "failed",
				})
				failedCount++
				continue
			}
		} else {
			// Tidak ada schedule, tetap dibuat dengan note Tidak ada jadwal
			note = s.GenerateNote(nil, req.Date, req.ShiftID)
		}

		attendance := entity.Attendance{
			Date:    &req.Date,
			ShiftID: &req.ShiftID,

			ScheduleID: scheduleID,
			HostID:     &hostID,

			CheckedInAt: helpers.TimeNow(),

			Status: "present",
			Note:   note,
		}

		_, err = repoWithTx.Create(&attendance)
		if err != nil {
			results = append(results, params.AttendanceCheckInResult{
				HostName: hostName,
				Message:  "Gagal menyimpan attendance",
				Status:   "failed",
			})
			failedCount++
			continue
		}

		results = append(results, params.AttendanceCheckInResult{
			HostName: hostName,
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

func (s *AttendanceServiceImpl) GenerateNote(schedule *models.Schedule, attendanceDate time.Time, shiftID uint) string {
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
