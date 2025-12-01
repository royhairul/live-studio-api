package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/repository"
	"github.com/royhairul/live-studio-api/internal/pkg/timehandler"
	"gorm.io/gorm"

	orderparams "github.com/royhairul/live-studio-api/internal/domains/order/params"
	productparams "github.com/royhairul/live-studio-api/internal/domains/product/params"

	shopeeservice "github.com/royhairul/live-studio-api/internal/clients/shopee/service"
	accountservice "github.com/royhairul/live-studio-api/internal/domains/account/service"
	orderservice "github.com/royhairul/live-studio-api/internal/domains/order/service"
	productservice "github.com/royhairul/live-studio-api/internal/domains/product/service"
)

type TransactionServiceImpl struct {
	repository repository.TransactionRepository
	shopeeSvc  shopeeservice.ShopeeCheckoutService
	accountSvc accountservice.AccountService
	orderSvc   orderservice.OrderService
	productSvc productservice.ProductService
	options    params.TransactionFilter
}

func NewTransactionService(
	repository repository.TransactionRepository,
	shopeeSvc shopeeservice.ShopeeCheckoutService,
	accountSvc accountservice.AccountService,
	orderSvc orderservice.OrderService,
	productSvc productservice.ProductService,
) TransactionService {
	return &TransactionServiceImpl{
		repository: repository,
		shopeeSvc:  shopeeSvc,
		accountSvc: accountSvc,
		orderSvc:   orderSvc,
		productSvc: productSvc,
		options:    params.TransactionFilter{},
	}
}

// WithID implements TransactionService.
func (s *TransactionServiceImpl) WithID(id string) TransactionService {
	instance := *s
	instance.options.ID = &id
	return &instance
}

// WithStatus implements TransactionService.
func (s *TransactionServiceImpl) WithStatus(status string) TransactionService {
	instance := *s
	instance.options.Status = &status
	return &instance
}

// WithAccountID implements TransactionService.
func (s *TransactionServiceImpl) WithAccountID(accountID string) TransactionService {
	instance := *s
	instance.options.AccountID = &accountID
	return &instance
}

// WithStudioID implements TransactionService.
func (s *TransactionServiceImpl) WithStudioID(studioID string) TransactionService {
	instance := *s
	instance.options.StudioID = &studioID
	return &instance
}

// WithDate implements TransactionService.
func (s *TransactionServiceImpl) WithDate(startTime time.Time, endTime time.Time) TransactionService {
	instance := *s
	instance.options.StartTime = &startTime
	instance.options.EndTime = &endTime
	return &instance
}

// Create implements TransactionService.
func (s *TransactionServiceImpl) Create(ctx context.Context, req params.CreateTransactionRequest) ([]*params.CreatedTransactionResponse, error) {
	t, err := timehandler.ParseDate(req.Date)
	if err != nil {
		return nil, err
	}

	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	endTime := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.UTC)

	accounts, err := s.accountSvc.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []*params.CreatedTransactionResponse

	for _, account := range accounts {
		responseItem := params.CreatedTransactionResponse{
			AccountID:      account.ID,
			AccountName:    account.Name,
			NewTransaction: 0,
		}

		if account.Platform == "Shopee" {
			transactions, err := s.shopeeSvc.GetByRangeTime(account.Cookie, startTime, endTime)
			if err != nil {
				return nil, err
			}

			for _, tx := range transactions.List {
				exists, err := s.repository.FindOne(ctx, params.TransactionFilter{UniqueID: &tx.CheckoutID})
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						log.Println("not found true")
						exists = nil
					} else {
						fmt.Printf("error checking transaction %s: %v\n", tx.CheckoutID, err)
						continue
					}
				}

				if exists != nil && exists.UniqueID == tx.CheckoutID {
					fmt.Printf("transaction %s was created, skip\n", tx.CheckoutID)
					continue // Sudah ada, skip
				}

				transaction := entity.Transaction{
					UniqueID:               tx.CheckoutID,
					Status:                 tx.CheckoutStatus,
					TotalPurchase:          tx.TotalBrandCommission,
					TotalCommission:        tx.EstimatedTotalCommission,
					TotalCommissionWithMCN: tx.EstimatedTotalCommissionWithMCN,
					PurchaseTime:           timehandler.ParseInt64Date(tx.PurchaseTime),
					CompleteTime:           timehandler.ParseInt64Date(tx.CheckoutCompleteTime),
					AccountID:              account.ID,
				}

				created, err := s.repository.Create(ctx, &transaction)
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
					responseItem.NewTransaction++
				}
			}
		}
		response = append(response, &responseItem)
	}
	return response, nil
}

// Update implements TransactionService.
func (s *TransactionServiceImpl) Update(ctx context.Context, id string, req params.UpdateTransactionRequest) (*params.TransactionList, error) {
	panic("unimplemented")
}

// FindAll implements TransactionService.
func (s *TransactionServiceImpl) FindAll(ctx context.Context) (*params.TransactionResponse, error) {
	transactions, err := s.repository.FindAll(ctx, s.options)
	if err != nil {
		return nil, err
	}

	var list []params.TransactionList
	for _, tx := range transactions {
		list = append(list, *params.NewTransactionItem(tx))
	}
	return params.NewTransactionResponse(list), nil
}

// FindAllGrouped implements TransactionService.
func (s *TransactionServiceImpl) FindAllGrouped(ctx context.Context) ([]*params.TransactionGroupedResponse, error) {
	transactions, err := s.repository.FindAll(ctx, s.options)
	if err != nil {
		return nil, err
	}
	results := GroupTransactions(transactions)
	return results, nil
}

// FindOne implements TransactionService.
func (s *TransactionServiceImpl) FindOne(ctx context.Context) (*params.TransactionList, error) {
	transaction, err := s.repository.FindOne(ctx, s.options)
	if err != nil {
		return nil, err
	}
	result := params.NewTransactionItem(transaction)
	return result, nil
}

// Delete implements TransactionService.
func (s *TransactionServiceImpl) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetTotalCommission implements TransactionService.
func (s *TransactionServiceImpl) GetTotalCommission(ctx context.Context) (*params.Commission, error) {
	transactions, err := s.FindAllGrouped(ctx)
	if err != nil {
		return &params.Commission{}, err
	}
	var result params.Commission
	for _, tx := range transactions {
		result.Total += tx.Commission.Total
		result.Paid += tx.Commission.Paid
		result.Pending += tx.Commission.Pending
	}
	return &result, nil
}
