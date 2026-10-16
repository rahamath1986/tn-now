package admin

import (
	"tn-now/server/internal/scraper"
)

// Re-export rewriter types and constructors from scraper package
type ArticleRewriteRequest = scraper.ArticleRewriteRequest
type RewrittenArticle = scraper.RewrittenArticle
type ArticleRewriter = scraper.ArticleRewriter

var NewArticleRewriter = scraper.NewArticleRewriter
var GetDefaultRewriter = scraper.GetDefaultRewriter
