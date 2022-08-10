# 开始每个盘子里有ai个寿司
# 每次随机一个盘子编号 吃掉那个盘子的一个寿司
# 求吃完所有盘子里寿司需要的步数的期望值
# n<=300
# 1<=ai<=3

# !状态如何定义:dfs(one, two, three) 表示还有几个盘剩下一个/二个/三个寿司

from collections import Counter
from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 概率dp

n = int(input())
nums = list(map(float, input().split()))
counter = Counter(nums)


@lru_cache(None)
def dfs(remain1: int, remain2: int, remain3: int) -> float:
    if remain1 == remain2 == remain3 == 0:
        return 0
    sum_ = remain1 + remain2 + remain3
    res = n / sum_  # 吃到寿司
    p1, p2, p3 = remain1 / sum_, remain2 / sum_, remain3 / sum_
    if remain1:
        res += p1 * dfs(remain1 - 1, remain2, remain3)
    if remain2:
        res += p2 * dfs(remain1, remain2 - 1, remain3)
    if remain3:
        res += p3 * dfs(remain1, remain2, remain3 - 1)
    return res


print(dfs(counter[1], counter[2], counter[3]))
