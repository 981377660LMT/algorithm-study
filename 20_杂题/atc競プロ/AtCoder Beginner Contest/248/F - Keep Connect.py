"""
2<=n<=3000 暗示dp
9e8<=p<=1e9 素数

n列 2*n个顶点 3*n-2条边的类似铁轨的图
删去i条边(i=1,2,...n-1) 剩下的图还是联通的 对每个i求方案数模p 

dp维度:当前的列数i、删去的边数j、第i列上下两个点是否连通
"""
from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def main() -> None:
    # n, MOD = map(int, input().split())

    # !down-top 的形式 TLE
    # @lru_cache(None)
    # def dfs(index: int, remain: int, isPreConnected: bool) -> int:
    #     """当前的列数index、还可以删去的边数remain、前一列上下两个点是否连通isPreConnected"""
    #     if remain < 0:
    #         return 0
    #     if index == n:
    #         return 1 if (remain == 0 and isPreConnected) else 0

    #     res = 0
    #     if not isPreConnected:  # 此时横不能删
    #         res += dfs(index + 1, remain, True)
    #         res += dfs(index + 1, remain - 1, False) if remain >= 1 else 0
    #     else:
    #         # 当前不删竖
    #         res += dfs(index + 1, remain, True)  # 不删
    #         res += 2 * dfs(index + 1, remain - 1, True) if remain >= 1 else 0  # 删任一个横

    #         # 当前删竖
    #         res += dfs(index + 1, remain - 1, True) if remain >= 1 else 0  # 只删1个竖
    #         res += (
    #             2 * dfs(index + 1, remain - 2, False) if remain >= 2 else 0
    #         )  # 删1个竖加任一个横

    #     return res % MOD

    # for i in range(1, n):
    #     print((dfs(1, i, True) + dfs(1, i - 1, False)) % MOD)
    # dfs.cache_clear()

    # !改成 top-down 的形式
    n, MOD = map(int, input().split())
    dp = [[0] * 2 for _ in range(n + 1)]  # 删的次数 ; 前一列两个点[不联通, 联通]
    dp[0][1], dp[1][0] = 1, 1
    for _ in range(1, n):
        ndp = [[0] * 2 for _ in range(n + 1)]
        for i in range(n):
            ndp[i][1] = (ndp[i][1] + dp[i][0]) % MOD
            if i >= 1:
                ndp[i][0] = (ndp[i][0] + dp[i - 1][0]) % MOD

            ndp[i][1] = (ndp[i][1] + dp[i][1]) % MOD
            if i >= 1:
                ndp[i][1] = (ndp[i][1] + 2 * dp[i - 1][1]) % MOD

            if i >= 1:
                ndp[i][1] = (ndp[i][1] + dp[i - 1][1]) % MOD
            if i >= 2:
                ndp[i][0] = (ndp[i][0] + 2 * dp[i - 2][1]) % MOD

        dp = ndp

    for i in range(1, n):
        print(dp[i][1])


# 7 15
if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
