# from typing import List, Tuple, Optional
# from collections import defaultdict, Counter


# MOD = int(1e9 + 7)
# INF = int(1e20)

# from functools import lru_cache

# # 1 <= n <= 109


# def cal(upper: int) -> int:
#     @lru_cache(None)
#     def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, visited: int, isOk: bool) -> int:
#         """当前在第pos位，hasLeadingZero表示有前导0，isLimit表示是否贴合上界"""
#         if pos == len(nums):
#             return isOk and not hasLeadingZero

#         res = 0
#         up = nums[pos] if isLimit else 9
#         for cur in range(up + 1):
#             if hasLeadingZero and cur == 0:
#                 res += dfs(pos + 1, True, (isLimit and cur == up), visited, isOk)
#             else:
#                 res += dfs(
#                     pos + 1,
#                     False,
#                     (isLimit and cur == up),
#                     (visited | (1 << cur)),
#                     (isOk and (visited & (1 << cur) == 0)),
#                 )
#         return res

#     nums = list(map(int, str(upper)))
#     res = dfs(0, True, True, 0, True)
#     dfs.cache_clear()
#     return res


# class Solution:
#     def countSpecialNumbers(self, n: int) -> int:
#         return cal(n) - cal(0)
