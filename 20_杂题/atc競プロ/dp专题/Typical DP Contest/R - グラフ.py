# 给出n个点的有向图 无重边无自环
# 选择两条`路径`染色 求染色结点个数最大值
# n<=300

# https://simezi-tan.hatenadiary.org/entry/20140909/1410211433
#  強連結成分分解 ＋ 推移閉包 ＋ 動的計画法
# 缩点成DAG后 dp[i][j] 表示两条路径从根节点出发 端点为i和j时答案
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n = int(input())
adjMatrix = []
for _ in range(n):
    row = list(map(int, input().split()))
    adjMatrix.append(row)
