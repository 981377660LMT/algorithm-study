# abc357-D-外观数列2
# f(n) 表示n个数字n相连，求f(n)模998244353的值
#
# !思路：先写出通项公式.
# 记n的十进制位数为m, 则 f(n) = n + n*10^m + n*10^(2m) + ... + n*10^((n-1)m)
# !等比数列化简得 f(n) = n * (10^(nm) - 1) / (10^m - 1)

MOD = 998244353


def solve(n: int) -> int:
    m = len(str(n))
    num = n * (pow(10, n * m, MOD) - 1) % MOD
    den = (pow(10, m, MOD) - 1) % MOD
    den = pow(den, MOD - 2, MOD)
    return num * den % MOD


if __name__ == "__main__":
    n = int(input())
    print(solve(n))
