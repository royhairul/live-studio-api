package params

type PerformaAccountResponse struct {
	AccountID         uint    `json:"account_id"`
	AccountName       string  `json:"account_name"`
	GMV               int64   `json:"gmv"`
	CommissionPaid    int64   `json:"commission_paid"`
	CommissionPending int64   `json:"commission_pending"`
	Ads               int64   `json:"ads"`
	Acos              float64 `json:"acos"`
	Roas              float64 `json:"roas"`
	Income            int64   `json:"income"`
}
