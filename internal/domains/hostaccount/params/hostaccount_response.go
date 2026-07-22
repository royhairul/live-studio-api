package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/hostaccount/entity"
)

// HostAccountSpan is the slice of a report window during which a host actually
// held an account. The aggregator uses it to decide which live sessions belong
// to the host: an assignment that started or ended mid-range must not pull in
// sessions from outside its own span.
type HostAccountSpan struct {
	AccountID uint
	From      time.Time
	To        time.Time
}

// Contains reports whether a session starting at t falls inside the span.
func (s HostAccountSpan) Contains(t time.Time) bool {
	return !t.Before(s.From) && !t.After(s.To)
}

type HostAccountResponse struct {
	ID uint `json:"id"`

	HostID   string `json:"host_id"`
	HostName string `json:"host_name"`

	AccountID   uint   `json:"account_id"`
	AccountName string `json:"account_name"`

	StudioID   uint   `json:"studio_id"`
	StudioName string `json:"studio_name"`

	ValidFrom time.Time  `json:"valid_from"`
	ValidTo   *time.Time `json:"valid_to"`
	IsActive  bool       `json:"is_active"`

	Note string `json:"note"`
}

func NewHostAccountResponse(data *entity.HostAccount) *HostAccountResponse {
	if data == nil {
		return nil
	}

	return &HostAccountResponse{
		ID:          data.ID,
		HostID:      data.HostID.String(),
		HostName:    data.Host.Name,
		AccountID:   data.AccountID,
		AccountName: data.Account.Name,
		StudioID:    uint(data.Account.StudioID),
		StudioName:  data.Account.Studio.Name,
		ValidFrom:   data.ValidFrom,
		ValidTo:     data.ValidTo,
		IsActive:    data.ValidTo == nil,
		Note:        data.Note,
	}
}
