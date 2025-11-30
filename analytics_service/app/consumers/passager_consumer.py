import time
import pika
import json
import asyncio
from app.models.raw_data import PassengerRawData
from pydantic import ValidationError
from app.core.database import get_connection
from app.api.websocket import manager

def publish_ws_update(data):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='34.233.205.241',
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

def process_passenger_message(data: dict):
    try:
        # Asignar ruta_id por defecto si falta
        data['timestamp'] = int(time.time())

        if 'ruta_id' not in data:
            data['ruta_id'] = 2  # Ruta "Suchi"
        passenger = PassengerRawData(**data)
        print("✅ Pasajero válido:", passenger)
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

        # Publicar en la cola para WebSocket
        publish_ws_update({
            "type": "passenger",
            "passenger_count_current": passenger.data.passenger_count_current
        })

        # (Opcional) Broadcast local si tienes clientes WebSocket en este proceso
        try:
            loop = asyncio.get_running_loop()
        except RuntimeError:
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
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

def start_passenger_consumer():
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='34.233.205.241',
            credentials=pika.PlainCredentials('guest', 'guest')
        )
    )
    channel = connection.channel()
    channel.queue_declare(queue='passenger_49269307234447', durable=True)

    def callback_passenger(ch, method, properties, body):
        try:
            data = json.loads(body)
            print("✅ Mensaje pasajero recibido:", data)
            process_passenger_message(data)
        except Exception as e:
            print("❌ Error procesando mensaje pasajero:", e)

    channel.basic_consume(queue='passenger_49269307234447', on_message_callback=callback_passenger, auto_ack=True)

    print(" [*] Esperando mensajes en passenger_49269307234447. Para salir presiona CTRL+C")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print(" [*] Consumidor detenido manualmente.")
        channel.stop_consuming()
        connection.close()