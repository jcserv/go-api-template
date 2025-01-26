package v1

import (
	"fmt"

	"github.com/jcserv/go-api-template/internal/repository"
)

type CreateBook struct {
	Title    string `json:"title"`
	AuthorID int32  `json:"author_id"`
}

func (c *CreateBook) Parse() (*repository.CreateBookParams, error) {
	if c.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	return &repository.CreateBookParams{
		Title:    c.Title,
		AuthorID: c.AuthorID,
	}, nil
}
