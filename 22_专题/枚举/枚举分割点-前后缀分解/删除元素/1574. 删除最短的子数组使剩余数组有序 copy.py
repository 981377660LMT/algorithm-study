# 删除子数组
# 请你删除一个子数组（可以为空），使得 arr 中剩下的元素是 非递减 的。
# 返回满足题目要求删除的最短子数组的长度。
# 必须连续，那就考虑滑动窗口。

# 滑窗/二分

from typing import List


class Solution:
    def findLengthOfShortestSubarray(self, arr: List[int]) -> int:
        ...
