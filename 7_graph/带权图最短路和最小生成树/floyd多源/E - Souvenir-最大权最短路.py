from typing import List, Tuple


# 给定一个n个顶点的有向图,每个点有一个权值
# 给定q个询问,每个询问给定起点s和终点t
# !问从s到t的最短路上的点权和的最大值
# 如果无法到达,输出Impossible,否则输出(距离,最大权值和)
# !n<=300 q<=n*(n-1)

INF = int(4e18)


def souvenir(
    n: int, links: List[str], prices: List[int], queries: List[Tuple[int, int]]
) -> List[Tuple[int, int]]:
    dist = [[INF] * n for _ in range(n)]
    priceSum = [[0] * n for _ in range(n)]  # i -> j の最大価格
    for i in range(n):
        dist[i][i] = 0
        priceSum[i][i] = prices[i]

    for i in range(n):
        for j in range(n):
            if links[i][j] == "Y":
                dist[i][j] = 1
                priceSum[i][j] = prices[i] + prices[j]

    for k in range(n):
        for i in range(n):
            for j in range(n):
                cand = dist[i][k] + dist[k][j]
                if cand < dist[i][j]:
                    dist[i][j] = cand
                    priceSum[i][j] = priceSum[i][k] + priceSum[k][j] - prices[k]
                elif cand == dist[i][j]:
                    tmp = priceSum[i][k] + priceSum[k][j] - prices[k]
                    if tmp > priceSum[i][j]:
                        priceSum[i][j] = tmp

    res = []
    for (start, target) in queries:
        if dist[start][target] == INF:
            res.append((INF, 0))
        else:
            res.append((dist[start][target], priceSum[start][target]))
    return res


if __name__ == "__main__":
    n = int(input())
    prices = list(map(int, input().split()))
    links = [input() for _ in range(n)]
    q = int(input())
    queries = []
    for _ in range(q):
        s, t = map(int, input().split())
        s, t = s - 1, t - 1
        queries.append((s, t))

    res = souvenir(n, links, prices, queries)
    for (dist, price) in res:
        if dist == INF:
            print("Impossible")
        else:
            print(dist, price)
