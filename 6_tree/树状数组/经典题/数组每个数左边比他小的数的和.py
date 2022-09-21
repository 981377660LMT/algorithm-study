"""对数组的每个元素 求出左边比他小的数的和
按照`元素大小排序`加入树状数组
然后查询当前元素左边比他小的数的和
"""


from typing import List
from BIT import BIT1


def leftSmallerSum(nums: List[int]) -> List[int]:
    """对数组的每个元素求出左边比他小的数的和"""
    n = len(nums)
    res = [0] * n
    arr = sorted([(i, num) for i, num in enumerate(nums, 1)], key=lambda x: x[1])
    bit = BIT1(int(1e9 + 10))
    for index, num in arr:
        res[index - 1] = bit.query(index - 1)
        bit.add(index, num)
    return res


print(leftSmallerSum([1, 2, 6, 4]))
