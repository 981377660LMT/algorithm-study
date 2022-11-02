# B.Broken Rounding
# !连续四舍五入

# 题目简述:将x每次改为一个与10**i的倍数相差最小的数字，重复k次，
# i的取值为1,2,3,...,k。
# x<=1e15
# k<=15

# 我们只需要快速找到两个符合的倍数就可以了。
# !这时候可以想到两个数: `floor(x/10**i)*10**i` 和 `ceil(x/10**i)*10**i`。
# 我们只需要找到这两个数与x的差谁最小，
# 并用那个数取代即可，其他
# 倍数显然比这两个数更不符合。


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    x, k = map(int, input().split())
    for i in range(1, k + 1):
        base = 10**i
        div = x // base
        lower, upper = div * base, (div + 1) * base
        if abs(x - lower) < abs(x - upper):
            x = lower
        else:
            x = upper
    print(x)
