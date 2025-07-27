import psycopg2
from datetime import datetime, date
from collections import Counter
from app.core.database import get_connection

def run_daily_etl(fecha: date = None, ruta_id: int = None):
    if fecha is None:
        fecha = datetime.now().date()

    print(f"🔎 Iniciando ETL para la fecha: {fecha}")

    try:
        conn = get_connection()
        cur = conn.cursor()

        # 1. Obtener rutas a procesar
        if ruta_id is not None:
            rutas = [ruta_id]
        else:
            cur.execute("SELECT id FROM ruta WHERE activa = true")
            rutas = [row[0] for row in cur.fetchall()]
        print(f"Rutas a procesar: {rutas}")

        if not rutas:
            print("⚠️ No hay rutas activas para procesar.")
            cur.close()
            conn.close()
            return

        for rid in rutas:
            print(f"Procesando ruta {rid}...")

            # 2. Obtener datos crudos de GPS y pasajeros para la ruta y fecha
            cur.execute("""
                SELECT latitude, longitude, speed_kmh, acceleration_ms2, turn_rate_dps, vehicle_state
                FROM raw_gps_data
                WHERE ruta_id = %s AND to_timestamp(timestamp) :: date = %s
            """, (rid, fecha))
            gps_rows = cur.fetchall()

            cur.execute("""
                SELECT event, nn_passenger_detected, confidence, passenger_count_delta,
                       passenger_count_total, passenger_count_current, timestamp
                FROM raw_passenger_data
                WHERE ruta_id = %s AND to_timestamp(timestamp) :: date = %s
            """, (rid, fecha))
            passenger_rows = cur.fetchall()

            print(f"  Datos GPS encontrados: {len(gps_rows)}")
            print(f"  Datos pasajeros encontrados: {len(passenger_rows)}")

            # 3. Calcular métricas

            # Velocidad promedio
            velocidad_promedio = (
                sum(row[2] for row in gps_rows) / len(gps_rows)
                if gps_rows else 0
            )

            # Pasajeros total (solo entradas válidas)
            pasajeros_total = (
                sum(row[3] for row in passenger_rows if row[0] == "ENTRY" and row[1])
                if passenger_rows else 0
            )

            # Ocupación máxima
            ocupacion_maxima = (
                max(row[5] for row in passenger_rows) if passenger_rows else 0
            )

            # Total de viajes (puedes mejorar la lógica, aquí es 1 si hay datos GPS)
            total_viajes = 1 if gps_rows else 0

            # Pasajeros promedio por viaje
            pasajeros_promedio_por_viaje = pasajeros_total / total_viajes if total_viajes else 0

            # Hora pico (franja horaria con más entradas de pasajeros)
            if passenger_rows:
                horas = [
                    datetime.fromtimestamp(row[6]).hour
                    for row in passenger_rows if row[0] == "ENTRY" and row[1]
                ]
                if horas:
                    hora_mas_frecuente = Counter(horas).most_common(1)[0][0]
                    hora_pico = f"{hora_mas_frecuente:02d}:00-{(hora_mas_frecuente+1)%24:02d}:00"
                else:
                    hora_pico = "Sin datos"
            else:
                hora_pico = "Sin datos"

            # Probabilidad de ocupación alta (porcentaje de registros con ocupación >= 80% de la máxima del día)
            if passenger_rows and ocupacion_maxima > 0:
                registros_alta = [row for row in passenger_rows if row[5] >= 0.8 * ocupacion_maxima]
                probabilidad_ocupacion_alta = len(registros_alta) / len(passenger_rows)
            else:
                probabilidad_ocupacion_alta = 0

            # Intervalo de confianza de velocidad (mínimo y máximo de las velocidades del día)
            if gps_rows:
                velocidades = [row[2] for row in gps_rows]
                intervalo_confianza_velocidad_min = min(velocidades)
                intervalo_confianza_velocidad_max = max(velocidades)
            else:
                intervalo_confianza_velocidad_min = 0
                intervalo_confianza_velocidad_max = 0

            # 4. Insertar o actualizar resumen diario
            cur.execute("""
                INSERT INTO resumen_diario_ruta (
                    fecha, ruta_id, pasajeros_total, pasajeros_promedio_por_viaje,
                    velocidad_promedio, hora_pico, total_viajes, ocupacion_maxima,
                    probabilidad_ocupacion_alta, intervalo_confianza_velocidad_min,
                    intervalo_confianza_velocidad_max
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (fecha, ruta_id) DO UPDATE SET
                    pasajeros_total = EXCLUDED.pasajeros_total,
                    pasajeros_promedio_por_viaje = EXCLUDED.pasajeros_promedio_por_viaje,
                    velocidad_promedio = EXCLUDED.velocidad_promedio,
                    hora_pico = EXCLUDED.hora_pico,
                    total_viajes = EXCLUDED.total_viajes,
                    ocupacion_maxima = EXCLUDED.ocupacion_maxima,
                    probabilidad_ocupacion_alta = EXCLUDED.probabilidad_ocupacion_alta,
                    intervalo_confianza_velocidad_min = EXCLUDED.intervalo_confianza_velocidad_min,
                    intervalo_confianza_velocidad_max = EXCLUDED.intervalo_confianza_velocidad_max
            """, (
                fecha, rid, pasajeros_total, pasajeros_promedio_por_viaje,
                velocidad_promedio, hora_pico, total_viajes, ocupacion_maxima,
                probabilidad_ocupacion_alta, intervalo_confianza_velocidad_min,
                intervalo_confianza_velocidad_max
            ))

            print(f"  Resumen diario insertado/actualizado para ruta {rid}")

        conn.commit()
        cur.close()
        conn.close()
        print(f"✅ ETL diario completado para {fecha}")

    except Exception as e:
        print(f"❌ Error durante el ETL: {e}")