package game

import (
	"testing"
)

func TestRomanKeyCardBlackwood(t *testing.T) {
	tests := []struct {
		name           string
		openerHand     []Card
		auction       []Bid
		expectedBid    Bid
		description    string
	}{
		{
			name: "Respond with 0 key cards without Queen",
			openerHand: []Card{
				// No Aces or King of trumps
				// Add some non-key cards to make a valid hand
				{Suit: Spades, Rank: Two},
				{Suit: Hearts, Rank: Three},
				{Suit: Diamonds, Rank: Four},
				{Suit: Clubs, Rank: Five},
				{Suit: Spades, Rank: Six},
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: North}, // Contract in spades
				{Level: 4, Strain: NoTrump, Position: South}, // 4NT Blackwood
			},
			expectedBid: NewBid(5, Clubs), // 0 key cards without Queen
			description: "Should bid 5♣ showing 0 key cards without Queen",
		},
		{
			name: "Respond with 1 key card with Queen",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: Queen}, // Queen of trumps
				// No other key cards
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: North},
				{Level: 4, Strain: NoTrump, Position: South},
			},
			expectedBid: NewBid(5, Spades), // 1 key card with Queen
			description: "Should bid 5♠ showing 1 key card with Queen",
		},
		{
			name: "Respond with 2 key cards without Queen",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: King},  // Key card (King of trumps)
				// No Queen of trumps
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: North},
				{Level: 4, Strain: NoTrump, Position: South},
			},
			expectedBid: NewBid(5, Hearts), // 2 key cards without Queen
			description: "Should bid 5♥ showing 2 key cards without Queen",
		},
		{
			name: "Respond with 3 key cards with Queen",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: King},  // Key card (King of trumps)
				{Suit: Hearts, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: Queen}, // Queen of trumps
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: North},
				{Level: 4, Strain: NoTrump, Position: South},
			},
			expectedBid: NewBid(5, Spades), // 3 key cards with Queen
			description: "Should bid 5♠ showing 3 key cards with Queen",
		},
		{
			name: "Respond with 4 key cards with Queen",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: King},  // Key card (King of trumps)
				{Suit: Hearts, Rank: Ace},   // Key card (Ace)
				{Suit: Diamonds, Rank: Ace}, // Key card (Ace)
				{Suit: Spades, Rank: Queen}, // Queen of trumps
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: North},
				{Level: 4, Strain: NoTrump, Position: South},
			},
			expectedBid: NewBid(5, Diamonds), // 4 key cards with Queen
			description: "Should bid 5♦ showing 4 key cards with Queen",
		},
		{
			name: "Standard Blackwood in NoTrump",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Hearts, Rank: Ace},   // Key card (Ace)
				// No King of trumps in NoTrump
			},
			auction: []Bid{
				{Level: 3, Strain: NoTrump, Position: North},
				{Level: 4, Strain: NoTrump, Position: South}, // 4NT Blackwood
			},
			expectedBid: NewBid(5, Hearts), // 2 Aces in standard Blackwood
			description: "In NoTrump, should count only Aces and bid 5♥ showing 2 Aces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the player with the test hand
			player := NewPlayer(North)
			player.Hand = NewHand(tt.openerHand)

			// Set up the auction
			auction := NewAuction()
			for _, bid := range tt.auction {
				auction.AddBid(bid)
			}

			// Get the response to Blackwood
			bid := player.MakeBid(auction)

			// Verify the response
			if bid.Level != tt.expectedBid.Level || bid.Strain != tt.expectedBid.Strain {
				t.Errorf("%s\nExpected: %s\nGot: %s", 
					tt.description, 
					tt.expectedBid, 
					bid)
			}
		})
	}
}

func TestBlackwoodInitiation(t *testing.T) {
	tests := []struct {
		name           string
		openerHand     []Card
		auction       []Bid
		expectedBid    Bid
		description    string
	}{
		{
			name: "Respond to 4NT with all key cards",
			openerHand: []Card{
				{Suit: Spades, Rank: Ace},   // Key card (Ace)
				{Suit: Spades, Rank: King},  // Key card (King of trumps)
				{Suit: Hearts, Rank: Ace},   // Key card (Ace)
				{Suit: Diamonds, Rank: Ace}, // Key card (Ace)
				{Suit: Spades, Rank: Queen}, // Queen of trumps
				// Add more cards to make a complete hand
				{Suit: Hearts, Rank: Two},
				{Suit: Diamonds, Rank: Three},
			},
			auction: []Bid{
				{Level: 4, Strain: Spades, Position: South},
				{Level: 4, Strain: NoTrump, Position: North}, // Partner bids 4NT
			},
			expectedBid: NewBid(5, Diamonds), // 4 key cards with Queen
			description: "Should respond 5♦ showing 4 key cards with Queen",
		},
		{
			name: "Respond to 4NT with no key cards",
			openerHand: []Card{
				// No key cards
				{Suit: Spades, Rank: Two},
				{Suit: Hearts, Rank: Three},
				{Suit: Diamonds, Rank: Four},
				{Suit: Clubs, Rank: Five},
				// Add more cards to make a complete hand
				{Suit: Spades, Rank: Six},
				{Suit: Hearts, Rank: Seven},
			},
			auction: []Bid{
				{Level: 4, Strain: Hearts, Position: South},
				{Level: 4, Strain: NoTrump, Position: North}, // Partner bids 4NT
			},
			expectedBid: NewBid(5, Clubs), // 0 key cards without Queen
			description: "Should respond 5♣ showing 0 key cards without Queen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the player with the test hand
			player := NewPlayer(South)
			player.Hand = NewHand(tt.openerHand)

			// Set up the auction
			auction := NewAuction()
			for _, bid := range tt.auction {
				auction.AddBid(bid)
			}

			// Get the player's bid
			bid := player.MakeBid(auction)

			// Verify the bid matches the expected bid
			if bid.Level != tt.expectedBid.Level || bid.Strain != tt.expectedBid.Strain {
				t.Errorf("%s\nExpected: %s\nGot: %s", 
					tt.description, 
					tt.expectedBid, 
					bid)
			}
		})
	}
}
