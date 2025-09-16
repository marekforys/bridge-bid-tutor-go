package game

import (
	"testing"
)

func TestAI_PolishClub(t *testing.T) {
	// Test case: Opener has 18+ HCP, should open 1C.
	t.Run("Opener opens 1C with strong hand", func(t *testing.T) {
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen},
			{Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
			{Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Queen},
			{Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack},
		})

		auction := NewAuction()
		bid := opener.MakeBid(auction)

		if bid.Level != 1 || bid.Strain != Clubs {
			t.Errorf("Expected 1C opening with strong hand, but got %s", bid)
		}
	})

	// Test case: Responder has 0-6 HCP, should respond 1D (negative).
	t.Run("Responder gives 1D negative response", func(t *testing.T) {
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Two}, {Suit: Spades, Rank: Three}, {Suit: Spades, Rank: Four},
			{Suit: Hearts, Rank: Two}, {Suit: Hearts, Rank: Three}, {Suit: Hearts, Rank: Four},
			{Suit: Diamonds, Rank: Two}, {Suit: Diamonds, Rank: Three}, {Suit: Diamonds, Rank: Four},
			{Suit: Clubs, Rank: Two}, {Suit: Clubs, Rank: Three}, {Suit: Clubs, Rank: Four}, {Suit: Clubs, Rank: Five},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // Partner opens 1C

		bid := responder.MakeBid(auction)

		if bid.Level != 1 || bid.Strain != Diamonds {
			t.Errorf("Expected 1D negative response, but got %s", bid)
		}
	})
}

func TestAI_Stayman(t *testing.T) {
	// Test case: Responder has 8+ HCP and a 4-card major, should bid 2C (Stayman).
	t.Run("Responder bids Stayman", func(t *testing.T) {
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			{Suit: Hearts, Rank: Two}, {Suit: Hearts, Rank: Three}, {Suit: Hearts, Rank: Four},
			{Suit: Diamonds, Rank: Two}, {Suit: Diamonds, Rank: Three}, {Suit: Diamonds, Rank: Four},
			{Suit: Clubs, Rank: Two}, {Suit: Clubs, Rank: Three}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // Partner opens 1NT

		bid := responder.MakeBid(auction)

		if bid.Level != 2 || bid.Strain != Clubs {
			t.Errorf("Expected Stayman 2C bid, but got %s", bid)
		}
	})

	// Test case: Opener has a 4-card major and must respond to Stayman.
	t.Run("Opener responds to Stayman", func(t *testing.T) {
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			{Suit: Hearts, Rank: Two}, {Suit: Hearts, Rank: Three}, {Suit: Hearts, Rank: Four},
			{Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Queen},
			{Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North})      // We open 1NT
		auction.AddBid(Bid{Level: 2, Strain: Clubs, Position: South}) // Partner bids 2C (Stayman)

		bid := opener.MakeBid(auction)

		if bid.Level != 2 || bid.Strain != Spades {
			t.Errorf("Expected 2S response to Stayman, but got %s", bid)
		}
	})
}
