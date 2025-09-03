'use client';

import { Asset } from '@/app/assets/page';
import { X } from 'lucide-react';

interface AssetQRModalProps {
  show: boolean;
  onClose: () => void;
  onDownload: () => void;
  asset: Asset | null;
  qrCodeDataUrl: string;
}

export const AssetQRModal = ({ show, onClose, onDownload, asset, qrCodeDataUrl }: AssetQRModalProps) => {
  if (!show || !asset) return null;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-10 mx-auto p-5 border w-4/5 max-w-4xl shadow-lg rounded-md bg-white">
        <div className="mt-3">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900">Asset QR Code</h3>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          
          <div className="flex gap-8">
            {/* Left side - QR Code */}
            <div className="flex flex-col items-center">
              <div className="bg-white p-4 border rounded-lg">
                {qrCodeDataUrl ? (
                  <img
                    src={qrCodeDataUrl}
                    alt="QR Code"
                    className="w-48 h-48"
                  />
                ) : (
                  <div className="w-48 h-48 flex items-center justify-center bg-gray-50 border-2 border-dashed border-gray-300 rounded-lg">
                    <div className="text-center">
                      <div className="w-16 h-16 mx-auto mb-2 bg-gray-200 rounded-lg flex items-center justify-center">
                        <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V6a1 1 0 00-1-1H5a1 1 0 00-1 1v1a1 1 0 001 1zm12 0h2a1 1 0 001-1V6a1 1 0 00-1-1h-2a1 1 0 00-1 1v1a1 1 0 001 1zM5 20h2a1 1 0 001-1v-1a1 1 0 00-1-1H5a1 1 0 00-1 1v1a1 1 0 001 1z" />
                        </svg>
                      </div>
                      <p className="text-sm text-gray-500">Generating QR Code...</p>
                    </div>
                  </div>
                )}
              </div>
              <div className="mt-2 text-center">
                <span className="text-sm font-mono text-gray-600">
                  {asset.serial_number}
                </span>
              </div>
            </div>
            
            {/* Right side - Company Info */}
            <div className="flex-1 flex flex-col justify-center">
              <div className="text-left">
                <div className="text-xs text-gray-500 mb-1">PROPERTY OF</div>
                <div className="text-2xl font-bold text-gray-800 mb-2">SAMS Corporation</div>
                <div className="text-sm text-gray-600 mb-4">Smart Asset Management System</div>
                <div className="text-sm text-gray-600 mb-6">(021) 555-0123</div>
                
                <div className="bg-green-100 text-green-800 px-3 py-2 rounded-md text-sm font-mono text-left">
                  Asset ID: {asset.id?.substring(0, 8).toUpperCase()}
                </div>
              </div>
            </div>
          </div>
          
          <div className="flex justify-end gap-3 mt-6">
            <button
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              Close
            </button>
            <button
              onClick={onDownload}
              className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
            >
              Download QR
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
