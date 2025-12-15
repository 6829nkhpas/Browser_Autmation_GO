package behavior

import (
	"crypto/rand"
	"math/big"
	"time"
)

// WaitHuman adds a random delay to simulate human thinking/reaction time
func WaitHuman(minMs, maxMs int) {
	if minMs >= maxMs {
		time.Sleep(time.Duration(minMs) * time.Millisecond)
		return
	}

	delay := GetRandomInRange(minMs, maxMs)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// WaitReading simulates time spent reading content
// Estimate based on word count (average reading speed: 200-250 WPM)
func WaitReading(wordCount int) {
	if wordCount <= 0 {
		WaitHuman(500, 1000)
		return
	}

	// Assume 200 WPM = ~300ms per word
	baseTime := wordCount * 300

	// Add variation (+/- 30%)
	variation := GetRandomInRange(-baseTime/3, baseTime/3)
	totalTime := baseTime + variation

	// Cap at reasonable limits
	if totalTime < 500 {
		totalTime = 500
	}
	if totalTime > 10000 {
		totalTime = 10000
	}

	time.Sleep(time.Duration(totalTime) * time.Millisecond)
}

// WaitThinking simulates decision-making pause
func WaitThinking() {
	WaitHuman(1000, 3000)
}

// WaitAfterClick simulates time after clicking (waiting for response/load)
func WaitAfterClick() {
	WaitHuman(800, 1500)
}

// WaitBeforeAction simulates preparation time before action
func WaitBeforeAction() {
	WaitHuman(500, 1200)
}

// WaitLong simulates extended pause (distraction, multitasking)
func WaitLong() {
	WaitHuman(3000, 8000)
}

// GetRandomInRange returns a random integer between min and max (inclusive)
func GetRandomInRange(min, max int) int {
	if min >= max {
		return min
	}

	rang := int64(max - min + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(rang))
	if err != nil {
		return min
	}

	return min + int(n.Int64())
}

// GetRandomDuration returns a random duration between min and max
func GetRandomDuration(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}

	minMs := min.Milliseconds()
	maxMs := max.Milliseconds()

	randomMs := GetRandomInRange(int(minMs), int(maxMs))
	return time.Duration(randomMs) * time.Millisecond
}

// ShouldTakeBreak returns true with given probability (0-100)
func ShouldTakeBreak(probabilityPercent int) bool {
	if probabilityPercent <= 0 {
		return false
	}
	if probabilityPercent >= 100 {
		return true
	}

	return GetRandomInRange(0, 100) < probabilityPercent
}
