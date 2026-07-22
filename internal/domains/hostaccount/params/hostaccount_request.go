package params

// CreateHostAccountRequest assigns an account to a host.
//
// ValidFrom defaults to today when omitted; ValidTo omitted means the assignment
// stays open. Dates are YYYY-MM-DD, matching the rest of the API.
type CreateHostAccountRequest struct {
	HostID    string `json:"host_id" validate:"required,uuid4"`
	AccountID uint   `json:"account_id" validate:"required"`
	ValidFrom string `json:"valid_from"`
	ValidTo   string `json:"valid_to"`
	Note      string `json:"note"`
}

// UpdateHostAccountRequest patches an assignment. Every field is optional so a
// caller can close an assignment by sending valid_to alone.
type UpdateHostAccountRequest struct {
	HostID    *string `json:"host_id" validate:"omitempty,uuid4"`
	AccountID *uint   `json:"account_id"`
	ValidFrom *string `json:"valid_from"`
	ValidTo   *string `json:"valid_to"`
	Note      *string `json:"note"`
}
