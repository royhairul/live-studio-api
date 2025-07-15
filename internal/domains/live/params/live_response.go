package params

type LiveResponse struct {
	AccountName string `json:"name"`
	Relive      int    `json:"relive"`
	Total       int    `json:"total"`
	ReportLive  any    `json:"reportLive"`
}
