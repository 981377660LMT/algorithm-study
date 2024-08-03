from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 m x n 的二进制矩阵 grid 。

# 如果矩阵中一行或者一列从前往后与从后往前读是一样的，那么我们称这一行或者这一列是 回文 的。

# 你可以将 grid 中任意格子的值 翻转 ，也就是将格子里的值从 0 变成 1 ，或者从 1 变成 0 。


# 请你返回 最少 翻转次数，使得矩阵中 所有 行和列都是 回文的 ，且矩阵中 1 的数目可以被 4 整除 。
def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


# !分组+dp
class Solution:
    def minFlips(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])

        # !四个角
        counter = [0] * 2
        diff1Info = []
        for i in range(ROW // 2):
            for j in range(COL // 2):
                counter[grid[i][j]] += 1
                counter[grid[i][~j]] += 1
                counter[grid[~i][j]] += 1
                counter[grid[~i][~j]] += 1
                v0, v1 = counter
                diff1Info.append((v0, v1))
                counter[0] = 0
                counter[1] = 0
        diff1 = sum(4 - max2(v0, v1) for v0, v1 in diff1Info)

        # !中心线
        diff2 = 0
        centerOnes = 0
        if ROW & 1:
            centerOnes += sum(grid[ROW // 2])
            for i in range(COL // 2):
                if grid[ROW // 2][i] != grid[ROW // 2][~i]:
                    diff2 += 1
        if COL & 1:
            centerOnes += sum(grid[i][COL // 2] for i in range(ROW))
            for i in range(ROW // 2):
                if grid[i][COL // 2] != grid[~i][COL // 2]:
                    diff2 += 1
        if ROW & 1 and COL & 1:
            centerOnes -= grid[ROW // 2][COL // 2]

        dp = [False] * 4
        dp[centerOnes & 3] = True
        for _ in range(diff2):
            ndp = [False] * 4
            # +1/-1
            for i in range(4):
                if dp[i]:
                    ndp[(i + 1) & 3] = True
                    ndp[(i - 1) & 3] = True
            dp = ndp
        cost = INF
        for i in range(4):
            if dp[i]:
                cost = min(cost, i)
        return diff1 + diff2 + cost


# grid = [[1,0,0],[0,1,0],[0,0,1]]

# print(Solution().minFlips([[1, 0, 0], [0, 1, 0], [0, 0, 1]]))  # 3
# print(Solution().minFlips([[0, 1], [0, 1], [0, 0]]))  # 3
# print(Solution().minFlips([[1], [1]]))  # 3
print(Solution().minFlips([[0], [1], [1], [1], [1]]))  # 2
print(Solution().minFlips([[1], [1], [1]]))  # 2
# [[1],[1],[1],[0]]
