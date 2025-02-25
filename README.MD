# Net HTTP Boilerplate

This project is a boilerplate for building a RESTful API using Go, Chi router, GORM for database interactions, and JWT for authentication. It provides a structured and scalable foundation for developing web applications.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Features

- RESTful API with CRUD operations
- JWT-based authentication
- Middleware for logging, CORS, and request ID
- GORM for database interactions
- Structured project layout
- Environment variable configuration

## Getting Started

### Prerequisites

- Go 1.23.1 or higher
- PostgreSQL

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/net-http-boilerplate.git
   cd net-http-boilerplate
   ```
2. Copy the example environment file and configure it:
    ```sh
    cp .env.example .env
    ```
3. Install the dependencies
    ```sh
    go mod download
    ```

### Running the Application

1. Build the application
    ```sh
    make build
    ```
2. Run the aplication 
    ```sh
    make run
    ```
3. The server will start on the port specified in command (default is 8080)

## Project Structure
```
.
├── cmd/
│   ├── api.go
│   └── root.go
├── docs/
│   └── README.md
├── internal/
│   ├── api/
│   │   ├── middleware.go
│   │   └── resp/
│   │       └── writer.go
│   ├── auth/
│   │   └── middleware.go
│   ├── category/
│   │   ├── dto.go
│   │   ├── errors.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── service.go
│   ├── config/
│   │   └── config.go
│   ├── entity/
│   │   ├── Category.go
│   │   ├── Health.go
│   │   ├── Post.go
│   │   └── User.go
│   ├── jwt/
│   │   └── jwt.go
│   ├── post/
│   │   ├── dto.go
│   │   ├── error.go
│   │   ├── handler.go
│   │   ├── repo.go
│   │   └── service.go
│   ├── postgres/
│   │   ├── gorm.go
│   │   └── migrate.go
│   └── user/
│       ├── dto.go
│       ├── errors.go
│       ├── handler.go
│       ├── repo.go
│       └── service.go
├── main.go
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Configuration
The application uses environment variables for configuration. Create a `.env` file in the root directory and configure the following variables:
```.env
# Postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=posts_db

# JWT
JWT_SECRET=key
JWT_ISSUER=issuer
```

## API Endpoints

### Authentication

- `POST /users/register` - Register a new user
- `POST /users/login` - Login and get JWT tokens
- `POST /users/refresh` - Refresh JWT token

### Posts

- `GET /posts`- Get all posts
- `POST /posts` - Create a new post
- `GET /posts/{id}` - Get a post by ID
- `PUT /posts/{id}` - Update a post by ID
- `DELETE /posts/{id}` - Delete a post by ID

### Categories

- `POST /categories` - Create a new category
- `GET /categories` - Get all categories
- `PUT /categories/{id}` - Update a category
- `DELETE /categories/{id}` - Delete a category

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License
This project is licensed under the MIT License.