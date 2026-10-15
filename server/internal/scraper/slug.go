package scraper

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// commonTamilWordReplacements maps frequent Tamil news keywords to their standard Romanized forms
// for clean, high-CTR, SEO-optimized URL slugs.
var commonTamilWordReplacements = []struct {
	tamil string
	roman string
}{
	{"தமிழ்நாடு", "tamilnadu"},
	{"சென்னை", "chennai"},
	{"மதுரை", "madurai"},
	{"கோவை", "coimbatore"},
	{"கோயம்புத்தூர்", "coimbatore"},
	{"திருச்சி", "trichy"},
	{"திருச்சிராப்பள்ளி", "tiruchirappalli"},
	{"சேலம்", "salem"},
	{"திருநெல்வேலி", "tirunelveli"},
	{"தூத்துக்குடி", "thoothukudi"},
	{"தஞ்சாவூர்", "thanjavur"},
	{"வேலூர்", "vellore"},
	{"ஈரோடு", "erode"},
	{"திண்டுக்கல்", "dindigul"},
	{"கன்னியாகுமரி", "kanyakumari"},
	{"முதலமைச்சர்", "cm"},
	{"முதல்வர்", "cm"},
	{"ஸ்டாலின்", "stalin"},
	{"அரசாங்கம்", "govt"},
	{"அரசு", "arasu"},
	{"செய்திகள்", "seithigal"},
	{"செய்தி", "seithi"},
	{"அறிவிப்பு", "arivippu"},
	{"மழை", "mazhai"},
	{"கனமழை", "kanamazhai"},
	{"வெள்ளம்", "vellam"},
	{"புயல்", "puyal"},
	{"பள்ளிகளுக்கு", "schools"},
	{"பள்ளிகள்", "schools"},
	{"பள்ளிக்கு", "school"},
	{"பள்ளி", "school"},
	{"கல்லூரிகளுக்கு", "colleges"},
	{"கல்லூரிகள்", "colleges"},
	{"கல்லூரிக்கு", "college"},
	{"கல்லூரி", "college"},
	{"மாணவர்களுக்கு", "students"},
	{"மாணவர்கள்", "students"},
	{"விடுமுறை", "holiday"},
	{"தங்கம்", "gold"},
	{"விலை", "rate"},
	{"பெட்ரோல்", "petrol"},
	{"டீசல்", "diesel"},
	{"பட்ஜெட்", "budget"},
	{"தேர்தல்", "election"},
	{"வாக்குப்பதிவு", "voting"},
	{"காவல்துறை", "police"},
	{"விபத்து", "accident"},
	{"போக்குவரத்து", "traffic"},
	{"ரயில்", "train"},
	{"மெட்ரோ", "metro"},
	{"விமானம்", "flight"},
	{"விளையாட்டு", "sports"},
	{"சினிமா", "cinema"},
	{"அரசியல்", "politics"},
	{"வணிகம்", "business"},
}

func init() {
	sort.SliceStable(commonTamilWordReplacements, func(i, j int) bool {
		return len([]rune(commonTamilWordReplacements[i].tamil)) > len([]rune(commonTamilWordReplacements[j].tamil))
	})
}

