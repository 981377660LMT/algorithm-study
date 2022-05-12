# 一个数组，求有多少对(i,j)满足nums[i]与nums[j]大于他们中间的所有数 (山谷)
# 枚举中间数 有序容器

from typing import List
from sortedcontainers import SortedList


def solve1(nums: List[int]) -> int:
    left = SortedList()
    right = SortedList(nums)
    res = 0
    for num in nums:
        count1, count2 = len(left) - left.bisect_right(num), len(right) - right.bisect_right(num)
        res += count1 * count2
        left.add(num)
        right.discard(num)
    return res


def solve2(nums: List[int]) -> int:
    res = 0
    for i in range(len(nums)):
        for j in range(i + 1, len(nums)):
            for k in range(j + 1, len(nums)):
                if nums[i] > nums[j] < nums[k]:
                    res += 1
    return res


print(solve2([3, 2, 1, 2, 3]))
print(solve1([3, 2, 1, 2, 3]))

