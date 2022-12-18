# E - Choose Two and Eat One-广义的最小生成树
# 箱子里有n个球 进行n-1次操作:
# !每次取出两个球,吃掉其中一个,剩下的球放回箱子 获得的分数为 (x^y+y^x)%MOD
# 求最大的分数
# n<=500 1<=nums[i]<=M-1<=1e9
# !O(n^2logk) 最大生成树

# !总结与反思
# 0.图的生成树:从若干条边(广义的边)中选n-1个对,使得这n个点连通,且这些对的得分最大
#   也可以理解为从叶子节点开始不断删除点和选取边,直到只剩下一个点
# 1.发现有后效性时,就不要想dp了,换个角度思考问题(dp不对,换成图论)
# !2.对二维矩阵不敏感 只想到了网格图 没想到邻接矩阵
# !3.奇怪的表达式一般是边权(代价)


from itertools import combinations
from typing import List
from 模板 import kruskal


def chooseTwoAndEatOne(nums: List[int], mod: int) -> int:
    n = len(nums)
    edges = []
    for i, j in combinations(range(n), 2):
        x, y = nums[i], nums[j]
        score = (pow(x, y, mod) + pow(y, x, mod)) % mod
        edges.append((i, j, -score))
    res, _ = kruskal(n, edges)
    return -res


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    nums = list(map(int, input().split()))
    print(chooseTwoAndEatOne(nums, m))
