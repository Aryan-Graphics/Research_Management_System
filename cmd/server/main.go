package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"research-management-system/internal/config"
	"research-management-system/internal/database"
	"research-management-system/internal/handlers"
	"research-management-system/internal/middleware"
)

func main() {
	cfg := config.Load()

	// Initialize Store: Try MySQL first, gracefully fallback to In-Memory
	var store database.Store
	mysqlStore, err := database.NewMySQLStore(cfg)
	if err == nil {
		store = mysqlStore
		log.Println("[Database] Connected successfully to MySQL!")
	} else {
		log.Printf("[Database Warning] Could not connect to MySQL (%v)", err)
		log.Println("[Database] Falling back to high-fidelity In-Memory Store with pre-seeded academic data.")
		log.Println("[Database Info] To use persistent MySQL: Ensure MySQL service is running and configure MYSQL_USER, MYSQL_PASSWORD, MYSQL_DB.")
		store = database.NewMemoryStore()
	}

	// Handlers
	authHandler := handlers.NewAuthHandler(store)
	projectHandler := handlers.NewProjectHandler(store)
	pubHandler := handlers.NewPublicationHandler(store)
	milestoneHandler := handlers.NewMilestoneHandler(store)
	statsHandler := handlers.NewStatsHandler(store)
	researcherHandler := handlers.NewResearcherHandler(store)

	mux := http.NewServeMux()

	// Page routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" {
			http.ServeFile(w, r, filepath.Join("static", "index.html"))
			return
		}
		if path == "/login" || path == "/login.html" {
			http.ServeFile(w, r, filepath.Join("static", "login.html"))
			return
		}
		if path == "/dashboard" || path == "/dashboard.html" {
			http.ServeFile(w, r, filepath.Join("static", "dashboard.html"))
			return
		}

		// Fallback static files
		filePath := filepath.Join("static", filepath.Clean(strings.TrimPrefix(path, "/")))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}

		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	// Static assets handler (/static/...)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Auth API
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/me", middleware.RequireAuth(authHandler.Me))

	// Projects API
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.RequireAuth(projectHandler.HandleProjects)(w, r)
		} else {
			projectHandler.HandleProjects(w, r)
		}
	})
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			middleware.RequireAuth(projectHandler.HandleProjectByID)(w, r)
		} else {
			projectHandler.HandleProjectByID(w, r)
		}
	})

	// Publications API
	mux.HandleFunc("/api/publications", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.RequireAuth(pubHandler.HandlePublications)(w, r)
		} else {
			pubHandler.HandlePublications(w, r)
		}
	})
	mux.HandleFunc("/api/publications/", middleware.RequireAuth(pubHandler.HandlePublicationByID))

	// Milestones API
	mux.HandleFunc("/api/milestones", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.RequireAuth(milestoneHandler.HandleMilestones)(w, r)
		} else {
			milestoneHandler.HandleMilestones(w, r)
		}
	})
	mux.HandleFunc("/api/milestones/", middleware.RequireAuth(milestoneHandler.HandleMilestoneByID))

	// Researchers API
	mux.HandleFunc("/api/researchers", researcherHandler.List)

	// Stats & Activities API
	mux.HandleFunc("/api/stats", statsHandler.GetStats)
	mux.HandleFunc("/api/activities", statsHandler.GetActivities)

	// Health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","mode":"%s"}`, store.GetMode())
	})

	addr := ":" + cfg.Port
	fmt.Println("==================================================================")
	fmt.Printf("🚀 Research Management System Server running at: http://localhost:%s\n", cfg.Port)
	fmt.Printf("   - Front Page:  http://localhost:%s/\n", cfg.Port)
	fmt.Printf("   - Login Page:  http://localhost:%s/login\n", cfg.Port)
	fmt.Printf("   - Dashboard:   http://localhost:%s/dashboard\n", cfg.Port)
	fmt.Printf("   - Data Store:  %s\n", store.GetMode())
	fmt.Println("==================================================================")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
