# 求从S点到T点 经过X点偶数次 走的路径长度为k 的方案数。
# !dp 很明显的三个维度 记忆化搜索
# !但是记忆化搜索TLE了 必须要dp数组写

# n,m,k<=2000

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    # @lru_cache(None)
    # def dfs(cur: int, remain: int, mod: int) -> int:
    #     if remain == 0:
    #         return int(cur == end and mod == 0)
    #     res = 0
    #     for next in adjMap[cur]:
    #         res += dfs(next, remain - 1, (mod ^ ( next==x ))
    #         res %= MOD
    #     return res

    # n, m, k, start, end, x = map(int, input().split())
    # adjMap = defaultdict(set)
    # for _ in range(m):
    #     u, v = map(int, input().split())
    #     adjMap[u].add(v)
    #     adjMap[v].add(u)
    # res = dfs(start, k, 0)
    # dfs.cache_clear()
    # print(res)

    n, m, k, start, end, x = map(int, input().split())
    adjList = [[] for _ in range(n + 1)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList[u].append(v)
        adjList[v].append(u)

    # !配るDP 开三维数组TLE了 换成滚动数组优化才过 (局部性原理)
    dp = [[0, 0] for _ in range(n + 1)]
    dp[start][0] = 1
    for _ in range(k):
        ndp = [[0, 0] for _ in range(n + 1)]
        for pre in range(1, n + 1):
            for next in adjList[pre]:
                # !与 1 异或 表示 flip TLE
                ndp[next][0 ^ (next == x)] += dp[pre][0]
                ndp[next][1 ^ (next == x)] += dp[pre][1]
                ndp[next][0 ^ (next == x)] %= MOD
                ndp[next][1 ^ (next == x)] %= MOD
        dp = ndp
    print(dp[end][0])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
