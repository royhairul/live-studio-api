package params

// ShopeeFinanceCommissionReportResponse merepresentasikan respons laporan komisi Shopee
// yang ditampilkan pada modul Finance di aplikasi.
type ShopeeFinanceCommissionReportResponse struct {
	List []ShopeeFinanceCommissionList `json:"list"`
	// IsTerminated           bool            `json:"is_terminated"`
	// IncomeBreakdown        IncomeBreakdown `json:"income_breakdown"`
	// PendingPayoutAmount    int64           `json:"pending_payout_amount"`
	// PendingPayoutAmountDis string          `json:"pending_payout_amount_dis"`
}

// ShopeeFinanceCommissionList merepresentasikan data ringkasan pembayaran komisi
// pada laporan keuangan Shopee (Finance > Laporan Komisi).
type ShopeeFinanceCommissionList struct {
	// Payment and Validation Information
	PaymentShopeeAdjustmentList []interface{} `json:"payment_shopee_adjustment_list"`
	ValidationID                string        `json:"validation_id"`
	IsMcn                       int           `json:"is_mcn"`
	WithVat                     bool          `json:"with_vat"`
	WithWht                     bool          `json:"with_wht"`
	WithSst                     bool          `json:"with_sst"`
	ValidationMonth             int           `json:"validation_month"`
	PaymentMonth                int           `json:"payment_month"`
	ValidationPeriodStartTime   int64         `json:"validation_period_start_time"`
	ValidationPeriodEndTime     int64         `json:"validation_period_end_time"`
	ValidationReviewTime        int64         `json:"validation_review_time"`
	ValidationPayoutStatus      int           `json:"validation_payout_status"`
	OverallValidationStatus     int           `json:"overall_validation_status"`

	// Invoice & Payment Status
	InvoiceStatus                  int    `json:"invoice_status"`
	InvoiceDismissReason           string `json:"invoice_dismiss_reason"`
	ServiceFeeInvoiceStatus        int    `json:"service_fee_invoice_status"`
	ServiceFeeInvoiceDismissReason string `json:"service_fee_invoice_dismiss_reason"`
	PaymentStatus                  int    `json:"payment_status"`
	PaymentRejectedReason          string `json:"payment_rejected_reason"`
	PaidFailedReasonType           int    `json:"paid_failed_reason_type"`
	PaymentFileStatus              int    `json:"payment_file_status"`
	UnableToPayDetailStatus        int    `json:"unable_to_pay_detail_status"`
	TransferStatus                 int    `json:"transfer_status"`

	// Period and Cycle
	OrderCompletedPeriodStartTime int64 `json:"order_completed_period_start_time"`
	OrderCompletedPeriodEndTime   int64 `json:"order_completed_period_end_time"`
	PaymentPeriodStartTime        int64 `json:"payment_period_start_time"`
	PaymentPeriodEndTime          int64 `json:"payment_period_end_time"`
	SettlementCycle               int   `json:"settlement_cycle"`
	Isoweek                       int   `json:"isoweek"`

	// Amounts
	TotalPaymentAmount      int64  `json:"total_payment_amount"`
	TotalPaymentAmountDis   string `json:"total_payment_amount_dis"`
	EligibleTotalAmount     int64  `json:"eligible_total_amount"`
	EligibleTotalAmountDis  string `json:"eligible_total_amount_dis"`
	BrandCommission         int64  `json:"brand_commission"`
	BrandCommissionDis      string `json:"brand_commission_dis"`
	ShopeeCommission        int64  `json:"shopee_commission"`
	ShopeeCommissionDis     string `json:"shopee_commission_dis"`
	BonusCommission         int64  `json:"bonus_commission"`
	BonusCommissionDis      string `json:"bonus_commission_dis"`
	McnManagementFee        int64  `json:"mcn_management_fee"`
	McnManagementFeeDis     string `json:"mcn_management_fee_dis"`
	PaymentShopeeAdjustment int64  `json:"payment_shopee_adjustment"`

	// Identifiers and Metadata
	PayoutID       string `json:"payout_id"`
	AccountType    int    `json:"account_type"`
	PaymentChannel int    `json:"payment_channel"`

	// Timestamps
	PayoutCreatedTime    int64 `json:"payout_created_time"`
	PaymentCompletedTime int64 `json:"payment_completed_time"`
	PaymentTime          int64 `json:"payment_time"`
}

