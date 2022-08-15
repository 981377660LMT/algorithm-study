# 找到最长的子序列使得相邻两项符号交替
class Solution:
    def solve(self, nums):
        up, down = 0, 0
        for num in nums:
            if num > 0:
                up = max(up, down + 1)
            elif num < 0:
                down = max(down, up + 1)
        return max(up, down)
