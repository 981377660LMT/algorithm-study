import bisect
from typing import List

# 与9_排序和搜索\力扣加加\插入排序二分\493. 翻转对.py 思路一样
# 右侧即i<j 即求nums[i]>nums[j]的个数
# 反转数组保证i<j 再用二分查找位置


def countSmaller(nums: List[int]):
    res = []
    visited = []
    for i in reversed(nums):
        # 每次遍历开始比较i与j的关系,二分法看看是第几位
        index = bisect.bisect_left(visited, i)
        res.append(index)
        bisect.insort(visited, i)

    return res[::-1]


print(countSmaller([5, 2, 6, 1]))
