package behavior

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"unicode"

	"github.com/go-rod/rod"
)

// TypeHumanLike types text with human-like behavior:
// - Variable typing speed
// - Occasional typos with corrections
// - Thinking pauses at punctuation
// - Realistic WPM (40-80)
func TypeHumanLike(elem *rod.Element, text string) error {
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// Occasional typo (5% chance)
		if shouldMakeTypo() && i < len(runes)-1 {
			// Type wrong character
			wrongChar := getRandomTypo(char)
			if err := typeChar(elem, wrongChar); err != nil {
				return err
			}

			// Short pause to "notice" the typo
			WaitHuman(100, 300)

			// Backspace to delete typo
			if err := elem.Input(string(rune(8))); err != nil {
				return fmt.Errorf("backspace failed: %w", err)
			}

			WaitHuman(50, 150)
		}

		// Type the correct character
		if err := typeChar(elem, char); err != nil {
			return err
		}

		// Variable delay between characters
		delay := getTypingDelay(char, i, len(runes))
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	return nil
}

// typeChar types a single character
func typeChar(elem *rod.Element, char rune) error {
	// Use Input for single character
	err := elem.Input(string(char))
	if err != nil {
		return fmt.Errorf("failed to type char '%c': %w", char, err)
	}
	return nil
}

// getTypingDelay returns delay in milliseconds for typing a character
// Realistic typing speed: 40-80 WPM = 120-200ms per character
func getTypingDelay(char rune, position, total int) int {
	baseDelay := GetRandomInRange(80, 200)

	// Longer pause after punctuation (thinking time)
	if isPunctuation(char) {
		baseDelay += GetRandomInRange(100, 400)
	}

	// Longer pause after spaces (between words)
	if char == ' ' {
		baseDelay += GetRandomInRange(20, 80)
	}

	// Slightly faster in the middle of text (flow state)
	if position > total/4 && position < 3*total/4 {
		baseDelay = int(float64(baseDelay) * 0.9)
	}

	// Occasional longer pause (thinking/looking at screen)
	if shouldPause() {
		baseDelay += GetRandomInRange(200, 600)
	}

	return baseDelay
}

// isPunctuation checks if character is punctuation
func isPunctuation(char rune) bool {
	return unicode.IsPunct(char) || char == '.' || char == ',' || char == '!' || char == '?'
}

// shouldMakeTypo returns true 5% of the time
func shouldMakeTypo() bool {
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return false
	}
	return n.Int64() < 5
}

// shouldPause returns true 10% of the time
func shouldPause() bool {
	n, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return false
	}
	return n.Int64() < 10
}

// getRandomTypo returns a typo for the given character
// Simulates keyboard proximity errors
func getRandomTypo(char rune) rune {
	// Keyboard layout proximity map (QWERTY)
	proximityMap := map[rune]string{
		'a': "sqwz",
		'b': "vghn",
		'c': "xdfv",
		'd': "sfcre",
		'e': "rwd",
		'f': "gdvcrt",
		'g': "fhbvty",
		'h': "gjnbyu",
		'i': "uok",
		'j': "hknum",
		'k': "jlmio",
		'l': "kop",
		'm': "njk",
		'n': "bhjm",
		'o': "ipkl",
		'p': "ol",
		'q': "wa",
		'r': "etfd",
		's': "awedxz",
		't': "ryfg",
		'u': "yhji",
		'v': "cfgb",
		'w': "qesa",
		'x': "zsdc",
		'y': "tugh",
		'z': "asx",
	}

	lowerChar := unicode.ToLower(char)
	nearby, exists := proximityMap[lowerChar]

	if !exists || len(nearby) == 0 {
		// If no proximity mapping, just return adjacent character
		return char + 1
	}

	// Pick random nearby character
	idx := GetRandomInRange(0, len(nearby)-1)
	typo := rune(nearby[idx])

	// Preserve case
	if unicode.IsUpper(char) {
		typo = unicode.ToUpper(typo)
	}

	return typo
}

// TypePassword types a password with consistent timing (no typos for passwords)
func TypePassword(elem *rod.Element, password string) error {
	for _, char := range password {
		if err := typeChar(elem, char); err != nil {
			return err
		}
		// Consistent timing for passwords
		WaitHuman(80, 150)
	}
	return nil
}

// ClearInput clears an input field naturally
func ClearInput(elem *rod.Element) error {
	// Select all
	if err := elem.SelectAllText(); err != nil {
		return fmt.Errorf("select all failed: %w", err)
	}

	WaitHuman(100, 200)

	// Delete
	if err := elem.Input(""); err != nil {
		return fmt.Errorf("clear failed: %w", err)
	}

	return nil
}
