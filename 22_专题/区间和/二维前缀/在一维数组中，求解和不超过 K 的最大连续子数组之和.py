from typing import List
from sortedcontainers import SortedList


def max_subarray_sum(nums: List[int], k: int) -> int:
    res = float('-inf')
    n = len(nums)
    sum = [0] * (n + 1)
    for i in range(1, n + 1):
        sum[i] = sum[i - 1] + nums[i - 1]

    treeset = SortedList([0])
    for i in range(1, n + 1):
        pre = sum[i]
        left = treeset.bisect_left(pre - k)
        if left < len(treeset):
            res = max(res, pre - treeset[left])
        treeset.add(pre)

    return res


print(max_subarray_sum([1, 2, 3, 4], 8))
print(max_subarray_sum([1, -2, 3, 4], 5))
