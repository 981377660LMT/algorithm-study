from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 你一个下标从 0 开始、长度为 n 的数组 nums 。一开始，所有元素都是 未染色 （值为 0 ）的。

# 给你一个二维整数数组 queries ，其中 queries[i] = [indexi, colori] 。

# 对于每个操作，你需要将数组 nums 中下标为 indexi 的格子染色为 colori 。

# 请你返回一个长度与 queries 相等的数组 answer ，其中 answer[i]是前 i 个操作 之后 ，相邻元素颜色相同的数目。


# 更正式的，answer[i] 是执行完前 i 个操作后，0 <= j < n - 1 的下标 j 中，满足 nums[j] == nums[j + 1] 且 nums[j] != 0 的数目。class FrequencyTracker:
class Solution:
    def colorTheArray(self, n: int, queries: List[List[int]]) -> List[int]:
        res = []
        same = 0
        cur = [0] * n

        def remove(i: int) -> None:
            nonlocal same
            if i > 0 and cur[i] == cur[i - 1] != 0:
                same -= 1
            if i < n - 1 and cur[i] == cur[i + 1] != 0:
                same -= 1

        def add(i: int, c: int) -> None:
            nonlocal same
            if i > 0 and c == cur[i - 1] != 0:
                same += 1
            if i < n - 1 and c == cur[i + 1] != 0:
                same += 1

        for i, c in queries:
            remove(i)
            cur[i] = c
            add(i, c)
            res.append(same)

        return res
