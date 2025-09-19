const API = {
  createSession: async () => {
    const res = await fetch('/api/sessions', { method: 'POST' });
    if (!res.ok) throw new Error('Failed to create session');
    return res.json();
  },
  getSession: async (id) => {
    const res = await fetch(`/api/sessions/${id}`);
    if (!res.ok) throw new Error('Failed to fetch session');
    return res.json();
  },
  postBid: async (id, position, bid) => {
    const res = await fetch(`/api/sessions/${id}/bid`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ position, bid })
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
  }
};

const el = (id) => document.getElementById(id);

function render(state) {
  if (!state) return;
  el('sessionId').textContent = state.id;
  el('dealer').textContent = state.dealer;
  el('complete').textContent = String(state.complete);
  // Show current turn (same as dealer)
  const current = state.dealer;
  const turnEl = el('currentTurn');
  if (turnEl) turnEl.textContent = current;

  // Auto-select current dealer as position to reduce mistakes
  if (el('position')) {
    el('position').value = current;
  }

  // Update compact bid history
  const lastBids = { N: '-', E: '-', S: '-', W: '-' };
  state.auction.forEach(bid => {
    lastBids[bid.position.charAt(0)] = bid.pass ? 'Pass' : 
      (bid.redouble ? 'XX' : (bid.double ? 'X' : `${bid.level}${bid.strain}`));
  });
  
  // Update the bid history display
  Object.entries(lastBids).forEach(([pos, bid]) => {
    const el = document.getElementById(`bid-${pos}`);
    if (el) el.textContent = bid;
  });

  const players = state.players.map(p => {
    const isDealer = p.position === state.dealer;
    const cls = `player${isDealer ? ' dealer' : ''}`;
    const badge = isDealer ? '<span class="badge">Current</span>' : '';
    return `<div class="${cls}" style="margin-bottom:8px">
      <div><b>${p.position}</b> — HCP: ${p.hcp} ${badge}</div>
      <div><span class="suit-spades">♠</span> ${p.spades}</div>
      <div><span class="suit-hearts">♥</span> ${p.hearts}</div>
      <div><span class="suit-diamonds">♦</span> ${p.diamonds}</div>
      <div><span class="suit-clubs">♣</span> ${p.clubs}</div>
    </div>`;
  }).join('');
  el('players').innerHTML = players;

  const tbody = el('auction').querySelector('tbody');
  tbody.innerHTML = state.auction.length > 0 
    ? state.auction.map(a => `<tr><td>${a.position}</td><td>${
        a.pass ? 'Pass' : (a.redouble ? 'XX' : (a.double ? 'X' : `${a.level}${a.strain}`))
      }</td></tr>`).join('')
    : '<tr><td colspan="2">No bids yet</td></tr>';

  // Update bid button enabled/disabled state
  updateBidAvailability(state);
}

