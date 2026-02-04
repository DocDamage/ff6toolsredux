package main

import (
	"ffvi_editor/io"
	"fmt"
)

func main() {
	fmt.Println("Testing ROM loading with correct ROM...")

	// Test auto-detection
	extractor, err := io.NewROMSpriteExtractor("")
	if err != nil {
		fmt.Printf("❌ Failed to load ROM: %v\n", err)
		return
	}

	fmt.Println("✅ ROM loaded successfully!")

	// Test extracting a few characters
	characters := []string{"Terra", "Locke", "Cyan", "Shadow", "Edgar"}
	for i, name := range characters {
		sprite, err := extractor.ExtractCharacterSprite(i)
		if err == nil && sprite != nil && len(sprite.Data) > 0 {
			fmt.Printf("✅ %s: %d bytes\n", name, len(sprite.Data))
		} else {
			fmt.Printf("❌ %s: Failed (%v)\n", name, err)
		}
	}

	fmt.Println("\n🎉 ROM loading test complete!")
}
