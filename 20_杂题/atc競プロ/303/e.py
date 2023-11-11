from collections import Counter, deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 以下の条件を満たす
# k+1 頂点
# k 辺のグラフをレベル
# k (k≥2) の星と呼びます。

# ある
# 1 つの頂点が、他の
# k 個の頂点と
# 1 本ずつ辺で結ばれている。それ以外の辺は存在しない。
# 高橋君は、はじめ何個かの星からなるグラフを持っていました。そして、以下の手続きを全てのグラフの頂点が連結になるまでくり返し行いました。

# 持っているグラフの頂点から二つの頂点を選ぶ。このとき、選んだ二つの頂点の次数は共に
# 1 であり、かつ選んだ二つの頂点は非連結である必要がある。選んだ二つの頂点を結ぶ辺を張る。
# その後、高橋君は手続きが終了した後のグラフの頂点に、適当に
# 1 から
# N の番号を付けました。このグラフは木となっており、これを
# T と呼びます。
# T には
# N−1 本の辺があり、
# i 番目の辺は
# u
# i
# ​
#   と
# v
# i
# ​
#   を結んでいました。

# その後高橋君は、はじめ持っていた星の個数とレベルを忘れてしまいました。
# T の情報からはじめ持っていた星の個数とレベルを求めてください。

# 星图
if __name__ == "__main__":
    n = int(input())
    adjList = [[] for _ in range(n)]
    deg = [0] * n
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        adjList[u].append(v)
        adjList[v].append(u)
        deg[u] += 1
        deg[v] += 1
    leaves = [i for i in range(n) if deg[i] == 1]
    for i in leaves:
        for next in adjList[i]:
            deg[next] += 1
    star = [d - 1 for d in deg if d >= 3]
    print(*sorted(star))
