"""
点分治解法
https://leetcode.cn/problems/number-of-good-paths/solution/dian-fen-zhi-jie-fa-by-vclip-nr5b/
"""

# # 给定一个有 N 个点（编号 0,1,…,N−1）的树，每条边都有一个权值（不超过 1000）。
# # 树上两个节点 x 与 y 之间的路径长度就是路径上各条边的权值之和。
# # 求长度不超过 K 的路径(树上距离小于等于k的对数)有多少条。
# # n<=1e4


# from collections import defaultdict
# from typing import DefaultDict, List, Set, Tuple


# def main(n: int, adjMap: DefaultDict[int, Set[Tuple[int, int]]], k: int) -> int:
#     """求长度不超过 K 的路径有多少条

#     暴力:n^2 从每个点出发bfs
#     点分治:nlogn^2
#     我们找到这棵树的重心G，把这棵树分为若干个子树，那么发现满足条件的点对只有3种情况：
#     1.点对在某个子树中（直接递归求解）
#     2.两个点所构成的路径经过了重心G，但你会发现这两个点一定不能在同一个子树中。
#     所以我们处理出当前这棵树中每个点的d值，didi表示点ii到重心G的距离。
#     那么只需要用di+dj≤kdi+dj≤k这样(i,j)(i,j)的数量减去di+dj≤kdi+dj≤k且满足i,ji,j在同一个子树中的数量
#     而你会发现，后者可以在递归子树中处理。
#     3.这条路经的一个端点是G，那么实质上和2.是一种情况，再加入一个dG=0dG=0即可
#     """

#     def getCenter() -> List[int]:
#         """求重心"""

#         def dfs(cur: int, pre: int) -> None:
#             if cur in removed:
#                 return
#             subsize[cur] = 1
#             for next, _ in adjMap[cur]:
#                 if next == pre:
#                     continue
#                 dfs(next, cur)
#                 subsize[cur] += subsize[next]
#                 maxsize[cur] = max(maxsize[cur], subsize[next])
#             maxsize[cur] = max(maxsize[cur], n - subsize[cur])
#             if maxsize[cur] <= n / 2:
#                 res.append(cur)

#         res = []
#         maxsize = [0] * n  # 最大连通块大小
#         subsize = [0] * n  # 子树大小
#         dfs(0, -1)
#         return res

#     def cal(root: int) -> int:
#         """处理root所在的树"""
#         if root in removed:
#             return 0
#         res = 0
#         center = getCenter()[0]

#     removed = set()
#     return cal(0)


# while True:
#     n, k = map(int, input().split())
#     if n == 0 and k == 0:
#         break
#     adjMap = defaultdict(set)
#     for _ in range(n - 1):
#         u, v, w = map(int, input().split())
#         adjMap[u].add((v, w))
#         adjMap[v].add((u, w))
#     print(main(n, adjMap, k))

# https://www.acwing.com/activity/content/problem/content/2974/


# todo
