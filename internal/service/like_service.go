package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"

	"github.com/google/uuid"
)

type LikeService struct {
	repo *repository.LikeRepository
}

func NewLikeService(repo *repository.LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

func (s *LikeService) LikePost(postID, userID uuid.UUID) error {
	like := &models.Like{
		PostID: postID,
		UserID: userID,
	}

	return s.repo.Create(like)
}

func (s *LikeService) LikeComment(commentID, userID uuid.UUID) error {
	like := &models.Like{
		CommentID: commentID,
		UserID:    userID,
	}

	return s.repo.Create(like)
}

func (s *LikeService) UnlikeComment(commentID, userID uuid.UUID) error {
	return s.repo.Delete(commentID, userID)
}

func (s *LikeService) UnlikePost(postID, userID uuid.UUID) error {
	return s.repo.Delete(postID, userID)
}

func (s *LikeService) GetLikesCount(commentID uuid.UUID) (int64, error) {
	return s.repo.CountLikes(commentID)
}
