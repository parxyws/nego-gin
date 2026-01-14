package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (repo *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Where("email = ?", user.Email).Omit("avatar_url", "phone_num").FirstOrCreate(user)

		if dbRes.RowsAffected == 0 {
			return gorm.ErrRegistered
		}

		if dbRes.Error != nil {
			return fmt.Errorf("UserRepository.CreateUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Model(&domain.User{}).Where("user_id = ?", user.UserID).Updates(user)

		if dbRes.RowsAffected == 0 {
			return gorm.ErrInvalidData
		}

		if dbRes.Error != nil {
			return fmt.Errorf("UserRepository.CreateUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepositoryImpl) UpdateSingleColumnUser(ctx context.Context, entity *domain.User, column string) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Model(&domain.User{}).Where("user_id = ?", entity.UserID).Update(column, entity.LastLogin)
		if dbRes.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if dbRes.Error != nil {
			return fmt.Errorf("UserRepository.UpdateSingleColumnUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *UserRepositoryImpl) DeleteUser(ctx context.Context, user *domain.User) error {
	dbTx := repo.DB.WithContext(ctx)
	return dbTx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(user).Error; err != nil {
			return fmt.Errorf("UserRepository.DeleteUser - %w", err)
		}
		return nil
	})
}

func (repo *UserRepositoryImpl) ReadByUsername(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)

	if err := dbTx.Where("username = ?", user.Username).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByUsername - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadByEmail(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Where("email = ?", user.Email).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByEmail - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadById(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Preload("Products").Preload("Auctions").Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadById - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadByIdMinimal(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByIdMinimal - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadAllByRoles(ctx context.Context, userID string, sortOrderDesc bool, limit int, createdAt string) ([]domain.User, error) {
	var users []domain.User
	dbTx := repo.DB.WithContext(ctx)

	orderPattern := clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: sortOrderDesc},
			{Column: clause.Column{Name: "user_id"}, Desc: sortOrderDesc},
		},
	}

	query := dbTx.Model(&domain.User{}).Where("is_verified IS NOT NULL").Clauses(orderPattern)

	if createdAt != "" {
		query = query.Where(`
			(created_at < ?)
			OR (created_at = ? AND user_id < ?)
		`, createdAt, createdAt, userID)
	}

	if err := query.Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadAllByRoles - %w", err)
	}

	return users, nil
}
