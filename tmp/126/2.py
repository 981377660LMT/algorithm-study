from heapq import heapify
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的正整数数组 nums 。

# 同时给你一个长度为 m 的二维操作数组 queries ，其中 queries[i] = [indexi, ki] 。

# 一开始，数组中的所有元素都 未标记 。

# 你需要依次对数组执行 m 次操作，第 i 次操作中，你需要执行：


# 如果下标 indexi 对应的元素还没标记，那么标记这个元素。
# 然后标记 ki 个数组中还没有标记的 最小 元素。如果有元素的值相等，那么优先标记它们中下标较小的。如果少于 ki 个未标记元素存在，那么将它们全部标记。
# 请你返回一个长度为 m 的数组 answer ，其中 answer[i]是第 i 次操作后数组中还没标记元素的 和 。
class Solution:
    def unmarkedSumArray(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        allSum = sum(nums)
        curSum = 0
        visited = [False] * len(nums)
        sl = SortedList((nums[i], i) for i in range(len(nums)))

        res = []
        for index, k in queries:
            if not visited[index]:
                visited[index] = True
                curSum += nums[index]
                sl.remove((nums[index], index))
            while k > 0 and sl:
                cur, i = sl[0]
                if not visited[i]:
                    visited[i] = True
                    curSum += cur
                    k -= 1
                sl.remove((cur, i))
            res.append(allSum - curSum)
        return res
