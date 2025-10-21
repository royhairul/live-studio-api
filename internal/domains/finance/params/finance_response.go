package params

import (
	"fmt"
	"time"

	shopeeparam "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	accountparam "github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
)

type FinanceResponse struct {
	AccountID      string `json:"id"`
	AccountName    string `json:"name"`
	OrderDate      string `json:"order_date"`
	ValidationDate string `json:"validation_date"`
	Commission     int64  `json:"commission"`
	PaymentStatus  string `json:"payment_status"`
	PaymentMethod  string `json:"payment_method"`
	PaymentDate    string `json:"payment_date"`
}

func GetPaymentMethod(value int) string {
	switch value {
	case 1:
		return "Transfer Bank"
	case 2:
		return "ShopeePay"
	default:
		return "-"
	}
}

func GetPaymentStatusLabel(status int) string {
	switch status {
	case 4:
		return "Sudah Dibayar"
	case 9:
		return "Menunggu Validasi"
	case 10:
		return "Menunggu Dibayar"
	default:
		return "Unknown"
	}
}

func ParsePaymentDate(paymentTime int64) string {
	t, err := time.Parse("20060102", fmt.Sprint(paymentTime))
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func NewFinanceResponse(account accountparam.AccountResponse, commission shopeeparam.ShopeeFinanceCommissionList) *FinanceResponse {
	return &FinanceResponse{
		AccountID:      account.UniqueID,
		AccountName:    account.Name,
		Commission:     commission.TotalPaymentAmount,
		OrderDate:      timehandler.FormatInt64Date(commission.OrderCompletedPeriodEndTime),
		ValidationDate: timehandler.FormatInt64Date(commission.ValidationReviewTime),
		PaymentDate:    ParsePaymentDate(commission.PaymentTime),
		PaymentStatus:  GetPaymentStatusLabel(commission.PaymentStatus),
		PaymentMethod:  GetPaymentMethod(commission.PaymentChannel),
	}
}
