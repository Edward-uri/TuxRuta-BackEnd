from fastapi import APIRouter, Query
from datetime import date
from app.core.database import get_connection
from app.services.etl import run_daily_etl
from calendar import monthrange

router = APIRouter()

@router.get("/metricas-diarias")
def metricas_diarias(fecha: date = Query(...), ruta_id: int = Query(...)):
  
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT * FROM resumen_diario_ruta
        WHERE fecha = %s AND ruta_id = %s
    """, (fecha, ruta_id))
    row = cur.fetchone()
    cur.close()
    conn.close()
    if row:
        return {
            "fecha": row[1],
            "ruta_id": row[2],
            "pasajeros_total": row[3],
            "pasajeros_promedio_por_viaje": float(row[4]),
            "velocidad_promedio": float(row[5]),
            "hora_pico": row[6],
            "total_viajes": row[7],
            "ocupacion_maxima": float(row[8]),
            "probabilidad_ocupacion_alta": float(row[9]),
            "intervalo_confianza_velocidad_min": float(row[10]),
            "intervalo_confianza_velocidad_max": float(row[11]),
        }
    else:
        run_daily_etl(fecha)
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("""
            SELECT * FROM resumen_diario_ruta
            WHERE fecha = %s AND ruta_id = %s
        """, (fecha, ruta_id))
        row = cur.fetchone()
        cur.close()
        conn.close()
        if row:
            return {
                "fecha": row[1],
                "ruta_id": row[2],
                "pasajeros_total": row[3],
                "pasajeros_promedio_por_viaje": float(row[4]),
                "velocidad_promedio": float(row[5]),
                "hora_pico": row[6],
                "total_viajes": row[7],
                "ocupacion_maxima": float(row[8]),
                "probabilidad_ocupacion_alta": float(row[9]),
                "intervalo_confianza_velocidad_min": float(row[10]),
                "intervalo_confianza_velocidad_max": float(row[11]),
            }
        else:
            return {"error": "No hay datos para esa fecha y ruta."}
        
@router.get("/metricas-mensuales")
def metricas_mensuales(
    año: int = Query(..., ge=2020),
    mes: int = Query(..., ge=1, le=12),
    ruta_id: int = Query(...)
):
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT * FROM comparativa_mensual
        WHERE año = %s AND mes = %s AND ruta_id = %s
    """, (año, mes, ruta_id))
    row = cur.fetchone()
    cur.close()
    conn.close()
    if row:
        return {
            "id": row[0],
            "mes": row[1],
            "año": row[2],
            "ruta_id": row[3],
            "pasajeros_promedio_dia": float(row[4]),
            "pasajeros_total_mes": row[5],
            "mejor_dia_semana": row[6],
            "peor_rendimiento_dia": row[7],
            "velocidad_promedio_mes": float(row[8]),
            "probabilidad_ocupacion_alta": float(row[9]),
            "intervalo_confianza_velocidad_min": float(row[10]),
            "intervalo_confianza_velocidad_max": float(row[11]),
        }
    else:
        from app.services.etl_monthly import run_monthly_etl
        run_monthly_etl(año, mes)
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("""
            SELECT * FROM comparativa_mensual
            WHERE año = %s AND mes = %s AND ruta_id = %s
        """, (año, mes, ruta_id))
        row = cur.fetchone()
        cur.close()
        conn.close()
        if row:
            return {
                "id": row[0],
                "mes": row[1],
                "año": row[2],
                "ruta_id": row[3],
                "pasajeros_promedio_dia": float(row[4]),
                "pasajeros_total_mes": row[5],
                "mejor_dia_semana": row[6],
                "peor_rendimiento_dia": row[7],
                "velocidad_promedio_mes": float(row[8]),
                "probabilidad_ocupacion_alta": float(row[9]),
                "intervalo_confianza_velocidad_min": float(row[10]),
                "intervalo_confianza_velocidad_max": float(row[11]),
            }
        else:
            return {"error": "No hay datos para ese mes, año y ruta."}
        

@router.get("/metricas-semanales")
def metricas_semanales(
    año: int = Query(..., ge=2020),
    mes: int = Query(..., ge=1, le=12),
    ruta_id: int = Query(...)
):
    _, last_day = monthrange(año, mes)
    rangos = [
        (1, 7),
        (8, 14),
        (15, 21),
        (22, last_day)
    ]
    resultados = []

    conn = get_connection()
    cur = conn.cursor()

    for inicio, fin in rangos:
        fecha_inicio = date(año, mes, inicio)
        fecha_fin = date(año, mes, fin)
        cur.execute("""
            SELECT
                fecha, pasajeros_total, pasajeros_promedio_por_viaje, velocidad_promedio, hora_pico,
                total_viajes, ocupacion_maxima, probabilidad_ocupacion_alta,
                intervalo_confianza_velocidad_min, intervalo_confianza_velocidad_max
            FROM resumen_diario_ruta
            WHERE ruta_id = %s AND fecha BETWEEN %s AND %s
        """, (ruta_id, fecha_inicio, fecha_fin))
        rows = cur.fetchall()
        if not rows:
            resultados.append({
                "rango": f"{inicio}-{fin}",
                "fecha_inicio": fecha_inicio.isoformat(),
                "fecha_fin": fecha_fin.isoformat(),
                "mensaje": "Sin datos"
            })
            continue

        dias = len(rows)
        pasajeros_total = sum(r[1] for r in rows)
        pasajeros_promedio_por_viaje = sum(r[2] for r in rows) / dias if dias else 0
        velocidad_promedio = sum(r[3] for r in rows) / dias if dias else 0
        from collections import Counter
        horas_pico = [r[4] for r in rows if r[4] != "Sin datos"]
        hora_pico = Counter(horas_pico).most_common(1)[0][0] if horas_pico else "Sin datos"
        total_viajes = sum(r[5] for r in rows)
        ocupacion_maxima = max(r[6] for r in rows)
        probabilidad_ocupacion_alta = sum(r[7] for r in rows) / dias if dias else 0
        intervalo_confianza_velocidad_min = min(r[8] for r in rows)
        intervalo_confianza_velocidad_max = max(r[9] for r in rows)

        resultados.append({
            "rango": f"{inicio}-{fin}",
            "fecha_inicio": fecha_inicio.isoformat(),
            "fecha_fin": fecha_fin.isoformat(),
            "ruta_id": ruta_id,
            "pasajeros_total": pasajeros_total,
            "pasajeros_promedio_por_viaje": pasajeros_promedio_por_viaje,
            "velocidad_promedio": velocidad_promedio,
            "hora_pico": hora_pico,
            "total_viajes": total_viajes,
            "ocupacion_maxima": ocupacion_maxima,
            "probabilidad_ocupacion_alta": probabilidad_ocupacion_alta,
            "intervalo_confianza_velocidad_min": intervalo_confianza_velocidad_min,
            "intervalo_confianza_velocidad_max": intervalo_confianza_velocidad_max
        })

    cur.close()
    conn.close()
    return resultados