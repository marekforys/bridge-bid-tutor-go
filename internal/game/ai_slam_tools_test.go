package game

import "testing"

// Tests for Puppet Stayman (3C) and Gerber (4C) over opener's 2NT after 1C-1D-2NT
func TestAI_SlamToolsOverTwoNT(t *testing.T) {
    // Helper: build an auction up to 1C-1D-2NT with opener on North and responder on South
    buildAuction := func() *Auction {
        a := NewAuction()
        a.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        a.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        a.AddBid(Bid{Level: 2, Strain: 4, Position: North}) // 2NT by opener (18-19 bal)
        return a
    }

    // Puppet: opener with 5 hearts should answer 3H
    t.Run("Puppet: opener shows 5 hearts as 3H", func(t *testing.T) {
        opener := NewPlayer(North)
        // 18-19 balanced with 5 hearts
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen},
            {Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Nine}, {Suit: Hearts, Rank: Seven}, {Suit: Hearts, Rank: Three},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Eight},
        })
        auction := buildAuction()
        // Responder bids Puppet 3C
        auction.AddBid(Bid{Level: 3, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Hearts {
            t.Fatalf("Expected 3H response to Puppet, got %s", bid)
        }
    })

    // Puppet: opener with 5 spades should answer 3S
    t.Run("Puppet: opener shows 5 spades as 3S", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Seven}, {Suit: Spades, Rank: Three},
            {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Nine},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Eight},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 3, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Spades {
            t.Fatalf("Expected 3S response to Puppet, got %s", bid)
        }
    })

    // Puppet: opener with no 5M should answer 3D
    t.Run("Puppet: opener denies 5M with 3D", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Nine},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Six}, {Suit: Hearts, Rank: Four},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            {Suit: Clubs, Rank: Ace},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 3, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Diamonds {
            t.Fatalf("Expected 3D denial to Puppet, got %s", bid)
        }
    })

    // Gerber: 0 aces -> 4D
    t.Run("Gerber: 0 aces -> 4D", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Queen},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Queen},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 4, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 4 || bid.Strain != Diamonds {
            t.Fatalf("Expected 4D for 0/4 aces, got %s", bid)
        }
    })

    // Gerber: 1 ace -> 4H
    t.Run("Gerber: 1 ace -> 4H", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Queen},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Queen},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack}, {Suit: Clubs, Rank: Ten},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 4, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 4 || bid.Strain != Hearts {
            t.Fatalf("Expected 4H for 1 ace, got %s", bid)
        }
    })

    // Gerber: 2 aces -> 4S
    t.Run("Gerber: 2 aces -> 4S", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace},
            {Suit: Hearts, Rank: Ace},
            {Suit: Diamonds, Rank: King}, {Suit: Diamonds, Rank: Queen},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 4, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 4 || bid.Strain != Spades {
            t.Fatalf("Expected 4S for 2 aces, got %s", bid)
        }
    })

    // Gerber: 3 aces -> 4NT
    t.Run("Gerber: 3 aces -> 4NT", func(t *testing.T) {
        opener := NewPlayer(North)
        opener.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace},
            {Suit: Hearts, Rank: Ace},
            {Suit: Diamonds, Rank: Ace},
            {Suit: Clubs, Rank: King}, {Suit: Clubs, Rank: Queen}, {Suit: Clubs, Rank: Jack},
        })
        auction := buildAuction()
        auction.AddBid(Bid{Level: 4, Strain: Clubs, Position: South})

        bid := opener.MakeBid(auction)
        if bid.Level != 4 || bid.Strain != 4 {
            t.Fatalf("Expected 4NT for 3 aces, got %s", bid)
        }
    })
}
