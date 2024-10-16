package repository

import (
	"GoVersi/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) FindByPostID(postID uuid.UUID) ([]*models.Comment, error) {
	var comments []*models.Comment
	if err := r.db.Where("post_id = ?", postID.String()).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Post{}, "id = ?", id).Error
}
