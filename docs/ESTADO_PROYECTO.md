# POS Colegio — Estado del proyecto

## 1. Fuente de verdad y reglas de trabajo

- Repositorio: `backend-pos-colegio`.
- Rama de trabajo: `desarrollo`.
- Rama de producción: `main`.
- Base de datos remota: Supabase PostgreSQL.
- DDL de referencia: `database/ddl/001_estructura.sql`.
- Backend: Go, con separación Handler → Service → Repository.
- Flujo de trabajo: modificar → probar → revisar → commit → push.
- No repetir validaciones ya confirmadas, salvo que cambie el código,
  la base de datos, el entorno o la prueba necesite repetirse.
- Antes de afirmar que un cambio está terminado, distinguir entre
  validación automatizada, prueba manual y revisión de estructura.

## 2. Estado confirmado del backend

### Infraestructura

- [x] Configuración del backend.
- [x] Conexión a Supabase mediante `SUPABASE_DB_URL`.
- [x] Servidor HTTP con Chi.
- [x] Middleware de identificador de petición, recuperación de errores y registro.
- [x] Endpoint de salud `GET /health`.
- [x] Separación Handler → Service → Repository.

### Autenticación y sesiones

- [x] `POST /auth/login` valida las credenciales del empleado.
- [x] La contraseña se verifica con bcrypt.
- [x] El inicio de sesión crea una sesión y devuelve un token.
- [x] En la base de datos se almacena el hash SHA-256 del token, no el token original.
- [x] La duración de sesión se configura mediante `SESSION_EXPIRATION_HOURS`.
- [x] `GET /auth/session` valida el token Bearer.
- [x] Se comprueba que la sesión no esté revocada ni expirada.
- [x] La respuesta de sesión incluye los datos del empleado, rol y permisos.
- [x] La respuesta de consulta de sesión no devuelve el token.

## 3. Pruebas ya confirmadas

### Automatizadas y de compilación

- [x] `gofmt` aplicado a los archivos Go modificados.
- [x] `go test ./...` ejecutado correctamente; los paquetes reportaron
  `[no test files]`. Esto confirma que el comando termina correctamente,
  pero todavía no existen pruebas automatizadas propias.
- [x] `go build ./...` ejecutado correctamente.
- [x] `git diff --check` ejecutado correctamente.

### Pruebas funcionales manuales

- [x] Inicio de sesión válido con el superadministrador de Supabase.
- [x] Consulta de sesión con token válido: HTTP 200.
- [x] Token inválido: HTTP 401.
- [x] Ausencia de token: HTTP 401.
- [x] Sesión revocada: HTTP 401.
- [x] Sesión expirada: HTTP 401.

No guardar tokens reales, contraseñas ni secretos en este documento.

## 4. Estado de la base de datos y del DDL

- [x] La tabla `sesion` existe en el DDL de referencia.
- [x] Se verificó que el DDL contiene 46 sentencias `CREATE TABLE`.
- [x] Se revisó la numeración de las secciones del DDL.
- [x] Se verificó funcionalmente el uso de sesiones en Supabase.
- [ ] No consta una validación completa de ejecución del DDL desde una
  base de datos vacía.
- [ ] No consta una comparación exhaustiva de todas las tablas y
  restricciones del DDL contra el esquema remoto de Supabase.

La existencia de `sesion` y las pruebas funcionales no sustituyen las
dos validaciones pendientes anteriores.

## 5. Cambios pendientes de revisar e integrar

Al iniciar este registro, `git status` mostraba cambios pendientes en:

- `backend/go/.env.example`
- `backend/go/cmd/server/main.go`
- `backend/go/go.mod`
- `backend/go/go.sum`
- `backend/go/internal/config/config.go`
- `backend/go/internal/handler/auth.go`
- `backend/go/internal/repository/auth.go`
- `backend/go/internal/server/server.go`
- `backend/go/internal/service/auth.go`
- `database/ddl/001_estructura.sql`
- `backend/go/internal/middleware/auth.go` (archivo nuevo)

- [ ] Revisar el diff completo de los 11 archivos.
- [ ] Confirmar que `.env.example` no contiene secretos.
- [ ] Repetir las validaciones finales si la revisión produce cambios.
- [ ] Crear el commit de autenticación y sesiones.
- [ ] Enviar el commit a `origin/desarrollo`.

Este documento registra el estado antes de integrar esos cambios; no
significa que ya estén confirmados en Git.

## 6. Siguiente etapa

Después de revisar e integrar autenticación y sesiones:

1. Definir y documentar el mecanismo de autenticación que consumirá el frontend.
2. Implementar y probar la autorización basada en permisos de `rol_permiso`.
3. Incorporar pruebas automatizadas para autenticación y sesiones.
4. Formalizar el procedimiento reproducible de preparación de la base de datos.
5. Continuar con los módulos de negocio siguiendo el esquema vigente.

## 7. Regla para actualizar este documento

Actualizar este archivo cuando una tarea cambie de estado y después de
confirmar sus pruebas. Registrar por separado lo implementado, lo probado,
lo pendiente y la siguiente tarea. No marcar como completado algo que solo
se haya inspeccionado o cuya validación no conste.
