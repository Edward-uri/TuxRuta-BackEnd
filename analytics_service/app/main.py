import schedule
import time
import sys
from app.services.etl import run_daily_etl

def job():
    print("⏰ Ejecutando ETL diario...")
    run_daily_etl()

schedule.every().day.at("04:36").do(job)

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "run-now":
        print("🚀 Ejecutando ETL manualmente...")
        run_daily_etl()
    else:
        print("⏳ Esperando para ejecutar ETL diario cada día a las 04:33...")
        while True:
            schedule.run_pending()
            time.sleep(60)