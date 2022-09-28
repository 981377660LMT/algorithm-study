# 剑指 Offer 45. 把数组排成最小的数

from functools import cmp_to_key
from typing import List

# !给定一组非负整数，重新排列它们的顺序使之组成一个最大的整数。
# 数组组成最大数(拼接最大序)
# cmp_to_key 将compare函数转换成key


def mergeSort(s1: str, s2: str) -> int:
    """拼接两个字符串，字典序最小"""
    return -1 if s1 + s2 < s2 + s1 else 1


def toMax(nums: List[int]) -> str:
    """拼接最大序"""
    arr = list(map(str, nums))
    arr = sorted(arr, key=cmp_to_key(lambda s1, s2: mergeSort(s1, s2)), reverse=True)
    return "".join(arr)


def toMin(nums: List[int]) -> str:
    """拼接最小序"""
    arr = list(map(str, nums))
    arr = sorted(arr, key=cmp_to_key(lambda s1, s2: mergeSort(s1, s2)))
    return "".join(arr)


print(toMax([10, 1, 2]))  # 2110
print(toMax([3, 30, 34, 5, 9]))  # 9534330
