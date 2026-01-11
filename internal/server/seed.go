package server

import (
	"errors"
	"fmt"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seeder struct {
	DB  *gorm.DB
	cfg *config.Config
}

func NewSeeder(db *gorm.DB, cfg *config.Config) *Seeder {
	return &Seeder{DB: db, cfg: cfg}
}

func (s *Seeder) Seed() error {
	// Check if superadmin user already exists before doing anything
	var existingUser domain.User
	err := s.DB.Where("username = ?", s.cfg.Admin.User).First(&existingUser).Error

	if err == nil {
		// User exists, check if they have superadmin role
		var userRoleCount int64
		s.DB.Model(&domain.UserRole{}).
			Joins("JOIN roles ON user_roles.role_id = roles.role_id").
			Where("user_roles.user_id = ? AND roles.role_name = ?",
				existingUser.UserID, "superadmin").
			Count(&userRoleCount)

		if userRoleCount > 0 {
			// Superadmin exists with proper role, skip all seeding
			return nil
		}
		// User exists but doesn't have superadmin role, continue seeding
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	// User doesn't exist or doesn't have superadmin role, proceed with seeding

	// Start transaction for all seeding operations
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Seed roles first
		if err := s.seedRoles(tx); err != nil {
			return fmt.Errorf("failed to seed roles: %w", err)
		}

		// Seed user
		if err := s.seedSuperadmin(tx, existingUser.UserID); err != nil {
			return fmt.Errorf("failed to seed user: %w", err)
		}

		return nil
	})
}

func (s *Seeder) seedRoles(tx *gorm.DB) error {
	roles := []domain.Role{
		{RoleName: "superadmin", Description: "Full system access with all permissions"},
		{RoleName: "admin", Description: "Administrative access to manage users and content"},
		{RoleName: "customer", Description: "Standard user access"},
		{RoleName: "vendor", Description: "Standard seller access"},
	}

	for _, role := range roles {
		// Use Upsert approach (create if not exists, update if exists)
		result := tx.Where(domain.Role{RoleName: role.RoleName}).
			Assign(domain.Role{Description: role.Description}).
			FirstOrCreate(&role)

		if result.Error != nil {
			return fmt.Errorf("failed to seed role %s: %w", role.RoleName, result.Error)
		}
	}

	return nil
}

func (s *Seeder) seedSuperadmin(tx *gorm.DB, existingUserID string) error {
	// Get superadmin role
	var superadminRole domain.Role
	if err := tx.Where("role_name = ?", "superadmin").First(&superadminRole).Error; err != nil {
		return fmt.Errorf("superadmin role not found: %w", err)
	}

	var user domain.User

	if existingUserID != "" {
		// User exists but doesn't have superadmin role, get the existing user
		if err := tx.Where("user_id = ?", existingUserID).First(&user).Error; err != nil {
			return fmt.Errorf("failed to find existing user: %w", err)
		}
	} else {
		// Create new superadmin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.cfg.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		user = domain.User{
			Email:        s.cfg.Admin.Email,
			Username:     s.cfg.Admin.User,
			PasswordHash: string(hashedPassword),
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Check if user already has superadmin role
	var existingUserRole domain.UserRole
	err := tx.Where("user_id = ? AND role_id = ?", user.UserID, superadminRole.RoleID).
		First(&existingUserRole).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Assign superadmin role
		userRole := &domain.UserRole{
			UserID: user.UserID,
			RoleID: superadminRole.RoleID,
		}

		if err := tx.Create(userRole).Error; err != nil {
			return fmt.Errorf("failed to assign superadmin role: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing user role: %w", err)
	}

	// Also ensure customer role is assigned (if needed)
	var customerRole domain.Role
	if err := tx.Where("role_name = ?", "customer").First(&customerRole).Error; err != nil {
		return fmt.Errorf("customer role not found: %w", err)
	}

	var existingCustomerRole domain.UserRole
	err = tx.Where("user_id = ? AND role_id = ?", user.UserID, customerRole.RoleID).
		First(&existingCustomerRole).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		customerUserRole := &domain.UserRole{
			UserID: user.UserID,
			RoleID: customerRole.RoleID,
		}

		if err := tx.Create(customerUserRole).Error; err != nil {
			return fmt.Errorf("failed to assign customer role: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existing customer role: %w", err)
	}

	return nil
}
