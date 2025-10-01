package repository

import (
	"log"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/attendance/entity"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/params"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceRepositoryImpl struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{DB: db}
}

func (r *AttendanceRepositoryImpl) BuildQuery(filter params.AttendanceFilter) *gorm.DB {
	query := r.DB.Model(entity.Attendance{}).Preload(clause.Associations)

	if filter.AccountID != nil {
		query = query.Where("account_id", filter.AccountID)
	}

	if filter.HostID != nil {
		query = query.Where("host_id", filter.HostID)
	}

	if filter.ShiftID != nil {
		query = query.Where("shift_id", filter.ShiftID)
	}

	if filter.StudioID != nil {
		log.Println("Filter Studio ID:", *filter.StudioID)
		query = query.Where("studio_id", filter.StudioID)
	}

	if filter.StartTime != nil && filter.EndTime != nil {
		query = query.Where("date::date BETWEEN ? AND ?", filter.StartTime, filter.EndTime)
	}

	return query
}

func (r *AttendanceRepositoryImpl) FindAll(filter params.AttendanceFilter) ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	if err := r.BuildQuery(filter).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *AttendanceRepositoryImpl) FindOne(filter params.AttendanceFilter) (*entity.Attendance, error) {
	var attendance *entity.Attendance
	if err := r.BuildQuery(filter).First(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *AttendanceRepositoryImpl) Create(attendance *entity.Attendance) (*entity.Attendance, error) {
	if err := r.DB.Create(attendance).Error; err != nil {
		return nil, err
	}

	err := r.DB.
		Preload("Schedule").
		Preload("Host").
		Preload("Shift").
		Preload("Studio").
		First(&attendance).Error
	if err != nil {
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
	err := r.DB.
		Preload("Schedule").
		Preload("Host").
		Preload("Shift").
		Preload("Studio").
		First(&attendance, "id = ?", id).Error
	if err != nil {
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
		Preload("Studio").
		Where("checked_out_at IS NULL").
		Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	return attendances, nil
}

// FindUncheckedOutByStudio implements AttendanceRepository.
func (r *AttendanceRepositoryImpl) FindUncheckedOutByStudio(studioID string, date *time.Time) (*entity.Attendance, error) {
	var attendance entity.Attendance
	err := r.DB.Where("studio_id = ?", studioID).
		Where("date::date = ?", date).
		Where("checked_out_at IS NULL").
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepositoryImpl) FindByScheduleID(id uint) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := r.DB.Where("schedule_id = ?", id).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// FindByHostShiftAndDate implements AttendanceRepository.
func (r *AttendanceRepositoryImpl) FindByHostShiftAndDate(hostID string, shiftID uint, date time.Time) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := r.DB.Where("host_id = ? AND date = ? AND shift_id = ?", hostID, date, shiftID).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// FindAllByDateRange implements AttendanceRepository.
func (r *AttendanceRepositoryImpl) FindAllByDateRange(startTime *time.Time, endTime *time.Time) ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	err := r.DB.
		Preload("Schedule").
		Preload("Shift").
		Preload("Host").
		Preload("Studio").
		Where("date::date BETWEEN ? AND ?", startTime, endTime).
		Find(&attendances).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// FindAllByHostID implements AttendanceRepository.
func (r *AttendanceRepositoryImpl) FindAllByHostID(id string) ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance
	err := r.DB.
		Preload("Schedule").
		Preload("Shift").
		Preload("Host").
		Preload("Studio").
		Find(&attendances, "host_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// FindByAccountID implements AttendanceRepository.
func (r *AttendanceRepositoryImpl) FindByAccountID(id uint) (*entity.Attendance, error) {
	var attendance entity.Attendance
	err := r.DB.Where("account_id = ?", id).
		Where("checked_out_at IS NULL").
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}
