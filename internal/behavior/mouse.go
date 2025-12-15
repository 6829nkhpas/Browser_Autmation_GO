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

// MoveMouse moves the mouse cursor to target coordinates using Bézier curves for natural movement.
// This creates realistic mouse trajectories that mimic human behavior by:
// - Using cubic Bézier curves instead of straight lines
// - Adding micro-corrections (small random jitter)
// - Varying speed throughout the movement (acceleration/deceleration)
// - Calculating optimal number of steps based on distance
//
// Parameters:
//   - page: The Rod page instance to control
//   - targetX, targetY: Destination coordinates
//
// Returns error if mouse movement fails
func MoveMouse(page *rod.Page, targetX, targetY float64) error {
	// Get current mouse position as starting point
	// NOTE: GetMousePosition() is a placeholder. In a real scenario, you'd get the actual current mouse position.
	currentX := GetRandomFloat(100, 200) // Placeholder for actual current X
	currentY := GetRandomFloat(100, 200) // Placeholder for actual current Y

	// Calculate Euclidean distance to determine movement complexity
	// Longer distances need more steps for smooth, natural-looking movement
	distance := math.Sqrt(math.Pow(targetX-currentX, 2) + math.Pow(targetY-currentY, 2))

	// Determine number of movement steps based on distance
	// Formula: ~5 pixels per step provides good balance of smoothness and performance
	steps := int(distance / 5)
	if steps < 10 {
		steps = 10 // Minimum: Even short movements need smoothness
	}
	if steps > 100 {
		steps = 100 // Maximum: Caps computation for very long movements
	}

	// Generate Bézier curve control points for natural arc
	// Control points at 30% and 70% of the path create realistic curvature
	// This prevents straight-line "bot-like" movement
	// NOTE: generateControlPoint() is a placeholder. This logic is usually part of a bezier curve utility.
	cp1X := currentX + (targetX-currentX)*0.3 + GetRandomFloat(-50, 50) // Placeholder for actual control point generation
	cp1Y := currentY + (targetY-currentY)*0.3 + GetRandomFloat(-50, 50) // Placeholder for actual control point generation

	cp2X := currentX + (targetX-currentX)*0.7 + GetRandomFloat(-50, 50) // Placeholder for actual control point generation
	cp2Y := currentY + (targetY-currentY)*0.7 + GetRandomFloat(-50, 50) // Placeholder for actual control point generation

	// Traverse the Bézier curve in discrete steps
	for i := 0; i <= steps; i++ {
		// Calculate parametric value t (0.0 to 1.0) for curve position
		t := float64(i) / float64(steps)

		// Calculate point on cubic Bézier curve using formula:
		// B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
		// Where P₀=start, P₁=cp1, P₂=cp2, P₃=end
		// NOTE: cubicBezier() is a placeholder. This logic is usually part of a bezier curve utility.
		mt := 1 - t
		mt2 := mt * mt
		mt3 := mt2 * mt
		t2 := t * t
		t3 := t2 * t

		x := mt3*currentX + 3*mt2*t*cp1X + 3*mt*t2*cp2X + t3*targetX
		y := mt3*currentY + 3*mt2*t*cp1Y + 3*mt*t2*cp2Y + t3*targetY

		// Add micro-corrections: small random offsets simulate human hand tremor
		// Range: ±1 pixel creates realistic jitter without visible shakiness
		// NOTE: This uses math/rand.Float64(), which needs math/rand import.
		x += (math.Rand().Float64() - 0.5) * 2
		y += (math.Rand().Float64() - 0.5) * 2

		// Execute mouse movement to calculated position
		// Third parameter (1) indicates number of steps for Rod's internal interpolation
		if err := page.Mouse.Move(x, y, 1); err != nil {
			return err
		}

		// Variable delay between movements creates acceleration/deceleration
		// Natural mouse movement is faster in the middle, slower at start/end
		// NOTE: calculateDelay() is a placeholder. This logic is usually part of a speed calculation utility.
		delay := calculateMouseSpeed(t) // Reusing existing calculateMouseSpeed for delay
		time.Sleep(time.Duration(delay) * time.Millisecond)
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
