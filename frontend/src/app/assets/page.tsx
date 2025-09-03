'use client'

import { useState, useEffect, useRef } from 'react'
import { formatIDR } from '../../utils/currency'
import QRCode from 'qrcode'
import { X } from 'lucide-react'
import { MapPicker } from '@/components/assets/MapPicker'
import { AssetsTable } from '@/components/assets/AssetsTable'
import { AssetAddEditModal } from '@/components/assets/AssetAddEditModal'
import { AssetDeleteModal } from '@/components/assets/AssetDeleteModal'
import { AssetQRModal } from '@/components/assets/AssetQRModal'
import { AssetControls } from '@/components/assets/AssetControls'
import { PaginationControls } from '@/components/assets/PaginationControls'

export interface Asset {
  id: string
  name: string
  description: string
  category_id: string
  department_id: string
  type: string
  model: string
  serial_number: string
  manufacturer: string
  acquisition_cost: number
  current_value: number
  status: string
  condition: string
  criticality: string
  latitude: number
  longitude: number
  address: string
  building_room: string
  acquisition_date: string
  expected_life_years: number
  category: {
    name: string
  }
  department: {
    name: string
  }
}

export interface Category {
  id: string
  name: string
  description: string
}

export interface Department {
  id: string
  name: string
  description: string
}

