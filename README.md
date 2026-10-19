# 🔬 ScholarSync - Research Management System (RMS)

A modern, responsive, and full-featured **Research Management System** built with **HTML5, CSS3, JavaScript, Golang, and MySQL 8.0**.

Designed for universities, laboratories, research institutes, and academic consortiums to manage the complete lifecycle of scientific inquiry, grant funding, deliverables, and publications.

---

## 🌟 Key Features

### 1. 🌐 Front Page (Public Portal)
- **Live Institutional Metrics**: Dynamic counters for active projects, peer citations, and total grant funding.
- **Interactive Project Explorer**: Real-time search across titles, abstracts, and keywords, with multi-faculty filtering (AI, Quantum, Biotech, Clean Energy).
- **Publications Archive**: Showcase peer-reviewed papers with direct DOI links and citation counts.
- **Academic Highlights**: Institutional faculties, capabilities, and open science initiatives.

### 2. 🔐 Authentication & Access Control
- **Dual Flow**: Sign In and Registration for academic scholars and administrators.
- **One-Click Demo Profiles**: Instant 1-click test login for:
  - 👑 **Dr. Aris Thorne (Admin / PI)**: `admin@research.org` / `password123`
  - ⚛️ **Dr. Elena Rostova (Principal Investigator)**: `elena@research.org` / `password123`
  - 🧬 **Dr. Marcus Vance (Senior Researcher)**: `marcus@research.org` / `password123`
- **Security**: Cryptographic password hashing (`bcrypt`) and signed token authorization.

### 3. 📊 Research Dashboard
- **Executive Analytics**: KPI metric cards and interactive canvas charts (budget allocation by faculty and project status distribution).
- **Project Governance (Full CRUD)**:
  - Initiate new research projects with budgets, timelines, and abstracts.
  - Edit project milestones and sponsor grant metadata.
  - Comprehensive Project Dossier modal showing linked publications and deliverables.
- **Publications Repository**: Add peer-reviewed manuscripts, DOI identifiers, and track citations.
- **Milestones & Deliverables Tracker**: Interactive timeline checklist with 1-click status transitions (`Pending` ↔ `In Progress` ↔ `Achieved`).
- **Faculty Directory**: Academic roster across participating departments.
- **Institutional Audit Trail**: Activity log tracking milestone updates, project completions, and grant allocations.

---

## 🛠️ Architecture & Tech Stack

| Layer | Technology | Details |
|---|---|---|
| **Frontend** | HTML5, CSS3, Modern JavaScript (ES6+) | Vanilla JS with responsive CSS variables, glassmorphism, Canvas charts |
| **Backend** | Golang (Go 1.27) | `net/http`, modular RESTful architecture, custom token middleware |
| **Database** | MySQL 8.0 / In-Memory Fallback | Fully relational MySQL schema (`sql/schema.sql`) with automatic in-memory fallback |

---

## 🚀 Quick Start Guide

### 1. Run the Server
In PowerShell or Terminal:

```powershell
# Run directly with Go:
go run cmd/server/main.go

# Or execute the compiled binary:
.\server.exe
```

The server will automatically start on `http://localhost:8080`.

### 2. Open in Your Browser
- **Front Page**: [http://localhost:8080/](http://localhost:8080/)
- **Sign In / Register**: [http://localhost:8080/login](http://localhost:8080/login)
- **Research Dashboard**: [http://localhost:8080/dashboard](http://localhost:8080/dashboard)

---

## 🗄️ MySQL Database Setup

The application is configured to connect to MySQL on `localhost:3306` with database `research_db`.

### Automatic Zero-Config Fallback:
If MySQL is not currently running or configured, the system **automatically falls back to a high-fidelity In-Memory Store** seeded with real research data. This allows you to explore the entire UI, create, edit, and delete projects without needing any setup!

### To connect to persistent MySQL:
1. Ensure your MySQL 8.0 service is running:
   ```powershell
   Start-Service MySQL80
   ```
2. Run the SQL schema to create tables and seed data:
   ```powershell
   mysql -u root -p < sql/schema.sql
   ```
3. Set your environment variables (or modify `.env.example`):
   ```powershell
   $env:MYSQL_USER="root"
   $env:MYSQL_PASSWORD="your_password"
   $env:MYSQL_DB="research_db"
   .\server.exe
   ```

The dashboard displays a live badge in the top bar: **"Live MySQL 8.0"** or **"In-Memory Mode"**.

---

## 📁 Project Structure

```text
Research Management System/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server, routing, and graceful DB fallback
├── internal/
│   ├── config/              # Environment configuration loader
│   ├── database/
│   │   ├── store.go         # Common repository interface
│   │   ├── mysql.go         # MySQL 8.0 driver, schema migration, and queries
│   │   └── memory.go        # Thread-safe in-memory store with seeded research data
│   ├── handlers/            # REST API endpoints (Auth, Projects, Pubs, Milestones, Stats)
│   ├── middleware/          # Token authentication & route guard
│   └── models/              # Academic entities (User, Project, Publication, Milestone, Activity)
├── sql/
│   └── schema.sql           # Complete MySQL DDL & initial seed records
├── static/
│   ├── index.html           # Public front / landing page
│   ├── login.html           # Authentication portal (Sign In & Sign Up)
│   ├── dashboard.html       # Full Research Portal Dashboard
│   ├── css/
│   │   ├── style.css        # Core typography, buttons, badges, modals, toasts
│   │   ├── landing.css      # Front page layout & hero
│   │   ├── auth.css         # Sign in / Register glassmorphism cards
│   │   └── dashboard.css    # Sidebar, charts, tables, and timeline styles
│   └── js/
│       ├── api.js           # API client, token management, toasts
│       ├── landing.js       # Live search, filters, and metrics
│       ├── auth.js          # Authentication & 1-click demo profiles
│       └── dashboard.js     # Project CRUD, canvas charts, milestones, publications
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── server.exe               # Compiled executable binary
└── README.md
```
