package service

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/repository"

	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	orderparams "github.com/royhairul/live-studio-api/internal/domains/order/params"
	orderservice "github.com/royhairul/live-studio-api/internal/domains/order/service"
	productparams "github.com/royhairul/live-studio-api/internal/domains/product/params"
	productservice "github.com/royhairul/live-studio-api/internal/domains/product/service"
)

type TransactionServiceImpl struct {
	repository repository.TransactionRepository
	shopeeSvc  shopeeservice.ShopeeCheckoutService
	accountSvc accountservice.AccountService
	orderSvc   orderservice.OrderService
	productSvc productservice.ProductService
}

func NewTransactionService(repository repository.TransactionRepository, shopeeSvc shopeeservice.ShopeeCheckoutService, accountSvc accountservice.AccountService, orderSvc orderservice.OrderService, productSvc productservice.ProductService) TransactionService {
	return &TransactionServiceImpl{repository, shopeeSvc, accountSvc, orderSvc, productSvc}
}

// Create implements TransactionService.
func (s *TransactionServiceImpl) Create(req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error) {
	t, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.UTC)

	accounts, err := s.accountSvc.FindAll()
	if err != nil {
		return nil, err
	}

	var response []*params.CreatedTransactionResponse

	for _, account := range accounts {

		responseItem := params.CreatedTransactionResponse{
			AccountID:   account.ID,
			AccountName: account.Name,
			Total:       0,
		}

		if account.Platform == "Shopee" {
			transactions, err := s.shopeeSvc.GetByRangeTime(account.Cookie, startTime, endTime)
			if err != nil {
				return nil, err
			}

			for _, tx := range transactions.List {
				exists, err := s.repository.FindByUniqueID(tx.CheckoutID)
				if err != nil {
					fmt.Printf("error checking transaction %s: %v\n", tx.CheckoutID, err)
					continue
				}

				if exists != nil && exists.UniqueID == tx.CheckoutID {
					fmt.Printf("transaction %s was created, skip\n", tx.CheckoutID)
					continue // Sudah ada, skip
				}

				transaction := entity.Transaction{
					UniqueID:                        tx.CheckoutID,
					Status:                          tx.CheckoutStatus,
					EstimatedTotalCommission:        tx.EstimatedTotalCommission,
					EstimatedTotalCommissionWithMCN: tx.EstimatedTotalCommissionWithMCN,
					PurchaseTime:                    tx.PurchaseTime,
					CompleteTime:                    tx.CheckoutCompleteTime,
					AccountID:                       account.ID,
				}

				created, err := s.repository.Create(&transaction)
				if err != nil {
					return nil, err
				}

				for _, order := range tx.Orders {

					createdOrder, err := s.orderSvc.Create(orderparams.CreateOrderRequest{
						SerialNumber:  order.OrderSN,
						Status:        order.OrderStatus,
						CompleteTime:  order.CompleteTime,
						TransactionID: created.ID,
					})
					if err != nil {
						return nil, err
					}

					for _, product := range createdOrder.Products {
						_, err := s.productSvc.Create(productparams.CreateProductRequest{
							UniqueID: product.UniqueID,
							Name:     product.Name,
							ShopID:   product.ShopID,
							ShopName: product.ShopName,
						})
						if err != nil {
							return nil, err
						}
					}
				}

				if created != nil {
					responseItem.Total++
				}

			}
		}

		response = append(response, &responseItem)
	}

	return response, nil
}

// Update implements TransactionService.
func (s *TransactionServiceImpl) Update(id string, req params.UpdateTransactionRequest) (*params.TransactionResponse, error) {
	panic("unimplemented")
}

// FindAll implements TransactionService.
func (s *TransactionServiceImpl) FindAll() ([]*params.TransactionResponse, error) {
	transactions, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	grouped := make(map[uint]*params.TransactionResponse)

	for _, tx := range transactions {
		if _, exists := grouped[tx.AccountID]; !exists {
			grouped[tx.AccountID] = &params.TransactionResponse{
				AccountID:   tx.Account.ID,
				AccountName: tx.Account.Name,
				Total:       0,
				List:        []params.TransactionDetailResponse{},
			}
		}

		grouped[tx.AccountID].Total++
		grouped[tx.AccountID].List = append(grouped[tx.AccountID].List, *params.NewTransactionDetailResponse(tx))
	}

	var results []*params.TransactionResponse
	for _, res := range grouped {
		results = append(results, res)
	}

	return results, nil
}

// FindByID implements TransactionService.
func (s *TransactionServiceImpl) FindByID(id string) (*params.TransactionResponse, error) {
	panic("unimplemented")
}

// Delete implements TransactionService.
func (s *TransactionServiceImpl) Delete(id string) error {
	panic("unimplemented")
}
