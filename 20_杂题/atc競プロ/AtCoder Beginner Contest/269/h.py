# 给定一颗以0为根节点的树
# 输出大小为k的集合个数
# 集合中任意两个树结点没有祖孙关系


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    adjList = [[] for _ in range(n)]
    parents = list(map(int, input().split()))
    for cur, pre in enumerate(parents, 1):
        pre -= 1
        adjList[pre].append(cur)
    print(adjList)
