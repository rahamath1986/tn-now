package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ScrapedItem represents an extracted news or video piece
type ScrapedItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SourceURL   string    `json:"sourceUrl"`
	District    string    `json:"district"`
	Category    string    `json:"category"`
	ContentType string    `json:"contentType"` // 'TEXT_STORY' | 'VIDEO_LINK' | 'PHOTO'
	ImageURLs   []string  `json:"imageUrls"`
	VideoURL    string    `json:"videoUrl"`
	VideoID     string    `json:"videoId"`
	VideoType   string    `json:"videoType"` // 'youtube' | 'direct' | 'none'
	Status      string    `json:"status"`    // 'PENDING'
	PublishedAt time.Time `json:"publishedAt"`
}

// ScrapeResult summarizes the ingestion run
type ScrapeResult struct {
	SourceURL   string         `json:"sourceUrl"`
	TotalFound  int            `json:"totalFound"`
	StagedCount int            `json:"stagedCount"`
	DuplicateCount int         `json:"duplicateCount"`
	Items       []*ScrapedItem `json:"items"`
}

// XML RSS representation
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL  string `xml:"url,attr"`
		Type string `xml:"type,attr"`
	} `xml:"enclosure"`
	MediaContent struct {
		URL    string `xml:"url,attr"`
		Medium string `xml:"medium,attr"`
	} `xml:"content"`
	MediaThumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
}

// Atom feed representation (e.g. YouTube channels)
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID        string `xml:"id"`
	VideoID   string `xml:"videoId"`
	Title     string `xml:"title"`
	Link      struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Published string `xml:"published"`
	MediaGroup struct {
		Description string `xml:"description"`
		Thumbnail   struct {
			URL string `xml:"url,attr"`
		} `xml:"thumbnail"`
	} `xml:"group"`
}

// All 40 official administrative districts of Tamil Nadu
var tnDistricts = []string{
	"Ariyalur", "Chengalpattu", "Chennai", "Coimbatore", "Cuddalore",
	"Dharmapuri", "Dindigul", "Erode", "Kallakurichi", "Kanchipuram",
	"Kanyakumari", "Karur", "Krishnagiri", "Madurai", "Mayiladuthurai",
	"Nagapattinam", "Namakkal", "Nilgiris", "Perambalur", "Pudukkottai",
	"Ramanathapuram", "Ranipet", "Salem", "Sivaganga", "Tenkasi",
	"Thanjavur", "Theni", "Thoothukudi", "Tiruchirappalli", "Tirunelveli",
	"Tirupattur", "Tiruppur", "Tiruvallur", "Tiruvannamalai", "Tiruvarur",
	"Vellore", "Viluppuram", "Virudhunagar", "Trichy", "Tuticorin",
}

// Tamil script to official English district mapping
var tnTamilDistrictMap = map[string]string{
	"சென்னை":          "Chennai",
	"மதுரை":           "Madurai",
	"கோவை":           "Coimbatore",
	"கோயம்புத்தூர்":     "Coimbatore",
	"திருச்சி":          "Tiruchirappalli",
	"திருச்சிராப்பள்ளி": "Tiruchirappalli",
	"சேலம்":           "Salem",
	"தஞ்சாவூர்":        "Thanjavur",
	"தஞ்சை":           "Thanjavur",
	"நெல்லை":          "Tirunelveli",
	"திருநெல்வேலி":     "Tirunelveli",
	"கன்னியாகுமரி":     "Kanyakumari",
	"திண்டுக்கல்":      "Dindigul",
	"ஈரோடு":           "Erode",
	"வேலூர்":           "Vellore",
	"திருப்பூர்":         "Tiruppur",
	"தூத்துக்குடி":      "Thoothukudi",
	"கடலூர்":          "Cuddalore",
	"காஞ்சிபுரம்":      "Kanchipuram",
	"புதுக்கோட்டை":     "Pudukkottai",
	"பெரம்பலூர்":       "Perambalur",
	"சிவகங்கை":         "Sivaganga",
	"ராமநாதபுரம்":      "Ramanathapuram",
	"நாகப்பட்டினம்":     "Nagapattinam",
	"ராணிப்பேட்டை":     "Ranipet",
	"திருப்பத்தூர்":      "Tirupattur",
	"திருவண்ணாமலை":     "Tiruvannamalai",
	"திருவாரூர்":        "Tiruvarur",
	"கள்ளக்குறிச்சி":    "Kallakurichi",
	"தருமபுரி":         "Dharmapuri",
	"கிருஷ்ணகிரி":       "Krishnagiri",
	"தேனி":            "Theni",
	"தென்காசி":         "Tenkasi",
	"நாமக்கல்":         "Namakkal",
	"நீலகிரி":          "Nilgiris",
	"அரியலூர்":         "Ariyalur",
	"செங்கல்பட்டு":      "Chengalpattu",
	"திருவள்ளூர்":       "Tiruvallur",
	"விழுப்புரம்":       "Viluppuram",
	"விருதுநகர்":        "Virudhunagar",
	"மயிலாடுதுறை":      "Mayiladuthurai",
}

// Curated high-resolution landmark/landscape imagery for each Tamil Nadu district and events
// Curated official geographic locator maps for each Tamil Nadu district, Tamil Nadu State, India (National), and World (International)
var districtMapImageMap = map[string]string{
	"Tamil Nadu":      "https://upload.wikimedia.org/wikipedia/commons/thumb/2/2c/Tamil_Nadu_districts_map.svg/960px-Tamil_Nadu_districts_map.svg.png",
	"National":        "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b3/India_map_with_states_and_union_territories.svg/960px-India_map_with_states_and_union_territories.svg.png",
	"International":   "https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/World_map_-_low_resolution.svg/960px-World_map_-_low_resolution.svg.png",
	"Madurai":         "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Madurai_in_Tamil_Nadu_%28India%29.svg/960px-Madurai_in_Tamil_Nadu_%28India%29.svg.png",
	"Chennai":         "https://upload.wikimedia.org/wikipedia/commons/thumb/4/48/Chennai_in_Tamil_Nadu_%28India%29.svg/960px-Chennai_in_Tamil_Nadu_%28India%29.svg.png",
	"Coimbatore":      "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Coimbatore_in_Tamil_Nadu_%28India%29.svg/960px-Coimbatore_in_Tamil_Nadu_%28India%29.svg.png",
	"Salem":           "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ee/Salem_in_Tamil_Nadu_%28India%29.svg/960px-Salem_in_Tamil_Nadu_%28India%29.svg.png",
	"Thanjavur":       "https://upload.wikimedia.org/wikipedia/commons/thumb/7/78/Thanjavur_in_Tamil_Nadu_%28India%29.svg/960px-Thanjavur_in_Tamil_Nadu_%28India%29.svg.png",
	"Kanyakumari":     "https://upload.wikimedia.org/wikipedia/commons/thumb/5/5e/Kanyakumari_in_Tamil_Nadu_%28India%29.svg/960px-Kanyakumari_in_Tamil_Nadu_%28India%29.svg.png",
	"Vellore":         "https://upload.wikimedia.org/wikipedia/commons/thumb/4/44/Vellore_in_Tamil_Nadu_%28India%29.svg/960px-Vellore_in_Tamil_Nadu_%28India%29.svg.png",
	"Tiruchirappalli": "https://upload.wikimedia.org/wikipedia/commons/thumb/4/49/Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg/960px-Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg.png",
	"Trichy":          "https://upload.wikimedia.org/wikipedia/commons/thumb/4/49/Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg/960px-Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg.png",
	"Dindigul":        "https://upload.wikimedia.org/wikipedia/commons/thumb/5/52/Dindigul_in_Tamil_Nadu_%28India%29.svg/960px-Dindigul_in_Tamil_Nadu_%28India%29.svg.png",
	"Erode":           "https://upload.wikimedia.org/wikipedia/commons/thumb/6/60/Erode_in_Tamil_Nadu_%28India%29.svg/960px-Erode_in_Tamil_Nadu_%28India%29.svg.png",
	"Tirunelveli":     "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Tirunelveli_in_Tamil_Nadu_%28India%29.svg/960px-Tirunelveli_in_Tamil_Nadu_%28India%29.svg.png",
	"Thoothukudi":     "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Thoothukudi_in_Tamil_Nadu_%28India%29.svg/960px-Thoothukudi_in_Tamil_Nadu_%28India%29.svg.png",
	"Tuticorin":       "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Thoothukudi_in_Tamil_Nadu_%28India%29.svg/960px-Thoothukudi_in_Tamil_Nadu_%28India%29.svg.png",
	"Tiruppur":        "https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Tiruppur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruppur_in_Tamil_Nadu_%28India%29.svg.png",
	"Nilgiris":        "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6c/Nilgiris_in_Tamil_Nadu_%28India%29.svg/960px-Nilgiris_in_Tamil_Nadu_%28India%29.svg.png",
	"Cuddalore":       "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Cuddalore_in_Tamil_Nadu_%28India%29.svg/960px-Cuddalore_in_Tamil_Nadu_%28India%29.svg.png",
	"Kanchipuram":     "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Kanchipuram_in_Tamil_Nadu_%28India%29.svg/960px-Kanchipuram_in_Tamil_Nadu_%28India%29.svg.png",
	"Chengalpattu":    "https://upload.wikimedia.org/wikipedia/commons/thumb/1/18/Chengalpattu_in_Tamil_Nadu_%28India%29.svg/960px-Chengalpattu_in_Tamil_Nadu_%28India%29.svg.png",
	"Tiruvallur":      "https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Tiruvallur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvallur_in_Tamil_Nadu_%28India%29.svg.png",
	"Pudukkottai":     "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d3/Pudukkottai_in_Tamil_Nadu_%28India%29.svg/960px-Pudukkottai_in_Tamil_Nadu_%28India%29.svg.png",
	"Sivaganga":       "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Sivaganga_in_Tamil_Nadu_%28India%29.svg/960px-Sivaganga_in_Tamil_Nadu_%28India%29.svg.png",
	"Ramanathapuram":  "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Ramanathapuram_in_Tamil_Nadu_%28India%29.svg/960px-Ramanathapuram_in_Tamil_Nadu_%28India%29.svg.png",
	"Nagapattinam":    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Nagapattinam_in_Tamil_Nadu_%28India%29.svg/960px-Nagapattinam_in_Tamil_Nadu_%28India%29.svg.png",
	"Ranipet":         "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/Ranipet_in_Tamil_Nadu_%28India%29.svg/960px-Ranipet_in_Tamil_Nadu_%28India%29.svg.png",
	"Tirupattur":      "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/Tirupattur_in_Tamil_Nadu_%28India%29.svg/960px-Tirupattur_in_Tamil_Nadu_%28India%29.svg.png",
	"Tiruvannamalai":  "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Tiruvannamalai_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvannamalai_in_Tamil_Nadu_%28India%29.svg.png",
	"Tiruvarur":       "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cb/Tiruvarur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvarur_in_Tamil_Nadu_%28India%29.svg.png",
	"Kallakurichi":    "https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Kallakurichi_in_Tamil_Nadu_%28India%29.svg/960px-Kallakurichi_in_Tamil_Nadu_%28India%29.svg.png",
	"Dharmapuri":      "https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Dharmapuri_in_Tamil_Nadu_%28India%29.svg/960px-Dharmapuri_in_Tamil_Nadu_%28India%29.svg.png",
	"Krishnagiri":     "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4a/Krishnagiri_in_Tamil_Nadu_%28India%29.svg/960px-Krishnagiri_in_Tamil_Nadu_%28India%29.svg.png",
	"Theni":           "https://upload.wikimedia.org/wikipedia/commons/thumb/5/50/Theni_in_Tamil_Nadu_%28India%29.svg/960px-Theni_in_Tamil_Nadu_%28India%29.svg.png",
	"Tenkasi":         "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Tenkasi_in_Tamil_Nadu_%28India%29.svg/960px-Tenkasi_in_Tamil_Nadu_%28India%29.svg.png",
	"Namakkal":        "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c5/Namakkal_in_Tamil_Nadu_%28India%29.svg/960px-Namakkal_in_Tamil_Nadu_%28India%29.svg.png",
	"Ariyalur":        "https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Ariyalur_in_Tamil_Nadu_%28India%29.svg/960px-Ariyalur_in_Tamil_Nadu_%28India%29.svg.png",
	"Perambalur":      "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e5/Perambalur_in_Tamil_Nadu_%28India%29.svg/960px-Perambalur_in_Tamil_Nadu_%28India%29.svg.png",
	"Viluppuram":      "https://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/Viluppuram_in_Tamil_Nadu_%28India%29.svg/960px-Viluppuram_in_Tamil_Nadu_%28India%29.svg.png",
	"Virudhunagar":    "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Virudhunagar_in_Tamil_Nadu_%28India%29.svg/960px-Virudhunagar_in_Tamil_Nadu_%28India%29.svg.png",
	"Mayiladuthurai":  "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/Mayiladuthurai_in_Tamil_Nadu_%28India%29.svg/960px-Mayiladuthurai_in_Tamil_Nadu_%28India%29.svg.png",
	"Karur":           "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Karur_in_Tamil_Nadu_%28India%29.svg/960px-Karur_in_Tamil_Nadu_%28India%29.svg.png",
}

func getFallbackImage(district, category string) string {
	dClean := strings.TrimSpace(district)
	if img, ok := districtMapImageMap[dClean]; ok && img != "" {
		return img
	}
	for k, v := range districtMapImageMap {
		if strings.EqualFold(k, dClean) {
			return v
		}
	}
	return districtMapImageMap["Tamil Nadu"]
}

func GetFallbackDistrictImageExported(district, category string) string {
	return getFallbackImage(district, category)
}

