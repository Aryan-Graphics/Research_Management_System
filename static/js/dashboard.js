/**
 * ScholarSync - Complete Dashboard Controller
 */

let currentUser = null;
let allProjects = [];

document.addEventListener('DOMContentLoaded', async () => {
  currentUser = API.getUser();
  if (!currentUser) {
    window.location.href = '/login';
    return;
  }

  renderUserProfile();
  setupNavigation();
  setupModals();
  setupProjectHandlers();
  setupPublicationHandlers();
  setupMilestoneHandlers();
  
  // Initial load
  await loadDashboardOverview();
});

/* ==========================================================================
   User Profile & Auth State
   ========================================================================== */
function renderUserProfile() {
  const avatarEl = document.getElementById('user-avatar');
  const nameEl = document.getElementById('user-name');
  const roleEl = document.getElementById('user-role');
  const topbarUserEl = document.getElementById('topbar-user-name');

  if (avatarEl && currentUser.avatar) avatarEl.src = currentUser.avatar;
  if (nameEl) nameEl.textContent = currentUser.name;
  if (roleEl) roleEl.textContent = currentUser.role;
  if (topbarUserEl) topbarUserEl.textContent = currentUser.name;

  const logoutBtn = document.getElementById('logout-btn');
  if (logoutBtn) {
    logoutBtn.addEventListener('click', async () => {
      try {
        await API.post('/api/auth/logout', {});
      } catch (e) {}
      API.setToken(null);
      API.setUser(null);
      window.location.href = '/login';
    });
  }
}

/* ==========================================================================
   Sidebar Navigation
   ========================================================================== */
function setupNavigation() {
  const navItems = document.querySelectorAll('.sidebar-item');
  const tabContents = document.querySelectorAll('.tab-content');
  const pageTitle = document.getElementById('current-page-title');

  navItems.forEach(item => {
    item.addEventListener('click', () => {
      const target = item.dataset.section;
      if (!target) return;

      navItems.forEach(i => i.classList.remove('active'));
      tabContents.forEach(c => c.classList.remove('active'));

      item.classList.add('active');
      const targetEl = document.getElementById(`section-${target}`);
      if (targetEl) targetEl.classList.add('active');

      if (pageTitle) {
        const text = item.querySelector('span')?.textContent || 'Dashboard';
        pageTitle.textContent = text;
      }

      // Lazy load section data
      switch (target) {
        case 'overview':
          loadDashboardOverview();
          break;
        case 'projects':
          loadProjects();
          break;
        case 'publications':
          loadPublications();
          break;
        case 'milestones':
          loadMilestones();
          break;
        case 'team':
          loadTeam();
          break;
        case 'activity':
          loadActivities();
          break;
      }
    });
  });
}

/* ==========================================================================
   Dashboard Overview & Charts
   ========================================================================== */
async function loadDashboardOverview() {
  try {
    const stats = await API.get('/api/stats');
    if (!stats) return;

    // Database mode badge
    const dbBadge = document.getElementById('db-mode-badge');
    const dbText = document.getElementById('db-mode-text');
    if (dbBadge && dbText) {
      dbText.textContent = stats.database_mode;
      if (!stats.database_mode.toLowerCase().includes('mysql')) {
        dbBadge.classList.add('fallback');
      } else {
        dbBadge.classList.remove('fallback');
      }
    }

    // KPIs
    document.getElementById('kpi-projects').textContent = stats.total_projects;
    document.getElementById('kpi-active-projects').textContent = stats.active_projects;
    document.getElementById('kpi-pubs').textContent = stats.total_publications;
    document.getElementById('kpi-funding').textContent = formatCurrency(stats.total_budget);
    document.getElementById('kpi-milestones').textContent = stats.pending_milestones;

    // Render Canvas Charts
    renderStatusChart(stats.projects_by_status || {});
    renderBudgetChart(stats.budget_by_department || {});

    // Recent activity on overview
    loadRecentActivitiesWidget();
  } catch (err) {
    showToast('Failed to load dashboard statistics', 'error');
  }
}

