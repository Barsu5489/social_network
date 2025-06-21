import { writable } from 'svelte/store';

export const user = writable(null);
export const isAuthenticated = writable(false);

const API_URL = 'http://localhost:3000/api';

export async function register({ firstName, lastName, email, password }) {
  const response = await fetch(`${API_URL}/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      FirstName: firstName,
      LastName: lastName,
      Email: email,
      PasswordHash: password,
      Nickname: "",
      DateOfBirth: "",
      AboutMe: "",
      AvatarURL: "",
      IsPrivate: false,
      CreatedAt: 0,
      UpdatedAt: 0,
      DeletedAt: null,
    })
  });

  if (!response.ok) {
    throw new Error('Failed to register');
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
    throw new Error('Invalid credentials');
  }

  const data = await response.json();
  user.set(data.data);
  isAuthenticated.set(true);

  return data;
}

export async function logout() {
  const response = await fetch(`${API_URL}/logout`, {
    method: 'POST'
  });

  if (!response.ok) {
    throw new Error('Failed to logout');
  }

  user.set(null);
  isAuthenticated.set(false);
}
