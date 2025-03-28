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


