# Library Management REST API

## Description
A simple REST API to manage books, members, and borrowing transactions in a library.

## Tech Stack
- Go
- Gin
- PostgreSQL
- sqlx
- Docker Compose
- Ai

## Endpoints
### Utility
- GET /api/v1/health

### Books
`GET`   : `/books`
`GET`   : `/books/:id`
`POST`  : `/books`
`PUT`   : `/books/:id`
`DELETE`: `/books/:id`

### Members
`GET`    : `/members`
`GET`    : `/members/:id`
`POST`   : `/members`
`PUT`    : `/members/:id`
`DELETE` : `/members/:id`

### Borrowing Transactions
`GET`    : `/borrowings`
`GET`    : `/borrowings/:id`
`POST`   : `/borrowings`
`POST`   : `/borrowings/:id/return`

---

## 🛠️ Getting Started
1. **Clone the repo**
2. **Setup Environment**: Copy `.env.example` to `.env` and fill in your database details.
3. **Run with Docker**:
   ```bash
   docker-compose up -d
   ```
4. **Run the API**:
   ```bash
   go run cmd/api/main.go
   ```
   *(Or use `air` for live reloading)*

---