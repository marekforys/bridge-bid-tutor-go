package game


// Position represents a player's position at the table
type Position int

const (
	North Position = iota
	East
	South
	West
)

// String returns the string representation of a Position.
func (p Position) String() string {
	switch p {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "Unknown"
	}
}

// Player represents a bridge player
type Player struct {
	Position Position
	Hand     *Hand
}

// NewPlayer creates a new player with the given position
func NewPlayer(pos Position) *Player {
	return &Player{
		Position: pos,
		Hand:     &Hand{},
	}
}

// Deal deals cards to the player
func (p *Player) Deal(cards []Card) {
	p.Hand = NewHand(cards)
}

// Add a helper to get the partner's position
func (p *Player) Partner() *Player {
	// This is a placeholder. In a real game, you'd have a reference
	// to the other players, probably from the Game struct.
	return nil
}

// Partner returns the position of the player's partner.
func (p Position) Partner() Position {
	return (p + 2) % 4
}

// MakeBid determines the bid for a computer player.
func (p *Player) MakeBid(auction *Auction) Bid {
	hcp, distribution := p.Hand.Evaluate()

	// Find our last bid and our partner's last bid.
	var myLastBid, partnerLastBid *Bid
	for i := len(auction.Bids) - 1; i >= 0; i-- {
		b := &auction.Bids[i]
		if myLastBid == nil && b.Position == p.Position {
			myLastBid = b
		}
		if partnerLastBid == nil && b.Position == p.Position.Partner() {
			partnerLastBid = b
		}
		if myLastBid != nil && partnerLastBid != nil {
			break // Found both, can stop searching.
		}
	}

	// Determine the bidding context.
	partnerIsLastBidder := false
	if len(auction.Bids) > 0 {
		lastBidder := auction.Bids[len(auction.Bids)-1].Position
		partnerIsLastBidder = lastBidder == p.Position.Partner()
	}

	if myLastBid == nil {
		// We haven't bid yet.
		if partnerLastBid == nil || !partnerIsLastBidder {
			// It's our turn to open for the partnership.
			return p.makeOpeningBid(auction, hcp, distribution)
		} else {
			// Partner opened, and it's our turn to respond.
			return p.makeResponseBid(auction, partnerLastBid, hcp, distribution)
		}
	} else {
		// We have bid before.
		if partnerIsLastBidder {
			// Partner just responded to our bid, so we must rebid.
			return p.makeRebid(auction, myLastBid, partnerLastBid, hcp, distribution)
		}
		// An opponent has bid. For now, we will just pass.
		// More advanced competitive bidding logic would go here.
	}

	// Default case, should not be reached in normal play.
	return NewPass()
}

