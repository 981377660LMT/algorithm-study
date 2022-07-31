INF = int(1e20)


class Solution:
    def solve(self, nums):
        n = len(nums)
        pos, neg = -INF, -INF
        res = 0
        for i in range(n):
            if nums[i] == 0:
                pos, neg = -INF, -INF
            if nums[i] > 0:
                pos, neg = pos + 1, neg + 1
                pos = max(pos, 1)
            if nums[i] < 0:
                pos, neg = neg + 1, pos + 1
                neg = max(neg, 1)
            res = max(res, pos)
        return res
