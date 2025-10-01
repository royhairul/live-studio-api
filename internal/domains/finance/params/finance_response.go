package params

type FinanceResponse struct {
	AccountName string `json:"name"`
	Total       int    `json:"total"`
	ReportLive  any    `json:"reportLive"`
}

type CommissionTotalResponse struct {
	TotalGMV        uint                       `json:"total_gmv"`
	TotalCommission uint                       `json:"total_commission"`
	TotalIncome     uint                       `json:"total_income"`
	List            []CommissionStudioResponse `json:"list"`
}
type CommissionStudioResponse struct {
	StudioName      string `json:"studio_name"`
	TotalGMV        uint   `json:"total_gmv"`
	TotalCommission uint   `json:"total_commission"`
	TotalIncome     uint   `json:"total_income"`
}

type CommissionStudioDetailResponse struct {
	StudioName        string `json:"studio_name"`
	AccountName       string `json:"account_name"`
	GMV               uint   `json:"gmv"`
	CommissionSuccess uint   `json:"commission_success"`
	CommissionPending uint   `json:"commission_pending"`
	ACOS              uint   `json:"acos"`
	ROAS              uint   `json:"roas"`
	Ads               uint   `json:"ads"`
	Income            uint   `json:"income"`
}
