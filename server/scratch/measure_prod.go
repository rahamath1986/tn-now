package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Measurement struct {
	URL            string
	Status         int
	TTFB           time.Duration
	TotalTime      time.Duration
	RawHTMLBytes   int
	GzipBytes      int
	RequestCount   int
	CacheControl   string
	HeroImageFound bool
	HeroLCPAttrs   string
	ArticlePreload string
}

func gzipSize(data []byte) int {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, _ = w.Write(data)
	_ = w.Close()
	return b.Len()
}

func measureURL(u string) Measurement {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	start := time.Now()
	resp, err := client.Do(req)
	ttfb := time.Since(start)

	if err != nil {
		return Measurement{URL: u, Status: 0, TTFB: ttfb}
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	totalTime := time.Since(start)

	m := Measurement{
		URL:          u,
		Status:       resp.StatusCode,
		TTFB:         ttfb,
		TotalTime:    totalTime,
		RawHTMLBytes: len(bodyBytes),
		GzipBytes:    gzipSize(bodyBytes),
		CacheControl: resp.Header.Get("Cache-Control"),
	}

	htmlStr := string(bodyBytes)

	// Count scripts
	scripts := strings.Count(htmlStr, "<script")
	styles := strings.Count(htmlStr, "<style") + strings.Count(htmlStr, "rel=\"stylesheet\"")
	imgs := strings.Count(htmlStr, "<img")

	m.RequestCount = 1 + scripts + styles + imgs

	// Check Hero image
	if strings.Contains(htmlStr, `id="heroImage"`) {
		m.HeroImageFound = true
		idx := strings.Index(htmlStr, `id="heroImage"`)
		end := idx + 200
		if end > len(htmlStr) {
			end = len(htmlStr)
		}
		m.HeroLCPAttrs = htmlStr[idx:end]
	}

	// Check preload in head
	if strings.Contains(htmlStr, `rel="preload" as="image"`) {
		m.ArticlePreload = "Found"
	}

	return m
}

func main() {
	urls := []string{
		"https://www.tn24.in/",
		"https://www.tn24.in/?post=41672925-f198-4895-bafb-293335fbcda9",
		"https://www.tn24.in/district/chennai",
		"https://www.tn24.in/district/madurai",
		"https://www.tn24.in/category/politics",
		"https://www.tn24.in/assets/brand/tn24-logo.svg?v=20260908d",
		"https://www.tn24.in/favicon.svg",
	}

	for _, u := range urls {
		m := measureURL(u)
		fmt.Printf("=== %s ===\n", m.URL)
		fmt.Printf("Status: %d | TTFB: %v | Total: %v\n", m.Status, m.TTFB, m.TotalTime)
		fmt.Printf("HTML Raw: %d bytes (%.1f KB) | Gzip: %d bytes (%.1f KB)\n", m.RawHTMLBytes, float64(m.RawHTMLBytes)/1024, m.GzipBytes, float64(m.GzipBytes)/1024)
		fmt.Printf("Cache-Control: %s\n", m.CacheControl)
		fmt.Printf("Estimated Sub-requests: %d\n", m.RequestCount)
		if m.HeroImageFound {
			fmt.Printf("Hero Image snippet: %s\n", m.HeroLCPAttrs)
		}
		if m.ArticlePreload != "" {
			fmt.Printf("Article Preload: %s\n", m.ArticlePreload)
		}
		fmt.Println()
	}
}
