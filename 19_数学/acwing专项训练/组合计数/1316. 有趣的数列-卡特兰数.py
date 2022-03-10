MOD = 5000011

n, k = map(int, input().split())

# 我们称一个长度为 2n 的数列是有趣的，当且仅当该数列满足以下三个条件：

# 它是从 1 到 2n 共 2n 个整数的一个排列 {ai}；
# 所有的奇数项满足 a1<a3<⋯<a2n−1 ，所有的偶数项满足 a2<a4<⋯<a2n；
# 任意相邻的两项 a2i−1 与 a2i (1≤i≤n) 满足奇数项小于偶数项，即：a2i−1<a2i。
# 任务是：对于给定的 n，请求出有多少个不同的长度为 2n 的有趣的数列。

# from math import comb, factorial

# n = int(input())
# a = factorial(2 * n)
# b = factorial(n)
# print(a // b // b // (n + 1))
# print(comb(2 * n, n) // (n + 1))


def getPs(n):
    for i in range(2, n + 1):
        if not st[i]:
            ps.append(i)
        j = 0
        while ps[j] <= n // i:
            st[ps[j] * i] = True
            if not i % ps[j]:
                break
            j += 1


def get(a, p):
    cnt = 0
    while a:
        a //= p
        cnt += a

    return cnt


def qmi(a, k, p):
    res = 1
    while k:
        if k & 1:
            res = res * a % p
        k >>= 1
        a = a * a % p
    return res


def C(a, b):

    res = 1
    for p in ps:
        if p > a:
            break
        v = get(a, p) - get(a - b, p) - get(b, p)
        res = res * qmi(p, v, mod) % mod

    return res


N = 2000010
n, mod = map(int, input().split())
ps = []
st = [False] * N
getPs(N - 1)

print((C(2 * n, n) - C(2 * n, n - 1)) % mod)

