/**
 * GRO Unified Header System v1.0 - JavaScript Behaviors
 *
 * Provides:
 * - Theme toggle (light/dark)
 * - Agent runtime health check
 * - Mobile menu toggle
 * - Catalog status check
 * - Automatic active link detection
 *
 * Usage:
 * <script src="[path]/gro-header.js"></script>
 *
 * Or initialize manually:
 * GROHeader.init({ agentRuntimeUrl: 'http://localhost:8081', ... })
 */

(function() {
  'use strict';

  const GROHeader = {
    config: {
      agentRuntimeUrl: 'http://localhost:8081',
      catalogUrl: 'http://localhost:3000',
      healthCheckInterval: 30000, // 30 seconds
      catalogCheckInterval: 60000, // 60 seconds
      themeStorageKey: 'gro-theme',
      debug: false
    },

    state: {
      agentStatus: 'checking',
      catalogCount: '...',
      theme: 'light',
      mobileMenuOpen: false
    },

    /**
     * Initialize the header
     */
    init: function(options = {}) {
      // Merge config
      this.config = { ...this.config, ...options };

      // Log init if debug
      if (this.config.debug) {
        console.log('[GROHeader] Initializing...', this.config);
      }

      // Initialize theme
      this.initTheme();

      // Initialize event listeners
      this.initEventListeners();

      // Start health checks
      this.startHealthChecks();

      // Set active nav link
      this.setActiveNavLink();

      if (this.config.debug) {
        console.log('[GROHeader] Initialized successfully');
      }
    },

    /**
     * Initialize theme from localStorage or system preference
     */
    initTheme: function() {
      // Check localStorage first
      const savedTheme = localStorage.getItem(this.config.themeStorageKey);

      // Fall back to system preference
      const prefersDark = window.matchMedia &&
                          window.matchMedia('(prefers-color-scheme: dark)').matches;

      this.state.theme = savedTheme || (prefersDark ? 'dark' : 'light');

      // Apply theme
      this.applyTheme(this.state.theme);
    },

    /**
     * Apply theme to document
     */
    applyTheme: function(theme) {
      const isDark = theme === 'dark';

      // Set body class
      document.body.classList.toggle('dark', isDark);
      document.body.setAttribute('data-theme', theme);

      // Update theme icon
      const themeIcon = document.getElementById('themeIcon');
      if (themeIcon) {
        themeIcon.textContent = isDark ? '☀' : '🌙';
      }

      // Also check for .gro-header__theme-toggle
      const themeToggles = document.querySelectorAll('.gro-header__theme-toggle');
      themeToggles.forEach(toggle => {
        const icon = toggle.querySelector('span');
        if (icon) {
          icon.textContent = isDark ? '☀' : '🌙';
        }
      });

      this.state.theme = theme;

      if (this.config.debug) {
        console.log('[GROHeader] Theme applied:', theme);
      }
    },

    /**
     * Toggle theme
     */
    toggleTheme: function() {
      const newTheme = this.state.theme === 'dark' ? 'light' : 'dark';
      this.applyTheme(newTheme);
      localStorage.setItem(this.config.themeStorageKey, newTheme);

      if (this.config.debug) {
        console.log('[GROHeader] Theme toggled to:', newTheme);
      }
    },

    /**
     * Initialize all event listeners
     */
    initEventListeners: function() {
      // Theme toggle - support multiple selectors
      const themeToggles = [
        '#themeBadge',
        '.gro-header__theme-toggle',
        '[data-gro-theme-toggle]'
      ];

      themeToggles.forEach(selector => {
        const elements = document.querySelectorAll(selector);
        elements.forEach(el => {
          el.addEventListener('click', (e) => {
            e.preventDefault();
            this.toggleTheme();
          });
        });
      });

      // Mobile hamburger menu
      const hamburger = document.querySelector('.gro-header__hamburger');
      if (hamburger) {
        hamburger.addEventListener('click', () => {
          this.toggleMobileMenu();
        });
      }

      // Also support legacy hamburger ID
      const legacyHamburger = document.getElementById('hamburger');
      if (legacyHamburger) {
        legacyHamburger.addEventListener('click', () => {
          this.toggleMobileMenu();
        });
      }

      // Close mobile menu when clicking outside
      document.addEventListener('click', (e) => {
        const menu = document.querySelector('.gro-header__mobile-menu');
        const hamburger = document.querySelector('.gro-header__hamburger');

        if (menu && hamburger &&
            !menu.contains(e.target) &&
            !hamburger.contains(e.target) &&
            menu.classList.contains('active')) {
          this.toggleMobileMenu();
        }
      });

      if (this.config.debug) {
        console.log('[GROHeader] Event listeners initialized');
      }
    },

    /**
     * Toggle mobile menu
     */
    toggleMobileMenu: function() {
      this.state.mobileMenuOpen = !this.state.mobileMenuOpen;

      // Update menu
      const menu = document.querySelector('.gro-header__mobile-menu');
      if (menu) {
        menu.classList.toggle('active', this.state.mobileMenuOpen);
      }

      // Also support legacy mobile-menu ID
      const legacyMenu = document.getElementById('mobile-menu');
      if (legacyMenu) {
        legacyMenu.classList.toggle('active', this.state.mobileMenuOpen);
      }

      // Update hamburger aria-expanded
      const hamburger = document.querySelector('.gro-header__hamburger') ||
                       document.getElementById('hamburger');
      if (hamburger) {
        hamburger.setAttribute('aria-expanded', this.state.mobileMenuOpen);
      }

      if (this.config.debug) {
        console.log('[GROHeader] Mobile menu toggled:', this.state.mobileMenuOpen);
      }
    },

    /**
     * Start health check polling
     */
    startHealthChecks: function() {
      // Check agent runtime health
      this.checkAgentHealth();
      setInterval(() => this.checkAgentHealth(), this.config.healthCheckInterval);

      // Check catalog status
      this.checkCatalogStatus();
      setInterval(() => this.checkCatalogStatus(), this.config.catalogCheckInterval);
    },

    /**
     * Check agent runtime health
     */
    checkAgentHealth: async function() {
      try {
        const response = await fetch(`${this.config.agentRuntimeUrl}/health`, {
          method: 'GET',
          cache: 'no-cache'
        });

        if (response.ok) {
          this.updateAgentStatus('connected');
        } else {
          this.updateAgentStatus('disconnected');
        }
      } catch (error) {
        this.updateAgentStatus('disconnected');

        if (this.config.debug) {
          console.warn('[GROHeader] Agent health check failed:', error);
        }
      }
    },

    /**
     * Update agent status UI
     */
    updateAgentStatus: function(status) {
      this.state.agentStatus = status;

      // Update status badge - support multiple selectors
      const selectors = [
        '#agent-status',
        '.gro-header__badge--connected',
        '[data-gro-agent-status]'
      ];

      selectors.forEach(selector => {
        const badges = document.querySelectorAll(selector);
        badges.forEach(badge => {
          // Update classes
          badge.classList.toggle('connected', status === 'connected');
          badge.classList.toggle('disconnected', status === 'disconnected');
          badge.classList.toggle('checking', status === 'checking');

          // Update text
          const textEl = badge.querySelector('span:not(.gro-header__badge-dot)');
          if (textEl && !textEl.classList.contains('gro-header__badge-dot')) {
            textEl.textContent = status === 'connected' ? 'Connected' :
                                status === 'disconnected' ? 'Offline' : 'Checking...';
          }

          // Update colors via style (for legacy support)
          if (status === 'connected') {
            badge.style.background = 'rgba(34, 197, 94, 0.2)';
            badge.style.color = 'rgb(74, 222, 128)';
          } else if (status === 'disconnected') {
            badge.style.background = 'rgba(239, 68, 68, 0.2)';
            badge.style.color = 'rgb(248, 113, 113)';
          } else {
            badge.style.background = 'rgba(234, 179, 8, 0.2)';
            badge.style.color = 'rgb(252, 211, 77)';
          }
        });
      });

      if (this.config.debug) {
        console.log('[GROHeader] Agent status updated:', status);
      }
    },

    /**
     * Check catalog status
     */
    checkCatalogStatus: async function() {
      try {
        const response = await fetch(`${this.config.catalogUrl}/health`, {
          method: 'GET',
          cache: 'no-cache'
        });

        if (response.ok) {
          const data = await response.json();
          const count = data.entry_count || data.count || '...';
          this.updateCatalogCount(count);
        } else {
          this.updateCatalogCount('...');
        }
      } catch (error) {
        this.updateCatalogCount('...');

        if (this.config.debug) {
          console.warn('[GROHeader] Catalog check failed:', error);
        }
      }
    },

    /**
     * Update catalog count display
     */
    updateCatalogCount: function(count) {
      this.state.catalogCount = count;

      // Update catalog count elements
      const selectors = [
        '#catalogEntryCount',
        '[data-gro-catalog-count]'
      ];

      selectors.forEach(selector => {
        const elements = document.querySelectorAll(selector);
        elements.forEach(el => {
          el.textContent = count;
        });
      });

      if (this.config.debug) {
        console.log('[GROHeader] Catalog count updated:', count);
      }
    },

    /**
     * Set active nav link based on current URL
     */
    setActiveNavLink: function() {
      const currentUrl = window.location.origin + window.location.pathname;
      const currentPort = window.location.port;

      // Get all nav links
      const navLinks = document.querySelectorAll('.gro-header__nav-link, .verb-link');

      navLinks.forEach(link => {
        const linkUrl = link.getAttribute('href');

        // Check if link matches current location
        const isActive = linkUrl === currentUrl ||
                        linkUrl.includes(`:${currentPort}`) ||
                        (linkUrl.includes('localhost:8080') && currentPort === '8080') ||
                        (linkUrl.includes('localhost:5173') && currentPort === '5173') ||
                        (linkUrl.includes('localhost:8084') && currentPort === '8084') ||
                        (linkUrl.includes('localhost:3000') && currentPort === '3000') ||
                        (linkUrl.includes('localhost:8083') && currentPort === '8083');

        // Update active class
        if (isActive) {
          link.classList.add('gro-header__nav-link--active', 'active');
        } else {
          link.classList.remove('gro-header__nav-link--active', 'active');
        }
      });

      if (this.config.debug) {
        console.log('[GROHeader] Active nav link set for port:', currentPort);
      }
    },

    /**
     * Public method to manually refresh health checks
     */
    refresh: function() {
      this.checkAgentHealth();
      this.checkCatalogStatus();

      if (this.config.debug) {
        console.log('[GROHeader] Manual refresh triggered');
      }
    },

    /**
     * Public method to get current state
     */
    getState: function() {
      return { ...this.state };
    }
  };

  // Auto-initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      GROHeader.init();
    });
  } else {
    GROHeader.init();
  }

  // Expose GROHeader globally for manual control
  window.GROHeader = GROHeader;

  // Support legacy global function names
  window.toggleTheme = () => GROHeader.toggleTheme();

})();
