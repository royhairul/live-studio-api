package params

import (
	"fmt"
	"time"

	shopeeparam "github.com/royhairul/live-studio-api/internal/clients/shopee/params"
	accountparam "github.com/royhairul/live-studio-api/internal/domains/account/params"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
)

type FinanceMetric struct {
	Total        int64   `json:"total"`
	Pending      int64   `json:"pending"`
	Paid         int64   `json:"paid"`
	PendingRatio float64 `json:"pending_ratio"`
	PaidRatio    float64 `json:"paid_ratio"`
}
type FinanceResponse struct {
	Metric FinanceMetric `json:"metric"`
	List   []FinanceItem `json:"list"`
}

type FinanceItem struct {
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

func NewFinanceResponse(items []FinanceItem) *FinanceResponse {
	var metric FinanceMetric

	for _, item := range items {
		metric.Total += item.Commission

		switch item.PaymentStatus {
		case "Sudah Dibayar":
			metric.Paid += item.Commission
		case "Menunggu Dibayar", "Menunggu Validasi":
			metric.Pending += item.Commission
		}
	}

	// Hitung rasio
	if metric.Total > 0 {
		metric.PaidRatio = float64(metric.Paid) / float64(metric.Total) * 100
		metric.PendingRatio = float64(metric.Pending) / float64(metric.Total) * 100
	}

	return &FinanceResponse{
		Metric: metric,
		List:   items,
	}
}

func NewFinanceItem(account accountparam.AccountResponse, commission shopeeparam.ShopeeFinanceCommissionList) *FinanceItem {
	return &FinanceItem{
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
