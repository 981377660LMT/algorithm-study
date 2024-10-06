from typing import List


def supersetZeta(c1: List[int]):
    """超集的前缀和变换.

    c2[v] = c1[s1] + c1[s2] + ... + c1[s_n] (si & v == v)
    """
    n = len(c1)
    assert n & (n - 1) == 0, "n must be a power of 2"
    i = 1
    while i < n:
        for j in range(n):
            if not j & i:
                c1[j] += c1[j | i]
        i <<= 1


def supersetMobius(c2: List[int]):
    """超集的前缀和逆变换."""
    n = len(c2)
    assert n & (n - 1) == 0, "n must be a power of 2"
    i = 1
    while i < n:
        for j in range(n):
            if not j & i:
                c2[j] -= c2[j | i]
        i <<= 1


def subsetZeta(c1: List[int]):
    """子集的前缀和变换.
    c2[v] = c1[s1] + c1[s2] + ... + c1[v] (v & si == si)
    """
    n = len(c1)
    assert n & (n - 1) == 0, "n must be a power of 2"
    i = 1
    while i < n:
        for j in range(n):
            if not j & i:
                c1[j | i] += c1[j]
        i <<= 1


def subsetMobius(c2: List[int]):
    """子集的前缀和逆变换."""
    n = len(c2)
    assert n & (n - 1) == 0, "n must be a power of 2"
    i = 1
    while i < n:
        for j in range(n):
            if not j & i:
                c2[j | i] -= c2[j]
        i <<= 1


if __name__ == "__main__":
    nums = [1, 2, 3, 4]
    subsetSum = [0] * (1 << len(nums))
    for i in range(len(nums)):
        subsetSum[1 << i] = nums[i]
    subsetZeta(subsetSum)
    print(subsetSum)
