# O(n^2)
# 不如归并排序
import bisect
from typing import List


def reversePairs(nums: List[int]) -> int:
    res = 0
    double = []
    for i in reversed(nums):
        # 每次遍历开始比较i与2j的关系,二分法看看是第几位
        res += bisect.bisect_left(double, i)
        bisect.insort(double, i)
    return res


print(reversePairs([2, 4, 3, 5, 1]))

