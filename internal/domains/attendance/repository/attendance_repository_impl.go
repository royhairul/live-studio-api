package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"gorm.io/gorm"
)

type AttendanceRepositoryImpl struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{DB: db}
}

func (r *AttendanceRepositoryImpl) WithTx(tx *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{
		DB: tx,
	}
}

func (r *AttendanceRepositoryImpl) BeginTransaction() *gorm.DB {
	return r.DB.Begin()
}

func (r *AttendanceRepositoryImpl) FindAll() ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	if err := r.DB.Preload("Host").Preload("Host.Studio").Preload("Shift").Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *AttendanceRepositoryImpl) Create(attendance *entity.Attendance) (*entity.Attendance, error) {
	if err := r.DB.Create(attendance).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Schedule").Preload("Host").Preload("Shift").First(attendance, attendance.ID).Error; err != nil {
		return nil, err
	}

	return attendance, nil
}

func (r *AttendanceRepositoryImpl) Save(attendance *entity.Attendance) error {
	if err := r.DB.Preload("Schedule").Save(attendance).Error; err != nil {
		return err
	}

	return nil
}

func (r *AttendanceRepositoryImpl) Delete(id string) error {
	if err := r.DB.Where("id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *AttendanceRepositoryImpl) FindByID(id uint) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := r.DB.Where("id = ?", id).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepositoryImpl) FindUncheckedOutByHost() ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	err := r.DB.
		Preload("Schedule").
		Preload("Shift").
		Preload("Host").
		Where("checked_out_at IS NULL").
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *AttendanceRepositoryImpl) FindByScheduleID(id uint) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := r.DB.Where("schedule_id = ?", id).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}
