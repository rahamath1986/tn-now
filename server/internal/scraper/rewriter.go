package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// ArticleRewriteRequest contains article text and optional language metadata
type ArticleRewriteRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	District    string `json:"district,omitempty"`
	Category    string `json:"category,omitempty"`
	Language    string `json:"language,omitempty"`
}

// RewrittenArticle contains the cleanly rewritten journalistic content
type RewrittenArticle struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Lead        string `json:"lead,omitempty"`
	Language    string `json:"language"`
	Slug        string `json:"slug"`
	IsRewritten bool   `json:"isRewritten"`
	Engine      string `json:"engine"`
}

// ArticleRewriter handles high-fidelity journalistic rewriting with zero language contamination
type ArticleRewriter struct {
	apiKey     string
	httpClient *http.Client
}

var defaultGlobalRewriter *ArticleRewriter

// GetDefaultRewriter returns or initializes the shared rewriter instance
func GetDefaultRewriter() *ArticleRewriter {
	if defaultGlobalRewriter == nil {
		key := os.Getenv("GEMINI_API_KEY")
		defaultGlobalRewriter = NewArticleRewriter(key)
	}
	return defaultGlobalRewriter
}

// NewArticleRewriter creates a new article rewriter instance
func NewArticleRewriter(apiKey string) *ArticleRewriter {
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	return &ArticleRewriter{
		apiKey: strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

// SetAPIKey updates the Gemini API key dynamically
func (r *ArticleRewriter) SetAPIKey(key string) {
	r.apiKey = strings.TrimSpace(key)
}

// RewriteArticle rewrites an article with strict language preservation and zero Tanglish
func (r *ArticleRewriter) RewriteArticle(ctx context.Context, req ArticleRewriteRequest) (*RewrittenArticle, error) {
	cleanTitle := decodeAndCleanHTML(req.Title)
	cleanDesc := decodeAndCleanHTML(req.Description)

	if cleanTitle == "" && cleanDesc == "" {
		return nil, fmt.Errorf("article headline and content cannot both be empty")
	}
	if cleanTitle == "" {
		cleanTitle = cleanDesc
		if len([]rune(cleanTitle)) > 80 {
			cleanTitle = string([]rune(cleanTitle)[:80]) + "..."
		}
	}

	targetLang := strings.ToLower(strings.TrimSpace(req.Language))
	if targetLang == "" || targetLang == "auto" {
		targetLang = DetectLanguage(cleanTitle + " " + cleanDesc)
	}
	if strings.HasPrefix(targetLang, "ta") {
		targetLang = "ta"
	} else if strings.HasPrefix(targetLang, "en") {
		targetLang = "en"
	} else {
		targetLang = DetectLanguage(cleanTitle + " " + cleanDesc)
	}

	req.Title = cleanTitle
	req.Description = cleanDesc
	req.Language = targetLang

	apiKey := r.apiKey
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	if apiKey != "" {
		r.apiKey = apiKey
		candidateModels := []string{"gemini-flash-lite-latest", "gemini-3.5-flash", "gemini-flash-latest", "gemini-2.5-flash"}
		for _, model := range candidateModels {
			res, err := r.callGeminiAPI(ctx, model, req)
			if err == nil && res != nil {
				if validateLanguageIntegrity(targetLang, res.Title, res.Description) {
					res.Slug = GenerateArticleSlug(res.Title, req.ID)
					res.IsRewritten = true
					res.Engine = model
					return res, nil
				}
			}
		}
	}

	heuristicRes := r.HeuristicRewrite(req)
	heuristicRes.Slug = GenerateArticleSlug(heuristicRes.Title, req.ID)
	heuristicRes.IsRewritten = true
	heuristicRes.Engine = "heuristic-cleaner"
	return heuristicRes, nil
}

func (r *ArticleRewriter) callGeminiAPI(ctx context.Context, model string, req ArticleRewriteRequest) (*RewrittenArticle, error) {
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, r.apiKey)

	systemPrompt := buildSystemPrompt(req.Language)
	userContent := fmt.Sprintf("Original Headline: %s\nOriginal Article Content:\n%s\nDistrict: %s\nEditorial Category: %s",
		req.Title, req.Description, req.District, req.Category)

	payload := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": systemPrompt},
			},
		},
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": userContent},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":       0.2,
			"response_mime_type": "application/json",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var geminiRes struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiRes); err != nil {
		return nil, err
	}

	if len(geminiRes.Candidates) == 0 || len(geminiRes.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty candidates returned from gemini")
	}

	rawJSON := geminiRes.Candidates[0].Content.Parts[0].Text
	rawJSON = strings.TrimSpace(rawJSON)
	if strings.HasPrefix(rawJSON, "```json") {
		rawJSON = strings.TrimPrefix(rawJSON, "```json")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
	} else if strings.HasPrefix(rawJSON, "```") {
		rawJSON = strings.TrimPrefix(rawJSON, "```")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
	}
	rawJSON = strings.TrimSpace(rawJSON)

	var parsed struct {
		Title string `json:"title"`
		Lead  string `json:"lead"`
		Body  string `json:"body"`
	}
	if err := json.Unmarshal([]byte(rawJSON), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse gemini json output: %w", err)
	}

	parsed.Title = decodeAndCleanHTML(parsed.Title)
	parsed.Lead = decodeAndCleanHTML(parsed.Lead)
	parsed.Body = decodeAndCleanHTML(parsed.Body)

	if parsed.Title == "" {
		parsed.Title = req.Title
	}

	fullDesc := parsed.Body
	if fullDesc == "" {
		fullDesc = parsed.Lead
	}
	if fullDesc == "" {
		fullDesc = req.Description
	}

	return &RewrittenArticle{
		ID:          req.ID,
		Title:       parsed.Title,
		Description: fullDesc,
		Lead:        parsed.Lead,
		Language:    req.Language,
	}, nil
}

