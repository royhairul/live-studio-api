package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/target/entity"
	"github.com/royhairul/live-studio-api/internal/domains/target/params"
	"github.com/royhairul/live-studio-api/internal/domains/target/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/constants"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"

	studioservice "github.com/royhairul/live-studio-api/internal/domains/studio/service"

	performaagg "github.com/royhairul/live-studio-api/internal/aggregator/performa"
)

type TargetServiceImpl struct {
	repository  repository.TargetRepository
	studioSvc   studioservice.StudioService
	performaAgg performaagg.PerformaAggregator
}

func NewTargetService(
	repository repository.TargetRepository,
	studioSvc studioservice.StudioService,
	performaAgg performaagg.PerformaAggregator,
) TargetService {
	return &TargetServiceImpl{
		repository,
		studioSvc,
		performaAgg,
	}
}

// Create implements TargetService.
func (s *TargetServiceImpl) Create(ctx context.Context, req params.CreateTargetRequest) (*params.CreatedTargetResponse, error) {
	parsedTime, err := time.Parse(constants.LayoutMMYY, req.Date)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse date: %v", err)
	}

	target := entity.Target{
		Date:         parsedTime,
		TargetGMV:    req.TargetGMV,
		TargetIncome: req.TargetIncome,
		StudioID:     req.StudioID,
	}

	created, err := s.repository.Create(ctx, &target)
	if err != nil {
		return nil, err
	}

	return params.NewCreatedTargetResponse(created), nil
}

// CreateOrUpdate implements TargetService.
func (s *TargetServiceImpl) CreateOrUpdate(ctx context.Context, req params.CreateTargetRequest) (*params.CreatedTargetResponse, error) {
	parsedTime, err := time.Parse(constants.LayoutMMYY, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	// Check target is created or not found
	exist, err := s.repository.FindByStudioAndDate(ctx, fmt.Sprintf("%d", req.StudioID), parsedTime)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var target *entity.Target
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// if not found, create new target
		target, err = s.repository.Create(ctx, &entity.Target{
			Date:         parsedTime,
			TargetGMV:    req.TargetGMV,
			TargetIncome: req.TargetIncome,
			StudioID:     req.StudioID,
		})
	} else {
		// if already exits, update target
		exist.TargetGMV = req.TargetGMV
		exist.TargetIncome = req.TargetIncome
		target, err = s.repository.Update(ctx, exist)
	}

	if err != nil {
		return nil, err
	}

	return params.NewCreatedTargetResponse(target), nil
}

// Update implements TargetService.
func (s *TargetServiceImpl) Update(ctx context.Context, id string, req params.UpdateTargetRequest) (*params.UpdatedTargetResponse, error) {
	target, err := s.repository.FindByID(ctx, id)
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
		parsedDate, err := timehandler.ParseDate(*req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD: %v", err)
		}
		target.Date = *parsedDate
	}

	// simpan update
	updated, err := s.repository.Update(ctx, target)
	if err != nil {
		return nil, err
	}

	return params.NewUpdatedTargetResponse(updated), nil
}

// FindAll implements TargetService.
func (s *TargetServiceImpl) FindAll(ctx context.Context, req params.TargetRequest) ([]*params.TargetResponse, error) {
	var (
		start, end  time.Time
		yearInt     int
		monthInt    int
		location    = time.Now().Location()
		dateDisplay string
	)

	// 🗓️ Jika month & year dikosongkan → gunakan bulan sekarang
	if req.Month == "" || req.Year == "" {
		now := time.Now()
		yearInt, monthInt = now.Year(), int(now.Month())
		start = time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, location)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		dateDisplay = fmt.Sprintf("%s %d", start.Month(), start.Year())
	} else {
		// 🗓️ Parse input bulan & tahun
		dateStr := fmt.Sprintf("%s %s", req.Month, req.Year)
		parsedTime, err := time.Parse(constants.LayoutMMYY, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %v", err)
		}

		yearInt = parsedTime.Year()
		monthInt = int(parsedTime.Month())
		start = time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, location)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		dateDisplay = fmt.Sprintf("%s %d", start.Month(), start.Year())
	}

	// Ambil semua studio
	studios, err := s.studioSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Ambil semua target
	targets, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Buat map target untuk bulan tsb
	targetMap := make(map[uint]*entity.Target)
	for _, t := range targets {
		if t.Date.Year() == yearInt && int(t.Date.Month()) == monthInt {
			targetMap[t.StudioID] = t
		}
	}

	// Gunakan performa aggregator
	studioGMV := make(map[uint]int64)
	studioIncome := make(map[uint]int64)

	for _, studio := range studios {
		performaList, total, err := s.performaAgg.CalculateByStudio(ctx, fmt.Sprint(studio.ID), &start, &end)
		if err != nil {
			log.Printf("warn: gagal ambil performa studioID=%d: %v", studio.ID, err)
			continue
		}

		studioGMV[studio.ID] = total.GMV
		studioIncome[studio.ID] = total.Income

		_ = performaList // bisa digunakan jika ingin breakdown
	}

	// Bangun response akhir
	var responses []*params.TargetResponse
	for _, studio := range studios {
		var targetGMV, targetIncome int64
		if t, ok := targetMap[studio.ID]; ok {
			targetGMV = t.TargetGMV
			targetIncome = t.TargetIncome
		}

		responses = append(responses, &params.TargetResponse{
			StudioID:   fmt.Sprint(studio.ID),
			StudioName: studio.Name,
			Date:       dateDisplay,
			GMV: params.Metric{
				Target: targetGMV,
				Real:   studioGMV[studio.ID],
				Ratio:  CalcRatio(studioGMV[studio.ID], targetGMV),
			},
			Income: params.Metric{
				Target: targetIncome,
				Real:   studioIncome[studio.ID],
				Ratio:  CalcRatio(studioIncome[studio.ID], targetIncome),
			},
		})
	}

	return responses, nil
}

func (s *TargetServiceImpl) FindByID(ctx context.Context, id string) (*params.TargetResponse, error) {
	// Ambil target by ID
	target, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Ambil studio terkait
	studio, err := s.studioSvc.FindByID(ctx, fmt.Sprintf("%d", target.StudioID))
	if err != nil {
		return nil, err
	}

	// Range bulan sesuai target.Date
	year, month, _ := target.Date.Date()
	location := target.Date.Location()

	start := time.Date(year, month, 1, 0, 0, 0, 0, location)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Gunakan Performa Aggregator untuk studio ini
	performaList, total, err := s.performaAgg.CalculateByStudio(ctx, fmt.Sprint(studio.ID), &start, &end)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate performa: %v", err)
	}

	_ = performaList // jika ingin breakdown akun, bisa digunakan di response detail

	// Bentuk response
	res := &params.TargetResponse{
		StudioID:   fmt.Sprint(studio.ID),
		StudioName: studio.Name,
		Date:       fmt.Sprintf("%s %d", start.Month(), start.Year()),
		GMV: params.Metric{
			Target: target.TargetGMV,
			Real:   total.GMV,
			Ratio:  CalcRatio(total.GMV, target.TargetGMV),
		},
		Income: params.Metric{
			Target: target.TargetIncome,
			Real:   total.Income,
			Ratio:  CalcRatio(total.Income, target.TargetIncome),
		},
	}

	return res, nil
}

// Delete implements TargetService.
func (s *TargetServiceImpl) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
