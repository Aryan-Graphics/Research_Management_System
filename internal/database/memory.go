package database

import (
	"errors"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"research-management-system/internal/models"
)

type MemoryStore struct {
	mu           sync.RWMutex
	users        map[int]*models.User
	projects     map[int]*models.Project
	publications map[int]*models.Publication
	milestones   map[int]*models.Milestone
	activities   []models.Activity

	nextUserID        int
	nextProjectID     int
	nextPublicationID int
	nextMilestoneID   int
	nextActivityID    int
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users:        make(map[int]*models.User),
		projects:     make(map[int]*models.Project),
		publications: make(map[int]*models.Publication),
		milestones:   make(map[int]*models.Milestone),
		activities:   make([]models.Activity, 0),

		nextUserID:        1,
		nextProjectID:     1,
		nextPublicationID: 1,
		nextMilestoneID:   1,
		nextActivityID:    1,
	}

	store.seedInitialData()
	return store
}

func (s *MemoryStore) GetMode() string {
	return "In-Memory (Zero-Config Mode - Start MySQL service anytime)"
}

func (s *MemoryStore) seedInitialData() {
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashStr := string(passHash)

	// Users
	u1 := &models.User{
		ID:          s.nextUserID,
		Name:        "Dr. Aris Thorne",
		Email:       "admin@research.org",
		Password:    hashStr,
		Role:        "Admin",
		Department:  "Computer Science & AI",
		Institution: "Global Institute of Advanced Sciences",
		Avatar:      "https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=250&q=80",
		CreatedAt:   time.Now().Add(-120 * 24 * time.Hour),
	}
	s.users[u1.ID] = u1
	s.nextUserID++

	u2 := &models.User{
		ID:          s.nextUserID,
		Name:        "Dr. Elena Rostova",
		Email:       "elena@research.org",
		Password:    hashStr,
		Role:        "Principal Investigator",
		Department:  "Quantum Physics",
		Institution: "European Quantum Initiative",
		Avatar:      "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&w=250&q=80",
		CreatedAt:   time.Now().Add(-100 * 24 * time.Hour),
	}
	s.users[u2.ID] = u2
	s.nextUserID++

	u3 := &models.User{
		ID:          s.nextUserID,
		Name:        "Dr. Marcus Vance",
		Email:       "marcus@research.org",
		Password:    hashStr,
		Role:        "Senior Researcher",
		Department:  "Biotechnology & Genomics",
		Institution: "Oxford Bio-Sciences Lab",
		Avatar:      "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&w=250&q=80",
		CreatedAt:   time.Now().Add(-80 * 24 * time.Hour),
	}
	s.users[u3.ID] = u3
	s.nextUserID++

	u4 := &models.User{
		ID:          s.nextUserID,
		Name:        "Dr. Sarah Lin",
		Email:       "sarah@research.org",
		Password:    hashStr,
		Role:        "Researcher",
		Department:  "Renewable Energy",
		Institution: "Pacific Clean Tech Center",
		Avatar:      "https://images.unsplash.com/photo-1580489944761-15a19d654956?auto=format&fit=crop&w=250&q=80",
		CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
	}
	s.users[u4.ID] = u4
	s.nextUserID++

	// Projects
	p1 := &models.Project{
		ID:                 s.nextProjectID,
		Code:               "RES-AI-2025",
		Title:              "Neuro-Symbolic Reasoning in Autonomous Robotic Surgery",
		Abstract:           "Investigating hybrid neural networks and symbolic deductive logic to achieve zero-shot fault tolerance in robotic micro-surgical procedures.",
		LeadResearcherID:   u1.ID,
		LeadResearcherName: u1.Name,
		Department:         "Computer Science & AI",
		Status:             "In Progress",
		Budget:             420000.00,
		GrantSource:        "National Science Foundation (NSF)",
		StartDate:          "2025-01-15",
		EndDate:            "2027-06-30",
		Tags:               "Artificial Intelligence, Robotics, Healthcare",
		CreatedAt:          time.Now().Add(-90 * 24 * time.Hour),
		UpdatedAt:          time.Now().Add(-2 * 24 * time.Hour),
	}
	s.projects[p1.ID] = p1
	s.nextProjectID++

	p2 := &models.Project{
		ID:                 s.nextProjectID,
		Code:               "RES-QP-2024",
		Title:              "Room-Temperature Topological Superconductors",
		Abstract:           "Synthesis and characterization of layered transition-metal dichalcogenides to exhibit macroscopic quantum coherence at ambient pressures.",
		LeadResearcherID:   u2.ID,
		LeadResearcherName: u2.Name,
		Department:         "Quantum Physics",
		Status:             "In Progress",
		Budget:             680000.00,
		GrantSource:        "Department of Energy (DOE)",
		StartDate:          "2024-03-01",
		EndDate:            "2026-12-31",
		Tags:               "Quantum, Nanotechnology, Superconductivity",
		CreatedAt:          time.Now().Add(-180 * 24 * time.Hour),
		UpdatedAt:          time.Now().Add(-5 * 24 * time.Hour),
	}
	s.projects[p2.ID] = p2
	s.nextProjectID++

	p3 := &models.Project{
		ID:                 s.nextProjectID,
		Code:               "RES-BIO-2025",
		Title:              "CRISPR-Cas14 Epitranscriptomic Modification for Drought Resilience",
		Abstract:           "Targeted RNA editing in staple crop genomes to enhance osmotic stress tolerance without permanent genomic disruption.",
		LeadResearcherID:   u3.ID,
		LeadResearcherName: u3.Name,
		Department:         "Biotechnology & Genomics",
		Status:             "Under Review",
		Budget:             310000.00,
		GrantSource:        "Bill & Melinda Gates Agricultural Fund",
		StartDate:          "2025-04-01",
		EndDate:            "2027-03-31",
		Tags:               "Genomics, Agriculture, CRISPR",
		CreatedAt:          time.Now().Add(-45 * 24 * time.Hour),
		UpdatedAt:          time.Now().Add(-10 * 24 * time.Hour),
	}
	s.projects[p3.ID] = p3
	s.nextProjectID++

	p4 := &models.Project{
		ID:                 s.nextProjectID,
		Code:               "RES-EN-2024",
		Title:              "Solid-State Perovskite Solar Cells with Over 31% Efficiency",
		Abstract:           "Developing dual-junction perovskite-silicon tandem cells with atomic layer passivation to surpass theoretical Shockley-Queisser limits.",
		LeadResearcherID:   u4.ID,
		LeadResearcherName: u4.Name,
		Department:         "Renewable Energy",
		Status:             "Completed",
		Budget:             250000.00,
		GrantSource:        "Clean Horizon Grant",
		StartDate:          "2024-01-10",
		EndDate:            "2025-09-30",
		Tags:               "Solar, Photovoltaics, Clean Tech",
		CreatedAt:          time.Now().Add(-240 * 24 * time.Hour),
		UpdatedAt:          time.Now().Add(-30 * 24 * time.Hour),
	}
	s.projects[p4.ID] = p4
	s.nextProjectID++

	// Publications
	pub1Id := 1
	s.publications[s.nextPublicationID] = &models.Publication{
		ID:                  s.nextPublicationID,
		ProjectID:           &pub1Id,
		Title:               "Hybrid Neuro-Symbolic Guidance for Laparoscopic Interventions",
		Authors:             "A. Thorne, K. Patel, S. Zhang",
		JournalOrConference: "Nature Machine Intelligence",
		Year:                2025,
		DOI:                 "10.1038/s42256-025-0012",
		Citations:           38,
		PDFLink:             "https://example.com/papers/neuro-symbolic.pdf",
		CreatedAt:           time.Now().Add(-40 * 24 * time.Hour),
	}
	s.nextPublicationID++

	pub2Id := 2
	s.publications[s.nextPublicationID] = &models.Publication{
		ID:                  s.nextPublicationID,
		ProjectID:           &pub2Id,
		Title:               "Evidence of Topological Surface States in Bi2Se3 Superlattices",
		Authors:             "E. Rostova, H. Weber, J. Dupont",
		JournalOrConference: "Physical Review Letters",
		Year:                2024,
		DOI:                 "10.1103/PhysRevLett.132.086401",
		Citations:           94,
		PDFLink:             "https://example.com/papers/topological-superlattice.pdf",
		CreatedAt:           time.Now().Add(-95 * 24 * time.Hour),
	}
	s.nextPublicationID++

	pub3Id := 3
	s.publications[s.nextPublicationID] = &models.Publication{
		ID:                  s.nextPublicationID,
		ProjectID:           &pub3Id,
		Title:               "Epitranscriptomic Base Editing Without Off-Target Cas Cleavage",
		Authors:             "M. Vance, L. Gomez",
		JournalOrConference: "Science Advances",
		Year:                2025,
		DOI:                 "10.1126/sciadv.adf9812",
		Citations:           17,
		PDFLink:             "https://example.com/papers/epitranscriptomic.pdf",
		CreatedAt:           time.Now().Add(-20 * 24 * time.Hour),
	}
	s.nextPublicationID++

	pub4Id := 4
	s.publications[s.nextPublicationID] = &models.Publication{
		ID:                  s.nextPublicationID,
		ProjectID:           &pub4Id,
		Title:               "Atomic Layer Passivation in 31.4% Perovskite-Silicon Tandems",
		Authors:             "S. Lin, M. K. Chen",
		JournalOrConference: "Joule / Cell Press",
		Year:                2024,
		DOI:                 "10.1016/j.joule.2024.08.004",
		Citations:           142,
		PDFLink:             "https://example.com/papers/perovskite-tandem.pdf",
		CreatedAt:           time.Now().Add(-110 * 24 * time.Hour),
	}
	s.nextPublicationID++

	// Milestones
	s.milestones[s.nextMilestoneID] = &models.Milestone{
		ID:          s.nextMilestoneID,
		ProjectID:   p1.ID,
		Title:       "In-Vitro Phantom Testing & Validation",
		Description: "Complete preclinical validation on synthetic organ phantoms with real-time kinesthetic feedback.",
		DueDate:     "2025-11-30",
		Status:      "In Progress",
		CreatedAt:   time.Now().Add(-60 * 24 * time.Hour),
	}
	s.nextMilestoneID++

	s.milestones[s.nextMilestoneID] = &models.Milestone{
		ID:          s.nextMilestoneID,
		ProjectID:   p1.ID,
		Title:       "Algorithmic Safety Proof Verification",
		Description: "Formal mathematical proof of invariant safety bounds under edge conditions.",
		DueDate:     "2025-08-15",
		Status:      "Achieved",
		CreatedAt:   time.Now().Add(-90 * 24 * time.Hour),
	}
	s.nextMilestoneID++

	s.milestones[s.nextMilestoneID] = &models.Milestone{
		ID:          s.nextMilestoneID,
		ProjectID:   p2.ID,
		Title:       "Cryogenic Dilution Refrigerator Magnetometry",
		Description: "Measure magnetic susceptibility down to 15 mK under 14 Tesla magnetic fields.",
		DueDate:     "2025-12-20",
		Status:      "In Progress",
		CreatedAt:   time.Now().Add(-40 * 24 * time.Hour),
	}
	s.nextMilestoneID++

	s.milestones[s.nextMilestoneID] = &models.Milestone{
		ID:          s.nextMilestoneID,
		ProjectID:   p3.ID,
		Title:       "Greenhouse Controlled Drought Trials",
		Description: "Expose genetically modified specimens to 21-day regulated moisture restriction.",
		DueDate:     "2026-02-28",
		Status:      "Pending",
		CreatedAt:   time.Now().Add(-15 * 24 * time.Hour),
	}
	s.nextMilestoneID++

	// Activities
	s.activities = append(s.activities, models.Activity{
		ID:        s.nextActivityID,
		UserName:  "Dr. Aris Thorne",
		Action:    "Milestone Achieved",
		Details:   "Algorithmic Safety Proof Verification marked as Achieved for project RES-AI-2025",
		CreatedAt: time.Now().Add(-2 * 24 * time.Hour),
	})
	s.nextActivityID++

	s.activities = append(s.activities, models.Activity{
		ID:        s.nextActivityID,
		UserName:  "Dr. Sarah Lin",
		Action:    "Project Completed",
		Details:   "Marked project RES-EN-2024 as Completed with total 31.4% efficiency milestone reached",
		CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
	})
	s.nextActivityID++

	s.activities = append(s.activities, models.Activity{
		ID:        s.nextActivityID,
		UserName:  "Dr. Marcus Vance",
		Action:    "Publication Added",
		Details:   "Published 'Epitranscriptomic Base Editing' in Science Advances",
		CreatedAt: time.Now().Add(-4 * time.Hour),
	})
	s.nextActivityID++

	s.activities = append(s.activities, models.Activity{
		ID:        s.nextActivityID,
		UserName:  "Dr. Elena Rostova",
		Action:    "Budget Allocation",
		Details:   "Received milestone tranche 2 funding ($220,000) from DOE",
		CreatedAt: time.Now().Add(-30 * time.Minute),
	})
	s.nextActivityID++
}

