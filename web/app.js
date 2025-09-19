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
  el('sessionId').textContent = state.id;
  el('dealer').textContent = state.dealer;
  el('complete').textContent = String(state.complete);

  const players = state.players.map(p => {
    return `<div style="margin-bottom:8px">
      <div><b>${p.position}</b> — HCP: ${p.hcp}</div>
      <div>♠ ${p.spades}</div>
      <div>♥ ${p.hearts}</div>
      <div>♦ ${p.diamonds}</div>
      <div>♣ ${p.clubs}</div>
    </div>`;
  }).join('');
  el('players').innerHTML = players;

  const rows = state.auction.map(a => `<tr><td>${a.position}</td><td>${
    a.pass ? 'Pass' : (a.redouble ? 'XX' : (a.double ? 'X' : `${a.level}${a.strain}`))
  }</td></tr>`).join('');
  el('auction').innerHTML = rows || '<tr><td colspan="2">No bids yet</td></tr>';
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
    const bid = el('bid').value.trim();
    if (!bid) {
      el('message').textContent = 'Enter a bid';
      return;
    }
    try {
      const state = await API.postBid(sessionId, position, bid);
      lastState = state;
      render(state);
      el('bid').value = '';
      el('message').textContent = 'Bid accepted';
    } catch (e) {
      el('message').textContent = e.message;
    }
  });
}

window.addEventListener('DOMContentLoaded', main);
