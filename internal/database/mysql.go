package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"research-management-system/internal/config"
	"research-management-system/internal/models"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(cfg *config.Config) (*MySQLStore, error) {
	// First check/create database if needed
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&multiStatements=true",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort)

	tempDB, err := sql.Open("mysql", rootDSN)
	if err == nil {
		tempDB.SetConnMaxLifetime(10 * time.Second)
		if pingErr := tempDB.Ping(); pingErr == nil {
			_, _ = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.MySQLDB))
			tempDB.Close()
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&multiStatements=true",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDB)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	store := &MySQLStore{db: db}
	if err := store.autoMigrateAndSeed(); err != nil {
		log.Printf("[MySQL Warning] Migration/Seed notice: %v", err)
	}

	return store, nil
}

func (s *MySQLStore) GetMode() string {
	return "Live MySQL 8.0"
}

func (s *MySQLStore) autoMigrateAndSeed() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			email VARCHAR(191) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'Researcher',
			department VARCHAR(150) NOT NULL,
			institution VARCHAR(200) NOT NULL DEFAULT 'Global Institute of Advanced Sciences',
			avatar VARCHAR(255) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS projects (
			id INT AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(50) NOT NULL UNIQUE,
			title VARCHAR(255) NOT NULL,
			abstract TEXT NOT NULL,
			lead_researcher_id INT NOT NULL,
			lead_researcher_name VARCHAR(150) NOT NULL,
			department VARCHAR(150) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'In Progress',
			budget DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
			grant_source VARCHAR(200) DEFAULT '',
			start_date VARCHAR(20) NOT NULL,
			end_date VARCHAR(20) NOT NULL,
			tags VARCHAR(255) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (lead_researcher_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS publications (
			id INT AUTO_INCREMENT PRIMARY KEY,
			project_id INT NULL,
			title VARCHAR(255) NOT NULL,
			authors VARCHAR(255) NOT NULL,
			journal_or_conference VARCHAR(255) NOT NULL,
			year INT NOT NULL,
			doi VARCHAR(100) DEFAULT '',
			citations INT DEFAULT 0,
			pdf_link VARCHAR(255) DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS milestones (
			id INT AUTO_INCREMENT PRIMARY KEY,
			project_id INT NOT NULL,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			due_date VARCHAR(20) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'Pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS activities (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_name VARCHAR(150) NOT NULL,
			action VARCHAR(100) NOT NULL,
			details TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}

	// Check if users exist, if not seed default accounts
	var userCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err == nil && userCount == 0 {
		hash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy" // password123
		_, _ = s.db.Exec(`INSERT INTO users (id, name, email, password_hash, role, department, institution, avatar) VALUES
			(1, 'Dr. Aris Thorne', 'admin@research.org', ?, 'Admin', 'Computer Science & AI', 'Global Institute of Advanced Sciences', 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=250&q=80'),
			(2, 'Dr. Elena Rostova', 'elena@research.org', ?, 'Principal Investigator', 'Quantum Physics', 'European Quantum Initiative', 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&w=250&q=80'),
			(3, 'Dr. Marcus Vance', 'marcus@research.org', ?, 'Senior Researcher', 'Biotechnology & Genomics', 'Oxford Bio-Sciences Lab', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=250&q=80'),
			(4, 'Dr. Sarah Lin', 'sarah@research.org', ?, 'Researcher', 'Renewable Energy', 'Pacific Clean Tech Center', 'https://images.unsplash.com/photo-1580489944761-15a19d654956?auto=format&fit=crop&w=250&q=80')`,
			hash, hash, hash, hash)

		_, _ = s.db.Exec(`INSERT INTO projects (id, code, title, abstract, lead_researcher_id, lead_researcher_name, department, status, budget, grant_source, start_date, end_date, tags) VALUES
			(1, 'RES-AI-2025', 'Neuro-Symbolic Reasoning in Autonomous Robotic Surgery', 'Investigating hybrid neural networks and symbolic deductive logic to achieve zero-shot fault tolerance in robotic micro-surgical procedures.', 1, 'Dr. Aris Thorne', 'Computer Science & AI', 'In Progress', 420000.00, 'National Science Foundation (NSF)', '2025-01-15', '2027-06-30', 'Artificial Intelligence, Robotics, Healthcare'),
			(2, 'RES-QP-2024', 'Room-Temperature Topological Superconductors', 'Synthesis and characterization of layered transition-metal dichalcogenides to exhibit macroscopic quantum coherence at ambient pressures.', 2, 'Dr. Elena Rostova', 'Quantum Physics', 'In Progress', 680000.00, 'Department of Energy (DOE)', '2024-03-01', '2026-12-31', 'Quantum, Nanotechnology, Superconductivity'),
			(3, 'RES-BIO-2025', 'CRISPR-Cas14 Epitranscriptomic Modification for Drought Resilience', 'Targeted RNA editing in staple crop genomes to enhance osmotic stress tolerance without permanent genomic disruption.', 3, 'Dr. Marcus Vance', 'Biotechnology & Genomics', 'Under Review', 310000.00, 'Bill & Melinda Gates Agricultural Fund', '2025-04-01', '2027-03-31', 'Genomics, Agriculture, CRISPR'),
			(4, 'RES-EN-2024', 'Solid-State Perovskite Solar Cells with Over 31% Efficiency', 'Developing dual-junction perovskite-silicon tandem cells with atomic layer passivation to surpass theoretical Shockley-Queisser limits.', 4, 'Dr. Sarah Lin', 'Renewable Energy', 'Completed', 250000.00, 'Clean Horizon Grant', '2024-01-10', '2025-09-30', 'Solar, Photovoltaics, Clean Tech')`)

		_, _ = s.db.Exec(`INSERT INTO publications (id, project_id, title, authors, journal_or_conference, year, doi, citations, pdf_link) VALUES
			(1, 1, 'Hybrid Neuro-Symbolic Guidance for Laparoscopic Interventions', 'A. Thorne, K. Patel, S. Zhang', 'Nature Machine Intelligence', 2025, '10.1038/s42256-025-0012', 38, 'https://example.com/papers/neuro-symbolic.pdf'),
			(2, 2, 'Evidence of Topological Surface States in Bi2Se3 Superlattices', 'E. Rostova, H. Weber, J. Dupont', 'Physical Review Letters', 2024, '10.1103/PhysRevLett.132.086401', 94, 'https://example.com/papers/topological.pdf'),
			(3, 3, 'Epitranscriptomic Base Editing Without Off-Target Cas Cleavage', 'M. Vance, L. Gomez', 'Science Advances', 2025, '10.1126/sciadv.adf9812', 17, 'https://example.com/papers/crispr.pdf'),
			(4, 4, 'Atomic Layer Passivation in 31.4% Perovskite-Silicon Tandems', 'S. Lin, M. K. Chen', 'Joule / Cell Press', 2024, '10.1016/j.joule.2024.08.004', 142, 'https://example.com/papers/solar.pdf')`)

		_, _ = s.db.Exec(`INSERT INTO milestones (id, project_id, title, description, due_date, status) VALUES
			(1, 1, 'In-Vitro Phantom Testing & Validation', 'Complete preclinical validation on synthetic organ phantoms with real-time kinesthetic feedback.', '2025-11-30', 'In Progress'),
			(2, 1, 'Algorithmic Safety Proof Verification', 'Formal mathematical proof of invariant safety bounds under edge conditions.', '2025-08-15', 'Achieved'),
			(3, 2, 'Cryogenic Dilution Refrigerator Magnetometry', 'Measure magnetic susceptibility down to 15 mK under 14 Tesla magnetic fields.', '2025-12-20', 'In Progress'),
			(4, 3, 'Greenhouse Controlled Drought Trials', 'Expose genetically modified specimens to 21-day regulated moisture restriction.', '2026-02-28', 'Pending')`)

		_, _ = s.db.Exec(`INSERT INTO activities (id, user_name, action, details) VALUES
			(1, 'Dr. Aris Thorne', 'Milestone Achieved', 'Algorithmic Safety Proof Verification marked as Achieved for project RES-AI-2025'),
			(2, 'Dr. Sarah Lin', 'Project Completed', 'Marked project RES-EN-2024 as Completed with total 31.4% efficiency milestone reached'),
			(3, 'Dr. Marcus Vance', 'Publication Added', 'Published Epitranscriptomic Base Editing in Science Advances'),
			(4, 'Dr. Elena Rostova', 'Budget Allocation', 'Received milestone tranche 2 funding ($220,000) from DOE')`)
	}

	return nil
}

func (s *MySQLStore) GetUserByEmail(email string) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, name, email, password_hash, role, department, institution, avatar, created_at FROM users WHERE email = ?", email)
	var u models.User
	var avatar sql.NullString
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Department, &u.Institution, &avatar, &u.CreatedAt); err != nil {
		return nil, err
	}
	u.Avatar = avatar.String
	return &u, nil
}