func buildSystemPrompt(lang string) string {
	if lang == "ta" {
		return `You are an elite Senior Chief News Editor (செய்தி ஆசிரியர்) at a leading Tamil media publication.
Your job is to rewrite scraped raw news articles into clean, high-journalistic Tamil (தூய செய்தி ஊடகத் தமிழ் / செந்தமிழ்).

STRICT MANDATORY RULES FOR ZERO LANGUAGE ISSUES:
1. ZERO TANGLISH (ABSOLUTELY NO LATIN/ENGLISH SCRIPT FOR TAMIL):
   - Never write Tamil phonetically in English (e.g., NEVER write "Chennai-la mazhai peithathu").
   - Every single word in "title", "lead", and "body" MUST be in authentic Unicode Tamil script.
2. STRICT LANGUAGE PRESERVATION:
   - This article is Tamil. Under NO circumstances should you translate it into English.
   - Keep the entire article 100% in Tamil.
3. ENTITY & FACT ACCURACY:
   - Preserve all proper nouns, political leaders, district names, government institutions, dates, percentages, amounts, and currency figures (e.g., ₹500 கோடி) with 100% precision.
   - NEVER invent, alter, or hallucinate facts.
4. REMOVE ALL BOILERPLATE & JUNK:
   - Strip out promotional text, "subscribe to our channel", "source: PTI", "follow us on twitter", and publisher links.
5. JOURNALISTIC STRUCTURE:
   - "title": Engaging, punchy, objective headline in pure Tamil without sensational clickbait.
   - "lead": 1-2 sentence crisp lead paragraph.
   - "body": 2-3 well-structured paragraphs in formal journalistic Tamil.

Respond ONLY with valid JSON in this exact structure:
{
  "title": "தூய தமிழில் தலைப்பு",
  "lead": "சுருக்கமான தொடக்க உரை",
  "body": "முழுமையான செய்தி விவரம் பத்திகளாக"
}`
	}

	return `You are an elite Senior Bureau News Editor at a reputable news wire (Reuters / AP style).
Your job is to rewrite scraped raw news articles into crisp, high-standard, professional, objective English journalism.

STRICT MANDATORY RULES FOR ZERO LANGUAGE ISSUES:
1. PURE ENGLISH PRESERVATION:
   - Every sentence in "title", "lead", and "body" MUST be in clean, grammatically flawless English.
   - Do NOT mix in Tanglish, Hindi, or untranslated snippets unless it is an official proper name.
2. ENTITY & FACT ACCURACY:
   - Preserve all personal names, official titles, districts, numbers, dates, currency amounts, and quotes exactly as reported.
   - NEVER hallucinate or alter factual context.
3. REMOVE ALL BOILERPLATE & JUNK:
   - Strip out promotional text, "subscribe to our newsletter", advertisement phrases, and social media calls.
4. JOURNALISTIC STRUCTURE:
   - "title": Punchy, engaging, objective headline in AP journalistic style.
   - "lead": 1-2 sentence informative lead.
   - "body": 2-3 well-structured, clear journalistic paragraphs.

Respond ONLY with valid JSON in this exact structure:
{
  "title": "Objective journalistic headline",
  "lead": "1-2 sentence lead summary",
  "body": "Full journalistic story body in clear paragraphs"
}`
}

