# # For how many positive integers at most N is the product of the digits at most K
# # 1-n中 各位上的积 不超过k的小于等于n的数的个数


# from functools import lru_cache


# @lru_cache(None)
# def cal(upper: int, k: int) -> int:
#     @lru_cache(None)
#     def dfs(pos: int, pre: int, isLimit: bool) -> int:
#         """当前在第pos位，前一位为pre，isLimit表示是否贴合上界"""
#         if pos == len(nums):
#             return 1

#         res = 0
#         up = nums[pos] if isLimit else 1
#         for cur in range(up + 1):
#             if pre == cur == 1:
#                 continue
#             res += dfs(pos + 1, cur, (isLimit and cur == up))
#         return res

#     nums = list(map(int, str(upper)))
#     return dfs(0, -1, True)


# class Solution:
#     def main(self, n: int, k: int) -> int:
#         """For how many positive integers at most N is the product of the digits at most K"""
#         return cal(n, k) - cal(0, k)


# print(Solution().main(n=13, k=2))  # 5
# # print(Solution().main(n=100, k=80))  # 99
# # print(Solution().main(n=1000000000000000000, k=1000000000))  # 841103275147365677
