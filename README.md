# Library Management REST API

## Description
A simple REST API to manage books, members, and borrowing transactions in a library.

## Tech Stack
- Go
- Gin
- PostgreSQL
- sqlx
- Docker Compose

## Actors
- Admin / Librarian

## Main Features
### Books
- Create book
- Get all books
- Get book detail
- Update book
- Delete book

### Members
- Create member
- Get all members
- Get member detail
- Update member
- Delete member

### Borrowings
- Create borrowing
- Return borrowed book
- Get all borrowings
- Get borrowing detail

## Entities
### Books
- id
- title
- author
- isbn
- published_year
- stock
- created_at
- updated_at

### Members
- id
- name
- email
- phone
- created_at
- updated_at

### Borrowings
- id
- member_id
- book_id
- borrow_date
- due_date
- returned_at
- status
- created_at
- updated_at

## Business Rules
- Book can only be borrowed if stock > 0
- Borrowing reduces stock by 1
- Returning increases stock by 1
- Member must exist before borrowing
- Book must exist before borrowing
- Borrowing status can be: borrowed, returned
- Returned borrowing cannot be returned again

## Endpoints
### Utility
- GET /api/v1/health

### Books
- GET /api/v1/books
- GET /api/v1/books/:id
- POST /api/v1/books
- PUT /api/v1/books/:id
- DELETE /api/v1/books/:id

### Members
- GET /api/v1/members
- GET /api/v1/members/:id
- POST /api/v1/members
- PUT /api/v1/members/:id
- DELETE /api/v1/members/:id

### Borrowings
- GET /api/v1/borrowings
- GET /api/v1/borrowings/:id
- POST /api/v1/borrowings
- POST /api/v1/borrowings/:id/return

## Out of Scope
- Authentication
- Authorization
- Pagination
- Search and filtering
- Categories
- Fines and overdue penalties
- Reservation system
- Cover image upload
- Automated testing