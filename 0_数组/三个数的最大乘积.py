from typing import List


def maxTwoProduct(sortedArray: List[int]) -> int:
    """数组中两个数的最大乘积"""
    if len(sortedArray) < 2:
        raise Exception("invalid")
    return max(sortedArray[-1] * sortedArray[-2], sortedArray[0] * sortedArray[1])


def maxThreeProduct(sortedArray: List[int]) -> int:
    """数组中三个数的最大乘积:排序后最后三个数的乘积和最后一个数与开头两个数的乘积."""
    if len(sortedArray) < 3:
        raise Exception("invalid")
    return max(
        sortedArray[-1] * sortedArray[-2] * sortedArray[-3],
        sortedArray[0] * sortedArray[1] * sortedArray[-1],
    )
