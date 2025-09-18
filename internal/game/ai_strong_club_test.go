package game

import "testing"

// Strong 1C continuations after 1C–1D (18+ HCP routes)
func TestAI_StrongClubContinuations(t *testing.T) {
    // 1C (strong) - 1D; balanced 18-19 -> 2NT
    t.Run("Strong balanced 18-19 rebids 2NT", func(t *testing.T) {
        opener := NewPlayer(North)
        // 18-19 HCP balanced: AQx, KQx, Kxx, Axx ~ 18 HCP
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Seven},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Six},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Eight}, {Suit: Clubs, Rank: Three},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North}) // 1C strong
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South}) // 1D waiting/negative

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != 4 { // 2NT
            t.Fatalf("Expected 2NT rebid, got %s", bid)
        }
    })

    // 1C (strong) - 1D; strong with long clubs and unbalanced -> 2C
    t.Run("Strong with long clubs rebids 2C", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Make it clearly unbalanced: 6+ clubs and a singleton somewhere
            {Suit: Spades, Rank: Ace},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Seven}, {Suit: Diamonds, Rank: Three},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack}, {Suit: Clubs, Rank: Ten}, {Suit: Clubs, Rank: Five}, {Suit: Clubs, Rank: Four},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Clubs {
            t.Fatalf("Expected 2C rebid, got %s", bid)
        }
    })

    // 1C (strong) - 1D; strong with 4-card heart suit -> 2H
    t.Run("Strong with 4 hearts rebids 2H", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Seven},
            {Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Nine}, {Suit: Hearts, Rank: Four},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Three},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Hearts {
            t.Fatalf("Expected 2H rebid, got %s", bid)
        }
    })

    // 1C (strong) - 1D; strong with long diamonds and unbalanced -> 2D default
    t.Run("Strong with diamonds rebids 2D", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            // Unbalanced with 5+ diamonds
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Queen}, {Suit: Diamonds, Rank: Jack}, {Suit: Diamonds, Rank: Ten}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Nine}, {Suit: Clubs, Rank: Four},
        })

        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Diamonds {
            t.Fatalf("Expected 2D rebid, got %s", bid)
        }
    })
}
