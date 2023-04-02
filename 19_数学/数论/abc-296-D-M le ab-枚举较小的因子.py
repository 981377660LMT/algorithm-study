# 分解为两个数的乘积的最小整数
# https://atcoder.jp/contests/abc296/tasks/abc296_d
# 暴力枚举因数

# 给定n<=1e12,m<=1e12
# !要求在 1 ~ n 中找到两个数 a,b 使得在满足 a * b ≥ m 的条件下，
# 同时 a*b 最小。
# 不存在则输出-1


# 根号=> 枚举较小的因子
# !设1<=a<=b<=n 从1到n枚举a,则b为ceil(m/a),每次更新a*b的最小值
# !因为要保证a<=b 所以 a>b就break 那么a最多枚举 sqrt(m)+1次

from math import ceil


INF = int(1e18)


def solve(n: int, m: int) -> int:
    res = INF
    for small in range(1, n + 1):
        big = ceil(m / small)
        if 1 <= big <= n:
            res = min(res, small * big)
        if small > big:
            break
    return res if res != INF else -1


if __name__ == "__main__":
    n, m = map(int, input().split())
    print(solve(n, m))
