package admin

import (
	"fmt"
	"html"
)

// RenderLoginPage produces the high-contrast, premium dark themed login page for TN24 Admin Console
func RenderLoginPage(errorMsg string) string {
	errHtml := ""
	if errorMsg != "" {
		errHtml = fmt.Sprintf(`
            <div class="error-banner" id="errorBanner">
                <span class="error-icon">⚠️</span>
                <span>%s</span>
            </div>
        `, html.EscapeString(errorMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ta">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TN24 Console &mdash; நிர்வாகி உள்நுழைவு (Admin Login)</title>
    <link rel="icon" type="image/svg+xml" href="/portal/assets/brand/tn24-icon.svg">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-base: #030712;
            --bg-card: rgba(15, 23, 42, 0.85);
            --border-glow: rgba(56, 189, 248, 0.25);
            --primary-cyan: #38bdf8;
            --primary-blue: #0284c7;
            --accent-red: #ef4444;
            --text-main: #f8fafc;
            --text-muted: #94a3b8;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        }

        body {
            background-color: var(--bg-base);
            background-image: 
                radial-gradient(circle at 15%% 20%%, rgba(56, 189, 248, 0.08) 0%%, transparent 40%%),
                radial-gradient(circle at 85%% 80%%, rgba(239, 68, 68, 0.06) 0%%, transparent 40%%),
                linear-gradient(180deg, #070d18 0%%, #030712 100%%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
            color: var(--text-main);
        }

        .login-card {
            background: var(--bg-card);
            backdrop-filter: blur(16px);
            -webkit-backdrop-filter: blur(16px);
            border: 1px solid var(--border-glow);
            box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.7), 0 0 30px rgba(56, 189, 248, 0.1);
            border-radius: 16px;
            width: 100%%;
            max-width: 420px;
            padding: 36px 32px;
            text-align: center;
            position: relative;
            overflow: hidden;
            animation: fadeIn 0.3s ease-out;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(12px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .login-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 3px;
            background: linear-gradient(90deg, #38bdf8, #6366f1, #ef4444);
        }

        .brand-header {
            margin-bottom: 24px;
        }

        .brand-logo-img {
            height: 48px;
            display: inline-block;
            margin-bottom: 12px;
            filter: drop-shadow(0 4px 12px rgba(56, 189, 248, 0.3));
        }

        .login-title {
            font-size: 20px;
            font-weight: 800;
            color: #ffffff;
            letter-spacing: -0.3px;
            margin-bottom: 4px;
        }

        .login-sub {
            font-size: 12px;
            color: var(--text-muted);
            font-weight: 500;
        }

        .error-banner {
            background: rgba(239, 68, 68, 0.15);
            border: 1px solid rgba(239, 68, 68, 0.4);
            color: #fca5a5;
            padding: 10px 14px;
            border-radius: 8px;
            font-size: 12px;
            font-weight: 600;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 8px;
            text-align: left;
            animation: shake 0.3s ease-in-out;
        }

        @keyframes shake {
            0%%, 100%% { transform: translateX(0); }
            25%% { transform: translateX(-4px); }
            75%% { transform: translateX(4px); }
        }

        .form-group {
            margin-bottom: 18px;
            text-align: left;
        }

        .form-label {
            display: block;
            font-size: 11.5px;
            font-weight: 700;
            color: #cbd5e1;
            margin-bottom: 6px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .input-wrap {
            position: relative;
            display: flex;
            align-items: center;
        }

        .input-icon {
            position: absolute;
            left: 14px;
            color: #64748b;
            font-size: 14px;
            pointer-events: none;
        }

        .form-input {
            width: 100%%;
            background: rgba(2, 6, 23, 0.7);
            border: 1px solid #334155;
            border-radius: 8px;
            padding: 12px 14px 12px 38px;
            color: #ffffff;
            font-size: 14px;
            outline: none;
            transition: all 0.2s;
        }

        .form-input:focus {
            border-color: var(--primary-cyan);
            box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
            background: rgba(15, 23, 42, 0.9);
        }

        .password-toggle {
            position: absolute;
            right: 12px;
            background: transparent;
            border: none;
            color: #64748b;
            cursor: pointer;
            padding: 4px;
            font-size: 14px;
            transition: color 0.15s;
        }

        .password-toggle:hover {
            color: var(--primary-cyan);
        }

        .btn-submit {
            width: 100%%;
            background: linear-gradient(135deg, #0284c7 0%%, #0369a1 100%%);
            border: 1px solid rgba(56, 189, 248, 0.4);
            color: #ffffff;
            padding: 12px 16px;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 800;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            transition: all 0.2s;
            margin-top: 10px;
            box-shadow: 0 4px 16px rgba(2, 132, 199, 0.3);
        }

        .btn-submit:hover {
            background: linear-gradient(135deg, #0369a1 0%%, #075985 100%%);
            transform: translateY(-1px);
            box-shadow: 0 6px 20px rgba(2, 132, 199, 0.45);
        }

        .btn-submit:active {
            transform: translateY(0);
        }

        .login-footer {
            margin-top: 24px;
            padding-top: 18px;
            border-top: 1px solid rgba(255, 255, 255, 0.07);
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 11.5px;
            color: #64748b;
        }

        .footer-link {
            color: var(--primary-cyan);
            text-decoration: none;
            font-weight: 600;
            transition: color 0.15s;
        }

        .footer-link:hover {
            color: #7dd3fc;
            text-decoration: underline;
        }

        .badge-live {
            display: inline-flex;
            align-items: center;
            gap: 5px;
            background: rgba(74, 222, 128, 0.12);
            color: #4ade80;
            padding: 2px 8px;
            border-radius: 12px;
            font-weight: 700;
            font-size: 10px;
            border: 1px solid rgba(74, 222, 128, 0.25);
        }

        .live-dot {
            width: 6px;
            height: 6px;
            border-radius: 50%%;
            background: #4ade80;
            box-shadow: 0 0 6px #4ade80;
            animation: pulse 1.5s infinite;
        }

        @keyframes pulse {
            0%%, 100%% { opacity: 1; transform: scale(1); }
            50%% { opacity: 0.4; transform: scale(0.85); }
        }
    </style>
</head>
<body>

    <div class="login-card">
        <div class="brand-header">
            <a href="/portal">
                <img src="/portal/assets/brand/tn24-logo.svg" alt="TN24 Logo" class="brand-logo-img">
            </a>
            <h1 class="login-title">நிர்வாக மையம்</h1>
            <p class="login-sub">TN24 Management &amp; Editorial Console Access</p>
        </div>

        %s

        <form method="POST" action="/admin/login" id="loginForm">
            <div class="form-group">
                <label class="form-label" for="username">பயனர் பெயர் (Username)</label>
                <div class="input-wrap">
                    <span class="input-icon">👤</span>
                    <input type="text" id="username" name="username" class="form-input" placeholder="Enter username" required autofocus autocomplete="username">
                </div>
            </div>

            <div class="form-group">
                <label class="form-label" for="password">கடவுச்சொல் (Password)</label>
                <div class="input-wrap">
                    <span class="input-icon">🔒</span>
                    <input type="password" id="password" name="password" class="form-input" placeholder="Enter password" required autocomplete="current-password">
                    <button type="button" class="password-toggle" id="togglePasswordBtn" onclick="togglePasswordVisibility()" title="Show/Hide Password">👁️</button>
                </div>
            </div>

            <button type="submit" class="btn-submit" id="submitBtn">
                <span>உள்நுழைக (Sign In)</span>
                <span>➔</span>
            </button>
        </form>

        <div class="login-footer">
            <span class="badge-live"><span class="live-dot"></span> TN24 24/7 Engine</span>
            <a href="/portal" class="footer-link">பொது தளம் (Portal) ↗</a>
        </div>
    </div>

    <script>
        function togglePasswordVisibility() {
            const pwd = document.getElementById('password');
            const btn = document.getElementById('togglePasswordBtn');
            if (pwd.type === 'password') {
                pwd.type = 'text';
                btn.textContent = '🙈';
            } else {
                pwd.type = 'password';
                btn.textContent = '👁️';
            }
        }

        document.getElementById('loginForm').addEventListener('submit', function(e) {
            const btn = document.getElementById('submitBtn');
            btn.innerHTML = '<span>சரிபார்க்கிறது...</span>';
            btn.style.opacity = '0.7';
            btn.style.pointerEvents = 'none';
        });
    </script>
</body>
</html>`, errHtml)
}
