# n个物品 买k个 求价格不超过P的方案数
# k<=N<=40

# !半分全列挙 时间复杂度 n*2^(n/2)
from bisect import bisect_right
from itertools import combinations
import sys


input = sys.stdin.readline


n, k, p = map(int, input().split())
costs = list(map(int, input().split()))
left, right = costs[: n // 2], costs[n // 2 :]
counter = [[] for _ in range(k + 1)]

# 1. 折半搜左边 排序
for i in range(k + 1):
    for sub in combinations(left, i):
        counter[i].append(sum(sub))

for nums in counter:
    nums.sort()

# 2. 折半枚举右边 二分左边
res = 0
for i in range(k + 1):
    for sub in combinations(right, i):
        sum_ = sum(sub)
        res += bisect_right(counter[k - i], p - sum_)

print(res)
