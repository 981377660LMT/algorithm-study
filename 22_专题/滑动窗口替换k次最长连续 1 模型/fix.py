from typing import Sequence, TypeVar


def max2(a: int, b: int) -> int:
    return a if a > b else b


T = TypeVar("T")


def fix(arr: Sequence[T], target: T, k: int) -> int:
    """
    允许将数组中的任意字符替换为target字符k次，求target字符的最大连续长度.

    :param arr: 源字符串
    :param target: 关心的字符
    :param k: 可替换k次
    :return: target 最大连续长度
    """
    left = 0
    res = 0
    for right in range(len(arr)):
        k -= arr[right] != target
        while k < 0:
            k += arr[left] != target
            left += 1
        res = max2(res, right - left + 1)
    return res
