'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { ArrowLeft, MapPin, DollarSign, Calendar, Settings, Shield, Building, Tag, Hash, User, Clock, FileText, X } from 'lucide-react';
import AssetMap from '@/components/assets/AssetMap';
import { AssetAddEditModal } from '@/components/assets/AssetAddEditModal';
import { MapPicker } from '@/components/assets/MapPicker';

interface Asset {
  id: string;
  name: string;
  description: string;
  category_id: string;
  department_id: string;
  type: string;
  model: string;
  serial_number: string;
  manufacturer: string;
  acquisition_cost: number;
  current_value: number;
  status: string;
  condition: string;
  criticality: string;
  latitude?: number;
  longitude?: number;
  address: string;
  building_room: string;
  acquisition_date: string;
  expected_life_years: number;
  depreciation_rate: number;
  category?: any;
  department?: any;
}

export default function AssetDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [asset, setAsset] = useState<Asset | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Edit modal state
  const [showEditModal, setShowEditModal] = useState(false);
  const [editFormData, setEditFormData] = useState<Partial<Asset>>({});
  const [categories, setCategories] = useState<any[]>([]);
  const [departments, setDepartments] = useState<any[]>([]);

  // Map Picker State
  const [showMapPicker, setShowMapPicker] = useState(false);

  useEffect(() => {
    if (params.id) {
      fetchAsset(params.id as string);
      fetchRelatedData();
    }
  }, [params.id]);

  const fetchRelatedData = async () => {
    try {
      const [categoriesRes, departmentsRes] = await Promise.all([
        fetch('http://localhost:8080/api/v1/categories'),
        fetch('http://localhost:8080/api/v1/departments'),
      ]);
      
      const categoriesData = await categoriesRes.json();
      const departmentsData = await departmentsRes.json();

      setCategories(categoriesData.data || []);
      setDepartments(departmentsData.data || []);
    } catch (error) {
      console.error('Error fetching related data:', error);
    }
  }

  const fetchAsset = async (id: string) => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${id}`);
      if (!response.ok) {
        throw new Error('Failed to fetch asset');
      }
      const data = await response.json();
      setAsset(data.data);
      console.log('Asset data fetched:', data.data);
    } catch (err) {
      setError('Asset not found');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const handleEditAsset = () => {
    if (asset) {
      setEditFormData({
        name: asset.name,
        description: asset.description,
        category_id: asset.category_id,
        department_id: asset.department_id,
        type: asset.type,
        model: asset.model,
        serial_number: asset.serial_number,
        manufacturer: asset.manufacturer,
        acquisition_cost: asset.acquisition_cost,
        current_value: asset.current_value,
        status: asset.status,
        condition: asset.condition,
        criticality: asset.criticality,
        latitude: asset.latitude,
        longitude: asset.longitude,
        address: asset.address,
        building_room: asset.building_room,
        acquisition_date: asset.acquisition_date,
        expected_life_years: asset.expected_life_years
      });
      setShowEditModal(true);
    }
  };

  const handleSaveEdit = async (formData: Partial<Asset>) => {
    if (!asset) return;
    
    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${asset.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });
      
      if (response.ok) {
        setShowEditModal(false);
        // Refresh the asset data
        fetchAsset(asset.id);
      }
    } catch (error) {
      console.error('Error updating asset:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-green-600"></div>
      </div>
    );
  }

  if (error || !asset) {
    return (
      <div className="text-center">
        <h1 className="text-2xl font-bold text-red-600">{error || 'Asset could not be loaded.'}</h1>
      </div>
    );
  }

  return (
    <div className="p-1 md:p-4">
      <div className="mb-4">
        <button
          onClick={() => router.push('/assets')}
          className="inline-flex items-center text-sm font-medium text-gray-600 hover:text-gray-900"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Assets List
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          {/* Main Info */}
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
             <div className="flex flex-col md:flex-row md:items-start md:justify-between">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{asset.name}</h1>
                <div className="flex items-center text-sm text-gray-500 mt-2 space-x-4">
                  <span className="inline-flex items-center">
                    <Tag className="w-4 h-4 mr-1.5" />
                    {asset.category?.name || 'Uncategorized'}
                  </span>
                  <span className="inline-flex items-center">
                    <Building className="w-4 h-4 mr-1.5" />
                    {asset.department?.name || 'No Department'}
                  </span>
                </div>
              </div>
              <div className="mt-4 md:mt-0 md:ml-6 text-left md:text-right">
                <span className={`px-3 py-1 text-xs font-medium rounded-full ${
                    asset.status === 'active' ? 'bg-green-100 text-green-800' :
                    asset.status === 'maintenance' ? 'bg-yellow-100 text-yellow-800' :
                    asset.status === 'inactive' ? 'bg-gray-100 text-gray-800' : 'bg-red-100 text-red-800'
                  }`}>{asset.status.charAt(0).toUpperCase() + asset.status.slice(1)}</span>
              </div>
            </div>
            {asset.description && <p className="text-sm text-gray-600 mt-4">{asset.description}</p>}
          </div>

           {/* Financial Information */}
           <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <DollarSign className="w-5 h-5 mr-2 text-green-600" />
                Financial Information
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div className="text-center p-4 bg-green-50 rounded-lg">
                  <p className="text-sm font-medium text-gray-500">Acquisition Cost</p>
                  <p className="mt-1 text-2xl font-bold text-green-600">{formatCurrency(asset.acquisition_cost)}</p>
                </div>
                <div className="text-center p-4 bg-blue-50 rounded-lg">
                  <p className="text-sm font-medium text-gray-500">Current Value</p>
                  <p className="mt-1 text-2xl font-bold text-blue-600">{formatCurrency(asset.current_value)}</p>
                </div>
                <div className="text-center p-4 bg-purple-50 rounded-lg">
                  <p className="text-sm font-medium text-gray-500">Depreciation Rate</p>
                  <p className="mt-1 text-2xl font-bold text-purple-600">
                    {asset.acquisition_cost > 0 ? ((asset.acquisition_cost - asset.current_value) / asset.acquisition_cost * 100).toFixed(1) : 0}%
                  </p>
                </div>
              </div>
              
              {/* Financial Metrics */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-500">Total Depreciation</span>
                    <span className="text-lg font-semibold text-gray-900">
                      {formatCurrency(asset.acquisition_cost - asset.current_value)}
                    </span>
                  </div>
                  <div className="mt-2">
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-red-500 h-2 rounded-full" 
                        style={{ width: `${((asset.acquisition_cost - asset.current_value) / asset.acquisition_cost) * 100}%` }}
                      ></div>
                    </div>
                  </div>
                </div>
                <div className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-500">Value Retention</span>
                    <span className="text-lg font-semibold text-gray-900">
                      {((asset.current_value / asset.acquisition_cost) * 100).toFixed(1)}%
                    </span>
                  </div>
                  <div className="mt-2">
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-green-500 h-2 rounded-full" 
                        style={{ width: `${(asset.current_value / asset.acquisition_cost) * 100}%` }}
                      ></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Location Information */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <MapPin className="w-5 h-5 mr-2 text-green-600" />
                Location Information
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500">Address</label>
                  <p className="mt-1 text-sm text-gray-900">{asset.address || 'No address available'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">Building/Room</label>
                  <p className="mt-1 text-sm text-gray-900">{asset.building_room || 'N/A'}</p>
                </div>
                {asset.latitude && asset.longitude && (
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-500 mb-2">Coordinates</label>
                    <div className="bg-gray-50 p-3 rounded-lg font-mono text-sm">
                      {asset.latitude.toFixed(6)}, {asset.longitude.toFixed(6)}
                    </div>
                  </div>
                )}
              </div>
              
              {/* Asset Map */}
              {asset.latitude && asset.longitude ? (
                <AssetMap
                  latitude={asset.latitude}
                  longitude={asset.longitude}
                  assetName={asset.name}
                />
              ) : (
                <div className="bg-gray-50 rounded-lg h-32 flex items-center justify-center">
                  <p className="text-sm text-gray-500">No location coordinates available</p>
                </div>
              )}
            </div>

            {/* Lifecycle Information */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Clock className="w-5 h-5 mr-2 text-green-600" />
                Lifecycle Information
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500">Acquisition Date</label>
                  <p className="mt-1 text-sm text-gray-900">
                    {asset.acquisition_date ? formatDate(asset.acquisition_date) : 'N/A'}
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">Expected Lifespan</label>
                  <p className="mt-1 text-sm text-gray-900">{asset.expected_life_years} years</p>
                </div>
              </div>
            </div>

             {/* Compliance & Details */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Shield className="w-5 h-5 mr-2 text-green-600" />
                Compliance & Details
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-500">Serial Number</label>
                  <p className="mt-1 text-sm text-gray-900 font-mono">{asset.serial_number}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-500">Manufacturer</label>
                  <p className="mt-1 text-sm text-gray-900">{asset.manufacturer}</p>
                </div>
                 <div>
                  <label className="block text-sm font-medium text-gray-500">Model</label>
                  <p className="mt-1 text-sm text-gray-900">{asset.model}</p>
                </div>
                 <div>
                  <label className="block text-sm font-medium text-gray-500">Condition</label>
                  <p className="mt-1 text-sm text-gray-900 capitalize">{asset.condition}</p>
                </div>
                 <div>
                  <label className="block text-sm font-medium text-gray-500">Criticality</label>
                  <p className="mt-1 text-sm text-gray-900 capitalize">{asset.criticality}</p>
                </div>
              </div>
            </div>
            
        </div>

        {/* Right Sidebar */}
        <div className="space-y-6">
           {/* Quick Actions */}
           <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <Settings className="w-5 h-5 mr-2 text-green-600" />
                Quick Actions
              </h2>
              <div className="space-y-3">
                <button 
                  onClick={handleEditAsset}
                  className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
                >
                  Edit Asset
                </button>
                <button className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                  View QR Code
                </button>
                <button className="w-full px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 transition-colors">
                  Schedule Maintenance
                </button>
                <button
                  onClick={() => {
                    if (window.confirm('Are you sure you want to delete this asset? This action cannot be undone.')) {
                      fetch(`http://localhost:8080/api/v1/assets/${asset.id}`, {
                        method: 'DELETE',
                      }).then(response => {
                        if (response.ok) {
                          router.push('/assets');
                        } else {
                          alert('Failed to delete asset.');
                        }
                      });
                    }
                  }}
                  className="w-full px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
                >
                  Delete Asset
                </button>
              </div>
            </div>

            {/* Maintenance History */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                <FileText className="w-5 h-5 mr-2 text-green-600" />
                Maintenance History
              </h2>
               <div className="text-center text-sm text-gray-500 py-4">
                  No maintenance records found.
               </div>
            </div>
        </div>
      </div>

      {/* Edit Asset Modal */}
      <AssetAddEditModal
        show={showEditModal}
        onClose={() => setShowEditModal(false)}
        onSave={() => handleSaveEdit(editFormData)}
        asset={asset}
        formData={editFormData}
        setFormData={setEditFormData}
        categories={categories}
        departments={departments}
        openMapPicker={() => setShowMapPicker(true)}
      />

      {/* Map Picker Modal */}
      {showMapPicker && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60] p-4">
          <div className="bg-white rounded-lg p-6 w-full h-full max-w-4xl max-h-[90vh] flex flex-col">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Select Asset Location</h3>
              <button onClick={() => setShowMapPicker(false)} className="text-gray-500 hover:text-gray-700">
                <X className="w-6 h-6" />
              </button>
            </div>
            <MapPicker
              initialLocation={{ 
                lat: editFormData.latitude || asset.latitude || -6.2088, 
                lng: editFormData.longitude || asset.longitude || 106.8456 
              }}
              onLocationSelect={(details) => {
                setEditFormData(prev => ({ 
                  ...prev, 
                  latitude: details.lat, 
                  longitude: details.lng,
                  address: details.address,
                  building_room: details.building,
                }));
                setShowMapPicker(false);
              }}
              onClose={() => setShowMapPicker(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
}
