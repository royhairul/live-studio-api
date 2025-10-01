package params

type PerformaAccountResponse struct {
	CurrentPeriod  PeriodInfo                         `json:"current_period"`
	PreviousPeriod PeriodInfo                         `json:"previous_period"`
	Metrics        Metrics                            `json:"metrics"`
	List           []PerformaStudioDetailItemResponse `json:"list"`
}

type PerformaAccountDetailItemResponse struct {
	AccountID   uint    `json:"account_id"`
	AccountName string  `json:"account_name"`
	GMV         int64   `json:"gmv"`
	Commission  int64   `json:"commission"`
	Ads         int64   `json:"ads"`
	Acos        float64 `json:"acos"`
	Roas        float64 `json:"roas"`
	Income      int64   `json:"income"`
}
