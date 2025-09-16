package game

import "fmt"

// Bid represents a single bid in the auction
type Bid struct {
	Level    int  // 1-7
	Strain   Suit // Clubs, Diamonds, Hearts, Spades, or NoTrump
	Pass     bool // True if this is a pass
	Double   bool // True if this is a double
	Redouble bool // True if this is a redouble
}

// NewPass creates a new pass bid
func NewPass() Bid {
	return Bid{Pass: true}
}

// NewDouble creates a new double bid
func NewDouble() Bid {
	return Bid{Double: true}
}

// NewRedouble creates a new redouble bid
func NewRedouble() Bid {
	return Bid{Redouble: true}
}

// NewBid creates a new contract bid
func NewBid(level int, strain Suit) Bid {
	return Bid{
		Level:  level,
		Strain: strain,
	}
}

// String returns a string representation of the bid
func (b Bid) String() string {
	switch {
	case b.Pass:
		return "Pass"
	case b.Double:
		return "Double"
	case b.Redouble:
		return "Redouble"
	default:
		strainNames := []string{"♣", "♦", "♥", "♠", "NT"}
		strain := strainNames[b.Strain]
		if b.Strain == 4 { // No Trump
			strain = "NT"
		}
		return fmt.Sprintf("%d%s", b.Level, strain)
	}
}

// Auction represents the bidding sequence
type Auction struct {
	Bids []Bid
}

// NewAuction creates a new auction
func NewAuction() *Auction {
	return &Auction{
		Bids: make([]Bid, 0, 20), // Shouldn't need more than 20 bids
	}
}

// AddBid adds a bid to the auction
func (a *Auction) AddBid(bid Bid) {
	a.Bids = append(a.Bids, bid)
}

// LastNonPassBid returns the last bid that wasn't a pass
func (a *Auction) LastNonPassBid() (Bid, bool) {
	for i := len(a.Bids) - 1; i >= 0; i-- {
		if !a.Bids[i].Pass && !a.Bids[i].Double && !a.Bids[i].Redouble {
			return a.Bids[i], true
		}
	}
	return Bid{}, false
}

// IsValidBid checks if a bid is valid given the current auction state
func (a *Auction) IsValidBid(bid Bid) bool {
	if bid.Pass || bid.Double || bid.Redouble {
		return true
	}

	lastBid, found := a.LastNonPassBid()
	if !found {
		return true // First bid is always valid
	}

	// Calculate bid values for comparison
	bidValue := bid.Level*5 + int(bid.Strain)
	lastBidValue := lastBid.Level*5 + int(lastBid.Strain)

	return bidValue > lastBidValue
}

// IsAuctionComplete checks if the auction is complete
// The auction is complete when there are three consecutive passes after the opening bid
func (a *Auction) IsAuctionComplete() bool {
	if len(a.Bids) < 4 {
		return false
	}

	// Check last 3 bids are passes
	for i := len(a.Bids) - 1; i >= len(a.Bids)-3; i-- {
		if !a.Bids[i].Pass {
			return false
		}
	}

	// Make sure there was at least one non-pass bid
	for _, bid := range a.Bids {
		if !bid.Pass {
			return true
		}
	}

	return false
}
