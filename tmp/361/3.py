from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums ，以及整数 modulo 和整数 k 。

# 请你找出并统计数组中 趣味子数组 的数目。

# 如果 子数组 nums[l..r] 满足下述条件，则称其为 趣味子数组 ：

# 在范围 [l, r] 内，设 cnt 为满足 nums[i] % modulo == k 的索引 i 的数量。并且 cnt % modulo == k 。
# 以整数形式表示并返回趣味子数组的数目。


# 注意：子数组是数组中的一个连续非空的元素序列。
class Solution:
    def countInterestingSubarrays(self, nums: List[int], modulo: int, k: int) -> int:
        ok = [int(x % modulo == k) for x in nums]
        preSum = defaultdict(int, {0: 1})
        res, curSum = 0, 0

        for i, num in enumerate(ok):
            curSum += num
            curSum %= modulo
            res += preSum[(curSum - k) % modulo]
            preSum[curSum] += 1
        return res


# [11,12,21,31]
# 10
# 1
# [2,2,5]
# 3
# 2


# 5
print(Solution().countInterestingSubarrays([11, 12, 21, 31], 10, 1))


# 2
print(Solution().countInterestingSubarrays([2, 2, 5], 3, 2))
