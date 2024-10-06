# 所有数对的lcm之和模mod(最小公倍数之和).
# C - LCMs
# https://atcoder.jp/contests/agc038/tasks/agc038_c
# N<=2e5,A[i]<=1e6

from typing import List


MOD = 998244353

INV2 = pow(2, MOD - 2, MOD)


def lcmSum(nums: List[int]) -> int:
    upper = max(nums) + 1
    c1, c2 = [0] * upper, [0] * upper  # 每个数的所有倍数之和
    for v in nums:
        c1[v] += v
    for f in range(1, upper):
        for m in range(f, upper, f):
            c2[f] = (c2[f] + c1[m]) % MOD
    for i in range(1, upper):
        c2[i] = c2[i] * c2[i] % MOD
    for f in range(upper - 1, 0, -1):
        for m in range(f * 2, upper, f):
            c2[f] = (c2[f] - c2[m]) % MOD

    # 排除i=j
    for v in nums:
        c2[v] = (c2[v] - v * v) % MOD
    # 排除i>j
    for i in range(1, upper):
        c2[i] = c2[i] * INV2 % MOD

    res = 0
    for i in range(1, upper):
        inv = pow(i, MOD - 2, MOD)
        res += c2[i] * inv % MOD
        res %= MOD
    return res


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    print(lcmSum(A))
