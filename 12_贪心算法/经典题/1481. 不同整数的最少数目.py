from typing import List
from collections import Counter

# 现需要从数组中恰好移除 k 个元素，请找出移除后数组中不同整数的最少数目。
class Solution:
    def findLeastNumOfUniqueInts(self, arr: List[int], k: int) -> int:
        c = Counter(arr)
        arr = sorted(arr, key=lambda x: (c[x], x))
        return len(set(arr[k:]))


print(Solution().findLeastNumOfUniqueInts(arr=[4, 3, 1, 1, 3, 3, 2], k=3))
# 输出：2
# 解释：先移除 4、2 ，然后再移除两个 1 中的任意 1 个或者三个 3 中的任意 1 个，最后剩下 1 和 3 两种整数。
