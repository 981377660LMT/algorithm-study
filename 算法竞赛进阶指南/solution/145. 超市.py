# 超市里有 N 件商品，每件商品都有利润 pi 和过期时间 di，每天只能卖一件商品，
# 过期商品不能再卖。
# 求合理安排每天卖的商品的情况下，可以得到的最大收益是多少。
# 0≤N≤10000,
# 1≤pi,di≤10000

from heapq import heappop, heappush


def main() -> None:

    n, *rest = input().split()
    n = int(n)
    scores = list(map(int, rest[::2]))
    days = list(map(int, rest[1::2]))

    pq = []
    for score, day in zip(scores, days):
        heappush(pq, (day, -score))

    res = 0
    for d in range(10000 + 10):
        while pq and pq[0][0] <= d:
            heappop(pq)
        if pq:
            res += -heappop(pq)[1]
    print(res)


while True:
    try:
        main()
        input()
    except EOFError:
        break

