package repository

type UserRepositoryImpl struct {
	// TODO: add database instance
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		// TODO: initialize dependencies
	}
}