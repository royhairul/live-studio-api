package params

type AccountadsRequest struct {
	// TODO: add request fields
}

type CreateAccountadsRequest struct {
	AccountID uint   `json:"account_id"`
	Date      string `json:"date"`
	Ads       uint   `json:"ads"`
}

type UpdateAccountadsRequest struct {
	// TODO: add request fields
}
