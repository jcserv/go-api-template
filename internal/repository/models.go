// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

type Author struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID       int32  `json:"id"`
	Title    string `json:"title"`
	AuthorID int32  `json:"author_id"`
}
