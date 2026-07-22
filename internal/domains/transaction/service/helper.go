package service

import (
	"strings"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
)

func GroupTransactions(
	transactions []*entity.Transaction,
) []*params.TransactionGroupedResponse {
	grouped := make(map[uint]*params.TransactionGroupedResponse)

	for _, tx := range transactions {
		if _, exists := grouped[tx.AccountID]; !exists {
			grouped[tx.AccountID] = &params.TransactionGroupedResponse{
				AccountID:   tx.AccountID,
				AccountName: tx.Account.Name,
				Total:       0,
				Commission: params.Commission{
					Total:   0,
					Paid:    0,
					Pending: 0,
				},
				List: []params.TransactionList{},
			}
		}

		grouped[tx.AccountID].Total++
		grouped[tx.AccountID].Commission.Total += tx.TotalCommissionWithMCN

		// Dibandingkan tanpa peka huruf besar-kecil: Shopee mengirim
		// "Waiting for payment", sedangkan perbandingan persis sebelumnya
		// mencari "Waiting For Payment" sehingga tidak pernah cocok dan
		// paid/pending selalu nol.
		if strings.EqualFold(tx.Status, "Waiting For Payment") {
			grouped[tx.AccountID].Commission.Paid += tx.TotalCommissionWithMCN
		}
		if strings.EqualFold(tx.Status, "Pending") {
			grouped[tx.AccountID].Commission.Pending += tx.TotalCommissionWithMCN
		}
		grouped[tx.AccountID].List = append(grouped[tx.AccountID].List, *params.NewTransactionItem(tx))
	}

	var results []*params.TransactionGroupedResponse
	for _, res := range grouped {
		results = append(results, res)
	}
	return results
}
