import { writable } from 'svelte/store';
import { fetchWithSession } from './api.js';

export const user = writable(null);
export const isAuthenticated = writable(false);

const API_URL = 'http://localhost:3000/api';

export async function register(userData) {
  const response = await fetch(`${API_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Failed to register');
  }

  return response.json();
}

export async function login({ email, password }) {
  const response = await fetch(`${API_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    credentials: 'include', // ensures cookie is saved
    body: JSON.stringify({ email, password })
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || 'Invalid credentials');
  }

  const data = await response.json();
  user.set(data.data);
  isAuthenticated.set(true);

  return data;
}

export async function logout() {
  const response = await fetch(`${API_URL}/logout`, {
    method: 'POST',
    credentials: 'include'
  });

  if (!response.ok) {
    throw new Error('Failed to logout');
  }

  user.set(null);
  isAuthenticated.set(false);
}

// Universal profile fetcher
export async function getProfile(targetId = '') {
  try {
    const url = targetId
      ? `/profile?target_id=${encodeURIComponent(targetId)}`
      : `/profile`;

    const data = await fetchWithSession(url);

    if (!targetId) {
      user.set(data.data.user);
      isAuthenticated.set(true);
    }

    return data.data;
  } catch (error) {
    console.error('Error loading profile:', error.message, error.stack);

    if (!targetId) {
      user.set(null);
      isAuthenticated.set(false);
    }

    throw error;
  }
}

// Alias for fetching current user (no target ID)
export const getCurrentUser = () => getProfile();

export async function updateProfile(aboutMe) {
  try {
    const response = await fetchWithSession('/profile', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ about_me: aboutMe })
    });

    // Refresh profile on success
    await getProfile();
    return true;
  } catch (error) {
    console.error("Failed to update profile:", error);
    throw error;
  }
}