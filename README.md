# api-go

## Requisitos

Para ejecutar el proyecto hay que tener instalado:

- Go (versión 1.18)

En mi caso utilice nginx para correr el html


## Instalación y Ejecución

1. Clonar el repositorio, para verlo en el localhost correr el nginx para ver el frontend 

2. Instalar las dependencias necesarias:
   ```sh
   go mod tidy
   ```

3. Ejecutar el archivo go :
   ```sh
   go run main.go
   ```
4. La API se ejecutará en `http://localhost:8080`.

## Funcion de los endpoints utilizados
Crear un partido (POST /api/matches)

    Recibe un objeto JSON con los datos del partido.

    Guarda el partido en memoria y le asigna un ID único.

    Devuelve el partido creado con su ID asignado.

Obtener todos los partidos (GET /api/matches)

    Retorna un JSON con todos los partidos almacenados en memoria.

Obtener un partido por ID (GET /api/matches/{id})

    Busca un partido por su ID en la memoria.

    Si existe, lo devuelve en formato JSON.

    Si no existe, devuelve un error 404.

Actualizar un partido (PUT /api/matches/{id})

    Recibe un JSON con la nueva información del partido.

    Busca el partido por su ID y lo actualiza si existe.

    Si no existe, devuelve un error 404.

Eliminar un partido (DELETE /api/matches/{id})

    Busca el partido por su ID y lo elimina de la memoria.

    Si el partido no existe, devuelve un error 404.

## Dependencias utilizadas

### `github.com/gorilla/mux`
Se usa para manejar rutas dinámicas en la API, permitiendo utilizar parámetros en las URLs como `/{id}`

### `github.com/rs/cors`
Permite configurar CORS para que la API pueda ser accedida desde diferentes orígenes

## Implementación del Almacenamiento JSON

Los datos se guardan en el archivo matches.json
Cada modificación se sincroniza automáticamente con el archivo
Al iniciar el servicio, los datos se cargan desde el archivo

## Operaciones PATCH para Registro de Eventos
## Registrar Gol (PATCH /api/matches/{id}/goals)

Incrementa el contador de goles del equipo local.
Si el partido no existe, devuelve un error 404.

## Registrar Tarjeta Amarilla (PATCH /api/matches/{id}/yellowcards)

Incrementa el contador de tarjetas amarillas del partido.
Si el partido no existe, devuelve un error 404.

## Registrar Tarjeta Roja (PATCH /api/matches/{id}/redcards)

Incrementa el contador de tarjetas rojas del partido.
Si el partido no existe, devuelve un error 404.

## Establecer Tiempo Extra (PATCH /api/matches/{id}/extratime)

Establece la bandera de tiempo extra a true para el partido.
Si el partido no existe, devuelve un error 404.


## Resultado imagenes:
Listar los partidos:

![WhatsApp Image 2025-03-27 at 9 18 10 PM](https://github.com/user-attachments/assets/17dcdd89-af9c-4b85-8ac2-f6b77d9c34ff)

Crear nuevo partido:
![WhatsApp Image 2025-03-27 at 9 19 13 PM](https://github.com/user-attachments/assets/3412564c-58a3-4605-9702-074a906f98d5)

![WhatsApp Image 2025-03-27 at 9 19 33 PM](https://github.com/user-attachments/assets/e39f9efa-f18e-4982-a414-40c4c3d19591)

Actualizar partido: 
![WhatsApp Image 2025-03-27 at 9 20 23 PM](https://github.com/user-attachments/assets/6a8b50a3-5456-4129-b4e9-d09ac15fac48)

![WhatsApp Image 2025-03-27 at 9 20 39 PM](https://github.com/user-attachments/assets/5e360f7c-90f5-4fbd-84c8-c5cc6626b66f)

Buscar partidos:
![WhatsApp Image 2025-03-27 at 9 20 54 PM](https://github.com/user-attachments/assets/9407386c-3ddf-4510-bfe5-4d2eec7ed1c4)

Borrar partidos:
![WhatsApp Image 2025-03-27 at 9 21 26 PM](https://github.com/user-attachments/assets/87f3befa-ec6d-4df6-83b9-be0ef4cf4652)








