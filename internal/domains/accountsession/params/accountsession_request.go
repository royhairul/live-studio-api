package params

type AccountsessionRequest struct {
	// TODO: add request fields
}

type CreateAccountsessionRequest struct {
	// TODO: add request fields
	AccountID     uint `json:"account_id"`
	AttendanceID  uint `json:"attendance_id"`
	StudioID      uint `json:"studio_id"`
	GMVSalesStart uint `json:"gmv_sales_start"`
	GMVPaidStart  uint `json:"gmv_paid_start"`
}

type UpdateEndSessionRequest struct {
	// TODO: add request fields
	GMVSalesEnd uint `json:"gmv_sales_end"`
	GMVPaidEnd  uint `json:"gmv_paid_end"`
}
