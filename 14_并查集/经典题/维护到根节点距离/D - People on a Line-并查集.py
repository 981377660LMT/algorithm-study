# 每个判断为 right-left=dist
# 问所有的判断是否无矛盾

# !dfs或bfs求出每个点到每个组的根的距离,再逐一检验
from UnionFindMapWithDist import UnionFindMapWithDist2
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    uf = UnionFindMapWithDist2()
    for _ in range(m):
        left, right, weight = map(int, input().split())
        left, right = left - 1, right - 1
        uf.add(left).add(right)
        if not uf.isConnected(left, right):
            uf.union(left, right, weight)
        else:
            dist1, dist2 = uf.distToRoot[left], uf.distToRoot[right]
            if dist1 - weight != dist2:
                print("No")
                exit(0)

    print("Yes")
