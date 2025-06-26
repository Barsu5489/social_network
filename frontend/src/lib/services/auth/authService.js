import { writable } from 'svelte/store';

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
