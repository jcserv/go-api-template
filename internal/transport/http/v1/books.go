package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jcserv/go-api-template/internal/transport/http/httputil"
	"github.com/jcserv/go-api-template/internal/utils/log"
)

func (a *API) CreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input CreateBook
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Debug(ctx, err.Error())
			httputil.BadRequest(w, fmt.Errorf("unable to decode request body"))
			return
		}

		params, err := input.Parse()
		if err != nil {
			httputil.BadRequest(w, err)
			return
		}

		_, err = a.bookService.CreateBook(ctx, params)
		if err != nil {
			httputil.BadRequest(w, err)
			return
		}

		httputil.NoContent(w)
	}
}
