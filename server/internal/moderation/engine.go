package moderation

import (
	"strings"
)

type EvaluationResult struct {
	ToxicityScore float64  `json:"toxicityScore"`
	FlagsRaised   []string `json:"flagsRaised"`
	Decision      string   `json:"decision"` // 'SAFE' | 'QUARANTINE' | 'REJECT'
}

var blacklistedSpamKeywords = []string{
	"buy cheap", "click here now", "free bitcoin", "guaranteed profit", "whatsapp number", "crypto bonus",
}

var blacklistedToxKeywords = []string{
	"kill", "violence", "bomb", "riot", "hate", "attack", "murder", "abuse",
}

// EvaluateText evaluates content text for toxicity and blacklisted keywords
func EvaluateText(title, body string) *EvaluationResult {
	combined := strings.ToLower(title + " " + body)
	var flags []string
	score := 0.0

	// Check spam patterns
	for _, word := range blacklistedSpamKeywords {
		if strings.Contains(combined, word) {
			flags = append(flags, "SPAM_KEYWORD:"+word)
			score += 0.35
		}
	}

	// Check toxicity/violence patterns
	for _, word := range blacklistedToxKeywords {
		if strings.Contains(combined, word) {
			flags = append(flags, "TOXIC_KEYWORD:"+word)
			score += 0.45
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	decision := "SAFE"
	if score >= 0.7 {
		decision = "REJECT"
	} else if score >= 0.3 {
		decision = "QUARANTINE"
	}

	return &EvaluationResult{
		ToxicityScore: score,
		FlagsRaised:   flags,
		Decision:      decision,
	}
}
