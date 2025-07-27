import psycopg2
from datetime import datetime, date
from collections import defaultdict
from app.core.database import get_connection

def run_monthly_etl(year: int = None, month: int = None, ruta_id: int = None):
    if year is None or month is None:
        today = datetime.now().date()
        year = today.year
        month = today.month

    print(f"🔎 Iniciando ETL mensual para {year}-{month:02d}")

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

            # 2. Obtener datos diarios de ese mes y ruta
            cur.execute("""
                SELECT fecha, pasajeros_total, velocidad_promedio, ocupacion_maxima,
                       probabilidad_ocupacion_alta, intervalo_confianza_velocidad_min, intervalo_confianza_velocidad_max
                FROM resumen_diario_ruta
                WHERE ruta_id = %s AND EXTRACT(YEAR FROM fecha) = %s AND EXTRACT(MONTH FROM fecha) = %s
            """, (rid, year, month))
            rows = cur.fetchall()

            if not rows:
                print(f"  No hay datos diarios para ruta {rid} en este mes.")
                continue

            dias = len(rows)
            pasajeros_total_mes = sum(r[1] for r in rows)
            pasajeros_promedio_dia = pasajeros_total_mes / dias if dias else 0
            velocidad_promedio_mes = sum(r[2] for r in rows) / dias if dias else 0
            ocupacion_maxima_mes = max(r[3] for r in rows)
            dias_activos = dias

            # Mejor día de la semana (por promedio de pasajeros)
            dias_semana = defaultdict(list)
            for r in rows:
                dia_semana = r[0].isoweekday()
                dias_semana[dia_semana].append(r[1])
            mejor_dia_semana = max(dias_semana, key=lambda d: sum(dias_semana[d])/len(dias_semana[d]))

            # Peor rendimiento día (fecha con menos pasajeros)
            peor_rendimiento_dia = min(rows, key=lambda r: r[1])[0]

            # Probabilidad ocupación alta (promedio mensual)
            probabilidad_ocupacion_alta = sum(r[4] for r in rows) / dias if dias else 0

            # Intervalos de confianza de velocidad
            intervalo_confianza_velocidad_min = min(r[5] for r in rows)
            intervalo_confianza_velocidad_max = max(r[6] for r in rows)

            # 3. Insertar o actualizar comparativa mensual
            cur.execute("""
                INSERT INTO comparativa_mensual (
                    mes, año, ruta_id, pasajeros_promedio_dia, pasajeros_total_mes,
                    mejor_dia_semana, peor_rendimiento_dia, velocidad_promedio_mes,
                    probabilidad_ocupacion_alta, intervalo_confianza_velocidad_min,
                    intervalo_confianza_velocidad_max
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (año, mes, ruta_id) DO UPDATE SET
                    pasajeros_promedio_dia = EXCLUDED.pasajeros_promedio_dia,
                    pasajeros_total_mes = EXCLUDED.pasajeros_total_mes,
                    mejor_dia_semana = EXCLUDED.mejor_dia_semana,
                    peor_rendimiento_dia = EXCLUDED.peor_rendimiento_dia,
                    velocidad_promedio_mes = EXCLUDED.velocidad_promedio_mes,
                    probabilidad_ocupacion_alta = EXCLUDED.probabilidad_ocupacion_alta,
                    intervalo_confianza_velocidad_min = EXCLUDED.intervalo_confianza_velocidad_min,
                    intervalo_confianza_velocidad_max = EXCLUDED.intervalo_confianza_velocidad_max
            """, (
                month, year, rid, pasajeros_promedio_dia, pasajeros_total_mes,
                mejor_dia_semana, peor_rendimiento_dia, velocidad_promedio_mes,
                probabilidad_ocupacion_alta, intervalo_confianza_velocidad_min,
                intervalo_confianza_velocidad_max
            ))

            print(f"  Comparativa mensual insertada/actualizada para ruta {rid}")

        conn.commit()
        cur.close()
        conn.close()
        print(f"✅ ETL mensual completado para {year}-{month:02d}")

    except Exception as e:
        print(f"❌ Error en ETL mensual: {e}")