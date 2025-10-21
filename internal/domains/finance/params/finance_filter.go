package params

import "time"

type FinanceFilter struct {
	UniqueID      *string
	StartDate     *time.Time
	EndDate       *time.Time
	StudioID      *string
	Status        *string
	PaymentMethod *string
}
