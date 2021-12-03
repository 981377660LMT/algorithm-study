from typing import List
from collections import Counter


# 给你一个整数数组 arr。你可以从中选出一个整数集合，并删除这些整数在数组中的每次出现。
# 返回 至少 能删除数组中的一半整数的整数集合的最小大小。
# len(arr)为偶数
class Solution:
    def minSetSize(self, arr: List[int]) -> int:
        top = Counter(arr).most_common()
        res = remove = 0
        for _, count in top:
            res += 1
            remove += count
            if remove >= len(arr) // 2:
                return res


print(Solution().minSetSize(arr=[3, 3, 3, 3, 5, 5, 5, 2, 2, 7]))
# 输出：2
# 解释：选择 {3,7} 使得结果数组为 [5,5,5,2,2]、长度为 5（原数组长度的一半）。
# 大小为 2 的可行集合有 {3,5},{3,2},{5,2}。
# 选择 {2,7} 是不可行的，它的结果数组为 [3,3,3,3,5,5,5]，新数组长度大于原数组的二分之一。

