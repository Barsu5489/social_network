// stores/user.js
import { writable } from 'svelte/store';

const storedUser = localStorage.getItem('user');
export const user = writable(storedUser ? JSON.parse(storedUser) : null);
export const isAuthenticated = writable(!!storedUser);

user.subscribe(value => {
  if (value) {
    localStorage.setItem('user', JSON.stringify(value));
  } else {
    localStorage.removeItem('user');
  }
});