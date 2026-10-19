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

type MilestoneHandler struct {
	store database.Store
}

func NewMilestoneHandler(store database.Store) *MilestoneHandler {
	return &MilestoneHandler{store: store}
}

func (h *MilestoneHandler) HandleMilestones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var projectID *int
		if pIDStr := r.URL.Query().Get("project_id"); pIDStr != "" {
			if id, err := strconv.Atoi(pIDStr); err == nil {
				projectID = &id
			}
		}
		milestones, err := h.store.GetMilestones(projectID)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		sendJSON(w, http.StatusOK, milestones)

	case http.MethodPost:
		claims := middleware.GetUserClaims(r.Context())
		if claims == nil {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		var m models.Milestone
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		if m.Title == "" || m.ProjectID == 0 {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Title and ProjectID are required"})
			return
		}

		if m.Status == "" {
			m.Status = "Pending"
		}
		if m.DueDate == "" {
			m.DueDate = "2025-12-31"
		}

		if err := h.store.CreateMilestone(&m); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		user, _ := h.store.GetUserByID(claims.UserID)
		name := "Researcher"
		if user != nil {
			name = user.Name
		}
		_ = h.store.AddActivity(name, "Milestone Created", fmt.Sprintf("Added milestone '%s'", m.Title))

		sendJSON(w, http.StatusCreated, m)

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *MilestoneHandler) HandleMilestoneByID(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIDFromURL(r)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid milestone ID"})
		return
	}

	if r.Method == http.MethodPut || r.Method == http.MethodPatch {
		claims := middleware.GetUserClaims(r.Context())
		if claims == nil {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			return
		}

		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		if err := h.store.UpdateMilestoneStatus(id, req.Status); err != nil {
			sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		user, _ := h.store.GetUserByID(claims.UserID)
		name := "Researcher"
		if user != nil {
			name = user.Name
		}
		_ = h.store.AddActivity(name, "Milestone Updated", fmt.Sprintf("Updated milestone #%d status to %s", id, req.Status))

		sendJSON(w, http.StatusOK, map[string]string{"message": "Milestone status updated", "status": req.Status})
		return
	}

	http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
}
