import { useEffect } from "react";
import { useMap, Polyline } from "react-leaflet";
import L from "leaflet";

export function RoutingControl({ route }: {
    route: [number, number][]
}) {
    const map = useMap()

    // Zoom and pan the map so the whole route is visible. Without this the view
    // stays on the default center, which hides longer routes off-screen and
    // buries short ones under the pickup/destination markers.
    useEffect(() => {
        if (!route || route.length < 2) return

        map.fitBounds(L.latLngBounds(route), { padding: [48, 48] })
    }, [map, route])

    if (!route) return null

    return <Polyline positions={route} color="blue" />
}