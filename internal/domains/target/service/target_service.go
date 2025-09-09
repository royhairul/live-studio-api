package service

import "github.com/royhairul/live-studio-api/internal/domains/target/params"

type TargetService interface {
	FindAll() ([]*params.TargetResponse, error)
	FindByID(id string) (*params.TargetResponse, error)
	Create(req params.CreateTargetRequest) (*params.CreatedTargetResponse, error)
	Update(id string, req params.UpdateTargetRequest) (*params.TargetResponse, error)
	Delete(id string) error
}
