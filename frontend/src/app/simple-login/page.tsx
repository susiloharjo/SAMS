'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function SimpleLogin() {
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('user.1001');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const router = useRouter();

  const handleLogin = async () => {
    setLoading(true);
    setMessage('Attempting login...');

    try {
      console.log('Making request to:', 'http://localhost:8080/api/v1/auth/login');
      
      const response = await fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      console.log('Response status:', response.status);
      console.log('Response ok:', response.ok);

      const data = await response.json();
      console.log('Response data:', data);

      if (response.ok && data.data) {
        // Store tokens directly
        localStorage.clear(); // Clear everything first
        localStorage.setItem('access_token', data.data.access_token);
        localStorage.setItem('refresh_token', data.data.refresh_token);
        localStorage.setItem('user', JSON.stringify(data.data.user));

        setMessage('Login successful! Redirecting...');
        
        // Force redirect
        setTimeout(() => {
          window.location.href = '/assets';
        }, 1000);
      } else {
        setMessage('Login failed: ' + (data.message || 'Unknown error'));
      }
    } catch (error) {
      console.error('Login error:', error);
      setMessage('Network error: ' + error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow">
        <h1 className="text-2xl font-bold text-center">Simple Login Test</h1>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <button
            onClick={handleLogin}
            disabled={loading}
            className="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
          
          <button
            onClick={() => localStorage.clear()}
            className="w-full py-2 px-4 bg-red-600 text-white rounded-md hover:bg-red-700"
          >
            Clear Storage
          </button>
        </div>
        
        <div className="mt-4 p-4 bg-gray-100 rounded">
          <h3 className="font-medium">Status:</h3>
          <p className="text-sm">{message || 'Ready to login'}</p>
        </div>
      </div>
    </div>
  );
}
