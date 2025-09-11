package service

import (
	"fmt"
	"log"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
	"github.com/royhairul/live-studio-api/internal/domains/target/params"
	"github.com/royhairul/live-studio-api/internal/domains/target/repository"

	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"
	transactionservice "github.com/royhairul/live-studio-api/internal/domains/transaction/service"
)

type TargetServiceImpl struct {
	repository        repository.TargetRepository
	studioSvc         studioservice.StudioService
	transactionSvc    transactionservice.TransactionService
	atttendanceSvc    attendanceservice.AttendanceService
	accountsessionSvc accountsessionservice.AccountsessionService
}

func NewTargetService(
	repository repository.TargetRepository,
	studioSvc studioservice.StudioService,
	transactionSvc transactionservice.TransactionService,
	attendanceSvc attendanceservice.AttendanceService,
	accountsessionSvc accountsessionservice.AccountsessionService,
) TargetService {
	return &TargetServiceImpl{repository, studioSvc, transactionSvc, attendanceSvc, accountsessionSvc}
}

// Create implements TargetService.
func (s *TargetServiceImpl) Create(req params.CreateTargetRequest) (*params.CreatedTargetResponse, error) {
	parsedTime, err := time.Parse("January 2006", req.Date)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse date: %v", err)
	}

	target := entity.Target{
		Date:         parsedTime,
		TargetGMV:    req.TargetGMV,
		TargetIncome: req.TargetIncome,
		StudioID:     req.StudioID,
	}

	created, err := s.repository.Create(&target)
	if err != nil {
		return nil, err
	}

	result := params.CreatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", created.StudioID),
		StudioName:   created.Studio.Name,
		Date:         fmt.Sprintf("%v", created.Date),
		TargetGMV:    created.TargetGMV,
		TargetIncome: created.TargetGMV,
	}

	return &result, nil
}

// CreateOrUpdate implements TargetService.
func (s *TargetServiceImpl) CreateOrUpdate(req params.CreateTargetRequest) (*params.CreatedTargetResponse, error) {
	parsedTime, err := time.Parse("January 2006", req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	// Check apakah sudah ada target dengan date + studio
	exist, err := s.repository.FindByStudioAndDate(fmt.Sprintf("%d", req.StudioID), parsedTime)
	if err != nil {
		return nil, err
	}

	var target *entity.Target
	if exist != nil {
		// update data existing
		exist.TargetGMV = req.TargetGMV
		exist.TargetIncome = req.TargetIncome
		updated, err := s.repository.Update(exist)
		if err != nil {
			return nil, err
		}
		target = updated
	} else {
		// create baru
		newTarget := entity.Target{
			Date:         parsedTime,
			TargetGMV:    req.TargetGMV,
			TargetIncome: req.TargetIncome,
			StudioID:     req.StudioID,
		}
		created, err := s.repository.Create(&newTarget)
		if err != nil {
			return nil, err
		}
		target = created
	}

	result := params.CreatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", target.StudioID),
		StudioName:   target.Studio.Name,
		Date:         target.Date.Format("2006-01-02"), // lebih konsisten
		TargetGMV:    target.TargetGMV,
		TargetIncome: target.TargetIncome,
	}

	return &result, nil
}

// Update implements TargetService.
func (s *TargetServiceImpl) Update(id string, req params.UpdateTargetRequest) (*params.UpdatedTargetResponse, error) {
	target, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.StudioID != 0 {
		target.StudioID = req.StudioID
	}
	if req.TargetGMV != nil {
		target.TargetGMV = *req.TargetGMV
	}
	if req.TargetIncome != nil {
		target.TargetIncome = *req.TargetIncome
	}
	if req.Date != nil && *req.Date != "" {
		// contoh format: "2025-09-01"
		parsedDate, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD: %v", err)
		}
		target.Date = parsedDate
	}

	// simpan update
	updated, err := s.repository.Update(target)
	if err != nil {
		return nil, err
	}

	result := params.UpdatedTargetResponse{
		StudioID:     fmt.Sprintf("%d", updated.StudioID),
		StudioName:   updated.Studio.Name,
		Date:         fmt.Sprintf("%v", updated.Date),
		TargetGMV:    updated.TargetGMV,
		TargetIncome: updated.TargetGMV,
	}

	return &result, nil
}

