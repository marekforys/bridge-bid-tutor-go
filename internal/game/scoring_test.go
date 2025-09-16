package game

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name        string
		contract    Bid
		vul         Vulnerability
		wantScore   Score
	}{
		{
			name:     "Part-score in minors",
			contract: Bid{Level: 2, Strain: Clubs},
			vul:      NotVulnerable,
			wantScore: Score{TrickScore: 40, BonusScore: 50, TotalScore: 90},
		},
		{
			name:     "Game in majors",
			contract: Bid{Level: 4, Strain: Hearts},
			vul:      NotVulnerable,
			wantScore: Score{TrickScore: 120, BonusScore: 300, TotalScore: 420, MadeGame: true},
		},
		{
			name:     "Game in NoTrump, vulnerable",
			contract: Bid{Level: 3, Strain: 4},
			vul:      Vulnerable,
			wantScore: Score{TrickScore: 100, BonusScore: 500, TotalScore: 600, MadeGame: true},
		},
		{
			name:     "Small slam in minors, vulnerable",
			contract: Bid{Level: 6, Strain: Diamonds},
			vul:      Vulnerable,
			wantScore: Score{TrickScore: 120, BonusScore: 500 + 750, TotalScore: 1370, MadeGame: true, MadeSlam: true},
		},
		{
			name:     "Grand slam in majors, doubled",
			contract: Bid{Level: 7, Strain: Spades, Double: true},
			vul:      NotVulnerable,
			wantScore: Score{TrickScore: 420, BonusScore: 300 + 1000 + 50, TotalScore: 1770, MadeGame: true, MadeSlam: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScore := CalculateScore(tt.contract, tt.vul)
			if gotScore.TotalScore != tt.wantScore.TotalScore {
				t.Errorf("CalculateScore() TotalScore = %v, want %v", gotScore.TotalScore, tt.wantScore.TotalScore)
			}
			if gotScore.MadeGame != tt.wantScore.MadeGame {
				t.Errorf("CalculateScore() MadeGame = %v, want %v", gotScore.MadeGame, tt.wantScore.MadeGame)
			}
			if gotScore.MadeSlam != tt.wantScore.MadeSlam {
				t.Errorf("CalculateScore() MadeSlam = %v, want %v", gotScore.MadeSlam, tt.wantScore.MadeSlam)
			}
		})
	}
}
