import psycopg2
from datetime import datetime, date
from app.core.database import get_connection

def run_daily_etl(fecha: date = None):
    if fecha is None:
        fecha = datetime.now().date()

    print(f"🔎 Iniciando ETL para la fecha: {fecha}")

    try:
        conn = get_connection()
        cur = conn.cursor()

        # 1. Obtener todas las rutas activas
        cur.execute("SELECT id FROM ruta WHERE activa = true")
        rutas = [row[0] for row in cur.fetchall()]
        print(f"Rutas activas encontradas: {rutas}")

        if not rutas:
            print("⚠️ No hay rutas activas para procesar.")
            cur.close()
            conn.close()
            return

        for ruta_id in rutas:
            print(f"Procesando ruta {ruta_id}...")

            # 2. Obtener datos crudos de GPS y pasajeros para la ruta y fecha
            cur.execute("""
                SELECT latitude, longitude, speed_kmh, acceleration_ms2, turn_rate_dps, vehicle_state
                FROM raw_gps_data
                WHERE ruta_id = %s AND to_timestamp(timestamp) :: date = %s
            """, (ruta_id, fecha))
            gps_rows = cur.fetchall()

            cur.execute("""
                SELECT event, nn_passenger_detected, confidence, passenger_count_delta,
                       passenger_count_total, passenger_count_current
                FROM raw_passenger_data
                WHERE ruta_id = %s AND to_timestamp(timestamp) :: date = %s
            """, (ruta_id, fecha))
            passenger_rows = cur.fetchall()

            print(f"  Datos GPS encontrados: {len(gps_rows)}")
            print(f"  Datos pasajeros encontrados: {len(passenger_rows)}")

            # 3. Calcular métricas (ejemplo básico, puedes mejorar)
            velocidad_promedio = (
                sum(row[2] for row in gps_rows) / len(gps_rows)
                if gps_rows else 0
            )
            pasajeros_total = (
                sum(row[3] for row in passenger_rows if row[0] == "ENTRY" and row[1])
                if passenger_rows else 0
            )
            ocupacion_maxima = (
                max(row[5] for row in passenger_rows) if passenger_rows else 0
            )
            total_viajes = 1 if gps_rows else 0  # Mejora: detectar viajes reales
            hora_pico = "00:00-01:00"  # Mejora: calcular hora pico real
            probabilidad_ocupacion_alta = 0  # Mejora: calcula % de tiempo con ocupación alta
            intervalo_confianza_velocidad_min = velocidad_promedio  # Mejora: calcula real
            intervalo_confianza_velocidad_max = velocidad_promedio  # Mejora: calcula real

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
                fecha, ruta_id, pasajeros_total,
                pasajeros_total / total_viajes if total_viajes else 0,
                velocidad_promedio, hora_pico, total_viajes,
                ocupacion_maxima, probabilidad_ocupacion_alta,
                intervalo_confianza_velocidad_min, intervalo_confianza_velocidad_max
            ))

            print(f"  Resumen diario insertado/actualizado para ruta {ruta_id}")

        conn.commit()
        cur.close()
        conn.close()
        print(f"✅ ETL diario completado para {fecha}")

    except Exception as e:
        print(f"❌ Error en ETL: {e}")