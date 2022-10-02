# 给定一个数字串，向其中随意放入加号组成表达式，问所有情况的表达式值的和是多少
# len(s)<=2e5
# 例如:
# 1234
# !1234, 123+4, 12+34, 12+3+4, 1+234, 1+23+4, 1+2+34, 1+2+3+4
# 和为 1736

# !计算贡献
# 1*1/2 + 10*1/4 + 100*1/8 + ...
# !∑ci*(2**i-1)*10**(n-i) (i=1..n)


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
pow2 = [1]
pow10 = [1]
for _ in range(int(2e5 + 10)):
    pow2.append(pow2[-1] * 2 % MOD)
    pow10.append(pow10[-1] * 10 % MOD)


if __name__ == "__main__":
    s = input()
    n = len(s)
    if n == 1:
        print(s)
        exit(0)

    res, curSum = 0, pow2[n - 1]
    for i in range(n - 1, -1, -1):
        num = int(s[i])

    res *= pow(2, n - 2, MOD)
    res %= MOD
    print(res)
