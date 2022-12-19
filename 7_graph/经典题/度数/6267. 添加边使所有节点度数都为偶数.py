# 6267. 添加边使所有节点度数为偶数

from collections import Counter
from itertools import chain, permutations
from typing import List

# 给你一个有 n 个节点的 无向 图，节点编号为 1 到 n 。再给你整数 n 和一个二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示节点 ai 和 bi 之间有一条边。图不一定连通。
# 你可以给图中添加 至多 两条额外的边（也可以一条边都不添加），使得图中没有重边也没有自环。
# 如果添加额外的边后，可以使得图中所有点的度数都是偶数，返回 true ，否则返回 false 。
# 点的度数是连接一个点的边的数目。


# 1. 奇数度数的顶点为偶数个 只需要讨论0 2 4个奇数度数的顶点
# 2. 0 => ok
#    2 => 两个奇数度数的顶点之间相连或者两个奇数度数的顶点和一个偶数度数的顶点之间相连  4 => 4个奇数度数的顶点之间相连
#    4 => 四个奇数度数的顶点之间相连的情况


class Solution:
    def isPossible(self, n: int, edges: List[List[int]]) -> bool:
        adjList = [set() for _ in range(n)]
        for u, v in edges:
            u, v = u - 1, v - 1
            adjList[u].add(v)
            adjList[v].add(u)

        odd, even = [], []
        for i in range(n):
            if len(adjList[i]) % 2 == 1:
                odd.append(i)
            else:
                even.append(i)

        if len(odd) == 0:
            return True

        if len(odd) == 2:  # 还可以一个偶顶点连两
            u, v = odd
            if v not in adjList[u]:
                return True
            # 枚举所有的偶顶点
            for w in even:
                if w not in adjList[u] and w not in adjList[v]:
                    return True
            return False

        if len(odd) == 4:
            for a, b, c, d in permutations(odd):
                if a not in adjList[b] and c not in adjList[d]:
                    return True
            return False

        return False


assert Solution().isPossible(n=5, edges=[[1, 2], [2, 3], [3, 4], [4, 2], [1, 4], [2, 5]])


# 错误解法:
# class Solution2:
#     def isPossible(self, n: int, edges: List[List[int]]) -> bool:
#         c = [x for x in Counter(chain(*edges)).values() if x & 1]
#         return len(c) <= 4 and not (n % 2 == 0 and n - 1 in c) and c.count(n - 2) < 3
# generate random cases, until find a case that solution1 is not equal to solution2
# import random

# count = 0
# while True:
#     n = random.randint(1, 5)
#     edges = []
#     for _ in range(random.randint(0, n * (n - 1) // 2)):
#         u, v = random.randint(1, n), random.randint(1, n)
#         if u != v and [u, v] not in edges and [v, u] not in edges:
#             edges.append([u, v])
#     s1 = Solution().isPossible(n, edges)
#     s2 = Solution2().isPossible(n, edges)

#     count += 1
#     if s1 != s2:
#         print(n, edges)
#         print(s1, s2)
#         print(count)
#         break
# # a = [[4, 1], [1, 3], [4, 3], [1, 5], [5, 3], [2, 5], [2, 1]]
# # a = [sorted(x) for x in a]
# # a.sort()
# # print(a)
