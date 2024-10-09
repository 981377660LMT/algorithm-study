from typing import List

MOD = int(1e9 + 7)


def countSubsetGcd(nums: List[int]) -> List[int]:
    """gcd为i的子集个数.

    >>> countSubsetGcd([1, 2, 3])
    [0, 5, 1, 1]
    """
    if not nums:
        return [0]
    upper = max(nums) + 1
    c1, c2 = [0] * upper, [0] * upper
    for v in nums:
        c1[v] += 1
    for f in range(1, upper):
        for m in range(f, upper, f):
            c2[f] = (c2[f] + c1[m]) % MOD
    for i in range(1, upper):
        c2[i] = (pow(2, c2[i], MOD) - 1) % MOD  # !gcd为i的子集个数
    for f in range(upper - 1, 0, -1):
        for m in range(f * 2, upper, f):
            c2[f] = (c2[f] - c2[m]) % MOD

    return c2


def _bruteForce(nums: List[int]) -> List[int]:
    from math import gcd

    counter = [0] * (max(nums) + 1)
    for i in range(1, 1 << len(nums)):
        sub = []
        for j in range(len(nums)):
            if i & (1 << j):
                sub.append(nums[j])
        gcd_ = 0
        for s in sub:
            gcd_ = gcd(gcd_, s)
        counter[gcd_] += 1
    return counter


if __name__ == "__main__":
    import random

    for _ in range(100):
        n = 1 + random.randint(0, 10)
        nums = [random.randint(1, 100) for _ in range(n)]
        assert countSubsetGcd(nums) == _bruteForce(nums)
    print("Done!")

    # assert countSubsetGcd([1, 2, 3, 4, 5]) == [0, 5, 4, 2, 1, 0]
    print(countSubsetGcd([1, 2, 3]))  # [0, 5, 4, 2, 1, 0]
