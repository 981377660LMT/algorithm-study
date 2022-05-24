"""
换根dp：枚举每个点作为源点，计算整个树的最大流量
"""
# https://www.acwing.com/solution/content/3484/

from collections import defaultdict
from functools import lru_cache


T = int(input())
for _ in range(T):
    n = int(input())
    adjMap = defaultdict(lambda: defaultdict(int))
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        adjMap[u][v] = w