// type Bill struct {
// 	ValidatedBrandCommission        int64  `json:"validated_brand_commission"`
// 	ValidatedBrandCommissionDis     string `json:"validated_brand_commission_dis"`
// 	EligibleShopeeCommission        int64  `json:"eligible_shopee_commission"`
// 	EligibleShopeeCommissionDis     string `json:"eligible_shopee_commission_dis"`
// 	EligibleTotalCommission         int64  `json:"eligible_total_commission"`
// 	EligibleTotalCommissionDis      string `json:"eligible_total_commission_dis"`
// 	OrderPlacedMonth                int    `json:"order_placed_month"`
// 	OrderQuantity                   int    `json:"order_quantity"`
// 	WhtAmountDis                    string `json:"wht_amount_dis"`
// 	VatAmountDis                    string `json:"vat_amount_dis"`
// 	WhtAmount                       int64  `json:"wht_amount"`
// 	VatAmount                       int64  `json:"vat_amount"`
// 	SstAmountDis                    string `json:"sst_amount_dis"`
// 	SstAmount                       int64  `json:"sst_amount"`
// 	ValidationStatus                int    `json:"validation_status"`
// 	RejectedReason                  string `json:"rejected_reason"`
// 	ReceiptStatus                   int    `json:"receipt_status"`
// 	InvoiceStatus                   int    `json:"invoice_status"`
// 	InvoiceDismissReason            string `json:"invoice_dismiss_reason"`
// 	TotalPaymentAmount              int64  `json:"total_payment_amount"`
// 	TotalPaymentAmountDis           string `json:"total_payment_amount_dis"`
// 	BillType                        int    `json:"bill_type"`
// 	BonusDescription                string `json:"bonus_description"`
// 	BrandTaxationBasedAmount        int64  `json:"brand_taxation_based_amount"`
// 	BrandTaxationBasedAmountDis     string `json:"brand_taxation_based_amount_dis"`
// 	ShopeeTaxationBasedAmount       int64  `json:"shopee_taxation_based_amount"`
// 	ShopeeTaxationBasedAmountDis    string `json:"shopee_taxation_based_amount_dis"`
// 	TotalTaxationBasedAmount        int64  `json:"total_taxation_based_amount"`
// 	TotalTaxationBasedAmountDis     string `json:"total_taxation_based_amount_dis"`
// 	PayableCommission               int64  `json:"payable_commission"`
// 	PayableCommissionDis            string `json:"payable_commission_dis"`
// 	HistoryTotalCommissionAmount    int64  `json:"history_total_commission_amount"`
// 	HistoryTotalCommissionAmountDis string `json:"history_total_commission_amount_dis"`
// 	EligiblePppFee                  int64  `json:"eligible_ppp_fee"`
// 	EligiblePppFeeDis               string `json:"eligible_ppp_fee_dis"`
// 	CompletePeriodStartTime         int64  `json:"complete_period_start_time"`
// 	CompletePeriodEndTime           int64  `json:"complete_period_end_time"`
// }

// type ServiceFee struct {
// 	PaymentServiceFeeTaxMap map[string]PaymentServiceFeeTax `json:"payment_service_fee_tax_map"`
// 	ValidationID            string                          `json:"validation_id"`
// 	OrderPlacedMonth        int                             `json:"order_placed_month"`
// 	BillSource              int                             `json:"bill_source"`
// 	ServiceFeeRate          int                             `json:"service_fee_rate"`
// 	ServiceFeeCap           int64                           `json:"service_fee_cap"`
// 	ServiceFee              int64                           `json:"service_fee"`
// 	ServiceFeeRounded       int64                           `json:"service_fee_rounded"`
// 	PreTaxServiceFee        int64                           `json:"pre_tax_service_fee"`
// 	PreTaxServiceFeeRounded int64                           `json:"pre_tax_service_fee_rounded"`
// }

