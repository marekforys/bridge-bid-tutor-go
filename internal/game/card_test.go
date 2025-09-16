package game

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	if len(deck) != 52 {
		t.Errorf("Expected deck length of 52, but got %d", len(deck))
	}

	// Check for uniqueness
	cardMap := make(map[Card]bool)
	for _, card := range deck {
		if cardMap[card] {
			t.Errorf("Duplicate card found: %s", card)
		}
		cardMap[card] = true
	}
}

func TestDeck_Shuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	deck2.Shuffle()

	if len(deck1) != len(deck2) {
		t.Fatalf("Shuffled deck has different length")
	}

	// It's statistically improbable for a shuffled deck to be the same as a new one.
	// We'll just check if they are not identical.
	same := true
	for i := range deck1 {
		if deck1[i] != deck2[i] {
			same = false
			break
		}
	}

	if same {
		t.Error("Shuffled deck is identical to a new deck")
	}
}
