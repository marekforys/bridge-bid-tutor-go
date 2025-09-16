package game

// Vulnerability represents the vulnerability status of a partnership.
type Vulnerability bool

const (
	NotVulnerable Vulnerability = false
	Vulnerable    Vulnerability = true
)

// Score represents the score for a contract.
type Score struct {
	TrickScore int
	BonusScore int
	TotalScore int
	MadeGame   bool
	MadeSlam   bool
}

// CalculateScore calculates the score for a made contract.
// This is a simplified version that assumes the contract is made exactly.
func CalculateScore(contract Bid, vulnerability Vulnerability) Score {
	var score Score
	tricks := contract.Level

	// Calculate trick score
	switch contract.Strain {
	case Clubs, Diamonds:
		score.TrickScore = tricks * 20
	case Hearts, Spades:
		score.TrickScore = tricks * 30
	case 4: // NoTrump
		score.TrickScore = 40 + (tricks-1)*30
	}

	if contract.Double {
		score.TrickScore *= 2
	}
	if contract.Redouble {
		score.TrickScore *= 4
	}

	// Calculate bonus score
	if score.TrickScore >= 100 {
		score.MadeGame = true
		if vulnerability {
			score.BonusScore += 500
		} else {
			score.BonusScore += 300
		}
	} else {
		score.BonusScore += 50 // Part-score bonus
	}

	// Slam bonus
	if tricks == 6 { // Small slam
		if vulnerability {
			score.BonusScore += 750
		} else {
			score.BonusScore += 500
		}
		score.MadeSlam = true
	} else if tricks == 7 { // Grand slam
		if vulnerability {
			score.BonusScore += 1500
		} else {
			score.BonusScore += 1000
		}
		score.MadeSlam = true
	}

	// Doubled/Redoubled bonus for making the contract
	if contract.Double {
		score.BonusScore += 50
	}
	if contract.Redouble {
		score.BonusScore += 100
	}

	score.TotalScore = score.TrickScore + score.BonusScore
	return score
}
