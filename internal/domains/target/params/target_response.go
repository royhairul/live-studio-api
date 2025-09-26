package params

type Metric struct {
	Real   int64   `json:"real"`
	Target int64   `json:"target"`
	Ratio  float64 `json:"ratio"`
}

type CreatedTargetResponse struct {
	StudioID     string `json:"studio_id"`
	StudioName   string `json:"studio_name"`
	Date         string `json:"date"`
	TargetGMV    int64  `json:"target_gmv"`
	TargetIncome int64  `json:"target_income"`
}

type UpdatedTargetResponse struct {
	StudioID     string `json:"studio_id"`
	StudioName   string `json:"studio_name"`
	Date         string `json:"date"`
	TargetGMV    int64  `json:"target_gmv"`
	TargetIncome int64  `json:"target_income"`
}

type TargetResponse struct {
	// TODO: add response fields
	StudioID   string `json:"studio_id"`
	StudioName string `json:"studio_name"`
	Date       string `json:"date"`
	GMV        Metric `json:"gmv"`
	Income     Metric `json:"income"`
}
