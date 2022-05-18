from functools import lru_cache
from typing import List


class Solution:
    def canReach(self, arr: List[int], start: int) -> bool:
        """注意每个点只要访问一次，所以用visited来标记是否访问过。"""

        def dfs(index: int) -> bool:
            if index < 0 or index >= len(arr) or arr[index] == -1:
                return False
            if arr[index] == 0:
                return True
            tmp = arr[index]
            arr[index] = -1
            return dfs(index + tmp) or dfs(index - tmp)

        return dfs(start)

