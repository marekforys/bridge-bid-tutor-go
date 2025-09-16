package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yourusername/bridge-bid-tutor-go/internal/game"
)

func main() {
	fmt.Println("Welcome to Bridge Bidding Tutor!")
	fmt.Println("------------------------------")

	// Initialize game
	game := NewGame()

	// Start the game loop
	if err := game.Start(); err != nil {
		log.Fatalf("Error starting game: %v", err)
	}
}

// Game represents the main game state
type Game struct {
	Deck    game.Deck
	Players []*game.Player
	Auction *game.Auction
	Dealer  game.Position
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Initialize deck and shuffle
	deck := game.NewDeck()
	deck.Shuffle()

	// Create players
	players := make([]*game.Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = game.NewPlayer(game.Position(i))
	}

	// Deal cards
	for i := 0; i < 52; i++ {
		players[i%4].Hand.Cards = append(players[i%4].Hand.Cards, deck[i])
	}

	// Sort each player's hand
	for _, p := range players {
		p.Hand.Sort()
	}

	return &Game{
		Deck:    deck,
		Players: players,
		Auction: game.NewAuction(),
		Dealer:  game.North, // First dealer is North
	}
}

// Start begins the game loop
func (g *Game) Start() error {
	// Game loop
	for !g.Auction.IsAuctionComplete() {
		// Get current player
		currentPlayer := g.Players[g.Dealer]

		// Display game state
		g.displayGameState(currentPlayer)

		// Get bid from player or AI
		var bid game.Bid
		if currentPlayer.IsHuman() {
			// Human player's turn
			prompt := promptui.Prompt{
				Label: "Enter your bid (e.g., '1H', 'pass', 'double')",
				Validate: func(input string) error {
					parsedBid, err := parseBid(input)
					if err != nil {
						return err
					}
					if !g.Auction.IsValidBid(parsedBid) {
						return fmt.Errorf("must be higher than previous bid")
					}
					return nil
				},
			}

			result, err := prompt.Run()
			if err != nil {
				// Handle user interruption (e.g., Ctrl+C)
				if err == promptui.ErrInterrupt {
					fmt.Println("\nGame aborted. Goodbye!")
					return nil
				}
				return fmt.Errorf("prompt failed: %w", err)
			}

			bid, _ = parseBid(result) // We can ignore the error here because validation already passed

		} else {
			// AI's turn
			bid = currentPlayer.MakeBid(g.Auction)
			fmt.Printf("\n%s bids: %s\n", currentPlayer.Position, bid)
		}

		// Add bid to auction
		bid.Position = currentPlayer.Position
		g.Auction.AddBid(bid)
		g.Dealer = (g.Dealer + 1) % 4

		// Add a newline for better readability
		fmt.Println()
	}

	// Auction is complete
	fmt.Println("\nAuction complete!")
	fmt.Println("Final contract:", g.Auction.Bids[len(g.Auction.Bids)-4])

	return nil
}

// displayGameState shows the current game state to the player
func (g *Game) displayGameState(currentPlayer *game.Player) {
	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Show auction history
	fmt.Println("Auction:")
	for _, bid := range g.Auction.Bids {
		fmt.Printf("%s: %s\n", bid.Position, bid)
	}
	fmt.Println()

	// Show player's hand
	if currentPlayer.IsHuman() {
		hcp, _ := currentPlayer.Hand.Evaluate()
		fmt.Printf("Your hand (HCP: %d):\n", hcp)
		fmt.Println("Spades:", currentPlayer.Hand.GetSuit(game.Spades))
		fmt.Println("Hearts:", currentPlayer.Hand.GetSuit(game.Hearts))
		fmt.Println("Diamonds:", currentPlayer.Hand.GetSuit(game.Diamonds))
		fmt.Println("Clubs:", currentPlayer.Hand.GetSuit(game.Clubs))
		fmt.Println()
	}
}

// parseBid converts a string input into a Bid struct
func parseBid(input string) (game.Bid, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "pass", "p":
		return game.NewPass(), nil
	case "double", "dbl", "x":
		return game.NewDouble(), nil
	case "redouble", "rdbl", "xx":
		return game.NewRedouble(), nil
	}

	// Parse contract bid (e.g., "1H", "3NT")
	if len(input) < 2 {
		return game.Bid{}, fmt.Errorf("invalid bid format")
	}

	// Parse level
	level := int(input[0] - '0')
	if level < 1 || level > 7 {
		return game.Bid{}, fmt.Errorf("bid level must be between 1 and 7")
	}

	// Parse strain
	strain := input[1:]
	var s game.Suit
	switch strings.ToUpper(strain) {
	case "C":
		s = game.Clubs
	case "D":
		s = game.Diamonds
	case "H":
		s = game.Hearts
	case "S":
		s = game.Spades
	case "NT", "N":
		s = 4 // No Trump
	default:
		return game.Bid{}, fmt.Errorf("invalid suit: %s", strain)
	}

	return game.NewBid(level, s), nil
}
