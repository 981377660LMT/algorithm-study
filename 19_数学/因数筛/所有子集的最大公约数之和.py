# 所有子集的最大公约数之和
# n<=1e5 nums[i]<=1e5
# !枚举因子

# !如何去重?
# 例如6的时候:
# {1: [6], 2: [6], 3: [6], 6: [6]}

# 1. 按照从大到小的顺序枚举答案(gcd_)
# 2. 减去那些不合法的gcd_的倍数

# 对于每个整数x，计算最大公约数为x的倍数的集合是很容易的。
# 但是，如果直接相加，例如当最大公约数为6时，将重复计算x = 1, 2, 3, 6的情况。
# 这时可以使用容斥原理。在这种类型的问题中，按照从大到小的顺序逐个考虑x，
# 并逐步计算“最大公约数恰好为x的情况”的数量将是简单的。
# 在考虑x时，对于x的(大于等于2x的)倍数，其“恰好的值”已经被计算出来，
# 因此可以通过简单的减法来计算关于x的“恰好的值”。

# !最大公约数i 的每一个因子 都计算一遍贡献，总和恰好为i，这样就不会重复计算了


from collections import Counter
from math import gcd
from typing import List

MOD = int(1e9 + 7)


def sumOfGcd(nums: List[int]) -> int:
    """计算并返回 nums 的所有 非空 子序列的 最大公约数的 和 。
    => 计算每个(约)数作为gcd的贡献.
    """
    if len(nums) == 0:
        return 0

    counter = Counter(nums)
    max_ = max(counter)
    multiCounter = [0] * (max_ + 1)  # !每个约数的倍数数量
    for fac in range(1, max_ + 1):
        for mul in range(fac, max_ + 1, fac):
            multiCounter[fac] += counter[mul]
            multiCounter[fac] %= MOD
    subCounter = [(pow(2, c, MOD) - 1) % MOD for c in multiCounter]  # !每个约数的非空子集数量
    for gcd_ in range(max_, 0, -1):  # 从大到小遍历，减去倍数的子集数
        for mul in range(gcd_ * 2, max_ + 1, gcd_):
            subCounter[gcd_] -= subCounter[mul]
            subCounter[gcd_] %= MOD
    res = 0
    for gcd_ in range(1, max_ + 1):
        res += gcd_ * subCounter[gcd_]
        res %= MOD
    return res


def bruteForce(nums: List[int]) -> int:
    """暴力枚举所有子序列的最大公约数之和"""
    res = 0
    for i in range(1, 1 << len(nums)):
        sub = []
        for j in range(len(nums)):
            if i & (1 << j):
                sub.append(nums[j])
        res += gcd(*sub)
    return res


if __name__ == "__main__":
    import random

    for _ in range(100):
        n = 1 + random.randint(0, 10)
        nums = [random.randint(1, 100) for _ in range(n)]
        assert sumOfGcd(nums) == bruteForce(nums)
    print("Done!")