async function main() {
  let sessionId = null;
  let lastState = null;

  async function refresh() {
    if (!sessionId) return;
    try {
      lastState = await API.getSession(sessionId);
      render(lastState);
      el('message').textContent = '';
    } catch (e) {
      el('message').textContent = e.message;
    }
  }

  el('newSessionBtn').addEventListener('click', async () => {
    try {
      const state = await API.createSession();
      sessionId = state.id;
      lastState = state;
      render(state);
      el('message').textContent = 'New session created';
    } catch (e) {
      el('message').textContent = e.message;
    }
  });

  el('refreshBtn').addEventListener('click', refresh);

  el('sendBidBtn').addEventListener('click', async () => {
    if (!sessionId) {
      el('message').textContent = 'Create a session first';
      return;
    }
    const position = el('position').value;
    let bid = el('bid').value.trim();
    if (!bid) {
      el('message').textContent = 'Enter a bid';
      return;
    }
    if (!isValidBidFormat(bid)) {
      el('message').textContent = 'Invalid bid format. Use 1-7 + C/D/H/S/NT (e.g., 1C, 2NT) or Pass, X, XX';
      return;
    }
    // Normalize: uppercase and map single trailing N to NT
    bid = normalizeBid(bid);
    try {
      const state = await API.postBid(sessionId, position, bid);
      lastState = state;
      render(state);
      // Auto-advance position to next dealer (server returns updated dealer)
      if (el('position')) {
        el('position').value = state.dealer;
      }
      // Clear and focus bid input for next entry
      el('bid').value = '';
      el('bid').focus();
      el('message').textContent = 'Bid accepted';
      updateBidAvailability(state);
    } catch (e) {
      el('message').textContent = e.message;
    }
  });

  // Recompute availability and update input validation when user types
  el('position').addEventListener('change', () => updateBidAvailability(lastState));
  
  const bidInput = el('bid');
  bidInput.addEventListener('input', () => {
    const bid = bidInput.value.trim();
    updateBidAvailability(lastState);
    
    // Update input validation classes in real-time
    if (bid.length === 0) {
      bidInput.classList.remove('valid', 'invalid');
    } else {
      const isValid = isValidBidFormat(bid);
      bidInput.classList.toggle('valid', isValid);
      bidInput.classList.toggle('invalid', !isValid);
    }
  });
  
  // Clear validation state when input loses focus if empty
  bidInput.addEventListener('blur', () => {
    if (bidInput.value.trim() === '') {
      bidInput.classList.remove('valid', 'invalid');
    }
  });
}

window.addEventListener('DOMContentLoaded', main);

function updateBidAvailability(state) {
  const btn = el('sendBidBtn');
  if (!btn) return;
  if (!state) { btn.disabled = true; return; }
  const pos = el('position').value;
  const dealer = state.dealer;
  const complete = !!state.complete;
  const bidValue = (el('bid').value || '').trim();
  const validFormat = bidValue.length > 0 && isValidBidFormat(bidValue);
  const can = !complete && pos === dealer && validFormat;
  btn.disabled = !can;
  const msg = el('message');
  if (complete) {
    msg.textContent = 'Auction is complete';
  } else if (pos !== dealer) {
    msg.textContent = `It is ${dealer}'s turn to bid`;
  } else if (bidValue && !validFormat) {
    msg.textContent = 'Invalid bid format. Use 1-7 + C/D/H/S/NT (e.g., 1C, 2NT) or Pass, X, XX';
  } else {
    // keep current message if any
  }
}

// Client-side validation: accepts contract bids (1-7 + C/D/H/S/NT/N) and special tokens Pass/P, Double/DBL/X, Redouble/RDBL/XX
function isValidBidFormat(input) {
  const s = (input || '').trim().toLowerCase();
  if (!s) return false;
  // special tokens
  if (s === 'pass' || s === 'p' || s === 'double' || s === 'dbl' || s === 'x' || s === 'redouble' || s === 'rdbl' || s === 'xx') {
    return true;
  }
  // contract like 1c, 2nt, 3h, 4s
  const m = s.match(/^([1-7])(c|d|h|s|nt|n)$/i);
  return !!m;
}

function normalizeBid(input) {
  let s = (input || '').trim().toUpperCase();
  if (!s) return s;
  // Special tokens normalization
  if (s === 'P') return 'Pass';
  if (s === 'X' || s === 'DBL') return 'X';
  if (s === 'XX' || s === 'RDBL') return 'XX';
  if (s === 'PASS' || s === 'DOUBLE' || s === 'REDOUBLE') return s.charAt(0) + s.slice(1).toLowerCase();
  // Map 1N -> 1NT, etc.
  const m = s.match(/^([1-7])(C|D|H|S|N|NT)$/);
  if (m) {
    const lvl = m[1];
    let strain = m[2];
    if (strain === 'N') strain = 'NT';
    return `${lvl}${strain}`;
  }
  return s;
}
