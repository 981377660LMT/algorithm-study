from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的整数数组 nums 和一个正整数 k 。

# 一个数组的 能量值 定义为：

# 如果 所有 元素都是依次 连续 且 上升 的，那么能量值为 最大 的元素。
# 否则为 -1 。
# 你需要求出 nums 中所有长度为 k 的 子数组 的能量值。


# 请你返回一个长度为 n - k + 1 的整数数组 results ，其中 results[i] 是子数组 nums[i..(i + k - 1)] 的能量值。
class Solution:
    def resultsArray(self, nums: List[int], k: int) -> List[int]:
        preSum = [0]
        for i in range(len(nums) - 1):
            cur = (nums[i + 1] - nums[i]) == 1
            preSum.append(preSum[-1] + cur)

        def check(start):
            return preSum[start + k - 1] - preSum[start] == k - 1

        res = []
        for i in range(len(nums) - k + 1):
            if check(i):
                res.append(nums[i + k - 1])
            else:
                res.append(-1)
        return res
