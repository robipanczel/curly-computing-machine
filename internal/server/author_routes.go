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

func (h *Server) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	authorRequest := database.AuthorRequest{}

	err := render.Bind(r, &authorRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorID, err := h.db.CreateAuthor(r.Context(), authorRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID primitive.ObjectID `json:"id"`
	}{
		ID: *authorID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Server) GetAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := primitive.ObjectIDFromHex(chi.URLParam(r, "author_id"))
	if err != nil {
		http.Error(w, "invalid author_id", http.StatusBadRequest)
		return
	}

	author, err := h.db.GetAuthor(r.Context(), authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if author == nil {
		http.Error(w, fmt.Errorf("no author with this ID").Error(), http.StatusNotFound)
		return
	}

	render.Render(w, r, author)
}
