package params

type ShopeeCheckoutResponse struct {
	TotalCount int                   `json:"total_count"`
	List       []*ShopeeCheckoutList `json:"list"`
}

type ShopeeCheckoutList struct {
	PurchaseTime                    int64                       `json:"purchase_time"`
	CheckoutID                      string                      `json:"checkout_id"`
	CheckoutStatus                  string                      `json:"checkout_status"`
	ConversionStatus                int                         `json:"conversion_status"`
	CheckoutCompleteTime            int64                       `json:"checkout_complete_time"`
	AffiliateID                     int64                       `json:"affiliate_id"`
	AffiliateName                   string                      `json:"affiliate_name"`
	UserStatus                      string                      `json:"user_status"`
	GrossCommission                 int64                       `json:"gross_commission"`
	CappedCommission                int64                       `json:"capped_commission"`
	TotalBrandCommission            int64                       `json:"total_brand_commission"`
	EstimatedTotalCommissionWithMCN int64                       `json:"estimated_total_commission_with_mcn"`
	EstimatedTotalCommission        int64                       `json:"estimated_total_commission"`
	UtmContent                      string                      `json:"utm_content"`
	Device                          string                      `json:"device"`
	Referrer                        string                      `json:"referrer"` // bisa diubah ke struct jika ingin diparse
	Orders                          []ShopeeOrder               `json:"orders"`
	ClickTime                       int64                       `json:"click_time"`
	ClickID                         string                      `json:"click_id"`
	ProductType                     string                      `json:"product_type"`
	InternalSource                  string                      `json:"internal_source"`
	IndirectSource                  string                      `json:"indirect_source"`
	DirectSource                    string                      `json:"direct_source"`
	LastExternalSource              string                      `json:"last_external_source"`
	FirstExternalSource             string                      `json:"first_external_source"`
	IsShopeeCapped                  bool                        `json:"is_shopee_capped"`
	AttributionType                 int                         `json:"attribution_type"`
	EstimatedValidationMonth        string                      `json:"estimated_validation_month"`
	ReportPaymentValidationInfo     ReportPaymentValidationInfo `json:"report_payment_validation_info"`
	AffiliateNetCommission          string                      `json:"affiliate_net_commission"`
	McnManagementFeeCommission      string                      `json:"mcn_management_fee_commission"`
	McnAgreementID                  string                      `json:"mcn_agreement_id"`
	CampaignMcnID                   string                      `json:"campaign_mcn_id"`
	CampaignMcnName                 string                      `json:"campaign_mcn_name"`
	LinkedMcnID                     string                      `json:"linked_mcn_id"`
	LinkedMcnName                   string                      `json:"linked_mcn_name"`
	LinkedMcnCommissionRate         string                      `json:"linked_mcn_commission_rate"`
	Tenant                          int                         `json:"tenant"`
	AppType                         int                         `json:"app_type"`
}

type ShopeeOrder struct {
	OrderSN                string            `json:"order_sn"`
	OrderID                string            `json:"order_id"`
	OrderStatus            string            `json:"order_status"`
	ShopType               int               `json:"shop_type"`
	CancelReason           string            `json:"cancel_reason"`
	DisplayOrderStatus     int               `json:"display_order_status"`
	CompleteTime           int64             `json:"complete_time"`
	FraudCompleteTime      int64             `json:"fraud_complete_time"`
	AffiliateTransactionID string            `json:"affiliate_transaction_id"`
	ShopeeOrderStatus      int               `json:"shopee_order_status"`
	Items                  []ShopeeOrderItem `json:"items"`
}

type ShopeeOrderItem struct {
	ItemStatus                      string `json:"item_status"`
	DisplayItemStatus               string `json:"display_item_status"`
	AffiliateItemStatus             int    `json:"affiliate_item_status"`
	ShopID                          int64  `json:"shop_id"`
	ShopName                        string `json:"shop_name"`
	PromotionID                     string `json:"promotion_id"`
	ModelID                         string `json:"model_id"`
	ItemID                          int64  `json:"item_id"`
	ItemName                        string `json:"item_name"`
	ItemPrice                       int64  `json:"item_price"`
	ActualAmount                    int64  `json:"actual_amount"`
	RefundedAmount                  int64  `json:"refunded_amount"`
	Qty                             int    `json:"qty"`
	ImgCode                         string `json:"img_code"`
	ItemCommission                  int64  `json:"item_commission"`
	CappedBrandCommission           int64  `json:"capped_brand_commission"`
	GlobalCategoryLv1ID             int    `json:"global_category_lv1_id"`
	GlobalCategoryLv2ID             int    `json:"global_category_lv2_id"`
	GlobalCategoryLv3ID             int    `json:"global_category_lv3_id"`
	GlobalCategoryLv1Name           string `json:"global_category_lv1_name"`
	GlobalCategoryLv2Name           string `json:"global_category_lv2_name"`
	GlobalCategoryLv3Name           string `json:"global_category_lv3_name"`
	IsFraud                         int    `json:"is_fraud"`
	FraudReason                     string `json:"fraud_reason"`
	FraudStatus                     int    `json:"fraud_status"`
	BrandCommissionRate             int    `json:"brand_commission_rate"`
	PlatformCommissionRate          int    `json:"platform_commission_rate"`
	AttributionType                 int    `json:"attribution_type"`
	Channel                         int    `json:"channel"`
	CampaignMCNBrandGrossCommission string `json:"campaign_mcn_brand_gross_commission"`
	CampaignType                    int    `json:"campaign_type"`
}

type ReportPaymentValidationInfo struct {
	ValidationCycle                    int    `json:"validation_cycle"`
	EstimateValidationMonth            string `json:"estimate_validation_month"`
	EstimateValidationISOWeek          int    `json:"estimate_validation_isoweek"`
	OrderEstimateValidationPeriodStart int64  `json:"order_estimate_validation_period_start"`
	OrderEstimateValidationPeriodEnd   int64  `json:"order_estimate_validation_period_end"`
}
