from typing import Callable


def getKth0(left: int, right: int, countNgt: Callable[[int], int], kth: int) -> int:
    """
    给定二分答案的区间[left,right], 求第kth小的答案.
    countNgt: 答案不超过mid时, 满足条件的个数.
    kth从0开始.
    """
    while left <= right:
        mid = left + (right - left) // 2
        if countNgt(mid) <= kth:
            left = mid + 1
        else:
            right = mid - 1
    return right + 1


def getKth1(left: int, right: int, countNgt: Callable[[int], int], kth: int) -> int:
    """
    给定二分答案的区间[left,right], 求第kth小的答案.
    countNgt: 答案不超过mid时, 满足条件的个数.
    kth从1开始.
    """
    while left <= right:
        mid = left + (right - left) // 2
        if countNgt(mid) < kth:
            left = mid + 1
        else:
            right = mid - 1
    return left


EPS = 1e-12


def getKth0Float64(left: float, right: float, countNgt: Callable[[float], int], kth: int) -> float:
    """
    给定二分答案的区间[left,right], 求第kth小的答案.
    countNgt: 答案不超过mid时, 满足条件的个数.
    kth从0开始.
    """
    while left <= right:
        mid = left + (right - left) / 2
        if countNgt(mid) <= kth:
            left = mid + EPS
        else:
            right = mid - EPS
    return right + EPS


def getKth1Float64(left: float, right: float, countNgt: Callable[[float], int], kth: int) -> float:
    """
    给定二分答案的区间[left,right], 求第kth小的答案.
    countNgt: 答案不超过mid时, 满足条件的个数.
    kth从1开始.
    """
    while left <= right:
        mid = left + (right - left) / 2
        if countNgt(mid) < kth:
            left = mid + EPS
        else:
            right = mid - EPS
    return left


if __name__ == "__main__":
    # https://leetcode.cn/problems/kth-smallest-number-in-multiplication-table/submissions/495841768/
    class Solution:
        def findKthNumber(self, m: int, n: int, k: int) -> int:
            def countNGT(mid: int) -> int:
                """有多少个不超过mid的候选"""
                res = 0
                for r in range(1, m + 1):
                    res += min(mid // r, n)
                return res

            assert getKth0(0, n * m, countNGT, k - 1) == getKth1(0, n * m, countNGT, k)
            return getKth0(0, n * m, countNGT, k - 1)
