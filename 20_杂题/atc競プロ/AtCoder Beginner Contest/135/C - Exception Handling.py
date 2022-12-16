# C - Exception Handling
# 移除每个元素后的数组最大值 (前缀+后缀里的最大值)


from itertools import accumulate
from typing import List

INF = int(1e18)


def removeNumMax(nums: List[int]) -> List[int]:
    preMax = [-INF] + list(accumulate(nums, max))
    sufMax = ([-INF] + list(accumulate(nums[::-1], max)))[::-1]
    return [max(preMax[i], sufMax[i + 1]) for i in range(len(nums))]  # type: ignore


n = int(input())
nums = [int(input()) for _ in range(n)]
print(*removeNumMax(nums), sep="\n")
