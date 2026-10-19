/**
 * ScholarSync - Authentication Controller
 */

document.addEventListener('DOMContentLoaded', () => {
  // If already logged in, redirect to dashboard
  if (API.getToken() && API.getUser()) {
    window.location.href = '/dashboard';
    return;
  }

  setupTabs();
  setupDemoChips();
  setupForms();
});

function setupTabs() {
  const tabBtns = document.querySelectorAll('.auth-tab-btn');
  const forms = document.querySelectorAll('.auth-form');

  tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      tabBtns.forEach(b => b.classList.remove('active'));
      forms.forEach(f => f.classList.remove('active'));

      btn.classList.add('active');
      const target = btn.dataset.tab;
      const targetForm = document.getElementById(`${target}-form`);
      if (targetForm) {
        targetForm.classList.add('active');
      }
    });
  });
}

function setupDemoChips() {
  const chips = document.querySelectorAll('.demo-chip');
  chips.forEach(chip => {
    chip.addEventListener('click', () => {
      const email = chip.dataset.email;
      const password = chip.dataset.password || 'password123';

      const emailInput = document.getElementById('login-email');
      const passInput = document.getElementById('login-password');
      if (emailInput && passInput) {
        emailInput.value = email;
        passInput.value = password;
        showToast(`Selected demo credentials for: ${email}`, 'info', 2000);
      }
    });
  });
}

function setupForms() {
  // Login Form
  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('login-email').value.trim();
      const password = document.getElementById('login-password').value;
      const btn = loginForm.querySelector('button[type="submit"]');

      if (!email || !password) {
        showToast('Please enter both email and password', 'warning');
        return;
      }

      btn.disabled = true;
      btn.textContent = 'Authenticating...';

      try {
        const res = await API.post('/api/auth/login', { email, password });
        API.setToken(res.token);
        API.setUser(res.user);
        showToast(`Welcome back, ${res.user.name}!`, 'success');
        setTimeout(() => {
          window.location.href = '/dashboard';
        }, 600);
      } catch (err) {
        showToast(err.message || 'Login failed. Please check your credentials.', 'error');
        btn.disabled = false;
        btn.textContent = 'Sign In to Portal';
      }
    });
  }

  // Register Form
  const registerForm = document.getElementById('register-form');
  if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('reg-name').value.trim();
      const email = document.getElementById('reg-email').value.trim();
      const password = document.getElementById('reg-password').value;
      const role = document.getElementById('reg-role').value;
      const department = document.getElementById('reg-dept').value.trim();
      const institution = document.getElementById('reg-inst').value.trim();
      const btn = registerForm.querySelector('button[type="submit"]');

      if (!name || !email || !password) {
        showToast('Please fill in all required fields', 'warning');
        return;
      }

      btn.disabled = true;
      btn.textContent = 'Creating Account...';

      try {
        const res = await API.post('/api/auth/register', {
          name,
          email,
          password,
          role,
          department,
          institution
        });
        API.setToken(res.token);
        API.setUser(res.user);
        showToast('Account created successfully! Redirecting...', 'success');
        setTimeout(() => {
          window.location.href = '/dashboard';
        }, 600);
      } catch (err) {
        showToast(err.message || 'Registration failed', 'error');
        btn.disabled = false;
        btn.textContent = 'Create Academic Account';
      }
    });
  }
}
