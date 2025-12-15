#!/bin/bash

# LinkedIn Automation Bot - Setup Script
# This script installs required system dependencies for Chrome/Chromium

echo "=== LinkedIn Automation Bot - Dependency Installer ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "This script needs sudo privileges to install system packages."
    echo "You may be prompted for your password."
    echo ""
fi

echo "Installing Chrome/Chromium dependencies..."
echo ""

sudo apt-get update

sudo apt-get install -y \
    libnss3 \
    libnspr4 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxfixes3 \
    libxrandr2 \
    libgbm1 \
    libasound2

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Dependencies installed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Make sure .env file is configured with your LinkedIn credentials"
    echo "2. Run the bot: ./bin/linkedin-bot"
    echo ""
else
    echo ""
    echo "❌ Failed to install some dependencies."
    echo "Please check the error messages above."
    exit 1
fi
