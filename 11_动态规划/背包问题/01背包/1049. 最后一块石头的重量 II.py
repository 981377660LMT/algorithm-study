from functools import lru_cache


class Solution:
    def solve(self, nums):
        @lru_cache(None)
        def dfs(index: int, curSum: int) -> int:
            if index >= len(nums):
                return curSum

            cur = nums[index]
            #  If the sum is not negative we keep going with the reduced sum, otherwise we just skip the number
            if curSum - 2 * cur < 0:
                return dfs(index + 1, curSum)
            return min(dfs(index + 1, curSum - 2 * cur), dfs(index + 1, curSum))

        sum_ = sum(nums)
        return dfs(0, sum_)


class Solution2:
    def solve(self, nums):
        dp = set([0])
        for num in nums:
            ndp = set()
            for pre in dp:
                ndp.add(pre + num)
                ndp.add(pre - num)
            dp = ndp

        res = int(1e20)
        for num in dp:
            if num >= 0:
                res = min(res, num)
        return res


print(Solution().solve(nums=[1, 2, 5]))

