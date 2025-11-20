package repository

import (
	"errors"

	"github.com/tsutsumi389/real-time-auction/backend/internal/domain"
	"gorm.io/gorm"
)

// AdminRepository handles database operations for Admin entities
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository creates a new AdminRepository instance
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByEmail finds an admin by email address
func (r *AdminRepository) FindByEmail(email string) (*domain.Admin, error) {
	var admin domain.Admin
	result := r.db.Where("email = ?", email).First(&admin)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &admin, nil
}

// FindByID finds an admin by ID
func (r *AdminRepository) FindByID(id int64) (*domain.Admin, error) {
	var admin domain.Admin
	result := r.db.First(&admin, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error condition)
		}
		return nil, result.Error
	}

	return &admin, nil
}

// Create creates a new admin record
func (r *AdminRepository) Create(admin *domain.Admin) error {
	return r.db.Create(admin).Error
}

// Update updates an existing admin record
func (r *AdminRepository) Update(admin *domain.Admin) error {
	return r.db.Save(admin).Error
}

// Delete soft deletes an admin by setting status to deleted
func (r *AdminRepository) Delete(id int64) error {
	return r.db.Model(&domain.Admin{}).
		Where("id = ?", id).
		Update("status", domain.StatusDeleted).Error
}