func (s *MySQLStore) GetUserByID(id int) (*models.User, error) {
	row := s.db.QueryRow("SELECT id, name, email, password_hash, role, department, institution, avatar, created_at FROM users WHERE id = ?", id)
	var u models.User
	var avatar sql.NullString
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Department, &u.Institution, &avatar, &u.CreatedAt); err != nil {
		return nil, err
	}
	u.Avatar = avatar.String
	return &u, nil
}

func (s *MySQLStore) CreateUser(u *models.User) error {
	res, err := s.db.Exec("INSERT INTO users (name, email, password_hash, role, department, institution, avatar) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.Name, u.Email, u.Password, u.Role, u.Department, u.Institution, u.Avatar)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		u.ID = int(id)
	}
	return nil
}

func (s *MySQLStore) GetUsers() ([]models.UserSafe, error) {
	rows, err := s.db.Query("SELECT id, name, email, role, department, institution, avatar, created_at FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserSafe
	for rows.Next() {
		var u models.UserSafe
		var avatar sql.NullString
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Department, &u.Institution, &avatar, &u.CreatedAt); err != nil {
			return nil, err
		}
		u.Avatar = avatar.String
		users = append(users, u)
	}
	return users, nil
}

func (s *MySQLStore) GetProjects(search, department, status string) ([]models.Project, error) {
	query := "SELECT id, code, title, abstract, lead_researcher_id, lead_researcher_name, department, status, budget, grant_source, start_date, end_date, tags, created_at, updated_at FROM projects WHERE 1=1"
	var args []interface{}

	if department != "" && department != "All Departments" {
		query += " AND department = ?"
		args = append(args, department)
	}
	if status != "" && status != "All Statuses" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if search != "" {
		query += " AND (title LIKE ? OR abstract LIKE ? OR code LIKE ? OR lead_researcher_name LIKE ? OR tags LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}
	query += " ORDER BY id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var tags, grantSource sql.NullString
		if err := rows.Scan(&p.ID, &p.Code, &p.Title, &p.Abstract, &p.LeadResearcherID, &p.LeadResearcherName,
			&p.Department, &p.Status, &p.Budget, &grantSource, &p.StartDate, &p.EndDate, &tags, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Tags = tags.String
		p.GrantSource = grantSource.String
		projects = append(projects, p)
	}
	return projects, nil
}

func (s *MySQLStore) GetProjectByID(id int) (*models.Project, error) {
	row := s.db.QueryRow("SELECT id, code, title, abstract, lead_researcher_id, lead_researcher_name, department, status, budget, grant_source, start_date, end_date, tags, created_at, updated_at FROM projects WHERE id = ?", id)
	var p models.Project
	var tags, grantSource sql.NullString
	if err := row.Scan(&p.ID, &p.Code, &p.Title, &p.Abstract, &p.LeadResearcherID, &p.LeadResearcherName,
		&p.Department, &p.Status, &p.Budget, &grantSource, &p.StartDate, &p.EndDate, &tags, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	p.Tags = tags.String
	p.GrantSource = grantSource.String
	return &p, nil
}

func (s *MySQLStore) CreateProject(p *models.Project) error {
	res, err := s.db.Exec(`INSERT INTO projects (code, title, abstract, lead_researcher_id, lead_researcher_name, department, status, budget, grant_source, start_date, end_date, tags) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Code, p.Title, p.Abstract, p.LeadResearcherID, p.LeadResearcherName, p.Department, p.Status, p.Budget, p.GrantSource, p.StartDate, p.EndDate, p.Tags)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		p.ID = int(id)
	}
	return nil
}

func (s *MySQLStore) UpdateProject(p *models.Project) error {
	_, err := s.db.Exec(`UPDATE projects SET title = ?, abstract = ?, department = ?, status = ?, budget = ?, grant_source = ?, start_date = ?, end_date = ?, tags = ? WHERE id = ?`,
		p.Title, p.Abstract, p.Department, p.Status, p.Budget, p.GrantSource, p.StartDate, p.EndDate, p.Tags, p.ID)
	return err
}

func (s *MySQLStore) DeleteProject(id int) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}

func (s *MySQLStore) GetPublications(projectID *int) ([]models.Publication, error) {
	query := "SELECT id, project_id, title, authors, journal_or_conference, year, doi, citations, pdf_link, created_at FROM publications"
	var args []interface{}
	if projectID != nil {
		query += " WHERE project_id = ?"
		args = append(args, *projectID)
	}
	query += " ORDER BY year DESC, id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pubs []models.Publication
	for rows.Next() {
		var pub models.Publication
		var pID sql.NullInt64
		var doi, pdf sql.NullString
		if err := rows.Scan(&pub.ID, &pID, &pub.Title, &pub.Authors, &pub.JournalOrConference, &pub.Year, &doi, &pub.Citations, &pdf, &pub.CreatedAt); err != nil {
			return nil, err
		}
		if pID.Valid {
			idInt := int(pID.Int64)
			pub.ProjectID = &idInt
		}
		pub.DOI = doi.String
		pub.PDFLink = pdf.String
		pubs = append(pubs, pub)
	}
	return pubs, nil
}

func (s *MySQLStore) CreatePublication(pub *models.Publication) error {
	res, err := s.db.Exec("INSERT INTO publications (project_id, title, authors, journal_or_conference, year, doi, citations, pdf_link) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		pub.ProjectID, pub.Title, pub.Authors, pub.JournalOrConference, pub.Year, pub.DOI, pub.Citations, pub.PDFLink)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		pub.ID = int(id)
	}
	return nil
}

func (s *MySQLStore) DeletePublication(id int) error {
	_, err := s.db.Exec("DELETE FROM publications WHERE id = ?", id)
	return err
}

func (s *MySQLStore) GetMilestones(projectID *int) ([]models.Milestone, error) {
	query := "SELECT id, project_id, title, description, due_date, status, created_at FROM milestones"
	var args []interface{}
	if projectID != nil {
		query += " WHERE project_id = ?"
		args = append(args, *projectID)
	}
	query += " ORDER BY due_date ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []models.Milestone
	for rows.Next() {
		var m models.Milestone
		var desc sql.NullString
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Title, &desc, &m.DueDate, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.Description = desc.String
		milestones = append(milestones, m)
	}
	return milestones, nil
}

func (s *MySQLStore) CreateMilestone(m *models.Milestone) error {
	res, err := s.db.Exec("INSERT INTO milestones (project_id, title, description, due_date, status) VALUES (?, ?, ?, ?, ?)",
		m.ProjectID, m.Title, m.Description, m.DueDate, m.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		m.ID = int(id)
	}
	return nil
}

func (s *MySQLStore) UpdateMilestoneStatus(id int, status string) error {
	_, err := s.db.Exec("UPDATE milestones SET status = ? WHERE id = ?", status, id)
	return err
}

func (s *MySQLStore) GetActivities(limit int) ([]models.Activity, error) {
	rows, err := s.db.Query("SELECT id, user_name, action, details, created_at FROM activities ORDER BY id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acts []models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.UserName, &a.Action, &a.Details, &a.CreatedAt); err != nil {
			return nil, err
		}
		acts = append(acts, a)
	}
	return acts, nil
}

func (s *MySQLStore) AddActivity(userName, action, details string) error {
	_, err := s.db.Exec("INSERT INTO activities (user_name, action, details) VALUES (?, ?, ?)", userName, action, details)
	return err
}

func (s *MySQLStore) GetStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
		ProjectsByStatus:   make(map[string]int),
		BudgetByDepartment: make(map[string]float64),
		PublicationsByYear: make(map[int]int),
		DatabaseMode:       s.GetMode(),
	}

	_ = s.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&stats.TotalProjects)
	_ = s.db.QueryRow("SELECT COUNT(*) FROM projects WHERE status = 'In Progress'").Scan(&stats.ActiveProjects)
	_ = s.db.QueryRow("SELECT COUNT(*) FROM projects WHERE status = 'Completed'").Scan(&stats.CompletedProjects)
	_ = s.db.QueryRow("SELECT COALESCE(SUM(budget), 0) FROM projects").Scan(&stats.TotalBudget)

	_ = s.db.QueryRow("SELECT COUNT(*), COALESCE(SUM(citations), 0) FROM publications").Scan(&stats.TotalPublications, &stats.TotalCitations)
	_ = s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalResearchers)
	_ = s.db.QueryRow("SELECT COUNT(*) FROM milestones WHERE status IN ('Pending', 'In Progress')").Scan(&stats.PendingMilestones)

	// Status breakdown
	rows, err := s.db.Query("SELECT status, COUNT(*) FROM projects GROUP BY status")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var st string
			var count int
			if err := rows.Scan(&st, &count); err == nil {
				stats.ProjectsByStatus[st] = count
			}
		}
	}

	// Budget by department
	rowsDept, err := s.db.Query("SELECT department, COALESCE(SUM(budget), 0) FROM projects GROUP BY department")
	if err == nil {
		defer rowsDept.Close()
		for rowsDept.Next() {
			var dept string
			var budget float64
			if err := rowsDept.Scan(&dept, &budget); err == nil {
				stats.BudgetByDepartment[dept] = budget
			}
		}
	}

	// Pubs by year
	rowsPubs, err := s.db.Query("SELECT year, COUNT(*) FROM publications GROUP BY year")
	if err == nil {
		defer rowsPubs.Close()
		for rowsPubs.Next() {
			var yr, count int
			if err := rowsPubs.Scan(&yr, &count); err == nil {
				stats.PublicationsByYear[yr] = count
			}
		}
	}

	return stats, nil
}
