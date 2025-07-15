package schedule

import (
	"fmt"

	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/helpers"
	"github.com/royhairul/live-studio-api/models"
)

func GetShiftAll() ([]models.ScheduleShift, error) {
	var shifts []models.ScheduleShift
	if err := database.DB.Find(&shifts).Error; err != nil {
		return nil, err
	}

	return shifts, nil
}

func CreateShift(dto *dto.CreateShiftDTO) error {
	startTime, err := helpers.ParseTime(dto.StartTime)
	if err != nil {
		return fmt.Errorf("error start time: %w", err)
	}

	endTime, err := helpers.ParseTime(dto.EndTime)
	if err != nil {
		return fmt.Errorf("error end time: %w", err)
	}

	shift := &models.ScheduleShift{
		Name:      dto.Name,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := database.DB.Create(shift).Error; err != nil {
		return err
	}

	return nil
}

func GetShiftByID(id string) (*models.ScheduleShift, error) {
	var shift models.ScheduleShift
	if err := database.DB.Where("id = ?", id).First(&shift).Error; err != nil {
		return nil, fmt.Errorf("shift not found: %w", err)
	}

	return &shift, nil
}

func UpdateShift(id string, dto *dto.UpdateShiftDTO) error {
	shift, err := GetShiftByID(id)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	startTime, _ := helpers.ParseTime(*dto.StartTime)
	endTime, _ := helpers.ParseTime(*dto.EndTime)

	shift.Name = *dto.Name
	shift.StartTime = startTime
	shift.EndTime = endTime

	database.DB.Updates(&shift)

	return nil
}

func DeleteShift(id string) error {
	shift, err := GetShiftByID(id)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if err := database.DB.Delete(&shift).Error; err != nil {
		return fmt.Errorf("failed to delete shift: %w", err)
	}

	return nil
}
