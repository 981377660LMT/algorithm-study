# abc349-F - Subsequence LCM
# https://atcoder.jp/contests/abc349/tasks/abc349_f
# lcm为M的子序列个数
# https://atcoder.jp/contests/abc349/submissions/52267196

# zeta变换

mod = 998244353


def fact(n):
    res = n
    a = []
    i = 2
    while i * i <= res:
        if res % i == 0:
            cnt = 0
            while res % i == 0:
                cnt += 1
                res //= i
            a.append((i, cnt))
        i += 1
    if res != 1:
        a.append((res, 1))
    return a


N, M = map(int, input().split())
A = list(map(int, input().split()))
f = fact(M)
K = len(f)
cnt = [0] * (1 << K)
for v in A:
    if M % v != 0:
        continue
    bit = 0
    for j in range(K):
        p, e = f[j]
        if v % (p**e) == 0:
            bit |= 1 << j
    cnt[bit] += 1

# ゼータ変換
for v in range(K):
    for j in range(1 << K):
        if (j >> v) & 1 == 0:
            cnt[j | (1 << v)] += cnt[j]

h = [pow(2, cnt[i], mod) for i in range(1 << K)]
# メビウス変換
for v in range(K):
    for j in range(1 << K):
        if (j >> v) & 1 == 0:
            h[j | (1 << v)] -= h[j]
            h[j | (1 << v)] %= mod

res = h[(1 << K) - 1]
if M == 1:
    res -= 1
print(res % mod)
