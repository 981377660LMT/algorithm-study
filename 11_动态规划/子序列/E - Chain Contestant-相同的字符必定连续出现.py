# 给出一个只包括A~J的字符串，定义一种子序列为：
# !在这个子序列中，相同的字符必定连续出现，求出这样的子序列有多少个。
# n<=1000

# !子序列 每个位置选还是不选 dfs(index,visited,pre)

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# print(ord("A")) # 65

if __name__ == "__main__":

    # !dfs TLE
    # @lru_cache(None)
    # def dfs(index: int, visited: int, pre: int) -> int:
    #     if index == n:
    #         return 1

    #     # jump
    #     res = dfs(index + 1, visited, pre)

    #     # not jump
    #     cur = ord(s[index]) - 65
    #     if pre == cur or not visited & (1 << cur):
    #         res += dfs(index + 1, visited | (1 << cur), cur)
    #     return res % MOD

    # n = int(input())
    # s = input()
    # print(dfs(0, 0, 0) - 1)

    # !dp
    n = int(input())
    s = input()
    dp = [[0] * 10 for _ in range(1 << 10)]  # dp[visited][pre]
    dp[0][0] = 1
    for i in range(n):
        ndp = [[0] * 10 for _ in range(1 << 10)]
        cur = ord(s[i]) - 65
        for visited in range(1 << 10):
            for pre in range(10):
                # jump
                ndp[visited][pre] += dp[visited][pre]
                # not jump
                if (pre == cur) or (not visited & (1 << cur)):
                    ndp[visited | (1 << cur)][cur] += dp[visited][pre]
                    ndp[visited | (1 << cur)][cur] %= MOD
        dp = ndp

    res = 0
    for nums in dp:
        res = (res + sum(nums)) % MOD
    print(res - 1)
