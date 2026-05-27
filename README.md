# 📦 Sistema POS (Point of Sale) - Proyecto Gestión de proyecto de software

Un sistema de Punto de Venta modular de alto rendimiento. Este proyecto utiliza una arquitectura basada en dominios (Clean Architecture) para el backend y una interfaz reactiva para el frontend, todo estandarizado a través de contenedores Docker para garantizar consistencia en el desarrollo y despliegue.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-0088CC?style=for-the-badge&logo=go&logoColor=white)
![Svelte](https://img.shields.io/badge/Svelte-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)

## 📋 Requisitos Previos

Para desarrollar en este proyecto **no necesitas instalar Go, Node.js ni PostgreSQL en tu sistema operativo local**. Toda la infraestructura está contenerizada. Solo necesitas:

* [Docker](https://www.docker.com/products/docker-desktop/) (Engine activo)
* Docker Compose (v2 o superior)
* Git

---

## 🚀 Entorno de Desarrollo (Quick Start)

Sigue estos pasos para levantar el entorno local de desarrollo por primera vez:

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/kotonitro/GPS_2026_1.git](https://github.com/kotonitro/GPS_2026_1.git)
   cd GPS_2026_1

2. **Construir y levantar la infraestructura:**
Este comando descarga las imágenes, instala las dependencias de los contenedores y levanta la base de datos PostgreSQL junto con el entorno de programación en segundo plano.
    ```bash
   docker compose up --build -d

3. **Ingresar al entorno de desarrollo:**
Para ejecutar comandos de Go (como compilar, instalar librerías o correr el servidor), debes entrar al contenedor:
    ```Bash
    # Usa 'sh' o 'bash' dependiendo de la configuración de tu terminal
    docker compose exec backend bash

4. **Sincronizar dependencias e iniciar el servidor (Dentro del contenedor):**
    ```bash
    cd backend
    go mod tidy
    go run cmd/main.go

## Arquitectura del Proyecto

El backend está estructurado siguiendo el Patrón Repositorio y separación de responsabilidades, dividiendo la lógica en Base de Dato y Controlador HTTP.

    📦 Raíz del Proyecto
    ┣ 📂 backend/
    ┃ ┣ 📂 cmd/
    ┃ ┃ └── main.go                 # Punto de entrada y registro de rutas Gin
    ┃ ┣ 📂 internal/
    ┃ ┃ ┣ 📂 database/              # Conexión Singleton a PostgreSQL
    ┃ ┃ ┣ 📂 inventario/            # Dominio modular (Ejemplo)
    ┃ ┃ ┃ ┣ 📜 controladores.go     # Capa HTTP (Recibe JSON, responde JSON)
    ┃ ┃ ┃ ┣ 📜 modelos.go           # Estructuras de datos (Structs y Tags GORM)
    ┃ ┃ ┃ └── 📜 repositorio.go     # Capa de Persistencia (Consultas GORM/SQL)
    ┃ ┃ ┗ 📂 ventas/                # Otros dominios...
    ┃ ┣ 📜 go.mod                   # Gestor de dependencias de Go
    ┃ ┗ 📜 Dockerfile.dev           # Receta del contenedor de desarrollo
    ┣ 📂 frontend/                  # Interfaz de usuario en Svelte
    ┣ 📜 docker-compose.yml         # Orquestador de servicios
    ┗ 📜 .gitignore                 # Exclusión de binarios y datos locales

## Reglas de Contribución para el Equipo
Para mantener la estabilidad del código base, todos los desarrolladores deben adherirse a las siguientes normas:

GORM para persistencia: El 90% de las operaciones a la base de datos se manejan mediante el ORM (GORM). Las consultas en SQL puro solo se reservan para reportes estadísticos complejos.

Aislamiento de Controladores: El archivo controladores.go no debe contener llamadas directas a GORM (db.Create, db.Find). Toda interacción con la base de datos se delega al repositorio.go.

Manejo de Errores: En Go no se ignoran los errores. Todo error devuelto por el Repositorio debe ser capturado en el Controlador y respondido al cliente con el código HTTP correspondiente (ej: 400, 500).

Sincronización de Módulos: Si instalas una nueva librería (ej: go get ...), es obligatorio ejecutar go mod tidy antes de realizar tu commit para evitar desincronización en el archivo go.sum del equipo.

Ramas de Git: Trabajar siempre en ramas dedicadas (feature/modulo-clientes, fix/bug-inventario) y realizar Merge a main únicamente cuando el código ha sido probado localmente en el contenedor.


## Comandos de Mantenimiento Docker
Ver registros en tiempo real: docker compose logs -f

Apagar el entorno al terminar el día: docker compose down

Apagar el entorno y BORRAR la base de datos local (Reset): docker compose down -v