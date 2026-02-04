package global

import (
	"os"
	"testing"
	"time"
)

// TestConstants tests global constants
func TestConstants(t *testing.T) {
	if WindowWidth != 725 {
		t.Errorf("WindowWidth = %d, want 725", WindowWidth)
	}
	if WindowHeight != 800 {
		t.Errorf("WindowHeight = %d, want 800", WindowHeight)
	}
}

// TestSaveFileType tests SaveFileType constants
func TestSaveFileType(t *testing.T) {
	if PC != 0 {
		t.Errorf("PC = %d, want 0", PC)
	}
	if PS != 1 {
		t.Errorf("PS = %d, want 1", PS)
	}
}

// TestPWD tests PWD variable initialization
func TestPWD(t *testing.T) {
	// PWD should be set during init
	if PWD == "" {
		t.Error("PWD is empty after initialization")
	}

	// PWD should be a valid directory
	info, err := os.Stat(PWD)
	if err != nil {
		t.Errorf("PWD is not a valid path: %v", err)
	}
	if !info.IsDir() {
		t.Error("PWD is not a directory")
	}
}

// TestNowToTicks tests the NowToTicks function
func TestNowToTicks(t *testing.T) {
	ticks1 := NowToTicks()
	time.Sleep(10 * time.Millisecond)
	ticks2 := NowToTicks()

	// Ticks should be positive
	if ticks1 == 0 {
		t.Error("NowToTicks() returned 0")
	}

	// Later time should have higher tick count
	if ticks2 <= ticks1 {
		t.Errorf("Ticks not increasing: ticks1=%d, ticks2=%d", ticks1, ticks2)
	}
}

// TestNowToTicksConsistency tests that NowToTicks produces reasonable values
func TestNowToTicksConsistency(t *testing.T) {
	before := time.Now()
	ticks := NowToTicks()
	after := time.Now()

	// Convert back to approximate time (rough check)
	// The formula is: uint64(float64(time.Now().UnixNano())*0.01) + uint64(60*60*24*365*1970*10000000)
	// This should produce a large number representing ticks since some epoch

	if ticks == 0 {
		t.Error("NowToTicks() returned 0")
	}

	// Just verify it's in a reasonable range (positive and very large)
	// The offset accounts for time since 1970 in 100-nanosecond intervals
	if ticks < uint64(60*60*24*365*1970*10000000) {
		t.Error("NowToTicks() returned unexpectedly small value")
	}

	_ = before
	_ = after
}

// TestSaveFileTypeValues tests that SaveFileType values are sequential
func TestSaveFileTypeValues(t *testing.T) {
	// Test that iota worked correctly
	types := []SaveFileType{PC, PS}
	for i, typ := range types {
		if int(typ) != i {
			t.Errorf("SaveFileType at index %d has value %d, want %d", i, typ, i)
		}
	}
}

// BenchmarkNowToTicks benchmarks the NowToTicks function
func BenchmarkNowToTicks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NowToTicks()
	}
}
