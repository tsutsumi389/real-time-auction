package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/tsutsumi389/real-time-auction/internal/domain"
	"github.com/tsutsumi389/real-time-auction/internal/repository"
)

var (
	ErrAdminNotFound   = errors.New("admin not found")
	ErrInvalidStatus   = errors.New("invalid status value")
	ErrInvalidRole     = errors.New("invalid role value")
	ErrInvalidPage     = errors.New("page must be greater than 0")
	ErrInvalidLimit    = errors.New("limit must be between 1 and 100")
	ErrInvalidSortMode = errors.New("invalid sort mode")
)

// AdminService handles admin-related business logic
type AdminService struct {
	adminRepo repository.AdminRepositoryInterface
}

// NewAdminService creates a new AdminService instance
func NewAdminService(adminRepo repository.AdminRepositoryInterface) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

// GetAdminByID retrieves a single admin by ID
func (s *AdminService) GetAdminByID(id int64) (*domain.Admin, error) {
	admin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	if admin == nil {
		return nil, ErrAdminNotFound
	}

	return admin, nil
}

// GetAdminList retrieves a paginated list of admins with filters
func (s *AdminService) GetAdminList(req *domain.AdminListRequest) (*domain.AdminListResponse, error) {
	// Validate and set defaults
	if err := s.validateAdminListRequest(req); err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.adminRepo.CountAdminsWithFilters(req)
	if err != nil {
		return nil, fmt.Errorf("failed to count admins: %w", err)
	}

	// Get admins
	admins, err := s.adminRepo.FindAdminsWithFilters(req)
	if err != nil {
		return nil, fmt.Errorf("failed to find admins: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	// Build response
	response := &domain.AdminListResponse{
		Admins: admins,
		Pagination: domain.Pagination{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// UpdateAdminStatus updates the status of an admin account
func (s *AdminService) UpdateAdminStatus(id int64, status domain.AdminStatus) (*domain.Admin, error) {
	// Validate status value
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	// Check if admin exists
	admin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	if admin == nil {
		return nil, ErrAdminNotFound
	}

	// Update status
	if err := s.adminRepo.UpdateAdminStatus(id, status); err != nil {
		return nil, fmt.Errorf("failed to update admin status: %w", err)
	}

	// Fetch updated admin
	updatedAdmin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated admin: %w", err)
	}

	return updatedAdmin, nil
}

// validateAdminListRequest validates and sets defaults for AdminListRequest
func (s *AdminService) validateAdminListRequest(req *domain.AdminListRequest) error {
	// Set default page if not specified
	if req.Page <= 0 {
		req.Page = 1
	}

	// Set default limit if not specified
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// Validate limit range
	if req.Limit > 100 {
		return ErrInvalidLimit
	}

	// Validate role if specified
	if req.Role != "" && !isValidRole(req.Role) {
		return ErrInvalidRole
	}

	// Validate statuses if specified
	for _, status := range req.Status {
		if !isValidStatus(status) {
			return ErrInvalidStatus
		}
	}

	// Validate sort mode if specified
	if req.Sort != "" && !isValidSortMode(req.Sort) {
		return ErrInvalidSortMode
	}

	return nil
}

// isValidRole checks if the role is valid
func isValidRole(role domain.AdminRole) bool {
	return role == domain.RoleSystemAdmin || role == domain.RoleAuctioneer
}

// isValidStatus checks if the status is valid
func isValidStatus(status domain.AdminStatus) bool {
	return status == domain.StatusActive ||
		status == domain.StatusSuspended ||
		status == domain.StatusDeleted
}

// isValidSortMode checks if the sort mode is valid
func isValidSortMode(sort string) bool {
	validSorts := []string{
		"id_asc", "id_desc",
		"email_asc", "email_desc",
		"role_asc", "role_desc",
		"status_asc", "status_desc",
		"created_at_asc", "created_at_desc",
	}

	for _, validSort := range validSorts {
		if sort == validSort {
			return true
		}
	}

	return false
}
