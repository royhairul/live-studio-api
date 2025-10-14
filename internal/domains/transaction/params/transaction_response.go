package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"

	orderparams "github.com/royhairul/live-studio-api/internal/domains/order/params"
)

type Commission struct {
	Total   int64 `json:"total"`
	Pending int64 `json:"pending"`
	Paid    int64 `json:"paid"`
}

type TransactionResponse struct {
	// TODO: add response fields
	ID                              int64                       `json:"id"`
	AccountID                       uint                        `json:"account_id"`
	AccountName                     string                      `json:"account_name"`
	UniqueID                        string                      `json:"unique_id"`
	Status                          string                      `json:"status"`
	EstimatedTotalCommission        int64                       `json:"est_total_commission"`
	EstimatedTotalCommissionWithMCN int64                       `json:"est_total_commission_with_mcn"`
	PurchaseTime                    *time.Time                  `json:"purchase_time"`
	CompleteTime                    *time.Time                  `json:"complete_time"`
	Orders                          []orderparams.OrderResponse `json:"orders,omitempty"`
}

type TransactionGroupedResponse struct {
	// TODO: add response fields
	AccountID   uint                  `json:"account_id"`
	AccountName string                `json:"account_name"`
	Total       int                   `json:"total"`
	Commission  Commission            `json:"commission"`
	List        []TransactionResponse `json:"list"`
}

type CreatedTransactionResponse struct {
	AccountID      uint   `json:"account_id"`
	AccountName    string `json:"account_name"`
	NewTransaction int    `json:"new_transaction"`
}

func NewTransactionResponse(transaction *entity.Transaction) *TransactionResponse {
	return &TransactionResponse{
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
