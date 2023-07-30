from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个由 正 整数组成的数组 nums 。

# 如果数组中的某个子数组满足下述条件，则称之为 完全子数组 ：

# 子数组中 不同 元素的数目等于整个数组不同元素的数目。
# 返回数组中 完全子数组 的数目。


# 子数组 是数组中的一个连续非空序列。
class Solution:
    def countCompleteSubarrays(self, nums: List[int]) -> int:
        # 枚举子数组
        res = 0
        count = len(set(nums))
        for left in range(len(nums)):
            s = set()
            for right in range(left, len(nums)):
                s.add(nums[right])
                if len(s) == count:
                    res += 1
        return res


# originPrint = print
# log = False
# def print(*args, **kwargs):
#     if log:
#         originPrint(*args, **kwargs)
# class Solution:
#     def countCompleteSubarrays(self, nums: List[int]) -> int:
#         n = len(nums)
#         if len(set(nums)) == 1:
#             return n*(n+1)//2
#         d = {num:[] for num in set(nums)}
#         # 最晚的起始位置
#         # 最早的结束位置
#         for i, num in enumerate(nums):
#             d[num].append(i)
#         # print(d)
#         l = d.values()
#         print(l)
#         res = 0
#         for i in range(n):
#             # 从 i 开始
#             m = i # 最小的结束位置
#             for v in l:
#                 if i>v[-1]: break
#                 j = v[bisect_left(v, i)]
#                 m = max(m, j)
#                 print(v, j, m)
#             else:
#                 res += n-m
#                 print(i, m)
#         return res
