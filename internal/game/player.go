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

// MakeBid determines the bid for a computer player
func (p *Player) MakeBid(auction *Auction) Bid {
	hcp, distribution := p.Hand.Evaluate()

	// Find the last bid from our partner
	var partnerBid *Bid
	for i := len(auction.Bids) - 1; i >= 0; i-- {
		if auction.Bids[i].Position == p.Position.Partner() {
			partnerBid = &auction.Bids[i]
			break
		}
	}

	// --- Response Bidding Logic ---
	if partnerBid != nil && !partnerBid.Pass {
		// Partner has opened, let's respond
		if hcp >= 6 && hcp <= 9 {
			// Support partner's suit if we have a fit (3+ cards)
			if distribution[partnerBid.Strain] >= 3 {
				bid := NewBid(partnerBid.Level+1, partnerBid.Strain)
				if auction.IsValidBid(bid) {
					return bid
				}
			}
		}
		// More response logic can be added here...
	}

	// --- Opening Bidding Logic ---

	// Rule of 20: If HCP + length of two longest suits is >= 20, open.
	var suitLengths []int
	for s := Clubs; s <= Spades; s++ {
		suitLengths = append(suitLengths, distribution[s])
	}
	sort.Sort(sort.Reverse(sort.IntSlice(suitLengths)))
	ruleOf20 := hcp + suitLengths[0] + suitLengths[1]

	if ruleOf20 >= 20 {
		// Find the longest suit to open with
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

	// Default to passing if no other bid is made
	return NewPass()
}

// IsHuman returns true if the player is human
func (p *Player) IsHuman() bool {
	// For now, only South is human
	return p.Position == South
}
