package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/royhairul/live-studio-api/helpers/timehandler"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/entity"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/repository"
	"gorm.io/gorm"

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
	t, err := timehandler.ParseDate(req.Date)
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
					UniqueID:                        tx.CheckoutID,
					Status:                          tx.CheckoutStatus,
					EstimatedTotalCommission:        tx.EstimatedTotalCommission,
					EstimatedTotalCommissionWithMCN: tx.EstimatedTotalCommissionWithMCN,
					PurchaseTime:                    timehandler.ParseInt64Date(tx.PurchaseTime),
					CompleteTime:                    timehandler.ParseInt64Date(tx.CheckoutCompleteTime),
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

// FindAllByStatus implements TransactionService.
func (s *TransactionServiceImpl) FindAllByStatus(status string) ([]*params.TransactionResponse, error) {
	transactions, err := s.repository.FindAllByStatus(status)
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

// FindAllByAccount implements TransactionService.
func (s *TransactionServiceImpl) FindByAccount(accountID string) (*params.TransactionResponse, error) {
	account, err := s.accountSvc.FindById(accountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repository.FindAllByAccount(accountID)
	if err != nil {
		return nil, err
	}

	list := []params.TransactionDetailResponse{}

	result := &params.TransactionResponse{
		AccountID:   account.ID,
		AccountName: account.Name,
		Total:       len(transactions),
		List:        list,
	}

	return result, nil
}

// FindAllByAccountAndStatus implements TransactionService.
func (s *TransactionServiceImpl) FindAllByAccountAndStatus(accountID string, status string) (*params.TransactionResponse, error) {
	account, err := s.accountSvc.FindById(accountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repository.FindAllByAccount(accountID)
	if err != nil {
		return nil, err
	}

	txPaid, err := s.repository.FindAllByAccountStatus(accountID, "Waiting for payment")
	if err != nil {
		return nil, err
	}

	txPending, err := s.repository.FindAllByAccountStatus(accountID, "Pending")
	if err != nil {
		return nil, err
	}

	list := []params.TransactionDetailResponse{}

	totalAll := 0
	for _, tx := range transactions {
		list = append(list, *params.NewTransactionDetailResponse(tx))
		totalAll += int(tx.EstimatedTotalCommission)
	}

	totalPaid := 0
	for _, tx := range txPaid {
		totalPaid += int(tx.EstimatedTotalCommission)
	}

	totalPending := 0
	for _, tx := range txPending {
		totalPending += int(tx.EstimatedTotalCommission)
	}

	commission := params.Commission{
		Total:   uint(totalAll),
		Pending: uint(totalPending),
		Paid:    uint(totalPaid),
	}

	result := &params.TransactionResponse{
		AccountID:   account.ID,
		AccountName: account.Name,
		Total:       len(transactions),
		Commission:  commission,
		List:        list,
	}

	return result, nil
}

// FindAllByAccountStatusDate implements TransactionService.
func (s *TransactionServiceImpl) FindAllByAccountStatusDate(accountID string, status string, startDate *time.Time, endDate *time.Time) (*params.TransactionResponse, error) {
	account, err := s.accountSvc.FindById(accountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repository.FindAllByAccount(accountID)
	if err != nil {
		return nil, err
	}

	txPaid, err := s.repository.FindAllByAccountStatusDate(accountID, "Waiting for payment", startDate, endDate)
	if err != nil {
		return nil, err
	}

	txPaidCompleted, err := s.repository.FindAllByAccountStatusDate(accountID, "Complete", startDate, endDate)
	if err != nil {
		return nil, err
	}

	txPending, err := s.repository.FindAllByAccountStatusDate(accountID, "Pending", startDate, endDate)
	if err != nil {
		return nil, err
	}

	list := []params.TransactionDetailResponse{}

	totalAll := 0
	for _, tx := range transactions {
		list = append(list, *params.NewTransactionDetailResponse(tx))
		totalAll += int(tx.EstimatedTotalCommission)
	}

	totalPaid := 0
	for _, tx := range txPaid {
		totalPaid += int(tx.EstimatedTotalCommission)
	}
	for _, tx := range txPaidCompleted {
		totalPaid += int(tx.EstimatedTotalCommission)
	}

	totalPending := 0
	for _, tx := range txPending {
		totalPending += int(tx.EstimatedTotalCommission)
	}

	commission := params.Commission{
		Total:   uint(totalAll),
		Pending: uint(totalPending),
		Paid:    uint(totalPaid),
	}

	result := &params.TransactionResponse{
		AccountID:   account.ID,
		AccountName: account.Name,
		Total:       len(transactions),
		Commission:  commission,
		List:        list,
	}

	return result, nil
}

// FindAllByDate implements TransactionService.
func (s *TransactionServiceImpl) FindAllByDate(accountID string, startDate *time.Time, endDate *time.Time) ([]*params.TransactionResponse, error) {
	transactions, err := s.repository.FindAllByDate(startDate, endDate)
	if err != nil {
		return nil, err
	}

	grouped := make(map[uint]*params.TransactionResponse)

	for _, tx := range transactions {
		if _, exists := grouped[tx.AccountID]; !exists {
			grouped[tx.AccountID] = &params.TransactionResponse{
				AccountID:   tx.Account.ID,
				AccountName: tx.Account.Name,
				Commission: params.Commission{
					Total:   0,
					Pending: 0,
					Paid:    0,
				},
				Total: 0,
				List:  []params.TransactionDetailResponse{},
			}
		}

		grouped[tx.AccountID].Total++
		grouped[tx.AccountID].List = append(grouped[tx.AccountID].List, *params.NewTransactionDetailResponse(tx))

		// Commission
		if tx.Status == "Waiting for payment" {
			grouped[tx.AccountID].Commission.Paid += uint(tx.EstimatedTotalCommissionWithMCN)
		}
		if tx.Status == "Pending" {
			grouped[tx.AccountID].Commission.Pending += uint(tx.EstimatedTotalCommissionWithMCN)
		}
		grouped[tx.AccountID].Commission.Total += uint(tx.EstimatedTotalCommissionWithMCN)
	}

	results := []*params.TransactionResponse{}
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
