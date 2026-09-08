package portal

func RenderPortalPage() string {
	return `<!DOCTYPE html>
<html lang="ta">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TN24 &mdash; தமிழ்நாட்டின் முதன்மை 24/7 செய்தி தளம் | Tamil Nadu Breaking News Live Today</title>

    <!-- ═══════════════════════════════════════════════════════
         CORE SEO META TAGS — TN24 Tamil News Portal
         Keywords researched for max search volume in Tamil Nadu
    ═══════════════════════════════════════════════════════ -->
    <meta name="description" content="TN24 — Tamil Nadu's #1 live news portal. Breaking news today in Tamil, live updates from all 38 districts: Chennai, Coimbatore, Madurai, Salem, Trichy. Politics, Cinema, Sports, TNPSC, Government Jobs & more. 24/7 உண்மைச் செய்திகள்.">

    <!-- Expanded High-Volume Keywords (Tanglish + Tamil Script + English) -->
    <meta name="keywords" content="
        TN24, TN24 News, TN24 Live, tn24now,
        Tamil news, Tamil Nadu news, Tamil news live, Tamil news today,
        Breaking news Tamil, breaking news today Tamil Nadu,
        Tamil Nadu breaking news, live news Tamil Nadu,
        தமிழ்நாடு செய்திகள், தமிழ் செய்திகள், இன்றைய செய்திகள்,
        நேரலை செய்திகள், சமீபத்திய செய்திகள், முக்கிய செய்திகள்,
        Tamil news live today, latest Tamil news, today Tamil news,
        Chennai news, Chennai latest news, Chennai breaking news,
        சென்னை செய்திகள், சென்னை தமிழ் செய்திகள்,
        Coimbatore news, Madurai news, Salem news, Trichy news, Tirunelveli news,
        Erode news, Vellore news, Thanjavur news, Tiruppur news, Dindigul news,
        கோவை செய்திகள், மதுரை செய்திகள், சேலம் செய்திகள், திருச்சி செய்திகள்,
        Tamil Nadu politics, Tamil Nadu political news, DMK news, AIADMK news,
        தமிழக அரசியல், அரசியல் செய்திகள், தமிழ்நாடு அரசு செய்திகள்,
        Tamil cinema news, kollywood news, Tamil movie news, new Tamil movies,
        தமிழ் சினிமா செய்திகள், கொல்லிவுட் செய்திகள், திரை செய்திகள்,
        Tamil sports news, IPL Tamil news, cricket Tamil, CSK news,
        விளையாட்டு செய்திகள், கிரிக்கெட் செய்திகள்,
        TNPSC, TNPSC news, Tamil Nadu government jobs, TN govt jobs,
        தமிழக அரசு வேலைவாய்ப்பு, அரசு வேலை செய்திகள்,
        Tamil weather news, Chennai weather, Tamil Nadu weather,
        வானிலை அறிக்கை, தமிழ்நாடு வானிலை,
        Tamil business news, Tamil Nadu economy, share market Tamil,
        வணிக செய்திகள், பங்குச் சந்தை செய்திகள்,
        Tamil technology news, tech news Tamil, mobile news Tamil,
        தொழில்நுட்ப செய்திகள்,
        Tamil education news, school news Tamil Nadu, university news Tamil,
        கல்வி செய்திகள்,
        Tamil crime news, police news Tamil Nadu, court news Tamil,
        குற்றம் செய்திகள், நீதிமன்ற செய்திகள்,
        Tamil agriculture news, farmer news Tamil Nadu, rain news Tamil,
        விவசாய செய்திகள், மழை செய்திகள்,
        Pongal news, Tamil festival news, Tamil calendar 2025,
        பொங்கல் செய்திகள், தமிழ் திருவிழா,
        TN24 digital news, TN24 portal, TN24 online
    ">

    <meta name="author" content="TN24 News Media Network">
    <meta name="publisher" content="TN24 Digital News">
    <meta name="copyright" content="TN24 News Media Network 2024">
    <meta name="robots" content="index, follow, max-snippet:-1, max-image-preview:large, max-video-preview:-1">
    <meta name="googlebot" content="index, follow, max-snippet:-1, max-image-preview:large, max-video-preview:-1">
    <meta name="googlebot-news" content="index, follow">

    <!-- Language & Locale -->
    <meta name="language" content="Tamil">
    <meta http-equiv="content-language" content="ta, en-IN">
    <link rel="alternate" hreflang="ta" href="https://tn24now.in/portal">
    <link rel="alternate" hreflang="en-IN" href="https://tn24now.in/portal">
    <link rel="alternate" hreflang="x-default" href="https://tn24now.in/portal">

    <!-- Geographic / Local SEO -->
    <meta name="geo.region" content="IN-TN">
    <meta name="geo.placename" content="Tamil Nadu, India">
    <meta name="geo.position" content="11.1271;78.6569">
    <meta name="ICBM" content="11.1271, 78.6569">

    <!-- Google News Specific -->
    <meta name="news_keywords" content="Tamil Nadu news, breaking news Tamil, Chennai news, TNPSC, Tamil cinema, Tamil politics, DMK, AIADMK, IPL, cricket Tamil, government jobs Tamil Nadu">
    <meta name="syndication-source" content="https://tn24now.in/portal">
    <meta name="original-source" content="https://tn24now.in/portal">

    <!-- Canonical URL -->
    <link rel="canonical" href="https://tn24now.in/portal">

    <!-- Open Graph (Facebook / WhatsApp / LinkedIn / Telegram) -->
    <meta property="og:type" content="website">
    <meta property="og:site_name" content="TN24 News">
    <meta property="og:locale" content="ta_IN">
    <meta property="og:locale:alternate" content="en_IN">
    <meta property="og:title" content="TN24 — Tamil Nadu's #1 Live News Portal | 24/7 Breaking Tamil News">
    <meta property="og:description" content="Get live breaking news from Tamil Nadu — Politics, Cinema, Sports, TNPSC, Government Jobs, Weather & more. 38 மாவட்டங்களின் உடனடி செய்திகள் 24/7.">
    <meta property="og:url" content="https://tn24now.in/portal">
    <meta property="og:image" content="https://tn24now.in/portal/assets/brand/tn24-profile.jpg">
    <meta property="og:image:width" content="1024">
    <meta property="og:image:height" content="1024">
    <meta property="og:image:type" content="image/jpeg">
    <meta property="og:image:alt" content="TN24 — Tamil Nadu Live News Network">

    <!-- Twitter / X Card Meta Tags -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:site" content="@TN24Now">
    <meta name="twitter:creator" content="@TN24Now">
    <meta name="twitter:title" content="TN24 — Tamil Nadu Breaking News Live 24/7">
    <meta name="twitter:description" content="Live Tamil Nadu news — Politics, Cinema, Sports, TNPSC & Government Jobs. 38 districts covered 24/7. தமிழ்நாட்டின் #1 செய்தி தளம்.">
    <meta name="twitter:image" content="https://tn24now.in/portal/assets/brand/tn24-profile.jpg">
    <meta name="twitter:image:alt" content="TN24 Tamil News Network">

    <!-- Structured Data: Schema.org JSON-LD -->
    <script type="application/ld+json">
    {
      "@context": "https://schema.org",
      "@graph": [
        {
          "@type": "NewsMediaOrganization",
          "@id": "https://tn24now.in/#organization",
          "name": "TN24",
          "alternateName": ["TN24 News", "TN24 Tamil News", "TN24 Now"],
          "description": "Tamil Nadu's #1 digital news portal delivering 24/7 live breaking news in Tamil and English from all 38 districts.",
          "url": "https://tn24now.in/portal",
          "logo": {
            "@type": "ImageObject",
            "url": "https://tn24now.in/portal/assets/brand/tn24-profile.jpg",
            "width": 1024,
            "height": 1024
          },
          "foundingDate": "2024",
          "areaServed": {
            "@type": "State",
            "name": "Tamil Nadu",
            "containedInPlace": {
              "@type": "Country",
              "name": "India"
            }
          },
          "contactPoint": {
            "@type": "ContactPoint",
            "telephone": "+91-8124395082",
            "contactType": "News Desk",
            "email": "tn24now@gmail.com",
            "areaServed": "IN-TN",
            "availableLanguage": ["Tamil", "English"]
          },
          "sameAs": [
            "https://twitter.com/TN24Now",
            "https://facebook.com/TN24Now",
            "https://instagram.com/TN24Now",
            "https://youtube.com/@TN24Now"
          ]
        },
        {
          "@type": "WebSite",
          "@id": "https://tn24now.in/#website",
          "url": "https://tn24now.in/portal",
          "name": "TN24 — Tamil Nadu Breaking News 24/7",
          "description": "Live Tamil Nadu news portal covering breaking news, politics, cinema, sports, TNPSC, government jobs from all 38 districts.",
          "inLanguage": ["ta", "en-IN"],
          "publisher": {
            "@id": "https://tn24now.in/#organization"
          },
          "potentialAction": {
            "@type": "SearchAction",
            "target": "https://tn24now.in/portal?q={search_term_string}",
            "query-input": "required name=search_term_string"
          }
        },
        {
          "@type": "BreadcrumbList",
          "itemListElement": [
            {
              "@type": "ListItem",
              "position": 1,
              "name": "TN24 Home",
              "item": "https://tn24now.in/portal"
            }
          ]
        },
        {
          "@type": "SiteNavigationElement",
          "name": ["Latest News", "Politics", "Cinema", "Sports", "Business", "Technology", "District News", "Government Jobs"],
          "url": [
            "https://tn24now.in/portal",
            "https://tn24now.in/portal#politics",
            "https://tn24now.in/portal#cinema",
            "https://tn24now.in/portal#sports",
            "https://tn24now.in/portal#business",
            "https://tn24now.in/portal#technology",
            "https://tn24now.in/portal#districts",
            "https://tn24now.in/portal#government-jobs"
          ]
        }
      ]
    }
    </script>

    <link rel="icon" type="image/svg+xml" href="/portal/assets/brand/tn24-icon.svg">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@600;700&family=Inter:wght@400;500;600;700;800&family=Merriweather:ital,wght@0,300;0,400;0,700;1,400&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-body: #f0f2f5;
            --bg-card: #ffffff;
            --navy-header: #0a1118;
            --navy-dark: #0f1c2e;
            --nav-blue: #1557bf;
            --nav-blue-hover: #104499;
            --nav-blue-active: #0c3577;
            --badge-blue: #1c75bc;
            --red-subscribe: #ef3e36;
            --red-subscribe-hover: #d32f2f;
            --text-dark: #1e293b;
            --text-muted: #64748b;
            --text-light: #94a3b8;
            --border-color: #e2e8f0;
            --border-subtle: #edf2f7;
            --green-trend: #10b981;
            --font-main: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            --font-headline: 'Inter', -apple-system, sans-serif;
            --font-serif: 'Merriweather', Georgia, serif;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: var(--font-main);
            background-color: var(--bg-body);
            color: var(--text-dark);
            line-height: 1.45;
            -webkit-font-smoothing: antialiased;
        }

        a { text-decoration: none; color: inherit; transition: color 0.15s; }
        a:hover { color: var(--nav-blue); }
        button { cursor: pointer; font-family: inherit; border: none; }

        /* ---------------- TOP BANNER ---------------- */
        .top-banner {
            background-color: #070d18;
            background-image: 
                radial-gradient(rgba(56, 189, 248, 0.12) 1px, transparent 1px),
                linear-gradient(180deg, #090e17 0%, #111d2e 100%);
            color: #cbd5e1;
            border-bottom: 1px solid #1e293b;
            position: relative;
            background-size: 24px 24px, 100% 100%;
        }
        .top-banner-inner {
            max-width: 1320px;
            margin: 0 auto;
            padding: 14px 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .top-left {
            display: flex;
            align-items: center;
            gap: 16px;
            font-size: 11px;
            font-weight: 600;
            letter-spacing: 0.5px;
            text-transform: uppercase;
            color: #94a3b8;
        }
        .social-icons {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .social-icons a {
            color: #94a3b8;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 22px;
            height: 22px;
            border-radius: 4px;
            transition: all 0.2s;
        }
        .social-icons a:hover {
            color: #ffffff;
            background: rgba(255,255,255,0.1);
        }

        /* BRAND LOGO */
        .brand-logo {
            display: flex;
            align-items: center;
            gap: 12px;
            color: #ffffff;
        }
        .logo-symbol {
            width: 38px;
            height: 38px;
            border-radius: 50%;
            background: conic-gradient(from 180deg, #38bdf8, #818cf8, #f43f5e, #38bdf8);
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 2px;
            box-shadow: 0 0 15px rgba(56, 189, 248, 0.4);
        }
        .logo-symbol-inner {
            width: 100%;
            height: 100%;
            background: #090e17;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .logo-symbol-inner svg {
            width: 20px;
            height: 20px;
            fill: none;
            stroke: #38bdf8;
            stroke-width: 2.5;
        }
        .brand-text {
            font-size: 26px;
            font-weight: 800;
            letter-spacing: 1.5px;
            color: #ffffff;
            font-family: 'Inter', sans-serif;
            text-transform: lowercase;
            display: flex;
            align-items: baseline;
            gap: 6px;
        }
        .brand-text span {
            font-size: 10px;
            text-transform: uppercase;
            letter-spacing: 1px;
            background: rgba(56, 189, 248, 0.15);
            color: #38bdf8;
            padding: 2px 6px;
            border-radius: 4px;
            font-weight: 700;
        }

        .top-right {
            display: flex;
            align-items: center;
            gap: 20px;
            font-size: 11px;
            font-weight: 600;
        }
        .top-user-links {
            display: flex;
            align-items: center;
            gap: 12px;
            color: #94a3b8;
            letter-spacing: 0.5px;
        }
        .top-user-links a:hover { color: #ffffff; }
        .lang-select {
            background: transparent;
            color: #94a3b8;
            border: none;
            font-size: 11px;
            font-weight: 600;
            cursor: pointer;
            outline: none;
        }
        .lang-select option { background: #111d2e; color: #fff; }
        .btn-subscribe {
            background: var(--red-subscribe);
            color: #ffffff;
            font-size: 12px;
            font-weight: 700;
            padding: 8px 18px;
            border-radius: 4px;
            display: flex;
            align-items: center;
            gap: 6px;
            letter-spacing: 0.4px;
            transition: all 0.2s;
            box-shadow: 0 4px 12px rgba(239, 62, 54, 0.35);
        }
        .btn-subscribe:hover {
            background: var(--red-subscribe-hover);
            transform: translateY(-1px);
        }

        /* ---------------- PRIMARY NAV BAR ---------------- */
        .primary-nav {
            background: var(--nav-blue);
            color: #ffffff;
            box-shadow: 0 2px 8px rgba(0,0,0,0.12);
            position: sticky;
            top: 0;
            z-index: 100;
        }
        .nav-inner {
            max-width: 1320px;
            margin: 0 auto;
            padding: 0 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .nav-links {
            display: flex;
            align-items: center;
            list-style: none;
        }
        .nav-item {
            position: relative;
        }
        .nav-link {
            display: inline-flex;
            align-items: center;
            gap: 5px;
            color: #ffffff;
            font-size: 12px;
            font-weight: 700;
            letter-spacing: 0.8px;
            text-transform: uppercase;
            padding: 14px 16px;
            transition: background 0.15s;
            cursor: pointer;
        }
        .nav-link:hover, .nav-item.active .nav-link {
            background: var(--nav-blue-hover);
            color: #ffffff;
        }
        .nav-link svg {
            width: 10px;
            height: 10px;
            opacity: 0.8;
            stroke-width: 3;
        }

        /* Nav Dropdown Menu */
        .nav-dropdown {
            position: absolute;
            top: 100%;
            left: 0;
            background: #ffffff;
            min-width: 180px;
            box-shadow: 0 8px 24px rgba(0,0,0,0.15);
            border-radius: 0 0 6px 6px;
            display: none;
            flex-direction: column;
            padding: 6px 0;
            z-index: 150;
            border-top: 2px solid var(--nav-blue);
        }
        .nav-item:hover .nav-dropdown { display: flex; }
        .nav-dropdown a {
            padding: 9px 18px;
            font-size: 12px;
            font-weight: 600;
            color: #334155;
            transition: background 0.15s;
        }
        .nav-dropdown a:hover {
            background: #f1f5f9;
            color: var(--nav-blue);
        }

        .nav-search-box {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .nav-search-input {
            background: rgba(255,255,255,0.15);
            border: 1px solid rgba(255,255,255,0.25);
            color: #ffffff;
            font-size: 12px;
            padding: 6px 12px;
            border-radius: 4px;
            outline: none;
            width: 140px;
            transition: width 0.2s, background 0.2s;
        }
        .nav-search-input::placeholder { color: rgba(255,255,255,0.7); }
        .nav-search-input:focus {
            width: 200px;
            background: rgba(255,255,255,0.25);
            border-color: #ffffff;
        }
        .nav-search-btn {
            background: transparent;
            color: #ffffff;
            padding: 6px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        /* ---------------- MARKET / TICKER BAR ---------------- */
        .market-ticker-bar {
            background: #ffffff;
            border-bottom: 1px solid var(--border-color);
            box-shadow: 0 1px 3px rgba(0,0,0,0.03);
        }
        .ticker-inner {
            max-width: 1320px;
            margin: 0 auto;
            padding: 8px 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 20px;
        }
        .ticker-selector-wrap {
            position: relative;
            min-width: 150px;
        }
        .ticker-select-btn {
            background: #f8fafc;
            border: 1px solid var(--border-color);
            padding: 7px 12px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 700;
            color: #1e293b;
            display: flex;
            align-items: center;
            justify-content: space-between;
            width: 100%;
        }
        .ticker-select-btn svg { width: 12px; height: 12px; }

        /* Sparkline */
        .sparkline-box {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        .sparkline-svg {
            width: 140px;
            height: 28px;
        }

        .ticker-stats {
            display: flex;
            align-items: center;
            gap: 28px;
            flex: 1;
            justify-content: flex-end;
        }
        .ticker-stat-item {
            display: flex;
            flex-direction: column;
            font-size: 11px;
        }
        .stat-label {
            font-size: 9px;
            font-weight: 700;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .stat-val {
            font-size: 12px;
            font-weight: 800;
            color: #0f172a;
        }
        .stat-val.up { color: var(--green-trend); }

        /* ---------------- MAIN 3-COLUMN LAYOUT ---------------- */
        .main-container {
            max-width: 1320px;
            margin: 20px auto 40px auto;
            padding: 0 20px;
            display: grid;
            grid-template-columns: 210px minmax(0, 1fr) 280px;
            gap: 24px;
        }

        /* ---------------- COLUMN 1: LEFT SIDEBAR ---------------- */
        .left-col {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }
        .mini-card {
            display: flex;
            flex-direction: column;
            background: #ffffff;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            overflow: hidden;
            transition: transform 0.15s, box-shadow 0.15s;
            cursor: pointer;
        }
        .mini-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(0,0,0,0.06);
        }
        .mini-thumb-wrap {
            position: relative;
            width: 100%;
            height: 115px;
            background: #e2e8f0;
            overflow: hidden;
        }
        .mini-thumb-wrap img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
        .badge-pill {
            position: absolute;
            bottom: 6px;
            left: 6px;
            background: var(--badge-blue);
            color: #ffffff;
            font-size: 9px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 3px 6px;
            border-radius: 2px;
        }
        .mini-content {
            padding: 10px;
        }
        .mini-title {
            font-size: 12px;
            font-weight: 700;
            color: #0f172a;
            line-height: 1.35;
            margin-bottom: 6px;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
        }
        .mini-meta {
            font-size: 9px;
            font-weight: 600;
            color: var(--text-muted);
            text-transform: uppercase;
        }

        /* Left Press Releases Section */
        .section-header-compact {
            font-size: 11px;
            font-weight: 800;
            text-transform: uppercase;
            letter-spacing: 0.8px;
            color: #0f172a;
            border-bottom: 2px solid #0f172a;
            padding-bottom: 6px;
            margin-top: 4px;
            margin-bottom: 12px;
        }
        .press-list {
            display: flex;
            flex-direction: column;
            gap: 12px;
        }
        .press-item {
            cursor: pointer;
            border-bottom: 1px solid var(--border-subtle);
            padding-bottom: 10px;
        }
        .press-title {
            font-size: 11px;
            font-weight: 700;
            color: #1e293b;
            line-height: 1.35;
            margin-bottom: 4px;
        }
        .press-item:hover .press-title { color: var(--nav-blue); }
        .press-meta {
            font-size: 9px;
            font-weight: 600;
            color: var(--text-muted);
            text-transform: uppercase;
        }

        /* Left Banner Ad 200x200 */
        .ad-banner-square {
            width: 100%;
            height: 200px;
            background: linear-gradient(135deg, #a5b4fc 0%, #f472b6 100%);
            border-radius: 4px;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            color: rgba(255,255,255,0.9);
            font-size: 11px;
            font-weight: 700;
            letter-spacing: 1px;
            text-transform: uppercase;
            box-shadow: inset 0 0 20px rgba(0,0,0,0.08);
            position: relative;
        }
        .ad-tag {
            position: absolute;
            top: 6px;
            right: 6px;
            font-size: 8px;
            letter-spacing: 0.5px;
            color: rgba(255,255,255,0.7);
            font-weight: 600;
        }

        /* ---------------- COLUMN 2: CENTER MAIN CONTENT ---------------- */
        .center-col {
            display: flex;
            flex-direction: column;
            gap: 24px;
        }

        /* Lead Hero Article */
        .hero-lead-card {
            background: #ffffff;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            overflow: hidden;
            position: relative;
            box-shadow: 0 4px 16px rgba(0,0,0,0.04);
            cursor: default;
        }
        .hero-image-wrap {
            position: relative;
            width: 100%;
            height: 380px;
            background: #0f172a;
            overflow: hidden;
            cursor: pointer;
        }
        .hero-image-wrap img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
            transition: transform 0.4s ease;
        }
        .hero-lead-card:hover .hero-image-wrap img {
            transform: scale(1.02);
        }
        .hero-overlay {
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            background: linear-gradient(0deg, rgba(10,17,24,0.95) 0%, rgba(10,17,24,0.7) 60%, transparent 100%);
            padding: 30px 24px 16px 24px;
            color: #ffffff;
        }
        .hero-title {
            font-size: 24px;
            font-weight: 800;
            line-height: 1.25;
            color: #ffffff;
            margin-bottom: 8px;
            text-shadow: 0 2px 4px rgba(0,0,0,0.6);
        }
        .hero-desc {
            font-size: 13px;
            color: #cbd5e1;
            line-height: 1.45;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
        }

        /* Hero Sub-Teasers (3 items strip) */
        .hero-teasers-strip {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            background: #111d2e;
            border-top: 1px solid rgba(255,255,255,0.1);
        }
        .hero-teaser-item {
            padding: 14px 16px;
            border-right: 1px solid rgba(255,255,255,0.08);
            color: #ffffff;
            cursor: pointer;
            transition: background 0.15s;
        }
        .hero-teaser-item:last-child { border-right: none; }
        .hero-teaser-item:hover { background: rgba(255,255,255,0.05); }
        .teaser-badge-meta {
            display: flex;
            align-items: center;
            gap: 6px;
            margin-bottom: 6px;
        }
        .teaser-badge {
            background: var(--badge-blue);
            color: #ffffff;
            font-size: 8px;
            font-weight: 700;
            padding: 2px 5px;
            border-radius: 2px;
            text-transform: uppercase;
        }
        .teaser-date {
            font-size: 9px;
            color: #94a3b8;
            font-weight: 600;
            text-transform: uppercase;
        }
        .teaser-title {
            font-size: 11px;
            font-weight: 700;
            line-height: 1.35;
            color: #f1f5f9;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
        }

        /* Horizontal Leaderboard Banner (728x90) */
        .ad-banner-horizontal {
            width: 100%;
            height: 90px;
            background: linear-gradient(90deg, #a78bfa 0%, #f472b6 50%, #fb923c 100%);
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ffffff;
            font-size: 13px;
            font-weight: 800;
            letter-spacing: 2px;
            text-transform: uppercase;
            box-shadow: inset 0 0 20px rgba(0,0,0,0.06);
            position: relative;
        }

        /* 2-Column Article Grid */
        .editorial-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 24px;
        }
        .article-card {
            background: #ffffff;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            overflow: hidden;
            display: flex;
            flex-direction: column;
            transition: transform 0.15s, box-shadow 0.15s;
            cursor: pointer;
        }
        .article-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 18px rgba(0,0,0,0.06);
        }
        .article-card-thumb {
            position: relative;
            width: 100%;
            height: 160px;
            background: #e2e8f0;
            overflow: hidden;
        }
        .article-card-thumb img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
        }
        .article-card-body {
            padding: 14px;
            display: flex;
            flex-direction: column;
            flex: 1;
        }
        .article-card-title {
            font-size: 14px;
            font-weight: 800;
            color: #0f172a;
            line-height: 1.35;
            margin-bottom: 8px;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
        }
        .article-card-meta {
            font-size: 9px;
            font-weight: 700;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.3px;
            margin-bottom: 8px;
        }
        .article-card-desc {
            font-size: 12px;
            color: #475569;
            line-height: 1.45;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            margin-bottom: 12px;
            flex: 1;
        }

        /* Mid-stream Ad Slot (336x280) */
        .ad-banner-mid {
            background: linear-gradient(135deg, #93c5fd 0%, #c084fc 100%);
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ffffff;
            font-size: 12px;
            font-weight: 800;
            letter-spacing: 1.5px;
            text-transform: uppercase;
            min-height: 260px;
            position: relative;
        }

        /* ---------------- COLUMN 3: RIGHT SIDEBAR ---------------- */
        .right-col {
            display: flex;
            flex-direction: column;
            gap: 24px;
        }

        /* Calculator & Live Market Rates Widget */
        .calculator-widget {
            background: #0d1a2d;
            color: #ffffff;
            border-radius: 8px;
            padding: 16px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.18);
            border: 1px solid #1e293b;
        }
        .widget-title-navy {
            font-size: 11px;
            font-weight: 800;
            text-transform: uppercase;
            letter-spacing: 0.8px;
            color: #38bdf8;
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 12px;
        }
        .market-pulse-badge {
            font-size: 9px;
            font-weight: 700;
            color: #4ade80;
            background: rgba(74, 222, 128, 0.12);
            border: 1px solid rgba(74, 222, 128, 0.3);
            border-radius: 10px;
            padding: 2px 7px;
            display: inline-flex;
            align-items: center;
            gap: 4px;
        }
        .market-pulse-dot {
            width: 6px;
            height: 6px;
            background: #4ade80;
            border-radius: 50%;
            animation: pulseDot 1.5s infinite;
        }
        @keyframes pulseDot {
            0%, 100% { opacity: 1; transform: scale(1); }
            50% { opacity: 0.4; transform: scale(0.8); }
        }
        .market-ticker-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 6px;
            margin-bottom: 14px;
        }
        .market-ticker-card {
            background: rgba(15, 23, 42, 0.7);
            border: 1px solid rgba(56, 189, 248, 0.15);
            border-radius: 6px;
            padding: 7px 9px;
        }
        .market-ticker-label {
            font-size: 9.5px;
            font-weight: 700;
            color: #94a3b8;
            text-transform: uppercase;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .market-ticker-val {
            font-size: 12px;
            font-weight: 800;
            color: #f8fafc;
            margin-top: 2px;
            letter-spacing: 0.2px;
        }
        .calc-rows {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        .calc-row {
            background: #13233a;
            border-radius: 6px;
            padding: 9px 10px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 11px;
            border: 1px solid rgba(56, 189, 248, 0.18);
        }
        .calc-left {
            display: flex;
            align-items: center;
            gap: 6px;
        }
        .calc-input {
            width: 44px;
            background: rgba(15, 23, 42, 0.8);
            border: 1px solid #334155;
            border-radius: 4px;
            color: #ffffff;
            font-size: 12px;
            font-weight: 700;
            text-align: center;
            outline: none;
            padding: 3px 4px;
        }
        .calc-input:focus {
            border-color: #38bdf8;
            box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
        }
        .calc-name {
            font-weight: 700;
            color: #38bdf8;
        }
        .calc-eq {
            color: #94a3b8;
            font-size: 11px;
            font-weight: 600;
        }
        .calc-val {
            font-weight: 800;
            color: #f8fafc;
            font-size: 12.5px;
            letter-spacing: 0.3px;
        }
        .calc-unit-select {
            background: #0f1f33;
            color: #38bdf8;
            border: 1px solid rgba(56, 189, 248, 0.3);
            font-size: 10px;
            font-weight: 700;
            border-radius: 4px;
            padding: 3px 6px;
            cursor: pointer;
            outline: none;
        }
        .calc-unit-select:focus {
            border-color: #38bdf8;
        }

        /* ICO / Event Calendar Widget */
        .widget-box {
            background: #ffffff;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            padding: 16px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.02);
        }
        .widget-title-white {
            font-size: 11px;
            font-weight: 800;
            text-transform: uppercase;
            letter-spacing: 0.8px;
            color: #0f172a;
            border-bottom: 2px solid #0f172a;
            padding-bottom: 6px;
            margin-bottom: 14px;
        }
        .event-list {
            display: flex;
            flex-direction: column;
            gap: 14px;
        }
        .event-item {
            display: flex;
            align-items: flex-start;
            gap: 12px;
        }
        .event-logo-badge {
            width: 38px;
            height: 38px;
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ffffff;
            font-weight: 800;
            font-size: 14px;
            flex-shrink: 0;
        }
        .event-logo-hdac { background: #002c6c; }
        .event-logo-lion { background: #f59e0b; }
        .event-logo-shield { background: #ef4444; }
        .event-info {
            display: flex;
            flex-direction: column;
        }
        .event-date {
            font-size: 9px;
            font-weight: 700;
            color: var(--text-muted);
            text-transform: uppercase;
        }
        .event-name {
            font-size: 12px;
            font-weight: 800;
            color: #0f172a;
            line-height: 1.3;
        }
        .event-sub {
            font-size: 10px;
            color: #64748b;
            line-height: 1.35;
        }
        .btn-view-all {
            margin-top: 14px;
            width: 100%;
            background: #f1f5f9;
            color: #334155;
            font-size: 10px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 8px 12px;
            border-radius: 4px;
            text-align: center;
            display: block;
            transition: all 0.15s;
        }
        .btn-view-all:hover {
            background: #e2e8f0;
            color: var(--nav-blue);
        }

        /* 250x250 Ad Banner */
        .ad-banner-250 {
            width: 100%;
            height: 250px;
            background: linear-gradient(135deg, #93c5fd 0%, #f472b6 100%);
            border-radius: 4px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ffffff;
            font-size: 12px;
            font-weight: 800;
            letter-spacing: 1px;
            text-transform: uppercase;
            position: relative;
        }

        /* Most Read (Numbered 1-5) */
        .most-read-list {
            display: flex;
            flex-direction: column;
            gap: 12px;
        }
        .most-read-item {
            display: flex;
            align-items: flex-start;
            gap: 12px;
            cursor: pointer;
        }
        .num-box {
            width: 28px;
            height: 28px;
            background: #f1f5f9;
            color: #475569;
            font-size: 13px;
            font-weight: 800;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 3px;
            flex-shrink: 0;
            border: 1px solid var(--border-color);
        }
        .most-read-title {
            font-size: 12px;
            font-weight: 700;
            color: #1e293b;
            line-height: 1.35;
        }
        .most-read-item:hover .most-read-title { color: var(--nav-blue); }

        /* Latest Tweets */
        .tweet-card {
            background: #f8fafc;
            border: 1px solid var(--border-subtle);
            border-radius: 4px;
            padding: 12px;
        }
        .tweet-text {
            font-size: 11px;
            color: #334155;
            line-height: 1.45;
            margin-bottom: 8px;
        }
        .tweet-meta {
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 9px;
            font-weight: 700;
            color: var(--text-muted);
            text-transform: uppercase;
        }
        .tweet-actions {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .tweet-actions a { color: var(--nav-blue); }

        /* ---------------- ARTICLE READER MODAL ---------------- */
        .modal-backdrop {
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(10, 17, 24, 0.75);
            backdrop-filter: blur(4px);
            z-index: 999;
            display: none;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .modal-card {
            background: #ffffff;
            border-radius: 8px;
            max-width: 780px;
            width: 100%;
            max-height: 90vh;
            overflow-y: auto;
            box-shadow: 0 20px 50px rgba(0,0,0,0.3);
            display: flex;
            flex-direction: column;
            position: relative;
            animation: modalFadeIn 0.2s ease-out;
        }
        @keyframes modalFadeIn {
            from { opacity: 0; transform: scale(0.96); }
            to { opacity: 1; transform: scale(1); }
        }
        .modal-header {
            padding: 20px 24px 14px 24px;
            border-bottom: 1px solid var(--border-color);
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 16px;
        }
        .modal-close-btn {
            background: #f1f5f9;
            color: #475569;
            border-radius: 50%;
            width: 32px;
            height: 32px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 18px;
            font-weight: 700;
            transition: all 0.15s;
        }
        .modal-close-btn:hover { background: #e2e8f0; color: #000; }
        .modal-body {
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 18px;
        }
        .modal-media {
            width: 100%;
            max-height: 380px;
            background: #0f172a;
            border-radius: 6px;
            overflow: hidden;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .modal-media img {
            width: 100%;
            max-height: 380px;
            object-fit: cover;
        }
        .modal-media iframe {
            width: 100%;
            height: 380px;
            border: none;
        }
        .modal-meta-bar {
            display: flex;
            align-items: center;
            gap: 10px;
            flex-wrap: wrap;
            font-size: 11px;
            color: var(--text-muted);
            font-weight: 600;
        }
        .modal-title {
            font-size: 22px;
            font-weight: 800;
            color: #0f172a;
            line-height: 1.3;
        }
        .modal-text {
            font-size: 15px;
            line-height: 1.7;
            color: #334155;
            white-space: pre-wrap;
        }
        .modal-actions {
            padding: 16px 24px;
            background: #f8fafc;
            border-top: 1px solid var(--border-color);
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
        }
        .modal-btn-row {
            display: flex;
            gap: 8px;
            align-items: center;
        }
        .btn-modal-action {
            background: var(--nav-blue);
            color: #ffffff;
            font-size: 12px;
            font-weight: 700;
            padding: 8px 16px;
            border-radius: 4px;
            display: inline-flex;
            align-items: center;
            gap: 6px;
            cursor: pointer;
            border: none;
            transition: all 0.15s;
        }
        .btn-modal-action:hover { background: var(--nav-blue-hover); }
        .btn-share-whatsapp {
            background: #25D366 !important;
            color: #ffffff !important;
        }
        .btn-share-whatsapp:hover { background: #1ebe5d !important; }
        .btn-share-native {
            background: #0284c7 !important;
            color: #ffffff !important;
        }
        .btn-share-native:hover { background: #0369a1 !important; }
        .btn-copy-link {
            background: #475569 !important;
            color: #ffffff !important;
        }
        .btn-copy-link:hover { background: #334155 !important; }
        .btn-source-link {
            background: #1e293b !important;
            color: #94a3b8 !important;
            border: 1px solid #334155 !important;
        }
        .btn-source-link:hover { color: #f8fafc !important; }
        .modal-back-btn {
            display: none;
            background: #f1f5f9;
            color: #0f172a;
            border: none;
            border-radius: 50%;
            width: 36px;
            height: 36px;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            flex-shrink: 0;
            transition: background 0.15s;
        }
        .modal-back-btn:hover { background: #e2e8f0; }

        /* Toast */
        .toast-notify {
            position: fixed;
            bottom: 24px;
            right: 24px;
            background: #0f172a;
            color: #ffffff;
            padding: 12px 20px;
            border-radius: 6px;
            font-size: 13px;
            font-weight: 600;
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
            display: none;
            z-index: 2000;
        }

        /* ---------------- RESPONSIVE / MOBILE ENHANCEMENTS ---------------- */
        @media (max-width: 1080px) {
            .main-container {
                grid-template-columns: minmax(0, 1fr) 280px;
            }
            .left-col {
                display: none;
            }
        }
        @media (max-width: 768px) {
            /* Compact Header */
            .top-banner {
                padding: 8px 12px;
            }
            .top-banner-inner {
                flex-direction: row;
                justify-content: space-between;
                align-items: center;
                gap: 8px;
            }
            .top-left {
                font-size: 11px;
                font-weight: 700;
                color: #64748b;
            }
            .brand-logo img {
                height: 34px !important;
            }
            .social-icons {
                display: none;
            }

            /* Horizontally Scrollable Nav Bar with Smooth Touch Momentum */
            .primary-nav {
                position: sticky;
                top: 0;
                z-index: 100;
            }
            .nav-inner {
                padding: 0 8px;
                overflow-x: auto;
                -webkit-overflow-scrolling: touch;
                scrollbar-width: none;
            }
            .nav-inner::-webkit-scrollbar {
                display: none;
            }
            .nav-links {
                display: flex;
                flex-wrap: nowrap;
                white-space: nowrap;
                gap: 2px;
            }
            .nav-link {
                padding: 9px 12px;
                font-size: 12px;
                letter-spacing: 0.2px;
                border-radius: 4px;
            }

            /* Single-line Ticker */
            .ticker-ribbon {
                padding: 6px 10px;
            }
            .ticker-inner {
                flex-direction: row;
                align-items: center;
                gap: 8px;
            }
            .ticker-badge {
                font-size: 9px;
                padding: 3px 6px;
                flex-shrink: 0;
            }
            .ticker-stats {
                display: none;
            }

            /* 1-Column Feed */
            .main-container {
                grid-template-columns: 1fr;
                padding: 10px 8px;
                gap: 14px;
            }
            .right-col {
                order: 3;
            }

            /* Hero Card Mobile Optimizations */
            .hero-feature-card {
                border-radius: 10px;
                margin-bottom: 12px;
            }
            .hero-image-wrap {
                height: auto;
                aspect-ratio: 16/9;
                max-height: 220px;
            }
            .hero-overlay {
                padding: 14px;
            }
            .hero-title {
                font-size: 17px !important;
                line-height: 1.4 !important;
                margin-bottom: 6px;
            }
            .hero-desc {
                font-size: 12.5px !important;
                line-height: 1.4 !important;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                overflow: hidden;
            }

            /* Teasers & Grid Cards */
            .hero-teasers-strip {
                grid-template-columns: 1fr;
                gap: 8px;
            }
            .editorial-grid {
                grid-template-columns: 1fr;
                gap: 10px;
            }
            .editorial-card {
                padding: 12px;
                border-radius: 8px;
            }

            /* App-like Full Screen Article Reader on Mobile */
            .modal-backdrop {
                padding: 0 !important;
                align-items: flex-start !important;
            }
            .modal-card {
                max-width: 100% !important;
                width: 100% !important;
                height: 100vh !important;
                max-height: 100vh !important;
                border-radius: 0 !important;
                box-shadow: none !important;
                display: flex !important;
                flex-direction: column !important;
            }
            .modal-header {
                padding: 12px 14px !important;
                position: sticky !important;
                top: 0 !important;
                background: #ffffff !important;
                z-index: 20 !important;
                border-bottom: 1px solid #e2e8f0 !important;
                gap: 10px !important;
            }
            .modal-back-btn {
                display: flex !important;
            }
            .modal-close-btn {
                width: 32px !important;
                height: 32px !important;
                font-size: 20px !important;
            }
            .modal-title {
                font-size: 18px !important;
                line-height: 1.4 !important;
            }
            .modal-body {
                padding: 14px !important;
                flex: 1 !important;
                overflow-y: auto !important;
                -webkit-overflow-scrolling: touch !important;
                gap: 14px !important;
            }
            .modal-media {
                max-height: 230px !important;
                border-radius: 8px !important;
            }
            .modal-media img, .modal-media iframe {
                max-height: 230px !important;
                height: 230px !important;
            }
            .modal-text {
                font-size: 16px !important;
                line-height: 1.75 !important;
                color: #1e293b !important;
            }
            .modal-actions {
                position: sticky !important;
                bottom: 0 !important;
                background: #ffffff !important;
                border-top: 1px solid #e2e8f0 !important;
                padding: 10px 12px !important;
                flex-direction: column !important;
                gap: 8px !important;
                box-shadow: 0 -4px 16px rgba(0,0,0,0.08) !important;
                z-index: 20 !important;
            }
            .modal-btn-row {
                display: grid !important;
                grid-template-columns: 1.4fr 1fr 0.9fr 0.9fr !important;
                gap: 6px !important;
                width: 100% !important;
            }
            .btn-modal-action {
                justify-content: center !important;
                padding: 9px 4px !important;
                font-size: 12px !important;
                border-radius: 6px !important;
            }
        }
    </style>
</head>
<body>

    <!-- TOP BANNER -->
    <header class="top-banner">
        <div class="top-banner-inner">
            <div class="top-left">
                <span id="currentDateDisplay">செப்டம்பர் 8, 2026</span>
            </div>

            <a href="/portal" class="brand-logo" style="text-decoration:none;">
                <img src="/portal/assets/brand/tn24-logo.svg?v=20260908d" alt="TN24" style="height:46px; display:block;" />
            </a>

            <div class="top-right">
                <div class="social-icons">
                    <a href="javascript:void(0)" title="Facebook"><svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><path d="M9 8H6v4h3v12h5V12h3.642L18 8h-4V6.333C14 5.374 14.5 5 15.5 5H18V0h-3.808C10.595 0 9 1.582 9 4.615V8z"/></svg></a>
                    <a href="javascript:void(0)" title="Twitter / X"><svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg></a>
                    <a href="javascript:void(0)" title="RSS Feed"><svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><circle cx="6.18" cy="17.82" r="2.18"/><path d="M4 4.44v2.83c7.03 0 12.73 5.7 12.73 12.73h2.83c0-8.59-6.97-15.56-15.56-15.56zm0 5.66v2.83c3.9 0 7.07 3.17 7.07 7.07h2.83c0-5.47-4.43-9.9-9.9-9.9z"/></svg></a>
                    <a href="javascript:void(0)" title="Telegram"><svg width="12" height="12" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.562 8.161c-.18.718-1.5 6.305-2.174 9.141-.285 1.2-.727 1.4-1.156 1.439-.933.086-1.642-.617-2.546-1.21-1.415-.929-2.215-1.507-3.589-2.413-1.587-1.047-.558-1.623.346-2.564.237-.246 4.343-3.981 4.422-4.321.01-.043.018-.204-.078-.29-.096-.085-.237-.056-.34-.033-.146.033-2.476 1.573-6.99 4.622-.662.455-1.261.678-1.798.666-.592-.013-1.73-.334-2.578-.609-1.04-.338-1.868-.517-1.796-1.091.037-.299.434-.605 1.189-.918 4.654-2.028 7.759-3.364 9.314-4.009 4.434-1.841 5.356-2.161 5.957-2.172.132-.002.427.031.618.187.161.132.206.311.228.436.022.126.049.414.027.64z"/></svg></a>
                </div>
            </div>
        </div>
    </header>

    <!-- PRIMARY NAVIGATION RIBBON -->
    <nav class="primary-nav">
        <div class="nav-inner">
            <ul class="nav-links">
                <li class="nav-item active" id="nav-item-home">
                    <a href="javascript:void(0)" onclick="filterCategory('All')" class="nav-link">முகப்பு</a>
                </li>
                <li class="nav-item" id="nav-item-districts">
                    <a href="javascript:void(0)" class="nav-link">மாவட்டங்கள் <svg viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M6 9l6 6 6-6"/></svg></a>
                    <div class="nav-dropdown" id="districtDropdownMenu">
                        <a href="javascript:void(0)" onclick="filterDistrict('All Districts')">அனைத்து மாவட்டங்கள் (38+)</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Chennai')">சென்னை</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Coimbatore')">கோயம்புத்தூர்</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Madurai')">மதுரை</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Tiruchirappalli')">திருச்சிராப்பள்ளி</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Salem')">சேலம்</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Tirunelveli')">திருநெல்வேலி</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Erode')">ஈரோடு</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Vellore')">வேலூர்</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Thanjavur')">தஞ்சாவூர்</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Kanyakumari')">கன்னியாகுமரி</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Ranipet')">ராணிப்பேட்டை</a>
                        <a href="javascript:void(0)" onclick="filterDistrict('Dindigul')">திண்டுக்கல்</a>
                    </div>
                </li>
                <li class="nav-item" id="nav-item-viral">
                    <a href="javascript:void(0)" onclick="filterViral()" class="nav-link" style="color: #f97316; font-weight: 700;">🔥 வைரல் செய்திகள்</a>
                </li>
                <li class="nav-item" id="nav-item-sports">
                    <a href="javascript:void(0)" onclick="filterCategory('Sports')" class="nav-link" style="color: #38bdf8; font-weight: 700;">⚽ விளையாட்டு</a>
                </li>
                <li class="nav-item" id="nav-item-categories">
                    <a href="javascript:void(0)" class="nav-link">பிரிவுகள் <svg viewBox="0 0 24 24" fill="none" stroke="currentColor"><path d="M6 9l6 6 6-6"/></svg></a>
                    <div class="nav-dropdown">
                        <a href="javascript:void(0)" onclick="filterCategory('News')">செய்திகள் &amp; பொதுமக்கள்</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Politics')">அரசியல் &amp; ஆட்சி முறை</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Sports')">விளையாட்டு &amp; கிரிக்கெட்</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Technical')">தொழில்நுட்பம் &amp; ஏஐ</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Business')">வணிகம் &amp; பங்குச்சந்தை</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Entertainment')">சினிமா &amp; கலை உலகம்</a>
                        <a href="javascript:void(0)" onclick="filterCategory('Crime')">குற்ற நிகழ்வுகள் &amp; சட்டம்</a>
                    </div>
                </li>
                <li class="nav-item" id="nav-item-events">
                    <a href="javascript:void(0)" onclick="scrollToSection('calendarWidget')" class="nav-link">📅 நிகழ்வுகள்</a>
                </li>
                <li class="nav-item" id="nav-item-news">
                    <a href="javascript:void(0)" onclick="filterCategory('News')" class="nav-link">செய்திகள்</a>
                </li>
            </ul>

            <div class="nav-search-box">
                <input type="text" id="searchInput" class="nav-search-input" placeholder="செய்திகளைத் தேடுக..." onkeyup="handleSearchKey(event)">
                <button class="nav-search-btn" onclick="executeSearch()" title="தேடுக">
                    <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/></svg>
                </button>
            </div>
        </div>
    </nav>

    <!-- MARKET & REGION TICKER BAR -->
    <div class="market-ticker-bar">
        <div class="ticker-inner">
            <div class="ticker-selector-wrap">
                <button class="ticker-select-btn" id="tickerSelectBtn" onclick="toggleTickerDropdown()">
                    <span id="selectedDistrictLabel">தமிழ்நாடு - அனைத்து வட்டாரங்கள்</span>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M6 9l6 6 6-6"/></svg>
                </button>
            </div>

            <!-- Feed View Mode Toggle -->
            <div style="display: flex; gap: 4px; align-items: center; margin-right: 12px;">
                <button type="button" id="portal-view-feed-btn" class="action-btn" style="font-size: 11px; padding: 2px 8px; background: #0284c7; color: #fff; border-color: #38bdf8; font-weight: 700; border-radius: 4px;" onclick="setPortalViewMode('feed')">📰 செய்தி ஓட்டம்</button>
                <button type="button" id="portal-view-grouped-btn" class="action-btn" style="font-size: 11px; padding: 2px 8px; color: #94a3b8; border: 1px solid #334155; background: transparent; font-weight: 600; border-radius: 4px;" onclick="setPortalViewMode('grouped')">📑 தொகுப்பு</button>
            </div>

            <!-- Animated Trend Sparkline -->
            <div class="sparkline-box">
                <svg class="sparkline-svg" viewBox="0 0 140 28">
                    <path d="M 0 22 Q 25 24, 40 16 T 70 18 T 95 10 T 120 14 L 140 4" fill="none" stroke="#10b981" stroke-width="2.2" stroke-linecap="round"/>
                    <circle cx="140" cy="4" r="3" fill="#10b981"/>
                </svg>
            </div>

            <div class="ticker-stats">
                <div class="ticker-stat-item">
                    <span class="stat-label">நேரலை 2026</span>
                    <span class="stat-val up" id="tickerLiveCount">143 செய்திகள்</span>
                </div>
                <div class="ticker-stat-item">
                    <span class="stat-label">சரிபார்க்கப்பட்டவை</span>
                    <span class="stat-val">98.4%</span>
                </div>
                <div class="ticker-stat-item">
                    <span class="stat-label">அதிவேக வைரல்</span>
                    <span class="stat-val up" id="tickerViralCount">42 பதிவுகள்</span>
                </div>
                <div class="ticker-stat-item">
                    <span class="stat-label">வாசகர்கள் எண்ணிக்கை</span>
                    <span class="stat-val">1.85M நேரலை</span>
                </div>
            </div>
        </div>
    </div>

    <!-- MAIN 3-COLUMN EDITORIAL CONTAINER -->
    <main class="main-container">

        <!-- COLUMN 1: LEFT SIDEBAR -->
        <aside class="left-col">
            <div id="leftFeedCards" style="display: flex; flex-direction: column; gap: 16px;">
                <!-- Dynamically populated mini cards -->
            </div>

            <div class="section-header-compact">அரசு அறிவிப்புகள் &amp; அறிக்கைகள்</div>
            <div class="press-list" id="pressReleasesList">
                <!-- Dynamically populated press releases -->
            </div>

            <!-- Square Ad Banner 200x200 -->
            <div id="ad-banner-square-slot" style="margin-top: 14px;">
                <div class="ad-banner-square" style="padding:0; border:none; background:transparent;">
                    <img src="/portal/assets/brand/tn24-square.svg" style="width:200px; height:200px; border-radius:8px; display:block;" alt="TN24 Direct" />
                </div>
            </div>
        </aside>

        <!-- COLUMN 2: CENTER MAIN CONTENT -->
        <section class="center-col" id="portal-feed-anchor">

            <!-- Active Feed Filter Bar Banner -->
            <div id="activeFilterBar" style="display: none; background: linear-gradient(90deg, rgba(15,23,42,0.95), rgba(30,41,59,0.9)); border: 1px solid rgba(56,189,248,0.4); border-radius: 10px; padding: 12px 18px; margin-bottom: 18px; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.4);">
                <div style="display: flex; align-items: center; gap: 12px;">
                    <span id="activeFilterIcon" style="font-size: 24px;">🔥</span>
                    <div>
                        <div id="activeFilterTitle" style="font-size: 14px; font-weight: 800; color: #38bdf8; letter-spacing: 0.5px;">செயலில் உள்ள செய்திப் பிரிவு</div>
                        <div id="activeFilterDesc" style="font-size: 11px; color: #94a3b8; margin-top: 2px;">நேரலை செய்தி ஓட்டம் வடிகட்டப்பட்டுள்ளது</div>
                    </div>
                </div>
                <button type="button" onclick="filterCategory('All')" style="background: rgba(56,189,248,0.15); border: 1px solid #38bdf8; color: #38bdf8; font-size: 12px; font-weight: 700; padding: 6px 14px; border-radius: 6px; cursor: pointer;" title="அனைத்து செய்திகளையும் பார்க்க">
                    ✕ அனைத்து செய்திகள்
                </button>
            </div>

            <!-- Lead Hero Article -->
            <article class="hero-lead-card" id="heroCard">
                <div class="hero-image-wrap" onclick="openHeroArticle(event)">
                    <img id="heroImage" src="https://images.unsplash.com/photo-1585829365295-ab7cd400c167?w=1200" alt="தலைப்புச் செய்தி">
                    <div class="hero-overlay">
                        <span class="badge-pill" id="heroBadge">தலைப்புச் செய்தி</span>
                        <h1 class="hero-title" id="heroTitle">தமிழ்நாடு முக்கிய செய்தி நிகழ்வுகள்</h1>
                        <p class="hero-desc" id="heroDesc">தமிழ்நாடு மற்றும் வட்டார முக்கிய நிகழ்வுகள் குறித்த விரிவான கள நிலவரம் மற்றும் நேரடி செய்தி தொகுப்பு.</p>
                        <div id="heroDateMeta" style="margin-top: 8px; font-size: 11px; color: #cbd5e1; display: flex; align-items: center; gap: 8px; font-weight: 500;"></div>
                    </div>
                </div>

                <!-- Hero Sub-Teasers (3 Items) -->
                <div class="hero-teasers-strip" id="heroTeasers">
                    <!-- Populated dynamically -->
                </div>
            </article>

            <!-- Horizontal Leaderboard Banner (728x90) -->
            <div id="ad-banner-header-slot" style="margin: 18px 0;">
                <div class="ad-banner-horizontal" style="padding:0; border:none; background:transparent;">
                    <img src="/portal/assets/brand/tn24-header.svg?v=20260908c" style="width:100%; height:90px; border-radius:8px; display:block;" alt="TN24 Live News" />
                </div>
            </div>

            <!-- 2-Column Editorial Card Grid -->
            <div class="editorial-grid" id="editorialGrid">
                <!-- Populated dynamically -->
            </div>

        </section>

        <!-- COLUMN 3: RIGHT SIDEBAR -->
        <aside class="right-col">

            <!-- Live Market Rates & Currency Calculator Widget -->
            <div class="calculator-widget">
                <div class="widget-title-navy">
                    <span>💰 தங்கம், வெள்ளி &amp; நாணய கால்குலேட்டர்</span>
                    <span class="market-pulse-badge"><span class="market-pulse-dot"></span> நேரலை TN</span>
                </div>

                <!-- Live Quick Glance Rate Cards (Gold, Silver, Currency against INR) -->
                <div class="market-ticker-grid">
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>🥇 தங்கம் 22K (1 சவரன்)</span><span style="color:#f59e0b;">+0.3%</span></div>
                        <div class="market-ticker-val" id="ticker-gold-22k">₹ 54,000</div>
                    </div>
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>👑 தங்கம் 24K (10 கிராம்)</span><span style="color:#f59e0b;">+0.3%</span></div>
                        <div class="market-ticker-val" id="ticker-gold-24k">₹ 73,650</div>
                    </div>
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>🥈 வெள்ளி (1 கிலோ)</span><span style="color:#38bdf8;">+0.5%</span></div>
                        <div class="market-ticker-val" id="ticker-silver">₹ 94,500</div>
                    </div>
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>💵 அமெரிக்க டாலர் (USD)</span><span style="color:#4ade80;">நேரலை</span></div>
                        <div class="market-ticker-val" id="ticker-usd">₹ 87.40</div>
                    </div>
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>🇦🇪 யுஏஇ திர்ஹாம் (AED)</span><span style="color:#4ade80;">நேரலை</span></div>
                        <div class="market-ticker-val" id="ticker-aed">₹ 23.80</div>
                    </div>
                    <div class="market-ticker-card">
                        <div class="market-ticker-label"><span>💶 யூரோ (EUR)</span><span style="color:#4ade80;">நேரலை</span></div>
                        <div class="market-ticker-val" id="ticker-eur">₹ 95.10</div>
                    </div>
                </div>

                <!-- Interactive Calculators -->
                <div class="calc-rows">
                    <!-- Gold Calculator -->
                    <div class="calc-row">
                        <div class="calc-left">
                            <input type="number" step="any" min="0.1" class="calc-input" value="1" id="calcGoldQty" oninput="runCalculator()">
                            <select class="calc-unit-select" id="calcGoldUnit" onchange="runCalculator()">
                                <option value="sovereign">1 சவரன் (8 கிராம்)</option>
                                <option value="gram">1 கிராம்</option>
                                <option value="10g">10 கிராம்</option>
                                <option value="100g">100 கிராம்</option>
                            </select>
                            <select class="calc-unit-select" id="calcGoldPurity" onchange="runCalculator()">
                                <option value="22k">22K ஆபரணத் தங்கம்</option>
                                <option value="24k">24K சுத்தத் தங்கம்</option>
                            </select>
                            <span class="calc-eq">=</span>
                        </div>
                        <span class="calc-val" id="calcGoldOut">₹ 54,000</span>
                    </div>

                    <!-- Silver Calculator -->
                    <div class="calc-row">
                        <div class="calc-left">
                            <input type="number" step="any" min="1" class="calc-input" value="10" id="calcSilverQty" oninput="runCalculator()">
                            <select class="calc-unit-select" id="calcSilverUnit" onchange="runCalculator()">
                                <option value="gram">கிராம்</option>
                                <option value="kg">கிலோ</option>
                            </select>
                            <span class="calc-name" style="color:#cbd5e1; font-size:10px;">சுத்த வெள்ளி (999)</span>
                            <span class="calc-eq">=</span>
                        </div>
                        <span class="calc-val" id="calcSilverOut">₹ 945.00</span>
                    </div>

                    <!-- Currency against INR Calculator -->
                    <div class="calc-row">
                        <div class="calc-left">
                            <input type="number" step="any" min="1" class="calc-input" value="100" id="calcCurrQty" oninput="runCalculator()">
                            <select class="calc-unit-select" id="calcCurrType" onchange="runCalculator()">
                                <option value="USD">USD ($ - அமெரிக்க டாலர்)</option>
                                <option value="AED">AED (د.إ - யுஏஇ திர்ஹாம்)</option>
                                <option value="EUR">EUR (€ - யூரோ)</option>
                                <option value="GBP">GBP (£ - பிரிட்டிஷ் பவுண்ட்)</option>
                                <option value="SGD">SGD (S$ - சிங்கப்பூர் டாலர்)</option>
                                <option value="SAR">SAR (﷼ - சவுதி ரியால்)</option>
                                <option value="KWD">KWD (KD - குவைத் தினார்)</option>
                                <option value="CAD">CAD (C$ - கனடிய டாலர்)</option>
                                <option value="AUD">AUD (A$ - ஆஸ்திரேலிய டாலர்)</option>
                                <option value="MYR">MYR (RM - மலேசிய ரிங்கிட்)</option>
                                <option value="QAR">QAR (QR - கத்தார் ரியால்)</option>
                            </select>
                            <span class="calc-eq">=</span>
                        </div>
                        <span class="calc-val" id="calcCurrOut">₹ 8,740.00</span>
                    </div>
                </div>
            </div>

            <!-- Tamil Nadu Real-time Events & Governance Summits Calendar Widget -->
            <div class="widget-box" id="calendarWidget" style="background: #0d1a2d; border: 1px solid #1e293b; color: #fff;">
                <div class="widget-title-white" style="color: #38bdf8; letter-spacing: 0.5px; border-bottom: 1px solid #1e293b; display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; padding-bottom: 6px;">
                    <span style="display: flex; align-items: center; gap: 6px;">📅 தமிழ்நாடு மாநாடுகள் &amp; நிகழ்வுகள்</span>
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <button type="button" onclick="loadPortalFeed(); showToast('தமிழ்நாடு மாநாடுகள் புதுப்பிக்கப்படுகின்றன...');" style="background: transparent; border: none; color: #38bdf8; cursor: pointer; font-size: 13px; padding: 0;" title="நிகழ்வுகளைப் புதுப்பிக்க">🔄</button>
                        <span style="font-size: 9px; font-weight: 700; color: #4ade80; background: rgba(74, 222, 128, 0.12); padding: 2px 7px; border-radius: 4px; border: 1px solid rgba(74,222,128,0.3); display: inline-flex; align-items: center; gap: 4px;">
                            <span class="market-pulse-dot"></span> நேரலை 2026
                        </span>
                    </div>
                </div>
                <div class="event-list" id="calendarEventsList">
                    <!-- Populated dynamically via real-time feed -->
                </div>
                <a href="javascript:void(0)" onclick="openAllEventsModal()" class="btn-view-all" style="background: rgba(56,189,248,0.1); color: #38bdf8; border: 1px solid rgba(56,189,248,0.25); text-align: center; border-radius: 6px; margin-top: 10px; display: block;">முழு மாநில மாநாட்டு காலண்டர் ↗</a>
            </div>

            <!-- 250x250 Ad Banner -->
            <div id="ad-banner-sidebar-slot" style="margin-bottom: 20px;">
                <div class="ad-banner-250" style="padding:0; border:none; background:transparent;">
                    <img src="/portal/assets/brand/tn24-sidebar.svg?v=20260908c" onerror="this.onerror=null; this.src='/admin/api/maps/svg?district=Tamil%20Nadu'" style="width:100%; max-width:250px; height:250px; border-radius:8px; display:block;" alt="TN24 Prime" />
                </div>
            </div>

            <!-- Most Read (Numbered 1 to 5) -->
            <div class="widget-box">
                <div class="widget-title-white">அதிகம் வாசிக்கப்பட்டவை</div>
                <div class="most-read-list" id="mostReadList">
                    <!-- Populated dynamically -->
                </div>
            </div>

            <!-- Latest Tweets -->
            <div class="widget-box">
                <div class="widget-title-white">சமூக ஊடக நேரலை</div>
                <div class="tweet-card">
                    <p class="tweet-text">தமிழ்நாட்டின் அனைத்து 38 மாவட்டங்களின் நேரடி செய்திகள் மற்றும் கள நிலவரங்கள் TN24 தளத்தில் உடனுக்குடன் நேரலையாக வழங்கப்படுகிறது.</p>
                    <div class="tweet-meta">
                        <span>நேரலை பதிவு</span>
                        <div class="tweet-actions">
                            <a href="javascript:void(0)" onclick="showToast('பகிரப்பட்டது!')">பகிர்</a>
                            &bull;
                            <a href="javascript:void(0)" onclick="showToast('விருப்பத்தில் சேர்க்கப்பட்டது!')">விருப்பம்</a>
                        </div>
                    </div>
                </div>
            </div>

        </aside>

    </main>

    <!-- FOOTER BRANDING -->
    <footer style="margin-top: 40px; background: #090e17; border-top: 1px solid #1e293b; color: #94a3b8; padding: 36px 20px 24px; text-align: center;">
        <div style="max-width: 1320px; margin: 0 auto; display: flex; flex-direction: column; align-items: center; gap: 14px;">
            <img src="/portal/assets/brand/tn24-logo.svg?v=20260908d" onerror="this.onerror=null; this.src='/admin/api/maps/svg?district=Tamil%20Nadu'" alt="TN24" style="height: 50px; display: block;" />
            <p style="font-size: 13px; color: #cbd5e1; max-width: 650px; margin: 0 auto; line-height: 1.6;">
                <strong style="color: #ffffff;">TN24 &mdash; தமிழ்நாட்டின் முதன்மை 24/7 டிஜிட்டல் செய்தி &amp; நேரடி தகவல் தளம்</strong><br>
                தமிழ்நாட்டின் அனைத்து 38 மாவட்டங்கள், இந்தியா மற்றும் உலகளாவிய முக்கிய நிகழ்வுகள் உடனுக்குடன் நேரலையாக.
            </p>
            <div style="display: flex; gap: 16px; flex-wrap: wrap; justify-content: center; font-size: 12px; margin-top: 2px;">
                <a href="mailto:tn24now@gmail.com" style="color: #38bdf8; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; background: rgba(56,189,248,0.1); padding: 5px 12px; border-radius: 6px; border: 1px solid rgba(56,189,248,0.25);">
                    📧 <strong>tn24now@gmail.com</strong>
                </a>
                <a href="tel:+918124395082" style="color: #4ade80; text-decoration: none; display: inline-flex; align-items: center; gap: 5px; background: rgba(74,222,128,0.1); padding: 5px 12px; border-radius: 6px; border: 1px solid rgba(74,222,128,0.25);">
                    📞 <strong>+91 81243 95082</strong>
                </a>
            </div>
            <div style="font-size: 11px; color: #64748b; margin-top: 6px;">
                &copy; 2026 TN24 டிஜிட்டல் செய்தி ஊடக வலையமைப்பு. அனைத்து உரிமைகளும் பாதுகாக்கப்பட்டவை.
            </div>
        </div>
    </footer>

    <!-- ARTICLE READER MODAL -->
    <div class="modal-backdrop" id="articleModal" onclick="handleBackdropClick(event)">
        <div class="modal-card">
            <div class="modal-header">
                <button class="modal-back-btn" onclick="closeArticleModal()" title="பின்செல்க">
                    <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
                </button>
                <div style="flex: 1; min-width: 0;">
                    <div class="modal-meta-bar">
                        <span class="badge-pill" id="modalCategoryBadge" style="position: static;">செய்திகள்</span>
                        <span id="modalDistrict">தமிழ்நாடு</span>
                        &bull;
                        <span id="modalDate">செப் 8, 2026</span>
                    </div>
                    <h2 class="modal-title" id="modalTitle">செய்தித் தலைப்பு</h2>
                </div>
                <button class="modal-close-btn" onclick="closeArticleModal()">&times;</button>
            </div>
            <div class="modal-body">
                <div class="modal-media" id="modalMediaContainer">
                    <img id="modalImage" src="" onerror="this.onerror=null; this.src='/admin/api/maps/svg?district=Tamil%20Nadu'" alt="செய்திப் படம்">
                </div>
                <div class="modal-text" id="modalDescription">முழு செய்தி விவரம் இங்கே காண்பிக்கப்படும்...</div>
            </div>
            <div class="modal-actions">
                <span style="font-size: 11px; font-weight: 700; color: #64748b;" id="modalAuthor">TN24 செய்திக் குழு</span>
                <div class="modal-btn-row">
                    <button class="btn-modal-action btn-share-whatsapp" onclick="shareToWhatsApp()" title="WhatsApp-ல் பகிர்க">
                        <svg width="15" height="15" fill="currentColor" viewBox="0 0 24 24"><path d="M12.031 0C5.397 0 .017 5.38.017 12.014c0 2.117.553 4.185 1.603 6.007L0 24l6.166-1.618a11.96 11.96 0 005.865 1.53h.005c6.634 0 12.014-5.38 12.014-12.014 0-3.208-1.25-6.224-3.52-8.495A11.939 11.939 0 0012.031 0zm-.005 21.966h-.004c-1.802 0-3.568-.485-5.107-1.398l-.366-.217-3.795.996 1.013-3.7-.238-.379a9.986 9.986 0 01-1.534-5.254C2.001 6.475 6.495 1.98 12.026 1.98c2.673 0 5.187 1.042 7.078 2.934a9.948 9.948 0 012.928 7.076c0 5.532-4.494 10.026-10.006 10.026z"/></svg>
                        WhatsApp
                    </button>
                    <button class="btn-modal-action btn-share-native" onclick="shareNative()" title="பிற செயலிகளில் பகிர்க">
                        <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/></svg>
                        பகிர்க
                    </button>
                    <button class="btn-modal-action btn-copy-link" onclick="copyShareLink()" title="இணைப்பை நகலெடு">
                        🔗 நகல்
                    </button>
                    <button class="btn-modal-action btn-source-link" id="modalSourceBtn" onclick="visitOriginalSource()" title="அசல் செய்தி இணைப்பு">
                        மூலம் ↗
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- EVENT DETAILS MODAL -->
    <div class="modal-backdrop" id="eventModal" onclick="handleEventBackdropClick(event)">
        <div class="modal-card" style="max-width: 580px;">
            <div class="modal-header">
                <div>
                    <div class="modal-meta-bar">
                        <span class="badge-pill" id="eventModalCategoryBadge" style="position: static; background: #0284c7;">மாநாடு</span>
                        <span id="eventModalDistrict" style="color: #38bdf8; font-weight: 700;">சென்னை</span>
                        &bull;
                        <span id="eventModalCountdown" style="color: #4ade80; font-weight: 700;">இன்னும் 4 நாட்களில்</span>
                    </div>
                    <h2 class="modal-title" id="eventModalTitle" style="font-size: 19px; margin-top: 6px;">நிகழ்வுத் தலைப்பு</h2>
                </div>
                <button class="modal-close-btn" onclick="closeEventModal()">&times;</button>
            </div>
            <div class="modal-body" style="padding: 20px;">
                <div style="background: rgba(15,23,42,0.8); border: 1px solid #334155; border-radius: 8px; padding: 14px; margin-bottom: 16px;">
                    <div style="display: flex; gap: 10px; align-items: center; margin-bottom: 8px;">
                        <span style="font-size: 16px;">🗓️</span>
                        <div>
                            <div style="font-size: 10px; color: #94a3b8; font-weight: 700; text-transform: uppercase;">அதிகாரப்பூர்வ தேதி &amp; நேரம்</div>
                            <div id="eventModalDate" style="font-size: 13px; font-weight: 800; color: #f8fafc;">செப் 12 - 14, 2026</div>
                        </div>
                    </div>
                    <div style="display: flex; gap: 10px; align-items: center;">
                        <span style="font-size: 16px;">📍</span>
                        <div>
                            <div style="font-size: 10px; color: #94a3b8; font-weight: 700; text-transform: uppercase;">இடம் &amp; முகவரி</div>
                            <div id="eventModalVenue" style="font-size: 12px; font-weight: 600; color: #cbd5e1;">சென்னை வர்த்தக மையம் • நந்தம்பாக்கம், சென்னை</div>
                        </div>
                    </div>
                </div>

                <div style="font-size: 13px; color: #cbd5e1; line-height: 1.65; margin-bottom: 20px;" id="eventModalDesc">
                    முழு நிகழ்ச்சி விவரங்கள்.
                </div>

                <div style="display: flex; gap: 10px;">
                    <a id="eventGoogleCalBtn" href="#" target="_blank" class="action-btn" style="flex: 1; text-align: center; background: #0284c7; color: #fff; border-color: #38bdf8; font-weight: 700; padding: 10px 14px; border-radius: 6px; text-decoration: none;">
                        📅 கூகிள் காலண்டரில் சேர்க்க
                    </a>
                    <button type="button" onclick="closeEventModal()" class="action-btn" style="padding: 10px 18px; border-radius: 6px;">மூடுக</button>
                </div>
            </div>
        </div>
    </div>

    <!-- FULL STATE CALENDAR & SUMMITS MODAL -->
    <div class="modal-backdrop" id="fullStateCalendarModal" onclick="handleFullCalendarBackdropClick(event)" style="display: none;">
        <div class="modal-card" style="max-width: 880px; max-height: 90vh; display: flex; flex-direction: column;">
            <div class="modal-header" style="background: linear-gradient(135deg, #091326 0%, #0f172a 100%); border-bottom: 1px solid #1e293b; padding: 18px 24px;">
                <div style="display: flex; align-items: center; gap: 12px;">
                    <div style="width: 42px; height: 42px; border-radius: 10px; background: rgba(56,189,248,0.15); border: 1px solid #38bdf8; display: flex; align-items: center; justify-content: center; font-size: 20px;">📅</div>
                    <div>
                        <div style="display: flex; align-items: center; gap: 8px;">
                            <h2 class="modal-title" style="font-size: 18px; margin: 0; color: #f8fafc;">தமிழ்நாடு மாநில மாநாடுகள் &amp; நிகழ்வுகள் 2026</h2>
                            <span style="font-size: 9px; font-weight: 800; color: #4ade80; background: rgba(74, 222, 128, 0.15); padding: 2px 8px; border-radius: 4px; border: 1px solid rgba(74,222,128,0.3);">நேரலை பட்டியல்</span>
                        </div>
                        <p style="font-size: 11px; color: #94a3b8; margin: 2px 0 0 0;">அனைத்து 38 மாவட்டங்களுக்கான அதிகாரப்பூர்வ மாநாட்டு கால அட்டவணை • தொழில், விவசாயம், மின்வாகனம், கலாச்சாரம் &amp; இளைஞர் நலம்</p>
                    </div>
                </div>
                <button class="modal-close-btn" onclick="closeFullCalendarModal()">&times;</button>
            </div>

            <!-- Filter and Search Bar -->
            <div style="padding: 12px 24px; background: #0b1322; border-bottom: 1px solid #1e293b; display: flex; flex-wrap: wrap; gap: 10px; align-items: center; justify-content: space-between;">
                <div style="display: flex; gap: 8px; align-items: center; flex: 1; min-width: 240px;">
                    <input type="text" id="fullCalendarSearch" placeholder="மாநாடுகள், மாவட்டங்கள் அல்லது தலைப்புகளைத் தேடுக..." oninput="filterFullCalendar()" style="width: 100%; background: #0f172a; border: 1px solid #334155; border-radius: 6px; padding: 8px 12px; color: #fff; font-size: 12px; outline: none;" />
                </div>
                <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap;">
                    <select id="fullCalendarCategoryFilter" onchange="filterFullCalendar()" style="background: #0f172a; border: 1px solid #334155; border-radius: 6px; padding: 8px 12px; color: #38bdf8; font-size: 12px; font-weight: 600; outline: none;">
                        <option value="ALL">அனைத்துப் பிரிவுகளும் (8)</option>
                        <option value="Governance & Industry">ஆட்சி முறை &amp; தொழில்துறை</option>
                        <option value="Agriculture & Innovation">வேளாண்மை &amp; புத்தாக்கம்</option>
                        <option value="Tech & Manufacturing">தொழில்நுட்பம் &amp; உற்பத்தி</option>
                        <option value="Arts & Literature">கலை &amp; இலக்கியம்</option>
                        <option value="Sports & Athletics">விளையாட்டு &amp; தடகளம்</option>
                        <option value="Commerce & Textile">வணிகம் &amp; ஜவுளி</option>
                        <option value="Aerospace & Defense">விண்வெளி &amp; பாதுகாப்பு</option>
                        <option value="Tourism & Heritage">சுற்றுலா &amp; பாரம்பரியம்</option>
                    </select>
                    <button type="button" onclick="loadPortalFeed(); showToast('தமிழ்நாடு மாநாடுகள் புதுப்பிக்கப்பட்டன');" class="action-btn" style="background: rgba(56,189,248,0.12); color: #38bdf8; border-color: rgba(56,189,248,0.3); font-size: 11px; padding: 7px 12px; border-radius: 6px;">🔄 புதுப்பிக்க</button>
                </div>
            </div>

            <!-- Events List Container -->
            <div class="modal-body" id="fullCalendarList" style="padding: 20px 24px; overflow-y: auto; flex: 1; display: grid; grid-template-columns: repeat(auto-fill, minmax(360px, 1fr)); gap: 14px; background: #070d19; max-height: calc(85vh - 160px);">
                <!-- Dynamically populated -->
            </div>

            <div style="padding: 12px 24px; background: #0b1322; border-top: 1px solid #1e293b; display: flex; justify-content: space-between; align-items: center;">
                <span id="fullCalendarCount" style="font-size: 11px; color: #94a3b8; font-weight: 600;">8 அதிகாரப்பூர்வ மாநாடுகள்</span>
                <button type="button" onclick="closeFullCalendarModal()" class="action-btn" style="padding: 8px 18px; border-radius: 6px; font-size: 12px;">மூடுக</button>
            </div>
        </div>
    </div>

    <!-- TOAST NOTIFICATION -->
    <div class="toast-notify" id="toastNotify">செய்தி அறிக்கை</div>

    <!-- JAVASCRIPT LOGIC -->
    <script>
        let portalData = null;
        let currentArticles = [];
        let activeHeroArticle = null;
        let currentModalArticle = null;
        let currentLanguage = 'all';
        let currentDistrict = '';
        let currentCategory = '';
        let currentQuery = '';
        let currentIsViral = false;
        let currentViewMode = 'feed'; // 'feed' | 'grouped'

        // On Load
        document.addEventListener('DOMContentLoaded', () => {
            updateDateString();
            checkUrlPost();
            loadPortalFeed();
            initMarketRates();
            runCalculator();
        });

        function updateDateString() {
            const monthsTa = ['ஜனவரி', 'பிப்ரவரி', 'மார்ச்', 'ஏப்ரல்', 'மே', 'ஜூன்', 'ஜூலை', 'ஆகஸ்ட்', 'செப்டம்பர்', 'அக்டோபர்', 'நவம்பர்', 'டிசம்பர்'];
            const daysTa = ['ஞாயிற்றுக்கிழமை', 'திங்கட்கிழமை', 'செவ்வாய்க்கிழமை', 'புதன்கிழமை', 'வியாழக்கிழமை', 'வெள்ளிக்கிழமை', 'சனிக்கிழமை'];
            const now = new Date();
            const dateEl = document.getElementById('currentDateDisplay');
            if (dateEl) {
                dateEl.textContent = daysTa[now.getDay()] + ', ' + monthsTa[now.getMonth()] + ' ' + now.getDate() + ', ' + now.getFullYear();
            }
        }

        function setPortalLanguage(lang) {
            currentLanguage = lang;
            ['all', 'ta', 'en'].forEach(l => {
                const btn = document.getElementById('portal-lang-' + l);
                if (!btn) return;
                if (l === lang) {
                    btn.style.background = '#0284c7';
                    btn.style.color = '#ffffff';
                    btn.style.borderColor = '#38bdf8';
                    btn.style.fontWeight = '800';
                } else {
                    btn.style.background = 'transparent';
                    btn.style.color = '#94a3b8';
                    btn.style.borderColor = 'transparent';
                    btn.style.fontWeight = '600';
                }
            });
            loadPortalFeed();
            showToast('மொழி வடிகட்டி: ' + (lang === 'all' ? 'அனைத்து மொழிகள்' : (lang === 'ta' ? 'தமிழ் (Tamil)' : 'ஆங்கிலம் (English)')));
        }

        function setPortalViewMode(mode) {
            currentViewMode = mode;
            const feedBtn = document.getElementById('portal-view-feed-btn');
            const groupedBtn = document.getElementById('portal-view-grouped-btn');
            if (mode === 'grouped') {
                if (groupedBtn) {
                    groupedBtn.style.background = '#0284c7';
                    groupedBtn.style.color = '#ffffff';
                    groupedBtn.style.borderColor = '#38bdf8';
                    groupedBtn.style.fontWeight = '800';
                }
                if (feedBtn) {
                    feedBtn.style.background = 'transparent';
                    feedBtn.style.color = '#94a3b8';
                    feedBtn.style.border = '1px solid #334155';
                    feedBtn.style.fontWeight = '600';
                }
            } else {
                if (feedBtn) {
                    feedBtn.style.background = '#0284c7';
                    feedBtn.style.color = '#ffffff';
                    feedBtn.style.borderColor = '#38bdf8';
                    feedBtn.style.fontWeight = '800';
                }
                if (groupedBtn) {
                    groupedBtn.style.background = 'transparent';
                    groupedBtn.style.color = '#94a3b8';
                    groupedBtn.style.border = '1px solid #334155';
                    groupedBtn.style.fontWeight = '600';
                }
            }
            if (portalData) {
                renderPortalUI(portalData);
            }
        }

        async function loadPortalFeed(district, category, query, isViral) {
            try {
                if (district !== undefined) currentDistrict = district;
                if (category !== undefined) currentCategory = category;
                if (query !== undefined) currentQuery = query;
                if (isViral !== undefined) currentIsViral = isViral;

                let url = '/api/portal/feed';
                const params = new URLSearchParams();
                if (currentDistrict && currentDistrict !== 'All Districts' && currentDistrict !== 'TN-ALL / REGIONAL') {
                    params.append('district', currentDistrict);
                }
                if (currentCategory && currentCategory !== 'All') {
                    params.append('category', currentCategory);
                }
                if (currentQuery) params.append('q', currentQuery);
                if (currentIsViral) params.append('viral', 'true');
                if (currentLanguage && currentLanguage !== 'all') {
                    params.append('lang', currentLanguage);
                }
                if (params.toString()) url += '?' + params.toString();

                const res = await fetch(url);
                const json = await res.json();
                if (json.success && json.data) {
                    portalData = json.data;
                    renderPortalUI(portalData);
                    checkUrlPost();
                }
            } catch (err) {
                console.error('Failed to load portal feed:', err);
                showToast('செய்திகளைப் பெற முடியவில்லை');
            }
        }

        function renderBannerHtml(slotName, bannerData) {
            if (bannerData && bannerData.type === 'article' && bannerData.article) {
                const art = bannerData.article;
                const fallback = '/admin/api/maps/svg?district=' + encodeURIComponent(art.district || 'Tamil Nadu');
                const thumb = art.thumbnail || fallback;
                if (slotName === 'header' || slotName === 'infeed') {
                    return '<div class="ad-banner-horizontal promoted-story" onclick="openArticleByID(\'' + art.id + '\', event)" style="cursor:pointer; display:flex; align-items:center; gap:16px; background:linear-gradient(90deg, #09111e, #1e293b); border:1px solid #38bdf8; border-radius:8px; padding:10px 16px; text-align:left; box-shadow:0 4px 16px rgba(56,189,248,0.15); width:100%; box-sizing:border-box;">' +
                        '<span class="ad-tag" style="background:#38bdf8; color:#000; font-weight:800;">சிறப்புச் செய்தி</span>' +
                        '<img src="' + thumb + '" onerror="this.onerror=null; this.src=\'' + fallback + '\'" style="width:110px; height:68px; object-fit:cover; border-radius:6px; flex-shrink:0;" alt="" />' +
                        '<div style="flex:1; overflow:hidden;">' +
                            '<div style="font-size:11px; font-weight:700; color:#38bdf8; text-transform:uppercase;">' + escapeHtml(art.category || 'செய்திகள்') + ' &bull; ' + escapeHtml(art.district || 'தமிழ்நாடு') + '</div>' +
                            '<div style="font-size:14px; font-weight:800; color:#ffffff; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">' + escapeHtml(art.title) + '</div>' +
                            '<div style="font-size:12px; color:#94a3b8; line-height:1.3; max-height:30px; overflow:hidden;">' + escapeHtml(art.description || '') + '</div>' +
                        '</div>' +
                        '<button type="button" style="padding:6px 14px; background:#0284c7; color:#fff; font-weight:700; font-size:12px; border-radius:4px; border:none; cursor:pointer; flex-shrink:0;">முழு விவரம் ▶</button>' +
                    '</div>';
                } else {
                    return '<div class="promoted-story-card" onclick="openArticleByID(\'' + art.id + '\', event)" style="cursor:pointer; background:linear-gradient(180deg, #09111e, #1e293b); border:1px solid #38bdf8; border-radius:8px; padding:12px; text-align:left; margin-bottom:16px; box-shadow:0 4px 16px rgba(56,189,248,0.15); width:100%; box-sizing:border-box;">' +
                        '<span class="ad-tag" style="background:#38bdf8; color:#000; font-weight:800; display:inline-block; margin-bottom:8px;">சிறப்புச் செய்தி</span>' +
                        '<img src="' + thumb + '" onerror="this.onerror=null; this.src=\'' + fallback + '\'" style="width:100%; height:120px; object-fit:cover; border-radius:6px; margin-bottom:8px;" alt="" />' +
                        '<div style="font-size:10px; font-weight:700; color:#38bdf8; text-transform:uppercase;">' + escapeHtml(art.district || 'தமிழ்நாடு') + '</div>' +
                        '<div style="font-size:13px; font-weight:800; color:#ffffff; margin-top:2px;">' + escapeHtml(art.title) + '</div>' +
                    '</div>';
                }
            }
            if (slotName === 'header') {
                return '<div class="ad-banner-horizontal" style="padding:0; border:none; background:transparent;"><img src="/portal/assets/brand/tn24-header.svg?v=20260908c" style="width:100%; height:90px; border-radius:8px; display:block;" alt="TN24 Live News" /></div>';
            }
            if (slotName === 'square') {
                return '<div class="ad-banner-square" style="padding:0; border:none; background:transparent;"><img src="/portal/assets/brand/tn24-square.svg?v=20260908c" style="width:200px; height:200px; border-radius:8px; display:block;" alt="TN24 Direct" /></div>';
            }
            if (slotName === 'sidebar') {
                return '<div class="ad-banner-250" style="padding:0; border:none; background:transparent;"><img src="/portal/assets/brand/tn24-sidebar.svg?v=20260908c" onerror="this.onerror=null; this.src=\'/admin/api/maps/svg?district=Tamil%20Nadu\'" style="width:100%; max-width:250px; height:250px; border-radius:8px; display:block;" alt="TN24 Prime" /></div>';
            }
            if (slotName === 'infeed') {
                return '<div class="ad-banner-mid" style="padding:0; border:none; background:transparent; max-width:100%;"><img src="/portal/assets/brand/tn24-infeed.svg?v=20260908c" style="width:100%; height:90px; border-radius:8px; display:block;" alt="TN24 Speed Desk" /></div>';
            }
            return '';
        }

        function renderPortalUI(data) {
            // Render Dynamic Banners
            const b = data.banners || {};
            const headerSlot = document.getElementById('ad-banner-header-slot');
            if (headerSlot) headerSlot.innerHTML = renderBannerHtml('header', b.header);

            const squareSlot = document.getElementById('ad-banner-square-slot');
            if (squareSlot) squareSlot.innerHTML = renderBannerHtml('square', b.square);

            const sidebarSlot = document.getElementById('ad-banner-sidebar-slot');
            if (sidebarSlot) sidebarSlot.innerHTML = renderBannerHtml('sidebar', b.sidebar);

            // Update Ticker Stats
            if (document.getElementById('tickerLiveCount')) {
                document.getElementById('tickerLiveCount').textContent = (data.totalCount || 143) + ' செய்திகள்';
            }
            if (document.getElementById('tickerViralCount')) {
                document.getElementById('tickerViralCount').textContent = (data.viralCount || 42) + ' பதிவுகள்';
            }

            // Update Active Feed Filter Banner
            const afb = document.getElementById('activeFilterBar');
            const afi = document.getElementById('activeFilterIcon');
            const aft = document.getElementById('activeFilterTitle');
            const afd = document.getElementById('activeFilterDesc');
            if (afb) {
                if (currentIsViral) {
                    afb.style.display = 'flex';
                    if (afi) afi.textContent = '🔥';
                    if (aft) aft.textContent = 'வைரல் ரேடார் & டிரெண்டிங் செய்திகள்';
                    if (afd) afd.textContent = 'தமிழ்நாட்டின் அதிக கவனம்பெற்ற அதிவேக வைரல் செய்திகள் (' + (data.centerArticles ? data.centerArticles.length : 0) + ' செய்திகள் நேரலை)';
                } else if (currentCategory && currentCategory !== 'All') {
                    afb.style.display = 'flex';
                    const catMeta = {
                        'news': { icon: '📰', title: 'செய்திகள் & பொதுமக்கள் களம்', desc: 'மாநில அளவிலான மக்கள் நலம், ஆட்சி முறை மற்றும் சமுதாய முக்கிய தகவல்கள்' },
                        'politics': { icon: '🏛️', title: 'அரசியல் & தேர்தல் களம்', desc: 'தமிழக அரசியல் களம், சட்டப்பேரவை நிகழ்வுகள் மற்றும் தலைவர்களின் அறிக்கைகள்' },
                        'sports': { icon: '⚽', title: 'விளையாட்டு & கிரிக்கெட் அரங்கம்', desc: 'கிரிக்கெட், கால்பந்து, துலீப் டிராபி மற்றும் முக்கிய விளையாட்டு தொடர்கள்' },
                        'technical': { icon: '💻', title: 'தொழில்நுட்பம் & ஏஐ களம்', desc: 'செயற்கை நுண்ணறிவு, மென்பொருள், விண்வெளி மற்றும் டிஜிட்டல் தொழில்நுட்பம்' },
                        'business': { icon: '💼', title: 'வணிகம் & பங்குச்சந்தை மையம்', desc: 'பொருளாதாரம், தொழில்துறை முதலீடுகள், சிப்காட் திட்டங்கள் மற்றும் வர்த்தகம்' },
                        'entertainment': { icon: '🎬', title: 'சினிமா & பொழுதுபோக்கு களம்', desc: 'திரைப்பட விமர்சனங்கள், புது ரிலீஸ், பிரபலங்கள் மற்றும் கலை உலக செய்திகள்' },
                        'crime': { icon: '🚨', title: 'குற்ற நிகழ்வுகள் & சட்டம் ஒழுங்கு', desc: 'காவல்துறை விசாரணை, நீதிமன்ற தீர்ப்புகள் மற்றும் மக்கள் பாதுகாப்பு செய்திகள்' }
                    };
                    const meta = catMeta[currentCategory.toLowerCase()] || { icon: '📂', title: currentCategory + ' பிரிவு', desc: 'தேர்ந்தெடுக்கப்பட்ட பிரிவுக்கான செய்திகள்' };
                    if (afi) afi.textContent = meta.icon;
                    if (aft) aft.textContent = meta.title;
                    if (afd) afd.textContent = meta.desc + ' (' + (data.centerArticles ? data.centerArticles.length : 0) + ' செய்திகள் நேரலை)';
                } else if (currentDistrict && currentDistrict !== 'All Districts' && currentDistrict !== 'TN-ALL / REGIONAL') {
                    afb.style.display = 'flex';
                    if (afi) afi.textContent = '📍';
                    if (aft) aft.textContent = currentDistrict + ' மாவட்டச் செய்திகள்';
                    if (afd) afd.textContent = currentDistrict + ' மாவட்ட நேரடி நிகழ்வுகள் மற்றும் கள நிலவரம் (' + (data.centerArticles ? data.centerArticles.length : 0) + ' செய்திகள் நேரலை)';
                } else if (currentQuery) {
                    afb.style.display = 'flex';
                    if (afi) afi.textContent = '🔍';
                    if (aft) aft.textContent = 'தேடல் முடிவுகள்: "' + currentQuery + '"';
                    if (afd) afd.textContent = 'தேடலுக்குப் பொருத்தமான செய்திகள்';
                } else {
                    afb.style.display = 'none';
                }
            }

            // Hero Story
            if (data.hero) {
                activeHeroArticle = data.hero;
                const heroImg = document.getElementById('heroImage');
                const heroTitle = document.getElementById('heroTitle');
                const heroDesc = document.getElementById('heroDesc');
                const heroBadge = document.getElementById('heroBadge');
                const heroDateMeta = document.getElementById('heroDateMeta');

                if (heroImg) {
                    const fallbackHero = '/admin/api/maps/svg?district=' + encodeURIComponent(data.hero.district || 'Tamil Nadu');
                    heroImg.onerror = function() {
                        this.onerror = null;
                        this.src = fallbackHero;
                    };
                    heroImg.src = data.hero.thumbnail || fallbackHero;
                }
                if (heroTitle) heroTitle.textContent = decodeHtml(data.hero.title);
                if (heroDesc) heroDesc.textContent = decodeHtml(data.hero.description);
                if (heroBadge) {
                    if (data.hero.isViral) {
                        heroBadge.innerHTML = '🔥 வைரல் முதன்மைச் செய்தி &bull; ' + escapeHtml((data.hero.category || 'செய்திகள்').toUpperCase());
                        heroBadge.style.background = 'linear-gradient(135deg, #e11d48, #f43f5e)';
                        heroBadge.style.boxShadow = '0 2px 10px rgba(225, 29, 72, 0.4)';
                    } else {
                        heroBadge.textContent = (data.hero.category || 'செய்திகள்') + ' • ' + (data.hero.district || 'தமிழ்நாடு');
                        heroBadge.style.background = '#0284c7';
                        heroBadge.style.boxShadow = 'none';
                    }
                }
                if (heroDateMeta) {
                    const heroSrc = getSourceDomain(data.hero.sourceUrl);
                    heroDateMeta.innerHTML = '<span>🕒 ' + formatDate(data.hero.createdAt) + '</span>' + 
                        (heroSrc ? '<span>&bull;</span><span style="color:#38bdf8; font-weight:600;">' + escapeHtml(heroSrc) + '</span>' : '') +
                        (data.hero.isViral ? '<span>&bull;</span><span style="color:#f43f5e; font-weight:700;">🔥 வைரல் செய்தி</span>' : '') +
                        '<span>&bull;</span><span style="color:#94a3b8;">நிருபர்: ' + escapeHtml(data.hero.author || 'TN24 செய்திக் குழு') + '</span>';
                }
            } else {
                activeHeroArticle = null;
                const heroTitle = document.getElementById('heroTitle');
                const heroDesc = document.getElementById('heroDesc');
                const heroBadge = document.getElementById('heroBadge');
                const heroDateMeta = document.getElementById('heroDateMeta');
                if (heroTitle) heroTitle.textContent = 'இந்தத் தேடலுக்கு செய்திகள் எதுவும் கிடைக்கவில்லை.';
                if (heroDesc) heroDesc.textContent = 'வேறு பிரிவை அல்லது மாவட்டத்தைத் தேர்வு செய்யவும்.';
                if (heroBadge) heroBadge.textContent = 'செய்திகள் இல்லை';
                if (heroDateMeta) heroDateMeta.innerHTML = '';
            }

            // Hero Sub-Teasers (3 Items)
            const teasersContainer = document.getElementById('heroTeasers');
            if (teasersContainer) {
                if (data.heroTeasers && data.heroTeasers.length > 0) {
                    teasersContainer.innerHTML = data.heroTeasers.map(function(item) {
                        const src = getSourceDomain(item.sourceUrl);
                        return '<div class="hero-teaser-item" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                            '<div class="teaser-badge-meta">' +
                                '<span class="teaser-badge">' + escapeHtml(item.category || 'செய்திகள்') + '</span>' +
                                '<span class="teaser-date">🕒 ' + formatDate(item.createdAt) + (src ? ' &bull; ' + escapeHtml(src) : '') + '</span>' +
                            '</div>' +
                            '<div class="teaser-title">' + escapeHtml(item.title) + '</div>' +
                        '</div>';
                    }).join('');
                } else {
                    teasersContainer.innerHTML = '';
                }
            }

            // Left Column Mini Cards
            const leftContainer = document.getElementById('leftFeedCards');
            if (leftContainer) {
                if (data.leftFeed && data.leftFeed.length > 0) {
                    leftContainer.innerHTML = data.leftFeed.map(function(item) {
                        const fallback = '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
                        const thumb = item.thumbnail || fallback;
                        const src = getSourceDomain(item.sourceUrl);
                        return '<div class="mini-card" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                            '<div class="mini-thumb-wrap">' +
                                '<img src="' + thumb + '" onerror="this.onerror=null; this.src=\'' + fallback + '\'" alt="" loading="lazy">' +
                                '<span class="badge-pill">' + escapeHtml(item.category || 'செய்திகள்') + '</span>' +
                            '</div>' +
                            '<div class="mini-content">' +
                                '<h4 class="mini-title">' + escapeHtml(item.title) + '</h4>' +
                                '<div class="mini-meta">🕒 ' + formatDate(item.createdAt) + (src ? ' &bull; ' + escapeHtml(src) : '') + '</div>' +
                            '</div>' +
                        '</div>';
                    }).join('');
                } else {
                    leftContainer.innerHTML = '<div style="padding:16px; color:var(--text-muted); font-size:12px; text-align:center;">செய்திகள் இல்லை.</div>';
                }
            }

            // Press Releases List
            const pressContainer = document.getElementById('pressReleasesList');
            if (pressContainer) {
                if (data.pressReleases && data.pressReleases.length > 0) {
                    pressContainer.innerHTML = data.pressReleases.map(function(item) {
                        const src = getSourceDomain(item.sourceUrl);
                        return '<div class="press-item" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                            '<div class="press-title">' + escapeHtml(item.title) + '</div>' +
                            '<div class="press-meta">🕒 ' + formatDate(item.createdAt) + (src ? ' &bull; ' + escapeHtml(src) : '') + ' &bull; ' + escapeHtml(item.district || 'அரசு செய்தி') + '</div>' +
                        '</div>';
                    }).join('');
                } else {
                    pressContainer.innerHTML = '<div style="padding:12px; color:var(--text-muted); font-size:12px; text-align:center;">அறிவிப்புகள் இல்லை.</div>';
                }
            }

            // Center Editorial 2-Column Grid
            const gridContainer = document.getElementById('editorialGrid');
            if (gridContainer) {
                if (data.centerArticles && data.centerArticles.length > 0) {
                    if (currentViewMode === 'grouped') {
                        // Group posts by Category dynamically
                        const groups = {};
                        data.centerArticles.forEach(item => {
                            const cat = (item.category || 'News').trim();
                            if (!groups[cat]) groups[cat] = [];
                            groups[cat].push(item);
                        });
                        const catIcons = {
                            'Sports': '⚽',
                            'Politics': '🏛️',
                            'Business': '💼',
                            'Technical': '💻',
                            'Entertainment': '🎬',
                            'Crime': '🚨',
                            'News': '📰'
                        };
                        const catTamilNames = {
                            'Sports': 'விளையாட்டு & கிரிக்கெட்',
                            'Politics': 'அரசியல் & ஆட்சி',
                            'Business': 'வணிகம் & சந்தை',
                            'Technical': 'தொழில்நுட்பம் & ஏஐ',
                            'Entertainment': 'சினிமா & பொழுதுபோக்கு',
                            'Crime': 'குற்ற நிகழ்வுகள் & சட்டம்',
                            'News': 'செய்திகள் & பொதுமக்கள்'
                        };
                        let html = '';
                        Object.keys(groups).sort().forEach((catKey, gIdx) => {
                            const items = groups[catKey];
                            const icon = catIcons[catKey] || '📌';
                            const titleTa = catTamilNames[catKey] || catKey;
                            html += '<div style="grid-column: 1 / -1; margin-top:' + (gIdx === 0 ? '0' : '24px') + '; margin-bottom:12px; display:flex; align-items:center; justify-content:space-between; border-bottom:2px solid #38bdf8; padding-bottom:6px;">' +
                                '<div style="font-size:16px; font-weight:800; color:#38bdf8; display:flex; align-items:center; gap:8px;">' +
                                    '<span>' + icon + '</span> <span>' + escapeHtml(titleTa) + '</span>' +
                                '</div>' +
                                '<span style="font-size:11px; background:#1e293b; color:#94a3b8; padding:3px 8px; border-radius:12px; font-weight:700;">' + items.length + ' செய்திகள்</span>' +
                            '</div>';
                            items.forEach(item => {
                                const fallback = '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
                                const thumb = item.thumbnail || fallback;
                                const isTa = item.language === 'ta';
                                const langBadge = isTa ? '<span style="background:#b45309; color:#fef3c7; font-size:9px; font-weight:800; padding:2px 6px; border-radius:4px;">🇮🇳 தமிழ்</span>' : '<span style="background:#1d4ed8; color:#dbeafe; font-size:9px; font-weight:800; padding:2px 6px; border-radius:4px;">🇬🇧 EN</span>';
                                const src = getSourceDomain(item.sourceUrl);
                                const sourceDist = src ? (escapeHtml(src) + ' &bull; ' + escapeHtml(item.district || 'தமிழ்நாடு')) : escapeHtml(item.district || 'தமிழ்நாடு');
                                html += '<div class="article-card" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                                    '<div class="article-card-thumb">' +
                                        '<img src="' + thumb + '" onerror="this.onerror=null; this.src=\'' + fallback + '\'" alt="" loading="lazy">' +
                                        '<div style="position:absolute; top:8px; left:8px; display:flex; gap:4px;">' +
                                            '<span class="badge-pill">' + escapeHtml(item.category || 'செய்திகள்') + '</span>' +
                                            langBadge +
                                        '</div>' +
                                    '</div>' +
                                    '<div class="article-card-body">' +
                                        '<h3 class="article-card-title">' + escapeHtml(item.title) + '</h3>' +
                                        '<div class="article-card-meta">🕒 ' + formatDate(item.createdAt) + ' &bull; ' + sourceDist + '</div>' +
                                        '<p class="article-card-desc">' + escapeHtml(item.description || 'வட்டார நிகழ்வுகள் குறித்த நேரடி செய்தி தொகுப்பு.') + '</p>' +
                                    '</div>' +
                                '</div>';
                            });
                        });
                        gridContainer.innerHTML = html;
                    } else {
                        // Standard Editorial Feed View
                        let html = '';
                        data.centerArticles.forEach(function(item, index) {
                            if (index === 4) {
                                html += renderBannerHtml('infeed', b.infeed);
                            }
                            const fallback = '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
                            const thumb = item.thumbnail || fallback;
                            const isTa = item.language === 'ta';
                            const langBadge = isTa ? '<span style="background:#b45309; color:#fef3c7; font-size:9px; font-weight:800; padding:2px 6px; border-radius:4px; margin-left:4px;">🇮🇳 தமிழ்</span>' : '<span style="background:#1d4ed8; color:#dbeafe; font-size:9px; font-weight:800; padding:2px 6px; border-radius:4px; margin-left:4px;">🇬🇧 EN</span>';
                            const src = getSourceDomain(item.sourceUrl);
                            const sourceDist = src ? (escapeHtml(src) + ' &bull; ' + escapeHtml(item.district || 'தமிழ்நாடு')) : escapeHtml(item.district || 'தமிழ்நாடு');
                            html += '<div class="article-card" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                                '<div class="article-card-thumb">' +
                                    '<img src="' + thumb + '" onerror="this.onerror=null; this.src=\'' + fallback + '\'" alt="" loading="lazy">' +
                                    '<span class="badge-pill">' + escapeHtml(item.category || 'செய்திகள்') + '</span>' +
                                '</div>' +
                                '<div class="article-card-body">' +
                                    '<h3 class="article-card-title">' + escapeHtml(item.title) + langBadge + '</h3>' +
                                    '<div class="article-card-meta">🕒 ' + formatDate(item.createdAt) + ' &bull; ' + sourceDist + '</div>' +
                                    '<p class="article-card-desc">' + escapeHtml(item.description || 'தமிழ்நாடு மற்றும் வட்டார நிகழ்வுகள் குறித்த நேரடி செய்தி தொகுப்பு.') + '</p>' +
                                '</div>' +
                            '</div>';
                        });
                        gridContainer.innerHTML = html;
                    }
                } else {
                    gridContainer.innerHTML = '<div style="padding:40px; color:var(--text-muted); font-size:14px; text-align:center; grid-column:1/-1;">செய்திகள் எதுவும் கிடைக்கவில்லை.</div>';
                }
            }

            // Right Column: Most Read (1 to 5)
            const mostReadContainer = document.getElementById('mostReadList');
            if (mostReadContainer) {
                if (data.mostRead && data.mostRead.length > 0) {
                    mostReadContainer.innerHTML = data.mostRead.map(function(item, idx) {
                        return '<div class="most-read-item" onclick="openArticleByID(\'' + item.id + '\', event)">' +
                            '<div class="num-box">' + (idx + 1) + '</div>' +
                            '<div class="most-read-title">' + escapeHtml(item.title) + '</div>' +
                        '</div>';
                    }).join('');
                } else {
                    mostReadContainer.innerHTML = '<div style="padding:12px; color:var(--text-muted); font-size:12px; text-align:center;">செய்திகள் இல்லை.</div>';
                }
            }

            // Right Column: Real-time TN Events & Summits Calendar
            const eventsContainer = document.getElementById('calendarEventsList');
            if (eventsContainer) {
                if (data.events && data.events.length > 0) {
                    eventsContainer.innerHTML = data.events.slice(0, 4).map(function(ev) {
                        return '<div class="event-item" onclick="openEventModal(\'' + ev.id + '\')" style="cursor:pointer; background:rgba(15,23,42,0.6); border:1px solid rgba(56,189,248,0.15); border-radius:6px; padding:10px; transition:all 0.15s; margin-bottom:8px;" onmouseover="this.style.borderColor=\'#38bdf8\'; this.style.transform=\'translateX(3px)\'" onmouseout="this.style.borderColor=\'rgba(56,189,248,0.15)\'; this.style.transform=\'none\'">' +
                            '<div class="event-logo-badge" style="background:' + (ev.badgeColor || '#0284c7') + '; color:#fff; font-weight:800; border-radius:6px;">' + escapeHtml(ev.badge) + '</div>' +
                            '<div class="event-info" style="flex:1; overflow:hidden;">' +
                                '<div style="display:flex; justify-content:space-between; align-items:center; gap:4px; margin-bottom:2px;">' +
                                    '<span class="event-date" style="color:#38bdf8; font-weight:700;">' + escapeHtml(ev.dateDisplay) + '</span>' +
                                    '<span style="font-size:9px; font-weight:700; color:#4ade80; background:rgba(74,222,128,0.12); padding:1px 5px; border-radius:3px;">' + escapeHtml(ev.countdown) + '</span>' +
                                '</div>' +
                                '<span class="event-name" style="color:#f8fafc; font-size:11.5px; line-height:1.35;">' + escapeHtml(ev.title) + '</span>' +
                                '<span class="event-sub" style="color:#94a3b8; font-size:10px; margin-top:3px;">📍 ' + escapeHtml(ev.venue) + '</span>' +
                            '</div>' +
                        '</div>';
                    }).join('');
                } else {
                    eventsContainer.innerHTML = '<div style="padding:12px; color:var(--text-muted); font-size:12px; text-align:center;">நிகழ்வுகள் எதுவும் இல்லை.</div>';
                }
            }
        }

        // Article Modal Open/Close
        function openHeroArticle(e) {
            if (e && e.stopPropagation) e.stopPropagation();
            if (activeHeroArticle) {
                showArticleModal(activeHeroArticle);
            }
        }

        function openArticleByID(id, e) {
            if (e && e.stopPropagation) e.stopPropagation();
            if (!portalData) return;
            const all = [
                ...(portalData.hero ? [portalData.hero] : []),
                ...(portalData.heroTeasers || []),
                ...(portalData.leftFeed || []),
                ...(portalData.pressReleases || []),
                ...(portalData.centerArticles || []),
                ...(portalData.mostRead || [])
            ];
            const found = all.find(item => item.id === id);
            if (found) {
                showArticleModal(found);
            }
        }

        function showArticleModal(item) {
            currentModalArticle = item;
            if (window.history && window.history.replaceState) {
                window.history.replaceState(null, '', '/portal?post=' + encodeURIComponent(item.id));
            }
            document.getElementById('modalTitle').textContent = decodeHtml(item.title);
            document.getElementById('modalCategoryBadge').textContent = (item.category || 'செய்திகள்').toUpperCase();
            const src = getSourceDomain(item.sourceUrl);
            document.getElementById('modalDistrict').textContent = item.district || 'தமிழ்நாடு';
            document.getElementById('modalDate').textContent = '🕒 ' + formatDate(item.createdAt) + (src ? ' • ' + src : '');
            document.getElementById('modalDescription').textContent = decodeHtml(item.description || 'விளக்க உரை கிடைக்கவில்லை.');
            document.getElementById('modalAuthor').textContent = (src ? src + ' • ' : '') + 'நிருபர்: ' + (item.author || 'TN24 செய்திக் குழு');

            const mediaContainer = document.getElementById('modalMediaContainer');
            const vUrl = item.videoUrl || '';
            if (vUrl && (vUrl.includes('bbc.com/ws/av-embeds') || vUrl.includes('/embed/') || (item.videoType === 'bbc'))) {
                const src = vUrl.includes('bbc.com/ws/av-embeds') ? vUrl : ('https://www.bbc.com/ws/av-embeds/articles/' + item.videoId + '/ta');
                mediaContainer.innerHTML = '<iframe src="' + src + '" allow="autoplay; fullscreen; encrypted-media" allowfullscreen style="width:100%; height:380px; border:0; display:block;"></iframe>';
            } else if (vUrl && (vUrl.endsWith('.mp4') || vUrl.endsWith('.webm') || vUrl.includes('.mp4?') || vUrl.includes('.webm?'))) {
                mediaContainer.innerHTML = '<video controls autoplay style="width:100%; max-height:380px; display:block;" src="' + vUrl + '"></video>';
            } else if (vUrl && (vUrl.includes('youtube.com') || vUrl.includes('youtu.be') || item.videoId)) {
                let ytId = item.videoId;
                if (!ytId && vUrl.includes('/shorts/')) {
                    ytId = vUrl.split('/shorts/')[1].split('?')[0];
                } else if (!ytId && vUrl.includes('v=')) {
                    ytId = vUrl.split('v=')[1].split('&')[0];
                }
                if (ytId && !ytId.startsWith('vid_') && !ytId.startsWith('p0')) {
                    mediaContainer.innerHTML = '<iframe src="https://www.youtube.com/embed/' + ytId + '?autoplay=1" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen style="width:100%; height:380px; border:0; display:block;"></iframe>';
                } else {
                    const thumb = item.thumbnail || ('/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu'));
                    mediaContainer.innerHTML = '<img src="' + thumb + '" alt="">';
                }
            } else {
                const thumb = item.thumbnail || ('/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu'));
                mediaContainer.innerHTML = '<img src="' + thumb + '" alt="">';
            }

            document.getElementById('articleModal').style.display = 'flex';
            document.body.style.overflow = 'hidden';
        }

        function closeArticleModal() {
            document.getElementById('articleModal').style.display = 'none';
            document.getElementById('modalMediaContainer').innerHTML = '';
            document.body.style.overflow = '';
            if (window.history && window.history.replaceState) {
                window.history.replaceState(null, '', '/portal');
            }
        }

        function handleBackdropClick(e) {
            if (e.target.id === 'articleModal') {
                closeArticleModal();
            }
        }

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                closeArticleModal();
                closeEventModal();
                closeFullCalendarModal();
            }
        });

        function visitOriginalSource() {
            if (currentModalArticle && currentModalArticle.sourceUrl) {
                window.open(currentModalArticle.sourceUrl, '_blank');
            } else {
                showToast('மூல செய்தி இணைப்பு கிடைக்கவில்லை');
            }
        }

        function shareToWhatsApp() {
            if (!currentModalArticle) return;
            const title = currentModalArticle.title || '';
            const url = window.location.origin + '/portal?post=' + encodeURIComponent(currentModalArticle.id);
            const text = encodeURIComponent('🔥 *' + title + '*\n\nமுழு செய்தி விவரம் படிக்க:\n' + url + '\n\n— TN24 News');
            window.open('https://api.whatsapp.com/send?text=' + text, '_blank');
        }

        function shareNative() {
            if (!currentModalArticle) return;
            const url = window.location.origin + '/portal?post=' + encodeURIComponent(currentModalArticle.id);
            if (navigator.share) {
                navigator.share({
                    title: currentModalArticle.title,
                    text: currentModalArticle.title + ' — TN24 News',
                    url: url
                }).catch(() => {});
            } else {
                copyShareLink();
            }
        }

        function copyShareLink() {
            if (currentModalArticle) {
                const shareUrl = window.location.origin + '/portal?post=' + encodeURIComponent(currentModalArticle.id);
                navigator.clipboard.writeText(shareUrl).then(() => {
                    showToast('செய்தி இணைப்பு நகலெடுக்கப்பட்டது!');
                }).catch(() => {
                    showToast('இணைப்பு: ' + shareUrl);
                });
            }
        }

        let urlPostChecked = false;
        function checkUrlPost() {
            if (urlPostChecked) return;
            const urlParams = new URLSearchParams(window.location.search);
            const postId = urlParams.get('post') || (typeof window.INITIAL_POST_ID !== 'undefined' ? window.INITIAL_POST_ID : '');
            if (!postId) return;

            urlPostChecked = true;

            // Check in portalData
            if (portalData) {
                const all = [
                    ...(portalData.hero ? [portalData.hero] : []),
                    ...(portalData.heroTeasers || []),
                    ...(portalData.leftFeed || []),
                    ...(portalData.pressReleases || []),
                    ...(portalData.centerArticles || []),
                    ...(portalData.mostRead || [])
                ];
                const found = all.find(item => item.id === postId);
                if (found) {
                    showArticleModal(found);
                    return;
                }
            }

            // Fetch from single post API
            fetch('/api/portal/post?id=' + encodeURIComponent(postId))
                .then(r => r.json())
                .then(res => {
                    if (res.success && res.data) {
                        showArticleModal(res.data);
                    }
                })
                .catch(() => {});
        }

        // Real-time Single Event Modal Logic
        function openEventModal(eventId) {
            if (!portalData || !portalData.events) return;
            const ev = portalData.events.find(e => e.id === eventId);
            if (!ev) return;

            document.getElementById('eventModalTitle').textContent = ev.title;
            const catBadge = document.getElementById('eventModalCategoryBadge');
            if (catBadge) {
                catBadge.textContent = (ev.category || 'மாநாடு').toUpperCase();
                catBadge.style.background = ev.badgeColor || '#0284c7';
            }
            document.getElementById('eventModalDistrict').textContent = ev.district;
            document.getElementById('eventModalCountdown').textContent = ev.countdown;
            document.getElementById('eventModalDate').textContent = ev.dateDisplay;
            document.getElementById('eventModalVenue').textContent = ev.venue;
            document.getElementById('eventModalDesc').textContent = ev.description;

            // Generate Google Calendar Link
            const calTitle = encodeURIComponent(ev.title);
            const calLocation = encodeURIComponent(ev.venue + ', ' + ev.district + ', Tamil Nadu');
            const calDetails = encodeURIComponent(ev.description + '\n\nOfficial Source: TN24 Digital News Network (https://tn24.news)');
            const calUrl = 'https://calendar.google.com/calendar/render?action=TEMPLATE&text=' + calTitle + '&location=' + calLocation + '&details=' + calDetails;
            const calBtn = document.getElementById('eventGoogleCalBtn');
            if (calBtn) calBtn.href = calUrl;

            document.getElementById('eventModal').style.display = 'flex';
        }

        function closeEventModal() {
            const m = document.getElementById('eventModal');
            if (m) m.style.display = 'none';
        }

        function handleEventBackdropClick(e) {
            if (e.target && e.target.id === 'eventModal') {
                closeEventModal();
            }
        }

        // Full State Calendar Modal Logic
        function openAllEventsModal() {
            if (!portalData || !portalData.events || portalData.events.length === 0) {
                showToast('தமிழ்நாடு 2026 மாநாடுகள் புதுப்பிக்கப்படுகின்றன...');
                loadPortalFeed().then(() => {
                    if (portalData && portalData.events) {
                        renderFullCalendar(portalData.events);
                        const m = document.getElementById('fullStateCalendarModal');
                        if (m) m.style.display = 'flex';
                    }
                });
                return;
            }
            renderFullCalendar(portalData.events);
            const m = document.getElementById('fullStateCalendarModal');
            if (m) m.style.display = 'flex';
        }

        function closeFullCalendarModal() {
            const m = document.getElementById('fullStateCalendarModal');
            if (m) m.style.display = 'none';
        }

        function handleFullCalendarBackdropClick(e) {
            if (e.target && e.target.id === 'fullStateCalendarModal') {
                closeFullCalendarModal();
            }
        }

        function filterFullCalendar() {
            if (!portalData || !portalData.events) return;
            const q = (document.getElementById('fullCalendarSearch') ? document.getElementById('fullCalendarSearch').value : '').toLowerCase().trim();
            const cat = (document.getElementById('fullCalendarCategoryFilter') ? document.getElementById('fullCalendarCategoryFilter').value : 'ALL');

            const filtered = portalData.events.filter(function(ev) {
                const matchesCat = (cat === 'ALL' || ev.category === cat);
                const matchesQuery = (!q || ev.title.toLowerCase().includes(q) || ev.district.toLowerCase().includes(q) || ev.venue.toLowerCase().includes(q) || ev.description.toLowerCase().includes(q));
                return matchesCat && matchesQuery;
            });

            renderFullCalendar(filtered);
        }

        function renderFullCalendar(events) {
            const list = document.getElementById('fullCalendarList');
            const countEl = document.getElementById('fullCalendarCount');
            if (!list) return;

            if (countEl) {
                countEl.textContent = events.length + ' அதிகாரப்பூர்வ மாநாடுகள் காண்பிக்கப்படுகின்றன';
            }

            if (events.length === 0) {
                list.innerHTML = '<div style="grid-column: 1 / -1; padding: 40px; text-align: center; color: #94a3b8; font-size: 13px;">நிகழ்வுகள் எதுவும் இல்லை.</div>';
                return;
            }

            list.innerHTML = events.map(function(ev) {
                const calTitle = encodeURIComponent(ev.title);
                const calLocation = encodeURIComponent(ev.venue + ', ' + ev.district + ', Tamil Nadu');
                const calDetails = encodeURIComponent(ev.description + '\n\nOfficial Schedule: TN24 2026 Summits (https://tn24.news)');
                const calUrl = 'https://calendar.google.com/calendar/render?action=TEMPLATE&text=' + calTitle + '&location=' + calLocation + '&details=' + calDetails;

                return '<div style="background: rgba(15,23,42,0.85); border: 1px solid #1e293b; border-radius: 10px; padding: 16px; display: flex; flex-direction: column; justify-content: space-between; gap: 12px; transition: all 0.2s;" onmouseover="this.style.borderColor=\'#38bdf8\'; this.style.transform=\'translateY(-2px)\'" onmouseout="this.style.borderColor=\'#1e293b\'; this.style.transform=\'none\'">' +
                    '<div>' +
                        '<div style="display: flex; justify-content: space-between; align-items: flex-start; gap: 8px; margin-bottom: 8px;">' +
                            '<div style="display: flex; align-items: center; gap: 8px;">' +
                                '<span style="background:' + (ev.badgeColor || '#0284c7') + '; color: #fff; font-size: 11px; font-weight: 800; padding: 3px 8px; border-radius: 5px;">' + escapeHtml(ev.badge) + '</span>' +
                                '<span style="font-size: 11px; font-weight: 700; color: #38bdf8;">' + escapeHtml(ev.category) + '</span>' +
                            '</div>' +
                            '<span style="font-size: 10px; font-weight: 700; color: #4ade80; background: rgba(74,222,128,0.12); padding: 2px 7px; border-radius: 4px; border: 1px solid rgba(74,222,128,0.3);">' + escapeHtml(ev.countdown) + '</span>' +
                        '</div>' +
                        '<h3 style="font-size: 14px; font-weight: 800; color: #f8fafc; margin: 0 0 6px 0; line-height: 1.4;">' + escapeHtml(ev.title) + '</h3>' +
                        '<div style="font-size: 12px; color: #38bdf8; font-weight: 700; margin-bottom: 4px;">🗓️ ' + escapeHtml(ev.dateDisplay) + '</div>' +
                        '<div style="font-size: 11.5px; color: #94a3b8; margin-bottom: 8px;">📍 ' + escapeHtml(ev.venue) + ' • <strong style="color:#e2e8f0;">' + escapeHtml(ev.district) + '</strong></div>' +
                        '<p style="font-size: 12px; color: #cbd5e1; line-height: 1.5; margin: 0;">' + escapeHtml(ev.description) + '</p>' +
                    '</div>' +
                    '<div style="display: flex; gap: 8px; border-top: 1px solid rgba(255,255,255,0.06); padding-top: 12px; margin-top: 4px;">' +
                        '<a href="' + calUrl + '" target="_blank" class="action-btn" style="flex: 1; text-align: center; text-decoration: none; font-size: 11px; font-weight: 700; background: #0284c7; color: #fff; border-color: #38bdf8; padding: 7px 10px; border-radius: 6px;">📅 கூகிள் காலண்டர்</a>' +
                        '<button type="button" onclick="openEventModal(\'' + ev.id + '\')" class="action-btn" style="font-size: 11px; padding: 7px 12px; border-radius: 6px;">விவரங்கள் ↗</button>' +
                    '</div>' +
                '</div>';
            }).join('');
        }

        function setActiveNavTab(id) {
            document.querySelectorAll('.nav-links .nav-item').forEach(el => el.classList.remove('active'));
            if (id) {
                const target = document.getElementById(id);
                if (target) target.classList.add('active');
            }
        }

        // Filtering
        function filterDistrict(name) {
            const isAll = (!name || name === 'All Districts' || name === 'TN-ALL / REGIONAL' || name === 'அனைத்து மாவட்டங்கள் (38+)');
            document.getElementById('selectedDistrictLabel').textContent = isAll ? 'தமிழ்நாடு - அனைத்து வட்டாரங்கள்' : name.toUpperCase();
            setActiveNavTab('nav-item-districts');
            currentCategory = '';
            currentQuery = '';
            currentIsViral = false;
            loadPortalFeed(isAll ? '' : name, '', '', false);
            showToast(isAll ? 'அனைத்து மாவட்ட செய்திகளும் காண்பிக்கப்படுகின்றன' : name + ' மாவட்டச் செய்திகள்');
        }

        function filterCategory(cat) {
            currentQuery = '';
            if (!cat || cat === 'All') {
                currentCategory = '';
                currentDistrict = '';
                currentIsViral = false;
                document.getElementById('selectedDistrictLabel').textContent = 'தமிழ்நாடு - அனைத்து வட்டாரங்கள்';
                setActiveNavTab('nav-item-home');
                loadPortalFeed('', '', '', false);
                showToast('அனைத்து செய்திகளும் புதுப்பிக்கப்பட்டன');
                return;
            }

            currentIsViral = false;
            if (cat.toLowerCase() === 'sports') {
                setActiveNavTab('nav-item-sports');
            } else if (cat.toLowerCase() === 'news') {
                setActiveNavTab('nav-item-news');
            } else {
                setActiveNavTab('nav-item-categories');
            }
            loadPortalFeed('', cat, '', false);
            showToast('செய்திப் பிரிவு: ' + cat);
        }

        function filterViral() {
            currentQuery = '';
            if (currentIsViral) {
                currentIsViral = false;
                setActiveNavTab('nav-item-home');
                loadPortalFeed(currentDistrict || '', currentCategory || '', '', false);
                showToast('வைரல் வடிகட்டி நீக்கப்பட்டது');
            } else {
                currentIsViral = true;
                setActiveNavTab('nav-item-viral');
                loadPortalFeed('', '', '', true);
                showToast('🔥 வைரல் செய்திகள் வரிசைப்படுத்தப்பட்டுள்ளன');
            }
        }

        function handleSearchKey(e) {
            if (e.key === 'Enter') executeSearch();
        }

        function executeSearch() {
            const q = document.getElementById('searchInput').value.trim();
            loadPortalFeed('', '', q, false);
            if (q) showToast('தேடப்படுகிறது: ' + q);
        }

        function toggleTickerDropdown() {
            const name = prompt('மாவட்டத்தைத் தேர்வு செய்யவும் (உதா: சென்னை, கோவை, மதுரை, திருச்சி, அல்லது அனைத்து மாவட்டங்கள்):');
            if (name) filterDistrict(name);
        }

        // Market & Forex Spot Benchmarks (against INR)
        const marketRates = {
            gold22kPerGram: 6750, // ₹6,750 / g (₹54,000 / Sovereign 8g)
            gold24kPerGram: 7365, // ₹7,365 / g (₹73,650 / 10g)
            silverPerGram: 94.50,  // ₹94.50 / g (₹94,500 / 1 Kg)
            currencies: {
                USD: 87.40,
                AED: 23.80,
                EUR: 95.10,
                GBP: 111.80,
                SGD: 66.30,
                SAR: 23.30,
                KWD: 285.20,
                CAD: 63.80,
                AUD: 57.50,
                MYR: 19.80,
                QAR: 24.00
            }
        };

        function initMarketRates() {
            fetch('https://open.er-api.com/v6/latest/USD')
                .then(res => res.json())
                .then(data => {
                    if (data && data.rates && data.rates.INR) {
                        const usdInr = data.rates.INR;
                        marketRates.currencies.USD = usdInr;
                        ['AED', 'EUR', 'GBP', 'SGD', 'SAR', 'KWD', 'CAD', 'AUD', 'MYR', 'QAR'].forEach(curr => {
                            if (data.rates[curr]) {
                                marketRates.currencies[curr] = usdInr / data.rates[curr];
                            }
                        });
                        const tUsd = document.getElementById('ticker-usd');
                        if (tUsd) tUsd.textContent = '₹ ' + marketRates.currencies.USD.toFixed(2);
                        const tAed = document.getElementById('ticker-aed');
                        if (tAed) tAed.textContent = '₹ ' + marketRates.currencies.AED.toFixed(2);
                        const tEur = document.getElementById('ticker-eur');
                        if (tEur) tEur.textContent = '₹ ' + marketRates.currencies.EUR.toFixed(2);
                        runCalculator();
                    }
                })
                .catch(() => {});
        }

        // Interactive Market & Currency against INR Calculator
        function runCalculator() {
            // 1. Gold Calculator
            const goldQtyInput = document.getElementById('calcGoldQty');
            const goldUnitInput = document.getElementById('calcGoldUnit');
            const goldPurityInput = document.getElementById('calcGoldPurity');
            const goldOut = document.getElementById('calcGoldOut');

            if (goldQtyInput && goldUnitInput && goldPurityInput && goldOut) {
                const qty = parseFloat(goldQtyInput.value) || 0;
                const unit = goldUnitInput.value;
                const purity = goldPurityInput.value;

                let gramWeight = qty;
                if (unit === 'sovereign') gramWeight = qty * 8;
                else if (unit === '10g') gramWeight = qty * 10;
                else if (unit === '100g') gramWeight = qty * 100;

                const ratePerGram = purity === '24k' ? marketRates.gold24kPerGram : marketRates.gold22kPerGram;
                const goldTotal = gramWeight * ratePerGram;
                goldOut.textContent = '₹ ' + goldTotal.toLocaleString('en-IN', { maximumFractionDigits: 0 });
            }

            // 2. Silver Calculator
            const silverQtyInput = document.getElementById('calcSilverQty');
            const silverUnitInput = document.getElementById('calcSilverUnit');
            const silverOut = document.getElementById('calcSilverOut');

            if (silverQtyInput && silverUnitInput && silverOut) {
                const qty = parseFloat(silverQtyInput.value) || 0;
                const unit = silverUnitInput.value;
                const gramWeight = unit === 'kg' ? qty * 1000 : qty;
                const silverTotal = gramWeight * marketRates.silverPerGram;
                silverOut.textContent = '₹ ' + silverTotal.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
            }

            // 3. Currency against INR Calculator
            const currQtyInput = document.getElementById('calcCurrQty');
            const currTypeInput = document.getElementById('calcCurrType');
            const currOut = document.getElementById('calcCurrOut');

            if (currQtyInput && currTypeInput && currOut) {
                const qty = parseFloat(currQtyInput.value) || 0;
                const curr = currTypeInput.value;
                const rate = marketRates.currencies[curr] || 87.40;
                const inrTotal = qty * rate;
                currOut.textContent = '₹ ' + inrTotal.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
            }
        }

        function scrollToSection(id) {
            const el = document.getElementById(id);
            if (el) el.scrollIntoView({ behavior: 'smooth' });
        }

        function showToast(msg) {
            const toast = document.getElementById('toastNotify');
            if (!toast) return;
            toast.textContent = msg;
            toast.style.display = 'block';
            setTimeout(() => { toast.style.display = 'none'; }, 3000);
        }

        function getSourceDomain(url) {
            if (!url) return '';
            try {
                const u = new URL(url);
                let host = u.hostname.replace(/^www\./, '');
                if (host.includes('thehindu.com')) return 'தி இந்து (The Hindu)';
                if (host.includes('bbc.com') || host.includes('bbci.co.uk')) return 'பிபிசி தமிழ் (BBC Tamil)';
                if (host.includes('oneindia.com')) return 'ஒன்இந்தியா தமிழ்';
                if (host.includes('news18.com')) return 'News18 தமிழ்நாடு';
                if (host.includes('dinamalar.com')) return 'தினமலர்';
                if (host.includes('dailythanthi.com')) return 'தினத்தந்தி';
                if (host.includes('polimernews.com')) return 'பாலிமர் செய்திகள்';
                if (host.includes('dinamani.com')) return 'தினமணி';
                if (host.includes('puthiyathalaimurai.com')) return 'புதிய தலைமுறை';
                if (host.includes('maalaimalar.com')) return 'மாலை மலர்';
                if (host.includes('vikatan.com')) return 'விகடன்';
                if (host.includes('google.com')) return 'கூகுள் செய்திகள்';
                if (host.includes('youtube.com') || host.includes('youtu.be')) return 'யூடியூப் நேரலை';
                return host;
            } catch(e) {
                return '';
            }
        }

        function formatDate(dateStr) {
            if (!dateStr) return 'சமீபத்தியது';
            const d = new Date(dateStr);
            if (isNaN(d.getTime())) return 'சமீபத்தியது';
            const months = ["ஜன", "பிப்", "மார்ச்", "ஏப்", "மே", "ஜூன்", "ஜூலை", "ஆக", "செப்", "அக்", "நவ", "டிச"];
            const m = months[d.getMonth()];
            const day = d.getDate();
            const yr = d.getFullYear();
            let hr = d.getHours();
            const min = d.getMinutes() < 10 ? '0' + d.getMinutes() : d.getMinutes();
            const ampm = hr >= 12 ? 'பி.ப' : 'மு.ப';
            hr = hr % 12;
            hr = hr ? hr : 12;
            const hrStr = hr < 10 ? '0' + hr : hr;
            return m + ' ' + day + ', ' + yr + ' ' + hrStr + ':' + min + ' ' + ampm;
        }

        function decodeHtml(str) {
            if (!str) return '';
            return String(str)
                .replace(/&#x([0-9a-fA-F]+);?/gi, (m, hex) => String.fromCharCode(parseInt(hex, 16)))
                .replace(/&#([0-9]+);?/g, (m, dec) => String.fromCharCode(parseInt(dec, 10)))
                .replace(/&quot;/g, '"')
                .replace(/&apos;/g, "'")
                .replace(/&amp;/g, '&')
                .replace(/&lt;/g, '<')
                .replace(/&gt;/g, '>')
                .replace(/&nbsp;/g, ' ');
        }

        function escapeHtml(str) {
            if (!str) return '';
            const clean = decodeHtml(str);
            return clean
                .replace(/&/g, '&amp;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;')
                .replace(/"/g, '&quot;')
                .replace(/'/g, '&#039;');
        }
    </script>
</body>
</html>`
}
