class Solution:
    def solve(self, nums):
        n = len(nums)
        pos, neg = -int(1e20), -int(1e20)
        res = 0
        for i in range(n):
            if nums[i] == 0:
                pos, neg = -int(1e20), -int(1e20)
            if nums[i] > 0:
                pos, neg = pos + 1, neg + 1
                pos = max(pos, 1)
            if nums[i] < 0:
                pos, neg = neg + 1, pos + 1
                neg = max(neg, 1)
            res = max(res, pos)
        return res
