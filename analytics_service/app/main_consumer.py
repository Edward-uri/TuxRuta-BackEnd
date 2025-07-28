from app.consumers.rabbit_consumer import start_gps_consumer
from app.consumers.passager_consumer import start_passenger_consumer
import threading

if __name__ == "__main__":
    t1 = threading.Thread(target=start_gps_consumer, daemon=True)
    t2 = threading.Thread(target=start_passenger_consumer, daemon=True)
    t1.start()
    t2.start()
    t1.join()
    t2.join()