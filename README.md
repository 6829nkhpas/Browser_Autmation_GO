# LinkedIn Human-Behavior Automation Framework

A production-grade Go automation framework that simulates authentic human behavior while avoiding bot-detection on LinkedIn.

## 🎯 Overview

This is **NOT** a scraper. This is a sophisticated behavioral browser automation system built with:
- **Clean Architecture** with strict separation of concerns
- **Advanced anti-detection** measures (8+ stealth techniques)
- **Human behavior simulation** (Bézier mouse, realistic typing, natural scrolling)
- **Ethical rate limiting** and business hours enforcement
- **Production-grade Go** with idiomatic patterns

## 🏗️ Architecture

```
linkedin-automation/
├── cmd/bot/              # Entry point
├── internal/
│   ├── app/              # Application lifecycle & logging
│   ├── config/           # Configuration management
│   ├── browser/          # Chrome + Rod setup
│   ├── stealth/          # Anti-detection techniques
│   ├── behavior/         # Human behavior engine
│   ├── auth/             # Authentication & cookies
│   ├── linkedin/         # UI selectors & constants
│   ├── scheduler/        # Rate limiting & business hours
│   └── store/            # Persistence layer
├── assets/templates/     # Message templates
└── data/                 # Cookies, state, logs
```

### Core Design Principle

**All user actions follow this pipeline:**
```
Intent → Behavior Simulation → UI Interaction → Result → State Update
```

**Rod is NEVER used directly** outside `browser/` and `behavior/` packages.

## ✨ Features

### Anti-Detection (Stealth Layer)
- ✅ Disable `navigator.webdriver`
- ✅ Randomized User-Agent (7+ realistic Chrome variants)
- ✅ Viewport randomization (common desktop resolutions)
- ✅ Language & timezone spoofing
- ✅ Permission API overrides
- ✅ Canvas fingerprint randomization
- ✅ WebGL fingerprint variation

### Human Behavior Simulation
- 🎯 **Bézier-curve mouse movement** with micro-corrections
- ⌨️ **Realistic typing**:
  - Variable speed (40-80 WPM)
  - Random typos with immediate corrections
  - Thinking pauses at punctuation
  - Keyboard proximity-based typo generation
- 📜 **Natural scrolling**:
  - Acceleration/deceleration curves
  - Reading pauses
  - Occasional scroll-back (re-reading)
- 🤔 **Cognitive delays**:
  - Thinking pauses (1-3s)
  - Reading time estimation
  - Pre-action preparation time

### Rate Limiting & Scheduling
- 📅 **Business hours only** (configurable timezone)
- 🔢 **Daily connection limits** (default: 30/day)
- 💬 **Hourly message limits** (default: 10/hour)
- ⏸️ **Automatic breaks** (5-20 minutes every 30-60 minutes)
- 🚫 **Graceful blocking** (never crashes on limits)

### Session Management
- 🍪 **Cookie persistence** across runs
- 🔐 **Automatic session restoration**
- 🛡️ **Security challenge detection** (CAPTCHA, OTP, checkpoints)
- ⚠️ **Graceful abort** on security challenges (no brute-force)

## 🚀 Quick Start

### Prerequisites
- Go ≥ 1.21
- Google Chrome

### Installation

```bash
# Clone the repository
git clone <repo-url>
cd GO_Browser_Automation

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### Configuration

Edit `.env`:

```env
# LinkedIn Credentials
LINKEDIN_EMAIL=your.email@example.com
LINKEDIN_PASSWORD=your_password_here

# Browser Settings
BROWSER_HEADLESS=false  # Set to true for headless mode
BROWSER_WIDTH=1366
BROWSER_HEIGHT=768

# Rate Limits
DAILY_CONNECTION_LIMIT=30
HOURLY_MESSAGE_LIMIT=10

# Business Hours (24-hour format)
BUSINESS_HOURS_START=9
BUSINESS_HOURS_END=17
BUSINESS_DAYS=Monday,Tuesday,Wednesday,Thursday,Friday

