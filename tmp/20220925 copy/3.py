from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 对于 k <= i < n - k 之间的一个下标 i ，如果它满足以下条件，我们就称它为一个 好 下标：

# 下标 i 之前 的 k 个元素是 非递增的 。
# 下标 i 之后 的 k 个元素是 非递减的 。
# 按 升序 返回所有好下标。


class Solution:
    def goodIndices(self, nums: List[int], k: int) -> List[int]:
        n = len(nums)
        # 连续非递增长度
        dp1 = [1] * n
        for i in range(1, n):
            if nums[i] <= nums[i - 1]:
                dp1[i] = dp1[i - 1] + 1

        # 连续非递减长度
        dp2 = [1] * n
        for i in range(1, n):
            if nums[i] >= nums[i - 1]:
                dp2[i] = dp2[i - 1] + 1

        res = []
        for i in range(k, n - k):
            if dp1[i - 1] >= k and dp2[i + k] >= k:
                res.append(i)
        return res


print(Solution().goodIndices(nums=[2, 1, 1, 1, 3, 4, 1], k=2))