function renderStatusChart(statusData) {
  const canvas = document.getElementById('statusChart');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  const labels = Object.keys(statusData);
  const values = Object.values(statusData);
  const total = values.reduce((a, b) => a + b, 0);

  if (total === 0) {
    ctx.fillStyle = '#94a3b8';
    ctx.font = '14px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText('No project status data', canvas.width / 2, canvas.height / 2);
    return;
  }

  const colors = ['#2563eb', '#10b981', '#f59e0b', '#8b5cf6', '#ef4444'];
  const centerX = canvas.width / 2;
  const centerY = canvas.height / 2 - 20;
  const radius = 70;
  const innerRadius = 42;

  let startAngle = -Math.PI / 2;

  values.forEach((val, i) => {
    const sliceAngle = (val / total) * (2 * Math.PI);
    const color = colors[i % colors.length];

    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, startAngle, startAngle + sliceAngle);
    ctx.arc(centerX, centerY, innerRadius, startAngle + sliceAngle, startAngle, true);
    ctx.closePath();
    ctx.fillStyle = color;
    ctx.fill();

    startAngle += sliceAngle;
  });

  // Legend at bottom
  let legendX = 20;
  const legendY = canvas.height - 30;
  ctx.font = '11px sans-serif';
  ctx.textAlign = 'left';

  labels.forEach((lbl, i) => {
    const color = colors[i % colors.length];
    ctx.fillStyle = color;
    ctx.fillRect(legendX, legendY, 10, 10);
    ctx.fillStyle = '#475569';
    ctx.fillText(`${lbl} (${values[i]})`, legendX + 14, legendY + 9);
    legendX += ctx.measureText(`${lbl} (${values[i]})`).width + 24;
  });
}

function renderBudgetChart(budgetData) {
  const canvas = document.getElementById('budgetChart');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  const departments = Object.keys(budgetData);
  const budgets = Object.values(budgetData);

  if (departments.length === 0) {
    ctx.fillStyle = '#94a3b8';
    ctx.font = '14px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText('No budget data recorded', canvas.width / 2, canvas.height / 2);
    return;
  }

  const maxBudget = Math.max(...budgets, 100000);
  const chartHeight = 160;
  const bottomY = 200;
  const barWidth = 40;
  const gap = 30;
  let startX = 40;

  departments.forEach((dept, i) => {
    const val = budgets[i];
    const barHeight = (val / maxBudget) * chartHeight;
    const x = startX + i * (barWidth + gap);
    const y = bottomY - barHeight;

    // Draw Bar with rounded top
    ctx.fillStyle = '#2563eb';
    ctx.fillRect(x, y, barWidth, barHeight);

    // Label value on top of bar
    ctx.fillStyle = '#1e293b';
    ctx.font = 'bold 10px sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText(`$${(val / 1000).toFixed(0)}k`, x + barWidth / 2, y - 6);

    // Department name underneath
    ctx.fillStyle = '#64748b';
    ctx.font = '11px sans-serif';
    const shortDept = dept.split(' ')[0];
    ctx.fillText(shortDept, x + barWidth / 2, bottomY + 16);
  });
}

async function loadRecentActivitiesWidget() {
  const container = document.getElementById('overview-activity-list');
  if (!container) return;

  try {
    const acts = await API.get('/api/activities?limit=5');
    if (!acts || acts.length === 0) {
      container.innerHTML = '<p style="color: #94a3b8;">No recent activities.</p>';
      return;
    }

    container.innerHTML = acts.map(a => `
      <div style="display:flex; justify-content:space-between; align-items:center; padding: 0.75rem 0; border-bottom: 1px solid #f1f5f9;">
        <div>
          <div style="font-size: 0.875rem; font-weight: 600; color: #1e293b;">${escapeHtml(a.user_name)}: ${escapeHtml(a.action)}</div>
          <div style="font-size: 0.8rem; color: #64748b;">${escapeHtml(a.details)}</div>
        </div>
        <div style="font-size: 0.75rem; color: #94a3b8;">${formatDate(a.created_at)}</div>
      </div>
    `).join('');
  } catch (e) {}
}

