"""A - Make it Zigzag

输入为1-2n的全排列
邻位交换使得原数组变为摆动数组(zigzag)
要求不能超过一半的操作次数
求最小操作次数、需要邻位交换的元素下标

考察奇偶位
循环奇(偶)数位的数字 看一下它和左右两边3个数哪个最大(小) 把最大(小) 的数字换到中间去就行了
特殊处理一下开始的第一个数(没有左边)
"""

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(nums: List[int]) -> List[int]:
    """返回需要邻位交换的的元素下标 最多交换n/2次"""
    n = len(nums)
    res = []

    for i in range(0, n, 2):  # 遍历偶数位 要让偶数位为谷底
        cand = -1
        min_ = nums[i]
        pre = nums[i - 1] if i - 1 >= 0 else INF
        next = nums[i + 1] if i + 1 < n else INF
        if pre < min_:
            min_ = pre
            cand = i - 1
        if next < min_:
            min_ = next
            cand = i + 1

        if cand != -1:
            res.append(min(cand, i))
            nums[i], nums[cand] = nums[cand], nums[i]

    return res


n = int(input())
perm = list(map(int, input().split()))
res = solve(perm)
print(len(res))
print(*[int(num) + 1 for num in res])
