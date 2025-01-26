package service

import (
	"context"

	"github.com/jcserv/go-api-template/internal/repository"
)

type Book struct {
	repo *repository.Queries
}

func NewBookService(repo *repository.Queries) *Book {
	return &Book{
		repo: repo,
	}
}

func (s *Book) CreateBook(ctx context.Context, arg *repository.CreateBookParams) (*repository.Book, error) {
	book, err := s.repo.CreateBook(ctx, *arg)
	return &book, err
}