// FindAll implements TargetService.
func (s *TargetServiceImpl) FindAll() ([]*params.TargetResponse, error) {
	studios, err := s.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	targets, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	// range bulan sekarang
	now := time.Now()
	year, month, _ := now.Date()
	location := now.Location()

	start := time.Date(year, month, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond) // akhir bulan

	// ambil semua absensi bulan ini
	attendances, err := s.atttendanceSvc.FindByDateRange(&start, &end)
	if err != nil {
		return nil, err
	}

	// kumpulkan semua data per studio
	studioGMV := make(map[uint]int64)    // GMV dari accountsession
	studioIncome := make(map[uint]int64) // Income dari transaksi

	for _, att := range attendances {
		sessions, err := s.accountsessionSvc.FindAllByAttendanceID(fmt.Sprintf("%d", att.ID))
		if err != nil {
			continue
		}

		for _, session := range sessions {
			// akumulasi GMV
			log.Println("studio id: ", session.StudioID)
			studioGMV[session.StudioID] += int64(session.GMVPaid)

			// ambil transaksi by account
			trxs, err := s.transactionSvc.FindAllByDate(fmt.Sprintf("%d", session.AccountID), &start, &end)
			if err != nil {
				continue
			}

			// akumulasi income
			for _, trx := range trxs {
				studioIncome[session.StudioID] += trx.Commission.Total // pastikan field trx.Amount ada
			}
		}
	}

	// debug log
	for studioID, gmv := range studioGMV {
		log.Printf("StudioID: %d | GMV: %d | Income: %d\n", studioID, gmv, studioIncome[studioID])
	}

	// mapping hasil ke response
	var responses []*params.TargetResponse
	for _, studio := range studios {
		// cari target bulan ini untuk studio
		var targetGMV, targetIncome int64
		for _, t := range targets {
			if t.StudioID == studio.ID &&
				t.Date.Year() == year &&
				t.Date.Month() == month {
				targetGMV = t.TargetGMV
				targetIncome = t.TargetIncome
				break
			}
		}

		// realisasi
		realGMV := studioGMV[studio.ID]
		realIncome := studioIncome[studio.ID]

		// bentuk response
		res := &params.TargetResponse{
			StudioID:   fmt.Sprintf("%d", studio.ID),
			StudioName: studio.Name,
			Date:       fmt.Sprintf("%s %d", start.Month(), start.Year()),
			GMV: params.Metric{
				Target: targetGMV,
				Real:   realGMV,
				Ratio:  CalcRatio(realGMV, targetGMV),
			},
			Income: params.Metric{
				Target: targetIncome,
				Real:   realIncome,
				Ratio:  CalcRatio(realIncome, targetIncome),
			},
		}

		responses = append(responses, res)
	}

	return responses, nil
}

