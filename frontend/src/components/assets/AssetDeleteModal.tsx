'use client';

import { Asset } from '@/app/assets/page';

interface AssetDeleteModalProps {
  show: boolean;
  onClose: () => void;
  onDelete: () => void;
  asset: Asset | null;
}

export const AssetDeleteModal = ({ show, onClose, onDelete, asset }: AssetDeleteModalProps) => {
  if (!show) return null;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div className="mt-3 text-center">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Delete Asset</h3>
          <p className="text-sm text-gray-500 mb-6">
            Are you sure you want to delete &quot;{asset?.name}&quot;? This action cannot be undone.
          </p>
          <div className="flex justify-center space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              onClick={onDelete}
              className="px-4 py-2 bg-red-600 border border-transparent rounded-md text-sm font-medium text-white hover:bg-red-700"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
