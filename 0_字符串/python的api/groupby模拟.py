# 模拟groupby的返回值
# 将连续子串分组打包输出


from typing import List


class Solution:
    def solve(self, nums: List[int]) -> List[List[int]]:
        res = []
        pre = ''
        # for


print(Solution().solve(nums=[4, 4, 1, 6, 6, 6, 1, 1, 1, 1]))
# [
#     [4, 4],
#     [1],
#     [6, 6, 6],
#     [1, 1, 1, 1]
# ]
