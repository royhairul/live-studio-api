package params

type ShopeeLiveFinanceRequest struct {
	Name     string `json:"name,omitempty"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy,omitempty"`
	Sort     string `json:"sort,omitempty"`
	TimeDim  string `json:"timeDim,omitempty"`
	EndDate  string `json:"endDate,omitempty"`
}
