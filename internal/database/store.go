package database

import (
	"research-management-system/internal/models"
)

// Store defines the interface for data access (both MySQL and Memory fallback)
type Store interface {
	// Users
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	GetUsers() ([]models.UserSafe, error)

	// Projects
	GetProjects(search, department, status string) ([]models.Project, error)
	GetProjectByID(id int) (*models.Project, error)
	CreateProject(project *models.Project) error
	UpdateProject(project *models.Project) error
	DeleteProject(id int) error

	// Publications
	GetPublications(projectID *int) ([]models.Publication, error)
	CreatePublication(pub *models.Publication) error
	DeletePublication(id int) error

	// Milestones
	GetMilestones(projectID *int) ([]models.Milestone, error)
	CreateMilestone(m *models.Milestone) error
	UpdateMilestoneStatus(id int, status string) error

	// Activities
	GetActivities(limit int) ([]models.Activity, error)
	AddActivity(userName, action, details string) error

	// Stats
	GetStats() (*models.DashboardStats, error)

	// Database Mode Indicator ("MySQL" or "In-Memory")
	GetMode() string
}
