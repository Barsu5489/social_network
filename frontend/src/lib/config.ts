// This file contains the base URL for the backend API.
// By centralizing it here, we can easily change it in one place if needed.

export const API_BASE_URL = process.env.NODE_ENV === 'production' 
  ? (typeof window !== 'undefined' ? 'http://localhost:8080' : 'http://backend:3000')
  : 'http://localhost:3000';
