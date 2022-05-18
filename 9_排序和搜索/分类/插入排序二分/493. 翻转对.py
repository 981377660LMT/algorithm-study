# 如果 i < j 且 nums[i] > 2*nums[j] 我们就将 (i, j) 称作一个重要翻转对。输出对数

# 思路:倒序遍历 nums(保证i<j)，遍历过的数我们将它 乘 2，放入 tb 中，同时保持 tb 是按升序排列的。

import bisect
from typing import List


def reversePairs(nums: List[int]) -> int:
    res = 0
    double = []
    for i in reversed(nums):
        # 每次遍历开始比较i与2j的关系,二分法看看是第几位
        res += bisect.bisect_left(double, i)
        # 倒序遍历,乘以二存起来
        j = 2 * i
        bisect.insort(double, j)

    return res


print(reversePairs([2, 4, 3, 5, 1]))

