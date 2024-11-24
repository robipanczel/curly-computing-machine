package server

import (
	"curly-computing-machine/internal/database"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Server) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.db.ListBooks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (h *Server) AddBook(w http.ResponseWriter, r *http.Request) {
	bookRequest := database.BookRequest{}

	err := render.Bind(r, &bookRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookID, err := h.db.AddBook(r.Context(), bookRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID primitive.ObjectID `json:"id"`
	}{
		ID: *bookID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Server) BorrowBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "book_id"))
	if err != nil {
		http.Error(w, "invalid book_id", http.StatusBadRequest)
		return
	}

	borrowerID, err := primitive.ObjectIDFromHex(r.URL.Query().Get("borrower_id"))
	if err != nil {
		http.Error(w, "invalid borrower_id", http.StatusBadRequest)
		return
	}

	err = h.db.BorrowBook(r.Context(), bookID, borrowerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
