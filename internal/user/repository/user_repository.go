package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, entity *domain.User) (*domain.User, error) {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("email = ?", entity.Email).Omit("avatar_url", "phone_num").FirstOrCreate(entity)

		if result.RowsAffected == 0 {
			return gorm.ErrRegistered
		}

		if result.Error != nil {
			return fmt.Errorf("UserRepository.CreateUser - %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, entity *domain.User) (*domain.User, error) {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.User{}).Where("id = ?", entity.UserID).Updates(entity)

		if result.RowsAffected == 0 {
			return gorm.ErrInvalidData
		}

		if result.Error != nil {
			return fmt.Errorf("UserRepository.CreateUser - %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *UserRepositoryImpl) DeleteUser(ctx context.Context, entity *domain.User) error {
	tx := u.DB.WithContext(ctx)
	return tx.Transaction(func(tx *gorm.DB) error {
		existingUser := new(domain.User)
		if err := tx.Where("email = ?", entity.Email).First(existingUser).Error; err != nil {
			return fmt.Errorf("UserRepository.DeleteUser - %w", err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(entity.PasswordHash)); err != nil {
			return fmt.Errorf("UserRepository.DeleteUser - %w", err)
		}

		if err := tx.Delete(existingUser).Error; err != nil {
			return fmt.Errorf("UserRepository.DeleteUser - %w", err)
		}

		return nil
	})
}

func (u *UserRepositoryImpl) ReadByUsername(ctx context.Context, entity *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	tx := u.DB.WithContext(ctx)

	if err := tx.Where("username = ?", entity.Username).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByUsername - %w", err)
	}

	return foundUser, nil
}

func (u *UserRepositoryImpl) ReadByEmail(ctx context.Context, entity *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	tx := u.DB.WithContext(ctx)
	if err := tx.Where("email = ?", entity.Email).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByEmail - %w", err)
	}

	return foundUser, nil
}

func (u *UserRepositoryImpl) ReadById(ctx context.Context, entity *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	tx := u.DB.WithContext(ctx)
	if err := tx.Model(&domain.User{}).Preload("Products").Preload("Auctions").Take(foundUser, "id = ?", entity.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadById - %w", err)
	}

	return foundUser, nil
}

func (u *UserRepositoryImpl) ReadAllByRoles(ctx context.Context, id string, sortOrderDesc bool, limit int, createdAt string) ([]domain.User, error) {
	var users []domain.User
	tx := u.DB.WithContext(ctx)

	orderPattern := clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: sortOrderDesc},
			{Column: clause.Column{Name: "id"}, Desc: sortOrderDesc},
		},
	}

	q := tx.Model(&domain.User{}).Where("is_verified IS NOT NULL").Clauses(orderPattern)

	if createdAt != "" {
		q = q.Where(`
			(created_at < ?)
			OR (created_at = ? AND id < ?)
		`, createdAt, createdAt, id)
	}

	if err := q.Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadAllByRoles - %w", err)
	}

	return users, nil
}
