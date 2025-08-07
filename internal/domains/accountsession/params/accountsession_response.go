package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
)

type AccountsessionResponse struct {
	// TODO: add response fields
	ID            uint       `json:"id"`
	AccountID     uint       `json:"account_id"`
	AccountName   string     `json:"account_name"`
	HostName      string     `json:"host_name"`
	Date          *time.Time `json:"date"`
	CheckIn       *time.Time `json:"check_in"`
	CheckOut      *time.Time `json:"check_out"`
	GMVSalesStart uint       `json:"gmv_sales_start"`
	GMVSalesEnd   uint       `json:"gmv_sales_end"`
	GMVPaidStart  uint       `json:"gmv_paid_start"`
	GMVPaidEnd    uint       `json:"gmv_paid_end"`
}

func NewAccountsessionResponse(accountsession *entity.Accountsession) *AccountsessionResponse {
	return &AccountsessionResponse{
		ID:            accountsession.ID,
		AccountID:     accountsession.Account.ID,
		AccountName:   accountsession.Account.Name,
		HostName:      accountsession.Attendance.Host.Name,
		Date:          accountsession.Attendance.Date,
		CheckIn:       accountsession.Attendance.CheckedInAt,
		CheckOut:      accountsession.Attendance.CheckedOutAt,
		GMVSalesStart: accountsession.GMVSalesStart,
		GMVSalesEnd:   accountsession.GMVSalesEnd,
		GMVPaidStart:  accountsession.GMVPaidStart,
		GMVPaidEnd:    accountsession.GMVPaidEnd,
	}
}
