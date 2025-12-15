# Quick Start Guide

## 🚀 Running the LinkedIn Automation Bot

### Prerequisites Installed ✅
- Go ≥ 1.21 ✅
- Project built ✅
- `.env` configured ✅

### 🔧 One-Time Setup Required

The bot needs Chrome/Chromium system dependencies. Run:

```bash
./setup.sh
```

Or manually install:
```bash
sudo apt-get install -y libnss3 libnspr4 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 libxkbcommon0 libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libgbm1 libasound2
```

### ▶️ Run the Bot

```bash
./bin/linkedin-bot
```

## What Happens

1. **Initialization** (~5 seconds)
   - Loads configuration from `.env`
   - Downloads/locates Chromium (first run only)
   - Launches browser with stealth techniques

2. **Login** (~10 seconds)
   - Attempts cookie restoration
   - Falls back to credentials if needed
   - Detects security challenges (CAPTCHA/OTP)

3. **Automation Flow** (example enabled)
   - Searches for "Software Engineer" profiles
   - Visits profiles naturally
   - Sends connection requests with personalized notes
   - Respects rate limits and business hours
   - Takes random breaks

4. **Graceful Shutdown**
   - Saves cookies
   - Logs statistics
   - Cleans up browser

## 📝 Configuration

Edit `.env` to customize:
- LinkedIn credentials
- Rate limits (daily connections, hourly messages)
- Business hours
- Search criteria
- Break frequency

## 🎮 Example Flow

The current implementation runs `runExampleAutomation()`:
- Searches for profiles
- Visits and extracts data
- Sends connection requests
- Follows up with messages after 1-3 days

To disable: Comment out line 161 in `internal/app/app.go`

## 📊 Logs

Check `logs/linkedin-bot-YYYY-MM-DD.log` for:
- Action history
- Error messages
- Rate limit events
- Security challenge detections

## ⚠️ Important

- First run downloads ~150MB Chromium
- Run during business hours configured in `.env`
- Respect LinkedIn's terms of service
- Monitor for security challenges

## 🛑 Stopping

Press `Ctrl+C` - the bot will shut down gracefully.
