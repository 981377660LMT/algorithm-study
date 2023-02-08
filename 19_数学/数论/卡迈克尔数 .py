# 卡迈克尔数
# カーマイケル関数 (Carmichael function)
# https://tjkendev.github.io/procon-library/python/prime/carmichael-function.html

# 卡迈克尔数λ(n)的定义:
# 对[1,n]内所有与n互素的数a都有 a^m ≡ 1 (mod n) ,
# 把这个最小的m称为卡迈克尔数λ(n).


def gcd(m, n):
    while n:
        m, n = n, m % n
    return m


def lcm(m, n):
    return m // gcd(m, n) * n


# calculate λ(x)
def carmichael(x):
    r = 1

    # p_i = 2
    b = 0
    while x & 1 == 0:
        b += 1
        x >>= 1
    if b > 1:
        r = 2 if b == 2 else 2 ** (b - 2)

    # p_i ≥ 3
    y = 3
    while y * y <= x:
        if x % y == 0:
            c = 0
            while x % y == 0:
                x //= y
                c += 1
            r = lcm(r, (y - 1) * y ** (c - 1))
        y += 1
    if x > 1:
        r = lcm(r, x - 1)
    return r


print(carmichael(10))
