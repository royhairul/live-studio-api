package schedule

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
)

func GenerateHostSchedule() {}

func GetHostScheduleAll() ([]models.Schedule, error) {
	var schedules []models.Schedule

	if err := database.DB.
		Preload("Host").
		Preload("Host.Studio"). // jika ingin memuat studio juga
		Preload("Shift").
		Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve host schedules: %w", err)
	}

	return schedules, nil
}

func GetHostScheduleByID(id string) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := database.DB.Preload("Host").Preload("Shift").First(&schedule, id).Error; err != nil {
		return nil, fmt.Errorf("host schedule not found: %w", err)
	}
	return &schedule, nil
}

func CreateHostSchedule(dto *dto.CreateHostScheduleDTO) (*models.Schedule, error) {
	// Ambil shift dari database
	var shift models.ScheduleShift
	if err := database.DB.Where("id = ?", dto.ShiftID).First(&shift).Error; err != nil {
		return nil, fmt.Errorf("shift not found: %w", err)
	}

	// Validasi host ID jika perlu
	var host models.Host
	if err := database.DB.Where("id = ?", dto.HostID).First(&host).Error; err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	// Ubah string ke time.Time
	dateParsed, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	// Buat schedule baru
	hostSchedule := &models.Schedule{
		HostID:    dto.HostID,
		ShiftID:   dto.ShiftID,
		Date:      dateParsed,
		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
	}

	// Simpan ke database
	if err := database.DB.Create(hostSchedule).Error; err != nil {
		return nil, fmt.Errorf("failed to create host schedule: %w", err)
	}

	// Ambil data lengkap dengan relasi Host dan Shift
	if err := database.DB.Preload("Host").Preload("Shift").First(hostSchedule, hostSchedule.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load schedule with relations: %w", err)
	}

	return hostSchedule, nil
}

func UpdateHostSchedule(id string, dto *dto.UpdateHostScheduleDTO) error {
	var schedule models.Schedule
	if err := database.DB.Where("id = ?", id).First(&schedule).Error; err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	// Ambil shift baru
	var shift models.ScheduleShift
	if err := database.DB.Where("id = ?", dto.ShiftID).First(&shift).Error; err != nil {
		return fmt.Errorf("shift not found: %w", err)
	}

	// Ubah format tanggal
	layout := "2006-01-02"
	dateParsed, err := time.Parse(layout, *dto.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	// Update field pada schedule
	schedule.Date = dateParsed
	schedule.StartTime = shift.StartTime
	schedule.EndTime = shift.EndTime

	if err := database.DB.Save(&schedule).Error; err != nil {
		return fmt.Errorf("failed to update host schedule: %w", err)
	}

	return nil
}

func DeleteHostSchedule(id string) error {
	var schedule models.Schedule
	if err := database.DB.Where("id = ?", id).First(&schedule).Error; err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	if err := database.DB.Delete(&schedule).Error; err != nil {
		return fmt.Errorf("failed to delete host schedule: %w", err)
	}

	return nil
}

func GetHostScheduleByHostID(hostID string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := database.DB.Where("host_id = ?", hostID).Preload("Shift").Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve host schedule by host ID: %w", err)
	}

	if len(schedules) == 0 {
		return nil, fmt.Errorf("no schedules found for host ID: %s", hostID)
	}

	return schedules, nil
}

func SwitchHostSchedule(dto *dto.SwitchHostScheduleDTO) error {
	var fromSchedule, toSchedule models.Schedule

	// Find the schedules for both hosts
	if err := database.DB.Where("id = ?", dto.FromScheduleID).First(&fromSchedule).Error; err != nil {
		return fmt.Errorf("from schedule not found: %w", err)
	}
	if err := database.DB.Where("id = ?", dto.ToScheduleID).First(&toSchedule).Error; err != nil {
		return fmt.Errorf("to schedule not found: %w", err)
	}

	// Swap the host IDs
	fromSchedule.HostID = dto.ToHostID
	toSchedule.HostID = dto.FromHostID

	// Save the changes
	if err := database.DB.Save(&fromSchedule).Error; err != nil {
		return fmt.Errorf("failed to update from schedule: %w", err)
	}
	if err := database.DB.Save(&toSchedule).Error; err != nil {
		return fmt.Errorf("failed to update to schedule: %w", err)
	}

	return nil
}