// Curated high-resolution official Wikimedia portraits for key public personalities, leaders, and sports icons
var personImageMap = map[string]string{
	// Leaders & Politicians - Tamil Nadu
	"ஸ்டாலின்":           "https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg",
	"மு.க.ஸ்டாலின்":       "https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg",
	"m.k. stalin":        "https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg",
	"stalin":             "https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg",

	"எடப்பாடி":          "https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg",
	"பழனிசாமி":           "https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg",
	"edappadi":           "https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg",
	"palaniswami":        "https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg",
	"eps":                "https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg",

	"விஜய்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg",
	"தளபதி விஜய்":       "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg",
	"தவெக":             "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg",
	"vijay":             "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg",
	"tvk":               "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg",

	"உதயநிதி":           "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Udhayanidhi_Stalin.jpg/800px-Udhayanidhi_Stalin.jpg",
	"udhayanidhi":        "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Udhayanidhi_Stalin.jpg/800px-Udhayanidhi_Stalin.jpg",

	"அண்ணாமலை":          "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/K._Annamalai_in_2023.jpg/800px-K._Annamalai_in_2023.jpg",
	"annamalai":          "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/K._Annamalai_in_2023.jpg/800px-K._Annamalai_in_2023.jpg",

	"சீமான்":            "https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Seeman_at_a_meeting.jpg/800px-Seeman_at_a_meeting.jpg",
	"seeman":            "https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Seeman_at_a_meeting.jpg/800px-Seeman_at_a_meeting.jpg",

	"பன்னீர்செல்வம்":       "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/O._Panneerselvam.jpg/800px-O._Panneerselvam.jpg",
	"ops":                "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/O._Panneerselvam.jpg/800px-O._Panneerselvam.jpg",
	"panneerselvam":      "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/O._Panneerselvam.jpg/800px-O._Panneerselvam.jpg",

	"தினகரன்":           "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/T._T._V._Dhinakaran.jpg/800px-T._T._V._Dhinakaran.jpg",
	"dhinakaran":         "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/T._T._V._Dhinakaran.jpg/800px-T._T._V._Dhinakaran.jpg",

	// Celebrities & Cinema
	"கமல்":              "https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Kamal_Haasan_in_2022.jpg/800px-Kamal_Haasan_in_2022.jpg",
	"கமல்ஹாசன்":         "https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Kamal_Haasan_in_2022.jpg/800px-Kamal_Haasan_in_2022.jpg",
	"kamal haasan":       "https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Kamal_Haasan_in_2022.jpg/800px-Kamal_Haasan_in_2022.jpg",

	"ரஜினி":             "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Rajinikanth_in_2023.jpg/800px-Rajinikanth_in_2023.jpg",
	"ரஜினிகாந்த்":         "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Rajinikanth_in_2023.jpg/800px-Rajinikanth_in_2023.jpg",
	"rajinikanth":        "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Rajinikanth_in_2023.jpg/800px-Rajinikanth_in_2023.jpg",

	"அஜித்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Ajith_Kumar.jpg/800px-Ajith_Kumar.jpg",
	"ajith":              "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Ajith_Kumar.jpg/800px-Ajith_Kumar.jpg",

	"சூர்யா":            "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Suriya_at_2D_Entertainment_office.jpg/800px-Suriya_at_2D_Entertainment_office.jpg",
	"suriya":             "https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Suriya_at_2D_Entertainment_office.jpg/800px-Suriya_at_2D_Entertainment_office.jpg",

	"விக்ரம்":            "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Vikram_at_Cobra_promotions.jpg/800px-Vikram_at_Cobra_promotions.jpg",
	"vikram":             "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Vikram_at_Cobra_promotions.jpg/800px-Vikram_at_Cobra_promotions.jpg",

	"தனுஷ்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Dhanush_at_The_Gray_Man_press_conference.jpg/800px-Dhanush_at_The_Gray_Man_press_conference.jpg",
	"dhanush":            "https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Dhanush_at_The_Gray_Man_press_conference.jpg/800px-Dhanush_at_The_Gray_Man_press_conference.jpg",

	"சிவகார்த்திகேயன்":    "https://upload.wikimedia.org/wikipedia/commons/thumb/3/30/Sivakarthikeyan_at_Doctor_success_meet.jpg/800px-Sivakarthikeyan_at_Doctor_success_meet.jpg",
	"sivakarthikeyan":    "https://upload.wikimedia.org/wikipedia/commons/thumb/3/30/Sivakarthikeyan_at_Doctor_success_meet.jpg/800px-Sivakarthikeyan_at_Doctor_success_meet.jpg",

	"நயன்தாரா":          "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Nayanthara_at_SIIMA_2016.jpg/800px-Nayanthara_at_SIIMA_2016.jpg",
	"nayanthara":         "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Nayanthara_at_SIIMA_2016.jpg/800px-Nayanthara_at_SIIMA_2016.jpg",

	"திரிஷா":            "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fa/Trisha_Krishnan_at_PS2_Promotions.jpg/800px-Trisha_Krishnan_at_PS2_Promotions.jpg",
	"trisha":             "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fa/Trisha_Krishnan_at_PS2_Promotions.jpg/800px-Trisha_Krishnan_at_PS2_Promotions.jpg",

	"ரஹ்மான்":           "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg",
	"ஏ.ஆர்.ரஹ்மான்":      "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg",
	"ar rahman":          "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg",
	"a.r. rahman":        "https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg",

	"அனிருத்":           "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Anirudh_Ravichander_at_Vikram_audio_launch.jpg/800px-Anirudh_Ravichander_at_Vikram_audio_launch.jpg",
	"anirudh":            "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Anirudh_Ravichander_at_Vikram_audio_launch.jpg/800px-Anirudh_Ravichander_at_Vikram_audio_launch.jpg",

	// National Leaders & Personalities
	"மோடி":              "https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg",
	"நரேந்திர மோடி":      "https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg",
	"modi":               "https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg",
	"narendra modi":      "https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg",

	"ராகுல்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Rahul_Gandhi_in_2023.jpg/800px-Rahul_Gandhi_in_2023.jpg",
	"ராகுல் காந்தி":      "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Rahul_Gandhi_in_2023.jpg/800px-Rahul_Gandhi_in_2023.jpg",
	"rahul gandhi":       "https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Rahul_Gandhi_in_2023.jpg/800px-Rahul_Gandhi_in_2023.jpg",

	"அமித் ஷா":          "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Amit_Shah_in_2023.jpg/800px-Amit_Shah_in_2023.jpg",
	"அமித்ஷா":           "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Amit_Shah_in_2023.jpg/800px-Amit_Shah_in_2023.jpg",
	"amit shah":          "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Amit_Shah_in_2023.jpg/800px-Amit_Shah_in_2023.jpg",

	"மம்தா":             "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg",
	"மம்தா பானர்ஜி":     "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg",
	"mamata":             "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg",
	"mamata banerjee":    "https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg",

	"நிர்மலா சீதாராமன்": "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9f/Nirmala_Sitharaman_in_2022.jpg/800px-Nirmala_Sitharaman_in_2022.jpg",
	"nirmala sitharaman": "https://upload.wikimedia.org/wikipedia/commons/thumb/9/9f/Nirmala_Sitharaman_in_2022.jpg/800px-Nirmala_Sitharaman_in_2022.jpg",

	"மாணிக்கம் தாகூர்":    "https://upload.wikimedia.org/wikipedia/commons/thumb/2/22/Manickam_Tagore.jpg/800px-Manickam_Tagore.jpg",
	"manickam tagore":    "https://upload.wikimedia.org/wikipedia/commons/thumb/2/22/Manickam_Tagore.jpg/800px-Manickam_Tagore.jpg",

	"திருமாவளவன்":        "https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Thol._Thirumavalavan.jpg/800px-Thol._Thirumavalavan.jpg",
	"thirumavalavan":     "https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Thol._Thirumavalavan.jpg/800px-Thol._Thirumavalavan.jpg",

	"கனிமொழி":           "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ec/Kanimozhi_Karunanidhi.jpg/800px-Kanimozhi_Karunanidhi.jpg",
	"kanimozhi":          "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ec/Kanimozhi_Karunanidhi.jpg/800px-Kanimozhi_Karunanidhi.jpg",

	"பவன் கல்யாண்":       "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ab/Pawan_Kalyan_in_2024.jpg/800px-Pawan_Kalyan_in_2024.jpg",
	"pawan kalyan":       "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ab/Pawan_Kalyan_in_2024.jpg/800px-Pawan_Kalyan_in_2024.jpg",

	"திரௌபதி முர்மு":     "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/Droupadi_Murmu_official_portrait.jpg/800px-Droupadi_Murmu_official_portrait.jpg",
	"droupadi murmu":     "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e6/Droupadi_Murmu_official_portrait.jpg/800px-Droupadi_Murmu_official_portrait.jpg",

	// Sports Personalities - Football & Global
	"மெஸ்ஸி":            "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg",
	"லியோனல் மெஸ்ஸி":    "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg",
	"messi":              "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg",
	"lionel messi":       "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg",

	"ரொனால்டோ":          "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg",
	"கிறிஸ்டியானோ ரொனால்டோ": "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg",
	"ronaldo":            "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg",
	"cristiano ronaldo":  "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg",

	"நெய்மர்":           "https://upload.wikimedia.org/wikipedia/commons/thumb/8/83/ISL_Final_2019-20_%28Neymar_cropped%29.jpg/800px-ISL_Final_2019-20_%28Neymar_cropped%29.jpg",
	"neymar":             "https://upload.wikimedia.org/wikipedia/commons/thumb/8/83/ISL_Final_2019-20_%28Neymar_cropped%29.jpg/800px-ISL_Final_2019-20_%28Neymar_cropped%29.jpg",

	// Sports Personalities - Cricket & Athletics
	"தோனி":              "https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg",
	"எம்.எஸ். தோனி":      "https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg",
	"dhoni":              "https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg",
	"m.s. dhoni":         "https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg",

	"கோலி":              "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg",
	"விராட் கோலி":        "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg",
	"kohli":              "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg",
	"virat kohli":        "https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg",

	"ரோஹித்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg",
	"ரோகித் சர்மா":       "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg",
	"rohit":              "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg",
	"rohit sharma":       "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg",

	"அஸ்வின்":            "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Ravichandran_Ashwin.jpg/800px-Ravichandran_Ashwin.jpg",
	"ashwin":             "https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Ravichandran_Ashwin.jpg/800px-Ravichandran_Ashwin.jpg",

	"சச்சின்":           "https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Sachin_Tendulkar_at_MRF_Promotion_Event.jpg/800px-Sachin_Tendulkar_at_MRF_Promotion_Event.jpg",
	"sachin":             "https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Sachin_Tendulkar_at_MRF_Promotion_Event.jpg/800px-Sachin_Tendulkar_at_MRF_Promotion_Event.jpg",

	// Global Leaders
	"டிரம்ப்":            "https://upload.wikimedia.org/wikipedia/commons/thumb/5/56/Donald_Trump_official_portrait.jpg/800px-Donald_Trump_official_portrait.jpg",
	"டொனால்ட் டிரம்ப்":     "https://upload.wikimedia.org/wikipedia/commons/thumb/5/56/Donald_Trump_official_portrait.jpg/800px-Donald_Trump_official_portrait.jpg",
	"trump":              "https://upload.wikimedia.org/wikipedia/commons/thumb/5/56/Donald_Trump_official_portrait.jpg/800px-Donald_Trump_official_portrait.jpg",

	"புதின்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8d/Vladimir_Putin_%282020-02-20%29.jpg/800px-Vladimir_Putin_%282020-02-20%29.jpg",
	"putin":              "https://upload.wikimedia.org/wikipedia/commons/thumb/8/8d/Vladimir_Putin_%282020-02-20%29.jpg/800px-Vladimir_Putin_%282020-02-20%29.jpg",

	"பைடன்":             "https://upload.wikimedia.org/wikipedia/commons/thumb/6/68/Joe_Biden_presidential_portrait.jpg/800px-Joe_Biden_presidential_portrait.jpg",
	"biden":              "https://upload.wikimedia.org/wikipedia/commons/thumb/6/68/Joe_Biden_presidential_portrait.jpg/800px-Joe_Biden_presidential_portrait.jpg",
}

// Curated high-resolution photography for topical news categories when no specific article image exists
var categoryImageMap = map[string]string{
	"sports":        "https://images.unsplash.com/photo-1508098682722-e99c43a406b2?w=1200&auto=format&fit=crop&q=80",
	"politics":      "https://images.unsplash.com/photo-1541872703-74c5e44368f9?w=1200&auto=format&fit=crop&q=80",
	"business":      "https://images.unsplash.com/photo-1590283603385-17ffb3a7f29f?w=1200&auto=format&fit=crop&q=80",
	"technical":     "https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&auto=format&fit=crop&q=80",
	"tech":          "https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&auto=format&fit=crop&q=80",
	"entertainment": "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=1200&auto=format&fit=crop&q=80",
	"cinema":        "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=1200&auto=format&fit=crop&q=80",
	"crime":         "https://images.unsplash.com/photo-1589829545856-d10d557cf95f?w=1200&auto=format&fit=crop&q=80",
	"civic":         "https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?w=1200&auto=format&fit=crop&q=80",
	"news":          "https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=1200&auto=format&fit=crop&q=80",
}

