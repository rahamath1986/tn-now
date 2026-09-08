package compliance

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
)

type Handler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewHandler(q *db.Queries, conn *pgxpool.Pool) *Handler {
	return &Handler{q: q, conn: conn}
}

type grievanceRequest struct {
	ComplainantEmail  string `json:"complainantEmail"`
	Category          string `json:"category"` // 'COPYRIGHT' | 'DEFAMATION' | 'PRIVACY_VIOLATION' | 'OBJECTIONABLE_CONTENT'
	Description       string `json:"description"`
	ContentID         string `json:"contentId,omitempty"`
	ComplainantUserID string `json:"complainantUserId,omitempty"`
}

type takedownRequest struct {
	ContentID   string `json:"contentId"`
	OrderRef    string `json:"orderRef"`
	Reason      string `json:"reason"`
	AdminUserID string `json:"adminUserId"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/compliance/grievance", h.HandleSubmitGrievance)
	mux.HandleFunc("/compliance/takedown", h.HandleLegalTakedown)
}

func (h *Handler) HandleSubmitGrievance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req grievanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if strings.TrimSpace(req.ComplainantEmail) == "" || strings.TrimSpace(req.Category) == "" || strings.TrimSpace(req.Description) == "" {
		writeError(w, http.StatusUnprocessableEntity, "ComplainantEmail, Category, and Description are required")
		return
	}

	ctx := r.Context()
	var cID pgtype.UUID
	if req.ContentID != "" {
		_ = cID.Scan(req.ContentID)
	}

	var uID pgtype.UUID
	if req.ComplainantUserID != "" {
		_ = uID.Scan(req.ComplainantUserID)
	}

	// Insert Grievance Record
	grievance, err := h.q.CreateGrievance(ctx, db.CreateGrievanceParams{
		ComplainantUserID: uID,
		ContactEmail:      req.ComplainantEmail,
		Category:          req.Category,
		Description:       req.Description,
		ContentID:         cID,
	})
	if err != nil {
		slog.Error("Failed to register grievance", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to register grievance")
		return
	}

	ticketRef := fmt.Sprintf("GRV-%s-%04d", time.Now().Format("20060102"), rand.Intn(10000))

	writeSuccess(w, http.StatusCreated, map[string]string{
		"grievanceId":   uuidToString(grievance.ID),
		"ticketRef":     ticketRef,
		"ackSla":        "24 Hours",
		"resolutionSla": "15 Days",
	}, "Grievance registered under Indian IT Rules 2021")
}

func (h *Handler) HandleLegalTakedown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req takedownRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ContentID == "" || req.OrderRef == "" || req.AdminUserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "ContentID, OrderRef, and AdminUserID are required")
		return
	}

	ctx := r.Context()
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback(ctx)

	qtx := h.q.WithTx(tx)
	cID := stringToUUID(req.ContentID)
	adminID := stringToUUID(req.AdminUserID)

	// Update Content Status to TAKEDOWN
	_, err = qtx.UpdateContentStatus(ctx, db.UpdateContentStatusParams{
		ID:               cID,
		Status:           "TAKEDOWN",
		ModerationStatus: "REJECTED",
		PublishedAt:      pgtype.Timestamptz{Valid: false},
	})
	if err != nil {
		slog.Error("Failed to update content status for legal takedown", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to perform takedown")
		return
	}

	// Record Audit Trail
	details := fmt.Sprintf("Legal Order Ref: %s, Reason: %s", req.OrderRef, req.Reason)
	_, err = qtx.CreateAuditLog(ctx, db.CreateAuditLogParams{
		ActorUserID:  adminID,
		Action:       "LEGAL_TAKEDOWN",
		TargetEntity: "content",
		TargetID:     cID,
		Details:      pgtype.Text{String: details, Valid: true},
	})
	if err != nil {
		slog.Error("Failed to create audit trail for takedown", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to write audit trail")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction commit failed")
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{
		"contentId": req.ContentID,
		"orderRef":  req.OrderRef,
		"status":    "TAKEDOWN",
	}, "Content successfully taken down pursuant to legal order")
}

func stringToUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	_ = u.Scan(s)
	return u
}

func uuidToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:])
}

func writeSuccess(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"message": message,
		"errors":  []interface{}{},
	})
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"data":    nil,
		"message": message,
		"errors":  []interface{}{},
	})
}
