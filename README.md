# Bridge Bidding Tutor in Go

A console-based bridge bidding simulator that helps you practice and learn bridge bidding conventions.

## Features

- Interactive command-line interface for bidding.
- AI opponents with support for modern bidding conventions.
- Hand evaluation (High Card Points and distribution).
- End-of-auction summary showing all four hands for review.

## Prerequisites

- Go 1.21 or higher

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/marekforys/bridge-bid-tutor-go.git
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

- You play as South (your hand will be displayed).
- The other three positions (North, East, West) are controlled by the computer.
- When it's your turn, enter your bid:
  - To bid: Enter the level followed by the suit (e.g., `1H`, `2NT`, `3C`).
  - To pass: Type `pass` or `p`.
  - To double: Type `double`, `dbl`, or `x`.
  - To redouble: Type `redouble`, `rdbl`, or `xx`.

## Bidding System: Polish Club

The AI uses the **Polish Club** system, a sophisticated and powerful bidding system. Here are the core principles of its implementation:

- **Artificial 1♣ Opening**: This is the cornerstone of the system. It shows one of two hand types:
  - **Weak**: 11-14 HCP with a balanced or semi-balanced hand, without a 5-card major.
  - **Strong**: 18+ HCP with any shape.
- **1♦ Negative Response**: The primary response to a `1♣` opening is `1♦`, which is artificial and shows a weak hand (typically 0-6 HCP).
- **5-Card Majors**: An opening bid of `1♥` or `1♠` guarantees a 5-card suit.
- **1NT Opening**: Shows a balanced hand with 15-17 HCP.
- **Conventions over 1NT**: The AI still uses **Stayman** and **Jacoby Transfers** in response to a `1NT` opening.

## End-of-Auction Review

At the end of each auction, the game displays all four hands, allowing you to review the bidding in the context of the full deal.

```
Auction complete!
Final contract: 1S
------------------------------

--- All Hands ---

North (HCP: 5)
  Spades:   J 3
  Hearts:   J 10 8 6 5
  Diamonds: J 8 4
  Clubs:    Q 9 8

East (HCP: 11)
  Spades:   A 9 2
  Hearts:   K 9 3
  Diamonds: 10 6
  Clubs:    K J 10 7 2

South (HCP: 11)
  Spades:   Q 4
  Hearts:   7 2
  Diamonds: K Q 5 3
  Clubs:    A 6 5 4 3

West (HCP: 13)
  Spades:   K 10 8 7 6 5
  Hearts:   A Q 4
  Diamonds: A 9 7 2
  Clubs:
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
