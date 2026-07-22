package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/entity"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/params"
	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
)

type HostAccountServiceImpl struct {
	repository repository.HostAccountRepository
	options    params.HostAccountFilter
}

func NewHostAccountService(repository repository.HostAccountRepository) HostAccountService {
	return &HostAccountServiceImpl{
		repository: repository,
		options:    params.HostAccountFilter{},
	}
}

func (s *HostAccountServiceImpl) clone() *HostAccountServiceImpl {
	options := s.options
	return &HostAccountServiceImpl{repository: s.repository, options: options}
}

func (s *HostAccountServiceImpl) WithID(id string) HostAccountService {
	next := s.clone()
	next.options.ID = &id
	return next
}

func (s *HostAccountServiceImpl) WithHostID(hostID string) HostAccountService {
	next := s.clone()
	next.options.HostID = &hostID
	return next
}

func (s *HostAccountServiceImpl) WithAccountID(accountID string) HostAccountService {
	next := s.clone()
	next.options.AccountID = &accountID
	return next
}

func (s *HostAccountServiceImpl) WithStudioID(studioID string) HostAccountService {
	next := s.clone()
	next.options.StudioID = &studioID
	return next
}

func (s *HostAccountServiceImpl) WithActiveOnly() HostAccountService {
	next := s.clone()
	next.options.OnlyActive = true
	return next
}

func (s *HostAccountServiceImpl) FindAll(ctx context.Context) ([]*params.HostAccountResponse, error) {
	items, err := s.repository.FindAll(ctx, s.options)
	if err != nil {
		return nil, err
	}

	results := make([]*params.HostAccountResponse, 0, len(items))
	for _, item := range items {
		results = append(results, params.NewHostAccountResponse(item))
	}
	return results, nil
}

func (s *HostAccountServiceImpl) FindOne(ctx context.Context) (*params.HostAccountResponse, error) {
	item, err := s.repository.FindOne(ctx, s.options)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewNotFoundError("host account assignment not found")
		}
		return nil, err
	}
	return params.NewHostAccountResponse(item), nil
}

func (s *HostAccountServiceImpl) Create(ctx context.Context, req params.CreateHostAccountRequest) (*params.HostAccountResponse, error) {
	hostID, err := uuid.Parse(req.HostID)
	if err != nil {
		return nil, errorhandler.NewBadRequestError("invalid host_id", err.Error())
	}

	validFrom, validTo, err := parseValidity(req.ValidFrom, req.ValidTo)
	if err != nil {
		return nil, err
	}

	// Serah terima di hari yang sama: kalau pemegang sebelumnya baru berhenti
	// tadi, mulai dari saat itu — bukan dari tengah malam, yang akan dianggap
	// bertumpuk dengan penugasan lama dan ditolak.
	if req.ValidFrom == "" {
		handover, err := s.latestHandover(ctx, req.AccountID, *validFrom)
		if err != nil {
			return nil, err
		}
		if handover != nil {
			validFrom = handover
		}
	}

	if err := s.assertNoConflict(ctx, req.AccountID, hostID, validFrom, validTo, nil); err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, &entity.HostAccount{
		HostID:    hostID,
		AccountID: req.AccountID,
		ValidFrom: *validFrom,
		ValidTo:   validTo,
		Note:      req.Note,
	})
	if err != nil {
		return nil, err
	}

	return params.NewHostAccountResponse(created), nil
}

func (s *HostAccountServiceImpl) Update(ctx context.Context, id string, req params.UpdateHostAccountRequest) (*params.HostAccountResponse, error) {
	existing, err := s.repository.FindOne(ctx, params.HostAccountFilter{ID: &id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorhandler.NewNotFoundError("host account assignment not found")
		}
		return nil, err
	}

	if req.HostID != nil {
		hostID, err := uuid.Parse(*req.HostID)
		if err != nil {
			return nil, errorhandler.NewBadRequestError("invalid host_id", err.Error())
		}
		existing.HostID = hostID
	}

	if req.AccountID != nil {
		existing.AccountID = *req.AccountID
	}

	if req.ValidFrom != nil {
		parsed, err := parseBoundary(*req.ValidFrom)
		if err != nil {
			return nil, errorhandler.NewBadRequestError("invalid valid_from", err.Error())
		}
		existing.ValidFrom = *parsed
	}

	if req.ValidTo != nil {
		// An empty string reopens the assignment; that is the only way to clear
		// valid_to, since omitting the field means "leave unchanged".
		if *req.ValidTo == "" {
			existing.ValidTo = nil
		} else {
			parsed, err := parseBoundary(*req.ValidTo)
			if err != nil {
				return nil, errorhandler.NewBadRequestError("invalid valid_to", err.Error())
			}
			existing.ValidTo = parsed
		}
	}

	if req.Note != nil {
		existing.Note = *req.Note
	}

	if existing.ValidTo != nil && !existing.ValidTo.After(existing.ValidFrom) {
		return nil, errorhandler.NewBadRequestError("invalid date range", "valid_to must be after valid_from")
	}

	if err := s.assertNoConflict(ctx, existing.AccountID, existing.HostID, &existing.ValidFrom, existing.ValidTo, &existing.ID); err != nil {
		return nil, err
	}

	updated, err := s.repository.Update(ctx, existing)
	if err != nil {
		return nil, err
	}

	return params.NewHostAccountResponse(updated), nil
}

