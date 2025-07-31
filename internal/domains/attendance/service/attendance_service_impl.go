package service

import (
	"fmt"
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

			// ShiftStartTime: attendance.Shift.StartTime,
			// ShiftEndTime:   attendance.Shift.EndTime,
			ShiftID:   attendance.ShiftID,
			ShiftName: attendance.Shift.Name,

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

			ShiftID:   attendance.ShiftID,
			ShiftName: attendance.Shift.Name,

			Note: attendance.Note,
		})
	}

	return results, nil
}

func (s *AttendanceServiceImpl) CheckIn(req params.AttendanceCheckInRequest) (*params.AttendanceResponse, error) {
	host, err := s.hostRepo.FindByID(req.HostID)
	if err != nil {
		return nil, err
	}

	// Cari attendance existing
	existingAttendance, err := s.repository.FindByHostShiftAndDate(req.HostID, req.ShiftID, req.Date)
	if err == nil && existingAttendance != nil {
		if existingAttendance.CheckedOutAt == nil {
			return nil, fmt.Errorf("host already checkout")
		}
		// else: sudah check-in & sudah checkout → boleh check-in lagi
	}

	note := s.GenerateNote(nil, req.Date, req.ShiftID)

	attendance := entity.Attendance{
		Date:        &req.Date,
		ShiftID:     req.ShiftID,
		HostID:      host.ID,
		StudioID:    req.StudioID,
		CheckedInAt: helpers.TimeNow(),
		Status:      "present",
		Note:        note,
	}

	created, err := s.repository.Create(&attendance)
	if err != nil {
		return nil, err
	}

	result := params.NewAttendanceResponse(created)
	return result, nil
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

// func (s *AttendanceServiceImpl) GenerateSummary(successCount, failedCount int, results []params.AttendanceCheckInResult) *params.AttendanceCheckInSummary {
// 	finalMessage := "Berhasil check-in semua host"
// 	if successCount == 0 {
// 		finalMessage = "Gagal check-in semua host"
// 	} else if failedCount > 0 {
// 		finalMessage = "Sebagian berhasil check-in"
// 	}

// 	return &params.AttendanceCheckInSummary{
// 		Message:      finalMessage,
// 		SuccessCount: successCount,
// 		FailedCount:  failedCount,
// 		Results:      results,
// 	}
// }
