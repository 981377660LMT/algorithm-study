# 莫比乌斯函数
# https://tjkendev.github.io/procon-library/python/prime/moebius-function.html


# 莫比乌斯函数 μ(n) 的定义:
# !1. n如果存在平方素因子(被p^2整除),μ(n)=0;
# !2. 如果n有k个素因子,μ(n)=(-1)^k.

# calculate μ(n): O(√N) μ(n)
def moebius(n):
    x = 2
    c = 0
    while x * x <= n:
        if n % x == 0:
            n //= x
            if n % x == 0:
                return 0
            c += 1
        x += 1
    if n > 1:
        c += 1
    return -1 if c % 2 else 1


# a table of μ(n) for n in [0, N]
# O(N log log N)
def moebius_table(n):
    (*p,) = range(n + 1)
    sq = int(n**0.5)
    for x in range(2, sq + 1):
        if p[x] == x:
            for y in range(x * x, n + 1, x):
                p[y] = x
            for y in range(x * x, n + 1, x * x):
                p[y] = 0
    res = [0] * (n + 1)
    res[0] = 0
    res[1] = 1
    for x in range(2, n + 1):
        res[x] = p[x] and -res[x // p[x]]
    return res


print(moebius(10))
# => "1"
print(moebius_table(20))
# => "[0, 1, -1, -1, 0, -1, 1, -1, 0, 0, 1, -1, 0, -1, 1, 1, 0, -1, 0, -1, 0]"
