# AtCoder Beginner Contest 246 A~G - pzr的文章 - 知乎
# https://zhuanlan.zhihu.com/p/492352451

# 给定一棵根为1的树，除了根，每个点都有一个权值w;
# 现在小明和小红从根节点开始按照如下规则玩一个游戏:
# 1.小红任意选择一个点，把这个点的权值变为0.
# 2.小明从当前点出发，可以走到任意一个儿子节点
# 3.然后小明可以决定是否结束游戏(如果小明在叶子节点则必须结束游戏)
# 最后小明获得的分数就是小明所在点的权值，小明希望获得的分数尽可能得高，小红希望小明获得的分数尽可能的低
# 假设两人都足够聪明的情况下(即总是做出对当前最有利的操作)，小明可以获得的最大分数是多少。

# !二分答案(离散化数组) + 树形 dp
from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def main() -> None:
    def check(mid: int) -> bool:
        def dfs(cur: int, pre: int) -> int:
            """(对方最优操作下)子树中权值>=mid的点的个数"""
            subtree = 0
            for next in adjMap[cur]:
                if next == pre:
                    continue
                subtree += dfs(next, cur)
            subtree = max(0, subtree - 1)  # 被对方移除了一个
            subtree += int(values[cur] >= mid)  # 根节点是否可以
            return subtree

        return dfs(0, -1) >= 1

    n = int(input())
    values = [-1] + list(map(int, input().split()))
    adjMap = defaultdict(set)
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u].add(v)
        adjMap[v].add(u)

    # !注意这里要离散化权值 不然python过不去
    # left, right = 1, int(1e9 + 7)
    # while left <= right:
    #     mid = (left + right) // 2
    #     if check(mid):
    #         left = mid + 1
    #     else:
    #         right = mid - 1

    allValues = sorted(set(values))
    left, right = 0, len(allValues) - 1  # !二分答案 allValues里的第几个值是答案
    while left <= right:
        mid = (left + right) // 2
        if check(allValues[mid]):
            left = mid + 1
        else:
            right = mid - 1
    print(allValues[right])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
