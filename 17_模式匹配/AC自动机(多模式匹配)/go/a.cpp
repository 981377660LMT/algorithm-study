#pragma once

template <size_t X = 26, char margin = 'a'>
struct Trie {
  struct Node {
    array<int, X> nxt;
    vector<int> idxs;
    int idx;
    char key;
    Node(char c) : idx(-1), key(c) { fill(nxt.begin(), nxt.end(), -1); }
  };

  vector<Node> st;

  Trie(char c = '$') { st.emplace_back(c); }

  inline int &next(int i, int j) { return st[i].nxt[j]; }

  void add(const string &s, int x) {
    int pos = 0;
    for (int i = 0; i < (int)s.size(); i++) {
      int k = s[i] - margin;
      if (~next(pos, k)) {
        pos = next(pos, k);
        continue;
      }
      int npos = st.size();
      next(pos, k) = npos;
      st.emplace_back(s[i]);
      pos = npos;
    }
    st[pos].idx = x;
    st[pos].idxs.emplace_back(x);
  }

  int find(const string &s) {
    int pos = 0;
    for (int i = 0; i < (int)s.size(); i++) {
      int k = s[i] - margin;
      if (next(pos, k) < 0) return -1;
      pos = next(pos, k);
    }
    return pos;
  }

  int move(int pos, char c) {
    assert(pos < (int)st.size());
    return pos < 0 ? -1 : next(pos, c - margin);
  }

  int size() const { return st.size(); }

  int idx(int pos) { return pos < 0 ? -1 : st[pos].idx; }

  vector<int> idxs(int pos) { return pos < 0 ? vector<int>() : st[pos].idxs; }
};