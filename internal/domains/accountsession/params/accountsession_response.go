package params

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/domains/accountsession/entity"
)

type AccountsessionResponse struct {
	// TODO: add response fields
	ID          uint       `json:"id"`
	AccountID   uint       `json:"account_id"`
	AccountName string     `json:"account_name"`
	StudioID    uint       `json:"studio_id"`
	StudioName  string     `json:"studio_name"`
	HostName    string     `json:"host_name"`
	Date        *time.Time `json:"date"`
	CheckIn     *time.Time `json:"check_in"`
	CheckOut    *time.Time `json:"check_out"`
	Duration    int64      `json:"duration"`
	GMVSales    uint       `json:"gmv_sales"`
	GMVPaid     uint       `json:"gmv_paid"`
}

func NewAccountsessionResponse(accountsession *entity.Accountsession) *AccountsessionResponse {
	return &AccountsessionResponse{
		ID:          accountsession.ID,
		AccountID:   accountsession.Account.ID,
		AccountName: accountsession.Account.Name,
		HostName:    accountsession.Attendance.Host.Name,
		Date:        accountsession.Attendance.Date,
		Duration:    accountsession.Attendance.Duration(),
		CheckIn:     accountsession.Attendance.CheckedInAt,
		CheckOut:    accountsession.Attendance.CheckedOutAt,
		GMVSales:    accountsession.TotalSales(),
		GMVPaid:     accountsession.TotalPaid(),
		StudioID:    accountsession.Studio.ID,
		StudioName:  accountsession.Studio.Name,
	}
}
