from typing import List
from functools import lru_cache
from itertools import accumulate


class Solution:
    def stoneGameV(self, stoneValue: List[int]) -> int:
        preSum = [0] + list(accumulate(stoneValue))

        # 左闭右开
        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left + 1 >= right:
                return 0

            res = -0x3FFFFFFF
            for i in range(left + 1, right):
                leftSum, rightSum = preSum[i] - preSum[left], preSum[right] - preSum[i]
                if leftSum > rightSum:
                    res = max(res, rightSum + dfs(i, right))
                elif leftSum < rightSum:
                    res = max(res, leftSum + dfs(left, i))
                else:
                    res = leftSum + max(dfs(left, i), dfs(i, right))
            return res

        return dfs(0, len(stoneValue))


print(Solution().stoneGameV(stoneValue=[6, 2, 3, 4, 5, 5]))
