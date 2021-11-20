# 注意有序
from typing import List


# 将 arr[i] 和 arr[i + len] 相比，如果相等，就说明出现次数超过 1 / 4
class Solution:
    def findSpecialInteger(self, arr: List[int]) -> int:
        interval = len(arr) // 4
        for i in range(0, len(arr) - interval):
            if arr[i] == arr[i + interval]:
                return arr[i]
        return -1


# 如果无序，则需要计数
