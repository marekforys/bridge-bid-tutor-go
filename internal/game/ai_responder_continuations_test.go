package game

import "testing"

// Responder continuations after 1C-1D and opener's strong rebid
func TestAI_ResponderContinuations_AfterStrongRebids(t *testing.T) {
    // 1C-1D-2NT; responder with 8 HCP -> 3NT
    t.Run("After 2NT, responder bids 3NT with 8+ HCP", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            // 9 HCP: QS(2) + KH(3) + AC(4)
            {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Seven},
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Six},
            {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Three},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: 4, Position: North}) // 2NT strong balanced

        bid := responder.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != 4 {
            t.Fatalf("Expected 3NT, got %s", bid)
        }
    })

    // 1C-1D-2C; responder with 4 hearts -> 2H
    t.Run("After 2C, responder shows 4-card heart suit", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Eight}, {Suit: Hearts, Rank: Four},
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Four},
            {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Three},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Clubs, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Hearts {
            t.Fatalf("Expected 2H, got %s", bid)
        }
    })

    // 1C-1D-2C; responder no major -> 2D waiting
    t.Run("After 2C, responder bids 2D waiting if no major", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Eight}, {Suit: Hearts, Rank: Four},
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Four},
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Three},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Clubs, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != Diamonds {
            t.Fatalf("Expected 2D waiting, got %s", bid)
        }
    })

    // 1C-1D-2H; responder with 3 hearts -> 3H
    t.Run("After 2H, responder raises with 3-card support", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: King}, {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Four},
            {Suit: Spades, Rank: Nine}, {Suit: Spades, Rank: Four},
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Hearts, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Hearts {
            t.Fatalf("Expected 3H raise, got %s", bid)
        }
    })

    // 1C-1D-2H; responder without support but with 6+ HCP -> 2NT
    t.Run("After 2H, responder bids 2NT with no support and 6+ HCP", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Hearts, Rank: Ten}, {Suit: Hearts, Rank: Four},
            {Suit: Spades, Rank: Queen}, {Suit: Spades, Rank: Seven},
            {Suit: Diamonds, Rank: Ace}, {Suit: Diamonds, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Hearts, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != 4 {
            t.Fatalf("Expected 2NT, got %s", bid)
        }
    })

    // 1C-1D-2D; responder with 3 diamonds and 6+ HCP -> 3D
    t.Run("After 2D, responder raises with diamond support", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Diamonds, Rank: Queen}, {Suit: Diamonds, Rank: Jack}, {Suit: Diamonds, Rank: Four},
            {Suit: Hearts, Rank: Ten},
            {Suit: Spades, Rank: Nine},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Diamonds, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 3 || bid.Strain != Diamonds {
            t.Fatalf("Expected 3D, got %s", bid)
        }
    })

    // 1C-1D-2D; responder without support but 6+ HCP -> 2NT
    t.Run("After 2D, responder bids 2NT without support", func(t *testing.T) {
        responder := NewPlayer(South)
        responder.Hand = NewHand([]Card{
            {Suit: Diamonds, Rank: Nine}, {Suit: Diamonds, Rank: Four},
            {Suit: Hearts, Rank: King},
            {Suit: Spades, Rank: Queen},
            {Suit: Clubs, Rank: Ace}, {Suit: Clubs, Rank: Ten},
        })
        auction := NewAuction()
        auction.AddBid(Bid{Level: 1, Strain: Clubs, Position: North})
        auction.AddBid(Bid{Level: 1, Strain: Diamonds, Position: South})
        auction.AddBid(Bid{Level: 2, Strain: Diamonds, Position: North})

        bid := responder.MakeBid(auction)
        if bid.Level != 2 || bid.Strain != 4 {
            t.Fatalf("Expected 2NT, got %s", bid)
        }
    })
}
