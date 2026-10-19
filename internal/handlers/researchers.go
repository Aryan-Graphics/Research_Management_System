package handlers

import (
	"net/http"

	"research-management-system/internal/database"
)

type ResearcherHandler struct {
	store database.Store
}

func NewResearcherHandler(store database.Store) *ResearcherHandler {
	return &ResearcherHandler{store: store}
}

func (h *ResearcherHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	users, err := h.store.GetUsers()
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	sendJSON(w, http.StatusOK, users)
}
