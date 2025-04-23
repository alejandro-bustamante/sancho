# Arquitectura del Backend - Proyecto Streamrip Wrapper

Este documento resume la arquitectura utilizada en el backend del proyecto que envuelve `streamrip` con una API REST amigable para el usuario, utilizando Go como lenguaje principal, SQLite como base de datos (con soporte ICU), `taglib` para manejo de metadatos y `gin` como router http.

## ✨ Objetivo
Proveer una estructura modular, escalable y mantenible, basada en un patrón **MVC extendido** (con servicios y repositorios), orientada a una API REST.

---

## 📂 Estructura General

```
/proyecto-streamrip
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Punto de entrada de la aplicación
│   ├── api/
│   │   └── routes.go                # Definición centralizada de rutas
│   ├── internal/
│   │   ├── config/                  # Configuración (env, puertos, etc)
│   │   ├── controller/              # Controladores HTTP (manejo de requests)
│   │   ├── model/                   # Estructuras de datos (entidades)
│   │   ├── repository/              # Abstracción e implementación de acceso a datos
│   │   ├── service/                 # Lógica de negocio y clientes externos (streamrip, taglib)
│   ├── migrations/                  # Scripts SQL para el esquema inicial
│   ├── pkg/util/                    # Utilidades compartidas
│   ├── Dockerfile                   # Imagen para entorno de producción
│   └── go.mod                       # Dependencias de Go
```

---

## 🔀 Patrón de Arquitectura

Basado en **MVC extendido**:

| Capa           | Descripción                                                                 |
|----------------|-------------------------------------------------------------------------------|
| Model          | Estructuras de datos puras (`model/Song.go`, `User.go`)                      |
| Controller     | Recibe peticiones HTTP, llama a servicios, devuelve respuesta (vista)        |
| View (implícita) | La "vista" es el `c.JSON(...)` o `c.String(...)` devuelto al cliente              |
| Service        | Lógica de negocio, orquestación, validaciones, llamadas a wrappers externos   |
| Repository     | Acceso a base de datos (con interfaces y backends concretos como SQLite)     |

---

## 💡 Roles de Carpetas Clave

### `/cmd/server/main.go`
- Punto de entrada.
- Inicializa la app, inyecta dependencias, arranca el router.
- Conviene mantenerlo dentro de `cmd/` por si se agregan otras apps (CLI, workers, etc.).

### `/api/routes.go`
- Centraliza el registro de rutas.
- Separa la definición de rutas de la lógica del controlador.
- Facilita la lectura y el versionado (ej. `/api/v1/...`).

### `/controller/`
- Define handlers HTTP.
- No contiene lógica de negocio ni acceso a datos.
- Devuelve respuestas JSON (actúa como "vista").

### `/service/`
- Encapsula la lógica de negocio.
- Usa los repositorios y wrappers.
- Ejemplos:
  - `SongService` decide si descargar una canción, guarda metadatos, evita duplicados.
  - `StreamripClient` es una interfaz para comunicar con el wrapper en Python.

### `/repository/`
- Define interfaces para acceso a datos.
- Implementaciones concretas, como `sqlite`, van en subcarpetas.

### `/model/`
- Estructuras de datos (entities, DTOs).
- No contiene lógica.

---

## ⚛️ Inyección de Dependencias

El proyecto favorece la inyección de dependencias manual:

### ✅ Ejemplo
```go
func main() {
    db := sqlite.NewConnection()
    songRepo := sqlite.NewSongRepo(db)
    streamripClient := streamrip.NewHttpClient("http://localhost:5000")
    songService := service.NewSongService(songRepo, streamripClient)
    songController := controller.NewSongController(songService)

    router := api.SetupRouter(songController)
    router.Run(":8080")
}
```

Esto:
- Permite mocks en tests.
- Desacopla implementaciones concretas.

---

## 🚀 Flujo General de una Descarga

```plaintext
frontend (JS)
   ↓
controller/song_controller.go → recibe la petición
   ↓
service/song_service.go → decide si hay que descargar o no
   ↓
   ├── consulta repo (evitar duplicados)
   └── llama a streamrip (cliente HTTP)
           ↓
       streamrip (Python) devuelve path y metadatos
           ↓
       se guardan datos en base de datos
           ↓
       se responde al frontend
```

---

## 🔍 Buenas Prácticas Clave

- **No acoples controladores a servicios o DB directamente.** Usa interfaces.
- **La lógica de negocio nunca debe estar en el controlador.**
- **El controlador actúa como puente HTTP → servicio.**
- **No usar modelos como repositorios.** Los modelos solo definen estructuras.
- **Separar lógica de negocio del acceso a datos.** Facilita testeo y cambios futuros.
- **Encapsular wrappers externos (como streamrip o taglib) en servicios.**
- **Centralizar las rutas en un solo archivo.** Mejora la organización.

---

## 🔗 Futuras extensiones

- Agregar workers en `cmd/worker/` para tareas async.
- Tests con `Testify` para controladores, servicios y repositorios.
- Agregar soporte para configuración via `.env` en `config/`.
- Usar `go generate` para mockear interfaces con `mockgen`.

---

## 🚫 Antipatrones a evitar

- Repetir lógica en controladores.
- Usar estructuras de modelo para validación compleja.
- Llamar a la base directamente desde el controlador.
- Acoplar la implementación concreta (ej. SQLite) al servicio.

---
