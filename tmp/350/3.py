from functools import lru_cache
from heapq import heappop, heappush, nlargest, nsmallest
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 、下标从 0 开始的整数数组 nums ，表示收集不同巧克力的成本。每个巧克力都对应一个不同的类型，最初，位于下标 i 的巧克力就对应第 i 个类型。

# 在一步操作中，你可以用成本 x 执行下述行为：


# 同时对于所有下标 0 <= i < n - 1 进行以下操作， 将下标 i 处的巧克力的类型更改为下标 (i + 1) 处的巧克力对应的类型。如果 i == n - 1 ，则该巧克力的类型将会变更为下标 0 处巧克力对应的类型。
# 假设你可以执行任意次操作，请返回收集所有类型巧克力所需的最小成本。


# 解释：最开始，巧克力的类型分别是 [0,1,2] 。我们可以用成本 1 购买第 1 个类型的巧克力。
# 接着，我们用成本 5 执行一次操作，巧克力的类型变更为 [2,0,1] 。我们可以用成本 1 购买第 0 个类型的巧克力。
# 然后，我们用成本 5 执行一次操作，巧克力的类型变更为 [1,2,0] 。我们可以用成本 1 购买第 2 个类型的巧克力。
# 因此，收集所有类型的巧克力需要的总成本是 (1 + 5 + 1 + 5 + 1) = 13 。可以证明这是一种最优方案。


# 看成巧克力不动,成本旋转


# 单调栈 然后前缀和算下贡献就是On
# !先确定了旋转的次数 cnt，我们对于每块巧克力只需要找到旋转过程中经历的价格最小值即可。
# !预处理每个位置需要旋转的次数
class Solution:
    def minCost(self, nums: List[int], x: int) -> int:
        """当前需要收集第index类型的巧克力,当前轮转了rotate次"""
        n = len(nums)
        min_ = min(nums)
        nums.extend([min_ + x] * n)
        nums.sort()
        return sum(nums[:n])


print(Solution().minCost(nums=[20, 1, 15], x=5))
# [15,150,56,69,214,203]
# 42
print(Solution().minCost(nums=[15, 150, 56, 69, 214, 203], x=6))
