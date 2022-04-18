from typing import List


class Solution:
    def longestWPI(self, hours: List[int]) -> int:

        hours = [1 if h > 8 else -1 for h in hours]
        preSum = {0: -1}
        res, curSum = 0, 0
        for i, num in enumerate(hours):
            curSum += num
            if curSum > 0:
                res = i + 1
            if curSum - 1 in preSum:
                res = max(res, i - preSum[curSum - 1])
            preSum.setdefault(curSum, i)
        return res


print(Solution().longestWPI([9, 9, 6, 0, 6, 6, 9]))
