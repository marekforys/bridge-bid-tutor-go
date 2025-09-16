package game

import (
	"testing"
)

func TestAuction_IsValidBid(t *testing.T) {
	tests := []struct {
		name    string
		auction *Auction
		bid     Bid
		want    bool
	}{
		{
			name:    "First bid is always valid",
			auction: NewAuction(),
			bid:     NewBid(1, Clubs),
			want:    true,
		},
		{
			name: "Higher bid is valid",
			auction: &Auction{Bids: []Bid{NewBid(1, Clubs)}},
			bid:     NewBid(1, Diamonds),
			want:    true,
		},
		{
			name: "Lower bid is invalid",
			auction: &Auction{Bids: []Bid{NewBid(1, Diamonds)}},
			bid:     NewBid(1, Clubs),
			want:    false,
		},
		{
			name: "Same bid is invalid",
			auction: &Auction{Bids: []Bid{NewBid(1, Spades)}},
			bid:     NewBid(1, Spades),
			want:    false,
		},
		{
			name: "Pass is always valid",
			auction: &Auction{Bids: []Bid{NewBid(7, Spades)}},
			bid:     NewPass(),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.auction.IsValidBid(tt.bid); got != tt.want {
				t.Errorf("Auction.IsValidBid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuction_IsAuctionComplete(t *testing.T) {
	tests := []struct {
		name    string
		auction *Auction
		want    bool
	}{
		{
			name:    "Not enough bids",
			auction: &Auction{Bids: []Bid{NewBid(1, Clubs), NewPass(), NewPass()}},
			want:    false,
		},
		{
			name: "Three passes after opening bid",
			auction: &Auction{Bids: []Bid{NewBid(1, Clubs), NewPass(), NewPass(), NewPass()}},
			want:    true,
		},
		{
			name: "Four passes in a row from start",
			auction: &Auction{Bids: []Bid{NewPass(), NewPass(), NewPass(), NewPass()}},
			want:    false, // This is a pass-out, not a completed auction with a contract.
		},
		{
			name: "Interrupted passes",
			auction: &Auction{Bids: []Bid{NewBid(1, Clubs), NewPass(), NewBid(1, Diamonds), NewPass(), NewPass()}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.auction.IsAuctionComplete(); got != tt.want {
				t.Errorf("Auction.IsAuctionComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}
