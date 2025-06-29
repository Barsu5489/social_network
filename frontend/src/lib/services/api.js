const API_URL = 'http://localhost:3000/api';

export async function fetchWithSession(path, options = {}) {
  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {})
    }
  });

  if (!response.ok) {
  	const errorText = await response.text();
  	console.error('API Request failed:', response.status, errorText);
  	throw new Error(`Request failed with status ${response.status}: ${errorText}` || 'Request failed');
  }

  const data = await response.json();
  console.log('API Response:', data);
  return data;
}