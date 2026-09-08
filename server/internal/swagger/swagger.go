package swagger

import (
	"net/http"
)

const docJSON = `{
  "openapi": "3.0.3",
  "info": {
    "title": "TN NOW Platform API Specification",
    "description": "Hyper-local Tamil Nadu Real-time Discovery Platform Backend API Services",
    "version": "1.0.0"
  },
  "paths": {
    "/health/live": {
      "get": {
        "summary": "Liveness check",
        "responses": { "200": { "description": "Active" } }
      }
    },
    "/health/ready": {
      "get": {
        "summary": "Readiness check",
        "responses": { "200": { "description": "Database and Redis connection status" } }
      }
    },
    "/auth/register": {
      "post": {
        "summary": "Register a new user profile",
        "responses": { "201": { "description": "User created" } }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "Authenticate user & issue JWT tokens",
        "responses": { "200": { "description": "JWT tokens issued" } }
      }
    },
    "/home": {
      "get": {
        "summary": "Retrieve aggregated home feed updates",
        "responses": { "200": { "description": "Home feed updates payload" } }
      }
    },
    "/content": {
      "get": {
        "summary": "Retrieve paginated category feeds",
        "responses": { "200": { "description": "Paginated content items list" } }
      }
    },
    "/video/resolve": {
      "post": {
        "summary": "Parse video URL and return embed metadata",
        "responses": { "200": { "description": "Video embed metadata" } }
      }
    },
    "/content/submit/video": {
      "post": {
        "summary": "Submit a video link post",
        "responses": { "201": { "description": "Submitted for moderation" } }
      }
    },
    "/moderation/queue": {
      "get": {
        "summary": "Retrieve pending moderation queue",
        "responses": { "200": { "description": "Pending queue items" } }
      }
    },
    "/contributor/profile": {
      "get": {
        "summary": "Retrieve contributor reputation profile",
        "responses": { "200": { "description": "Contributor stats and badges" } }
      }
    },
    "/social/like": {
      "post": {
        "summary": "Like or unlike content",
        "responses": { "200": { "description": "Like state updated" } }
      }
    },
    "/compliance/grievance": {
      "post": {
        "summary": "Register grievance complaint under IT Rules 2021",
        "responses": { "201": { "description": "Grievance ticket created" } }
      }
    },
    "/search": {
      "get": {
        "summary": "Search content updates across Tamil Nadu",
        "responses": { "200": { "description": "Matching content results" } }
      }
    }
  }
}`

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(docJSON))
	})
}
