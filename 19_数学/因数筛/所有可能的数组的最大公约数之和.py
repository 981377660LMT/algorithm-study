# https://www.geeksforgeeks.org/sum-of-gcd-of-all-possible-sequences/
# !所有可能的数组的最大公约数之和
# 对一个长为n的数组,每个数的范围在[1,k]之间
# 求这样的所有数组(k^n种)的最大公约数之和模1e9+7
# n<=1e5 k<=1e5
# O(klogn)

# 从大到小枚举每个gcd作为答案(计算贡献)
# !注意减去gcd的倍数(这些不是答案)

from itertools import product
from math import gcd
from random import randint


MOD = int(1e9 + 7)


def sumofGcd(n: int, k: int) -> int:
    counter = [0] * (k + 1)  # gcd为i的倍数的数组的个数
    for gcd_ in range(k, 0, -1):
        select = k // gcd_
        count = pow(select, n, MOD)
        bad = 0  # 重复计算
        for multi in range(gcd_ * 2, k + 1, gcd_):
            bad += counter[multi]
            bad %= MOD
        counter[gcd_] = (count - bad) % MOD
    res = 0
    for gcd_ in range(1, k + 1):
        res += counter[gcd_] * gcd_
        res %= MOD
    return res


if __name__ == "__main__":
    for _ in range(10):
        n = randint(1, 8)
        k = randint(1, 8)
        res = 0
        for nums in product(range(1, k + 1), repeat=n):
            res += gcd(*nums)
        if res != sumofGcd(n, k):
            print(n, k)
            print(res, sumofGcd(n, k))
            break
    print("pass")