/* ==========================================================================
   Projects Management
   ========================================================================== */
async function loadProjects() {
  const tbody = document.getElementById('projects-table-body');
  if (!tbody) return;

  const search = document.getElementById('dash-project-search')?.value.trim() || '';
  const dept = document.getElementById('dash-dept-filter')?.value || '';
  const status = document.getElementById('dash-status-filter')?.value || '';

  try {
    tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; padding:2rem; color: #64748b;">Loading projects...</td></tr>';
    
    let url = `/api/projects?search=${encodeURIComponent(search)}`;
    if (dept) url += `&department=${encodeURIComponent(dept)}`;
    if (status) url += `&status=${encodeURIComponent(status)}`;

    const projects = await API.get(url);
    allProjects = projects || [];

    if (allProjects.length === 0) {
      tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; padding:2rem; color: #64748b;">No projects found matching the filter criteria.</td></tr>';
      return;
    }

    tbody.innerHTML = allProjects.map(p => `
      <tr>
        <td style="font-weight: 700; color: #2563eb;">${escapeHtml(p.code)}</td>
        <td style="font-weight: 600; color: #0f172a; max-width: 280px;">
          ${escapeHtml(p.title)}
          <div style="font-size: 0.75rem; color: #64748b; font-weight: normal; margin-top: 2px;">
            Grant: ${escapeHtml(p.grant_source || 'Internal Institutional')}
          </div>
        </td>
        <td>${escapeHtml(p.lead_researcher_name)}</td>
        <td><span class="tag-chip">${escapeHtml(p.department)}</span></td>
        <td><span class="badge badge-${escapeHtml(p.status.toLowerCase().replace(/\s+/g, '-'))}">${escapeHtml(p.status)}</span></td>
        <td style="font-weight: 700; color: #0f172a;">${formatCurrency(p.budget)}</td>
        <td>
          <div style="display:flex; gap:0.4rem;">
            <button class="btn btn-secondary btn-sm" onclick="viewProjectDetails(${p.id})">Details</button>
            <button class="btn btn-secondary btn-sm" onclick="openEditProjectModal(${p.id})">Edit</button>
            <button class="btn btn-danger btn-sm" onclick="deleteProject(${p.id})">Delete</button>
          </div>
        </td>
      </tr>
    `).join('');
  } catch (err) {
    tbody.innerHTML = `<tr><td colspan="7" style="color:#ef4444; text-align:center;">Failed to load projects: ${escapeHtml(err.message)}</td></tr>`;
  }
}

function setupProjectHandlers() {
  const searchInput = document.getElementById('dash-project-search');
  const deptFilter = document.getElementById('dash-dept-filter');
  const statusFilter = document.getElementById('dash-status-filter');

  let debounce;
  if (searchInput) {
    searchInput.addEventListener('input', () => {
      clearTimeout(debounce);
      debounce = setTimeout(loadProjects, 300);
    });
  }
  if (deptFilter) deptFilter.addEventListener('change', loadProjects);
  if (statusFilter) statusFilter.addEventListener('change', loadProjects);

  // New Project Form Submit
  const form = document.getElementById('project-form');
  if (form) {
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const id = document.getElementById('project-id').value;
      const title = document.getElementById('project-title').value.trim();
      const code = document.getElementById('project-code').value.trim();
      const department = document.getElementById('project-dept').value;
      const status = document.getElementById('project-status').value;
      const budget = parseFloat(document.getElementById('project-budget').value) || 0;
      const grantSource = document.getElementById('project-grant').value.trim();
      const startDate = document.getElementById('project-start').value;
      const endDate = document.getElementById('project-end').value;
      const tags = document.getElementById('project-tags').value.trim();
      const abstract = document.getElementById('project-abstract').value.trim();

      const payload = {
        title,
        code,
        department,
        status,
        budget,
        grant_source: grantSource,
        start_date: startDate,
        end_date: endDate,
        tags,
        abstract
      };

      try {
        if (id) {
          await API.put(`/api/projects/${id}`, payload);
          showToast('Project updated successfully', 'success');
        } else {
          await API.post('/api/projects', payload);
          showToast('New project registered successfully', 'success');
        }
        closeModal('modal-project');
        loadProjects();
        loadDashboardOverview();
      } catch (err) {
        showToast(err.message || 'Error saving project', 'error');
      }
    });
  }
}

