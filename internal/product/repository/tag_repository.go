package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type TagRepositoryImpl struct {
	DB *gorm.DB
}

func NewTagRepositoryImpl(db *gorm.DB) product.TagRepository {
	return &TagRepositoryImpl{DB: db}
}

func (t *TagRepositoryImpl) CreateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("tag_name = ?", tag.TagName).FirstOrCreate(tag)

		if result.RowsAffected == 0 {
			return gorm.ErrRegistered
		}

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (t *TagRepositoryImpl) UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("tag_id = ?", tag.TagID).Updates(tag)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (t *TagRepositoryImpl) DeleteTag(ctx context.Context, tag *domain.Tag) error {
	tx := t.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("tag_id = ?", tag.TagID).Delete(tag)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (t *TagRepositoryImpl) ReadTagById(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	if result := tx.Where("tag_id = ?", tag.TagID).First(&tag); result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

func (t *TagRepositoryImpl) ReadAllTag(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	tx := t.DB.WithContext(ctx)
	if result := tx.Find(&tags); result.Error != nil {
		return nil, result.Error
	}

	return tags, nil
}