func (s *HostAccountServiceImpl) Delete(ctx context.Context, id string) error {
	if _, err := s.repository.FindOne(ctx, params.HostAccountFilter{ID: &id}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorhandler.NewNotFoundError("host account assignment not found")
		}
		return err
	}
	return s.repository.Delete(ctx, id)
}

// AccountIDsForHost implements HostAccountService.
func (s *HostAccountServiceImpl) AccountIDsForHost(ctx context.Context, hostID string, start, end time.Time) ([]params.HostAccountSpan, error) {
	items, err := s.repository.FindAll(ctx, params.HostAccountFilter{
		HostID:    &hostID,
		StartDate: &start,
		EndDate:   &end,
	})
	if err != nil {
		return nil, err
	}

	spans := make([]params.HostAccountSpan, 0, len(items))
	for _, item := range items {
		// Clamp to the report window: an assignment reaching outside it must not
		// drag in sessions from before it started or after it ended.
		from := item.ValidFrom
		if from.Before(start) {
			from = start
		}

		to := end
		if item.ValidTo != nil && item.ValidTo.Before(to) {
			to = *item.ValidTo
		}

		if to.Before(from) {
			continue
		}

		spans = append(spans, params.HostAccountSpan{
			AccountID: item.AccountID,
			From:      from,
			To:        to,
		})
	}

	return spans, nil
}

// latestHandover returns the moment the account was last released, when that is
// later than the proposed start. It returns nil when nothing was released after
// that point, leaving the caller's own start date in place.
func (s *HostAccountServiceImpl) latestHandover(ctx context.Context, accountID uint, notBefore time.Time) (*time.Time, error) {
	accountKey := fmt.Sprint(accountID)

	existing, err := s.repository.FindAll(ctx, params.HostAccountFilter{AccountID: &accountKey})
	if err != nil {
		return nil, err
	}

	var latest *time.Time
	for _, item := range existing {
		if item.ValidTo == nil || !item.ValidTo.After(notBefore) {
			continue
		}
		if latest == nil || item.ValidTo.After(*latest) {
			latest = item.ValidTo
		}
	}

	return latest, nil
}

// assertNoConflict rejects an assignment that would leave one account held by
// two different hosts at the same time — a live session cannot belong to two
// hosts, so the ambiguity is blocked at write time rather than guessed at read
// time. The same host holding the account across adjacent spans is fine.
func (s *HostAccountServiceImpl) assertNoConflict(
	ctx context.Context,
	accountID uint,
	hostID uuid.UUID,
	validFrom *time.Time,
	validTo *time.Time,
	excludeID *uint,
) error {
	accountKey := fmt.Sprint(accountID)
	filter := params.HostAccountFilter{AccountID: &accountKey, ExcludeID: excludeID}

	existing, err := s.repository.FindAll(ctx, filter)
	if err != nil {
		return err
	}

	// An open-ended assignment runs to the end of time for overlap purposes.
	end := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	if validTo != nil {
		end = *validTo
	}

	for _, item := range existing {
		if item.HostID == hostID {
			continue
		}
		if item.Overlaps(*validFrom, end) {
			return errorhandler.NewBadRequestError(
				"account already assigned",
				fmt.Sprintf("account is held by host %s for an overlapping period", item.Host.Name),
			)
		}
	}

	return nil
}

// parseBoundary accepts either a plain date (YYYY-MM-DD, what the rest of the
// API uses) or a full RFC3339 timestamp.
//
// The timestamp form matters for ending an assignment: closing one at today's
// midnight would place valid_to before the moment it started, so "stop now"
// has to be expressible to the second.
func parseBoundary(value string) (*time.Time, error) {
	if parsed, err := timehandler.ParseDate(value); err == nil {
		return parsed, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// startOfToday is the default valid_from. Midnight rather than time.Now() so
// that assigning an account part-way through the day still counts that day's
// earlier sessions, which is what assigning it in the UI is taken to mean.
func startOfToday() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

// parseValidity resolves the request dates. An empty valid_from means "starting
// today", which is what assigning an account in the UI normally means.
func parseValidity(from, to string) (*time.Time, *time.Time, error) {
	validFrom := startOfToday()
	if from != "" {
		parsed, err := parseBoundary(from)
		if err != nil {
			return nil, nil, errorhandler.NewBadRequestError("invalid valid_from", err.Error())
		}
		validFrom = *parsed
	}

	var validTo *time.Time
	if to != "" {
		parsed, err := parseBoundary(to)
		if err != nil {
			return nil, nil, errorhandler.NewBadRequestError("invalid valid_to", err.Error())
		}
		if !parsed.After(validFrom) {
			return nil, nil, errorhandler.NewBadRequestError("invalid date range", "valid_to must be after valid_from")
		}
		validTo = parsed
	}

	return &validFrom, validTo, nil
}