window.openNewProjectModal = function() {
  document.getElementById('project-modal-title').textContent = 'Initiate New Research Project';
  document.getElementById('project-id').value = '';
  document.getElementById('project-form').reset();
  openModal('modal-project');
};

window.openEditProjectModal = async function(id) {
  const p = allProjects.find(item => item.id === id);
  if (!p) return;

  document.getElementById('project-modal-title').textContent = 'Edit Research Project';
  document.getElementById('project-id').value = p.id;
  document.getElementById('project-title').value = p.title;
  document.getElementById('project-code').value = p.code;
  document.getElementById('project-dept').value = p.department;
  document.getElementById('project-status').value = p.status;
  document.getElementById('project-budget').value = p.budget;
  document.getElementById('project-grant').value = p.grant_source || '';
  document.getElementById('project-start').value = p.start_date;
  document.getElementById('project-end').value = p.end_date;
  document.getElementById('project-tags').value = p.tags || '';
  document.getElementById('project-abstract').value = p.abstract;

  openModal('modal-project');
};

window.deleteProject = async function(id) {
  if (!confirm('Are you sure you want to remove this project? This will also unlink associated milestones.')) {
    return;
  }
  try {
    await API.delete(`/api/projects/${id}`);
    showToast('Project successfully removed', 'success');
    loadProjects();
    loadDashboardOverview();
  } catch (err) {
    showToast(err.message || 'Failed to delete project', 'error');
  }
};

window.viewProjectDetails = async function(id) {
  try {
    const p = await API.get(`/api/projects/${id}`);
    const pubs = await API.get(`/api/publications?project_id=${id}`);
    const milestones = await API.get(`/api/milestones?project_id=${id}`);

    const detailsEl = document.getElementById('project-details-body');
    if (!detailsEl) return;

    detailsEl.innerHTML = `
      <div style="margin-bottom: 1.25rem;">
        <span class="project-code">${escapeHtml(p.code)}</span>
        <span class="badge badge-${escapeHtml(p.status.toLowerCase().replace(/\s+/g, '-'))}" style="margin-left: 0.5rem;">${escapeHtml(p.status)}</span>
        <h2 style="font-size: 1.35rem; color: #0f172a; margin-top: 0.5rem;">${escapeHtml(p.title)}</h2>
        <p style="color: #64748b; font-size: 0.9rem; margin-top: 0.5rem; line-height: 1.6;">${escapeHtml(p.abstract)}</p>
      </div>

      <div style="display:grid; grid-template-columns: 1fr 1fr; gap: 1rem; background: #f8fafc; padding: 1rem; border-radius: 8px; margin-bottom: 1.5rem; border: 1px solid #e2e8f0;">
        <div><strong>Lead PI:</strong> ${escapeHtml(p.lead_researcher_name)}</div>
        <div><strong>Department:</strong> ${escapeHtml(p.department)}</div>
        <div><strong>Budget Grant:</strong> ${formatCurrency(p.budget)}</div>
        <div><strong>Sponsor:</strong> ${escapeHtml(p.grant_source || 'Internal')}</div>
        <div><strong>Duration:</strong> ${formatDate(p.start_date)} &rarr; ${formatDate(p.end_date)}</div>
        <div><strong>Tags:</strong> ${escapeHtml(p.tags || 'General')}</div>
      </div>

      <h4 style="font-size: 1rem; margin-bottom: 0.75rem; color: #0f172a;">Project Milestones (${milestones.length})</h4>
      <div style="display:flex; flex-direction:column; gap:0.5rem; margin-bottom: 1.5rem;">
        ${milestones.length === 0 ? '<p style="color:#94a3b8; font-size:0.85rem;">No milestones set for this project.</p>' : milestones.map(m => `
          <div style="display:flex; justify-content:space-between; align-items:center; background:#fff; padding:0.6rem 0.9rem; border:1px solid #e2e8f0; border-radius:6px;">
            <div>
              <span style="font-weight:600; font-size:0.9rem;">${escapeHtml(m.title)}</span>
              <div style="font-size:0.75rem; color:#64748b;">Due: ${formatDate(m.due_date)}</div>
            </div>
            <span class="badge badge-${escapeHtml(m.status.toLowerCase().replace(/\s+/g, '-'))}">${escapeHtml(m.status)}</span>
          </div>
        `).join('')}
      </div>

      <h4 style="font-size: 1rem; margin-bottom: 0.75rem; color: #0f172a;">Associated Publications (${pubs.length})</h4>
      <div style="display:flex; flex-direction:column; gap:0.5rem;">
        ${pubs.length === 0 ? '<p style="color:#94a3b8; font-size:0.85rem;">No publications linked to this project yet.</p>' : pubs.map(pub => `
          <div style="background:#fff; padding:0.6rem 0.9rem; border:1px solid #e2e8f0; border-radius:6px;">
            <div style="font-weight:600; font-size:0.9rem;">${escapeHtml(pub.title)}</div>
            <div style="font-size:0.8rem; color:#64748b;">${escapeHtml(pub.authors)} &bull; <em>${escapeHtml(pub.journal_or_conference)} (${pub.year})</em></div>
          </div>
        `).join('')}
      </div>
    `;

    openModal('modal-project-details');
  } catch (err) {
    showToast('Failed to load project details', 'error');
  }
};

