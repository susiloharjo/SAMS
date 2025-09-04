'use client';

import { useState, useEffect, useCallback } from 'react';
import { CategoryAddEditModal } from '@/components/categories/CategoryAddEditModal';
import { api } from '@/utils/api';

// You will need to create these components, similar to the ones for Assets.
// import { CategoryControls } from '@/components/categories/CategoryControls';
// import { CategoriesTable } from '@/components/categories/CategoriesTable';
// import { CategoryDeleteModal } from '@/components/categories/CategoryDeleteModal';

interface Category {
  id: string;
  name: string;
  description: string;
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dataVersion, setDataVersion] = useState(0); // Used to trigger refetch

  // State for modals
  const [showAddEditModal, setShowAddEditModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);

  const fetchCategories = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.get('/api/v1/categories');
      const data = await response.json();
      setCategories(data.data || []);
    } catch (err) {
      setError('Could not load categories. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories, dataVersion]);

  const handleAddCategory = () => {
    setSelectedCategory(null);
    setShowAddEditModal(true);
  };

  const handleEditCategory = (category: Category) => {
    setSelectedCategory(category);
    setShowAddEditModal(true);
  };

  const handleDeleteCategory = (category: Category) => {
    setSelectedCategory(category);
    setShowDeleteModal(true);
  };

  const handleSave = async (categoryData: Partial<Category>) => {
    try {
      const response = selectedCategory
        ? await api.put(`/api/v1/categories/${selectedCategory.id}`, categoryData)
        : await api.post('/api/v1/categories', categoryData);

      if (response.ok) {
        setShowAddEditModal(false);
        setDataVersion(v => v + 1); // Trigger refetch
      } else {
        const errorData = await response.json();
        alert(`Failed to save category: ${errorData.message}`);
      }
    } catch (error) {
      console.error('Error saving category:', error);
      alert('An error occurred while saving the category.');
    }
  };
  
  const confirmDelete = async () => {
    if (!selectedCategory) return;

    try {
      const response = await api.delete(`/api/v1/categories/${selectedCategory.id}`);

      if (response.ok) {
        setShowDeleteModal(false);
        setDataVersion(v => v + 1); // Trigger refetch
      } else {
        const errorData = await response.json();
        alert(`Failed to delete category: ${errorData.message}`);
      }
    } catch (error) {
      console.error('Error deleting category:', error);
      alert('An error occurred while deleting the category.');
    }
  };

  return (
    <div className="p-4 md:p-6">
      <h1 className="text-2xl font-semibold text-gray-800 mb-4">Manage Categories</h1>
      
      <div className="mb-4">
        {/* Placeholder for CategoryControls */}
        <button onClick={handleAddCategory} className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700">
          Add Category
        </button>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        {loading && <p className="p-4">Loading...</p>}
        {error && <p className="p-4 text-red-600">{error}</p>}
        {!loading && !error && (
          // Placeholder for CategoriesTable
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
                {categories.map((cat) => (
                  <tr key={cat.id}>
                    <td className="px-6 py-4 whitespace-nowrap">{cat.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap">{cat.description}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <button onClick={() => handleEditCategory(cat)} className="text-indigo-600 hover:text-indigo-900">Edit</button>
                      <button onClick={() => handleDeleteCategory(cat)} className="text-red-600 hover:text-red-900 ml-4">Delete</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <CategoryAddEditModal
        show={showAddEditModal}
        onClose={() => setShowAddEditModal(false)}
        onSave={handleSave}
        category={selectedCategory}
      />

      {/* Placeholder for Delete Modal */}
      {showDeleteModal && (
         <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
            <h2 className="text-lg font-semibold mb-2">Confirm Deletion</h2>
            <p>Are you sure you want to delete the category "{selectedCategory?.name}"?</p>
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
