package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post *models.Post) error {
	// Aqui você pode colocar validações, regras de negócio, etc.
	return s.repo.Create(post)
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *PostService) UpdatePost(post *models.Post) error {
	return s.repo.Update(post)
}

func (s *PostService) DeletePost(id int) error {
	return s.repo.Delete(id)
}
