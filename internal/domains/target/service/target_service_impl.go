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

// Update implements TargetService.
func (s *TargetServiceImpl) Update(id string, req params.UpdateTargetRequest) (*params.TargetResponse, error) {
	panic("unimplemented")
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

func CalcRatio(real, target int64) float64 {
	if target == 0 {
		return 0
	}
	return (float64(real) / float64(target)) * 100
}

// FindByID implements TargetService.
func (s *TargetServiceImpl) FindByID(id string) (*params.TargetResponse, error) {
	panic("unimplemented")
}

// Delete implements TargetService.
func (s *TargetServiceImpl) Delete(id string) error {
	panic("unimplemented")
}
