package repositories

import (
	repositories "user-service/repositories/user"

	"gorm.io/gorm"
)

type Regsitry struct {
	db *gorm.DB
}

type IRepositoryRegistry interface {
	GetUser() repositories.IUserRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepositoryRegistry {
	return &Regsitry{db: db}
}

func (r *Regsitry) GetUser() repositories.IUserRepository {
	return repositories.NewUserRepository(r.db)
}