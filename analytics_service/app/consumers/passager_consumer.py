import time
import pika
import json
import asyncio
from app.models.raw_data import PassengerRawData
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
        print(f" Error publicando en WS queue: {e}")

def process_passenger_message(channel, data: dict):
    try:
    
        
        if 'data' not in data:
            nested_data = {
                "event": data.get("event_type", "UNKNOWN"),
                "sensor_distance_mm": data.get("sensor_distance_mm", 0),
                "nn_passenger_detected": data.get("nn_passenger_detected", False),
                "confidence": data.get("confidence", 0.0),
                "passenger_count_delta": data.get("passenger_delta", 0), # Mapeo de nombre
                "passenger_count_total": data.get("total_entries", 0), # Asumiendo total_entries es el total
                "passenger_count_current": data.get("current_count", 0) # Mapeo de nombre
            }
            data['data'] = nested_data

        # Asignar timestamp si falta
        if 'timestamp' not in data:
            data['timestamp'] = int(time.time())
        else:
            data['timestamp'] = int(data['timestamp'])

        if 'ruta_id' not in data:
            data['ruta_id'] = 2  # Ruta "Suchi" default

        passenger = PassengerRawData(**data)
        print(f" Pasajero válido: {passenger.device_id} (Count: {passenger.data.passenger_count_current})")

        # Usar el pool de conexiones
        try:
            with get_db_connection() as conn:
                with conn.cursor() as cur:
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
        except Exception as e:
            print(f" Error DB Pasajero: {e}")

        publish_ws_update(channel, {
            "type": "passenger",
            "passenger_count_current": passenger.data.passenger_count_current,
            "device_id": passenger.device_id,
            "ruta_id": passenger.ruta_id
        })

        try:
            loop = asyncio.get_running_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            
        asyncio.run_coroutine_threadsafe(
            manager.broadcast({
                "type": "passenger",
                "passenger_count_current": passenger.data.passenger_count_current,
                "device_id": passenger.device_id
            }),
            loop
        )

    except ValidationError as ve:
        print(f" Error de validación pasajero: {ve}")
    except Exception as e:
        print(f" Error procesando pasajero: {e}")

def start_passenger_consumer():
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
            
            channel.basic_qos(prefetch_count=20)
            
            queue_name = 'passenger_49269307234447'
            channel.queue_declare(queue=queue_name, durable=True)

            def callback_passenger(ch, method, properties, body):
                try:
                    data = json.loads(body)
                    process_passenger_message(ch, data)
                    ch.basic_ack(delivery_tag=method.delivery_tag)
                except Exception as e:
                    print(f"Error en callback pasajero: {e}")

            channel.basic_consume(queue=queue_name, on_message_callback=callback_passenger, auto_ack=False)

            print(f" Esperando mensajes en {queue_name}. Para salir presiona CTRL+C")
            channel.start_consuming()

        except pika.exceptions.AMQPConnectionError as e:
            print(f" Error de conexión RabbitMQ (Pasajero): {e}. Reintentando en 5s...")
            time.sleep(5)
        except Exception as e:
            print(f" Error inesperado (Pasajero): {e}. Reintentando en 5s...")
            time.sleep(5)