# Given an array arr, and an integer k,
# find the kth next greater element for any element arr[i],
# or -1 if doesn’t exist.

# input:
# arr = [1,4,2,5,3]
# k = 2

# output:
# [3, -1, 5, -1, -1]

# constraint
# len(arr) < 10^5, k < 50, arr[i] < 10^9


from typing import List
from sortedcontainers import SortedList


def findNextKthLarge(nums: List[int], k: int) -> List[int]:
    """对每个数,寻找右侧比自己`严格`大的数中的第k个(kth next greater)

    倒序遍历+SortedList二分查找
    """
    n = len(nums)
    res = [-1] * n
    sl = SortedList()
    for i in range(n - 1, -1, -1):
        cur = nums[i]
        pos = sl.bisect_right(cur)
        if pos + k - 1 < len(sl):
            res[i] = sl[pos + k - 1]  # type: ignore
        sl.add(cur)
    return res


assert findNextKthLarge([1, 4, 2, 2, 2], 2) == [2, -1, -1, -1, -1]
assert findNextKthLarge([1, 4, 2, 5, 3], 2) == [3, -1, 5, -1, -1]
