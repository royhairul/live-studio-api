package params

type PerformaStudioResponse struct {
	CurrentPeriod  PeriodInfo                   `json:"current_period"`
	PreviousPeriod PeriodInfo                   `json:"previous_period"`
	Metrics        Metrics                      `json:"metrics"`
	List           []PerformaStudioItemResponse `json:"list"`
}

type PerformaStudioItemResponse struct {
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	PerformaMetricItem
}

type PerformaStudioDetailResponse struct {
	StudioID       uint                               `json:"studio_id"`
	StudioName     string                             `json:"studio_name"`
	CurrentPeriod  PeriodInfo                         `json:"current_period"`
	PreviousPeriod PeriodInfo                         `json:"previous_period"`
	Metrics        Metrics                            `json:"metrics"`
	List           []PerformaStudioDetailItemResponse `json:"list"`
}

type PerformaStudioDetailItemResponse struct {
	AccountID   uint   `json:"account_id"`
	AccountName string `json:"account_name"`
	PerformaMetricItem
}
