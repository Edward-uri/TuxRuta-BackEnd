from fastapi import FastAPI
from app.api.metrics import router as metrics_router

app = FastAPI()
app.include_router(metrics_router, prefix="/api")