/* ==========================================================================
   Publications Management
   ========================================================================== */
async function loadPublications() {
  const tbody = document.getElementById('pubs-table-body');
  if (!tbody) return;

  try {
    tbody.innerHTML = '<tr><td colspan="6" style="text-align:center; padding:2rem; color: #64748b;">Loading papers repository...</td></tr>';
    const pubs = await API.get('/api/publications');

    if (!pubs || pubs.length === 0) {
      tbody.innerHTML = '<tr><td colspan="6" style="text-align:center; padding:2rem; color: #64748b;">No publications registered yet.</td></tr>';
      return;
    }

    tbody.innerHTML = pubs.map(pub => `
      <tr>
        <td style="font-weight:600; color: #0f172a; max-width: 320px;">
          ${escapeHtml(pub.title)}
          <div style="font-size:0.75rem; color:#64748b; margin-top:2px;">
            DOI: <a href="https://doi.org/${escapeHtml(pub.doi)}" target="_blank" style="color:#2563eb; text-decoration:underline;">${escapeHtml(pub.doi || 'In Press')}</a>
          </div>
        </td>
        <td>${escapeHtml(pub.authors)}</td>
        <td><em>${escapeHtml(pub.journal_or_conference)}</em></td>
        <td><span class="badge badge-active">${pub.year}</span></td>
        <td style="font-weight:700; color:#10b981;">${pub.citations} Citations</td>
        <td>
          <button class="btn btn-danger btn-sm" onclick="deletePublication(${pub.id})">Delete</button>
        </td>
      </tr>
    `).join('');
  } catch (err) {
    tbody.innerHTML = `<tr><td colspan="6" style="color:#ef4444; text-align:center;">Failed to load publications</td></tr>`;
  }
}

