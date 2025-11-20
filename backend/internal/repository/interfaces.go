package repository

import "github.com/tsutsumi389/real-time-auction/internal/domain"

// AdminRepositoryInterface defines the interface for admin repository operations
type AdminRepositoryInterface interface {
	FindByEmail(email string) (*domain.Admin, error)
	FindByID(id int64) (*domain.Admin, error)
	Create(admin *domain.Admin) error
	Update(admin *domain.Admin) error
	Delete(id int64) error
}