// type PaymentServiceFeeTax struct {
// 	FeeID               int64 `json:"fee_id"`
// 	BaseTaxAmount       int64 `json:"base_tax_amount"`
// 	FeeTaxType          int   `json:"fee_tax_type"`
// 	FeeTaxRate          int   `json:"fee_tax_rate"`
// 	FeeTaxDirection     int   `json:"fee_tax_direction"`
// 	FeeTaxAmount        int64 `json:"fee_tax_amount"`
// 	FeeTaxAmountRounded int64 `json:"fee_tax_amount_rounded"`
// }

// type IncomeBreakdown struct {
// 	SocialMedias              PlatformIncome `json:"social_medias"`
// 	ShopeeLive                PlatformIncome `json:"shopee_live"`
// 	ShopeeVideo               PlatformIncome `json:"shopee_video"`
// 	TotalPaymentAmount        int64          `json:"total_payment_amount"`
// 	TotalPaymentAmountDis     string         `json:"total_payment_amount_dis"`
// 	SocialMediasCommission    int64          `json:"social_medias_commission"`
// 	SocialMediasCommissionDis string         `json:"social_medias_commission_dis"`
// 	ShopeeLiveCommission      int64          `json:"shopee_live_commission"`
// 	ShopeeLiveCommissionDis   string         `json:"shopee_live_commission_dis"`
// 	ShopeeVideoCommission     int64          `json:"shopee_video_commission"`
// 	ShopeeVideoCommissionDis  string         `json:"shopee_video_commission_dis"`
// 	ServiceFee                int64          `json:"service_fee"`
// 	ServiceFeeDis             string         `json:"service_fee_dis"`
// }

// type PlatformIncome struct {
// 	SellerCommission Commission       `json:"seller_commission"`
// 	ShopeeCommission ShopeeCommission `json:"shopee_commission"`
// 	BonusCommission  BonusCommission  `json:"bonus_commission"`
// 	PppCommission    Commission       `json:"ppp_commission"`
// 	McnManagementFee Commission       `json:"mcn_management_fee"`
// 	ServiceFee       Commission       `json:"service_fee"`
// }

// type Commission struct {
// 	SellerComm          int64  `json:"seller_comm,omitempty"`
// 	SellerCommDis       string `json:"seller_comm_dis,omitempty"`
// 	ShopeeComm          int64  `json:"shopee_comm,omitempty"`
// 	ShopeeCommDis       string `json:"shopee_comm_dis,omitempty"`
// 	PppComm             int64  `json:"ppp_comm,omitempty"`
// 	PppCommDis          string `json:"ppp_comm_dis,omitempty"`
// 	McnManagementFee    int64  `json:"mcn_management_fee,omitempty"`
// 	McnManagementFeeDis string `json:"mcn_management_fee_dis,omitempty"`
// 	ServiceFee          int64  `json:"service_fee,omitempty"`
// 	ServiceFeeDis       string `json:"service_fee_dis,omitempty"`
// }

// type ShopeeCommission struct {
// 	ShopeeComm                int64  `json:"shopee_comm"`
// 	ShopeeCommDis             string `json:"shopee_comm_dis"`
// 	OrderedInSameShop         int64  `json:"ordered_in_same_shop"`
// 	OrderedInSameShopDis      string `json:"ordered_in_same_shop_dis"`
// 	OrderedInDifferentShop    int64  `json:"ordered_in_different_shop"`
// 	OrderedInDifferentShopDis string `json:"ordered_in_different_shop_dis"`
// }

// type BonusCommission struct {
// 	BonusType    []BonusType `json:"bonus_type"`
// 	BonusComm    int64       `json:"bonus_comm"`
// 	BonusCommDis string      `json:"bonus_comm_dis"`
// }

// type BonusType struct {
// 	Header          string `json:"header"`
// 	Content         int64  `json:"content"`
// 	ContentDis      string `json:"content_dis"`
// 	EventCampaignID int64  `json:"event_campaign_id"`
// }
