from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def answerQueries(self, nums: List[int], queries: List[int]) -> List[int]:
        """
        返回一个长度为 m 的数组 answer ,
        其中 answer[i] 是 nums 中 元素之和小于等于 queries[i] 的 子序列 的 最大 长度  。
        """
        ...
