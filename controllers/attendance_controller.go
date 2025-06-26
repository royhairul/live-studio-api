package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
)

func AttendanceCheckIn(c *gin.Context) {
	var req dto.CheckAttendanceRequestDTO

	// 1. Binding request body ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		// response := helpers.BuildErrorResponse("Failed to process request", err.Error(), helpers.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// 2. Cari Schedule berdasarkan host_id, shift_id, date
	hostID := req.HostID
	shiftID := req.ShiftID
	date := req.Date

	var schedule models.Schedule
	if err := database.DB.Where("host_id = ? AND shift_id = ? AND date = ?", hostID, shiftID, date).First(&schedule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// 3. Generate Note apakah sesuai jadwal atau tidak
	attendanceDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected YYYY-MM-DD"})
		return
	}

	now := time.Now()
	systemNote, _ := generateNote(schedule, attendanceDate, req.ShiftID)

	attendance := models.Attendance{
		ScheduleID:  req.ShiftID,
		CheckedInAt: &now,
		Status:      "present",
		Note:        systemNote,
	}

	if err := database.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan presensi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Presensi berhasil", "data": attendance})
}

func AttendanceCheckOut(c *gin.Context) {
	type AttendanceCheckOutRequest struct {
		ScheduleID uint `json:"schedule_id" binding:"required"`
		ShiftID    uint `json:"shift_id" binding:"required"`
	}

	var req AttendanceCheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari presensi dengan ScheduleID & ShiftID & belum checkout
	var attendance models.Attendance
	if err := database.DB.
		Where("schedule_id = ? AND checked_out_at IS NULL", req.ScheduleID).
		First(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found or already checked out"})
		return
	}

	// Ambil schedule untuk validasi note
	var schedule models.Schedule
	if err := database.DB.First(&schedule, req.ScheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	now := time.Now()
	systemNote, _ := generateNote(schedule, now, req.ShiftID)

	attendance.CheckedOutAt = &now
	attendance.Note = systemNote

	if err := database.DB.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan checkout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout berhasil", "data": attendance})
}

func generateNote(schedule models.Schedule, attendanceDate time.Time, shiftID uint) (string, bool) {
	expectedDate := schedule.Date.Format("2006-01-02")
	actualDate := attendanceDate.Format("2006-01-02")

	dateMatch := expectedDate == actualDate
	shiftMatch := schedule.ShiftID == shiftID

	if dateMatch && shiftMatch {
		return "Sesuai Jadwal", true
	}
	if dateMatch && !shiftMatch {
		return "Shift tidak sesuai", false
	}
	if !dateMatch && shiftMatch {
		return "Tanggal tidak sesuai", false
	}
	return "Tanggal dan Shift tidak sesuai", false
}
