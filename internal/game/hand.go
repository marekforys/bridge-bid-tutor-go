package game

import (
	"sort"
	"strings"
)

// Hand represents a player's hand of cards
type Hand struct {
	Cards []Card
}

// NewHand creates a new hand with the given cards
func NewHand(cards []Card) *Hand {
	h := &Hand{Cards: make([]Card, len(cards))}
	copy(h.Cards, cards)
	h.Sort()
	return h
}

// Sort sorts the cards in the hand by suit and rank
func (h *Hand) Sort() {
	sort.Slice(h.Cards, func(i, j int) bool {
		if h.Cards[i].Suit == h.Cards[j].Suit {
			return h.Cards[i].Rank > h.Cards[j].Rank // Higher ranks first
		}
		return h.Cards[i].Suit < h.Cards[j].Suit // Sort by suit
	})
}

// Evaluate evaluates the hand and returns the high card points and distribution
func (h *Hand) Evaluate() (hcp int, distribution map[Suit]int) {
	distribution = make(map[Suit]int)
	
	// Count high card points and distribution
	for _, card := range h.Cards {
		switch card.Rank {
		case Ace:
			hcp += 4
		case King:
			hcp += 3
		case Queen:
			hcp += 2
		case Jack:
			hcp += 1
		}
		distribution[card.Suit]++
	}
	
	// Add distribution points (optional, can be part of the bidding strategy)
	// for _, count := range distribution {
	//     if count == 2 {
	//         hcp += 1 // Doubleton
	//     } else if count == 1 {
	//         hcp += 2 // Singleton
	//     } else if count == 0 {
	//         hcp += 3 // Void
	//     }
	// }
	
	return hcp, distribution
}

// SuitCount returns the number of cards in the specified suit
func (h *Hand) SuitCount(s Suit) int {
	count := 0
	for _, card := range h.Cards {
		if card.Suit == s {
			count++
		}
	}
	return count
}

// HasStopper returns true if the hand has a stopper in the specified suit
// A stopper is typically considered as A, Kx, Qxx, or Jxxx
func (h *Hand) HasStopper(s Suit) bool {
	hasAce := false
	hasKing := false
	hasQueen := false
	count := 0

	for _, card := range h.Cards {
		if card.Suit == s {
			count++
			switch card.Rank {
			case Ace:
				hasAce = true
			case King:
				hasKing = true
			case Queen:
				hasQueen = true
			}
		}
	}

	return hasAce || (hasKing && count >= 2) || (hasQueen && count >= 3)
}

// GetSuit returns a string representation of cards in the specified suit
func (h *Hand) GetSuit(s Suit) string {
	var cards []string
	for _, card := range h.Cards {
		if card.Suit == s {
			cards = append(cards, card.RankString())
		}
	}
	return strings.Join(cards, " ")
}

// IsBalanced returns true if the hand is balanced.
// A balanced hand has no voids, no singletons, and at most one doubleton.
func (h *Hand) IsBalanced() bool {
	doubletonCount := 0
	for s := Clubs; s <= Spades; s++ {
		count := h.SuitCount(s)
		if count == 0 || count == 1 {
			return false // No voids or singletons
		}
		if count == 2 {
			doubletonCount++
		}
	}
	return doubletonCount <= 1
}
