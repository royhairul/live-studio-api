package {{feature}}

import (
	"{{Module}}/internal/domains/{{feature}}/controller"
	"{{Module}}/internal/domains/{{feature}}/repository"
	"{{Module}}/internal/domains/{{feature}}/service"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"{{feature}}",
	fx.Provide(
		repository.New{{Feature}}Repository,
		service.New{{Feature}}Service,
		controller.New{{Feature}}Controller,
	),
	fx.Invoke(RegisterRoutes),
)