function setupPublicationHandlers() {
  const form = document.getElementById('pub-form');
  if (form) {
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const title = document.getElementById('pub-title').value.trim();
      const authors = document.getElementById('pub-authors').value.trim();
      const venue = document.getElementById('pub-venue').value.trim();
      const year = parseInt(document.getElementById('pub-year').value) || 2025;
      const doi = document.getElementById('pub-doi').value.trim();
      const citations = parseInt(document.getElementById('pub-citations').value) || 0;
      const projectIDStr = document.getElementById('pub-project-id').value;
      const projectID = projectIDStr ? parseInt(projectIDStr) : null;

      try {
        await API.post('/api/publications', {
          title,
          authors,
          journal_or_conference: venue,
          year,
          doi,
          citations,
          project_id: projectID
        });
        showToast('Publication registered successfully', 'success');
        closeModal('modal-pub');
        form.reset();
        loadPublications();
        loadDashboardOverview();
      } catch (err) {
        showToast(err.message || 'Error adding publication', 'error');
      }
    });
  }
}

window.openNewPublicationModal = function() {
  // Populate projects dropdown
  const select = document.getElementById('pub-project-id');
  if (select) {
    select.innerHTML = '<option value="">-- Optional: Link to Research Project --</option>' +
      allProjects.map(p => `<option value="${p.id}">${escapeHtml(p.code)} - ${escapeHtml(p.title)}</option>`).join('');
  }
  openModal('modal-pub');
};

window.deletePublication = async function(id) {
  if (!confirm('Are you sure you want to delete this publication?')) return;
  try {
    await API.delete(`/api/publications/${id}`);
    showToast('Publication removed', 'success');
    loadPublications();
    loadDashboardOverview();
  } catch (err) {
    showToast(err.message || 'Failed to delete publication', 'error');
  }
};

/* ==========================================================================
   Milestones & Deliverables
   ========================================================================== */
async function loadMilestones() {
  const container = document.getElementById('milestones-container');
  if (!container) return;

  try {
    container.innerHTML = '<div style="text-align:center; padding:2rem; color:#64748b;">Loading milestones...</div>';
    const milestones = await API.get('/api/milestones');

    if (!milestones || milestones.length === 0) {
      container.innerHTML = '<div style="text-align:center; padding:2rem; color:#64748b;">No milestones found. Click "+ New Milestone" to create one.</div>';
      return;
    }

    container.innerHTML = milestones.map(m => `
      <div class="milestone-item">
        <div class="milestone-left">
          <input type="checkbox" class="milestone-check" ${m.status === 'Achieved' ? 'checked' : ''} onchange="toggleMilestone(${m.id}, this.checked)">
          <div>
            <div class="milestone-title ${m.status === 'Achieved' ? 'text-decoration: line-through;' : ''}">${escapeHtml(m.title)}</div>
            <div class="milestone-desc">${escapeHtml(m.description || 'No additional description')} &bull; <strong style="color:#2563eb;">Due: ${formatDate(m.due_date)}</strong></div>
          </div>
        </div>
        <div>
          <span class="badge badge-${escapeHtml(m.status.toLowerCase().replace(/\s+/g, '-'))}">${escapeHtml(m.status)}</span>
        </div>
      </div>
    `).join('');
  } catch (err) {
    container.innerHTML = `<div style="color:#ef4444; text-align:center;">Failed to load milestones</div>`;
  }
}

window.toggleMilestone = async function(id, isChecked) {
  const newStatus = isChecked ? 'Achieved' : 'In Progress';
  try {
    await API.put(`/api/milestones/${id}`, { status: newStatus });
    showToast(`Milestone marked as ${newStatus}`, 'success');
    loadMilestones();
    loadDashboardOverview();
  } catch (err) {
    showToast('Failed to update milestone status', 'error');
  }
};

function setupMilestoneHandlers() {
  const form = document.getElementById('milestone-form');
  if (form) {
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const projectID = parseInt(document.getElementById('milestone-project-id').value);
      const title = document.getElementById('milestone-title').value.trim();
      const dueDate = document.getElementById('milestone-due').value;
      const status = document.getElementById('milestone-status').value;
      const description = document.getElementById('milestone-desc').value.trim();

      if (!projectID || !title) {
        showToast('Please select a project and enter milestone title', 'warning');
        return;
      }

      try {
        await API.post('/api/milestones', {
          project_id: projectID,
          title,
          due_date: dueDate,
          status,
          description
        });
        showToast('Milestone created', 'success');
        closeModal('modal-milestone');
        form.reset();
        loadMilestones();
        loadDashboardOverview();
      } catch (err) {
        showToast(err.message || 'Error creating milestone', 'error');
      }
    });
  }
}

