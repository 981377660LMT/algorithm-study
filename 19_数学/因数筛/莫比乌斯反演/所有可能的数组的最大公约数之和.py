# E - Sum of gcd of Tuples (Hard)
# https://atcoder.jp/contests/abc162/tasks/abc162_e
# !所有可能的数组的最大公约数之和
# 对一个长为n的数组,每个数的范围在[1,k]之间
# 求这样的所有数组(k^n种)的最大公约数之和模1e9+7
# n<=1e5 k<=1e5
# O(klogn)

# 从大到小枚举每个gcd作为答案(计算贡献)
# !注意减去gcd的倍数(这些不是答案)


MOD = int(1e9 + 7)


def sumofGcd(n: int, k: int) -> int:
    counter = [0] * (k + 1)  # gcd为i的倍数的数组的个数
    for f in range(1, k + 1):
        c = k // f
        counter[f] = pow(c, n, MOD)
    for f in range(k, 0, -1):  # mobius
        for m in range(f * 2, k + 1, f):
            counter[f] = (counter[f] - counter[m]) % MOD
    res = 0
    for gcd_ in range(1, k + 1):
        res += counter[gcd_] * gcd_
        res %= MOD
    return res


if __name__ == "__main__":
    N, K = map(int, input().split())
    print(sumofGcd(N, K))
