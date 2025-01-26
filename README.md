# Dealls Test App

This is a backend application built with Go (Golang) for Dating site. It includes features such as user authentication, premium subscription handling, and swipe-based interactions. The repository contains robust unit and integration tests to ensure the reliability of the application. This is based on requirement for the minimum specification that needed to be fulfilled.  

---

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Running And Testing](#running-and-testing)
4. [Tech Stack](#tech-stack)
5. [API Endpoints](#api-endpoints)

---

## Installation

### Prerequisites

- Go 1.20 or above
- Git
- A terminal or command-line interface

### Steps

1. Clone this repository:

   ```bash
   git clone https://github.com/yourusername/dealls-test-app.git
   cd dealls-test-app

2. Install dependencies:

   ```bash
   git clone https://github.com/yourusername/dealls-test-app.git
   cd dealls-test-app

3. Ensure you have the `.env` file for environment variables
---

## Configuration

This application requires certain environment variables to be set for proper functionality. These variables can be defined in a `.env` file at the root of the project.

  ```env
  JWT_SECRET_KEY=your_secret_key
  ```

---

## Running And Testing

Run the application ( The server should start at http://localhost:8080.)

  ```bash
  go run ./cmd/main.go
  ```

Testing the application could use these command 

### Unit Testing
   ```bash
    go test ./test/unit_test/...
   ```

### Integration Tests
   ```bash
    go test ./test/integration/api_integration_test.go
   ```

### Postman Tests
1. Open postman and click `import`,
2. Drag or choose 2 files in the folder `test/postman_data`,
3. Make sure the environment `dealls-env` already active,
4. you Could try the collection `dealls-test`,for the protected API with `Authorization`, you could login to get `token` it will automatically register the `JWT Token` into the variable.

---
## Tech Stack
- **Language**: Go
- **Framework**: Gorilla Mux
- **Database**: in-memory Repository
- **Authentication**: JWT
- **Testing**: `testify`,`httptest`

---
## API Endpoints

### Public Endpoints

| Method | Endpoint    | Description          |
|--------|-------------|----------------------|
| POST   | `/signup`   | Register a new user  |
| POST   | `/login`    | Login and get a JWT token |

### Protected Endpoints

| Method | Endpoint           | Description                     |
|--------|--------------------|---------------------------------|
| POST   | `/purchase-premium`| Purchase premium subscription   |
| GET    | `/candidates`      | Get swipe candidates            |
| POST   | `/swipe`           | Swipe on a user                 |

> **Note:** Protected endpoints require a valid `Authorization` header with a JWT token.

