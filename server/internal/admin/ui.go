package admin

func RenderAdminDashboard() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TN24 Control Hub &mdash; Management &amp; Analytics Console</title>
    <link rel="icon" type="image/svg+xml" href="/portal/assets/brand/tn24-icon.svg">
    <style>
        :root {
            --bg-primary: #09090b;
            --bg-surface: #141417;
            --bg-card: #1c1c21;
            --border: #2e2e38;
            --text-primary: #f4f4f6;
            --text-muted: #a1a1aa;
            --brand-orange: #ff5722;
            --brand-orange-light: #ff784e;
            --color-success: #10b981;
            --color-warning: #f59e0b;
            --color-danger: #ef4444;
            --color-info: #38bdf8;
            --color-purple: #a855f7;
            --focus-ring: #ff784e;
        }

        * { box-sizing: border-box; margin: 0; padding: 0; }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.5;
            padding: 0;
            display: flex;
            min-height: 100vh;
        }

        .skip-link {
            position: absolute;
            top: -40px;
            left: 10px;
            background: var(--brand-orange);
            color: #fff;
            padding: 8px 16px;
            z-index: 1000;
            text-decoration: none;
            font-weight: 700;
            border-radius: 4px;
            transition: top 0.2s;
        }
        .skip-link:focus { top: 10px; outline: 3px solid #fff; }

        a:focus-visible, button:focus-visible, input:focus-visible, select:focus-visible, [tabindex]:focus-visible {
            outline: 3px solid var(--focus-ring) !important;
            outline-offset: 2px !important;
        }

        aside {
            width: 270px;
            background: var(--bg-surface);
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            padding: 24px 16px;
            flex-shrink: 0;
        }

        .brand-logo {
            display: flex;
            align-items: center;
            gap: 10px;
            padding-bottom: 24px;
            border-bottom: 1px solid var(--border);
            margin-bottom: 24px;
        }

        .brand-badge {
            background: linear-gradient(135deg, var(--brand-orange), #ea580c);
            color: #fff;
            font-weight: 900;
            font-size: 14px;
            padding: 6px 10px;
            border-radius: 8px;
            letter-spacing: 1px;
        }

        .brand-title {
            font-size: 18px;
            font-weight: 800;
            color: var(--text-primary);
            letter-spacing: -0.5px;
        }

        .nav-list {
            list-style: none;
            display: flex;
            flex-direction: column;
            gap: 6px;
        }

        .nav-btn {
            display: flex;
            align-items: center;
            gap: 12px;
            width: 100%;
            background: transparent;
            border: 1px solid transparent;
            color: var(--text-muted);
            padding: 10px 14px;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            text-align: left;
            transition: all 0.15s ease;
        }

        .nav-btn:hover {
            background: rgba(255, 255, 255, 0.04);
            color: var(--text-primary);
        }

        .nav-btn.active {
            background: rgba(255, 87, 34, 0.12);
            color: var(--brand-orange-light);
            border-color: rgba(255, 87, 34, 0.3);
        }

        .nav-btn.agent-btn {
            border-left: 3px solid var(--color-purple);
        }
        .nav-btn.agent-btn.active {
            background: rgba(168, 85, 247, 0.15);
            color: #c084fc;
            border-color: rgba(168, 85, 247, 0.4);
        }

        .sidebar-footer {
            margin-top: auto;
            padding-top: 20px;
            border-top: 1px solid var(--border);
            font-size: 12px;
            color: var(--text-muted);
        }

        main {
            flex-grow: 1;
            overflow-y: auto;
            padding: 32px 40px;
            display: flex;
            flex-direction: column;
            gap: 28px;
        }

        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--border);
        }

        .header-title h1 {
            font-size: 26px;
            font-weight: 800;
            color: var(--text-primary);
        }

        .header-title p {
            color: var(--text-muted);
            font-size: 14px;
            margin-top: 2px;
        }

        .live-status-pill {
            display: flex;
            align-items: center;
            gap: 8px;
            background: rgba(16, 185, 129, 0.1);
            color: var(--color-success);
            padding: 6px 14px;
            border-radius: 9999px;
            font-size: 13px;
            font-weight: 700;
            border: 1px solid rgba(16, 185, 129, 0.3);
        }

        .status-dot {
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--color-success);
            box-shadow: 0 0 8px var(--color-success);
        }

        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
            gap: 16px;
        }

        .metric-card {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 20px;
            display: flex;
            flex-direction: column;
            gap: 8px;
            transition: transform 0.15s ease, border-color 0.15s ease;
        }

        .metric-card:hover {
            transform: translateY(-2px);
            border-color: rgba(255, 87, 34, 0.4);
        }

        .metric-card .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 13px;
            font-weight: 600;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .metric-card .card-value {
            font-size: 32px;
            font-weight: 900;
            color: var(--text-primary);
            line-height: 1;
        }

        .metric-card .card-footer {
            font-size: 12px;
            display: flex;
            align-items: center;
            gap: 6px;
        }

        .badge-pill {
            display: inline-block;
            padding: 2px 8px;
            border-radius: 6px;
            font-size: 11px;
            font-weight: 700;
        }
        .badge-success { background: rgba(16, 185, 129, 0.15); color: var(--color-success); border: 1px solid rgba(16, 185, 129, 0.3); }
        .badge-warning { background: rgba(245, 158, 11, 0.15); color: var(--color-warning); border: 1px solid rgba(245, 158, 11, 0.3); }
        .badge-danger { background: rgba(239, 68, 68, 0.15); color: var(--color-danger); border: 1px solid rgba(239, 68, 68, 0.3); }
        .badge-info { background: rgba(56, 189, 248, 0.15); color: var(--color-info); border: 1px solid rgba(56, 189, 248, 0.3); }
        .badge-purple { background: rgba(168, 85, 247, 0.15); color: #c084fc; border: 1px solid rgba(168, 85, 247, 0.3); }

        .charts-grid {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: 20px;
        }

        @media (max-width: 992px) {
            body { flex-direction: column; }
            aside { width: 100%; }
            .charts-grid { grid-template-columns: 1fr; }
        }

        .chart-box {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 16px;
        }

        .chart-box h2 {
            font-size: 16px;
            font-weight: 700;
            color: var(--text-primary);
        }

        .bar-chart-row {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 12px;
        }

        .bar-chart-label {
            width: 110px;
            font-size: 13px;
            font-weight: 600;
            color: var(--text-primary);
        }

        .bar-chart-track {
            flex-grow: 1;
            background: rgba(255, 255, 255, 0.05);
            height: 18px;
            border-radius: 6px;
            overflow: hidden;
            display: flex;
        }

        .bar-chart-fill {
            height: 100%;
            background: linear-gradient(90deg, var(--brand-orange), #f97316);
            border-radius: 6px;
            transition: width 0.8s ease-in-out;
        }

        .bar-chart-value {
            width: 60px;
            text-align: right;
            font-size: 13px;
            font-weight: 700;
            color: var(--text-muted);
        }

        .table-box {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 20px;
            overflow-x: auto;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            text-align: left;
            font-size: 14px;
        }

        th {
            background: rgba(255, 255, 255, 0.03);
            color: var(--text-muted);
            font-weight: 700;
            padding: 12px 16px;
            border-bottom: 1px solid var(--border);
            text-transform: uppercase;
            font-size: 11px;
            letter-spacing: 0.5px;
        }

        td {
            padding: 14px 16px;
            border-bottom: 1px solid var(--border);
            color: var(--text-primary);
        }

        tr:hover td {
            background: rgba(255, 255, 255, 0.02);
        }

        .action-btn {
            background: transparent;
            border: 1px solid var(--border);
            color: var(--text-primary);
            padding: 6px 12px;
            border-radius: 6px;
            font-size: 12px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.15s ease;
        }
        .action-btn:hover {
            border-color: var(--brand-orange);
            color: var(--brand-orange-light);
        }
        .action-btn-danger:hover {
            border-color: var(--color-danger);
            color: var(--color-danger);
        }

        .btn-primary {
            background: var(--brand-orange);
            color: #fff;
            border: none;
            font-weight: 700;
            padding: 8px 16px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 13px;
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
        .btn-primary:hover { background: #ea580c; }

        .btn-purple {
            background: linear-gradient(135deg, #a855f7, #7e22ce);
            color: #fff;
            border: none;
            font-weight: 700;
            padding: 8px 18px;
            border-radius: 8px;
            cursor: pointer;
            font-size: 13px;
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
        .btn-purple:hover { opacity: 0.9; }

        .btn-approve {
            background: var(--color-success);
            color: #000;
            border: none;
            font-weight: 700;
            padding: 6px 12px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
        }
        .btn-approve:hover { background: #34d399; }

        /* Google Sign-In Card */
        .google-auth-card {
            max-width: 440px;
            margin: 40px auto;
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            padding: 36px 32px;
            text-align: center;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
        }

        .google-logo-badge {
            width: 52px;
            height: 52px;
            background: #fff;
            border-radius: 50%;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 20px;
            font-size: 26px;
            box-shadow: 0 4px 12px rgba(255, 255, 255, 0.2);
        }

        /* AI Agent Workspace */
        .agent-container {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        .agent-header-card {
            background: linear-gradient(135deg, rgba(168, 85, 247, 0.1), rgba(28, 28, 33, 0.8));
            border: 1px solid rgba(168, 85, 247, 0.3);
            border-radius: 12px;
            padding: 20px 24px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .agent-chat-box {
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            gap: 16px;
            min-height: 380px;
            max-height: 540px;
            overflow-y: auto;
        }

        .chat-bubble {
            padding: 14px 18px;
            border-radius: 10px;
            font-size: 14px;
            line-height: 1.6;
            max-width: 85%;
            word-break: break-word;
        }

        .chat-bubble.agent {
            background: rgba(255, 255, 255, 0.04);
            border: 1px solid var(--border);
            color: var(--text-primary);
            align-self: flex-start;
        }

        .chat-bubble.user {
            background: rgba(168, 85, 247, 0.2);
            border: 1px solid rgba(168, 85, 247, 0.4);
            color: #f3e8ff;
            align-self: flex-end;
        }

        .action-card {
            background: #121216;
            border: 1px solid rgba(168, 85, 247, 0.4);
            border-radius: 10px;
            padding: 14px 18px;
            margin-top: 12px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .quick-chips {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
            margin-top: 6px;
        }

        .chip {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid var(--border);
            color: var(--text-muted);
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.15s ease;
        }
        .chip:hover {
            border-color: var(--color-purple);
            color: #c084fc;
            background: rgba(168, 85, 247, 0.1);
        }

        .prompt-bar {
            display: flex;
            gap: 12px;
        }
        .prompt-input {
            flex-grow: 1;
            background: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 12px 16px;
            color: #fff;
            font-size: 14px;
        }

        /* Modal Dialog */
        .modal {
            display: none;
            position: fixed;
            z-index: 10000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.7);
            align-items: center;
            justify-content: center;
        }
        .modal-content {
            background: var(--bg-card);
            border: 1px solid var(--border);
            padding: 28px;
            border-radius: 12px;
            width: 90%;
            max-width: 500px;
            color: var(--text-primary);
        }

        /* Custom Confirmation & Rejection Modal Popup */
        .custom-modal-backdrop {
            display: none;
            position: fixed;
            inset: 0;
            z-index: 99999;
            background: rgba(3, 7, 18, 0.75);
            backdrop-filter: blur(8px);
            -webkit-backdrop-filter: blur(8px);
            align-items: center;
            justify-content: center;
            padding: 16px;
        }
        .custom-modal-card {
            background: #0f172a;
            border: 1px solid #1e293b;
            border-radius: 16px;
            width: 100%;
            max-width: 480px;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.8), 0 0 0 1px rgba(255, 255, 255, 0.05);
            overflow: hidden;
            animation: modalPop 0.22s cubic-bezier(0.16, 1, 0.3, 1);
        }
        @keyframes modalPop {
            from { opacity: 0; transform: scale(0.94) translateY(-10px); }
            to { opacity: 1; transform: scale(1) translateY(0); }
        }
        .reason-chip {
            font-size: 11px;
            padding: 4px 10px;
            border-radius: 6px;
            background: rgba(255, 255, 255, 0.05);
            color: #94a3b8;
            border: 1px solid rgba(255, 255, 255, 0.1);
            cursor: pointer;
            transition: all 0.15s ease;
            user-select: none;
        }
        .reason-chip:hover, .reason-chip.selected {
            background: rgba(245, 158, 11, 0.2);
            color: #f59e0b;
            border-color: #f59e0b;
        }
        .form-group {
            margin-bottom: 16px;
            display: flex;
            flex-direction: column;
            gap: 6px;
        }
        .form-group label {
            font-size: 13px;
            font-weight: 600;
            color: var(--text-muted);
        }
        .form-control {
            background: var(--bg-surface);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 10px;
            color: #fff;
            font-size: 14px;
        }

        #toast {
            position: fixed;
            bottom: 24px;
            right: 24px;
            background: #18181b;
            color: #fff;
            padding: 12px 20px;
            border-radius: 8px;
            border: 1px solid var(--brand-orange);
            box-shadow: 0 10px 25px rgba(0, 0, 0, 0.6);
            display: none;
            font-weight: 600;
            font-size: 14px;
            z-index: 9999;
        }

        .console-box {
            background: #0d0d11;
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 16px;
            font-family: monospace;
            font-size: 13px;
            color: #38bdf8;
            max-height: 240px;
            overflow-y: auto;
            white-space: pre-wrap;
        }

        .sr-only {
            position: absolute;
            width: 1px;
            height: 1px;
            padding: 0;
            margin: -1px;
            overflow: hidden;
            clip: rect(0, 0, 0, 0);
            white-space: nowrap;
            border-width: 0;
        }
    </style>
</head>
<body>
    <a href="#main-content" class="skip-link">Skip to main content</a>
    <div id="aria-live-status" role="status" aria-live="polite" class="sr-only"></div>

    <!-- Navigation Sidebar -->
    <aside aria-label="Main Navigation">
        <div class="brand-logo" style="padding-bottom: 18px; margin-bottom: 20px; border-bottom: 1px solid var(--border);">
            <a href="/portal" target="_blank" title="Open Live Portal in new tab" style="display: flex; align-items: center; text-decoration: none;">
                <img src="/portal/assets/brand/tn24-logo.svg?v=20260908d" alt="TN24" style="height: 38px; display: block;" />
            </a>
        </div>

        <nav>
            <ul class="nav-list" role="tablist">
                <li>
                    <button role="tab" id="tab-overview" aria-selected="true" aria-controls="panel-overview" class="nav-btn active" onclick="switchTab('overview')">
                        📊 Overview &amp; Stats
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-agent" aria-selected="false" aria-controls="panel-agent" class="nav-btn agent-btn" onclick="switchTab('agent')">
                        🤖 AI Agent
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-cron" aria-selected="false" aria-controls="panel-cron" class="nav-btn" onclick="switchTab('cron')">
                        ⏱️ Cron &amp; Background Jobs
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-moderation" aria-selected="false" aria-controls="panel-moderation" class="nav-btn" onclick="switchTab('moderation')">
                        🛡️ Content Moderation
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-grievances" aria-selected="false" aria-controls="panel-grievances" class="nav-btn" onclick="switchTab('grievances')">
                        ⚖️ Grievances (IT Rules)
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-contributors" aria-selected="false" aria-controls="panel-contributors" class="nav-btn" onclick="switchTab('contributors')">
                        ⭐ Contributor Trust
                    </button>
                </li>
                <li>
                    <button role="tab" id="tab-audit" aria-selected="false" aria-controls="panel-audit" class="nav-btn" onclick="switchTab('audit')">
                        📜 Compliance Audit Log
                    </button>
                </li>
            </ul>
        </nav>

        <div class="sidebar-footer" style="padding: 14px; background: rgba(15,23,42,0.85); border-radius: 8px; border: 1px solid rgba(56,189,248,0.2); margin-top: auto;">
            <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; padding-bottom: 8px; border-bottom: 1px solid rgba(255,255,255,0.08);">
                <div style="display: flex; align-items: center; gap: 6px;">
                    <span style="width: 7px; height: 7px; border-radius: 50%; background: #4ade80; box-shadow: 0 0 6px #4ade80;"></span>
                    <span style="font-size: 11px; font-weight: 800; color: #38bdf8;">👤 Admin Online</span>
                </div>
                <a href="/admin/logout" class="action-btn" style="color: #fca5a5; border-color: rgba(239,68,68,0.4); background: rgba(239,68,68,0.12); font-size: 10px; padding: 2px 7px; text-decoration: none; border-radius: 4px; font-weight: 700; display: inline-flex; align-items: center; gap: 3px;" title="Sign out of TN24 admin console">🚪 Logout</a>
            </div>
            <p><strong>District:</strong> Madurai Core (TN)</p>
            <p>Compliance: Indian IT Rules 2021</p>
            <p style="margin-top: 4px; color: var(--color-success);">● Live PostgreSQL &amp; Redis</p>
        </div>
    </aside>

    <!-- Main Content Panels -->
    <main id="main-content">
        <!-- Overview Panel -->
        <section id="panel-overview" role="tabpanel" aria-labelledby="tab-overview">
            <header>
                <div class="header-title">
                    <h1>Real-Time Database Operations &amp; Telemetry</h1>
                    <p>Live metrics from PostgreSQL tables and background workers across Tamil Nadu</p>
                </div>
                <div class="live-status-pill" aria-label="Infrastructure state: Online and connected">
                    <span class="status-dot"></span>
                    <span id="live-infra-badge">Connecting...</span>
                </div>
            </header>

            <div style="height: 20px;"></div>

            <!-- Top KPI Cards -->
            <div class="metrics-grid">
                <div class="metric-card">
                    <div class="card-header">
                        <span>Total Content</span>
                        <span>📦</span>
                    </div>
                    <div class="card-value" id="kpi-total">0</div>
                    <div class="card-footer">
                        <span class="badge-pill badge-info">Real DB Records</span>
                    </div>
                </div>

                <div class="metric-card">
                    <div class="card-header">
                        <span>Pending Review</span>
                        <span>⏳</span>
                    </div>
                    <div class="card-value" style="color: var(--color-warning);" id="kpi-pending">0</div>
                    <div class="card-footer">
                        <span class="badge-pill badge-warning">Moderation Intake</span>
                    </div>
                </div>

                <div class="metric-card">
                    <div class="card-header">
                        <span>Quarantined</span>
                        <span>🛡️</span>
                    </div>
                    <div class="card-value" style="color: var(--color-danger);" id="kpi-quarantine">0</div>
                    <div class="card-footer">
                        <span class="badge-pill badge-danger">Toxicity Flagged</span>
                    </div>
                </div>

                <div class="metric-card">
                    <div class="card-header">
                        <span>IT Rules SLA</span>
                        <span>⏱️</span>
                    </div>
                    <div class="card-value" style="color: var(--color-success);" id="kpi-sla">100%</div>
                    <div class="card-footer">
                        <span class="badge-pill badge-success">Statutory 24h/15d</span>
                    </div>
                </div>
            </div>

            <div style="height: 24px;"></div>

            <!-- Visual Charts -->
            <div class="charts-grid">
                <div class="chart-box">
                    <h2>📍 Geographic Distribution of Content (Real Districts)</h2>
                    <div id="district-chart">
                        <p style="color: var(--text-muted); font-size: 13px;">Loading district metrics from database...</p>
                    </div>
                </div>

                <div class="chart-box">
                    <h2>🎬 Content by Format</h2>
                    <div id="format-breakdown" style="display: flex; flex-direction: column; gap: 14px; margin-top: 10px;">
                        <p style="color: var(--text-muted); font-size: 13px;">Loading format breakdown...</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Dedicated AI Agent Panel -->
        <section id="panel-agent" role="tabpanel" aria-labelledby="tab-agent" style="display: none;">
            <!-- Genuine Google OAuth 2.0 Sign-In Gate -->
            <div id="agent-auth-gate" class="google-auth-card">
                <div class="google-logo-badge">G</div>
                <h2 style="font-size: 22px; margin-bottom: 8px;">Antigravity AI Agent</h2>
                <p style="color: var(--text-muted); font-size: 13px; margin-bottom: 24px;">
                    Protected administrative console. Requires authenticating with your genuine Google / Gmail account via Google Identity Services.
                </p>

                <!-- Official Google OAuth 2.0 Button -->
                <button type="button" class="btn-purple" onclick="initiateGenuineGoogleOAuth()" style="width: 100%; justify-content: center; padding: 14px; font-size: 14px; margin-bottom: 16px;">
                    <span style="font-size: 18px; margin-right: 6px;">🌐</span> Sign in with Google (OAuth 2.0)
                </button>

                <p style="font-size: 12px; color: var(--text-muted); line-height: 1.5;">
                    Redirects to <code>accounts.google.com</code> for cryptographically verified Google authentication.
                </p>

                <div style="margin-top: 20px; border-top: 1px solid var(--border); padding-top: 16px;">
                    <a href="/admin/auth/google" style="color: #38bdf8; font-size: 12px; text-decoration: none;">⚙️ Configure / Review Google Client Credentials &rarr;</a>
                </div>
            </div>

            <!-- Authenticated AI Agent Console -->
            <div id="agent-workspace" class="agent-container" style="display: none;">
                <header class="agent-header-card" style="background: linear-gradient(135deg, #13131a 0%, #1e1b2e 100%); border: 1px solid rgba(168, 85, 247, 0.35); padding: 18px 24px; border-radius: 12px;">
                    <div style="display: flex; align-items: center; gap: 16px;">
                        <div style="width: 48px; height: 48px; border-radius: 12px; background: linear-gradient(135deg, #a855f7, #06b6d4); display: flex; align-items: center; justify-content: center; font-size: 26px; box-shadow: 0 0 20px rgba(168, 85, 247, 0.4);">
                            🚀
                        </div>
                        <div>
                            <div style="display: flex; align-items: center; gap: 10px;">
                                <h2 style="font-size: 19px; font-weight: 800; color: #fff; margin: 0; letter-spacing: -0.02em;">Google Antigravity IDE • Operations Copilot</h2>
                                <span style="font-size: 11px; background: rgba(56, 189, 248, 0.15); color: #38bdf8; border: 1px solid rgba(56, 189, 248, 0.4); padding: 2px 8px; border-radius: 12px; font-weight: 700;">DeepMind Runtime</span>
                            </div>
                            <p id="operator-badge-text" style="color: #c084fc; font-size: 12px; font-weight: 600; margin-top: 4px;">
                                Verified Google Operator: rahamath1986@gmail.com
                            </p>
                        </div>
                    </div>
                    <div style="display: flex; align-items: center; gap: 10px;">
                        <span class="badge-pill" style="background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.4); font-weight: 700;">● Continuous Inference Active</span>
                        <button class="action-btn" onclick="handleSignOut()">Sign Out</button>
                    </div>
                </header>

                <div class="agent-chat-box" id="chat-messages">
                    <div class="chat-bubble agent">
                        <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px;">
                            <strong>🚀 Google Antigravity IDE Agent</strong>
                            <span style="font-size: 10px; background: rgba(168, 85, 247, 0.2); color: #c084fc; padding: 2px 6px; border-radius: 6px;">Google DeepMind</span>
                        </div>
                        Vanakkam! I am the <strong>Google Antigravity Agent</strong> running with continuous generative model inference.<br><br>
                        I am connected directly to your PostgreSQL database (40 districts, content pipelines), Redis cache, and cron job scheduler. What operational task would you like to coordinate today?
                    </div>
                </div>

                <div class="quick-chips">
                    <span class="chip" onclick="sendQuickPrompt('Tell me about yourself')">✨ Introduce Yourself</span>
                    <span class="chip" onclick="sendQuickPrompt('list all administrative districts')">📍 List 40 Districts</span>
                    <span class="chip" onclick="sendQuickPrompt('Run full telemetry and database diagnostics')">⚡ Full System Diagnostics</span>
                    <span class="chip" onclick="sendQuickPrompt('Triage pending submissions with Tamil NLP toxicity check')">🛡️ Triage Moderation</span>
                    <span class="chip" onclick="sendQuickPrompt('Check statutory IT Rules 2021 grievance SLA compliance')">⚖️ Grievance SLA</span>
                    <span class="chip" onclick="sendQuickPrompt('Scan external video links for broken embeds')">⏱️ Scan Video Links</span>
                </div>

                <form class="prompt-bar" onsubmit="handleSendPrompt(event)">
                    <input id="prompt-input" class="prompt-input" placeholder="Instruct the Antigravity IDE agent (e.g. 'Tell me about yourself', 'Audit Madurai feed')..." />
                    <button type="submit" class="btn-purple">Send</button>
                </form>
            </div>
        </section>

        <!-- Cron Jobs Panel -->
        <section id="panel-cron" role="tabpanel" aria-labelledby="tab-cron" style="display: none;">
            <header>
                <div class="header-title">
                    <h1>Scheduled Cron &amp; Background Job Manager</h1>
                    <p>Create, control, pause, and trigger live server background jobs with execution logs</p>
                </div>
                <button class="btn-primary" onclick="openNewJobModal()">➕ New Scheduled Job</button>
            </header>

            <div style="height: 20px;"></div>

            <div class="table-box">
                <table aria-label="Cron Jobs Management Table">
                    <thead>
                        <tr>
                            <th scope="col">Job Name</th>
                            <th scope="col">Interval</th>
                            <th scope="col">Type</th>
                            <th scope="col">Status</th>
                            <th scope="col">Last Run</th>
                            <th scope="col">Runs / Fails</th>
                            <th scope="col">Actions</th>
                        </tr>
                    </thead>
                    <tbody id="cron-table-body">
                        <tr><td colspan="7" style="color: var(--text-muted);">Loading active cron jobs...</td></tr>
                    </tbody>
                </table>
            </div>

            <div style="height: 24px;"></div>

            <div class="chart-box">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <h2>📜 Live Cron Execution Log Output</h2>
                    <button class="action-btn" onclick="fetchCronLogs()">Refresh Logs</button>
                </div>
                <div id="cron-console" class="console-box">Connecting to execution stream...</div>
            </div>
        </section>

        <!-- Moderation & Staged News Review Panel -->
        <section id="panel-moderation" role="tabpanel" aria-labelledby="tab-moderation" style="display: none;">
            <header style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px;">
                <div class="header-title">
                    <h1>📰 Staging &amp; Content Moderation Review Desk</h1>
                    <p>Review scraped news, video stories, and community submissions before publishing live to app users</p>
                </div>
                <div style="display: flex; gap: 8px;">
                    <button class="btn-approve" onclick="approveAllPending()">⚡ Approve All Safe</button>
                    <button class="action-btn" onclick="triggerDeduplicate()" style="background: rgba(99,102,241,0.2); border-color: rgba(99,102,241,0.4); color: #a5b4fc; font-weight: 600;" title="Auto-remove duplicate posts and retain only the latest version">🧹 Deduplicate (Keep Latest)</button>
                    <button class="action-btn" onclick="fetchPendingContent()">🔄 Refresh Queue</button>
                </div>
            </header>

            <div style="height: 16px;"></div>

            <!-- Scrape Any Site Link Input Box -->
            <div class="chart-box" style="margin-bottom: 20px; background: linear-gradient(135deg, rgba(16,185,129,0.06), rgba(59,130,246,0.06)); border: 1px solid rgba(16,185,129,0.25);">
                <div style="font-weight: 600; font-size: 14px; margin-bottom: 8px; display: flex; align-items: center; gap: 8px;">
                    <span style="color: #38bdf8;">🔍 Deep Scrape Any News / Media Site Link:</span>
                </div>
                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                    <input id="scrape-url-input" class="form-control" style="flex: 1; min-width: 300px;" placeholder="Enter website link, article URL, or RSS feed (e.g. https://news.google.com/rss/headlines/section/geo/Tamil%20Nadu)" />
                    <button id="btn-scrape-url" class="btn-primary" onclick="triggerSiteScrape()">⚡ Scrape &amp; Stage for Review</button>
                </div>
                <div style="font-size: 11px; color: var(--text-muted); margin-top: 6px;">
                    💡 Scrapes all Tamil Nadu related text stories, video embeds (YouTube / direct), and images. All items enter <strong>PENDING</strong> status for your verification before publication.
                </div>
            </div>

            <!-- Content Retention & Lifecycle Control Box -->
            <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: 10px; padding: 14px 18px; margin-bottom: 16px;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; flex-wrap: wrap; gap: 8px;">
                    <div style="display: flex; align-items: center; gap: 10px;">
                        <span style="font-size: 16px;">⏳</span>
                        <div>
                            <span style="font-weight: 700; font-size: 14px; color: #f8fafc;">Content Retention &amp; Lifecycle Policy</span>
                            <span style="background: rgba(56,189,248,0.15); color: #38bdf8; font-size: 11px; padding: 2px 8px; border-radius: 4px; margin-left: 8px; font-weight: 600;">Scraped Content</span>
                            <span style="background: rgba(16,185,129,0.15); color: #10b981; font-size: 11px; padding: 2px 8px; border-radius: 4px; margin-left: 6px; font-weight: 600;">🔒 Manual Posts Protected (Never Cleared)</span>
                        </div>
                    </div>
                    <div style="display: flex; align-items: center; gap: 10px; font-size: 12px; color: var(--text-muted);">
                        <span>Policy: <strong id="retention-current-badge" style="color: #38bdf8;">24 Hours</strong></span>
                        <span>•</span>
                        <span>Total Items: <strong id="retention-total-count" style="color: #f8fafc;">--</strong></span>
                        <span>•</span>
                        <span>Stale Items: <strong id="retention-stale-count" style="color: #f43f5e;">--</strong></span>
                    </div>
                </div>

                <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px; background: rgba(0,0,0,0.25); border: 1px solid rgba(255,255,255,0.05); border-radius: 8px; padding: 10px 14px;">
                    <div style="display: flex; align-items: center; gap: 10px; flex-wrap: wrap;">
                        <label style="font-size: 12px; color: #cbd5e1; font-weight: 600;">Keep Content For:</label>
                        <div style="display: flex; align-items: center; gap: 6px;">
                            <input id="retention-hours-input" type="number" min="1" max="720" value="24" style="width: 70px; background: #1e293b; border: 1px solid #38bdf8; color: #fff; border-radius: 6px; padding: 4px 8px; font-size: 13px; font-weight: 700; text-align: center;" />
                            <span style="font-size: 12px; color: #94a3b8;">hours</span>
                        </div>
                        <div style="display: flex; gap: 4px;">
                            <button type="button" class="action-btn" style="font-size: 11px; padding: 3px 8px;" onclick="setRetentionPreset(12)">12h</button>
                            <button type="button" class="action-btn" style="font-size: 11px; padding: 3px 8px; color: #38bdf8; border-color: #38bdf8;" onclick="setRetentionPreset(24)">24h (Default)</button>
                            <button type="button" class="action-btn" style="font-size: 11px; padding: 3px 8px;" onclick="setRetentionPreset(48)">48h</button>
                            <button type="button" class="action-btn" style="font-size: 11px; padding: 3px 8px;" onclick="setRetentionPreset(72)">72h</button>
                        </div>
                        <label style="display: flex; align-items: center; gap: 6px; font-size: 12px; color: #94a3b8; cursor: pointer; margin-left: 8px;">
                            <input id="retention-auto-cleanup" type="checkbox" checked style="cursor: pointer;" />
                            <span>Auto-purge background scheduler (every 15m)</span>
                        </label>
                        <button type="button" class="btn-primary" style="font-size: 12px; padding: 5px 12px;" onclick="saveRetentionSettings()">💾 Save Policy</button>
                    </div>

                    <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
                        <span id="retention-last-cleanup-note" style="font-size: 11px; color: var(--text-muted);">Last cleanup: --</span>
                        <button type="button" class="btn-secondary" style="font-size: 12px; padding: 5px 12px; background: rgba(99,102,241,0.15); color: #a5b4fc; border: 1px solid rgba(99,102,241,0.4); border-radius: 6px; cursor: pointer; font-weight: 600;" onclick="triggerDeduplicate()" title="Auto-remove duplicate posts and retain only the latest version">
                            🧹 Deduplicate (Keep Latest)
                        </button>
                        <button type="button" class="btn-danger" style="font-size: 12px; padding: 5px 14px; background: rgba(244,63,94,0.15); color: #f43f5e; border: 1px solid rgba(244,63,94,0.4); border-radius: 6px; cursor: pointer; font-weight: 600;" onclick="triggerRetentionCleanupNow()" title="Immediately purges all content items older than the configured hours across all statuses">
                            🧹 Clean Stale Content Now
                        </button>
                    </div>
                </div>
                <div style="font-size: 11px; color: var(--text-muted); margin-top: 8px;">
                    💡 Scraped news/videos in any status (Pending, Published, Rejected) older than the policy window are auto-purged. <strong>Manually created posts are permanently protected and will never be deleted.</strong>
                </div>
            </div>

            <!-- Portal Banners & TN24 Brand Showcase Manager Box -->
            <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: 10px; padding: 16px; margin-bottom: 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.2);">
                <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 12px; border-bottom: 1px solid #1e293b; padding-bottom: 12px; margin-bottom: 12px;">
                    <div style="display: flex; align-items: center; gap: 10px;">
                        <span style="font-size: 20px;">🌟</span>
                        <div>
                            <span style="font-weight: 800; font-size: 14px; color: #f8fafc;">TN24 Portal Banners &amp; Showcase Manager</span>
                            <div style="font-size: 11px; color: #94a3b8;">Assign any story to the Main Hero Banner or Advertisement Slots across the live TN24 portal</div>
                        </div>
                    </div>
                    <div style="display: flex; gap: 8px; align-items: center;">
                        <a href="/portal" target="_blank" class="action-btn" style="color: #38bdf8; border-color: #38bdf8; font-size: 12px; text-decoration: none; padding: 4px 10px;">🌐 Open TN24 Live Portal ↗</a>
                        <button type="button" class="action-btn" onclick="loadBannerConfig()" style="font-size: 12px; padding: 4px 8px;">🔄 Refresh Slots</button>
                    </div>
                </div>

                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 14px;">
                    <!-- Slot 1: Main Hero Banner -->
                    <div style="background: rgba(30, 41, 59, 0.6); border: 1px solid rgba(234, 179, 8, 0.4); border-radius: 8px; padding: 12px;">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                            <span style="font-size: 12px; font-weight: 800; color: #eab308;">⭐ MAIN HERO BANNER</span>
                            <button type="button" class="action-btn" style="font-size: 10px; padding: 2px 6px;" onclick="resetMainHeroBanner()" title="Reset to auto viral top story">↺ Auto</button>
                        </div>
                        <div id="banner-slot-main-display" style="font-size: 12px; color: #cbd5e1; min-height: 48px; display: flex; align-items: center; gap: 8px;">
                            <span style="color: var(--text-muted); font-size: 11px;">⚡ Auto: Top Viral Story</span>
                        </div>
                    </div>

                    <!-- Slot 2: Header Leaderboard (728x90) -->
                    <div style="background: rgba(30, 41, 59, 0.4); border: 1px solid #334155; border-radius: 8px; padding: 12px;">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                            <span style="font-size: 12px; font-weight: 700; color: #38bdf8;">📢 HEADER LEADERBOARD (728x90)</span>
                            <button type="button" class="action-btn" style="font-size: 10px; padding: 2px 6px;" onclick="resetAdSlot('header')" title="Reset to default contact banner">↺ Default</button>
                        </div>
                        <div id="banner-slot-header-display" style="font-size: 12px; color: #cbd5e1; min-height: 48px; display: flex; align-items: center; gap: 8px;">
                            <img src="/portal/assets/brand/tn24-header.svg" style="width:52px; height:22px; object-fit:cover; border-radius:4px; border:1px dashed #38bdf8;" />
                            <div style="flex:1; overflow:hidden;">
                                <div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (728x90)</div>
                                <div style="font-size:10px; color:var(--text-muted);">Active: Leaderboard Advertising Banner</div>
                            </div>
                        </div>
                    </div>

                    <!-- Slot 3: Sidebar Rectangle (250x250) -->
                    <div style="background: rgba(30, 41, 59, 0.4); border: 1px solid #334155; border-radius: 8px; padding: 12px;">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                            <span style="font-size: 12px; font-weight: 700; color: #38bdf8;">📢 SIDEBAR BANNER (250x250)</span>
                            <button type="button" class="action-btn" style="font-size: 10px; padding: 2px 6px;" onclick="resetAdSlot('sidebar')" title="Reset to contact banner">↺ Default</button>
                        </div>
                        <div id="banner-slot-sidebar-display" style="font-size: 12px; color: #cbd5e1; min-height: 48px; display: flex; align-items: center; gap: 8px;">
                            <img src="/portal/assets/brand/tn24-sidebar.svg" style="width:26px; height:26px; object-fit:cover; border-radius:4px; border:1px dashed #38bdf8;" />
                            <div style="flex:1; overflow:hidden;">
                                <div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (250x250)</div>
                                <div style="font-size:10px; color:var(--text-muted);">Active: Sidebar 250x250 Advertising Banner</div>
                            </div>
                        </div>
                    </div>

                    <!-- Slot 4: Left Square & In-Feed -->
                    <div style="background: rgba(30, 41, 59, 0.4); border: 1px solid #334155; border-radius: 8px; padding: 12px;">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                            <span style="font-size: 12px; font-weight: 700; color: #38bdf8;">📢 IN-FEED &amp; SQUARE SLOTS</span>
                            <button type="button" class="action-btn" style="font-size: 10px; padding: 2px 6px;" onclick="resetAdSlot('infeed')" title="Reset to contact banner">↺ Default</button>
                        </div>
                        <div id="banner-slot-infeed-display" style="font-size: 12px; color: #cbd5e1; min-height: 48px; display: flex; align-items: center; gap: 8px;">
                            <img src="/portal/assets/brand/tn24-infeed.svg" style="width:52px; height:22px; object-fit:cover; border-radius:4px; border:1px dashed #a855f7;" />
                            <div style="flex:1; overflow:hidden;">
                                <div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (In-Feed &amp; Square)</div>
                                <div style="font-size:10px; color:var(--text-muted);">Active: In-Feed 728x90 &amp; Square 200x200 Advertising</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; flex-wrap: wrap; gap: 10px;">
                <div style="display: flex; gap: 6px; align-items: center; flex-wrap: wrap;">
                    <button id="filter-btn-pending" class="action-btn" style="background: rgba(56,189,248,0.15); color: #38bdf8; border-color: #38bdf8;" onclick="filterContentStatus('PENDING')">🟡 Pending Review (<span id="mod-tab-count-pending">0</span>)</button>
                    <button id="filter-btn-viral" class="action-btn" style="color: #f97316; border-color: rgba(249,115,22,0.4);" onclick="toggleViralFilter()">🔥 Viral Priority</button>
                    <button id="filter-btn-published" class="action-btn" onclick="filterContentStatus('PUBLISHED')">🟢 Approved (Live) (<span id="mod-tab-count-published">0</span>)</button>
                    <button id="filter-btn-rejected" class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.3);" onclick="filterContentStatus('REJECTED')">🔴 Rejected (<span id="discarded-count-badge">0</span>)</button>
                    <button id="filter-btn-all" class="action-btn" onclick="filterContentStatus('ALL')">📋 All Contents (<span id="mod-tab-count-all">0</span>)</button>
                    <button id="btn-empty-trash" class="action-btn" style="display: none; color: #f43f5e; border-color: #f43f5e; background: rgba(244,63,94,0.15); font-weight: 700;" onclick="emptyTrash()">🗑️ Delete All Rejected</button>
                </div>

                <!-- Search Posts Input Box -->
                <div style="display: flex; align-items: center; gap: 6px; background: #18181b; border: 1px solid #3f3f4e; border-radius: 8px; padding: 4px 10px; min-width: 280px; flex: 1; max-width: 440px;">
                    <span style="color: #38bdf8; font-size: 14px;">🔍</span>
                    <input id="content-search-input" type="text" placeholder="Search posts by headline, keyword, district, or source..." style="background: transparent; border: none; color: #f8fafc; font-size: 13px; outline: none; width: 100%;" onkeyup="handleContentSearchKey(event)" />
                    <button id="content-search-clear-btn" onclick="clearContentSearch()" style="background: transparent; border: none; color: #94a3b8; cursor: pointer; font-size: 16px; display: none; padding: 0 4px;" title="Clear search">&times;</button>
                    <button class="action-btn" onclick="triggerContentSearch()" style="font-size: 11px; padding: 4px 10px; background: #38bdf8; color: #09090b; border-color: #38bdf8; font-weight: 700; border-radius: 6px;">Search</button>
                </div>

                <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
                    <button class="action-btn" onclick="triggerReclassifyAll()" style="font-size: 12px; padding: 6px 12px; color: #a855f7; border-color: rgba(168,85,247,0.4);" title="Re-scan and accurately classify all content geographically">
                        🔄 Auto-Reclassify All
                    </button>
                    <button class="action-btn" onclick="triggerRefetchAllText()" style="font-size: 12px; padding: 6px 12px; color: #10b981; border-color: rgba(16,185,129,0.4);" title="Visit original news source webpages to extract full unabridged text paragraphs and photos">
                        📖 Refetch Long Text
                    </button>
                    <button class="btn-primary" onclick="openManualContentModal()" style="font-size: 12px; padding: 6px 14px; background: linear-gradient(135deg, #38bdf8, #0284c7); display: inline-flex; align-items: center; gap: 6px;">
                        <span style="font-weight: 900; font-size: 14px;">+</span> Add Content Manually
                    </button>
                    <div style="font-size: 12px; color: var(--text-muted);" id="content-status-summary">
                        Showing pending review queue
                    </div>
                </div>
            </div>

            <!-- Multi-Filter & Grouping Command Bar (Language, District, Category, Grouping) -->
            <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: 8px; padding: 10px 14px; margin-bottom: 12px; display: flex; flex-wrap: wrap; gap: 12px; align-items: center; justify-content: space-between;">
                <div style="display: flex; align-items: center; gap: 12px; flex-wrap: wrap;">
                    <!-- Language Filter & Config -->
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="font-size: 11px; font-weight: 700; color: #94a3b8; text-transform: uppercase;">🌐 Language:</span>
                        <select id="mod-lang-select" onchange="setModLanguageFilter(this.value)" style="background: #1e293b; border: 1px solid #38bdf8; color: #38bdf8; font-weight: 600; padding: 4px 8px; border-radius: 6px; font-size: 12px; cursor: pointer;">
                            <option value="all">🌐 All Languages</option>
                            <option value="ta">🇮🇳 தமிழ் (Tamil)</option>
                            <option value="en">🇬🇧 English</option>
                            <option value="ta-en">🔄 Tanglish (TA-EN)</option>
                            <option value="hi">🇮🇳 हिन्दी (Hindi)</option>
                            <option value="ml">🇮🇳 മലയാളம் (Malayalam)</option>
                            <option value="te">🇮🇳 తెలుగు (Telugu)</option>
                            <option value="kn">🇮🇳 ಕನ್ನಡ (Kannada)</option>
                        </select>
                        <button type="button" class="action-btn" onclick="openLanguageConfigModal()" style="font-size: 11px; padding: 4px 8px; color: #38bdf8; border-color: rgba(56,189,248,0.4); display: inline-flex; align-items: center; gap: 4px;" title="Platform language & localization settings">
                            ⚙️ Config
                        </button>
                    </div>

                    <div style="width: 1px; height: 20px; background: #334155;"></div>

                    <!-- District Filter -->
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="font-size: 11px; font-weight: 700; color: #94a3b8; text-transform: uppercase;">📍 District:</span>
                        <select id="mod-district-select" onchange="applyModFilters()" style="background: #1e293b; border: 1px solid #334155; color: #f8fafc; padding: 4px 8px; border-radius: 6px; font-size: 12px; max-width: 160px; cursor: pointer;">
                            <option value="all">All Districts (38+)</option>
                            <option value="Chennai">Chennai</option>
                            <option value="Coimbatore">Coimbatore</option>
                            <option value="Madurai">Madurai</option>
                            <option value="Tiruchirappalli">Tiruchirappalli</option>
                            <option value="Salem">Salem</option>
                            <option value="Tirunelveli">Tirunelveli</option>
                            <option value="Erode">Erode</option>
                            <option value="Vellore">Vellore</option>
                            <option value="Thanjavur">Thanjavur</option>
                            <option value="Kanyakumari">Kanyakumari</option>
                            <option value="Dindigul">Dindigul</option>
                            <option value="Cuddalore">Cuddalore</option>
                            <option value="Kancheepuram">Kancheepuram</option>
                            <option value="Chengalpattu">Chengalpattu</option>
                            <option value="Tiruvallur">Tiruvallur</option>
                            <option value="Tiruppur">Tiruppur</option>
                            <option value="Ranipet">Ranipet</option>
                            <option value="Tirupattur">Tirupattur</option>
                            <option value="Tiruvannamalai">Tiruvannamalai</option>
                            <option value="Viluppuram">Viluppuram</option>
                            <option value="Kallakurichi">Kallakurichi</option>
                            <option value="Dharmapuri">Dharmapuri</option>
                            <option value="Krishnagiri">Krishnagiri</option>
                            <option value="Namakkal">Namakkal</option>
                            <option value="Nilgiris">Nilgiris</option>
                            <option value="Karur">Karur</option>
                            <option value="Perambalur">Perambalur</option>
                            <option value="Ariyalur">Ariyalur</option>
                            <option value="Nagapattinam">Nagapattinam</option>
                            <option value="Mayiladuthurai">Mayiladuthurai</option>
                            <option value="Tiruvarur">Tiruvarur</option>
                            <option value="Pudukkottai">Pudukkottai</option>
                            <option value="Sivaganga">Sivaganga</option>
                            <option value="Ramanathapuram">Ramanathapuram</option>
                            <option value="Virudhunagar">Virudhunagar</option>
                            <option value="Theni">Theni</option>
                            <option value="Tenkasi">Tenkasi</option>
                            <option value="Thoothukudi">Thoothukudi</option>
                            <option value="Tamil Nadu">Tamil Nadu (Statewide)</option>
                        </select>
                    </div>

                    <div style="width: 1px; height: 20px; background: #334155;"></div>

                    <!-- Category Filter -->
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="font-size: 11px; font-weight: 700; color: #94a3b8; text-transform: uppercase;">🏷️ Category:</span>
                        <select id="mod-category-select" onchange="applyModFilters()" style="background: #1e293b; border: 1px solid #334155; color: #f8fafc; padding: 4px 8px; border-radius: 6px; font-size: 12px; max-width: 150px; cursor: pointer;">
                            <option value="all">All Categories</option>
                            <option value="News">News &amp; Civic</option>
                            <option value="Politics">Politics &amp; Govt</option>
                            <option value="Sports">Sports &amp; Athletics</option>
                            <option value="Business">Business &amp; Markets</option>
                            <option value="Technical">Tech &amp; Science</option>
                            <option value="Entertainment">Cinema &amp; Arts</option>
                            <option value="Crime">Crime &amp; Law</option>
                        </select>
                    </div>

                    <div style="width: 1px; height: 20px; background: #334155;"></div>

                    <!-- Sort By Selector -->
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="font-size: 11px; font-weight: 800; color: #38bdf8; text-transform: uppercase;">⇅ Sort By:</span>
                        <select id="mod-sort-select" onchange="changeModSort(this.value)" style="background: #1e293b; border: 1px solid #38bdf8; color: #38bdf8; padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 700; cursor: pointer;">
                            <option value="date_desc">📅 Post Date: Newest First</option>
                            <option value="date_asc">📅 Post Date: Oldest First</option>
                            <option value="source_date">📰 Source &amp; Post Date (A-Z)</option>
                            <option value="source_desc">📰 Source &amp; Post Date (Z-A)</option>
                            <option value="viral">🔥 Viral Priority First</option>
                        </select>
                    </div>
                </div>

                <!-- Group By Selector -->
                <div style="display: flex; align-items: center; gap: 8px;">
                    <span style="font-size: 11px; font-weight: 800; color: #eab308; text-transform: uppercase;">📑 Group Posts By:</span>
                    <select id="mod-group-by-select" onchange="changeModGroupBy(this.value)" style="background: #1e293b; border: 1px solid #eab308; color: #fef08a; padding: 4px 10px; border-radius: 6px; font-size: 12px; font-weight: 700; cursor: pointer;">
                        <option value="none">Flat List (Default)</option>
                        <option value="source">📰 Group by Source Provider (BBC, Thanthi, Dinamalar...)</option>
                        <option value="language">🌐 Group by Language (Tamil / English)</option>
                        <option value="district">📍 Group by District (38 Districts)</option>
                        <option value="category">🏷️ Group by Category (Sports, Politics, etc.)</option>
                    </select>
                </div>
            </div>

            <!-- Active Search Filter Badge Bar -->
            <div id="active-search-badge-bar" style="display: none; align-items: center; justify-content: space-between; background: rgba(56,189,248,0.08); border: 1px solid rgba(56,189,248,0.25); border-radius: 8px; padding: 8px 14px; margin-bottom: 12px;">
                <div style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #e2e8f0;">
                    <span>🔍 Filtering posts matching: <strong id="active-search-query-text" style="color: #38bdf8;"></strong></span>
                    <span id="active-search-count-badge" style="background: rgba(56,189,248,0.2); color: #38bdf8; font-size: 11px; padding: 2px 8px; border-radius: 4px; font-weight: 700;"></span>
                </div>
                <button class="action-btn" style="font-size: 11px; padding: 3px 10px; color: #f43f5e; border-color: rgba(244,63,94,0.4);" onclick="clearContentSearch()">✕ Clear Search</button>
            </div>

            <!-- Bulk Selection & Action Bar -->
            <div id="bulk-action-bar" style="display: none; align-items: center; justify-content: space-between; background: rgba(56,189,248,0.1); border: 1px solid rgba(56,189,248,0.3); border-radius: 8px; padding: 10px 16px; margin-bottom: 12px; flex-wrap: wrap; gap: 8px;">
                <div style="display: flex; align-items: center; gap: 12px;">
                    <span id="selected-count-badge" style="font-weight: 700; color: #38bdf8; font-size: 13px;">0 items selected</span>
                    <button class="action-btn" style="font-size: 11px; padding: 3px 8px;" onclick="clearSelection()">Clear</button>
                </div>
                <div style="display: flex; gap: 8px;">
                    <button class="btn-approve" style="font-size: 12px; padding: 6px 14px;" onclick="approveSelectedContent()">✓ Approve Selected</button>
                    <button class="action-btn" style="color: #f59e0b; border-color: rgba(245,158,11,0.4); background: rgba(245,158,11,0.1); font-weight: 600; font-size: 12px; padding: 5px 12px;" onclick="rejectSelectedContent()">🚫 Reject Selected</button>
                    <button class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.4); background: rgba(244,63,94,0.1); font-weight: 600; font-size: 12px; padding: 5px 12px;" onclick="deleteSelectedContent()">🗑️ Delete Selected</button>
                </div>
            </div>

            <div class="table-box">
                <table aria-label="Content Moderation Table">
                    <thead>
                        <tr>
                            <th scope="col" style="width: 36px; text-align: center;">
                                <input type="checkbox" id="select-all-checkbox" onchange="toggleSelectAll(this.checked)" title="Select All" style="cursor: pointer; width: 16px; height: 16px;" />
                            </th>
                            <th scope="col" style="width: 120px;">Media</th>
                            <th scope="col">Headline &amp; Content</th>
                            <th scope="col">Format &amp; District</th>
                            <th scope="col" style="min-width: 175px; cursor: pointer; user-select: none;" onclick="toggleSourceDateSort()" title="Click to cycle sorting: Date Newest ➔ Source (A-Z) ➔ Date Oldest ➔ Source (Z-A)">
                                Source &amp; Post Date <span id="sort-source-date-icon" style="color: #38bdf8; font-size: 11px; margin-left: 4px; padding: 2px 4px; background: rgba(56,189,248,0.15); border-radius: 4px;">🕒 ▼</span>
                            </th>
                            <th scope="col" style="width: 200px;">Actions</th>
                        </tr>
                    </thead>
                    <tbody id="mod-table-body">
                        <tr><td colspan="6" style="color: var(--text-muted);">Loading content for review...</td></tr>
                    </tbody>
                </table>
            </div>

            <!-- Table Pagination Bar -->
            <div id="table-pagination-bar" style="display: flex; justify-content: space-between; align-items: center; margin-top: 14px; padding: 12px 6px; border-top: 1px solid rgba(255,255,255,0.08); flex-wrap: wrap; gap: 10px;">
                <div style="font-size: 12px; color: var(--text-muted);" id="pagination-info">
                    Showing items...
                </div>
                <div style="display: flex; align-items: center; gap: 10px;">
                    <span style="font-size: 12px; color: var(--text-muted);">Rows per page:</span>
                    <select id="pagination-limit" onchange="changePageLimit(this.value)" style="background: #1e293b; border: 1px solid #334155; color: #fff; padding: 4px 8px; border-radius: 6px; font-size: 12px; cursor: pointer;">
                        <option value="10">10</option>
                        <option value="15" selected>15</option>
                        <option value="25">25</option>
                        <option value="50">50</option>
                    </select>
                    <div style="display: flex; gap: 4px; align-items: center; margin-left: 8px;">
                        <button id="page-first-btn" class="action-btn" onclick="goToSpecificPage(1)" style="font-size: 12px; padding: 4px 8px;" title="First Page" disabled>⏮ First</button>
                        <button id="page-prev-btn" class="action-btn" onclick="changePage(-1)" style="font-size: 12px; padding: 4px 10px;" disabled>&larr; Prev</button>
                        <div style="display: flex; align-items: center; gap: 4px; background: #0f172a; border: 1px solid #334155; border-radius: 6px; padding: 2px 6px;">
                            <span style="font-size: 12px; color: #94a3b8;">Page</span>
                            <input type="number" id="jump-to-page-input" min="1" value="1" onkeydown="if(event.key==='Enter') goToSpecificPage(this.value)" style="width: 48px; background: #1e293b; border: 1px solid #475569; color: #38bdf8; font-weight: 700; text-align: center; border-radius: 4px; padding: 2px 4px; font-size: 12px;" />
                            <span style="font-size: 12px; color: #94a3b8;">of <span id="page-total-count" style="color: #f8fafc; font-weight: 600;">1</span></span>
                            <button type="button" class="action-btn" onclick="goToSpecificPage(document.getElementById('jump-to-page-input').value)" style="font-size: 11px; padding: 2px 8px; background: #38bdf8; color: #0f172a; border: none; font-weight: 700; border-radius: 4px; margin-left: 2px; cursor: pointer;">Go &rarr;</button>
                        </div>
                        <button id="page-next-btn" class="action-btn" onclick="changePage(1)" style="font-size: 12px; padding: 4px 10px;" disabled>Next &rarr;</button>
                        <button id="page-last-btn" class="action-btn" onclick="goToSpecificPage(currentContentTotalPages)" style="font-size: 12px; padding: 4px 8px;" title="Last Page" disabled>Last ⏭</button>
                    </div>
                </div>
            </div>
        </section>

        <!-- Grievances Panel -->
        <section id="panel-grievances" role="tabpanel" aria-labelledby="tab-grievances" style="display: none;">
            <header>
                <div class="header-title">
                    <h1>IT Rules 2021 Grievance Redressal Desk</h1>
                    <p>Statutory 24h receipt acknowledgment and 15-day resolution SLA tracker</p>
                </div>
            </header>
            <div style="height: 20px;"></div>
            <div class="table-box">
                <table aria-label="Grievances Management Table">
                    <thead>
                        <tr>
                            <th scope="col">Ticket Reference</th>
                            <th scope="col">Complainant</th>
                            <th scope="col">Category</th>
                            <th scope="col">Target Content</th>
                            <th scope="col">Status</th>
                        </tr>
                    </thead>
                    <tbody id="grievance-table-body">
                        <tr><td colspan="5" style="color: var(--text-muted);">No open compliance grievances recorded in database.</td></tr>
                    </tbody>
                </table>
            </div>
        </section>

        <!-- Contributors Panel -->
        <section id="panel-contributors" role="tabpanel" aria-labelledby="tab-contributors" style="display: none;">
            <header>
                <div class="header-title">
                    <h1>Contributor Reputation &amp; Badges</h1>
                    <p>Real contributor levels, points, and Founding Contributor 2026 badges</p>
                </div>
            </header>
            <div style="height: 20px;"></div>
            <div class="table-box">
                <table aria-label="Contributor Reputation Table">
                    <thead>
                        <tr>
                            <th scope="col">User ID</th>
                            <th scope="col">Level</th>
                            <th scope="col">Trust Score</th>
                            <th scope="col">Points</th>
                            <th scope="col">Approved Posts</th>
                            <th scope="col">Founding Badge</th>
                        </tr>
                    </thead>
                    <tbody id="contrib-table-body">
                        <tr><td colspan="6" style="color: var(--text-muted);">Loading contributor profiles from database...</td></tr>
                    </tbody>
                </table>
            </div>
        </section>

        <!-- Compliance Audit Log Panel -->
        <section id="panel-audit" role="tabpanel" aria-labelledby="tab-audit" style="display: none;">
            <header>
                <div class="header-title">
                    <h1>Immutable Compliance Audit Trail</h1>
                    <p>Real audit trail records directly from PostgreSQL audit_logs table</p>
                </div>
            </header>
            <div style="height: 20px;"></div>
            <div class="table-box">
                <table aria-label="Compliance Audit Log Trail">
                    <thead>
                        <tr>
                            <th scope="col">Timestamp</th>
                            <th scope="col">Actor</th>
                            <th scope="col">Action</th>
                            <th scope="col">Target</th>
                            <th scope="col">Details</th>
                        </tr>
                    </thead>
                    <tbody id="audit-table-body">
                        <tr><td colspan="5" style="color: var(--text-muted);">No administrative events recorded yet.</td></tr>
                    </tbody>
                </table>
            </div>
        </section>
    </main>

    <!-- Create New Job Modal -->
    <div id="jobModal" class="modal">
        <div class="modal-content">
            <h2 style="margin-bottom: 16px;">Schedule New Background Task</h2>
            <form onsubmit="handleCreateJobSubmit(event)">
                <div class="form-group">
                    <label for="job-name">Job Identifier / Name:</label>
                    <input id="job-name" class="form-control" placeholder="e.g. daily_cache_cleanup" required />
                </div>
                <div class="form-group">
                    <label for="job-type">Task Type:</label>
                    <select id="job-type" class="form-control">
                        <option value="DEAD_LINK_CHECKER">Dead Video Link Scanner</option>
                        <option value="AUTO_MODERATION">Auto-Moderation Toxicity Evaluator</option>
                        <option value="GRIEVANCE_SLA">Grievance SLA Compliance Check</option>
                        <option value="REPUTATION_UPDATE">Contributor Reputation Recalculation</option>
                        <option value="CUSTOM">Custom System Maintenance Task</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="job-interval">Schedule Interval:</label>
                    <select id="job-interval" class="form-control">
                        <option value="1m">Every 1 Minute (@every 1m)</option>
                        <option value="5m" selected>Every 5 Minutes (@every 5m)</option>
                        <option value="15m">Every 15 Minutes (@every 15m)</option>
                        <option value="1h">Every 1 Hour (@every 1h)</option>
                        <option value="6h">Every 6 Hours (@every 6h)</option>
                        <option value="24h">Daily (@every 24h)</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="job-desc">Description:</label>
                    <input id="job-desc" class="form-control" placeholder="Optional description..." />
                </div>
                <div style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px;">
                    <button type="button" class="action-btn" onclick="closeNewJobModal()">Cancel</button>
                    <button type="submit" class="btn-primary">Create &amp; Schedule</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Manage Sources Modal -->
    <div id="sourcesModal" class="modal">
        <div class="modal-content" style="max-width: 600px;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
                <h2 style="margin: 0; font-size: 18px; color: #38bdf8;">🔗 Manage Scraping Sources for TN Live News Cron</h2>
                <button type="button" class="action-btn" onclick="closeSourcesModal()">✕</button>
            </div>
            <p style="color: var(--text-muted); font-size: 13px; margin-bottom: 16px;">
                Add regional RSS feeds, media websites, and video channels that <code>tn_live_news_cron</code> will crawl during each background cycle.
            </p>

            <div style="margin-bottom: 16px;">
                <div style="font-weight: 600; font-size: 12px; margin-bottom: 6px; color: #94a3b8;">ACTIVE SOURCE TARGETS:</div>
                <div id="sources-list" style="max-height: 180px; overflow-y: auto; background: rgba(0,0,0,0.2); border: 1px solid #3f3f4e; border-radius: 6px; padding: 8px;">
                    <div style="color: var(--text-muted); font-size: 12px; padding: 6px;">Loading active sources...</div>
                </div>
            </div>

            <div style="margin-bottom: 16px;">
                <div style="font-weight: 600; font-size: 12px; margin-bottom: 6px; color: #94a3b8;">ADD NEW SOURCE URL:</div>
                <div style="display: flex; gap: 8px;">
                    <input id="new-source-url" class="form-control" style="flex: 1;" placeholder="Enter RSS feed, website URL, or YouTube channel feed..." />
                    <button class="btn-primary" onclick="addSourceFromModal()">+ Add Source</button>
                </div>
            </div>

            <div style="margin-bottom: 16px;">
                <div style="font-weight: 600; font-size: 12px; margin-bottom: 6px; color: #94a3b8;">POPULAR TAMIL NADU FEEDS (CLICK TO ADD):</div>
                <div style="display: flex; gap: 6px; flex-wrap: wrap;">
                    <button class="action-btn" style="font-size: 11px;" onclick="quickAddSource('https://www.thehindu.com/news/national/tamil-nadu/feeder/default.rss')">+ The Hindu Tamil Nadu RSS</button>
                    <button class="action-btn" style="font-size: 11px;" onclick="quickAddSource('https://feeds.bbci.co.uk/tamil/rss.xml')">+ BBC News தமிழ் RSS</button>
                    <button class="action-btn" style="font-size: 11px;" onclick="quickAddSource('https://tamil.oneindia.com/rss/tamil-news-fb.xml')">+ OneIndia Tamil News RSS</button>
                    <button class="action-btn" style="font-size: 11px;" onclick="quickAddSource('https://news.google.com/rss/search?q=Tamil+Nadu&hl=ta&gl=IN&ceid=IN:ta')">+ Google News Tamil Nadu RSS</button>
                </div>
            </div>

            <div style="display: flex; justify-content: space-between; align-items: center; gap: 10px; margin-top: 20px;">
                <button type="button" class="action-btn" style="color: #38bdf8; border-color: rgba(56,189,248,0.4); font-size: 11px;" onclick="resetSourcesToDefault()">🔄 Reset to 4 Verified Feeds</button>
                <button type="button" class="btn-primary" onclick="closeSourcesModal()">Done</button>
            </div>
        </div>
    </div>

    <!-- Custom Confirmation & Rejection Modal Popup -->
    <div id="customConfirmModal" class="custom-modal-backdrop" onclick="handleConfirmBackdropClick(event)">
        <div class="custom-modal-card" onclick="event.stopPropagation()">
            <!-- Modal Header / Icon Header -->
            <div style="padding: 24px 24px 16px; display: flex; gap: 16px; align-items: flex-start;">
                <div id="customConfirmIconWrap" style="width: 46px; height: 46px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 22px; flex-shrink: 0;">
                </div>
                <div style="flex: 1;">
                    <h3 id="customConfirmTitle" style="margin: 0 0 6px; font-size: 17px; font-weight: 700; color: #f8fafc;">Confirmation</h3>
                    <p id="customConfirmMessage" style="margin: 0; font-size: 13px; color: #94a3b8; line-height: 1.5;"></p>
                    
                    <!-- Dynamic Rejection Options (Visible only for rejection) -->
                    <div id="customRejectDetails" style="display: none; margin-top: 14px; padding-top: 14px; border-top: 1px solid rgba(255,255,255,0.08);">
                        <div style="font-size: 11px; font-weight: 600; color: #e2e8f0; margin-bottom: 8px;">Select Rejection Reason (Optional):</div>
                        <div style="display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 10px;">
                            <span class="reason-chip" onclick="selectRejectReason(this, 'Duplicate / Repeated')">Duplicate</span>
                            <span class="reason-chip" onclick="selectRejectReason(this, 'Low Quality / Incomplete')">Low Quality</span>
                            <span class="reason-chip" onclick="selectRejectReason(this, 'Irrelevant / Off-topic')">Off-topic</span>
                            <span class="reason-chip" onclick="selectRejectReason(this, 'Misleading / Unverified')">Unverified</span>
                            <span class="reason-chip" onclick="selectRejectReason(this, 'Outdated Stale News')">Outdated</span>
                        </div>
                        <input type="text" id="customRejectNote" placeholder="Additional rejection notes (optional)..." style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 7px 12px; color: #f8fafc; font-size: 12px; outline: none;" />
                    </div>
                </div>
            </div>

            <!-- Modal Footer Buttons -->
            <div style="padding: 14px 24px; background: rgba(0, 0, 0, 0.3); border-top: 1px solid rgba(255, 255, 255, 0.06); display: flex; justify-content: flex-end; gap: 10px;">
                <button type="button" id="customConfirmCancelBtn" class="action-btn" style="padding: 7px 18px; font-size: 12px; border-radius: 8px; background: #1e293b; color: #cbd5e1; border: 1px solid #334155; cursor: pointer; font-weight: 600;" onclick="closeCustomConfirm(false)">Cancel</button>
                <button type="button" id="customConfirmActionBtn" style="padding: 7px 20px; font-size: 12px; border-radius: 8px; border: none; cursor: pointer; font-weight: 700; display: inline-flex; align-items: center; gap: 6px;" onclick="handleCustomConfirmAction()"></button>
            </div>
        </div>
    </div>

    <!-- Content Detail View Modal (Consumer Live Web & Mobile Preview) -->
    <div id="contentViewModal" class="modal">
        <div class="modal-content" style="max-width: 820px; max-height: 94vh; overflow-y: auto; background: #030712; border: 1px solid #1f2937; padding: 18px;">
            
            <!-- Top Controls Bar -->
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; padding-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.08); flex-wrap: wrap; gap: 8px;">
                <div style="display: flex; align-items: center; gap: 10px;">
                    <span style="font-weight: 700; font-size: 14px; color: #38bdf8;">👁️ Consumer Live Web Preview</span>
                    <div style="display: flex; background: rgba(0,0,0,0.4); border: 1px solid #374151; border-radius: 6px; padding: 2px;">
                        <button id="view-mode-web" class="action-btn" style="padding: 3px 10px; font-size: 11px; background: rgba(56,189,248,0.2); color: #38bdf8; border: none;" onclick="setViewDeviceMode('web')">💻 Web Reader</button>
                        <button id="view-mode-mobile" class="action-btn" style="padding: 3px 10px; font-size: 11px; background: transparent; color: var(--text-muted); border: none;" onclick="setViewDeviceMode('mobile')">📱 Mobile App</button>
                    </div>
                </div>
                <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap;">
                    <div style="display: flex; align-items: center; gap: 6px;">
                        <span style="font-size: 11px; color: var(--text-muted);">Region:</span>
                        <select id="modal-district-select" onchange="updateCurrentItemDistrict(this.value)" style="background: #1e293b; border: 1px solid #38bdf8; color: #fff; padding: 3px 8px; border-radius: 6px; font-size: 11px; cursor: pointer;">
                        </select>
                    </div>
                    <button id="modal-viral-toggle-btn" class="action-btn" onclick="toggleCurrentItemViral()" style="font-size: 11px; padding: 4px 10px; cursor: pointer;">⚡ Mark Viral</button>
                    <button class="action-btn" style="color: #a855f7; border-color: rgba(168,85,247,0.4); font-size: 11px; padding: 4px 10px; cursor: pointer;" onclick="openEditContentModal(activeViewItemId)" title="Edit title, text, image, or banner placement">✏️ Edit</button>
                    <button id="modal-banner-toggle-btn" class="action-btn" style="color: #eab308; border-color: rgba(234,179,8,0.4); font-size: 11px; padding: 4px 10px; cursor: pointer;" onclick="toggleCurrentItemMainBanner()">⭐ Set Hero</button>
                    <button id="view-modal-approve-btn" class="btn-approve" style="padding: 6px 14px; font-size: 12px;">✓ Approve</button>
                    <button id="view-modal-reject-btn" class="action-btn" style="color: #f59e0b; border-color: rgba(245,158,11,0.4); background: rgba(245,158,11,0.1); font-weight: 600; font-size: 12px; padding: 5px 12px;">🚫 Reject</button>
                    <button id="view-modal-delete-btn" class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.4); background: rgba(244,63,94,0.1); font-weight: 600; font-size: 12px; padding: 5px 12px;">🗑️ Delete</button>
                    <button type="button" class="action-btn" onclick="closeContentViewModal()" style="font-size: 14px; padding: 3px 8px;">✕</button>
                </div>
            </div>

            <!-- Consumer Simulator Container -->
            <div id="consumer-preview-container" style="display: flex; justify-content: center; padding: 12px 0; background: #030712; min-height: 480px;">
                
                <!-- The Consumer Device Frame -->
                <div id="consumer-device-frame" style="width: 100%; max-width: 680px; transition: all 0.25s ease; background: #0f172a; border: 1px solid #1e293b; border-radius: 12px; overflow: hidden; box-shadow: 0 25px 50px -12px rgba(0,0,0,0.7);">
                    
                    <!-- TN24 Consumer App Bar -->
                    <div style="display: flex; justify-content: space-between; align-items: center; padding: 12px 18px; background: #090d16; border-bottom: 1px solid #1e293b;">
                        <div style="display: flex; align-items: center; gap: 8px;">
                            <div style="background: linear-gradient(135deg, #ef4444, #f97316); color: #fff; font-weight: 900; font-size: 12px; padding: 2px 7px; border-radius: 4px; letter-spacing: 0.5px;">TN24</div>
                            <span id="consumer-district-badge" style="background: rgba(56,189,248,0.15); color: #38bdf8; font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 12px;">📍 Madurai</span>
                        </div>
                        <div style="display: flex; align-items: center; gap: 8px; color: #94a3b8; font-size: 12px;">
                            <span id="consumer-category-badge" style="background: rgba(148,163,184,0.15); color: #cbd5e1; font-size: 11px; padding: 2px 8px; border-radius: 4px;">News</span>
                            <span style="font-size: 11px; color: #64748b;" id="consumer-time-badge">Today</span>
                        </div>
                    </div>

                    <!-- Article Body in Consumer View -->
                    <div style="padding: 20px;">
                        
                        <!-- Badges Row -->
                        <div style="display: flex; gap: 8px; align-items: center; margin-bottom: 12px; flex-wrap: wrap;">
                            <span id="consumer-type-badge" class="badge-pill badge-info">TEXT STORY</span>
                            <span id="consumer-status-pill" class="badge-pill badge-warning">STAGED PENDING</span>
                            <span style="font-size: 11px; color: var(--text-muted); margin-left: auto;">Source: <strong id="consumer-source-name" style="color: #cbd5e1;">Direct RSS</strong></span>
                        </div>

                        <!-- Media Display (Photo / Video) -->
                        <div id="consumer-media-box" style="margin-bottom: 16px;"></div>

                        <!-- Headline -->
                        <h2 id="consumer-title" style="font-size: 20px; font-weight: 700; color: #f8fafc; line-height: 1.4; margin-bottom: 14px;"></h2>

                        <!-- Full Unabridged Story Body -->
                        <div id="consumer-body" style="font-size: 14px; color: #cbd5e1; line-height: 1.7; margin-bottom: 20px; white-space: pre-wrap; max-height: 480px; overflow-y: auto; padding-right: 8px;"></div>

                        <!-- Full Web Source Link -->
                        <div style="margin-bottom: 18px;">
                            <a id="consumer-source-link" href="#" target="_blank" style="display: inline-flex; align-items: center; gap: 6px; font-size: 12px; color: #38bdf8; text-decoration: none; background: rgba(56,189,248,0.1); padding: 6px 12px; border-radius: 6px; border: 1px solid rgba(56,189,248,0.3);">
                                🌐 Open Original Publisher Story &rarr;
                            </a>
                        </div>

                        <!-- Social Engagement Bar -->
                        <div style="display: flex; justify-content: space-between; align-items: center; border-top: 1px solid #1e293b; padding-top: 14px; color: #94a3b8; font-size: 13px;">
                            <div style="display: flex; gap: 16px;">
                                <span style="display: flex; align-items: center; gap: 4px; color: #38bdf8;">👍 <strong>148</strong></span>
                                <span style="display: flex; align-items: center; gap: 4px;">💬 <strong>24</strong> comments</span>
                            </div>
                            <div style="display: flex; gap: 14px;">
                                <span style="cursor: pointer;">↗️ Share on WhatsApp</span>
                                <span style="cursor: pointer;">🔖 Save</span>
                            </div>
                        </div>

                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Edit Content & Broadcast Placement Modal -->
    <div id="editContentModal" class="modal" style="display: none; position: fixed; inset: 0; background: rgba(3, 7, 18, 0.85); backdrop-filter: blur(8px); z-index: 10001; align-items: center; justify-content: center; padding: 20px;">
        <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: 16px; width: 100%; max-width: 760px; max-height: 90vh; display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7);">
            <!-- Modal Header -->
            <div style="padding: 16px 20px; border-bottom: 1px solid #1e293b; display: flex; align-items: center; justify-content: space-between; background: rgba(15, 23, 42, 0.9);">
                <div style="display: flex; align-items: center; gap: 10px;">
                    <div style="width: 32px; height: 32px; border-radius: 8px; background: rgba(168, 85, 247, 0.15); border: 1px solid rgba(168, 85, 247, 0.3); display: flex; align-items: center; justify-content: center; font-size: 16px;">✏️</div>
                    <div>
                        <h3 style="margin: 0; font-size: 16px; font-weight: 700; color: #f8fafc;">Edit Article &amp; Broadcast Placement</h3>
                        <p style="margin: 0; font-size: 12px; color: #94a3b8;">Updates reflect live on the TN24 portal immediately</p>
                    </div>
                </div>
                <button type="button" class="action-btn" onclick="closeEditContentModal()" style="font-size: 14px; padding: 4px 10px;">✕</button>
            </div>

            <!-- Modal Form Body -->
            <div style="padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 14px;">
                <input type="hidden" id="edit-post-id" />

                <!-- Source & Post Date Meta Bar -->
                <div id="edit-post-source-date-bar" style="background: rgba(15, 23, 42, 0.7); border: 1px solid #334155; border-radius: 8px; padding: 8px 12px; font-size: 12px; color: #94a3b8; display: flex; align-items: center; justify-content: space-between;">
                    <span id="edit-post-source-info">🔗 Source: Web</span>
                    <span id="edit-post-date-info">🕒 Published: Recent</span>
                </div>

                <!-- Title -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Article Headline / Title *</label>
                    <input id="edit-post-title" type="text" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 10px 12px; color: #f8fafc; font-size: 14px; font-weight: 600;" placeholder="Enter headline..." />
                </div>

                <!-- District, Category & Language Row -->
                <div style="display: grid; grid-template-columns: 1.2fr 1fr 1fr; gap: 12px;">
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Geographic Region / District</label>
                        <select id="edit-post-district" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 9px 12px; color: #f8fafc; font-size: 13px;">
                        </select>
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Editorial Category</label>
                        <select id="edit-post-category" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 9px 12px; color: #f8fafc; font-size: 13px;">
                            <option value="News">News</option>
                            <option value="Politics">Politics</option>
                            <option value="Civic">Civic</option>
                            <option value="Business">Business</option>
                            <option value="Technical">Technical</option>
                            <option value="Entertainment">Entertainment</option>
                            <option value="Crime">Crime</option>
                            <option value="Sports">Sports</option>
                        </select>
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Article Language</label>
                        <select id="edit-post-language" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 9px 12px; color: #f8fafc; font-size: 13px;">
                            <option value="ta">🇮🇳 தமிழ் (Tamil)</option>
                            <option value="en">🇬🇧 English</option>
                            <option value="ta-en">🔄 Tanglish (TA-EN)</option>
                            <option value="hi">🇮🇳 हिन्दी (Hindi)</option>
                            <option value="ml">🇮🇳 മലയാളம் (Malayalam)</option>
                            <option value="te">🇮🇳 తెలుగు (Telugu)</option>
                            <option value="kn">🇮🇳 ಕನ್ನಡ (Kannada)</option>
                        </select>
                    </div>
                </div>

                <!-- Thumbnail URL & Live Preview -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Feature Photo / Thumbnail URL</label>
                    <div style="display: flex; gap: 10px; align-items: center;">
                        <input id="edit-post-thumbnail" type="text" oninput="updateEditThumbnailPreview()" style="flex: 1; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 9px 12px; color: #f8fafc; font-size: 13px;" placeholder="https://..." />
                        <div style="width: 64px; height: 42px; border-radius: 6px; overflow: hidden; background: #000; border: 1px solid #334155; flex-shrink: 0;">
                            <img id="edit-post-thumb-preview" src="" onerror="this.src='/admin/api/maps/svg?district=Tamil%20Nadu'" style="width: 100%; height: 100%; object-fit: cover;" alt="Preview" />
                        </div>
                    </div>
                </div>

                <!-- Full Story Description / Body -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Full Story Text / Description</label>
                    <textarea id="edit-post-description" rows="7" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 10px 12px; color: #f8fafc; font-size: 13px; line-height: 1.5; resize: vertical;" placeholder="Article text..."></textarea>
                </div>

                <!-- Banner & Priority Placement Controls -->
                <div style="background: rgba(30, 41, 59, 0.5); border: 1px solid #334155; border-radius: 10px; padding: 14px; display: flex; flex-direction: column; gap: 10px;">
                    <div style="font-size: 13px; font-weight: 700; color: #38bdf8;">🌟 Portal Placement &amp; Showcase Config</div>
                    <div style="display: flex; gap: 20px; align-items: center; flex-wrap: wrap;">
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-size: 13px; color: #f8fafc;">
                            <input id="edit-post-is-viral" type="checkbox" style="width: 16px; height: 16px; cursor: pointer;" />
                            <span>🔥 Viral Priority</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-size: 13px; color: #eab308; font-weight: 600;">
                            <input id="edit-post-is-main-banner" type="checkbox" style="width: 16px; height: 16px; cursor: pointer;" />
                            <span>⭐ Set as Portal Main Hero Banner</span>
                        </label>
                    </div>
                    <div style="display: flex; align-items: center; gap: 12px; margin-top: 4px;">
                        <label style="font-size: 12px; font-weight: 600; color: #cbd5e1;">Promote in Ad Banner Slot:</label>
                        <select id="edit-post-ad-slot" style="background: #0f172a; border: 1px solid #38bdf8; color: #f8fafc; border-radius: 6px; padding: 4px 10px; font-size: 12px;">
                            <option value="none">None (Standard Feed Story)</option>
                            <option value="header">Header Leaderboard (728x90)</option>
                            <option value="infeed">In-Feed Mid Billboard (728x90)</option>
                            <option value="sidebar">Sidebar Rectangle (250x250)</option>
                            <option value="square">Left Square Slot (200x200)</option>
                        </select>
                    </div>
                </div>
            </div>

            <!-- Modal Footer Buttons -->
            <div style="padding: 14px 20px; border-top: 1px solid #1e293b; display: flex; justify-content: flex-end; gap: 10px; background: rgba(15, 23, 42, 0.9);">
                <button type="button" class="action-btn" onclick="closeEditContentModal()">Cancel</button>
                <button type="button" class="btn-primary" onclick="saveEditContentForm()" style="font-weight: 700;">💾 Save &amp; Reflect Live</button>
            </div>
        </div>
    </div>

    <!-- Manual Content Creation Modal -->
    <div id="manualContentModal" class="modal">
        <div class="modal-content" style="max-width: 680px; max-height: 92vh; overflow-y: auto; background: #0f172a; border: 1px solid #1e293b; border-radius: 12px; padding: 22px;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; border-bottom: 1px solid rgba(255,255,255,0.08); padding-bottom: 12px;">
                <div>
                    <h2 style="margin: 0; font-size: 18px; color: #f8fafc; font-weight: 700;">✍️ Add Content / News Manually</h2>
                    <p style="margin: 4px 0 0 0; font-size: 12px; color: var(--text-muted);">Create a news story, video link, photo feature, or civic event</p>
                </div>
                <button type="button" class="action-btn" onclick="closeManualContentModal()" style="font-size: 16px; padding: 4px 10px;">✕</button>
            </div>

            <form id="manual-content-form" onsubmit="event.preventDefault();" style="display: flex; flex-direction: column; gap: 14px;">
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">HEADLINE / TITLE *</label>
                    <input id="manual-title" class="form-control" placeholder="Enter bold headline in Tamil or English..." required style="width: 100%; box-sizing: border-box;" />
                </div>

                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">CONTENT FORMAT</label>
                        <select id="manual-content-type" class="form-control" style="width: 100%; box-sizing: border-box; background: #1e293b; color: #fff;">
                            <option value="TEXT_STORY">📝 News Article / Story</option>
                            <option value="VIDEO_LINK">🎥 Video Link (YouTube)</option>
                            <option value="PHOTO">📸 Photo Feature</option>
                            <option value="EVENT">📅 Civic Event / Festival</option>
                        </select>
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">CATEGORY</label>
                        <select id="manual-category" class="form-control" style="width: 100%; box-sizing: border-box; background: #1e293b; color: #fff;">
                            <option value="News">News</option>
                            <option value="Civic Issues">Civic Issues</option>
                            <option value="Events &amp; Announcements">Events &amp; Announcements</option>
                            <option value="Viral Videos">Viral Videos</option>
                            <option value="Government Schemes">Government Schemes</option>
                            <option value="Crime &amp; Safety">Crime &amp; Safety</option>
                        </select>
                    </div>
                </div>

                <div style="display: grid; grid-template-columns: 1.2fr 1fr 1fr; gap: 12px;">
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">DISTRICT</label>
                        <select id="manual-district" class="form-control" style="width: 100%; box-sizing: border-box; background: #1e293b; color: #fff;">
                            <!-- Populated dynamically with all 40 TN districts -->
                        </select>
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">LANGUAGE</label>
                        <select id="manual-language" class="form-control" style="width: 100%; box-sizing: border-box; background: #1e293b; color: #fff;">
                            <option value="ta">🇮🇳 தமிழ் (Tamil)</option>
                            <option value="en">🇬🇧 English</option>
                            <option value="ta-en">🔄 Tanglish (TA-EN)</option>
                            <option value="hi">🇮🇳 हिन्दी (Hindi)</option>
                            <option value="ml">🇮🇳 മലയാളം (Malayalam)</option>
                            <option value="te">🇮🇳 తెలుగు (Telugu)</option>
                            <option value="kn">🇮🇳 ಕನ್ನಡ (Kannada)</option>
                        </select>
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">SOURCE ATTRIBUTION</label>
                        <input id="manual-source-url" class="form-control" placeholder="e.g. Editorial Desk, Dinamalar" style="width: 100%; box-sizing: border-box;" />
                    </div>
                </div>

                <div>
                    <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">FULL STORY BODY / ARTICLE CONTENT *</label>
                    <textarea id="manual-description" class="form-control" rows="5" placeholder="Write or paste full article paragraphs here..." style="width: 100%; box-sizing: border-box; font-family: inherit; line-height: 1.5; resize: vertical;" required></textarea>
                </div>

                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">IMAGE / BANNER URL (OPTIONAL)</label>
                        <input id="manual-image-url" class="form-control" placeholder="https://example.com/photo.jpg (leave blank for district landmark)" style="width: 100%; box-sizing: border-box;" />
                    </div>
                    <div>
                        <label style="display: block; font-size: 12px; font-weight: 600; color: #94a3b8; margin-bottom: 6px;">VIDEO / YOUTUBE URL (OPTIONAL)</label>
                        <input id="manual-video-url" class="form-control" placeholder="https://www.youtube.com/watch?v=..." style="width: 100%; box-sizing: border-box;" />
                    </div>
                </div>

                <div style="display: flex; justify-content: flex-end; gap: 10px; border-top: 1px solid rgba(255,255,255,0.08); padding-top: 16px; margin-top: 6px;">
                    <button type="button" class="action-btn" onclick="closeManualContentModal()">Cancel</button>
                    <button type="button" class="btn-primary" onclick="submitManualContent(false)" style="background: rgba(56,189,248,0.2); color: #38bdf8; border-color: rgba(56,189,248,0.4);">
                        📥 Save as Staged (Pending)
                    </button>
                    <button type="button" class="btn-approve" onclick="submitManualContent(true)" style="padding: 7px 18px;">
                        ✓ Publish Live Immediately
                    </button>
                </div>
            </form>
        </div>
    <!-- Language & Localization Configuration Modal -->
    <div id="languageConfigModal" class="modal" style="display: none; position: fixed; inset: 0; background: rgba(3, 7, 18, 0.85); backdrop-filter: blur(8px); z-index: 10001; align-items: center; justify-content: center; padding: 20px;">
        <div style="background: #0f172a; border: 1px solid #1e293b; border-radius: 16px; width: 100%; max-width: 580px; max-height: 90vh; display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.7);">
            <!-- Modal Header -->
            <div style="padding: 16px 20px; border-bottom: 1px solid #1e293b; display: flex; align-items: center; justify-content: space-between; background: rgba(15, 23, 42, 0.9);">
                <div style="display: flex; align-items: center; gap: 10px;">
                    <div style="width: 32px; height: 32px; border-radius: 8px; background: rgba(56, 189, 248, 0.15); border: 1px solid rgba(56, 189, 248, 0.3); display: flex; align-items: center; justify-content: center; font-size: 16px;">🌐</div>
                    <div>
                        <h3 style="margin: 0; font-size: 16px; font-weight: 700; color: #f8fafc;">Language &amp; Localization Configuration</h3>
                        <p style="margin: 0; font-size: 12px; color: #94a3b8;">Manage supported languages, default locale, and scraper filters</p>
                    </div>
                </div>
                <button type="button" class="action-btn" onclick="closeLanguageConfigModal()" style="font-size: 14px; padding: 4px 10px;">✕</button>
            </div>

            <!-- Modal Body -->
            <div style="padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 16px;">
                <!-- Default Language -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">Default Consumer App &amp; Portal Language</label>
                    <select id="cfg-default-language" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 10px 12px; color: #f8fafc; font-size: 13px;">
                        <option value="ta">🇮🇳 தமிழ் (Tamil) — Primary Statewide Default</option>
                        <option value="en">🇬🇧 English</option>
                        <option value="ta-en">🔄 Tanglish (Tamil + English Hybrid)</option>
                        <option value="hi">🇮🇳 हिन्दी (Hindi)</option>
                        <option value="ml">🇮🇳 മലയാളം (Malayalam)</option>
                        <option value="te">🇮🇳 తెలుగు (Telugu)</option>
                        <option value="kn">🇮🇳 ಕನ್ನಡ (Kannada)</option>
                    </select>
                </div>

                <!-- Enabled Languages Checkboxes -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 8px;">Active / Supported Languages in System</label>
                    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px; background: rgba(30, 41, 59, 0.4); border: 1px solid #334155; border-radius: 8px; padding: 12px;">
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="ta" checked style="accent-color: #38bdf8;" />
                            <span>🇮🇳 Tamil (ta)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="en" checked style="accent-color: #38bdf8;" />
                            <span>🇬🇧 English (en)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="ta-en" checked style="accent-color: #38bdf8;" />
                            <span>🔄 Tanglish (ta-en)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="hi" style="accent-color: #38bdf8;" />
                            <span>🇮🇳 Hindi (hi)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="ml" style="accent-color: #38bdf8;" />
                            <span>🇮🇳 Malayalam (ml)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="te" style="accent-color: #38bdf8;" />
                            <span>🇮🇳 Telugu (te)</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: #f8fafc; cursor: pointer;">
                            <input type="checkbox" class="cfg-lang-cb" value="kn" style="accent-color: #38bdf8;" />
                            <span>🇮🇳 Kannada (kn)</span>
                        </label>
                    </div>
                </div>

                <!-- Scraper Ingestion Policy -->
                <div>
                    <label style="display: block; font-size: 12px; font-weight: 700; color: #cbd5e1; margin-bottom: 6px;">News Ingestion Scraper Language Policy</label>
                    <select id="cfg-scraper-policy" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #334155; border-radius: 8px; padding: 10px 12px; color: #f8fafc; font-size: 13px;">
                        <option value="ALL">Ingest All Content (Detect &amp; Tag Language Automatically)</option>
                        <option value="ONLY_TAMIL">Strictly Tamil Only (Discard non-Tamil articles)</option>
                        <option value="ONLY_TAMIL_AND_ENGLISH">Tamil and English Only</option>
                    </select>
                </div>

                <div style="font-size: 12px; color: #94a3b8; background: rgba(56, 189, 248, 0.08); border-left: 3px solid #38bdf8; padding: 8px 12px; border-radius: 4px;">
                    💡 Language filters immediately update the moderation toolbar and default filters across TN24.
                </div>
            </div>

            <!-- Modal Footer -->
            <div style="padding: 14px 20px; border-top: 1px solid #1e293b; display: flex; justify-content: flex-end; gap: 10px; background: rgba(15, 23, 42, 0.9);">
                <button type="button" class="action-btn" onclick="closeLanguageConfigModal()">Cancel</button>
                <button type="button" class="btn-primary" onclick="saveLanguageConfig()" style="font-weight: 700; background: linear-gradient(135deg, #38bdf8, #0284c7);">💾 Save Language Settings</button>
            </div>
        </div>
    </div>

    <div id="toast" role="alert"></div>

    <script>
        let currentAuthToken = localStorage.getItem('tnnow_gauth_token');
        let currentOperatorEmail = localStorage.getItem('tnnow_operator_email');

        function switchTab(tabId) {
            document.querySelectorAll('.nav-btn').forEach(btn => {
                btn.classList.remove('active');
                btn.setAttribute('aria-selected', 'false');
            });
            const activeBtn = document.getElementById('tab-' + tabId);
            if (activeBtn) {
                activeBtn.classList.add('active');
                activeBtn.setAttribute('aria-selected', 'true');
            }

            document.querySelectorAll('main section').forEach(panel => {
                panel.style.display = 'none';
            });
            const targetPanel = document.getElementById('panel-' + tabId);
            if (targetPanel) {
                targetPanel.style.display = 'block';
            }

            if (tabId === 'agent') {
                checkAgentAuthState();
            } else if (tabId === 'moderation') {
                fetchPendingContent();
                updateDiscardedCountBadge();
            } else if (tabId === 'cron') {
                fetchCronJobs();
                fetchCronLogs();
            } else if (tabId === 'audit') {
                fetchAuditLogs();
            } else if (tabId === 'contributors') {
                fetchContributors();
            }

            announce('Switched to ' + tabId);
        }

        function showToast(msg) {
            const toast = document.getElementById('toast');
            toast.innerText = msg;
            toast.style.display = 'block';
            announce(msg);
            setTimeout(() => { toast.style.display = 'none'; }, 3000);
        }

        let customConfirmResolve = null;

        function showCustomConfirm(options) {
            return new Promise((resolve) => {
                customConfirmResolve = resolve;
                const modal = document.getElementById('customConfirmModal');
                const titleEl = document.getElementById('customConfirmTitle');
                const msgEl = document.getElementById('customConfirmMessage');
                const iconWrap = document.getElementById('customConfirmIconWrap');
                const actionBtn = document.getElementById('customConfirmActionBtn');
                const cancelBtn = document.getElementById('customConfirmCancelBtn');
                const rejectBox = document.getElementById('customRejectDetails');

                if (!modal) {
                    resolve(window.confirm(options.message || 'Are you sure?'));
                    return;
                }

                titleEl.textContent = options.title || 'Confirmation';
                msgEl.textContent = options.message || 'Are you sure you want to proceed?';
                cancelBtn.textContent = options.cancelText || 'Cancel';
                actionBtn.textContent = options.confirmText || 'Confirm';

                const type = options.type || 'warning';
                if (type === 'reject') {
                    iconWrap.style.background = 'rgba(245, 158, 11, 0.15)';
                    iconWrap.style.color = '#f59e0b';
                    iconWrap.style.border = '1px solid rgba(245, 158, 11, 0.3)';
                    iconWrap.innerHTML = '🚫';
                    actionBtn.style.background = '#f59e0b';
                    actionBtn.style.color = '#000';
                    if (rejectBox) rejectBox.style.display = 'block';
                    document.querySelectorAll('.reason-chip').forEach(c => c.classList.remove('selected'));
                    const noteInput = document.getElementById('customRejectNote');
                    if (noteInput) noteInput.value = '';
                } else if (type === 'delete') {
                    iconWrap.style.background = 'rgba(244, 63, 94, 0.15)';
                    iconWrap.style.color = '#f43f5e';
                    iconWrap.style.border = '1px solid rgba(244, 63, 94, 0.3)';
                    iconWrap.innerHTML = '🗑️';
                    actionBtn.style.background = '#f43f5e';
                    actionBtn.style.color = '#fff';
                    if (rejectBox) rejectBox.style.display = 'none';
                } else if (type === 'approve') {
                    iconWrap.style.background = 'rgba(34, 197, 94, 0.15)';
                    iconWrap.style.color = '#22c55e';
                    iconWrap.style.border = '1px solid rgba(34, 197, 94, 0.3)';
                    iconWrap.innerHTML = '✓';
                    actionBtn.style.background = '#22c55e';
                    actionBtn.style.color = '#fff';
                    if (rejectBox) rejectBox.style.display = 'none';
                } else {
                    iconWrap.style.background = 'rgba(56, 189, 248, 0.15)';
                    iconWrap.style.color = '#38bdf8';
                    iconWrap.style.border = '1px solid rgba(56, 189, 248, 0.3)';
                    iconWrap.innerHTML = '⚡';
                    actionBtn.style.background = '#38bdf8';
                    actionBtn.style.color = '#000';
                    if (rejectBox) rejectBox.style.display = 'none';
                }

                modal.style.display = 'flex';
            });
        }

        function selectRejectReason(el, reason) {
            document.querySelectorAll('.reason-chip').forEach(c => c.classList.remove('selected'));
            el.classList.add('selected');
            const noteInput = document.getElementById('customRejectNote');
            if (noteInput) noteInput.value = reason;
        }

        function closeCustomConfirm(result) {
            const modal = document.getElementById('customConfirmModal');
            if (modal) modal.style.display = 'none';
            if (customConfirmResolve) {
                const res = customConfirmResolve;
                customConfirmResolve = null;
                res(result);
            }
        }

        function handleCustomConfirmAction() {
            closeCustomConfirm(true);
        }

        function handleConfirmBackdropClick(e) {
            if (e.target.id === 'customConfirmModal') {
                closeCustomConfirm(false);
            }
        }

        function announce(msg) {
            const live = document.getElementById('aria-live-status');
            if (live) live.innerText = msg;
        }

        // Genuine Google OAuth 2.0 Handler
        function initiateGenuineGoogleOAuth() {
            showToast('Redirecting to Google Identity Services...');
            window.location.href = '/admin/auth/google';
        }

        // Check for Google OAuth callback parameters in URL
        const urlParams = new URLSearchParams(window.location.search);
        const gtoken = urlParams.get('gtoken');
        const gemail = urlParams.get('gemail');
        if (gtoken && gemail) {
            currentAuthToken = gtoken;
            currentOperatorEmail = gemail;
            localStorage.setItem('tnnow_gauth_token', currentAuthToken);
            localStorage.setItem('tnnow_operator_email', currentOperatorEmail);
            window.history.replaceState({}, document.title, window.location.pathname);
            setTimeout(() => {
                switchTab('agent');
                showToast('✓ Cryptographically verified with Google Account: ' + gemail);
            }, 300);
        }

        function checkAgentAuthState() {
            if (currentAuthToken && currentOperatorEmail) {
                document.getElementById('agent-auth-gate').style.display = 'none';
                document.getElementById('agent-workspace').style.display = 'flex';
                document.getElementById('operator-badge-text').innerText = 'Verified Google Operator: ' + currentOperatorEmail;
            } else {
                document.getElementById('agent-auth-gate').style.display = 'block';
                document.getElementById('agent-workspace').style.display = 'none';
            }
        }

        function handleSignOut() {
            currentAuthToken = null;
            currentOperatorEmail = null;
            localStorage.removeItem('tnnow_gauth_token');
            localStorage.removeItem('tnnow_operator_email');
            showToast('Signed out from Google account');
            checkAgentAuthState();
        }

        function promptConfigureGeminiKey() {
            const key = prompt('Enter your Google Gemini API Key to connect Antigravity model execution:');
            if (!key || !key.trim()) return;
            fetch('/admin/api/auth/configure-gemini', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ apiKey: key.trim() })
            })
            .then(res => res.json())
            .then(res => {
                showToast(res.message || 'Antigravity Key configured successfully');
            })
            .catch(() => showToast('Failed to save Antigravity key'));
        }

        // AI Agent Chat & Prompts with Antigravity IDE Interactive Features
        function sendQuickPrompt(promptText) {
            document.getElementById('prompt-input').value = promptText;
            handleSendPrompt(new Event('submit'));
        }

        function copyCodeBlock(btn) {
            const pre = btn.closest('pre');
            if (!pre) return;
            const code = pre.querySelector('code').innerText;
            navigator.clipboard.writeText(code).then(() => {
                btn.innerText = '✓ Copied';
                setTimeout(() => btn.innerText = '📋 Copy', 2000);
            });
        }

        function handleSendPrompt(e) {
            if (e) e.preventDefault();
            const input = document.getElementById('prompt-input');
            let prompt = input.value.trim();
            if (!prompt) return;

            // Interactive Antigravity Slash Commands
            if (prompt === '/help') {
                appendChatBubble('user', prompt);
                appendAgentResponse({
                    reply: "### 🚀 Google Antigravity IDE Interactive Slash Commands:\n\n" +
                           "• <code>/status</code> &mdash; Run comprehensive database and telemetry diagnostics\n" +
                           "• <code>/districts</code> &mdash; Query all 40 registered districts from PostgreSQL\n" +
                           "• <code>/cron</code> &mdash; Audit all background workers and execution schedules\n" +
                           "• <code>/mod</code> &mdash; Triage pending moderation queue submissions\n" +
                           "• <code>/clear</code> &mdash; Clear current chat console history\n" +
                           "• <code>/help</code> &mdash; Show this commands guide\n\n" +
                           "You can also type any prompt in natural language, English, or Tamil.",
                    thoughtTrace: [
                        { stepNumber: 1, title: "Antigravity CLI Dispatch", detail: "Resolved interactive slash command locally." }
                    ]
                });
                input.value = '';
                return;
            }

            if (prompt === '/clear') {
                document.getElementById('chat-messages').innerHTML = '';
                input.value = '';
                return;
            }

            if (prompt === '/status') prompt = 'Run full telemetry and database diagnostics';
            else if (prompt === '/districts') prompt = 'list all administrative districts';
            else if (prompt === '/cron') prompt = 'Show all scheduled cron background jobs';
            else if (prompt === '/mod') prompt = 'Triage pending submissions with Tamil NLP toxicity check';

            appendChatBubble('user', input.value.trim());
            input.value = '';

            // Interactive Typing Indicator
            const container = document.getElementById('chat-messages');
            const indicator = document.createElement('div');
            indicator.id = 'typing-indicator';
            indicator.className = 'chat-bubble agent';
            indicator.innerHTML = '<div style="display:flex;align-items:center;gap:10px;color:#c084fc;font-weight:600;font-size:13px;">' +
                '<span style="display:inline-block;width:8px;height:8px;border-radius:50%;background:#38bdf8;box-shadow:0 0 10px #38bdf8;animation:pulse 1s infinite alternate;"></span>' +
                'Google Antigravity IDE is reasoning...</div>';
            container.appendChild(indicator);
            container.scrollTop = container.scrollHeight;

            fetch('/admin/api/agent/prompt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + currentAuthToken
                },
                body: JSON.stringify({ prompt: prompt })
            })
            .then(res => res.json())
            .then(res => {
                const ind = document.getElementById('typing-indicator');
                if (ind) ind.remove();
                if (res.success && res.data) {
                    appendAgentResponse(res.data);
                } else {
                    appendChatBubble('agent', 'Error: ' + (res.message || 'Could not process prompt'));
                }
            })
            .catch(() => {
                const ind = document.getElementById('typing-indicator');
                if (ind) ind.remove();
                appendChatBubble('agent', 'Error connecting to Antigravity reasoning engine.');
            });
        }

        function appendChatBubble(sender, text) {
            const container = document.getElementById('chat-messages');
            const bubble = document.createElement('div');
            bubble.className = 'chat-bubble ' + sender;
            bubble.innerHTML = (sender === 'agent' ? '<strong>🚀 Google Antigravity IDE:</strong><br>' : '<strong>👤 You:</strong><br>') + text;
            container.appendChild(bubble);
            container.scrollTop = container.scrollHeight;
        }

        function formatMarkdown(text) {
            if (!text) return '';
            const t = String.fromCharCode(96);
            let html = text
                .replace(/&/g, '&amp;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;');
            // Code blocks with copy button
            html = html.replace(new RegExp(t+t+t+'([a-zA-Z0-9_-]*)\\n([\\s\\S]*?)'+t+t+t, 'g'), 
                '<div style="position:relative;margin:10px 0;"><button onclick="copyCodeBlock(this)" style="position:absolute;top:8px;right:8px;background:#272730;color:#38bdf8;border:1px solid #3f3f4e;padding:4px 8px;border-radius:4px;font-size:11px;font-weight:700;cursor:pointer;">📋 Copy</button><pre style="background:#141419;border:1px solid #272730;padding:14px;border-radius:8px;overflow-x:auto;font-family:monospace;font-size:13px;color:#e2e8f0;margin:0;"><code>$2</code></pre></div>');
            // Inline code
            html = html.replace(new RegExp(t+'([^'+t+']+)'+t, 'g'), '<code style="background:#272730; padding:2px 6px; border-radius:4px; font-size:12px; color:#38bdf8; font-family:monospace;">$1</code>');
            // Headers
            html = html.replace(/^### (.*$)/gm, '<h3 style="color:#fff;font-size:16px;margin:12px 0 6px 0;font-weight:700;">$1</h3>');
            html = html.replace(/^## (.*$)/gm, '<h2 style="color:#fff;font-size:17px;margin:14px 0 8px 0;font-weight:800;">$1</h2>');
            // Bold
            html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
            // Italic
            html = html.replace(/\*([^*]+)\*/g, '<em>$1</em>');
            // Bullet points
            html = html.replace(/^\s*•\s+(.*)$/gm, '<li style="margin-left: 18px; margin-bottom: 4px;">$1</li>');
            html = html.replace(/^\s*\*\s+(.*)$/gm, '<li style="margin-left: 18px; margin-bottom: 4px;">$1</li>');
            // Line breaks
            html = html.replace(/\n/g, '<br>');
            return html;
        }

        function appendAgentResponse(data) {
            const container = document.getElementById('chat-messages');
            const bubble = document.createElement('div');
            bubble.className = 'chat-bubble agent';
            
            let html = '<div style="display: flex; align-items: center; gap: 8px; margin-bottom: 10px;">' +
                '<strong>🚀 Google Antigravity IDE Agent</strong>' +
                '<span style="font-size: 10px; background: rgba(56, 189, 248, 0.15); color: #38bdf8; border: 1px solid rgba(56, 189, 248, 0.3); padding: 2px 6px; border-radius: 6px; font-weight: 700;">DeepMind Continuous Inference</span>' +
                '</div>';

            // Collapsible Antigravity Thought Trace
            if (data.thoughtTrace && data.thoughtTrace.length > 0) {
                html += '<details open style="background: rgba(168, 85, 247, 0.08); border: 1px solid rgba(168, 85, 247, 0.25); border-radius: 8px; padding: 10px 14px; margin-bottom: 14px; font-size: 12px; cursor: pointer;">' +
                    '<summary style="font-weight: 700; color: #c084fc; margin-bottom: 6px; outline: none; user-select: none;">🧠 ANTIGRAVITY REASONING TRACE (' + data.thoughtTrace.length + ' steps) <span style="font-size: 10px; opacity: 0.7; font-weight: normal;">(click to collapse)</span></summary>';
                data.thoughtTrace.forEach(step => {
                    html += '<div style="margin-bottom: 4px; padding-left: 8px;">• <strong>Step ' + step.stepNumber + ': ' + step.title + '</strong> &mdash; <span style="color: var(--text-muted);">' + step.detail + '</span></div>';
                });
                html += '</details>';
            }

            html += '<div style="line-height: 1.65; color: #f1f5f9;">' + formatMarkdown(data.reply) + '</div>';

            if (data.actionCards && data.actionCards.length > 0) {
                window.pendingActionCards = window.pendingActionCards || {};
                data.actionCards.forEach((card, idx) => {
                    const cardKey = card.id || ('card_' + Date.now() + '_' + idx);
                    window.pendingActionCards[cardKey] = card;
                    const label = (card.buttonLabel || 'Execute Action').replace(/^⚡\s*/, '');
                    html += '<div class="action-card">' +
                        '<div>' +
                            '<strong>' + card.title + '</strong><br>' +
                            '<span style="font-size: 12px; color: var(--text-muted);">' + card.description + '</span>' +
                        '</div>' +
                        '<button class="btn-primary" onclick="dispatchCardAction(\'' + cardKey + '\')">⚡ ' + label + '</button>' +
                    '</div>';
                });
            }

            bubble.innerHTML = html;
            container.appendChild(bubble);
            container.scrollTop = container.scrollHeight;
        }

        function dispatchCardAction(cardKey) {
            const card = (window.pendingActionCards || {})[cardKey];
            if (!card) {
                showToast('Action expired or unavailable');
                return;
            }
            executeAgentAction(card.actionType, card.payload);
        }

        function executeAgentAction(type, payload) {
            showToast('⚡ Antigravity executing: ' + type + '...');
            fetch('/admin/api/agent/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + currentAuthToken
                },
                body: JSON.stringify({ actionType: type, payload: payload })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success && res.data) {
                    showToast('✓ ' + res.data.result);
                    let chatDetail = '✅ <strong>Antigravity Action Executed on Server</strong>:<br>' +
                        '• <strong>Result:</strong> ' + res.data.result + '<br>' +
                        '• <strong>Server Duration:</strong> ' + res.data.durationMs + 'ms<br>' +
                        '• <strong>Audit Trail:</strong> Saved in PostgreSQL <code>audit_logs</code>';

                    if (type === 'CREATE_CRON') {
                        fetch('/admin/api/cron')
                            .then(cres => cres.json())
                            .then(cres => {
                                if (cres && cres.data && cres.data.length > 0) {
                                    let cronListHtml = '<div style="margin-top:12px;background:rgba(255,255,255,0.04);border:1px solid #3f3f4e;border-radius:8px;padding:12px;">' +
                                        '<div style="font-weight:700;color:#38bdf8;margin-bottom:8px;">📋 Active Background Cron Workers (' + cres.data.length + ' registered):</div>' +
                                        '<ul style="margin:0;padding-left:18px;font-size:12px;line-height:1.7;">';
                                    cres.data.forEach(j => {
                                        const isNew = (res.data.result && res.data.result.includes(j.name));
                                        cronListHtml += '<li><strong>' + j.name + '</strong> (<code>' + j.id + '</code>) &mdash; <span class="badge-pill badge-info">@every ' + j.scheduleInterval + '</span> ' + (isNew ? '<span class="badge-pill badge-success" style="font-size:10px;">NEW</span>' : '') + '</li>';
                                    });
                                    cronListHtml += '</ul>' +
                                        '<div style="margin-top:10px;"><button class="btn-primary" onclick="showSection(\'cron\')" style="padding:4px 10px;font-size:11px;">⏱️ Open Scheduled Cron Table &rarr;</button></div>' +
                                        '</div>';
                                    appendChatBubble('agent', chatDetail + cronListHtml);
                                } else {
                                    appendChatBubble('agent', chatDetail);
                                }
                            })
                            .catch(() => appendChatBubble('agent', chatDetail));
                    } else if (type === 'CREATE_CONTENT') {
                        chatDetail += '<div style="margin-top:10px;"><button class="btn-primary" onclick="switchTab(\'overview\')" style="padding:4px 10px;font-size:11px;">📊 View Live Feed & Metrics &rarr;</button></div>';
                        appendChatBubble('agent', chatDetail);
                    } else if (type === 'SCRAPE_URL') {
                        chatDetail += '<div style="margin-top:10px;"><button class="btn-primary" onclick="switchTab(\'moderation\')" style="padding:4px 10px;font-size:11px;">📰 Open Staging Review Queue &rarr;</button></div>';
                        appendChatBubble('agent', chatDetail);
                    } else {
                        appendChatBubble('agent', chatDetail);
                    }

                    fetchStats();
                    fetchCronJobs();
                } else {
                    showToast('✗ ' + (res.message || 'Execution failed'));
                }
            })
            .catch(() => showToast('✗ Failed to dispatch Antigravity action'));
        }

        // Live Database Stats Fetch
        function fetchStats() {
            fetch('/admin/api/stats')
                .then(res => res.json())
                .then(res => {
                    if (!res || !res.data) return;
                    const d = res.data;
                    const totElem = document.getElementById('kpi-total');
                    if (totElem && d.totalContent !== undefined) totElem.innerText = d.totalContent;
                    const pendElem = document.getElementById('kpi-pending');
                    if (pendElem && d.pendingModeration !== undefined) pendElem.innerText = d.pendingModeration;
                    const quarElem = document.getElementById('kpi-quarantine');
                    if (quarElem && d.quarantinedContent !== undefined) quarElem.innerText = d.quarantinedContent;
                    
                    if (d.grievanceStats) {
                        const sla = d.grievanceStats.slaComplianceRate !== undefined ? d.grievanceStats.slaComplianceRate : (d.grievanceStats.SLAComplianceRate || 100);
                        const slaElem = document.getElementById('kpi-sla');
                        if (slaElem) slaElem.innerText = Number(sla).toFixed(1) + '%';
                    }
                    if (d.infrastructure && d.infrastructure.postgresStatus) {
                        const infraBadge = document.getElementById('live-infra-badge');
                        if (infraBadge) infraBadge.innerText = d.infrastructure.postgresStatus;
                    }

                    // Render Real District Bar Chart
                    const chart = document.getElementById('district-chart');
                    if (chart) {
                        if (d.districtMetrics && d.districtMetrics.length > 0) {
                            let html = '';
                            d.districtMetrics.forEach(m => {
                                html += '<div class="bar-chart-row">' +
                                    '<span class="bar-chart-label">' + m.districtName + '</span>' +
                                    '<div class="bar-chart-track" role="progressbar" aria-valuenow="' + m.percentage + '" aria-valuemin="0" aria-valuemax="100">' +
                                        '<div class="bar-chart-fill" style="width: ' + Math.max(m.percentage, 4) + '%;"></div>' +
                                    '</div>' +
                                    '<span class="bar-chart-value">' + m.percentage + '% (' + m.count + ')</span>' +
                                '</div>';
                            });
                            chart.innerHTML = html;
                        } else {
                            chart.innerHTML = '<p style="color: var(--text-muted); font-size: 13px; padding: 12px 0;">No district content registered yet.</p>';
                        }
                    }

                    // Render Real Format Breakdown
                    const fbox = document.getElementById('format-breakdown');
                    if (fbox && d.formatBreakdown) {
                        const fb = d.formatBreakdown;
                        fbox.innerHTML = 
                            renderFormatRow('Video Links', '#ef4444', fb['VIDEO_LINK'] || 0) +
                            renderFormatRow('Photo Posts', '#38bdf8', fb['PHOTO'] || 0) +
                            renderFormatRow('Text Stories', '#10b981', fb['TEXT_STORY'] || 0) +
                            renderFormatRow('Local Events', '#f59e0b', fb['EVENT'] || 0);
                    }
                })
                .catch(() => {});
        }

        function renderFormatRow(name, color, count) {
            return '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                '<span style="display: flex; align-items: center; gap: 8px;">' +
                    '<span style="width: 12px; height: 12px; border-radius: 3px; background: ' + color + ';"></span>' +
                    '<strong>' + name + '</strong>' +
                '</span>' +
                '<span class="badge-pill badge-info">' + count + '</span>' +
            '</div>';
        }

        // Cron API Functions
        function fetchCronJobs() {
            fetch('/admin/api/cron')
                .then(res => res.json())
                .then(res => {
                    const tbody = document.getElementById('cron-table-body');
                    if (!res || !res.data || res.data.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="7" style="color: var(--text-muted);">No cron jobs active.</td></tr>';
                        return;
                    }
                    let rows = '';
                    res.data.forEach(job => {
                        const statusBadge = job.isActive
                            ? '<span class="badge-pill badge-success">ACTIVE</span>'
                            : '<span class="badge-pill badge-warning">PAUSED</span>';
                        const toggleLabel = job.isActive ? 'Pause' : 'Resume';
                        const lastRun = job.lastRunAt ? new Date(job.lastRunAt).toLocaleTimeString() : 'Never';

                        rows += '<tr>' +
                            '<td><strong>' + job.name + '</strong><br><code style="font-size: 11px; color: var(--text-muted);">' + job.id + '</code></td>' +
                            '<td><span class="badge-pill badge-info">@every ' + job.scheduleInterval + '</span></td>' +
                            '<td>' + job.jobType + '</td>' +
                            '<td>' + statusBadge + '</td>' +
                            '<td>' + lastRun + '</td>' +
                            '<td>' + job.runCount + ' / ' + job.failureCount + '</td>' +
                            '<td style="white-space: nowrap;">' +
                                '<button class="btn-approve" onclick="triggerCronJob(\'' + job.id + '\')">▶ Run</button> ' +
                                '<button class="action-btn" onclick="toggleCronJob(\'' + job.id + '\')">' + toggleLabel + '</button> ' +
                                (job.id === 'tn_live_news_cron' ? '<button class="action-btn" style="color: #38bdf8; border-color: rgba(56,189,248,0.4);" onclick="openSourcesModal(\'' + job.id + '\')">🔗 Sources (' + (job.sourceUrls ? job.sourceUrls.length : 1) + ')</button> ' : '') +
                                '<button class="action-btn" onclick="promptEditCronJob(\'' + job.id + '\', \'' + job.name.replace(/'/g, "\\'") + '\', \'' + job.scheduleInterval + '\')">✏️ Edit</button> ' +
                                '<button class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.3);" onclick="confirmDeleteCronJob(\'' + job.id + '\')">🗑️ Delete</button>' +
                            '</td>' +
                        '</tr>';
                    });
                    tbody.innerHTML = rows;
                })
                .catch(() => {});
        }

        function promptEditCronJob(jobId, currentName, currentInterval) {
            const newInterval = prompt('Edit Schedule Interval for "' + jobId + '" (e.g. 15m, 30m, 1h, 2h, 24h):', currentInterval);
            if (!newInterval) return;
            const newName = prompt('Edit Name for "' + jobId + '":', currentName) || currentName;

            showToast('Updating ' + jobId + '...');
            fetch('/admin/api/cron/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: jobId, name: newName, interval: newInterval })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Updated ' + jobId);
                    fetchCronJobs();
                } else {
                    showToast('✗ ' + (res.message || 'Update failed'));
                }
            })
            .catch(() => showToast('✗ Failed to update cron job'));
        }

        async function confirmDeleteCronJob(jobId) {
            const confirmed = await showCustomConfirm({
                title: 'Delete Background Worker',
                message: 'Are you sure you want to permanently delete background worker "' + jobId + '"?',
                type: 'delete',
                confirmText: '🗑️ Delete Worker'
            });
            if (!confirmed) return;
            showToast('Deleting ' + jobId + '...');
            fetch('/admin/api/cron/delete', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: jobId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Deleted ' + jobId);
                    fetchCronJobs();
                } else {
                    showToast('✗ ' + (res.message || 'Delete failed'));
                }
            })
            .catch(() => showToast('✗ Failed to delete cron job'));
        }

        function triggerCronJob(jobId) {
            showToast('⏳ Triggering ' + jobId + '...');
            fetch('/admin/api/cron/trigger', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: jobId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    const msg = (res.data && res.data.message) ? res.data.message : ('✓ ' + jobId + ' triggered in background');
                    showToast(msg);
                    fetchCronJobs();
                    fetchCronLogs();
                    // Poll logs and moderation queue every 3 seconds for 21 seconds to display live progress
                    let pollCount = 0;
                    const pollInterval = setInterval(() => {
                        pollCount++;
                        fetchCronLogs();
                        fetchCronJobs();
                        if (typeof fetchPendingContent === 'function') fetchPendingContent();
                        if (typeof fetchStats === 'function') fetchStats();
                        if (pollCount >= 7) clearInterval(pollInterval);
                    }, 3000);
                } else {
                    showToast('✗ ' + (res.message || 'Execution error'));
                }
            })
            .catch(() => showToast('✗ Failed to trigger job (network error)'));
        }

        function toggleCronJob(jobId) {
            fetch('/admin/api/cron/toggle', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: jobId })
            })
            .then(res => res.json())
            .then(res => {
                showToast('State updated for ' + jobId);
                fetchCronJobs();
            })
            .catch(() => showToast('Failed to update state'));
        }

        function fetchCronLogs() {
            fetch('/admin/api/cron/logs')
                .then(res => res.json())
                .then(res => {
                    const box = document.getElementById('cron-console');
                    if (!res || !res.data || res.data.length === 0) {
                        box.innerText = 'No execution logs recorded yet. Trigger a job above to view output.';
                        return;
                    }
                    let text = '';
                    res.data.forEach(l => {
                        text += '[' + new Date(l.executedAt).toLocaleTimeString() + '] [' + l.status + '] (' + l.durationMs + 'ms) ' + l.jobId + ': ' + l.message + '\n';
                    });
                    box.innerText = text;
                })
                .catch(() => {});
        }

        let activeSourcesJobId = 'tn_live_news_cron';

        function openSourcesModal(jobId) {
            activeSourcesJobId = jobId || 'tn_live_news_cron';
            document.getElementById('sourcesModal').style.display = 'flex';
            fetchSourcesList();
        }

        function closeSourcesModal() {
            document.getElementById('sourcesModal').style.display = 'none';
        }

        function fetchSourcesList() {
            fetch('/admin/api/cron/sources?jobId=' + activeSourcesJobId)
                .then(res => res.json())
                .then(res => {
                    const list = document.getElementById('sources-list');
                    if (!res || !res.data || res.data.length === 0) {
                        list.innerHTML = '<div style="color: var(--text-muted); font-size: 12px; padding: 6px;">No custom sources configured. Default The Hindu Tamil Nadu feed is active.</div>';
                        return;
                    }
                    let html = '';
                    res.data.forEach(url => {
                        const enc = encodeURIComponent(url);
                        html += '<div style="display: flex; justify-content: space-between; align-items: center; padding: 8px 10px; border-bottom: 1px solid rgba(255,255,255,0.06); gap: 8px;">' +
                            '<span style="font-size: 12px; color: #f8fafc; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1;" title="' + url + '">🔗 ' + url + '</span>' +
                            '<div style="display: flex; gap: 4px; flex-shrink: 0;">' +
                                '<button class="action-btn" style="color: #38bdf8; padding: 2px 8px; font-size: 11px;" onclick="modifySourceInModal(\'' + enc + '\')">✏️ Modify</button>' +
                                '<button class="action-btn" style="color: #10b981; padding: 2px 8px; font-size: 11px;" onclick="testScrapeSource(\'' + enc + '\')">⚡ Scrape Now</button>' +
                                '<button class="action-btn" style="color: #f43f5e; padding: 2px 8px; font-size: 11px;" onclick="removeSourceFromModal(\'' + enc + '\')">🗑️ Remove</button>' +
                            '</div>' +
                        '</div>';
                    });
                    list.innerHTML = html;
                })
                .catch(() => {});
        }

        function modifySourceInModal(encodedOldUrl) {
            const oldUrl = decodeURIComponent(encodedOldUrl);
            const newUrl = prompt('Modify source URL for ' + activeSourcesJobId + ':', oldUrl);
            if (!newUrl || newUrl.trim() === '' || newUrl === oldUrl) return;
            showToast('Updating source...');
            fetch('/admin/api/cron/sources/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: activeSourcesJobId, oldUrl: oldUrl, newUrl: newUrl.trim() })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    fetchSourcesList();
                    fetchCronJobs();
                } else {
                    showToast('✗ ' + (res.message || 'Update failed'));
                }
            })
            .catch(() => showToast('✗ Failed to update source'));
        }

        function testScrapeSource(encodedUrl) {
            const url = decodeURIComponent(encodedUrl);
            closeSourcesModal();
            switchTab('moderation');
            const input = document.getElementById('scrape-url-input');
            if (input) input.value = url;
            triggerSiteScrape();
        }

        function addSourceFromModal() {
            const input = document.getElementById('new-source-url');
            const url = input.value.trim();
            if (!url) return;
            showToast('Adding source...');
            fetch('/admin/api/cron/sources/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: activeSourcesJobId, url: url })
            })
            .then(res => res.json())
            .then(res => {
                input.value = '';
                showToast('✓ ' + res.message);
                fetchSourcesList();
                fetchCronJobs();
            })
            .catch(() => showToast('✗ Failed to add source'));
        }

        function quickAddSource(url) {
            document.getElementById('new-source-url').value = url;
            addSourceFromModal();
        }

        function resetSourcesToDefault() {
            if (!confirm('Reset sources for ' + activeSourcesJobId + ' to the 4 verified regional feeds (The Hindu, BBC Tamil, OneIndia, Google News TN)?')) return;
            showToast('Resetting sources...');
            fetch('/admin/api/cron/sources/reset', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: activeSourcesJobId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    fetchSourcesList();
                    fetchCronJobs();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to reset sources'));
                }
            })
            .catch(() => showToast('✗ Failed to reset sources'));
        }

        function removeSourceFromModal(encodedUrl) {
            const url = decodeURIComponent(encodedUrl);
            showToast('Removing source...');
            fetch('/admin/api/cron/sources/remove', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ jobId: activeSourcesJobId, url: url })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + res.message);
                fetchSourcesList();
                fetchCronJobs();
            })
            .catch(() => showToast('✗ Failed to remove source'));
        }

        function openNewJobModal() {
            document.getElementById('jobModal').style.display = 'flex';
        }
        function closeNewJobModal() {
            document.getElementById('jobModal').style.display = 'none';
        }

        function handleCreateJobSubmit(e) {
            e.preventDefault();
            const name = document.getElementById('job-name').value;
            const type = document.getElementById('job-type').value;
            const interval = document.getElementById('job-interval').value;
            const desc = document.getElementById('job-desc').value;

            fetch('/admin/api/cron/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    id: name.toLowerCase().replace(/[^a-z0-9_]/g, '_'),
                    name: name,
                    jobType: type,
                    scheduleInterval: interval,
                    description: desc
                })
            })
            .then(res => res.json())
            .then(res => {
                closeNewJobModal();
                showToast('✓ New scheduled cron job created');
                fetchCronJobs();
            })
            .catch(() => showToast('Failed to create job'));
        }

        function fetchAuditLogs() {
            fetch('/admin/api/audit')
                .then(res => res.json())
                .then(res => {
                    const tbody = document.getElementById('audit-table-body');
                    if (!res || !res.data || res.data.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="5" style="color: var(--text-muted);">No audit logs recorded in database.</td></tr>';
                        return;
                    }
                    let rows = '';
                    res.data.forEach(a => {
                        rows += '<tr>' +
                            '<td>' + new Date(a.createdAt).toLocaleString() + '</td>' +
                            '<td><code>' + a.actorId + '</code></td>' +
                            '<td><span class="badge-pill badge-info">' + a.action + '</span></td>' +
                            '<td><code>' + a.targetEntity + ':' + a.targetId + '</code></td>' +
                            '<td>' + a.details + '</td>' +
                        '</tr>';
                    });
                    tbody.innerHTML = rows;
                })
                .catch(() => {});
        }

        function fetchContributors() {
            fetch('/admin/api/contributors')
                .then(res => res.json())
                .then(res => {
                    const tbody = document.getElementById('contrib-table-body');
                    if (!res || !res.data || res.data.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="6" style="color: var(--text-muted);">No contributors recorded yet.</td></tr>';
                        return;
                    }
                    let rows = '';
                    res.data.forEach(c => {
                        const badge = c.isFounding
                            ? '<span class="badge-pill badge-warning">★ Founding 2026</span>'
                            : '<span style="color: var(--text-muted); font-size: 12px;">Standard</span>';
                        rows += '<tr>' +
                            '<td><code>' + c.userId + '</code></td>' +
                            '<td><span class="badge-pill badge-info">' + c.level + '</span></td>' +
                            '<td><strong style="color: var(--color-success);">' + c.trustScore + ' / 100</strong></td>' +
                            '<td>' + c.points + '</td>' +
                            '<td>' + c.approvedCount + '</td>' +
                            '<td>' + badge + '</td>' +
                        '</tr>';
                    });
                    tbody.innerHTML = rows;
                })
                .catch(() => {});
        }

        // Staged News & Moderation Queue Functions
        let activeContentFilter = 'PENDING';
        let activeContentSearchQuery = '';
        let cachedContentMap = {};
        let currentContentPage = 1;
        let currentContentLimit = 15;
        let currentContentTotalPages = 1;
        let currentContentTotalCount = 0;

        function handleContentSearchKey(event) {
            const input = document.getElementById('content-search-input');
            const clearBtn = document.getElementById('content-search-clear-btn');
            if (input && clearBtn) {
                clearBtn.style.display = input.value.trim() ? 'inline-block' : 'none';
            }
            if (event.key === 'Enter') {
                triggerContentSearch();
            }
        }

        function triggerContentSearch() {
            const input = document.getElementById('content-search-input');
            activeContentSearchQuery = input ? input.value.trim() : '';
            currentContentPage = 1;

            const clearBtn = document.getElementById('content-search-clear-btn');
            if (clearBtn) clearBtn.style.display = activeContentSearchQuery ? 'inline-block' : 'none';

            const badgeBar = document.getElementById('active-search-badge-bar');
            const queryText = document.getElementById('active-search-query-text');
            if (activeContentSearchQuery) {
                if (badgeBar) badgeBar.style.display = 'flex';
                if (queryText) queryText.innerText = '"' + activeContentSearchQuery + '"';
            } else {
                if (badgeBar) badgeBar.style.display = 'none';
            }
            fetchPendingContent();
        }

        function clearContentSearch() {
            const input = document.getElementById('content-search-input');
            if (input) input.value = '';
            activeContentSearchQuery = '';
            currentContentPage = 1;

            const clearBtn = document.getElementById('content-search-clear-btn');
            if (clearBtn) clearBtn.style.display = 'none';

            const badgeBar = document.getElementById('active-search-badge-bar');
            if (badgeBar) badgeBar.style.display = 'none';

            fetchPendingContent();
        }

        function filterContentStatus(status) {
            activeContentFilter = status;
            currentContentPage = 1;

            isViralOnlyFilter = false;
            const viralBtn = document.getElementById('filter-btn-viral');
            if (viralBtn) {
                viralBtn.style.background = 'transparent';
                viralBtn.style.color = '#f97316';
                viralBtn.style.borderColor = 'rgba(249,115,22,0.4)';
            }

            ['pending', 'published', 'all', 'rejected'].forEach(s => {
                const btn = document.getElementById('filter-btn-' + s);
                if (btn) {
                    if (s === status.toLowerCase()) {
                        if (s === 'rejected') {
                            btn.style.background = 'rgba(244,63,94,0.2)';
                            btn.style.color = '#f43f5e';
                            btn.style.borderColor = '#f43f5e';
                        } else if (s === 'published') {
                            btn.style.background = 'rgba(34,197,94,0.15)';
                            btn.style.color = '#22c55e';
                            btn.style.borderColor = '#22c55e';
                        } else {
                            btn.style.background = 'rgba(56,189,248,0.15)';
                            btn.style.color = '#38bdf8';
                            btn.style.borderColor = '#38bdf8';
                        }
                    } else {
                        btn.style.background = 'transparent';
                        if (s === 'rejected') {
                            btn.style.color = '#f43f5e';
                            btn.style.borderColor = 'rgba(244,63,94,0.3)';
                        } else {
                            btn.style.color = 'var(--text-muted)';
                            btn.style.borderColor = '#3f3f4e';
                        }
                    }
                }
            });

            const emptyTrashBtn = document.getElementById('btn-empty-trash');
            if (emptyTrashBtn) {
                emptyTrashBtn.style.display = (status === 'REJECTED') ? 'inline-block' : 'none';
            }

            const summary = document.getElementById('content-status-summary');
            if (summary) {
                if (status === 'PENDING') summary.innerText = 'Showing pending review queue';
                else if (status === 'PUBLISHED') summary.innerText = 'Showing approved live articles';
                else if (status === 'REJECTED') summary.innerText = 'Showing rejected items queue';
                else summary.innerText = 'Showing all active content feed';
            }
            fetchPendingContent();
            updateDiscardedCountBadge();
        }

        const districtMapImageMap = {
            'Tamil Nadu':      'https://upload.wikimedia.org/wikipedia/commons/thumb/2/2c/Tamil_Nadu_districts_map.svg/960px-Tamil_Nadu_districts_map.svg.png',
            'National':        'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b3/India_map_with_states_and_union_territories.svg/960px-India_map_with_states_and_union_territories.svg.png',
            'International':   'https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/World_map_-_low_resolution.svg/960px-World_map_-_low_resolution.svg.png',
            'Madurai':         'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Madurai_in_Tamil_Nadu_%28India%29.svg/960px-Madurai_in_Tamil_Nadu_%28India%29.svg.png',
            'Chennai':         'https://upload.wikimedia.org/wikipedia/commons/thumb/4/48/Chennai_in_Tamil_Nadu_%28India%29.svg/960px-Chennai_in_Tamil_Nadu_%28India%29.svg.png',
            'Coimbatore':      'https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Coimbatore_in_Tamil_Nadu_%28India%29.svg/960px-Coimbatore_in_Tamil_Nadu_%28India%29.svg.png',
            'Salem':           'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ee/Salem_in_Tamil_Nadu_%28India%29.svg/960px-Salem_in_Tamil_Nadu_%28India%29.svg.png',
            'Thanjavur':       'https://upload.wikimedia.org/wikipedia/commons/thumb/7/78/Thanjavur_in_Tamil_Nadu_%28India%29.svg/960px-Thanjavur_in_Tamil_Nadu_%28India%29.svg.png',
            'Kanyakumari':     'https://upload.wikimedia.org/wikipedia/commons/thumb/5/5e/Kanyakumari_in_Tamil_Nadu_%28India%29.svg/960px-Kanyakumari_in_Tamil_Nadu_%28India%29.svg.png',
            'Vellore':         'https://upload.wikimedia.org/wikipedia/commons/thumb/4/44/Vellore_in_Tamil_Nadu_%28India%29.svg/960px-Vellore_in_Tamil_Nadu_%28India%29.svg.png',
            'Tiruchirappalli': 'https://upload.wikimedia.org/wikipedia/commons/thumb/4/49/Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg/960px-Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg.png',
            'Trichy':          'https://upload.wikimedia.org/wikipedia/commons/thumb/4/49/Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg/960px-Tiruchirappalli_in_Tamil_Nadu_%28India%29.svg.png',
            'Dindigul':        'https://upload.wikimedia.org/wikipedia/commons/thumb/5/52/Dindigul_in_Tamil_Nadu_%28India%29.svg/960px-Dindigul_in_Tamil_Nadu_%28India%29.svg.png',
            'Erode':           'https://upload.wikimedia.org/wikipedia/commons/thumb/6/60/Erode_in_Tamil_Nadu_%28India%29.svg/960px-Erode_in_Tamil_Nadu_%28India%29.svg.png',
            'Tirunelveli':     'https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Tirunelveli_in_Tamil_Nadu_%28India%29.svg/960px-Tirunelveli_in_Tamil_Nadu_%28India%29.svg.png',
            'Thoothukudi':     'https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Thoothukudi_in_Tamil_Nadu_%28India%29.svg/960px-Thoothukudi_in_Tamil_Nadu_%28India%29.svg.png',
            'Tuticorin':       'https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Thoothukudi_in_Tamil_Nadu_%28India%29.svg/960px-Thoothukudi_in_Tamil_Nadu_%28India%29.svg.png',
            'Tiruppur':        'https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Tiruppur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruppur_in_Tamil_Nadu_%28India%29.svg.png',
            'Nilgiris':        'https://upload.wikimedia.org/wikipedia/commons/thumb/6/6c/Nilgiris_in_Tamil_Nadu_%28India%29.svg/960px-Nilgiris_in_Tamil_Nadu_%28India%29.svg.png',
            'Cuddalore':       'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Cuddalore_in_Tamil_Nadu_%28India%29.svg/960px-Cuddalore_in_Tamil_Nadu_%28India%29.svg.png',
            'Kanchipuram':     'https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Kanchipuram_in_Tamil_Nadu_%28India%29.svg/960px-Kanchipuram_in_Tamil_Nadu_%28India%29.svg.png',
            'Chengalpattu':    'https://upload.wikimedia.org/wikipedia/commons/thumb/1/18/Chengalpattu_in_Tamil_Nadu_%28India%29.svg/960px-Chengalpattu_in_Tamil_Nadu_%28India%29.svg.png',
            'Tiruvallur':      'https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Tiruvallur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvallur_in_Tamil_Nadu_%28India%29.svg.png',
            'Pudukkottai':     'https://upload.wikimedia.org/wikipedia/commons/thumb/d/d3/Pudukkottai_in_Tamil_Nadu_%28India%29.svg/960px-Pudukkottai_in_Tamil_Nadu_%28India%29.svg.png',
            'Sivaganga':       'https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Sivaganga_in_Tamil_Nadu_%28India%29.svg/960px-Sivaganga_in_Tamil_Nadu_%28India%29.svg.png',
            'Ramanathapuram':  'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Ramanathapuram_in_Tamil_Nadu_%28India%29.svg/960px-Ramanathapuram_in_Tamil_Nadu_%28India%29.svg.png',
            'Nagapattinam':    'https://upload.wikimedia.org/wikipedia/commons/thumb/e/e0/Nagapattinam_in_Tamil_Nadu_%28India%29.svg/960px-Nagapattinam_in_Tamil_Nadu_%28India%29.svg.png',
            'Ranipet':         'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9e/Ranipet_in_Tamil_Nadu_%28India%29.svg/960px-Ranipet_in_Tamil_Nadu_%28India%29.svg.png',
            'Tirupattur':      'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4c/Tirupattur_in_Tamil_Nadu_%28India%29.svg/960px-Tirupattur_in_Tamil_Nadu_%28India%29.svg.png',
            'Tiruvannamalai':  'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Tiruvannamalai_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvannamalai_in_Tamil_Nadu_%28India%29.svg.png',
            'Tiruvarur':       'https://upload.wikimedia.org/wikipedia/commons/thumb/c/cb/Tiruvarur_in_Tamil_Nadu_%28India%29.svg/960px-Tiruvarur_in_Tamil_Nadu_%28India%29.svg.png',
            'Kallakurichi':    'https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Kallakurichi_in_Tamil_Nadu_%28India%29.svg/960px-Kallakurichi_in_Tamil_Nadu_%28India%29.svg.png',
            'Dharmapuri':      'https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Dharmapuri_in_Tamil_Nadu_%28India%29.svg/960px-Dharmapuri_in_Tamil_Nadu_%28India%29.svg.png',
            'Krishnagiri':     'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4a/Krishnagiri_in_Tamil_Nadu_%28India%29.svg/960px-Krishnagiri_in_Tamil_Nadu_%28India%29.svg.png',
            'Theni':           'https://upload.wikimedia.org/wikipedia/commons/thumb/5/50/Theni_in_Tamil_Nadu_%28India%29.svg/960px-Theni_in_Tamil_Nadu_%28India%29.svg.png',
            'Tenkasi':         'https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Tenkasi_in_Tamil_Nadu_%28India%29.svg/960px-Tenkasi_in_Tamil_Nadu_%28India%29.svg.png',
            'Namakkal':        'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c5/Namakkal_in_Tamil_Nadu_%28India%29.svg/960px-Namakkal_in_Tamil_Nadu_%28India%29.svg.png',
            'Ariyalur':        'https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Ariyalur_in_Tamil_Nadu_%28India%29.svg/960px-Ariyalur_in_Tamil_Nadu_%28India%29.svg.png',
            'Perambalur':      'https://upload.wikimedia.org/wikipedia/commons/thumb/e/e5/Perambalur_in_Tamil_Nadu_%28India%29.svg/960px-Perambalur_in_Tamil_Nadu_%28India%29.svg.png',
            'Viluppuram':      'https://upload.wikimedia.org/wikipedia/commons/thumb/d/d4/Viluppuram_in_Tamil_Nadu_%28India%29.svg/960px-Viluppuram_in_Tamil_Nadu_%28India%29.svg.png',
            'Virudhunagar':    'https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Virudhunagar_in_Tamil_Nadu_%28India%29.svg/960px-Virudhunagar_in_Tamil_Nadu_%28India%29.svg.png',
            'Mayiladuthurai':  'https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/Mayiladuthurai_in_Tamil_Nadu_%28India%29.svg/960px-Mayiladuthurai_in_Tamil_Nadu_%28India%29.svg.png',
            'Karur':           'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a2/Karur_in_Tamil_Nadu_%28India%29.svg/960px-Karur_in_Tamil_Nadu_%28India%29.svg.png',
        };

        const personImageMap = {
            'ஸ்டாலின்':           'https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg',
            'மு.க.ஸ்டாலின்':       'https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg',
            'stalin':             'https://upload.wikimedia.org/wikipedia/commons/9/9d/The_Chief_Minister_of_Tamil_Nadu%2C_Thiru_M.K._Stalin.jpg',
            'எடப்பாடி':          'https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg',
            'பழனிசாமி':           'https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg',
            'edappadi':           'https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg',
            'eps':                'https://upload.wikimedia.org/wikipedia/commons/1/1e/EdappadiKPalaniswami.jpg',
            'விஜய்':             'https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg',
            'vijay':             'https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg',
            'tvk':               'https://upload.wikimedia.org/wikipedia/commons/thumb/c/cd/Vijay_at_the_Nadigar_Sangam_Protest.jpg/800px-Vijay_at_the_Nadigar_Sangam_Protest.jpg',
            'தோனி':              'https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg',
            'dhoni':              'https://upload.wikimedia.org/wikipedia/commons/thumb/7/70/M.S._Dhoni_%28Prabal_Pardesi%29.jpg/800px-M.S._Dhoni_%28Prabal_Pardesi%29.jpg',
            'மம்தா':             'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg',
            'mamata':             'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4d/Official_portrait_of_Mamata_Banerjee.jpg/800px-Official_portrait_of_Mamata_Banerjee.jpg',
            'மோடி':              'https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg',
            'modi':               'https://upload.wikimedia.org/wikipedia/commons/b/ba/Narendra_Modi_Portrait_2026.jpg',
            'ராகுல்':             'https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Rahul_Gandhi_in_2023.jpg/800px-Rahul_Gandhi_in_2023.jpg',
            'rahul':              'https://upload.wikimedia.org/wikipedia/commons/thumb/9/97/Rahul_Gandhi_in_2023.jpg/800px-Rahul_Gandhi_in_2023.jpg',
            'அமித் ஷா':          'https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Amit_Shah_in_2023.jpg/800px-Amit_Shah_in_2023.jpg',
            'amit shah':          'https://upload.wikimedia.org/wikipedia/commons/thumb/f/f9/Amit_Shah_in_2023.jpg/800px-Amit_Shah_in_2023.jpg',
            'சீமான்':            'https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Seeman_at_a_meeting.jpg/800px-Seeman_at_a_meeting.jpg',
            'seeman':            'https://upload.wikimedia.org/wikipedia/commons/thumb/7/77/Seeman_at_a_meeting.jpg/800px-Seeman_at_a_meeting.jpg',
            'அண்ணாமலை':          'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/K._Annamalai_in_2023.jpg/800px-K._Annamalai_in_2023.jpg',
            'annamalai':          'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/K._Annamalai_in_2023.jpg/800px-K._Annamalai_in_2023.jpg',
            'உதயநிதி':           'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Udhayanidhi_Stalin.jpg/800px-Udhayanidhi_Stalin.jpg',
            'udhayanidhi':        'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Udhayanidhi_Stalin.jpg/800px-Udhayanidhi_Stalin.jpg',
            'கோலி':              'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg',
            'kohli':              'https://upload.wikimedia.org/wikipedia/commons/thumb/e/ef/Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg/800px-Virat_Kohli_during_the_India_vs_Aus_4th_Test_match_at_Narendra_Modi_Stadium_on_09_March_2023.jpg',
            'ரோஹித்':             'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg',
            'rohit':              'https://upload.wikimedia.org/wikipedia/commons/thumb/1/1d/Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg/800px-Rohit_Sharma_during_the_2019_Cricket_World_Cup.jpg',
            'டிரம்ப்':            'https://upload.wikimedia.org/wikipedia/commons/thumb/5/56/Donald_Trump_official_portrait.jpg/800px-Donald_Trump_official_portrait.jpg',
            'trump':              'https://upload.wikimedia.org/wikipedia/commons/thumb/5/56/Donald_Trump_official_portrait.jpg/800px-Donald_Trump_official_portrait.jpg',
            'புதின்':             'https://upload.wikimedia.org/wikipedia/commons/thumb/8/8d/Vladimir_Putin_%282020-02-20%29.jpg/800px-Vladimir_Putin_%282020-02-20%29.jpg',
            'putin':              'https://upload.wikimedia.org/wikipedia/commons/thumb/8/8d/Vladimir_Putin_%282020-02-20%29.jpg/800px-Vladimir_Putin_%282020-02-20%29.jpg',
            'கமல்':              'https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Kamal_Haasan_in_2022.jpg/800px-Kamal_Haasan_in_2022.jpg',
            'kamal':              'https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Kamal_Haasan_in_2022.jpg/800px-Kamal_Haasan_in_2022.jpg',
            'ரஜினி':             'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Rajinikanth_in_2023.jpg/800px-Rajinikanth_in_2023.jpg',
            'rajini':             'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b8/Rajinikanth_in_2023.jpg/800px-Rajinikanth_in_2023.jpg',
            'நிர்மலா':           'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9f/Nirmala_Sitharaman_in_2022.jpg/800px-Nirmala_Sitharaman_in_2022.jpg',
            'nirmala':            'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9f/Nirmala_Sitharaman_in_2022.jpg/800px-Nirmala_Sitharaman_in_2022.jpg',
            'மெஸ்ஸி':             'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg',
            'லியோனல் மெஸ்ஸி':    'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg',
            'messi':              'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c1/Lionel_Messi_20180626.jpg/800px-Lionel_Messi_20180626.jpg',
            'ரொனால்டோ':          'https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg',
            'ronaldo':            'https://upload.wikimedia.org/wikipedia/commons/thumb/8/8c/Cristiano_Ronaldo_2018.jpg/800px-Cristiano_Ronaldo_2018.jpg',
            'நெய்மர்':           'https://upload.wikimedia.org/wikipedia/commons/thumb/8/83/ISL_Final_2019-20_%28Neymar_cropped%29.jpg/800px-ISL_Final_2019-20_%28Neymar_cropped%29.jpg',
            'neymar':             'https://upload.wikimedia.org/wikipedia/commons/thumb/8/83/ISL_Final_2019-20_%28Neymar_cropped%29.jpg/800px-ISL_Final_2019-20_%28Neymar_cropped%29.jpg',
            'அஸ்வின்':            'https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Ravichandran_Ashwin.jpg/800px-Ravichandran_Ashwin.jpg',
            'ashwin':             'https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Ravichandran_Ashwin.jpg/800px-Ravichandran_Ashwin.jpg',
            'சச்சின்':           'https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Sachin_Tendulkar_at_MRF_Promotion_Event.jpg/800px-Sachin_Tendulkar_at_MRF_Promotion_Event.jpg',
            'sachin':             'https://upload.wikimedia.org/wikipedia/commons/thumb/2/25/Sachin_Tendulkar_at_MRF_Promotion_Event.jpg/800px-Sachin_Tendulkar_at_MRF_Promotion_Event.jpg',
            'சூர்யா':            'https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Suriya_at_2D_Entertainment_office.jpg/800px-Suriya_at_2D_Entertainment_office.jpg',
            'suriya':             'https://upload.wikimedia.org/wikipedia/commons/thumb/7/7b/Suriya_at_2D_Entertainment_office.jpg/800px-Suriya_at_2D_Entertainment_office.jpg',
            'விக்ரம்':            'https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Vikram_at_Cobra_promotions.jpg/800px-Vikram_at_Cobra_promotions.jpg',
            'vikram':             'https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Vikram_at_Cobra_promotions.jpg/800px-Vikram_at_Cobra_promotions.jpg',
            'தனுஷ்':             'https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Dhanush_at_The_Gray_Man_press_conference.jpg/800px-Dhanush_at_The_Gray_Man_press_conference.jpg',
            'dhanush':            'https://upload.wikimedia.org/wikipedia/commons/thumb/8/87/Dhanush_at_The_Gray_Man_press_conference.jpg/800px-Dhanush_at_The_Gray_Man_press_conference.jpg',
            'சிவகார்த்திகேயன்':    'https://upload.wikimedia.org/wikipedia/commons/thumb/3/30/Sivakarthikeyan_at_Doctor_success_meet.jpg/800px-Sivakarthikeyan_at_Doctor_success_meet.jpg',
            'sivakarthikeyan':    'https://upload.wikimedia.org/wikipedia/commons/thumb/3/30/Sivakarthikeyan_at_Doctor_success_meet.jpg/800px-Sivakarthikeyan_at_Doctor_success_meet.jpg',
            'ரஹ்மான்':           'https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg',
            'ar rahman':          'https://upload.wikimedia.org/wikipedia/commons/thumb/3/36/A._R._Rahman_at_NMACC_Gala.jpg/800px-A._R._Rahman_at_NMACC_Gala.jpg',
            'அனிருத்':           'https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Anirudh_Ravichander_at_Vikram_audio_launch.jpg/800px-Anirudh_Ravichander_at_Vikram_audio_launch.jpg',
            'anirudh':            'https://upload.wikimedia.org/wikipedia/commons/thumb/6/6f/Anirudh_Ravichander_at_Vikram_audio_launch.jpg/800px-Anirudh_Ravichander_at_Vikram_audio_launch.jpg',
        };

        const categoryImageMap = {
            'sports':        'https://images.unsplash.com/photo-1508098682722-e99c43a406b2?w=1200&auto=format&fit=crop&q=80',
            'politics':      'https://images.unsplash.com/photo-1541872703-74c5e44368f9?w=1200&auto=format&fit=crop&q=80',
            'business':      'https://images.unsplash.com/photo-1590283603385-17ffb3a7f29f?w=1200&auto=format&fit=crop&q=80',
            'technical':     'https://images.unsplash.com/photo-1518770660439-4636190af475?w=1200&auto=format&fit=crop&q=80',
            'entertainment': 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=1200&auto=format&fit=crop&q=80',
            'crime':         'https://images.unsplash.com/photo-1589829545856-d10d557cf95f?w=1200&auto=format&fit=crop&q=80',
            'civic':         'https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?w=1200&auto=format&fit=crop&q=80',
            'news':          'https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=1200&auto=format&fit=crop&q=80',
        };

        function getFallbackDistrictImage(district, category, title) {
            if (title) {
                const tLower = title.toLowerCase();
                for (const [name, img] of Object.entries(personImageMap)) {
                    if (tLower.includes(name.toLowerCase())) {
                        return img;
                    }
                }
            }
            const cLower = (category || '').toLowerCase();
            for (const [name, img] of Object.entries(categoryImageMap)) {
                if (cLower.includes(name)) {
                    return img;
                }
            }
            if (title) {
                const tLower = title.toLowerCase();
                if (tLower.includes('விளையாட்டு') || tLower.includes('கிரிக்கெட்') || tLower.includes('கால்பந்து') || tLower.includes('மெஸ்ஸி') || tLower.includes('cricket') || tLower.includes('football') || tLower.includes('sports') || tLower.includes('ipl') || tLower.includes('match') || tLower.includes('goal')) {
                    return categoryImageMap['sports'];
                }
                if (tLower.includes('அரசியல்') || tLower.includes('politics') || tLower.includes('தேர்தல்') || tLower.includes('அரசு')) {
                    return categoryImageMap['politics'];
                }
                if (tLower.includes('திரைப்பட') || tLower.includes('சினிமா') || tLower.includes('cinema') || tLower.includes('movie')) {
                    return categoryImageMap['entertainment'];
                }
                if (tLower.includes('வணிகம்') || tLower.includes('business') || tLower.includes('பொருளாதார')) {
                    return categoryImageMap['business'];
                }
                if (tLower.includes('தொழில்நுட்ப') || tLower.includes('tech') || tLower.includes('ai') || tLower.includes('இஸ்ரோ')) {
                    return categoryImageMap['technical'];
                }
                if (tLower.includes('குற்றம்') || tLower.includes('crime') || tLower.includes('கைது') || tLower.includes('போலீஸ்')) {
                    return categoryImageMap['crime'];
                }
            }
            const dClean = (district || 'Tamil Nadu').trim();
            if (districtMapImageMap[dClean]) {
                return districtMapImageMap[dClean];
            }
            for (const [k, v] of Object.entries(districtMapImageMap)) {
                if (k.toLowerCase() === dClean.toLowerCase()) return v;
            }
            return '/admin/api/maps/svg?district=' + encodeURIComponent(dClean);
        }

        let activeModLanguage = 'all';
        let activeModDistrict = 'all';
        let activeModCategory = 'all';
        let activeModGroupBy = 'none';
        let activeModSort = 'date_desc';

        function changeModSort(val) {
            activeModSort = val;
            updateSortIcons();
            currentContentPage = 1;
            fetchPendingContent();
        }

        function toggleSourceDateSort() {
            if (activeModSort === 'date_desc') {
                activeModSort = 'source_date';
            } else if (activeModSort === 'source_date') {
                activeModSort = 'date_asc';
            } else if (activeModSort === 'date_asc') {
                activeModSort = 'source_desc';
            } else {
                activeModSort = 'date_desc';
            }
            const sortSelect = document.getElementById('mod-sort-select');
            if (sortSelect) sortSelect.value = activeModSort;
            updateSortIcons();
            currentContentPage = 1;
            fetchPendingContent();
        }

        function updateSortIcons() {
            const icon = document.getElementById('sort-source-date-icon');
            if (!icon) return;
            if (activeModSort === 'date_desc') {
                icon.innerHTML = '🕒 ▼';
                icon.title = 'Sorted by Post Date (Newest first)';
            } else if (activeModSort === 'date_asc') {
                icon.innerHTML = '🕒 ▲';
                icon.title = 'Sorted by Post Date (Oldest first)';
            } else if (activeModSort === 'source_date') {
                icon.innerHTML = '📰 A→Z';
                icon.title = 'Sorted by Source (A-Z) & Post Date';
            } else if (activeModSort === 'source_desc') {
                icon.innerHTML = '📰 Z→A';
                icon.title = 'Sorted by Source (Z-A) & Post Date';
            } else {
                icon.innerHTML = '⇅';
                icon.title = 'Click to sort by Source & Post Date';
            }
        }

        function decodeHTMLEntities(str) {
            if (!str) return '';
            let s = String(str).replace(/&#x([0-9a-fA-F]+);?/g, function(match, hex) {
                return String.fromCharCode(parseInt(hex, 16));
            }).replace(/&#([0-9]+);?/g, function(match, dec) {
                return String.fromCharCode(parseInt(dec, 10));
            });
            const doc = new DOMParser().parseFromString(s, 'text/html');
            return doc.body.textContent || s;
        }

        function escapeHtml(str) {
            if (!str) return '';
            return String(str)
                .replace(/&/g, '&amp;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;')
                .replace(/"/g, '&quot;')
                .replace(/'/g, '&#039;');
        }

        function formatSourceDate(dateStr) {
            if (!dateStr) return 'Recent';
            const d = new Date(dateStr);
            if (isNaN(d.getTime())) return 'Recent';
            const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
            const m = months[d.getMonth()];
            const day = d.getDate();
            const yr = d.getFullYear();
            let hr = d.getHours();
            const min = d.getMinutes() < 10 ? '0' + d.getMinutes() : d.getMinutes();
            const ampm = hr >= 12 ? 'PM' : 'AM';
            hr = hr % 12;
            hr = hr ? hr : 12;
            const hrStr = hr < 10 ? '0' + hr : hr;
            return m + ' ' + day + ', ' + yr + ' ' + hrStr + ':' + min + ' ' + ampm;
        }

        function getSourceHostName(url) {
            if (!url) return 'Direct';
            try {
                const u = new URL(url);
                let host = u.hostname.replace(/^www\./, '').replace(/^m\./, '').toLowerCase();
                if (host.includes('bbc.com')) return 'BBC Tamil';
                if (host.includes('vikatan.com')) return 'Vikatan';
                if (host.includes('thehindu.com')) return 'The Hindu';
                if (host.includes('dinamalar.com')) return 'Dinamalar';
                if (host.includes('dailythanthi.com')) return 'Daily Thanthi';
                if (host.includes('dinamani.com')) return 'Dinamani';
                if (host.includes('maalaimalar.com')) return 'Maalai Malar';
                if (host.includes('puthiyathalaimurai.com')) return 'Puthiyathalaimurai';
                if (host.includes('news18.com')) return 'News18 Tamil';
                if (host.includes('dinakaran.com')) return 'Dinakaran';
                if (host.includes('newstamil.tv')) return 'News Tamil 24x7';
                if (host.includes('news7tamil.live')) return 'News7 Tamil';
                if (host.includes('indianexpress.com')) return 'Indian Express';
                if (host.includes('oneindia.com')) return 'OneIndia Tamil';
                if (host.includes('polimernews.com')) return 'Polimer News';
                if (host.includes('nakkheeran.in')) return 'Nakkheeran';
                if (host.includes('youtube.com') || host.includes('youtu.be')) return 'YouTube';
                return host;
            } catch(e) {
                return 'Web Source';
            }
        }

        function setModLanguageFilter(lang) {
            activeModLanguage = lang;
            const selectEl = document.getElementById('mod-lang-select');
            if (selectEl && selectEl.value !== lang) {
                selectEl.value = lang;
            }
            currentContentPage = 1;
            fetchPendingContent();
        }

        function openLanguageConfigModal() {
            fetch('/admin/api/settings/language')
                .then(res => res.json())
                .then(data => {
                    if (data && data.settings) {
                        const defLang = data.settings.default_language || 'ta';
                        const enabledLangs = (data.settings.enabled_languages || 'ta,en,ta-en').split(',').map(s => s.trim());
                        const policy = data.settings.scraper_language_policy || 'ALL';

                        const defSelect = document.getElementById('cfg-default-language');
                        if (defSelect) defSelect.value = defLang;

                        document.querySelectorAll('.cfg-lang-cb').forEach(cb => {
                            cb.checked = enabledLangs.includes(cb.value);
                        });

                        const polSelect = document.getElementById('cfg-scraper-policy');
                        if (polSelect) polSelect.value = policy;
                    }
                    document.getElementById('languageConfigModal').style.display = 'flex';
                })
                .catch(() => {
                    document.getElementById('languageConfigModal').style.display = 'flex';
                });
        }

        function closeLanguageConfigModal() {
            document.getElementById('languageConfigModal').style.display = 'none';
        }

        function saveLanguageConfig() {
            const defLang = document.getElementById('cfg-default-language').value;
            const enabledLangs = [];
            document.querySelectorAll('.cfg-lang-cb:checked').forEach(cb => {
                enabledLangs.push(cb.value);
            });
            if (enabledLangs.length === 0) {
                showToast('⚠️ Please enable at least one language');
                return;
            }
            if (!enabledLangs.includes(defLang)) {
                enabledLangs.push(defLang);
            }
            const policy = document.getElementById('cfg-scraper-policy').value;

            showToast('Saving language & localization settings...');
            fetch('/admin/api/settings/language', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    default_language: defLang,
                    enabled_languages: enabledLangs.join(','),
                    scraper_language_policy: policy
                })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Language settings saved successfully');
                    closeLanguageConfigModal();
                    syncLanguageToolbar(enabledLangs, defLang);
                    fetchPendingContent();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to save settings'));
                }
            })
            .catch(() => showToast('✗ Network error saving language settings'));
        }

        function syncLanguageToolbar(enabledLangs, defLang) {
            const langLabels = {
                'ta': '🇮🇳 தமிழ் (Tamil)',
                'en': '🇬🇧 English',
                'ta-en': '🔄 Tanglish (TA-EN)',
                'hi': '🇮🇳 हिन्दी (Hindi)',
                'ml': '🇮🇳 മലയാളം (Malayalam)',
                'te': '🇮🇳 తెలుగు (Telugu)',
                'kn': '🇮🇳 ಕನ್ನಡ (Kannada)'
            };
            const sel = document.getElementById('mod-lang-select');
            if (!sel) return;
            const currentVal = activeModLanguage || 'all';
            let opts = '<option value="all">🌐 All Languages</option>';
            enabledLangs.forEach(code => {
                const label = langLabels[code] || ('🌐 ' + code.toUpperCase());
                opts += '<option value="' + code + '">' + label + '</option>';
            });
            sel.innerHTML = opts;
            if (enabledLangs.includes(currentVal) || currentVal === 'all') {
                sel.value = currentVal;
            } else {
                sel.value = 'all';
                activeModLanguage = 'all';
            }
        }

        function loadLanguageConfig() {
            fetch('/admin/api/settings/language')
                .then(res => res.json())
                .then(data => {
                    if (data && data.settings) {
                        const defLang = data.settings.default_language || 'ta';
                        const enabledLangs = (data.settings.enabled_languages || 'ta,en,ta-en').split(',').map(s => s.trim());
                        syncLanguageToolbar(enabledLangs, defLang);
                    }
                })
                .catch(() => {});
        }

        function applyModFilters() {
            const distSelect = document.getElementById('mod-district-select');
            if (distSelect) activeModDistrict = distSelect.value;
            const catSelect = document.getElementById('mod-category-select');
            if (catSelect) activeModCategory = catSelect.value;
            currentContentPage = 1;
            fetchPendingContent();
        }

        function changeModGroupBy(val) {
            activeModGroupBy = val;
            fetchPendingContent();
        }

        function selectGroupItems(encodedGroupKey) {
            const groupKey = decodeURIComponent(encodedGroupKey);
            const checkboxes = document.querySelectorAll('tr[data-group="' + CSS.escape(groupKey) + '"] .content-row-checkbox');
            let anyUnchecked = false;
            checkboxes.forEach(cb => {
                if (!cb.checked) anyUnchecked = true;
            });
            checkboxes.forEach(cb => {
                cb.checked = anyUnchecked;
            });
            updateSelectionState();
        }

        let isViralOnlyFilter = false;

        function toggleViralFilter() {
            isViralOnlyFilter = !isViralOnlyFilter;
            const btn = document.getElementById('filter-btn-viral');
            if (btn) {
                if (isViralOnlyFilter) {
                    btn.style.background = 'rgba(249,115,22,0.25)';
                    btn.style.color = '#f97316';
                    btn.style.borderColor = '#f97316';
                } else {
                    btn.style.background = 'transparent';
                    btn.style.color = '#f97316';
                    btn.style.borderColor = 'rgba(249,115,22,0.4)';
                }
            }
            currentContentPage = 1;
            fetchPendingContent();
        }

        function renderContentTableRow(item, groupName = '') {
            cachedContentMap[item.id] = item;
            const fallbackThumb = getFallbackDistrictImage(item.district, item.category, item.title);
            let mediaHtml = '';
            if (item.videoId || item.videoUrl) {
                const ytThumb = item.videoId && !item.videoId.startsWith('vid_') && !item.videoId.startsWith('p0') ? ('https://img.youtube.com/vi/' + item.videoId + '/hqdefault.jpg') : '';
                const thumbUrl = item.thumbnail || ytThumb || fallbackThumb;
                mediaHtml = '<div style="position:relative; width:110px; height:64px; border-radius:6px; overflow:hidden; border:1px solid #38bdf8; cursor:pointer;" onclick="openContentViewModal(\'' + item.id + '\')">' +
                    '<img src="' + thumbUrl + '" onerror="this.onerror=null; this.src=\'' + fallbackThumb + '\';" style="width:100%; height:100%; object-fit:cover;" />' +
                    '<div style="position:absolute; inset:0; display:flex; align-items:center; justify-content:center; background:rgba(0,0,0,0.4); color:#fff; font-size:16px;">▶</div>' +
                '</div>';
            } else {
                const thumbUrl = item.thumbnail || fallbackThumb;
                const mapFallbackUrl = '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
                mediaHtml = '<img src="' + thumbUrl + '" onerror="this.onerror=null; this.src=\'' + mapFallbackUrl + '\';" style="width:110px; height:64px; object-fit:cover; border-radius:6px; border:1px solid #3f3f4e; cursor:pointer;" onclick="openContentViewModal(\'' + item.id + '\')" alt="Geographic Map" />';
            }

            const cleanTitle = decodeHTMLEntities(item.title);
            const cleanDesc = decodeHTMLEntities(item.description || '');

            const langCode = (item.language || '').toLowerCase();
            let langBadge = '';
            if (langCode === 'ta') {
                langBadge = '<span class="badge-pill" style="background: rgba(16,185,129,0.15); color: #10b981; border: 1px solid rgba(16,185,129,0.4); font-weight:700; font-size:10px;">🇮🇳 தமிழ்</span> ';
            } else if (langCode === 'en') {
                langBadge = '<span class="badge-pill" style="background: rgba(56,189,248,0.15); color: #38bdf8; border: 1px solid rgba(56,189,248,0.4); font-weight:700; font-size:10px;">🇬🇧 EN</span> ';
            } else if (langCode === 'ta-en' || langCode === 'tanglish') {
                langBadge = '<span class="badge-pill" style="background: rgba(168,85,247,0.15); color: #c084fc; border: 1px solid rgba(168,85,247,0.4); font-weight:700; font-size:10px;">🔄 TA-EN</span> ';
            } else if (langCode === 'hi') {
                langBadge = '<span class="badge-pill" style="background: rgba(245,158,11,0.15); color: #f59e0b; border: 1px solid rgba(245,158,11,0.4); font-weight:700; font-size:10px;">🇮🇳 हिन्दी</span> ';
            } else if (langCode) {
                langBadge = '<span class="badge-pill" style="background: rgba(148,163,184,0.15); color: #94a3b8; border: 1px solid rgba(148,163,184,0.4); font-weight:700; font-size:10px;">🌐 ' + escapeHtml(langCode.toUpperCase()) + '</span> ';
            } else {
                const isTamilAuto = (/[\u0B80-\u0BFF]/.test(cleanTitle + ' ' + cleanDesc));
                langBadge = isTamilAuto
                    ? '<span class="badge-pill" style="background: rgba(16,185,129,0.15); color: #10b981; border: 1px solid rgba(16,185,129,0.4); font-weight:700; font-size:10px;">🇮🇳 தமிழ்</span> '
                    : '<span class="badge-pill" style="background: rgba(56,189,248,0.15); color: #38bdf8; border: 1px solid rgba(56,189,248,0.4); font-weight:700; font-size:10px;">🇬🇧 EN</span> ';
            }

            const viralBadge = item.isViral 
                ? '<span class="badge-pill" style="background: rgba(249,115,22,0.2); color: #f97316; border: 1px solid rgba(249,115,22,0.5); font-weight: 700; font-size: 10px; padding: 2px 6px; letter-spacing: 0.5px;">🔥 VIRAL PRIORITY</span> ' 
                : '';

            const formatBadge = item.contentType === 'VIDEO_LINK'
                ? '<span class="badge-pill badge-warning">VIDEO</span>'
                : '<span class="badge-pill badge-info">TEXT</span>';

            const statusBadge = item.status === 'PUBLISHED'
                ? '<span class="badge-pill badge-success">APPROVED</span>'
                : (item.status === 'REJECTED'
                    ? '<span class="badge-pill badge-danger">REJECTED</span>'
                    : '<span class="badge-pill badge-warning">PENDING</span>');

            const srcName = getSourceHostName(item.sourceUrl);
            const dateDisplay = formatSourceDate(item.createdAt);
            const sourceCell = '<div style="font-size: 12px; font-weight: 600; color: #f1f5f9; white-space: nowrap;">' +
                (item.sourceUrl ? '<a href="' + item.sourceUrl + '" target="_blank" style="color:#38bdf8; text-decoration:none;" title="' + item.sourceUrl + '">🔗 ' + escapeHtml(srcName) + '</a>' : '<span style="color:var(--text-muted);">Direct</span>') +
                '</div>' +
                '<div style="font-size: 11px; color: #94a3b8; margin-top: 4px; white-space: nowrap;">' +
                '🕒 ' + dateDisplay +
                '</div>';

            const rowStyle = item.isViral ? 'background: rgba(249,115,22,0.03); border-left: 3px solid #f97316;' : '';

            const mainBannerBadge = item.isMainBanner 
                ? '<span class="badge-pill" style="background: rgba(234,179,8,0.2); color: #eab308; border: 1px solid rgba(234,179,8,0.6); font-weight: 800; font-size: 10px; padding: 2px 6px; letter-spacing: 0.5px; margin-right: 4px;">🌟 MAIN BANNER</span> ' 
                : '';

            const starBtn = '<button class="action-btn" style="color: ' + (item.isMainBanner ? '#eab308; border-color: #eab308; background: rgba(234,179,8,0.15);' : '#cbd5e1; border-color: #334155;') + ' margin-right: 4px;" onclick="quickToggleMainBanner(\'' + item.id + '\', ' + (!item.isMainBanner) + ')" title="' + (item.isMainBanner ? 'Current Portal Hero Banner (Click to unset)' : 'Set as Portal Main Hero Banner') + '">' + (item.isMainBanner ? '⭐ Hero: ON' : '☆ Set Hero') + '</button>';

            let rowActions = '<button class="action-btn" style="color: #38bdf8; border-color: rgba(56,189,248,0.4); margin-right: 4px;" onclick="openContentViewModal(\'' + item.id + '\')" title="Preview post">👁️ View</button>' +
                '<button class="action-btn" style="color: #a855f7; border-color: rgba(168,85,247,0.4); margin-right: 4px;" onclick="openEditContentModal(\'' + item.id + '\')" title="Edit post content and broadcast placement">✏️ Edit</button>' +
                starBtn;
            if (item.status === 'PENDING') {
                rowActions += '<button class="btn-approve" onclick="approveContent(\'' + item.id + '\')" style="margin-right: 4px;" title="Approve & publish live">✓ Approve</button>' +
                    '<button class="action-btn" style="color: #f59e0b; border-color: rgba(245,158,11,0.4); background: rgba(245,158,11,0.1); font-weight:600; margin-right: 4px;" onclick="rejectContent(\'' + item.id + '\')" title="Reject post">🚫 Reject</button>' +
                    '<button class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.4); background: rgba(244,63,94,0.1); font-weight:600;" onclick="deletePermanent(\'' + item.id + '\')" title="Permanently delete from database">🗑️ Delete</button>';
            } else if (item.status === 'REJECTED') {
                rowActions += '<button class="btn-approve" style="margin-right: 4px;" onclick="approveContent(\'' + item.id + '\')" title="Approve & publish live">✓ Approve</button>' +
                    '<button class="action-btn" style="color: #38bdf8; border-color: rgba(56,189,248,0.4); background: rgba(56,189,248,0.1); font-weight:600; margin-right: 4px;" onclick="restoreContent(\'' + item.id + '\')" title="Restore to pending review">↺ Restore</button>' +
                    '<button class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.4); background: rgba(244,63,94,0.1); font-weight:600;" onclick="deletePermanent(\'' + item.id + '\')" title="Permanently delete from database">🗑️ Delete</button>';
            } else {
                rowActions += '<button class="action-btn" style="color: #f59e0b; border-color: rgba(245,158,11,0.4); background: rgba(245,158,11,0.1); font-weight:600; margin-right: 4px;" onclick="rejectContent(\'' + item.id + '\')" title="Unpublish & Reject">🚫 Reject</button>' +
                    '<button class="action-btn" style="color: #f43f5e; border-color: rgba(244,63,94,0.4); background: rgba(244,63,94,0.1); font-weight:600;" onclick="deletePermanent(\'' + item.id + '\')" title="Permanently delete from database">🗑️ Delete</button>';
            }

            return '<tr style="' + rowStyle + '" data-group="' + (groupName ? encodeURIComponent(groupName) : '') + '">' +
                '<td style="text-align: center;"><input type="checkbox" class="content-row-checkbox" value="' + item.id + '" onchange="updateSelectionState()" style="cursor: pointer; width: 16px; height: 16px;" /></td>' +
                '<td>' + mediaHtml + '</td>' +
                '<td>' +
                    mainBannerBadge +
                    '<strong style="font-size: 14px; color: #f8fafc; cursor:pointer;" onclick="openContentViewModal(\'' + item.id + '\')">' + cleanTitle + '</strong><br>' +
                    '<span style="font-size: 12px; color: var(--text-muted); display:-webkit-box; -webkit-line-clamp:2; -webkit-box-orient:vertical; overflow:hidden;">' + cleanDesc + '</span>' +
                '</td>' +
                '<td>' +
                    langBadge + viralBadge + formatBadge + ' ' + statusBadge + '<br>' +
                    '<span style="font-size: 12px; color: #a5b4fc; font-weight:600;">📍 ' + (item.district || 'Tamil Nadu') + '</span> ' +
                    '<span style="font-size: 11px; color: #94a3b8; font-weight:500;">(' + (item.category || 'News') + ')</span>' +
                '</td>' +
                '<td>' + sourceCell + '</td>' +
                '<td style="white-space: nowrap;">' + rowActions + '</td>' +
            '</tr>';
        }

        function fetchPendingContent() {
            const tbody = document.getElementById('mod-table-body');
            if (tbody && (!tbody.children || tbody.children.length === 0)) {
                tbody.innerHTML = '<tr><td colspan="6" style="color: var(--text-muted); text-align: center; padding: 24px;">Loading content...</td></tr>';
            }

            const viralParam = isViralOnlyFilter ? '&viral=true' : '';
            const searchParam = activeContentSearchQuery ? ('&q=' + encodeURIComponent(activeContentSearchQuery)) : '';
            const langParam = activeModLanguage !== 'all' ? ('&lang=' + activeModLanguage) : '';
            const distParam = (activeModDistrict && activeModDistrict !== 'all') ? ('&district=' + encodeURIComponent(activeModDistrict)) : '';
            const catParam = (activeModCategory && activeModCategory !== 'all') ? ('&category=' + encodeURIComponent(activeModCategory)) : '';
            const sortParam = activeModSort ? ('&sort=' + encodeURIComponent(activeModSort)) : '';

            fetch('/admin/api/scraper/pending?status=' + activeContentFilter + '&page=' + currentContentPage + '&limit=' + currentContentLimit + viralParam + searchParam + langParam + distParam + catParam + sortParam)
                .then(res => res.json())
                .then(res => {
                    const tbody = document.getElementById('mod-table-body');
                    if (!tbody) return;
                    if (!res || !res.data || res.data.length === 0) {
                        let msg = activeContentSearchQuery 
                            ? 'No posts found matching "' + activeContentSearchQuery + '".'
                            : 'No items found for this view.';
                        if (!activeContentSearchQuery && activeContentFilter === 'PENDING') {
                            msg = '<div style="padding: 24px 12px; text-align: center;">' +
                                '<div style="font-size: 16px; font-weight: 700; color: #38bdf8; margin-bottom: 8px;">🎉 Queue Empty: All Pending Content Has Been Moderated</div>' +
                                '<div style="font-size: 13px; color: #94a3b8; margin-bottom: 16px;">All current items have been approved to live or rejected. You can view approved live articles or inspect all contents below.</div>' +
                                '<div style="display: flex; justify-content: center; gap: 10px; flex-wrap: wrap;">' +
                                    '<button class="action-btn" onclick="filterContentStatus(\'PUBLISHED\')" style="font-size: 12px; padding: 7px 16px; background: rgba(34,197,94,0.15); color: #22c55e; border-color: #22c55e; font-weight: 700; cursor: pointer;">🟢 View Approved (Live) Posts</button>' +
                                    '<button class="action-btn" onclick="filterContentStatus(\'ALL\')" style="font-size: 12px; padding: 7px 16px; background: rgba(56,189,248,0.15); color: #38bdf8; border-color: #38bdf8; font-weight: 700; cursor: pointer;">📋 View All Contents</button>' +
                                '</div>' +
                            '</div>';
                        }
                        tbody.innerHTML = '<tr><td colspan="6" style="color: var(--text-muted); text-align: center; padding: 24px;">' + msg + '</td></tr>';
                        const pagInfo = document.getElementById('pagination-info');
                        if (pagInfo) pagInfo.innerText = 'Showing 0 of 0 items';
                        const pagDisp = document.getElementById('page-current-display');
                        if (pagDisp) pagDisp.innerText = 'Page 1 of 1';
                        const prevBtn = document.getElementById('page-prev-btn');
                        if (prevBtn) prevBtn.disabled = true;
                        const nextBtn = document.getElementById('page-next-btn');
                        if (nextBtn) nextBtn.disabled = true;
                        const countBadge = document.getElementById('active-search-count-badge');
                        if (countBadge && activeContentSearchQuery) {
                            countBadge.innerText = '0 results found';
                        }
                        return;
                    }
                    cachedContentMap = {};
                    let rows = '';

                    // Client-side sort by Source and Post Date using friendly publisher names
                    if (res.data && res.data.length > 0) {
                        res.data.sort((a, b) => {
                            const srcA = getSourceHostName(a.sourceUrl || a.source_url || '').toLowerCase();
                            const srcB = getSourceHostName(b.sourceUrl || b.source_url || '').toLowerCase();
                            const dateA = new Date(a.createdAt || a.created_at || 0).getTime();
                            const dateB = new Date(b.createdAt || b.created_at || 0).getTime();
                            if (activeModSort === 'source_date' || activeModSort === 'source_asc') {
                                const cmp = srcA.localeCompare(srcB);
                                if (cmp !== 0) return cmp;
                                return dateB - dateA;
                            } else if (activeModSort === 'source_desc') {
                                const cmp = srcB.localeCompare(srcA);
                                if (cmp !== 0) return cmp;
                                return dateB - dateA;
                            } else if (activeModSort === 'date_asc') {
                                return dateA - dateB;
                            } else if (activeModSort === 'viral') {
                                if (a.isViral !== b.isViral) return b.isViral ? 1 : -1;
                                return dateB - dateA;
                            } else {
                                return dateB - dateA;
                            }
                        });
                    }

                    if (activeModGroupBy === 'none') {
                        // Standard flat view
                        res.data.forEach(item => {
                            rows += renderContentTableRow(item);
                        });
                    } else {
                        // Grouped view by Source, Language, District, or Category
                        const groups = {};
                        res.data.forEach(item => {
                            let gKey = '';
                            const cleanT = decodeHTMLEntities(item.title);
                            const cleanD = decodeHTMLEntities(item.description || '');
                            if (activeModGroupBy === 'source') {
                                const src = getSourceHostName(item.sourceUrl || item.source_url || '');
                                gKey = '📰 ' + (src || 'Direct Web');
                            } else if (activeModGroupBy === 'language') {
                                const isTa = item.language === 'ta' || (/[\u0B80-\u0BFF]/.test(cleanT + ' ' + cleanD));
                                gKey = isTa ? '🇮🇳 தமிழ் (Tamil Stories)' : '🇬🇧 English Stories';
                            } else if (activeModGroupBy === 'district') {
                                gKey = '📍 ' + (item.district || 'Tamil Nadu (Statewide)');
                            } else if (activeModGroupBy === 'category') {
                                const cat = item.category || 'News';
                                let icon = '🏷️';
                                const cLower = cat.toLowerCase();
                                if (cLower.includes('sport')) icon = '⚽';
                                else if (cLower.includes('politic')) icon = '🏛️';
                                else if (cLower.includes('business') || cLower.includes('finance')) icon = '💼';
                                else if (cLower.includes('tech')) icon = '💻';
                                else if (cLower.includes('entertain') || cLower.includes('cinema')) icon = '🎬';
                                else if (cLower.includes('crime')) icon = '⚖️';
                                else if (cLower.includes('civic')) icon = '🏙️';
                                gKey = icon + ' ' + cat;
                            }
                            if (!groups[gKey]) groups[gKey] = [];
                            groups[gKey].push(item);
                        });

                        Object.keys(groups).sort().forEach(gKey => {
                            const groupList = groups[gKey];
                            rows += '<tr style="background: linear-gradient(90deg, #1e293b, #0f172a); border-top: 2px solid #38bdf8; border-bottom: 1px solid #334155;">' +
                                '<td colspan="6" style="padding: 10px 16px;">' +
                                    '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                                        '<div style="font-weight: 800; font-size: 13px; color: #38bdf8; display: flex; align-items: center; gap: 8px;">' +
                                            '<span>' + gKey + '</span>' +
                                            '<span style="background: rgba(56,189,248,0.2); color: #38bdf8; font-size: 11px; padding: 2px 8px; border-radius: 4px; font-weight: 700;">' + groupList.length + ' posts</span>' +
                                        '</div>' +
                                        '<div style="font-size: 11px; color: #94a3b8;">' +
                                            '<button type="button" class="action-btn" style="font-size: 11px; padding: 2px 8px;" onclick="selectGroupItems(\'' + encodeURIComponent(gKey) + '\')">Select All in Group</button>' +
                                        '</div>' +
                                    '</div>' +
                                '</td>' +
                            '</tr>';
                            groupList.forEach(item => {
                                rows += renderContentTableRow(item, gKey);
                            });
                        });
                    }

                    tbody.innerHTML = rows;
                    clearSelection();

                    if (res.pagination) {
                        currentContentPage = res.pagination.page;
                        currentContentTotalPages = res.pagination.totalPages;
                        currentContentTotalCount = res.pagination.totalCount;

                        const countBadge = document.getElementById('active-search-count-badge');
                        if (countBadge && activeContentSearchQuery) {
                            countBadge.innerText = currentContentTotalCount + ' results found';
                        }

                        const start = currentContentTotalCount === 0 ? 0 : (currentContentPage - 1) * currentContentLimit + 1;
                        const end = Math.min(currentContentPage * currentContentLimit, currentContentTotalCount);
                        const pagInfo = document.getElementById('pagination-info');
                        if (pagInfo) pagInfo.innerText = 'Showing ' + start + ' - ' + end + ' of ' + currentContentTotalCount + ' items';
                        const pagDisp = document.getElementById('page-current-display');
                        if (pagDisp) pagDisp.innerText = 'Page ' + currentContentPage + ' of ' + currentContentTotalPages;
                        const pageInput = document.getElementById('jump-to-page-input');
                        if (pageInput) {
                            pageInput.value = currentContentPage;
                            pageInput.max = currentContentTotalPages || 1;
                        }
                        const totalEl = document.getElementById('page-total-count');
                        if (totalEl) totalEl.innerText = currentContentTotalPages || 1;

                        const firstBtn = document.getElementById('page-first-btn');
                        if (firstBtn) firstBtn.disabled = (currentContentPage <= 1);
                        const prevBtn = document.getElementById('page-prev-btn');
                        if (prevBtn) prevBtn.disabled = (currentContentPage <= 1);
                        const nextBtn = document.getElementById('page-next-btn');
                        if (nextBtn) nextBtn.disabled = (currentContentPage >= currentContentTotalPages);
                        const lastBtn = document.getElementById('page-last-btn');
                        if (lastBtn) lastBtn.disabled = (currentContentPage >= currentContentTotalPages);
                    }
                })
                .catch(() => {});
        }

        function changePage(delta) {
            const nextP = currentContentPage + delta;
            if (nextP >= 1 && nextP <= currentContentTotalPages) {
                currentContentPage = nextP;
                fetchPendingContent();
            }
        }

        function goToSpecificPage(target) {
            let p = parseInt(target, 10);
            if (isNaN(p)) return;
            if (p < 1) p = 1;
            if (currentContentTotalPages > 0 && p > currentContentTotalPages) p = currentContentTotalPages;
            if (p !== currentContentPage) {
                currentContentPage = p;
                fetchPendingContent();
            } else {
                const pageInput = document.getElementById('jump-to-page-input');
                if (pageInput) pageInput.value = currentContentPage;
            }
        }

        function changePageLimit(limitVal) {
            currentContentLimit = parseInt(limitVal, 10) || 15;
            currentContentPage = 1;
            fetchPendingContent();
        }

        const tnDistrictsList = [
            "Tamil Nadu", "Chennai", "Madurai", "Coimbatore", "Tiruchirappalli",
            "Salem", "Tirunelveli", "Thanjavur", "Vellore", "Dindigul", "Erode",
            "Tiruppur", "Thoothukudi", "Kanyakumari", "Cuddalore", "Kanchipuram",
            "Chengalpattu", "Tiruvallur", "Pudukkottai", "Sivaganga", "Ramanathapuram",
            "Nagapattinam", "Ranipet", "Tirupattur", "Tiruvannamalai", "Tiruvarur",
            "Kallakurichi", "Dharmapuri", "Krishnagiri", "Theni", "Tenkasi",
            "Namakkal", "Nilgiris", "Ariyalur", "Perambalur", "Viluppuram",
            "Virudhunagar", "Mayiladuthurai", "Karur", "National", "International"
        ];

        function openManualContentModal() {
            const distSelect = document.getElementById('manual-district');
            if (distSelect) {
                distSelect.innerHTML = tnDistrictsList.map(d => '<option value="' + d + '">' + d + '</option>').join('');
            }
            document.getElementById('manual-title').value = '';
            document.getElementById('manual-description').value = '';
            document.getElementById('manual-image-url').value = '';
            document.getElementById('manual-video-url').value = '';
            document.getElementById('manual-source-url').value = 'Manual Editorial Desk';
            document.getElementById('manualContentModal').style.display = 'flex';
        }

        function closeManualContentModal() {
            document.getElementById('manualContentModal').style.display = 'none';
        }

        function submitManualContent(publishLive) {
            const title = document.getElementById('manual-title').value.trim();
            const desc = document.getElementById('manual-description').value.trim();
            if (!title) {
                showToast('Please enter a headline');
                return;
            }
            if (!desc) {
                showToast('Please enter story body content');
                return;
            }

            const payload = {
                title: title,
                description: desc,
                contentType: document.getElementById('manual-content-type').value,
                category: document.getElementById('manual-category').value,
                district: document.getElementById('manual-district').value,
                language: (document.getElementById('manual-language') ? document.getElementById('manual-language').value : 'ta'),
                sourceUrl: document.getElementById('manual-source-url').value.trim(),
                imageUrl: document.getElementById('manual-image-url').value.trim(),
                videoUrl: document.getElementById('manual-video-url').value.trim(),
                publishLive: publishLive
            };

            showToast(publishLive ? 'Publishing live...' : 'Saving staged draft...');
            fetch('/admin/api/content/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + (publishLive ? 'Content published live to consumer app!' : 'Content saved to review queue!'));
                    closeManualContentModal();
                    fetchPendingContent();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to save content'));
                }
            })
            .catch(() => showToast('✗ Network error saving content'));
        }

        function toggleSelectAll(checked) {
            document.querySelectorAll('.content-row-checkbox').forEach(cb => {
                cb.checked = checked;
            });
            updateSelectionState();
        }

        function updateSelectionState() {
            const checkedBoxes = document.querySelectorAll('.content-row-checkbox:checked');
            const totalBoxes = document.querySelectorAll('.content-row-checkbox');
            const bar = document.getElementById('bulk-action-bar');
            const badge = document.getElementById('selected-count-badge');
            const selectAllBox = document.getElementById('select-all-checkbox');

            if (selectAllBox && totalBoxes.length > 0) {
                selectAllBox.checked = checkedBoxes.length === totalBoxes.length;
            }

            if (checkedBoxes.length > 0) {
                if (bar) bar.style.display = 'flex';
                if (badge) badge.innerText = checkedBoxes.length + ' item' + (checkedBoxes.length > 1 ? 's' : '') + ' selected';
            } else {
                if (bar) bar.style.display = 'none';
            }
        }

        function clearSelection() {
            document.querySelectorAll('.content-row-checkbox').forEach(cb => { cb.checked = false; });
            const selectAllBox = document.getElementById('select-all-checkbox');
            if (selectAllBox) selectAllBox.checked = false;
            const bar = document.getElementById('bulk-action-bar');
            if (bar) bar.style.display = 'none';
        }

        function getSelectedContentIds() {
            const ids = [];
            document.querySelectorAll('.content-row-checkbox:checked').forEach(cb => {
                ids.push(cb.value);
            });
            return ids;
        }

        async function approveSelectedContent() {
            const ids = getSelectedContentIds();
            if (ids.length === 0) {
                showToast('No items selected');
                return;
            }
            const confirmed = await showCustomConfirm({
                title: 'Approve & Publish (' + ids.length + ') Items',
                message: 'Are you sure you want to approve and publish ' + ids.length + ' selected posts live to public news portal and app users?',
                type: 'approve',
                confirmText: '✓ Approve (' + ids.length + ')'
            });
            if (!confirmed) return;
            showToast('Publishing ' + ids.length + ' selected items...');
            fetch('/admin/api/scraper/approve-batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentIds: ids })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    clearSelection();
                    fetchPendingContent();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Batch approval failed'));
                }
            })
            .catch(() => showToast('✗ Failed to approve selected items'));
        }

        async function rejectSelectedContent() {
            const ids = getSelectedContentIds();
            if (ids.length === 0) {
                showToast('No items selected');
                return;
            }
            const confirmed = await showCustomConfirm({
                title: 'Reject Selected Content (' + ids.length + ')',
                message: 'Are you sure you want to reject ' + ids.length + ' selected posts? They will be moved into the Rejected queue.',
                type: 'reject',
                confirmText: '🚫 Reject (' + ids.length + ')'
            });
            if (!confirmed) return;
            showToast('Rejecting ' + ids.length + ' selected items...');
            fetch('/admin/api/scraper/reject-batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentIds: ids })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    clearSelection();
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Batch reject failed'));
                }
            })
            .catch(() => showToast('✗ Failed to reject selected items'));
        }

        async function deleteSelectedContent() {
            const ids = getSelectedContentIds();
            if (ids.length === 0) {
                showToast('No items selected');
                return;
            }
            const confirmed = await showCustomConfirm({
                title: 'Permanently Delete (' + ids.length + ') Items',
                message: 'Are you sure you want to permanently delete ' + ids.length + ' selected items and their media from the database? This action cannot be undone.',
                type: 'delete',
                confirmText: '🗑️ Delete (' + ids.length + ')'
            });
            if (!confirmed) return;
            showToast('Deleting ' + ids.length + ' items permanently...');
            fetch('/admin/api/scraper/delete-batch', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentIds: ids })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    clearSelection();
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Batch delete failed'));
                }
            })
            .catch(() => showToast('✗ Failed to delete selected items'));
        }

        function setViewDeviceMode(mode) {
            const frame = document.getElementById('consumer-device-frame');
            const webBtn = document.getElementById('view-mode-web');
            const mobBtn = document.getElementById('view-mode-mobile');
            if (!frame) return;

            if (mode === 'mobile') {
                frame.style.maxWidth = '390px';
                frame.style.borderRadius = '32px';
                frame.style.border = '8px solid #334155';
                mobBtn.style.background = 'rgba(56,189,248,0.2)';
                mobBtn.style.color = '#38bdf8';
                webBtn.style.background = 'transparent';
                webBtn.style.color = 'var(--text-muted)';
            } else {
                frame.style.maxWidth = '680px';
                frame.style.borderRadius = '12px';
                frame.style.border = '1px solid #1e293b';
                webBtn.style.background = 'rgba(56,189,248,0.2)';
                webBtn.style.color = '#38bdf8';
                mobBtn.style.background = 'transparent';
                mobBtn.style.color = 'var(--text-muted)';
            }
        }

        let activeViewItemId = null;

        function openContentViewModal(itemId) {
            const item = cachedContentMap[itemId];
            if (!item) return;
            activeViewItemId = itemId;

            // Populate District Dropdown in Modal Top Bar
            const distSelect = document.getElementById('modal-district-select');
            if (distSelect) {
                distSelect.innerHTML = tnDistrictsList.map(d => {
                    const isSel = (d.toLowerCase() === (item.district || '').toLowerCase()) ? 'selected' : '';
                    return '<option value="' + d + '" ' + isSel + '>' + d + '</option>';
                }).join('');
            }

            // Viral toggle button in modal top bar
            const viralBtn = document.getElementById('modal-viral-toggle-btn');
            if (viralBtn) {
                if (item.isViral) {
                    viralBtn.style.background = 'rgba(249,115,22,0.25)';
                    viralBtn.style.color = '#f97316';
                    viralBtn.style.borderColor = '#f97316';
                    viralBtn.innerHTML = '🔥 Viral Priority: ON';
                } else {
                    viralBtn.style.background = 'transparent';
                    viralBtn.style.color = 'var(--text-muted)';
                    viralBtn.style.borderColor = 'var(--border)';
                    viralBtn.innerHTML = '⚡ Mark Viral';
                }
            }

            const bannerBtn = document.getElementById('modal-banner-toggle-btn');
            if (bannerBtn) {
                if (item.isMainBanner) {
                    bannerBtn.style.background = 'rgba(234,179,8,0.25)';
                    bannerBtn.style.color = '#eab308';
                    bannerBtn.style.borderColor = '#eab308';
                    bannerBtn.innerHTML = '⭐ Hero Banner: ON';
                } else {
                    bannerBtn.style.background = 'transparent';
                    bannerBtn.style.color = 'var(--text-muted)';
                    bannerBtn.style.borderColor = 'var(--border)';
                    bannerBtn.innerHTML = '☆ Set Hero';
                }
            }

            // Title & Content
            document.getElementById('consumer-title').innerText = item.title || 'Untitled Article';
            document.getElementById('consumer-body').innerText = item.description || 'No additional article text provided.';

            // Badges
            document.getElementById('consumer-district-badge').innerText = '📍 ' + (item.district || 'Tamil Nadu');
            document.getElementById('consumer-category-badge').innerText = item.category || 'News';
            document.getElementById('consumer-type-badge').innerText = item.contentType || 'TEXT STORY';
            
            const statusPill = document.getElementById('consumer-status-pill');
            if (item.status === 'PUBLISHED') {
                statusPill.className = 'badge-pill badge-success';
                statusPill.innerText = 'LIVE ON APP';
            } else {
                statusPill.className = 'badge-pill badge-warning';
                statusPill.innerText = 'STAGED PENDING';
            }

            document.getElementById('consumer-time-badge').innerText = '🕒 ' + formatSourceDate(item.createdAt);

            // Media: display BOTH video and image if available
            const mediaBox = document.getElementById('consumer-media-box');
            const fallbackThumb = getFallbackDistrictImage(item.district, item.category, item.title);
            const thumbUrl = item.thumbnail || fallbackThumb;

            let mediaContentHtml = '';
            const vUrl = item.videoUrl || '';
            if (vUrl && (vUrl.includes('bbc.com/ws/av-embeds') || vUrl.includes('/embed/') || (item.videoType === 'bbc'))) {
                const src = vUrl.includes('bbc.com/ws/av-embeds') ? vUrl : ('https://www.bbc.com/ws/av-embeds/articles/' + item.videoId + '/ta');
                mediaContentHtml += '<div style="margin-bottom: 14px; border-radius: 8px; overflow: hidden; background: #000;">' +
                    '<iframe style="width:100%; height:350px; border:0; display:block;" src="' + src + '" allow="autoplay; fullscreen; encrypted-media" allowfullscreen></iframe>' +
                '</div>';
            } else if (vUrl && (vUrl.endsWith('.mp4') || vUrl.endsWith('.webm') || vUrl.includes('.mp4?') || vUrl.includes('.webm?'))) {
                mediaContentHtml += '<div style="margin-bottom: 14px; border-radius: 8px; overflow: hidden; background: #000;">' +
                    '<video controls style="width:100%; max-height:350px; display:block;" src="' + vUrl + '"></video>' +
                '</div>';
            } else if (item.videoId && !item.videoId.startsWith('vid_') && !item.videoId.startsWith('p0')) {
                mediaContentHtml += '<div style="margin-bottom: 14px; border-radius: 8px; overflow: hidden; background: #000;">' +
                    '<iframe style="width:100%; height:330px; border:0; display:block;" src="https://www.youtube.com/embed/' + item.videoId + '?autoplay=0" allowfullscreen></iframe>' +
                '</div>';
            }
            if (thumbUrl) {
                const mapFallbackUrl = '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
                mediaContentHtml += '<div style="border-radius: 8px; overflow: hidden; max-height: 380px; background: #000; text-align: center;">' +
                    '<img src="' + thumbUrl + '" onerror="this.onerror=null; this.src=\'' + mapFallbackUrl + '\';" style="width:100%; max-height:380px; object-fit:cover; display:block;" alt="' + (item.district || 'Tamil Nadu') + ' Map" />' +
                '</div>';
            }

            mediaBox.style.display = 'block';
            mediaBox.innerHTML = mediaContentHtml;

            // Source Attribution
            document.getElementById('consumer-source-name').innerText = getSourceHostName(item.sourceUrl);
            const sourceEl = document.getElementById('consumer-source-link');
            if (item.sourceUrl) {
                sourceEl.style.display = 'inline-block';
                sourceEl.href = item.sourceUrl;
            } else {
                sourceEl.style.display = 'none';
            }

            // Buttons in top bar: Approve, Reject, Delete
            const approveBtn = document.getElementById('view-modal-approve-btn');
            if (approveBtn) {
                if (item.status === 'PUBLISHED') {
                    approveBtn.style.display = 'none';
                } else {
                    approveBtn.style.display = 'inline-block';
                    approveBtn.onclick = function() {
                        approveContent(item.id);
                        closeContentViewModal();
                    };
                }
            }

            const rejectBtn = document.getElementById('view-modal-reject-btn');
            if (rejectBtn) {
                if (item.status === 'REJECTED') {
                    rejectBtn.style.display = 'none';
                } else {
                    rejectBtn.style.display = 'inline-block';
                    rejectBtn.onclick = function() {
                        rejectContent(item.id);
                        closeContentViewModal();
                    };
                }
            }

            const deleteBtn = document.getElementById('view-modal-delete-btn');
            if (deleteBtn) {
                deleteBtn.style.display = 'inline-block';
                deleteBtn.onclick = function() {
                    deletePermanent(item.id);
                    closeContentViewModal();
                };
            }

            // Set default view mode to web
            setViewDeviceMode('web');
            document.getElementById('contentViewModal').style.display = 'flex';
        }

        function updateCurrentItemDistrict(newDist) {
            if (!activeViewItemId) return;
            showToast('Updating district to ' + newDist + '...');
            fetch('/admin/api/content/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: activeViewItemId, district: newDist })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ Updated to ' + newDist);
                if (cachedContentMap[activeViewItemId]) {
                    cachedContentMap[activeViewItemId].district = newDist;
                }
                const badge = document.getElementById('consumer-district-badge');
                if (badge) badge.innerText = '📍 ' + newDist;
                fetchPendingContent();
                fetchStats();
            })
            .catch(() => showToast('✗ Failed to update district'));
        }

        function toggleCurrentItemViral() {
            if (!activeViewItemId || !cachedContentMap[activeViewItemId]) return;
            const currentItem = cachedContentMap[activeViewItemId];
            const newViral = !currentItem.isViral;
            showToast(newViral ? 'Marking as Viral Priority...' : 'Removing viral priority...');
            fetch('/admin/api/content/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: activeViewItemId, isViral: newViral })
            })
            .then(res => res.json())
            .then(res => {
                showToast(newViral ? '🔥 Marked as Viral Priority' : '✓ Viral priority removed');
                currentItem.isViral = newViral;
                const viralBtn = document.getElementById('modal-viral-toggle-btn');
                if (viralBtn) {
                    if (newViral) {
                        viralBtn.style.background = 'rgba(249,115,22,0.25)';
                        viralBtn.style.color = '#f97316';
                        viralBtn.style.borderColor = '#f97316';
                        viralBtn.innerHTML = '🔥 Viral Priority: ON';
                    } else {
                        viralBtn.style.background = 'transparent';
                        viralBtn.style.color = 'var(--text-muted)';
                        viralBtn.style.borderColor = 'var(--border)';
                        viralBtn.innerHTML = '⚡ Mark Viral';
                    }
                }
                fetchPendingContent();
            })
            .catch(() => showToast('✗ Failed to update viral priority'));
        }

        async function triggerReclassifyAll() {
            const confirmed = await showCustomConfirm({
                title: 'Auto-Reclassify All Content',
                message: 'Re-scan and accurately classify all content items geographically across Tamil Nadu, National, and International?',
                type: 'warning',
                confirmText: '🔄 Reclassify All'
            });
            if (!confirmed) return;
            showToast('Reclassifying all content geographically...');
            fetch('/admin/api/content/reclassify-all', { method: 'POST' })
                .then(res => res.json())
                .then(res => {
                    showToast('✓ ' + res.message);
                    fetchPendingContent();
                    fetchStats();
                })
                .catch(() => showToast('✗ Reclassification failed'));
        }

        async function triggerRefetchAllText() {
            const confirmed = await showCustomConfirm({
                title: 'Refetch Full Text Paragraphs',
                message: 'Visit original news source webpages to extract full unabridged multi-paragraph text and photos for all articles?',
                type: 'warning',
                confirmText: '📖 Refetch Text'
            });
            if (!confirmed) return;
            showToast('Fetching unabridged text from source webpages...');
            fetch('/admin/api/content/refetch-text', { method: 'POST' })
                .then(res => res.json())
                .then(res => {
                    showToast('✓ ' + res.message);
                    fetchPendingContent();
                    fetchStats();
                })
                .catch(() => showToast('✗ Failed to refetch text from sources'));
        }

        function closeContentViewModal() {
            const mediaBox = document.getElementById('consumer-media-box');
            if (mediaBox) mediaBox.innerHTML = '';
            document.getElementById('contentViewModal').style.display = 'none';
        }

        function triggerSiteScrape() {
            const input = document.getElementById('scrape-url-input');
            const url = input.value.trim();
            if (!url) {
                showToast('Please enter a website or feed URL to scrape');
                return;
            }

            const btn = document.getElementById('btn-scrape-url');
            btn.disabled = true;
            btn.innerText = 'Scraping...';
            showToast('🔍 Scraping ' + url + '...');

            fetch('/admin/api/scraper/scrape', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url: url })
            })
            .then(res => res.json())
            .then(res => {
                btn.disabled = false;
                btn.innerText = '⚡ Scrape & Stage for Review';
                if (res.success) {
                    showToast('✓ ' + res.message);
                    input.value = '';
                    fetchPendingContent();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Scraping failed'));
                }
            })
            .catch(() => {
                btn.disabled = false;
                btn.innerText = '⚡ Scrape & Stage for Review';
                showToast('✗ Network error executing scraper');
            });
        }

        function approveContent(contentId) {
            showToast('Publishing content...');
            fetch('/admin/api/scraper/approve', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentId: contentId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Content published live to app users!');
                    fetchPendingContent();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Approval failed'));
                }
            })
            .catch(() => showToast('✗ Failed to approve content'));
        }

        async function rejectContent(contentId) {
            const confirmed = await showCustomConfirm({
                title: 'Reject News Article',
                message: 'Are you sure you want to reject this scraped item? It will be removed from the review queue and archived in the Rejected queue.',
                type: 'reject',
                confirmText: '🚫 Reject Article'
            });
            if (!confirmed) return;
            fetch('/admin/api/scraper/reject', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentId: contentId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('Rejected item moved to Rejected queue');
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ Failed to reject item');
                }
            })
            .catch(() => showToast('✗ Failed to reject item'));
        }

        function restoreContent(contentId) {
            fetch('/admin/api/scraper/restore', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentId: contentId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Restored item back to Pending queue');
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to restore item'));
                }
            })
            .catch(() => showToast('✗ Failed to restore item'));
        }

        async function deletePermanent(contentId) {
            const confirmed = await showCustomConfirm({
                title: 'Permanently Delete Item',
                message: 'Are you sure you want to permanently delete this content item and its media from the database? This action cannot be undone.',
                type: 'delete',
                confirmText: '🗑️ Delete Forever'
            });
            if (!confirmed) return;
            fetch('/admin/api/scraper/delete-permanent', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ contentId: contentId })
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ Item permanently deleted');
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to delete item'));
                }
            })
            .catch(() => showToast('✗ Failed to delete item'));
        }

        async function emptyTrash() {
            const confirmed = await showCustomConfirm({
                title: 'Delete All Rejected Items',
                message: 'Permanently purge ALL items currently in the Rejected queue? This action cannot be undone.',
                type: 'delete',
                confirmText: '🗑️ Delete All Rejected'
            });
            if (!confirmed) return;
            showToast('Purging all rejected items...');
            fetch('/admin/api/scraper/empty-trash', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    fetchPendingContent();
                    updateDiscardedCountBadge();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Failed to empty rejected items'));
                }
            })
            .catch(() => showToast('✗ Failed to empty rejected items'));
        }

        function updateDiscardedCountBadge() {
            fetch('/admin/api/stats')
                .then(res => res.json())
                .then(res => {
                    if (res && res.data) {
                        const d = res.data;
                        const pBadge = document.getElementById('mod-tab-count-pending');
                        if (pBadge) pBadge.innerText = d.pendingModeration !== undefined ? d.pendingModeration : 0;
                        const pubBadge = document.getElementById('mod-tab-count-published');
                        if (pubBadge) pubBadge.innerText = d.publishedContent !== undefined ? d.publishedContent : 0;
                        const rBadge = document.getElementById('discarded-count-badge');
                        if (rBadge) rBadge.innerText = d.rejectedContent !== undefined ? d.rejectedContent : 0;
                        const aBadge = document.getElementById('mod-tab-count-all');
                        if (aBadge) aBadge.innerText = d.totalContent !== undefined ? d.totalContent : 0;
                    }
                })
                .catch(() => {});
        }

        async function approveAllPending() {
            const confirmed = await showCustomConfirm({
                title: 'Approve & Publish All Pending',
                message: 'Approve and publish all pending staged items live to public news portal and app users?',
                type: 'approve',
                confirmText: '✓ Approve All Live'
            });
            if (!confirmed) return;
            showToast('Approving all items...');
            fetch('/admin/api/scraper/approve-all', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    fetchPendingContent();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Batch approval failed'));
                }
            })
            .catch(() => showToast('✗ Failed batch approval'));
        }

        async function triggerDeduplicate() {
            const confirmed = await showCustomConfirm({
                title: '🧹 Auto-Deduplicate Content',
                message: 'Scan all news posts and videos, retain only the newest/latest version of every duplicate group, and permanently purge older duplicates?',
                type: 'danger',
                confirmText: '🧹 Deduplicate (Keep Latest)'
            });
            if (!confirmed) return;
            showToast('Scanning and deduplicating content...');
            fetch('/admin/api/scraper/deduplicate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            })
            .then(res => res.json())
            .then(res => {
                if (res.success) {
                    showToast('✓ ' + res.message);
                    fetchPendingContent();
                    fetchRetentionSettings();
                    fetchStats();
                } else {
                    showToast('✗ ' + (res.message || 'Deduplication failed'));
                }
            })
            .catch(() => showToast('✗ Failed to trigger deduplication'));
        }

        function fetchRetentionSettings() {
            fetch('/admin/api/retention')
                .then(res => res.json())
                .then(res => {
                    if (res.success && res.data) {
                        const d = res.data;
                        const hoursInput = document.getElementById('retention-hours-input');
                        if (hoursInput) hoursInput.value = d.retentionHours || 24;
                        const badge = document.getElementById('retention-current-badge');
                        if (badge) badge.innerText = (d.retentionHours || 24) + ' Hours';
                        const totalEl = document.getElementById('retention-total-count');
                        if (totalEl) totalEl.innerText = d.totalContentCount !== undefined ? d.totalContentCount : '--';
                        const staleEl = document.getElementById('retention-stale-count');
                        if (staleEl) staleEl.innerText = (d.staleContentCount !== undefined ? d.staleContentCount : 0) + ' items';
                        const autoCheck = document.getElementById('retention-auto-cleanup');
                        if (autoCheck) autoCheck.checked = (d.autoCleanupEnabled !== false);
                        const noteEl = document.getElementById('retention-last-cleanup-note');
                        if (noteEl) {
                            if (d.lastCleanupAt) {
                                const dt = new Date(d.lastCleanupAt);
                                noteEl.innerText = 'Last cleanup: ' + dt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) + ' (' + (d.lastCleanupCount || 0) + ' purged)';
                            } else {
                                noteEl.innerText = 'Last cleanup: Auto-scheduled';
                            }
                        }
                    }
                })
                .catch(() => {});
        }

        function setRetentionPreset(hours) {
            const input = document.getElementById('retention-hours-input');
            if (input) input.value = hours;
            saveRetentionSettings();
        }

        function saveRetentionSettings() {
            const hours = parseInt(document.getElementById('retention-hours-input').value, 10) || 24;
            const autoEnabled = document.getElementById('retention-auto-cleanup').checked;
            showToast('Saving retention policy...');
            fetch('/admin/api/retention', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ retentionHours: hours, autoCleanupEnabled: autoEnabled })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + res.message);
                fetchRetentionSettings();
            })
            .catch(() => showToast('✗ Failed to save retention policy'));
        }

        async function triggerRetentionCleanupNow() {
            const hours = parseInt(document.getElementById('retention-hours-input').value, 10) || 24;
            const confirmed = await showCustomConfirm({
                title: 'Purge Stale Content Now',
                message: 'Immediately purge all scraped content items older than ' + hours + ' hours across ALL statuses (Pending, Published, Rejected)? Manually created posts are protected and will NOT be deleted.',
                type: 'delete',
                confirmText: '🧹 Clean Stale Now'
            });
            if (!confirmed) return;
            showToast('Purging stale content older than ' + hours + 'h...');
            fetch('/admin/api/retention/cleanup-now', { method: 'POST' })
                .then(res => res.json())
                .then(res => {
                    showToast('✓ ' + res.message);
                    fetchRetentionSettings();
                    fetchPendingContent();
                    fetchStats();
                })
                .catch(() => showToast('✗ Cleanup request failed'));
        }

        function loadBannerConfig() {
            fetch('/admin/api/banners/config')
                .then(res => res.json())
                .then(res => {
                    if (!res.success || !res.data) return;
                    const d = res.data;
                    // Main Hero
                    const mainEl = document.getElementById('banner-slot-main-display');
                    if (mainEl) {
                        if (d.portal_main_banner_id_item) {
                            const it = d.portal_main_banner_id_item;
                            mainEl.innerHTML = '<img src="' + (it.thumbnail || '/admin/api/maps/svg?district=Tamil%20Nadu') + '" style="width:40px; height:28px; object-fit:cover; border-radius:4px;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#fff; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">' + escapeHtml(it.title) + '</div>' +
                                '<div style="font-size:10px; color:#eab308;">📍 ' + escapeHtml(it.district) + ' &bull; Active Portal Hero</div></div>';
                        } else {
                            mainEl.innerHTML = '<span style="color:var(--text-muted); font-size:11px;">⚡ Auto: Top Viral Story</span>';
                        }
                    }
                    // Header Leaderboard
                    const hEl = document.getElementById('banner-slot-header-display');
                    if (hEl) {
                        if (d.portal_ad_banner_header_id_item) {
                            const it = d.portal_ad_banner_header_id_item;
                            hEl.innerHTML = '<img src="' + (it.thumbnail || '/admin/api/maps/svg?district=Tamil%20Nadu') + '" style="width:40px; height:28px; object-fit:cover; border-radius:4px;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#fff; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">' + escapeHtml(it.title) + '</div>' +
                                '<div style="font-size:10px; color:#38bdf8;">Promoted Story Slot</div></div>';
                        } else {
                            hEl.innerHTML = '<img src="/portal/assets/brand/tn24-header.svg" style="width:52px; height:22px; object-fit:cover; border-radius:4px; border:1px dashed #38bdf8;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (728x90)</div>' +
                                '<div style="font-size:10px; color:var(--text-muted);">Active: Leaderboard Advertising Banner</div></div>';
                        }
                    }
                    // Sidebar
                    const sEl = document.getElementById('banner-slot-sidebar-display');
                    if (sEl) {
                        if (d.portal_ad_banner_sidebar_id_item) {
                            const it = d.portal_ad_banner_sidebar_id_item;
                            sEl.innerHTML = '<img src="' + (it.thumbnail || '/admin/api/maps/svg?district=Tamil%20Nadu') + '" style="width:40px; height:28px; object-fit:cover; border-radius:4px;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#fff; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">' + escapeHtml(it.title) + '</div>' +
                                '<div style="font-size:10px; color:#38bdf8;">Promoted Story Slot</div></div>';
                        } else {
                            sEl.innerHTML = '<img src="/portal/assets/brand/tn24-sidebar.svg" style="width:26px; height:26px; object-fit:cover; border-radius:4px; border:1px dashed #38bdf8;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (250x250)</div>' +
                                '<div style="font-size:10px; color:var(--text-muted);">Active: Sidebar 250x250 Advertising Banner</div></div>';
                        }
                    }
                    // In-Feed / Square
                    const infEl = document.getElementById('banner-slot-infeed-display');
                    if (infEl) {
                        if (d.portal_ad_banner_infeed_id_item) {
                            const it = d.portal_ad_banner_infeed_id_item;
                            infEl.innerHTML = '<img src="' + (it.thumbnail || '/admin/api/maps/svg?district=Tamil%20Nadu') + '" style="width:40px; height:28px; object-fit:cover; border-radius:4px;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#fff; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">' + escapeHtml(it.title) + '</div>' +
                                '<div style="font-size:10px; color:#38bdf8;">Promoted Story Slot</div></div>';
                        } else {
                            infEl.innerHTML = '<img src="/portal/assets/brand/tn24-infeed.svg" style="width:52px; height:22px; object-fit:cover; border-radius:4px; border:1px dashed #a855f7;" />' +
                                '<div style="flex:1; overflow:hidden;"><div style="font-weight:700; color:#38bdf8; font-size:11px;">📢 Contact for Advertisement (In-Feed &amp; Square)</div>' +
                                '<div style="font-size:10px; color:var(--text-muted);">Active: In-Feed 728x90 &amp; Square 200x200 Advertising</div></div>';
                        }
                    }
                })
                .catch(() => {});
        }

        function resetMainHeroBanner() {
            showToast('Resetting main banner to auto...');
            fetch('/admin/api/banners/set-main', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: '' })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + res.message);
                loadBannerConfig();
                fetchPendingContent();
            })
            .catch(() => showToast('✗ Failed to reset main banner'));
        }

        function resetAdSlot(slot) {
            showToast('Resetting ' + slot + ' slot to dummy placeholder...');
            fetch('/admin/api/banners/config', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ resetSlot: slot })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + (res.message || 'Reset to dummy placeholder'));
                loadBannerConfig();
                fetchPendingContent();
            })
            .catch(() => showToast('✗ Failed to reset slot'));
        }

        function quickToggleMainBanner(id, enable) {
            showToast(enable ? 'Setting as Portal Main Hero Banner...' : 'Resetting main hero banner...');
            fetch('/admin/api/banners/set-main', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ id: enable ? id : '' })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + res.message);
                loadBannerConfig();
                fetchPendingContent();
            })
            .catch(() => showToast('✗ Failed to update main banner'));
        }

        function toggleCurrentItemMainBanner() {
            if (!activeViewItemId || !cachedContentMap[activeViewItemId]) return;
            const item = cachedContentMap[activeViewItemId];
            const newBanner = !item.isMainBanner;
            quickToggleMainBanner(item.id, newBanner);
            item.isMainBanner = newBanner;
            const btn = document.getElementById('modal-banner-toggle-btn');
            if (btn) {
                if (newBanner) {
                    btn.style.background = 'rgba(234,179,8,0.25)';
                    btn.style.color = '#eab308';
                    btn.style.borderColor = '#eab308';
                    btn.innerHTML = '⭐ Hero Banner: ON';
                } else {
                    btn.style.background = 'transparent';
                    btn.style.color = 'var(--text-muted)';
                    btn.style.borderColor = 'var(--border)';
                    btn.innerHTML = '☆ Set Hero';
                }
            }
        }

        function openEditContentModal(itemId) {
            const item = cachedContentMap[itemId];
            if (!item) return;

            document.getElementById('edit-post-id').value = item.id;
            document.getElementById('edit-post-title').value = item.title || '';
            document.getElementById('edit-post-description').value = item.description || '';
            document.getElementById('edit-post-thumbnail').value = item.thumbnail || '';

            const srcName = getSourceHostName(item.sourceUrl);
            const dateStr = formatSourceDate(item.createdAt);
            const srcInfoEl = document.getElementById('edit-post-source-info');
            if (srcInfoEl) {
                srcInfoEl.innerHTML = item.sourceUrl ? '🔗 Source: <a href="' + item.sourceUrl + '" target="_blank" style="color: #38bdf8; text-decoration: none;">' + escapeHtml(srcName) + '</a>' : '🔗 Source: Direct Entry';
            }
            const dateInfoEl = document.getElementById('edit-post-date-info');
            if (dateInfoEl) {
                dateInfoEl.textContent = '🕒 Published: ' + dateStr;
            }
            
            const thumbPreview = document.getElementById('edit-post-thumb-preview');
            if (thumbPreview) {
                thumbPreview.src = item.thumbnail || '/admin/api/maps/svg?district=' + encodeURIComponent(item.district || 'Tamil Nadu');
            }

            // District select
            const distSel = document.getElementById('edit-post-district');
            if (distSel) {
                distSel.innerHTML = tnDistrictsList.map(d => {
                    const isSel = (d.toLowerCase() === (item.district || '').toLowerCase()) ? 'selected' : '';
                    return '<option value="' + d + '" ' + isSel + '>' + d + '</option>';
                }).join('');
            }

            // Category select
            const catSel = document.getElementById('edit-post-category');
            if (catSel && item.category) {
                for (let i = 0; i < catSel.options.length; i++) {
                    if (catSel.options[i].value.toLowerCase() === item.category.toLowerCase()) {
                        catSel.selectedIndex = i;
                        break;
                    }
                }
            }

            // Language select
            const langSel = document.getElementById('edit-post-language');
            if (langSel) {
                const curLang = (item.language || 'ta').toLowerCase();
                langSel.value = curLang;
            }

            document.getElementById('edit-post-is-viral').checked = !!item.isViral;
            document.getElementById('edit-post-is-main-banner').checked = !!item.isMainBanner;
            document.getElementById('edit-post-ad-slot').value = 'none';

            document.getElementById('editContentModal').style.display = 'flex';
        }

        function closeEditContentModal() {
            document.getElementById('editContentModal').style.display = 'none';
        }

        function updateEditThumbnailPreview() {
            const val = document.getElementById('edit-post-thumbnail').value.trim();
            const preview = document.getElementById('edit-post-thumb-preview');
            if (preview) {
                preview.src = val || '/admin/api/maps/svg?district=Tamil%20Nadu';
            }
        }

        function saveEditContentForm() {
            const id = document.getElementById('edit-post-id').value;
            const title = document.getElementById('edit-post-title').value.trim();
            const desc = document.getElementById('edit-post-description').value.trim();
            const thumb = document.getElementById('edit-post-thumbnail').value.trim();
            const dist = document.getElementById('edit-post-district').value;
            const cat = document.getElementById('edit-post-category').value;
            const lang = document.getElementById('edit-post-language') ? document.getElementById('edit-post-language').value : 'ta';
            const isViral = document.getElementById('edit-post-is-viral').checked;
            const setMain = document.getElementById('edit-post-is-main-banner').checked;
            const adSlot = document.getElementById('edit-post-ad-slot').value;

            if (!title) {
                showToast('⚠️ Headline title cannot be empty');
                return;
            }

            showToast('Saving post updates and reflecting live...');
            fetch('/admin/api/content/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    id: id,
                    title: title,
                    description: desc,
                    thumbnail: thumb,
                    district: dist,
                    category: cat,
                    language: lang,
                    isViral: isViral,
                    setMainBanner: setMain,
                    adSlot: adSlot
                })
            })
            .then(res => res.json())
            .then(res => {
                showToast('✓ ' + res.message);
                closeEditContentModal();

                if (cachedContentMap[id]) {
                    cachedContentMap[id].title = title;
                    cachedContentMap[id].description = desc;
                    cachedContentMap[id].thumbnail = thumb;
                    cachedContentMap[id].district = dist;
                    cachedContentMap[id].category = cat;
                    cachedContentMap[id].language = lang;
                    cachedContentMap[id].isViral = isViral;
                    cachedContentMap[id].isMainBanner = setMain;
                }

                // If viewer is open on this item, update its fields
                if (activeViewItemId === id) {
                    const cTitle = document.getElementById('consumer-title');
                    if (cTitle) cTitle.innerText = title;
                    const cBody = document.getElementById('consumer-body');
                    if (cBody) cBody.innerText = desc;
                    const cDist = document.getElementById('consumer-district-badge');
                    if (cDist) cDist.innerText = '📍 ' + dist;
                    const cCat = document.getElementById('consumer-category-badge');
                    if (cCat) cCat.innerText = cat;
                }

                fetchPendingContent();
                loadBannerConfig();
                fetchStats();
            })
            .catch(() => showToast('✗ Failed to update post'));
        }

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                closeCustomConfirm(false);
                closeContentViewModal();
                closeEditContentModal();
                closeManualContentModal();
                closeLanguageConfigModal();
                closeSourcesModal();
            }
        });

        // Initial Real Data Fetch
        fetchStats();
        fetchRetentionSettings();
        loadBannerConfig();
        loadLanguageConfig();
        updateDiscardedCountBadge();
        if (window.location.hash) {
            const h = window.location.hash.substring(1);
            if (['dashboard', 'moderation', 'agent', 'cron', 'audit', 'contributors'].includes(h)) {
                switchTab(h);
            }
        }
    </script>
</body>
</html>`
}
