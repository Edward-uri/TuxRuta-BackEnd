from fastapi import WebSocket
import json

class WebSocketManager:
    def __init__(self):
        self.active_connections = []

    async def connect(self, websocket: WebSocket):
        await websocket.accept()
        self.active_connections.append(websocket)

    async def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)
        print("Conexiones activas:", len(self.active_connections))

    async def broadcast(self, data: dict):
        message = json.dumps(data)
        print("🔴 Enviando por WebSocket:", message)
        print("Conexiones activas en broadcast:", len(self.active_connections))  # <-- Agrega esto
        for connection in self.active_connections:
            try:
                await connection.send_text(message)
                print("✅ Mensaje enviado a un cliente WebSocket")
            except Exception as e:
                print("❌ Error enviando mensaje por WebSocket:", e)

manager = WebSocketManager()

async def websocket_endpoint(websocket: WebSocket):
    await manager.connect(websocket)
    try:
        while True:
            await websocket.receive_text()
    except Exception:
        pass
    finally:
        await manager.disconnect(websocket)