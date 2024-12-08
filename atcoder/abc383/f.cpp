#pragma GCC target("arch=skylake-avx512")
#pragma GCC target("avx2")
#pragma GCC optimize("O3")
#pragma GCC target("sse4")
#pragma GCC optimize("unroll-loops")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

#include <bits/stdc++.h>
using namespace std;

static const long long INF = -1000000000000000LL;

int main() {
    int N, X;
    long long K;
    cin >> N >> X >> K;

    struct Item {
        int p;
        long long u;
        int c;
    };
    vector<Item> items(N);
    int maxColor = 0;
    for (int i = 0; i < N; i++) {
        cin >> items[i].p >> items[i].u >> items[i].c;
        maxColor = max(maxColor, items[i].c);
    }

    vector<vector<pair<int,long long>>> color_groups(maxColor+1);
    for (auto &it : items) {
        color_groups[it.c].push_back({it.p, it.u});
    }

    vector<vector<long long>> all_color_dp;
    all_color_dp.reserve(maxColor);
    for (int c = 1; c <= maxColor; c++) {
        vector<long long> dp_color(X+1, INF);
        dp_color[0] = 0;
        for (auto &it : color_groups[c]) {
            int p = it.first;
            long long u = it.second;
            for (int w = X - p; w >= 0; w--) {
                if (dp_color[w] != INF) {
                    long long val = dp_color[w] + u;
                    if (val > dp_color[w+p]) {
                        dp_color[w+p] = val;
                    }
                }
            }
        }
        all_color_dp.push_back(move(dp_color));
    }

    auto build_frontier = [&](const vector<long long> &dp_color) {
        vector<pair<int,long long>> states;
        states.reserve(X+1);
        for (int w = 0; w <= X; w++) {
            if (dp_color[w] != INF) {
                states.push_back({w, dp_color[w]});
            }
        }
        if (states.empty()) {
            states.push_back({0,0});
        }

        sort(states.begin(), states.end(), [&](auto &a, auto &b){
            if (a.first == b.first) return a.second > b.second;
            return a.first < b.first;
        });

        {
            vector<pair<int,long long>> unique_states;
            unique_states.reserve(states.size());
            int prev_cost = -1;
            long long best_u = -1;
            for (auto &st : states) {
                if (st.first != prev_cost) {
                    unique_states.push_back(st);
                    prev_cost = st.first;
                    best_u = st.second;
                } else {
                    if (st.second > best_u) {
                        unique_states.back() = st;
                        best_u = st.second;
                    }
                }
            }
            states = move(unique_states);
        }

        vector<pair<int,long long>> frontier;
        frontier.reserve(states.size());
        long long best_util = -1;
        for (auto &st : states) {
            if (st.second > best_util) {
                frontier.push_back(st);
                best_util = st.second;
            }
        }
        return frontier;
    };

    vector<vector<pair<int,long long>>> frontiers;
    frontiers.reserve(all_color_dp.size());
    for (auto &dp_color : all_color_dp) {
        frontiers.push_back(build_frontier(dp_color));
    }

    vector<long long> dp(X+1, INF);
    dp[0] = 0;

    for (auto &frontier : frontiers) {
        vector<long long> new_dp(X+1, INF);
        for (int w = 0; w <= X; w++) {
            if (dp[w] != INF && dp[w] > new_dp[w]) {
                new_dp[w] = dp[w];
            }
        }

        for (auto &st : frontier) {
            int cst = st.first;
            long long util = st.second;
            long long color_bonus = (cst > 0) ? K : 0;
            for (int w = 0; w + cst <= X; w++) {
                if (dp[w] != INF) {
                    long long val = dp[w] + util + color_bonus;
                    if (val > new_dp[w+cst]) {
                        new_dp[w+cst] = val;
                    }
                }
            }
        }

        dp = move(new_dp);
    }

    long long ans = *max_element(dp.begin(), dp.end());
    cout << ans << "\n";

    return 0;
}
