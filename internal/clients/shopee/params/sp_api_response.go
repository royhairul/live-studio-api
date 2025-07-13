package params

type ShopeeApiResponse[T any] struct {
	Error    int    `json:"error"`
	ErrorMsg string `json:"error_msg"`
	Data     T      `json:"data"`
}
