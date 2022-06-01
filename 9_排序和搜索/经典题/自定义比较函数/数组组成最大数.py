from functools import cmp_to_key
from typing import List

# 给定一组非负整数，重新排列它们的顺序使之组成一个最大的整数。
# 数组组成最大数


# cmp_to_key 将compare函数转换成key


# 示例 1：


# 输入：[10,1,2]
# 输出：2110
# 示例 2：


# 输入：[3,30,34,5,9]
# 输出：9534330
def toMax(nums: List[int]) -> str:
    arr = list(map(str, nums))
    arr = sorted(arr, key=cmp_to_key(lambda s1, s2: int((s2 + s1)) - int((s1 + s2))))
    return ''.join(arr)


print(toMax([10, 1, 2]))
print(toMax([3, 30, 34, 5, 9]))

