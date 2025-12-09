from app.consumers.rabbit_consumer import start_gps_consumer
from app.consumers.passager_consumer import start_passenger_consumer
import threading
import time

def run_workers(target, count, name_prefix):
    threads = []
    for i in range(count):
        t = threading.Thread(target=target, daemon=True, name=f"{name_prefix}-{i+1}")
        t.start()
        threads.append(t)
        print(f" Started {name_prefix}-{i+1}")
    return threads

if __name__ == "__main__":
    print("Starting Analytics Consumers...")
    
    GPS_WORKERS = 20
    PASSENGER_WORKERS = 10

    gps_threads = run_workers(start_gps_consumer, GPS_WORKERS, "GPS-Worker")
    passenger_threads = run_workers(start_passenger_consumer, PASSENGER_WORKERS, "Passenger-Worker")

    print(f"Total active threads: {threading.active_count()}")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("Stopping consumers...")