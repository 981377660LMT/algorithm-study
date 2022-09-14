# 开始有n个盘子 每个盘子里有ai个寿司
# 每次随机一个盘子编号 吃掉那个盘子的一个寿司
# 求吃完所有盘子里寿司需要的步数的期望值
# n<=300
# !1<=ai<=3

# !状态如何定义:
# 利用寿司盘子无区别的特性 不需要存每个盘子剩下多少个 而是存剩下k个寿司的盘子有多少个
# (因为剩下0个的盘子数可以由剩下1，2，3个盘子的数量决定，所以只需关注剩下1，2，3个寿司的盘子的数量)
# !dfs(one, two, three) 表示还有几个盘剩下一个/二个/三个寿司 时 吃完的期望次数

# !如何转移:
# 下一次抽中剩0的盘子  概率p0=(n-i-j-k)/n  对应dp[i][j][k]
# 下一次抽中剩1的盘子  概率p1=(i)/n  对应dp[i-1][j][k]
# 下一次抽中剩2的盘子  概率p2=(j)/n  对应dp[i+1][j-1][k]
# 下一次抽中剩3的盘子  概率p3=(k)/n  对应dp[i][j+1][k-1]
# 则有
# !dp[i][j][k] = 1 + dp[i][j][k]*p0 + dp[i-1][j][k]*p1 + dp[i+1][j-1][k]*p2 + dp[i][j+1][k-1]*p3
# 式変形して自己ループを除去 =>
# !dp[i][j][k] = (dp[i−1][j][k]×i + dp[i+1][j−1][k]×j + dp[i][j+1][k−1]×k + N) / (i+j+k)

from collections import Counter
from functools import lru_cache
import sys
from typing import List, Tuple


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 期望dp+counter保存状态

n = int(input())
nums = list(map(float, input().split()))
counter = Counter(nums)

#########################################################################
# memo = [-1.0] * (n + 1) * (n + 1) * (n + 1)


# def dfs(remain1: int, remain2: int, remain3: int) -> float:
#     """counter保存状态

#     这道题的加强版为 九坤t4-筹码游戏-组合.py
#     """
#     if remain1 == remain2 == remain3 == 0:
#         return 0
#     hash_ = remain1 * n * n + remain2 * n + remain3
#     if memo[hash_] != -1:
#         return memo[hash_]

#     div = remain1 + remain2 + remain3
#     res = n / div
#     p1, p2, p3 = remain1 / div, remain2 / div, remain3 / div
#     if remain1:
#         res += p1 * dfs(remain1 - 1, remain2, remain3)
#     if remain2:
#         res += p2 * dfs(remain1 + 1, remain2 - 1, remain3)
#     if remain3:
#         res += p3 * dfs(remain1, remain2 + 1, remain3 - 1)
#     memo[hash_] = res
#     return res


# print(dfs(counter[1], counter[2], counter[3]))
########################################################################
# a, b, c = counter[1], counter[2], counter[3]
# dp = [[[0.0] * (n + 1) for _ in range(n + 1)] for _ in range(n + 1)]
# for i in range(c + 1):
#     for j in range(c + b + 1):
#         for k in range(n + 1):
#             if i + j + k > n:
#                 break
#             remain = i + j + k
#             if remain == 0:
#                 continue
#             dp[i][j][k] += n / remain
#             p1, p2, p3 = i / remain, j / remain, k / remain
#             if i:
#                 dp[i][j][k] += p1 * dp[i - 1][j + 1][k]
#             if j:
#                 dp[i][j][k] += p2 * dp[i][j - 1][k + 1]
#             if k:
#                 dp[i][j][k] += p3 * dp[i][j][k - 1]
# print(dp[c][b][a])


class Solution:
    def eatSushi(self, dishes: List[int]) -> float:
        """
        吃寿司,反向dp,counter元组保存每种寿司盘子的个数
        吃寿司是从有到无
        """

        @lru_cache(None)
        def dfs(cur: Tuple[int, ...], remain: int) -> float:
            """cur表示有cur[i]个盘子剩下i个寿司 remain表示还要吃的寿司个数"""
            if remain == 0:
                return 0

            lis = list(cur)
            step, p = 1, 0
            for i in range(1, len(cur)):  # 吃哪种盘子的寿司
                count = lis[i]
                if count:
                    lis[i] -= 1
                    lis[i - 1] += 1
                    step += (dfs(tuple(lis), remain - 1) / n) * count
                    p += count / n
                    lis[i] += 1
                    lis[i - 1] -= 1
            return step / p

        n = len(dishes)
        max_ = max(dishes)
        counter = [0] * (max_ + 1)
        for num in dishes:
            counter[num] += 1
        return dfs(tuple(counter), sum(dishes))  # 开始时的寿司盘子状态


# 给定每个盘子里的寿司个数
print(Solution().eatSushi([1, 1, 1]))
print(Solution().eatSushi([1, 2, 2, 0, 1, 0, 1, 2, 0, 1]))
