'use client';

import { useState, useEffect, useRef } from 'react';

// A new, self-contained, and working Map Picker Component
export const MapPicker = ({ onLocationSelect, onClose, initialLocation }: {
  onLocationSelect: (details: { lat: number, lng: number, address: string, building: string }) => void,
  onClose: () => void,
  initialLocation: { lat: number, lng: number } | null
}) => {
  const mapRef = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<google.maps.Map | null>(null);
  const [marker, setMarker] = useState<google.maps.marker.AdvancedMarkerElement | null>(null);
  const [selectedLocation, setSelectedLocation] = useState<{ lat: number, lng: number } | null>(initialLocation);
  const [addressDetails, setAddressDetails] = useState({ address: '', building: '' });
  const [searchQuery, setSearchQuery] = useState('');
  
  const geocodeLocation = (location: { lat: number, lng: number }) => {
    const geocoder = new google.maps.Geocoder();
    geocoder.geocode({ location }, (results, status) => {
      if (status === 'OK' && results && results[0]) {
        // Build address from components to avoid Plus Codes
        const addressComponents = results[0].address_components;
        const streetNumber = addressComponents.find(c => c.types.includes('street_number'))?.long_name || '';
        const route = addressComponents.find(c => c.types.includes('route'))?.long_name || '';
        const locality = addressComponents.find(c => c.types.includes('locality'))?.long_name || '';
        const administrativeArea = addressComponents.find(c => c.types.includes('administrative_area_level_1'))?.long_name || '';
        const postalCode = addressComponents.find(c => c.types.includes('postal_code'))?.long_name || '';
        const country = addressComponents.find(c => c.types.includes('country'))?.long_name || '';
        
        // Build clean address without Plus Code
        const streetAddress = [streetNumber, route].filter(Boolean).join(' ');
        const cityState = [locality, administrativeArea].filter(Boolean).join(', ');
        const address = [streetAddress, cityState, postalCode, country].filter(Boolean).join(', ');
        
        const premiseComponent = addressComponents.find(c => 
          c.types.includes('premise') || 
          c.types.includes('point_of_interest') || 
          c.types.includes('establishment')
        );
        const building = premiseComponent?.long_name || '';
        setAddressDetails({ address, building });
      } else {
        setAddressDetails({ address: 'Could not determine address', building: '' });
      }
    });
  };
  
  // Initialize map and geocode initial location
  useEffect(() => {
    if (mapRef.current && !map) {
      const initialCenter = initialLocation || { lat: -6.2088, lng: 106.8456 }; // Jakarta
      const newMap = new google.maps.Map(mapRef.current, {
        center: initialCenter,
        zoom: 12,
        mapId: 'SAMS_MAP_ID' // Required for Advanced Markers
      });
      setMap(newMap);
      if (initialLocation) {
        geocodeLocation(initialLocation);
      }
    }
  }, [mapRef, map, initialLocation]);

  // Handle map clicks
  useEffect(() => {
    if (!map) return;
    const clickListener = map.addListener('click', (e: google.maps.MapMouseEvent) => {
      if (e.latLng) {
        const newLocation = { lat: e.latLng.lat(), lng: e.latLng.lng() };
        setSelectedLocation(newLocation);
        geocodeLocation(newLocation);
      }
    });
    return () => { google.maps.event.removeListener(clickListener); }
  }, [map]);

  // Update marker
  useEffect(() => {
    if (map && selectedLocation) {
      if (marker) {
        marker.position = selectedLocation;
      } else {
        const newMarker = new google.maps.marker.AdvancedMarkerElement({
          position: selectedLocation,
          map,
          title: "Asset Location"
        });
        setMarker(newMarker);
      }
      map.panTo(selectedLocation);
    }
  }, [map, selectedLocation, marker]);

  const handleSearch = () => {
    if (!searchQuery.trim()) return;
    const geocoder = new google.maps.Geocoder();
    geocoder.geocode({ address: searchQuery }, (results, status) => {
      if (status === 'OK' && results && results[0]) {
        const location = results[0].geometry.location;
        const newLocation = { lat: location.lat(), lng: location.lng() };
        setSelectedLocation(newLocation);
        
        // Build address from components to avoid Plus Codes
        const addressComponents = results[0].address_components;
        const streetNumber = addressComponents.find(c => c.types.includes('street_number'))?.long_name || '';
        const route = addressComponents.find(c => c.types.includes('route'))?.long_name || '';
        const locality = addressComponents.find(c => c.types.includes('locality'))?.long_name || '';
        const administrativeArea = addressComponents.find(c => c.types.includes('administrative_area_level_1'))?.long_name || '';
        const postalCode = addressComponents.find(c => c.types.includes('postal_code'))?.long_name || '';
        const country = addressComponents.find(c => c.types.includes('country'))?.long_name || '';
        
        // Build clean address without Plus Code
        const streetAddress = [streetNumber, route].filter(Boolean).join(' ');
        const cityState = [locality, administrativeArea].filter(Boolean).join(', ');
        const address = [streetAddress, cityState, postalCode, country].filter(Boolean).join(', ');
        
        const premiseComponent = addressComponents.find(c => c.types.includes('premise') || c.types.includes('establishment'));
        const building = premiseComponent?.long_name || '';
        setAddressDetails({ address, building });
      } else {
        alert('Search failed: Location not found.');
      }
    });
  };

  const handleConfirm = () => {
    if (selectedLocation) {
      onLocationSelect({
        ...selectedLocation,
        ...addressDetails
      });
      onClose();
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex gap-2 mb-4">
        <input
          type="text"
          placeholder="Search for a location..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button onClick={handleSearch} className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
          Search
        </button>
      </div>
      <div ref={mapRef} className="w-full flex-grow border border-gray-300 rounded-md" />
      <div className="mt-4 p-2 bg-gray-100 rounded-md text-center">
          <p className="text-sm text-gray-800">
            Lat: {selectedLocation?.lat.toFixed(6) || 'N/A'}, Lng: {selectedLocation?.lng.toFixed(6) || 'N/A'}
          </p>
          {addressDetails.address && (
            <p className="text-xs text-gray-600 mt-1">
              {addressDetails.address}
            </p>
          )}
      </div>
      <div className="flex justify-end gap-2 mt-4">
        <button onClick={onClose} className="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50">
          Cancel
        </button>
        <button onClick={handleConfirm} disabled={!selectedLocation} className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400">
          Confirm Location
        </button>
      </div>
    </div>
  );
}
