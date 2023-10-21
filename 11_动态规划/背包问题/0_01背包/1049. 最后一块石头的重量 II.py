# 1049. 最后一块石头的重量 II


from typing import List

INF = int(1e20)


class Solution:
    def lastStoneWeightII(self, nums: List[int]) -> int:
        dp = set([0])
        for num in nums:
            ndp = set()
            for pre in dp:
                ndp.add(pre + num)
                ndp.add(pre - num)
            dp = ndp

        res = INF
        for num in dp:
            if num >= 0:
                res = min(res, num)
        return res


print(Solution().lastStoneWeightII(nums=[1, 2, 5]))
