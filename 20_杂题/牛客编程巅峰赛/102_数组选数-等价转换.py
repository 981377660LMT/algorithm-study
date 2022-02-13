# n,m<=100000

# 共鸣关系：同时选中+z 不同时选中-z
# 等价转换
class Solution:
    def wwork(self, n: int, m: int, nums: list[int], info: list[list[int]]):
        # write code here
        res = 0
        for i, j, score in info:
            nums[i - 1] += score
            nums[j - 1] += score
            res -= score
        return res + sum(num for num in nums if num > 0)