# Timezone
TIMEZONE=America/New_York
LANGUAGE=en-US
```

### Build & Run

```bash
# Build
go build -o bin/linkedin-bot cmd/bot/main.go

# Run
./bin/linkedin-bot
```

Or with go run:

```bash
go run cmd/bot/main.go
```

## 📚 Technical Details

### Stealth Techniques

#### 1. WebDriver Detection Bypass
Overrides `navigator.webdriver` and related automation properties.

#### 2. User-Agent Randomization
Rotates between 7 realistic Chrome user agents with version variations.

#### 3. Canvas Fingerprinting
Adds subtle, consistent noise to canvas rendering to avoid fingerprinting.

#### 4. WebGL Fingerprinting
Randomizes vendor/renderer strings within realistic bounds.

### Behavior Engine

All browser interactions go through the `BehaviorEngine`:

```go
type BehaviorEngine interface {
    Navigate(url string) error
    Click(selector string) error
    Type(selector, text string) error
    Scroll(pixels int) error
    Hover(selector string) error
    WaitHuman(minMs, maxMs int)
}
```

**Never call `page.Click()` directly!** Always use the behavior engine.

### Error Handling

- Every exported function returns explicit errors
- Errors wrapped with context using `fmt.Errorf`
- Structured logging for all events
- No silent failures
- No `panic()` except during startup

## 📊 Logging

Logs are written to `logs/linkedin-bot-YYYY-MM-DD.log`:

```
[INFO] 2024-01-15 10:30:45 Initializing LinkedIn Automation Framework...
[INFO] 2024-01-15 10:30:46 Launching Chrome browser...
[INFO] 2024-01-15 10:30:48 Applying stealth techniques...
[ACTION] [SUCCESS] LOGIN -> https://www.linkedin.com/feed/
[RATE_LIMIT] DAILY_CONNECTIONS: 5/30
```

## ⚠️ Ethical Usage

This framework enforces ethical automation:

- **Business hours only** - No 24/7 automation
- **Rate limits** - Respects LinkedIn's infrastructure
- **Manual intervention required** for security challenges
- **No brute-force** - Graceful abort on errors
- **Transparency** - All actions logged

## 🔧 Troubleshooting

### "Security challenge detected"
This is expected for first-time logins or suspicious activity. Run in non-headless mode (`BROWSER_HEADLESS=false`) and complete the challenge manually.

### "Failed to launch Chrome"
Ensure Chrome browser is installed and accessible in your PATH.

### "Login failed - still on login page"
Check your credentials in `.env`. Also ensure no security challenges are present.

## 📝 Development

### Code Quality Standards

- **Idiomatic Go** - Follow Go best practices
- **Clean Architecture** - Strict separation of concerns
- **No circular dependencies**
- **Explicit error handling**
- **Context-aware functions**
- **Max ~300 lines per file**
- **No magic numbers** (all in config)

### Testing

```bash
# Build to verify
go build ./...

# Run with verbose logging
BROWSER_HEADLESS=false go run cmd/bot/main.go
```

## 🤝 Contributing

This is a proof-of-concept system demonstrating:
- Advanced system design
- Anti-bot detection techniques
- Human behavior simulation
- Production Go patterns

## 📄 License

This project is for educational and professional skill demonstration purposes.

## ⚡ Performance

- **Initialization time**: ~3-5 seconds
- **Login time**: 5-10 seconds (with behavior simulation)
- **Memory usage**: ~150-200MB (Chrome + Go runtime)
- **Actions per hour**: 10-15 (with realistic delays)

## 🎓 Learning Resources

This project demonstrates:
- Chrome DevTools Protocol (CDP) integration
- Browser fingerprinting & evasion
- Béz ier curve mathematics
- Rate limiting algorithms
- Clean Architecture in Go
- Graceful shutdown patterns
- Structured logging

---

**Built with discipline, designed for scale.**
