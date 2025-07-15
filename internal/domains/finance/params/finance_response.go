package params

type FinanceResponse struct {
	AccountName string `json:"name"`
	Total       int    `json:"total"`
	ReportLive  any    `json:"reportLive"`
}
