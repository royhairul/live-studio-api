package params

import (
	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"

	orderparams "github.com/royhairul/live-studio-api/internal/domains/order/params"
)

type TransactionResponse struct {
	// TODO: add response fields
	AccountID   uint                        `json:"account_id"`
	AccountName string                      `json:"account_name"`
	Total       int                         `json:"total"`
	List        []TransactionDetailResponse `json:"list"`
}

type CreatedTransactionResponse struct {
	// TODO: add response fields
	AccountID   uint   `json:"account_id"`
	AccountName string `json:"account_name"`
	Total       int    `json:"total"`
}

type TransactionDetailResponse struct {
	// TODO: add response fields
	ID                              int64                       `json:"id"`
	UniqueID                        string                      `json:"unique_id"`
	Status                          string                      `json:"status"`
	EstimatedTotalCommission        int64                       `json:"est_total_comission"`
	EstimatedTotalCommissionWithMCN int64                       `json:"est_total_comission_with_mcn"`
	PurchaseTime                    int64                       `json:"purchase_time"`
	CompleteTime                    int64                       `json:"complete_time"`
	Orders                          []orderparams.OrderResponse `json:"orders"`
}

func NewTransactionDetailResponse(transaction *entity.Transaction) *TransactionDetailResponse {
	return &TransactionDetailResponse{
		ID:                       transaction.ID,
		UniqueID:                 transaction.UniqueID,
		Status:                   transaction.Status,
		EstimatedTotalCommission: transaction.EstimatedTotalCommission,
		PurchaseTime:             transaction.PurchaseTime,
		CompleteTime:             transaction.CompleteTime,
	}
}
