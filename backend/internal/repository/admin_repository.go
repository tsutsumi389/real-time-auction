package repository

import (
	"errors"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
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

// FindAdminsWithFilters retrieves admins with filters, sorting, and pagination
func (r *AdminRepository) FindAdminsWithFilters(req *domain.AdminListRequest) ([]domain.Admin, error) {
	var admins []domain.Admin

	query := r.db.Model(&domain.Admin{})

	// Apply keyword filter (email search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("email ILIKE ?", "%"+req.Keyword+"%")
	}

	// Apply role filter
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// Apply status filter (multiple statuses allowed)
	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	} else {
		// Default: show only active and suspended (exclude deleted)
		query = query.Where("status IN ?", []domain.AdminStatus{domain.StatusActive, domain.StatusSuspended})
	}

	// Apply sorting
	switch req.Sort {
	case "id_asc":
		query = query.Order("id ASC")
	case "id_desc":
		query = query.Order("id DESC")
	case "email_asc":
		query = query.Order("email ASC")
	case "email_desc":
		query = query.Order("email DESC")
	case "created_at_asc":
		query = query.Order("created_at ASC")
	case "created_at_desc":
		query = query.Order("created_at DESC")
	default:
		// Default sort: id ASC
		query = query.Order("id ASC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit).Offset(offset)

	result := query.Find(&admins)
	if result.Error != nil {
		return nil, result.Error
	}

	return admins, nil
}

// CountAdminsWithFilters counts the total number of admins matching the filters
func (r *AdminRepository) CountAdminsWithFilters(req *domain.AdminListRequest) (int64, error) {
	var count int64

	query := r.db.Model(&domain.Admin{})

	// Apply keyword filter (email search with ILIKE)
	if req.Keyword != "" {
		query = query.Where("email ILIKE ?", "%"+req.Keyword+"%")
	}

	// Apply role filter
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// Apply status filter (multiple statuses allowed)
	if len(req.Status) > 0 {
		query = query.Where("status IN ?", req.Status)
	} else {
		// Default: show only active and suspended (exclude deleted)
		query = query.Where("status IN ?", []domain.AdminStatus{domain.StatusActive, domain.StatusSuspended})
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// UpdateAdminStatus updates the status of an admin account
func (r *AdminRepository) UpdateAdminStatus(id int64, status domain.AdminStatus) error {
	return r.db.Model(&domain.Admin{}).
		Where("id = ?", id).
		Update("status", status).Error
}
