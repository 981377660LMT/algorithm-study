# 给定已经按 升序 排列、由不同整数组成的数组 arr，
# 返回满足 arr[i] == i 的最小索引 i
# 寻找最小的不动点
from typing import List


class Solution:
    def fixedPoint(self, arr: List[int]) -> int:
        left, right = 0, len(arr) - 1
        res = -1
        while left <= right:
            mid = (left + right) // 2
            if arr[mid] == mid:
                res = mid
                right = mid - 1
            elif arr[mid] < mid:
                left = mid + 1
            else:
                right = mid - 1
        return res


print(Solution().fixedPoint([-10, -5, 0, 3, 7]))
