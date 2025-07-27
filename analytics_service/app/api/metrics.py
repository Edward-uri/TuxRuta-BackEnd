from fastapi import APIRouter, Query
from datetime import date
from app.core.database import get_connection
from app.services.etl import run_daily_etl
from app.services.etl_monthly import run_monthly_etl
from calendar import monthrange

router = APIRouter()

@router.get("/metricas-diarias")
def metricas_diarias(fecha: date = Query(...), ruta_id: int = Query(...)):
    # Ejecuta ETL diario antes de mostrar datos
    run_daily_etl(fecha, ruta_id)
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
        "colectivo_id": row[3],
        "pasajeros_total": row[4],
        "pasajeros_promedio_por_viaje": float(row[5]) if row[5] is not None else 0.0,
        "velocidad_promedio": float(row[6]) if row[6] is not None else 0.0,
        "hora_pico": row[7],
        "total_viajes": row[8],
        "ocupacion_maxima": float(row[9]) if row[9] is not None else 0.0,
        "probabilidad_ocupacion_alta": float(row[10]) if row[10] is not None else 0.0,
        "intervalo_confianza_velocidad_min": float(row[11]) if row[11] is not None else 0.0,
        "intervalo_confianza_velocidad_max": float(row[12]) if row[12] is not None else 0.0,
    }
    else:
        return {"error": "No hay datos para esa fecha y ruta."}

@router.get("/metricas-mensuales")
def metricas_mensuales(
    año: int = Query(..., ge=2020),
    mes: int = Query(..., ge=1, le=12),
    ruta_id: int = Query(...)
):
    # Ejecuta ETL mensual antes de mostrar datos
    run_monthly_etl(año, mes, ruta_id)
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT * FROM comparativa_mensual
        WHERE año = %s AND mes = %s AND ruta_id = %s
    """, (año, mes, ruta_id))
    row = cur.fetchone()
    cur.close()
    conn.close()
    def safe_float(val):
        return float(val) if val is not None else 0.0
    if row:
        return {
            "id": row[0],
            "mes": row[1],
            "año": row[2],
            "ruta_id": row[3],
            "colectivo_id": row[4],
            "pasajeros_promedio_dia": safe_float(row[5]),
            "pasajeros_total_mes": row[6] if row[6] is not None else 0,
            "mejor_dia_semana": row[7],
            "peor_rendimiento_dia": row[8].isoformat() if row[8] else None,
            "velocidad_promedio_mes": safe_float(row[9]),
            "probabilidad_ocupacion_alta": safe_float(row[10]),
            "intervalo_confianza_velocidad_min": safe_float(row[11]),
            "intervalo_confianza_velocidad_max": safe_float(row[12]),
        }
    else:
        return {"error": "No hay datos para ese mes, año y ruta."}