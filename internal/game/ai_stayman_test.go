package game

import (
	"testing"
)

// TestStaymanConvention tests the Stayman convention after a 1NT opening
func TestStaymanConvention(t *testing.T) {
	// Test case: Responder has 8+ HCP and a 4-card major, should bid 2C (Stayman).
	t.Run("Responder with 4 hearts and 8 HCP uses Stayman", func(t *testing.T) {
		// Set up responder's hand: 8 HCP, 4 hearts
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			// 4 spades (A K Q J) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			// 3 hearts (K Q J) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Jack},
			// 3 diamonds (5 4 3) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 3 clubs (5 4 3) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four}, {Suit: Clubs, Rank: Three},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // Partner opens 1NT

		bid := responder.MakeBid(auction)

		// Should bid 2♣ (Stayman)
		if bid.Level != 2 || bid.Strain != Clubs {
			t.Errorf("Expected 2C (Stayman), got %s", bid)
		}
	})

	// Test case: Responder with 5-4 in majors should use transfer, not Stayman
	t.Run("Responder with 5-4 majors uses transfer", func(t *testing.T) {
		// Set up responder's hand: 5-4 in majors
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			// 5 spades (A K Q J 10) - 11 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack}, {Suit: Spades, Rank: Ten},
			// 4 hearts (K Q J 10) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, 
			{Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 2 diamonds (5 4) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // Partner opens 1NT

		bid := responder.MakeBid(auction)

		// Should bid 2H (transfer to spades) with 5-4 shape
		if bid.Level != 2 || bid.Strain != Hearts {
			t.Errorf("Expected 2H (transfer to spades) with 5-4 majors, got %s", bid)
		}
	})

	// Test case: Opener has both majors, shows hearts first (lower ranking suit)
	t.Run("Opener has both majors, shows hearts first", func(t *testing.T) {
		// Set up opener's hand with both majors (4-4)
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			// 4 spades (A K Q J) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			// 4 hearts (K Q J 10) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, 
			{Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 3 diamonds (5 4 3) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // 1NT opening
		auction.AddBid(Bid{Level: 2, Strain: 0, Position: South}) // 2♣ Stayman

		bid := opener.MakeBid(auction)

		// Should show hearts first (lower ranking suit) when holding both majors
		if bid.Level != 2 || bid.Strain != Hearts {
			t.Errorf("Expected 2H (showing hearts first with 4-4), got %s", bid)
		}
	})

	// Test case: Responder with invitational strength and 5-card major
	t.Run("Responder with invitational strength and 5-card major", func(t *testing.T) {
		// Set up responder's hand with invitational strength (10-11 HCP) and 5 hearts
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			// 5 hearts (A K Q J 10) - 11 HCP
			{Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: King}, 
			{Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 3 spades (K Q J) - 9 HCP
			{Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			// 3 diamonds (5 4 3) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // 1NT opening

		bid := responder.MakeBid(auction)

		// With 5-card major, should transfer instead of using Stayman
		if bid.Level != 2 || bid.Strain != Diamonds {
			t.Errorf("Expected 2D (transfer to hearts), got %s", bid)
		}
	})

	// Test case: Opener with 5-card spades and 4-card hearts shows hearts first
	t.Run("Opener with 5-card spades and 4-card hearts shows hearts first", func(t *testing.T) {
		// Set up opener's hand with 5 spades and 4 hearts
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			// 5 spades (A K Q J 10) - 11 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack}, {Suit: Spades, Rank: Ten},
			// 4 hearts (K Q J 10) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, 
			{Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 2 diamonds (5 4) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // 1NT opening
		auction.AddBid(Bid{Level: 2, Strain: 0, Position: South}) // 2♣ Stayman

		bid := opener.MakeBid(auction)

		// With 4+ hearts, should show hearts first regardless of longer spades
		if bid.Level != 2 || bid.Strain != Hearts {
			t.Errorf("Expected 2H (showing hearts first with 4+ hearts), got %s", bid)
		}
	})

	// Test case: Responder with 4-4 minors should not use Stayman
	t.Run("Responder with 4-4 minors should not use Stayman", func(t *testing.T) {
		// Set up responder's hand with 4-4 in minors, no major fit
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			// 3 spades (A K Q) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen},
			// 2 hearts (K Q) - 6 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
			// 4 diamonds (A 5 4 3) - 4 HCP
			{Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Five}, 
			{Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 4 clubs (A 5 4 3) - 4 HCP
			{Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Five}, 
			{Suit: Clubs, Rank: Four}, {Suit: Clubs, Rank: Three},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // 1NT opening

		bid := responder.MakeBid(auction)

		// Should not use Stayman without a 4-card major
		if bid.Strain == Clubs && bid.Level == 2 {
			t.Error("Should not use Stayman without a 4-card major")
		}
	})

	// Test case: Opener with 4 hearts and 5+ spades shows hearts first
	t.Run("Opener with 4 hearts and 5+ spades shows hearts first", func(t *testing.T) {
		// Set up opener's hand with 4 hearts and 5+ spades
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			// 5 spades (A K Q J 10) - 11 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack}, {Suit: Spades, Rank: Ten},
			// 4 hearts (K Q J 10) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, 
			{Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 2 diamonds (5 4) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // 1NT opening
		auction.AddBid(Bid{Level: 2, Strain: 0, Position: South}) // 2♣ Stayman

		bid := opener.MakeBid(auction)

		// With 4-4 in majors, should show hearts first (lower ranking suit)
		if bid.Level != 2 || bid.Strain != Hearts {
			t.Errorf("Expected 2H (showing hearts first with 4-4), got %s", bid)
		}
	})

	// Test case: Opener has a 4-card major and responds to Stayman.
	t.Run("Opener shows 4-card major", func(t *testing.T) {
		// Set up opener's hand: 4 spades
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			// 4 spades (A K Q J) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			// 3 hearts (K Q J) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Jack},
			// 3 diamonds (5 4 3) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 3 clubs (5 4 3) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four}, {Suit: Clubs, Rank: Three},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // We open 1NT
		auction.AddBid(Bid{Level: 2, Strain: 0, Position: South}) // Partner bids 2C (Stayman)

	bid := opener.MakeBid(auction)

	// Should show 4+ spades
	if bid.Level != 2 || bid.Strain != Spades {
		t.Errorf("Expected 2S (showing spades), got %s", bid)
	}
	})

	// Test case: Opener has no 4-card major
	t.Run("Opener denies 4-card major", func(t *testing.T) {
		// Set up opener's hand: no 4-card major
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			// 3 spades (A K Q) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen},
			// 3 hearts (K Q J) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Jack},
			// 4 diamonds (A 5 4 3) - 4 HCP
			{Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Five}, 
			{Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 3 clubs (5 4 3) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four}, {Suit: Clubs, Rank: Three},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // We open 1NT
		auction.AddBid(Bid{Level: 2, Strain: 0, Position: South}) // Partner bids 2C (Stayman)

	bid := opener.MakeBid(auction)

	// Should bid 2♦ (no 4-card major)
	if bid.Level != 2 || bid.Strain != Diamonds {
		t.Errorf("Expected 2D (no major), got %s", bid)
	}
	})

	// Test case: Responder has both majors
	t.Run("Responder with both majors uses Stayman", func(t *testing.T) {
		// Set up responder's hand: 4 spades and 4 hearts
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			// 4 spades (A K Q J) - 10 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, 
			{Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			// 4 hearts (K Q J 10) - 9 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, 
			{Suit: Hearts, Rank: Jack}, {Suit: Hearts, Rank: Ten},
			// 3 diamonds (5 4 3) - 0 HCP
			{Suit: Diamonds, Rank: Five}, {Suit: Diamonds, Rank: Four}, {Suit: Diamonds, Rank: Three},
			// 2 clubs (5 4) - 0 HCP
			{Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: 4, Position: North}) // Partner opens 1NT

	bid := responder.MakeBid(auction)

	// Should bid 2♣ (Stayman) to find the major fit
	if bid.Level != 2 || bid.Strain != Clubs {
		t.Errorf("Expected 2C (Stayman), got %s", bid)
	}
	})
}
