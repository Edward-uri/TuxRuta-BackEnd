from pydantic import BaseModel, Field
from typing import Optional

class GPSData(BaseModel):
    latitude: float
    longitude: float
    speed_kmh: float
    acceleration_ms2: float
    turn_rate_dps: float
    vehicle_state: str

class GPSRawData(BaseModel):
    timestamp: int
    device_id: str
    ruta_id: int
    sensor_type: str = Field(default="HYBRID_GPS_MPU")
    data: GPSData

class PassengerData(BaseModel):
    event: str
    sensor_distance_mm: int
    nn_passenger_detected: bool
    confidence: float
    passenger_count_delta: int
    passenger_count_total: int
    passenger_count_current: int

class PassengerRawData(BaseModel):
    timestamp: int
    device_id: str
    ruta_id: int
    sensor_type: str = Field(default="PASSENGER_COUNT")
    data: PassengerData