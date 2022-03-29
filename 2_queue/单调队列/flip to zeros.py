# 反转用异或描述
class Solution:
    def solve(self, nums):
        res = 0
        flipped = 0
        for num in nums:
            if num ^ flipped:
                flipped ^= 1
                res += 1
        return res


print(Solution().solve(nums=[1, 1, 0]))
