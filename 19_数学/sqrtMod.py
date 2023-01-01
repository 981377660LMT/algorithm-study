# sqrtMod 二次剩余
# 给定Y,P(P是素数)
# 求X^2=Y(mod P) 的一个解X


def cipolla(b: int, p: int) -> int:
    def mul_quad(a, b, theta, p):
        x = a[0] * b[0] + theta * a[1] * b[1]
        x %= p
        y = a[0] * b[1] + a[1] * b[0]
        y %= p
        return x, y

    def pow_quad(a, n, theta, p):
        res = (1, 0)
        while n:
            if n & 1:
                res = mul_quad(res, a, theta, p)
            a = mul_quad(a, a, theta, p)
            n >>= 1
        return res

    b %= p
    if p == 2:
        return b
    if b == 0:
        return 0
    if pow(b, (p - 1) // 2, p) != 1:
        return -1
    c = 0
    while pow((c * c - b) % p, (p - 1) // 2, p) == 1:
        c += 1
    theta = (c * c - b) % p
    res = pow_quad((c, 1), (p + 1) // 2, theta, p)
    return res[0]


t = int(input())
for _ in range(t):
    y, p = map(int, input().split())
    print(cipolla(y, p))
