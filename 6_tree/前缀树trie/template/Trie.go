/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 前缀树是一个树状的数据结构，
// 用于高效地存储和检索一系列字符串的前缀。
// 前缀树有许多应用，如自动补全和拼写检查

// https://maspypy.github.io/library/string/trie.hpp

package main

func main() {

}

// #include "alg/monoid/add.hpp"

// template <int sigma>
// struct Trie {
//   using ARR = array<int, sigma>;
//   int n_node;
//   vc<ARR> TO;
//   vc<int> parent;
//   vc<int> suffix_link;
//   vc<int> words;
//   vc<int> BFS; // BFS 順

//   Trie() {
//     n_node = 0;
//     new_node();
//   }

//   template <typename STRING>
//   int add(STRING S, int off) {
//     int v = 0;
//     for (auto&& ss: S) {
//       int s = ss - off;
//       assert(0 <= s && s < sigma);
//       if (TO[v][s] == -1) {
//         TO[v][s] = new_node();
//         parent.back() = v;
//       }
//       v = TO[v][s];
//     }
//     words.eb(v);
//     return v;
//   }

//   int add_char(int v, int c, int off) {
//     c -= off;
//     if (TO[v][c] != -1) return TO[v][c];
//     TO[v][c] = new_node();
//     parent.back() = v;
//     return TO[v][c];
//   }

//   void calc_suffix_link(bool upd_TO) {
//     suffix_link.assign(n_node, -1);
//     BFS.resize(n_node);
//     int p = 0, q = 0;
//     BFS[q++] = 0;
//     while (p < q) {
//       int v = BFS[p++];
//       FOR(s, sigma) {
//         int w = TO[v][s];
//         if (w == -1) continue;
//         BFS[q++] = w;
//         int f = suffix_link[v];
//         while (f != -1 && TO[f][s] == -1) f = suffix_link[f];
//         suffix_link[w] = (f == -1 ? 0 : TO[f][s]);
//       }
//     }
//     if (!upd_TO) return;
//     for (auto&& v: BFS) {
//       FOR(s, sigma) if (TO[v][s] == -1) {
//         int f = suffix_link[v];
//         TO[v][s] = (f == -1 ? 0 : TO[f][s]);
//       }
//     }
//   }

//   vc<int> calc_count() {
//     assert(!suffix_link.empty());
//     vc<int> count(n_node);
//     for (auto&& x: words) count[x]++;
//     for (auto&& v: BFS)
//       if (v) { count[v] += count[suffix_link[v]]; }
//     return count;
//   }

// private:
//   int new_node() {
//     parent.eb(-1);
//     TO.eb(ARR{});
//     fill(all(TO.back()), -1);
//     return n_node++;
//   }
// };
