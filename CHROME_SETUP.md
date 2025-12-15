# LinkedIn Automation Bot - Chrome Browser Configuration

## Current Situation

The bot needs a Chrome/Chromium browser to run. Rod (the browser automation library) will automatically download Chromium if one isn't found on your system.

## Options

### Option 1: Let Rod Download Chromium (Current)
**Pros**: Automatic, no manual installation needed  
**Cons**: Downloads ~150MB on first run, uses Rod's managed Chromium

The bot is configured to do this already.

### Option 2: Use Your Installed Chrome (If Available)
**Pros**: Uses your existing browser, familiar environment  
**Cons**: Requires Chrome to be installed in standard locations

**To use this**: Install Chrome via one of these methods:

```bash
# Method 1: Using snap (recommended for Ubuntu)
sudo snap install chromium

# Method 2: Download from Google and install .deb
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome-stable_current_amd64.deb
sudo apt-get install -f

# Method 3: Via apt repository
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
sudo apt update
sudo apt install google-chrome-stable
```

### Option 3: Specify Custom Chrome Path

If you have Chrome installed in a custom location, edit `.env`:

```bash
CHROME_BINARY_PATH=/path/to/your/chrome
```

## Current  Configuration

The bot's `chrome.go` already tries to find system Chrome in this order:
1. `google-chrome`
2. `google-chrome-stable`
3. `chromium`
4. `chromium-browser`
5. Falls back to Rod's auto-download

**Recommendation**: Let Rod download Chromium (easiest) OR install Chrome via snap if you prefer using system Chrome.
