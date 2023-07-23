from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始的整数数组 nums 和一个正整数 x 。

# 你 一开始 在数组的位置 0 处，你可以按照下述规则访问数组中的其他位置：

# 如果你当前在位置 i ，那么你可以移动到满足 i < j 的 任意 位置 j 。
# 对于你访问的位置 i ，你可以获得分数 nums[i] 。
# 如果你从位置 i 移动到位置 j 且 nums[i] 和 nums[j] 的 奇偶性 不同，那么你将失去分数 x 。
# 请你返回你能得到的 最大 得分之和。

# 注意 ，你一开始的分数为 nums[0] 。


class Solution:
    def maxScore(self, nums: List[int], x: int) -> int:
        # (even, odd)
        dp = [-INF, nums[0]] if nums[0] & 1 else [nums[0], -INF]
        for v in nums[1:]:
            ndp = [-INF, -INF]
            for pre in range(2):
                # 选v
                if (v & 1) == pre:
                    ndp[pre] = max(ndp[pre], dp[pre] + v)
                else:
                    ndp[pre ^ 1] = max(ndp[pre ^ 1], dp[pre] + v - x)
                # 不选v
                ndp[pre] = max(ndp[pre], dp[pre])
            dp = ndp

        return max(dp)


# nums = [2,3,6,1,9,2], x = 5
print(Solution().maxScore(nums=[2, 3, 6, 1, 9, 2], x=5))
