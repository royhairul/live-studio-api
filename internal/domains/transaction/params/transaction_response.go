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
		ID:                       transaction.ID,
		UniqueID:                 transaction.UniqueID,
		AccountID:                transaction.Account.ID,
		AccountName:              transaction.Account.Name,
		Status:                   transaction.Status,
		EstimatedTotalCommission: transaction.EstimatedTotalCommission,
		PurchaseTime:             transaction.PurchaseTime,
		CompleteTime:             transaction.CompleteTime,
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