func (s *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if strings.EqualFold(u.Email, email) {
			userCopy := *u
			return &userCopy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *MemoryStore) GetUserByID(id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, exists := s.users[id]; exists {
		userCopy := *u
		return &userCopy, nil
	}
	return nil, errors.New("user not found")
}

func (s *MemoryStore) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if strings.EqualFold(u.Email, user.Email) {
			return errors.New("email already registered")
		}
	}

	user.ID = s.nextUserID
	s.nextUserID++
	user.CreatedAt = time.Now()
	if user.Avatar == "" {
		user.Avatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=250&q=80"
	}

	s.users[user.ID] = user
	return nil
}

func (s *MemoryStore) GetUsers() ([]models.UserSafe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]models.UserSafe, 0, len(s.users))
	for _, u := range s.users {
		res = append(res, u.ToSafe())
	}
	return res, nil
}

func (s *MemoryStore) GetProjects(search, department, status string) ([]models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.Project
	searchLower := strings.ToLower(search)

	for _, p := range s.projects {
		if department != "" && department != "All Departments" && !strings.EqualFold(p.Department, department) {
			continue
		}
		if status != "" && status != "All Statuses" && !strings.EqualFold(p.Status, status) {
			continue
		}
		if search != "" {
			matches := strings.Contains(strings.ToLower(p.Title), searchLower) ||
				strings.Contains(strings.ToLower(p.Abstract), searchLower) ||
				strings.Contains(strings.ToLower(p.Code), searchLower) ||
				strings.Contains(strings.ToLower(p.LeadResearcherName), searchLower) ||
				strings.Contains(strings.ToLower(p.Tags), searchLower)
			if !matches {
				continue
			}
		}
		results = append(results, *p)
	}
	return results, nil
}

