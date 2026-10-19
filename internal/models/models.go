package models

import "time"

// User represents an academic or administrative user
type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	Role        string    `json:"role"` // 'Admin', 'Principal Investigator', 'Researcher', 'Reviewer'
	Department  string    `json:"department"`
	Institution string    `json:"institution"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserSafe is user representation without password hash
type UserSafe struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Department  string    `json:"department"`
	Institution string    `json:"institution"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *User) ToSafe() UserSafe {
	return UserSafe{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        u.Role,
		Department:  u.Department,
		Institution: u.Institution,
		Avatar:      u.Avatar,
		CreatedAt:   u.CreatedAt,
	}
}

// Project represents a research project
type Project struct {
	ID                 int       `json:"id"`
	Code               string    `json:"code"`
	Title              string    `json:"title"`
	Abstract           string    `json:"abstract"`
	LeadResearcherID   int       `json:"lead_researcher_id"`
	LeadResearcherName string    `json:"lead_researcher_name"`
	Department         string    `json:"department"`
	Status             string    `json:"status"` // 'Planning', 'In Progress', 'Under Review', 'Completed', 'Suspended'
	Budget             float64   `json:"budget"`
	GrantSource        string    `json:"grant_source"`
	StartDate          string    `json:"start_date"` // YYYY-MM-DD
	EndDate            string    `json:"end_date"`   // YYYY-MM-DD
	Tags               string    `json:"tags"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Publication represents a published paper or journal article
type Publication struct {
	ID                  int       `json:"id"`
	ProjectID           *int      `json:"project_id"`
	Title               string    `json:"title"`
	Authors             string    `json:"authors"`
	JournalOrConference string    `json:"journal_or_conference"`
	Year                int       `json:"year"`
	DOI                 string    `json:"doi"`
	Citations           int       `json:"citations"`
	PDFLink             string    `json:"pdf_link"`
	CreatedAt           time.Time `json:"created_at"`
}

// Milestone represents a project deadline or checkpoint
type Milestone struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"` // YYYY-MM-DD
	Status      string    `json:"status"`   // 'Pending', 'In Progress', 'Achieved'
	CreatedAt   time.Time `json:"created_at"`
}

// Activity represents an audit trail log
type Activity struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// DashboardStats holds summarized figures for dashboard cards and charts
type DashboardStats struct {
	TotalProjects      int                `json:"total_projects"`
	ActiveProjects     int                `json:"active_projects"`
	CompletedProjects  int                `json:"completed_projects"`
	TotalPublications  int                `json:"total_publications"`
	TotalCitations     int                `json:"total_citations"`
	TotalBudget        float64            `json:"total_budget"`
	TotalResearchers   int                `json:"total_researchers"`
	PendingMilestones  int                `json:"pending_milestones"`
	ProjectsByStatus   map[string]int     `json:"projects_by_status"`
	BudgetByDepartment map[string]float64 `json:"budget_by_department"`
	PublicationsByYear map[int]int        `json:"publications_by_year"`
	DatabaseMode       string             `json:"database_mode"` // "MySQL" or "In-Memory Fallback"
}

// LoginRequest defines credentials for authentication
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest defines input for account creation
type RegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Department  string `json:"department"`
	Institution string `json:"institution"`
}

// AuthResponse returns logged in user info and session token
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserSafe `json:"user"`
}
