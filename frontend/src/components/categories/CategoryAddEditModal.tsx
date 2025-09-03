'use client';

import { useEffect, useState } from 'react';
import { X } from 'lucide-react';

interface Category {
  id: string;
  name: string;
  description: string;
}

interface CategoryAddEditModalProps {
  show: boolean;
  onClose: () => void;
  onSave: (category: Partial<Category>) => void;
  category: Category | null;
}

export function CategoryAddEditModal({ show, onClose, onSave, category }: CategoryAddEditModalProps) {
  const [formData, setFormData] = useState({ name: '', description: '' });

  useEffect(() => {
    if (show && category) {
      setFormData({ name: category.name, description: category.description });
    } else {
      setFormData({ name: '', description: '' });
    }
  }, [show, category]);

  if (!show) {
    return null;
  }

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(formData);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold">{category ? 'Edit' : 'Add'} Category</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <X className="w-6 h-6" />
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">Name</label>
              <input
                type="text"
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500"
                required
              />
            </div>
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">Description</label>
              <textarea
                id="description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-green-500 focus:border-green-500"
              />
            </div>
          </div>
          <div className="mt-6 flex justify-end space-x-2">
            <button type="button" onClick={onClose} className="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300">Cancel</button>
            <button type="submit" className="px-4 py-2 rounded bg-green-600 text-white hover:bg-green-700">Save</button>
          </div>
        </form>
      </div>
    </div>
  );
}
