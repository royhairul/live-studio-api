package params

type ShopeeApiResponse[T any] struct {
	Error    int    `json:"error,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Code     int    `json:"code,omitempty"`
	Data     T      `json:"data"`
}

type ShopeeApiPaginationResult[T any] struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
	List     []T `json:"list"`
}
