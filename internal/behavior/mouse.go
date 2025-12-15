package behavior

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Point represents a 2D point
type Point struct {
	X float64
	Y float64
}

// MoveMouse moves the mouse to target coordinates using Bézier curve
// This creates natural, curved mouse movement like a real human
func MoveMouse(page *rod.Page, targetX, targetY float64) error {
	// Get current mouse position (start from a random position if unknown)
	startX := GetRandomFloat(100, 200)
	startY := GetRandomFloat(100, 200)

	// Generate control points for Bézier curve
	controlPoints := generateBezierControlPoints(
		Point{X: startX, Y: startY},
		Point{X: targetX, Y: targetY},
	)

	// Calculate number of steps based on distance
	distance := math.Sqrt(math.Pow(targetX-startX, 2) + math.Pow(targetY-startY, 2))
	steps := int(distance / 10) // ~10 pixels per step
	if steps < 10 {
		steps = 10
	}
	if steps > 100 {
		steps = 100
	}

	// Move mouse along the curve
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		// Calculate point on Bézier curve
		point := calculateBezierPoint(t, controlPoints)

		// Add micro-corrections (small random deviations)
		point.X += GetRandomFloat(-2, 2)
		point.Y += GetRandomFloat(-2, 2)

		// Move mouse to this point
		err := page.Mouse.Move(point.X, point.Y, 1)
		if err != nil {
			return fmt.Errorf("mouse move failed: %w", err)
		}

		// Variable speed - faster in middle, slower at start/end
		delay := calculateMouseSpeed(t)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Final position adjustment
	err := page.Mouse.Move(targetX, targetY, 1)
	if err != nil {
		return fmt.Errorf("final mouse move failed: %w", err)
	}

	return nil
}

// generateBezierControlPoints creates control points for a cubic Bézier curve
func generateBezierControlPoints(start, end Point) []Point {
	// Create two control points for smooth curve
	cp1X := start.X + (end.X-start.X)/3 + GetRandomFloat(-50, 50)
	cp1Y := start.Y + (end.Y-start.Y)/3 + GetRandomFloat(-50, 50)

	cp2X := start.X + 2*(end.X-start.X)/3 + GetRandomFloat(-50, 50)
	cp2Y := start.Y + 2*(end.Y-start.Y)/3 + GetRandomFloat(-50, 50)

	return []Point{
		start,
		{X: cp1X, Y: cp1Y},
		{X: cp2X, Y: cp2Y},
		end,
	}
}

// calculateBezierPoint calculates a point on a cubic Bézier curve
// t is between 0 and 1
func calculateBezierPoint(t float64, points []Point) Point {
	// Cubic Bézier formula: B(t) = (1-t)³P0 + 3(1-t)²tP1 + 3(1-t)t²P2 + t³P3
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t

	x := mt3*points[0].X +
		3*mt2*t*points[1].X +
		3*mt*t2*points[2].X +
		t3*points[3].X

	y := mt3*points[0].Y +
		3*mt2*t*points[1].Y +
		3*mt*t2*points[2].Y +
		t3*points[3].Y

	return Point{X: x, Y: y}
}

// calculateMouseSpeed returns delay in milliseconds based on position in movement
// Movements are slower at start and end, faster in the middle
func calculateMouseSpeed(t float64) int {
	// Use sine wave to create acceleration/deceleration
	// Slower at t=0 and t=1, faster at t=0.5
	speed := math.Sin(t * math.Pi)

	// Map to delay: higher speed = lower delay
	minDelay := 5
	maxDelay := 15
	delay := maxDelay - int(speed*float64(maxDelay-minDelay))

	return delay
}

// ClickAtPosition clicks at specific coordinates
func ClickAtPosition(page *rod.Page, x, y float64) error {
	// Move to position first
	if err := MoveMouse(page, x, y); err != nil {
		return err
	}

	WaitHuman(100, 300)

	// Perform click
	if err := page.Mouse.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("click failed: %w", err)
	}

	return nil
}

// GetRandomFloat returns a random float between min and max
func GetRandomFloat(min, max float64) float64 {
	if min >= max {
		return min
	}

	rang := max - min
	n, err := rand.Int(rand.Reader, big.NewInt(int64(rang*1000)))
	if err != nil {
		return min
	}

	return min + float64(n.Int64())/1000.0
}
