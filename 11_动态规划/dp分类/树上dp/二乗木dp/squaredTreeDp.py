# # https://ferin-15.github.io/program_contest_library/library/DP/squared_tree_dp.cpp.html
# // 2乗の木DPの雛形
# vector<ll> sz(n);
# vector<vector<mint>> dp(n);
# auto dfs = [&](auto &&self, ll v, ll p) -> void {
#     sz[v] = 1;
#     dp[v].resize(sz[v]+1);
#     dp[v][1] = 1;
#     for(auto to: g[v]) if(to!=p) {
#         self(self, to, v);
#         vector<mint> merged(2*(sz[v]+sz[to])+1);
#         REP(i, sz[v]*2+1) REP(j, sz[to]*2+1) {
#             merged[i+j] += dp[v][i] * dp[to][j];
#         }
#         sz[v] += sz[to];
#         dp[v] = move(merged);
#     }
#     dp[v][0] = 1;
# };


from typing import List, Tuple

MOD = int(1e9 + 7)


def solve(n: int, edges: List[Tuple[int, int]]) -> int:
    tree = [[] for _ in range(n)]
    for u, v in edges:
        tree[u].append(v)
        tree[v].append(u)
    subSize = [0] * n
    dp = [[] for _ in range(n)]

    def dfs(cur: int, pre: int) -> None:
        subSize[cur] = 1
        dp[cur] = [0] * (subSize[cur] + 1)
        dp[cur][1] = 1
        for next in tree[cur]:
            if next != pre:
                dfs(next, cur)
                merged = [0] * (2 * (subSize[cur] + subSize[next]) + 1)
                for i in range(2 * subSize[cur] + 1):
                    for j in range(2 * subSize[next] + 1):
                        merged[i + j] += dp[cur][i] * dp[next][j]
                        merged[i + j] %= MOD
                subSize[cur] += subSize[next]
                dp[cur] = merged
        dp[cur][0] = 1

    dfs(0, -1)
