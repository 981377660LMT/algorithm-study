# 注意这个值域连续比析合树严格,需要是1-k连续
# 小红拿到了一个排列，她想知道有多少区间满足，区间内部的数构成一个排列？
# !排列的定义：1到k，每个数都出现过且恰好出现一次，被称为一个长度为k的排列。例如[2,1,3],[4,3,2,1]都是排列。
# n<=2e5

# !从1往两边扩展并同时检查就好了
from collections import defaultdict
from typing import List


def solve(n: int, perm: List[int]) -> int:
    mp = defaultdict(int)
    for i, v in enumerate(perm):
        mp[v] = i
    left, right = mp[1], mp[1]  # in this range [left,right]
    res = 0
    for i in range(1, n + 1):
        left = min(left, mp[i])
        right = max(right, mp[i])
        if right - left + 1 == i:
            res += 1
    return res
