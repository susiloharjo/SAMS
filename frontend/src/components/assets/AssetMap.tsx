'use client';

import { useEffect, useRef, useState } from 'react';
import { Loader } from '@googlemaps/js-api-loader';
import { MapPin } from 'lucide-react';

interface AssetMapProps {
  latitude: number;
  longitude: number;
  assetName: string;
}

export default function AssetMap({ latitude, longitude, assetName }: AssetMapProps) {
  const mapRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState<'loading' | 'error' | 'ready'>('loading');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loader = new Loader({
      apiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || "",
      version: "weekly",
      libraries: ["marker"], // Add marker library
    });

    loader.load().then((google) => {
      if (mapRef.current) {
        const map = new google.maps.Map(mapRef.current, {
          center: { lat: latitude, lng: longitude },
          zoom: 15,
          mapId: 'SAMS_ASSET_MAP',
          disableDefaultUI: true,
          zoomControl: true,
        });

        // --- ADDING MARKER ---
        new google.maps.Marker({
          position: { lat: latitude, lng: longitude },
          map: map,
          title: assetName,
        });
        
        setStatus('ready');
      }
    }).catch(e => {
      console.error("Failed to load Google Maps", e);
      setError("Map script failed to load. Please check the API key and network connection.");
      setStatus('error');
    });
  }, [latitude, longitude, assetName]);

  return (
    <div className="relative w-full h-80 bg-gray-200 rounded-lg overflow-hidden border">
      {status === 'loading' && (
        <div className="absolute inset-0 flex flex-col items-center justify-center bg-gray-100">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600 mb-2"></div>
          <p className="text-sm text-gray-500">Loading map...</p>
        </div>
      )}
      {status === 'error' && (
        <div className="absolute inset-0 flex flex-col items-center justify-center bg-red-50 p-4">
          <MapPin className="w-8 h-8 text-red-500 mb-2" />
          <p className="text-sm text-red-700 font-semibold">Could not load map</p>
          <p className="text-xs text-red-600 text-center mt-1">{error}</p>
        </div>
      )}
      <div ref={mapRef} className={`w-full h-full transition-opacity duration-300 ${status === 'ready' ? 'opacity-100' : 'opacity-0'}`} />
    </div>
  );
}
