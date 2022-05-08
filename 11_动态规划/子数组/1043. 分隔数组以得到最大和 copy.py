from typing import List
from functools import lru_cache

# 枚举分割点


class Solution:
    def maxSumAfterPartitioning(self, arr: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= len(arr):
                return 0
            res = 0
            for i in range(1, k + 1):
                if index + i > len(arr):
                    continue
                sub = arr[index : index + i]
                subMax = max(sub)
                res = max(res, subMax * i + dfs(index + i))
            return res

        return dfs(0)


print(Solution().maxSumAfterPartitioning(arr=[1, 15, 7, 9, 2, 5, 10], k=3))
# 输出：84
# 解释：
# 因为 k=3 可以分隔成 [1,15,7] [9] [2,5,10]，结果为 [15,15,15,9,10,10,10]，和为 84，是该数组所有分隔变换后元素总和最大的。
# 若是分隔成 [1] [15,7,9] [2,5,10]，结果就是 [1, 15, 15, 15, 10, 10, 10] 但这种分隔方式的元素总和（76）小于上一种。
