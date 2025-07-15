package params

type ShopeeAccountResponse struct {
	ShopId   int    `json:"shopid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}
