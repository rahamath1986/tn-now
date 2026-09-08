package moderation

import (
	"testing"
)

func TestEvaluateText(t *testing.T) {
	// 1. Clean Content
	res := EvaluateText("Madurai Cultural Fest 2026", "Annual festival celebrations scheduled for this weekend.")
	if res.Decision != "SAFE" {
		t.Errorf("Expected decision SAFE, got %s", res.Decision)
	}
	if res.ToxicityScore != 0.0 {
		t.Errorf("Expected score 0.0, got %f", res.ToxicityScore)
	}

	// 2. Spam Content
	resSpam := EvaluateText("Guaranteed Profit", "Buy cheap coins click here now!")
	if resSpam.Decision == "SAFE" {
		t.Error("Expected spam content decision to not be SAFE")
	}
	if len(resSpam.FlagsRaised) == 0 {
		t.Error("Expected flags raised for spam content")
	}

	// 3. Toxic Content
	resToxic := EvaluateText("Riot and Kill Attack", "Violence and murder threats")
	if resToxic.Decision != "REJECT" {
		t.Errorf("Expected decision REJECT for violent content, got %s", resToxic.Decision)
	}
}
