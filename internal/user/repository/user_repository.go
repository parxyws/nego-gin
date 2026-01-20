package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/sirupsen/logrus"
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
			logrus.WithFields(logrus.Fields{"function": "UserRepository.CreateUser", "email": user.Email}).Warn("user already exists")
			return gorm.ErrRegistered
		}

		if dbRes.Error != nil {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.CreateUser", "email": user.Email}).Errorf("failed to create user: %v", dbRes.Error)
			return fmt.Errorf("UserRepository.CreateUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "UserRepository.CreateUser", "user_id": user.UserID}).Info("user created successfully")
	return user, nil
}

func (repo *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Model(&domain.User{}).Where("user_id = ?", user.UserID).Updates(user)

		if dbRes.RowsAffected == 0 {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateUser", "user_id": user.UserID}).Warn("user not found for update")
			return gorm.ErrInvalidData
		}

		if dbRes.Error != nil {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateUser", "user_id": user.UserID}).Errorf("failed to update user: %v", dbRes.Error)
			return fmt.Errorf("UserRepository.UpdateUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateUser", "user_id": user.UserID}).Info("user updated successfully")
	return user, nil
}

func (repo *UserRepositoryImpl) UpdateSingleColumnUser(ctx context.Context, entity *domain.User, column string) (*domain.User, error) {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Model(&domain.User{}).Where("user_id = ?", entity.UserID).Update(column, entity.LastLogin)
		if dbRes.RowsAffected == 0 {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateSingleColumnUser", "user_id": entity.UserID, "column": column}).Warn("user not found for single column update")
			return gorm.ErrRecordNotFound
		}

		if dbRes.Error != nil {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateSingleColumnUser", "user_id": entity.UserID, "column": column}).Errorf("failed to update single column: %v", dbRes.Error)
			return fmt.Errorf("UserRepository.UpdateSingleColumnUser - %w", dbRes.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "UserRepository.UpdateSingleColumnUser", "user_id": entity.UserID, "column": column}).Info("user single column updated successfully")
	return entity, nil
}

func (repo *UserRepositoryImpl) DeleteUser(ctx context.Context, user *domain.User) error {
	dbTx := repo.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(user).Error; err != nil {
			logrus.WithFields(logrus.Fields{"function": "UserRepository.DeleteUser", "user_id": user.UserID}).Errorf("failed to delete user: %v", err)
			return fmt.Errorf("UserRepository.DeleteUser - %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"function": "UserRepository.DeleteUser", "user_id": user.UserID}).Info("user deleted successfully")
	return nil
}

func (repo *UserRepositoryImpl) ReadByUsername(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)

	if err := dbTx.Where("username = ?", user.Username).First(foundUser).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserRepository.ReadByUsername", "username": user.Username}).Warnf("failed to read user by username: %v", err)
		return nil, fmt.Errorf("UserRepository.ReadByUsername - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadByEmail(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Where("email = ?", user.Email).First(foundUser).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserRepository.ReadByEmail", "email": user.Email}).Warnf("failed to read user by email: %v", err)
		return nil, fmt.Errorf("UserRepository.ReadByEmail - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadById(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Preload("Products").Preload("Auctions").Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserRepository.ReadById", "user_id": user.UserID}).Warnf("failed to read user by id: %v", err)
		return nil, fmt.Errorf("UserRepository.ReadById - %w", err)
	}

	return foundUser, nil
}

func (repo *UserRepositoryImpl) ReadByIdMinimal(ctx context.Context, user *domain.User) (*domain.User, error) {
	foundUser := new(domain.User)
	dbTx := repo.DB.WithContext(ctx)
	if err := dbTx.Model(&domain.User{}).Take(foundUser, "user_id = ?", user.UserID).Error; err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserRepository.ReadByIdMinimal", "user_id": user.UserID}).Warnf("failed to read user minimal by id: %v", err)
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
		logrus.WithFields(logrus.Fields{"function": "UserRepository.ReadAllByRoles", "limit": limit}).Errorf("failed to read all by roles: %v", err)
		return nil, fmt.Errorf("UserRepository.ReadAllByRoles - %w", err)
	}

	return users, nil
}
