'use client';

import { useState, useEffect, useCallback } from 'react';
import { DepartmentAddEditModal } from '@/components/departments/DepartmentAddEditModal';
import { api } from '@/utils/api';

interface Department {
  id: string;
  name: string;
  description: string;
}

export default function DepartmentsPage() {
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dataVersion, setDataVersion] = useState(0);

  // State for modals
  const [showAddEditModal, setShowAddEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedDepartment, setSelectedDepartment] = useState<Department | null>(null);

  const fetchDepartments = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.get('/api/v1/departments');
      const data = await response.json();
      setDepartments(data.data || []);
    } catch (err) {
      setError('Could not load departments. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDepartments();
  }, [fetchDepartments, dataVersion]);

  const handleAddDepartment = () => {
    setSelectedDepartment(null);
    setShowAddEditModal(true);
  };

  const handleEditDepartment = (department: Department) => {
    setSelectedDepartment(department);
    setShowAddEditModal(true);
  };

  const handleDeleteDepartment = (department: Department) => {
    setSelectedDepartment(department);
    setShowDeleteModal(true);
  };

  const handleSave = async (departmentData: Partial<Department>) => {
    try {
      const response = selectedDepartment
        ? await api.put(`/api/v1/departments/${selectedDepartment.id}`, departmentData)
        : await api.post('/api/v1/departments', departmentData);

      if (response.ok) {
        setShowAddEditModal(false);
        setDataVersion(v => v + 1);
      } else {
        const errorData = await response.json();
        alert(`Failed to save department: ${errorData.message}`);
      }
    } catch (error) {
      console.error('Error saving department:', error);
      alert('An error occurred while saving the department.');
    }
  };
  
  const confirmDelete = async () => {
    if (!selectedDepartment) return;

    try {
      const response = await api.delete(`/api/v1/departments/${selectedDepartment.id}`);

      if (response.ok) {
        setShowDeleteModal(false);
        setDataVersion(v => v + 1);
      } else {
        const errorData = await response.json();
        alert(`Failed to delete department: ${errorData.message}`);
      }
    } catch (error) {
      console.error('Error deleting department:', error);
      alert('An error occurred while deleting the department.');
    }
  };

  return (
    <div className="p-4 md:p-6">
      <h1 className="text-2xl font-semibold text-gray-800 mb-4">Manage Departments</h1>
      
      <div className="mb-4">
        <button onClick={handleAddDepartment} className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700">
          Add Department
        </button>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        {loading && <p className="p-4">Loading...</p>}
        {error && <p className="p-4 text-red-600">{error}</p>}
        {!loading && !error && (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {departments.map((dept) => (
                  <tr key={dept.id}>
                    <td className="px-6 py-4 whitespace-nowrap">{dept.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{dept.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button onClick={() => handleEditDepartment(dept)} className="text-indigo-600 hover:text-indigo-900">Edit</button>
                      <button onClick={() => handleDeleteDepartment(dept)} className="text-red-600 hover:text-red-900 ml-4">Delete</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <DepartmentAddEditModal
        show={showAddEditModal}
        onClose={() => setShowAddEditModal(false)}
        onSave={handleSave}
        department={selectedDepartment}
      />

      {/* Placeholder for Delete Modal */}
      {showDeleteModal && (
         <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
            <h2 className="text-lg font-semibold mb-2">Confirm Deletion</h2>
            <p>Are you sure you want to delete the department "{selectedDepartment?.name}"?</p>
            <div className="mt-6 flex justify-end space-x-2">
              <button onClick={() => setShowDeleteModal(false)} className="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300">Cancel</button>
              <button onClick={confirmDelete} className="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-700">Delete</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
