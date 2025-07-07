-- ===============================================
-- TuxRuta Database Schema
-- PostgreSQL 15.13-2
-- ===============================================

-- Crear base de datos (ejecutar como superusuario)
-- CREATE DATABASE tuxruta_db;
-- \c tuxruta_db;

-- Crear extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ===============================================
-- TABLA: usuario
-- ===============================================
CREATE TABLE usuario (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    rol VARCHAR(50) DEFAULT 'admin' NOT NULL,
    activo BOOLEAN DEFAULT true NOT NULL,
    ultimo_acceso TIMESTAMP NULL,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Índices para usuario
CREATE INDEX idx_usuario_email ON usuario(email);
CREATE INDEX idx_usuario_activo ON usuario(activo);

-- ===============================================
-- TABLA: ruta
-- ===============================================
CREATE TABLE ruta (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT,
    activa BOOLEAN DEFAULT true NOT NULL,
    creado_por INTEGER NOT NULL,
    modificado_por INTEGER,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modificado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_ruta_creado_por FOREIGN KEY (creado_por) REFERENCES usuario(id),
    CONSTRAINT fk_ruta_modificado_por FOREIGN KEY (modificado_por) REFERENCES usuario(id)
);

-- Índices para ruta
CREATE INDEX idx_ruta_activa ON ruta(activa);
CREATE INDEX idx_ruta_creado_por ON ruta(creado_por);
CREATE INDEX idx_ruta_nombre ON ruta(nombre);

-- ===============================================
-- TABLA: colectivo
-- ===============================================
CREATE TABLE colectivo (
    id SERIAL PRIMARY KEY,
    matricula VARCHAR(50) NOT NULL UNIQUE,
    ruta_id INTEGER NOT NULL,
    activo BOOLEAN DEFAULT true NOT NULL,
    creado_por INTEGER NOT NULL,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_colectivo_ruta FOREIGN KEY (ruta_id) REFERENCES ruta(id),
    CONSTRAINT fk_colectivo_creado_por FOREIGN KEY (creado_por) REFERENCES usuario(id)
);

-- Índices para colectivo
CREATE INDEX idx_colectivo_ruta_id ON colectivo(ruta_id);
CREATE INDEX idx_colectivo_activo ON colectivo(activo);
CREATE INDEX idx_colectivo_matricula ON colectivo(matricula);

-- ===============================================
-- TABLA: parada
-- ===============================================
CREATE TABLE parada (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    latitud DECIMAL(10, 8) NOT NULL,
    longitud DECIMAL(11, 8) NOT NULL,
    activa BOOLEAN DEFAULT true NOT NULL,
    creado_por INTEGER NOT NULL,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_parada_creado_por FOREIGN KEY (creado_por) REFERENCES usuario(id)
);

-- Índices para parada
CREATE INDEX idx_parada_activa ON parada(activa);
CREATE INDEX idx_parada_coordenadas ON parada(latitud, longitud);

-- ===============================================
-- TABLA: ruta_parada
-- ===============================================
CREATE TABLE ruta_parada (
    id SERIAL PRIMARY KEY,
    ruta_id INTEGER NOT NULL,
    parada_id INTEGER NOT NULL,
    orden INTEGER NOT NULL,
    modificado_por INTEGER,
    modificado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    
    -- Foreign Keys
    CONSTRAINT fk_ruta_parada_ruta FOREIGN KEY (ruta_id) REFERENCES ruta(id) ON DELETE CASCADE,
    CONSTRAINT fk_ruta_parada_parada FOREIGN KEY (parada_id) REFERENCES parada(id),
    CONSTRAINT fk_ruta_parada_modificado_por FOREIGN KEY (modificado_por) REFERENCES usuario(id),
    
    -- Unique constraint
    CONSTRAINT uq_ruta_parada_orden UNIQUE (ruta_id, orden),
    CONSTRAINT uq_ruta_parada_relacion UNIQUE (ruta_id, parada_id)
);

-- Índices para ruta_parada
CREATE INDEX idx_ruta_parada_ruta_orden ON ruta_parada(ruta_id, orden);

-- ===============================================
-- TABLA: ubicacion_colectivo
-- ===============================================
CREATE TABLE ubicacion_colectivo (
    id SERIAL PRIMARY KEY,
    colectivo_id INTEGER NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    latitud DECIMAL(10, 8) NOT NULL,
    longitud DECIMAL(11, 8) NOT NULL,
    velocidad DECIMAL(5, 2) DEFAULT 0.0,
    
    -- Foreign Keys
    CONSTRAINT fk_ubicacion_colectivo FOREIGN KEY (colectivo_id) REFERENCES colectivo(id)
);

-- Índices para ubicacion_colectivo
CREATE INDEX idx_ubicacion_colectivo_id ON ubicacion_colectivo(colectivo_id);
CREATE INDEX idx_ubicacion_fecha_hora ON ubicacion_colectivo(fecha_hora DESC);
CREATE INDEX idx_ubicacion_colectivo_fecha ON ubicacion_colectivo(colectivo_id, fecha_hora DESC);

-- ===============================================
-- TABLA: registro_pasajeros
-- ===============================================
CREATE TABLE registro_pasajeros (
    id SERIAL PRIMARY KEY,
    colectivo_id INTEGER NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    personas_abordo INTEGER NOT NULL DEFAULT 0,
    
    -- Foreign Keys
    CONSTRAINT fk_registro_colectivo FOREIGN KEY (colectivo_id) REFERENCES colectivo(id),
    
    -- Constraints
    CONSTRAINT chk_personas_abordo_positivo CHECK (personas_abordo >= 0)
);

-- Índices para registro_pasajeros
CREATE INDEX idx_registro_colectivo_id ON registro_pasajeros(colectivo_id);
CREATE INDEX idx_registro_fecha_hora ON registro_pasajeros(fecha_hora DESC);
CREATE INDEX idx_registro_colectivo_fecha ON registro_pasajeros(colectivo_id, fecha_hora DESC);

-- ===============================================
-- TABLA: alerta_evento
-- ===============================================
CREATE TABLE alerta_evento (
    id SERIAL PRIMARY KEY,
    colectivo_id INTEGER NOT NULL,
    fecha_hora TIMESTAMP NOT NULL,
    tipo_alerta VARCHAR(100) NOT NULL,
    descripcion TEXT,
    resuelto BOOLEAN DEFAULT false NOT NULL,
    resuelto_por INTEGER,
    fecha_resolucion TIMESTAMP,
    
    -- Foreign Keys
    CONSTRAINT fk_alerta_colectivo FOREIGN KEY (colectivo_id) REFERENCES colectivo(id),
    CONSTRAINT fk_alerta_resuelto_por FOREIGN KEY (resuelto_por) REFERENCES usuario(id)
);

-- Índices para alerta_evento
CREATE INDEX idx_alerta_colectivo_id ON alerta_evento(colectivo_id);
CREATE INDEX idx_alerta_resuelto ON alerta_evento(resuelto);
CREATE INDEX idx_alerta_fecha_hora ON alerta_evento(fecha_hora DESC);
CREATE INDEX idx_alerta_tipo ON alerta_evento(tipo_alerta);

-- ===============================================
-- TABLA: resumen_diario_ruta
-- ===============================================
CREATE TABLE resumen_diario_ruta (
    id SERIAL PRIMARY KEY,
    fecha DATE NOT NULL,
    ruta_id INTEGER NOT NULL,
    pasajeros_total INTEGER DEFAULT 0,
    pasajeros_promedio_por_viaje DECIMAL(8, 2) DEFAULT 0.0,
    velocidad_promedio DECIMAL(5, 2) DEFAULT 0.0,
    hora_pico VARCHAR(20),
    total_viajes INTEGER DEFAULT 0,
    total_alertas INTEGER DEFAULT 0,
    ocupacion_maxima DECIMAL(8, 2) DEFAULT 0.0,
    probabilidad_ocupacion_alta DECIMAL(5, 4) DEFAULT 0.0,
    intervalo_confianza_velocidad_min DECIMAL(5, 2) DEFAULT 0.0,
    intervalo_confianza_velocidad_max DECIMAL(5, 2) DEFAULT 0.0,
    
    -- Foreign Keys
    CONSTRAINT fk_resumen_diario_ruta FOREIGN KEY (ruta_id) REFERENCES ruta(id),
    
    -- Unique constraint
    CONSTRAINT uq_resumen_diario_fecha_ruta UNIQUE (fecha, ruta_id)
);

-- Índices para resumen_diario_ruta
CREATE INDEX idx_resumen_diario_fecha ON resumen_diario_ruta(fecha DESC);
CREATE INDEX idx_resumen_diario_ruta_id ON resumen_diario_ruta(ruta_id);
CREATE INDEX idx_resumen_diario_lookup ON resumen_diario_ruta(fecha, ruta_id);

-- ===============================================
-- TABLA: comparativa_mensual
-- ===============================================
CREATE TABLE comparativa_mensual (
    id SERIAL PRIMARY KEY,
    mes INTEGER NOT NULL,
    año INTEGER NOT NULL,
    ruta_id INTEGER NOT NULL,
    pasajeros_promedio_dia DECIMAL(8, 2) DEFAULT 0.0,
    pasajeros_total_mes INTEGER DEFAULT 0,
    mejor_dia_semana INTEGER,
    peor_rendimiento_dia DATE,
    crecimiento_porcentual DECIMAL(5, 2) DEFAULT 0.0,
    velocidad_promedio_mes DECIMAL(5, 2) DEFAULT 0.0,
    alertas_total_mes INTEGER DEFAULT 0,
    dias_activos INTEGER DEFAULT 0,
    probabilidad_ocupacion_alta DECIMAL(5, 4) DEFAULT 0.0,
    intervalo_confianza_velocidad_min DECIMAL(5, 2) DEFAULT 0.0,
    intervalo_confianza_velocidad_max DECIMAL(5, 2) DEFAULT 0.0,
    
    -- Foreign Keys
    CONSTRAINT fk_comparativa_mensual_ruta FOREIGN KEY (ruta_id) REFERENCES ruta(id),
    
    -- Constraints
    CONSTRAINT chk_mes_valido CHECK (mes >= 1 AND mes <= 12),
    CONSTRAINT chk_año_valido CHECK (año >= 2020),
    CONSTRAINT chk_mejor_dia_semana CHECK (mejor_dia_semana >= 1 AND mejor_dia_semana <= 7),
    
    -- Unique constraint
    CONSTRAINT uq_comparativa_mensual UNIQUE (año, mes, ruta_id)
);

-- Índices para comparativa_mensual
CREATE INDEX idx_comparativa_periodo ON comparativa_mensual(año, mes);
CREATE INDEX idx_comparativa_ruta_id ON comparativa_mensual(ruta_id);

-- ===============================================
-- TRIGGERS PARA ACTUALIZAR modified_en
-- ===============================================

-- Función para actualizar timestamp
CREATE OR REPLACE FUNCTION actualizar_modificado_en()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modificado_en = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para tabla ruta
CREATE TRIGGER trigger_ruta_modificado_en
    BEFORE UPDATE ON ruta
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_modificado_en();

-- Trigger para tabla ruta_parada
CREATE TRIGGER trigger_ruta_parada_modificado_en
    BEFORE UPDATE ON ruta_parada
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_modificado_en();

-- ===============================================
-- DATOS INICIALES
-- ===============================================

-- Usuario administrador por defecto
INSERT INTO usuario (email, password, rol, activo) 
VALUES 
    ('admin@tuxruta.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', true),
    ('secretaria@transporte.gov', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', true);

-- Comentarios sobre la estructura
COMMENT ON DATABASE tuxruta_db IS 'Base de datos para sistema TuxRuta - Gestión de transporte público';
COMMENT ON TABLE usuario IS 'Usuarios administrativos con acceso al sistema';
COMMENT ON TABLE ruta IS 'Rutas oficiales de transporte público';
COMMENT ON TABLE colectivo IS 'Unidades de transporte asociadas a rutas';
COMMENT ON TABLE parada IS 'Paradas físicas del sistema de transporte';
COMMENT ON TABLE ruta_parada IS 'Relación ordenada de paradas en cada ruta';
COMMENT ON TABLE ubicacion_colectivo IS 'Tracking GPS en tiempo real de colectivos';
COMMENT ON TABLE registro_pasajeros IS 'Conteo de pasajeros por colectivo';
COMMENT ON TABLE alerta_evento IS 'Eventos y alertas del sistema';
COMMENT ON TABLE resumen_diario_ruta IS 'Estadísticas diarias consolidadas por ruta';
COMMENT ON TABLE comparativa_mensual IS 'Análisis estadístico mensual por ruta';

-- ===============================================
-- VERIFICACIÓN FINAL
-- ===============================================



