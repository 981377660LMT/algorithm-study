# n<=1e5
# x,r<=1e8

# 一些圆c1 c2 ... ck 如果后一个圆严格在前一个圆的内部 那么这些圆就叫 'k-target'
# 第i个圆的圆心为(xi,0) 半径为ri
# !选择一些圆 组成最长的k-target子序列

# LIS+条件变形
# !注意到内含条件可变形为 abs(xi-xi+1) < (ri-ri+1)
# 即 `ri-ri+1 + xi-xi+1 > 0` 且 `ri-ri+1 - (xi-xi+1) > 0` (ri包含ri+1)
# !即 xi+ri>xi+1+ri+1 且 xi-ri<xi+1-ri+1
# !因此按照(x+r) 排序后 寻找(x-r)的LIS
# !注意 x+r相等时不能算作内含 所以此时让半径小的圆排前面 (类似于俄罗斯套娃那题)

import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


n = int(input())
circles = []
for _ in range(n):
    x, r = map(int, input().split())
    circles.append((x, r))
circles.sort(key=lambda x: (x[0] + x[1], (x[0] - x[1])), reverse=True)
nums = [x[0] - x[1] for x in circles]

from typing import List
from bisect import bisect_left, bisect_right


# LIS模板
def LIS(nums: List[int], isStrict=True) -> int:
    """求LIS长度"""
    n = len(nums)
    if n <= 1:
        return n

    res = [nums[0]]
    for i in range(1, n):
        pos = bisect_left(res, nums[i]) if isStrict else bisect_right(res, nums[i])
        if pos >= len(res):
            res.append(nums[i])
        else:
            res[pos] = nums[i]

    return len(res)


print(LIS(nums))


# TODO
