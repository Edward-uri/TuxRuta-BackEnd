from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.metrics import router as metrics_router
from app.api.websocket import websocket_endpoint

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