func (s *MemoryStore) GetProjectByID(id int) (*models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if p, exists := s.projects[id]; exists {
		pCopy := *p
		return &pCopy, nil
	}
	return nil, errors.New("project not found")
}

func (s *MemoryStore) CreateProject(p *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextProjectID
	s.nextProjectID++
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	s.projects[p.ID] = p
	return nil
}

func (s *MemoryStore) UpdateProject(p *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.projects[p.ID]
	if !exists {
		return errors.New("project not found")
	}

	existing.Title = p.Title
	existing.Abstract = p.Abstract
	existing.Department = p.Department
	existing.Status = p.Status
	existing.Budget = p.Budget
	existing.GrantSource = p.GrantSource
	existing.StartDate = p.StartDate
	existing.EndDate = p.EndDate
	existing.Tags = p.Tags
	existing.UpdatedAt = time.Now()

	return nil
}

func (s *MemoryStore) DeleteProject(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[id]; !exists {
		return errors.New("project not found")
	}
	delete(s.projects, id)

	// Remove project reference from publications
	for _, pub := range s.publications {
		if pub.ProjectID != nil && *pub.ProjectID == id {
			pub.ProjectID = nil
		}
	}
	// Delete associated milestones
	for mid, m := range s.milestones {
		if m.ProjectID == id {
			delete(s.milestones, mid)
		}
	}
	return nil
}

