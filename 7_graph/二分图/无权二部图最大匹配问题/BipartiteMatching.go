// https://ei1333.github.io/library/graph/flow/bipartite-matching.hpp

// package main
// /**
//  * @brief Bipartite-Matching(二部グラフの最大マッチング)
//  * @docs docs/bipartite-matching.md
//  */
//  struct BipartiteMatching {
//   vector< vector< int > > graph;
//   vector< int > alive, used, match;
//   int timestamp;

//   explicit BipartiteMatching(int n) : graph(n), alive(n, 1), used(n, 0), match(n, -1), timestamp(0) {}

//   void add_edge(int u, int v) {
//     graph[u].push_back(v);
//     graph[v].push_back(u);
//   }

//   bool augment(int idx) {
//     used[idx] = timestamp;
//     for(auto &to : graph[idx]) {
//       int to_match = match[to];
//       if(alive[to] == 0) continue;
//       if(to_match == -1 || (used[to_match] != timestamp && augment(to_match))) {
//         match[idx] = to;
//         match[to] = idx;
//         return true;
//       }
//     }
//     return false;
//   }

//   int bipartite_matching() {
//     int ret = 0;
//     for(int i = 0; i < (int) graph.size(); i++) {
//       if(alive[i] == 0) continue;
//       if(match[i] == -1) {
//         ++timestamp;
//         ret += augment(i);
//       }
//     }
//     return ret;
//   }

//   int add_vertex(int idx) {
//     alive[idx] = 1;
//     ++timestamp;
//     return augment(idx);
//   }

//   int erase_vertex(int idx) {
//     alive[idx] = 0;
//     if(match[idx] == -1) {
//       return 0;
//     }
//     match[match[idx]] = -1;
//     ++timestamp;
//     int ret = augment(match[idx]);
//     match[idx] = -1;
//     return ret - 1;
//   }

//   void output() const {
//     for(int i = 0; i < (int) graph.size(); i++) {
//       if(i < match[i]) {
//         cout << i << "-" << match[i] << endl;
//       }
//     }
//   }
// };