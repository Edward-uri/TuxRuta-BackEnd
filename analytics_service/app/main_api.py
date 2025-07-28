from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.metrics import router as metrics_router
from app.api.websocket import websocket_endpoint
import asyncio
main_loop = asyncio.get_event_loop()
from app.consumers.ws_consumer import start_ws_rabbitmq_consumer
app = FastAPI()


app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(metrics_router, prefix="/api")
app.add_api_websocket_route("/ws", websocket_endpoint)
start_ws_rabbitmq_consumer(main_loop)