window.openNewMilestoneModal = function() {
  const select = document.getElementById('milestone-project-id');
  if (select) {
    select.innerHTML = '<option value="">-- Select Target Research Project --</option>' +
      allProjects.map(p => `<option value="${p.id}">${escapeHtml(p.code)} - ${escapeHtml(p.title)}</option>`).join('');
  }
  openModal('modal-milestone');
};

/* ==========================================================================
   Team & Researchers
   ========================================================================== */
async function loadTeam() {
  const container = document.getElementById('team-grid');
  if (!container) return;

  try {
    container.innerHTML = '<div style="grid-column: 1/-1; text-align:center; padding:2rem; color:#64748b;">Loading faculty and researchers...</div>';
    const users = await API.get('/api/researchers');

    if (!users || users.length === 0) {
      container.innerHTML = '<div style="grid-column: 1/-1; text-align:center;">No researchers registered.</div>';
      return;
    }

    container.innerHTML = users.map(u => `
      <div class="researcher-card">
        <img src="${escapeHtml(u.avatar || 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=250&q=80')}" alt="${escapeHtml(u.name)}">
        <div class="researcher-name">${escapeHtml(u.name)}</div>
        <div class="researcher-dept">${escapeHtml(u.role)} &bull; ${escapeHtml(u.department)}</div>
        <div class="researcher-inst">${escapeHtml(u.institution)}</div>
        <div style="margin-top: 1rem; font-size: 0.8rem; color: #64748b;">
          ✉️ <a href="mailto:${escapeHtml(u.email)}" style="color:#2563eb;">${escapeHtml(u.email)}</a>
        </div>
      </div>
    `).join('');
  } catch (err) {
    container.innerHTML = '<div style="grid-column: 1/-1; color:#ef4444; text-align:center;">Failed to load team</div>';
  }
}

/* ==========================================================================
   Activity Log
   ========================================================================== */
async function loadActivities() {
  const container = document.getElementById('activity-timeline');
  if (!container) return;

  try {
    container.innerHTML = '<div style="padding:2rem; color:#64748b;">Loading audit logs...</div>';
    const acts = await API.get('/api/activities?limit=25');

    if (!acts || acts.length === 0) {
      container.innerHTML = '<div style="padding:2rem; color:#64748b;">No activities logged yet.</div>';
      return;
    }

    container.innerHTML = acts.map(a => `
      <div class="timeline-item">
        <div class="timeline-dot"></div>
        <div class="timeline-content">
          <div class="timeline-header">
            <div>
              <span class="timeline-user">${escapeHtml(a.user_name)}</span> &bull;
              <span class="timeline-action">${escapeHtml(a.action)}</span>
            </div>
            <span class="timeline-time">${formatDate(a.created_at)}</span>
          </div>
          <div class="timeline-detail">${escapeHtml(a.details)}</div>
        </div>
      </div>
    `).join('');
  } catch (err) {
    container.innerHTML = '<div style="color:#ef4444;">Failed to load activity logs</div>';
  }
}

/* ==========================================================================
   Modal Utilities
   ========================================================================== */
function setupModals() {
  document.querySelectorAll('.modal-overlay').forEach(modal => {
    modal.addEventListener('click', (e) => {
      if (e.target === modal) {
        closeModal(modal.id);
      }
    });
  });
}

window.openModal = function(id) {
  const el = document.getElementById(id);
  if (el) el.classList.add('active');
};

window.closeModal = function(id) {
  const el = document.getElementById(id);
  if (el) el.classList.remove('active');
};
