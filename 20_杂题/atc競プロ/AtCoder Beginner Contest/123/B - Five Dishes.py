# 预定料理的最短时间
# !只能在10的整数倍时刻点餐，求最短时间

from itertools import permutations
from math import ceil
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# ABC 丼： 調理時間 A 分
# ARC カレー： 調理時間 B 分
# AGC パスタ： 調理時間 C 分
# APC ラーメン： 調理時間 D 分
# ATC ハンバーグ： 調理時間 E 分
if __name__ == "__main__":
    nums = [int(input()) for _ in range(5)]

    res = INF
    for perm in permutations(nums):
        time = 0
        for cost in perm:
            # changeTime
            time = 10 * ceil(time / 10)
            time += cost
        res = min(res, time)

    print(res)
