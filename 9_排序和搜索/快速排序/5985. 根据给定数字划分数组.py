# 看似快速排序
# 不存在稳定且 in-place 的快速排序算法
# `稳定快排就是堆排序了`，比如 sort.Stable


class Solution:
    def pivotArray(self, nums: list[int], pivot: int) -> list[int]:
        return (
            [num for num in nums if num < pivot]
            + [num for num in nums if num == pivot]
            + [num for num in nums if num > pivot]
        )

