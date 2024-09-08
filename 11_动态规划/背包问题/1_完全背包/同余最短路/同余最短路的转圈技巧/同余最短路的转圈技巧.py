# TODO
# https://www.cnblogs.com/alex-wei/p/17531487.html
# https://zhuanlan.zhihu.com/p/672216458
# 理解：体积很大的多重背包问题
#
# 有趣的是，同余最短路不应该从最短路的角度考虑。
# !其本质上是根据单调性值域定义域互换后将完全背包转化为体积模m意义下的完全背包。
# 普通完全背包的转移是有向无环图，而环上完全背包转移成环，这让我们想到最短路。
# !但因为一个点不会经过它自己，对应原问题就是对于一个物品，不会使得它的总体积为基准物品体积的倍数，
# 所以，我们可以将完全背包转化为类多重背包问题。
#
# for(int i = 0, lim = __gcd(v[i], m); i < lim; i++)
#   for(int j = i, c = 0; c < 2; c += j == i) {
#     int p = (j + v[i]) % m;
#     f[p] = min(f[p], f[j] + v[i]), j = p;
#   }


from math import gcd
from typing import Iterable, List, Tuple

INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


def modShortestPath(coeffs: Iterable[int]) -> Tuple[int, List[int]]:
    """确定线性组合∑ai*xi的可能取到的值(ai非负)

    Args:
        coeffs (List[int]): 非负整数系数,最小的非零ai称为base

    Returns:
        Tuple[int, List[int]]: base, dist
        base (int): 最小的非零ai
        dist (List[int]): dist[i]记录的是最小的x,满足x=i(mod base)且x能被系数coeffs线性表出(xi非负)
        如果不存在这样的x,则dist[i]为INF
        如果coeff全为0,则返回空数组
    """
    coeffs = [v for v in coeffs if v > 0]
    if not coeffs:
        return 0, []

    base = min(coeffs)
    dp = [INF] * base  # dp[i]表示模base余数为i时，最小的k
    dp[0] = 0
    for v in coeffs:
        cycle = gcd(base, v)  # 在这些环上转移
        for j in range(cycle):
            cur = j
            count = 0
            while count < 2:  # 转两圈，涵盖从每个点出发的情况
                next = (cur + v) % base
                dp[next] = min2(dp[next], dp[cur] + v)
                cur = next
                count += cur == j
    return base, dp


# P2371 [国家集训队] 墨墨的等式
# https://www.luogu.com.cn/problem/P2371
# 给定n个系数coeffs和上下界lower,upper
# !对于 lower<=k<=upper 求有多少个k能够满足
# !a0*x0+a1*x1+...+an*xn=k
# n<=12 0<=ai<=5e5 1<=lower<=upper<=2^63-1
# !时间复杂度：O(n*ai)
def p2371() -> None:
    _, lower, upper = map(int, input().split())
    coeffs = list(map(int, input().split()))
    coeffs = [v for v in coeffs if v != 0]
    if not coeffs:
        print(0)
        return

    base, dp = modShortestPath(coeffs)

    res = 0
    for i in range(base):
        if upper >= dp[i]:
            res += (upper - dp[i]) // base + 1
        if lower > dp[i]:
            res -= (lower - dp[i] - 1) // base + 1
    print(res)


# P2662 牛场围栏(求最大的不能被线性表出的数)
# https://www.luogu.com.cn/problem/P2662
# 给一堆系数，求最大的不能被线性表出的数。
# 如果任何数可以被表出或者这个最大值不存在，输出-1
# n<100,ai<3000
# !从起点开始有一些点无法到达(即有一整个剩余系都不能被表出)
def p2662() -> None:
    n, cut = map(int, input().split())  # 木料的种类和每根木料削去的最大值
    sticks = list(map(int, input().split()))  # 第i根木料的原始长度
    coeffs = set()
    for s in sticks:
        for c in range(min(cut, s) + 1):
            coeffs.add(s - c)

    base, dist = modShortestPath(coeffs)
    if any(v == INF for v in dist):  # 这个剩余类不能被表出
        print(-1)
        exit(0)

    print(max(dist) - base)


# P9140 [THUPC 2023 初赛] 背包 (超大完全背包)
# https://www.luogu.com.cn/problem/P9140
#
# 有n个物品,第i种物品单个体积为vi,价值为ci
# q次询问,每次给出背包的容积,你需要选择若干个物品,
# 每种物品可以选任意多个
# !在选出物品的体积的和恰好为 V 的前提下最大化选出物品的价值的和。
# 若不存在体积和恰好为 V 的方案，输出 -1
# 为了体现你解决 NP-Hard 问题的能力，V 会远大于 vi，详见数据范围部分。
# n<=50 vi<=1e5 ci<=1e5 q<=1e5
# !1e11<=V<=1e12
# O(n*max(vi)
def p9140() -> None:
    n, q = map(int, input().split())
    goods = list(tuple(map(int, input().split())) for _ in range(n))  # !vi,ci

    baseV, baseC = 1, 0  # 性价比最高的物品
    gcd_ = 0
    for v, c in goods:
        if baseV * c > v * baseC:
            baseV, baseC = v, c
        gcd_ = gcd(gcd_, v)
    dp = [-INF] * baseV
    dp[0] = 0
    for v, c in goods:
        cycle = gcd(baseV, v)
        for j in range(cycle):
            cur = j
            count = 0
            while count < 2:
                next = (cur + v) % baseV
                # !选这个物品，就要抛弃  baseC * ((cur + v) // baseV 的价值
                dp[next] = max2(dp[next], dp[cur] + c - baseC * ((cur + v) // baseV))
                cur = next
                count += cur == j

    for _ in range(q):
        cap = int(input())
        if cap % gcd_:
            print(-1)
        else:
            print(dp[cap % baseV] + cap // baseV * baseC)


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    # p2371()
    # p2662()
    # p9140()
