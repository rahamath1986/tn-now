# API Specification - TN NOW

All API endpoints are prefixed with `/api/v1`.

---

## 1. Global Response and Error Formats

### Standard Success Response

For single objects or actions:
```json
{
  "success": true,
  "data": {
    "id": "b9683c3e-8fbf-4a39-ae7f-fb18dcf92a9c",
    "title": "Viral Video in Madurai"
  },
  "message": "Action completed successfully",
  "errors": []
}
```

### Standard Paginated List Response

For lists of items, cursor pagination is mandatory to prevent performance degradation as datasets scale.
```json
{
  "success": true,
  "data": [
    {
      "id": "b9683c3e-8fbf-4a39-ae7f-fb18dcf92a9c",
      "title": "Viral Video in Madurai"
    }
  ],
  "pagination": {
    "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNi0wOC0zMFQxMDozMjozMVoifQ==",
    "hasMore": true
  }
}
```

### Standard Error Response

```json
{
  "success": false,
  "data": null,
  "message": "Validation failed",
  "errors": [
    {
      "code": "INVALID_URL",
      "field": "url",
      "message": "URL is not from a supported platform"
    }
  ]
}
```

---

## 2. API Endpoints Reference

### 2.1 Authentication & Profile
- `POST /auth/register` - Create a new user profile.
- `POST /auth/login` - Authenticate with email/password and obtain JWT.
- `POST /auth/refresh` - Rotate expired JWT token using refresh tokens.
- `POST /auth/logout` - Revoke current session.
- `GET /contributors/me` - Retrieve current contributor stats and points.

### 2.2 Home & Feed Discovery
- `GET /home` - Retrieve unified homepage feed (Trending, Latest, Events).
- `GET /trending` - Paginated trending content (with time-decay scores).
- `GET /nearby` - Distance radius/district queries (requires coordinates or location parameter).

### 2.3 Submissions
- `POST /video-links/submit` - Validate allow-list and resolve oEmbed metadata for review card.
- `POST /submissions` - Submit Photo Post, Text Story, Event, or finalized Video Link.
- `POST /uploads/initiate` - Secure signed Cloudflare R2 / MinIO URL for direct asset upload.
- `POST /uploads/complete` - Finalize file upload metadata registration.

### 2.4 Interactions
- `POST /content/{id}/like` - Toggle like.
- `POST /content/{id}/bookmark` - Toggle bookmark.
- `POST /content/{id}/comments` - Create comment.
- `GET /content/{id}/comments` - Retrieve comments.
- `POST /content/{id}/report` - Report content violation.

### 2.5 Compliance
- `POST /grievances` - Submit grievance case under IT Rules 2021.
- `GET /grievances/{id}` - Track SLA resolution status.
