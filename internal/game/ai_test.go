package game

import (
	"testing"
)

func TestAI_PolishClub(t *testing.T) {
	// Test case: Opener has a weak, balanced hand and rebids 1NT.
	t.Run("Opener rebids 1NT after 1D response", func(t *testing.T) {
		opener := NewPlayer(North)
		// 12 HCP, balanced (4-3-3-3), no 5-card major
		opener.Hand = NewHand([]Card{
			// Spades (4): A Q 7 6 -> 6 HCP
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Seven}, {Suit: Spades, Rank: Six},
			// Hearts (3): K 8 5 -> 3 HCP
			{Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Eight}, {Suit: Hearts, Rank: Five},
			// Diamonds (3): K 9 4 -> 3 HCP
			{Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
			// Clubs (3): 3 2 J -> 1 HCP (J) making total 13 HCP; adjust to 12 by using 10 instead of J
			{Suit: Clubs, Rank: Ten}, {Suit: Clubs, Rank: Three}, {Suit: Clubs, Rank: Two},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})    // We open 1C
		auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South}) // Partner responds 1D

		bid := opener.MakeBid(auction)

		if bid.Level != 1 || bid.Strain != 4 { // 1NT
			t.Errorf("Expected 1NT rebid, but got %s", bid)
		}
	})
}

func TestAI_StrongClubContinuations_Duplicate(t *testing.T) {
    t.Skip("superseded by ai_strong_club_test.go")
    // 1C (strong) - 1D; balanced 18-19 -> 2NT
    t.Run("Strong balanced 18-19 rebids 2NT", func(t *testing.T) {
        opener := NewPlayer(North)
        // 18-19 HCP balanced: e.g., A K Q in two suits, KQ elsewhere
        opener.Hand = NewHand([]Card{
            // Spades: A Q 7
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Seven},
            // Hearts: K Q 6
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Six},
            // Diamonds: K 9 4
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            // Clubs: A 8 3
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Eight}, {Suit: Clubs, Rank: Three},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C strong
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South}) // 1D waiting/negative

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != 4 { // 2NT
            t.Errorf("Expected 2NT rebid, got %s", bid)
        }
    })

    // 1C (strong) - 1D; strong with long clubs -> 2C
    t.Run("Strong with long clubs rebids 2C", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Spades: A 7
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Seven},
            // Hearts: K Q 9
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Nine},
            // Diamonds: A 7 3
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Seven}, {Suit: Diamonds, Rank: Three},
            // Clubs: K Q J 10 5 (5+ clubs)
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack}, {Suit: Clubs, Rank: Ten}, {Suit: Clubs, Rank: Five},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Clubs {
            t.Errorf("Expected 2C rebid, got %s", bid)
        }
    })

    // 1C (strong) - 1D; strong with 4-card heart suit -> 2H
    t.Run("Strong with 4 hearts rebids 2H", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Spades: A K 7
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Seven},
            // Hearts: A Q 9 4 (4 hearts)
            {Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Nine}, {Suit: Hearts, Rank: Four},
            // Diamonds: K 9
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine},
            // Clubs: K Q 3
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Three},
        })

	})

	// Test case: Responder has 7+ HCP and a 4-card major, should bid the major.
	t.Run("Responder bids major with positive hand", func(t *testing.T) {
		responder := NewPlayer(South)
		responder.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Jack},
			{Suit: Hearts, Rank: Two}, {Suit: Hearts, Rank: Three}, {Suit: Hearts, Rank: Four},
			{Suit: Diamonds, Rank: Two}, {Suit: Diamonds, Rank: Three}, {Suit: Diamonds, Rank: Four},
			{Suit: Clubs, Rank: Two}, {Suit: Clubs, Rank: Three}, {Suit: Clubs, Rank: Four},
		})

		auction := NewAuction()
		auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})

		bid := responder.MakeBid(auction)
		if bid.Level != 1 || bid.Strain != Spades {
			t.Errorf("Expected 1S response, but got %s", bid)
		}
	})

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
		auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})

		bid := responder.MakeBid(auction)
		if bid.Level != 1 || bid.Strain != Diamonds {
			t.Errorf("Expected 1D negative response, but got %s", bid)
		}
	})
}

