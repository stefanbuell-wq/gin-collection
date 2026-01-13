// Authentication JavaScript Module
// Handles login, registration, logout, and session management

const API_BASE = 'api/';

/**
 * Register a new user
 */
async function register(username, email, password, fullName = null) {
    const response = await fetch(`${API_BASE}?action=register`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username,
            email,
            password,
            full_name: fullName
        })
    });

    const data = await response.json();

    if (!response.ok || data.error) {
        throw new Error(data.error || 'Registrierung fehlgeschlagen');
    }

    return data.user;
}

/**
 * Login user
 */
async function login(usernameOrEmail, password) {
    const response = await fetch(`${API_BASE}?action=login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: usernameOrEmail,
            email: usernameOrEmail,
            password
        })
    });

    const data = await response.json();

    if (!response.ok || data.error) {
        throw new Error(data.error || 'Login fehlgeschlagen');
    }

    // Store user info in sessionStorage
    sessionStorage.setItem('user', JSON.stringify(data.user));
    sessionStorage.setItem('authenticated', 'true');

    return data.user;
}

/**
 * Logout user
 */
async function logout() {
    try {
        await fetch(`${API_BASE}?action=logout`, {
            method: 'POST'
        });
    } catch (error) {
        console.error('Logout error:', error);
    }

    // Clear session storage
    sessionStorage.removeItem('user');
    sessionStorage.removeItem('authenticated');

    // Redirect to login
    window.location.href = 'login.html';
}

/**
 * Check if user is authenticated
 */
function isAuthenticated() {
    return sessionStorage.getItem('authenticated') === 'true';
}

/**
 * Get current user from session
 */
function getCurrentUser() {
    const userJson = sessionStorage.getItem('user');
    return userJson ? JSON.parse(userJson) : null;
}

/**
 * Verify authentication with server
 */
async function verifyAuth() {
    try {
        const response = await fetch(`${API_BASE}?action=me`);
        const data = await response.json();

        if (response.ok && data.user) {
            // Update session storage with fresh data
            sessionStorage.setItem('user', JSON.stringify(data.user));
            sessionStorage.setItem('authenticated', 'true');
            return data.user;
        } else {
            // Not authenticated
            sessionStorage.removeItem('user');
            sessionStorage.removeItem('authenticated');
            return null;
        }
    } catch (error) {
        console.error('Auth verification error:', error);
        return null;
    }
}

/**
 * Require authentication - redirect to login if not authenticated
 */
async function requireAuth() {
    const user = await verifyAuth();

    if (!user) {
        window.location.href = 'login.html';
        return null;
    }

    return user;
}

/**
 * Make authenticated API request
 */
async function apiRequest(endpoint, options = {}) {
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json'
        }
    };

    const mergedOptions = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...options.headers
        }
    };

    const response = await fetch(`${API_BASE}${endpoint}`, mergedOptions);

    // Handle 401 Unauthorized
    if (response.status === 401) {
        sessionStorage.removeItem('user');
        sessionStorage.removeItem('authenticated');
        window.location.href = 'login.html';
        throw new Error('Nicht authentifiziert');
    }

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || 'API-Anfrage fehlgeschlagen');
    }

    return data;
}

// Export functions for use in other scripts
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        register,
        login,
        logout,
        isAuthenticated,
        getCurrentUser,
        verifyAuth,
        requireAuth,
        apiRequest
    };
}
