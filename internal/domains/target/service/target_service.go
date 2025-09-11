package service

import "github.com/royhairul/live-studio-api/internal/domains/target/params"

type TargetService interface {
	FindAll() ([]*params.TargetResponse, error)
	FindAllByDate(month, year string) ([]*params.TargetResponse, error)
	FindByID(id string) (*params.TargetResponse, error)
	Create(req params.CreateTargetRequest) (*params.CreatedTargetResponse, error)
	CreateOrUpdate(req params.CreateTargetRequest) (*params.CreatedTargetResponse, error)
	Update(id string, req params.UpdateTargetRequest) (*params.UpdatedTargetResponse, error)
	Delete(id string) error
}
