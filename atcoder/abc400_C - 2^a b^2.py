# abc400_C - 2^a b^2
# https://atcoder.jp/contests/abc400/tasks/abc400_c
# !求1-n中形如2^a*b^2的数的个数.
# 1<=n<=1e18
#
# !枚举a, b为奇数个数.


from math import isqrt

if __name__ == "__main__":
    N = int(input())
    res = 0
    for a in range(1, N.bit_length()):
        x = N // (1 << a)
        sqrt_ = isqrt(x)
        res += (1 + sqrt_) // 2  # [1, sqrt_]中奇数个数
    print(res)
