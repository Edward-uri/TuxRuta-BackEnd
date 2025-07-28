import time
import pika
import json
import asyncio
from app.models.raw_data import GPSRawData
from pydantic import ValidationError
from app.core.database import get_connection
from app.api.websocket import manager

def publish_ws_update(data):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='13.219.25.247',
            credentials=pika.PlainCredentials('guest', 'guest')
        )
    )
    channel = connection.channel()
    channel.queue_declare(queue='ws_passenger_updates', durable=True)
    channel.basic_publish(
        exchange='',
        routing_key='ws_passenger_updates',
        body=json.dumps(data)
    )
    connection.close()

def process_gps_message(data: dict):
    try:
        # Asignar ruta_id por defecto si falta
        data['timestamp'] = int(time.time())
        if 'ruta_id' not in data:
            data['ruta_id'] = 1

        # Si falta 'data', crea el diccionario vacío
        if 'data' not in data:
            data['data'] = {}

        # Convertir timestamp a int si es floatya
        if 'timestamp' in data:
            data['timestamp'] = int(data['timestamp'])

        # Si speed_kmh es None, ponlo en 0.0
        if data['data'].get('speed_kmh') is None:
            data['data']['speed_kmh'] = 0.0

        for key in ['latitude', 'longitude', 'acceleration_ms2', 'turn_rate_dps']:
            if data['data'].get(key) is None:
                data['data'][key] = 0.0

        if data['data'].get('vehicle_state') is None:
            data['data']['vehicle_state'] = "DESCONOCIDO"

        gps = GPSRawData(**data)
        print("✅ GPS válido:", gps)
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

        # Publicar en la cola para WebSocket
        publish_ws_update({
            "type": "gps",
            "speed_kmh": gps.data.speed_kmh,
            "latitude": gps.data.latitude,
            "longitude": gps.data.longitude
        })

        # (Opcional) Broadcast local si tienes clientes WebSocket en este proceso
        try:
            loop = asyncio.get_running_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
        asyncio.run_coroutine_threadsafe(
            manager.broadcast({
                "type": "gps",
                "speed_kmh": gps.data.speed_kmh,
                "latitude": gps.data.latitude,
                "longitude": gps.data.longitude
            }),
            loop
        )
    except ValidationError as ve:
        print("❌ Error de validación GPS:", ve)
    except Exception as e:
        print("❌ Error guardando GPS en BD:", e)

def start_gps_consumer():
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='13.219.25.247',
            credentials=pika.PlainCredentials('guest', 'guest')
        )
    )
    channel = connection.channel()
    channel.queue_declare(queue='hybrid_49269307234447', durable=True)

    def callback_gps(ch, method, properties, body):
        try:
            data = json.loads(body)
            process_gps_message(data)
            print("✅ Mensaje GPS recibido:", data)
        except Exception as e:
            print("❌ Error procesando mensaje GPS:", e)

    channel.basic_consume(queue='hybrid_49269307234447', on_message_callback=callback_gps, auto_ack=True)

    print(" [*] Esperando mensajes en hybrid_49269307234447. Para salir presiona CTRL+C")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print(" [*] Consumidor detenido manualmente.")