func (s *MemoryStore) GetPublications(projectID *int) ([]models.Publication, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.Publication
	for _, pub := range s.publications {
		if projectID != nil {
			if pub.ProjectID == nil || *pub.ProjectID != *projectID {
				continue
			}
		}
		results = append(results, *pub)
	}
	return results, nil
}

func (s *MemoryStore) CreatePublication(pub *models.Publication) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pub.ID = s.nextPublicationID
	s.nextPublicationID++
	pub.CreatedAt = time.Now()

	s.publications[pub.ID] = pub
	return nil
}

func (s *MemoryStore) DeletePublication(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.publications[id]; !exists {
		return errors.New("publication not found")
	}
	delete(s.publications, id)
	return nil
}

func (s *MemoryStore) GetMilestones(projectID *int) ([]models.Milestone, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []models.Milestone
	for _, m := range s.milestones {
		if projectID != nil && m.ProjectID != *projectID {
			continue
		}
		results = append(results, *m)
	}
	return results, nil
}

func (s *MemoryStore) CreateMilestone(m *models.Milestone) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m.ID = s.nextMilestoneID
	s.nextMilestoneID++
	m.CreatedAt = time.Now()

	s.milestones[m.ID] = m
	return nil
}

func (s *MemoryStore) UpdateMilestoneStatus(id int, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, exists := s.milestones[id]
	if !exists {
		return errors.New("milestone not found")
	}
	m.Status = status
	return nil
}