func TestAI_PolishClub_Continuations(t *testing.T) {
	// 1C - 1H; opener raises to 3H with 4-card support
	t.Run("Opener raises to 3H with 4-card support", func(t *testing.T) {
		opener := NewPlayer(North)
		opener.Hand = NewHand([]Card{
			{Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Eight}, {Suit: Spades, Rank: Seven},
			{Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Seven}, {Suit: Hearts, Rank: Four},
			{Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Eight}, {Suit: Diamonds, Rank: Three},
			{Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Nine}, {Suit: Clubs, Rank: Five},
		})
        opener.Hand = NewHand([]Card{
            // Spades
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Eight}, {Suit: Spades, Rank: Seven},
            // Hearts (4): A K 7 4 -> 7 HCP
            {Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Seven}, {Suit: Hearts, Rank: Four},
            // Diamonds
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Eight}, {Suit: Diamonds, Rank: Three},
            // Clubs
            {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Nine}, {Suit: Clubs, Rank: Five},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C
        auction.AddBid(Bid{Level: 1, Strain: Hearts, Position: South}) // 1H positive

        bid := opener.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Hearts {
            t.Errorf("Expected 3H raise, got %s", bid)
        }
    })

    // 1C - 1S; opener no support, balanced minimum -> 1NT
    t.Run("Opener rebids 1NT with balanced minimum and no support", func(t *testing.T) {
        opener := NewPlayer(North)
        // Balanced 12 HCP, spade support <=2
        opener.Hand = NewHand([]Card{
            // Spades: Q 6 -> 2 HCP
            {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Six},
            // Hearts: K 7 5 -> 3 HCP
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Seven}, {Suit: Hearts, Rank: Five},
            // Diamonds: K 9 4 -> 3 HCP
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            // Clubs: A 8 3 -> 4 HCP
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Eight}, {Suit: Clubs, Rank: Three},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C
        auction.AddBid(Bid{Level: 1, Strain: Spades, Position: South}) // 1S

        bid := opener.MakeBid(auction)
        if bid.Level != 1 || bid.Strain != 4 { // 1NT
            t.Errorf("Expected 1NT rebid, got %s", bid)
        }
    })

    // 1C - 1S; opener unbalanced with long clubs -> 2C
    t.Run("Opener rebids 2C to show long clubs", func(t *testing.T) {
        opener := NewPlayer(North)
        // 5+ clubs, unbalanced, not balanced
        opener.Hand = NewHand([]Card{
            // Spades: 9 4
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Four},
            // Hearts: 9 4 2
            {Suit: Hearts, Rank: Nine}, {Suit: Hearts, Rank: Four}, {Suit: Hearts, Rank: Two},
            // Diamonds: A 7 3 -> 4 HCP
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Seven}, {Suit: Diamonds, Rank: Three},
            // Clubs: K Q J 9 5 -> 7 HCP
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack}, {Suit: Clubs, Rank: Nine}, {Suit: Clubs, Rank: Five},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C
        auction.AddBid(Bid{Level: 1, Strain: Spades, Position: South}) // 1S

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Clubs {
            t.Errorf("Expected 2C rebid, got %s", bid)
        }
    })

    // 1C - 1NT; opener with long clubs -> 2C
    t.Run("Opener follows 1NT response with 2C when long clubs", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Spades
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Four},
            // Hearts
            {Suit: Hearts, Rank: Nine}, {Suit: Hearts, Rank: Four}, {Suit: Hearts, Rank: Two},
            // Diamonds
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Seven}, {Suit: Diamonds, Rank: Three},
            // Clubs 5+
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack}, {Suit: Clubs, Rank: Nine}, {Suit: Clubs, Rank: Five},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C
        auction.AddBid(Bid{Level: 1, Strain: 4, Position: South})     // 1NT

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Clubs {
            t.Errorf("Expected 2C follow-up, got %s", bid)
        }
    })

    // 1NT opener; responder uses transfer 2D; opener must bid 2H
    t.Run("Opener completes Jacoby transfer 2D->2H", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Balanced with at least 2 hearts
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Ten}, {Suit: Spades, Rank: Nine},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Eight},
            {Suit: Diamonds, Rank: Queen}, {Suit: Diamonds, Rank: Ten}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten}, {Suit: Clubs, Rank: Nine},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: 4, Position: North})      // 1NT
        auction.AddBid(Bid{Level: 2, Strain: Diamonds, Position: South}) // 2D transfer

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Hearts {
            t.Errorf("Expected 2H completion of transfer, got %s", bid)
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
