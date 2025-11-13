package attendance

import (
	"github.com/royhairul/live-studio-api/internal/domains/attendance/controller"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/repository"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"attendance",
	fx.Provide(
		repository.NewAttendanceRepository,
		service.NewAttendanceService,
		controller.NewAttendanceController,
	),
	fx.Invoke(RegisterRoutes),
)
