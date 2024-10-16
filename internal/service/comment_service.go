package services

import (
	"GoVersi/internal/models"
	"GoVersi/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// Ajustando a função CreateComment para incluir o PostID
func (s *CommentService) CreateComment(content, imageURL string, postID, authorID uuid.UUID) (*models.Comment, error) {
	comment := &models.Comment{
		Content:  content,
		ImageURL: imageURL,
		PostID:   postID.String(), // Associando o comentário ao PostID
		AuthorID: authorID,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// Função para buscar todos os comentários de um post específico
func (s *CommentService) GetCommentsByPostID(postID uuid.UUID) ([]*models.Comment, error) {
	comments, err := s.repo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) GetCommentByID(id uuid.UUID) (*models.Comment, error) {
	comment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (s *CommentService) UpdateComment(commentID uuid.UUID, updatedData *models.Comment) (*models.Comment, error) {
	// Busca o comment existente
	existingComment, err := s.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	// Atualiza os campos com as novas informações, mantendo o que não mudou
	existingComment.Content = updatedData.Content
	existingComment.UpdatedAt = time.Now()

	// Salva a atualização
	if err := s.repo.Update(existingComment); err != nil {
		return nil, err
	}
	return existingComment, nil
}

func (s *CommentService) DeleteComment(id uuid.UUID) error {
	// Pode fazer validações antes de deletar, se necessário
	return s.repo.Delete(id)
}
