package service

import (
	"context"
	"fmt"
	"log"

	"github.com/royhairul/live-studio-api/internal/domains/host/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"

	"github.com/royhairul/live-studio-api/internal/domains/host/params"
	"github.com/royhairul/live-studio-api/internal/domains/host/repository"

	accountsessionservice "github.com/royhairul/live-studio-api/internal/domains/accountsession/service"
	attendanceservice "github.com/royhairul/live-studio-api/internal/domains/attendance/service"
)

type HostServiceImpl struct {
	repository        repository.HostRepository
	attendanceSvc     attendanceservice.AttendanceService
	accountSessionSvc accountsessionservice.AccountsessionService
}

func NewHostService(
	repository repository.HostRepository,
	attendanceSvc attendanceservice.AttendanceService,
	accountSessionSvc accountsessionservice.AccountsessionService,
) HostService {
	return &HostServiceImpl{repository, attendanceSvc, accountSessionSvc}
}

// Create implements HostService.
func (h *HostServiceImpl) Create(ctx context.Context, req params.CreateHostRequest) (*params.HostResponse, error) {
	host := entity.Host{
		Name:     req.Name,
		Phone:    req.Phone,
		StudioID: req.StudioID,
	}

	created, err := h.repository.Create(ctx, &host)
	if err != nil {
		return nil, err
	}

	result := params.NewHostResponse(created)

	return result, nil
}

// Update implements HostService.
func (h *HostServiceImpl) Update(ctx context.Context, id string, hostReq params.UpdateHostRequest) (*params.HostResponse, error) {
	host, err := h.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if hostReq.Name != nil {
		host.Name = *hostReq.Name
	}

	if hostReq.Phone != nil {
		host.Phone = *hostReq.Phone
	}

	if hostReq.StudioID != nil {
		host.StudioID = *hostReq.StudioID
		host.Studio = studioentity.Studio{}
	}

	saved, err := h.repository.Update(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to update host : %w", err)
	}

	log.Println(saved)

	result := params.NewHostResponse(saved)

	return result, nil
}

// FindAll implements HostService.
func (h *HostServiceImpl) FindAll(ctx context.Context) ([]*params.HostResponse, error) {
	hosts, err := h.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var results []*params.HostResponse
	for _, host := range hosts {
		results = append(results, params.NewHostResponse(host))
	}

	return results, nil
}

// FindByID implements HostService.
func (h *HostServiceImpl) FindByID(ctx context.Context, id string) (*params.HostResponse, error) {
	host, err := h.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := params.NewHostResponse(host)

	return result, nil
}

// Delete implements HostService.
func (h *HostServiceImpl) Delete(ctx context.Context, id string) error {
	if err := h.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// FindAllGroupedByStudio implements HostService.
func (h *HostServiceImpl) FindAllGroupedByStudio(ctx context.Context) ([]*params.HostGroupedByStudioResponse, error) {
	hosts, err := h.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	groupMap := make(map[uint]*params.HostGroupedByStudioResponse)
	for _, host := range hosts {
		group, exists := groupMap[host.StudioID]
		if !exists {
			group = &params.HostGroupedByStudioResponse{
				StudioID:   host.StudioID,
				StudioName: host.Studio.Name,
				Hosts:      []params.HostResponse{},
			}
			groupMap[host.StudioID] = group
		}

		group.Hosts = append(group.Hosts, *params.NewHostResponse(host))
	}

	var results []*params.HostGroupedByStudioResponse
	for _, group := range groupMap {
		results = append(results, group)
	}

	return results, nil
}
