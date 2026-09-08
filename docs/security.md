# Security & API Protection Specification - TN NOW

This document details threat models and defensive configurations built into the TN NOW codebase.

---

## 1. Authentication Security

1. **Password Hashing**: Stored using Argon2id or bcrypt (minimum work factor = 12).
2. **Access Tokens**: Short-lived JWT access tokens (15-minute expiry).
3. **Refresh Tokens**: Long-lived refresh tokens stored securely. Rotation is enforced upon every refresh call to prevent token theft replay attacks.
4. **Flutter Client Storage**: Tokens are stored using Android Keystore / iOS Keychain wrapper packages (`flutter_secure_storage`).

---

## 2. Server-Side Input Validation

1. **SSRF Prevention on Link Submission**:
   - The Go backend resolves URLs strictly against a pre-registered whitelist (`ALLOWED_VIDEO_DOMAINS`):
     - `*.youtube.com`, `youtu.be`, `*.instagram.com`, `instagram.com`, `*.facebook.com`, `facebook.com`.
   - Domain resolving must verify the target IP. Loopback (`127.0.0.1`, `::1`), private ranges (`10.0.0.0/8`, `192.168.0.0/16`, `172.16.0.0/12`), and link-local ranges (`169.254.169.254`) are explicitly blocked to prevent Server-Side Request Forgery.
2. **File Security (Photo/Story Uploads)**:
   - File uploads bypass the Go API and go direct-to-storage via signed URLs to prevent resource starvation.
   - Post-upload callback checks verify headers, magic bytes (validating JPEG/PNG structure), and enforce a maximum file size constraint (5MB per photo).

---

## 3. Rate Limiting Rules (Redis-backed)

Rate limits are applied at the Go routing middleware layer:
- **Auth Endpoints (`/auth/login`, `/auth/register`)**: 5 requests per minute per IP.
- **Link Submit (`/video-links/submit`)**: 10 requests per minute per User ID.
- **Create Post (`/submissions`)**: 20 requests per hour per User ID.
- **Generic GET Endpoints**: 100 requests per minute per client IP.

---

## 4. Role-Based Access Control (RBAC)

Exclusively verified server-side. Users cannot inject roles.
- `USER`: Read-only guest and basic interactions (like, follow, bookmark).
- `CONTRIBUTOR`: Submissions access.
- `MODERATOR`: Admin queue management access.
- `ADMIN`: User management, audit logs, override rights.
- `SUPER_ADMIN`: Configuration changes, DB overrides.
