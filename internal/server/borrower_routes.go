package server

import (
	"curly-computing-machine/internal/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Server) CreateBorrower(w http.ResponseWriter, r *http.Request) {
	borrowerRequest := database.BorrowerRequest{}

	err := render.Bind(r, &borrowerRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	borrowerID, err := h.db.CreateBorrower(r.Context(), borrowerRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID primitive.ObjectID `json:"id"`
	}{
		ID: *borrowerID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Server) GetBorrower(w http.ResponseWriter, r *http.Request) {
	borrowerID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "borrower_id"))
	if err != nil {
		http.Error(w, "invalid borrower_id", http.StatusBadRequest)
		return
	}

	borrower, err := h.db.GetBorrower(r.Context(), borrowerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if borrower == nil {
		http.Error(w, fmt.Errorf("no borrower with this ID").Error(), http.StatusNotFound)
		return
	}

	render.Render(w, r, borrower)
}

func (h *Server) BorrowedBooks(w http.ResponseWriter, r *http.Request) {
	borrowerID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "borrower_id"))
	if err != nil {
		http.Error(w, "invalid borrower_id", http.StatusBadRequest)
		return
	}

	books, err := h.db.BorrowedBooks(r.Context(), borrowerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}
