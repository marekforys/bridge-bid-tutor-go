package game

import (
	"fmt"
	"math/rand"
)

// Suit represents the four suits in a deck of cards
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

// Rank represents the rank of a card
type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

// Card represents a playing card with a suit and rank
type Card struct {
	Suit Suit
	Rank Rank
}

// String returns a string representation of the card
func (c Card) String() string {
	suitSymbols := []string{"C", "D", "H", "S"}
	rankSymbols := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	return fmt.Sprintf("%s%s", rankSymbols[c.Rank-2], suitSymbols[c.Suit])
}

// RankString returns a string representation of the card's rank.
func (c Card) RankString() string {
	rankSymbols := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	return rankSymbols[c.Rank-2]
}

// Deck represents a deck of 52 playing cards
type Deck []Card

// NewDeck creates and returns a new shuffled deck of 52 cards
func NewDeck() Deck {
	var deck Deck
	for s := Clubs; s <= Spades; s++ {
		for r := Two; r <= Ace; r++ {
			deck = append(deck, Card{Suit: s, Rank: r})
		}
	}
	return deck
}

// Shuffle randomizes the order of cards in the deck
func (d Deck) Shuffle() {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

// Deal deals a specified number of cards from the deck
func (d *Deck) Deal(n int) []Card {
	if n > len(*d) {
		n = len(*d)
	}
	cards := (*d)[:n]
	*d = (*d)[n:]
	return cards
}
