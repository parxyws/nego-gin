package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/product"
	"github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

type TagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) product.TagRepository {
	return &TagRepository{DB: db}
}

func (t *TagRepository) CreateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	result := tx.Where("tag_name = ?", tag.TagName).FirstOrCreate(tag)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return tag, nil // It already exists (FirstOrCreate), so we just return the found tag
	}

	return tag, nil
}

func (t *TagRepository) UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	result := tx.Where("tag_id = ?", tag.TagID).Updates(tag)

	if result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

func (t *TagRepository) DeleteTag(ctx context.Context, tag *domain.Tag) error {
	tx := t.DB.WithContext(ctx)

	result := tx.Where("tag_id = ?", tag.TagID).Delete(tag)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *TagRepository) ReadTagById(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	tx := t.DB.WithContext(ctx)

	if result := tx.Where("tag_id = ?", tag.TagID).First(&tag); result.Error != nil {
		return nil, result.Error
	}

	return tag, nil
}

func (t *TagRepository) ReadAllTag(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	tx := t.DB.WithContext(ctx)
	if result := tx.Find(&tags); result.Error != nil {
		return nil, result.Error
	}

	return tags, nil
}
