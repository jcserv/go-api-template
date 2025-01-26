package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jcserv/go-api-template/internal/repository"
	"github.com/jcserv/go-api-template/internal/test"
	"github.com/jcserv/go-api-template/internal/test/mocks"
	"go.uber.org/mock/gomock"
)

func TestIntegration_CreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	r := mux.NewRouter()
	mockBookService := mocks.NewMockIBook(ctrl)
	a := NewAPI(NewDependencies(mockBookService))
	a.RegisterRoutes(r)

	t.Run("should return 204 when successful", func(t *testing.T) {
		mockBookService.EXPECT().CreateBook(gomock.Any(), gomock.Any()).Return(&repository.Book{}, nil)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/books",
			test.GetBody(map[string]interface{}{
				"title":     "The Martian",
				"author_id": 1,
			}),
		)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204 but got %d", w.Code)
		}
	})

	t.Run("should return 400 when title is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/books",
			test.GetBody(map[string]interface{}{
				"author_id": 1,
			}),
		)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 but got %d", w.Code)
		}
	})
}
