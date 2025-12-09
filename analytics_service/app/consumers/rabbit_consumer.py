import time
import pika
import json
import asyncio
from app.models.raw_data import GPSRawData
from pydantic import ValidationError
from app.core.database import get_db_connection
from app.api.websocket import manager

def publish_ws_update(channel, data):
    try:
        channel.queue_declare(queue='ws_passenger_updates', durable=True)
        channel.basic_publish(
            exchange='',
            routing_key='ws_passenger_updates',
            body=json.dumps(data)
        )
    except Exception as e:
        print(f"Error publicando en WS queue: {e}")

def process_gps_message(channel, data: dict):
    try:
        # Validación básica de campos requeridos
        if 'device_id' not in data:
            print(" Mensaje sin device_id, ignorando.")
            return

        # Asignar timestamp si falta
        if 'timestamp' not in data:
            data['timestamp'] = int(time.time())
        else:
            data['timestamp'] = int(data['timestamp'])

  
        if 'ruta_id' not in data:
            
             data['ruta_id'] = 2 

        # Si falta 'data', crea el diccionario vacío
        if 'data' not in data:
            data['data'] = {}

        # Valores por defecto para campos numéricos
        for key in ['speed_kmh', 'latitude', 'longitude', 'acceleration_ms2', 'turn_rate_dps']:
            if data['data'].get(key) is None:
                data['data'][key] = 0.0

        if data['data'].get('vehicle_state') is None:
            data['data']['vehicle_state'] = "DESCONOCIDO"

        gps = GPSRawData(**data)
        print(f"GPS válido: {gps.device_id} (Ruta: {gps.ruta_id})")

        # Usar el pool de conexiones
        try:
            with get_db_connection() as conn:
                with conn.cursor() as cur:
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
        except Exception as e:
            print(f"Error DB: {e}")

        publish_ws_update(channel, {
            "type": "gps",
            "speed_kmh": gps.data.speed_kmh,
            "latitude": gps.data.latitude,
            "longitude": gps.data.longitude,
            "device_id": gps.device_id,
            "ruta_id": gps.ruta_id
        })

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
                "longitude": gps.data.longitude,
                "device_id": gps.device_id
            }),
            loop
        )

    except ValidationError as ve:
        print(f"Error de validación GPS: {ve}")
    except Exception as e:
        print(f"Error procesando GPS: {e}")

def start_gps_consumer():
    while True:
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(
                    host='34.233.205.241',
                    credentials=pika.PlainCredentials('guest', 'guest'),
                    heartbeat=600,
                    blocked_connection_timeout=300
                )
            )
            channel = connection.channel()
            
            channel.basic_qos(prefetch_count=50)
            
            queue_name = 'hybrid_49269307234447'
            channel.queue_declare(queue=queue_name, durable=True)

            def callback_gps(ch, method, properties, body):
                try:
                    data = json.loads(body)
                    process_gps_message(ch, data)
                    # Confirmación manual
                    ch.basic_ack(delivery_tag=method.delivery_tag)
                except Exception as e:
                    print(f" Error en callback: {e}")
              

            channel.basic_consume(queue=queue_name, on_message_callback=callback_gps, auto_ack=False)

            print(f" Esperando mensajes en {queue_name}. Para salir presiona CTRL+C")
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            print(f"Error de conexión RabbitMQ: {e}. Reintentando en 5s...")
            time.sleep(5)
        except Exception as e:
            print(f"Error inesperado: {e}. Reintentando en 5s...")
            time.sleep(5)