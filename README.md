# TN NOW

> Everything happening around Tamil Nadu, right now.

**TN NOW** is a real-time discovery platform for Tamil Nadu, focusing on local viral videos, trending updates, native posts, text stories, local events, and legal compliance. Launching initially in **Madurai** and scaling across all 38 Tamil Nadu districts.

---

## 1. Core Architecture & Tech Stack

- **Mobile Client**: Flutter (Material 3 Dark Mode, Riverpod state management, GoRouter navigation, Dio client with JWT token auto-refresh).
- **Backend API**: Go (Native `http.ServeMux`, `pgx/v5` PostgreSQL driver, `sqlc` compiled queries, Redis sliding window rate limiter, HMAC SHA-256 JWT auth).
- **Database & Cache**: PostgreSQL 16 (with `tsvector` full-text search) & Redis 7 (rate limiting & token blacklists).
- **Media & Storage**: Cloudflare R2 / MinIO (Docker emulated S3 storage).
- **Containerization**: Docker & Docker Compose (`docker-compose.yml`).

---

## 2. Platform Feature Overview (Phases 0 – 15 Complete)

- **Phase 0 (Product Specifications)**: Architecture, DB schema, security, moderation, legal compliance, and roadmap specs under `docs/`.
- **Phase 1 (Repository Foundation)**: Flutter app under `mobile/` and Go server under `server/` with API, Worker, and Migration CLI entry points.
- **Phase 2 (Flutter Foundation)**: Theme system (M3 dark mode), Dio network client with token interceptors, Riverpod state management, and GoRouter route registry.
- **Phase 3 (Go Foundation)**: RequestID, slog JSON logger, panic recovery, CORS, Redis sliding-window rate limiter, and OpenAPI doc endpoints.
- **Phase 4 (Database Layer)**: PostgreSQL schema migrations (`users`, `profiles`, `content`, `video_links`, `photos`, `stories`, `events`, `grievances`, `audit_logs`, `moderation_results`, `likes`, `bookmarks`, `follows`, `comments`).
- **Phase 5 (Authentication)**: Register, Login, Token Refresh (one-time rotation), and Logout endpoints with bcrypt password hashing and HMAC SHA-256 JWTs.
- **Phase 6 (Home Feed)**: Location-aware updates, category carousel, trending posts, and upcoming events preview (`GET /home`).
- **Phase 7 (Category Feed)**: Infinite scrolling cursor-paginated content feeds (`GET /content`).
- **Phase 8 (Video / Embed Engine)**: Embed resolver for YouTube, Instagram Reels, Facebook, and X (Twitter) URLs (`POST /video/resolve`) with background dead-link checker worker (`cmd/worker`).
- **Phase 9 (Submission Flow)**: Multi-format user posting forms (`VIDEO_LINK`, `PHOTO`, `TEXT_STORY`, `EVENT`) with moderation intake queues.
- **Phase 10 (Moderation Core)**: Automated text toxicity evaluator, Tamil/English keyword blacklists, state machine transitions (`SAFE`, `QUARANTINE`, `REJECT`), and admin queue review endpoints.
- **Phase 11 (Trust & Reputation)**: Contributor progression levels (`New`, `Trusted`, `Verified`, `Core`), trust score progress meters, and golden `Founding Contributor 2026` badges.
- **Phase 12 (Social & Engagement)**: Likes, Bookmarks, Follows, Comments thread modal sheet, and optimistic UI counters.
- **Phase 13 (Grievances & Compliance)**: Indian IT Rules 2021 compliant grievance redressal intake (`GRV-YYYYMMDD-XXXX` ticket ref with 24h receipt ACK and 15d resolution SLA) and legal takedown handlers.
- **Phase 14 (Search & Discovery)**: Full-text search engine (`GET /search`) with district filter chips (Madurai, Chennai, Coimbatore, Trichy, Salem).
- **Phase 15 (Hardening & Docker)**: HTTP security response headers (`nosniff`, `DENY`), OpenAPI 3.0 specs (`/swagger/doc.json`), multi-container Docker Compose setup, and 100% clean test audits.

---

## 3. Quick Start

### Prerequisites
- [Flutter SDK](https://docs.flutter.dev/get-started/install)
- [Go 1.22+](https://go.dev/doc/install)
- [Docker & Docker Compose](https://docs.docker.com/desktop/)

### Launch Infrastructure & API Server
```bash
# Start Docker infrastructure (PostgreSQL, Redis, MinIO)
make docker-up

# Run database migrations
make migrate

# Start Go backend server
make dev-server
```

### Launch Mobile Application
```bash 
cd mobile
flutter run
```

### Run Full Test Suite
```bash
# Backend Go Tests
cd server && go test -v ./...

# Mobile Flutter Tests
cd mobile && flutter test
```
