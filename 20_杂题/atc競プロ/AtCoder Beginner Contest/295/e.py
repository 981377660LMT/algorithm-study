from bisect import bisect_left, bisect_right
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 0 以上
# M 以下の整数からなる長さ
# N の数列
# A=(A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#  ) があります。

# 今からすぬけくんが以下の操作 1, 2 を順に行います。

# A
# i
# ​
#  =0 を満たすそれぞれの
# i について、
# 1 以上
# M 以下の整数を独立かつ一様ランダムに選び、
# A
# i
# ​
#   をその整数で置き換える。
# A を昇順に並び替える。
# すぬけくんが操作 1, 2 を行ったあとの
# A
# K
# ​
#   の期待値を
# mod 998244353 で出力してください。


fac = [1]
ifac = [1]
for i in range(1, int(1e4) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    nums = [int(i) for i in input().split()]
    # 枚举答案
    notZero = [x for x in nums if x != 0]
    notZero.sort()
    zeroCount = n - len(notZero)

    res = 0
    for kth in range(1, m + 1):
        # <kth的恰好k-1个数
        rawLess = bisect_left(notZero, kth)
        if rawLess > k - 1:
            continue
        left = k - 1 - rawLess  # 还要选几个数严格小于kth [1,kth-1]
        count1 = C(zeroCount, left) * pow(kth - 1, left, MOD) % MOD
        right = zeroCount - left  # 还要选几个数大于等于kth [kth,m]
        count2 = C(zeroCount, right) * pow(m - kth + 1, right, MOD) % MOD
        res += count1 * count2 % MOD
        # 减去不合法的(没选kth)

    # 所有可能
    all_ = pow(zeroCount, m, MOD)
    print(all_, zeroCount, res)
    print((res * pow(all_, MOD - 2, MOD)) % MOD)