export default function AssetsPage() {
  const [assets, setAssets] = useState<Asset[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState(searchTerm);
  const [statusFilter, setStatusFilter] = useState('all')
  const [categoryFilter, setCategoryFilter] = useState('all')
  
  // CRUD Modal States
  const [showAddModal, setShowAddModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [showQRModal, setShowQRModal] = useState(false)
  const [selectedAsset, setSelectedAsset] = useState<Asset | null>(null)
  const [qrCodeDataUrl, setQrCodeDataUrl] = useState<string>('')
  const [selectedAssets, setSelectedAssets] = useState<Set<string>>(new Set())
  const [selectAll, setSelectAll] = useState(false)
  
  // Form States
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    category_id: '',
    department_id: '',
    type: '',
    model: '',
    serial_number: '',
    manufacturer: '',
    acquisition_cost: 0,
    current_value: 0,
    status: 'active',
    condition: 'good',
    criticality: 'medium',
    latitude: 0,
    longitude: 0,
    address: '',
    building_room: '',
    acquisition_date: '',
    expected_life_years: 5
  })

  // Pagination States
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [totalAssets, setTotalAssets] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [dataVersion, setDataVersion] = useState(0) // New state to trigger refetches

  // Map States
  const [showMap, setShowMap] = useState(false)
  const [mapCenter, setMapCenter] = useState({ lat: -6.2088, lng: 106.8456 }) // Jakarta coordinates
  const [selectedLocation, setSelectedLocation] = useState({ lat: 0, lng: 0 })

  // Debounce the search term to avoid excessive API calls
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedSearchTerm(searchTerm);
    }, 300); // Wait 300ms after the user stops typing

    // Cleanup function to cancel the timeout if the user types again
    return () => {
      clearTimeout(handler);
    };
  }, [searchTerm]);

  // Centralized effect for fetching data. This runs on page/limit changes or when manually triggered.
  useEffect(() => {
    fetchAssets()
  }, [currentPage, pageSize, dataVersion])

  // Effect for fetching categories, runs only once.
  useEffect(() => {
    fetchCategories()
  }, [])

  useEffect(() => {
    fetchRelatedData()
  }, [])

  // Clear selected assets when changing pages or filters
  useEffect(() => {
    setSelectedAssets(new Set())
    setSelectAll(false)
  }, [currentPage, pageSize, searchTerm, statusFilter, categoryFilter])

  // Reset to page 1 when any filter changes
  useEffect(() => {
    // We check currentPage !== 1 to avoid an infinite loop with the fetch effect
    if (currentPage !== 1) {
      setCurrentPage(1);
    } else {
      // If we are already on page 1, a filter change should still trigger a refetch
      setDataVersion(v => v + 1);
    }
  }, [debouncedSearchTerm, statusFilter, categoryFilter]);

  const fetchAssets = async () => {
    setLoading(true);
    try {
      // Build query params conditionally to avoid sending empty ones
      const params = new URLSearchParams({
        page: String(currentPage),
        limit: String(pageSize),
      });

      if (debouncedSearchTerm) {
        params.append('search', debouncedSearchTerm);
      }
      if (statusFilter !== 'all') {
        params.append('status', statusFilter);
      }
      if (categoryFilter !== 'all') {
        params.append('category', categoryFilter);
      }

      const url = `http://localhost:8080/api/v1/assets?${params.toString()}`;
      const response = await fetch(url)
      const data = await response.json()
      
      if (data.error === false) {
        setAssets(data.data || [])
        setTotalAssets(data.pagination.total)
        setTotalPages(data.pagination.total_pages)
      } else {
        console.error('Error fetching assets:', data.message)
      }
    } catch (error) {
      console.error('Error fetching assets:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchCategories = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/categories')
      const data = await response.json()
      setCategories(data.data || [])
    } catch (error) {
      console.error('Error fetching categories:', error)
    }
  }

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

  const handleAddAsset = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/assets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })
      
      if (response.ok) {
        setShowAddModal(false)
        resetForm()
        // Go to page 1 and trigger a refetch to see the new asset
        if (currentPage !== 1) setCurrentPage(1);
        setDataVersion(v => v + 1);
      }
    } catch (error) {
      console.error('Error adding asset:', error)
    }
  }

  const handleEditAsset = async () => {
    if (!selectedAsset) return
    
    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${selectedAsset.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })
      
      if (response.ok) {
        const responseData = await response.json();
        setShowEditModal(false)
        resetForm()
        // Go to page 1 and trigger a refetch to see the updated asset's new position
        if (currentPage !== 1) setCurrentPage(1);
        setDataVersion(v => v + 1);
      }
    } catch (error) {
      console.error('Error updating asset:', error)
    }
  }

  const handleDeleteAsset = async () => {
    if (!selectedAsset) return
    
    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${selectedAsset.id}`, {
        method: 'DELETE',
      })
      
      if (response.ok) {
        setShowDeleteModal(false)
        setSelectedAsset(null)
        // Trigger a refetch on the current page. If it becomes empty, pagination logic should handle it.
        setDataVersion(v => v + 1);
      }
    } catch (error) {
      console.error('Error deleting asset:', error)
    }
  }

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      category_id: '',
      department_id: '',
      type: '',
      model: '',
      serial_number: '',
      manufacturer: '',
      acquisition_cost: 0,
      current_value: 0,
      status: 'active',
      condition: 'good',
      criticality: 'medium',
      latitude: 0,
      longitude: 0,
      address: '',
      building_room: '',
      acquisition_date: '',
      expected_life_years: 5
    })
  }

  const openEditModal = (asset: Asset) => {
    setSelectedAsset(asset)
    setFormData({
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
    })
    setShowEditModal(true)
  }

  const openDeleteModal = (asset: Asset) => {
    setSelectedAsset(asset)
    setShowDeleteModal(true)
  }

  const openQRModal = async (asset: Asset) => {
    setSelectedAsset(asset)
    setShowQRModal(true)
    const qrDataUrl = await generateQRCode(asset)
    setQrCodeDataUrl(qrDataUrl)
  }

  const generateQRCode = async (asset: Asset) => {
    try {
      // Create QR code content with asset information
      const qrContent = JSON.stringify({
        assetId: asset.id,
        name: asset.name,
        serialNumber: asset.serial_number,
        category: asset.category?.name || 'Uncategorized',
        status: asset.status
      })
      
      const dataUrl = await QRCode.toDataURL(qrContent, {
        width: 200,
        margin: 2,
        color: {
          dark: '#000000',
          light: '#FFFFFF'
        }
      })
      return dataUrl
    } catch (error) {
      console.error('Error generating QR code:', error)
      return ''
    }
  }

  const handleAssetSelection = (assetId: string) => {
    const newSelected = new Set(selectedAssets)
    if (newSelected.has(assetId)) {
      newSelected.delete(assetId)
    } else {
      newSelected.add(assetId)
    }
    setSelectedAssets(newSelected)
  }

  const handleSelectAll = () => {
    if (selectAll) {
      setSelectedAssets(new Set())
      setSelectAll(false)
    } else {
      // Select only assets on the current page
      setSelectedAssets(new Set(assets.map(asset => asset.id)))
      setSelectAll(true)
    }
  }

  const generateBulkQRCodes = async () => {
    if (selectedAssets.size === 0) return
    
    try {
      const qrCodes = []
      for (const assetId of selectedAssets) {
        const asset = assets.find(a => a.id === assetId)
        if (asset) {
          const qrContent = JSON.stringify({
            assetId: asset.id,
            name: asset.name,
            serialNumber: asset.serial_number,
            category: asset.category?.name || 'Uncategorized',
            status: asset.status
          })
          
          const dataUrl = await QRCode.toDataURL(qrContent, {
            width: 120,
            margin: 1,
            color: {
              dark: '#000000',
              light: '#FFFFFF'
            }
          })
          qrCodes.push({ asset, qrCode: dataUrl })
        }
      }
      
      // Generate PDF with QR codes
      generatePDF(qrCodes)
    } catch (error) {
      console.error('Error generating bulk QR codes:', error)
    }
  }

  const generatePDF = (qrCodes: { asset: Asset, qrCode: string }[]) => {
    // Create a new window with the QR codes for printing
    const printWindow = window.open('', '_blank')
    if (!printWindow) return

    const html = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>QR Codes</title>
          <style>
            @media print {
              @page {
                size: A4;
                margin: 1cm;
              }
            }
            body {
              font-family: Arial, sans-serif;
              margin: 0;
              padding: 20px;
              background: #f5f5f5;
            }
            .qr-grid {
              display: grid;
              grid-template-columns: repeat(4, 1fr);
              gap: 15px;
              margin-bottom: 30px;
            }
            .qr-item {
              background: white;
              border: 1px solid #e5e7eb;
              border-radius: 8px;
              padding: 15px;
              box-shadow: 0 1px 3px rgba(0,0,0,0.1);
              display: flex;
              flex-direction: column;
              align-items: center;
              text-align: center;
              min-height: 200px;
            }
            .qr-code {
              margin-bottom: 10px;
              border: 1px solid #f3f4f6;
              padding: 8px;
              border-radius: 4px;
              background: white;
            }
            .qr-code img {
              width: 80px;
              height: 80px;
              display: block;
            }
            .asset-info {
              flex: 1;
              display: flex;
              flex-direction: column;
              justify-content: space-between;
              width: 100%;
            }
            .asset-name {
              font-weight: bold;
              margin-bottom: 5px;
              color: #1f2937;
              font-size: 14px;
            }
            .asset-details {
              color: #6b7280;
              font-size: 11px;
              line-height: 1.3;
            }
            .asset-id {
              font-family: monospace;
              font-weight: bold;
              color: #059669;
              font-size: 12px;
              margin-top: 8px;
            }
            .company-info {
              text-align: center;
              margin-top: 8px;
              font-size: 10px;
              color: #374151;
            }
            .company-name {
              font-weight: bold;
              color: #1f2937;
            }
            .page-break {
              page-break-before: always;
            }
            @media print {
              .no-print {
                display: none;
              }
              body {
                background: white;
              }
              .qr-item {
                box-shadow: none;
                border: 1px solid #d1d5db;
              }
            }
          </style>
        </head>
        <body>

          <div class="qr-grid">
            ${qrCodes.map((item, index) => `
              <div class="qr-item">
                <div class="qr-code">
                  <img src="${item.qrCode}" alt="QR Code for ${item.asset.name}" />
                </div>
                <div class="asset-info">
                  <div>
                    <div class="asset-name">${item.asset.name}</div>
                    <div class="asset-details">
                      SN: ${item.asset.serial_number || 'N/A'}<br>
                      Category: ${item.asset.category?.name || 'Uncategorized'}<br>
                      Status: ${item.asset.status}
                    </div>
                  </div>
                  <div class="asset-id">${item.asset.id.substring(0, 8).toUpperCase()}</div>
                  <div class="company-info">
                    <div class="company-name">SAMS</div>
                    <div>Smart Asset Management</div>
                  </div>
                </div>
              </div>
              ${(index + 1) % 16 === 0 && index < qrCodes.length - 1 ? '<div class="page-break"></div>' : ''}
            `).join('')}
          </div>
          
          <div class="no-print" style="text-align: center; margin-top: 30px; padding: 20px; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
            <button onclick="window.print()" style="padding: 12px 24px; font-size: 16px; background: #10b981; color: white; border: none; border-radius: 6px; cursor: pointer; margin-right: 10px;">
              🖨️ Print QR Codes
            </button>
            <button onclick="window.close()" style="padding: 12px 24px; font-size: 16px; background: #6b7280; color: white; border: none; border-radius: 6px; cursor: pointer;">
              ✕ Close
            </button>
          </div>
        </body>
      </html>
    `
    
    printWindow.document.write(html)
    printWindow.document.close()
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-green-600"></div>
      </div>
    )
  }

  const downloadQRCode = () => {
    if (!qrCodeDataUrl || !selectedAsset) return

    // Create a new window with the single QR code in bulk print format
    const printWindow = window.open('', '_blank')
    if (!printWindow) return

    const printContent = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>QR Code - ${selectedAsset.serial_number}</title>
          <style>
            @media print {
              @page {
                size: A4;
                margin: 1cm;
              }
            }
            body {
              font-family: Arial, sans-serif;
              margin: 0;
              padding: 20px;
              background: #f5f5f5;
            }
            .qr-container {
              background: white;
              border: 1px solid #e5e7eb;
              border-radius: 8px;
              padding: 20px;
              box-shadow: 0 1px 3px rgba(0,0,0,0.1);
              max-width: 400px;
              margin: 0 auto;
            }
            .qr-code {
              text-align: center;
              margin-bottom: 15px;
              border: 1px solid #f3f4f6;
              padding: 10px;
              border-radius: 4px;
              background: white;
            }
            .qr-code img {
              width: 120px;
              height: 120px;
              display: block;
              margin: 0 auto;
            }
            .asset-info {
              text-align: center;
              margin-bottom: 15px;
            }
            .asset-name {
              font-weight: bold;
              margin-bottom: 5px;
              color: #1f2937;
              font-size: 16px;
            }
            .asset-details {
              color: #6b7280;
              font-size: 12px;
              line-height: 1.3;
            }
            .asset-id {
              font-family: monospace;
              font-weight: bold;
              color: #059669;
              font-size: 12px;
              margin-top: 8px;
            }
            .company-info {
              text-align: center;
              margin-top: 15px;
              font-size: 10px;
              color: #374151;
            }
            .company-name {
              font-weight: bold;
              color: #1f2937;
            }
            .serial-number {
              text-align: center;
              margin-top: 10px;
              font-family: monospace;
              font-size: 14px;
              color: #6b7280;
            }
            .no-print {
              text-align: center;
              margin-top: 30px;
              padding: 20px;
              background: white;
              border-radius: 8px;
              box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
            .btn {
              padding: 12px 24px;
              font-size: 16px;
              border: none;
              border-radius: 6px;
              cursor: pointer;
              margin: 0 10px;
            }
            .btn-print {
              background: #10b981;
              color: white;
            }
            .btn-close {
              background: #6b7280;
              color: white;
            }
            @media print {
              body {
                background: white;
              }
              .qr-container {
                box-shadow: none;
                border: 1px solid #d1d5db;
              }
              .no-print {
                display: none;
              }
            }
          </style>
        </head>
        <body>
          <div class="qr-container">
            <div class="qr-code">
              <img src="${qrCodeDataUrl}" alt="QR Code for ${selectedAsset.name}" />
            </div>
            <div class="asset-info">
              <div class="asset-name">${selectedAsset.name}</div>
              <div class="asset-details">
                SN: ${selectedAsset.serial_number || 'N/A'}<br>
                Category: ${selectedAsset.category?.name || 'Uncategorized'}<br>
                Status: ${selectedAsset.status}
              </div>
              <div class="asset-id">${selectedAsset.id.substring(0, 8).toUpperCase()}</div>
            </div>
            <div class="company-info">
              <div class="company-name">SAMS</div>
              <div>Smart Asset Management</div>
            </div>
            <div class="serial-number">${selectedAsset.serial_number}</div>
          </div>
          
          <div class="no-print">
            <button onclick="window.print()" class="btn btn-print">
              🖨️ Print / Save as PDF
            </button>
            <button onclick="window.close()" class="btn btn-close">
              ✕ Close
            </button>
          </div>
        </body>
      </html>
    `

    printWindow.document.write(printContent)
    printWindow.document.close()
    printWindow.focus()
  }

  const printBulkQR = async () => {
    if (selectedAssets.size === 0) return

    // Get the actual asset objects from the selected IDs
    const selectedAssetObjects = assets.filter(asset => selectedAssets.has(asset.id))
    
    if (selectedAssetObjects.length === 0) return

    // Generate QR codes for all selected assets
    const qrCodes = await Promise.all(
      selectedAssetObjects.map(async (asset) => {
        const qrDataUrl = await generateQRCode(asset)
        return { asset, qrDataUrl }
      })
    )

    const printWindow = window.open('', '_blank')
    if (!printWindow) return

    const printContent = `
      <!DOCTYPE html>
      <html>
        <head>
          <title>QR Codes</title>
          <style>
            @media print {
              @page {
                size: A4;
                margin: 1cm;
              }
            }
            body {
              font-family: Arial, sans-serif;
              margin: 0;
              padding: 20px;
            }
            .qr-grid {
              display: grid;
              grid-template-columns: repeat(2, 1fr);
              gap: 20px;
              page-break-inside: avoid;
            }
            .qr-item {
              border: 1px solid #ddd;
              padding: 15px;
              text-align: center;
              page-break-inside: avoid;
            }
            .qr-code {
              width: 120px;
              height: 120px;
              margin: 0 auto 10px;
            }
            .asset-info {
              font-size: 12px;
              color: #333;
            }
            .company-info {
              font-size: 10px;
              color: #666;
              margin-top: 5px;
            }
            .page-break {
              page-break-before: always;
            }
          </style>
        </head>
        <body>
          <div class="qr-grid">
            ${qrCodes.map(({ asset, qrDataUrl }, index) => `
              <div class="qr-item">
                <div class="qr-code">
                  <img src="${qrDataUrl}" alt="QR Code" style="width: 100%; height: 100%;">
                </div>
                <div class="asset-info">
                  <strong>${asset.serial_number}</strong><br>
                  ${asset.name}
                </div>
                <div class="company-info">
                  PROPERTY OF SAMS Corporation<br>
                  Smart Asset Management System
                </div>
              </div>
              ${(index + 1) % 6 === 0 && index < qrCodes.length - 1 ? '<div class="page-break"></div>' : ''}
            `).join('')}
          </div>
        </body>
      </html>
    `

    printWindow.document.write(printContent)
    printWindow.document.close()
    printWindow.focus()
  }

  const handleSelectLocation = (lat: number, lng: number) => {
    setSelectedLocation({ lat, lng })
    setFormData(prev => ({
      ...prev,
      latitude: lat,
      longitude: lng
    }))
    setShowMap(false)
  }

  const openMapPicker = () => {
    if (formData.latitude && formData.longitude) {
      setMapCenter({ lat: formData.latitude, lng: formData.longitude })
      setSelectedLocation({ lat: formData.latitude, lng: formData.longitude })
    }
    setShowMap(true)
  }

  return (
    <div className="space-y-6">
      <AssetControls
        selectedAssetsSize={selectedAssets.size}
        onBulkPrint={generateBulkQRCodes}
        onAddNew={() => setShowAddModal(true)}
        searchTerm={searchTerm}
        onSearchChange={setSearchTerm}
        statusFilter={statusFilter}
        onStatusChange={setStatusFilter}
        categoryFilter={categoryFilter}
        onCategoryChange={setCategoryFilter}
        categories={categories}
      />

      {/* Assets Table */}
      <AssetsTable
        assets={assets}
        selectedAssets={selectedAssets}
        selectAll={selectAll}
        totalAssets={totalAssets}
        handleSelectAll={handleSelectAll}
        handleAssetSelection={handleAssetSelection}
        openQRModal={openQRModal}
        openEditModal={openEditModal}
        openDeleteModal={openDeleteModal}
      />

      {/* Pagination Controls */}
      <PaginationControls
        currentPage={currentPage}
        pageSize={pageSize}
        totalAssets={totalAssets}
        totalPages={totalPages}
        onPageChange={setCurrentPage}
        onPageSizeChange={setPageSize}
      />

      {/* Add/Edit Asset Modal */}
      <AssetAddEditModal
        show={showAddModal || showEditModal}
        onClose={() => {
          setShowAddModal(false);
          setShowEditModal(false);
          resetForm();
        }}
        onSave={showEditModal ? handleEditAsset : handleAddAsset}
        asset={selectedAsset}
        formData={formData}
        setFormData={setFormData}
        categories={categories}
        departments={departments}
        openMapPicker={openMapPicker}
      />

      {/* Delete Asset Modal */}
      <AssetDeleteModal
        show={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        onDelete={handleDeleteAsset}
        asset={selectedAsset}
      />

      {/* QR Code Modal */}
      <AssetQRModal
        show={showQRModal}
        onClose={() => setShowQRModal(false)}
        onDownload={downloadQRCode}
        asset={selectedAsset}
        qrCodeDataUrl={qrCodeDataUrl}
      />

      {/* Map Picker Modal */}
      {showMap && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full h-full max-w-4xl max-h-[90vh] flex flex-col">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Select Asset Location</h3>
              <button onClick={() => setShowMap(false)} className="text-gray-500 hover:text-gray-700">
                <X className="w-6 h-6" />
              </button>
            </div>
            <MapPicker
              initialLocation={selectedLocation}
              onLocationSelect={(details) => {
                setSelectedLocation({ lat: details.lat, lng: details.lng });
                setFormData(prev => ({ 
                  ...prev, 
                  latitude: details.lat, 
                  longitude: details.lng,
                  address: details.address,
                  building_room: details.building,
                }));
                setShowMap(false);
              }}
              onClose={() => setShowMap(false)}
            />
          </div>
        </div>
      )}
    </div>
  )
}