// makeOpeningBid handles the logic for making an opening bid using Polish Club principles.
func (p *Player) makeOpeningBid(auction *Auction, hcp int, distribution map[Suit]int) Bid {
	isBalanced := p.Hand.IsBalanced()

	// --- Polish Club 1♣ Opening ---
	// 1. Strong hand: 18+ HCP, any shape.
	// 2. Weak balanced: 11-14 HCP, no 5-card major.
	openOneClub := false
	if hcp >= 18 {
		openOneClub = true
	} else if hcp >= 11 && hcp <= 14 && isBalanced && distribution[Hearts] < 5 && distribution[Spades] < 5 {
		openOneClub = true
	}
	if openOneClub {
		bid := NewBid(1, Clubs)
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	// --- 1NT Opening (15-17 HCP, balanced) ---
	if isBalanced && hcp >= 15 && hcp <= 17 {
		bid := NewBid(1, 4) // 4 represents NoTrump
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	// --- Major Suit Openings (5+ cards, 11-17 HCP) ---
	if hcp >= 11 && hcp <= 17 {
		if distribution[Spades] >= 5 {
			bid := NewBid(1, Spades)
			if auction.IsValidBid(bid) {
				return bid
			}
		}
		if distribution[Hearts] >= 5 {
			bid := NewBid(1, Hearts)
			if auction.IsValidBid(bid) {
				return bid
			}
		}
	}

	// --- 1♦ Opening (Natural, 4+ cards, 11-17 HCP) ---
	if hcp >= 11 && hcp <= 17 && distribution[Diamonds] >= 4 {
		bid := NewBid(1, Diamonds)
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	return NewPass()
}

// makeResponseBid handles the logic for responding to a partner's bid using Polish Club principles.
func (p *Player) makeResponseBid(auction *Auction, partnerBid *Bid, hcp int, distribution map[Suit]int) Bid {
	// --- Responses to 1♣ Opening ---
	if partnerBid.Level == 1 && partnerBid.Strain == Clubs {
		// Negative response: 0-6 HCP.
		if hcp <= 6 {
			bid := NewBid(1, Diamonds)
			if auction.IsValidBid(bid) {
				return bid
			}
		}
		// Positive responses: 7+ HCP.
		if hcp >= 7 {
			// Show a 4+ card major if available.
			if distribution[Spades] >= 4 {
				bid := NewBid(1, Spades)
				if auction.IsValidBid(bid) {
					return bid
				}
			}
			if distribution[Hearts] >= 4 {
				bid := NewBid(1, Hearts)
				if auction.IsValidBid(bid) {
					return bid
				}
			}
			// Balanced hand with no major.
			if hcp <= 10 && p.Hand.IsBalanced() {
				bid := NewBid(1, 4) // 1NT
				if auction.IsValidBid(bid) {
					return bid
				}
			}
		}
	}


	// --- Responses to 1NT Opening ---
	if partnerBid.Level == 1 && partnerBid.Strain == 4 { // Partner opened 1NT
		// Jacoby Transfers: Check for a 5-card major and 6+ HCP.
		if hcp >= 6 {
			if distribution[Hearts] >= 5 {
				bid := NewBid(2, Diamonds) // Transfer to Hearts
				if auction.IsValidBid(bid) {
					return bid
				}
			}
			if distribution[Spades] >= 5 {
				bid := NewBid(2, Hearts) // Transfer to Spades
				if auction.IsValidBid(bid) {
					return bid
				}
			}
		}

		// Stayman Convention: Check for 8+ HCP and a 4-card major.
		if hcp >= 8 && (distribution[Hearts] >= 4 || distribution[Spades] >= 4) {
			bid := NewBid(2, Clubs)
			if auction.IsValidBid(bid) {
				return bid
			}
		}
	}

	// --- Responses by Opener (after partner's response) ---
	// Response to Jacoby Transfer.
	if partnerBid.Level == 2 && partnerBid.Strain == Diamonds { // Transfer to Hearts
		return NewBid(2, Hearts)
	}
	if partnerBid.Level == 2 && partnerBid.Strain == Hearts { // Transfer to Spades
		return NewBid(2, Spades)
	}

	// Response to Stayman.
	if partnerBid.Level == 2 && partnerBid.Strain == Clubs {
		if distribution[Hearts] >= 4 {
			return NewBid(2, Hearts)
		}
		if distribution[Spades] >= 4 {
			return NewBid(2, Spades)
		}
		return NewBid(2, Diamonds) // No 4-card major.
	}

	// Simple support logic for suit openings.
	if hcp >= 6 && hcp <= 9 && distribution[partnerBid.Strain] >= 3 {
		bid := NewBid(partnerBid.Level+1, partnerBid.Strain)
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	return NewPass()
}


// makeRebid handles the logic for making a rebid after our partner has responded.
func (p *Player) makeRebid(auction *Auction, myLastBid, partnerLastBid *Bid, hcp int, distribution map[Suit]int) Bid {
	// --- Opener's Rebid after 1♣ - 1♦ ---
	if myLastBid.Level == 1 && myLastBid.Strain == Clubs && partnerLastBid.Level == 1 && partnerLastBid.Strain == Diamonds {
		// We opened 1♣ and partner responded with a negative 1♦.
		// Now we must clarify our opening.
		if hcp >= 11 && hcp <= 14 {
			// This was the "weak" variant of the 1♣ opening.
			if p.Hand.IsBalanced() {
				return NewBid(1, 4) // Rebid 1NT to show a balanced 11-14 HCP.
			}
			// Add logic for other weak, unbalanced hand types here.
		}
		// Add logic for strong hand (18+ HCP) rebids here.
	}

	// --- Opener's Rebid after 1NT opening ---
	if myLastBid.Level == 1 && myLastBid.Strain == 4 { // We opened 1NT
		// Respond to Jacoby transfers from partner.
		if partnerLastBid.Level == 2 && partnerLastBid.Strain == Diamonds { // Transfer to Hearts
			return NewBid(2, Hearts)
		}
		if partnerLastBid.Level == 2 && partnerLastBid.Strain == Hearts { // Transfer to Spades
			return NewBid(2, Spades)
		}

		// Respond to Stayman (2C): show a 4-card major if present, else 2D.
		if partnerLastBid.Level == 2 && partnerLastBid.Strain == Clubs {
			if distribution[Hearts] >= 4 {
				return NewBid(2, Hearts)
			}
			if distribution[Spades] >= 4 {
				return NewBid(2, Spades)
			}
			return NewBid(2, Diamonds)
		}
	}

	return NewPass()
}

// IsHuman returns true if the player is human
func (p *Player) IsHuman() bool {
	// For now, only South is human
	return p.Position == South
}
