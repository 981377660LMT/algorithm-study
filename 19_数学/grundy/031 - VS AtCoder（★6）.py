# 有n堆石子 每堆里有蓝色和白色石头Bi和Wi个

# 两人进行取石子博弈
# 每轮有两种选择:
# 1. 某一堆加w个蓝石头，移除1个白石头 (w>=1)
# 2. 移除 1<=k<=b//2 个蓝石头  (b>=2)
# 无法行动的人输

# n<=1e5
# w,b<=50

# nim游戏
# grundy数应用
# 异或不为0 先手胜

from functools import lru_cache
from typing import List, Set
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def mex(nums: Set[int]) -> int:
    res = 0
    while res in nums:
        res += 1
    return res


@lru_cache(None)
def grundy(w: int, b: int) -> int:
    """状态(w,b)时的grundy数"""
    nexts = set()
    if w >= 1:
        nexts.add(grundy(w - 1, b + w))
    if b >= 2:
        for i in range(1, b // 2 + 1):
            nexts.add(grundy(w, b - i))
    return mex(nexts)


n = int(input())
white = tuple(map(int, input().split()))
black = tuple(map(int, input().split()))

g = 0
for w, b in zip(white, black):
    g ^= grundy(w, b)  # 初始的grundy数


if g:
    print("First")
else:
    print("Second")
