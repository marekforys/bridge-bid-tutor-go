package game

import "sort"

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

	// Determine if our partner has bid before.
	var partnerBid *Bid
	for i := len(auction.Bids) - 1; i >= 0; i-- {
		bid := auction.Bids[i]
		if bid.Position == p.Position.Partner() {
			if !bid.Pass {
				partnerBid = &bid
				break
			}
		} else if bid.Position == p.Position {
			// We have bid since our partner's last bid.
			break
		}
	}

	// If partner has made a bid, we respond.
	if partnerBid != nil {
		return p.makeResponseBid(auction, partnerBid, hcp, distribution)
	}

	// If we are the first to bid in the partnership, we open.
	return p.makeOpeningBid(auction, hcp, distribution)
}

// makeOpeningBid handles the logic for making an opening bid.
func (p *Player) makeOpeningBid(auction *Auction, hcp int, distribution map[Suit]int) Bid {
	// Check for a balanced hand and 15-17 HCP to open 1NT.
	isBalanced := p.Hand.IsBalanced()
	if isBalanced && hcp >= 15 && hcp <= 17 {
		bid := NewBid(1, 4) // 4 represents NoTrump
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	// Rule of 20 for other openings.
	var suitLengths []int
	for s := Clubs; s <= Spades; s++ {
		suitLengths = append(suitLengths, distribution[s])
	}
	sort.Sort(sort.Reverse(sort.IntSlice(suitLengths)))
	if hcp+suitLengths[0]+suitLengths[1] >= 20 {
		// Find the longest suit to open.
		// In case of a tie, we can add more sophisticated rules later.
		longestSuit := Clubs
		maxLength := 0
		for s := Clubs; s <= Spades; s++ {
			if distribution[s] > maxLength {
				maxLength = distribution[s]
				longestSuit = s
			}
		}
		bid := NewBid(1, longestSuit)
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	return NewPass()
}

// makeResponseBid handles the logic for responding to a partner's bid.
func (p *Player) makeResponseBid(auction *Auction, partnerBid *Bid, hcp int, distribution map[Suit]int) Bid {
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

	// --- Responses by 1NT Opener ---
	// This is a simplified check. A full implementation would verify we opened 1NT.
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

	// Simple support logic (can be expanded).
	if hcp >= 6 && hcp <= 9 && distribution[partnerBid.Strain] >= 3 {
		bid := NewBid(partnerBid.Level+1, partnerBid.Strain)
		if auction.IsValidBid(bid) {
			return bid
		}
	}

	return NewPass()
}


// IsHuman returns true if the player is human
func (p *Player) IsHuman() bool {
	// For now, only South is human
	return p.Position == South
}