func (s *MemoryStore) GetActivities(limit int) ([]models.Activity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.activities)
	start := total - limit
	if start < 0 {
		start = 0
	}

	results := make([]models.Activity, 0, total-start)
	for i := total - 1; i >= start; i-- {
		results = append(results, s.activities[i])
	}
	return results, nil
}

func (s *MemoryStore) AddActivity(userName, action, details string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	act := models.Activity{
		ID:        s.nextActivityID,
		UserName:  userName,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}
	s.nextActivityID++
	s.activities = append(s.activities, act)
	return nil
}

func (s *MemoryStore) GetStats() (*models.DashboardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := &models.DashboardStats{
		TotalProjects:      len(s.projects),
		TotalPublications:  len(s.publications),
		TotalResearchers:   len(s.users),
		ProjectsByStatus:   make(map[string]int),
		BudgetByDepartment: make(map[string]float64),
		PublicationsByYear: make(map[int]int),
		DatabaseMode:       s.GetMode(),
	}

	for _, p := range s.projects {
		if p.Status == "In Progress" {
			stats.ActiveProjects++
		} else if p.Status == "Completed" {
			stats.CompletedProjects++
		}
		stats.TotalBudget += p.Budget
		stats.ProjectsByStatus[p.Status]++
		stats.BudgetByDepartment[p.Department] += p.Budget
	}

	for _, pub := range s.publications {
		stats.TotalCitations += pub.Citations
		stats.PublicationsByYear[pub.Year]++
	}

	for _, m := range s.milestones {
		if m.Status == "Pending" || m.Status == "In Progress" {
			stats.PendingMilestones++
		}
	}

	return stats, nil
}
