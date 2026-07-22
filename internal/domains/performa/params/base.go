package params

type PeriodInfo struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Days  int    `json:"days"`
}

type Metric struct {
	Total int64 `json:"total"`
	Diff  int64 `json:"diff"`
	Ratio int64 `json:"ratio"`
}

type Metrics struct {
	Commission Metric `json:"commission"`
	Income     Metric `json:"income"`
	GMV        Metric `json:"gmv"`
	Ads        Metric `json:"ads"`
}

// PerformaLiveMetrics are the engagement figures shared by the host, account,
// and studio pages. They come from the lives table, which is the only place
// Shopee's viewer and click data is stored.
//
// No omitempty: the dashboard needs a stable 0 rather than a missing key, and 0
// is a meaningful value here — Shopee reports null for several of these.
type PerformaLiveMetrics struct {
	// CTR is product clicks over views, as a fraction (0.15 = 15%).
	CTR float64 `json:"ctr"`
	// ConversionRate is Shopee's own figure, also a fraction.
	ConversionRate float64 `json:"conversion_rate"`
	// ActiveViewers is the sum of engaged unique viewers.
	ActiveViewers int `json:"active_viewers"`
}

type PerformaMetricItem struct {
	GMV        int64   `json:"gmv"`
	Commission int64   `json:"commission"`
	Ads        int64   `json:"ads"`
	Income     int64   `json:"income"`
	Acos       float64 `json:"acos,omitempty"`
	Roas       float64 `json:"roas,omitempty"`

	PerformaLiveMetrics
}
