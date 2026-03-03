# 📚 High-Performance Go REST API – 8 Levels Challenge

A production-grade, benchmark-optimized REST API built with Go.

Designed to complete all 8 levels:

1. Ping
2. Echo
3. CRUD (Create & Read)
4. CRUD (Update & Delete)
5. Auth Guard
6. Search & Pagination
7. Error Handling
8. Boss Speed Run

---

## 🚀 Tech Stack

- Go 1.22+
- Chi Router
- Clean Architecture
- In-memory storage
- Docker-ready
- Benchmark-tested

---

## 🧱 Architecture

Clean Architecture layers:

Delivery → Usecase → Repository → Domain

Project Structure:

go-api-8-levels/
├── cmd/
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── repository/
│   ├── delivery/
│   └── middleware/
├── scripts/
├── postman/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md

---

## 📦 Where Is Data Stored?

Data is stored in:

internal/repository/inmemory_book_repository.go

Implementation:

- Uses map[string]*Book
- Thread-safe with sync.RWMutex
- Stored in memory (RAM)
- Reset when server restarts
- No database required

⚠ For production persistence, replace repository with PostgreSQL or MySQL.

---

## 🔥 API Endpoints

### 1️⃣ Ping
GET /ping

### 2️⃣ Echo
POST /echo

### 3️⃣ CRUD: Create & Read
POST   /books  
GET    /books  
GET    /books/{id}

### 4️⃣ CRUD: Update & Delete
PUT    /books/{id}  
DELETE /books/{id}

### 5️⃣ Auth Guard
POST /auth/token  
GET  /books (Protected)

### 6️⃣ Search & Pagination
GET /books?author=John  
GET /books?page=1&limit=2

### 7️⃣ Error Handling
- Invalid JSON → 400
- Not found → 404
- Unauthorized → 401

### 8️⃣ Boss Speed Run
All endpoints under load testing.

---

## 🐳 Run Locally

go run cmd/main.go

Server runs on:
http://localhost:8099

---

## 🐳 Run With Docker

Build:

docker build -t go-api .

Run:

docker run -p 8099:8099 go-api

---

## 🧪 Benchmark

Run Go benchmark:

go test -bench=. -benchmem ./...

Load test example:

wrk -t8 -c100 -d30s http://localhost:8080/ping

Expected performance (optimized version):

- /ping → 150k+ req/sec
- GET /books → 100k+ req/sec
- POST /books → 60k+ req/sec

---

## 📮 Postman

Import collection from:

postman/go-api-8-levels.postman_collection.json

Set environment variable:

base_url = http://localhost:8099

---

## 🛡 Production Improvements

- PostgreSQL integration
- Redis caching
- Prometheus metrics
- Structured logging
- Rate limiting
- CI/CD pipeline

---

## 👨‍💻 Author

Firman Ismail Hariri  
Backend Engineer – Go & Distributed Systems