# Sancho

**Sancho** es un wrapper para [streamrip](https://github.com/nathom/streamrip) que permite **descargar y gestionar música para múltiples usuarios** desde una interfaz web moderna.  
Cada usuario tiene su propia carpeta con _symlinks_ hacia una carpeta global, evitando duplicados y optimizando el uso de espacio.

---

## ✨ Características

- 🎵 **Descarga de música** desde múltiples fuentes usando _streamrip_.
- 👥 **Soporte multiusuario** con carpetas personales y enlaces simbólicos a la biblioteca global.
- 📁 **Evita duplicados** mediante un sistema de enlaces simbólicos.
- 🌐 **Interfaz web SPA** construida con **Svelte + SvelteKit**.
- ⚡ **Backend en Go** usando **Gin**, **golang-migrate** y **sqlc**.
- 🗄️ **Base de datos SQLite** ligera y sin dependencias externas.
- 🐳 **Despliegue en Docker** con un solo contenedor.
- 🔒 Control de acceso por usuario.

---

## 📦 Requisitos

- [Docker](https://www.docker.com/)
- (Opcional para desarrollo) [Go](https://go.dev/) >= 1.22
- (Opcional para desarrollo) [Node.js](https://nodejs.org/) >= 20 y [pnpm](https://pnpm.io/)

---

## 🚀 Instalación

### Ejecutar Sancho con Docker

```bash
#!/bin/bash
docker run -d -p 5400:5400 \
  --name=sancho \
  --volume /folder/to/config:/data \
  --volume /folder/to/libraries:/sancho \
  --env QOBUZ_PASSWORD_OR_TOKEN="your_qobuz_token" \
  --env QOBUZ_USER_ID="qobuz_app_id" \
  --restart=unless-stopped \
  alebdc/sancho:latest
```
