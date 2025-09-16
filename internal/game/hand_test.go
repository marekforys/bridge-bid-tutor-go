package game

import (
	"reflect"
	"testing"
)

func TestHand_Evaluate(t *testing.T) {
	tests := []struct {
		name           string
		hand           Hand
		wantHCP        int
		wantDistribution map[Suit]int
	}{
		{
			name: "Basic hand evaluation",
			hand: Hand{
				Cards: []Card{
					{Suit: Spades, Rank: Ace},
					{Suit: Spades, Rank: King},
					{Suit: Hearts, Rank: Queen},
					{Suit: Diamonds, Rank: Jack},
					{Suit: Clubs, Rank: Two},
				},
			},
			wantHCP:        10, // 4 (A) + 3 (K) + 2 (Q) + 1 (J)
			wantDistribution: map[Suit]int{Spades: 2, Hearts: 1, Diamonds: 1, Clubs: 1},
		},
		{
			name: "Hand with no points",
			hand: Hand{
				Cards: []Card{
					{Suit: Spades, Rank: Two},
					{Suit: Spades, Rank: Three},
					{Suit: Hearts, Rank: Four},
					{Suit: Diamonds, Rank: Five},
					{Suit: Clubs, Rank: Six},
				},
			},
			wantHCP:        0,
			wantDistribution: map[Suit]int{Spades: 2, Hearts: 1, Diamonds: 1, Clubs: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHCP, gotDistribution := tt.hand.Evaluate()
			if gotHCP != tt.wantHCP {
				t.Errorf("Hand.Evaluate() gotHCP = %v, want %v", gotHCP, tt.wantHCP)
			}
			if !reflect.DeepEqual(gotDistribution, tt.wantDistribution) {
				t.Errorf("Hand.Evaluate() gotDistribution = %v, want %v", gotDistribution, tt.wantDistribution)
			}
		})
	}
}
