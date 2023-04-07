# https://yukicoder.me/problems/no/665
# 求1^k+2^k+3^k+...+n^k (mod m)的值
# polynomial_geometrical_sum
# 多项式级数k次幂求和
# https://yukicoder.me/submissions/624620


"""
D[i] = f(i)*r**i (f は d 次以下の多項式)のとき、
\sum_0^{n-1} D[i] mod P (n = infty も含む)を計算する
"""

MOD = int(1e9 + 7)
SIZE = 10**5 + 5

fac = [0] * SIZE
ifac = [0] * SIZE
fac[0] = fac[1] = 1
ifac[0] = ifac[1] = 1
for i in range(2, SIZE):
    fac[i] = fac[i - 1] * i % MOD
ifac[-1] = pow(fac[-1], MOD - 2, MOD)
for i in range(SIZE - 1, 0, -1):
    ifac[i - 1] = ifac[i] * i % MOD


def C(n, r):  # nCk mod MOD の計算
    if 0 <= r <= n:
        return (fac[n] * ifac[r] % MOD) * ifac[n - r] % MOD
    return 0


def Lagrange_interpolation(a, t):
    n = len(a) - 1
    t %= MOD
    if 0 <= t <= n:
        return a[t]
    rprod = [1] * (n + 2)
    r = 1
    for i in range(n + 2):
        rprod[n + 1 - i] = r
        r = r * (t - n + i) % MOD
    ans, lprod = 0, 1
    for i, ai in enumerate(a):
        bunsi = lprod * rprod[i + 1] % MOD
        bunbo = ifac[i] * ifac[n - i] % MOD * (-1 if (n - i) % 2 else 1)
        ans += bunsi * bunbo % MOD * ai % MOD
        lprod = lprod * (t - i) % MOD
    return ans % MOD


"""
D[i] = f(i)*r**i (f は d 次以下の多項式)のとき、
\sum_0^{n-1} D[i] mod P を計算する
"""


def polynomial_geometrical_sum(r, d, D, n):
    assert len(D) >= d + 2
    if n <= 0:
        return 0
    r %= MOD
    if r == 0:
        return 1 if d == 0 else 0
    # r==1 なら、累積和の数列は d+1 次式。これを補間して n-1 での値を求める
    if r == 1:
        for i in range(1, d + 2):
            D[i] = (D[i - 1] + D[i]) % MOD
        return Lagrange_interpolation(D, n - 1)
    # n が小さい場合、愚直累積和
    if n <= d + 2:
        return sum(D[:n]) % MOD
    # そうでない場合、累積和の数列は c + g(i)r^i (g(i): d次式)
    # c は極限値。g(i) を補間して g(n-1) を求める
    c = polynomial_geometrical_sum_infty(r, d, D)
    for i in range(1, d + 2):
        D[i] = (D[i - 1] + D[i]) % MOD
    R, rinv = 1, pow(r, MOD - 2, MOD)
    for i in range(d + 2):
        D[i] = (D[i] - c) * R % MOD
        R = R * rinv % MOD
    return (c + Lagrange_interpolation(D, n - 1) * pow(r, n - 1, MOD)) % MOD


"""
D[i] = f(i)*r**i (f は d 次以下の多項式)のとき、
\sum^infty D[i] mod P を計算する
"""


def polynomial_geometrical_sum_infty(r, d, D):
    r %= MOD
    assert r % MOD != 1
    # D と (1-rX)^{d+1} の畳み込みの d 次以下の和 が答え
    acc = res = 0
    R = 1
    for i in range(d + 1):
        acc = (acc + C(d + 1, i) * (1 - i % 2 * 2) * R) % MOD
        res = (res + acc * D[d - i]) % MOD
        R = R * r % MOD
    return res * pow(1 - r, MOD - d - 2, MOD) % MOD


"""
D = [(i**d)*(r**i) for i in range(d+2)] のテーブルを返す
長さ d+2 が返ることに注意
"""


def make_table_monomial_times_geometric(r, d):
    if d == 0:
        return [1, r]
    D = [0] * (d + 2)
    D[1] = 1
    for p in range(2, d + 2):
        if D[p] == 0:  # 素数のとき
            pd = pow(p, d, MOD)
            for j in range(1, (d + 1) // p + 1):
                D[p * j] = (1 if D[j] == 0 else D[j]) * pd % MOD
    if r == 1:
        return D
    R = r
    for i in range(1, d + 2):
        D[i] = D[i] * R % MOD
        R = R * r % MOD
    return D


##########################################################
##########################################################
n, k = map(int, input().split())
D = make_table_monomial_times_geometric(1, k)
print(polynomial_geometrical_sum(1, k, D, n + 1))
