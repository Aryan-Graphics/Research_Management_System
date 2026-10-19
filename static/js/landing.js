/**
 * ScholarSync - Front / Landing Page Controller
 */

document.addEventListener('DOMContentLoaded', () => {
  loadStats();
  loadPublicProjects();
  loadFeaturedPublications();
  setupEventListeners();
  checkAuthStatus();
});

function checkAuthStatus() {
  const user = API.getUser();
  const navActions = document.getElementById('nav-actions');
  if (user && navActions) {
    navActions.innerHTML = `
      <a href="/dashboard" class="btn btn-primary">Go to Dashboard &rarr;</a>
    `;
  }
}

async function loadStats() {
  try {
    const stats = await API.get('/api/stats');
    if (stats) {
      animateValue('stat-projects', stats.total_projects || 0);
      animateValue('stat-publications', stats.total_publications || 0);
      animateValue('stat-citations', stats.total_citations || 0);
      const budgetMil = ((stats.total_budget || 0) / 1000000).toFixed(1);
      const elBudget = document.getElementById('stat-funding');
      if (elBudget) elBudget.textContent = `$${budgetMil}M`;

      // Database badge in footer
      const dbIndicator = document.getElementById('footer-db-status');
      if (dbIndicator && stats.database_mode) {
        dbIndicator.textContent = stats.database_mode;
      }
    }
  } catch (err) {
    console.warn('Could not load live stats:', err);
  }
}

function animateValue(id, target) {
  const el = document.getElementById(id);
  if (!el) return;
  let current = 0;
  const step = Math.max(1, Math.floor(target / 20));
  const timer = setInterval(() => {
    current += step;
    if (current >= target) {
      current = target;
      clearInterval(timer);
    }
    el.textContent = current.toLocaleString();
  }, 35);
}

async function loadPublicProjects() {
  const container = document.getElementById('projects-container');
  if (!container) return;

  const searchInput = document.getElementById('search-input');
  const deptSelect = document.getElementById('dept-filter');
  const statusSelect = document.getElementById('status-filter');

  const search = searchInput ? searchInput.value.trim() : '';
  const department = deptSelect ? deptSelect.value : '';
  const status = statusSelect ? statusSelect.value : '';

  try {
    container.innerHTML = '<div style="grid-column: 1/-1; text-align: center; padding: 3rem; color: #64748b;">Loading active research initiatives...</div>';
    
    let url = `/api/projects?search=${encodeURIComponent(search)}`;
    if (department) url += `&department=${encodeURIComponent(department)}`;
    if (status) url += `&status=${encodeURIComponent(status)}`;

    const projects = await API.get(url);
    if (!projects || projects.length === 0) {
      container.innerHTML = `
        <div style="grid-column: 1/-1; text-align: center; padding: 3rem; background: #fff; border-radius: 12px; border: 1px solid #e2e8f0;">
          <p style="font-size: 1.1rem; color: #475569; font-weight: 600;">No research projects found matching your criteria.</p>
          <p style="font-size: 0.9rem; color: #94a3b8; margin-top: 0.5rem;">Try adjusting your search query or department filter.</p>
        </div>
      `;
      return;
    }

    container.innerHTML = projects.map(p => `
      <div class="project-card">
        <div>
          <div class="project-card-header">
            <span class="project-code">${escapeHtml(p.code)}</span>
            <span class="badge badge-${escapeHtml(p.status.toLowerCase().replace(/\s+/g, '-'))}">${escapeHtml(p.status)}</span>
          </div>
          <h3 class="project-title">${escapeHtml(p.title)}</h3>
          <p class="project-abstract">${escapeHtml(p.abstract)}</p>
        </div>

        <div>
          <div class="tags-wrap">
            ${(p.tags || '').split(',').filter(t => t.trim()).map(t => `<span class="tag-chip">${escapeHtml(t.trim())}</span>`).join('')}
          </div>

          <div class="project-meta">
            <div class="meta-row">
              <span>Principal Investigator:</span>
              <span class="meta-val">${escapeHtml(p.lead_researcher_name || 'Assigned Lead')}</span>
            </div>
            <div class="meta-row">
              <span>Department:</span>
              <span class="meta-val">${escapeHtml(p.department)}</span>
            </div>
            <div class="meta-row">
              <span>Grant Budget:</span>
              <span class="meta-val" style="color: #2563eb;">${formatCurrency(p.budget)}</span>
            </div>
          </div>
        </div>
      </div>
    `).join('');
  } catch (err) {
    container.innerHTML = `<div style="grid-column: 1/-1; color: #ef4444; text-align: center;">Failed to load projects: ${escapeHtml(err.message)}</div>`;
  }
}

async function loadFeaturedPublications() {
  const tbody = document.getElementById('publications-tbody');
  if (!tbody) return;

  try {
    const pubs = await API.get('/api/publications');
    if (!pubs || pubs.length === 0) {
      tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color: #64748b;">No publications registered yet.</td></tr>';
      return;
    }

    tbody.innerHTML = pubs.slice(0, 5).map(pub => `
      <tr>
        <td style="font-weight: 600; color: #0f172a; max-width: 320px;">
          ${escapeHtml(pub.title)}
          <div style="font-size: 0.775rem; color: #64748b; font-weight: normal; margin-top: 2px;">
            DOI: <a href="https://doi.org/${escapeHtml(pub.doi)}" target="_blank" style="color: #2563eb; text-decoration: underline;">${escapeHtml(pub.doi || 'In Press')}</a>
          </div>
        </td>
        <td style="color: #475569;">${escapeHtml(pub.authors)}</td>
        <td style="color: #475569; font-style: italic;">${escapeHtml(pub.journal_or_conference)}</td>
        <td><span class="badge badge-active">${pub.year}</span></td>
        <td style="font-weight: 700; color: #10b981;">${pub.citations} Citations</td>
      </tr>
    `).join('');
  } catch (err) {
    console.warn('Failed to load publications:', err);
  }
}

function setupEventListeners() {
  const searchInput = document.getElementById('search-input');
  const deptSelect = document.getElementById('dept-filter');
  const statusSelect = document.getElementById('status-filter');

  let debounceTimer;
  if (searchInput) {
    searchInput.addEventListener('input', () => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(loadPublicProjects, 300);
    });
  }

  if (deptSelect) {
    deptSelect.addEventListener('change', loadPublicProjects);
  }

  if (statusSelect) {
    statusSelect.addEventListener('change', loadPublicProjects);
  }
}
