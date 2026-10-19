-- Research Management System (RMS)
-- MySQL Schema & Seed Data

CREATE DATABASE IF NOT EXISTS `research_db` 
  CHARACTER SET utf8mb4 
  COLLATE utf8mb4_unicode_ci;

USE `research_db`;

-- 1. Users Table
CREATE TABLE IF NOT EXISTS `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(150) NOT NULL,
    `email` VARCHAR(191) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL,
    `role` VARCHAR(50) NOT NULL DEFAULT 'Researcher', -- 'Admin', 'Principal Investigator', 'Researcher', 'Reviewer'
    `department` VARCHAR(150) NOT NULL,
    `institution` VARCHAR(200) NOT NULL DEFAULT 'Global Institute of Advanced Sciences',
    `avatar` VARCHAR(255) DEFAULT '',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Research Projects Table
CREATE TABLE IF NOT EXISTS `projects` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `code` VARCHAR(50) NOT NULL UNIQUE,
    `title` VARCHAR(255) NOT NULL,
    `abstract` TEXT NOT NULL,
    `lead_researcher_id` INT NOT NULL,
    `lead_researcher_name` VARCHAR(150) NOT NULL,
    `department` VARCHAR(150) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'In Progress', -- 'Planning', 'In Progress', 'Under Review', 'Completed', 'Suspended'
    `budget` DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    `grant_source` VARCHAR(200) DEFAULT '',
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `tags` VARCHAR(255) DEFAULT '',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`lead_researcher_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Publications Table
CREATE TABLE IF NOT EXISTS `publications` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT NULL,
    `title` VARCHAR(255) NOT NULL,
    `authors` VARCHAR(255) NOT NULL,
    `journal_or_conference` VARCHAR(255) NOT NULL,
    `year` INT NOT NULL,
    `doi` VARCHAR(100) DEFAULT '',
    `citations` INT DEFAULT 0,
    `pdf_link` VARCHAR(255) DEFAULT '',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. Milestones Table
