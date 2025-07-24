from fastapi import FastAPI, WebSocket
from app.api.metrics import router as metrics_router
from app.api.websocket import websocket_endpoint

app = FastAPI()
app.include_router(metrics_router, prefix="/api")
app.add_api_websocket_route("/ws", websocket_endpoint)