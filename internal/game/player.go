package game

// Position represents a player's position at the table
type Position int

const (
	North Position = iota
	East
	South
	West
)

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

// MakeBid determines the bid for a computer player
func (p *Player) MakeBid(auction *Auction) Bid {
	hcp, distribution := p.Hand.Evaluate()

	// Very basic bidding system for the computer
	switch {
	case hcp >= 20:
		// Strong hand, open 2♣
		return NewBid(2, Clubs)
	case hcp >= 15:
		// Open 1 of a suit
		// Find the longest suit
		longestSuit := Clubs
		maxLength := 0
		for s := Clubs; s <= Spades; s++ {
			if count := distribution[s]; count > maxLength {
				maxLength = count
				longestSuit = s
			}
		}
		return NewBid(1, longestSuit)
	default:
		// Pass with weak hands
		return NewPass()
	}
}

// IsHuman returns true if the player is human
func (p *Player) IsHuman() bool {
	// For now, only South is human
	return p.Position == South
}
