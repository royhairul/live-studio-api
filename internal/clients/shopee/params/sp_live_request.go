package params

import "fmt"

// ShopeeLiveHistoryRequest maps to the query string of
// GET /supply/api/lm/sellercenter/liveList/v2
type ShopeeLiveHistoryRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`
	OrderBy  string `json:"orderBy"`
	Sort     string `json:"sort"`
	TimeDim  string `json:"timeDim"`
	EndDate  string `json:"endDate"`
}

// WithDefaults fills the values Shopee expects when the caller leaves them empty.
func (r ShopeeLiveHistoryRequest) WithDefaults() ShopeeLiveHistoryRequest {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}
	if r.TimeDim == "" {
		r.TimeDim = "1m"
	}
	return r
}

// ToQuery builds the query map. Every key is always sent, including the empty
// ones, because Shopee expects them to be present.
func (r ShopeeLiveHistoryRequest) ToQuery() map[string]string {
	return map[string]string{
		"page":     fmt.Sprint(r.Page),
		"pageSize": fmt.Sprint(r.PageSize),
		"name":     r.Name,
		"orderBy":  r.OrderBy,
		"sort":     r.Sort,
		"timeDim":  r.TimeDim,
		"endDate":  r.EndDate,
	}
}