CREATE TABLE IF NOT EXISTS `milestones` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `project_id` INT NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `description` TEXT,
    `due_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'Pending', -- 'Pending', 'In Progress', 'Achieved'
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. Activity Log Table
CREATE TABLE IF NOT EXISTS `activities` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `user_name` VARCHAR(150) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `details` TEXT NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ========================================================
-- Seed Initial Data (Default Passwords: "password123")
-- bcrypt hash for "password123":
-- $2a$10$7vB.wO7lX7.v3uFv1g8OquHk8c9p1S6k4Y2H0k9e1Z3e5t7j8y9w2
-- ========================================================

INSERT INTO `users` (`id`, `name`, `email`, `password_hash`, `role`, `department`, `institution`) VALUES
(1, 'Dr. Aris Thorne', 'admin@research.org', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Admin', 'Computer Science & AI', 'Global Institute of Advanced Sciences'),
(2, 'Dr. Elena Rostova', 'elena@research.org', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Principal Investigator', 'Quantum Physics', 'European Quantum Initiative'),
(3, 'Dr. Marcus Vance', 'marcus@research.org', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Senior Researcher', 'Biotechnology & Genomics', 'Oxford Bio-Sciences Lab'),
(4, 'Dr. Sarah Lin', 'sarah@research.org', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Researcher', 'Renewable Energy', 'Pacific Clean Tech Center')
ON DUPLICATE KEY UPDATE `email`=`email`;

INSERT INTO `projects` (`id`, `code`, `title`, `abstract`, `lead_researcher_id`, `lead_researcher_name`, `department`, `status`, `budget`, `grant_source`, `start_date`, `end_date`, `tags`) VALUES
(1, 'RES-AI-2025', 'Neuro-Symbolic Reasoning in Autonomous Robotic Surgery', 'Investigating hybrid neural networks and symbolic deductive logic to achieve zero-shot fault tolerance in robotic micro-surgical procedures.', 1, 'Dr. Aris Thorne', 'Computer Science & AI', 'In Progress', 420000.00, 'National Science Foundation (NSF)', '2025-01-15', '2027-06-30', 'Artificial Intelligence, Robotics, Healthcare'),
(2, 'RES-QP-2024', 'Room-Temperature Topological Superconductors', 'Synthesis and characterization of layered transition-metal dichalcogenides to exhibit macroscopic quantum coherence at ambient pressures.', 2, 'Dr. Elena Rostova', 'Quantum Physics', 'In Progress', 680000.00, 'Department of Energy (DOE)', '2024-03-01', '2026-12-31', 'Quantum, Nanotechnology, Superconductivity'),
(3, 'RES-BIO-2025', 'CRISPR-Cas14 Epitranscriptomic Modification for Drought Resilience', 'Targeted RNA editing in staple crop genomes to enhance osmotic stress tolerance without permanent genomic disruption.', 3, 'Dr. Marcus Vance', 'Biotechnology & Genomics', 'Under Review', 310000.00, 'Bill & Melinda Gates Agricultural Fund', '2025-04-01', '2027-03-31', 'Genomics, Agriculture, CRISPR'),
(4, 'RES-EN-2024', 'Solid-State Perovskite Solar Cells with Over 31% Efficiency', 'Developing dual-junction perovskite-silicon tandem cells with atomic layer passivation to surpass theoretical Shockley-Queisser limits.', 4, 'Dr. Sarah Lin', 'Renewable Energy', 'Completed', 250000.00, 'Clean Horizon Grant', '2024-01-10', '2025-09-30', 'Solar, Photovoltaics, Clean Tech')
ON DUPLICATE KEY UPDATE `code`=`code`;

INSERT INTO `publications` (`id`, `project_id`, `title`, `authors`, `journal_or_conference`, `year`, `doi`, `citations`, `pdf_link`) VALUES
(1, 1, 'Hybrid Neuro-Symbolic Guidance for Laparoscopic Interventions', 'A. Thorne, K. Patel, S. Zhang', 'Nature Machine Intelligence', 2025, '10.1038/s42256-025-0012', 38, '#'),
(2, 2, 'Evidence of Topological Surface States in Bi2Se3 Superlattices', 'E. Rostova, H. Weber, J. Dupont', 'Physical Review Letters', 2024, '10.1103/PhysRevLett.132.086401', 94, '#'),
(3, 3, 'Epitranscriptomic Base Editing Without Off-Target Cas Cleavage', 'M. Vance, L. Gomez', 'Science Advances', 2025, '10.1126/sciadv.adf9812', 17, '#'),
(4, 4, 'Atomic Layer Passivation in 31.4% Perovskite-Silicon Tandems', 'S. Lin, M. K. Chen', 'Joule / Cell Press', 2024, '10.1016/j.joule.2024.08.004', 142, '#')
ON DUPLICATE KEY UPDATE `title`=`title`;

INSERT INTO `milestones` (`id`, `project_id`, `title`, `description`, `due_date`, `status`) VALUES
(1, 1, 'In-Vitro Phantom Testing & Validation', 'Complete preclinical validation on synthetic organ phantoms with real-time kinesthetic feedback.', '2025-11-30', 'In Progress'),
(2, 1, 'Algorithmic Safety Proof Verification', 'Formal mathematical proof of invariant safety bounds under edge conditions.', '2025-08-15', 'Achieved'),
(3, 2, 'Cryogenic Dilution Refrigerator Magnetometry', 'Measure magnetic susceptibility down to 15 mK under 14 Tesla magnetic fields.', '2025-12-20', 'In Progress'),
(4, 3, 'Greenhouse Controlled Drought Trials', 'Expose genetically modified specimens to 21-day regulated moisture restriction.', '2026-02-28', 'Pending'),
(5, 4, '1000-Hour Continuous Damp-Heat Endurance Test', 'Complete IEC 61215 accelerated environmental degradation certification testing.', '2025-05-15', 'Achieved')
ON DUPLICATE KEY UPDATE `title`=`title`;

INSERT INTO `activities` (`id`, `user_name`, `action`, `details`, `created_at`) VALUES
(1, 'Dr. Aris Thorne', 'Milestone Achieved', 'Algorithmic Safety Proof Verification marked as Achieved for project RES-AI-2025', NOW() - INTERVAL 2 DAY),
(2, 'Dr. Sarah Lin', 'Project Completed', 'Marked project RES-EN-2024 as Completed with total 31.4% efficiency milestone reached', NOW() - INTERVAL 1 DAY),
(3, 'Dr. Marcus Vance', 'Publication Added', 'Published "Epitranscriptomic Base Editing" in Science Advances', NOW() - INTERVAL 4 HOUR),
(4, 'Dr. Elena Rostova', 'Budget Allocation', 'Received milestone tranche 2 funding ($220,000) from DOE', NOW() - INTERVAL 30 MINUTE)
ON DUPLICATE KEY UPDATE `action`=`action`;
