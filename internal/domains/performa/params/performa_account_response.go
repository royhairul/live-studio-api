package params

type PerformaAccountResponse struct {
	CurrentPeriod  PeriodInfo                         `json:"current_period"`
	PreviousPeriod PeriodInfo                         `json:"previous_period"`
	Metrics        Metrics                            `json:"metrics"`
	List           []PerformaStudioDetailItemResponse `json:"list"`
}

type PerformaAccountDetailItemResponse struct {
	AccountID   uint   `json:"account_id"`
	AccountName string `json:"account_name"`
	PerformaMetricItem
}
