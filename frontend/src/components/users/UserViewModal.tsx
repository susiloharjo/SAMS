'use client';

import { X, User, Calendar, Building, Shield } from 'lucide-react';

interface User {
  id: string;
  username: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  department_id?: string;
  department?: {
    id: string;
    name: string;
  };
  is_active: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

interface UserViewModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User;
}

export default function UserViewModal({
  isOpen,
  onClose,
  user,
}: UserViewModalProps) {
  if (!isOpen) return null;

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const getRoleBadgeColor = (role: string) => {
    switch (role) {
      case 'admin':
        return 'bg-red-100 text-red-800';
      case 'manager':
        return 'bg-blue-100 text-blue-800';
      case 'user':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusBadgeColor = (isActive: boolean) => {
    return isActive
      ? 'bg-green-100 text-green-800'
      : 'bg-red-100 text-red-800';
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div className="flex items-center gap-2">
            <User className="w-5 h-5 text-blue-600" />
            <h2 className="text-lg font-semibold text-gray-900">
              User Details
            </h2>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {/* User Avatar and Basic Info */}
          <div className="flex items-start gap-6 mb-6">
            <div className="w-20 h-20 bg-blue-100 rounded-full flex items-center justify-center">
              <span className="text-blue-600 font-semibold text-2xl">
                {user.first_name.charAt(0)}{user.last_name.charAt(0)}
              </span>
            </div>
            
            <div className="flex-1">
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                {user.first_name} {user.last_name}
              </h3>
              <p className="text-gray-600 mb-1">@{user.username}</p>
              <p className="text-gray-600 mb-3">{user.email}</p>
              
              <div className="flex gap-3">
                <span className={`inline-flex px-3 py-1 text-sm font-semibold rounded-full ${getRoleBadgeColor(user.role)}`}>
                  <Shield className="w-4 h-4 mr-1" />
                  {user.role}
                </span>
                <span className={`inline-flex px-3 py-1 text-sm font-semibold rounded-full ${getStatusBadgeColor(user.is_active)}`}>
                  {user.is_active ? 'Active' : 'Inactive'}
                </span>
              </div>
            </div>
          </div>

          {/* User Details Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Personal Information */}
            <div className="space-y-4">
              <h4 className="font-medium text-gray-900 border-b border-gray-200 pb-2">
                Personal Information
              </h4>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  Full Name
                </label>
                <p className="text-gray-900">{user.first_name} {user.last_name}</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  Username
                </label>
                <p className="text-gray-900">@{user.username}</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  Email Address
                </label>
                <p className="text-gray-900">{user.email}</p>
              </div>
            </div>

            {/* System Information */}
            <div className="space-y-4">
              <h4 className="font-medium text-gray-900 border-b border-gray-200 pb-2">
                System Information
              </h4>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  User ID
                </label>
                <p className="text-gray-900 font-mono text-sm">{user.id}</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  Role
                </label>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRoleBadgeColor(user.role)}`}>
                  {user.role}
                </span>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-500 mb-1">
                  Status
                </label>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusBadgeColor(user.is_active)}`}>
                  {user.is_active ? 'Active' : 'Inactive'}
                </span>
              </div>
            </div>
          </div>

          {/* Department Information */}
          <div className="mt-6 space-y-4">
            <h4 className="font-medium text-gray-900 border-b border-gray-200 pb-2">
              Department Information
            </h4>
            
            <div className="flex items-center gap-2">
              <Building className="w-4 h-4 text-gray-400" />
              <span className="text-gray-900">
                {user.department?.name || 'No Department Assigned'}
              </span>
            </div>
          </div>

          {/* Timestamps */}
          <div className="mt-6 space-y-4">
            <h4 className="font-medium text-gray-900 border-b border-gray-200 pb-2">
              Timestamps
            </h4>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-gray-400" />
                <div>
                  <label className="block text-sm font-medium text-gray-500 mb-1">
                    Created
                  </label>
                  <p className="text-gray-900 text-sm">{formatDate(user.created_at)}</p>
                </div>
              </div>
              
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-gray-400" />
                <div>
                  <label className="block text-sm font-medium text-gray-500 mb-1">
                    Last Updated
                  </label>
                  <p className="text-gray-900 text-sm">{formatDate(user.updated_at)}</p>
                </div>
              </div>
              
              {user.last_login && (
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4 text-gray-400" />
                  <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                      Last Login
                    </label>
                    <p className="text-gray-900 text-sm">{formatDate(user.last_login)}</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Actions */}
          <div className="mt-8 pt-6 border-t border-gray-200">
            <button
              onClick={onClose}
              className="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
