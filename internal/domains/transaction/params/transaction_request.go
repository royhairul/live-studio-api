package params

type TransactionRequest struct {
	// TODO: add request fields
}

type CreateTransactionRequest struct {
	// TODO: add request fields
	Date string `json:"date" validate:"required"`
}

type UpdateTransactionRequest struct {
	// TODO: add request fields
}
