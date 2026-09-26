/* ============================================================
   DDL DEFINITIVO - SISTEMA POS COLEGIO
   PostgreSQL

   Fuente de verdad:
   Esquema LIVE de Supabase reconstruido y validado.

   Total de tablas: 45

   Este archivo contiene únicamente estructura:
   - Extensiones
   - Tablas
   - Restricciones
   - Claves foráneas
   - Índices

   No contiene:
   - INSERT
   - Datos de prueba
   - Carga de productos
   - Carga de variantes
   - Carga de inventario
   - Triggers propios
   - Funciones propias
   - RLS
   - Policies
   ============================================================ */


/* ============================================================
   0. EXTENSIONES
   ============================================================ */

CREATE EXTENSION IF NOT EXISTS pgcrypto;


/* ============================================================
   1. CATALOGOS BASE
   ============================================================ */

CREATE TABLE categoria (
    idcategoria UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(100) NOT NULL,

    CONSTRAINT uq_categoria_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_categoria_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE talla (
    idtalla UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(50) NOT NULL,

    CONSTRAINT uq_talla_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_talla_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE color (
    idcolor UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(50) NOT NULL,

    CONSTRAINT uq_color_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_color_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE ubicacion (
    idubicacion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(100) NOT NULL,

    CONSTRAINT uq_ubicacion_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_ubicacion_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE forma_pago (
    idforma_pago UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(100) NOT NULL,

    CONSTRAINT uq_forma_pago_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_forma_pago_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE tipo_descuento (
    idtipo_descuento UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(100) NOT NULL,

    CONSTRAINT uq_tipo_descuento_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_tipo_descuento_descripcion
        CHECK (btrim(descripcion) <> '')
);


/* ============================================================
   2. SEGURIDAD Y ACCESO
   ============================================================ */

CREATE TABLE rol (
    idrol UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(100) NOT NULL,

    CONSTRAINT uq_rol_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_rol_descripcion
        CHECK (btrim(descripcion) <> '')
);


CREATE TABLE permiso (
    idpermiso UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    descripcion VARCHAR(150) NOT NULL,

    CONSTRAINT uq_permiso_descripcion
        UNIQUE (descripcion),

    CONSTRAINT ck_permiso_descripcion
        CHECK (btrim(descripcion) <> '')
);


/* ============================================================
   3. CLIENTES
   ============================================================ */

CREATE TABLE cliente (
    idcliente UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo VARCHAR(30) NOT NULL,
    dni VARCHAR(30) NOT NULL,
    nombre VARCHAR(150) NOT NULL,
    correo VARCHAR(150) NOT NULL,

    CONSTRAINT uq_cliente_codigo
        UNIQUE (codigo),

    CONSTRAINT uq_cliente_dni
        UNIQUE (dni),

    CONSTRAINT ck_cliente_codigo
        CHECK (btrim(codigo) <> ''),

    CONSTRAINT ck_cliente_dni
        CHECK (btrim(dni) <> ''),

    CONSTRAINT ck_cliente_nombre
        CHECK (btrim(nombre) <> ''),

    CONSTRAINT ck_cliente_correo
        CHECK (btrim(correo) <> '')
);


/* ============================================================
   4. PRODUCTOS
   ============================================================ */

CREATE TABLE producto (
    idproducto UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo VARCHAR(50) NOT NULL,
    descripcion VARCHAR(200) NOT NULL,
    idcategoria UUID NOT NULL,

    CONSTRAINT uq_producto_codigo
        UNIQUE (codigo),

    CONSTRAINT fk_producto_categoria
        FOREIGN KEY (idcategoria)
        REFERENCES categoria (idcategoria)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_producto_codigo
        CHECK (btrim(codigo) <> ''),

    CONSTRAINT ck_producto_descripcion
        CHECK (btrim(descripcion) <> '')
);


/* ============================================================
   5. VARIANTES DE PRODUCTO
   ============================================================ */

CREATE TABLE variante_producto (
    idvariante_producto UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idproducto UUID NOT NULL,
    idtalla UUID NOT NULL,
    idcolor UUID NOT NULL,
    codigo VARCHAR(50) NOT NULL,
    descripcion VARCHAR(200) NOT NULL,
    precio_venta NUMERIC(14,2) NOT NULL,
    genero VARCHAR(20) NOT NULL,

    CONSTRAINT uq_variante_producto_codigo
        UNIQUE (codigo),

    CONSTRAINT uq_variante_producto_combinacion
        UNIQUE (
            idproducto,
            idtalla,
            idcolor,
            genero
        ),

    CONSTRAINT fk_variante_producto_producto
        FOREIGN KEY (idproducto)
        REFERENCES producto (idproducto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_variante_producto_talla
        FOREIGN KEY (idtalla)
        REFERENCES talla (idtalla)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_variante_producto_color
        FOREIGN KEY (idcolor)
        REFERENCES color (idcolor)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_variante_producto_codigo
        CHECK (btrim(codigo) <> ''),

    CONSTRAINT ck_variante_producto_descripcion
        CHECK (btrim(descripcion) <> ''),

    CONSTRAINT ck_variante_producto_precio
        CHECK (precio_venta >= 0),

    CONSTRAINT ck_variante_producto_genero
        CHECK (genero IN ('MUJER', 'UNISEX'))
);


/* ============================================================
   6. RELACION ROL - PERMISO
   ============================================================ */

CREATE TABLE rol_permiso (
    idrol_permiso UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idrol UUID NOT NULL,
    idpermiso UUID NOT NULL,

    CONSTRAINT uq_rol_permiso
        UNIQUE (idrol, idpermiso),

    CONSTRAINT fk_rol_permiso_rol
        FOREIGN KEY (idrol)
        REFERENCES rol (idrol)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_rol_permiso_permiso
        FOREIGN KEY (idpermiso)
        REFERENCES permiso (idpermiso)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   7. EMPLEADOS
   ============================================================ */

CREATE TABLE empleado (
    idempleado UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(150) NOT NULL,
    telefono VARCHAR(30) NOT NULL,
    correo VARCHAR(150) NOT NULL,
    idrol UUID NOT NULL,
    password TEXT NOT NULL,
    estado VARCHAR(30) NOT NULL,

    CONSTRAINT uq_empleado_correo
        UNIQUE (correo),

    CONSTRAINT fk_empleado_rol
        FOREIGN KEY (idrol)
        REFERENCES rol (idrol)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_empleado_nombre
        CHECK (btrim(nombre) <> ''),

    CONSTRAINT ck_empleado_telefono
        CHECK (btrim(telefono) <> ''),

    CONSTRAINT ck_empleado_correo
        CHECK (btrim(correo) <> ''),

    CONSTRAINT ck_empleado_password
        CHECK (btrim(password) <> ''),

    CONSTRAINT ck_empleado_estado
        CHECK (btrim(estado) <> '')
);


/* ============================================================
   8. LUGARES DE ENTREGA
   ============================================================ */

CREATE TABLE lugar_entrega (
    idlugar_entrega UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idcliente UUID NOT NULL,
    descripcion VARCHAR(100) NOT NULL,
    direccion VARCHAR(300) NOT NULL,
    es_predeterminado BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT fk_lugar_entrega_cliente
        FOREIGN KEY (idcliente)
        REFERENCES cliente (idcliente)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_lugar_entrega_descripcion
        CHECK (btrim(descripcion) <> ''),

    CONSTRAINT ck_lugar_entrega_direccion
        CHECK (btrim(direccion) <> '')
);


/* ============================================================
   9. VENTAS
   ============================================================ */

CREATE TABLE venta (
    idventa UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    idcliente UUID NOT NULL,
    idempleado UUID NOT NULL,
    estado VARCHAR(30) NOT NULL,

    CONSTRAINT fk_venta_cliente
        FOREIGN KEY (idcliente)
        REFERENCES cliente (idcliente)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_venta_empleado
        FOREIGN KEY (idempleado)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_venta_estado
        CHECK (btrim(estado) <> '')
);


/* ============================================================
   10. DETALLE DE VENTA
   ============================================================ */

CREATE TABLE detalle_venta (
    iddetalle_venta UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idventa UUID NOT NULL,
    idvariante_producto UUID NOT NULL,
    cantidad INTEGER NOT NULL,
    precio_venta NUMERIC(14,2) NOT NULL,
    subtotal NUMERIC(14,2) NOT NULL,

    CONSTRAINT fk_detalle_venta_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_venta_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_venta_cantidad
        CHECK (cantidad > 0),

    CONSTRAINT ck_detalle_venta_precio
        CHECK (precio_venta >= 0),

    CONSTRAINT ck_detalle_venta_subtotal
        CHECK (subtotal >= 0),

    CONSTRAINT ck_detalle_venta_calculo
        CHECK (subtotal = cantidad * precio_venta)
);


/* ============================================================
   11. TIPOS Y DESCUENTOS
   ============================================================ */

CREATE TABLE descuento (
    iddescuento UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idtipo_descuento UUID NOT NULL,
    descripcion VARCHAR(200) NOT NULL,
    tipo_valor VARCHAR(20) NOT NULL,
    valor NUMERIC(14,2) NOT NULL,
    fecha_inicio TIMESTAMPTZ NOT NULL,
    fecha_fin TIMESTAMPTZ NOT NULL,
    prioridad INTEGER NOT NULL,
    acumulable BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT fk_descuento_tipo
        FOREIGN KEY (idtipo_descuento)
        REFERENCES tipo_descuento (idtipo_descuento)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_descuento_descripcion
        CHECK (btrim(descripcion) <> ''),

    CONSTRAINT ck_descuento_tipo_valor
        CHECK (tipo_valor IN ('PORCENTAJE', 'VALOR')),

    CONSTRAINT ck_descuento_valor
        CHECK (valor >= 0),

    CONSTRAINT ck_descuento_porcentaje
        CHECK (
            tipo_valor <> 'PORCENTAJE'
            OR valor <= 100
        ),

    CONSTRAINT ck_descuento_fechas
        CHECK (fecha_fin >= fecha_inicio),

    CONSTRAINT ck_descuento_prioridad
        CHECK (prioridad >= 0)
);


/* ============================================================
   12. DESCUENTO - VARIANTE
   ============================================================ */

CREATE TABLE descuento_variante (
    iddescuento_variante UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    iddescuento UUID NOT NULL,
    idvariante_producto UUID NOT NULL,

    CONSTRAINT uq_descuento_variante
        UNIQUE (iddescuento, idvariante_producto),

    CONSTRAINT fk_descuento_variante_descuento
        FOREIGN KEY (iddescuento)
        REFERENCES descuento (iddescuento)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_descuento_variante_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   13. DESCUENTO - CLIENTE
   ============================================================ */

CREATE TABLE descuento_cliente (
    iddescuento_cliente UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    iddescuento UUID NOT NULL,
    idcliente UUID NOT NULL,

    CONSTRAINT uq_descuento_cliente
        UNIQUE (iddescuento, idcliente),

    CONSTRAINT fk_descuento_cliente_descuento
        FOREIGN KEY (iddescuento)
        REFERENCES descuento (iddescuento)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_descuento_cliente_cliente
        FOREIGN KEY (idcliente)
        REFERENCES cliente (idcliente)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   14. APLICACION DE DESCUENTO
   ============================================================ */

CREATE TABLE aplicacion_descuento (
    idaplicacion_descuento UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idventa UUID NOT NULL,
    iddescuento UUID NOT NULL,
    tipo_valor_aplicado VARCHAR(20) NOT NULL,
    valor_aplicado NUMERIC(14,2) NOT NULL,

    CONSTRAINT fk_aplicacion_descuento_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_aplicacion_descuento_descuento
        FOREIGN KEY (iddescuento)
        REFERENCES descuento (iddescuento)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_aplicacion_descuento_tipo
        CHECK (tipo_valor_aplicado IN ('PORCENTAJE', 'VALOR')),

    CONSTRAINT ck_aplicacion_descuento_valor
        CHECK (valor_aplicado >= 0),

    CONSTRAINT ck_aplicacion_descuento_porcentaje
        CHECK (
            tipo_valor_aplicado <> 'PORCENTAJE'
            OR valor_aplicado <= 100
        )
);


/* ============================================================
   15. PAGOS
   ============================================================ */

CREATE TABLE pago (
    idpago UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idventa UUID NOT NULL,
    idforma_pago UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valor NUMERIC(14,2) NOT NULL,

    CONSTRAINT fk_pago_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_pago_forma_pago
        FOREIGN KEY (idforma_pago)
        REFERENCES forma_pago (idforma_pago)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_pago_valor
        CHECK (valor > 0)
);


/* ============================================================
   16. FACTURA
   ============================================================ */

CREATE TABLE factura (
    idfactura UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idventa UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    numero VARCHAR(50) NOT NULL,

    CONSTRAINT uq_factura_venta
        UNIQUE (idventa),

    CONSTRAINT uq_factura_numero
        UNIQUE (numero),

    CONSTRAINT fk_factura_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_factura_numero
        CHECK (btrim(numero) <> '')
);


/* ============================================================
   17. RESERVAS
   ============================================================ */

CREATE TABLE reserva (
    idreserva UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idcliente UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    abono_minimo_aplicado NUMERIC(14,2) NOT NULL,
    estado VARCHAR(30) NOT NULL,

    CONSTRAINT fk_reserva_cliente
        FOREIGN KEY (idcliente)
        REFERENCES cliente (idcliente)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_reserva_abono_minimo
        CHECK (abono_minimo_aplicado >= 0),

    CONSTRAINT ck_reserva_estado
        CHECK (btrim(estado) <> '')
);


/* ============================================================
   18. DETALLE DE RESERVA
   ============================================================ */

CREATE TABLE detalle_reserva (
    iddetalle_reserva UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idreserva UUID NOT NULL,
    idvariante_producto UUID NOT NULL,
    cantidad INTEGER NOT NULL,

    CONSTRAINT fk_detalle_reserva_reserva
        FOREIGN KEY (idreserva)
        REFERENCES reserva (idreserva)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_reserva_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_reserva_cantidad
        CHECK (cantidad > 0)
);


/* ============================================================
   19. INVENTARIO
   ============================================================ */

CREATE TABLE inventario (
    idinventario UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idvariante_producto UUID NOT NULL,
    idubicacion UUID NOT NULL,
    stock INTEGER NOT NULL,

    CONSTRAINT uq_inventario_variante_ubicacion
        UNIQUE (idvariante_producto, idubicacion),

    CONSTRAINT fk_inventario_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_inventario_ubicacion
        FOREIGN KEY (idubicacion)
        REFERENCES ubicacion (idubicacion)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_inventario_stock
        CHECK (stock >= 0)
);


/* ============================================================
   20. ASIGNACION DE INVENTARIO
   ============================================================ */

CREATE TABLE asignacion_inventario (
    idasignacion_inventario UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idvariante_producto UUID NOT NULL,
    idubicacion UUID NOT NULL,
    cantidad_inicial INTEGER NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(30) NOT NULL,

    CONSTRAINT fk_asignacion_inventario_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_asignacion_inventario_ubicacion
        FOREIGN KEY (idubicacion)
        REFERENCES ubicacion (idubicacion)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_asignacion_inventario_cantidad
        CHECK (cantidad_inicial > 0),

    CONSTRAINT ck_asignacion_inventario_estado
        CHECK (btrim(estado) <> '')
);


/* ============================================================
   21. ASIGNACION A VENTA
   ============================================================ */

CREATE TABLE asignacion_venta (
    idasignacion_inventario UUID PRIMARY KEY,
    idventa UUID NOT NULL,

    CONSTRAINT fk_asignacion_venta_asignacion
        FOREIGN KEY (idasignacion_inventario)
        REFERENCES asignacion_inventario (idasignacion_inventario)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_asignacion_venta_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   22. ASIGNACION A RESERVA
   ============================================================ */

CREATE TABLE asignacion_reserva (
    idasignacion_inventario UUID PRIMARY KEY,
    idreserva UUID NOT NULL,

    CONSTRAINT fk_asignacion_reserva_asignacion
        FOREIGN KEY (idasignacion_inventario)
        REFERENCES asignacion_inventario (idasignacion_inventario)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_asignacion_reserva_reserva
        FOREIGN KEY (idreserva)
        REFERENCES reserva (idreserva)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   23. MOVIMIENTOS DE INVENTARIO
   ============================================================ */

CREATE TABLE movimiento_inventario (
    idmovimiento_inventario UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idinventario UUID NOT NULL,
    tipo VARCHAR(30) NOT NULL,
    cantidad INTEGER NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    idempleado UUID NOT NULL,
    causa VARCHAR(200) NOT NULL,

    CONSTRAINT fk_movimiento_inventario_inventario
        FOREIGN KEY (idinventario)
        REFERENCES inventario (idinventario)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_movimiento_inventario_empleado
        FOREIGN KEY (idempleado)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_movimiento_inventario_tipo
        CHECK (btrim(tipo) <> ''),

    CONSTRAINT ck_movimiento_inventario_cantidad
        CHECK (cantidad > 0),

    CONSTRAINT ck_movimiento_inventario_causa
        CHECK (btrim(causa) <> '')
);


/* ============================================================
   24. MOVIMIENTOS DE ASIGNACION
   ============================================================ */

CREATE TABLE movimiento_asignacion (
    idmovimiento_asignacion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idasignacion_inventario UUID NOT NULL,
    tipo VARCHAR(30) NOT NULL,
    cantidad INTEGER NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    idempleado UUID NOT NULL,

    CONSTRAINT fk_movimiento_asignacion_asignacion
        FOREIGN KEY (idasignacion_inventario)
        REFERENCES asignacion_inventario (idasignacion_inventario)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_movimiento_asignacion_empleado
        FOREIGN KEY (idempleado)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_movimiento_asignacion_tipo
        CHECK (btrim(tipo) <> ''),

    CONSTRAINT ck_movimiento_asignacion_cantidad
        CHECK (cantidad > 0)
);


/* ============================================================
   25. PROVEEDORES
   ============================================================ */

CREATE TABLE proveedor (
    idproveedor UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(150) NOT NULL,
    telefono VARCHAR(30) NOT NULL,
    correo VARCHAR(150) NOT NULL,

    CONSTRAINT ck_proveedor_nombre
        CHECK (btrim(nombre) <> ''),

    CONSTRAINT ck_proveedor_telefono
        CHECK (btrim(telefono) <> ''),

    CONSTRAINT ck_proveedor_correo
        CHECK (btrim(correo) <> '')
);


/* ============================================================
   26. COMPRAS
   ============================================================ */

CREATE TABLE compra (
    idcompra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idproveedor UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(30) NOT NULL,

    CONSTRAINT fk_compra_proveedor
        FOREIGN KEY (idproveedor)
        REFERENCES proveedor (idproveedor)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_compra_estado
        CHECK (btrim(estado) <> '')
);


/* ============================================================
   27. DETALLE DE COMPRA
   ============================================================ */

CREATE TABLE detalle_compra (
    iddetalle_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idcompra UUID NOT NULL,
    idvariante_producto UUID NOT NULL,
    cantidad INTEGER NOT NULL,
    precio_compra NUMERIC(14,2) NOT NULL,
    subtotal NUMERIC(14,2) NOT NULL,

    CONSTRAINT fk_detalle_compra_compra
        FOREIGN KEY (idcompra)
        REFERENCES compra (idcompra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_compra_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_compra_cantidad
        CHECK (cantidad > 0),

    CONSTRAINT ck_detalle_compra_precio
        CHECK (precio_compra >= 0),

    CONSTRAINT ck_detalle_compra_subtotal
        CHECK (subtotal >= 0),

    CONSTRAINT ck_detalle_compra_calculo
        CHECK (subtotal = cantidad * precio_compra)
);


/* ============================================================
   28. RECEPCION DE COMPRA
   ============================================================ */

CREATE TABLE recepcion_compra (
    idrecepcion_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idcompra UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_recepcion_compra_compra
        FOREIGN KEY (idcompra)
        REFERENCES compra (idcompra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   29. DETALLE DE RECEPCION
   ============================================================ */

CREATE TABLE detalle_recepcion_compra (
    iddetalle_recepcion_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idrecepcion_compra UUID NOT NULL,
    idvariante_producto UUID NOT NULL,
    cantidad_recibida INTEGER NOT NULL,
    cantidad_aceptada INTEGER NOT NULL,
    cantidad_rechazada INTEGER NOT NULL,

    CONSTRAINT fk_detalle_recepcion_compra_recepcion
        FOREIGN KEY (idrecepcion_compra)
        REFERENCES recepcion_compra (idrecepcion_compra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_recepcion_compra_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_recepcion_recibida
        CHECK (cantidad_recibida > 0),

    CONSTRAINT ck_detalle_recepcion_aceptada
        CHECK (cantidad_aceptada >= 0),

    CONSTRAINT ck_detalle_recepcion_rechazada
        CHECK (cantidad_rechazada >= 0),

    CONSTRAINT ck_detalle_recepcion_total
        CHECK (
            cantidad_aceptada + cantidad_rechazada = cantidad_recibida
        )
);


/* ============================================================
   30. DISTRIBUCION DE RECEPCION
   ============================================================ */

CREATE TABLE distribucion_recepcion (
    iddistribucion_recepcion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    iddetalle_recepcion_compra UUID NOT NULL,
    idubicacion UUID NOT NULL,
    cantidad INTEGER NOT NULL,

    CONSTRAINT fk_distribucion_recepcion_detalle
        FOREIGN KEY (iddetalle_recepcion_compra)
        REFERENCES detalle_recepcion_compra (iddetalle_recepcion_compra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_distribucion_recepcion_ubicacion
        FOREIGN KEY (idubicacion)
        REFERENCES ubicacion (idubicacion)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_distribucion_recepcion_cantidad
        CHECK (cantidad > 0)
);


/* ============================================================
   31. DEVOLUCIONES DE COMPRA
   ============================================================ */

CREATE TABLE devolucion_compra (
    iddevolucion_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idcompra UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_devolucion_compra_compra
        FOREIGN KEY (idcompra)
        REFERENCES compra (idcompra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
);


/* ============================================================
   32. DETALLE DE DEVOLUCION
   ============================================================ */

CREATE TABLE detalle_devolucion_compra (
    iddetalle_devolucion_compra UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    iddevolucion_compra UUID NOT NULL,
    idvariante_producto UUID NOT NULL,
    cantidad INTEGER NOT NULL,

    CONSTRAINT fk_detalle_devolucion_compra_devolucion
        FOREIGN KEY (iddevolucion_compra)
        REFERENCES devolucion_compra (iddevolucion_compra)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_devolucion_compra_variante
        FOREIGN KEY (idvariante_producto)
        REFERENCES variante_producto (idvariante_producto)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_devolucion_compra_cantidad
        CHECK (cantidad > 0)
);


/* ============================================================
   33. ENTREGAS
   ============================================================ */

CREATE TABLE entrega (
    identrega UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idventa UUID NOT NULL,
    idlugar_entrega UUID NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    direccion_historica VARCHAR(300) NOT NULL,

    CONSTRAINT fk_entrega_venta
        FOREIGN KEY (idventa)
        REFERENCES venta (idventa)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_entrega_lugar
        FOREIGN KEY (idlugar_entrega)
        REFERENCES lugar_entrega (idlugar_entrega)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_entrega_direccion_historica
        CHECK (btrim(direccion_historica) <> '')
);


/* ============================================================
   34. DETALLE DE ENTREGA
   ============================================================ */

CREATE TABLE detalle_entrega (
    iddetalle_entrega UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identrega UUID NOT NULL,
    idasignacion_inventario UUID NOT NULL,
    cantidad INTEGER NOT NULL,

    CONSTRAINT fk_detalle_entrega_entrega
        FOREIGN KEY (identrega)
        REFERENCES entrega (identrega)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_detalle_entrega_asignacion
        FOREIGN KEY (idasignacion_inventario)
        REFERENCES asignacion_inventario (idasignacion_inventario)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_detalle_entrega_cantidad
        CHECK (cantidad > 0)
);


/* ============================================================
   35. OFERTAS
   ============================================================ */

CREATE TABLE oferta (
    idoferta UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo VARCHAR NOT NULL,
    descripcion VARCHAR NOT NULL,
    iddescuento UUID NOT NULL,

    CONSTRAINT fk_oferta_descuento
        FOREIGN KEY (iddescuento)
        REFERENCES descuento (iddescuento)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_oferta_titulo
        CHECK (btrim(titulo) <> ''),

    CONSTRAINT ck_oferta_descripcion
        CHECK (btrim(descripcion) <> '')
);


/* ============================================================
   36. CONFIGURACION DEL SISTEMA
   ============================================================ */

CREATE TABLE configuracion (
    idconfiguracion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clave VARCHAR(100) NOT NULL,
    valor TEXT NOT NULL,
    descripcion VARCHAR(200) NOT NULL,

    CONSTRAINT uq_configuracion_clave
        UNIQUE (clave),

    CONSTRAINT ck_configuracion_clave
        CHECK (btrim(clave) <> ''),

    CONSTRAINT ck_configuracion_valor
        CHECK (btrim(valor) <> ''),

    CONSTRAINT ck_configuracion_descripcion
        CHECK (btrim(descripcion) <> '')
);


/* ============================================================
   37. SOLICITUDES DE AUTORIZACION
   ============================================================ */

CREATE TABLE solicitud_autorizacion (
    idsolicitud_autorizacion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tipo_accion VARCHAR(50) NOT NULL,
    entidad VARCHAR(100) NOT NULL,
    identificador VARCHAR(100) NOT NULL,
    idempleado_solicitante UUID NOT NULL,
    idempleado_responsable UUID NOT NULL,
    fecha_solicitud TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    descripcion_solicitud VARCHAR(500) NOT NULL,
    tipo_movimiento VARCHAR(20) NOT NULL DEFAULT 'NOAPLICA',
    cantidad_solicitada INTEGER NOT NULL DEFAULT 0,
    estado VARCHAR(20) NOT NULL DEFAULT 'PENDIENTE',

    CONSTRAINT fk_solicitud_autorizacion_solicitante
        FOREIGN KEY (idempleado_solicitante)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT fk_solicitud_autorizacion_responsable
        FOREIGN KEY (idempleado_responsable)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_solicitud_autorizacion_tipo_accion
        CHECK (
            tipo_accion IN (
                'AUTORIZAR_AJUSTE_INVENTARIO',
                'AUTORIZAR_ANULACION_ENTREGA',
                'AUTORIZAR_ANULACION_PAGO',
                'AUTORIZAR_ANULACION_VENTA',
                'AUTORIZAR_DESCUENTO',
                'AUTORIZAR_DESCUENTO_ESPECIAL',
                'AUTORIZAR_OFERTA',
                'AUTORIZAR_OFERTA_ESPECIAL'
            )
        ),

    CONSTRAINT ck_solicitud_autorizacion_entidad
        CHECK (btrim(entidad) <> ''),

    CONSTRAINT ck_solicitud_autorizacion_identificador
        CHECK (btrim(identificador) <> ''),

    CONSTRAINT ck_solicitud_autorizacion_descripcion
        CHECK (btrim(descripcion_solicitud) <> ''),

    CONSTRAINT ck_solicitud_autorizacion_tipo_movimiento
        CHECK (
            tipo_movimiento IN (
                'ENTRADA',
                'SALIDA',
                'NOAPLICA'
            )
        ),

    CONSTRAINT ck_solicitud_autorizacion_cantidad
        CHECK (cantidad_solicitada >= 0),

    CONSTRAINT ck_solicitud_autorizacion_estado
        CHECK (
            estado IN (
                'PENDIENTE',
                'APROBADA',
                'RECHAZADA',
                'CANCELADA'
            )
        )
);


/* ============================================================
   38. AUDITORIA
   ============================================================ */

CREATE TABLE auditoria (
    idauditoria UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempleado UUID NOT NULL,
    accion VARCHAR(50) NOT NULL,
    entidad VARCHAR(100) NOT NULL,
    identificador VARCHAR(100) NOT NULL,
    fecha TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    detalle TEXT NOT NULL,

    CONSTRAINT fk_auditoria_empleado
        FOREIGN KEY (idempleado)
        REFERENCES empleado (idempleado)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT ck_auditoria_accion
        CHECK (btrim(accion) <> ''),

    CONSTRAINT ck_auditoria_entidad
        CHECK (btrim(entidad) <> ''),

    CONSTRAINT ck_auditoria_identificador
        CHECK (btrim(identificador) <> ''),

    CONSTRAINT ck_auditoria_detalle
        CHECK (btrim(detalle) <> '')
);


/* ============================================================
   39. STAGING DE VARIANTES
   ============================================================ */

CREATE TABLE stg_variante_producto (
    id_fuente INTEGER NOT NULL,
    codigo VARCHAR(50) NOT NULL,
    descripcion VARCHAR(200) NOT NULL,
    color VARCHAR(50) NOT NULL,
    genero VARCHAR(20) NOT NULL,
    talla VARCHAR(50) NOT NULL,
    precio_venta NUMERIC(14,2) NOT NULL,

    CONSTRAINT pk_stg_variante_producto
        PRIMARY KEY (id_fuente),

    CONSTRAINT uq_stg_variante_producto_codigo
        UNIQUE (codigo)
);


/* ============================================================
   INDICES
   ============================================================ */


/* ------------------------------------------------------------
   PRODUCTO
   ------------------------------------------------------------ */

CREATE INDEX idx_producto_idcategoria
    ON producto (idcategoria);


/* ------------------------------------------------------------
   ROL - PERMISO
   ------------------------------------------------------------ */

CREATE INDEX idx_rol_permiso_idrol
    ON rol_permiso (idrol);

CREATE INDEX idx_rol_permiso_idpermiso
    ON rol_permiso (idpermiso);


/* ------------------------------------------------------------
   EMPLEADO
   ------------------------------------------------------------ */

CREATE INDEX idx_empleado_idrol
    ON empleado (idrol);


/* ------------------------------------------------------------
   VARIANTE PRODUCTO
   ------------------------------------------------------------ */

CREATE INDEX idx_variante_producto_idproducto
    ON variante_producto (idproducto);

CREATE INDEX idx_variante_producto_idtalla
    ON variante_producto (idtalla);

CREATE INDEX idx_variante_producto_idcolor
    ON variante_producto (idcolor);


/* ------------------------------------------------------------
   LUGAR DE ENTREGA
   ------------------------------------------------------------ */

CREATE INDEX idx_lugar_entrega_idcliente
    ON lugar_entrega (idcliente);


/* ------------------------------------------------------------
   VENTA
   ------------------------------------------------------------ */

CREATE INDEX idx_venta_idcliente
    ON venta (idcliente);

CREATE INDEX idx_venta_idempleado
    ON venta (idempleado);

CREATE INDEX idx_venta_fecha
    ON venta (fecha);


/* ------------------------------------------------------------
   DETALLE DE VENTA
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_venta_idventa
    ON detalle_venta (idventa);

CREATE INDEX idx_detalle_venta_idvariante
    ON detalle_venta (idvariante_producto);


/* ------------------------------------------------------------
   DESCUENTO
   ------------------------------------------------------------ */

CREATE INDEX idx_descuento_idtipo
    ON descuento (idtipo_descuento);

CREATE INDEX idx_descuento_fecha
    ON descuento (fecha_inicio, fecha_fin);


/* ------------------------------------------------------------
   DESCUENTO - VARIANTE
   ------------------------------------------------------------ */

CREATE INDEX idx_descuento_variante_iddescuento
    ON descuento_variante (iddescuento);

CREATE INDEX idx_descuento_variante_idvariante
    ON descuento_variante (idvariante_producto);


/* ------------------------------------------------------------
   DESCUENTO - CLIENTE
   ------------------------------------------------------------ */

CREATE INDEX idx_descuento_cliente_iddescuento
    ON descuento_cliente (iddescuento);

CREATE INDEX idx_descuento_cliente_idcliente
    ON descuento_cliente (idcliente);


/* ------------------------------------------------------------
   APLICACION DE DESCUENTO
   ------------------------------------------------------------ */

CREATE INDEX idx_aplicacion_descuento_idventa
    ON aplicacion_descuento (idventa);

CREATE INDEX idx_aplicacion_descuento_iddescuento
    ON aplicacion_descuento (iddescuento);


/* ------------------------------------------------------------
   PAGO
   ------------------------------------------------------------ */

CREATE INDEX idx_pago_idventa
    ON pago (idventa);

CREATE INDEX idx_pago_idforma_pago
    ON pago (idforma_pago);

CREATE INDEX idx_pago_fecha
    ON pago (fecha);


/* ------------------------------------------------------------
   FACTURA
   ------------------------------------------------------------ */

CREATE INDEX idx_factura_fecha
    ON factura (fecha);


/* ------------------------------------------------------------
   RESERVA
   ------------------------------------------------------------ */

CREATE INDEX idx_reserva_idcliente
    ON reserva (idcliente);

CREATE INDEX idx_reserva_fecha
    ON reserva (fecha);


/* ------------------------------------------------------------
   DETALLE RESERVA
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_reserva_idreserva
    ON detalle_reserva (idreserva);

CREATE INDEX idx_detalle_reserva_idvariante
    ON detalle_reserva (idvariante_producto);


/* ------------------------------------------------------------
   INVENTARIO
   ------------------------------------------------------------ */

CREATE INDEX idx_inventario_idvariante
    ON inventario (idvariante_producto);

CREATE INDEX idx_inventario_idubicacion
    ON inventario (idubicacion);


/* ------------------------------------------------------------
   ASIGNACION INVENTARIO
   ------------------------------------------------------------ */

CREATE INDEX idx_asignacion_inventario_idvariante
    ON asignacion_inventario (idvariante_producto);

CREATE INDEX idx_asignacion_inventario_idubicacion
    ON asignacion_inventario (idubicacion);

CREATE INDEX idx_asignacion_inventario_fecha
    ON asignacion_inventario (fecha);


/* ------------------------------------------------------------
   ASIGNACION VENTA
   ------------------------------------------------------------ */

CREATE INDEX idx_asignacion_venta_idventa
    ON asignacion_venta (idventa);


/* ------------------------------------------------------------
   ASIGNACION RESERVA
   ------------------------------------------------------------ */

CREATE INDEX idx_asignacion_reserva_idreserva
    ON asignacion_reserva (idreserva);


/* ------------------------------------------------------------
   MOVIMIENTO INVENTARIO
   ------------------------------------------------------------ */

CREATE INDEX idx_movimiento_inventario_idinventario
    ON movimiento_inventario (idinventario);

CREATE INDEX idx_movimiento_inventario_idempleado
    ON movimiento_inventario (idempleado);

CREATE INDEX idx_movimiento_inventario_fecha
    ON movimiento_inventario (fecha);


/* ------------------------------------------------------------
   MOVIMIENTO ASIGNACION
   ------------------------------------------------------------ */

CREATE INDEX idx_movimiento_asignacion_idasignacion
    ON movimiento_asignacion (idasignacion_inventario);

CREATE INDEX idx_movimiento_asignacion_idempleado
    ON movimiento_asignacion (idempleado);

CREATE INDEX idx_movimiento_asignacion_fecha
    ON movimiento_asignacion (fecha);


/* ------------------------------------------------------------
   COMPRA
   ------------------------------------------------------------ */

CREATE INDEX idx_compra_idproveedor
    ON compra (idproveedor);

CREATE INDEX idx_compra_fecha
    ON compra (fecha);


/* ------------------------------------------------------------
   DETALLE COMPRA
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_compra_idcompra
    ON detalle_compra (idcompra);

CREATE INDEX idx_detalle_compra_idvariante
    ON detalle_compra (idvariante_producto);


/* ------------------------------------------------------------
   RECEPCION COMPRA
   ------------------------------------------------------------ */

CREATE INDEX idx_recepcion_compra_idcompra
    ON recepcion_compra (idcompra);

CREATE INDEX idx_recepcion_compra_fecha
    ON recepcion_compra (fecha);


/* ------------------------------------------------------------
   DETALLE RECEPCION
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_recepcion_compra_idrecepcion
    ON detalle_recepcion_compra (idrecepcion_compra);

CREATE INDEX idx_detalle_recepcion_compra_idvariante
    ON detalle_recepcion_compra (idvariante_producto);


/* ------------------------------------------------------------
   DISTRIBUCION RECEPCION
   ------------------------------------------------------------ */

CREATE INDEX idx_distribucion_recepcion_iddetalle
    ON distribucion_recepcion (iddetalle_recepcion_compra);

CREATE INDEX idx_distribucion_recepcion_idubicacion
    ON distribucion_recepcion (idubicacion);


/* ------------------------------------------------------------
   DEVOLUCION COMPRA
   ------------------------------------------------------------ */

CREATE INDEX idx_devolucion_compra_idcompra
    ON devolucion_compra (idcompra);

CREATE INDEX idx_devolucion_compra_fecha
    ON devolucion_compra (fecha);


/* ------------------------------------------------------------
   DETALLE DEVOLUCION COMPRA
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_devolucion_compra_iddevolucion
    ON detalle_devolucion_compra (iddevolucion_compra);

CREATE INDEX idx_detalle_devolucion_compra_idvariante
    ON detalle_devolucion_compra (idvariante_producto);


/* ------------------------------------------------------------
   ENTREGA
   ------------------------------------------------------------ */

CREATE INDEX idx_entrega_idventa
    ON entrega (idventa);

CREATE INDEX idx_entrega_idlugar
    ON entrega (idlugar_entrega);

CREATE INDEX idx_entrega_fecha
    ON entrega (fecha);


/* ------------------------------------------------------------
   DETALLE ENTREGA
   ------------------------------------------------------------ */

CREATE INDEX idx_detalle_entrega_identrega
    ON detalle_entrega (identrega);

CREATE INDEX idx_detalle_entrega_idasignacion
    ON detalle_entrega (idasignacion_inventario);


/* ------------------------------------------------------------
   AUDITORIA
   ------------------------------------------------------------ */

CREATE INDEX idx_auditoria_idempleado
    ON auditoria (idempleado);

CREATE INDEX idx_auditoria_fecha
    ON auditoria (fecha);

CREATE INDEX idx_auditoria_entidad_identificador
    ON auditoria (entidad, identificador);


/* ------------------------------------------------------------
   OFERTA
   ------------------------------------------------------------ */

/*
   El esquema LIVE actualmente no posee un índice secundario
   sobre oferta.iddescuento.
*/


/* ------------------------------------------------------------
   PROVEEDOR
   ------------------------------------------------------------ */

/*
   El esquema LIVE actualmente no posee índices secundarios
   sobre proveedor.
*/


/* ------------------------------------------------------------
   STAGING VARIANTE PRODUCTO
   ------------------------------------------------------------ */

/*
   Los índices son creados automáticamente por:
   - pk_stg_variante_producto
   - uq_stg_variante_producto_codigo
*/


/* ============================================================
   INDICE UNICO PARCIAL
   Un solo lugar de entrega predeterminado por cliente.
   ============================================================ */

CREATE UNIQUE INDEX uq_lugar_entrega_predeterminado
    ON lugar_entrega (idcliente)
    WHERE es_predeterminado = TRUE;


/* ============================================================
   FIN DEL DDL DEFINITIVO
   ============================================================ */
