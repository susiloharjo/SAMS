// API utility functions with automatic authentication

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Get authentication headers
export const getAuthHeaders = (): HeadersInit => {
  const token = localStorage.getItem('access_token');
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
};

// Check if token is expired
export const isTokenExpired = (token: string): boolean => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const currentTime = Math.floor(Date.now() / 1000);
    return payload.exp < currentTime;
  } catch (error) {
    return true;
  }
};

// Refresh token silently
export const refreshTokenSilently = async (): Promise<boolean> => {
  try {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) {
      return false;
    }

    const response = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    const data = await response.json();

    if (response.ok && !data.error) {
      localStorage.setItem('access_token', data.data.access_token);
      localStorage.setItem('refresh_token', data.data.refresh_token);
      return true;
    }
    return false;
  } catch (error) {
    console.error('Token refresh failed:', error);
    return false;
  }
};

// Authenticated fetch with automatic token refresh
export const authenticatedFetch = async (
  url: string,
  options: RequestInit = {}
): Promise<Response> => {
  // Get token from localStorage (should be set during login)
  let token = localStorage.getItem('access_token');
  
  // Check if token is expired and refresh if needed
  if (token && isTokenExpired(token)) {
    console.log('Token expired, attempting to refresh...');
    const refreshed = await refreshTokenSilently();
    if (refreshed) {
      console.log('Token refreshed successfully');
      token = localStorage.getItem('access_token');
    } else {
      console.log('Token refresh failed, clearing auth data');
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      
      // Don't redirect here to avoid loops, just let the request fail
      // The AuthContext will handle redirection on its next check
    }
  }
  
  // Add authorization header with token
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
    ...options.headers,
  };

  // Make the request with the token
  const response = await fetch(url, {
    ...options,
    headers,
  });

  return response;
};

// API request helpers
export const api = {
  get: (endpoint: string) => authenticatedFetch(`${API_BASE_URL}${endpoint}`),
  
  post: (endpoint: string, data: any) => authenticatedFetch(`${API_BASE_URL}${endpoint}`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  put: (endpoint: string, data: any) => authenticatedFetch(`${API_BASE_URL}${endpoint}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  delete: (endpoint: string) => authenticatedFetch(`${API_BASE_URL}${endpoint}`, {
    method: 'DELETE',
  }),
};
