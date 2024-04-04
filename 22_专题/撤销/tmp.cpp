#include <bits/stdc++.h>
#define finish(x) return cout << x << endl, 0
typedef long long ll;
typedef long double ldb;
const int md = 1e9 + 7;
const ll inf = 2e18;
const int OO = 1;
const int OOO = 1;
using namespace std;

// used in the original data structure to represent a change in memory cell
struct change {
	int t, value;
	int *cell;
	change() {}
	change(int _time, int _value, int *_cell) {
		t = _time;
		value = _value;
		cell = _cell;
	}
	void undo() {
		*cell = value;
	}
};

// dsu with normal undo, allowes to query bipartite-ness.
struct dsu {
	int n;
	vector<int> p, r;
	vector<int> up;
	int nxtupd;
	vector<change> stack;
	int first_violate;
	dsu() {}
	dsu(int sz) {
		n = sz;
		p.resize(n);
		r.resize(n, 0);
		up.resize(n, 0);
		nxtupd = 0;
		first_violate = -1;
		for (int i = 0; i < n; i++) p[i] = i;
	}
	void upd(int *cell, int value) {
		stack.push_back(change(nxtupd, *cell, cell));
		*cell = value;
	}
	pair<int, int> find(int x) {
		if (x == p[x]) return{ x, 0 };
		pair<int, int> tmp = find(p[x]);
		tmp.second ^= up[x];
		return tmp;
	}
	bool is_bipartite() {
		return first_violate == -1; // no violation
	}
	void unite(int x, int y) {
		nxtupd++;
		int col = 1;
		pair<int, int> tmp = find(x);
		x = tmp.first, col ^= tmp.second;
		tmp = find(y);
		y = tmp.first, col ^= tmp.second;
		if (x == y) {
			if (col == 1 && first_violate == -1) first_violate = nxtupd;
			upd(&p[0], p[0]); // mark this update's existence.
			return;
		}
		if (r[x] < r[y]) {
			upd(&p[x], y);
			upd(&up[x], col);
		}
		else {
			upd(&p[y], x);
			upd(&up[y], col);
			if (r[x] == r[y]) {
				upd(&r[x], r[x] + 1);
			}
		}
		return;
	}
	void undo() {
		if (!stack.size()) return;
		int t = stack.back().t;
		if (first_violate == t) first_violate = -1;
		while (stack.size() && stack.back().t == t) {
			stack.back().undo();
			stack.pop_back();
		}
	}
};

// I'd want this 'struct update' to be more generic, and support maintaining different kinds of updates...
// but for now I'm not sure how. So this is a specific impl for this problem ('update' only relates to 'unite').
struct update {
	char type; // 'A' or 'B'
	int x, y;
	update() {}
	update(int xx, int yy) {
		x = xx;
		y = yy;
		type = 'B';
	}
};

struct dsuqueue {
	dsu D;
	vector<update> S;
	int bottom; // bottom of the stack: S[0..bottom-1] is entirely B's.
	dsuqueue() {}
	dsuqueue(int sz) {
		D = dsu(sz);
		bottom = 0;
	}
	// utility:
	void advance_bottom() {
		while (bottom < S.size() && S[bottom].type == 'B') bottom++;
	}
	void fix() {
		if (!S.size() || S.back().type == 'A') return;
		advance_bottom();
		vector<update> saveB, saveA;
		saveB.push_back(S.back());
		S.pop_back(), D.undo();
		while (saveA.size() != saveB.size() && S.size() > bottom) {
			if (S.back().type == 'A')
				saveA.push_back(S.back());
			else
				saveB.push_back(S.back());
			S.pop_back(), D.undo();
		}
		// reverse saveA and saveB so their relative order is maintained
		reverse(saveA.begin(), saveA.end());
		reverse(saveB.begin(), saveB.end());
		for (const update &u : saveB) {
			S.push_back(u);
			D.unite(u.x, u.y);
		}
		for (const update &u : saveA) {
			S.push_back(u);
			D.unite(u.x, u.y);
		}
		advance_bottom();
	}
	void reverse_updates() {
		for (int i = 0; i < S.size(); i++)
			D.undo();
		for (int i = (int)S.size() - 1; i >= 0; i--) {
			D.unite(S[i].x, S[i].y);
			S[i].type = 'A';
		}
		reverse(S.begin(), S.end());
		bottom = 0;
	}
	void undo() {
		advance_bottom();
		if (bottom == S.size()) {
			// no more A's, let's reverse and begin again.
			reverse_updates();
		}
		fix();
		D.undo();
		S.pop_back();
	}

	// begin copying all functions from the original dsu.
	bool is_bipartite() {
		return D.is_bipartite();
	}
	void unite(int x, int y) {
		D.unite(x, y);
		S.push_back(update(x, y));
	}
};

int n, m, q;
vector<pair<int, int>> e;
vector<int> R; // for each i, R[i] would be the minimum such that query (i, R[i]) is answered with NO (the graph is bipartite)

int main() {
	ios::sync_with_stdio(0), cin.tie(0);
	cin >> n >> m >> q;
	dsuqueue D(n);
	e.resize(m);
	for (auto &i : e) {
		cin >> i.first >> i.second;
		--i.first, --i.second;
	}

	R.resize(m);
	for (int i = 0; i < m; i++)
		D.unite(e[i].first, e[i].second);
	int nxt = 0;
	for (int i = 0; i < m; i++) {
		while (!D.is_bipartite() && nxt < m) {
			D.undo();
			nxt++;
		}
		if (D.is_bipartite())
			R[i] = nxt - 1;
		else
			R[i] = md;

		D.unite(e[i].first, e[i].second);
	}

	while (q--) {
		int l, r;
		cin >> l >> r;
		--l, --r;
		if (R[l] <= r) cout << "NO\n";
		else cout << "YES\n";
	}
}