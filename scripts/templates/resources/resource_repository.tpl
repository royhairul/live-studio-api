import (
	"{{Module}}/internal/domains/{{feature}}/entity"
)

type {{Feature}}Repository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.{{Feature}}, error)
	FindByID(id string) (*entity.{{Feature}}, error)
	Create(data *entity.{{Feature}}) (*entity.{{Feature}}, error)
	Update(data *entity.{{Feature}}) (*entity.{{Feature}}, error)
	Delete(id string) error
}