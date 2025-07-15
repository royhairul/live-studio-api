package schedule

import (
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/controller"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/repository"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/service"

	shiftrepository "github.com/royhairul/live-studio-api/internal/domains/shift/repository"
)

func ProvideScheduleController() controller.ScheduleController {
	validate := validator.New()
	scheduleRepo := repository.NewScheduleRepository(database.DB)
	shiftRepo := shiftrepository.NewShiftRepository(database.DB)

	scheduleService := service.NewScheduleService(scheduleRepo, shiftRepo)

	return controller.NewScheduleController(scheduleService, validate)
}
