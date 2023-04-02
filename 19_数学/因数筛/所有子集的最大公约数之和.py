# 所有子集的最大公约数之和
# n<=1e5 nums[i]<=1e5
# !枚举因子

# !如何去重?
# 例如6的时候:
# {1: [6], 2: [6], 3: [6], 6: [6]}

# 1. 按照从大到小的顺序枚举答案(gcd_)
# 2. 减去那些不合法的gcd_的倍数

# 各整数 x に対し、最大公約数が x の倍数となるような集合を数え上げるのは簡単です。しかし、このまま
# 足し合わせると、例えば最大公約数が 6 の場合は x = 1, 2, 3, 6 で重複して数えられてしまうことになります。
# ここで包除原理を用います。このとき、このタイプの問題では、x を大きい順に見ていき、順次「最大公約
# 数がちょうど x であるような場合の数」を求めていくのが簡単でしょう。x について見るとき、x の (2x 以
# 上の) 倍数についてはその「ちょうどの値」が計算されているので、単純な引き算で x についての「ちょうど
# の値」を求めることができます。


from collections import Counter
from math import gcd
from typing import List

MOD = int(1e9 + 7)

# counter = [0] * (k + 1)  # gcd为i的数组的个数
# for gcd_ in range(k, 0, -1):
#     select = k // gcd_
#     count = pow(select, n, MOD)
#     bad = 0  # 重复计算
#     for multi in range(gcd_ * 2, k + 1, gcd_):
#         bad += counter[multi]
#         bad %= MOD
#     counter[gcd_] = (count - bad) % MOD
# res = 0
# for gcd_ in range(1, k + 1):
#     res += counter[gcd_] * gcd_
#     res %= MOD
# return res


def sumOfGcd(nums: List[int]) -> int:
    """计算并返回 nums 的所有 非空 子序列的 最大公约数的 和 。"""
    if len(nums) == 0:
        return 0

    counter = Counter(nums)
    max_ = max(counter)
    multiSum = [0] * (max_ + 1)  # 数组中i的倍数的数的和
    for fac in range(1, max_ + 1):
        for mul in range(fac, max_ + 1, fac):
            multiSum[fac] += counter[mul] * mul
            multiSum[fac] %= MOD

    # gcd为i的倍数的子集的个数
    for gcd_ in range(max_, 0, -1):
        bad = 0
        for multi in range(gcd_ * 2, max_ + 1, gcd_):
            bad += multiSum[multi]
            bad %= MOD
        multiSum[gcd_] = (multiSum[gcd_] - bad) % MOD

    res = 0
    for gcd_ in range(1, max_ + 1):
        res += multiSum[gcd_]
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


print(bruteForce(nums=[6, 10, 3]))
print(bruteForce(nums=[6]))
print(sumOfGcd(nums=[6, 6]))
print(sumOfGcd(nums=[6, 10, 3]))
