package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcserv/go-api-template/internal/service"
)

const (
	APIV1URLPath = "/api/v1/"
	books        = APIV1URLPath + "books"
	book         = books + "/{id}"
)

type API struct {
	bookService service.IBook
}

func NewAPI(deps *Dependencies) *API {
	return &API{
		bookService: deps.bookService,
	}
}

type Dependencies struct {
	bookService service.IBook
}

func NewDependencies(bookService service.IBook) *Dependencies {
	return &Dependencies{
		bookService: bookService,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	r.HandleFunc(books, a.CreateBook()).Methods(http.MethodPost)
	// r.HandleFunc(book, a.ReadBook()).Methods(http.MethodGet)
	// r.HandleFunc(book, a.ReadBooks()).Methods(http.MethodGet)
	// r.HandleFunc(book, a.UpdateBook()).Methods(http.MethodPut)
	// r.HandleFunc(book, a.DeleteBook()).Methods(http.MethodDelete)
}