func (s *TargetServiceImpl) FindAllByDate(month, year string) ([]*params.TargetResponse, error) {
	// Format jadi "September 2025"
	dateStr := fmt.Sprintf("%s %s", month, year)

	parsedTime, err := time.Parse("January 2006", dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	yearInt := parsedTime.Year()
	monthInt := int(parsedTime.Month())

	// Tentukan range tanggal bulan tsb
	location := time.Now().Location()
	start := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Ambil semua studio
	studios, err := s.studioSvc.FindAll()
	if err != nil {
		return nil, err
	}

	// Ambil target sesuai bulan & tahun
	targets, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	targetMap := make(map[uint]*entity.Target)
	for _, t := range targets {
		if t.Date.Year() == yearInt && int(t.Date.Month()) == monthInt {
			targetMap[t.StudioID] = t
		}
	}

	// Ambil semua absensi bulan tsb
	attendances, err := s.atttendanceSvc.FindByDateRange(&start, &end)
	if err != nil {
		return nil, err
	}

	// Kumpulkan data per studio
	studioGMV := make(map[uint]int64)
	studioIncome := make(map[uint]int64)

	for _, att := range attendances {
		sessions, err := s.accountsessionSvc.FindAllByAttendanceID(fmt.Sprint(att.ID))
		if err != nil {
			log.Printf("warn: gagal ambil sessions untuk attendanceID=%d: %v", att.ID, err)
			continue
		}

		for _, session := range sessions {
			studioGMV[session.StudioID] += int64(session.GMVPaid)

			trxs, err := s.transactionSvc.FindAllByDate(fmt.Sprint(session.AccountID), &start, &end)
			if err != nil {
				log.Printf("warn: gagal ambil transaksi untuk accountID=%d: %v", session.AccountID, err)
				continue
			}

			for _, trx := range trxs {
				studioIncome[session.StudioID] += trx.Commission.Total
			}
		}
	}

	// Mapping hasil ke response
	var responses []*params.TargetResponse
	for _, studio := range studios {
		var targetGMV, targetIncome int64
		if t, ok := targetMap[studio.ID]; ok {
			targetGMV = t.TargetGMV
			targetIncome = t.TargetIncome
		}

		realGMV := studioGMV[studio.ID]
		realIncome := studioIncome[studio.ID]

		responses = append(responses, &params.TargetResponse{
			StudioID:   fmt.Sprint(studio.ID),
			StudioName: studio.Name,
			Date:       fmt.Sprintf("%s %d", start.Month(), start.Year()),
			GMV: params.Metric{
				Target: targetGMV,
				Real:   realGMV,
				Ratio:  CalcRatio(realGMV, targetGMV),
			},
			Income: params.Metric{
				Target: targetIncome,
				Real:   realIncome,
				Ratio:  CalcRatio(realIncome, targetIncome),
			},
		})
	}

	return responses, nil
}

// FindByID implements TargetService.
func (s *TargetServiceImpl) FindByID(id string) (*params.TargetResponse, error) {
	// ambil target by ID
	target, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// ambil studio terkait
	studio, err := s.studioSvc.FindByID(fmt.Sprintf("%d", target.StudioID))
	if err != nil {
		return nil, err
	}

	// range bulan sesuai target.Date
	year, month, _ := target.Date.Date()
	location := target.Date.Location()

	start := time.Date(year, month, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// ambil absensi bulan target
	attendances, err := s.atttendanceSvc.FindByDateRange(&start, &end)
	if err != nil {
		return nil, err
	}

	// kumpulkan data realisasi
	var realGMV, realIncome int64
	for _, att := range attendances {
		sessions, err := s.accountsessionSvc.FindAllByAttendanceID(fmt.Sprint(att.ID))
		if err != nil {
			continue
		}

		for _, session := range sessions {
			if session.StudioID != target.StudioID {
				continue
			}

			// akumulasi GMV
			realGMV += int64(session.GMVPaid)

			// transaksi by account
			trxs, err := s.transactionSvc.FindAllByDate(fmt.Sprint(session.AccountID), &start, &end)
			if err != nil {
				continue
			}

			for _, trx := range trxs {
				realIncome += trx.Commission.Total
			}
		}
	}

	// bentuk response
	res := &params.TargetResponse{
		StudioID:   fmt.Sprint(studio.ID),
		StudioName: studio.Name,
		Date:       fmt.Sprintf("%s %d", start.Month(), start.Year()),
		GMV: params.Metric{
			Target: target.TargetGMV,
			Real:   realGMV,
			Ratio:  CalcRatio(realGMV, target.TargetGMV),
		},
		Income: params.Metric{
			Target: target.TargetIncome,
			Real:   realIncome,
			Ratio:  CalcRatio(realIncome, target.TargetIncome),
		},
	}

	return res, nil
}

// Delete implements TargetService.
func (s *TargetServiceImpl) Delete(id string) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}

	return nil
}
