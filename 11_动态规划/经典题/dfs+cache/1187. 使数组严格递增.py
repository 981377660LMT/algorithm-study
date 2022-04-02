from typing import List
from functools import lru_cache
from bisect import bisect_right

# 每一步「操作」中，你可以分别从 arr1 和 arr2 中各选出一个索引，
# 分别为 i 和 j，0 <= i < arr1.length 和 0 <= j < arr2.length，
# 然后进行赋值运算 arr1[i] = arr2[j]。
# 如果无法让 arr1 严格递增，请返回 -1。

# 返回使 arr1 严格递增所需要的最小「操作」数（可能为 0）。
# 1 <= arr1.length, arr2.length <= 2000

# dfs(i, prev): "i" represents index in arr1. "prev" represents the previous element in arr1 after swap (or maybe not swap).
# 每遍历一个arr1中的元素，就要看是否需要交换，需要比pre大；并作为pre记录进行下一次dfs

INF = 0x7FFFFFFF


class Solution:
    def makeArrayIncreasing(self, arr1: List[int], arr2: List[int]) -> int:
        @lru_cache(None)
        def dfs(i: int, pre: int) -> int:
            if i >= len(arr1):
                return 0

            pos = bisect_right(store, pre)
            swap = 1 + dfs(i + 1, store[pos]) if pos < len(store) else INF
            noswap = dfs(i + 1, arr1[i]) if arr1[i] > pre else INF
            return min(swap, noswap)

        store = sorted(set(arr2))
        res = dfs(0, -INF)
        return res if res != INF else -1


print(Solution().makeArrayIncreasing(arr1=[1, 5, 3, 6, 7], arr2=[1, 3, 2, 4]))
# 输出：1
# 解释：用 2 来替换 5，之后 arr1 = [1, 2, 3, 6, 7]。
