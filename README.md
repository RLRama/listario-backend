# Listario (backend API)

A RESTful API built with [Go](https://go.dev/), [Iris](https://www.iris-go.com/), and [GORM](https://gorm.io/) to serve as the backend for a full-stack note-taking web app.

## Features

- User system (JWT authentication)
- PostgreSQL database support with GORM
- RESTful API design
- Environment-based configuration
- CRUD operations for notes (currently under development)

## Tech stack

- **Backend**: Go, Iris, GORM
- **Database**: PostgreSQL
- **Auth**: JWT
- **Environment management**: `godotenv`

## Installation

### Prerequisites

- Go (>= 1.24)
- An empty PostgreSQL database

### Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/RLRama/listario-backend.git
   cd listario-backend

2. Copy the example environment file:

   ```sh
   cp .env.example .env
   ```
   > Then update this file with your own credentials.

3. Update frontend origin string

   - In the function `newApp` in the **main.go** file, specifically the `AllowedOrigins` field

5. Install dependencies with

   ```sh
   go mod tidy
   ```

6. Then start the server (this applies auto migrations)

   ```sh
   go run .
   ```

7. Deployment

   #### Build and run

   ```sh
   # build the binary
   go build -o app
   
   # execution
   ./app
   ```

### API endpoints

- Refer to [API_DOCS](API_DOCS.md) to see endpoints and cURL examples of usage.
