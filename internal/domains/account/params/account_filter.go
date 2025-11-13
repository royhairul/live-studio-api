package params

type AccountFilter struct {
	ID       *string
	UniqueID *string
	StudioID *string

	IncludeDeleted bool
}