func detectCategoryImage(category, title string) string {
	catLower := strings.ToLower(category)
	for k, img := range categoryImageMap {
		if strings.Contains(catLower, k) {
			return img
		}
	}
	tLower := strings.ToLower(title)
	if strings.Contains(tLower, "விளையாட்டு") || strings.Contains(tLower, "கிரிக்கெட்") || strings.Contains(tLower, "கால்பந்து") ||
		strings.Contains(tLower, "cricket") || strings.Contains(tLower, "football") || strings.Contains(tLower, "sports") ||
		strings.Contains(tLower, "ipl") || strings.Contains(tLower, "fifa") || strings.Contains(tLower, "olympic") ||
		strings.Contains(tLower, "goal") || strings.Contains(tLower, "trophy") || strings.Contains(tLower, "match") {
		return categoryImageMap["sports"]
	}
	if strings.Contains(tLower, "அரசியல்") || strings.Contains(tLower, "தேர்தல்") || strings.Contains(tLower, "அரசு") ||
		strings.Contains(tLower, "முதல்வர்") || strings.Contains(tLower, "அமைச்சர்") || strings.Contains(tLower, "politics") ||
		strings.Contains(tLower, "parliament") || strings.Contains(tLower, "assembly") || strings.Contains(tLower, "minister") {
		return categoryImageMap["politics"]
	}
	if strings.Contains(tLower, "திரைப்பட") || strings.Contains(tLower, "சினிமா") || strings.Contains(tLower, "நடிக") ||
		strings.Contains(tLower, "cinema") || strings.Contains(tLower, "movie") || strings.Contains(tLower, "trailer") ||
		strings.Contains(tLower, "box office") {
		return categoryImageMap["entertainment"]
	}
	if strings.Contains(tLower, "வணிகம்") || strings.Contains(tLower, "பொருளாதார") || strings.Contains(tLower, "பங்குச்சந்தை") ||
		strings.Contains(tLower, "ரூபாய்") || strings.Contains(tLower, "gold") || strings.Contains(tLower, "business") ||
		strings.Contains(tLower, "market") || strings.Contains(tLower, "sensex") {
		return categoryImageMap["business"]
	}
	if strings.Contains(tLower, "தொழில்நுட்ப") || strings.Contains(tLower, "இஸ்ரோ") || strings.Contains(tLower, "விண்கலம்") ||
		strings.Contains(tLower, "tech") || strings.Contains(tLower, "ai") || strings.Contains(tLower, "software") ||
		strings.Contains(tLower, "space") || strings.Contains(tLower, "isro") {
		return categoryImageMap["technical"]
	}
	if strings.Contains(tLower, "குற்றம்") || strings.Contains(tLower, "கைது") || strings.Contains(tLower, "போலீஸ்") ||
		strings.Contains(tLower, "crime") || strings.Contains(tLower, "police") || strings.Contains(tLower, "arrest") ||
		strings.Contains(tLower, "court") {
		return categoryImageMap["crime"]
	}
	return ""
}

func detectPersonImage(title, body string) string {
	titleLower := strings.ToLower(title)
	for kw, img := range personImageMap {
		if strings.Contains(titleLower, strings.ToLower(kw)) {
			return img
		}
	}
	bodySample := strings.ToLower(body)
	if len(bodySample) > 400 {
		bodySample = bodySample[:400]
	}
	for kw, img := range personImageMap {
		if strings.Contains(bodySample, strings.ToLower(kw)) {
			return img
		}
	}
	return ""
}

// GetFallbackImageWithPerson resolves fallback photography hierarchically:
// 1. Person mentioned in article headline/text (e.g. Lionel Messi, Stalin, Vijay, Dhoni)
// 2. Category photography dynamically (e.g. Sports stadium, Politics, Business, Cinema)
// 3. Geographic District SVG / State fallback
func GetFallbackImageWithPerson(title, district, category string) string {
	if pImg := detectPersonImage(title, ""); pImg != "" {
		return pImg
	}
	if cImg := detectCategoryImage(category, title); cImg != "" {
		return cImg
	}
	return getFallbackImage(district, category)
}

var tamilScriptRegex = regexp.MustCompile(`[\x{0B80}-\x{0BFF}]`)

// DetectLanguage returns 'ta' for Tamil script content or 'en' for English content
func DetectLanguage(text string) string {
	if tamilScriptRegex.MatchString(text) {
		return "ta"
	}
	return "en"
}

func DetectDistrictExported(text string) string {
	return detectDistrict(text)
}

func DetectCategoryExported(text, urlStr string) string {
	return detectCategory(text, urlStr)
}

func isViralContent(title, description, urlStr, contentType string) bool {
	if contentType == "VIDEO_LINK" || strings.Contains(urlStr, "youtube.com") || strings.Contains(urlStr, "youtu.be") {
		return true
	}
	combined := strings.ToLower(title + " " + description + " " + urlStr)
	viralKws := []string{
		"வைரல்", "பரபரப்பு", "பகீர்", "அதிர்ச்சி", "ஷாக்", "டிரெண்டிங்", "காணொளி", "வீடியோ",
		"viral", "trending", "shocking", "video", "spotted", "breaking", "buzz", "controversy", "sensational",
	}
	for _, kw := range viralKws {
		if strings.Contains(combined, kw) {
			return true
		}
	}
	return false
}

func IsViralContentExported(title, description, urlStr, contentType string) bool {
	return isViralContent(title, description, urlStr, contentType)
}

// DBMu guarantees thread-safe serialization for PostgreSQL queries across concurrent HTTP requests
var DBMu sync.Mutex

// ScrapeAndStage fetches content from siteURL, parses text/images/video, and stages into PostgreSQL with status='PENDING'
func ScrapeAndStage(ctx context.Context, conn *pgxpool.Pool, targetURL string) (*ScrapeResult, error) {
	if strings.TrimSpace(targetURL) == "" {
		return nil, fmt.Errorf("target URL cannot be empty")
	}

	slog.Info("Scraper starting fetch", slog.String("url", targetURL))

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   25 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ta,en-US;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d fetching %s", resp.StatusCode, targetURL)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	bodyStr := string(bodyBytes)

	var items []*ScrapedItem

	// 1. Try parsing as RSS XML
	if strings.Contains(contentType, "xml") || strings.Contains(bodyStr, "<rss") {
		var rss rssFeed
		if err := xml.Unmarshal(bodyBytes, &rss); err == nil && len(rss.Channel.Items) > 0 {
			slog.Info("Parsed RSS channel items", slog.Int("count", len(rss.Channel.Items)), slog.String("url", targetURL))
			for _, it := range rss.Channel.Items {
				scraped := parseRSSItem(it)
				if scraped != nil {
					items = append(items, scraped)
				}
			}
		}
	}

	// 2. Try parsing as Atom XML (YouTube channel or blog feed)
	if len(items) == 0 && (strings.Contains(contentType, "xml") || strings.Contains(bodyStr, "<feed")) {
		var atom atomFeed
		if err := xml.Unmarshal(bodyBytes, &atom); err == nil && len(atom.Entries) > 0 {
			slog.Info("Parsed Atom feed entries", slog.Int("count", len(atom.Entries)), slog.String("url", targetURL))
			for _, e := range atom.Entries {
				scraped := parseAtomEntry(e)
				if scraped != nil {
					items = append(items, scraped)
				}
			}
		}
	}

	// 3. Fallback: Parse HTML webpage for news articles, og tags, and video embeds
	if len(items) == 0 {
		if isSingleArticle(targetURL, bodyStr) {
			single := parseHTMLPage(targetURL, bodyStr)
			if single != nil {
				items = append(items, single)
			}
			parsedItems := parseHTMLPageMultiple(targetURL, bodyStr)
			for _, pit := range parsedItems {
				if pit.SourceURL != targetURL && (single == nil || pit.Title != single.Title) {
					items = append(items, pit)
				}
			}
		} else {
			parsedItems := parseHTMLPageMultiple(targetURL, bodyStr)
			if len(parsedItems) > 0 {
				slog.Info("Parsed HTML page multiple stories", slog.Int("count", len(parsedItems)), slog.String("url", targetURL))
				items = append(items, parsedItems...)
			} else {
				single := parseHTMLPage(targetURL, bodyStr)
				if single != nil {
					items = append(items, single)
				}
			}
		}

		// Crawl navigation links matching all sections of added sources
		navLinks := extractMatchingNavLinks(targetURL, bodyStr)
		if len(navLinks) == 0 && isSingleArticle(targetURL, bodyStr) {
			// If scraping a single article directly, also discover sections from source base URL
			if u, err := url.Parse(targetURL); err == nil {
				baseRoot := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				if strings.Contains(u.Path, "/tamil") {
					baseRoot += "/tamil"
				}
				if rootReq, err := http.NewRequestWithContext(ctx, http.MethodGet, baseRoot, nil); err == nil {
					rootReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
					if rootResp, err := client.Do(rootReq); err == nil {
						if rootBytes, err := io.ReadAll(rootResp.Body); err == nil && len(rootBytes) > 0 {
							navLinks = extractMatchingNavLinks(baseRoot, string(rootBytes))
						}
						rootResp.Body.Close()
					}
				}
			}
		}

		if len(navLinks) > 0 {
			slog.Info("Crawling matching section & district nav links", slog.Int("count", len(navLinks)), slog.String("url", targetURL))
			for _, navURL := range navLinks {
				if len(items) >= 40 {
					break
				}
				subReq, err := http.NewRequestWithContext(ctx, http.MethodGet, navURL, nil)
				if err != nil {
					continue
				}
				subReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
				subResp, err := client.Do(subReq)
				if err != nil {
					continue
				}
				subBytes, err := io.ReadAll(subResp.Body)
				subResp.Body.Close()
				if err == nil && len(subBytes) > 0 {
					subItems := parseHTMLPageMultiple(navURL, string(subBytes))
					for _, sit := range subItems {
						if len(items) >= 40 {
							break
						}
						isDup := false
						for _, itm := range items {
							if itm.SourceURL == sit.SourceURL || itm.Title == sit.Title {
								isDup = true
								break
							}
						}
						if !isDup {
							items = append(items, sit)
						}
					}
				}
			}
		}
	}

	slog.Info("Total scraped items extracted", slog.Int("count", len(items)), slog.String("url", targetURL))

	// Stage extracted items into PostgreSQL content table with status = 'PENDING'
	result := &ScrapeResult{
		SourceURL:  targetURL,
		TotalFound: len(items),
		Items:      items,
	}

	if conn == nil {
		result.StagedCount = len(items)
		return result, nil
	}

	// 1. Enrich scraped items in parallel (cap at 25 items per run to ensure fast execution under 5s)
	// NO DATABASE LOCK HELD during network I/O so other HTTP requests are never blocked!
	if len(items) > 25 {
		items = items[:25]
	}

	enrichItem := func(item *ScrapedItem) {
		// Fetch full article webpage to extract complete unabridged text body, all images, and any video
		if strings.HasPrefix(item.SourceURL, "http") && !strings.Contains(item.SourceURL, "youtube.com") && !strings.Contains(item.SourceURL, "youtu.be") {
			fullTxt, articleImgs, subVidID, subVidURL, subVidType, pubDate := FetchFullTextAndMediaExported(item.SourceURL)
			if item.PublishedAt.IsZero() && !pubDate.IsZero() {
				item.PublishedAt = pubDate
			}
			if len(fullTxt) > len(item.Description) || len(item.Description) < 300 {
				if fullTxt != "" {
					item.Description = fullTxt
				}
			}
			// Add extracted article images, prioritizing real images over fallback SVG map
			if len(articleImgs) > 0 {
				var newImages []string
				seenImg := make(map[string]bool)
				for _, aimg := range articleImgs {
					if !seenImg[aimg] {
						seenImg[aimg] = true
						newImages = append(newImages, aimg)
					}
				}
				for _, ex := range item.ImageURLs {
					if !strings.Contains(ex, "/maps/svg") && !seenImg[ex] {
						seenImg[ex] = true
						newImages = append(newImages, ex)
					}
				}
				item.ImageURLs = newImages
			}
			// If article has video, add video too
			if (item.VideoID == "" || item.VideoURL == "") && subVidURL != "" {
				item.VideoID = subVidID
				item.VideoURL = subVidURL
				item.VideoType = subVidType
				item.ContentType = "VIDEO_LINK"
			}
		}

		// Re-classify District & Category using full unabridged article content
		item.District = detectDistrict(item.Title + " " + item.Description + " " + item.SourceURL)
		item.Category = detectCategory(item.Title + " " + item.Description, item.SourceURL)

		// If no image available, try dynamically with person or category specified in article
		if len(item.ImageURLs) == 0 {
			if personImg := detectPersonImage(item.Title, item.Description); personImg != "" {
				item.ImageURLs = append(item.ImageURLs, personImg)
			} else if catImg := detectCategoryImage(item.Category, item.Title+" "+item.Description); catImg != "" {
				item.ImageURLs = append(item.ImageURLs, catImg)
			}
		}

		if item.PublishedAt.IsZero() {
			if pt, ok := ParsePublishedTime(item.Description); ok {
				item.PublishedAt = pt
			}
		}
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 6) // Max 6 concurrent network fetches
	for _, itm := range items {
		wg.Add(1)
		go func(it *ScrapedItem) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			enrichItem(it)
		}(itm)
	}
	wg.Wait()

	// 2. Fast Database Staging: lock DBMu ONLY for database inserts (takes < 25ms total)
	DBMu.Lock()
	defer DBMu.Unlock()

	dbCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ensure required columns exist
	_, _ = conn.Exec(dbCtx, `
		ALTER TABLE content ADD COLUMN IF NOT EXISTS is_viral BOOLEAN DEFAULT FALSE;
		ALTER TABLE content ADD COLUMN IF NOT EXISTS language VARCHAR(20) DEFAULT 'ta';
	`)

	// Ensure system author exists
	var authorID string
	_ = conn.QueryRow(dbCtx, "SELECT id FROM users LIMIT 1").Scan(&authorID)
	if authorID == "" {
		authorID = "00000000-0000-0000-0000-000000000001"
		_, _ = conn.Exec(dbCtx, `
			INSERT INTO users (id, username, email, password_hash, role)
			VALUES ($1, 'scraper_system', 'scraper@tnnow.in', 'system_hash', 'ADMIN')
			ON CONFLICT (id) DO NOTHING
		`, authorID)
	}

	for _, item := range items {
		trimmedTitle := strings.TrimSpace(item.Title)
		if len([]rune(trimmedTitle)) < 10 || strings.Contains(trimmedTitle, ".css-") || strings.Contains(trimmedTitle, "{") {
			continue
		}

		// Skip ancient archival posts older than 48 hours (e.g. historical links from 2022/2023)
		if !item.PublishedAt.IsZero() && item.PublishedAt.Before(time.Now().Add(-48*time.Hour)) {
			continue
		}

		// Deduplication check: canonical URL, normalized title, or external video ID
		normTitle := NormalizeTitle(item.Title)
		canonURL := CanonicalizeURL(item.SourceURL)

		var existingID string
		_ = conn.QueryRow(dbCtx, `
			SELECT c.id FROM content c
			LEFT JOIN video_links vl ON c.id = vl.content_id
			WHERE c.source_url = $1
			   OR (c.title = $2 AND $2 != '')
			   OR (vl.external_video_id = $3 AND $3 != '')
			   OR (c.source_url = $4 AND $4 != '')
			   OR (LENGTH($5) >= 10 AND LOWER(REGEXP_REPLACE(c.title, '[\.,!?:;''"()\[\]/|\-–—\s]+', ' ', 'g')) = $5)
			LIMIT 1
		`, item.SourceURL, item.Title, item.VideoID, canonURL, normTitle).Scan(&existingID)

		if existingID != "" {
			result.DuplicateCount++
			item.ID = existingID
			continue
		}

		// Resolve District UUID (fallback to 'Tamil Nadu', never default to Madurai)
		var districtID string
		_ = conn.QueryRow(dbCtx, "SELECT id FROM districts WHERE name ILIKE $1 LIMIT 1", item.District).Scan(&districtID)
		if districtID == "" {
			_ = conn.QueryRow(dbCtx, "SELECT id FROM districts WHERE name = 'Tamil Nadu' LIMIT 1").Scan(&districtID)
		}
		if districtID == "" {
			_ = conn.QueryRow(dbCtx, "SELECT id FROM districts ORDER BY name LIMIT 1").Scan(&districtID)
		}

		// Resolve Category UUID
		var categoryID string
		_ = conn.QueryRow(dbCtx, "SELECT id FROM categories WHERE name ILIKE $1 LIMIT 1", "%"+item.Category+"%").Scan(&categoryID)
		if categoryID == "" {
			_ = conn.QueryRow(dbCtx, "SELECT id FROM categories WHERE name = 'News' LIMIT 1").Scan(&categoryID)
		}
		if categoryID == "" {
			_ = conn.QueryRow(dbCtx, "SELECT id FROM categories ORDER BY name LIMIT 1").Scan(&categoryID)
		}

		createdAt := time.Now()
		if !item.PublishedAt.IsZero() {
			createdAt = item.PublishedAt
		}

		isViral := isViralContent(item.Title, item.Description, item.SourceURL, item.ContentType)
		var newContentID string
		err := conn.QueryRow(dbCtx, `
			INSERT INTO content (
				author_user_id, district_id, category_id, title, description,
				content_type, source_type, source_url, status, moderation_status,
				verification_status, is_viral, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, 'USER', $7, 'PENDING', 'UNMODERATED', 'SOURCE_IDENTIFIED', $8, $9, $9)
			RETURNING id
		`, authorID, districtID, categoryID, item.Title, item.Description, item.ContentType, item.SourceURL, isViral, createdAt).Scan(&newContentID)

		if err != nil {
			slog.Error("Failed to insert scraped item into PostgreSQL", slog.String("err", err.Error()), slog.String("title", item.Title))
			continue
		}

		item.ID = newContentID
		result.StagedCount++

		// If video detected, stage child row in video_links
		if item.ContentType == "VIDEO_LINK" && item.VideoID != "" {
			thumb := ""
			if len(item.ImageURLs) > 0 {
				thumb = item.ImageURLs[0]
			}
			_, _ = conn.Exec(dbCtx, `
				INSERT INTO video_links (
					content_id, platform, external_video_id, canonical_url,
					thumbnail_url, link_status, last_checked_at, created_at
				)
				VALUES ($1, $2, $3, $4, $5, 'ALIVE', NOW(), NOW())
				ON CONFLICT (external_video_id) DO NOTHING
			`, newContentID, item.VideoType, item.VideoID, item.VideoURL, thumb)
		}

		// Persist images into stories and photos tables so thumbnail queries succeed
		if len(item.ImageURLs) > 0 {
			_, _ = conn.Exec(dbCtx, `
				INSERT INTO stories (content_id, headline, body, single_photo_url)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (content_id) DO NOTHING
			`, newContentID, item.Title, item.Description, item.ImageURLs[0])

			_, _ = conn.Exec(dbCtx, `
				INSERT INTO photos (content_id, photo_count, photo_urls)
				VALUES ($1, $2, $3)
				ON CONFLICT (content_id) DO NOTHING
			`, newContentID, len(item.ImageURLs), item.ImageURLs)
		}

		// If this is an Event, stage child row in events table
		if strings.Contains(strings.ToLower(item.Category), "event") || item.ContentType == "EVENT" {
			_, _ = conn.Exec(dbCtx, `
				INSERT INTO events (content_id, event_name, event_date, organizer_name)
				VALUES ($1, $2, NOW() + INTERVAL '1 day', $3)
				ON CONFLICT (content_id) DO NOTHING
			`, newContentID, item.Title, item.District+" Civic & Event Desk")
		}
	}

	// Enforce configured retention policy immediately after staging so stale items are purged
	var hoursStr string
	_ = conn.QueryRow(dbCtx, "SELECT value FROM system_settings WHERE key = 'content_retention_hours'").Scan(&hoursStr)
	hours, _ := strconv.Atoi(hoursStr)
	if hours <= 0 {
		hours = 24
	}
	_, _ = conn.Exec(dbCtx, `
		DELETE FROM content
		WHERE created_at < NOW() - make_interval(hours => $1)
		  AND source_type != 'MANUAL_ENTRY'
		  AND status = 'PUBLISHED'
	`, hours)

	// Automatically run deduplication cleanup after staging to retain only the latest version
	_, _, _ = DeduplicateContentLocked(dbCtx, conn)

	return result, nil
}

