# Bridge Bidding Tutor in Go

A console-based bridge bidding simulator that helps you practice and learn bridge bidding conventions.

## Features

- Interactive command-line interface
- Basic AI opponents with simple bidding strategies
- Hand evaluation (High Card Points)
- Support for all standard bridge bids, including passes, doubles, and redoubles

## Prerequisites

- Go 1.21 or higher

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bridge-bid-tutor-go.git
   cd bridge-bid-tutor-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

1. Build and run the game:
   ```bash
   go run cmd/bridge/main.go
   ```

2. Follow the on-screen instructions to place your bids.

## How to Play

- You play as South (your hand will be displayed)
- The other three positions (North, East, West) are controlled by the computer
- When it's your turn, enter your bid:
  - To bid: Enter the level followed by the suit (e.g., `1H`, `2NT`, `3C`)
  - To pass: Type `pass` or `p`
  - To double: Type `double`, `dbl`, or `x`
  - To redouble: Type `redouble`, `rdbl`, or `xx`

## Bidding Conventions

The AI uses a simple bidding system:
- 20+ HCP: Open 2♣ (strong hand)
- 15-19 HCP: Open 1 of a suit (longest suit)
- < 15 HCP: Pass

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
