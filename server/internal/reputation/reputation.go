package reputation

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"tn-now/server/internal/db"
)

type ContributorStats struct {
	UserID                 string     `json:"userId"`
	Level                  string     `json:"level"`
	Points                 int32      `json:"points"`
	TrustScore             int32      `json:"trustScore"`
	ApprovedCount          int32      `json:"approvedCount"`
	RejectedCount          int32      `json:"rejectedCount"`
	ReportCount            int32      `json:"reportCount"`
	IsFoundingContributor  bool       `json:"isFoundingContributor"`
	FoundingBadgeGrantedAt *time.Time `json:"foundingBadgeGrantedAt,omitempty"`
}

// CalculateLevel determines contributor level based on approved posts and trust score
func CalculateLevel(approvedCount, trustScore int32) string {
	if approvedCount >= 50 && trustScore >= 80 {
		return "Core Contributor"
	}
	if approvedCount >= 20 && trustScore >= 50 {
		return "Verified Contributor"
	}
	if approvedCount >= 5 && trustScore >= 20 {
		return "Trusted Contributor"
	}
	return "New Contributor"
}

// GetContributorStats retrieves contributor profile statistics
func GetContributorStats(ctx context.Context, q *db.Queries, userID string) (*ContributorStats, error) {
	var uID pgtype.UUID
	_ = uID.Scan(userID)

	profile, err := q.GetContributorProfile(ctx, uID)
	if err != nil {
		return nil, err
	}

	var badgeDate *time.Time
	if profile.FoundingBadgeGrantedAt.Valid {
		t := profile.FoundingBadgeGrantedAt.Time
		badgeDate = &t
	}

	return &ContributorStats{
		UserID:                 userID,
		Level:                  profile.Level,
		Points:                 profile.Points,
		TrustScore:             profile.TrustScore,
		ApprovedCount:          profile.ApprovedCount,
		RejectedCount:          profile.RejectedCount,
		ReportCount:            profile.ReportCount,
		IsFoundingContributor:  profile.IsFoundingContributor,
		FoundingBadgeGrantedAt: badgeDate,
	}, nil
}
