package product

import (
	"github.com/royhairul/live-studio-api/internal/domains/product/controller"
	"github.com/royhairul/live-studio-api/internal/domains/product/repository"
	"github.com/royhairul/live-studio-api/internal/domains/product/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"product",
	fx.Provide(
		repository.NewProductRepository,
		service.NewProductService,
		controller.NewProductController,
	),
	fx.Invoke(RegisterRoutes),
)
