// 最早完成时间等于最晚完成时间=t ,且所
// 有其他节点要么最晚完成时间小于t ,要么
// 最早完成时间大于t

// https://www.youtube.com/watch?v=q22FRC3fc4w
// Fear has fallen down on me, I’m buried in too deep
// I stumble to my feet, tears frozen to my cheek
// Dreams that bring me light, they tell me to survive
// I’m reaching for a life beyond the peak

#ifndef ONLINE_JUDGE
#include "templates/debug.hpp"
#else
#define debug(...)
#endif

#include <bits/stdc++.h>
using namespace std;
using i64 = int64_t;
using u64 = uint64_t;

// #define int i64
mt19937_64 rng(std::random_device{}());
void solve() {
    int n = 10, m = 40;
    // cin >> n >> m;
    vector<vector<int>> adj(n);
    set<pair<int, int>> es;
    for (int i = 0; i < m; i++) {
        int u = rng() % n, v = rng() % n;
        if (u > v) swap(u, v);
        while (u == v || es.count({u, v})) {
            u = rng() % n, v = rng() % n;
            if (u > v) swap(u, v);
        }
        // int u, v;
        // cin >> u >> v;
        adj[u].push_back(v);
        es.insert({u, v});
    }
    vector<vector<int>> link(n, vector<int>(n, 0));
    for (int i = 0; i < n; i++) {
        link[i][i] = 1;
        for (int j : adj[i]) {
            link[i][j] = 1;
        }
    }
    for (int k = 0; k < n; k++) {
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                if (link[i][k] && link[k][j]) {
                    link[i][j] = 1;
                }
            }
        }
    }
    vector<int> std_ans(n, 0);
    for (int i = 0; i < n; i++) {
        int cnt = 0;
        for (int j = 0; j < n; j++) {
            if (link[i][j] || link[j][i]) {
                cnt++;
            }
        }
        if (cnt == n) std_ans[i] = 1;
    }
    // get min and max time for each node
    vector<int> min_time(n, 0), max_time(n, 0);
    int tot = 0;
    for (int i = 0; i < n; i++) {
        for (int j : adj[i]) {
            min_time[j] = max(min_time[j], min_time[i] + 1);
            tot = max(tot, min_time[j]);
        }
    }
    max_time = vector<int>(n, tot);
    for (int i = n - 1; i >= 0; i--) {
        for (int j : adj[i]) {
            max_time[i] = min(max_time[i], max_time[j] - 1);
        }
    }
    vector<int> my_ans(n, 1);
    for (int i = 0; i < n; i++) {
        bool feas = true;
        if (min_time[i] != max_time[i]) feas = false;
        int t = min_time[i];
        for (int j = 0; j < n; j++) {
            if (j != i && (max_time[j] >= t && min_time[j] <= t)) feas = false;
        }
        my_ans[i] = feas;
    }
    for (int i = 0; i < n; i++) {
        if (std_ans[i] != my_ans[i]) {
            cout << "NO\n";
            return;
        }
    }
    // debug(my_ans, std_ans);
    cout << "YES\n";
}
#undef int

// 大胆猜想，小心求证 Make bold hypotheses and verify carefully
// - You REALLY need some key observations...
// - Don't trust seemaxgly trival conclusions
// - Do something instead of nothing and stay organized
// - Don't get stuck on one approach
// - Formalization is the death of intuition

int main() {
    cin.tie(nullptr);
    ios::sync_with_stdio(false);
    int t = 1;
    cin >> t;
    while (t--) {
        solve();
    }
}