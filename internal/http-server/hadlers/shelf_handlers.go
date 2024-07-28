package handlers

import (
	"encoding/json"
	"myproject/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ShelfHandler struct {
	shelfRepo repository.ShelfRepository
}

func NewShelfHandler(shelfRepo repository.ShelfRepository) *ShelfHandler {
	return &ShelfHandler{shelfRepo: shelfRepo}
}

func (h *ShelfHandler) CreateShelf(w http.ResponseWriter, r *http.Request){
	shelf, err := h.shelfRepo.AddShelf()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelf)
}

func (h *ShelfHandler) GetShelf(w http.ResponseWriter, r *http.Request){
	shelfId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid shelf ID", http.StatusBadRequest)
		return
	}

	shelf, err := h.shelfRepo.GetShelf(shelfId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelf)
}

func (h *ShelfHandler) DeleteShelf(w http.ResponseWriter, r *http.Request){
	shelfId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid shelf ID", http.StatusBadRequest)
		return
	}
	err = h.shelfRepo.DeleteShelf(shelfId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
