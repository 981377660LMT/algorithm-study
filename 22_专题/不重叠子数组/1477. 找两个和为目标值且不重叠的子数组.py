from typing import List
from itertools import accumulate


# 请你在 arr 中找 两个互不重叠的子数组 且它们的和都等于 target 。
# 可能会有多种方案，请你返回满足要求的两个子数组长度和的 最小值 。
# 请返回满足要求的最小长度和，如果无法找到这样的两个子数组，请返回 -1 。

# 1 <= arr.length <= 10^5

# 思路：哈希表+前缀和
# 记录每个位置的子数组最短长度
INF = 0x7FFFFFFF


class Solution:
    def minSumOfLengths(self, arr: List[int], target: int) -> int:
        lookup = {0: -1}
        # 之前看的和为target的最短子数组长度；记录，维护看过的最小值
        shortest_len_record = [INF] * len(arr)
        res = shortest_len = INF
        for i, running_sum in enumerate(accumulate(arr)):
            if running_sum - target in lookup:
                end = lookup[running_sum - target]
                if end != -1:
                    # i - end表示当前候选 best_till[end]示前面一个最短的target子数组
                    res = min(res, i - end + shortest_len_record[end])
                shortest_len = min(shortest_len, i - end)
            shortest_len_record[i] = shortest_len
            lookup[running_sum] = i
            print(lookup, shortest_len_record)
        return -1 if res == INF else res


print(Solution().minSumOfLengths([3, 4, 7, 7], 7))

