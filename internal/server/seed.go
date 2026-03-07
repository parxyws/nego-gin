package server

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Seeder struct {
	DB  *gorm.DB
	cfg *config.Config
	app *gin.Engine
}

func NewSeeder(db *gorm.DB, cfg *config.Config, app *gin.Engine) *Seeder {
	return &Seeder{DB: db, cfg: cfg, app: app}
}

func (s *Seeder) Seed() error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Seed permissions first (as they're required for roles)
		if err := s.seedPermission(tx); err != nil {
			return fmt.Errorf("failed to seed permissions: %w", err)
		}

		// Seed roles
		if err := s.seedRoles(tx); err != nil {
			return fmt.Errorf("failed to seed roles: %w", err)
		}

		// Seed superadmin user
		if err := s.seedSuperadmin(tx); err != nil {
			return fmt.Errorf("failed to seed superadmin: %w", err)
		}

		return nil
	})
}

func (s *Seeder) seedPermission(tx *gorm.DB) error {
	var existingPermissions []domain.Permission
	if err := tx.Find(&existingPermissions).Error; err != nil {
		return fmt.Errorf("failed to check existing permissions: %w", err)
	}

	if len(existingPermissions) > 0 {
		return nil
	}

	resourceMap := make(map[string]bool)
	accesses := []string{"create", "delete", "read", "update"}

	for _, r := range s.app.Routes() {
		parts := strings.Split(r.Path, "/")
		if len(parts) < 4 {
			continue
		}
		resourceMap[parts[3]] = true
	}

	var permissions []domain.Permission
	for resource := range resourceMap {
		for _, access := range accesses {
			permissions = append(permissions, domain.Permission{
				Name:        fmt.Sprintf("%s:%s", resource, access),
				Resource:    resource,
				Action:      access,
				Description: fmt.Sprintf("Permission to perform %s operations on %s resources.", access, resource),
			})
		}
	}

	// Only create permissions if we have any to create
	if len(permissions) > 0 {
		if err := tx.Create(permissions).Error; err != nil {
			return fmt.Errorf("failed to seed permissions: %w", err)
		}
	}

	return nil
}
func (s *Seeder) seedRoles(tx *gorm.DB) error {
	roles := []domain.Role{
		{RoleName: "superadmin", Description: "Full system access with all permissions"},
		{RoleName: "admin", Description: "Administrative access to manage users and content"},
		{RoleName: "customer", Description: "Standard user access"},
		{RoleName: "vendor", Description: "Standard seller access"},
	}

	for _, role := range roles {
		// Check if role already exists
		var existingRole domain.Role
		err := tx.Where("role_name = ?", role.RoleName).First(&existingRole).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Role doesn't exist, create it
			if err := tx.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.RoleName, err)
			}
		} else if err != nil {
			// Some other error occurred
			return fmt.Errorf("failed to check role %s: %w", role.RoleName, err)
		}
		// If no error, role already exists, skip
	}

	return nil
}

func (s *Seeder) seedSuperadmin(tx *gorm.DB) error {
	// Check if superadmin user already exists
	var existingUser domain.User
	err := tx.Where("username = ?", s.cfg.Admin.User).First(&existingUser).Error

	var user domain.User

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User doesn't exist, create new superadmin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.cfg.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0))

		user = domain.User{
			UserID:       id.String(),
			Email:        s.cfg.Admin.Email,
			Username:     s.cfg.Admin.User,
			PasswordHash: string(hashedPassword),
			FirstName:    s.cfg.Admin.FirstName,
			LastName:     s.cfg.Admin.LastName,
			IsVerified:   sql.NullTime{Time: time.Now(), Valid: true},
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		// Some other error occurred
		return fmt.Errorf("failed to check existing user: %w", err)
	} else {
		// User exists, use the existing user
		user = existingUser
	}

	// Get superadmin role
	var superadminRole domain.Role
	if err := tx.Where("role_name = ?", "superadmin").First(&superadminRole).Error; err != nil {
		return fmt.Errorf("superadmin role not found: %w", err)
	}

	// Check if user already has superadmin role
	var existingUserRole domain.UserRole
	err = tx.Where("user_id = ? AND role_id = ?", user.UserID, superadminRole.RoleID).
		First(&existingUserRole).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Superadmin role not assigned, create it
		userRole := &domain.UserRole{
			UserID: user.UserID,
			RoleID: superadminRole.RoleID,
		}

		if err := tx.Create(userRole).Error; err != nil {
			return fmt.Errorf("failed to assign superadmin role: %w", err)
		}
	} else if err != nil {
		// Some other error occurred
		return fmt.Errorf("failed to check existing user role: %w", err)
	}

	return nil
}
