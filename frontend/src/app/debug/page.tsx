'use client';

import React, { useState, useEffect } from 'react';

export default function DebugPage() {
  const [tokenInfo, setTokenInfo] = useState<{
    token: string | null;
    isValid: boolean;
    decodedToken: any;
    error: string | null;
  }>({
    token: null,
    isValid: false,
    decodedToken: null,
    error: null
  });

  useEffect(() => {
    // Get token from localStorage
    const token = localStorage.getItem('access_token');
    
    try {
      if (token) {
        // Decode token (JWT is in format header.payload.signature)
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
          return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        const decodedToken = JSON.parse(jsonPayload);
        
        // Check if token is expired
        const currentTime = Math.floor(Date.now() / 1000);
        const isExpired = decodedToken.exp < currentTime;
        
        setTokenInfo({
          token,
          isValid: !isExpired,
          decodedToken,
          error: isExpired ? 'Token is expired' : null
        });
      } else {
        setTokenInfo({
          token: null,
          isValid: false,
          decodedToken: null,
          error: 'No token found in localStorage'
        });
      }
    } catch (error) {
      setTokenInfo({
        token: token,
        isValid: false,
        decodedToken: null,
        error: `Error decoding token: ${error instanceof Error ? error.message : String(error)}`
      });
    }
  }, []);

  const makeTestRequest = async () => {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/assets`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const result = await response.text();
      alert(`Status: ${response.status}\nResponse: ${result}`);
    } catch (error) {
      alert(`Error: ${error instanceof Error ? error.message : String(error)}`);
    }
  };

  const clearStorage = () => {
    localStorage.clear();
    sessionStorage.clear();
    window.location.reload();
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100 p-4">
      <div className="w-full max-w-3xl bg-white shadow-lg rounded-lg p-6">
        <h1 className="text-3xl font-bold mb-6 text-center">Authentication Debug</h1>
        
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-2">Token Status</h2>
          <div className="bg-gray-50 p-4 rounded border">
            <p><strong>Token Present:</strong> {tokenInfo.token ? 'Yes' : 'No'}</p>
            <p><strong>Token Valid:</strong> {tokenInfo.isValid ? 'Yes' : 'No'}</p>
            {tokenInfo.error && (
              <p className="text-red-600"><strong>Error:</strong> {tokenInfo.error}</p>
            )}
          </div>
        </div>

        {tokenInfo.token && (
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-2">Token Details</h2>
            <div className="bg-gray-50 p-4 rounded border overflow-auto">
              <p className="text-xs break-all"><strong>Raw Token:</strong> {tokenInfo.token}</p>
              {tokenInfo.decodedToken && (
                <div className="mt-2">
                  <p><strong>User ID:</strong> {tokenInfo.decodedToken.user_id}</p>
                  <p><strong>Username:</strong> {tokenInfo.decodedToken.username}</p>
                  <p><strong>Role:</strong> {tokenInfo.decodedToken.role}</p>
                  <p><strong>Issued At:</strong> {new Date(tokenInfo.decodedToken.iat * 1000).toLocaleString()}</p>
                  <p><strong>Expires At:</strong> {new Date(tokenInfo.decodedToken.exp * 1000).toLocaleString()}</p>
                  <p><strong>Current Time:</strong> {new Date().toLocaleString()}</p>
                </div>
              )}
            </div>
          </div>
        )}

        <div className="flex flex-col space-y-2 sm:flex-row sm:space-y-0 sm:space-x-2 justify-center">
          <button 
            onClick={makeTestRequest}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
          >
            Test API Request
          </button>
          <button 
            onClick={clearStorage}
            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-opacity-50"
          >
            Clear Storage
          </button>
        </div>
      </div>
    </div>
  );
}