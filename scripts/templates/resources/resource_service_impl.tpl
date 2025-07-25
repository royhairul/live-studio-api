import (
	"fmt"
	"log"

    "{{Module}}/internal/domains/{{feature}}/entity"
    "{{Module}}/internal/domains/{{feature}}/params"
    "{{Module}}/internal/domains/{{feature}}/repository"
)

type {{Feature}}ServiceImpl struct {
	repository repository.{{Feature}}Repository
}

func New{{Feature}}Service(repository repository.{{Feature}}Repository) {{Feature}}Service {
	return &{{Feature}}ServiceImpl{repository}
}

// Create implements {{Feature}}Service.
func (s *{{Feature}}ServiceImpl) Create(req params.Create{{Feature}}Request) (*params.{{Feature}}Response, error) {
    panic("unimplemented")
}

// Update implements {{Feature}}Service.
func (s *{{Feature}}ServiceImpl) Update(id string, req params.Update{{Feature}}Request) (*params.{{Feature}}Response, error) {
    panic("unimplemented")
}

// FindAll implements {{Feature}}Service.
func (s *{{Feature}}ServiceImpl) FindAll() ([]*params.{{Feature}}Response, error) {
    panic("unimplemented")
}

// FindByID implements {{Feature}}Service.
func (s *{{Feature}}ServiceImpl) FindByID(id string) (*params.{{Feature}}Response, error) {
    panic("unimplemented")
}

// Delete implements {{Feature}}Service.
func (s *{{Feature}}ServiceImpl) Delete(id string) error {
    panic("unimplemented")
}