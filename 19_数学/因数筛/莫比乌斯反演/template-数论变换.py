from typing import List


def multipleZeta(c1: List[int]):
    """倍数的前缀和变换.

    c2[v] = c1[v] + c1[2*v] + c1[3*v] + ... + c1[k*v]
    """
    upper = len(c1)
    for f in range(1, upper):
        for m in range(2 * f, upper, f):
            c1[f] += c1[m]


def multipleMobius(c2: List[int]):
    """倍数的前缀和逆变换."""
    upper = len(c2)
    for f in range(upper - 1, 0, -1):
        for m in range(2 * f, upper, f):
            c2[f] -= c2[m]


def divisorZeta(c1: List[int]):
    """因子的前缀和变换.

    c2[v] = c1[f1] + c1[f2] + ... + c1[f_n] (v % fi == 0)
    """
    upper = len(c1)
    for f in range(upper - 1, 0, -1):
        for m in range(2 * f, upper, f):
            c1[m] += c1[f]


def divisorMobius(c2: List[int]):
    """因子的前缀和逆变换."""
    upper = len(c2)
    for f in range(1, upper):
        for m in range(2 * f, upper, f):
            c2[m] -= c2[f]


if __name__ == "__main__":
    # 3312. 查询排序后的最大公约数
    # https://leetcode.cn/problems/sorted-gcd-pair-queries/description/
    from bisect import bisect_right
    from itertools import accumulate
    from typing import List

    class Solution:
        def gcdValues(self, nums: List[int], queries: List[int]) -> List[int]:
            upper = max(nums) + 1
            c = [0] * upper
            for v in nums:
                c[v] += 1
            multipleZeta(c)
            for i in range(1, upper):
                c[i] = c[i] * (c[i] - 1) // 2
            multipleMobius(c)

            presum = list(accumulate(c))
            return [bisect_right(presum, q) for q in queries]
