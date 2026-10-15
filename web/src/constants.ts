export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8084';
export const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL ?? 'ws://localhost:8084/ws';

// Default map view — Addis Ababa, Ethiopia
// (approx. 9°01′N 38°45′E, around the city center / Meskel Square)
export const MAP_CENTER: [number, number] = [9.0192, 38.7525];
// Detailed Esri street basemap, served from the same public ArcGIS endpoint the
// app already used and requiring no API key. Unlike the previous plain gray
// base layer, it renders street names, place labels and points of interest.
// (CARTO basemaps now stamp an "API key required" watermark on keyless tiles,
// so they are no longer used here.)
export const MAP_TILE_URL =
  'https://server.arcgisonline.com/ArcGIS/rest/services/World_Street_Map/MapServer/tile/{z}/{y}/{x}';
export const MAP_TILE_ATTRIBUTION =
  "Tiles &copy; <a href='https://www.esri.com/'>Esri</a>, HERE, Garmin, FAO, NOAA, USGS, OpenStreetMap contributors";
export const MAP_MAX_ZOOM = 19;
