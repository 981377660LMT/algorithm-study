# enumerateSubstring(子串遍历)
from typing import Any, Generator, List, Tuple


def enumerateSubstring(arr) -> int:
    n = len(arr)
    res = 0
    for i in range(n):
        # !子串arr[i:j+1]左端点为i，这里可以做一些预处理
        ...
        for j in range(i, n):
            print(arr[i : j + 1])
            res += 1
    return res


# def enumerateSubstringByCount(arr: List[Any], target: Any):
#     n = len(arr)
#     nexts = [0] * (n + 1)
#     nexts[n] = n
#     for i in range(n - 1, -1, -1):
#         nexts[i] = i if arr[i] == target else nexts[i + 1]
#     for i in range(n):
#         j, count = i, 1 if arr[i] == target else 0
#         while j != n:
#             ...


def enumerateSubstringGroupByTargetCount(
    arr, target, start
) -> Generator[Tuple[int, int, int], None, None]:
    """遍历左端点为start的子串,每一项形如(left,right,count)."""
    ...


if __name__ == "__main__":
    enumerateSubstring("abc")

    class Solution:
        # 100348. 统计 1 显著的字符串的数量-根号值域
        # https://leetcode.cn/problems/count-the-number-of-substrings-with-dominant-ones/solutions/2860181/mei-ju-by-tsreaper-x830/
        # 给你一个二进制字符串 s。请你统计并返回其中 1 显著 的子字符串的数量。
        # !如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串
        # !n<=4e4
        def numberOfSubstrings(self, s: str) -> int:
            ...
