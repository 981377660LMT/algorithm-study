# https://kobejean.github.io/cp-library/cp_library/math/conv/_superset_mobius_fn.py

from typing import List


def fwt_xor(arr: List[int], inv=False) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            a0, a1 = arr[m ^ b], arr[m]
            arr[m ^ b], arr[m] = a0 + a1, a0 - a1
            m = m + 1 | b
    if inv:
        for i in range(n):
            arr[i] //= n


def fwt_or(arr: List[int], inv=False) -> None:
    if not inv:
        _subset_zeta_transform(arr)
    else:
        _subset_moebius_transform(arr)


def fwt_and(arr: List[int], inv=False) -> None:
    if not inv:
        _superset_zeta_transform(arr)
    else:
        _superset_moebius_transform(arr)


def _subset_zeta_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m] += arr[m ^ b]
            m = m + 1 | b


def _subset_moebius_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m] -= arr[m ^ b]
            m = m + 1 | b


def _superset_zeta_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m ^ b] += arr[m]
            m = m + 1 | b


def _superset_moebius_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m ^ b] -= arr[m]
            m = m + 1 | b


if __name__ == "__main__":

    class Solution:
        # 982. 按位与为零的三元组
        # https://leetcode.cn/problems/triples-with-bitwise-and-equal-to-zero/description/
        def countTriplets(self, nums: List[int]) -> int:
            U = 1 << (max(nums).bit_length())
            counter = [0] * U
            for num in nums:
                counter[num] += 1

            fwt_and(counter)
            for i, v in enumerate(counter):
                counter[i] *= v * v
            fwt_and(counter, inv=True)
            return counter[0]

        # 3514. 不同 XOR 三元组的数目 II-fwt
        # https://leetcode.cn/problems/number-of-unique-xor-triplets-ii/solutions/3649377/mei-ju-fu-oulogu-fwt-zuo-fa-pythonjavacg-69r3/
        def uniqueXorTriplets(self, nums: List[int]) -> int:
            """O(UlogU)."""
            U = 1 << (max(nums).bit_length())
            counter = [0] * U
            for num in nums:
                counter[num] += 1

            fwt_xor(counter)
            for i, v in enumerate(counter):
                counter[i] *= v * v
            fwt_xor(counter, inv=True)

            return sum(v > 0 for v in counter)
