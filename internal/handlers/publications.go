package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"research-management-system/internal/database"
	"research-management-system/internal/middleware"
	"research-management-system/internal/models"
)

type PublicationHandler struct {
	store database.Store
}

func NewPublicationHandler(store database.Store) *PublicationHandler {
	return &PublicationHandler{store: store}
}

func (h *PublicationHandler) HandlePublications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var projectID *int
		if pIDStr := r.URL.Query().Get("project_id"); pIDStr != "" {
			if id, err := strconv.Atoi(pIDStr); err == nil {
				projectID = &id
			}
		}
		pubs, err := h.store.GetPublications(projectID)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		sendJSON(w, http.StatusOK, pubs)

	case http.MethodPost:
		claims := middleware.GetUserClaims(r.Context())
		if claims == nil {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		var pub models.Publication
		if err := json.NewDecoder(r.Body).Decode(&pub); err != nil {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		if pub.Title == "" || pub.Authors == "" || pub.JournalOrConference == "" {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Title, Authors, and Venue are required"})
			return
		}

		if pub.Year == 0 {
			pub.Year = 2025
		}

		if err := h.store.CreatePublication(&pub); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		user, _ := h.store.GetUserByID(claims.UserID)
		name := "Researcher"
		if user != nil {
			name = user.Name
		}
		_ = h.store.AddActivity(name, "Publication Added", fmt.Sprintf("Registered paper: %s in %s", pub.Title, pub.JournalOrConference))

		sendJSON(w, http.StatusCreated, pub)

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *PublicationHandler) HandlePublicationByID(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIDFromURL(r)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid publication ID"})
		return
	}

	if r.Method == http.MethodDelete {
		claims := middleware.GetUserClaims(r.Context())
		if claims == nil {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		if err := h.store.DeletePublication(id); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		sendJSON(w, http.StatusOK, map[string]string{"message": "Publication removed"})
		return
	}

	http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
}
