package params

type AccountadsRequest struct {
	// TODO: add request fields
}

type CreateAccountadsRequest struct {
	AccountID uint   `json:"account_id" validate:"required"`
	Date      string `json:"date" validate:"required"`
	Ads       uint   `json:"ads" vaildate:"required"`
}

type UpdateAccountadsRequest struct {
	AccountID *uint   `json:"account_id" validate:"required"`
	Date      *string `json:"date" validate:"required"`
	Ads       *uint   `json:"ads" vaildate:"required"`
}
