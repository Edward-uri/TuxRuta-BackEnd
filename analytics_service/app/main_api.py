from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.metrics import router as metrics_router
from app.api.websocket import websocket_endpoint
import asyncio
from app.consumers.ws_consumer import start_ws_rabbitmq_consumer
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: Iniciar consumidor
    loop = asyncio.get_running_loop()
    start_ws_rabbitmq_consumer(loop)
    yield
    # Shutdown: Aquí podrías cerrar conexiones si fuera necesario

app = FastAPI(lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=False,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(metrics_router, prefix="/api")
app.add_api_websocket_route("/ws", websocket_endpoint)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("app.main_api:app", host="0.0.0.0", port=8000, reload=True)