from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# カード 1, カード 2, …, カード N の N 枚のカードがあり、 カード i (1≤i≤N) には整数 A
# i
# ​
#   が書かれています。

# K=1,2,…,N について、次の問題を解いてください。

# カード 1, カード 2, …, カード K の K 枚のカードが入っている袋があります。
# 次の操作を 2 回繰り返し、記録された数を順に x,y とします。

# 袋から無作為にカードを 1 枚取り出し、カードに書かれている数を記録する。その後、カードを 袋の中に戻す 。

# max(x,y) の値の期待値を mod998244353 で出力してください（注記参照）。
# ただし、max(x,y) で x と y のうち小さくない方の値を表します。

fac = [1]
ifac = [1]
for i in range(1, int(4e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def A(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return (fac[n] * ifac[n - k]) % MOD


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    counter = defaultdict(int)
    for num in nums:
        counter[num] += 1

    keys = sorted(counter.keys())
    for k in range(1, n + 1):
        all_ = k * k
        inv = pow(all_, MOD - 2, MOD)
        # 枚举最大值是哪个数
        res = 0
        preSum = 0
        for max_ in keys:
            preSum += counter[max_]
            cur = counter[max_]
            count = k * k - ((preSum - cur) * (preSum - cur)) - (k - preSum) * (k - preSum)
            res += count * max_
            res %= MOD
        print(res * inv % MOD)
