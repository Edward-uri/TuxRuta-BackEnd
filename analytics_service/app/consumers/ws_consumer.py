import threading
import pika
import asyncio
import json
from app.api.websocket import manager 

def ws_rabbitmq_consumer(loop):
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host='13.219.25.247',
            credentials=pika.PlainCredentials('guest', 'guest')
        )
    )
    channel = connection.channel()
    channel.queue_declare(queue='ws_passenger_updates', durable=True)
    
    def callback(ch, method, properties, body):
        data = json.loads(body)
        print("🟢 Recibido de RabbitMQ para WS:", data)
        print("Conexiones activas antes de broadcast:", len(manager.active_connections))
        asyncio.run_coroutine_threadsafe(manager.broadcast(data), loop)

    channel.basic_consume(queue='ws_passenger_updates', on_message_callback=callback, auto_ack=True)
    print(" [*] WebSocket suscrito a ws_passenger_updates")
    channel.start_consuming()

def start_ws_rabbitmq_consumer(loop):
    t = threading.Thread(target=ws_rabbitmq_consumer, args=(loop,), daemon=True)
    t.start()