// TransliterateTamilToLatin converts Tamil Unicode text to phonetic English/Latin characters.
func TransliterateTamilToLatin(input string) string {
	var sb strings.Builder
	runes := []rune(input)
	n := len(runes)

	for i := 0; i < n; i++ {
		r := runes[i]

		// Independent vowels
		switch r {
		case '\u0B85': // அ
			sb.WriteString("a")
			continue
		case '\u0B86': // ஆ
			sb.WriteString("aa")
			continue
		case '\u0B87': // இ
			sb.WriteString("i")
			continue
		case '\u0B88': // ஈ
			sb.WriteString("ee")
			continue
		case '\u0B89': // உ
			sb.WriteString("u")
			continue
		case '\u0B8A': // ஊ
			sb.WriteString("oo")
			continue
		case '\u0B8E': // எ
			sb.WriteString("e")
			continue
		case '\u0B8F': // ஏ
			sb.WriteString("ae")
			continue
		case '\u0B90': // ஐ
			sb.WriteString("ai")
			continue
		case '\u0B92': // ஒ
			sb.WriteString("o")
			continue
		case '\u0B93': // ஓ
			sb.WriteString("oa")
			continue
		case '\u0B94': // ஔ
			sb.WriteString("au")
			continue
		case '\u0B83': // ஃ
			sb.WriteString("k")
			continue
		}

		// Consonants
		var baseConsonant string
		switch r {
		case '\u0B95': // க
			baseConsonant = "k"
		case '\u0B99': // ங
			baseConsonant = "ng"
		case '\u0B9A': // ச
			baseConsonant = "s"
		case '\u0B9E': // ஞ
			baseConsonant = "gn"
		case '\u0B9F': // ட
			baseConsonant = "t"
		case '\u0BA3': // ண
			baseConsonant = "n"
		case '\u0BA4': // த
			baseConsonant = "th"
		case '\u0BA8': // ந
			baseConsonant = "n"
		case '\u0BAA': // ப
			baseConsonant = "p"
		case '\u0BAE': // ம
			baseConsonant = "m"
		case '\u0BAF': // ய
			baseConsonant = "y"
		case '\u0BB0': // ர
			baseConsonant = "r"
		case '\u0BB2': // ல
			baseConsonant = "l"
		case '\u0BB5': // வ
			baseConsonant = "v"
		case '\u0BB4': // ழ
			baseConsonant = "zh"
		case '\u0BB3': // ள
			baseConsonant = "l"
		case '\u0BB1': // ற
			baseConsonant = "r"
		case '\u0BA9': // ன
			baseConsonant = "n"
		case '\u0B9C': // ஜ
			baseConsonant = "j"
		case '\u0BB7': // ஷ
			baseConsonant = "sh"
		case '\u0BB8': // ஸ
			baseConsonant = "s"
		case '\u0BB9': // ஹ
			baseConsonant = "h"
		}

		if baseConsonant != "" {
			// Check if next rune is a vowel sign or pulli
			if i+1 < n {
				nextR := runes[i+1]
				switch nextR {
				case '\u0BCD': // ் (virama / pulli - mutes inherent vowel)
					sb.WriteString(baseConsonant)
					i++ // consume pulli
					continue
				case '\u0BBE': // ா
					sb.WriteString(baseConsonant + "aa")
					i++
					continue
				case '\u0BBF': // ி
					sb.WriteString(baseConsonant + "i")
					i++
					continue
				case '\u0BC0': // ீ
					sb.WriteString(baseConsonant + "ee")
					i++
					continue
				case '\u0BC1': // ு
					sb.WriteString(baseConsonant + "u")
					i++
					continue
				case '\u0BC2': // ூ
					sb.WriteString(baseConsonant + "oo")
					i++
					continue
				case '\u0BC6': // ெ
					sb.WriteString(baseConsonant + "e")
					i++
					continue
				case '\u0BC7': // ே
					sb.WriteString(baseConsonant + "ae")
					i++
					continue
				case '\u0BC8': // ை
					sb.WriteString(baseConsonant + "ai")
					i++
					continue
				case '\u0BCA': // ொ
					sb.WriteString(baseConsonant + "o")
					i++
					continue
				case '\u0BCB': // ோ
					sb.WriteString(baseConsonant + "oa")
					i++
					continue
				case '\u0BCC': // ௌ
					sb.WriteString(baseConsonant + "au")
					i++
					continue
				}
			}
			// Standalone consonant has inherent 'a' sound
			sb.WriteString(baseConsonant + "a")
			continue
		}

		// Keep ASCII letters, digits, and spaces
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			sb.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			sb.WriteRune(' ')
		}
	}

	return sb.String()
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
var multiHyphenRegex = regexp.MustCompile(`-+`)

// ExtractShortID extracts the first 8 characters of a UUID or ID string.
func ExtractShortID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	parts := strings.Split(id, "-")
	if len(parts) > 0 && len(parts[0]) >= 6 {
		return parts[0]
	}
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

// GenerateArticleSlug converts an article title and unique ID into a clean,
// SEO-optimized, human-readable ASCII slug (e.g. "chennai-rain-update-holiday-3f82a1bc").
func GenerateArticleSlug(title, id string) string {
	cleaned := CleanHTML(title)

	// Step 1: Apply common high-frequency Tamil news keywords
	for _, rep := range commonTamilWordReplacements {
		cleaned = strings.ReplaceAll(cleaned, rep.tamil, rep.roman)
	}

	// Step 2: Transliterate any remaining Tamil characters to phonetic Latin
	transliterated := TransliterateTamilToLatin(cleaned)

	// Step 3: Replace all non-alphanumeric characters with hyphens
	slug := nonAlphanumericRegex.ReplaceAllString(strings.ToLower(transliterated), "-")
	slug = multiHyphenRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	// Step 4: Cap text portion at 60 characters at a word boundary
	if len(slug) > 60 {
		trimmed := slug[:60]
		if lastHyphen := strings.LastIndex(trimmed, "-"); lastHyphen > 25 {
			slug = trimmed[:lastHyphen]
		} else {
			slug = trimmed
		}
	}

	shortID := ExtractShortID(id)
	if slug == "" {
		if shortID != "" {
			return "news-" + shortID
		}
		return "news-article"
	}

	if shortID != "" {
		return slug + "-" + shortID
	}

	return slug
}
