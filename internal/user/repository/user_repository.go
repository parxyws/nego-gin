package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	dbRes := dbTx.Where("email = ?", user.Email).Omit("avatar_url", "phone_num").FirstOrCreate(user)

	if dbRes.RowsAffected == 0 {
		return nil, gorm.ErrRegistered
	}

	if dbRes.Error != nil {
		return nil, fmt.Errorf("UserRepository.CreateUser - %w", dbRes.Error)
	}

	return user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	dbRes := dbTx.Model(&domain.User{}).Where("user_id = ?", user.UserID).Updates(user)

	if dbRes.RowsAffected == 0 {
		return nil, gorm.ErrInvalidData
	}

	if dbRes.Error != nil {
		return nil, fmt.Errorf("UserRepository.UpdateUser - %w", dbRes.Error)
	}

	return user, nil
}

func (repo *UserRepository) UpdateSingleColumnUser(ctx context.Context, entity *domain.User, column string) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	dbRes := dbTx.Model(&domain.User{}).Where("user_id = ?", entity.UserID).Update(column, entity.LastLogin)
	if dbRes.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if dbRes.Error != nil {
		return nil, fmt.Errorf("UserRepository.UpdateSingleColumnUser - %w", dbRes.Error)
	}

	return entity, nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, user *domain.User) error {
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Delete(user).Error; err != nil {
		return fmt.Errorf("UserRepository.DeleteUser - %w", err)
	}
	return nil
}

func (repo *UserRepository) ReadByUsername(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)

	if err := dbTx.Where("username = ?", user.Username).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByUsername - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepository) ReadByEmail(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Where("email = ?", user.Email).First(foundUser).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByEmail - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepository) ReadById(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Preload("UserRoles.Role").Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadById - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepository) ReadByIdMinimal(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.ReadByIdMinimal - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepository) ReadAllByRoles(ctx context.Context, userID string, sortOrderDesc bool, limit int, createdAt string) ([]domain.User, error) {
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
