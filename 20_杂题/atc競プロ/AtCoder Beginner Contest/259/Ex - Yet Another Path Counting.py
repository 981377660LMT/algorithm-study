# AtCoder Beginner Contest 259 - SGColin的文章 - 知乎
# https://zhuanlan.zhihu.com/p/539701972

# 给定一个矩阵 Anxn (1≤n ≤400)，从某个格子出发，每次可以向右或向下走。
# !问起点终点的数字相同的路径有多少条?

# 分情况
# 暴力+dp两种算法的结合

# !注意:product/combinations会比正常循环慢一些


from collections import defaultdict
import sys
import os
from typing import List, Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353

Point = Tuple[int, int]
DIR2 = [[0, 1], [1, 0]]
MAX = 800
C = [[0] * (MAX + 5) for _ in range((MAX + 5))]
for i in range(MAX + 2):
    C[i][0] = 1
    for j in range(1, i + 1):
        C[i][j] = (C[i - 1][j - 1] + C[i - 1][j]) % MOD


def main() -> None:
    def strategy1(points: List[Point]) -> int:
        """brute force 枚举起始点"""
        res = len(points)  # 自己走到自己
        for i in range(len(points)):
            sr, sc = points[i]
            for j in range(i + 1, len(points)):
                er, ec = points[j]
                if er >= sr and ec >= sc:
                    res += C[er - sr + ec - sc][ec - sc]
                    res %= MOD
        return res

    def strategy2(points: List[Point]) -> int:
        """dp (当前横坐标，当前纵坐标) 为终点的路径数"""

        dp = [[0] * n for _ in range(n)]
        for r, c in points:
            dp[r][c] = 1
        for r in range(n):
            for c in range(n):
                # !注意这里不要用 for 循环
                if r + 1 < n:
                    dp[r + 1][c] += dp[r][c]
                    dp[r + 1][c] %= MOD
                if c + 1 < n:
                    dp[r][c + 1] += dp[r][c]
                    dp[r][c + 1] %= MOD
        res = 0
        for r, c in points:
            res += dp[r][c]
        return res % MOD

    n = int(input())
    matrix = []
    counter = defaultdict(list)
    for r in range(n):
        row = tuple(map(int, input().split()))
        for c, num in enumerate(row):
            counter[num].append((r, c))
        matrix.append(row)

    res = 0
    for points in counter.values():
        if len(points) <= n:
            res += strategy1(points)  # !枚举+组合数 最多O(n^3)的计算量
            res %= MOD
        else:
            res += strategy2(points)  # !dp 最多O(n^3)的计算量 因为最多n种数字
            res %= MOD

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
