# 最后函数要平移到一起 考察两个数组的diff(力扣也有类似题)

from collections import Counter
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

n = int(input())
nums1 = list(map(int, input().split()))
nums2 = list(map(int, input().split()))

# 增量相同或者减量相同才不会互相嘲笑
# 修改TOM分数 至少多少关
diff = [a - b for a, b in zip(nums1, nums2)]
counter = Counter(diff)
max_ = max(counter.values())
print(n - max_)

# import sys
# from typing import Set, Tuple


# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")

# n = int(input())
# string = input()
# adjList = [[] for _ in range(n)]
# for _ in range(n - 1):
#     u, v = map(int, input().split())
#     u, v = u - 1, v - 1
#     adjList[u].append(v)
#     adjList[v].append(u)


# def dfs(cur: int, pre: int) -> Tuple[Set[str], int]:
#     if len(adjList[cur]) == 1 and adjList[cur][0] == pre:
#         return set([string[cur]]), 0
#     subTree, res = set(), 0
#     for next in adjList[cur]:
#         if next == pre:
#             continue
#         nextTree, nextRes = dfs(next, cur)
#         res += nextRes
#         if len(nextTree) > len(subTree):
#             subTree, nextTree = nextTree, subTree
#         subTree |= nextTree

#     tmp = len(subTree)
#     subTree.add(string[cur])
#     return subTree, res + tmp


# print(dfs(0, -1)[1])
