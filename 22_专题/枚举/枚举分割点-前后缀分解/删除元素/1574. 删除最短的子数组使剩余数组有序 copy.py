# 删除子数组
# 请你删除一个子数组（可以为空），使得 arr 中剩下的元素是 非递减 的。
# 返回满足题目要求删除的最短子数组的长度。
# 必须连续，那就考虑滑动窗口。

# !滑窗/二分

from typing import List


class Solution:
    def findLengthOfShortestSubarray(self, arr: List[int]) -> int:
        n = len(arr)
        i, j = 0, n - 1
        while i + 1 < n and arr[i] <= arr[i + 1]:
            i += 1
        while j - 1 >= 0 and arr[j - 1] <= arr[j]:
            j -= 1
        res = min(j, n - 1 - i)
        if res == 0:
            return 0

        right = j
        for left in range(i + 1):
            while right < n and arr[left] > arr[right]:
                right += 1
            res = min(res, right - left - 1)
        return res