// ParsePublishedTime extracts and parses publication date from strings, datelines, or timestamps
func ParsePublishedTime(raw string) (time.Time, bool) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, false
	}
	loc, _ := time.LoadLocation("Asia/Kolkata")
	if loc == nil {
		loc = time.FixedZone("IST", 5*3600+1800)
	}

	trimmed := strings.TrimSpace(raw)
	directLayouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}
	for _, l := range directLayouts {
		if t, err := time.Parse(l, trimmed); err == nil {
			return t, true
		}
	}

	dateRegexes := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:Published|Updated)\s*[-:]?\s*([A-Za-z]+\s+\d{1,2},\s+\d{4}\s+\d{1,2}:\d{2}\s*(?:am|pm)\s*(?:IST)?)`),
		regexp.MustCompile(`(?i)([A-Za-z]+\s+\d{1,2},\s+\d{4}\s+\d{1,2}:\d{2}\s*(?:am|pm)\s*(?:IST)?)`),
	}

	for _, re := range dateRegexes {
		if m := re.FindStringSubmatch(raw); len(m) > 1 {
			dtStr := regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(m[1]), " ")
			upper := strings.ReplaceAll(strings.ReplaceAll(dtStr, "am", "AM"), "pm", "PM")
			patterns := []string{
				"January 02, 2006 03:04 PM IST",
				"January 02, 2006 3:04 PM IST",
				"January 2, 2006 03:04 PM IST",
				"January 2, 2006 3:04 PM IST",
				"January 02, 2006 03:04 PM",
				"January 02, 2006 3:04 PM",
				"January 2, 2006 03:04 PM",
				"January 2, 2006 3:04 PM",
				"Jan 02, 2006 03:04 PM IST",
				"Jan 02, 2006 3:04 PM IST",
				"Jan 2, 2006 03:04 PM IST",
				"Jan 2, 2006 3:04 PM IST",
			}
			for _, p := range patterns {
				if t, err := time.ParseInLocation(p, upper, loc); err == nil {
					return t, true
				}
			}
		}
	}

	return time.Time{}, false
}

var brandSuffixes = []string{
	"- bbc news தமிழ்", "- bbc news tamil", "- bbc தமிழ்", "- பிபிசி தமிழ்", "| bbc news தமிழ்",
	"- தினமலர்", "| தினமலர்", "- dinamalar", "| dinamalar",
	"- தினத்தந்தி", "| தினத்தந்தி", "- daily thanthi", "| daily thanthi",
	"- தினகரன்", "| தினகரன்", "- dinakaran", "| dinakaran",
	"- தினமணி", "| தினமணி", "- dinamani", "| dinamani",
	"- மாலை மலர்", "| மாலை மலர்", "- maalai malar", "| maalai malar",
	"- the hindu", "| the hindu", "- இந்து தமிழ் திசை", "| இந்து தமிழ் திசை",
	"- news18", "| news18", "- polimer", "- thanthi tv", "- puthiyathalaimurai",
	"- oneindia tamil", "| oneindia tamil", "- vikatan", "| vikatan",
	"- nakkheeran", "| nakkheeran", "- news7 tamil", "| news7 tamil",
}

var puncRegex = regexp.MustCompile(`[\.,!?:;'"“”‘’()\[\]/|\-–—\s]+`)

// CanonicalizeURL strips tracking queries, fragments, trailing slashes, and host prefixes
func CanonicalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Host = strings.ToLower(u.Host)
	u.Host = strings.TrimPrefix(u.Host, "www.")
	u.Host = strings.TrimPrefix(u.Host, "m.")
	u.Path = strings.TrimRight(u.Path, "/")
	u.Fragment = ""

	q := u.Query()
	for k := range q {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "utm_") || strings.HasPrefix(lk, "at_") ||
			lk == "ref" || lk == "ref_src" || lk == "ref_url" || lk == "fbclid" ||
			lk == "gclid" || lk == "ocid" || lk == "pfrom" || lk == "amp" ||
			lk == "_ga" || lk == "_gl" {
			q.Del(k)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// NormalizeTitle strips HTML entities, brand suffixes, punctuation, and excess spaces
func NormalizeTitle(t string) string {
	s := strings.ToLower(CleanHTML(t))
	for _, b := range brandSuffixes {
		s = strings.TrimSuffix(s, b)
	}
	s = puncRegex.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// DeduplicateContent clusters duplicate posts across the system, retains ONLY the latest version,
// merges essential status/photos/body into the retained version, and purges older duplicates.
func DeduplicateContent(ctx context.Context, conn *pgxpool.Pool) (prunedCount int64, clusterCount int64, err error) {
	if conn == nil {
		return 0, 0, nil
	}
	DBMu.Lock()
	defer DBMu.Unlock()
	return DeduplicateContentLocked(ctx, conn)
}

// DeduplicateContentLocked runs deduplication while DBMu is already held by caller
func DeduplicateContentLocked(ctx context.Context, conn *pgxpool.Pool) (prunedCount int64, clusterCount int64, err error) {
	if conn == nil {
		return 0, 0, nil
	}

	dbCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	type postRow struct {
		id        string
		title     string
		normTitle string
		sourceURL string
		canonURL  string
		videoID   string
		status    string
		modStatus string
		isViral   bool
		photoURL  string
		body      string
		desc      string
		createdAt time.Time
		updatedAt time.Time
	}

	rows, err := conn.Query(dbCtx, `
		SELECT c.id, c.title, COALESCE(c.source_url, ''),
		       COALESCE(vl.external_video_id, ''),
		       c.status, c.moderation_status, COALESCE(c.is_viral, false),
		       COALESCE(s.single_photo_url, ''),
		       COALESCE(s.body, ''),
		       COALESCE(c.description, ''),
		       c.created_at, c.updated_at
		FROM content c
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN stories s ON c.id = s.content_id
		ORDER BY COALESCE(c.created_at, c.updated_at) DESC, c.updated_at DESC
	`)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query content for deduplication: %w", err)
	}
	defer rows.Close()

	var allPosts []*postRow
	for rows.Next() {
		var p postRow
		if err := rows.Scan(&p.id, &p.title, &p.sourceURL, &p.videoID,
			&p.status, &p.modStatus, &p.isViral,
			&p.photoURL, &p.body, &p.desc,
			&p.createdAt, &p.updatedAt); err != nil {
			continue
		}
		p.normTitle = NormalizeTitle(p.title)
		p.canonURL = CanonicalizeURL(p.sourceURL)
		allPosts = append(allPosts, &p)
	}
	rows.Close()

	parent := make(map[int]int)
	find := func(i int) int {
		root := i
		for parent[root] != root {
			root = parent[root]
		}
		curr := i
		for curr != root {
			nxt := parent[curr]
			parent[curr] = root
			curr = nxt
		}
		return root
	}
	union := func(i, j int) {
		rootI := find(i)
		rootJ := find(j)
		if rootI != rootJ {
			parent[rootJ] = rootI
		}
	}

	for i := range allPosts {
		parent[i] = i
	}

	urlMap := make(map[string]int)
	videoMap := make(map[string]int)
	titleMap := make(map[string]int)

	for i, p := range allPosts {
		if p.canonURL != "" && p.canonURL != "http://" && p.canonURL != "https://" {
			if existingIdx, exists := urlMap[p.canonURL]; exists {
				union(existingIdx, i)
			} else {
				urlMap[p.canonURL] = i
			}
		}

		if p.videoID != "" {
			if existingIdx, exists := videoMap[p.videoID]; exists {
				union(existingIdx, i)
			} else {
				videoMap[p.videoID] = i
			}
		}

		if len([]rune(p.normTitle)) >= 8 {
			if existingIdx, exists := titleMap[p.normTitle]; exists {
				union(existingIdx, i)
			} else {
				titleMap[p.normTitle] = i
			}
		}
	}

	clusters := make(map[int][]*postRow)
	for i, p := range allPosts {
		r := find(i)
		clusters[r] = append(clusters[r], p)
	}

	var duplicateIDs []string

	for _, group := range clusters {
		if len(group) <= 1 {
			continue
		}
		clusterCount++

		sort.Slice(group, func(i, j int) bool {
			if !group[i].createdAt.Equal(group[j].createdAt) {
				return group[i].createdAt.After(group[j].createdAt)
			}
			return group[i].updatedAt.After(group[j].updatedAt)
		})

		keeper := group[0]
		duplicates := group[1:]

		keeperUpdated := false
		keeperStatus := keeper.status
		keeperModStatus := keeper.modStatus
		keeperIsViral := keeper.isViral
		keeperPhoto := keeper.photoURL
		keeperBody := keeper.body
		keeperDesc := keeper.desc

		for _, d := range duplicates {
			duplicateIDs = append(duplicateIDs, d.id)

			if d.status == "PUBLISHED" && keeperStatus != "PUBLISHED" {
				keeperStatus = "PUBLISHED"
				keeperModStatus = "APPROVED"
				keeperUpdated = true
			}
			if d.isViral && !keeperIsViral {
				keeperIsViral = true
				keeperUpdated = true
			}
			if (keeperPhoto == "" || strings.Contains(keeperPhoto, "/maps/svg")) && d.photoURL != "" && !strings.Contains(d.photoURL, "/maps/svg") {
				keeperPhoto = d.photoURL
				keeperUpdated = true
			}
			if len([]rune(keeperBody)) < 300 && len([]rune(d.body)) > len([]rune(keeperBody)) {
				keeperBody = d.body
				if len([]rune(keeperDesc)) < 300 {
					keeperDesc = d.desc
				}
				keeperUpdated = true
			}
		}

		if keeperUpdated {
			_, _ = conn.Exec(dbCtx, `
				UPDATE content
				SET status = $1, moderation_status = $2, is_viral = $3, description = $4, updated_at = NOW()
				WHERE id = $5
			`, keeperStatus, keeperModStatus, keeperIsViral, keeperDesc, keeper.id)

			if keeperPhoto != "" {
				_, _ = conn.Exec(dbCtx, `
					UPDATE stories SET single_photo_url = $1, body = $2 WHERE content_id = $3
				`, keeperPhoto, keeperBody, keeper.id)
				_, _ = conn.Exec(dbCtx, `
					UPDATE photos SET photo_urls = ARRAY[$1] WHERE content_id = $2
				`, keeperPhoto, keeper.id)
			} else if keeperBody != "" {
				_, _ = conn.Exec(dbCtx, `
					UPDATE stories SET body = $1 WHERE content_id = $2
				`, keeperBody, keeper.id)
			}
		}
	}

	if len(duplicateIDs) > 0 {
		tag, delErr := conn.Exec(dbCtx, `DELETE FROM content WHERE id = ANY($1)`, duplicateIDs)
		if delErr != nil {
			return 0, clusterCount, fmt.Errorf("failed to delete duplicate posts: %w", delErr)
		}
		prunedCount = tag.RowsAffected()
		slog.Info("Deduplication complete: retained latest post in each cluster", slog.Int64("pruned", prunedCount), slog.Int64("clusters", clusterCount))
	}

	return prunedCount, clusterCount, nil
}

// MigrateDatesAndEnforceRetentionPolicy migrates real publication dates from text and purges stale items according to config
func MigrateDatesAndEnforceRetentionPolicy(ctx context.Context, conn *pgxpool.Pool) {
	if conn == nil {
		return
	}
	DBMu.Lock()
	defer DBMu.Unlock()

	rows, err := conn.Query(ctx, `
		SELECT c.id, c.description, s.body, c.created_at
		FROM content c
		LEFT JOIN stories s ON s.content_id = c.id
		WHERE c.source_type != 'MANUAL_ENTRY'
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	type itemFix struct {
		id        string
		cleanDesc string
		cleanBody string
		realDate  time.Time
	}
	var fixes []itemFix

	cleanJunk := func(txt string) string {
		lines := strings.Split(txt, "\n")
		var filtered []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			lower := strings.ToLower(trimmed)
			if strings.Contains(lower, "subscribed with another email") ||
				strings.Contains(lower, "logout and login") ||
				strings.HasPrefix(lower, "published -") ||
				strings.HasPrefix(lower, "updated -") ||
				strings.HasPrefix(lower, "published:") ||
				strings.HasPrefix(lower, "updated:") {
				continue
			}
			filtered = append(filtered, line)
		}
		return strings.TrimSpace(strings.Join(filtered, "\n\n"))
	}

	for rows.Next() {
		var id string
		var desc string
		var body *string
		var curCreated time.Time
		if err := rows.Scan(&id, &desc, &body, &curCreated); err != nil {
			continue
		}
		bodyStr := ""
		if body != nil {
			bodyStr = *body
		}
		combined := desc + " " + bodyStr
		realDate, hasDate := ParsePublishedTime(combined)

		cleanDesc := cleanJunk(desc)
		cleanBody := cleanJunk(bodyStr)

		if hasDate || cleanDesc != desc || cleanBody != bodyStr {
			fixes = append(fixes, itemFix{
				id:        id,
				cleanDesc: cleanDesc,
				cleanBody: cleanBody,
				realDate:  realDate,
			})
		}
	}
	rows.Close()

	for _, f := range fixes {
		if !f.realDate.IsZero() {
			_, _ = conn.Exec(ctx, "UPDATE content SET created_at = $1, description = $2 WHERE id = $3", f.realDate, f.cleanDesc, f.id)
			_, _ = conn.Exec(ctx, "UPDATE stories SET body = $1 WHERE content_id = $2", f.cleanBody, f.id)
		} else {
			if f.cleanDesc != "" {
				_, _ = conn.Exec(ctx, "UPDATE content SET description = $1 WHERE id = $2", f.cleanDesc, f.id)
			}
			if f.cleanBody != "" {
				_, _ = conn.Exec(ctx, "UPDATE stories SET body = $1 WHERE content_id = $2", f.cleanBody, f.id)
			}
		}
	}

	var hoursStr string
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'content_retention_hours'").Scan(&hoursStr)
	hours, _ := strconv.Atoi(hoursStr)
	if hours <= 0 {
		hours = 24
	}

	tag, err := conn.Exec(ctx, `
		DELETE FROM content
		WHERE created_at < NOW() - make_interval(hours => $1)
		  AND source_type != 'MANUAL_ENTRY'
		  AND status = 'PUBLISHED'
	`, hours)
	if err == nil {
		purged := tag.RowsAffected()
		slog.Info("Retention policy enforced: pruned stale posts older than config", slog.Int("hours", hours), slog.Int64("purged", purged))
		_, _ = conn.Exec(ctx, `
			INSERT INTO system_settings (key, value, updated_at)
			VALUES ('last_cleanup_at', $1, NOW())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
		`, time.Now().Format(time.RFC3339))
		_, _ = conn.Exec(ctx, `
			INSERT INTO system_settings (key, value, updated_at)
			VALUES ('last_cleanup_count', $1, NOW())
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
		`, fmt.Sprintf("%d", purged))
	}

	// Also purge duplicates automatically during retention cycle
	_, _, _ = DeduplicateContentLocked(ctx, conn)
}

func parseRSSItem(it rssItem) *ScrapedItem {
	title := cleanHTML(it.Title)
	if title == "" {
		return nil
	}

	desc := cleanHTML(it.Description)
	district := detectDistrict(title + " " + desc + " " + it.Link)
	category := detectCategory(title+" "+desc, it.Link)

	publishedAt, _ := ParsePublishedTime(it.PubDate)

	item := &ScrapedItem{
		Title:       title,
		Description: desc,
		SourceURL:   it.Link,
		District:    district,
		Category:    category,
		ContentType: "TEXT_STORY",
		Status:      "PENDING",
		PublishedAt: publishedAt,
	}

	// Extract images
	if it.MediaThumbnail.URL != "" {
		item.ImageURLs = append(item.ImageURLs, it.MediaThumbnail.URL)
	}
	if it.MediaContent.URL != "" && (it.MediaContent.Medium == "image" || strings.Contains(it.MediaContent.URL, ".jpg") || strings.Contains(it.MediaContent.URL, ".png")) {
		item.ImageURLs = append(item.ImageURLs, it.MediaContent.URL)
	}

	// Check for videos
	if strings.Contains(it.Link, "youtube.com") || strings.Contains(it.Link, "youtu.be") {
		item.ContentType = "VIDEO_LINK"
		item.VideoType = "youtube"
		item.VideoURL = it.Link
		item.VideoID = extractYouTubeID(it.Link)
		item.Category = "Viral Videos"
		if len(item.ImageURLs) == 0 && item.VideoID != "" {
			item.ImageURLs = append(item.ImageURLs, fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", item.VideoID))
		}
	} else if strings.Contains(it.Enclosure.Type, "video") && it.Enclosure.URL != "" {
		item.ContentType = "VIDEO_LINK"
		item.VideoType = "direct"
		item.VideoURL = it.Enclosure.URL
		item.Category = "Viral Videos"
	}

	// Never leave No Media!
	if len(item.ImageURLs) == 0 {
		item.ImageURLs = append(item.ImageURLs, getFallbackImage(district, category))
	}

	return item
}

func parseAtomEntry(e atomEntry) *ScrapedItem {
	title := cleanHTML(e.Title)
	if title == "" {
		return nil
	}

	link := e.Link.Href
	desc := cleanHTML(e.MediaGroup.Description)
	district := detectDistrict(title + " " + desc)

	videoID := e.VideoID
	if videoID == "" && strings.Contains(e.ID, "yt:video:") {
		videoID = strings.TrimPrefix(e.ID, "yt:video:")
	}

	publishedAt, _ := ParsePublishedTime(e.Published)

	item := &ScrapedItem{
		Title:       title,
		Description: desc,
		SourceURL:   link,
		District:    district,
		Category:    "Viral Videos",
		ContentType: "VIDEO_LINK",
		VideoURL:    link,
		VideoID:     videoID,
		VideoType:   "youtube",
		Status:      "PENDING",
		PublishedAt: publishedAt,
	}

	if e.MediaGroup.Thumbnail.URL != "" {
		item.ImageURLs = append(item.ImageURLs, e.MediaGroup.Thumbnail.URL)
	} else if videoID != "" {
		item.ImageURLs = append(item.ImageURLs, fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID))
	}

	if len(item.ImageURLs) == 0 {
		item.ImageURLs = append(item.ImageURLs, getFallbackImage(district, "Viral Videos"))
	}

	return item
}

// extractMetaContent searches for meta tag content regardless of attribute order or extra attributes (e.g. data-react-helmet)
func extractMetaContent(htmlContent string, namesOrProperties ...string) string {
	for _, prop := range namesOrProperties {
		// Pattern 1: property/name="..." ... content="..."
		re1 := regexp.MustCompile(`(?i)<meta\s+[^>]*?(?:property|name)=["']` + regexp.QuoteMeta(prop) + `["'][^>]*?content=["']([^"']*)["']`)
		if m := re1.FindStringSubmatch(htmlContent); len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			return strings.TrimSpace(m[1])
		}
		// Pattern 2: content="..." ... property/name="..."
		re2 := regexp.MustCompile(`(?i)<meta\s+[^>]*?content=["']([^"']*)["'][^>]*?(?:property|name)=["']` + regexp.QuoteMeta(prop) + `["']`)
		if m := re2.FindStringSubmatch(htmlContent); len(m) > 1 && strings.TrimSpace(m[1]) != "" {
			return strings.TrimSpace(m[1])
		}
	}
	return ""
}

// extractBBCArticleID extracts the article/clip identifier from a BBC URL
func extractBBCArticleID(u string) string {
	parts := strings.Split(strings.TrimRight(u, "/"), "/")
	if len(parts) > 0 {
		clean := strings.Split(parts[len(parts)-1], "#")[0]
		return strings.Split(clean, "?")[0]
	}
	return ""
}

// extractBBCVideoEmbed detects BBC embedded media player URLs or aresMedia clip IDs
func extractBBCVideoEmbed(targetURL, htmlContent string) (string, string, string) {
	// 1. Direct BBC av-embeds URL in the HTML page
	reEmbed := regexp.MustCompile(`(?i)https?://(?:www\.)?bbc\.com/ws/av-embeds/[^\s"'<>]+`)
	if m := reEmbed.FindString(htmlContent); m != "" {
		clipID := ""
		clipRe := regexp.MustCompile(`/(p[a-z0-9]{7,8})/`)
		if cm := clipRe.FindStringSubmatch(m); len(cm) > 1 {
			clipID = cm[1]
		} else {
			clipID = "bbc_" + extractBBCArticleID(targetURL)
		}
		return clipID, m, "bbc"
	}

	// 2. BBC aresMedia clip ID embedded in Next.js state or JSON
	clipRe := regexp.MustCompile(`"format"\s*:\s*"video"[^}]*?"id"\s*:\s*"([a-zA-Z0-9]+)"`)
	if cm := clipRe.FindStringSubmatch(htmlContent); len(cm) > 1 {
		clipID := cm[1]
		artID := extractBBCArticleID(targetURL)
		embedURL := fmt.Sprintf("https://www.bbc.com/ws/av-embeds/articles/%s/%s/ta", artID, clipID)
		return clipID, embedURL, "bbc"
	}

	// 3. Fallback: if this is a BBC watch page
	if strings.Contains(targetURL, "bbc.com") && strings.Contains(targetURL, "/watch/") {
		artID := extractBBCArticleID(targetURL)
		if artID != "" {
			return "bbc_" + artID, targetURL, "bbc"
		}
	}

	return "", "", ""
}

// cleanSourceTitle strips trailing source branding such as " - BBC News தமிழ்" or " - The Hindu"
func cleanSourceTitle(title string) string {
	title = cleanHTML(title)
	re := regexp.MustCompile(`(?i)\s*[-|–—]\s*(?:BBC News தமிழ்|BBC News|BBC|The Hindu|Dinamalar|தினமலர்|Puthiyathalaimurai|புதிய தலைமுறை|Maalai Malar|மாலை மலர்|மாவை மலர்|News18|Oneindia).*$`)
	title = re.ReplaceAllString(title, "")
	return strings.TrimSpace(title)
}

// isSingleArticle determines if a given URL or page content is a dedicated article/video rather than a portal index
func isSingleArticle(targetURL, htmlContent string) bool {
	uLower := strings.ToLower(targetURL)
	if strings.Contains(uLower, "/articles/") || strings.Contains(uLower, "/watch/") || strings.HasSuffix(uLower, ".html") {
		return true
	}
	ogType := strings.ToLower(extractMetaContent(htmlContent, "og:type"))
	if ogType == "article" || strings.Contains(ogType, "video") {
		return true
	}
	return false
}

// isBoilerplateAnchor filters out accessibility "skip to content" links, navigation fragments, and UI buttons
func isBoilerplateAnchor(text, href string) bool {
	lowerHref := strings.ToLower(strings.TrimSpace(href))
	if strings.HasPrefix(lowerHref, "#") || strings.Contains(lowerHref, "#content") || strings.Contains(lowerHref, "#main") || strings.HasPrefix(lowerHref, "javascript:") {
		return true
	}
	lowerText := strings.ToLower(strings.TrimSpace(text))
	if strings.Contains(lowerText, "உள்ளடக்கத்துக்குத் தாண்டிச் செல்க") ||
		strings.Contains(lowerText, "உள்ளடக்கத்திற்கு செல்ல") ||
		strings.Contains(lowerText, "தாண்டிச் செல்க") ||
		strings.Contains(lowerText, "செல்லவும்") ||
		strings.Contains(lowerText, "skip to content") ||
		strings.Contains(lowerText, "skip to main") ||
		strings.Contains(lowerText, "முகப்புபார்க்ககேட்கபிரிவுகள்") ||
		(strings.Contains(lowerText, "முகப்பு") && len([]rune(text)) < 8) ||
		(strings.Contains(lowerText, "பிரிவுகள்") && len([]rune(text)) < 10) ||
		strings.Contains(lowerText, "cookie") ||
		strings.Contains(lowerText, "privacy") ||
		strings.Contains(lowerText, "terms of use") ||
		strings.Contains(lowerText, "தனியுரிமை") ||
		strings.Contains(lowerText, "விதிமுறைகள்") ||
		strings.Contains(lowerText, "about the bbc") {
		return true
	}
	return false
}

func extractFullTextAndMedia(targetURL, htmlContent string) (string, []string, string, string, string, time.Time) {
	baseParsed, _ := url.Parse(targetURL)

	videoID := ""
	videoURL := ""
	videoType := "none"

	// 1. Detect BBC Video Embed first if available
	if bbcClipID, bbcEmbedURL, bbcType := extractBBCVideoEmbed(targetURL, htmlContent); bbcEmbedURL != "" {
		videoID = bbcClipID
		videoURL = bbcEmbedURL
		videoType = bbcType
	}

	// 2. Detect YouTube video if present
	if videoURL == "" {
		ytRegex := regexp.MustCompile(`(?i)(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:watch\?v=|embed\/)|youtu\.be\/)([a-zA-Z0-9_-]{11})`)
		if m := ytRegex.FindStringSubmatch(htmlContent); len(m) > 1 {
			videoID = m[1]
			videoURL = fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
			videoType = "youtube"
		}
	}

	// 3. Detect generic HTML5 / OpenGraph video embed
	if videoURL == "" {
		if ogVid := extractMetaContent(htmlContent, "og:video:secure_url", "og:video:url", "og:video", "twitter:player"); ogVid != "" {
			videoURL = ogVid
			videoID = "vid_" + extractBBCArticleID(ogVid)
			videoType = "other"
		}
	}
	if videoURL == "" {
		videoTagRegex := regexp.MustCompile(`(?i)<(?:video|source)\s+[^>]*src=["']([^"']+\.(?:mp4|webm|m3u8)[^"']*)["']`)
		if m := videoTagRegex.FindStringSubmatch(htmlContent); len(m) > 1 {
			rawVid := m[1]
			if baseParsed != nil {
				if ref, err := url.Parse(rawVid); err == nil {
					rawVid = baseParsed.ResolveReference(ref).String()
				}
			}
			videoURL = rawVid
			videoID = "vid_" + extractBBCArticleID(videoURL)
			videoType = "other"
		}
	}

	// 4. Extract all article images
	var images []string
	seenImg := make(map[string]bool)

	// Primary high-res article image from OpenGraph / Twitter meta
	primaryImg := extractMetaContent(htmlContent, "og:image", "og:image:url", "og:image:secure_url", "twitter:image", "twitter:image:src")
	if primaryImg != "" && !seenImg[primaryImg] {
		seenImg[primaryImg] = true
		images = append(images, primaryImg)
	}

	// Secondary article img tags
	imgRegex := regexp.MustCompile(`(?i)<img\s+[^>]*src=["']([^"']+)["']`)
	for _, m := range imgRegex.FindAllStringSubmatch(htmlContent, -1) {
		if len(m) < 2 {
			continue
		}
		rawSrc := m[1]
		lower := strings.ToLower(rawSrc)
		if strings.Contains(lower, "icon") || strings.Contains(lower, "logo") || strings.Contains(lower, "pixel") ||
			strings.Contains(lower, "ad.") || strings.Contains(lower, "1x1") || strings.Contains(lower, "avatar") {
			continue
		}
		if !(strings.Contains(lower, ".jpg") || strings.Contains(lower, ".jpeg") || strings.Contains(lower, ".png") || strings.Contains(lower, ".webp")) {
			continue
		}
		fullImgURL := rawSrc
		if baseParsed != nil {
			ref, err := url.Parse(rawSrc)
			if err == nil {
				fullImgURL = baseParsed.ResolveReference(ref).String()
			}
		}
		if !seenImg[fullImgURL] {
			seenImg[fullImgURL] = true
			images = append(images, fullImgURL)
		}
		if len(images) >= 8 {
			break
		}
	}

	// 5. Extract Published Date from HTML meta or time tags
	publishedAt := time.Time{}
	if metaDate := extractMetaContent(htmlContent, "article:published_time", "og:published_time", "publish-date", "pubdate", "article:modified_time"); metaDate != "" {
		if t, ok := ParsePublishedTime(metaDate); ok {
			publishedAt = t
		}
	}
	if publishedAt.IsZero() {
		timeTagRegex := regexp.MustCompile(`(?i)<time\s+[^>]*datetime=["']([^"']+)["']`)
		if m := timeTagRegex.FindStringSubmatch(htmlContent); len(m) > 1 {
			if t, ok := ParsePublishedTime(m[1]); ok {
				publishedAt = t
			}
		}
	}

	// 6. Extract Full Article Body Text
	cleanedHTML := htmlContent
	for _, tag := range []string{"script", "style", "nav", "header", "footer", "aside", "noscript", "svg", "form", "button", "iframe"} {
		re := regexp.MustCompile(`(?is)<` + tag + `[^>]*>.*?</` + tag + `>`)
		cleanedHTML = re.ReplaceAllString(cleanedHTML, "")
	}

	var paragraphs []string
	pRegex := regexp.MustCompile(`(?is)<p[^>]*>(.*?)</p>`)
	for _, m := range pRegex.FindAllStringSubmatch(cleanedHTML, -1) {
		if len(m) < 2 {
			continue
		}
		txt := strings.TrimSpace(cleanHTML(m[1]))
		lower := strings.ToLower(txt)

		// Filter out navigation boilerplate, skip links, subscription spam and datelines
		if strings.Contains(lower, "உள்ளடக்கத்துக்குத் தாண்டிச் செல்க") ||
			strings.Contains(lower, "உள்ளடக்கத்திற்கு செல்ல") ||
			strings.Contains(lower, "skip to content") ||
			strings.Contains(lower, "முகப்புபார்க்ககேட்கபிரிவுகள்") ||
			strings.Contains(lower, "subscribed with another email") ||
			strings.Contains(lower, "logout and login") ||
			strings.HasPrefix(lower, "published -") ||
			strings.HasPrefix(lower, "updated -") ||
			strings.HasPrefix(lower, "published:") ||
			strings.HasPrefix(lower, "updated:") {
			if publishedAt.IsZero() {
				if t, ok := ParsePublishedTime(txt); ok {
					publishedAt = t
				}
			}
			continue
		}

		if len([]rune(txt)) > 15 &&
			!strings.Contains(lower, "all rights reserved") &&
			!strings.Contains(lower, "subscribe to") &&
			!strings.Contains(lower, "subscription") &&
			!strings.Contains(lower, "advertisement") &&
			!strings.Contains(lower, "whatsapp") &&
			!strings.Contains(lower, "follow us on") &&
			!strings.Contains(lower, "copyright") &&
			!strings.Contains(lower, "terms of service") &&
			!strings.Contains(lower, "privacy policy") &&
			!strings.Contains(lower, "choose language") &&
			!strings.Contains(lower, "newsletter") &&
			!strings.Contains(lower, "sign in") &&
			!strings.Contains(lower, "photo credit:") &&
			!strings.Contains(lower, "news and reviews from the world") &&
			!strings.Contains(lower, "looking at world affairs") &&
			!strings.Contains(lower, "download of the top") &&
			!strings.Contains(lower, "books of the week") &&
			!strings.Contains(lower, "stories from beyond the binary") {
			paragraphs = append(paragraphs, txt)
		}
		if len(paragraphs) >= 100 {
			break
		}
	}

	fullText := ""
	if len(paragraphs) > 0 {
		fullText = strings.Join(paragraphs, "\n\n")
	} else {
		// Fallback to og:description
		if ogDesc := extractMetaContent(htmlContent, "og:description", "description", "twitter:description"); ogDesc != "" {
			fullText = cleanHTML(ogDesc)
		}
	}

	if publishedAt.IsZero() {
		if t, ok := ParsePublishedTime(htmlContent); ok {
			publishedAt = t
		}
	}

	return fullText, images, videoID, videoURL, videoType, publishedAt
}

func FetchFullTextAndMediaExported(sourceURL string) (string, []string, string, string, string, time.Time) {
	if !strings.HasPrefix(sourceURL, "http") || strings.Contains(sourceURL, "youtube.com") || strings.Contains(sourceURL, "youtu.be") {
		return "", nil, "", "", "", time.Time{}
	}
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest(http.MethodGet, sourceURL, nil)
	if err != nil {
		return "", nil, "", "", "", time.Time{}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, "", "", "", time.Time{}
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil || len(bytes) == 0 {
		return "", nil, "", "", "", time.Time{}
	}
	fullText, images, videoID, videoURL, videoType, publishedAt := extractFullTextAndMedia(sourceURL, string(bytes))
	return fullText, images, videoID, videoURL, videoType, publishedAt
}

func parseHTMLPage(targetURL, htmlContent string) *ScrapedItem {
	// Extract clean title from og:title, twitter:title, or title tag
	title := extractMetaContent(htmlContent, "og:title", "twitter:title")
	if title == "" {
		titleTagRegex := regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
		if m := titleTagRegex.FindStringSubmatch(htmlContent); len(m) > 1 {
			title = m[1]
		}
	}
	title = cleanSourceTitle(title)
	if title == "" {
		title = "Tamil Nadu News Report: " + targetURL
	}

	// Extract full text content, all images, and any video
	fullText, images, videoID, videoURL, videoType, publishedAt := extractFullTextAndMedia(targetURL, htmlContent)
	if fullText == "" {
		fullText = "Civic update scraped from " + targetURL
	}

	district := detectDistrict(title + " " + fullText + " " + targetURL)
	category := detectCategory(title+" "+fullText, targetURL)

	if len(images) == 0 {
		if personImg := detectPersonImage(title, fullText); personImg != "" {
			images = append(images, personImg)
		} else if catImg := detectCategoryImage(category, title+" "+fullText); catImg != "" {
			images = append(images, catImg)
		} else {
			images = append(images, getFallbackImage(district, category))
		}
	}

	// Check if video link or text story
	contentType := "TEXT_STORY"
	if videoURL != "" || videoID != "" {
		contentType = "VIDEO_LINK"
		if videoType == "" || videoType == "none" {
			if strings.Contains(videoURL, "youtube.com") || strings.Contains(videoURL, "youtu.be") {
				videoType = "youtube"
			} else if strings.Contains(videoURL, "bbc.com") || strings.Contains(videoURL, "av-embeds") {
				videoType = "bbc"
			} else {
				videoType = "other"
			}
		}
	}

	return &ScrapedItem{
		Title:       title,
		Description: fullText,
		SourceURL:   targetURL,
		District:    district,
		Category:    category,
		ContentType: contentType,
		ImageURLs:   images,
		VideoURL:    videoURL,
		VideoID:     videoID,
		VideoType:   videoType,
		Status:      "PENDING",
		PublishedAt: publishedAt,
	}
}

func parseHTMLPageMultiple(targetURL, htmlContent string) []*ScrapedItem {
	var results []*ScrapedItem
	seen := make(map[string]bool)

	// Look for article links <a ... href="..." ...>headline text</a>
	aTagRegex := regexp.MustCompile(`(?i)<a\s+[^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := aTagRegex.FindAllStringSubmatch(htmlContent, -1)

	baseParsed, _ := url.Parse(targetURL)

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		rawHref := strings.TrimSpace(m[1])
		rawInner := cleanHTML(m[2])

		// Skip accessibility "skip to content" links, jump fragments, and UI boilerplate
		if isBoilerplateAnchor(rawInner, rawHref) {
			continue
		}

		anchorInner := cleanSourceTitle(rawInner)

		// Filter anchors to headline length
		if len([]rune(anchorInner)) < 15 || len([]rune(anchorInner)) > 250 {
			continue
		}

		// Resolve absolute URL
		fullURL := rawHref
		if baseParsed != nil {
			ref, err := url.Parse(rawHref)
			if err == nil {
				resolved := baseParsed.ResolveReference(ref)
				if resolved.Host == baseParsed.Host || strings.Contains(resolved.Host, "youtube.com") || strings.Contains(resolved.Host, "youtu.be") {
					fullURL = resolved.String()
				} else {
					continue
				}
			}
		}

		if fullURL == targetURL || seen[fullURL] || seen[anchorInner] {
			continue
		}
		seen[fullURL] = true
		seen[anchorInner] = true

		district := detectDistrict(anchorInner + " " + fullURL)
		category := detectCategory(anchorInner, fullURL)
		contentType := "TEXT_STORY"
		videoURL := ""
		videoID := ""
		videoType := "none"
		var images []string

		// Check if link is YouTube
		if strings.Contains(fullURL, "youtube.com") || strings.Contains(fullURL, "youtu.be") {
			contentType = "VIDEO_LINK"
			videoID = extractYouTubeID(fullURL)
			videoType = "youtube"
			videoURL = fullURL
			category = "Viral Videos"
			if videoID != "" {
				images = append(images, fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID))
			}
		} else if strings.Contains(fullURL, "bbc.com") && strings.Contains(fullURL, "/watch/") {
			// Check if link is BBC Watch page
			contentType = "VIDEO_LINK"
			clipID, embedURL, bType := extractBBCVideoEmbed(fullURL, "")
			videoID = clipID
			videoURL = embedURL
			videoType = bType
			category = "Viral Videos"
		}

		// Check if anchor markup contains an image
		if len(images) == 0 {
			imgRegex := regexp.MustCompile(`(?i)<img\s+[^>]*src=["']([^"']+)["']`)
			imgMatch := imgRegex.FindStringSubmatch(m[0])
			if len(imgMatch) > 1 && !strings.Contains(imgMatch[1], "icon") && !strings.Contains(imgMatch[1], "logo") {
				imgHref := imgMatch[1]
				if baseParsed != nil {
					ref, err := url.Parse(imgHref)
					if err == nil {
						imgHref = baseParsed.ResolveReference(ref).String()
					}
				}
				images = append(images, imgHref)
			}
		}

		// Guarantee NO item is left without media!
		if len(images) == 0 {
			images = append(images, getFallbackImage(district, category))
		}

		var desc string
		if DetectLanguage(anchorInner) == "ta" {
			desc = fmt.Sprintf("%s - %s மற்றும் தமிழ்நாடு வட்டார முக்கிய நிகழ்வுகள் குறித்த விரிவான கள நிலவரம் மற்றும் நேரடி செய்தி தொகுப்பு.", anchorInner, district)
		} else {
			desc = fmt.Sprintf("%s. Comprehensive on-ground coverage and latest news updates from %s.", anchorInner, district)
		}

		results = append(results, &ScrapedItem{
			Title:       anchorInner,
			Description: desc,
			SourceURL:   fullURL,
			District:    district,
			Category:    category,
			ContentType: contentType,
			ImageURLs:   images,
			VideoURL:    videoURL,
			VideoID:     videoID,
			VideoType:   videoType,
			Status:      "PENDING",
		})

		if len(results) >= 30 {
			break
		}
	}

	return results
}

var internationalKeywords = []string{
	"பாங்காக்", "தாய்லாந்து", "பட்டாயா", "அமெரிக்கா", "வாஷிங்டன்", "நியூயார்க்", "ரஷ்யா", "உக்ரைன்",
	"இஸ்ரேல்", "காசா", "பாலஸ்தீனம்", "சீனா", "இலங்கை", "கொழும்பு", "பாகிஸ்தான்", "வங்கதேசம்",
	"டாக்கா", "துபாய்", "யுஏஇ", "லண்டன்", "இங்கிலாந்து", "பிரிட்டன்", "ஜப்பான்", "டோக்கியோ",
	"மலேசியா", "சிங்கப்பூர்", "கனடா", "ஆஸ்திரேலியா", "ஜெர்மனி", "பிரான்ஸ்", "பாரிஸ்",
	"bangkok", "thailand", "pattaya", "usa", "united states", "washington", "new york", "russia", "moscow",
	"ukraine", "kyiv", "israel", "gaza", "palestine", "china", "beijing", "sri lanka", "colombo",
	"pakistan", "islamabad", "bangladesh", "dhaka", "dubai", "uae", "london", "uk", "britain",
	"japan", "tokyo", "malaysia", "singapore", "canada", "australia", "germany", "france", "paris",
	"international", "world news", "global",
}

var nationalKeywords = []string{
	"லக்னோ", "உத்தரப்பிரதேசம்", "உபி", "சஹரான்பூர்", "கொல்கத்தா", "மேற்கு வங்கம்", "மேற்குவங்க",
	"மணிப்பூர்", "இம்பால்", "டெல்லி", "புதுடெல்லி", "மும்பை", "மகாராஷ்டிரா", "பெங்களூரு", "கர்நாடகா",
	"ஹைதராபாத்", "ஆந்திரா", "தெலங்கானா", "கேரளா", "திருவனந்தபுரம்", "கொச்சி", "கோழிக்கோடு",
	"பிகார்", "பாட்னா", "குஜராத்", "அகமதாபாத்", "ராஜஸ்தான்", "ஜெய்ப்பூர்", "ஒடிசா", "புவனேஸ்வர்",
	"பஞ்சாப்", "சண்டிகர்", "ஹரியானா", "காஷ்மீர்", "ஸ்ரீநகர்", "அசாம்", "குவாஹாட்டி", "மத்தியப் பிரதேசம்",
	"போபால்", "ஜார்க்கண்ட்", "ராஞ்சி", "சத்தீஸ்கர்", "ராய்ப்பூர்", "கோவா",
	"lucknow", "uttar pradesh", "saharanpur", "kolkata", "west bengal", "manipur", "imphal",
	"delhi", "new delhi", "mumbai", "maharashtra", "bengaluru", "bangalore", "karnataka",
	"hyderabad", "andhra", "telangana", "kerala", "trivandrum", "kochi", "bihar", "patna",
	"gujarat", "ahmedabad", "rajasthan", "jaipur", "odisha", "punjab", "chandigarh",
	"haryana", "kashmir", "srinagar", "assam", "bhopal", "national news", "india news",
}

// Major Tamil Nadu town and taluk mappings
var tnTownDistrictMap = map[string]string{
	"சென்னை":          "Chennai",
	"தாம்பரம்":         "Chennai",
	"ஆவடி":            "Chennai",
	"கிண்டி":           "Chennai",
	"அடையாறு":         "Chennai",
	"மயிலாப்பூர்":       "Chennai",
	"வேளச்சேரி":        "Chennai",
	"வளசரவாக்கம்":      "Chennai",
	"madras":          "Chennai",
	"chennai":         "Chennai",
	
	"மதுரை":           "Madurai",
	"மேலூர்":           "Madurai",
	"உசிலம்பட்டி":      "Madurai",
	"திருமங்கலம்":      "Madurai",
	"madurai":         "Madurai",

	"கோவை":           "Coimbatore",
	"கோயம்புத்தூர்":     "Coimbatore",
	"பொள்ளாச்சி":       "Coimbatore",
	"மேட்டுப்பாளையம்":   "Coimbatore",
	"coimbatore":      "Coimbatore",

	"திருச்சி":          "Tiruchirappalli",
	"திருச்சிராப்பள்ளி": "Tiruchirappalli",
	"ஸ்ரீரங்கம்":        "Tiruchirappalli",
	"மணப்பாறை":        "Tiruchirappalli",
	"trichy":          "Tiruchirappalli",
	"tiruchirappalli": "Tiruchirappalli",

	"சேலம்":           "Salem",
	"ஆத்தூர்":          "Salem",
	"மேட்டூர்":          "Salem",
	"salem":           "Salem",

	"தஞ்சாவூர்":        "Thanjavur",
	"தஞ்சை":           "Thanjavur",
	"கும்பகோணம்":       "Thanjavur",
	"பட்டுக்கோட்டை":     "Thanjavur",
	"thanjavur":       "Thanjavur",
	"tanjore":         "Thanjavur",

	"நெல்லை":          "Tirunelveli",
	"திருநெல்வேலி":     "Tirunelveli",
	"பாளையங்கோட்டை":     "Tirunelveli",
	"tirunelveli":     "Tirunelveli",

	"கன்னியாகுமரி":     "Kanyakumari",
	"நாகர்கோவில்":      "Kanyakumari",
	"மார்த்தாண்டம்":     "Kanyakumari",
	"குமரி":            "Kanyakumari",
	"kanyakumari":     "Kanyakumari",
	"nagercoil":       "Kanyakumari",

	"திண்டுக்கல்":      "Dindigul",
	"பழனி":            "Dindigul",
	"கொடைக்கானல்":      "Dindigul",
	"dindigul":        "Dindigul",
	"palani":          "Dindigul",
	"kodaikanal":      "Dindigul",

	"வேலூர்":           "Vellore",
	"காட்பாடி":          "Vellore",
	"குடியாத்தம்":       "Vellore",
	"vellore":         "Vellore",

	"ஈரோடு":           "Erode",
	"கோபிசெட்டிபாளையம்": "Erode",
	"பவானி":           "Erode",
	"erode":           "Erode",

	"திருப்பூர்":         "Tiruppur",
	"அவிநாசி":          "Tiruppur",
	"தாராபுரம்":         "Tiruppur",
	"tiruppur":        "Tiruppur",

	"தூத்துக்குடி":      "Thoothukudi",
	"கோவில்பட்டி":      "Thoothukudi",
	"திருச்செந்தூர்":     "Thoothukudi",
	"thoothukudi":     "Thoothukudi",
	"tuticorin":       "Thoothukudi",

	"கடலூர்":          "Cuddalore",
	"சிதம்பரம்":        "Cuddalore",
	"விருத்தாசலம்":     "Cuddalore",
	"நெய்வேலி":         "Cuddalore",
	"cuddalore":       "Cuddalore",
	"neyveli":         "Cuddalore",

	"காஞ்சிபுரம்":      "Kanchipuram",
	"காஞ்சி":           "Kanchipuram",
	"ஸ்ரீபெரும்புதூர்":   "Kanchipuram",
	"kanchipuram":     "Kanchipuram",

	"செங்கல்பட்டு":      "Chengalpattu",
	"மகாபலிபுரம்":      "Chengalpattu",
	"மாமல்லபுரம்":      "Chengalpattu",
	"chengalpattu":    "Chengalpattu",

	"திருவள்ளூர்":       "Tiruvallur",
	"பொன்னேரி":         "Tiruvallur",
	"திருத்தணி":        "Tiruvallur",
	"tiruvallur":      "Tiruvallur",

	"புதுக்கோட்டை":     "Pudukkottai",
	"அறந்தாங்கி":        "Pudukkottai",
	"pudukkottai":     "Pudukkottai",

	"சிவகங்கை":         "Sivaganga",
	"காரைக்குடி":       "Sivaganga",
	"தேவகோட்டை":       "Sivaganga",
	"sivaganga":       "Sivaganga",
	"karaikudi":       "Sivaganga",

	"ராமநாதபுரம்":      "Ramanathapuram",
	"ராமேஸ்வரம்":      "Ramanathapuram",
	"பரமக்குடி":        "Ramanathapuram",
	"ramanathapuram":  "Ramanathapuram",
	"rameswaram":      "Ramanathapuram",

	"நாகப்பட்டினம்":     "Nagapattinam",
	"நாகை":            "Nagapattinam",
	"வேதாரண்யம்":      "Nagapattinam",
	"nagapattinam":    "Nagapattinam",

	"ராணிப்பேட்டை":     "Ranipet",
	"ஆற்காடு":          "Ranipet",
	"ranipet":         "Ranipet",

	"திருப்பத்தூர்":      "Tirupattur",
	"வாணியம்பாடி":      "Tirupattur",
	"ஆம்பூர்":          "Tirupattur",
	"tirupattur":      "Tirupattur",
	"ambur":           "Tirupattur",

	"திருவண்ணாமலை":     "Tiruvannamalai",
	"ஆரணி":            "Tiruvannamalai",
	"tiruvannamalai":  "Tiruvannamalai",

	"திருவாரூர்":        "Tiruvarur",
	"மன்னார்குடி":       "Tiruvarur",
	"tiruvarur":       "Tiruvarur",

	"கள்ளக்குறிச்சி":    "Kallakurichi",
	"kallakurichi":    "Kallakurichi",

	"தருமபுரி":         "Dharmapuri",
	"ஒகேனக்கல்":        "Dharmapuri",
	"dharmapuri":      "Dharmapuri",

	"கிருஷ்ணகிரி":       "Krishnagiri",
	"ஓசூர்":           "Krishnagiri",
	"krishnagiri":     "Krishnagiri",
	"hosur":           "Krishnagiri",

	"தேனி":            "Theni",
	"போடி":            "Theni",
	"பெரியகுளம்":       "Theni",
	"theni":           "Theni",

	"தென்காசி":         "Tenkasi",
	"குற்றாலம்":        "Tenkasi",
	"tenkasi":         "Tenkasi",

	"நாமக்கல்":         "Namakkal",
	"திருச்செங்கோடு":   "Namakkal",
	"namakkal":        "Namakkal",

	"நீலகிரி":          "Nilgiris",
	"ஊட்டி":           "Nilgiris",
	"குன்னூர்":         "Nilgiris",
	"nilgiris":        "Nilgiris",
	"ooty":            "Nilgiris",

	"அரியலூர்":         "Ariyalur",
	"ariyalur":        "Ariyalur",

	"பெரம்பலூர்":       "Perambalur",
	"perambalur":      "Perambalur",

	"விழுப்புரம்":       "Viluppuram",
	"திண்டிவனம்":       "Viluppuram",
	"viluppuram":      "Viluppuram",

	"விருதுநகர்":        "Virudhunagar",
	"சிவகாசி":          "Virudhunagar",
	"ராஜபாளையம்":      "Virudhunagar",
	"virudhunagar":    "Virudhunagar",
	"sivakasi":        "Virudhunagar",

	"மயிலாடுதுறை":      "Mayiladuthurai",
	"mayiladuthurai":  "Mayiladuthurai",

	"கரூர்":            "Karur",
	"karur":           "Karur",
}

func detectDistrict(text string) string {
	lower := strings.ToLower(text)

	// 1. High-priority check for Tamil Nadu districts and town datelines
	for kw, district := range tnTownDistrictMap {
		if strings.Contains(text, kw) || strings.Contains(lower, strings.ToLower(kw)) {
			return district
		}
	}

	// 2. Check for International locations
	for _, kw := range internationalKeywords {
		if strings.Contains(text, kw) || strings.Contains(lower, kw) {
			return "International"
		}
	}

	// 3. Check for National locations (Other states / cities outside TN)
	for _, kw := range nationalKeywords {
		if strings.Contains(text, kw) || strings.Contains(lower, kw) {
			return "National"
		}
	}

	// 4. Default: State-level Tamil Nadu news (NEVER blind Madurai)
	return "Tamil Nadu"
}

func detectCategory(text, urlStr string) string {
	cleanURL := urlStr
	if u, err := url.Parse(urlStr); err == nil {
		cleanURL = u.Path
	}

	// We prioritize the primary title/heading and URL slug over deep body text
	// to avoid false positives on metaphoric words (e.g. 'அரசியல் விளையாட்டு' or 'விளையாட்டு அல்ல')
	firstChunk := text
	if len(text) > 300 {
		firstChunk = text[:300]
	}
	primary := strings.ToLower(firstChunk + " " + cleanURL)

	// 1. Crime & Safety (High priority for police, murder, arrest, accidents)
	crimeKeywords := []string{
		"கொலை", "கைது", "குற்றம்", "மோசடி", "தாக்குதல்", "கொள்ளை", "போலீஸ்", "நீதிமன்றம்", "சிறை",
		"தீண்டாமை", "பலி", "விபத்து", "தற்கொலை", "சடலம்", "ஆணவக்கொலை", "தூக்கு", "விசாரணை",
		"crime", "police", "arrest", "murder", "theft", "scam", "fraud", "court", "jail", "assault", "robbery", "accident", "encounter",
	}
	for _, kw := range crimeKeywords {
		if strings.Contains(primary, kw) {
			return "Crime"
		}
	}

	// 2. Politics & Governance (High priority for elections, parties, assembly, budget, cabinet)
	politicsKeywords := []string{
		"அரசியல்", "தேர்தல்", "திமுக", "அதிமுக", "பாஜக", "தவெக", "காங்கிரஸ்", "கம்யூனிஸ்ட்",
		"கூட்டாட்சி", "சட்டமன்ற", "நாடாளுமன்ற", "ஆளுநர்", "அமைச்சர்", "முதல்வர்", "பிரதமர்",
		"ஜனாதிபதி", "ஆட்சி", "பட்ஜெட்", "எம்.பி", "எம்.எல்.ஏ", "கட்சிகள்", "மத்திய அரசு", "மாநில அரசு",
		"politics", "election", "bjp", "dmk", "aiadmk", "tvk", "congress", "parliament", "assembly",
		"governor", "minister", "budget", "modi", "stalin", "note this point", "pt explainer",
	}
	for _, kw := range politicsKeywords {
		if strings.Contains(primary, kw) {
			return "Politics"
		}
	}

	// 3. Cinema & Entertainment
	cinemaKeywords := []string{
		"சினிமா", "திரைப்பட", "நடிகர்", "நடிகை", "இயக்குநர்", "படம்", "பாடல்", "டிரெய்லர்",
		"டீசர்", "பிக்பாஸ்", "பிக் பாஸ்", "ஓடிடி", "ரிலீஸ்", "பாக்ஸ் ஆபீஸ்", "கோலிவுட்",
		"cinema", "movie", "film", "actor", "actress", "kollywood", "trailer", "teaser", "box office", "biggboss",
	}
	for _, kw := range cinemaKeywords {
		if strings.Contains(primary, kw) {
			return "Entertainment"
		}
	}

	// 4. Sports & Athletics (Must have genuine sporting indicators, avoiding metaphors)
	sportsKeywords := []string{
		"கிரிக்கெட்", "கால்பந்து", "துலீப்", "ஐபிஎல்", "விக்கெட்", "டெஸ்ட் போட்டி", "ஒருநாள் போட்டி",
		"டி20", "செஸ்", "கபடி", "சாம்பியன்ஷிப்", "தடகள", "ஒலிம்பிக்", "டென்னிஸ்", "ஹாக்கி",
		"பேட்மிண்டன்", "விளையாட்டு செய்தி", "விளையாட்டு அரங்கம்", "விளையாட்டுப் போட்டி", "ரேஸிங்",
		"cricket", "football", "ipl", "bcci", "icc", "fifa", "tennis", "hockey", "badminton", "kabaddi", "wimbledon", "athlete", "tournament",
	}
	for _, kw := range sportsKeywords {
		if strings.Contains(primary, kw) {
			return "Sports"
		}
	}
	// Cricket century / runs check (ensuring not weather heat century like 'வெயில் சதம்')
	if (strings.Contains(primary, "ரன்கள்") || strings.Contains(primary, "சதம்")) && !strings.Contains(primary, "வெயில்") && !strings.Contains(primary, "வெப்ப") {
		return "Sports"
	}
	if strings.Contains(primary, "விளையாட்டு") && !strings.Contains(primary, "விளையாட்டு அல்ல") && !strings.Contains(primary, "அரசியல் விளையாட்டு") {
		return "Sports"
	}

	// 5. Business & Economy
	businessKeywords := []string{
		"பொருளாதாரம்", "வணிகம்", "பங்குச்சந்தை", "சென்செக்ஸ்", "நிப்டி", "தங்கம் விலை", "வெள்ளி விலை",
		"ரூபாய் மதிப்பு", "வங்கி", "முதலீடு", "விலை உயர்வு", "பணவீக்கம்",
		"business", "market", "economy", "sensex", "nifty", "gold price", "silver price", "stock market", "sipcot", "rbi", "inflation",
	}
	for _, kw := range businessKeywords {
		if strings.Contains(primary, kw) {
			return "Business"
		}
	}

	// 6. Technology & AI
	techKeywords := []string{
		"தொழில்நுட்ப", "செயற்கை நுண்ணறிவு", "இஸ்ரோ", "நாசா", "விண்கலம்", "செயற்கைக்கோள்", "சாப்ட்வேர்", "ஸ்மார்ட்போன்",
		"technology", "software", "artificial intelligence", "generative ai", "robotics", "isro", "nasa", "satellite", "gadget", "cybersecurity", "crypto", "bitcoin",
	}
	for _, kw := range techKeywords {
		if strings.Contains(primary, kw) {
			return "Technical"
		}
	}

	// 7. Events & Announcements
	eventKeywords := []string{
		"திருவிழா", "நிகழ்வு", "விழா", "ஆலோசனை கூட்டம்", "முகாம்", "கண்காட்சி", "பூஜை",
		"festival", "thiruvizha", "temple", "celebration", "exhibition", "expo", "conference", "jallikattu", "pongal",
	}
	for _, kw := range eventKeywords {
		if strings.Contains(primary, kw) {
			return "Events"
		}
	}

	// 8. Government Schemes
	if strings.Contains(primary, "scheme") || strings.Contains(primary, "govt") || strings.Contains(primary, "திட்டம்") || strings.Contains(primary, "அரசு நலத்திட்டம்") {
		return "Politics"
	}

	return "News"
}

func extractMatchingNavLinks(baseURL string, htmlContent string) []string {
	var matched []string
	seen := make(map[string]bool)

	aTagRegex := regexp.MustCompile(`(?i)<a\s+[^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := aTagRegex.FindAllStringSubmatch(htmlContent, -1)
	baseParsed, _ := url.Parse(baseURL)

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		rawHref := strings.TrimSpace(m[1])
		innerText := cleanHTML(m[2])

		// Skip accessibility "skip to content" links, jump fragments, and UI boilerplate
		if isBoilerplateAnchor(innerText, rawHref) {
			continue
		}

		combined := strings.ToLower(rawHref + " " + innerText)

		isMatch := false
		// Match general news, state, districts
		if strings.Contains(combined, "district") || strings.Contains(combined, "tamil-nadu") || strings.Contains(combined, "tamilnadu") ||
			strings.Contains(combined, "cities") || strings.Contains(combined, "state") ||
			strings.Contains(combined, "news") || strings.Contains(combined, "செய்தி") ||
			strings.Contains(combined, "மாவட்டம்") || strings.Contains(combined, "தமிழ்நாடு") {
			isMatch = true
		}
		// Match national / India
		if strings.Contains(combined, "india") || strings.Contains(combined, "national") ||
			strings.Contains(combined, "இந்தியா") || strings.Contains(combined, "தேசம்") {
			isMatch = true
		}
		// Match international / world
		if strings.Contains(combined, "world") || strings.Contains(combined, "international") ||
			strings.Contains(combined, "உலகம்") || strings.Contains(combined, "உலக") {
			isMatch = true
		}
		// Match cinema / entertainment
		if strings.Contains(combined, "cinema") || strings.Contains(combined, "entertainment") ||
			strings.Contains(combined, "திரைப்படம்") || strings.Contains(combined, "சினிமா") || strings.Contains(combined, "திரை") {
			isMatch = true
		}
		// Match sports
		if strings.Contains(combined, "sport") || strings.Contains(combined, "cricket") ||
			strings.Contains(combined, "விளையாட்டு") || strings.Contains(combined, "கிரிக்கெட்") {
			isMatch = true
		}
		// Match business / economy
		if strings.Contains(combined, "business") || strings.Contains(combined, "economy") || strings.Contains(combined, "finance") ||
			strings.Contains(combined, "வணிகம்") || strings.Contains(combined, "பொருளாதாரம்") {
			isMatch = true
		}
		// Match science & tech
		if strings.Contains(combined, "science") || strings.Contains(combined, "tech") ||
			strings.Contains(combined, "அறிவியல்") || strings.Contains(combined, "தொழில்நுட்பம்") {
			isMatch = true
		}
		// Match health & lifestyle
		if strings.Contains(combined, "health") || strings.Contains(combined, "lifestyle") ||
			strings.Contains(combined, "உடல்நலம்") || strings.Contains(combined, "வாழ்வியல்") || strings.Contains(combined, "மருத்துவம்") {
			isMatch = true
		}
		// Match video / watch
		if strings.Contains(combined, "watch") || strings.Contains(combined, "video") ||
			strings.Contains(combined, "காணொளி") || strings.Contains(combined, "வீடியோ") {
			isMatch = true
		}
		// Match topics / categories
		if strings.Contains(combined, "/topics/") || strings.Contains(combined, "/category/") || strings.Contains(combined, "/categories/") ||
			strings.Contains(combined, "பிரிவுகள்") {
			isMatch = true
		}
		// Match events & spiritual
		if strings.Contains(combined, "event") || strings.Contains(combined, "festival") ||
			strings.Contains(combined, "aanmeegam") || strings.Contains(combined, "திருவிழா") || strings.Contains(combined, "நிகழ்வு") || strings.Contains(combined, "ஆன்மீகம்") {
			isMatch = true
		}
		// Match specific districts in URL
		for _, d := range tnDistricts {
			if strings.Contains(combined, strings.ToLower(d)) {
				isMatch = true
				break
			}
		}

		if !isMatch {
			continue
		}

		fullURL := rawHref
		if baseParsed != nil {
			ref, err := url.Parse(rawHref)
			if err == nil {
				resolved := baseParsed.ResolveReference(ref)
				if resolved.Host == baseParsed.Host {
					fullURL = resolved.String()
				} else {
					continue
				}
			}
		}

		if fullURL == baseURL || seen[fullURL] || strings.Contains(fullURL, "#") || strings.Contains(fullURL, "javascript:") ||
			strings.Contains(fullURL, "/login") || strings.Contains(fullURL, "/register") || strings.Contains(fullURL, "/subscribe") {
			continue
		}
		seen[fullURL] = true
		matched = append(matched, fullURL)

		if len(matched) >= 15 {
			break
		}
	}

	return matched
}

func extractYouTubeID(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	if u.Host == "youtu.be" {
		return strings.TrimPrefix(u.Path, "/")
	}
	return u.Query().Get("v")
}

var (
	unclosedHexEntityRegex = regexp.MustCompile(`(?i)&#x([0-9a-f]+);?`)
	unclosedDecEntityRegex = regexp.MustCompile(`&#([0-9]+);?`)
	htmlTagRegex           = regexp.MustCompile(`<[^>]*>`)
)

// CleanHTML normalizes, cleans, and decodes HTML entities (including unclosed hex entities like &#x27; and &#x27)
func CleanHTML(s string) string {
	if s == "" {
		return ""
	}
	// 1. Decode numeric hex entities (both closed like &#x27; and unclosed like &#x27)
	s = unclosedHexEntityRegex.ReplaceAllStringFunc(s, func(m string) string {
		hexStr := strings.TrimPrefix(strings.TrimPrefix(m, "&#x"), "&#X")
		hexStr = strings.TrimSuffix(hexStr, ";")
		if val, err := strconv.ParseInt(hexStr, 16, 32); err == nil && val > 0 {
			return string(rune(val))
		}
		return m
	})
	// 2. Decode numeric decimal entities (both closed like &#39; and unclosed like &#39)
	s = unclosedDecEntityRegex.ReplaceAllStringFunc(s, func(m string) string {
		decStr := strings.TrimPrefix(m, "&#")
		decStr = strings.TrimSuffix(decStr, ";")
		if val, err := strconv.ParseInt(decStr, 10, 32); err == nil && val > 0 {
			return string(rune(val))
		}
		return m
	})
	// 3. Standard library unescape: handles &quot;, &amp;, &apos;, &lt;, &gt;, &eacute;, etc.
	s = html.UnescapeString(s)

	// 4. Common named entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&ndash;", "–")
	s = strings.ReplaceAll(s, "&mdash;", "—")
	s = strings.ReplaceAll(s, "&lsquo;", "‘")
	s = strings.ReplaceAll(s, "&rsquo;", "’")
	s = strings.ReplaceAll(s, "&ldquo;", "“")
	s = strings.ReplaceAll(s, "&rdquo;", "”")
	s = strings.ReplaceAll(s, "&hellip;", "…")
	s = strings.ReplaceAll(s, "&bull;", "•")

	// 5. Strip any HTML tags
	s = htmlTagRegex.ReplaceAllString(s, "")

	// 6. Clean up extra whitespace
	fields := strings.Fields(s)
	return strings.TrimSpace(strings.Join(fields, " "))
}

func cleanHTML(s string) string {
	return CleanHTML(s)
}
