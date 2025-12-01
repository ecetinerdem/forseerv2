# Forseer

Forseer is a modern portfolio management application built with a **Go backend** and a **React frontend**. It allows users to manage their stock portfolios, add/remove stocks, and monitor portfolio performance in real-time.

## Table of Contents

* [Tech Stack](#tech-stack)
* [Features](#features)
* [Architecture](#architecture)
* [Getting Started](#getting-started)
* [API Documentation](#api-documentation)
* [License](#license)

## Tech Stack

### Backend

* **Language:** Go
* **Framework:** Chi
* **Database:** NeonPostgres
* **Caching:** Redis
* **Email Service:** SendGrid
* **Cloud Hosting:** Google Cloud Platform (GCP)
* **API Documentation:** Swagger

### Frontend

* **Framework:** React
* **Bundler:** Vite
* **Language:** TypeScript
* **Styling:** TailwindCSS
* **Routing:** React Router DOM
* **Hosting:** Vercel

### Other Tools

* **Authentication & Authorization:** JWT tokens, Basic Auth
* **Rate Limiting:** Custom middleware
* **Logging:** Uber Zap
* **Env Management:** Custom env loader

## Features

* User registration and login with JWT authentication
* Portfolio creation, update, and deletion
* Add, update, and remove stocks from portfolios
* Portfolio search functionality
* Real-time caching with Redis
* Email notifications using SendGrid
* Swagger-based API documentation

## Architecture

```
Frontend (Vite + React + TypeScript)
       |
       v
Backend (Go + Chi)
       |
       v
Database (NeonPostgres)
       |
       +--> Cache (Redis)
       |
       +--> Mailer (SendGrid)
```

* **Frontend** is hosted on Vercel for global CDN delivery.
* **Backend** is deployed on Google Cloud Platform with Swagger docs accessible at `/v1/swagger/`.
* **Redis** is used for caching frequently accessed data.
* **SendGrid** handles transactional emails.
* **Chi** provides lightweight and fast HTTP routing and middleware support.
* **NeonPostgres** is the managed cloud database backend.

## Getting Started

### Prerequisites

* Node.js >= 20
* Go >= 1.20
* PostgreSQL (NeonPostgres) credentials
* Redis instance (optional but recommended)
* SendGrid API key

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

### Backend Setup

```bash
cd backend
go mod tidy
go run main.go
```

### Environment Variables

Create a `.env` file in the backend folder:

```env
PORT=8080
DB_URL=postgres://username:password@host:port/dbname
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
SENDGRID_API_KEY=your_sendgrid_api_key
JWT_SECRET=your_jwt_secret
CORS_ALLOWED_ORIGIN=http://localhost:5174
```

### Running the App

* Start backend server: `go run main.go`
* Start frontend: `npm run dev`
* Visit `http://localhost:5174` to use the app

## API Documentation

Swagger documentation is available at:

```
http://localhost:8080/v1/swagger/index.html
```

## License

MIT License
