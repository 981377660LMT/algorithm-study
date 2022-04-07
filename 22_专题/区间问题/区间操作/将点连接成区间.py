# Contiguous Intervals
# 将点连接成区间
class Solution:
    def solve(self, nums):
        res = []
        for num in sorted(nums):
            if not res or res[-1][1] < num - 1:
                res.append([num, num])
            else:
                res[-1][1] = num
        return res


print(Solution().solve(nums=[1, 3, 2, 7, 6]))
# [
#     [1, 3],
#     [6, 7]
# ]
