# 211. 计算系数
# 给定一个多项式 (ax+by)^k，请求出多项式展开后 x^n*y^m 项的系数。
MOD = 10007

a, b, k, n, m = map(int, input().split())

# x^n*y^m 的系数是  Cnk*a^n*b^m
# 分母中每个数都需要做一次快速幂，因此总时间复杂度是 O(nlogn)。

res = pow(a, n, MOD) * pow(b, m, MOD)

for i in range(k, k - n, -1):
    res = res * i % MOD
    res = res * pow(k - i + 1, MOD - 2, MOD) % MOD
print(res)
