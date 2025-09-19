package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	gamepkg "github.com/marekforys/bridge-bid-tutor-go/internal/game"
	"github.com/google/uuid"
)

// Server holds HTTP state and session store
type Server struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// Session captures a single table's state
type Session struct {
	ID      string              `json:"id"`
	Players []*gamepkg.Player   `json:"-"`
	Auction *gamepkg.Auction    `json:"-"`
	Dealer  gamepkg.Position    `json:"-"`
}

// New constructs a new Server
func New() *Server {
	return &Server{sessions: make(map[string]*Session)}
}

// RegisterRoutes attaches handlers to the mux
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/sessions/", s.handleSessionByID)
	mux.HandleFunc("/api/evaluate-bid", s.handleEvaluateBid)
}

// handleSessions manages collection endpoints
// POST /api/sessions -> create a new session
func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sess := s.newSession()
	s.sessPut(sess)

	writeJSON(w, http.StatusCreated, s.serializeSession(sess))
}

// handleSessionByID manages single-session endpoints
// GET /api/sessions/{id}
// POST /api/sessions/{id}/bid
func (s *Server) handleSessionByID(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[0]

	sess, ok := s.sessGet(id)
	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Route based on remaining path
	if len(parts) == 1 {
		// /api/sessions/{id}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, s.serializeSession(sess))
		return
	}

	action := parts[1]
	switch action {
	case "bid":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s.handlePostBid(w, r, sess)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// newSession constructs a new session with a shuffled deck, dealt hands, and a fresh auction
func (s *Server) newSession() *Session {
	deck := gamepkg.NewDeck()
	deck.Shuffle()

	players := make([]*gamepkg.Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = gamepkg.NewPlayer(gamepkg.Position(i))
	}

	// Deal 52 cards round-robin
	for i := 0; i < 52; i++ {
		players[i%4].Hand.Cards = append(players[i%4].Hand.Cards, deck[i])
	}
	// Sort hands
	for _, p := range players {
		p.Hand.Sort()
	}

	id := uuid.New().String()
	return &Session{
		ID:      id,
		Players: players,
		Auction: gamepkg.NewAuction(),
		Dealer:  gamepkg.North,
	}
}

// handlePostBid submits a bid for the current dealer of the session
// Expects JSON: {"position":"North|East|South|West","bid":"3H|Pass|2NT|X|XX"}
func (s *Server) handlePostBid(w http.ResponseWriter, r *http.Request, sess *Session) {
	setCORSHeaders(w)
	var req struct {
		Position string `json:"position"`
		Bid      string `json:"bid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	pos, err := parsePosition(req.Position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Determine whose turn it is
	current := sess.Players[sess.Dealer]
	if current.Position != pos {
		http.Error(w, fmt.Sprintf("it's %s's turn", current.Position), http.StatusConflict)
		return
	}

	// Parse and validate bid
	bid, err := parseBid(req.Bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !sess.Auction.IsValidBid(bid) {
		http.Error(w, "invalid bid relative to auction", http.StatusBadRequest)
		return
	}

	bid.Position = current.Position
	sess.Auction.AddBid(bid)
	sess.Dealer = (sess.Dealer + 1) % 4

	writeJSON(w, http.StatusOK, s.serializeSession(sess))
}

// Serialization helpers
func (s *Server) serializeSession(sess *Session) map[string]any {
	bids := make([]map[string]any, 0, len(sess.Auction.Bids))
	for _, b := range sess.Auction.Bids {
		bids = append(bids, map[string]any{
			"position": b.Position.String(),
			"level":    b.Level,
			"strain":   s.strainString(b.Strain),
			"double":   b.Double,
			"redouble": b.Redouble,
			"pass":     b.Pass,
		})
	}

	players := []map[string]any{}
	for _, p := range sess.Players {
		hcp, _ := p.Hand.Evaluate()
		players = append(players, map[string]any{
			"position": p.Position.String(),
			"hcp":      hcp,
			"spades":   p.Hand.GetSuit(gamepkg.Spades),
			"hearts":   p.Hand.GetSuit(gamepkg.Hearts),
			"diamonds": p.Hand.GetSuit(gamepkg.Diamonds),
			"clubs":    p.Hand.GetSuit(gamepkg.Clubs),
		})
	}

	return map[string]any{
		"id":      sess.ID,
		"dealer":  sess.Dealer.String(),
		"players": players,
		"auction": bids,
		"complete": sess.Auction.IsAuctionComplete(),
	}
}

func (s *Server) strainString(strain gamepkg.Suit) string {
	switch strain {
	case gamepkg.Clubs:
		return "C"
	case gamepkg.Diamonds:
		return "D"
	case gamepkg.Hearts:
		return "H"
	case gamepkg.Spades:
		return "S"
	case 4:
		return "NT"
	default:
		return fmt.Sprintf("%d", strain)
	}
}

// parseBid parses human text into a game bid
func parseBid(input string) (gamepkg.Bid, error) {
	in := strings.TrimSpace(strings.ToLower(input))
	switch in {
	case "pass", "p":
		return gamepkg.NewPass(), nil
	case "double", "dbl", "x":
		return gamepkg.NewDouble(), nil
	case "redouble", "rdbl", "xx":
		return gamepkg.NewRedouble(), nil
	}
	if len(in) < 2 {
		return gamepkg.Bid{}, fmt.Errorf("invalid bid format")
	}
	lvl := int(in[0] - '0')
	if lvl < 1 || lvl > 7 {
		return gamepkg.Bid{}, fmt.Errorf("bid level must be between 1 and 7")
	}
	s := strings.ToUpper(in[1:])
	var suit gamepkg.Suit
	switch s {
	case "C":
		suit = gamepkg.Clubs
	case "D":
		suit = gamepkg.Diamonds
	case "H":
		suit = gamepkg.Hearts
	case "S":
		suit = gamepkg.Spades
	case "NT", "N":
		suit = 4
	default:
		return gamepkg.Bid{}, fmt.Errorf("invalid suit: %s", s)
	}
	return gamepkg.NewBid(lvl, suit), nil
}

func parsePosition(pos string) (gamepkg.Position, error) {
	s := strings.ToLower(strings.TrimSpace(pos))
	switch s {
	case "north":
		return gamepkg.North, nil
	case "east":
		return gamepkg.East, nil
	case "south":
		return gamepkg.South, nil
	case "west":
		return gamepkg.West, nil
	default:
		return 0, fmt.Errorf("invalid position: %s", pos)
	}
}

// session store helpers
func (s *Server) sessPut(sess *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
}
func (s *Server) sessGet(id string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[id]
	return sess, ok
}

// handleEvaluateBid evaluates a bid and provides feedback
func (s *Server) handleEvaluateBid(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
		Position  string `json:"position"`
		Bid       string `json:"bid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Get the session
	s.mu.RLock()
	sess, ok := s.sessions[req.SessionID]
	s.mu.RUnlock()

	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Parse position
	pos, err := parsePosition(req.Position)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the bid
	bid, err := parseBid(req.Bid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the player
	var player *gamepkg.Player
	for _, p := range sess.Players {
		if p.Position == pos {
			player = p
			break
		}
	}

	if player == nil {
		http.Error(w, "player not found", http.StatusNotFound)
		return
	}

	// Get the AI's recommended bid
	recommendedBid := player.MakeBid(sess.Auction)

	// Check if the bid is the same as recommended
	isRecommended := false
	if bid.Pass && recommendedBid.Pass {
		isRecommended = true
	} else if bid.Double && recommendedBid.Double {
		isRecommended = true
	} else if bid.Redouble && recommendedBid.Redouble {
		isRecommended = true
	} else if bid.Level == recommendedBid.Level && bid.Strain == recommendedBid.Strain {
		isRecommended = true
	}

	// Prepare response
	response := map[string]interface{}{
		"isRecommended": isRecommended,
		"recommendedBid":  recommendedBid.String(),
	}

	// Add explanation if bid is not recommended
	if !isRecommended {
		hcp, _ := player.Hand.Evaluate()
		explanation := fmt.Sprintf("With %d HCP, the recommended bid is %s", hcp, recommendedBid.String())
		response["explanation"] = explanation
	}

	writeJSON(w, http.StatusOK, response)
}

// writeJSON is a helper to encode responses
func writeJSON(w http.ResponseWriter, status int, payload any) {
	setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// setCORSHeaders sets permissive CORS headers for browser clients
func setCORSHeaders(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
