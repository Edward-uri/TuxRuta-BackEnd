from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.metrics import router as metrics_router
from app.api.websocket import websocket_endpoint

app = FastAPI()
app.port = 9000
# Configura los orígenes permitidos
origins = [
    "http://localhost:4200",
    "https://9jw825mb-9000.usw3.devtunnels.ms",
    # agrega aquí otros orígenes si los necesitas
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(metrics_router, prefix="/api")
app.add_api_websocket_route("/ws", websocket_endpoint)