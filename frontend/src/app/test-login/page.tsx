'use client';

import { useState } from 'react';

export default function TestLogin() {
  const [result, setResult] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const testBackendConnection = async () => {
    setLoading(true);
    setResult('Testing backend connection...');
    
    try {
      const response = await fetch('http://localhost:8080/health');
      const data = await response.json();
      setResult(prev => prev + '\n✅ Backend health: ' + JSON.stringify(data));
    } catch (error) {
      setResult(prev => prev + '\n❌ Backend health failed: ' + error);
    }
    
    setLoading(false);
  };

  const testLogin = async () => {
    setLoading(true);
    setResult('Testing login...');
    
    try {
      const response = await fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: 'admin',
          password: 'user.1001'
        }),
      });
      
      const data = await response.json();
      setResult(prev => prev + '\n✅ Login response: ' + JSON.stringify(data, null, 2));
    } catch (error) {
      setResult(prev => prev + '\n❌ Login failed: ' + error);
    }
    
    setLoading(false);
  };

  const testLoginSimple = async () => {
    setLoading(true);
    setResult('Testing simple login...');
    
    try {
      const response = await fetch('http://localhost:8080/api/v1/test-login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({}),
      });
      
      const data = await response.json();
      setResult(prev => prev + '\n✅ Simple login response: ' + JSON.stringify(data, null, 2));
    } catch (error) {
      setResult(prev => prev + '\n❌ Simple login failed: ' + error);
    }
    
    setLoading(false);
  };

  return (
    <div className="min-h-screen p-8 bg-gray-50">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">SAMS Login Test</h1>
        
        <div className="space-y-4 mb-8">
          <button
            onClick={testBackendConnection}
            disabled={loading}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            Test Backend Connection
          </button>
          
          <button
            onClick={testLoginSimple}
            disabled={loading}
            className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50 ml-4"
          >
            Test Simple Login
          </button>
          
          <button
            onClick={testLogin}
            disabled={loading}
            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 ml-4"
          >
            Test Real Login
          </button>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-bold mb-4">Test Results:</h2>
          <pre className="whitespace-pre-wrap text-sm bg-gray-100 p-4 rounded">
            {result || 'Click a button to test...'}
          </pre>
        </div>
      </div>
    </div>
  );
}
