package service

import (
	"context"

	"github.com/jcserv/go-api-template/internal/repository"
)

type IBook interface {
	CreateBook(ctx context.Context, arg *repository.CreateBookParams) (*repository.Book, error)
	// ReadBook(ctx context.Context, id int32) (*repository.Book, error)
	// ReadBooks(ctx context.Context) ([]*repository.Book, error)
	// UpdateBook(ctx context.Context, arg *repository.UpdateBookParams) error
	// DeleteBook(ctx context.Context, id int32) error
}
