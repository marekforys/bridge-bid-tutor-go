package game

import "testing"

// Tests for responder follow-ups after Puppet answers over 2NT
// Sequence: 1C (North) - 1D (South) - 2NT (North) - 3C (South, Puppet) - 3H/3S/3D (North) - ? (South)
func TestAI_PuppetFollowUps(t *testing.T) {
    buildBase := func() *Auction {
        a := NewAuction()
        a.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        a.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        a.AddBid(Bid{Level: 2, Strain: 4, Position: North})
        a.AddBid(Bid{Level: 3, Strain: Clubs, Position: South}) // Puppet
        return a
    }

    // 3H reply: with 3+ hearts and 13+ HCP -> 6H
    t.Run("3H reply -> 6H with 13+ and 3+ hearts", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: Ace}, {Suit: Hearts, Rank: Queen}, {Suit: Hearts, Rank: Four}, // 3+ hearts, 7 HCP
            {Suit: Spades, Rank: King}, // 3 HCP
            {Suit: Diamonds, Rank: King}, // 3 HCP
            {Suit: Clubs, Rank: Ace}, // 4 HCP
        }) // total 17 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Hearts, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 6 || bid.Strain != Hearts {
            t.Fatalf("Expected 6H, got %s", bid)
        }
    })

    // 3H reply: with 3+ hearts and 8-12 HCP -> 4H
    t.Run("3H reply -> 4H with 8-12 and 3+ hearts", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Five}, // 3 hearts, 3 HCP
            {Suit: Spades, Rank: Queen}, // 2 HCP
            {Suit: Diamonds, Rank: Ace}, // 4 HCP
        }) // total 9 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Hearts, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 4 || bid.Strain != Hearts {
            t.Fatalf("Expected 4H, got %s", bid)
        }
    })

    // 3H reply: with <3 hearts, prefer 3NT with values
    t.Run("3H reply -> 3NT with values and no fit", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: King}, // 3 HCP, only 1 heart
            {Suit: Spades, Rank: Queen}, // 2 HCP
            {Suit: Diamonds, Rank: Ace}, // 4 HCP
        }) // 9 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Hearts, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 3 || bid.Strain != 4 {
            t.Fatalf("Expected 3NT, got %s", bid)
        }
    })

    // 3S reply: with 3+ spades and 13+ HCP -> 6S
    t.Run("3S reply -> 6S with 13+ and 3+ spades", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Ace}, {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Four}, // 3 spades, 6 HCP
            {Suit: Hearts, Rank: King}, // 3 HCP
            {Suit: Diamonds, Rank: Ace}, // 4 HCP
            {Suit: Clubs, Rank: King}, // 3 HCP
        }) // 16 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Spades, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 6 || bid.Strain != Spades {
            t.Fatalf("Expected 6S, got %s", bid)
        }
    })

    // 3S reply: with 3+ spades and 8-12 HCP -> 4S
    t.Run("3S reply -> 4S with 8-12 and 3+ spades", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Spades, Rank: King}, {Suit: Spades, Rank: Ten}, {Suit: Spades, Rank: Five}, // 3 spades, 3 HCP
            {Suit: Hearts, Rank: Queen}, // 2 HCP
            {Suit: Diamonds, Rank: Ace}, // 4 HCP
        }) // 9 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Spades, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 4 || bid.Strain != Spades {
            t.Fatalf("Expected 4S, got %s", bid)
        }
    })

    // 3S reply: with <3 spades -> 3NT with values
    t.Run("3S reply -> 3NT with values and no fit", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Spades, Rank: King}, // 3 HCP, only 1 spade
            {Suit: Hearts, Rank: Queen}, // 2 HCP
            {Suit: Diamonds, Rank: Ace}, // 4 HCP
        }) // 9 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Spades, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 3 || bid.Strain != 4 {
            t.Fatalf("Expected 3NT, got %s", bid)
        }
    })

    // 3D denial: responder with 8+ HCP -> 3NT
    t.Run("3D denial -> 3NT with 8+ HCP", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Spades, Rank: King}, // 3
            {Suit: Hearts, Rank: Queen}, // 2
            {Suit: Diamonds, Rank: Ace}, // 4
        }) // 9 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Diamonds, Position: North})
        bid := responder.MakeBid(a)
        if bid.Level != 3 || bid.Strain != 4 {
            t.Fatalf("Expected 3NT, got %s", bid)
        }
    })

    // 3D denial: responder with <8 HCP -> Pass
    t.Run("3D denial -> Pass with <8 HCP", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Spades, Rank: Jack}, // 1 HCP
            {Suit: Hearts, Rank: Jack}, // 1 HCP
            {Suit: Diamonds, Rank: Queen}, // 2 HCP
            {Suit: Clubs, Rank: Queen}, // 2 HCP
        }) // 6 HCP
        a := buildBase()
        a.AddBid(Bid{Level: 3, Strain: Diamonds, Position: North})
        bid := responder.MakeBid(a)
        if !bid.Pass {
            t.Fatalf("Expected Pass, got %s", bid)
        }
    })
}
