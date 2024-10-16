package repository

import (
	"GoVersi/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(like *models.Like) error {
	return r.db.Create(like).Error
}

func (r *LikeRepository) Delete(commentID, userID uuid.UUID) error {
	return r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&models.Like{}).Error
}

func (r *LikeRepository) CountLikes(commentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("comment_id = ?", commentID).Count(&count).Error
	return count, err
}
