import "{{Module}}/internal/domains/{{feature}}/params"

type {{Feature}}Service interface {
	FindAll() ([]*params.{{Feature}}Response, error)
	FindByID(id string) (*params.{{Feature}}Response, error)
	Create(req params.Create{{Feature}}Request) (*params.{{Feature}}Response, error)
	Update(id string, req params.Update{{Feature}}Request) (*params.{{Feature}}Response, error)
	Delete(id string) error
}