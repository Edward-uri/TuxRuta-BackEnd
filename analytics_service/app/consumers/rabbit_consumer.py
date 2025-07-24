import pika
import json
import asyncio
from app.models.raw_data import GPSRawData, PassengerRawData
from pydantic import ValidationError
from app.core.database import get_connection
from app.api.websocket import manager

def process_gps_message(data: dict):
    try:
        gps = GPSRawData(**data)
        print("✅ GPS válido:", gps)
        # Guardar en BD
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("""
            INSERT INTO raw_gps_data (
                timestamp, device_id, ruta_id, latitude, longitude,
                speed_kmh, acceleration_ms2, turn_rate_dps, vehicle_state
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            gps.timestamp, gps.device_id, gps.ruta_id,
            gps.data.latitude, gps.data.longitude, gps.data.speed_kmh,
            gps.data.acceleration_ms2, gps.data.turn_rate_dps, gps.data.vehicle_state
        ))
        conn.commit()
        cur.close()
        conn.close()

        # Enviar velocidad por WebSocket
        loop = asyncio.get_event_loop()
        asyncio.run_coroutine_threadsafe(
            manager.broadcast({
                "type": "gps",
                "speed_kmh": gps.data.speed_kmh
            }),
            loop
        )

    except ValidationError as ve:
        print("❌ Error de validación GPS:", ve)
    except Exception as e:
        print("❌ Error guardando GPS en BD:", e)

def process_passenger_message(data: dict):
    try:
        passenger = PassengerRawData(**data)
        print("✅ Pasajero válido:", passenger)
        # Guardar en BD
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("""
            INSERT INTO raw_passenger_data (
                timestamp, device_id, ruta_id, event, sensor_distance_mm,
                nn_passenger_detected, confidence, passenger_count_delta,
                passenger_count_total, passenger_count_current
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
        """, (
            passenger.timestamp, passenger.device_id, passenger.ruta_id,
            passenger.data.event, passenger.data.sensor_distance_mm,
            passenger.data.nn_passenger_detected, passenger.data.confidence,
            passenger.data.passenger_count_delta, passenger.data.passenger_count_total,
            passenger.data.passenger_count_current
        ))
        conn.commit()
        cur.close()
        conn.close()

        # Enviar conteo de pasajeros por WebSocket
        loop = asyncio.get_event_loop()
        asyncio.run_coroutine_threadsafe(
            manager.broadcast({
                "type": "passenger",
                "passenger_count_current": passenger.data.passenger_count_current
            }),
            loop
        )

    except ValidationError as ve:
        print("❌ Error de validación pasajero:", ve)
    except Exception as e:
        print("❌ Error guardando pasajero en BD:", e)