package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"research-management-system/internal/database"
	"research-management-system/internal/middleware"
	"research-management-system/internal/models"
)

type ProjectHandler struct {
	store database.Store
}

func NewProjectHandler(store database.Store) *ProjectHandler {
	return &ProjectHandler{store: store}
}

func (h *ProjectHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *ProjectHandler) HandleProjectByID(w http.ResponseWriter, r *http.Request) {
	id, err := middleware.GetIDFromURL(r)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r, id)
	case http.MethodPut:
		h.Update(w, r, id)
	case http.MethodDelete:
		h.Delete(w, r, id)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	department := q.Get("department")
	status := q.Get("status")

	projects, err := h.store.GetProjects(search, department, status)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	sendJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
	project, err := h.store.GetProjectByID(id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
		return
	}
	sendJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	user, err := h.store.GetUserByID(claims.UserID)
	if err != nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "User not found"})
		return
	}

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if p.Title == "" || p.Department == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Title and Department are required"})
		return
	}

	if p.Code == "" {
		p.Code = fmt.Sprintf("RES-%s-%d", strings.ToUpper(p.Department[:3]), user.ID*100+int(p.Budget)%99+1)
	}
	if p.Status == "" {
		p.Status = "In Progress"
	}
	if p.StartDate == "" {
		p.StartDate = "2025-01-01"
	}
	if p.EndDate == "" {
		p.EndDate = "2027-12-31"
	}

	p.LeadResearcherID = user.ID
	p.LeadResearcherName = user.Name

	if err := h.store.CreateProject(&p); err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_ = h.store.AddActivity(user.Name, "Project Initiated", fmt.Sprintf("Created project [%s] %s", p.Code, p.Title))

	sendJSON(w, http.StatusCreated, p)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	existing, err := h.store.GetProjectByID(id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
		return
	}

	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	p.ID = id
	if p.Title == "" {
		p.Title = existing.Title
	}
	if p.Department == "" {
		p.Department = existing.Department
	}
	if p.Status == "" {
		p.Status = existing.Status
	}
	if p.StartDate == "" {
		p.StartDate = existing.StartDate
	}
	if p.EndDate == "" {
		p.EndDate = existing.EndDate
	}

	if err := h.store.UpdateProject(&p); err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	user, _ := h.store.GetUserByID(claims.UserID)
	userName := "Staff"
	if user != nil {
		userName = user.Name
	}
	_ = h.store.AddActivity(userName, "Project Updated", fmt.Sprintf("Updated project %s (%s)", existing.Code, p.Title))

	sendJSON(w, http.StatusOK, p)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
	claims := middleware.GetUserClaims(r.Context())
	if claims == nil {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return
	}

	existing, err := h.store.GetProjectByID(id)
	if err != nil {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "Project not found"})
		return
	}

	if err := h.store.DeleteProject(id); err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	user, _ := h.store.GetUserByID(claims.UserID)
	userName := "Staff"
	if user != nil {
		userName = user.Name
	}
	_ = h.store.AddActivity(userName, "Project Archived", fmt.Sprintf("Deleted research project: %s", existing.Title))

	sendJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully", "id": strconv.Itoa(id)})
}