func validateLanguageIntegrity(lang, title, body string) bool {
	combined := title + " " + body
	if strings.TrimSpace(combined) == "" {
		return false
	}

	if lang == "ta" {
		totalLetters := 0
		tamilChars := 0
		for _, r := range combined {
			if unicode.IsLetter(r) {
				totalLetters++
				if r >= 0x0B80 && r <= 0x0BFF {
					tamilChars++
				}
			}
		}
		if totalLetters == 0 {
			return false
		}
		ratio := float64(tamilChars) / float64(totalLetters)
		return ratio >= 0.40
	}

	totalLetters := 0
	latinLetters := 0
	for _, r := range combined {
		if unicode.IsLetter(r) {
			totalLetters++
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				latinLetters++
			}
		}
	}
	if totalLetters > 0 {
		ratio := float64(latinLetters) / float64(totalLetters)
		return ratio >= 0.70
	}

	return true
}

func (r *ArticleRewriter) HeuristicRewrite(req ArticleRewriteRequest) *RewrittenArticle {
	title := decodeAndCleanHTML(req.Title)
	desc := decodeAndCleanHTML(req.Description)

	title = stripBrandingSuffixes(title)
	desc = stripBoilerplate(desc)
	title = strings.TrimSpace(title)
	title = strings.TrimRight(title, ".,;-")

	lead := ""
	paragraphs := strings.Split(desc, "\n\n")
	if len(paragraphs) > 0 && strings.TrimSpace(paragraphs[0]) != "" {
		lead = strings.TrimSpace(paragraphs[0])
		if len([]rune(lead)) > 200 {
			lead = string([]rune(lead)[:200]) + "..."
		}
	}

	return &RewrittenArticle{
		ID:          req.ID,
		Title:       title,
		Description: desc,
		Lead:        lead,
		Language:    req.Language,
	}
}

func stripBrandingSuffixes(title string) string {
	re := regexp.MustCompile(`(?i)\s*[-|–—•]\s*(?:BBC News தமிழ்|BBC News|BBC|The Hindu|Dinamalar|தினமலர்|Puthiyathalaimurai|புதிய தலைமுறை|Maalai Malar|மாலை மலர்|News18|Oneindia|Tamil Oneindia|Vikatan|விகடன்|Daily Thanthi|தினத்தந்தி|ABP Nadu|News7 Tamil).*$`)
	title = re.ReplaceAllString(title, "")
	return strings.TrimSpace(title)
}

func stripBoilerplate(body string) string {
	lines := strings.Split(body, "\n")
	var cleanedLines []string

	junkPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)subscribe\s+to\s+our\s+channel`),
		regexp.MustCompile(`(?i)follow\s+us\s+on`),
		regexp.MustCompile(`(?i)also\s+read:`),
		regexp.MustCompile(`(?i)click\s+here\s+to\s+read`),
		regexp.MustCompile(`(?i)source:\s*(?:pti|ani|ians|reuters)`),
		regexp.MustCompile(`(?:தொடர்ந்து படிக்க|முழு செய்தியை படிக்க|எங்களை பின்தொடர|சந்தாதாரராக)`),
		regexp.MustCompile(`(?:செய்திகளை உடனுக்குடன் பெற|வாட்ஸ்அப் சேனலில் இணையவும்)`),
		regexp.MustCompile(`(?:பதிவிறக்கம் செய்ய இங்கே கிளிக் செய்யவும்)`),
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		isJunk := false
		for _, pattern := range junkPatterns {
			if pattern.MatchString(trimmed) {
				isJunk = true
				break
			}
		}
		if !isJunk {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	result := strings.Join(cleanedLines, "\n\n")
	reNewlines := regexp.MustCompile(`\n{3,}`)
	result = reNewlines.ReplaceAllString(result, "\n\n")
	return strings.TrimSpace(result)
}

func decodeAndCleanHTML(s string) string {
	if s == "" {
		return ""
	}
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "&#x27;", "'")
	s = strings.ReplaceAll(s, "&#x27", "'")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&#39", "'")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&ldquo;", "\"")
	s = strings.ReplaceAll(s, "&rdquo;", "\"")
	s = strings.ReplaceAll(s, "&lsquo;", "'")
	s = strings.ReplaceAll(s, "&rsquo;", "'")
	s = strings.ReplaceAll(s, "&ndash;", "-")
	s = strings.ReplaceAll(s, "&mdash;", "—")

	reTags := regexp.MustCompile(`<[^>]*>`)
	s = reTags.ReplaceAllString(s, " ")

	reSpaces := regexp.MustCompile(`[ \t\r\f]+`)
	s = reSpaces.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}
