package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(post *models.Post, authorID uuid.UUID) error {
	// Validações podem ser feitas aqui, como:
	if post.Title == "" || post.Content == "" {
		return errors.New("title and content cannot be empty")
	}

	// Define o AuthorID
	post.AuthorID = authorID

	// Chama o repositório para criar o post
	return s.repo.Create(post)
}

func (s *PostService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *PostService) UpdatePost(postID uuid.UUID, updatedData *models.Post) (*models.Post, error) {
	// Busca o post existente
	existingPost, err := s.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// Atualiza os campos com as novas informações, mantendo o que não mudou
	existingPost.Title = updatedData.Title
	existingPost.Content = updatedData.Content
	existingPost.Topic = updatedData.Topic
	existingPost.UpdatedAt = time.Now()

	// Salva a atualização
	if err := s.repo.Update(existingPost); err != nil {
		return nil, err
	}
	return existingPost, nil
}

func (s *PostService) DeletePost(id uuid.UUID) error {
	// Pode fazer validações antes de deletar, se necessário
	return s.repo.Delete(id)
}
