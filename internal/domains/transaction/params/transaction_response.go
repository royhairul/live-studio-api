package params

import (
	"time"

	orderparams "github.com/royhairul/live-studio-api/internal/domains/order/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
)

// Bagian existing
type Commission struct {
	Total   int64 `json:"total"`
	Pending int64 `json:"pending"`
	Paid    int64 `json:"paid"`
}

<<<<<<< HEAD
type CommissionMetric struct {
	Total        int64 `json:"total"`
	Pending      int64 `json:"pending"`
	Paid         int64 `json:"paid"`
	PendingRatio int   `json:"pending_ratio"`
	PaidRatio    int   `json:"paid_ratio"`
}

type TransactionMetric struct {
	TotalPurchase          int64 `json:"total_purchase"`
	TotalCommissionWithMCN int64 `json:"total_commission_with_mcn"`
	TotalCommission        int64 `json:"total_commission"`
}

type TransactionResponse struct {
	Metric TransactionMetric `json:"metric"`
	List   []TransactionList `json:"list"`
}

type TransactionList struct {
	ID                     int64                       `json:"id"`
	AccountID              uint                        `json:"account_id"`
	AccountName            string                      `json:"account_name"`
	AccountStudio          string                      `json:"account_studio"`
	UniqueID               string                      `json:"unique_id"`
	Status                 string                      `json:"status"`
	TotalPurchase          int64                       `json:"total_purchase"`
	TotalCommissionWithMCN int64                       `json:"total_commission_with_mcn"`
	TotalCommission        int64                       `json:"total_commission"`
	PurchaseTime           *time.Time                  `json:"purchase_time"`
	CompleteTime           *time.Time                  `json:"complete_time,omitempty"`
	Orders                 []orderparams.OrderResponse `json:"orders,omitempty"`
=======
// 🔹 Tambahan: Rasio komisi
type CommissionRatio struct {
	Pending float64 `json:"pending_ratio"`
	Paid    float64 `json:"paid_ratio"`
}

type TransactionResponse struct {
	Commission      Commission        `json:"commission"`
	CommissionRatio CommissionRatio   `json:"commission_ratio"`
	List            []TransactionList `json:"list"`
}

type TransactionList struct {
	ID                              int64                       `json:"id"`
	AccountID                       uint                        `json:"account_id"`
	AccountName                     string                      `json:"account_name"`
	UniqueID                        string                      `json:"unique_id"`
	Status                          string                      `json:"status"`
	EstimatedTotalCommission        int64                       `json:"est_total_commission"`
	EstimatedTotalCommissionWithMCN int64                       `json:"est_total_commission_with_mcn"`
	PurchaseTime                    *time.Time                  `json:"purchase_time"`
	CompleteTime                    *time.Time                  `json:"complete_time,omitempty"`
	Orders                          []orderparams.OrderResponse `json:"orders,omitempty"`
>>>>>>> 05936c56bc49c4c6ae6fe3260a47122fce7656fe
}

type TransactionGroupedResponse struct {
	AccountID   uint              `json:"account_id"`
	AccountName string            `json:"account_name"`
	Total       int               `json:"total"`
	Commission  Commission        `json:"commission"`
	List        []TransactionList `json:"list"`
}

type CreatedTransactionResponse struct {
	AccountID      uint   `json:"account_id"`
	AccountName    string `json:"account_name"`
	NewTransaction int    `json:"new_transaction"`
}

func NewTransactionItem(transaction *entity.Transaction) *TransactionList {
	return &TransactionList{
<<<<<<< HEAD
		ID:                     transaction.ID,
		UniqueID:               transaction.UniqueID,
		AccountID:              transaction.Account.ID,
		AccountStudio:          transaction.Account.Studio.Name,
		AccountName:            transaction.Account.Name,
		Status:                 transaction.Status,
		TotalPurchase:          transaction.TotalPurchase,
		TotalCommission:        transaction.TotalCommission,
		TotalCommissionWithMCN: transaction.TotalCommissionWithMCN,
		PurchaseTime:           transaction.PurchaseTime,
		CompleteTime:           transaction.CompleteTime,
	}
}

func NewTransactionResponse(list []TransactionList) *TransactionResponse {
	var totalPurchase int64
	var totalCommission int64
	var totalCommissionWithMCN int64

	for _, item := range list {
		totalPurchase += item.TotalPurchase
		totalCommission += item.TotalCommission
		totalCommissionWithMCN += item.TotalCommissionWithMCN
	}

	return &TransactionResponse{
		Metric: TransactionMetric{
			TotalPurchase:          totalPurchase,
			TotalCommission:        totalCommission,
			TotalCommissionWithMCN: totalCommissionWithMCN,
		},
		List: list,
=======
		ID:                       transaction.ID,
		UniqueID:                 transaction.UniqueID,
		AccountID:                transaction.Account.ID,
		AccountName:              transaction.Account.Name,
		Status:                   transaction.Status,
		EstimatedTotalCommission: transaction.EstimatedTotalCommission,
		PurchaseTime:             transaction.PurchaseTime,
		CompleteTime:             transaction.CompleteTime,
>>>>>>> 05936c56bc49c4c6ae6fe3260a47122fce7656fe
	}
}

// 🔹 Helper untuk buat TransactionResponse dengan rasio
func NewTransactionResponse(commission Commission, list []TransactionList) *TransactionResponse {
	total := commission.Total
	var pendingRatio, paidRatio float64

	if total > 0 {
		pendingRatio = float64(commission.Pending) / float64(total) * 100
		paidRatio = float64(commission.Paid) / float64(total) * 100
	}

	return &TransactionResponse{
		Commission: commission,
		CommissionRatio: CommissionRatio{
			Pending: round2(pendingRatio),
			Paid:    round2(paidRatio),
		},
		List: list,
	}
}

func round2(v float64) float64 {
	return float64(int(v*100)) / 100
}
